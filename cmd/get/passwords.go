package get

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/Driver-C/tryssh/pkg/utils"
	"github.com/spf13/cobra"
)

// NewPasswordsCommand creates and returns the cobra command for retrieving password entries.
func NewPasswordsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "passwords <password>",
		Short:   "Get alternative passwords",
		Long:    "Get alternative passwords",
		Aliases: []string{"password", "pass", "pwd"},
		Run: func(_ *cobra.Command, args []string) {
			var password string
			if len(args) > 0 {
				password = args[0]
			}
			configuration, err := config.LoadConfig()
			if err != nil {
				utils.Fatalln(err)
			}
			controller := control.NewGetController(control.TypePasswords, password, configuration)
			controller.ExecuteGet()
		},
	}
	return cmd
}
