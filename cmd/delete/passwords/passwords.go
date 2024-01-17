package passwords

import (
	"github.com/Driver-C/tryssh/config"
	"github.com/Driver-C/tryssh/control/delete"
	"github.com/spf13/cobra"
)

const deleteType = "passwords"

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
			deleteCtl := delete.NewDeleteController(deleteType, password, configuration)
			deleteCtl.ExecuteDelete()
		},
	}
	return passwordsCmd
}
