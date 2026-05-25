package encrypt

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helper to reset injectable functions after each test
func resetMocks() {
	readPasswordFn = defaultReadPassword
	fatalFn = defaultFatal
	os.Unsetenv("TRYSSH_MASTER_KEY")
	utils.ClearMasterKey()
}

// setupTempConfig creates a temp config dir and overrides defaults. Returns cleanup fn.
func setupTempConfig(t *testing.T, yamlContent string) {
	t.Helper()
	dir := t.TempDir()
	configDir := filepath.Join(dir, ".tryssh")
	require.NoError(t, os.MkdirAll(configDir, 0755))
	configPath := filepath.Join(configDir, config.ConfigFileName)
	knownHostsPath := filepath.Join(dir, "known_hosts")
	require.NoError(t, os.WriteFile(configPath, []byte(yamlContent), 0600))

	origConfigPath := config.DefaultConfigPath
	origKnownHostsPath := config.DefaultKnownHostsPath
	config.DefaultConfigPath = configPath
	config.DefaultKnownHostsPath = knownHostsPath
	t.Cleanup(func() {
		config.DefaultConfigPath = origConfigPath
		config.DefaultKnownHostsPath = origKnownHostsPath
	})
}

const emptyConfigYAML = "main:\n  ports: []\n  users: []\n  passwords: []\n  keys: []\nserverList: []\n"
const configWithPasswordsYAML = "main:\n  ports: [\"22\"]\n  users: [\"root\"]\n  passwords: [\"secret1\", \"secret2\"]\n  keys: []\nserverList:\n  - ip: \"10.0.0.1\"\n    port: \"22\"\n    user: \"root\"\n    password: \"secret3\"\n    alias: \"\"\n"

func TestNewEncryptCommand_Structure(t *testing.T) {
	cmd := NewEncryptCommand()
	assert.Equal(t, "encrypt", cmd.Use)
	assert.Equal(t, "Encrypt passwords in the configuration file", cmd.Short)
	assert.NotEmpty(t, cmd.Long)
}

func TestNewEncryptCommand_Run(t *testing.T) {
	defer resetMocks()
	setupTempConfig(t, emptyConfigYAML)

	os.Setenv("TRYSSH_MASTER_KEY", "testpassword123")
	var buf bytes.Buffer
	fatalFn = func(args ...interface{}) { fmt.Fprint(&buf, args...) }

	cmd := NewEncryptCommand()
	cmd.Run(cmd, []string{})
}

// --- countPlaintextPasswords ---

func TestCountPlaintextPasswords_AllPlaintext(t *testing.T) {
	c := &config.MainConfig{}
	c.Main.Passwords = []string{"pass1", "pass2"}
	c.ServerLists = []config.ServerListConfig{{IP: "10.0.0.1", Password: "secret"}}
	assert.Equal(t, 3, countPlaintextPasswords(c))
}

func TestCountPlaintextPasswords_Mixed(t *testing.T) {
	os.Setenv("TRYSSH_MASTER_KEY", "testpassword123")
	defer resetMocks()

	key, err := utils.GetMasterKey()
	require.NoError(t, err)
	defer utils.ClearMasterKey()

	encPass, err := utils.Encrypt("encrypted", key)
	require.NoError(t, err)

	c := &config.MainConfig{}
	c.Main.Passwords = []string{"plaintext", encPass}
	c.ServerLists = []config.ServerListConfig{{IP: "10.0.0.1", Password: encPass}}
	assert.Equal(t, 1, countPlaintextPasswords(c))
}

func TestCountPlaintextPasswords_Empty(t *testing.T) {
	assert.Equal(t, 0, countPlaintextPasswords(&config.MainConfig{}))
}

func TestCountPlaintextPasswords_EmptyPasswordsOnly(t *testing.T) {
	c := &config.MainConfig{}
	c.Main.Passwords = []string{""}
	c.ServerLists = []config.ServerListConfig{{IP: "10.0.0.1", Password: ""}}
	assert.Equal(t, 0, countPlaintextPasswords(c))
}

func TestCountPlaintextPasswords_AllEncrypted(t *testing.T) {
	os.Setenv("TRYSSH_MASTER_KEY", "testpassword123")
	defer resetMocks()

	key, err := utils.GetMasterKey()
	require.NoError(t, err)
	defer utils.ClearMasterKey()

	encPass, err := utils.Encrypt("secret", key)
	require.NoError(t, err)

	c := &config.MainConfig{}
	c.Main.Passwords = []string{encPass}
	c.ServerLists = []config.ServerListConfig{{IP: "10.0.0.1", Password: encPass}}
	assert.Equal(t, 0, countPlaintextPasswords(c))
}

// --- runEncrypt via env var path ---

func TestRunEncrypt_EnvVar_NoPasswords(t *testing.T) {
	defer resetMocks()
	setupTempConfig(t, emptyConfigYAML)

	os.Setenv("TRYSSH_MASTER_KEY", "testpassword123")
	runEncrypt()
}

func TestRunEncrypt_EnvVar_WithPasswords(t *testing.T) {
	defer resetMocks()
	setupTempConfig(t, configWithPasswordsYAML)

	os.Setenv("TRYSSH_MASTER_KEY", "testpassword123")
	runEncrypt()

	configPath := config.DefaultConfigPath
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)
	content := string(data)
	assert.Contains(t, content, "enc:")
	assert.NotContains(t, content, "secret1")
}

// --- runEncrypt interactive path (injected readPasswordFn) ---

