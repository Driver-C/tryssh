package create

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/Driver-C/tryssh/pkg/utils"
	"github.com/spf13/cobra"
)

// NewKeysCommand creates and returns the cobra command for creating a key file path entry.
func NewKeysCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keys <keyFilePath>",
		Args:    cobra.ExactArgs(1),
		Short:   "Create an alternative key file path",
		Long:    "Create an alternative key file path",
		Aliases: []string{"key"},
		Run: func(_ *cobra.Command, args []string) {
			keyPath := args[0]
			configuration, err := config.LoadConfig()
			if err != nil {
				utils.Fatalln(err)
			}
			controller := control.NewCreateController(control.TypeKeys, keyPath, configuration)
			controller.ExecuteCreate()
		},
	}
	return cmd
}
