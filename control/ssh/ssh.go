package ssh

import (
	"os"
	"tryssh/config"
	"tryssh/launcher/ssh"
	"tryssh/utils"
)

type Controller struct {
	targetIp      string
	configuration *config.MainConfig
}

func (sc *Controller) TryLogin() {
	targetServer, cacheIndex, isFound := config.SelectServerCache(sc.targetIp, sc.configuration)

	if isFound {
		utils.Logger.Infof("The cache for %s is found, which will be used to try.\n", sc.targetIp)
		sc.tryLoginWithCache(sc.targetIp, targetServer, cacheIndex)
	} else {
		utils.Logger.Warnf("The cache for %s could not be found. The password will be guessed.\n\n", sc.targetIp)
		sc.tryLoginWithoutCache(sc.targetIp)
	}
	utils.Logger.Fatalln("There is no password combination that can log in successfully\n")
}

func (sc *Controller) tryLoginWithCache(targetIp string, targetServer *config.ServerListConfig, cacheIndex int) {
	lan := &ssh.Launcher{SshConnector: *config.GetSshConnectorFromConfig(targetServer)}
	if lan.Launch() {
		os.Exit(0)
	} else {
		utils.Logger.Errorln("Failed to log in with cached information. The password will be guessed again\n\n")
		if config.DeleteServerCache(cacheIndex, sc.configuration) {
			utils.Logger.Infoln("Delete server cache successful.\n")
			sc.tryLoginWithoutCache(targetIp)
		}
	}
}

func (sc *Controller) tryLoginWithoutCache(targetIp string) {
	combinations := config.GenerateCombination(targetIp, sc.configuration)
	launchers := ssh.NewSshLaunchersByCombinations(combinations)
	for _, lan := range launchers {
		if err := lan.TryToConnect(); err == nil {
			utils.Logger.Infoln("Login succeeded. The cache will be added.\n")
			if config.AddServerCache(config.GetConfigFromSshConnector(&lan.SshConnector), sc.configuration) {
				utils.Logger.Infoln("Cache updated.\n\n")
				lan.Launch()
			} else {
				utils.Logger.Errorln("Cache update failed.\n\n")
			}
			os.Exit(0)
		}
	}
}

func NewSshController(targetIp string, configuration *config.MainConfig) *Controller {
	return &Controller{
		targetIp:      targetIp,
		configuration: configuration,
	}
}
