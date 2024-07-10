package create

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/spf13/cobra"
)

func NewPortsCommand() *cobra.Command {
	portsCmd := &cobra.Command{
		Use:     "ports <port>",
		Args:    cobra.ExactArgs(1),
		Short:   "Create a alternate port",
		Long:    "Create a alternate port",
		Aliases: []string{"port", "po"},
		Run: func(cmd *cobra.Command, args []string) {
			port := args[0]
			configuration := config.LoadConfig()
			createCtl := control.NewCreateController(control.TypePorts, port, configuration)
			createCtl.ExecuteCreate()
		},
	}
	return portsCmd
}
