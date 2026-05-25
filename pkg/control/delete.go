package control

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/utils"
)

// DeleteController handles deletion of configuration entries such as users, ports,
// passwords, keys, and server caches.
type DeleteController struct {
	deleteType    string
	deleteContent string
	configuration *config.MainConfig
}

// ExecuteDelete removes the configured entry from the main configuration.
func (dc *DeleteController) ExecuteDelete() {
	switch dc.deleteType {
	case TypeUsers:
		if newContents := dc.searchAndDelete(dc.configuration.Main.Users); newContents != nil {
			dc.configuration.Main.Users = newContents
			dc.updateConfig()
		} else {
			utils.Warnf("No matching username: %s\n", dc.deleteContent)
		}
	case TypePorts:
		if newContents := dc.searchAndDelete(dc.configuration.Main.Ports); newContents != nil {
			dc.configuration.Main.Ports = newContents
			dc.updateConfig()
		} else {
			utils.Warnf("No matching port: %s\n", dc.deleteContent)
		}
	case TypePasswords:
		if newContents := dc.searchAndDelete(dc.configuration.Main.Passwords); newContents != nil {
			dc.configuration.Main.Passwords = newContents
			dc.updateConfig()
		} else {
			utils.Warnln("No matching password")
		}
	case TypeKeys:
		if newContents := dc.searchAndDelete(dc.configuration.Main.Keys); newContents != nil {
			dc.configuration.Main.Keys = newContents
			dc.updateConfig()
		} else {
			utils.Warnf("No matching key: %s\n", dc.deleteContent)
		}
	case TypeCaches:
		if dc.deleteContent == "" {
			utils.Errorln("IP address cannot be empty characters")
			return
		}
		var deleteCount int
		for i := len(dc.configuration.ServerLists) - 1; i >= 0; i-- {
			if dc.configuration.ServerLists[i].IP == dc.deleteContent {
				dc.configuration.ServerLists = append(dc.configuration.ServerLists[:i],
					dc.configuration.ServerLists[i+1:]...)
				deleteCount++
			}
		}
		if deleteCount > 0 {
			dc.updateConfig()
		} else {
			utils.Warnf("No matching cache: %s\n", dc.deleteContent)
		}
	}
}

func (dc *DeleteController) searchAndDelete(contents []string) []string {
	for index, content := range contents {
		if dc.deleteContent == content {
			return append(contents[:index], contents[index+1:]...)
		}
	}
	return nil
}

func (dc *DeleteController) updateConfig() {
	if err := config.UpdateConfig(dc.configuration); err == nil {
		utils.Infof("Delete %s: %s completed.\n", dc.deleteType, dc.deleteContent)
	} else {
		utils.Errorf("Delete %s: %s failed.\n", dc.deleteType, dc.deleteContent)
	}
}

// NewDeleteController creates a new DeleteController for the specified type and content.
func NewDeleteController(deleteType string, deleteContent string,
	configuration *config.MainConfig) *DeleteController {
	return &DeleteController{
		deleteType:    deleteType,
		deleteContent: deleteContent,
		configuration: configuration,
	}
}
