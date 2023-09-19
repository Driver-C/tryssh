package unset

import (
	"github.com/spf13/cobra"
	"tryssh/config"
	"tryssh/control/alias"
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
