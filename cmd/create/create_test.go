package create

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func TestNewCreateCommand_Structure(t *testing.T) {
	cmd := NewCreateCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "create [command]", cmd.Use)
	assert.Contains(t, cmd.Short, "Create alternative")
	assert.Equal(t, "Create alternative username, port number, password, and login cache information", cmd.Long)
}

func TestNewCreateCommand_Aliases(t *testing.T) {
	cmd := NewCreateCommand()

	expectedAliases := []string{"cre", "crt", "add"}
	assert.Equal(t, expectedAliases, cmd.Aliases)
}

func TestNewCreateCommand_Subcommands(t *testing.T) {
	cmd := NewCreateCommand()

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

func TestNewCreateCommand_SubcommandCount(t *testing.T) {
	cmd := NewCreateCommand()
	assert.Len(t, cmd.Commands(), 5, "create command should have exactly 5 subcommands")
}

// --- Users command ---

func TestNewUsersCommand_Structure(t *testing.T) {
	cmd := NewUsersCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "users <username>", cmd.Use)
	assert.Equal(t, "Create an alternative username", cmd.Short)
	assert.Equal(t, "Create an alternative username", cmd.Long)
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
	cmd.Run(cmd, []string{"testuser"})
}

// --- Ports command ---

func TestNewPortsCommand_Structure(t *testing.T) {
	cmd := NewPortsCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "ports <port>", cmd.Use)
	assert.Equal(t, "Create an alternative port", cmd.Short)
	assert.Equal(t, "Create an alternative port", cmd.Long)
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
	cmd.Run(cmd, []string{"2222"})
}

// --- Passwords command ---

func TestNewPasswordsCommand_Structure(t *testing.T) {
	cmd := NewPasswordsCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "passwords", cmd.Use)
	assert.Equal(t, "Create an alternative password", cmd.Short)
	assert.Equal(t, "Create an alternative password (interactive prompt)", cmd.Long)
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
	assert.Equal(t, "Create an alternative key file path", cmd.Short)
	assert.Equal(t, "Create an alternative key file path", cmd.Long)
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

// --- Caches command ---

func TestNewCachesCommand_Structure(t *testing.T) {
	cmd := NewCachesCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "caches <cache>", cmd.Use)
	assert.Equal(t, "Create an alternative cache", cmd.Short)
	assert.Equal(t, "Create an alternative cache", cmd.Long)
	assert.Equal(t, []string{"cache"}, cmd.Aliases)
	assert.NotNil(t, cmd.Run)
}

func TestNewCachesCommand_Flags(t *testing.T) {
	cmd := NewCachesCommand()

	expectedFlags := []struct {
		name      string
		shorthand string
	}{
		{"ip", "i"},
		{"user", "u"},
		{"port", "P"},
		{"pwd", "p"},
		{"alias", "a"},
	}

	for _, f := range expectedFlags {
		flag := cmd.Flags().Lookup(f.name)
		assert.NotNil(t, flag, "flag %q should exist", f.name)
		assert.Equal(t, f.shorthand, flag.Shorthand, "flag %q shorthand mismatch", f.name)
	}
}

func TestNewCachesCommand_RequiredFlags(t *testing.T) {
	cmd := NewCachesCommand()

	requiredFlags := []string{"ip", "user", "port", "pwd"}
	for _, name := range requiredFlags {
		flag := cmd.Flags().Lookup(name)
		assert.NotNil(t, flag, "flag %q should exist", name)
		annotations := flag.Annotations
		if annotations != nil {
			_, hasRequired := annotations[cobra.BashCompOneRequiredFlag]
			assert.True(t, hasRequired, "flag %q should be marked as required", name)
		}
	}
}

func TestNewCachesCommand_AliasFlagNotRequired(t *testing.T) {
	cmd := NewCachesCommand()

	flag := cmd.Flags().Lookup("alias")
	assert.NotNil(t, flag)
	if flag.Annotations != nil {
		_, hasRequired := flag.Annotations[cobra.BashCompOneRequiredFlag]
		assert.False(t, hasRequired, "alias flag should not be required")
	}
}

func TestNewCachesCommand_FlagDefaultValues(t *testing.T) {
	cmd := NewCachesCommand()

	flags := []string{"ip", "user", "port", "pwd", "alias"}
	for _, name := range flags {
		val, err := cmd.Flags().GetString(name)
		assert.NoError(t, err)
		assert.Equal(t, "", val, "flag %q default should be empty string", name)
	}
}

