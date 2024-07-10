package get

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/spf13/cobra"
)

func NewUsersCommand() *cobra.Command {
	usersCmd := &cobra.Command{
		Use:     "users <username>",
		Short:   "Get alternate usernames",
		Long:    "Get alternate usernames",
		Aliases: []string{"user", "usr"},
		Run: func(cmd *cobra.Command, args []string) {
			var username string
			if len(args) > 0 {
				username = args[0]
			}
			configuration := config.LoadConfig()
			getCtl := control.NewGetController(control.TypeUsers, username, configuration)
			getCtl.ExecuteGet()
		},
	}
	return usersCmd
}
