package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTempDir(t *testing.T) {
	dir := TempDir(t)
	assert.DirExists(t, dir)
}

func TestCreateTestConfigFile(t *testing.T) {
	dir := t.TempDir()
	content := "main:\n  ports:\n  - \"22\""
	path := CreateTestConfigFile(t, dir, content)
	assert.Equal(t, filepath.Join(dir, "tryssh.db"), path)

	result := ReadFile(t, path)
	assert.Equal(t, content, result)
}

func TestCreateTestConfigFile_EmptyContent(t *testing.T) {
	dir := t.TempDir()
	path := CreateTestConfigFile(t, dir, "")
	assert.FileExists(t, path)
}

func TestCreateTestKnownHosts(t *testing.T) {
	dir := t.TempDir()
	content := "192.168.1.1 ssh-ed25519 AAAA..."
	path := CreateTestKnownHosts(t, dir, content)
	assert.Equal(t, filepath.Join(dir, "known_hosts"), path)

	result := ReadFile(t, path)
	assert.Equal(t, content, result)
}

func TestCreateTestKnownHosts_EmptyContent(t *testing.T) {
	dir := t.TempDir()
	path := CreateTestKnownHosts(t, dir, "")
	assert.FileExists(t, path)
}

func TestReadFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	require.NoError(t, os.WriteFile(path, []byte("hello world"), 0644))

	result := ReadFile(t, path)
	assert.Equal(t, "hello world", result)
}
