package main

import (
	"github.com/Driver-C/tryssh/cmd"
	"github.com/Driver-C/tryssh/pkg/utils"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			utils.Logger.Errorln(err)
		}
	}()

	rootCmd := cmd.NewTrysshCommand()
	if err := rootCmd.Execute(); err != nil {
		utils.Logger.Errorln(err)
	}
}
