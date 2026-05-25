package control

import (
	"testing"
	"time"

	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/stretchr/testify/assert"
)

func newTestMainConfig() *config.MainConfig {
	c := &config.MainConfig{}
	c.Main.Users = []string{"root", "admin"}
	c.Main.Ports = []string{"22", "2222"}
	c.Main.Passwords = []string{"password123", "admin123"}
	c.Main.Keys = []string{}
	c.ServerLists = []config.ServerListConfig{
		{
			IP:       "192.168.1.1",
			Port:     "22",
			User:     "root",
			Password: "password123",
			Alias:    "server1",
		},
		{
			IP:       "192.168.1.2",
			Port:     "22",
			User:     "admin",
			Password: "admin123",
			Alias:    "server2",
		},
	}
	return c
}

func TestNewSSHController(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewSSHController("192.168.1.1", cfg)
	assert.NotNil(t, ctrl)
	assert.Equal(t, "192.168.1.1", ctrl.targetIP)
	assert.Equal(t, cfg, ctrl.configuration)
}

func TestSSHController_TryLogin_WithCacheFound(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewSSHController("192.168.1.1", cfg)

	// TryLogin calls into SSH launchers which need actual SSH connections.
	// We test the configuration resolution and cache lookup logic.
	ctrl.targetIP = config.ResolveAlias(ctrl.targetIP, ctrl.configuration)
	assert.Equal(t, "192.168.1.1", ctrl.targetIP)

	server, idx, found := config.SelectServerCache("root", ctrl.targetIP, ctrl.configuration)
	assert.True(t, found)
	assert.NotNil(t, server)
	assert.Equal(t, 0, idx)
	assert.Equal(t, "192.168.1.1", server.IP)
	assert.Equal(t, "root", server.User)
}

func TestSSHController_TryLogin_CacheNotFound(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewSSHController("10.0.0.1", cfg)

	server, idx, found := config.SelectServerCache("", ctrl.targetIP, ctrl.configuration)
	assert.False(t, found)
	assert.Nil(t, server)
	assert.Equal(t, 0, idx)
}

func TestSSHController_TryLogin_AliasResolution(t *testing.T) {
	cfg := newTestMainConfig()

	// Test that alias resolves to IP
	resolved := config.ResolveAlias("server1", cfg)
	assert.Equal(t, "192.168.1.1", resolved)

	// Test that unknown alias returns the original string
	resolved = config.ResolveAlias("unknown", cfg)
	assert.Equal(t, "unknown", resolved)

	// Test that IP returns IP unchanged
	resolved = config.ResolveAlias("192.168.1.1", cfg)
	assert.Equal(t, "192.168.1.1", resolved)
}

func TestSSHController_TryLogin_UserSpecified(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewSSHController("192.168.1.1", cfg)

	// When user is specified, SelectServerCache filters by user
	_, _, found := config.SelectServerCache("admin", ctrl.targetIP, ctrl.configuration)
	assert.False(t, found, "192.168.1.1 has user=root, not admin")

	server, idx, found := config.SelectServerCache("root", ctrl.targetIP, ctrl.configuration)
	assert.True(t, found)
	assert.Equal(t, 0, idx)
	assert.Equal(t, "root", server.User)
}

func TestSSHController_TryLogin_EmptyUser(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewSSHController("192.168.1.2", cfg)

	// When user is empty, SelectServerCache returns any matching IP
	server, idx, found := config.SelectServerCache("", ctrl.targetIP, ctrl.configuration)
	assert.True(t, found)
	assert.Equal(t, 1, idx)
	assert.Equal(t, "192.168.1.2", server.IP)
}

func TestSSHController_FieldDefaults(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewSSHController("192.168.1.1", cfg)

	assert.Equal(t, false, ctrl.cacheIsFound)
	assert.Equal(t, 0, ctrl.cacheIndex)
	assert.Equal(t, 0, ctrl.concurrency)
	assert.Equal(t, time.Duration(0), ctrl.sshTimeout)
}
