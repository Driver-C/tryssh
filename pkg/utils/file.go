package utils

import (
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	configFileMode = 0644
)

func FileYamlMarshalAndWrite(path string, conf interface{}) bool {
	// Create a directory if it does not exist
	dirPath := filepath.Dir(path)
	if _, err := os.Stat(dirPath); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				Logger.Fatalln("Directory creation failed: ", err)
			}
		} else {
			Logger.Fatalln("An error occurred while searching for the directory: ", dirPath)
		}
	}

	confData, err := yaml.Marshal(conf)
	if err != nil {
		Logger.Fatalln("Configuration file marshal failed: ", err)
	} else {
		err := os.WriteFile(path, confData, configFileMode)
		if err != nil {
			Logger.Fatalln("Configuration file writing failed: ", err)
		}
	}
	return true
}

func ReadFile(filePath string) ([]byte, bool) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		Logger.Errorln("Error reading file: ", err)
		return nil, false
	}
	return content, true
}

func CheckFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateFile(filePath string, perm fs.FileMode) bool {
	file, err := os.Create(filePath)
	if err != nil {
		Logger.Errorln("Create file error: ", err)
		return false
	}
	if err := file.Chmod(perm); err != nil {
		Logger.Errorln("Chmod error: ", err)
		return false
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			Logger.Fatalln("Failed to close file after creating it: ", err)
		}
	}(file)
	return true
}

func UpdateFile(filePath string, fileContent []byte, perm fs.FileMode) bool {
	if err := os.WriteFile(filePath, fileContent, perm); err != nil {
		Logger.Errorln("File writing failed: ", err)
		return false
	}
	return true
}
