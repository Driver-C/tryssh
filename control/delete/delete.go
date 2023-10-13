package delete

import (
	"github.com/Driver-C/tryssh/config"
	"github.com/Driver-C/tryssh/utils"
)

const (
	typeUsers     = "users"
	typePorts     = "ports"
	typePasswords = "passwords"
	typeCaches    = "caches"
)

type Controller struct {
	deleteType    string
	deleteContent string
	configuration *config.MainConfig
}

func (dc Controller) ExecuteDelete() {
	switch dc.deleteType {
	case typeUsers:
		contents := dc.configuration.Main.Users
		if newContents := dc.searchAndDelete(contents); newContents != nil {
			dc.configuration.Main.Users = newContents
			dc.updateConfig()
		} else {
			utils.Logger.Warnf("No matching username: %s\n", dc.deleteContent)
		}
	case typePorts:
		contents := dc.configuration.Main.Ports
		if newContents := dc.searchAndDelete(contents); newContents != nil {
			dc.configuration.Main.Ports = newContents
			dc.updateConfig()
		} else {
			utils.Logger.Warnf("No matching port: %s\n", dc.deleteContent)
		}
	case typePasswords:
		contents := dc.configuration.Main.Passwords
		if newContents := dc.searchAndDelete(contents); newContents != nil {
			dc.configuration.Main.Passwords = newContents
			dc.updateConfig()
		} else {
			utils.Logger.Warnf("No matching password: %s\n", dc.deleteContent)
		}
	case typeCaches:
		// dc.deleteContent is ipAddress
		var deleteCount int
		if dc.deleteContent != "" {
			for index, server := range dc.configuration.ServerLists {
				if server.Ip == dc.deleteContent {
					dc.configuration.ServerLists = append(dc.configuration.ServerLists[:index],
						dc.configuration.ServerLists[index+1:]...)
					dc.updateConfig()
					deleteCount++
				}
			}
			if deleteCount == 0 {
				utils.Logger.Warnf("No matching cache: %s\n", dc.deleteContent)
			}
		} else {
			utils.Logger.Errorln("IP address cannot be empty characters")
		}
	}
}

func (dc Controller) searchAndDelete(contents []string) []string {
	for index, content := range contents {
		if dc.deleteContent == content {
			contents = append(contents[:index], contents[index+1:]...)
			return contents
		}
	}
	return nil
}

func (dc Controller) updateConfig() {
	if config.UpdateConfig(dc.configuration) {
		utils.Logger.Infof("Delete %s: %s completed.\n", dc.deleteType, dc.deleteContent)
	} else {
		utils.Logger.Errorf("Delete %s: %s failed.\n", dc.deleteType, dc.deleteContent)
	}
}

func NewDeleteController(deleteType string, deleteContent string,
	configuration *config.MainConfig) *Controller {
	return &Controller{
		deleteType:    deleteType,
		deleteContent: deleteContent,
		configuration: configuration,
	}
}
