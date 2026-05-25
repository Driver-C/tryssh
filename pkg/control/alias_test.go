package control

import (
	"bytes"
	"os"
	"testing"

	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/stretchr/testify/assert"
)

func setupAliasConfig(t *testing.T) (*config.MainConfig, string) {
	t.Helper()
	tmpDir := t.TempDir()
	configPath := tmpDir + "/.tryssh/tryssh.db"
	knownHostsPath := tmpDir + "/.tryssh/known_hosts"

	// Override default paths
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

func TestNewAliasController(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewAliasController("192.168.1.1", cfg, "myalias")
	assert.NotNil(t, ctrl)
	assert.Equal(t, "192.168.1.1", ctrl.targetIP)
	assert.Equal(t, "myalias", ctrl.alias)
	assert.Equal(t, cfg, ctrl.configuration)
}

func TestAliasController_SetAlias_Success(t *testing.T) {
	cfg, _ := setupAliasConfig(t)

	ctrl := NewAliasController("192.168.1.1", cfg, "newalias")
	result := ctrl.SetAlias()
	assert.True(t, result)
}

func TestAliasController_SetAlias_DuplicateAlias(t *testing.T) {
	cfg, _ := setupAliasConfig(t)

	// "server1" is already set as alias for 192.168.1.1
	ctrl := NewAliasController("192.168.1.2", cfg, "server1")
	result := ctrl.SetAlias()
	assert.False(t, result)
}

func TestAliasController_SetAlias_DuplicateAliasPrintsList(t *testing.T) {
	cfg, _ := setupAliasConfig(t)

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	ctrl := NewAliasController("192.168.1.2", cfg, "server1")
	ctrl.SetAlias()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()
	assert.Contains(t, output, "server1")
}

func TestAliasController_ListAlias_All(t *testing.T) {
	cfg := newTestMainConfig()

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	ctrl := NewAliasController("", cfg, "")
	ctrl.ListAlias()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()
	assert.Contains(t, output, "server1")
	assert.Contains(t, output, "server2")
	assert.Contains(t, output, "192.168.1.1")
	assert.Contains(t, output, "192.168.1.2")
}

func TestAliasController_ListAlias_Specific(t *testing.T) {
	cfg := newTestMainConfig()

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	ctrl := NewAliasController("", cfg, "server1")
	ctrl.ListAlias()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()
	assert.Contains(t, output, "server1")
	assert.Contains(t, output, "192.168.1.1")
	assert.NotContains(t, output, "server2")
}

func TestAliasController_ListAlias_NoneSet(t *testing.T) {
	cfg := newTestMainConfig()
	cfg.ServerLists = []config.ServerListConfig{
		{IP: "10.0.0.1", Port: "22", User: "root", Password: "pass"},
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	ctrl := NewAliasController("", cfg, "")
	ctrl.ListAlias()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()
	assert.Empty(t, output)
}

func TestAliasController_UnsetAlias(t *testing.T) {
	cfg, _ := setupAliasConfig(t)

	ctrl := NewAliasController("192.168.1.1", cfg, "server1")
	result := ctrl.UnsetAlias()
	assert.True(t, result)
}

func TestAliasController_UnsetAlias_NoMatch(t *testing.T) {
	cfg, _ := setupAliasConfig(t)

	ctrl := NewAliasController("192.168.1.1", cfg, "nonexistent")
	result := ctrl.UnsetAlias()
	assert.False(t, result) // No matching alias, should return false
}

func TestAliasController_SetAlias_NoMatchingIp(t *testing.T) {
	cfg, _ := setupAliasConfig(t)

	ctrl := NewAliasController("10.0.0.99", cfg, "uniquealias")
	result := ctrl.SetAlias()
	assert.False(t, result) // No matching IP, should return false
}

func TestAliasController_SetAlias_MultipleMatchingIps(t *testing.T) {
	cfg, _ := setupAliasConfig(t)

	// Add another server entry with the same IP
	cfg.ServerLists = append(cfg.ServerLists, config.ServerListConfig{
		IP:       "192.168.1.1",
		Port:     "2222",
		User:     "admin",
		Password: "pass",
	})

	ctrl := NewAliasController("192.168.1.1", cfg, "multiAlias")
	result := ctrl.SetAlias()
	assert.True(t, result)
}
