package create

import (
	"github.com/Driver-C/tryssh/cmd/create/caches"
	"github.com/Driver-C/tryssh/cmd/create/passwords"
	"github.com/Driver-C/tryssh/cmd/create/ports"
	"github.com/Driver-C/tryssh/cmd/create/users"
	"github.com/spf13/cobra"
)

func NewCreateCommand() *cobra.Command {
	createCmd := &cobra.Command{
		Use:     "create [command]",
		Short:   "Create alternate username, port number, password, and login cache information",
		Long:    "Create alternate username, port number, password, and login cache information",
		Aliases: []string{"cre", "crt", "add"},
	}
	createCmd.AddCommand(users.NewUsersCommand())
	createCmd.AddCommand(ports.NewPortsCommand())
	createCmd.AddCommand(passwords.NewPasswordsCommand())
	createCmd.AddCommand(caches.NewCachesCommand())
	return createCmd
}
