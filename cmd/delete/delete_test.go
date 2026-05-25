package delete

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

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
	configData := []byte("main:\n  ports: [\"22\"]\n  users: [\"root\"]\n  passwords: [\"testpass\"]\n  keys: [\"/home/user/.ssh/id_rsa\"]\nserverList:\n  - ip: \"192.168.1.1\"\n    port: \"22\"\n    user: \"root\"\n    password: \"testpass\"\n    key: \"\"\n    alias: \"myserver\"\n")
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

func TestNewDeleteCommand_Structure(t *testing.T) {
	cmd := NewDeleteCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "delete [command]", cmd.Use)
	assert.Contains(t, cmd.Short, "Delete alternative")
	assert.Equal(t, "Delete alternative username, port number, password, and login cache information", cmd.Long)
}

func TestNewDeleteCommand_Aliases(t *testing.T) {
	cmd := NewDeleteCommand()

	assert.Equal(t, []string{"del"}, cmd.Aliases)
}

func TestNewDeleteCommand_Subcommands(t *testing.T) {
	cmd := NewDeleteCommand()

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

func TestNewDeleteCommand_SubcommandCount(t *testing.T) {
	cmd := NewDeleteCommand()
	assert.Len(t, cmd.Commands(), 5, "delete command should have exactly 5 subcommands")
}

// --- Users command ---

func TestNewUsersCommand_Structure(t *testing.T) {
	cmd := NewUsersCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "users <username>", cmd.Use)
	assert.Equal(t, "Delete an alternative username", cmd.Short)
	assert.Equal(t, "Delete an alternative username", cmd.Long)
	assert.Equal(t, []string{"user", "usr"}, cmd.Aliases)
	assert.NotNil(t, cmd.Run)
}

func TestNewUsersCommand_ArgsValidation(t *testing.T) {
	cmd := NewUsersCommand()

	err := cmd.Args(cmd, []string{})
	assert.Error(t, err)

	err = cmd.Args(cmd, []string{"root"})
	assert.NoError(t, err)

	err = cmd.Args(cmd, []string{"root", "extra"})
	assert.Error(t, err)
}

func TestNewUsersCommand_Run(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewUsersCommand()
	cmd.Run(cmd, []string{"root"})
}

func TestNewUsersCommand_Run_NonExistentUser(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewUsersCommand()
	cmd.Run(cmd, []string{"nonexistent"})
}

// --- Ports command ---

func TestNewPortsCommand_Structure(t *testing.T) {
	cmd := NewPortsCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "ports <port>", cmd.Use)
	assert.Equal(t, "Delete an alternative port", cmd.Short)
	assert.Equal(t, "Delete an alternative port", cmd.Long)
	assert.Equal(t, []string{"port", "po"}, cmd.Aliases)
	assert.NotNil(t, cmd.Run)
}

func TestNewPortsCommand_ArgsValidation(t *testing.T) {
	cmd := NewPortsCommand()

	err := cmd.Args(cmd, []string{})
	assert.Error(t, err)

	err = cmd.Args(cmd, []string{"22"})
	assert.NoError(t, err)

	err = cmd.Args(cmd, []string{"22", "extra"})
	assert.Error(t, err)
}

func TestNewPortsCommand_Run(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewPortsCommand()
	cmd.Run(cmd, []string{"22"})
}

func TestNewPortsCommand_Run_NonExistentPort(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewPortsCommand()
	cmd.Run(cmd, []string{"9999"})
}

// --- Passwords command ---

func TestNewPasswordsCommand_Structure(t *testing.T) {
	cmd := NewPasswordsCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "passwords", cmd.Use)
	assert.Equal(t, "Delete an alternative password", cmd.Short)
	assert.Equal(t, "Delete an alternative password (interactive prompt)", cmd.Long)
	assert.Equal(t, []string{"password", "pass", "pwd"}, cmd.Aliases)
	assert.NotNil(t, cmd.Run)
}

func TestNewPasswordsCommand_ArgsValidation(t *testing.T) {
	cmd := NewPasswordsCommand()

	err := cmd.Args(cmd, []string{})
	assert.NoError(t, err)

	err = cmd.Args(cmd, []string{"mypass"})
	assert.Error(t, err)
}

// --- Keys command ---

func TestNewKeysCommand_Structure(t *testing.T) {
	cmd := NewKeysCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "keys <keyFilePath>", cmd.Use)
	assert.Equal(t, "Delete an alternative key file path", cmd.Short)
	assert.Equal(t, "Delete an alternative key file path", cmd.Long)
	assert.Equal(t, []string{"key"}, cmd.Aliases)
	assert.NotNil(t, cmd.Run)
}

func TestNewKeysCommand_ArgsValidation(t *testing.T) {
	cmd := NewKeysCommand()

	err := cmd.Args(cmd, []string{})
	assert.Error(t, err)

	err = cmd.Args(cmd, []string{"/path/to/key"})
	assert.NoError(t, err)

	err = cmd.Args(cmd, []string{"/path/to/key", "extra"})
	assert.Error(t, err)
}

func TestNewKeysCommand_Run(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewKeysCommand()
	cmd.Run(cmd, []string{"/home/user/.ssh/id_rsa"})
}

func TestNewKeysCommand_Run_NonExistentKey(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewKeysCommand()
	cmd.Run(cmd, []string{"/nonexistent/key"})
}

// --- Caches command ---

func TestNewCachesCommand_Structure(t *testing.T) {
	cmd := NewCachesCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "caches <ipAddress>", cmd.Use)
	assert.Equal(t, "Delete an alternative cache", cmd.Short)
	assert.Equal(t, "Delete an alternative cache", cmd.Long)
	assert.Equal(t, []string{"cache"}, cmd.Aliases)
	assert.NotNil(t, cmd.Run)
}

func TestNewCachesCommand_ArgsValidation(t *testing.T) {
	cmd := NewCachesCommand()

	err := cmd.Args(cmd, []string{})
	assert.Error(t, err)

	err = cmd.Args(cmd, []string{"192.168.1.1"})
	assert.NoError(t, err)

	err = cmd.Args(cmd, []string{"192.168.1.1", "extra"})
	assert.Error(t, err)
}

func TestNewCachesCommand_Run(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewCachesCommand()
	cmd.Run(cmd, []string{"192.168.1.1"})
}

func TestNewCachesCommand_Run_NonExistentIP(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewCachesCommand()
	cmd.Run(cmd, []string{"99.99.99.99"})
}

func TestNewDeleteCommand_NoRun(t *testing.T) {
	cmd := NewDeleteCommand()
	// Parent command has no Run function
	assert.Nil(t, cmd.Run)
}


func TestNewPasswordsCommand_RunReadPasswordFails(t *testing.T) {
	if os.Getenv("TEST_DELETE_PASSWORD_CMD") == "1" {
		cmd := NewPasswordsCommand()
		cmd.Run(cmd, []string{})
		return
	}
	subCmd := exec.Command(os.Args[0], "-test.run=TestNewPasswordsCommand_RunReadPasswordFails")
	subCmd.Env = append(os.Environ(), "TEST_DELETE_PASSWORD_CMD=1")
	subCmd.Stdin = strings.NewReader("test\n")
	output, err := subCmd.CombinedOutput()
	assert.Error(t, err)
	assert.Contains(t, string(output), "Failed to read password")
}
