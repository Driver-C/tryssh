package create

import (
	"github.com/spf13/cobra"
)

func NewCreateCommand() *cobra.Command {
	createCmd := &cobra.Command{
		Use:     "create [command]",
		Short:   "Create alternate username, port number, password, and login cache information",
		Long:    "Create alternate username, port number, password, and login cache information",
		Aliases: []string{"cre", "crt", "add"},
	}
	createCmd.AddCommand(NewUsersCommand())
	createCmd.AddCommand(NewPortsCommand())
	createCmd.AddCommand(NewPasswordsCommand())
	createCmd.AddCommand(NewCachesCommand())
	createCmd.AddCommand(NewKeysCommand())
	return createCmd
}
