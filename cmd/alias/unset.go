package alias

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/Driver-C/tryssh/pkg/utils"
	"github.com/spf13/cobra"
)

// NewAliasUnsetCommand creates and returns the cobra command for unsetting an alias.
func NewAliasUnsetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unset <alias>",
		Args:  cobra.ExactArgs(1),
		Short: "Unset the alias",
		Long:  "Unset the alias",
		Run: func(_ *cobra.Command, args []string) {
			aliasContent := args[0]
			configuration, err := config.LoadConfig()
			if err != nil {
				utils.Fatalln(err)
			}
			controller := control.NewAliasController("", configuration, aliasContent)
			controller.UnsetAlias()
		},
	}
	return cmd
}
