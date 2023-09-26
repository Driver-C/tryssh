package utils

import (
	"gopkg.in/yaml.v3"
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
