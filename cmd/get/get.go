package get

import (
	"github.com/spf13/cobra"
)

func NewGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [command]",
		Short: "Get alternative username, port number, password, and login cache information",
		Long:  "Get alternative username, port number, password, and login cache information",
	}
	cmd.AddCommand(NewUsersCommand())
	cmd.AddCommand(NewPortsCommand())
	cmd.AddCommand(NewPasswordsCommand())
	cmd.AddCommand(NewCachesCommand())
	cmd.AddCommand(NewKeysCommand())
	return cmd
}
