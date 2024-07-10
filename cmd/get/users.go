package get

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/spf13/cobra"
)

func NewUsersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "users <username>",
		Short:   "Get alternative usernames",
		Long:    "Get alternative usernames",
		Aliases: []string{"user", "usr"},
		Run: func(cmd *cobra.Command, args []string) {
			var username string
			if len(args) > 0 {
				username = args[0]
			}
			configuration := config.LoadConfig()
			controller := control.NewGetController(control.TypeUsers, username, configuration)
			controller.ExecuteGet()
		},
	}
	return cmd
}
