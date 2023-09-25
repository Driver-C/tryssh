package caches

import (
	"github.com/spf13/cobra"
	"tryssh/config"
	"tryssh/control/delete"
)

const deleteType = "caches"

func NewCachesCommand() *cobra.Command {
	cachesCmd := &cobra.Command{
		Use:   "caches <cacheIndex>",
		Args:  cobra.ExactArgs(1),
		Short: "Delete a alternate cache",
		Long:  "Delete a alternate cache",
		Run: func(cmd *cobra.Command, args []string) {
			cacheIndex := args[0]
			configuration := config.LoadConfig()
			deleteCtl := delete.NewDeleteController(deleteType, cacheIndex, configuration)
			deleteCtl.ExecuteDelete()
		},
	}
	return cachesCmd
}
