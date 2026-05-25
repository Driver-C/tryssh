package control

import (
	"encoding/json"
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/utils"
)

// CacheContent represents the JSON structure used when creating a new cache entry.
type CacheContent struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Alias    string `json:"alias"`
}

// CreateController handles creation of configuration entries such as users, ports,
// passwords, keys, and server caches.
type CreateController struct {
	createType    string
	createContent string
	configuration *config.MainConfig
}

// ExecuteCreate creates the configured entry in the main configuration.
func (cc *CreateController) ExecuteCreate() {
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
		if cc.createCaches() {
			cc.updateConfig()
		}
	}
}

func (cc *CreateController) updateConfig() {
	displayContent := cc.createContent
	if cc.createType == TypePasswords {
		displayContent = utils.MaskSecret(displayContent)
	}
	if err := config.UpdateConfig(cc.configuration); err == nil {
		utils.Infof("Create %s: %s completed.\n", cc.createType, displayContent)
	} else {
		utils.Errorf("Create %s: %s failed.\n", cc.createType, displayContent)
	}
}

func (cc *CreateController) createCaches() bool {
	var newCache CacheContent
	if err := json.Unmarshal([]byte(cc.createContent), &newCache); err != nil {
		utils.Errorln("Cache's JSON unmarshal failed.")
		return false
	}
	cc.configuration.ServerLists = append(cc.configuration.ServerLists,
		config.ServerListConfig{
			IP:       newCache.IP,
			Port:     newCache.Port,
			User:     newCache.User,
			Password: newCache.Password,
			Alias:    newCache.Alias,
		},
	)
	return true
}

// NewCreateController creates a new CreateController for the specified type and content.
func NewCreateController(createType string, createContent string,
	configuration *config.MainConfig) *CreateController {
	return &CreateController{
		createType:    createType,
		createContent: createContent,
		configuration: configuration,
	}
}
