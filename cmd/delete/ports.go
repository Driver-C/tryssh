package delete

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/Driver-C/tryssh/pkg/utils"
	"github.com/spf13/cobra"
)

// NewPortsCommand creates and returns the cobra command for deleting a port entry.
func NewPortsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ports <port>",
		Args:    cobra.ExactArgs(1),
		Short:   "Delete an alternative port",
		Long:    "Delete an alternative port",
		Aliases: []string{"port", "po"},
		Run: func(_ *cobra.Command, args []string) {
			port := args[0]
			configuration, err := config.LoadConfig()
			if err != nil {
				utils.Fatalln(err)
			}
			controller := control.NewDeleteController(control.TypePorts, port, configuration)
			controller.ExecuteDelete()
		},
	}
	return cmd
}
