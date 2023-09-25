package passwords

import (
	"github.com/spf13/cobra"
	"tryssh/config"
	"tryssh/control/get"
)

const getType = "passwords"

func NewPasswordsCommand() *cobra.Command {
	passwordsCmd := &cobra.Command{
		Use:   "passwords <password>",
		Short: "Get alternate passwords",
		Long:  "Get alternate passwords",
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
