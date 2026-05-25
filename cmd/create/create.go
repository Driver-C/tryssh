// Package create provides commands for creating configuration entries (users, ports, passwords, keys, caches).
package create

import (
	"github.com/spf13/cobra"
)

// NewCreateCommand creates and returns the cobra command for creating configuration entries.
func NewCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [command]",
		Short:   "Create alternative username, port number, password, and login cache information",
		Long:    "Create alternative username, port number, password, and login cache information",
		Aliases: []string{"cre", "crt", "add"},
	}
	cmd.AddCommand(NewUsersCommand())
	cmd.AddCommand(NewPortsCommand())
	cmd.AddCommand(NewPasswordsCommand())
	cmd.AddCommand(NewCachesCommand())
	cmd.AddCommand(NewKeysCommand())
	return cmd
}
