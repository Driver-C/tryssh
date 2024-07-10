package alias

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/spf13/cobra"
)

func NewAliasUnsetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unset <alias>",
		Args:  cobra.ExactArgs(1),
		Short: "Unset the alias",
		Long:  "Unset the alias",
		Run: func(cmd *cobra.Command, args []string) {
			aliasContent := args[0]
			configuration := config.LoadConfig()
			controller := control.NewAliasController("", configuration, aliasContent)
			controller.UnsetAlias()
		},
	}
	return cmd
}
