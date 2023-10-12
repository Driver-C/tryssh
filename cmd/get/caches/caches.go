package caches

import (
	"github.com/spf13/cobra"
	"tryssh/config"
	"tryssh/control/get"
)

const getType = "caches"

func NewCachesCommand() *cobra.Command {
	cachesCmd := &cobra.Command{
		Use:   "caches <ipAddress>",
		Short: "Get alternate caches by ipAddress",
		Long:  "Get alternate caches by ipAddress",
		Run: func(cmd *cobra.Command, args []string) {
			var ipAddress string
			if len(args) > 0 {
				ipAddress = args[0]
			}
			configuration := config.LoadConfig()
			getCtl := get.NewGetController(getType, ipAddress, configuration)
			getCtl.ExecuteGet()
		},
	}
	return cachesCmd
}
