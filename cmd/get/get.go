package get

import (
	"github.com/Driver-C/tryssh/cmd/get/caches"
	"github.com/Driver-C/tryssh/cmd/get/passwords"
	"github.com/Driver-C/tryssh/cmd/get/ports"
	"github.com/Driver-C/tryssh/cmd/get/users"
	"github.com/spf13/cobra"
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
	getCmd.AddCommand(NewKeysCommand())
	return getCmd
}
