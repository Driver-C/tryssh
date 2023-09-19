package alias

import (
	"github.com/spf13/cobra"
	"tryssh/cmd/alias/list"
	"tryssh/cmd/alias/set"
	"tryssh/cmd/alias/unset"
)

func NewAliasCommand() *cobra.Command {
	aliasCmd := &cobra.Command{
		Use:   "alias <subCommand> [flags]",
		Short: "Set, unset, and list aliases, aliases can be used to log in to servers",
		Long:  "Set, unset, and list aliases, aliases can be used to log in to servers",
	}
	aliasCmd.AddCommand(list.NewAliasListCommand())
	aliasCmd.AddCommand(set.NewAliasSetCommand())
	aliasCmd.AddCommand(unset.NewAliasUnsetCommand())
	return aliasCmd
}
