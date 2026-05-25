// Package main is the entry point for the tryssh CLI application.
package main

import (
	"os"

	"github.com/Driver-C/tryssh/cmd"
	"github.com/Driver-C/tryssh/pkg/utils"
)

func main() {
	os.Exit(run())
}

func run() int {
	defer func() {
		if err := recover(); err != nil {
			utils.Errorln(err)
		}
	}()

	rootCmd := cmd.NewTrysshCommand()
	if err := rootCmd.Execute(); err != nil {
		utils.Errorln(err)
		return 1
	}
	return 0
}
