package ssh

import (
	"context"
	"os"
	"sync"
	"time"
	"tryssh/config"
	"tryssh/launcher/ssh"
	"tryssh/utils"
)

type Controller struct {
	targetIp      string
	configuration *config.MainConfig
	cacheIsFound  bool
	cacheIndex    int
	concurrency   int
	sshTimeout    time.Duration
}

// TryLogin Functional entrance
func (sc *Controller) TryLogin(user string, concurrency int, sshTimeout time.Duration) {
	// Set timeout
	sc.sshTimeout = sshTimeout
	// Set concurrency
	sc.concurrency = concurrency
	// Obtain the real address based on the alias
	sc.searchAliasExistsOrNot()
	var targetServer *config.ServerListConfig
	targetServer, sc.cacheIndex, sc.cacheIsFound = config.SelectServerCache(user, sc.targetIp, sc.configuration)
	if user != "" {
		utils.Logger.Infof("Specify the username \"%s\" to attempt to login to the server.\n", user)
	}
	if sc.cacheIsFound {
		utils.Logger.Infof("The cache for %s is found, which will be used to try.\n", sc.targetIp)
		sc.tryLoginWithCache(user, targetServer)
	} else {
		utils.Logger.Warnf("The cache for %s could not be found. Start trying to login.\n\n", sc.targetIp)
		sc.tryLoginWithoutCache(user)
	}
	utils.Logger.Fatalln("There is no password combination that can log in successfully\n")
}

func (sc *Controller) tryLoginWithCache(user string, targetServer *config.ServerListConfig) {
	lan := &ssh.Launcher{SshConnector: *config.GetSshConnectorFromConfig(targetServer)}
	if lan.Launch() {
		os.Exit(0)
	} else {
		utils.Logger.Errorln("Failed to log in with cached information. Start trying to login again.\n\n")
		sc.tryLoginWithoutCache(user)
	}
}

func (sc *Controller) tryLoginWithoutCache(user string) {
	combinations := config.GenerateCombination(sc.targetIp, user, sc.configuration)
	launchers := ssh.NewSshLaunchersByCombinations(combinations, sc.sshTimeout)
	hitLaunchers := concurrencyTryToConnect(sc.concurrency, launchers)
	if hitLaunchers != nil {
		utils.Logger.Infoln("Login succeeded. The cache will be added.\n")
		// Determine if the login attempt was successful after the old cache login failed.
		// If so, delete the old cache information that cannot be logged in after the login attempt is successful
		if sc.cacheIsFound {
			utils.Logger.Infoln("The old cache will be deleted.\n")
			sc.configuration.ServerLists = append(
				sc.configuration.ServerLists[:sc.cacheIndex], sc.configuration.ServerLists[sc.cacheIndex+1:]...)
		}
		newServerCache := config.GetConfigFromSshConnector(&hitLaunchers[0].SshConnector)
		sc.configuration.ServerLists = append(sc.configuration.ServerLists, *newServerCache)
		if config.UpdateConfig(sc.configuration) {
			utils.Logger.Infoln("Cache added.\n\n")
			hitLaunchers[0].Launch()
		} else {
			utils.Logger.Errorln("Cache added failed.\n\n")
		}
		os.Exit(0)
	}
}

func (sc *Controller) searchAliasExistsOrNot() {
	for _, server := range sc.configuration.ServerLists {
		if server.Alias == sc.targetIp {
			sc.targetIp = server.Ip
		}
	}
}

func concurrencyTryToConnect(concurrency int, launchers []*ssh.Launcher) []*ssh.Launcher {
	var hitLaunchers []*ssh.Launcher
	var mutex sync.Mutex
	launchersChan := make(chan *ssh.Launcher)
	ctx, cancelFunc := context.WithCancel(context.Background())
	// Producer
	go func(ctx context.Context, launchersChan chan<- *ssh.Launcher, launchers []*ssh.Launcher) {
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
			launchersChan <-chan *ssh.Launcher, cwg *sync.WaitGroup, mutex *sync.Mutex) {
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

func NewSshController(targetIp string, configuration *config.MainConfig) *Controller {
	return &Controller{
		targetIp:      targetIp,
		configuration: configuration,
	}
}
