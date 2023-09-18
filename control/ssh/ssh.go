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

func (sc *Controller) TryLogin(user string) {
	targetServer, cacheIndex, isFound := config.SelectServerCache(user, sc.targetIp, sc.configuration)
	if user != "" {
		utils.Logger.Infof("Specify the username \"%s\" to attempt to login to the server.\n", user)
	}
	if isFound {
		utils.Logger.Infof("The cache for %s is found, which will be used to try.\n", sc.targetIp)
		sc.tryLoginWithCache(sc.targetIp, user, targetServer, cacheIndex)
	} else {
		utils.Logger.Warnf("The cache for %s could not be found. Start trying to login.\n\n", sc.targetIp)
		sc.tryLoginWithoutCache(sc.targetIp, user)
	}
	utils.Logger.Fatalln("There is no password combination that can log in successfully\n")
}

func (sc *Controller) tryLoginWithCache(targetIp string, user string,
	targetServer *config.ServerListConfig, cacheIndex int) {
	lan := &ssh.Launcher{SshConnector: *config.GetSshConnectorFromConfig(targetServer)}
	if lan.Launch() {
		os.Exit(0)
	} else {
		utils.Logger.Errorln("Failed to log in with cached information. Start trying to login again.\n\n")
		if config.DeleteServerCache(cacheIndex, sc.configuration) {
			utils.Logger.Infoln("Delete server cache successful.\n")
			sc.tryLoginWithoutCache(targetIp, user)
		}
	}
}

func (sc *Controller) tryLoginWithoutCache(targetIp string, user string) {
	combinations := config.GenerateCombination(targetIp, user, sc.configuration)
	launchers := ssh.NewSshLaunchersByCombinations(combinations)
	for _, lan := range launchers {
		if err := lan.TryToConnect(); err == nil {
			utils.Logger.Infoln("Login succeeded. The cache will be added.\n")
			newServerCache := config.GetConfigFromSshConnector(&lan.SshConnector)
			if config.AddServerCache(newServerCache, sc.configuration) {
				utils.Logger.Infoln("Cache added.\n\n")
				lan.Launch()
			} else {
				utils.Logger.Errorln("Cache added failed.\n\n")
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
