package control

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/stretchr/testify/assert"
)

func setupPruneConfig(t *testing.T) (*config.MainConfig, string) {
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

func TestNewPruneController(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewPruneController(cfg, true, 5*time.Second, 3)
	assert.NotNil(t, ctrl)
	assert.Equal(t, true, ctrl.auto)
	assert.Equal(t, 5*time.Second, ctrl.sshTimeout)
	assert.Equal(t, 3, ctrl.concurrency)
	assert.Equal(t, cfg, ctrl.configuration)
}

func TestPruneController_AutoMode(t *testing.T) {
	cfg, _ := setupPruneConfig(t)

	ctrl := NewPruneController(cfg, true, 1*time.Second, 2)

	// In auto mode, concurrencyDeleteCache is called.
	// Since we can't actually SSH, all caches will be marked for deletion.
	ctrl.PruneCaches()

	assert.Equal(t, 0, len(cfg.ServerLists))
}

func TestPruneController_AutoMode_EmptyServerList(t *testing.T) {
	cfg, _ := setupPruneConfig(t)
	cfg.ServerLists = []config.ServerListConfig{}

	ctrl := NewPruneController(cfg, true, 1*time.Second, 2)
	ctrl.PruneCaches()

	assert.Equal(t, 0, len(cfg.ServerLists))
}

func TestPruneController_AutoMode_Concurrency(t *testing.T) {
	cfg, _ := setupPruneConfig(t)

	// Add multiple servers
	for i := 3; i <= 10; i++ {
		cfg.ServerLists = append(cfg.ServerLists, config.ServerListConfig{
			IP:       "10.0.0." + string(rune('0'+i)),
			Port:     "22",
			User:     "root",
			Password: "pass",
		})
	}

	ctrl := NewPruneController(cfg, true, 1*time.Nanosecond, 5)
	ctrl.PruneCaches()

	// All connections fail with nanosecond timeout, so all should be removed
	assert.Equal(t, 0, len(cfg.ServerLists))
}

func TestPruneController_InteractiveMode_DeleteOne(t *testing.T) {
	cfg, _ := setupPruneConfig(t)
	// Only one server to keep test simple
	cfg.ServerLists = []config.ServerListConfig{
		{
			IP:       "192.168.1.1",
			Port:     "22",
			User:     "root",
			Password: "password123",
		},
	}

	ctrl := NewPruneController(cfg, false, 1*time.Nanosecond, 1)

	// Mock stdin with "yes" response
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	go func() {
		w.WriteString("yes\n")
		w.Close()
	}()

	// Suppress output
	oldStdout := os.Stdout
	_, wOut, _ := os.Pipe()
	os.Stdout = wOut

	ctrl.PruneCaches()

	wOut.Close()
	os.Stdout = oldStdout
	os.Stdin = oldStdin

	assert.Equal(t, 0, len(cfg.ServerLists))
}

func TestPruneController_InteractiveMode_KeepOne(t *testing.T) {
	cfg, _ := setupPruneConfig(t)
	cfg.ServerLists = []config.ServerListConfig{
		{
			IP:       "192.168.1.1",
			Port:     "22",
			User:     "root",
			Password: "password123",
		},
	}

	ctrl := NewPruneController(cfg, false, 1*time.Nanosecond, 1)

	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	go func() {
		w.WriteString("no\n")
		w.Close()
	}()

	oldStdout := os.Stdout
	_, wOut, _ := os.Pipe()
	os.Stdout = wOut

	ctrl.PruneCaches()

	wOut.Close()
	os.Stdout = oldStdout
	os.Stdin = oldStdin

	assert.Equal(t, 1, len(cfg.ServerLists))
}

func TestPruneController_interactiveDeleteCache_Yes(t *testing.T) {
	cfg := newTestMainConfig()

	ctrl := NewPruneController(cfg, false, 1*time.Second, 1)

	server := config.ServerListConfig{
		IP:       "10.0.0.1",
		Port:     "22",
		User:     "root",
		Password: "pass",
	}

	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	go func() {
		w.WriteString("yes\n")
		w.Close()
	}()

	oldStdout := os.Stdout
	_, wOut, _ := os.Pipe()
	os.Stdout = wOut

	result := ctrl.interactiveDeleteCache(server)

	wOut.Close()
	os.Stdout = oldStdout
	os.Stdin = oldStdin

	assert.True(t, result)
}

func TestPruneController_interactiveDeleteCache_No(t *testing.T) {
	cfg := newTestMainConfig()

	ctrl := NewPruneController(cfg, false, 1*time.Second, 1)

	server := config.ServerListConfig{
		IP:       "10.0.0.1",
		Port:     "22",
		User:     "root",
		Password: "pass",
	}

	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	go func() {
		w.WriteString("no\n")
		w.Close()
	}()

	oldStdout := os.Stdout
	_, wOut, _ := os.Pipe()
	os.Stdout = wOut

	result := ctrl.interactiveDeleteCache(server)

	wOut.Close()
	os.Stdout = oldStdout
	os.Stdin = oldStdin

	assert.False(t, result)
}

func TestPruneController_interactiveDeleteCache_InvalidThenYes(t *testing.T) {
	cfg := newTestMainConfig()

	ctrl := NewPruneController(cfg, false, 1*time.Second, 1)

	server := config.ServerListConfig{
		IP:       "10.0.0.1",
		Port:     "22",
		User:     "root",
		Password: "pass",
	}

	// Mock stdin with invalid input then "yes"
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	go func() {
		w.WriteString("invalid\n")
		w.WriteString("yes\n")
		w.Close()
	}()

	// Suppress output
	oldStdout := os.Stdout
	_, wOut, _ := os.Pipe()
	os.Stdout = wOut

	result := ctrl.interactiveDeleteCache(server)

	wOut.Close()
	os.Stdout = oldStdout
	os.Stdin = oldStdin

	assert.True(t, result)
}

func TestPruneController_ConcurrencyDeleteCache(t *testing.T) {
	cfg, _ := setupPruneConfig(t)

	ctrl := NewPruneController(cfg, true, 1*time.Nanosecond, 3)

	newList := ctrl.concurrencyDeleteCache()

	// All connections fail, so nothing should survive
	assert.Equal(t, 0, len(newList))
}

func TestPruneController_ConcurrencyDeleteCache_Empty(t *testing.T) {
	cfg, _ := setupPruneConfig(t)
	cfg.ServerLists = []config.ServerListConfig{}

	ctrl := NewPruneController(cfg, true, 1*time.Second, 3)

	newList := ctrl.concurrencyDeleteCache()
	assert.Equal(t, 0, len(newList))
}

func TestPruneController_PrintsPrompt(t *testing.T) {
	cfg := newTestMainConfig()

	ctrl := NewPruneController(cfg, false, 1*time.Second, 1)

	server := config.ServerListConfig{
		IP:       "10.0.0.1",
		Port:     "22",
		User:     "root",
		Password: "pass",
	}

	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	go func() {
		w.WriteString("yes\n")
		w.Close()
	}()

	var buf bytes.Buffer
	oldStdout := os.Stdout
	rOut, wOut, _ := os.Pipe()
	os.Stdout = wOut

	// Copy output in background
	done := make(chan struct{})
	go func() {
		io.Copy(&buf, rOut)
		close(done)
	}()

	ctrl.interactiveDeleteCache(server)

	wOut.Close()
	<-done
	os.Stdout = oldStdout
	os.Stdin = oldStdin

	output := buf.String()
	assert.True(t, strings.Contains(output, "yes/no"), "should contain prompt")
}
