package control

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/launcher"
	"github.com/Driver-C/tryssh/pkg/utils"
	"time"
)

// SSHController manages SSH login attempts using cached credentials or credential combinations.
type SSHController struct {
	targetIP      string
	configuration *config.MainConfig
	cacheIsFound  bool
	cacheIndex    int
	concurrency   int
	sshTimeout    time.Duration
}

// TryLogin attempts to log in to the target server, first using cached credentials
// and then by trying all credential combinations.
func (sc *SSHController) TryLogin(user string, concurrency int, sshTimeout time.Duration) {
	sc.sshTimeout = sshTimeout
	sc.concurrency = concurrency
	sc.targetIP = config.ResolveAlias(sc.targetIP, sc.configuration)
	var targetServer *config.ServerListConfig
	targetServer, sc.cacheIndex, sc.cacheIsFound = config.SelectServerCache(user, sc.targetIP, sc.configuration)
	if user != "" {
		utils.Infof("Specify the username \"%s\" to attempt to login to the server.\n", user)
	}
	if sc.cacheIsFound {
		utils.Infof("The cache for %s is found, which will be used to try.\n", sc.targetIP)
		sc.tryLoginWithCache(user, targetServer)
	} else {
		utils.Warnf("The cache for %s could not be found. Start trying to login.\n\n", sc.targetIP)
		sc.tryLoginWithoutCache(user)
	}
}

func (sc *SSHController) tryLoginWithCache(user string, targetServer *config.ServerListConfig) {
	lan := &launcher.SSHLauncher{SSHConnector: *launcher.GetSSHConnectorFromConfig(targetServer)}
	lan.SSHTimeout = sshClientTimeoutWhenLogin
	if !lan.Launch() {
		utils.Errorf("Failed to log in with cached information. Start trying to login again.\n\n")
		sc.tryLoginWithoutCache(user)
	}
}

func (sc *SSHController) tryLoginWithoutCache(user string) {
	combinations := config.GenerateCombination(sc.targetIP, user, sc.configuration)
	launchers := launcher.NewSSHLaunchersByCombinations(combinations, sc.sshTimeout)
	connectors := make([]launcher.Connector, len(launchers))
	for i, l := range launchers {
		connectors[i] = l
	}
	hitLaunchers := ConcurrencyTryToConnect(sc.concurrency, connectors)
	if len(hitLaunchers) > 0 {
		utils.Infoln("Login succeeded. The cache will be added.")
		hitLauncher := hitLaunchers[0].(*launcher.SSHLauncher)
		newServerCache := launcher.GetConfigFromSSHConnector(&hitLauncher.SSHConnector)
		if sc.cacheIsFound {
			newServerCache.Alias = sc.configuration.ServerLists[sc.cacheIndex].Alias
			utils.Infoln("The old cache will be deleted.")
			sc.configuration.ServerLists = append(
				sc.configuration.ServerLists[:sc.cacheIndex], sc.configuration.ServerLists[sc.cacheIndex+1:]...)
		}
		sc.configuration.ServerLists = append(sc.configuration.ServerLists, *newServerCache)
		if err := config.UpdateConfig(sc.configuration); err == nil {
			utils.Infoln("Cache added.")
			if hitLauncher.SSHTimeout < sshClientTimeoutWhenLogin {
				hitLauncher.SSHTimeout = sshClientTimeoutWhenLogin
			}
			if !hitLauncher.Launch() {
				utils.Errorf("Login failed.\n")
			}
		} else {
			utils.Errorf("Cache added failed.\n\n")
		}
	} else {
		utils.Errorf("There is no password combination that can log in.\n")
	}
}

// NewSSHController creates a new SSHController for the given target IP and configuration.
func NewSSHController(targetIP string, configuration *config.MainConfig) *SSHController {
	return &SSHController{
		targetIP:      targetIP,
		configuration: configuration,
	}
}
