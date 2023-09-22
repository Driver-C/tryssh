package passwords

import (
	"github.com/spf13/cobra"
	"tryssh/config"
	createControl "tryssh/control/create"
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
			createCtl := createControl.NewCreateController(createType, password, configuration)
			createCtl.ExecuteCreate()
		},
	}
	return passwordsCmd
}
