package version

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVersionCommand_Structure(t *testing.T) {
	cmd := NewVersionCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "version", cmd.Use)
	assert.Equal(t, "Print the client version information for the current context", cmd.Short)
	assert.Equal(t, "Print the client version information for the current context", cmd.Long)
	assert.NotNil(t, cmd.Run)
}

func TestNewVersionCommand_AllFieldsSet(t *testing.T) {
	// Save and restore original values
	origVersion := Version
	origBuildGoVersion := BuildGoVersion
	origBuildTime := BuildTime
	defer func() {
		Version = origVersion
		BuildGoVersion = origBuildGoVersion
		BuildTime = origBuildTime
	}()

	Version = "1.2.3"
	BuildGoVersion = "go1.25.0"
	BuildTime = "2024-01-01"

	cmd := NewVersionCommand()

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd.Run(cmd, []string{})

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	assert.Contains(t, output, "Version: 1.2.3")
	assert.Contains(t, output, "GoVersion: go1.25.0")
	assert.Contains(t, output, "BuildTime: 2024-01-01")
}

func TestNewVersionCommand_EmptyFields(t *testing.T) {
	// Save and restore original values
	origVersion := Version
	origBuildGoVersion := BuildGoVersion
	origBuildTime := BuildTime
	defer func() {
		Version = origVersion
		BuildGoVersion = origBuildGoVersion
		BuildTime = origBuildTime
	}()

	Version = ""
	BuildGoVersion = ""
	BuildTime = ""

	cmd := NewVersionCommand()

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd.Run(cmd, []string{})

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	assert.Contains(t, output, "(dev)",
		"output should show (dev) when all fields are empty")
}

func TestNewVersionCommand_PartialFieldsSet(t *testing.T) {
	// Save and restore original values
	origVersion := Version
	origBuildGoVersion := BuildGoVersion
	origBuildTime := BuildTime
	defer func() {
		Version = origVersion
		BuildGoVersion = origBuildGoVersion
		BuildTime = origBuildTime
	}()

	Version = "0.9.0"
	BuildGoVersion = ""
	BuildTime = ""

	cmd := NewVersionCommand()

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd.Run(cmd, []string{})

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	assert.Contains(t, output, "Version: 0.9.0")
	assert.NotContains(t, output, "GoVersion:")
	assert.NotContains(t, output, "BuildTime:")
}

func TestNewVersionCommand_OnlyGoVersionSet(t *testing.T) {
	origVersion := Version
	origBuildGoVersion := BuildGoVersion
	origBuildTime := BuildTime
	defer func() {
		Version = origVersion
		BuildGoVersion = origBuildGoVersion
		BuildTime = origBuildTime
	}()

	Version = ""
	BuildGoVersion = "go1.21.0"
	BuildTime = ""

	cmd := NewVersionCommand()

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd.Run(cmd, []string{})

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	assert.Contains(t, output, "GoVersion: go1.21.0")
	assert.NotContains(t, output, "BuildTime:")
	// Verify there's no standalone "Version:" line
	for _, line := range strings.Split(output, "\n") {
		assert.False(t, strings.HasPrefix(line, "Version:"),
			"should not have a Version: line, got: %q", line)
	}
}
