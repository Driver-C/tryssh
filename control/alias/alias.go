package alias

import (
	"tryssh/config"
	"tryssh/utils"
)

type Controller struct {
	targetIp      string
	configuration *config.MainConfig
	alias         string
}

func (ac *Controller) SetAlias() {
	var beSetCount int
	for index, server := range ac.configuration.ServerLists {
		if server.Ip == ac.targetIp {
			aliasServerList := ac.getServerListFromAlias()
			if len(aliasServerList) != 0 {
				ac.ListAlias()
				utils.Logger.Fatalf(
					"The alias \"%s\" has already been set, try another alias or delete it and set again.\n",
					ac.alias)
			}
			ac.configuration.ServerLists[index].Alias = ac.alias
			utils.Logger.Infof(
				"The server %s@%s:%s's alias \"%s\" will be set.\n",
				server.User, ac.targetIp, server.Port, ac.alias)
			beSetCount++
		}
	}
	if config.UpdateConfig(ac.configuration) {
		utils.Logger.Infof("%d cache infomation has been changed.\n", beSetCount)
	} else {
		utils.Logger.Fatalln("Main config update failed.")
	}
}

func (ac *Controller) ListAlias() {
	var aliasCount int
	for _, server := range ac.configuration.ServerLists {
		if ac.alias == "" {
			if server.Alias != "" {
				utils.Logger.Infof("Alias: %s Server: %s\n", server.Alias, server.Ip)
				aliasCount++
			}
		} else {
			if server.Alias == ac.alias {
				utils.Logger.Infof("Alias: %s Server: %s\n", server.Alias, server.Ip)
				aliasCount++
			}
		}
	}
	if aliasCount == 0 {
		utils.Logger.Infoln("No aliases were found that have been set.")
	}
}

func (ac *Controller) UnsetAlias() {
	var beUnsetCount int
	for index, server := range ac.configuration.ServerLists {
		if server.Alias == ac.alias {
			ac.configuration.ServerLists[index].Alias = ""
			utils.Logger.Infof(
				"The server %s@%s:%s's alias \"%s\" will be unset.\n",
				server.User, ac.targetIp, server.Port, ac.alias)
			beUnsetCount++
		}
	}
	if config.UpdateConfig(ac.configuration) {
		utils.Logger.Infof("%d cache infomation has been changed.\n", beUnsetCount)
	} else {
		utils.Logger.Fatalln("Main config update failed.")
	}
}

func (ac *Controller) getServerListFromAlias() []config.ServerListConfig {
	var aliasServerList []config.ServerListConfig
	for _, server := range ac.configuration.ServerLists {
		if server.Alias == ac.alias && ac.alias != "" {
			aliasServerList = append(aliasServerList, server)
		}
	}
	return aliasServerList
}

func NewAliasController(targetIp string, configuration *config.MainConfig, alias string) *Controller {
	return &Controller{
		targetIp:      targetIp,
		configuration: configuration,
		alias:         alias,
	}
}
