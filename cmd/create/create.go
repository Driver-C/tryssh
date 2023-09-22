package create

import (
	"github.com/spf13/cobra"
	"tryssh/cmd/create/caches"
	"tryssh/cmd/create/passwords"
	"tryssh/cmd/create/ports"
	"tryssh/cmd/create/users"
)

func NewCreateCommand() *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "create [command]",
		Short: "Create alternate username, port number, password, and login cache information",
		Long:  "Create alternate username, port number, password, and login cache information",
	}
	createCmd.AddCommand(users.NewUsersCommand())
	createCmd.AddCommand(ports.NewPortsCommand())
	createCmd.AddCommand(passwords.NewPasswordsCommand())
	createCmd.AddCommand(caches.NewCachesCommand())
	return createCmd
}
