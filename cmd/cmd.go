package cmd

import (
	"github.com/Driver-C/tryssh/cmd/alias"
	"github.com/Driver-C/tryssh/cmd/create"
	"github.com/Driver-C/tryssh/cmd/delete"
	"github.com/Driver-C/tryssh/cmd/get"
	"github.com/Driver-C/tryssh/cmd/scp"
	"github.com/Driver-C/tryssh/cmd/ssh"
	"github.com/Driver-C/tryssh/cmd/version"
	"github.com/spf13/cobra"
)

func NewTrysshCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tryssh [command]",
		Short: "A command line ssh terminal tool.",
		Long:  "A command line ssh terminal tool.",
	}
	rootCmd.AddCommand(version.NewVersionCommand())
	rootCmd.AddCommand(ssh.NewSshCommand())
	rootCmd.AddCommand(scp.NewScpCommand())
	rootCmd.AddCommand(alias.NewAliasCommand())
	rootCmd.AddCommand(create.NewCreateCommand())
	rootCmd.AddCommand(delete.NewDeleteCommand())
	rootCmd.AddCommand(get.NewGetCommand())
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	return rootCmd
}
