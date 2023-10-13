package caches

import (
	"github.com/Driver-C/tryssh/config"
	"github.com/Driver-C/tryssh/control/delete"
	"github.com/spf13/cobra"
)

const deleteType = "caches"

func NewCachesCommand() *cobra.Command {
	cachesCmd := &cobra.Command{
		Use:   "caches <ipAddress>",
		Args:  cobra.ExactArgs(1),
		Short: "Delete a alternate cache",
		Long:  "Delete a alternate cache",
		Run: func(cmd *cobra.Command, args []string) {
			ipAddress := args[0]
			configuration := config.LoadConfig()
			deleteCtl := delete.NewDeleteController(deleteType, ipAddress, configuration)
			deleteCtl.ExecuteDelete()
		},
	}
	return cachesCmd
}
