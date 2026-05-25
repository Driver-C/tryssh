package control

import (
	"os"
	"testing"

	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/stretchr/testify/assert"
)

func setupDeleteConfig(t *testing.T) (*config.MainConfig, string) {
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

func TestNewDeleteController(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewDeleteController(TypeUsers, "root", cfg)
	assert.NotNil(t, ctrl)
	assert.Equal(t, TypeUsers, ctrl.deleteType)
	assert.Equal(t, "root", ctrl.deleteContent)
	assert.Equal(t, cfg, ctrl.configuration)
}

func TestDeleteController_ExecuteDelete_Users(t *testing.T) {
	cfg, _ := setupDeleteConfig(t)

	ctrl := NewDeleteController(TypeUsers, "root", cfg)
	ctrl.ExecuteDelete()

	assert.NotContains(t, cfg.Main.Users, "root")
	assert.Contains(t, cfg.Main.Users, "admin")
}

func TestDeleteController_ExecuteDelete_Users_NotFound(t *testing.T) {
	cfg, _ := setupDeleteConfig(t)

	originalUsers := make([]string, len(cfg.Main.Users))
	copy(originalUsers, cfg.Main.Users)

	ctrl := NewDeleteController(TypeUsers, "nonexistent", cfg)
	ctrl.ExecuteDelete()

	assert.Equal(t, originalUsers, cfg.Main.Users)
}

func TestDeleteController_ExecuteDelete_Ports(t *testing.T) {
	cfg, _ := setupDeleteConfig(t)

	ctrl := NewDeleteController(TypePorts, "22", cfg)
	ctrl.ExecuteDelete()

	assert.NotContains(t, cfg.Main.Ports, "22")
	assert.Contains(t, cfg.Main.Ports, "2222")
}

func TestDeleteController_ExecuteDelete_Ports_NotFound(t *testing.T) {
	cfg, _ := setupDeleteConfig(t)

	originalPorts := make([]string, len(cfg.Main.Ports))
	copy(originalPorts, cfg.Main.Ports)

	ctrl := NewDeleteController(TypePorts, "9999", cfg)
	ctrl.ExecuteDelete()

	assert.Equal(t, originalPorts, cfg.Main.Ports)
}

func TestDeleteController_ExecuteDelete_Passwords(t *testing.T) {
	cfg, _ := setupDeleteConfig(t)

	ctrl := NewDeleteController(TypePasswords, "password123", cfg)
	ctrl.ExecuteDelete()

	assert.NotContains(t, cfg.Main.Passwords, "password123")
	assert.Contains(t, cfg.Main.Passwords, "admin123")
}

func TestDeleteController_ExecuteDelete_Passwords_NotFound(t *testing.T) {
	cfg, _ := setupDeleteConfig(t)

	originalPasswords := make([]string, len(cfg.Main.Passwords))
	copy(originalPasswords, cfg.Main.Passwords)

	ctrl := NewDeleteController(TypePasswords, "wrongpass", cfg)
	ctrl.ExecuteDelete()

	assert.Equal(t, originalPasswords, cfg.Main.Passwords)
}

func TestDeleteController_ExecuteDelete_Keys(t *testing.T) {
	cfg, _ := setupDeleteConfig(t)
	cfg.Main.Keys = []string{"/path/to/key1", "/path/to/key2"}

	ctrl := NewDeleteController(TypeKeys, "/path/to/key1", cfg)
	ctrl.ExecuteDelete()

	assert.NotContains(t, cfg.Main.Keys, "/path/to/key1")
	assert.Contains(t, cfg.Main.Keys, "/path/to/key2")
}

func TestDeleteController_ExecuteDelete_Keys_NotFound(t *testing.T) {
	cfg, _ := setupDeleteConfig(t)
	cfg.Main.Keys = []string{"/path/to/key1"}

	ctrl := NewDeleteController(TypeKeys, "/path/to/nonexistent", cfg)
	ctrl.ExecuteDelete()

	assert.Equal(t, []string{"/path/to/key1"}, cfg.Main.Keys)
}

func TestDeleteController_ExecuteDelete_Caches(t *testing.T) {
	cfg, _ := setupDeleteConfig(t)

	originalLen := len(cfg.ServerLists)
	ctrl := NewDeleteController(TypeCaches, "192.168.1.1", cfg)
	ctrl.ExecuteDelete()

	assert.Equal(t, originalLen-1, len(cfg.ServerLists))
	for _, s := range cfg.ServerLists {
		assert.NotEqual(t, "192.168.1.1", s.IP)
	}
}

func TestDeleteController_ExecuteDelete_Caches_NotFound(t *testing.T) {
	cfg, _ := setupDeleteConfig(t)

	originalLen := len(cfg.ServerLists)
	ctrl := NewDeleteController(TypeCaches, "10.0.0.99", cfg)
	ctrl.ExecuteDelete()

	assert.Equal(t, originalLen, len(cfg.ServerLists))
}

func TestDeleteController_ExecuteDelete_Caches_EmptyIp(t *testing.T) {
	cfg, _ := setupDeleteConfig(t)

	originalLen := len(cfg.ServerLists)
	ctrl := NewDeleteController(TypeCaches, "", cfg)
	ctrl.ExecuteDelete()

	assert.Equal(t, originalLen, len(cfg.ServerLists))
}

func TestDeleteController_SearchAndDelete_Found(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewDeleteController(TypeUsers, "root", cfg)

	// searchAndDelete uses dc.deleteContent which is "root"
	contents := []string{"a", "root", "c", "root", "d"}
	result := ctrl.searchAndDelete(contents)

	assert.NotNil(t, result)
	// First "root" should be removed
	assert.Equal(t, []string{"a", "c", "root", "d"}, result)
}

func TestDeleteController_SearchAndDelete_FindsCorrectItem(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewDeleteController(TypeUsers, "target", cfg)

	contents := []string{"a", "target", "b", "target", "c"}
	result := ctrl.searchAndDelete(contents)

	assert.NotNil(t, result)
	// First "target" should be removed
	assert.Equal(t, []string{"a", "b", "target", "c"}, result)
}

func TestDeleteController_SearchAndDelete_NotFound(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewDeleteController(TypeUsers, "missing", cfg)

	contents := []string{"a", "b", "c"}
	result := ctrl.searchAndDelete(contents)

	assert.Nil(t, result)
}

func TestDeleteController_SearchAndDelete_EmptySlice(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewDeleteController(TypeUsers, "anything", cfg)

	contents := []string{}
	result := ctrl.searchAndDelete(contents)

	assert.Nil(t, result)
}

func TestDeleteController_updateConfig(t *testing.T) {
	cfg, configPath := setupDeleteConfig(t)

	ctrl := NewDeleteController(TypeUsers, "root", cfg)
	ctrl.updateConfig()

	// Config file should be created
	assert.FileExists(t, configPath)

	data, err := os.ReadFile(configPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
}

func TestDeleteController_ExecuteDelete_Caches_SingleMatch(t *testing.T) {
	cfg, _ := setupDeleteConfig(t)

	// Ensure only one entry matches the IP
	ctrl := NewDeleteController(TypeCaches, "192.168.1.2", cfg)
	ctrl.ExecuteDelete()

	// Should remove exactly one
	assert.Equal(t, 1, len(cfg.ServerLists))
	assert.Equal(t, "192.168.1.1", cfg.ServerLists[0].IP)
}
