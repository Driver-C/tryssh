package set

import (
	"github.com/Driver-C/tryssh/config"
	"github.com/Driver-C/tryssh/control/alias"
	"github.com/spf13/cobra"
)

func NewAliasSetCommand() *cobra.Command {
	aliasSetCmd := &cobra.Command{
		Use:   "set <alias> [flags]",
		Args:  cobra.ExactArgs(1),
		Short: "Set an alias for the specified server address",
		Long:  "Set an alias for the specified server address",
		Run: func(cmd *cobra.Command, args []string) {
			aliasContent := args[0]
			targetAddress, _ := cmd.Flags().GetString("target")
			configuration := config.LoadConfig()
			aliasController := alias.NewAliasController(targetAddress, configuration, aliasContent)
			aliasController.SetAlias()
		},
	}
	aliasSetCmd.Flags().StringP(
		"target", "t", "", "Set the alias for the target server address")
	err := aliasSetCmd.MarkFlagRequired("target")
	if err != nil {
		return nil
	}
	return aliasSetCmd
}
