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
			var versionContent string
			if TrysshVersion != "" {
				versionContent += fmt.Sprintf("TrysshVersion: %s\n", TrysshVersion)
			}
			if BuildGoVersion != "" {
				versionContent += fmt.Sprintf("GoVersion: %s\n", BuildGoVersion)
			}
			if BuildTime != "" {
				versionContent += fmt.Sprintf("BuildTime: %s\n", BuildTime)
			}
			fmt.Printf(versionContent)
		},
	}
	return versionCmd
}
