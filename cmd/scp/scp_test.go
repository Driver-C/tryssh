package scp

import (
	"os"
	"path/filepath"
	"strings"
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

func TestNewScpCommand_Structure(t *testing.T) {
	cmd := NewScpCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "scp <source> <destination>", cmd.Use)
	assert.Equal(t, "Upload/Download file to/from the server through SSH protocol", cmd.Short)
	assert.Equal(t, "Upload/Download file to/from the server through SSH protocol", cmd.Long)
	assert.NotNil(t, cmd.Run)
}

func TestNewScpCommand_Example(t *testing.T) {
	cmd := NewScpCommand()

	assert.NotEmpty(t, cmd.Example)
	assert.Contains(t, cmd.Example, "tryssh scp")
	assert.Contains(t, cmd.Example, "192.168.1.1")
	assert.Contains(t, cmd.Example, "-r")
}

func TestNewScpCommand_Flags(t *testing.T) {
	cmd := NewScpCommand()

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

	// Test recursive flag
	recursiveFlag := cmd.Flags().Lookup("recursive")
	assert.NotNil(t, recursiveFlag, "recursive flag should exist")
	assert.Equal(t, "r", recursiveFlag.Shorthand)
	assert.Equal(t, "false", recursiveFlag.DefValue)

	// Test timeout flag
	timeoutFlag := cmd.Flags().Lookup("timeout")
	assert.NotNil(t, timeoutFlag, "timeout flag should exist")
	assert.Equal(t, "t", timeoutFlag.Shorthand)
	assert.Equal(t, "1s", timeoutFlag.DefValue)
}

func TestNewScpCommand_ArgsValidation(t *testing.T) {
	cmd := NewScpCommand()

	// No args should fail
	err := cmd.Args(cmd, []string{})
	assert.Error(t, err)

	// One arg should fail (needs exactly 2)
	err = cmd.Args(cmd, []string{"source"})
	assert.Error(t, err)

	// Two args should pass
	err = cmd.Args(cmd, []string{"./file.txt", "192.168.1.1:/root/"})
	assert.NoError(t, err)

	// Three args should fail
	err = cmd.Args(cmd, []string{"a", "b", "c"})
	assert.Error(t, err)
}

func TestNewScpCommand_FlagValues(t *testing.T) {
	cmd := NewScpCommand()

	concurrencyVal, err := cmd.Flags().GetInt("concurrency")
	assert.NoError(t, err)
	assert.Equal(t, 8, concurrencyVal)

	timeoutVal, err := cmd.Flags().GetDuration("timeout")
	assert.NoError(t, err)
	assert.Equal(t, 1*time.Second, timeoutVal)

	recursiveVal, err := cmd.Flags().GetBool("recursive")
	assert.NoError(t, err)
	assert.False(t, recursiveVal)

	userVal, err := cmd.Flags().GetString("user")
	assert.NoError(t, err)
	assert.Equal(t, "", userVal)
}

func TestNewScpCommand_FlagParsing(t *testing.T) {
	cmd := NewScpCommand()

	err := cmd.Flags().Set("user", "root")
	assert.NoError(t, err)
	userVal, _ := cmd.Flags().GetString("user")
	assert.Equal(t, "root", userVal)

	err = cmd.Flags().Set("concurrency", "4")
	assert.NoError(t, err)
	concurrencyVal, _ := cmd.Flags().GetInt("concurrency")
	assert.Equal(t, 4, concurrencyVal)

	err = cmd.Flags().Set("recursive", "true")
	assert.NoError(t, err)
	recursiveVal, _ := cmd.Flags().GetBool("recursive")
	assert.True(t, recursiveVal)

	err = cmd.Flags().Set("timeout", "3s")
	assert.NoError(t, err)
	timeoutVal, _ := cmd.Flags().GetDuration("timeout")
	assert.Equal(t, 3*time.Second, timeoutVal)
}

func TestNewScpCommand_ExampleContainsDownloadAndUpload(t *testing.T) {
	cmd := NewScpCommand()

	// Should contain download examples (remote to local)
	assert.True(t, strings.Contains(cmd.Example, "192.168.1.1:/root") ||
		strings.Contains(cmd.Example, "Download"),
		"example should contain download usage")

	// Should contain upload examples (local to remote)
	assert.True(t, strings.Contains(cmd.Example, "./test.txt") ||
		strings.Contains(cmd.Example, "Upload"),
		"example should contain upload usage")
}

func TestNewScpCommand_NoAliases(t *testing.T) {
	cmd := NewScpCommand()
	assert.Empty(t, cmd.Aliases, "scp command should have no aliases")
}

func TestNewScpCommand_Run_Download(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewScpCommand()
	_ = cmd.Flags().Set("timeout", "100ms")
	_ = cmd.Flags().Set("concurrency", "1")
	// TryCopy will fail (no server), but the Run closure is exercised
	cmd.Run(cmd, []string{"192.168.255.255:/root/file.txt", "/tmp/"})
}

func TestNewScpCommand_Run_Upload(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewScpCommand()
	_ = cmd.Flags().Set("timeout", "100ms")
	_ = cmd.Flags().Set("concurrency", "1")
	_ = cmd.Flags().Set("user", "root")
	cmd.Run(cmd, []string{"/tmp/file.txt", "192.168.255.255:/root/"})
}

func TestNewScpCommand_Run_Recursive(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewScpCommand()
	_ = cmd.Flags().Set("recursive", "true")
	_ = cmd.Flags().Set("timeout", "100ms")
	cmd.Run(cmd, []string{"/tmp/dir", "192.168.255.255:/root/"})
}

func TestNewScpCommand_UserFlagDescription(t *testing.T) {
	cmd := NewScpCommand()
	userFlag := cmd.Flags().Lookup("user")
	assert.Contains(t, userFlag.Usage, "username")
}

func TestNewScpCommand_RecursiveFlagDescription(t *testing.T) {
	cmd := NewScpCommand()
	recursiveFlag := cmd.Flags().Lookup("recursive")
	assert.Contains(t, recursiveFlag.Usage, "Recursively")
}
