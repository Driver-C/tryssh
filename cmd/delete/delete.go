package delete

import (
	"github.com/spf13/cobra"
)

func NewDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete [command]",
		Short:   "Delete alternative username, port number, password, and login cache information",
		Long:    "Delete alternative username, port number, password, and login cache information",
		Aliases: []string{"del"},
	}
	cmd.AddCommand(NewUsersCommand())
	cmd.AddCommand(NewPortsCommand())
	cmd.AddCommand(NewPasswordsCommand())
	cmd.AddCommand(NewCachesCommand())
	cmd.AddCommand(NewKeysCommand())
	return cmd
}
