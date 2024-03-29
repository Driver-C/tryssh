package create

import (
	"encoding/json"
	"github.com/Driver-C/tryssh/config"
	"github.com/Driver-C/tryssh/utils"
)

const (
	typeUsers     = "users"
	typePorts     = "ports"
	typePasswords = "passwords"
	typeCaches    = "caches"
)

type CacheContent struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Alias    string `json:"alias"`
}

type Controller struct {
	createType    string
	createContent string
	configuration *config.MainConfig
}

func (cc Controller) ExecuteCreate() {
	switch cc.createType {
	case typeUsers:
		cc.configuration.Main.Users = utils.RemoveDuplicate(
			append(cc.configuration.Main.Users, cc.createContent))
		cc.updateConfig()
	case typePorts:
		cc.configuration.Main.Ports = utils.RemoveDuplicate(
			append(cc.configuration.Main.Ports, cc.createContent))
		cc.updateConfig()
	case typePasswords:
		cc.configuration.Main.Passwords = utils.RemoveDuplicate(
			append(cc.configuration.Main.Passwords, cc.createContent))
		cc.updateConfig()
	case typeCaches:
		cc.createCaches()
		cc.updateConfig()
	}
}

func (cc Controller) updateConfig() {
	if config.UpdateConfig(cc.configuration) {
		utils.Logger.Infof("Create %s: %s completed.\n", cc.createType, cc.createContent)
	} else {
		utils.Logger.Errorf("Create %s: %s failed.\n", cc.createType, cc.createContent)
	}
}

func (cc Controller) createCaches() {
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
	configuration *config.MainConfig) *Controller {
	return &Controller{
		createType:    createType,
		createContent: createContent,
		configuration: configuration,
	}
}
