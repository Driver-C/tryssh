package unset

import (
	"github.com/Driver-C/tryssh/config"
	"github.com/Driver-C/tryssh/control/alias"
	"github.com/spf13/cobra"
)

func NewAliasUnsetCommand() *cobra.Command {
	aliasUnsetCmd := &cobra.Command{
		Use:   "unset <alias>",
		Args:  cobra.ExactArgs(1),
		Short: "Unset the alias",
		Long:  "Unset the alias",
		Run: func(cmd *cobra.Command, args []string) {
			aliasContent := args[0]
			configuration := config.LoadConfig()
			aliasController := alias.NewAliasController("", configuration, aliasContent)
			aliasController.UnsetAlias()
		},
	}
	return aliasUnsetCmd
}
