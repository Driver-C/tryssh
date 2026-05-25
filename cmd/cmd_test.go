package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTrysshCommand(t *testing.T) {
	rootCmd := NewTrysshCommand()

	assert.NotNil(t, rootCmd)
	assert.Equal(t, "tryssh [command]", rootCmd.Use)
	assert.Equal(t, "A command line ssh terminal tool.", rootCmd.Short)
	assert.Equal(t, "A command line ssh terminal tool.", rootCmd.Long)
}

func TestNewTrysshCommand_Subcommands(t *testing.T) {
	rootCmd := NewTrysshCommand()

	expectedSubcommands := []string{
		"version",
		"ssh",
		"scp",
		"alias",
		"create",
		"delete",
		"get",
		"prune",
		"encrypt",
	}

	for _, name := range expectedSubcommands {
		found := false
		for _, sub := range rootCmd.Commands() {
			if sub.Name() == name {
				found = true
				break
			}
		}
		assert.True(t, found, "expected subcommand %q to be registered", name)
	}
}

func TestNewTrysshCommand_DisableDefaultCompletion(t *testing.T) {
	rootCmd := NewTrysshCommand()

	assert.True(t, rootCmd.CompletionOptions.DisableDefaultCmd,
		"CompletionOptions.DisableDefaultCmd should be true")
}

func TestNewTrysshCommand_SubcommandCount(t *testing.T) {
	rootCmd := NewTrysshCommand()

	assert.Len(t, rootCmd.Commands(), 9,
		"root command should have exactly 9 subcommands")
}
