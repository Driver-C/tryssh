package alias

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/spf13/cobra"
)

func NewAliasSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <alias> [flags]",
		Args:  cobra.ExactArgs(1),
		Short: "Set an alias for the specified server address",
		Long:  "Set an alias for the specified server address",
		Run: func(cmd *cobra.Command, args []string) {
			aliasContent := args[0]
			targetAddress, _ := cmd.Flags().GetString("target")
			configuration := config.LoadConfig()
			controller := control.NewAliasController(targetAddress, configuration, aliasContent)
			controller.SetAlias()
		},
	}
	cmd.Flags().StringP(
		"target", "t", "", "Set the alias for the target server address")
	err := cmd.MarkFlagRequired("target")
	if err != nil {
		return nil
	}
	return cmd
}
