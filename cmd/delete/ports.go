package delete

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/spf13/cobra"
)

func NewPortsCommand() *cobra.Command {
	portsCmd := &cobra.Command{
		Use:     "ports <port>",
		Args:    cobra.ExactArgs(1),
		Short:   "Delete a alternate port",
		Long:    "Delete a alternate port",
		Aliases: []string{"port", "po"},
		Run: func(cmd *cobra.Command, args []string) {
			port := args[0]
			configuration := config.LoadConfig()
			deleteCtl := control.NewDeleteController(control.TypePorts, port, configuration)
			deleteCtl.ExecuteDelete()
		},
	}
	return portsCmd
}
