package list

import (
	"github.com/spf13/cobra"
	"tryssh/config"
	"tryssh/control/alias"
)

func NewAliasListCommand() *cobra.Command {
	aliasListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all alias",
		Long:  "List all alias",
		Run: func(cmd *cobra.Command, args []string) {
			configuration := config.LoadConfig()
			aliasController := alias.NewAliasController("", configuration, "")
			aliasController.ListAlias()
		},
	}
	return aliasListCmd
}
