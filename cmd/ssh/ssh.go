// Package ssh provides the command for connecting to servers via SSH protocol.
package ssh

import (
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/Driver-C/tryssh/pkg/utils"
	"github.com/spf13/cobra"
	"time"
)

const (
	concurrency = 8
	sshTimeout  = 1 * time.Second
)

// NewSSHCommand creates and returns the cobra command for SSH connections.
func NewSSHCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssh <ipAddress>",
		Args:  cobra.ExactArgs(1),
		Short: "Connect to the server through SSH protocol",
		Long:  "Connect to the server through SSH protocol",
		Run: func(cmd *cobra.Command, args []string) {
			user, _ := cmd.Flags().GetString("user")
			concurrencyOpt, _ := cmd.Flags().GetInt("concurrency")
			timeout, _ := cmd.Flags().GetDuration("timeout")
			targetIP := args[0]
			configuration, err := config.LoadConfig()
			if err != nil {
				utils.Fatalln(err)
			}
			controller := control.NewSSHController(targetIP, configuration)
			controller.TryLogin(user, concurrencyOpt, timeout)
		},
	}
	cmd.Flags().StringP(
		"user", "u", "", "Specify a username to attempt to login to the server,\n"+
			"if the specified username does not exist, try logging in using that username")
	cmd.Flags().IntP(
		"concurrency", "c", concurrency, "Number of multiple requests to perform at a time")
	cmd.Flags().DurationP("timeout", "t", sshTimeout,
		"SSH timeout when attempting to log in. It can be \"1s\" or \"1m\" or other duration")
	return cmd
}
