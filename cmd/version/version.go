// Package version provides version information for the tryssh CLI.
package version

import (
	"fmt"
	"github.com/spf13/cobra"
)

// Version holds the build version, Go version, and build time set at link time.
var (
	// Version is the application version string set at build time.
	Version        string
	BuildGoVersion string
	BuildTime      string
)

// NewVersionCommand creates and returns the cobra command for displaying version information.
func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the client version information for the current context",
		Long:  "Print the client version information for the current context",
		Run: func(_ *cobra.Command, _ []string) {
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
			if versionContent == "" {
					versionContent = "Version: (dev)\n"
				}
				fmt.Printf("%s", versionContent)
		},
	}
	return cmd
}
