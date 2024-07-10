package delete

import (
	"github.com/Driver-C/tryssh/config"
	"github.com/Driver-C/tryssh/control/delete"
	"github.com/spf13/cobra"
)

const deleteType = "keys"

func NewKeysCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keys <keyFilePath>",
		Args:    cobra.ExactArgs(1),
		Short:   "Delete a alternate key file path",
		Long:    "Delete a alternate key file path",
		Aliases: []string{"key"},
		Run: func(cmd *cobra.Command, args []string) {
			keyPath := args[0]
			configuration := config.LoadConfig()
			ctl := delete.NewDeleteController(deleteType, keyPath, configuration)
			ctl.ExecuteDelete()
		},
	}
	return cmd
}
