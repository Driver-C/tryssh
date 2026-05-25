package control

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGetController(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewGetController(TypeUsers, "", cfg)
	assert.NotNil(t, ctrl)
	assert.Equal(t, TypeUsers, ctrl.getType)
	assert.Equal(t, "", ctrl.getContent)
	assert.Equal(t, cfg, ctrl.configuration)
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestGetController_ExecuteGet_Users_All(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewGetController(TypeUsers, "", cfg)

	output := captureOutput(func() {
		ctrl.ExecuteGet()
	})

	assert.Contains(t, output, "INDEX")
	assert.Contains(t, output, "USER")
	assert.Contains(t, output, "root")
	assert.Contains(t, output, "admin")
}

func TestGetController_ExecuteGet_Users_Specific(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewGetController(TypeUsers, "root", cfg)

	output := captureOutput(func() {
		ctrl.ExecuteGet()
	})

	assert.Contains(t, output, "root")
	// "admin" should not appear as a data line (it may appear in header)
	lines := strings.Split(output, "\n")
	dataLines := []string{}
	for _, line := range lines {
		if strings.Contains(line, "admin") && !strings.Contains(line, "INDEX") {
			dataLines = append(dataLines, line)
		}
	}
	assert.Empty(t, dataLines, "admin should not appear in search results for 'root'")
}

func TestGetController_ExecuteGet_Ports_All(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewGetController(TypePorts, "", cfg)

	output := captureOutput(func() {
		ctrl.ExecuteGet()
	})

	assert.Contains(t, output, "INDEX")
	assert.Contains(t, output, "PORT")
	assert.Contains(t, output, "22")
	assert.Contains(t, output, "2222")
}

func TestGetController_ExecuteGet_Ports_Specific(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewGetController(TypePorts, "22", cfg)

	output := captureOutput(func() {
		ctrl.ExecuteGet()
	})

	assert.Contains(t, output, "22")
}

func TestGetController_ExecuteGet_Passwords_All(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewGetController(TypePasswords, "", cfg)

	output := captureOutput(func() {
		ctrl.ExecuteGet()
	})

	assert.Contains(t, output, "INDEX")
	assert.Contains(t, output, "PASSWORD")
	assert.Contains(t, output, "****")
	assert.Contains(t, output, "****")
	assert.NotContains(t, output, "password123")
	assert.NotContains(t, output, "admin123")
}

func TestGetController_ExecuteGet_Passwords_Specific(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewGetController(TypePasswords, "admin123", cfg)

	output := captureOutput(func() {
		ctrl.ExecuteGet()
	})

	assert.Contains(t, output, "****")
	assert.NotContains(t, output, "admin123")
}

func TestGetController_ExecuteGet_Keys_All(t *testing.T) {
	cfg := newTestMainConfig()
	cfg.Main.Keys = []string{"/path/to/key1", "/path/to/key2"}
	ctrl := NewGetController(TypeKeys, "", cfg)

	output := captureOutput(func() {
		ctrl.ExecuteGet()
	})

	assert.Contains(t, output, "INDEX")
	assert.Contains(t, output, "KEY")
	assert.Contains(t, output, "/path/to/key1")
	assert.Contains(t, output, "/path/to/key2")
}

func TestGetController_ExecuteGet_Keys_Specific(t *testing.T) {
	cfg := newTestMainConfig()
	cfg.Main.Keys = []string{"/path/to/key1", "/path/to/key2"}
	ctrl := NewGetController(TypeKeys, "/path/to/key1", cfg)

	output := captureOutput(func() {
		ctrl.ExecuteGet()
	})

	assert.Contains(t, output, "/path/to/key1")
}

func TestGetController_ExecuteGet_Caches_All(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewGetController(TypeCaches, "", cfg)

	output := captureOutput(func() {
		ctrl.ExecuteGet()
	})

	assert.Contains(t, output, "INDEX")
	assert.Contains(t, output, "CACHE")
	assert.Contains(t, output, "192.168.1.1")
	assert.Contains(t, output, "192.168.1.2")
}

func TestGetController_ExecuteGet_Caches_SpecificIp(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewGetController(TypeCaches, "192.168.1.1", cfg)

	output := captureOutput(func() {
		ctrl.ExecuteGet()
	})

	assert.Contains(t, output, "192.168.1.1")
	// Should contain only the matching IP
	lines := strings.Split(output, "\n")
	found192_2 := false
	for _, line := range lines {
		if strings.Contains(line, "192.168.1.2") {
			found192_2 = true
		}
	}
	assert.False(t, found192_2, "192.168.1.2 should not appear when searching for 192.168.1.1")
}

func TestGetController_ExecuteGet_Caches_NotFound(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewGetController(TypeCaches, "10.0.0.99", cfg)

	output := captureOutput(func() {
		ctrl.ExecuteGet()
	})

	// Header should be printed but no data lines
	assert.Contains(t, output, "INDEX")
	assert.Contains(t, output, "CACHE")
	// No 10.0.0.99 in output since it doesn't exist
	assert.NotContains(t, output, "10.0.0.99")
}

func TestGetController_SearchAndPrint_All(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewGetController(TypeUsers, "", cfg)

	contents := []string{"alpha", "beta", "gamma"}
	output := captureOutput(func() {
		ctrl.searchAndPrint(contents, false)
	})

	assert.Contains(t, output, "alpha")
	assert.Contains(t, output, "beta")
	assert.Contains(t, output, "gamma")
}

func TestGetController_SearchAndPrint_Specific(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewGetController(TypeUsers, "beta", cfg)

	contents := []string{"alpha", "beta", "gamma"}
	output := captureOutput(func() {
		ctrl.searchAndPrint(contents, false)
	})

	assert.Contains(t, output, "beta")
	assert.NotContains(t, output, "alpha")
}

func TestGetController_SearchAndPrint_EmptySlice(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewGetController(TypeUsers, "anything", cfg)

	output := captureOutput(func() {
		ctrl.searchAndPrint([]string{}, false)
	})

	assert.Empty(t, strings.TrimSpace(output))
}

func TestGetController_SearchAndPrint_NotFound(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewGetController(TypeUsers, "missing", cfg)

	contents := []string{"alpha", "beta", "gamma"}
	output := captureOutput(func() {
		ctrl.searchAndPrint(contents, false)
	})

	// No match found, so nothing printed
	assert.Empty(t, strings.TrimSpace(output))
}
