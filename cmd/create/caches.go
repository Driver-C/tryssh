package create

import (
	"encoding/json"
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/Driver-C/tryssh/pkg/utils"
	"github.com/spf13/cobra"
)

func NewCachesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "caches <cache>",
		Short:   "Create an alternative cache",
		Long:    "Create an alternative cache",
		Aliases: []string{"cache"},
		Run: func(cmd *cobra.Command, args []string) {
			newIp, _ := cmd.Flags().GetString("ip")
			newUser, _ := cmd.Flags().GetString("user")
			newPort, _ := cmd.Flags().GetString("port")
			newPassword, _ := cmd.Flags().GetString("pwd")
			newAlias, _ := cmd.Flags().GetString("alias")
			newCacheContent := control.CacheContent{
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
			controller := control.NewCreateController(control.TypeCaches, string(contentJson), configuration)
			controller.ExecuteCreate()
		},
	}
	cmd.Flags().StringP("ip", "i", "", "The ipaddress of the cache to be added")
	cmd.Flags().StringP("user", "u", "", "The username of the cache to be added")
	cmd.Flags().StringP("port", "P", "", "The port of the cache to be added")
	cmd.Flags().StringP("pwd", "p", "", "The password of the cache to be added")
	cmd.Flags().StringP("alias", "a", "", "The alias of the cache to be added")

	if err := cmd.MarkFlagRequired("ip"); err != nil {
		utils.Logger.Errorln("Flag: ip must be set.")
		return nil
	}
	if err := cmd.MarkFlagRequired("user"); err != nil {
		utils.Logger.Errorln("Flag: user must be set.")
		return nil
	}
	if err := cmd.MarkFlagRequired("port"); err != nil {
		utils.Logger.Errorln("Flag: port must be set.")
		return nil
	}
	if err := cmd.MarkFlagRequired("pwd"); err != nil {
		utils.Logger.Errorln("Flag: password must be set.")
		return nil
	}
	return cmd
}
