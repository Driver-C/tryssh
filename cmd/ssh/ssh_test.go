package ssh

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

func TestNewSSHCommand_Structure(t *testing.T) {
	cmd := NewSSHCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "ssh <ipAddress>", cmd.Use)
	assert.Equal(t, "Connect to the server through SSH protocol", cmd.Short)
	assert.Equal(t, "Connect to the server through SSH protocol", cmd.Long)
	assert.NotNil(t, cmd.Run)
}

func TestNewSSHCommand_Flags(t *testing.T) {
	cmd := NewSSHCommand()

	// Test user flag
	userFlag := cmd.Flags().Lookup("user")
	assert.NotNil(t, userFlag, "user flag should exist")
	assert.Equal(t, "u", userFlag.Shorthand)
	assert.Equal(t, "", userFlag.DefValue)

	// Test concurrency flag
	concurrencyFlag := cmd.Flags().Lookup("concurrency")
	assert.NotNil(t, concurrencyFlag, "concurrency flag should exist")
	assert.Equal(t, "c", concurrencyFlag.Shorthand)
	assert.Equal(t, "8", concurrencyFlag.DefValue)

	// Test timeout flag
	timeoutFlag := cmd.Flags().Lookup("timeout")
	assert.NotNil(t, timeoutFlag, "timeout flag should exist")
	assert.Equal(t, "t", timeoutFlag.Shorthand)
	assert.Equal(t, "1s", timeoutFlag.DefValue)
}

func TestNewSSHCommand_ArgsValidation(t *testing.T) {
	cmd := NewSSHCommand()

	// No args should fail
	err := cmd.Args(cmd, []string{})
	assert.Error(t, err)

	// One arg should pass
	err = cmd.Args(cmd, []string{"192.168.1.1"})
	assert.NoError(t, err)

	// Two args should fail
	err = cmd.Args(cmd, []string{"192.168.1.1", "extra"})
	assert.Error(t, err)
}

func TestNewSSHCommand_FlagValues(t *testing.T) {
	cmd := NewSSHCommand()

	// Test default concurrency value
	concurrencyVal, err := cmd.Flags().GetInt("concurrency")
	assert.NoError(t, err)
	assert.Equal(t, 8, concurrencyVal)

	// Test default timeout value
	timeoutVal, err := cmd.Flags().GetDuration("timeout")
	assert.NoError(t, err)
	assert.Equal(t, 1*time.Second, timeoutVal)

	// Test default user value
	userVal, err := cmd.Flags().GetString("user")
	assert.NoError(t, err)
	assert.Equal(t, "", userVal)
}

func TestNewSSHCommand_FlagParsing(t *testing.T) {
	cmd := NewSSHCommand()

	// Set flags and verify
	err := cmd.Flags().Set("user", "root")
	assert.NoError(t, err)
	userVal, _ := cmd.Flags().GetString("user")
	assert.Equal(t, "root", userVal)

	err = cmd.Flags().Set("concurrency", "16")
	assert.NoError(t, err)
	concurrencyVal, _ := cmd.Flags().GetInt("concurrency")
	assert.Equal(t, 16, concurrencyVal)

	err = cmd.Flags().Set("timeout", "5s")
	assert.NoError(t, err)
	timeoutVal, _ := cmd.Flags().GetDuration("timeout")
	assert.Equal(t, 5*time.Second, timeoutVal)
}

func TestNewSSHCommand_Run_DefaultFlags(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewSSHCommand()
	// TryLogin will attempt connections with empty credentials and fail,
	// but the Run closure body will be exercised
	cmd.Run(cmd, []string{"192.168.255.255"})
}

func TestNewSSHCommand_Run_WithUser(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewSSHCommand()
	_ = cmd.Flags().Set("user", "testuser")
	_ = cmd.Flags().Set("timeout", "100ms")
	_ = cmd.Flags().Set("concurrency", "1")
	cmd.Run(cmd, []string{"192.168.255.255"})
}

func TestNewSSHCommand_NoAliases(t *testing.T) {
	cmd := NewSSHCommand()
	assert.Empty(t, cmd.Aliases, "ssh command should have no aliases")
}

func TestNewSSHCommand_NoExample(t *testing.T) {
	cmd := NewSSHCommand()
	assert.Empty(t, cmd.Example, "ssh command should have no example")
}

func TestNewSSHCommand_UserFlagDescription(t *testing.T) {
	cmd := NewSSHCommand()
	userFlag := cmd.Flags().Lookup("user")
	assert.Contains(t, userFlag.Usage, "username")
}

func TestNewSSHCommand_ConcurrencyFlagDescription(t *testing.T) {
	cmd := NewSSHCommand()
	concurrencyFlag := cmd.Flags().Lookup("concurrency")
	assert.Contains(t, concurrencyFlag.Usage, "multiple requests")
}

func TestNewSSHCommand_TimeoutFlagDescription(t *testing.T) {
	cmd := NewSSHCommand()
	timeoutFlag := cmd.Flags().Lookup("timeout")
	assert.Contains(t, timeoutFlag.Usage, "timeout")
}
