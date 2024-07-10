package delete

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/spf13/cobra"
)

func NewCachesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "caches <ipAddress>",
		Args:    cobra.ExactArgs(1),
		Short:   "Delete an alternative cache",
		Long:    "Delete an alternative cache",
		Aliases: []string{"cache"},
		Run: func(cmd *cobra.Command, args []string) {
			ipAddress := args[0]
			configuration := config.LoadConfig()
			controller := control.NewDeleteController(control.TypeCaches, ipAddress, configuration)
			controller.ExecuteDelete()
		},
	}
	return cmd
}
