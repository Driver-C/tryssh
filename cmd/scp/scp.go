package scp

import (
	"github.com/Driver-C/tryssh/config"
	"github.com/Driver-C/tryssh/control/scp"
	"github.com/spf13/cobra"
	"time"
)

const (
	concurrency = 8
	sshTimeout  = 1 * time.Second
)

var scpExample = `# Download test.txt file from 192.168.1.1 and place it under ./
tryssh scp 192.168.1.1:/root/test.txt ./
# Upload test.txt file to 192.168.1.1 and place it under /root/
tryssh scp ./test.txt 192.168.1.1:/root/
# Download test.txt file from 192.168.1.1 and rename it to test2.txt and place it under ./
tryssh scp 192.168.1.1:/root/test.txt ./test2.txt

# Download testDir directory from 192.168.1.1 and place it under ~/Downloads/
tryssh scp -r 192.168.1.1:/root/testDir ~/Downloads/
# Upload testDir directory to 192.168.1.1 and rename it to testDir2 and place it under /root/
tryssh scp -r ~/Downloads/testDir 192.168.1.1:/root/testDir2`

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
			recursive, _ := cmd.Flags().GetBool("recursive")
			configuration := config.LoadConfig()
			scpControl := scp.NewScpController(source, destination, configuration)
			scpControl.TryCopy(user, concurrencyOpt, recursive, timeout)
		},
	}
	scpCmd.Flags().StringP(
		"user", "u", "", "Specify a username to attempt to login to the server,\n"+
			"if the specified username does not exist, try logging in using that username")
	scpCmd.Flags().IntP(
		"concurrency", "c", concurrency, "Number of multiple requests to perform at a time")
	scpCmd.Flags().BoolP("recursive", "r", false, "Recursively copy entire directories")
	scpCmd.Flags().DurationP("timeout", "t", sshTimeout,
		"SSH timeout when attempting to log in. It can be \"1s\" or \"1m\" or other duration")
	return scpCmd
}
