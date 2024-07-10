package create

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/spf13/cobra"
)

func NewPasswordsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "passwords <password>",
		Args:    cobra.ExactArgs(1),
		Short:   "Create an alternative password",
		Long:    "Create an alternative password",
		Aliases: []string{"password", "pass", "pwd"},
		Run: func(cmd *cobra.Command, args []string) {
			password := args[0]
			configuration := config.LoadConfig()
			controller := control.NewCreateController(control.TypePasswords, password, configuration)
			controller.ExecuteCreate()
		},
	}
	return cmd
}
