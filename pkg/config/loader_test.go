package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Driver-C/tryssh/pkg/utils"
	"github.com/Driver-C/tryssh/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestLoadConfigFromPath_ValidConfig(t *testing.T) {
	dir := testutil.TempDir(t)
	configContent := `main:
  ports:
  - "22"
  - "2222"
  users:
  - root
  - admin
  passwords:
  - pass1
  keys:
  - /path/to/key
serverList:
- ip: 192.168.1.1
  port: "22"
  user: root
  password: pass1
  alias: server1
`
	configPath := testutil.CreateTestConfigFile(t, dir, configContent)
	knownHostsPath := filepath.Join(dir, "known_hosts")

	conf, err := LoadConfigFromPath(configPath, knownHostsPath)
	assert.NoError(t, err)
	assert.NotNil(t, conf)

	assert.Equal(t, []string{"22", "2222"}, conf.Main.Ports)
	assert.Equal(t, []string{"root", "admin"}, conf.Main.Users)
	assert.Equal(t, []string{"pass1"}, conf.Main.Passwords)
	assert.Equal(t, []string{"/path/to/key"}, conf.Main.Keys)
	assert.Len(t, conf.ServerLists, 1)
	assert.Equal(t, "192.168.1.1", conf.ServerLists[0].IP)
	assert.Equal(t, "22", conf.ServerLists[0].Port)
	assert.Equal(t, "root", conf.ServerLists[0].User)
	assert.Equal(t, "pass1", conf.ServerLists[0].Password)
	assert.Equal(t, "server1", conf.ServerLists[0].Alias)

	// known_hosts should have been created
	assert.FileExists(t, knownHostsPath)
}

func TestLoadConfigFromPath_MissingConfig_GeneratesNew(t *testing.T) {
	dir := testutil.TempDir(t)
	// No config file created -- path points to a non-existent file
	configPath := filepath.Join(dir, ".tryssh", ConfigFileName)
	knownHostsPath := filepath.Join(dir, "known_hosts")

	conf, err := LoadConfigFromPath(configPath, knownHostsPath)
	assert.NoError(t, err)
	assert.NotNil(t, conf)

	// A new empty config should have been generated
	assert.FileExists(t, configPath)
	assert.Empty(t, conf.Main.Ports)
	assert.Empty(t, conf.Main.Users)
	assert.Empty(t, conf.Main.Passwords)
	assert.Empty(t, conf.Main.Keys)
	assert.Empty(t, conf.ServerLists)

	// known_hosts should also have been created
	assert.FileExists(t, knownHostsPath)
}

func TestLoadConfigFromPath_InvalidYAML(t *testing.T) {
	dir := testutil.TempDir(t)
	configContent := `invalid: [yaml: content`
	configPath := testutil.CreateTestConfigFile(t, dir, configContent)
	knownHostsPath := filepath.Join(dir, "known_hosts")

	conf, err := LoadConfigFromPath(configPath, knownHostsPath)
	assert.Error(t, err)
	assert.Nil(t, conf)
	assert.Contains(t, err.Error(), "parsing failed")
}

func TestLoadConfigFromPath_KnownHostsCreation(t *testing.T) {
	dir := testutil.TempDir(t)
	configContent := `main:
  ports: []
  users: []
  passwords: []
  keys: []
`
	configPath := testutil.CreateTestConfigFile(t, dir, configContent)
	knownHostsPath := filepath.Join(dir, "known_hosts")

	// known_hosts does not exist yet
	assert.NoFileExists(t, knownHostsPath)

	conf, err := LoadConfigFromPath(configPath, knownHostsPath)
	assert.NoError(t, err)
	assert.NotNil(t, conf)
	assert.FileExists(t, knownHostsPath)
}

func TestLoadConfigFromPath_KnownHostsAlreadyExists(t *testing.T) {
	dir := testutil.TempDir(t)
	configContent := `main:
  ports: []
  users: []
  passwords: []
  keys: []
`
	configPath := testutil.CreateTestConfigFile(t, dir, configContent)
	existingContent := "existing-host ssh-rsa AAAA...\n"
	knownHostsPath := testutil.CreateTestKnownHosts(t, dir, existingContent)

	conf, err := LoadConfigFromPath(configPath, knownHostsPath)
	assert.NoError(t, err)
	assert.NotNil(t, conf)

	// Existing known_hosts should not be overwritten
	data := testutil.ReadFile(t, knownHostsPath)
	assert.Equal(t, existingContent, data)
}

