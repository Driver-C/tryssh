package ports

import (
	"github.com/spf13/cobra"
	"tryssh/config"
	"tryssh/control/delete"
)

const deleteType = "ports"

func NewPortsCommand() *cobra.Command {
	portsCmd := &cobra.Command{
		Use:   "ports <port>",
		Args:  cobra.ExactArgs(1),
		Short: "Delete a alternate port",
		Long:  "Delete a alternate port",
		Run: func(cmd *cobra.Command, args []string) {
			port := args[0]
			configuration := config.LoadConfig()
			deleteCtl := delete.NewDeleteController(deleteType, port, configuration)
			deleteCtl.ExecuteDelete()
		},
	}
	return portsCmd
}
