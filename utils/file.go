package utils

import (
	"gopkg.in/yaml.v3"
	"os"
)

const (
	configFileMode = 0644
)

func FileYamlMarshalAndWrite(path string, conf interface{}) bool {
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