func TestLoadConfigFromPath_ConfigFileUnreadable(t *testing.T) {
	dir := testutil.TempDir(t)
	configContent := `main:
  ports: ["22"]
  users: ["root"]
  passwords: []
  keys: []
`
	configPath := testutil.CreateTestConfigFile(t, dir, configContent)
	knownHostsPath := filepath.Join(dir, "known_hosts")

	// Remove read permission
	err := os.Chmod(configPath, 0000)
	assert.NoError(t, err)
	defer os.Chmod(configPath, 0644)

	conf, err := LoadConfigFromPath(configPath, knownHostsPath)
	assert.Error(t, err)
	assert.Nil(t, conf)
	assert.Contains(t, err.Error(), "load failed")
}

func TestUpdateConfigAtPath(t *testing.T) {
	dir := testutil.TempDir(t)
	configPath := filepath.Join(dir, ".tryssh", ConfigFileName)
	err := os.MkdirAll(filepath.Dir(configPath), 0755)
	assert.NoError(t, err)

	conf := &MainConfig{}
	conf.Main.Ports = []string{"22", "2222"}
	conf.Main.Users = []string{"root"}
	conf.Main.Passwords = []string{"secret"}
	conf.Main.Keys = []string{"/home/user/.ssh/id_rsa"}
	conf.ServerLists = []ServerListConfig{
		{
			IP:       "10.0.0.1",
			Port:     "22",
			User:     "root",
			Password: "secret",
			Alias:    "myserver",
		},
	}

	err = UpdateConfigAtPath(configPath, conf)
	assert.NoError(t, err)

	// Verify the written file can be parsed back
	data, err := os.ReadFile(configPath)
	assert.NoError(t, err)

	var loaded MainConfig
	err = yaml.Unmarshal(data, &loaded)
	assert.NoError(t, err)
	assert.Equal(t, conf.Main.Ports, loaded.Main.Ports)
	assert.Equal(t, conf.Main.Users, loaded.Main.Users)
	assert.Equal(t, conf.Main.Passwords, loaded.Main.Passwords)
	assert.Equal(t, conf.Main.Keys, loaded.Main.Keys)
	assert.Len(t, loaded.ServerLists, 1)
	assert.Equal(t, "10.0.0.1", loaded.ServerLists[0].IP)
	assert.Equal(t, "myserver", loaded.ServerLists[0].Alias)
}

func TestUpdateConfigAtPath_InvalidPath(t *testing.T) {
	// Use a path in a non-existent directory that cannot be created
	// (e.g., under /proc on Linux or /dev on macOS if restricted)
	conf := &MainConfig{}
	conf.Main.Ports = []string{"22"}

	err := UpdateConfigAtPath("/nonexistent_root_dir/subdir/tryssh.db", conf)
	assert.Error(t, err)
}

func TestGenerateConfig(t *testing.T) {
	dir := testutil.TempDir(t)
	configPath := filepath.Join(dir, ".tryssh", ConfigFileName)

	err := generateConfig(configPath)
	assert.NoError(t, err)
	assert.FileExists(t, configPath)

	// The generated file should be valid YAML representing an empty MainConfig
	data, err := os.ReadFile(configPath)
	assert.NoError(t, err)

	var conf MainConfig
	err = yaml.Unmarshal(data, &conf)
	assert.NoError(t, err)
	assert.Empty(t, conf.Main.Ports)
	assert.Empty(t, conf.Main.Users)
	assert.Empty(t, conf.Main.Passwords)
	assert.Empty(t, conf.Main.Keys)
	assert.Empty(t, conf.ServerLists)
}

func TestDefaultPaths(t *testing.T) {
	configPath, knownHostsPath := DefaultPaths()
	assert.NotEmpty(t, configPath)
	assert.NotEmpty(t, knownHostsPath)
	assert.Contains(t, configPath, ConfigDirName)
	assert.Contains(t, configPath, ConfigFileName)
	assert.Contains(t, knownHostsPath, ConfigDirName)
	assert.Contains(t, knownHostsPath, KnownHostsFileName)
}

func TestLoadConfigFromPath_EmptyConfigFile(t *testing.T) {
	dir := testutil.TempDir(t)
	configPath := testutil.CreateTestConfigFile(t, dir, "")
	knownHostsPath := filepath.Join(dir, "known_hosts")

	conf, err := LoadConfigFromPath(configPath, knownHostsPath)
	assert.NoError(t, err)
	assert.NotNil(t, conf)
	assert.Empty(t, conf.Main.Ports)
	assert.Empty(t, conf.Main.Users)
}

