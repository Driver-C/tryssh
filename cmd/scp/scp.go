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
		Short:   "Upload/Download file to/from the server through SSH protocol",
		Long:    "Upload/Download file to/from the server through SSH protocol",
		Example: scpExample,
		Run: func(cmd *cobra.Command, args []string) {
			source := args[0]
			destination := args[1]
			user, _ := cmd.Flags().GetString("user")
			configuration := config.LoadConfig()
			scpControl := scp.NewScpController(source, destination, configuration)
			scpControl.TryCopy(user)
		},
	}
	scpCmd.Flags().StringP(
		"user", "u", "", "Specify a username to attempt to login to the server,\n"+
			"if the specified username does not exist, try logging in using that username")
	return scpCmd
}
