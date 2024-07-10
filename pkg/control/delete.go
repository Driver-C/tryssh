package control

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/utils"
)

type DeleteController struct {
	deleteType    string
	deleteContent string
	configuration *config.MainConfig
}

func (dc DeleteController) ExecuteDelete() {
	switch dc.deleteType {
	case TypeUsers:
		contents := dc.configuration.Main.Users
		if newContents := dc.searchAndDelete(contents); newContents != nil {
			dc.configuration.Main.Users = newContents
			dc.updateConfig()
		} else {
			utils.Logger.Warnf("No matching username: %s\n", dc.deleteContent)
		}
	case TypePorts:
		contents := dc.configuration.Main.Ports
		if newContents := dc.searchAndDelete(contents); newContents != nil {
			dc.configuration.Main.Ports = newContents
			dc.updateConfig()
		} else {
			utils.Logger.Warnf("No matching port: %s\n", dc.deleteContent)
		}
	case TypePasswords:
		contents := dc.configuration.Main.Passwords
		if newContents := dc.searchAndDelete(contents); newContents != nil {
			dc.configuration.Main.Passwords = newContents
			dc.updateConfig()
		} else {
			utils.Logger.Warnf("No matching password: %s\n", dc.deleteContent)
		}
	case TypeKeys:
		contents := dc.configuration.Main.Keys
		if newContents := dc.searchAndDelete(contents); newContents != nil {
			dc.configuration.Main.Keys = newContents
			dc.updateConfig()
		} else {
			utils.Logger.Warnf("No matching key: %s\n", dc.deleteContent)
		}
	case TypeCaches:
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

func (dc DeleteController) searchAndDelete(contents []string) []string {
	for index, content := range contents {
		if dc.deleteContent == content {
			contents = append(contents[:index], contents[index+1:]...)
			return contents
		}
	}
	return nil
}

func (dc DeleteController) updateConfig() {
	if config.UpdateConfig(dc.configuration) {
		utils.Logger.Infof("Delete %s: %s completed.\n", dc.deleteType, dc.deleteContent)
	} else {
		utils.Logger.Errorf("Delete %s: %s failed.\n", dc.deleteType, dc.deleteContent)
	}
}

func NewDeleteController(deleteType string, deleteContent string,
	configuration *config.MainConfig) *DeleteController {
	return &DeleteController{
		deleteType:    deleteType,
		deleteContent: deleteContent,
		configuration: configuration,
	}
}
