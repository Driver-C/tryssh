package config

import (
	"github.com/schwarmco/go-cartesian-product"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"tryssh/launcher"
	"tryssh/utils"
)

const (
	configPath = "/usr/local/etc/tryssh.yaml"
)

// MainConfig 主配置文件
type MainConfig struct {
	Main struct {
		Ports     []string `yaml:"ports,flow"`
		Users     []string `yaml:"users,flow"`
		Passwords []string `yaml:"passwords,flow"`
	} `yaml:"main"`
	ServerLists []ServerListConfig `yaml:"serverList"`
}

// ServerListConfig 服务器信息缓存列表
type ServerListConfig struct {
	Ip       string `yaml:"ip"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// GetSshConnectorFromConfig 通过ServerListConfig获取SshConnector
func GetSshConnectorFromConfig(conf *ServerListConfig) *launcher.SshConnector {
	return &launcher.SshConnector{
		Ip:       conf.Ip,
		Port:     conf.Port,
		User:     conf.User,
		Password: conf.Password,
	}
}

// GetConfigFromSshConnector 通过SshConnector获取ServerListConfig
func GetConfigFromSshConnector(tgt *launcher.SshConnector) *ServerListConfig {
	return &ServerListConfig{
		Ip:       tgt.Ip,
		Port:     tgt.Port,
		User:     tgt.User,
		Password: tgt.Password,
	}
}

// checkFileIsExist 检查文件是否存在
func checkFileIsExist(filename string) (exist bool) {
	exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return
}

// generateConfig 生成初始配置文件（直接覆盖）
func generateConfig() {
	utils.Logger.Infoln("Generating configuration file.\n")
	_ = utils.FileYamlMarshalAndWrite(configPath, &MainConfig{})
	utils.Logger.Infoln("Generating configuration file successful.\n")
	utils.Logger.Warnln("Main setting is empty. " +
		"You need to configure the main configuration before running again.\n")
	os.Exit(0)
}

// LoadConfig 加载配置文件
func LoadConfig() (c *MainConfig) {
	c = new(MainConfig)

	if checkFileIsExist(configPath) {
		conf, err := ioutil.ReadFile(configPath)
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

// SelectServerCache 在serverLists配置中查询是否存在缓存
func SelectServerCache(ip string, conf *MainConfig) (*ServerListConfig, int, bool) {
	for index, server := range conf.ServerLists {
		if server.Ip == ip {
			return &server, index, true
		}
	}
	return nil, 0, false
}

// AddServerCache 新增缓存
func AddServerCache(server *ServerListConfig, conf *MainConfig) (writeRes bool) {
	conf.ServerLists = append(conf.ServerLists, *server)
	writeRes = utils.FileYamlMarshalAndWrite(configPath, conf)
	return
}

// DeleteServerCache 删除缓存
func DeleteServerCache(oldIndex int, conf *MainConfig) (writeRes bool) {
	conf.ServerLists = append(conf.ServerLists[0:oldIndex], conf.ServerLists[oldIndex+1:]...)
	writeRes = utils.FileYamlMarshalAndWrite(configPath, conf)
	return
}

// GenerateCombination 生成所有端口、用户、密码组合的对象
func GenerateCombination(ip string, conf *MainConfig) (combinations chan []interface{}) {
	ips := []interface{}{ip}
	ports := utils.InterfaceSlice(conf.Main.Ports)
	users := utils.InterfaceSlice(conf.Main.Users)
	passwords := utils.InterfaceSlice(conf.Main.Passwords)
	// 生成组合 参数顺序不可变
	combinations = cartesian.Iter(ips, ports, users, passwords)
	return
}
