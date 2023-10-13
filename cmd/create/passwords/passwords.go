package passwords

import (
	"github.com/Driver-C/tryssh/config"
	"github.com/Driver-C/tryssh/control/create"
	"github.com/spf13/cobra"
)

const createType = "passwords"

func NewPasswordsCommand() *cobra.Command {
	passwordsCmd := &cobra.Command{
		Use:   "passwords <password>",
		Args:  cobra.ExactArgs(1),
		Short: "Create a alternate password",
		Long:  "Create a alternate password",
		Run: func(cmd *cobra.Command, args []string) {
			password := args[0]
			configuration := config.LoadConfig()
			createCtl := create.NewCreateController(createType, password, configuration)
			createCtl.ExecuteCreate()
		},
	}
	return passwordsCmd
}
