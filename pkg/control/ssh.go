package control

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/launcher"
	"github.com/Driver-C/tryssh/pkg/utils"
	"time"
)

type SshController struct {
	targetIp      string
	configuration *config.MainConfig
	cacheIsFound  bool
	cacheIndex    int
	concurrency   int
	sshTimeout    time.Duration
}

// TryLogin Functional entrance
func (sc *SshController) TryLogin(user string, concurrency int, sshTimeout time.Duration) {
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
}

func (sc *SshController) tryLoginWithCache(user string, targetServer *config.ServerListConfig) {
	lan := &launcher.SshLauncher{SshConnector: *launcher.GetSshConnectorFromConfig(targetServer)}
	// Set default timeout time
	lan.SshTimeout = sshClientTimeoutWhenLogin
	if !lan.Launch() {
		utils.Logger.Errorf("Failed to log in with cached information. Start trying to login again.\n\n")
		sc.tryLoginWithoutCache(user)
	}
}

func (sc *SshController) tryLoginWithoutCache(user string) {
	combinations := config.GenerateCombination(sc.targetIp, user, sc.configuration)
	launchers := launcher.NewSshLaunchersByCombinations(combinations, sc.sshTimeout)
	connectors := make([]launcher.Connector, len(launchers))
	for i, l := range launchers {
		connectors[i] = l
	}
	hitLaunchers := ConcurrencyTryToConnect(sc.concurrency, connectors)
	if len(hitLaunchers) > 0 {
		utils.Logger.Infoln("Login succeeded. The cache will be added.\n")
		hitLauncher := hitLaunchers[0].(*launcher.SshLauncher)
		// The new server cache information
		newServerCache := launcher.GetConfigFromSshConnector(&hitLauncher.SshConnector)
		// Determine if the login attempt was successful after the old cache login failed.
		// If so, delete the old cache information that cannot be logged in after the login attempt is successful
		if sc.cacheIsFound {
			// Sync outdated cache's alias
			newServerCache.Alias = sc.configuration.ServerLists[sc.cacheIndex].Alias

			utils.Logger.Infoln("The old cache will be deleted.\n")
			sc.configuration.ServerLists = append(
				sc.configuration.ServerLists[:sc.cacheIndex], sc.configuration.ServerLists[sc.cacheIndex+1:]...)
		}
		sc.configuration.ServerLists = append(sc.configuration.ServerLists, *newServerCache)
		if config.UpdateConfig(sc.configuration) {
			utils.Logger.Infoln("Cache added.\n\n")
			// If the timeout time is less than sshClientTimeoutWhenLogin during login,
			// change to sshClientTimeoutWhenLogin
			if hitLauncher.SshTimeout < sshClientTimeoutWhenLogin {
				hitLauncher.SshTimeout = sshClientTimeoutWhenLogin
			}
			if !hitLauncher.Launch() {
				utils.Logger.Errorf("Login failed.\n")
			}
		} else {
			utils.Logger.Errorf("Cache added failed.\n\n")
		}
	} else {
		utils.Logger.Errorf("There is no password combination that can log in.\n")
	}
}

func (sc *SshController) searchAliasExistsOrNot() {
	for _, server := range sc.configuration.ServerLists {
		if server.Alias == sc.targetIp {
			sc.targetIp = server.Ip
		}
	}
}

func NewSshController(targetIp string, configuration *config.MainConfig) *SshController {
	return &SshController{
		targetIp:      targetIp,
		configuration: configuration,
	}
}
