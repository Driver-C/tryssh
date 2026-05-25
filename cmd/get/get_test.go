package get

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/stretchr/testify/assert"
)

// setupTempConfig creates a temporary config file and overrides the default paths.
// Returns a cleanup function that must be called when the test is done.
func setupTempConfig(t *testing.T) func() {
	t.Helper()
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "tryssh.db")
	knownHostsPath := filepath.Join(tmpDir, "known_hosts")

	// Write an empty config file
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

// captureOutput captures stdout during the execution of fn.
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

func TestNewGetCommand_Structure(t *testing.T) {
	cmd := NewGetCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "get [command]", cmd.Use)
	assert.Contains(t, cmd.Short, "Get alternative")
}

func TestNewGetCommand_Subcommands(t *testing.T) {
	cmd := NewGetCommand()

	expectedSubcommands := []string{"users", "ports", "passwords", "caches", "keys"}
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

func TestNewGetCommand_SubcommandCount(t *testing.T) {
	cmd := NewGetCommand()
	assert.Len(t, cmd.Commands(), 5, "get command should have exactly 5 subcommands")
}

func TestNewGetCommand_NoAliases(t *testing.T) {
	cmd := NewGetCommand()
	assert.Empty(t, cmd.Aliases, "get command should have no aliases")
}

func TestNewGetCommand_Long(t *testing.T) {
	cmd := NewGetCommand()
	assert.Equal(t, "Get alternative username, port number, password, and login cache information", cmd.Long)
}

// --- Users command ---

func TestNewUsersCommand_Structure(t *testing.T) {
	cmd := NewUsersCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "users <username>", cmd.Use)
	assert.Equal(t, "Get alternative usernames", cmd.Short)
	assert.Equal(t, []string{"user", "usr"}, cmd.Aliases)
	assert.NotNil(t, cmd.Run)
}

func TestNewUsersCommand_Long(t *testing.T) {
	cmd := NewUsersCommand()
	assert.Equal(t, "Get alternative usernames", cmd.Long)
}

func TestNewUsersCommand_Run_NoArgs(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewUsersCommand()
	output := captureOutput(t, func() {
		cmd.Run(cmd, []string{})
	})
	assert.Contains(t, output, "INDEX\tUSER")
}

func TestNewUsersCommand_Run_WithArg(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewUsersCommand()
	output := captureOutput(t, func() {
		cmd.Run(cmd, []string{"testuser"})
	})
	assert.Contains(t, output, "INDEX\tUSER")
}

// --- Ports command ---

func TestNewPortsCommand_Structure(t *testing.T) {
	cmd := NewPortsCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "ports <port>", cmd.Use)
	assert.Equal(t, "Get alternative ports", cmd.Short)
	assert.Equal(t, []string{"port", "po"}, cmd.Aliases)
	assert.NotNil(t, cmd.Run)
}

func TestNewPortsCommand_Long(t *testing.T) {
	cmd := NewPortsCommand()
	assert.Equal(t, "Get alternative ports", cmd.Long)
}

func TestNewPortsCommand_Run_NoArgs(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewPortsCommand()
	output := captureOutput(t, func() {
		cmd.Run(cmd, []string{})
	})
	assert.Contains(t, output, "INDEX\tPORT")
}

func TestNewPortsCommand_Run_WithArg(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewPortsCommand()
	output := captureOutput(t, func() {
		cmd.Run(cmd, []string{"22"})
	})
	assert.Contains(t, output, "INDEX\tPORT")
}

// --- Passwords command ---

func TestNewPasswordsCommand_Structure(t *testing.T) {
	cmd := NewPasswordsCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "passwords <password>", cmd.Use)
	assert.Equal(t, "Get alternative passwords", cmd.Short)
	assert.Equal(t, []string{"password", "pass", "pwd"}, cmd.Aliases)
	assert.NotNil(t, cmd.Run)
}

func TestNewPasswordsCommand_Long(t *testing.T) {
	cmd := NewPasswordsCommand()
	assert.Equal(t, "Get alternative passwords", cmd.Long)
}

func TestNewPasswordsCommand_Run_NoArgs(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewPasswordsCommand()
	output := captureOutput(t, func() {
		cmd.Run(cmd, []string{})
	})
	assert.Contains(t, output, "INDEX\tPASSWORD")
}

func TestNewPasswordsCommand_Run_WithArg(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewPasswordsCommand()
	output := captureOutput(t, func() {
		cmd.Run(cmd, []string{"mypassword"})
	})
	assert.Contains(t, output, "INDEX\tPASSWORD")
}

// --- Keys command ---

func TestNewKeysCommand_Structure(t *testing.T) {
	cmd := NewKeysCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "keys <keyFilePath>", cmd.Use)
	assert.Equal(t, "Get alternative key file path", cmd.Short)
	assert.Equal(t, []string{"key"}, cmd.Aliases)
	assert.NotNil(t, cmd.Run)
}

func TestNewKeysCommand_Long(t *testing.T) {
	cmd := NewKeysCommand()
	assert.Equal(t, "Get alternative key file path", cmd.Long)
}

func TestNewKeysCommand_Run_NoArgs(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewKeysCommand()
	output := captureOutput(t, func() {
		cmd.Run(cmd, []string{})
	})
	assert.Contains(t, output, "INDEX\tKEY")
}

func TestNewKeysCommand_Run_WithArg(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewKeysCommand()
	output := captureOutput(t, func() {
		cmd.Run(cmd, []string{"/path/to/key"})
	})
	assert.Contains(t, output, "INDEX\tKEY")
}

// --- Caches command ---

func TestNewCachesCommand_Structure(t *testing.T) {
	cmd := NewCachesCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "caches [ipAddress]", cmd.Use)
	assert.Equal(t, "Get alternative caches by ipAddress", cmd.Short)
	assert.Equal(t, []string{"cache"}, cmd.Aliases)
	assert.NotNil(t, cmd.Run)
}

func TestNewCachesCommand_Long(t *testing.T) {
	cmd := NewCachesCommand()
	assert.Equal(t, "Get alternative caches by ipAddress", cmd.Long)
}

func TestNewCachesCommand_Run_NoArgs(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewCachesCommand()
	output := captureOutput(t, func() {
		cmd.Run(cmd, []string{})
	})
	assert.Contains(t, output, "INDEX\tCACHE")
}

func TestNewCachesCommand_Run_WithArg(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewCachesCommand()
	output := captureOutput(t, func() {
		cmd.Run(cmd, []string{"192.168.1.1"})
	})
	assert.Contains(t, output, "INDEX\tCACHE")
}
