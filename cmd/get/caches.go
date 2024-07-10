package get

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/spf13/cobra"
)

func NewCachesCommand() *cobra.Command {
	cachesCmd := &cobra.Command{
		Use:     "caches <ipAddress>",
		Short:   "Get alternate caches by ipAddress",
		Long:    "Get alternate caches by ipAddress",
		Aliases: []string{"cache"},
		Run: func(cmd *cobra.Command, args []string) {
			var ipAddress string
			if len(args) > 0 {
				ipAddress = args[0]
			}
			configuration := config.LoadConfig()
			getCtl := control.NewGetController(control.TypeCaches, ipAddress, configuration)
			getCtl.ExecuteGet()
		},
	}
	return cachesCmd
}
