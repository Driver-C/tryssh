package control

import (
	"os"
	"testing"

	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/stretchr/testify/assert"
)

func setupCreateConfig(t *testing.T) (*config.MainConfig, string) {
	t.Helper()
	tmpDir := t.TempDir()
	configPath := tmpDir + "/.tryssh/tryssh.db"
	knownHostsPath := tmpDir + "/.tryssh/known_hosts"

	originalConfigPath := config.DefaultConfigPath
	originalKnownHostsPath := config.DefaultKnownHostsPath
	config.DefaultConfigPath = configPath
	config.DefaultKnownHostsPath = knownHostsPath
	t.Cleanup(func() {
		config.DefaultConfigPath = originalConfigPath
		config.DefaultKnownHostsPath = originalKnownHostsPath
	})

	cfg := newTestMainConfig()
	return cfg, configPath
}

func TestNewCreateController(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewCreateController(TypeUsers, "root", cfg)
	assert.NotNil(t, ctrl)
	assert.Equal(t, TypeUsers, ctrl.createType)
	assert.Equal(t, "root", ctrl.createContent)
	assert.Equal(t, cfg, ctrl.configuration)
}

func TestCreateController_ExecuteCreate_Users(t *testing.T) {
	cfg, _ := setupCreateConfig(t)

	ctrl := NewCreateController(TypeUsers, "newuser", cfg)
	ctrl.ExecuteCreate()

	assert.Contains(t, cfg.Main.Users, "newuser")
}

func TestCreateController_ExecuteCreate_Users_Duplicate(t *testing.T) {
	cfg, _ := setupCreateConfig(t)

	ctrl := NewCreateController(TypeUsers, "root", cfg)
	ctrl.ExecuteCreate()

	// "root" already exists; RemoveDuplicate should keep only one
	count := 0
	for _, u := range cfg.Main.Users {
		if u == "root" {
			count++
		}
	}
	assert.Equal(t, 1, count)
}

func TestCreateController_ExecuteCreate_Ports(t *testing.T) {
	cfg, _ := setupCreateConfig(t)

	ctrl := NewCreateController(TypePorts, "2222", cfg)
	ctrl.ExecuteCreate()

	// "2222" already exists
	count := 0
	for _, p := range cfg.Main.Ports {
		if p == "2222" {
			count++
		}
	}
	assert.Equal(t, 1, count)
}

func TestCreateController_ExecuteCreate_NewPort(t *testing.T) {
	cfg, _ := setupCreateConfig(t)

	ctrl := NewCreateController(TypePorts, "3333", cfg)
	ctrl.ExecuteCreate()

	assert.Contains(t, cfg.Main.Ports, "3333")
}

func TestCreateController_ExecuteCreate_Passwords(t *testing.T) {
	cfg, _ := setupCreateConfig(t)

	ctrl := NewCreateController(TypePasswords, "newpass", cfg)
	ctrl.ExecuteCreate()

	assert.Contains(t, cfg.Main.Passwords, "newpass")
}

func TestCreateController_ExecuteCreate_Passwords_Duplicate(t *testing.T) {
	cfg, _ := setupCreateConfig(t)

	ctrl := NewCreateController(TypePasswords, "password123", cfg)
	ctrl.ExecuteCreate()

	count := 0
	for _, p := range cfg.Main.Passwords {
		if p == "password123" {
			count++
		}
	}
	assert.Equal(t, 1, count)
}

func TestCreateController_ExecuteCreate_Keys(t *testing.T) {
	cfg, _ := setupCreateConfig(t)

	ctrl := NewCreateController(TypeKeys, "/path/to/key", cfg)
	ctrl.ExecuteCreate()

	assert.Contains(t, cfg.Main.Keys, "/path/to/key")
}

func TestCreateController_ExecuteCreate_Caches_ValidJSON(t *testing.T) {
	cfg, _ := setupCreateConfig(t)

	cacheJSON := `{"ip":"10.0.0.5","port":"22","user":"testuser","password":"testpass","alias":"testalias"}`
	ctrl := NewCreateController(TypeCaches, cacheJSON, cfg)
	ctrl.ExecuteCreate()

	found := false
	for _, s := range cfg.ServerLists {
		if s.IP == "10.0.0.5" {
			found = true
			assert.Equal(t, "22", s.Port)
			assert.Equal(t, "testuser", s.User)
			assert.Equal(t, "testpass", s.Password)
			assert.Equal(t, "testalias", s.Alias)
			break
		}
	}
	assert.True(t, found, "new cache should be added to server lists")
}

func TestCreateController_ExecuteCreate_Caches_InvalidJSON(t *testing.T) {
	cfg, _ := setupCreateConfig(t)

	originalLen := len(cfg.ServerLists)
	ctrl := NewCreateController(TypeCaches, "not-valid-json", cfg)
	ctrl.ExecuteCreate()

	// Server lists should not change on invalid JSON
	assert.Equal(t, originalLen, len(cfg.ServerLists))
}

func TestCreateController_ExecuteCreate_Caches_EmptyJSON(t *testing.T) {
	cfg, _ := setupCreateConfig(t)

	cacheJSON := `{"ip":"","port":"","user":"","password":"","alias":""}`
	ctrl := NewCreateController(TypeCaches, cacheJSON, cfg)
	ctrl.ExecuteCreate()

	// Empty JSON is valid and should still add the cache entry
	assert.True(t, len(cfg.ServerLists) > 2)
}

func TestCreateController_updateConfig(t *testing.T) {
	cfg, configPath := setupCreateConfig(t)

	ctrl := NewCreateController(TypeUsers, "testuser", cfg)
	ctrl.updateConfig()

	// Config file should be created
	assert.FileExists(t, configPath)

	// Verify file is readable
	data, err := os.ReadFile(configPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
}

func TestCacheContent_JSONParsing(t *testing.T) {
	jsonStr := `{"ip":"10.0.0.1","port":"22","user":"root","password":"pass","alias":"test"}`

	// Test the struct definition via the exported ExecuteCreate
	cfg, _ := setupCreateConfig(t)
	ctrl := NewCreateController(TypeCaches, jsonStr, cfg)
	ctrl.ExecuteCreate()

	found := false
	for _, s := range cfg.ServerLists {
		if s.IP == "10.0.0.1" && s.Alias == "test" {
			found = true
		}
	}
	assert.True(t, found)
}
