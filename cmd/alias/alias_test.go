package alias

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// setupTempConfig creates a temporary config file and overrides the default paths.
func setupTempConfig(t *testing.T) func() {
	t.Helper()
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "tryssh.db")
	knownHostsPath := filepath.Join(tmpDir, "known_hosts")

	err := os.MkdirAll(tmpDir, 0755)
	assert.NoError(t, err)
	configData := []byte("main:\n  ports: []\n  users: []\n  passwords: []\n  keys: []\nserverList: []\n")
	err = os.WriteFile(configPath, configData, 0600)
	assert.NoError(t, err)

	origConfigPath := config.DefaultConfigPath
	origKnownHostsPath := config.DefaultKnownHostsPath
	config.DefaultConfigPath = configPath
	config.DefaultKnownHostsPath = knownHostsPath

	return func() {
		config.DefaultConfigPath = origConfigPath
		config.DefaultKnownHostsPath = origKnownHostsPath
	}
}

// setupTempConfigWithServers creates a temp config with server entries for alias testing.
func setupTempConfigWithServers(t *testing.T) func() {
	t.Helper()
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "tryssh.db")
	knownHostsPath := filepath.Join(tmpDir, "known_hosts")

	err := os.MkdirAll(tmpDir, 0755)
	assert.NoError(t, err)
	configData := []byte("main:\n  ports: [\"22\"]\n  users: [\"root\"]\n  passwords: [\"testpass\"]\n  keys: []\nserverList:\n  - ip: \"192.168.1.1\"\n    port: \"22\"\n    user: \"root\"\n    password: \"testpass\"\n    key: \"\"\n    alias: \"myserver\"\n  - ip: \"10.0.0.1\"\n    port: \"22\"\n    user: \"admin\"\n    password: \"adminpass\"\n    key: \"\"\n    alias: \"\"\n")
	err = os.WriteFile(configPath, configData, 0600)
	assert.NoError(t, err)

	origConfigPath := config.DefaultConfigPath
	origKnownHostsPath := config.DefaultKnownHostsPath
	config.DefaultConfigPath = configPath
	config.DefaultKnownHostsPath = knownHostsPath

	return func() {
		config.DefaultConfigPath = origConfigPath
		config.DefaultKnownHostsPath = origKnownHostsPath
	}
}

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestNewAliasCommand_Structure(t *testing.T) {
	cmd := NewAliasCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "alias <subCommand> [flags]", cmd.Use)
	assert.Contains(t, cmd.Short, "Set, unset, and list aliases")
}

func TestNewAliasCommand_Long(t *testing.T) {
	cmd := NewAliasCommand()
	assert.Equal(t, "Set, unset, and list aliases, aliases can be used to log in to servers", cmd.Long)
}

func TestNewAliasCommand_Subcommands(t *testing.T) {
	cmd := NewAliasCommand()

	expectedSubcommands := []string{"list", "set", "unset"}
	for _, name := range expectedSubcommands {
		found := false
		for _, sub := range cmd.Commands() {
			if sub.Name() == name {
				found = true
				break
			}
		}
		assert.True(t, found, "expected subcommand %q to be registered", name)
	}
}

func TestNewAliasCommand_SubcommandCount(t *testing.T) {
	cmd := NewAliasCommand()
	assert.Len(t, cmd.Commands(), 3, "alias command should have exactly 3 subcommands")
}

// --- List command ---

func TestNewAliasListCommand_Structure(t *testing.T) {
	cmd := NewAliasListCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "list", cmd.Use)
	assert.Equal(t, "List all alias", cmd.Short)
	assert.Equal(t, "List all alias", cmd.Long)
	assert.Contains(t, cmd.Aliases, "ls")
	assert.NotNil(t, cmd.Run)
}

func TestNewAliasListCommand_Run(t *testing.T) {
	cleanup := setupTempConfigWithServers(t)
	defer cleanup()

	cmd := NewAliasListCommand()
	output := captureOutput(t, func() {
		cmd.Run(cmd, []string{})
	})
	assert.Contains(t, output, "Alias: myserver")
	assert.Contains(t, output, "192.168.1.1")
}

func TestNewAliasListCommand_Run_EmptyConfig(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewAliasListCommand()
	// Should not panic when no aliases exist
	output := captureOutput(t, func() {
		cmd.Run(cmd, []string{})
	})
	// No aliases, so no output to stdout (logs go to logger)
	assert.NotContains(t, output, "Alias:")
}

// --- Set command ---

func TestNewAliasSetCommand_Structure(t *testing.T) {
	cmd := NewAliasSetCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "set <alias> [flags]", cmd.Use)
	assert.Equal(t, "Set an alias for the specified server address", cmd.Short)
	assert.Equal(t, "Set an alias for the specified server address", cmd.Long)
	assert.NotNil(t, cmd.Run)
}

