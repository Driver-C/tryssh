package config

import (
	"github.com/schwarmco/go-cartesian-product"
	"gopkg.in/yaml.v3"
	"os"
	"tryssh/launcher"
	"tryssh/utils"
)

const (
	configPath = "/usr/local/etc/tryssh.yaml"
)

// MainConfig Main config
type MainConfig struct {
	Main struct {
		Ports     []string `yaml:"ports,flow"`
		Users     []string `yaml:"users,flow"`
		Passwords []string `yaml:"passwords,flow"`
	} `yaml:"main"`
	ServerLists []ServerListConfig `yaml:"serverList"`
}

// ServerListConfig Server information cache list
type ServerListConfig struct {
	Ip       string `yaml:"ip"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Alias    string `yaml:"alias"`
}

// GetSshConnectorFromConfig Get SshConnector by ServerListConfig
func GetSshConnectorFromConfig(conf *ServerListConfig) *launcher.SshConnector {
	return &launcher.SshConnector{
		Ip:       conf.Ip,
		Port:     conf.Port,
		User:     conf.User,
		Password: conf.Password,
	}
}

// GetConfigFromSshConnector Get ServerListConfig by SshConnector
func GetConfigFromSshConnector(tgt *launcher.SshConnector) *ServerListConfig {
	return &ServerListConfig{
		Ip:       tgt.Ip,
		Port:     tgt.Port,
		User:     tgt.User,
		Password: tgt.Password,
	}
}

func checkFileIsExist(filename string) (exist bool) {
	exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return
}

// generateConfig Generate initial configuration file (force overwrite)
func generateConfig() {
	utils.Logger.Infoln("Generating configuration file.\n")
	_ = utils.FileYamlMarshalAndWrite(configPath, &MainConfig{})
	utils.Logger.Infoln("Generating configuration file successful.\n")
	utils.Logger.Warnln("Main setting is empty. " +
		"You need to configure the main configuration before running again.\n")
	os.Exit(0)
}

func LoadConfig() (c *MainConfig) {
	c = new(MainConfig)

	if checkFileIsExist(configPath) {
		conf, err := os.ReadFile(configPath)
		if err != nil {
			utils.Logger.Fatalln("Configuration file load failed: ", err)
		}
		unmarshalErr := yaml.Unmarshal(conf, c)
		if unmarshalErr != nil {
			utils.Logger.Fatalln("Configuration file parsing failed: ", unmarshalErr)
		} else {
			if len(c.Main.Ports) == 0 || len(c.Main.Users) == 0 || len(c.Main.Passwords) == 0 {
				utils.Logger.Warnln("Main setting is empty. " +
					"You need to configure the main configuration before running again.\n")
				os.Exit(0)
			}
		}
	} else {
		utils.Logger.Infoln("Configuration file cannot be found, it will be generated automatically.\n")
		generateConfig()
	}
	return
}

// SelectServerCache Search cache from server list
func SelectServerCache(user string, ip string, conf *MainConfig) (*ServerListConfig, int, bool) {
	for index, server := range conf.ServerLists {
		if server.Ip == ip {
			if user != "" {
				if server.User == user {
					return &server, index, true
				}
			} else {
				return &server, index, true
			}
		}
	}
	return nil, 0, false
}

func UpdateConfig(conf *MainConfig) (writeRes bool) {
	writeRes = utils.FileYamlMarshalAndWrite(configPath, conf)
	return
}

// GenerateCombination Generate objects for all port, user, and password combinations
func GenerateCombination(ip string, user string, conf *MainConfig) (combinations chan []interface{}) {
	ips := []interface{}{ip}
	users := []interface{}{user}
	ports := utils.InterfaceSlice(conf.Main.Ports)
	if user == "" {
		users = utils.InterfaceSlice(conf.Main.Users)
	}
	passwords := utils.InterfaceSlice(conf.Main.Passwords)
	// Generate combinations with immutable parameter order
	combinations = cartesian.Iter(ips, ports, users, passwords)
	return
}
