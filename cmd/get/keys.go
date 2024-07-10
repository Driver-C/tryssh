package get

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/spf13/cobra"
)

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
			ctl := control.NewGetController(control.TypeKeys, keyPath, configuration)
			ctl.ExecuteGet()
		},
	}
	return cmd
}
