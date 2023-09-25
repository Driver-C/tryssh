package get

import (
	"github.com/spf13/cobra"
	"tryssh/cmd/get/caches"
	"tryssh/cmd/get/passwords"
	"tryssh/cmd/get/ports"
	"tryssh/cmd/get/users"
)

func NewGetCommand() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get [command]",
		Short: "Get alternate username, port number, password, and login cache information",
		Long:  "Get alternate username, port number, password, and login cache information",
	}
	getCmd.AddCommand(users.NewUsersCommand())
	getCmd.AddCommand(ports.NewPortsCommand())
	getCmd.AddCommand(passwords.NewPasswordsCommand())
	getCmd.AddCommand(caches.NewCachesCommand())
	return getCmd
}
