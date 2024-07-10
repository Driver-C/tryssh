package create

import (
	"github.com/spf13/cobra"
)

func NewCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [command]",
		Short:   "Create alternative username, port number, password, and login cache information",
		Long:    "Create alternative username, port number, password, and login cache information",
		Aliases: []string{"cre", "crt", "add"},
	}
	cmd.AddCommand(NewUsersCommand())
	cmd.AddCommand(NewPortsCommand())
	cmd.AddCommand(NewPasswordsCommand())
	cmd.AddCommand(NewCachesCommand())
	cmd.AddCommand(NewKeysCommand())
	return cmd
}
