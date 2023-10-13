package alias

import (
	"github.com/Driver-C/tryssh/cmd/alias/list"
	"github.com/Driver-C/tryssh/cmd/alias/set"
	"github.com/Driver-C/tryssh/cmd/alias/unset"
	"github.com/spf13/cobra"
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
