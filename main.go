package main

import (
	"tryssh/cmd"
	"tryssh/utils"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			utils.Logger.Errorln(err)
		}
	}()

	rootCmd := cmd.NewTrysshCommand()
	if err := rootCmd.Execute(); err != nil {
		utils.Logger.Fatalln(err)
	}
}
