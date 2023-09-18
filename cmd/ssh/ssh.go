package ssh

import (
	"github.com/spf13/cobra"
	"tryssh/config"
	"tryssh/control/ssh"
)

func NewSshCommand() *cobra.Command {
	sshCmd := &cobra.Command{
		Use:   "ssh <ipAddress>",
		Args:  cobra.ExactArgs(1),
		Short: "Connect to the server through SSH protocol",
		Long:  "Connect to the server through SSH protocol",
		Run: func(cmd *cobra.Command, args []string) {
			user, _ := cmd.Flags().GetString("user")
			targetIp := args[0]
			configuration := config.LoadConfig()
			sshControl := ssh.NewSshController(targetIp, configuration)
			sshControl.TryLogin(user)
		},
	}
	sshCmd.Flags().StringP(
		"user", "u", "", "Specify a username to attempt to login to the server,\n"+
			"if the specified username does not exist, try logging in using that username")
	return sshCmd
}
