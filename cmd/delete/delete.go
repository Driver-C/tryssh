package delete

import (
	"github.com/spf13/cobra"
	"tryssh/cmd/delete/caches"
	"tryssh/cmd/delete/passwords"
	"tryssh/cmd/delete/ports"
	"tryssh/cmd/delete/users"
)

func NewDeleteCommand() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "delete [command]",
		Short: "Delete alternate username, port number, password, and login cache information",
		Long:  "Delete alternate username, port number, password, and login cache information",
	}
	deleteCmd.AddCommand(users.NewUsersCommand())
	deleteCmd.AddCommand(ports.NewPortsCommand())
	deleteCmd.AddCommand(passwords.NewPasswordsCommand())
	deleteCmd.AddCommand(caches.NewCachesCommand())
	return deleteCmd
}
