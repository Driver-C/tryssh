package get

import (
	"github.com/spf13/cobra"
)

func NewGetCommand() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get [command]",
		Short: "Get alternate username, port number, password, and login cache information",
		Long:  "Get alternate username, port number, password, and login cache information",
	}
	getCmd.AddCommand(NewUsersCommand())
	getCmd.AddCommand(NewPortsCommand())
	getCmd.AddCommand(NewPasswordsCommand())
	getCmd.AddCommand(NewCachesCommand())
	getCmd.AddCommand(NewKeysCommand())
	return getCmd
}
