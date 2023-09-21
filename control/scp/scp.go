package scp

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"
	"tryssh/config"
	"tryssh/launcher/scp"
	"tryssh/utils"
)

const sshClientTimeoutWhenLogin = 5 * time.Second

type Controller struct {
	source        string
	destination   string
	configuration *config.MainConfig
	cacheIsFound  bool
	cacheIndex    int
	destIp        string
	concurrency   int
	sshTimeout    time.Duration
}

// TryCopy Functional entrance
func (cc *Controller) TryCopy(user string, concurrency int, sshTimeout time.Duration) {
	// Set timeout
	cc.sshTimeout = sshTimeout
	// Set concurrency
	cc.concurrency = concurrency
	if strings.Contains(cc.source, ":") {
		cc.destIp = strings.Split(cc.source, ":")[0]
		remotePath := strings.Split(cc.source, ":")[1]
		// Obtain the real address based on the alias
		cc.searchAliasExistsOrNot()
		// Reassemble remote server address and file path
		cc.source = strings.Join([]string{cc.destIp, remotePath}, ":")
	} else if strings.Contains(cc.destination, ":") {
		cc.destIp = strings.Split(cc.destination, ":")[0]
		remotePath := strings.Split(cc.destination, ":")[1]
		// Obtain the real address based on the alias
		cc.searchAliasExistsOrNot()
		// Reassemble remote server address and file path
		cc.destination = strings.Join([]string{cc.destIp, remotePath}, ":")
	} else {
		return
	}
	// Obtain the real address based on the alias
	cc.searchAliasExistsOrNot()
	// Reassemble remote server address and file path
	var targetServer *config.ServerListConfig
	targetServer, cc.cacheIndex, cc.cacheIsFound = config.SelectServerCache(user, cc.destIp, cc.configuration)

	if cc.cacheIsFound {
		utils.Logger.Infof("The cache for %s is found, which will be used to try.\n", cc.destIp)
		cc.tryCopyWithCache(user, targetServer)
	} else {
		utils.Logger.Warnf("The cache for %s could not be found. Start trying to login.\n\n", cc.destIp)
		cc.tryCopyWithoutCache(user)
	}
	utils.Logger.Fatalln("There is no password combination that can log in successfully\n")
}

func (cc *Controller) tryCopyWithCache(user string, targetServer *config.ServerListConfig) {
	lan := &scp.Launcher{
		SshConnector: *config.GetSshConnectorFromConfig(targetServer),
		Src:          cc.source,
		Dest:         cc.destination,
	}
	// Set default timeout time
	lan.SshTimeout = sshClientTimeoutWhenLogin
	if lan.Launch() {
		os.Exit(0)
	} else {
		utils.Logger.Errorln("Failed to log in with cached information. Start trying to login again.\n\n")
		cc.tryCopyWithoutCache(user)
	}
}

func (cc *Controller) tryCopyWithoutCache(user string) {
	combinations := config.GenerateCombination(cc.destIp, user, cc.configuration)
	launchers := scp.NewScpLaunchersByCombinations(combinations, cc.source, cc.destination, cc.sshTimeout)
	hitLaunchers := concurrencyTryToConnect(cc.concurrency, launchers)
	if hitLaunchers != nil {
		utils.Logger.Infoln("Login succeeded. The cache will be added.\n")
		// Determine if the login attempt was successful after the old cache login failed.
		// If so, delete the old cache information that cannot be logged in after the login attempt is successful
		if cc.cacheIsFound {
			utils.Logger.Infoln("The old cache will be deleted.\n")
			cc.configuration.ServerLists = append(
				cc.configuration.ServerLists[:cc.cacheIndex], cc.configuration.ServerLists[cc.cacheIndex+1:]...)
		}
		newServerCache := config.GetConfigFromSshConnector(&hitLaunchers[0].SshConnector)
		cc.configuration.ServerLists = append(cc.configuration.ServerLists, *newServerCache)
		if config.UpdateConfig(cc.configuration) {
			utils.Logger.Infoln("Cache added.\n\n")
			// If the timeout time is less than sshClientTimeoutWhenLogin during login,
			// change to sshClientTimeoutWhenLogin
			if hitLaunchers[0].SshTimeout < sshClientTimeoutWhenLogin {
				hitLaunchers[0].SshTimeout = sshClientTimeoutWhenLogin
			}
			hitLaunchers[0].Launch()
		} else {
			utils.Logger.Errorln("Cache added failed.\n\n")
		}
		os.Exit(0)
	}
}

func (cc *Controller) searchAliasExistsOrNot() {
	for _, server := range cc.configuration.ServerLists {
		if server.Alias == cc.destIp {
			cc.destIp = server.Ip
		}
	}
}

func concurrencyTryToConnect(concurrency int, launchers []*scp.Launcher) []*scp.Launcher {
	var hitLaunchers []*scp.Launcher
	var mutex sync.Mutex
	// If the number of launchers is less than the set concurrency, change the concurrency to the number of launchers
	if concurrency > len(launchers) {
		concurrency = len(launchers)
	}
	launchersChan := make(chan *scp.Launcher)
	ctx, cancelFunc := context.WithCancel(context.Background())
	// Producer
	go func(ctx context.Context, launchersChan chan<- *scp.Launcher, launchers []*scp.Launcher) {
		for _, launcherP := range launchers {
			select {
			case <-ctx.Done():
				break
			default:
				launchersChan <- launcherP
			}
		}
		close(launchersChan)
	}(ctx, launchersChan, launchers)
	// Consumer
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(ctx context.Context, cancelFunc context.CancelFunc,
			launchersChan <-chan *scp.Launcher, cwg *sync.WaitGroup, mutex *sync.Mutex) {
			defer cwg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case launcherP, ok := <-launchersChan:
					if !ok {
						return
					}
					if err := launcherP.TryToConnect(); err == nil {
						mutex.Lock()
						hitLaunchers = append(hitLaunchers, launcherP)
						mutex.Unlock()
						cancelFunc()
					}
				}
			}
		}(ctx, cancelFunc, launchersChan, &wg, &mutex)
	}
	wg.Wait()
	cancelFunc()
	return hitLaunchers
}

func NewScpController(source string, destination string, configuration *config.MainConfig) *Controller {
	return &Controller{
		source:        source,
		destination:   destination,
		configuration: configuration,
	}
}
