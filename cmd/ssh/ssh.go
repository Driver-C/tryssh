package ssh

import (
	"github.com/spf13/cobra"
	"time"
	"tryssh/config"
	"tryssh/control/ssh"
)

const (
	concurrency = 8
	sshTimeout  = 1 * time.Second
)

func NewSshCommand() *cobra.Command {
	sshCmd := &cobra.Command{
		Use:   "ssh <ipAddress>",
		Args:  cobra.ExactArgs(1),
		Short: "Connect to the server through SSH protocol",
		Long:  "Connect to the server through SSH protocol",
		Run: func(cmd *cobra.Command, args []string) {
			user, _ := cmd.Flags().GetString("user")
			concurrencyOpt, _ := cmd.Flags().GetInt("concurrency")
			timeout, _ := cmd.Flags().GetDuration("timeout")
			targetIp := args[0]
			configuration := config.LoadConfig()
			sshControl := ssh.NewSshController(targetIp, configuration)
			sshControl.TryLogin(user, concurrencyOpt, timeout)
		},
	}
	sshCmd.Flags().StringP(
		"user", "u", "", "Specify a username to attempt to login to the server,\n"+
			"if the specified username does not exist, try logging in using that username")
	sshCmd.Flags().IntP(
		"concurrency", "c", concurrency, "Number of multiple requests to perform at a time")
	sshCmd.Flags().DurationP("timeout", "t", sshTimeout,
		"SSH timeout when attempting to log in. It can be \"1s\" or \"1m\" or other duration")
	return sshCmd
}
