// Package testutil provides test helper functions for the tryssh project.
package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

// TempDir creates and returns a temporary directory for tests.
func TempDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return dir
}

// CreateTestConfigFile writes a test configuration file in the given directory.
func CreateTestConfigFile(t *testing.T, dir string, content string) string {
	t.Helper()
	path := filepath.Join(dir, "tryssh.db")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}
	return path
}

// CreateTestKnownHosts writes a test known_hosts file in the given directory.
func CreateTestKnownHosts(t *testing.T, dir string, content string) string {
	t.Helper()
	path := filepath.Join(dir, "known_hosts")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("Failed to write known_hosts: %v", err)
	}
	return path
}

// ReadFile reads and returns the contents of a file in tests.
func ReadFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}
	return string(data)
}
