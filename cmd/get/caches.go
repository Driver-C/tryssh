// Package get provides commands for querying configuration entries.
package get

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/Driver-C/tryssh/pkg/utils"
	"github.com/spf13/cobra"
)

// NewCachesCommand creates and returns the cobra command for retrieving cache entries.
func NewCachesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "caches [ipAddress]",
		Short:   "Get alternative caches by ipAddress",
		Long:    "Get alternative caches by ipAddress",
		Aliases: []string{"cache"},
		Run: func(_ *cobra.Command, args []string) {
			var ipAddress string
			if len(args) > 0 {
				ipAddress = args[0]
			}
			configuration, err := config.LoadConfig()
			if err != nil {
				utils.Fatalln(err)
			}
			controller := control.NewGetController(control.TypeCaches, ipAddress, configuration)
			controller.ExecuteGet()
		},
	}
	return cmd
}