func TestNewCachesCommand_FlagParsing(t *testing.T) {
	cmd := NewCachesCommand()

	_ = cmd.Flags().Set("ip", "192.168.1.1")
	_ = cmd.Flags().Set("user", "root")
	_ = cmd.Flags().Set("port", "22")
	_ = cmd.Flags().Set("pwd", "secret")
	_ = cmd.Flags().Set("alias", "myserver")

	ipVal, _ := cmd.Flags().GetString("ip")
	assert.Equal(t, "192.168.1.1", ipVal)
	userVal, _ := cmd.Flags().GetString("user")
	assert.Equal(t, "root", userVal)
	portVal, _ := cmd.Flags().GetString("port")
	assert.Equal(t, "22", portVal)
	pwdVal, _ := cmd.Flags().GetString("pwd")
	assert.Equal(t, "secret", pwdVal)
	aliasVal, _ := cmd.Flags().GetString("alias")
	assert.Equal(t, "myserver", aliasVal)
}

func TestNewCachesCommand_Run(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewCachesCommand()
	_ = cmd.Flags().Set("ip", "192.168.1.1")
	_ = cmd.Flags().Set("user", "root")
	_ = cmd.Flags().Set("port", "22")
	_ = cmd.Flags().Set("pwd", "secret")
	_ = cmd.Flags().Set("alias", "myserver")

	cmd.Run(cmd, []string{})
}

func TestNewCachesCommand_Run_WithoutAlias(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewCachesCommand()
	_ = cmd.Flags().Set("ip", "10.0.0.1")
	_ = cmd.Flags().Set("user", "admin")
	_ = cmd.Flags().Set("port", "2222")
	_ = cmd.Flags().Set("pwd", "pass123")

	cmd.Run(cmd, []string{})
}

func TestNewUsersCommand_Run_DuplicateUser(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	// First creation
	cmd := NewUsersCommand()
	cmd.Run(cmd, []string{"testuser"})

	// Second creation with same name (should still work, controller handles dedup)
	cmd2 := NewUsersCommand()
	cmd2.Run(cmd2, []string{"testuser"})
}

func TestNewPortsCommand_Run_DuplicatePort(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewPortsCommand()
	cmd.Run(cmd, []string{"22"})

	cmd2 := NewPortsCommand()
	cmd2.Run(cmd2, []string{"22"})
}

func TestNewKeysCommand_Run_DuplicateKey(t *testing.T) {
	cleanup := setupTempConfig(t)
	defer cleanup()

	cmd := NewKeysCommand()
	cmd.Run(cmd, []string{"/home/user/.ssh/id_rsa"})

	cmd2 := NewKeysCommand()
	cmd2.Run(cmd2, []string{"/home/user/.ssh/id_rsa"})
}

func TestNewCachesCommand_NoArgsConstraint(t *testing.T) {
	cmd := NewCachesCommand()
	// caches command doesn't have an explicit Args validator
	// (it reads everything from flags)
	assert.Nil(t, cmd.Args)
}


func TestNewPasswordsCommand_RunReadPasswordFails(t *testing.T) {
	if os.Getenv("TEST_PASSWORD_CMD") == "1" {
		cmd := NewPasswordsCommand()
		cmd.Run(cmd, []string{})
		return
	}
	subCmd := exec.Command(os.Args[0], "-test.run=TestNewPasswordsCommand_RunReadPasswordFails")
	subCmd.Env = append(os.Environ(), "TEST_PASSWORD_CMD=1")
	subCmd.Stdin = strings.NewReader("test\n")
	output, err := subCmd.CombinedOutput()
	assert.Error(t, err)
	assert.Contains(t, string(output), "Failed to read password")
}

func TestNewPasswordsCommand_StructureDetails(t *testing.T) {
	cmd := NewPasswordsCommand()
	assert.Equal(t, "passwords", cmd.Use)
	assert.Equal(t, []string{"password", "pass", "pwd"}, cmd.Aliases)
	assert.Equal(t, "Create an alternative password", cmd.Short)
	assert.NotNil(t, cmd.Run)
}

func TestNewPasswordsCommand_NoArgsRequired(t *testing.T) {
	cmd := NewPasswordsCommand()
	err := cmd.Args(cmd, []string{})
	assert.NoError(t, err)
}
