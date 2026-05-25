package get

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/Driver-C/tryssh/pkg/utils"
	"github.com/spf13/cobra"
)

// NewUsersCommand creates and returns the cobra command for retrieving username entries.
func NewUsersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "users <username>",
		Short:   "Get alternative usernames",
		Long:    "Get alternative usernames",
		Aliases: []string{"user", "usr"},
		Run: func(_ *cobra.Command, args []string) {
			var username string
			if len(args) > 0 {
				username = args[0]
			}
			configuration, err := config.LoadConfig()
			if err != nil {
				utils.Fatalln(err)
			}
			controller := control.NewGetController(control.TypeUsers, username, configuration)
			controller.ExecuteGet()
		},
	}
	return cmd
}
