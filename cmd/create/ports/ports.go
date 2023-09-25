package ports

import (
	"github.com/spf13/cobra"
	"tryssh/config"
	"tryssh/control/create"
)

const createType = "ports"

func NewPortsCommand() *cobra.Command {
	portsCmd := &cobra.Command{
		Use:   "ports <port>",
		Args:  cobra.ExactArgs(1),
		Short: "Create a alternate port",
		Long:  "Create a alternate port",
		Run: func(cmd *cobra.Command, args []string) {
			port := args[0]
			configuration := config.LoadConfig()
			createCtl := create.NewCreateController(createType, port, configuration)
			createCtl.ExecuteCreate()
		},
	}
	return portsCmd
}
