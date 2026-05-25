package control

import (
	"testing"
	"time"

	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/stretchr/testify/assert"
)

// setupSSHTestConfig creates a config with temp file paths for write operations.
func setupSSHTestConfig(t *testing.T) (*config.MainConfig, string) {
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

	cfg := &config.MainConfig{}
	cfg.Main.Users = []string{"root"}
	cfg.Main.Ports = []string{"22"}
	cfg.Main.Passwords = []string{"testpass"}
	cfg.Main.Keys = []string{}
	cfg.ServerLists = []config.ServerListConfig{
		{
			IP:       "192.168.1.1",
			Port:     "22",
			User:     "root",
			Password: "testpass",
			Alias:    "server1",
		},
	}
	return cfg, configPath
}

func TestTryLogin_CacheFound_LaunchFails(t *testing.T) {
	cfg, _ := setupSSHTestConfig(t)

	ctrl := NewSSHController("192.168.1.1", cfg)
	// TryLogin will find cache for 192.168.1.1, attempt to Launch (which fails because
	// no real SSH server), then fall through to tryLoginWithoutCache which also fails.
	// The test verifies the function completes without panicking and that fields are set.
	ctrl.TryLogin("root", 1, 1*time.Nanosecond)

	assert.Equal(t, "192.168.1.1", ctrl.targetIP)
	assert.Equal(t, 1, ctrl.concurrency)
	assert.Equal(t, 1*time.Nanosecond, ctrl.sshTimeout)
	// Cache should have been found (user matches)
	assert.True(t, ctrl.cacheIsFound)
}

func TestTryLogin_CacheNotFound(t *testing.T) {
	cfg, _ := setupSSHTestConfig(t)

	ctrl := NewSSHController("10.0.0.99", cfg)
	ctrl.TryLogin("", 1, 1*time.Nanosecond)

	assert.Equal(t, "10.0.0.99", ctrl.targetIP)
	assert.False(t, ctrl.cacheIsFound)
}

func TestTryLogin_WithAlias(t *testing.T) {
	cfg, _ := setupSSHTestConfig(t)

	ctrl := NewSSHController("server1", cfg)
	ctrl.TryLogin("", 1, 1*time.Nanosecond)

	// Alias should be resolved to IP
	assert.Equal(t, "192.168.1.1", ctrl.targetIP)
	assert.True(t, ctrl.cacheIsFound)
}

func TestTryLogin_WithUserFilter(t *testing.T) {
	cfg, _ := setupSSHTestConfig(t)

	// Cache has user=root, search with user=admin should not find cache
	ctrl := NewSSHController("192.168.1.1", cfg)
	ctrl.TryLogin("admin", 1, 1*time.Nanosecond)

	assert.Equal(t, "192.168.1.1", ctrl.targetIP)
	// No cache match because user differs
	assert.False(t, ctrl.cacheIsFound)
}

func TestTryLogin_EmptyUser(t *testing.T) {
	cfg, _ := setupSSHTestConfig(t)

	ctrl := NewSSHController("192.168.1.1", cfg)
	ctrl.TryLogin("", 1, 1*time.Nanosecond)

	// Empty user should find any cache matching IP
	assert.Equal(t, "192.168.1.1", ctrl.targetIP)
	assert.True(t, ctrl.cacheIsFound)
}

func TestTryLogin_NoCombinations(t *testing.T) {
	cfg, _ := setupSSHTestConfig(t)
	// Remove all credentials so GenerateCombination produces nothing useful
	cfg.Main.Users = []string{}
	cfg.Main.Ports = []string{}
	cfg.Main.Passwords = []string{}
	cfg.Main.Keys = []string{}

	ctrl := NewSSHController("10.0.0.99", cfg)
	ctrl.TryLogin("", 1, 1*time.Nanosecond)

	assert.False(t, ctrl.cacheIsFound)
	// Should complete without panic even with no combinations
}

func TestTryLogin_ConcurrencyZero(t *testing.T) {
	cfg, _ := setupSSHTestConfig(t)

	ctrl := NewSSHController("10.0.0.99", cfg)
	// concurrency=0 should be handled (ConcurrencyTryToConnect treats <1 as 1)
	ctrl.TryLogin("", 0, 1*time.Nanosecond)

	assert.Equal(t, 0, ctrl.concurrency)
}

func TestTryLogin_MultipleCachesForSameIP(t *testing.T) {
	cfg, _ := setupSSHTestConfig(t)

	// Add another cache for same IP with different user
	cfg.ServerLists = append(cfg.ServerLists, config.ServerListConfig{
		IP:       "192.168.1.1",
		Port:     "2222",
		User:     "admin",
		Password: "adminpass",
		Alias:    "server1-admin",
	})

	ctrl := NewSSHController("192.168.1.1", cfg)
	ctrl.TryLogin("", 1, 1*time.Nanosecond)

	assert.Equal(t, "192.168.1.1", ctrl.targetIP)
	// Should find the first matching cache
	assert.True(t, ctrl.cacheIsFound)
	assert.Equal(t, 0, ctrl.cacheIndex)
}

func TestTryLogin_CacheIndexCorrect(t *testing.T) {
	cfg, _ := setupSSHTestConfig(t)

	// Add a cache at index 1 and search for that IP
	cfg.ServerLists = append(cfg.ServerLists, config.ServerListConfig{
		IP:       "10.0.0.5",
		Port:     "22",
		User:     "deploy",
		Password: "deploy123",
	})

	ctrl := NewSSHController("10.0.0.5", cfg)
	ctrl.TryLogin("deploy", 1, 1*time.Nanosecond)

	assert.Equal(t, "10.0.0.5", ctrl.targetIP)
	assert.True(t, ctrl.cacheIsFound)
	assert.Equal(t, 1, ctrl.cacheIndex)
}

func TestTryLogin_ConcurrencyGreaterThanConnectors(t *testing.T) {
	cfg, _ := setupSSHTestConfig(t)
	// Only 1 user, 1 port, 1 password = 1 combination
	cfg.Main.Users = []string{"testuser"}
	cfg.Main.Ports = []string{"2222"}
	cfg.Main.Passwords = []string{"testpass"}
	cfg.Main.Keys = []string{}

	ctrl := NewSSHController("10.0.0.99", cfg)
	// concurrency=10 but only 1 connector should work fine
	ctrl.TryLogin("", 10, 1*time.Nanosecond)

	assert.Equal(t, 10, ctrl.concurrency)
}

func TestNewSSHController_NilConfig(t *testing.T) {
	ctrl := NewSSHController("host", nil)
	assert.NotNil(t, ctrl)
	assert.Equal(t, "host", ctrl.targetIP)
	assert.Nil(t, ctrl.configuration)
}
