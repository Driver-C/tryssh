package delete

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/spf13/cobra"
)

func NewPasswordsCommand() *cobra.Command {
	passwordsCmd := &cobra.Command{
		Use:     "passwords <password>",
		Args:    cobra.ExactArgs(1),
		Short:   "Delete a alternate password",
		Long:    "Delete a alternate password",
		Aliases: []string{"password", "pass", "pwd"},
		Run: func(cmd *cobra.Command, args []string) {
			password := args[0]
			configuration := config.LoadConfig()
			deleteCtl := control.NewDeleteController(control.TypePasswords, password, configuration)
			deleteCtl.ExecuteDelete()
		},
	}
	return passwordsCmd
}