func TestLoadConfig(t *testing.T) {
	// Save and restore defaults
	origConfigPath := DefaultConfigPath
	origKnownHostsPath := DefaultKnownHostsPath
	defer func() {
		DefaultConfigPath = origConfigPath
		DefaultKnownHostsPath = origKnownHostsPath
	}()

	dir := testutil.TempDir(t)
	configDir := filepath.Join(dir, ConfigDirName)
	err := os.MkdirAll(configDir, 0755)
	assert.NoError(t, err)

	configPath := filepath.Join(configDir, ConfigFileName)
	knownHostsPath := filepath.Join(dir, "known_hosts")
	DefaultConfigPath = configPath
	DefaultKnownHostsPath = knownHostsPath

	configContent := `main:
  ports:
  - "22"
  users:
  - root
  passwords:
  - pass1
  keys: []
`
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)

	conf, err := LoadConfig()
	assert.NoError(t, err)
	assert.NotNil(t, conf)
	assert.Equal(t, []string{"22"}, conf.Main.Ports)
	assert.Equal(t, []string{"root"}, conf.Main.Users)
}

func TestLoadConfigFromPath_KnownHostsCreationFailure(t *testing.T) {
	dir := testutil.TempDir(t)
	configContent := `main:
  ports: []
  users: []
  passwords: []
  keys: []
`
	configPath := testutil.CreateTestConfigFile(t, dir, configContent)

	// Create a read-only directory so that creating a file inside it will fail
	readOnlyDir := filepath.Join(dir, "readonly")
	err := os.MkdirAll(readOnlyDir, 0555)
	assert.NoError(t, err)
	knownHostsPath := filepath.Join(readOnlyDir, "known_hosts")

	conf, err := LoadConfigFromPath(configPath, knownHostsPath)
	assert.Error(t, err)
	assert.Nil(t, conf)
	assert.Contains(t, err.Error(), "known_hosts")
}

func TestUpdateConfig(t *testing.T) {
	// Save and restore the default config path
	origPath := DefaultConfigPath
	defer func() { DefaultConfigPath = origPath }()

	dir := testutil.TempDir(t)
	configPath := filepath.Join(dir, ConfigDirName, ConfigFileName)
	DefaultConfigPath = configPath

	conf := &MainConfig{}
	conf.Main.Ports = []string{"22"}
	conf.Main.Users = []string{"testuser"}
	conf.Main.Passwords = []string{"testpass"}
	conf.Main.Keys = []string{}

	err := UpdateConfig(conf)
	assert.NoError(t, err)
	assert.FileExists(t, configPath)
}


func TestServerListConfig_String(t *testing.T) {
	tests := []struct {
		name     string
		config   ServerListConfig
		expected string
	}{
		{
			name:     "all fields populated",
			config:   ServerListConfig{IP: "192.168.1.1", Port: "22", User: "root", Password: "secret123", Key: "/home/user/.ssh/id_rsa", Alias: "myserver"},
			expected: "root@192.168.1.1:22",
		},
		{
			name:     "empty password and key",
			config:   ServerListConfig{IP: "10.0.0.1", Port: "22", User: "admin", Password: "", Key: "", Alias: ""},
			expected: "<empty>",
		},
		{
			name:     "short password",
			config:   ServerListConfig{IP: "10.0.0.2", Port: "2222", User: "root", Password: "ab", Key: ""},
			expected: "****",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.String()
			assert.Contains(t, result, tt.expected)
			// Password should never appear in plaintext
			if tt.config.Password != "" {
				assert.NotContains(t, result, tt.config.Password)
			}
			if tt.config.Key != "" {
				assert.NotContains(t, result, tt.config.Key)
			}
		})
	}
}

func TestEncryptConfigForSave_WithPasswords(t *testing.T) {
	// Test that encryptConfigForSave encrypts passwords when a master key is available
	// Since we can't easily set up a master key in unit tests, test the no-key path
	conf := &MainConfig{}
	conf.Main.Ports = []string{"22"}
	conf.Main.Users = []string{"root"}
	conf.Main.Passwords = []string{"secret"}
	conf.ServerLists = []ServerListConfig{
		{IP: "10.0.0.1", Port: "22", User: "root", Password: "cached_pass"},
	}

	// Without a master key, should return the same config
	result, err := encryptConfigForSave(conf)
	assert.NoError(t, err)
	assert.Equal(t, conf.Main.Passwords, result.Main.Passwords, "without master key, passwords unchanged")
}

