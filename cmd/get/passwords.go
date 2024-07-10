package get

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/spf13/cobra"
)

func NewPasswordsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "passwords <password>",
		Short:   "Get alternative passwords",
		Long:    "Get alternative passwords",
		Aliases: []string{"password", "pass", "pwd"},
		Run: func(cmd *cobra.Command, args []string) {
			var password string
			if len(args) > 0 {
				password = args[0]
			}
			configuration := config.LoadConfig()
			controller := control.NewGetController(control.TypePasswords, password, configuration)
			controller.ExecuteGet()
		},
	}
	return cmd
}
