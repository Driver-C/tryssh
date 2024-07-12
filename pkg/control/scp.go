package control

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/launcher"
	"github.com/Driver-C/tryssh/pkg/utils"
	"strings"
	"time"
)

type ScpController struct {
	source        string
	destination   string
	configuration *config.MainConfig
	cacheIsFound  bool
	cacheIndex    int
	destIp        string
	concurrency   int
	sshTimeout    time.Duration
	recursive     bool
}

// TryCopy Functional entrance
func (cc *ScpController) TryCopy(user string, concurrency int, recursive bool, sshTimeout time.Duration) {
	// Set timeout
	cc.sshTimeout = sshTimeout
	// Set concurrency
	cc.concurrency = concurrency
	// Set recursive or not
	cc.recursive = recursive
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
}

func (cc *ScpController) tryCopyWithCache(user string, targetServer *config.ServerListConfig) {
	lan := &launcher.ScpLauncher{
		SshConnector: *launcher.GetSshConnectorFromConfig(targetServer),
		Src:          cc.source,
		Dest:         cc.destination,
		Recursive:    cc.recursive,
	}
	// Set default timeout time
	lan.SshTimeout = sshClientTimeoutWhenLogin
	if !lan.Launch() {
		utils.Logger.Errorf("Failed to log in with cached information. Start trying to login again.\n\n")
		cc.tryCopyWithoutCache(user)
	}
}

func (cc *ScpController) tryCopyWithoutCache(user string) {
	combinations := config.GenerateCombination(cc.destIp, user, cc.configuration)
	launchers := launcher.NewScpLaunchersByCombinations(combinations, cc.source, cc.destination,
		cc.recursive, cc.sshTimeout)
	connectors := make([]launcher.Connector, len(launchers))
	for i, l := range launchers {
		connectors[i] = l
	}
	hitLaunchers := ConcurrencyTryToConnect(cc.concurrency, connectors)
	if len(hitLaunchers) > 0 {
		utils.Logger.Infoln("Login succeeded. The cache will be added.\n")
		hitLauncher := hitLaunchers[0].(*launcher.ScpLauncher)
		// The new server cache information
		newServerCache := launcher.GetConfigFromSshConnector(&hitLauncher.SshConnector)
		// Determine if the login attempt was successful after the old cache login failed.
		// If so, delete the old cache information that cannot be logged in after the login attempt is successful
		if cc.cacheIsFound {
			// Sync outdated cache's alias
			newServerCache.Alias = cc.configuration.ServerLists[cc.cacheIndex].Alias

			utils.Logger.Infoln("The old cache will be deleted.\n")
			cc.configuration.ServerLists = append(
				cc.configuration.ServerLists[:cc.cacheIndex], cc.configuration.ServerLists[cc.cacheIndex+1:]...)
		}
		cc.configuration.ServerLists = append(cc.configuration.ServerLists, *newServerCache)
		if config.UpdateConfig(cc.configuration) {
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

func (cc *ScpController) searchAliasExistsOrNot() {
	for _, server := range cc.configuration.ServerLists {
		if server.Alias == cc.destIp {
			cc.destIp = server.Ip
		}
	}
}

func NewScpController(source string, destination string, configuration *config.MainConfig) *ScpController {
	return &ScpController{
		source:        source,
		destination:   destination,
		configuration: configuration,
	}
}
