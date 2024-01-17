package list

import (
	"github.com/Driver-C/tryssh/config"
	"github.com/Driver-C/tryssh/control/alias"
	"github.com/spf13/cobra"
)

func NewAliasListCommand() *cobra.Command {
	aliasListCmd := &cobra.Command{
		Use:     "list",
		Short:   "List all alias",
		Long:    "List all alias",
		Aliases: []string{"ls"},
		Run: func(cmd *cobra.Command, args []string) {
			configuration := config.LoadConfig()
			aliasController := alias.NewAliasController("", configuration, "")
			aliasController.ListAlias()
		},
	}
	return aliasListCmd
}
