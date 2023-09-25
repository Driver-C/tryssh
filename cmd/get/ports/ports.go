package ports

import (
	"github.com/spf13/cobra"
	"tryssh/config"
	"tryssh/control/get"
)

const getType = "ports"

func NewPortsCommand() *cobra.Command {
	portsCmd := &cobra.Command{
		Use:   "ports <port>",
		Short: "Get alternate ports",
		Long:  "Get alternate ports",
		Run: func(cmd *cobra.Command, args []string) {
			var port string
			if len(args) > 0 {
				port = args[0]
			}
			configuration := config.LoadConfig()
			getCtl := get.NewGetController(getType, port, configuration)
			getCtl.ExecuteGet()
		},
	}
	return portsCmd
}
