package alias

import (
	"github.com/spf13/cobra"
)

func NewAliasCommand() *cobra.Command {
	aliasCmd := &cobra.Command{
		Use:   "alias <subCommand> [flags]",
		Short: "Set, unset, and list aliases, aliases can be used to log in to servers",
		Long:  "Set, unset, and list aliases, aliases can be used to log in to servers",
	}
	aliasCmd.AddCommand(NewAliasListCommand())
	aliasCmd.AddCommand(NewAliasSetCommand())
	aliasCmd.AddCommand(NewAliasUnsetCommand())
	return aliasCmd
}
