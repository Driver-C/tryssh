package version

import (
	"fmt"
	"github.com/spf13/cobra"
)

const VERSION = "0.1.0"

func NewVersionCommand() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the client version information for the current context",
		Long:  "Print the client version information for the current context",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(VERSION)
		},
	}
	return versionCmd
}
