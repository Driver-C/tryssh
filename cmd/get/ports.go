package get

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/spf13/cobra"
)

func NewPortsCommand() *cobra.Command {
	portsCmd := &cobra.Command{
		Use:     "ports <port>",
		Short:   "Get alternate ports",
		Long:    "Get alternate ports",
		Aliases: []string{"port", "po"},
		Run: func(cmd *cobra.Command, args []string) {
			var port string
			if len(args) > 0 {
				port = args[0]
			}
			configuration := config.LoadConfig()
			getCtl := control.NewGetController(control.TypePorts, port, configuration)
			getCtl.ExecuteGet()
		},
	}
	return portsCmd
}
