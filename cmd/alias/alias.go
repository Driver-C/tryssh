package alias

import (
	"github.com/spf13/cobra"
)

func NewAliasCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alias <subCommand> [flags]",
		Short: "Set, unset, and list aliases, aliases can be used to log in to servers",
		Long:  "Set, unset, and list aliases, aliases can be used to log in to servers",
	}
	cmd.AddCommand(NewAliasListCommand())
	cmd.AddCommand(NewAliasSetCommand())
	cmd.AddCommand(NewAliasUnsetCommand())
	return cmd
}
