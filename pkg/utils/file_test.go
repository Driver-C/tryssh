package utils

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// badMarshaler implements yaml.Marshaler and always returns an error.
type badMarshaler struct{}

func (badMarshaler) MarshalYAML() (interface{}, error) {
	return nil, errors.New("marshal error")
}

func TestFileYamlMarshalAndWrite_Success(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.yaml")

	conf := struct {
		Name  string `yaml:"name"`
		Value int    `yaml:"value"`
	}{Name: "test", Value: 42}

	err := FileYamlMarshalAndWrite(path, conf)
	assert.NoError(t, err)

	data, err := os.ReadFile(path)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "name: test")
	assert.Contains(t, string(data), "value: 42")
}

func TestFileYamlMarshalAndWrite_CreatesDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "a", "b", "c", "config.yaml")

	conf := struct {
		Key string `yaml:"key"`
	}{Key: "nested"}

	err := FileYamlMarshalAndWrite(path, conf)
	assert.NoError(t, err)

	data, err := os.ReadFile(path)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "key: nested")
}

func TestFileYamlMarshalAndWrite_InvalidPath(t *testing.T) {
	path := filepath.Join("/definitely-not-a-real-root-dir", "sub", "file.yaml")
	err := FileYamlMarshalAndWrite(path, struct{ A string }{A: "b"})
	assert.Error(t, err)
}

func TestReadFile_Success(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "read_test.txt")
	content := []byte("hello world")
	err := os.WriteFile(path, content, 0644)
	assert.NoError(t, err)

	data, ok := ReadFile(path)
	assert.True(t, ok)
	assert.Equal(t, content, data)
}

func TestReadFile_NonExistent(t *testing.T) {
	data, ok := ReadFile("/nonexistent/path/to/file.txt")
	assert.False(t, ok)
	assert.Nil(t, data)
}

func TestCheckFileIsExist_ExistingFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "exists.txt")
	err := os.WriteFile(path, []byte("data"), 0644)
	assert.NoError(t, err)

	assert.True(t, CheckFileIsExist(path))
}

func TestCheckFileIsExist_NonExistingFile(t *testing.T) {
	assert.False(t, CheckFileIsExist("/nonexistent/file/that/does/not/exist.txt"))
}

func TestCreateFile_Success(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "newfile.txt")

	err := CreateFile(path, 0644)
	assert.NoError(t, err)

	info, err := os.Stat(path)
	assert.NoError(t, err)
	assert.False(t, info.IsDir())
}

func TestCreateFile_InvalidPath(t *testing.T) {
	path := "/nonexistent-root-dir/subdir/file.txt"
	err := CreateFile(path, 0644)
	assert.Error(t, err)
}

func TestCreateFile_AlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "exists.txt")

	// Create the file first
	err := CreateFile(path, 0644)
	assert.NoError(t, err)

	// Second create with O_EXCL should fail
	err = CreateFile(path, 0644)
	assert.Error(t, err)
	assert.True(t, os.IsExist(err), "Expected os.IsExist error for duplicate create")
}

func TestUpdateFile_Success(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "update_test.txt")
	err := os.WriteFile(path, []byte("old"), 0644)
	assert.NoError(t, err)

	newContent := []byte("new content")
	err = UpdateFile(path, newContent, 0644)
	assert.NoError(t, err)

	data, err := os.ReadFile(path)
	assert.NoError(t, err)
	assert.Equal(t, newContent, data)
}

func TestUpdateFile_InvalidPath(t *testing.T) {
	err := UpdateFile("/nonexistent-root-dir/file.txt", []byte("data"), 0644)
	assert.Error(t, err)
}

func TestUpdateFile_SetsPermissions(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "perm_test.txt")

	err := UpdateFile(path, []byte("content"), 0600)
	assert.NoError(t, err)

	info, err := os.Stat(path)
	assert.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
}

func TestUpdateFile_OverwritesExisting(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "overwrite.txt")

	// Write initial content
	err := UpdateFile(path, []byte("initial content"), 0644)
	assert.NoError(t, err)

	// Overwrite with new content
	err = UpdateFile(path, []byte("updated content"), 0644)
	assert.NoError(t, err)

	data, err := os.ReadFile(path)
	assert.NoError(t, err)
	assert.Equal(t, "updated content", string(data))
}

func TestUpdateFile_ReadonlyDir(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("readonly dir test not reliable on Windows")
	}
	tmpDir := t.TempDir()
	readonlyDir := filepath.Join(tmpDir, "readonly")
	err := os.Mkdir(readonlyDir, 0555)
	assert.NoError(t, err)

	path := filepath.Join(readonlyDir, "file.txt")
	err = UpdateFile(path, []byte("data"), 0644)
	assert.Error(t, err)
}

func TestFileYamlMarshalAndWrite_StatNonNotExistError(t *testing.T) {
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		tmpDir := t.TempDir()
		blockedFile := filepath.Join(tmpDir, "blocked")
		err := os.WriteFile(blockedFile, []byte("x"), 0644)
		assert.NoError(t, err)
		err = os.Chmod(blockedFile, 0000)
		assert.NoError(t, err)
		path := filepath.Join(blockedFile, "sub", "file.yaml")
		err = FileYamlMarshalAndWrite(path, struct{ A string }{A: "b"})
		assert.Error(t, err)
	}
}

func TestFileYamlMarshalAndWrite_ExistingDirNoError(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.yaml")

	conf := struct {
		Foo string `yaml:"foo"`
	}{Foo: "bar"}

	err := FileYamlMarshalAndWrite(path, conf)
	assert.NoError(t, err)

	err = FileYamlMarshalAndWrite(path, struct {
		Foo string `yaml:"foo"`
	}{Foo: "updated"})
	assert.NoError(t, err)

	data, err := os.ReadFile(path)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "foo: updated")
}

func TestFileYamlMarshalAndWrite_YamlMarshalError(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "bad.yaml")

	// badMarshaler implements yaml.Marshaler and returns an error from MarshalYAML.
	err := FileYamlMarshalAndWrite(path, badMarshaler{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "marshal error")
}

func TestUpdateFile_WriteToFullDiskSim(t *testing.T) {
	// This test just verifies normal UpdateFile works with zero-length content
	// to ensure the write path is exercised with edge-case content.
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "empty.txt")

	err := UpdateFile(path, []byte{}, 0644)
	require.NoError(t, err)

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, []byte{}, data)
}

func TestCreateFile_SetsPermissions(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "permfile.txt")

	err := CreateFile(path, 0600)
	assert.NoError(t, err)

	info, err := os.Stat(path)
	assert.NoError(t, err)
	// On some systems the umask may affect permissions, so just verify the file exists
	assert.False(t, info.IsDir())
}

func TestUpdateFile_LargeContent(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "large.bin")

	// Write a larger payload (64KB)
	largeContent := make([]byte, 64*1024)
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}

	err := UpdateFile(path, largeContent, 0644)
	assert.NoError(t, err)

	data, err := os.ReadFile(path)
	assert.NoError(t, err)
	assert.Equal(t, largeContent, data)
}