func TestRunEncrypt_Interactive_EmptyPassword(t *testing.T) {
	defer resetMocks()

	var fatalMsg string
	fatalFn = func(args ...interface{}) { fatalMsg = fmt.Sprint(args...) }
	readPasswordFn = func(_ string) ([]byte, error) { return []byte{}, nil }

	runEncrypt()
	assert.Contains(t, fatalMsg, "empty")
}

func TestRunEncrypt_Interactive_ShortPassword(t *testing.T) {
	defer resetMocks()

	var fatalMsg string
	fatalFn = func(args ...interface{}) { fatalMsg = fmt.Sprint(args...) }
	readPasswordFn = func(_ string) ([]byte, error) { return []byte("abc"), nil }

	runEncrypt()
	assert.Contains(t, fatalMsg, "at least 4 characters")
}

func TestRunEncrypt_Interactive_PasswordMismatch(t *testing.T) {
	defer resetMocks()

	callCount := 0
	var fatalMsg string
	fatalFn = func(args ...interface{}) { fatalMsg = fmt.Sprint(args...) }
	readPasswordFn = func(_ string) ([]byte, error) {
		callCount++
		if callCount == 1 {
			return []byte("password1"), nil
		}
		return []byte("password2"), nil
	}

	runEncrypt()
	assert.Contains(t, fatalMsg, "do not match")
}

func TestRunEncrypt_Interactive_ReadError(t *testing.T) {
	defer resetMocks()

	var fatalMsg string
	fatalFn = func(args ...interface{}) { fatalMsg = fmt.Sprint(args...) }
	readPasswordFn = func(_ string) ([]byte, error) { return nil, errors.New("read error") }

	runEncrypt()
	assert.Contains(t, fatalMsg, "Failed to read password")
}

func TestRunEncrypt_Interactive_ConfirmReadError(t *testing.T) {
	defer resetMocks()

	callCount := 0
	var fatalMsg string
	fatalFn = func(args ...interface{}) { fatalMsg = fmt.Sprint(args...) }
	readPasswordFn = func(_ string) ([]byte, error) {
		callCount++
		if callCount == 1 {
			return []byte("password1"), nil
		}
		return nil, errors.New("confirm error")
	}

	runEncrypt()
	assert.Contains(t, fatalMsg, "Failed to read password")
}

func TestRunEncrypt_Interactive_Success(t *testing.T) {
	defer resetMocks()
	setupTempConfig(t, configWithPasswordsYAML)

	callCount := 0
	readPasswordFn = func(_ string) ([]byte, error) {
		callCount++
		return []byte("testpassword123"), nil
	}

	runEncrypt()

	configPath := config.DefaultConfigPath
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)
	content := string(data)
	assert.Contains(t, content, "enc:")
	assert.NotContains(t, content, "secret1")
}

// --- executeEncrypt error paths ---

func TestExecuteEncrypt_MixedPlaintextAndEncrypted(t *testing.T) {
	os.Setenv("TRYSSH_MASTER_KEY", "testpassword123")
	defer resetMocks()

	key, err := utils.GetMasterKey()
	require.NoError(t, err)
	defer utils.ClearMasterKey()

	encPass, err := utils.Encrypt("already-encrypted", key)
	require.NoError(t, err)

	setupTempConfig(t, "main:\n  ports: []\n  users: []\n  passwords: [\"plaintext\", \""+encPass+"\"]\n  keys: []\nserverList: []\n")

	executeEncrypt()

	data, err := os.ReadFile(config.DefaultConfigPath)
	require.NoError(t, err)
	assert.Equal(t, 2, strings.Count(string(data), "enc:"))
}

func TestExecuteEncrypt_LoadConfigError(t *testing.T) {
	defer resetMocks()

	// Point to non-existent path
	origConfigPath := config.DefaultConfigPath
	config.DefaultConfigPath = "/nonexistent_dir/subdir/tryssh.db"
	defer func() { config.DefaultConfigPath = origConfigPath }()

	os.Setenv("TRYSSH_MASTER_KEY", "testpassword123")

	var fatalMsg string
	fatalFn = func(args ...interface{}) { fatalMsg = fmt.Sprint(args...) }

	executeEncrypt()
	assert.NotEmpty(t, fatalMsg)
}

func TestExecuteEncrypt_UpdateConfigError(t *testing.T) {
	defer resetMocks()

	dir := t.TempDir()
	configDir := filepath.Join(dir, ".tryssh")
	require.NoError(t, os.MkdirAll(configDir, 0755))
	configPath := filepath.Join(configDir, config.ConfigFileName)
	require.NoError(t, os.WriteFile(configPath, []byte(configWithPasswordsYAML), 0600))
	knownHostsPath := filepath.Join(dir, "known_hosts")

	// First, load config normally to get a valid config object
	origConfigPath := config.DefaultConfigPath
	origKnownHostsPath := config.DefaultKnownHostsPath
	config.DefaultConfigPath = configPath
	config.DefaultKnownHostsPath = knownHostsPath

	os.Setenv("TRYSSH_MASTER_KEY", "testpassword123")

	var fatalCalled bool
	fatalFn = func(_ ...interface{}) {
		fatalCalled = true
	}

	_ = fatalCalled

	// Skip: UpdateConfig error path covered in pkg/config tests.
	config.DefaultConfigPath = origConfigPath
	config.DefaultKnownHostsPath = origKnownHostsPath
	t.Skip("UpdateConfig error path covered in pkg/config tests")
}
