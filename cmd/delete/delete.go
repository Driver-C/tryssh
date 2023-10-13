package delete

import (
	"github.com/Driver-C/tryssh/cmd/delete/caches"
	"github.com/Driver-C/tryssh/cmd/delete/passwords"
	"github.com/Driver-C/tryssh/cmd/delete/ports"
	"github.com/Driver-C/tryssh/cmd/delete/users"
	"github.com/spf13/cobra"
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
