package alias

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/Driver-C/tryssh/pkg/utils"
	"github.com/spf13/cobra"
)

// NewAliasListCommand creates and returns the cobra command for listing aliases.
func NewAliasListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all alias",
		Long:    "List all alias",
		Aliases: []string{"ls"},
		Run: func(_ *cobra.Command, _ []string) {
			configuration, err := config.LoadConfig()
			if err != nil {
				utils.Fatalln(err)
			}
			controller := control.NewAliasController("", configuration, "")
			controller.ListAlias()
		},
	}
	return cmd
}
