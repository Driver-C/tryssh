package config

import (
	"github.com/Driver-C/tryssh/utils"
	"github.com/schwarmco/go-cartesian-product"
	"gopkg.in/yaml.v3"
	"os"
	"os/user"
	"path/filepath"
)

const (
	configFileName     = "tryssh.db"
	configDirName      = ".tryssh"
	knownHostsFileName = "known_hosts"
)

var (
	configPath     string
	KnownHostsPath string
)

func init() {
	if usr, err := user.Current(); err != nil {
		utils.Logger.Warnf("Unable to obtain current user information: %s, "+
			"Will use the current directory as the configuration file directory.", err)
		configPath = filepath.Join("./", configDirName, configFileName)
		KnownHostsPath = filepath.Join("./", configDirName, knownHostsFileName)
	} else {
		configPath = filepath.Join(usr.HomeDir, configDirName, configFileName)
		KnownHostsPath = filepath.Join(usr.HomeDir, configDirName, knownHostsFileName)
	}
}

// MainConfig Main config
type MainConfig struct {
	Main struct {
		Ports     []string `yaml:"ports,flow"`
		Users     []string `yaml:"users,flow"`
		Passwords []string `yaml:"passwords,flow"`
		Keys      []string `yaml:"keys,flow"`
	} `yaml:"main"`
	ServerLists []ServerListConfig `yaml:"serverList"`
}

// ServerListConfig Server information cache list
type ServerListConfig struct {
	Ip       string `yaml:"ip"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Key      string `yaml:"key"`
	Alias    string `yaml:"alias"`
}

// generateConfig Generate initial configuration file (force overwrite)
func generateConfig() {
	utils.Logger.Infoln("Generating configuration file.\n")
	_ = utils.FileYamlMarshalAndWrite(configPath, &MainConfig{})
	utils.Logger.Infoln("Generating configuration file successful.\n")
	utils.Logger.Warnln("Main setting is empty. " +
		"You need to create some users, ports and passwords before running again.\n")
}

func LoadConfig() (c *MainConfig) {
	c = new(MainConfig)

	if utils.CheckFileIsExist(configPath) {
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
					"You need to create some users, ports and passwords before running again.\n")
			}
		}
	} else {
		utils.Logger.Infoln("Configuration file cannot be found, it will be generated automatically.\n")
		generateConfig()
	}

	// known_hosts
	if !utils.CheckFileIsExist(KnownHostsPath) {
		// Default permission is 0600
		if !utils.CreateFile(KnownHostsPath, 0600) {
			utils.Logger.Fatalln("The known_hosts file creation failed")
		}
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
	keys := utils.InterfaceSlice(conf.Main.Keys)
	// Generate combinations with immutable parameter order
	combinations = cartesian.Iter(ips, ports, users, passwords, keys)
	return
}
