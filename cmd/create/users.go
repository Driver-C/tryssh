package create

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/spf13/cobra"
)

func NewUsersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "users <username>",
		Args:    cobra.ExactArgs(1),
		Short:   "Create an alternative username",
		Long:    "Create an alternative username",
		Aliases: []string{"user", "usr"},
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]
			configuration := config.LoadConfig()
			controller := control.NewCreateController(control.TypeUsers, username, configuration)
			controller.ExecuteCreate()
		},
	}
	return cmd
}
