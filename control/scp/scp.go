package scp

import (
	"os"
	"strings"
	"tryssh/config"
	"tryssh/launcher/scp"
	"tryssh/utils"
)

type Controller struct {
	source        string
	destination   string
	configuration *config.MainConfig
	cacheIsFound  bool
	cacheIndex    int
	destIp        string
}

// TryCopy Functional entrance
func (cc *Controller) TryCopy(user string) {
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
	if lan.Launch() {
		os.Exit(0)
	} else {
		utils.Logger.Errorln("Failed to log in with cached information. Start trying to login again.\n\n")
		cc.tryCopyWithoutCache(user)
	}
}

func (cc *Controller) tryCopyWithoutCache(user string) {
	combinations := config.GenerateCombination(cc.destIp, user, cc.configuration)
	launchers := scp.NewScpLaunchersByCombinations(combinations, cc.source, cc.destination)
	for _, lan := range launchers {
		if err := lan.TryToConnect(); err == nil {
			utils.Logger.Infoln("Login succeeded. The cache will be added.\n")
			// Determine if the login attempt was successful after the old cache login failed.
			// If so, delete the old cache information that cannot be logged in after the login attempt is successful
			if cc.cacheIsFound {
				utils.Logger.Infoln("The old cache will be deleted.\n")
				config.DeleteServerCache(cc.cacheIndex, cc.configuration)
			}
			newServerCache := config.GetConfigFromSshConnector(&lan.SshConnector)
			if config.AddServerCache(newServerCache, cc.configuration) {
				utils.Logger.Infoln("Cache added.\n\n")
				lan.Launch()
			} else {
				utils.Logger.Errorln("Cache added failed.\n\n")
			}
			os.Exit(0)
		}
	}
}

func (cc *Controller) searchAliasExistsOrNot() {
	for _, server := range cc.configuration.ServerLists {
		if server.Alias == cc.destIp {
			cc.destIp = server.Ip
		}
	}
}

func NewScpController(source string, destination string, configuration *config.MainConfig) *Controller {
	return &Controller{
		source:        source,
		destination:   destination,
		configuration: configuration,
	}
}
