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
}

// TryCopy Functional entrance
func (cc *Controller) TryCopy(user string) {
	var destIp string
	if strings.Contains(cc.source, ":") {
		destIp = strings.Split(cc.source, ":")[0]
	} else if strings.Contains(cc.destination, ":") {
		destIp = strings.Split(cc.destination, ":")[0]
	}
	var targetServer *config.ServerListConfig
	targetServer, cc.cacheIndex, cc.cacheIsFound = config.SelectServerCache(user, destIp, cc.configuration)

	if cc.cacheIsFound {
		utils.Logger.Infof("The cache for %s is found, which will be used to try.\n", destIp)
		cc.tryCopyWithCache(destIp, user, targetServer)
	} else {
		utils.Logger.Warnf("The cache for %s could not be found. Start trying to login.\n\n", destIp)
		cc.tryCopyWithoutCache(destIp, user)
	}
	utils.Logger.Fatalln("There is no password combination that can log in successfully\n")
}

func (cc *Controller) tryCopyWithCache(destIp string, user string, targetServer *config.ServerListConfig) {
	lan := &scp.Launcher{
		SshConnector: *config.GetSshConnectorFromConfig(targetServer),
		Src:          cc.source,
		Dest:         cc.destination,
	}
	if lan.Launch() {
		os.Exit(0)
	} else {
		utils.Logger.Errorln("Failed to log in with cached information. Start trying to login again.\n\n")
		cc.tryCopyWithoutCache(destIp, user)
	}
}

func (cc *Controller) tryCopyWithoutCache(destIp string, user string) {
	combinations := config.GenerateCombination(destIp, user, cc.configuration)
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

func NewScpController(source string, destination string, configuration *config.MainConfig) *Controller {
	return &Controller{
		source:        source,
		destination:   destination,
		configuration: configuration,
	}
}
