package get

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/spf13/cobra"
)

func NewCachesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "caches <ipAddress>",
		Short:   "Get alternative caches by ipAddress",
		Long:    "Get alternative caches by ipAddress",
		Aliases: []string{"cache"},
		Run: func(cmd *cobra.Command, args []string) {
			var ipAddress string
			if len(args) > 0 {
				ipAddress = args[0]
			}
			configuration := config.LoadConfig()
			controller := control.NewGetController(control.TypeCaches, ipAddress, configuration)
			controller.ExecuteGet()
		},
	}
	return cmd
}
