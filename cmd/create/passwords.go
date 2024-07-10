package create

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/spf13/cobra"
)

func NewPasswordsCommand() *cobra.Command {
	passwordsCmd := &cobra.Command{
		Use:     "passwords <password>",
		Args:    cobra.ExactArgs(1),
		Short:   "Create a alternate password",
		Long:    "Create a alternate password",
		Aliases: []string{"password", "pass", "pwd"},
		Run: func(cmd *cobra.Command, args []string) {
			password := args[0]
			configuration := config.LoadConfig()
			createCtl := control.NewCreateController(control.TypePasswords, password, configuration)
			createCtl.ExecuteCreate()
		},
	}
	return passwordsCmd
}
