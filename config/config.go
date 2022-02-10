package config

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"tryssh/target"
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

// GetTargetFromConfig 通过ServerListConfig获取SshTarget
func GetTargetFromConfig(conf *ServerListConfig) *target.SshTarget {
	return &target.SshTarget{
		Ip:       conf.Ip,
		Port:     conf.Port,
		User:     conf.User,
		Password: conf.Password,
	}
}

// GetConfigFromTarget 通过SshTarget获取ServerListConfig
func GetConfigFromTarget(tgt *target.SshTarget) *ServerListConfig {
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
	log.Infoln("Generating configuration file.\n")
	_ = utils.FileYamlMarshalAndWrite(configPath, &MainConfig{})
	log.Infoln("Generating configuration file successful.\n")
	log.Warnln("Main setting is empty. " +
		"You need to configure the main configuration before running again.\n")
	os.Exit(0)
}

// LoadConfig 加载配置文件
func LoadConfig() (c *MainConfig) {
	c = new(MainConfig)

	if checkFileIsExist(configPath) {
		conf, err := ioutil.ReadFile(configPath)
		if err != nil {
			log.Fatalln("Configuration file load failed: ", err)
		}
		unmarshalErr := yaml.Unmarshal(conf, c)
		if unmarshalErr != nil {
			log.Fatalln("Configuration file parsing failed: ", unmarshalErr)
		} else {
			if len(c.Main.Ports) == 0 || len(c.Main.Users) == 0 || len(c.Main.Passwords) == 0 {
				log.Warnln("Main setting is empty. " +
					"You need to configure the main configuration before running again.\n")
				os.Exit(0)
			}
		}
	} else {
		log.Infoln("Configuration file cannot be found, it will be generated automatically.\n")
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
