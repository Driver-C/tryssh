package create

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/spf13/cobra"
)

func NewPortsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ports <port>",
		Args:    cobra.ExactArgs(1),
		Short:   "Create an alternative port",
		Long:    "Create an alternative port",
		Aliases: []string{"port", "po"},
		Run: func(cmd *cobra.Command, args []string) {
			port := args[0]
			configuration := config.LoadConfig()
			controller := control.NewCreateController(control.TypePorts, port, configuration)
			controller.ExecuteCreate()
		},
	}
	return cmd
}
