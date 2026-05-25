// Package control implements the business logic for SSH/SCP operations.
package control

import (
	"fmt"

	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/utils"
)

// AliasController manages alias operations for server cache entries.
type AliasController struct {
	targetIP      string
	configuration *config.MainConfig
	alias         string
}

// SetAlias assigns the configured alias to all server entries matching the target IP.
func (ac *AliasController) SetAlias() bool {
	aliasServerList := config.FindAlias(ac.alias, ac.configuration)
	if len(aliasServerList) != 0 {
		ac.ListAlias()
		utils.Errorf(
			"The alias \"%s\" has already been set, try another alias or delete it and set again.\n",
			ac.alias)
		return false
	}
	var beSetCount int
	for index, server := range ac.configuration.ServerLists {
		if server.IP == ac.targetIP {
			ac.configuration.ServerLists[index].Alias = ac.alias
			utils.Infof(
				"The server %s@%s:%s's alias \"%s\" will be set.\n",
				server.User, server.IP, server.Port, ac.alias)
			beSetCount++
		}
	}
	if beSetCount == 0 {
		utils.Warnf("No matching server for IP: %s\n", ac.targetIP)
		return false
	}
	if err := config.UpdateConfig(ac.configuration); err == nil {
		utils.Infof("%d cache information has been changed.\n", beSetCount)
		return true
	}
	utils.Errorln("Main config update failed.")
	return false
}

// ListAlias prints all configured aliases, or only those matching the configured alias filter.
func (ac *AliasController) ListAlias() {
	var aliasCount int
	for _, server := range ac.configuration.ServerLists {
		if ac.alias == "" {
			if server.Alias != "" {
				fmt.Printf("Alias: %s\tServer: %s\n", server.Alias, server.IP)
				aliasCount++
			}
		} else {
			if server.Alias == ac.alias {
				fmt.Printf("Alias: %s\tServer: %s\n", server.Alias, server.IP)
				aliasCount++
			}
		}
	}
	if aliasCount == 0 {
		utils.Infoln("No aliases were found that have been set.")
	}
}

// UnsetAlias removes the configured alias from all matching server entries.
func (ac *AliasController) UnsetAlias() bool {
	var beUnsetCount int
	for index, server := range ac.configuration.ServerLists {
		if server.Alias == ac.alias {
			ac.configuration.ServerLists[index].Alias = ""
			utils.Infof(
				"The server %s@%s:%s's alias \"%s\" will be unset.\n",
				server.User, server.IP, server.Port, ac.alias)
			beUnsetCount++
		}
	}
	if beUnsetCount == 0 {
		utils.Warnf("No matching alias: %s\n", ac.alias)
		return false
	}
	if err := config.UpdateConfig(ac.configuration); err == nil {
		utils.Infof("%d cache information has been changed.\n", beUnsetCount)
		return true
	}
	utils.Errorln("Main config update failed.")
	return false
}

// NewAliasController creates a new AliasController for the given target IP, configuration, and alias.
func NewAliasController(targetIP string, configuration *config.MainConfig, alias string) *AliasController {
	return &AliasController{
		targetIP:      targetIP,
		configuration: configuration,
		alias:         alias,
	}
}
