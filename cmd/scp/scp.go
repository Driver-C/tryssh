package scp

import (
	"github.com/spf13/cobra"
	"tryssh/config"
	"tryssh/control/scp"
)

var scpExample = `
# download file
	tryssh scp 192.168.1.1:/root/test.txt ./
# upload file
	tryssh scp ./test.txt 192.168.1.1:/root/`

func NewScpCommand() *cobra.Command {
	scpCmd := &cobra.Command{
		Use:     "scp <source> <destination>",
		Args:    cobra.ExactArgs(2),
		Short:   "Copy file to the server through SSH protocol",
		Long:    "Copy file to the server through SSH protocol",
		Example: scpExample,
		Run: func(cmd *cobra.Command, args []string) {
			source := args[0]
			destination := args[1]
			configuration := config.LoadConfig()
			scpControl := scp.NewScpController(source, destination, configuration)
			scpControl.TryCopy()
		},
	}
	return scpCmd
}
