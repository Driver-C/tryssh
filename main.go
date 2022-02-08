package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
	"tryssh/config"
	"tryssh/launcher"
	"tryssh/target"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln(err)
		}
	}()

	// 加载命令行参数
	args := flagsParse()
	targetIp := args[0]

	// 加载配置文件
	configuration := config.LoadConfig()

	// Start
	tryLogin(targetIp, configuration)
}

// generateCombination 生成所有端口、用户、密码组合的对象
func generateCombination(ip string, conf *config.MainConfig) (lan []*launcher.SshLauncher) {
	for _, port := range conf.Main.Ports {
		for _, user := range conf.Main.Users {
			for _, password := range conf.Main.Passwords {
				lan = append(lan, &launcher.SshLauncher{
					SshTarget: target.SshTarget{Ip: ip, Port: port, User: user, Password: password},
				})
			}
		}
	}
	return
}

func tryLogin(targetIp string, configuration *config.MainConfig) {
	// 获取连接器对象
	launchers := generateCombination(targetIp, configuration)
	// 搜索缓存
	targetServer, cacheIndex, isFound := config.SelectServerCache(targetIp, configuration)

	// 执行登陆逻辑
	if isFound {
		lan := &launcher.SshLauncher{
			SshTarget: *config.GetTargetFromConfig(targetServer),
		}
		// 判断缓存连接是否成功，不成功则重新尝试连接
		log.Infof("The cache for %s is found, which will be used to try.\n", targetIp)
		if lan.Launch() {
			os.Exit(0)
		} else {
			log.Errorln("Failed to log in with cached information. The password will be guessed again\n\n")
			for _, lan := range launchers {
				if lan.Launch() {
					log.Infoln("Login succeeded. The cache will be updated.\n")
					if config.UpdateServerCache(
						cacheIndex, config.GetConfigFromTarget(&lan.SshTarget), configuration) {
						log.Infoln("Cache updated.\n")
					} else {
						log.Errorln("Cache update failed.\n\n")
					}
					os.Exit(0)
				}
			}
			log.Fatalln("There is no password combination that can log in successfully\n")
		}
	} else {
		log.Warnf("The cache for %s could not be found. The password will be guessed.\n\n", targetIp)
		for _, lan := range launchers {
			if lan.Launch() {
				log.Infoln("Login succeeded. The cache will be added.\n")
				if config.AddServerCache(config.GetConfigFromTarget(&lan.SshTarget), configuration) {
					log.Infoln("Cache updated.\n")
				} else {
					log.Errorln("Cache update failed.\n\n")
				}
				os.Exit(0)
			}
		}
		log.Fatalln("There is no password combination that can log in successfully\n")
	}
}

// flagsParse 解析命令行参数
func flagsParse() []string {
	flag.Usage = func() {
		log.Infoln("Usage: tryssh [IpAddress]\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		log.Errorln("Wrong number of parameters. There can only be one\n\n")
		flag.Usage()
		os.Exit(1)
	}
	return args
}
