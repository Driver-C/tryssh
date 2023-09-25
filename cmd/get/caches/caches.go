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
			var cache string
			if len(args) > 0 {
				cache = args[0]
			}
			configuration := config.LoadConfig()
			getCtl := get.NewGetController(getType, cache, configuration)
			getCtl.ExecuteGet()
		},
	}
	return cachesCmd
}
