package version

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	TrysshVersion  string
	BuildGoVersion string
	BuildTime      string
)

func NewVersionCommand() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the client version information for the current context",
		Long:  "Print the client version information for the current context",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("TrysshVersion: %s, GoVersion: %s, BuildTime: %s\n",
				TrysshVersion, BuildGoVersion, BuildTime)
		},
	}
	return versionCmd
}
