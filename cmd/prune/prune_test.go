package prune

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Driver-C/tryssh/pkg/config"
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

func TestNewPruneCommand_Structure(t *testing.T) {
	cmd := NewPruneCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "prune", cmd.Use)
	assert.Contains(t, cmd.Short, "Check if all current caches are available")
	assert.Equal(t, "Check if all current caches are available and clear the ones that are not available", cmd.Long)
	assert.NotNil(t, cmd.Run)
}

func TestNewPruneCommand_Flags(t *testing.T) {
	cmd := NewPruneCommand()

	// Test auto flag
	autoFlag := cmd.Flags().Lookup("auto")
	assert.NotNil(t, autoFlag, "auto flag should exist")
	assert.Equal(t, "a", autoFlag.Shorthand)
	assert.Equal(t, "false", autoFlag.DefValue)

	// Test concurrency flag
	concurrencyFlag := cmd.Flags().Lookup("concurrency")
	assert.NotNil(t, concurrencyFlag, "concurrency flag should exist")
	assert.Equal(t, "c", concurrencyFlag.Shorthand)
	assert.Equal(t, "8", concurrencyFlag.DefValue)

	// Test timeout flag
	timeoutFlag := cmd.Flags().Lookup("timeout")
	assert.NotNil(t, timeoutFlag, "timeout flag should exist")
	assert.Equal(t, "t", timeoutFlag.Shorthand)
	assert.Equal(t, "2s", timeoutFlag.DefValue)
}

func TestNewPruneCommand_FlagValues(t *testing.T) {
	cmd := NewPruneCommand()

	autoVal, err := cmd.Flags().GetBool("auto")
	assert.NoError(t, err)
	assert.False(t, autoVal)

	concurrencyVal, err := cmd.Flags().GetInt("concurrency")
	assert.NoError(t, err)
	assert.Equal(t, 8, concurrencyVal)

	timeoutVal, err := cmd.Flags().GetDuration("timeout")
	assert.NoError(t, err)
	assert.Equal(t, 2*time.Second, timeoutVal)
}

func TestNewPruneCommand_FlagParsing(t *testing.T) {
	cmd := NewPruneCommand()

	err := cmd.Flags().Set("auto", "true")
	assert.NoError(t, err)
	autoVal, _ := cmd.Flags().GetBool("auto")
	assert.True(t, autoVal)

	err = cmd.Flags().Set("concurrency", "4")
	assert.NoError(t, err)
	concurrencyVal, _ := cmd.Flags().GetInt("concurrency")
	assert.Equal(t, 4, concurrencyVal)

	err = cmd.Flags().Set("timeout", "10s")
	assert.NoError(t, err)
	timeoutVal, _ := cmd.Flags().GetDuration("timeout")
	assert.Equal(t, 10*time.Second, timeoutVal)
}

func TestNewPruneCommand_NoArgsRequired(t *testing.T) {
	cmd := NewPruneCommand()

	assert.Nil(t, cmd.Args, "prune command should accept zero args by default")
}

func TestNewPruneCommand_Run_DefaultFlags(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewPruneCommand()
	// This will call PruneCaches() with empty server list, which should be safe
	cmd.Run(cmd, []string{})
}

func TestNewPruneCommand_Run_WithAutoFlag(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewPruneCommand()
	_ = cmd.Flags().Set("auto", "true")
	_ = cmd.Flags().Set("concurrency", "4")
	_ = cmd.Flags().Set("timeout", "500ms")
	cmd.Run(cmd, []string{})
}

func TestNewPruneCommand_Run_WithCustomFlags(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewPruneCommand()
	_ = cmd.Flags().Set("auto", "false")
	_ = cmd.Flags().Set("concurrency", "1")
	_ = cmd.Flags().Set("timeout", "1s")
	cmd.Run(cmd, []string{})
}

func TestNewPruneCommand_Example(t *testing.T) {
	cmd := NewPruneCommand()
	// prune command doesn't have an example
	assert.Empty(t, cmd.Example)
}

func TestNewPruneCommand_NoAliases(t *testing.T) {
	cmd := NewPruneCommand()
	assert.Empty(t, cmd.Aliases, "prune command should have no aliases")
}
