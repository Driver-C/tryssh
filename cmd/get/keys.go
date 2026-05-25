package get

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/Driver-C/tryssh/pkg/utils"
	"github.com/spf13/cobra"
)

// NewKeysCommand creates and returns the cobra command for retrieving key file path entries.
func NewKeysCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keys <keyFilePath>",
		Short:   "Get alternative key file path",
		Long:    "Get alternative key file path",
		Aliases: []string{"key"},
		Run: func(_ *cobra.Command, args []string) {
			var keyPath string
			if len(args) > 0 {
				keyPath = args[0]
			}
			configuration, err := config.LoadConfig()
			if err != nil {
				utils.Fatalln(err)
			}
			controller := control.NewGetController(control.TypeKeys, keyPath, configuration)
			controller.ExecuteGet()
		},
	}
	return cmd
}
