package delete

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/spf13/cobra"
)

func NewUsersCommand() *cobra.Command {
	usersCmd := &cobra.Command{
		Use:     "users <username>",
		Args:    cobra.ExactArgs(1),
		Short:   "Delete a alternate username",
		Long:    "Delete a alternate username",
		Aliases: []string{"user", "usr"},
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]
			configuration := config.LoadConfig()
			deleteCtl := control.NewDeleteController(control.TypeUsers, username, configuration)
			deleteCtl.ExecuteDelete()
		},
	}
	return usersCmd
}