func TestNewAliasSetCommand_TargetFlag(t *testing.T) {
	cmd := NewAliasSetCommand()

	targetFlag := cmd.Flags().Lookup("target")
	assert.NotNil(t, targetFlag, "target flag should exist")
	assert.Equal(t, "t", targetFlag.Shorthand)
	assert.Equal(t, "", targetFlag.DefValue)
}

func TestNewAliasSetCommand_TargetFlagRequired(t *testing.T) {
	cmd := NewAliasSetCommand()

	annotations := cmd.Flags().Lookup("target").Annotations
	assert.NotNil(t, annotations, "target flag should have annotations")
	_, hasRequired := annotations[cobra.BashCompOneRequiredFlag]
	assert.True(t, hasRequired, "target flag should be marked as required")
}

func TestNewAliasSetCommand_TargetFlagParsing(t *testing.T) {
	cmd := NewAliasSetCommand()

	err := cmd.Flags().Set("target", "192.168.1.1")
	assert.NoError(t, err)
	val, _ := cmd.Flags().GetString("target")
	assert.Equal(t, "192.168.1.1", val)
}

func TestNewAliasSetCommand_ArgsValidation(t *testing.T) {
	cmd := NewAliasSetCommand()

	err := cmd.Args(cmd, []string{})
	assert.Error(t, err)

	err = cmd.Args(cmd, []string{"myalias"})
	assert.NoError(t, err)

	err = cmd.Args(cmd, []string{"a", "b"})
	assert.Error(t, err)
}

func TestNewAliasSetCommand_Run(t *testing.T) {
	cleanup := setupTempConfigWithServers(t)
	defer cleanup()

	cmd := NewAliasSetCommand()
	_ = cmd.Flags().Set("target", "10.0.0.1")
	cmd.Run(cmd, []string{"newalias"})
}

func TestNewAliasSetCommand_Run_NoMatchingIP(t *testing.T) {
	cleanup := setupTempConfigWithServers(t)
	defer cleanup()

	cmd := NewAliasSetCommand()
	_ = cmd.Flags().Set("target", "99.99.99.99")
	cmd.Run(cmd, []string{"nope"})
}

func TestNewAliasSetCommand_Run_DuplicateAlias(t *testing.T) {
	cleanup := setupTempConfigWithServers(t)
	defer cleanup()

	cmd := NewAliasSetCommand()
	_ = cmd.Flags().Set("target", "10.0.0.1")
	// "myserver" is already used as an alias in the config
	cmd.Run(cmd, []string{"myserver"})
}

// --- Unset command ---

func TestNewAliasUnsetCommand_Structure(t *testing.T) {
	cmd := NewAliasUnsetCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "unset <alias>", cmd.Use)
	assert.Equal(t, "Unset the alias", cmd.Short)
	assert.Equal(t, "Unset the alias", cmd.Long)
	assert.NotNil(t, cmd.Run)
}

func TestNewAliasUnsetCommand_ArgsValidation(t *testing.T) {
	cmd := NewAliasUnsetCommand()

	err := cmd.Args(cmd, []string{})
	assert.Error(t, err)

	err = cmd.Args(cmd, []string{"myalias"})
	assert.NoError(t, err)

	err = cmd.Args(cmd, []string{"a", "b"})
	assert.Error(t, err)
}

func TestNewAliasUnsetCommand_Run(t *testing.T) {
	cleanup := setupTempConfigWithServers(t)
	defer cleanup()

	cmd := NewAliasUnsetCommand()
	cmd.Run(cmd, []string{"myserver"})
}

func TestNewAliasUnsetCommand_Run_NoMatchingAlias(t *testing.T) {
	cleanup := setupTempConfigWithServers(t)
	defer cleanup()

	cmd := NewAliasUnsetCommand()
	cmd.Run(cmd, []string{"nonexistent"})
}

// --- Full command hierarchy tests ---

func TestNewAliasCommand_NoAliases(t *testing.T) {
	cmd := NewAliasCommand()
	assert.Empty(t, cmd.Aliases, "alias command should have no aliases")
}

func TestNewAliasSetCommand_DefaultTargetValue(t *testing.T) {
	cmd := NewAliasSetCommand()
	val, err := cmd.Flags().GetString("target")
	assert.NoError(t, err)
	assert.Equal(t, "", val)
}

func TestNewAliasListCommand_NoArgs(t *testing.T) {
	cmd := NewAliasListCommand()
	// list command has no Args constraint set, so it accepts any args
	// Verify it doesn't have an explicit Args validator
	assert.Nil(t, cmd.Args)
}
