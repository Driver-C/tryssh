package control

import (
	"testing"
	"time"

	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/stretchr/testify/assert"
)

// setupSCPTestConfig creates a config with temp file paths for write operations.
func setupSCPTestConfig(t *testing.T) (*config.MainConfig, string) {
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

func TestTryCopy_SourceRemote_CacheFound(t *testing.T) {
	cfg, _ := setupSCPTestConfig(t)

	ctrl := NewScpController("192.168.1.1:/remote/file.txt", "/local/dest.txt", cfg)
	ctrl.TryCopy("", 1, false, 1*time.Nanosecond)

	// destIP should be resolved from source
	assert.Equal(t, "192.168.1.1", ctrl.destIP)
	assert.Equal(t, "192.168.1.1:/remote/file.txt", ctrl.source)
	assert.Equal(t, "/local/dest.txt", ctrl.destination)
	// Cache should be found for 192.168.1.1 with any user
	assert.True(t, ctrl.cacheIsFound)
}

func TestTryCopy_DestRemote_CacheFound(t *testing.T) {
	cfg, _ := setupSCPTestConfig(t)

	ctrl := NewScpController("/local/file.txt", "192.168.1.1:/remote/path/", cfg)
	ctrl.TryCopy("", 1, false, 1*time.Nanosecond)

	assert.Equal(t, "192.168.1.1", ctrl.destIP)
	assert.Equal(t, "/local/file.txt", ctrl.source)
	assert.Equal(t, "192.168.1.1:/remote/path/", ctrl.destination)
	assert.True(t, ctrl.cacheIsFound)
}

func TestTryCopy_SourceRemote_WithAlias(t *testing.T) {
	cfg, _ := setupSCPTestConfig(t)

	ctrl := NewScpController("server1:/remote/file.txt", "/local/dest.txt", cfg)
	ctrl.TryCopy("", 1, false, 1*time.Nanosecond)

	// Alias should be resolved to IP
	assert.Equal(t, "192.168.1.1", ctrl.destIP)
	assert.Equal(t, "192.168.1.1:/remote/file.txt", ctrl.source)
}

func TestTryCopy_DestRemote_WithAlias(t *testing.T) {
	cfg, _ := setupSCPTestConfig(t)

	ctrl := NewScpController("/local/file.txt", "server1:/remote/path/", cfg)
	ctrl.TryCopy("", 1, false, 1*time.Nanosecond)

	assert.Equal(t, "192.168.1.1", ctrl.destIP)
	assert.Equal(t, "192.168.1.1:/remote/path/", ctrl.destination)
}

func TestTryCopy_NoValidRemotePath(t *testing.T) {
	cfg, _ := setupSCPTestConfig(t)

	ctrl := NewScpController("/local/file.txt", "/local/dest.txt", cfg)
	ctrl.TryCopy("", 1, false, 1*time.Nanosecond)

	// Neither source nor destination is a valid remote path
	assert.Equal(t, "", ctrl.destIP)
	assert.False(t, ctrl.cacheIsFound)
}

func TestTryCopy_IPv6Source(t *testing.T) {
	cfg, _ := setupSCPTestConfig(t)

	// Add a cache for the IPv6 address
	cfg.ServerLists = append(cfg.ServerLists, config.ServerListConfig{
		IP:       "::1",
		Port:     "22",
		User:     "root",
		Password: "testpass",
	})

	ctrl := NewScpController("[::1]:/remote/file.txt", "/local/dest.txt", cfg)
	ctrl.TryCopy("", 1, false, 1*time.Nanosecond)

	assert.Equal(t, "::1", ctrl.destIP)
	assert.Equal(t, "[::1]:/remote/file.txt", ctrl.source)
}

func TestTryCopy_IPv6Dest(t *testing.T) {
	cfg, _ := setupSCPTestConfig(t)

	cfg.ServerLists = append(cfg.ServerLists, config.ServerListConfig{
		IP:       "fe80::1",
		Port:     "22",
		User:     "root",
		Password: "testpass",
	})

	ctrl := NewScpController("/local/file.txt", "[fe80::1]:/remote/path/", cfg)
	ctrl.TryCopy("", 1, false, 1*time.Nanosecond)

	assert.Equal(t, "fe80::1", ctrl.destIP)
	assert.Equal(t, "[fe80::1]:/remote/path/", ctrl.destination)
}

func TestTryCopy_CacheNotFound(t *testing.T) {
	cfg, _ := setupSCPTestConfig(t)

	ctrl := NewScpController("10.0.0.99:/remote/file.txt", "/local/dest.txt", cfg)
	ctrl.TryCopy("", 1, false, 1*time.Nanosecond)

	assert.Equal(t, "10.0.0.99", ctrl.destIP)
	assert.False(t, ctrl.cacheIsFound)
}

func TestTryCopy_WithUserFilter(t *testing.T) {
	cfg, _ := setupSCPTestConfig(t)

	// Cache has user=root, search with user=admin should not find cache
	ctrl := NewScpController("192.168.1.1:/remote/file.txt", "/local/dest.txt", cfg)
	ctrl.TryCopy("admin", 1, false, 1*time.Nanosecond)

	assert.Equal(t, "192.168.1.1", ctrl.destIP)
	assert.False(t, ctrl.cacheIsFound)
}

func TestTryCopy_RecursiveFlag(t *testing.T) {
	cfg, _ := setupSCPTestConfig(t)

	ctrl := NewScpController("192.168.1.1:/remote/dir", "/local/dir", cfg)
	ctrl.TryCopy("", 1, true, 1*time.Nanosecond)

	assert.True(t, ctrl.recursive)
	assert.Equal(t, 1, ctrl.concurrency)
	assert.Equal(t, 1*time.Nanosecond, ctrl.sshTimeout)
}

func TestTryCopy_ConcurrencyAndTimeout(t *testing.T) {
	cfg, _ := setupSCPTestConfig(t)

	ctrl := NewScpController("192.168.1.1:/remote/file", "/local/file", cfg)
	ctrl.TryCopy("", 5, false, 10*time.Second)

	assert.Equal(t, 5, ctrl.concurrency)
	assert.Equal(t, 10*time.Second, ctrl.sshTimeout)
}

func TestTryCopy_SourceParsedFirst(t *testing.T) {
	cfg, _ := setupSCPTestConfig(t)

	// When source is remote, it should be parsed first and destination left alone
	ctrl := NewScpController("192.168.1.1:/remote/file.txt", "/local/path", cfg)
	ctrl.TryCopy("", 1, false, 1*time.Nanosecond)

	assert.Equal(t, "192.168.1.1", ctrl.destIP)
	// source should be reformatted with resolved IP
	assert.Equal(t, "192.168.1.1:/remote/file.txt", ctrl.source)
	assert.Equal(t, "/local/path", ctrl.destination)
}

func TestTryCopy_DestParsedWhenSourceNotRemote(t *testing.T) {
	cfg, _ := setupSCPTestConfig(t)

	ctrl := NewScpController("/local/file.txt", "192.168.1.1:/remote/path", cfg)
	ctrl.TryCopy("", 1, false, 1*time.Nanosecond)

	assert.Equal(t, "192.168.1.1", ctrl.destIP)
	assert.Equal(t, "/local/file.txt", ctrl.source)
	assert.Equal(t, "192.168.1.1:/remote/path", ctrl.destination)
}

func TestTryCopy_NoCredentials(t *testing.T) {
	cfg, _ := setupSCPTestConfig(t)
	cfg.Main.Users = []string{}
	cfg.Main.Ports = []string{}
	cfg.Main.Passwords = []string{}
	cfg.Main.Keys = []string{}

	ctrl := NewScpController("10.0.0.99:/remote/file", "/local/file", cfg)
	// Should complete without panic even with no credentials
	ctrl.TryCopy("", 1, false, 1*time.Nanosecond)

	assert.Equal(t, "10.0.0.99", ctrl.destIP)
	assert.False(t, ctrl.cacheIsFound)
}

func TestNewScpController_Fields(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewScpController("src", "dst", cfg)

	assert.Equal(t, "src", ctrl.source)
	assert.Equal(t, "dst", ctrl.destination)
	assert.Equal(t, cfg, ctrl.configuration)
	assert.Equal(t, "", ctrl.destIP)
	assert.Equal(t, false, ctrl.cacheIsFound)
	assert.Equal(t, 0, ctrl.cacheIndex)
	assert.Equal(t, 0, ctrl.concurrency)
}
