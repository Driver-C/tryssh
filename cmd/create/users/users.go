package users

import (
	"github.com/spf13/cobra"
	"tryssh/config"
	createControl "tryssh/control/create"
)

const createType = "users"

func NewUsersCommand() *cobra.Command {
	usersCmd := &cobra.Command{
		Use:   "users <username>",
		Args:  cobra.ExactArgs(1),
		Short: "Create a alternate username",
		Long:  "Create a alternate username",
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]
			configuration := config.LoadConfig()
			createCtl := createControl.NewCreateController(createType, username, configuration)
			createCtl.ExecuteCreate()
		},
	}
	return usersCmd
}
