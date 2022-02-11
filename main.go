package main

import (
	log "github.com/sirupsen/logrus"
	"tryssh/cmd"
	"tryssh/config"
	"tryssh/control"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln(err)
		}
	}()

	// 生成命令参数结构
	flags := cmd.NewFlagsParse()

	// 加载配置文件
	configuration := config.LoadConfig()

	// Start
	action := control.Action{
		Flags:         flags,
		Configuration: configuration,
	}
	action.TryLogin()
}
