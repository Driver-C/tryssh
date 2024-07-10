package get

import (
	"github.com/Driver-C/tryssh/config"
	"github.com/Driver-C/tryssh/control/get"
	"github.com/spf13/cobra"
)

const getType = "keys"

func NewKeysCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keys <keyFilePath>",
		Short:   "Get alternate key file path",
		Long:    "Get alternate key file path",
		Aliases: []string{"key"},
		Run: func(cmd *cobra.Command, args []string) {
			var keyPath string
			if len(args) > 0 {
				keyPath = args[0]
			}
			configuration := config.LoadConfig()
			ctl := get.NewGetController(getType, keyPath, configuration)
			ctl.ExecuteGet()
		},
	}
	return cmd
}
