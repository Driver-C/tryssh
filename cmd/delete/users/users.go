package users

import (
	"github.com/Driver-C/tryssh/config"
	"github.com/Driver-C/tryssh/control/delete"
	"github.com/spf13/cobra"
)

const deleteType = "users"

func NewUsersCommand() *cobra.Command {
	usersCmd := &cobra.Command{
		Use:   "users <username>",
		Args:  cobra.ExactArgs(1),
		Short: "Delete a alternate username",
		Long:  "Delete a alternate username",
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]
			configuration := config.LoadConfig()
			deleteCtl := delete.NewDeleteController(deleteType, username, configuration)
			deleteCtl.ExecuteDelete()
		},
	}
	return usersCmd
}
