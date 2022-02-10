package main

import (
	"flag"
	"fmt"
	"github.com/schwarmco/go-cartesian-product"
	log "github.com/sirupsen/logrus"
	"os"
	"tryssh/config"
	"tryssh/launcher"
	"tryssh/utils"
)

const (
	VERSION = "1.0.3"
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
func generateCombination(ip string, conf *config.MainConfig) (combinations chan []interface{}) {
	ips := []interface{}{ip}
	ports := utils.InterfaceSlice(conf.Main.Ports)
	users := utils.InterfaceSlice(conf.Main.Users)
	passwords := utils.InterfaceSlice(conf.Main.Passwords)
	// 生成组合 参数顺序不可变
	combinations = cartesian.Iter(ips, ports, users, passwords)
	return
}

// tryLogin 程序入口
func tryLogin(targetIp string, configuration *config.MainConfig) {
	// 搜索缓存
	targetServer, cacheIndex, isFound := config.SelectServerCache(targetIp, configuration)

	// 执行登陆逻辑
	if isFound {
		log.Infof("The cache for %s is found, which will be used to try.\n", targetIp)
		tryLoginWithCache(targetIp, configuration, targetServer, cacheIndex)
	} else {
		log.Warnf("The cache for %s could not be found. The password will be guessed.\n\n", targetIp)
		tryLoginWithoutCache(targetIp, configuration)
	}
	log.Fatalln("There is no password combination that can log in successfully\n")
}

// tryLoginWithCache 尝试用缓存连接，判断缓存连接是否成功，不成功则重新猜密码
func tryLoginWithCache(targetIp string, configuration *config.MainConfig,
	targetServer *config.ServerListConfig, cacheIndex int) {
	lan := &launcher.SshLauncher{
		SshTarget: *config.GetTargetFromConfig(targetServer),
	}
	if lan.Launch() {
		os.Exit(0)
	} else {
		log.Errorln("Failed to log in with cached information. The password will be guessed again\n\n")
		if config.DeleteServerCache(cacheIndex, configuration) {
			log.Infoln("Delete server cache successful.\n")
			tryLoginWithoutCache(targetIp, configuration)
		}
	}
}

// tryLoginWithoutCache 猜密码登陆
func tryLoginWithoutCache(targetIp string, configuration *config.MainConfig) {
	// 获取连接器对象
	combinations := generateCombination(targetIp, configuration)
	launchers := launcher.NewSshLaunchersByCombinations(combinations)
	for _, lan := range launchers {
		if err := lan.TryToConnect(); err == nil {
			log.Infoln("Login succeeded. The cache will be added.\n")
			if config.AddServerCache(config.GetConfigFromTarget(&lan.SshTarget), configuration) {
				log.Infoln("Cache updated.\n\n")
				lan.Launch()
			} else {
				log.Errorln("Cache update failed.\n\n")
			}
			os.Exit(0)
		}
	}
}

// flagsParse 解析命令行参数
func flagsParse() []string {
	// 参数定义
	ver := flag.Bool("v", false, "Show Version")
	flag.Usage = func() {
		fmt.Println("Usage: tryssh [IpAddress]")
		flag.PrintDefaults()
	}
	flag.Parse()

	// 命名参数处理
	if *ver {
		log.Infof("Tryssh Version: %s", VERSION)
	}

	// 未命名参数处理
	args := flag.Args()
	if len(args) != 1 {
		log.Errorln("Wrong number of parameters. There can only be one\n\n")
		flag.Usage()
		os.Exit(1)
	}
	return args
}
