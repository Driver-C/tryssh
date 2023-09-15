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
			targetIp := args[0]
			configuration := config.LoadConfig()
			sshControl := ssh.NewSshController(targetIp, configuration)
			sshControl.TryLogin()
		},
	}
	return sshCmd
}
