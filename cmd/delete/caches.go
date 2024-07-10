package delete

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/spf13/cobra"
)

func NewCachesCommand() *cobra.Command {
	cachesCmd := &cobra.Command{
		Use:     "caches <ipAddress>",
		Args:    cobra.ExactArgs(1),
		Short:   "Delete a alternate cache",
		Long:    "Delete a alternate cache",
		Aliases: []string{"cache"},
		Run: func(cmd *cobra.Command, args []string) {
			ipAddress := args[0]
			configuration := config.LoadConfig()
			deleteCtl := control.NewDeleteController(control.TypeCaches, ipAddress, configuration)
			deleteCtl.ExecuteDelete()
		},
	}
	return cachesCmd
}
