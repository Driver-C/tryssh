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
	cacheIsFound  bool
	cacheIndex    int
}

// TryLogin Functional entrance
func (sc *Controller) TryLogin(user string) {
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
	launchers := ssh.NewSshLaunchersByCombinations(combinations)
	for _, lan := range launchers {
		if err := lan.TryToConnect(); err == nil {
			utils.Logger.Infoln("Login succeeded. The cache will be added.\n")
			// Determine if the login attempt was successful after the old cache login failed.
			// If so, delete the old cache information that cannot be logged in after the login attempt is successful
			if sc.cacheIsFound {
				utils.Logger.Infoln("The old cache will be deleted.\n")
				config.DeleteServerCache(sc.cacheIndex, sc.configuration)
			}
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

func (sc *Controller) searchAliasExistsOrNot() {
	for _, server := range sc.configuration.ServerLists {
		if server.Alias == sc.targetIp {
			sc.targetIp = server.Ip
		}
	}
}

func NewSshController(targetIp string, configuration *config.MainConfig) *Controller {
	return &Controller{
		targetIp:      targetIp,
		configuration: configuration,
	}
}
