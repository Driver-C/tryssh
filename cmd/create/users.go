package create

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/spf13/cobra"
)

func NewUsersCommand() *cobra.Command {
	usersCmd := &cobra.Command{
		Use:     "users <username>",
		Args:    cobra.ExactArgs(1),
		Short:   "Create a alternate username",
		Long:    "Create a alternate username",
		Aliases: []string{"user", "usr"},
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]
			configuration := config.LoadConfig()
			createCtl := control.NewCreateController(control.TypeUsers, username, configuration)
			createCtl.ExecuteCreate()
		},
	}
	return usersCmd
}