func TestEncryptConfigForSave_EncryptsPassword(t *testing.T) {
	os.Unsetenv("TRYSSH_MASTER_KEY")
	clearMasterKeyForTest()

	// Set env var and get key BEFORE clearing cache again
	os.Setenv("TRYSSH_MASTER_KEY", "testpassword123")
	defer os.Unsetenv("TRYSSH_MASTER_KEY")
	key, err := utils.GetMasterKey()
	require.NoError(t, err)

	conf := &MainConfig{}
	conf.Main.Ports = []string{"22"}
	conf.Main.Passwords = []string{"secret", ""}
	conf.ServerLists = []ServerListConfig{
		{IP: "10.0.0.1", Port: "22", User: "root", Password: "cached"},
	}

	result, encErr := encryptConfigForSave(conf)
	require.NoError(t, encErr)
	require.NotNil(t, result)
	// Original should not be modified
	assert.Equal(t, "secret", conf.Main.Passwords[0])
	// Result should be encrypted
	assert.True(t, utils.IsEncrypted(result.Main.Passwords[0]))
	assert.False(t, utils.IsEncrypted(result.Main.Passwords[1]), "empty string should not be encrypted")
	assert.True(t, utils.IsEncrypted(result.ServerLists[0].Password))

	// Verify decryption round-trip
	dec, decErr := utils.Decrypt(result.Main.Passwords[0], key)
	assert.NoError(t, decErr)
	assert.Equal(t, "secret", dec)
}

func TestDecryptConfig_WithEncryptedData(t *testing.T) {
	os.Setenv("TRYSSH_MASTER_KEY", "testpassword123")
	defer os.Unsetenv("TRYSSH_MASTER_KEY")
	clearMasterKeyForTest()

	key := deriveTestKey(t, "testpassword123")
	encPass, err := utils.Encrypt("mysecret", key)
	assert.NoError(t, err)

	conf := &MainConfig{}
	conf.Main.Ports = []string{"22"}
	conf.Main.Passwords = []string{encPass}
	conf.ServerLists = []ServerListConfig{
		{IP: "10.0.0.1", Port: "22", User: "root", Password: encPass},
	}

	err = decryptConfig(conf)
	assert.NoError(t, err)
	assert.Equal(t, "mysecret", conf.Main.Passwords[0])
	assert.Equal(t, "mysecret", conf.ServerLists[0].Password)
}

func TestDecryptConfig_WrongKey(t *testing.T) {
	// Set master key to "testpassword123"
	os.Setenv("TRYSSH_MASTER_KEY", "testpassword123")
	defer os.Unsetenv("TRYSSH_MASTER_KEY")
	clearMasterKeyForTest()

	// Encrypt with a DIFFERENT key derived from "otherpassword123"
	// We need to get a key for "otherpassword123" without messing up env var
	os.Setenv("TRYSSH_MASTER_KEY", "otherpassword123")
	clearMasterKeyForTest()
	otherKey, err := utils.GetMasterKey()
	require.NoError(t, err)

	encPass, encErr := utils.Encrypt("mysecret", otherKey)
	require.NoError(t, encErr)

	// Now switch back to the "correct" master key
	os.Setenv("TRYSSH_MASTER_KEY", "testpassword123")
	clearMasterKeyForTest()

	conf := &MainConfig{}
	conf.Main.Passwords = []string{encPass}

	decErr := decryptConfig(conf)
	assert.Error(t, decErr)
	assert.Contains(t, decErr.Error(), "decrypt")
}

func TestDecryptConfig_NoMasterKey(t *testing.T) {
	os.Unsetenv("TRYSSH_MASTER_KEY")
	clearMasterKeyForTest()

	// Plaintext config — no prompt needed, no error
	conf := &MainConfig{}
	conf.Main.Passwords = []string{"plaintext"}

	err := decryptConfig(conf)
	assert.NoError(t, err)
	assert.Equal(t, "plaintext", conf.Main.Passwords[0])
}

func TestDecryptConfig_EncryptedButNoKey(t *testing.T) {
	os.Unsetenv("TRYSSH_MASTER_KEY")
	clearMasterKeyForTest()

	// Config with encrypted content but no master key available — should error
	conf := &MainConfig{}
	conf.Main.Passwords = []string{"enc:AAAAAA=="}

	err := decryptConfig(conf)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no master key")
}

func clearMasterKeyForTest() {
	utils.ClearMasterKey()
}

func deriveTestKey(t *testing.T, password string) []byte {
	t.Helper()
	key, err := deriveTestKeyBytes([]byte(password))
	assert.NoError(t, err)
	return key
}

func deriveTestKeyBytes(password []byte) ([]byte, error) {
	os.Setenv("TRYSSH_MASTER_KEY", string(password))
	defer os.Unsetenv("TRYSSH_MASTER_KEY")
	return utils.GetMasterKey()
}
