// Package delete provides commands for deleting configuration entries.
package delete

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/Driver-C/tryssh/pkg/utils"
	"github.com/spf13/cobra"
)

// NewCachesCommand creates and returns the cobra command for deleting a cache entry.
func NewCachesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "caches <ipAddress>",
		Args:    cobra.ExactArgs(1),
		Short:   "Delete an alternative cache",
		Long:    "Delete an alternative cache",
		Aliases: []string{"cache"},
		Run: func(_ *cobra.Command, args []string) {
			ipAddress := args[0]
			configuration, err := config.LoadConfig()
			if err != nil {
				utils.Fatalln(err)
			}
			controller := control.NewDeleteController(control.TypeCaches, ipAddress, configuration)
			controller.ExecuteDelete()
		},
	}
	return cmd
}
