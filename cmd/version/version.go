package version

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	Version        string
	BuildGoVersion string
	BuildTime      string
)

func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the client version information for the current context",
		Long:  "Print the client version information for the current context",
		Run: func(cmd *cobra.Command, args []string) {
			var versionContent string
			if Version != "" {
				versionContent += fmt.Sprintf("Version: %s\n", Version)
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
	return cmd
}
