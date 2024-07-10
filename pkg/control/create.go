package control

import (
	"encoding/json"
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/utils"
)

type CacheContent struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Alias    string `json:"alias"`
}

type CreateController struct {
	createType    string
	createContent string
	configuration *config.MainConfig
}

func (cc CreateController) ExecuteCreate() {
	switch cc.createType {
	case TypeUsers:
		cc.configuration.Main.Users = utils.RemoveDuplicate(
			append(cc.configuration.Main.Users, cc.createContent))
		cc.updateConfig()
	case TypePorts:
		cc.configuration.Main.Ports = utils.RemoveDuplicate(
			append(cc.configuration.Main.Ports, cc.createContent))
		cc.updateConfig()
	case TypePasswords:
		cc.configuration.Main.Passwords = utils.RemoveDuplicate(
			append(cc.configuration.Main.Passwords, cc.createContent))
		cc.updateConfig()
	case TypeKeys:
		cc.configuration.Main.Keys = utils.RemoveDuplicate(
			append(cc.configuration.Main.Keys, cc.createContent))
		cc.updateConfig()
	case TypeCaches:
		cc.createCaches()
		cc.updateConfig()
	}
}

func (cc CreateController) updateConfig() {
	if config.UpdateConfig(cc.configuration) {
		utils.Logger.Infof("Create %s: %s completed.\n", cc.createType, cc.createContent)
	} else {
		utils.Logger.Errorf("Create %s: %s failed.\n", cc.createType, cc.createContent)
	}
}

func (cc CreateController) createCaches() {
	var newCache CacheContent
	if err := json.Unmarshal([]byte(cc.createContent), &newCache); err == nil {
		cc.configuration.ServerLists = append(cc.configuration.ServerLists,
			config.ServerListConfig{
				Ip:       newCache.Ip,
				Port:     newCache.Port,
				User:     newCache.User,
				Password: newCache.Password,
				Alias:    newCache.Alias,
			},
		)
	} else {
		utils.Logger.Errorln("Cache's JSON unmarshal failed.")
	}
}

func NewCreateController(createType string, createContent string,
	configuration *config.MainConfig) *CreateController {
	return &CreateController{
		createType:    createType,
		createContent: createContent,
		configuration: configuration,
	}
}
