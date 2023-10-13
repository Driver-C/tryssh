package users

import (
	"github.com/Driver-C/tryssh/config"
	"github.com/Driver-C/tryssh/control/get"
	"github.com/spf13/cobra"
)

const getType = "users"

func NewUsersCommand() *cobra.Command {
	usersCmd := &cobra.Command{
		Use:   "users <username>",
		Short: "Get alternate usernames",
		Long:  "Get alternate usernames",
		Run: func(cmd *cobra.Command, args []string) {
			var username string
			if len(args) > 0 {
				username = args[0]
			}
			configuration := config.LoadConfig()
			getCtl := get.NewGetController(getType, username, configuration)
			getCtl.ExecuteGet()
		},
	}
	return usersCmd
}
