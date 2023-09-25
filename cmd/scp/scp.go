package scp

import (
	"github.com/spf13/cobra"
	"time"
	"tryssh/config"
	"tryssh/control/scp"
)

const (
	concurrency = 8
	sshTimeout  = 1 * time.Second
)

var scpExample = `# download file
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
			concurrencyOpt, _ := cmd.Flags().GetInt("concurrency")
			timeout, _ := cmd.Flags().GetDuration("timeout")
			configuration := config.LoadConfig()
			scpControl := scp.NewScpController(source, destination, configuration)
			scpControl.TryCopy(user, concurrencyOpt, timeout)
		},
	}
	scpCmd.Flags().StringP(
		"user", "u", "", "Specify a username to attempt to login to the server,\n"+
			"if the specified username does not exist, try logging in using that username")
	scpCmd.Flags().IntP(
		"concurrency", "c", concurrency, "Number of multiple requests to perform at a time")
	scpCmd.Flags().DurationP("timeout", "t", sshTimeout,
		"SSH timeout when attempting to log in. It can be \"1s\" or \"1m\" or other duration")
	return scpCmd
}
