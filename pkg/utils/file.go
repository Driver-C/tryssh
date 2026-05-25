package utils

import (
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path/filepath"
)

// ConfigFileMode is the default file permission used for config files.
const ConfigFileMode = 0600

// FileYamlMarshalAndWrite marshals the given value to YAML and writes it atomically
// to the specified path, creating parent directories as needed.
func FileYamlMarshalAndWrite(path string, conf interface{}) error {
	dirPath := filepath.Dir(path)
	if _, err := os.Stat(dirPath); err != nil {
		if os.IsNotExist(err) {
			if mkdirErr := os.MkdirAll(dirPath, 0700); mkdirErr != nil {
				return mkdirErr
			}
		} else {
			return err
		}
	}

	confData, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	return UpdateFile(path, confData, ConfigFileMode)
}

// ReadFile reads the entire file and returns its contents.
func ReadFile(filePath string) ([]byte, bool) {
	content, err := os.ReadFile(filePath) //nolint:gosec // G304: path is from caller-provided config
	if err != nil {
		Errorln("Error reading file: ", err)
		return nil, false
	}
	return content, true
}

// CheckFileIsExist returns true if the file exists (including when unreadable due to permissions).
func CheckFileIsExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || !os.IsNotExist(err)
}

// CreateFile creates an empty file with the specified permissions atomically.
func CreateFile(filePath string, perm fs.FileMode) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, perm) //nolint:gosec // G304: path is from caller-provided config
	if err != nil {
		return err
	}
	return file.Close()
}

// UpdateFile writes the given content to the file with the specified permissions atomically
// using a temporary file and rename to prevent corruption on crash.
func UpdateFile(filePath string, fileContent []byte, perm fs.FileMode) error {
	dir := filepath.Dir(filePath)
	tmpFile, err := os.CreateTemp(dir, ".tryssh-tmp-*")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()

	if _, writeErr := tmpFile.Write(fileContent); writeErr != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpPath)
		return writeErr
	}
	if chmodErr := tmpFile.Chmod(perm); chmodErr != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpPath)
		return chmodErr
	}
	if closeErr := tmpFile.Close(); closeErr != nil {
		_ = os.Remove(tmpPath)
		return closeErr
	}
	if renameErr := os.Rename(tmpPath, filePath); renameErr != nil {
		_ = os.Remove(tmpPath)
		return renameErr
	}
	return nil
}
