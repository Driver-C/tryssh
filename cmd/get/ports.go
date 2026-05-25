package get

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/Driver-C/tryssh/pkg/utils"
	"github.com/spf13/cobra"
)

// NewPortsCommand creates and returns the cobra command for retrieving port entries.
func NewPortsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ports <port>",
		Short:   "Get alternative ports",
		Long:    "Get alternative ports",
		Aliases: []string{"port", "po"},
		Run: func(_ *cobra.Command, args []string) {
			var port string
			if len(args) > 0 {
				port = args[0]
			}
			configuration, err := config.LoadConfig()
			if err != nil {
				utils.Fatalln(err)
			}
			controller := control.NewGetController(control.TypePorts, port, configuration)
			controller.ExecuteGet()
		},
	}
	return cmd
}
