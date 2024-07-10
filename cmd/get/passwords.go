package get

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/spf13/cobra"
)

func NewPasswordsCommand() *cobra.Command {
	passwordsCmd := &cobra.Command{
		Use:     "passwords <password>",
		Short:   "Get alternate passwords",
		Long:    "Get alternate passwords",
		Aliases: []string{"password", "pass", "pwd"},
		Run: func(cmd *cobra.Command, args []string) {
			var password string
			if len(args) > 0 {
				password = args[0]
			}
			configuration := config.LoadConfig()
			getCtl := control.NewGetController(control.TypePasswords, password, configuration)
			getCtl.ExecuteGet()
		},
	}
	return passwordsCmd
}
