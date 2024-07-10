package create

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/spf13/cobra"
)

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
			ctl := control.NewCreateController(control.TypeKeys, keyPath, configuration)
			ctl.ExecuteCreate()
		},
	}
	return cmd
}
