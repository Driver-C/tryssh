package caches

import (
	"encoding/json"
	"github.com/Driver-C/tryssh/config"
	"github.com/Driver-C/tryssh/control/create"
	"github.com/Driver-C/tryssh/utils"
	"github.com/spf13/cobra"
)

const createType = "caches"

func NewCachesCommand() *cobra.Command {
	cachesCmd := &cobra.Command{
		Use:   "caches <cache>",
		Short: "Create a alternate cache",
		Long:  "Create a alternate cache",
		Run: func(cmd *cobra.Command, args []string) {
			newIp, _ := cmd.Flags().GetString("ip")
			newUser, _ := cmd.Flags().GetString("user")
			newPort, _ := cmd.Flags().GetString("port")
			newPassword, _ := cmd.Flags().GetString("pwd")
			newAlias, _ := cmd.Flags().GetString("alias")
			newCacheContent := create.CacheContent{
				Ip:       newIp,
				User:     newUser,
				Port:     newPort,
				Password: newPassword,
				Alias:    newAlias,
			}
			contentJson, err := json.Marshal(newCacheContent)
			if err != nil {
				utils.Logger.Errorln("Cache content JSON marshal failed.")
				return
			}
			configuration := config.LoadConfig()
			createCtl := create.NewCreateController(createType, string(contentJson), configuration)
			createCtl.ExecuteCreate()
		},
	}
	cachesCmd.Flags().StringP("ip", "i", "", "The ipaddress of the cache to be added")
	cachesCmd.Flags().StringP("user", "u", "", "The username of the cache to be added")
	cachesCmd.Flags().StringP("port", "P", "", "The port of the cache to be added")
	cachesCmd.Flags().StringP("pwd", "p", "", "The password of the cache to be added")
	cachesCmd.Flags().StringP("alias", "a", "", "The alias of the cache to be added")

	if err := cachesCmd.MarkFlagRequired("ip"); err != nil {
		utils.Logger.Errorln("Flag: ip must be set.")
		return nil
	}
	if err := cachesCmd.MarkFlagRequired("user"); err != nil {
		utils.Logger.Errorln("Flag: user must be set.")
		return nil
	}
	if err := cachesCmd.MarkFlagRequired("port"); err != nil {
		utils.Logger.Errorln("Flag: port must be set.")
		return nil
	}
	if err := cachesCmd.MarkFlagRequired("pwd"); err != nil {
		utils.Logger.Errorln("Flag: password must be set.")
		return nil
	}
	return cachesCmd
}
