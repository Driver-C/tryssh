package cmd

import (
	"github.com/spf13/cobra"
	"tryssh/cmd/alias"
	"tryssh/cmd/scp"
	"tryssh/cmd/ssh"
	"tryssh/cmd/version"
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
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	return rootCmd
}
