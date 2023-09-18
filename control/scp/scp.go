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
}

func (cc *Controller) TryCopy(user string) {
	var destIp string
	if strings.Contains(cc.source, ":") {
		destIp = strings.Split(cc.source, ":")[0]
	} else if strings.Contains(cc.destination, ":") {
		destIp = strings.Split(cc.destination, ":")[0]
	}
	targetServer, cacheIndex, isFound := config.SelectServerCache(user, destIp, cc.configuration)

	if isFound {
		utils.Logger.Infof("The cache for %s is found, which will be used to try.\n", destIp)
		cc.tryCopyWithCache(destIp, user, targetServer, cacheIndex)
	} else {
		utils.Logger.Warnf("The cache for %s could not be found. Start trying to login.\n\n", destIp)
		cc.tryCopyWithoutCache(destIp, user)
	}
	utils.Logger.Fatalln("There is no password combination that can log in successfully\n")
}

func (cc *Controller) tryCopyWithCache(destIp string, user string,
	targetServer *config.ServerListConfig, cacheIndex int) {
	lan := &scp.Launcher{
		SshConnector: *config.GetSshConnectorFromConfig(targetServer),
		Src:          cc.source,
		Dest:         cc.destination,
	}
	if lan.Launch() {
		os.Exit(0)
	} else {
		utils.Logger.Errorln("Failed to log in with cached information. Start trying to login again.\n\n")
		if config.DeleteServerCache(cacheIndex, cc.configuration) {
			utils.Logger.Infoln("Delete server cache successful.\n")
			cc.tryCopyWithoutCache(destIp, user)
		}
	}
}

func (cc *Controller) tryCopyWithoutCache(destIp string, user string) {
	combinations := config.GenerateCombination(destIp, user, cc.configuration)
	launchers := scp.NewScpLaunchersByCombinations(combinations, cc.source, cc.destination)
	for _, lan := range launchers {
		if err := lan.TryToConnect(); err == nil {
			utils.Logger.Infoln("Login succeeded. The cache will be added.\n")
			if config.AddServerCache(config.GetConfigFromSshConnector(&lan.SshConnector), cc.configuration) {
				utils.Logger.Infoln("Cache updated.\n\n")
				lan.Launch()
			} else {
				utils.Logger.Errorln("Cache update failed.\n\n")
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
