package passwords

import (
	"github.com/Driver-C/tryssh/config"
	"github.com/Driver-C/tryssh/control/get"
	"github.com/spf13/cobra"
)

const getType = "passwords"

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
			getCtl := get.NewGetController(getType, password, configuration)
			getCtl.ExecuteGet()
		},
	}
	return passwordsCmd
}
