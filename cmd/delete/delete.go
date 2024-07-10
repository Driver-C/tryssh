package delete

import (
	"github.com/spf13/cobra"
)

func NewDeleteCommand() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:     "delete [command]",
		Short:   "Delete alternate username, port number, password, and login cache information",
		Long:    "Delete alternate username, port number, password, and login cache information",
		Aliases: []string{"del"},
	}
	deleteCmd.AddCommand(NewUsersCommand())
	deleteCmd.AddCommand(NewPortsCommand())
	deleteCmd.AddCommand(NewPasswordsCommand())
	deleteCmd.AddCommand(NewCachesCommand())
	deleteCmd.AddCommand(NewKeysCommand())
	return deleteCmd
}
