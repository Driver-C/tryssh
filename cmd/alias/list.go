package alias

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/spf13/cobra"
)

func NewAliasListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all alias",
		Long:    "List all alias",
		Aliases: []string{"ls"},
		Run: func(cmd *cobra.Command, args []string) {
			configuration := config.LoadConfig()
			controller := control.NewAliasController("", configuration, "")
			controller.ListAlias()
		},
	}
	return cmd
}
