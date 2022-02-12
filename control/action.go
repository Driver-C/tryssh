package control

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"tryssh/cmd"
	"tryssh/config"
	"tryssh/launcher"
)

// Action 行动对象，一个对象代表一个操作流程
type Action struct {
	Flags         cmd.FlagsParse
	Configuration *config.MainConfig
}

func (a *Action) getTargetIp() (targetIp string) {
	if a.Flags.CopyFile {
		for _, v := range a.Flags.Args {
			if strings.Contains(v, ":") {
				targetIp = strings.Split(v, ":")[0]
			}
		}
	} else {
		targetIp = strings.Split(a.Flags.Args[0], ":")[0]
	}
	return
}

// TryLogin 逻辑入口
func (a *Action) TryLogin() {
	targetIp := a.getTargetIp()
	// 搜索缓存
	targetServer, cacheIndex, isFound := config.SelectServerCache(targetIp, a.Configuration)

	// 执行登陆逻辑
	if isFound {
		log.Infof("The cache for %s is found, which will be used to try.\n", targetIp)
		a.tryLoginWithCache(targetIp, targetServer, cacheIndex)
	} else {
		log.Warnf("The cache for %s could not be found. The password will be guessed.\n\n", targetIp)
		a.tryWithoutCache(targetIp)
	}
	log.Fatalln("There is no password combination that can log in successfully\n")
}

// tryWithoutCache 根据命令行参数执行猜密码登陆/传文件
func (a *Action) tryWithoutCache(targetIp string) {
	if a.Flags.CopyFile {
		a.tryCopyFileWithoutCache(targetIp)
	} else {
		a.tryLoginWithoutCache(targetIp)
	}
}

// tryLoginWithCache 尝试用缓存连接，判断缓存连接是否成功，不成功则重新猜密码
func (a *Action) tryLoginWithCache(targetIp string, targetServer *config.ServerListConfig, cacheIndex int) {
	lan := a.getConnectorWitchCache(targetServer)
	if lan.Launch() {
		os.Exit(0)
	} else {
		log.Errorln("Failed to log in with cached information. The password will be guessed again\n\n")
		if config.DeleteServerCache(cacheIndex, a.Configuration) {
			log.Infoln("Delete server cache successful.\n")
			a.tryWithoutCache(targetIp)
		}
	}
}

// tryLoginWithoutCache 猜密码登陆
func (a *Action) tryLoginWithoutCache(targetIp string) {
	// 获取连接器对象
	combinations := config.GenerateCombination(targetIp, a.Configuration)
	// SshLaunchers
	launchers := launcher.NewSshLaunchersByCombinations(combinations)
	for _, lan := range launchers {
		if err := lan.TryToConnect(); err == nil {
			log.Infoln("Login succeeded. The cache will be added.\n")
			if config.AddServerCache(config.GetConfigFromSshConnector(&lan.SshConnector), a.Configuration) {
				log.Infoln("Cache updated.\n\n")
				lan.Launch()
			} else {
				log.Errorln("Cache update failed.\n\n")
			}
			os.Exit(0)
		}
	}
}

// tryCopyFileWithoutCache 猜密码传输文件
func (a *Action) tryCopyFileWithoutCache(targetIp string) {
	// 获取连接器对象
	combinations := config.GenerateCombination(targetIp, a.Configuration)
	// ScpLaunchers
	launchers := launcher.NewScpLaunchersByCombinations(combinations, a.Flags.Args[0], a.Flags.Args[1])
	for _, lan := range launchers {
		if err := lan.TryToConnect(); err == nil {
			log.Infoln("Login succeeded. The cache will be added.\n")
			if config.AddServerCache(config.GetConfigFromSshConnector(&lan.SshConnector), a.Configuration) {
				log.Infoln("Cache updated.\n\n")
				lan.Launch()
			} else {
				log.Errorln("Cache update failed.\n\n")
			}
			os.Exit(0)
		}
	}
}

// getConnectorWitchCache 利用ServerListConfig获取单个launcher对象
func (a *Action) getConnectorWitchCache(targetServer *config.ServerListConfig) launcher.Connector {
	if a.Flags.CopyFile {
		return &launcher.ScpLauncher{
			SshConnector: *config.GetSshConnectorFromConfig(targetServer),
			Src:          a.Flags.Args[0],
			Dest:         a.Flags.Args[1],
		}
	} else {
		return &launcher.SshLauncher{
			SshConnector: *config.GetSshConnectorFromConfig(targetServer),
		}
	}
}
