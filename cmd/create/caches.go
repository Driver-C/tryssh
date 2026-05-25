package create

import (
	"encoding/json"
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/Driver-C/tryssh/pkg/utils"
	"github.com/spf13/cobra"
)

// NewCachesCommand creates and returns the cobra command for creating a cache entry.
func NewCachesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "caches <cache>",
		Short:   "Create an alternative cache",
		Long:    "Create an alternative cache",
		Aliases: []string{"cache"},
		Run: func(cmd *cobra.Command, _ []string) {
			newIP, _ := cmd.Flags().GetString("ip")
			newUser, _ := cmd.Flags().GetString("user")
			newPort, _ := cmd.Flags().GetString("port")
			newPassword, _ := cmd.Flags().GetString("pwd")
			newAlias, _ := cmd.Flags().GetString("alias")
			newCacheContent := control.CacheContent{
				IP:       newIP,
				User:     newUser,
				Port:     newPort,
				Password: newPassword,
				Alias:    newAlias,
			}
			contentJSON, err := json.Marshal(newCacheContent) //nolint:gosec // G117: password is needed for cache storage
			if err != nil {
				utils.Errorln("Cache content JSON marshal failed.")
				return
			}
			configuration, loadErr := config.LoadConfig()
			if loadErr != nil {
				utils.Fatalln(loadErr)
			}
			controller := control.NewCreateController(control.TypeCaches, string(contentJSON), configuration)
			controller.ExecuteCreate()
		},
	}
	cmd.Flags().StringP("ip", "i", "", "The ipaddress of the cache to be added")
	cmd.Flags().StringP("user", "u", "", "The username of the cache to be added")
	cmd.Flags().StringP("port", "P", "", "The port of the cache to be added")
	cmd.Flags().StringP("pwd", "p", "", "The password of the cache to be added")
	cmd.Flags().StringP("alias", "a", "", "The alias of the cache to be added")

	_ = cmd.MarkFlagRequired("ip")
	_ = cmd.MarkFlagRequired("user")
	_ = cmd.MarkFlagRequired("port")
	_ = cmd.MarkFlagRequired("pwd")
	return cmd
}
