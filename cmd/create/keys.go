package create

import (
	"github.com/Driver-C/tryssh/config"
	"github.com/Driver-C/tryssh/control/create"
	"github.com/spf13/cobra"
)

const createType = "keys"

func NewKeysCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keys <keyFilePath>",
		Args:    cobra.ExactArgs(1),
		Short:   "Create a alternate key file path",
		Long:    "Create a alternate key file path",
		Aliases: []string{"key"},
		Run: func(cmd *cobra.Command, args []string) {
			keyPath := args[0]
			configuration := config.LoadConfig()
			ctl := create.NewCreateController(createType, keyPath, configuration)
			ctl.ExecuteCreate()
		},
	}
	return cmd
}
