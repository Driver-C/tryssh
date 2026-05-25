package delete

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/Driver-C/tryssh/pkg/utils"
	"github.com/spf13/cobra"
)

// NewUsersCommand creates and returns the cobra command for deleting a username entry.
func NewUsersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "users <username>",
		Args:    cobra.ExactArgs(1),
		Short:   "Delete an alternative username",
		Long:    "Delete an alternative username",
		Aliases: []string{"user", "usr"},
		Run: func(_ *cobra.Command, args []string) {
			username := args[0]
			configuration, err := config.LoadConfig()
			if err != nil {
				utils.Fatalln(err)
			}
			controller := control.NewDeleteController(control.TypeUsers, username, configuration)
			controller.ExecuteDelete()
		},
	}
	return cmd
}
