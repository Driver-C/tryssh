package create

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/Driver-C/tryssh/pkg/utils"
	"github.com/spf13/cobra"
)

// NewPortsCommand creates and returns the cobra command for creating a port entry.
func NewPortsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ports <port>",
		Args:    cobra.ExactArgs(1),
		Short:   "Create an alternative port",
		Long:    "Create an alternative port",
		Aliases: []string{"port", "po"},
		Run: func(_ *cobra.Command, args []string) {
			port := args[0]
			configuration, err := config.LoadConfig()
			if err != nil {
				utils.Fatalln(err)
			}
			controller := control.NewCreateController(control.TypePorts, port, configuration)
			controller.ExecuteCreate()
		},
	}
	return cmd
}
