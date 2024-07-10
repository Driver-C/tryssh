package control

import (
	"fmt"
	"github.com/Driver-C/tryssh/pkg/config"
)

type GetController struct {
	getType       string
	getContent    string
	configuration *config.MainConfig
}

func (gc GetController) ExecuteGet() {
	switch gc.getType {
	case TypeUsers:
		fmt.Println("INDEX	USER")
		gc.searchAndPrint(gc.configuration.Main.Users)
	case TypePorts:
		fmt.Println("INDEX	PORT")
		gc.searchAndPrint(gc.configuration.Main.Ports)
	case TypePasswords:
		fmt.Println("INDEX	PASSWORD")
		gc.searchAndPrint(gc.configuration.Main.Passwords)
	case TypeKeys:
		fmt.Println("INDEX	KEY")
		gc.searchAndPrint(gc.configuration.Main.Keys)
	case TypeCaches:
		// gc.getContent is ipAddress
		fmt.Println("INDEX	CACHE")
		if gc.getContent != "" {
			for index, server := range gc.configuration.ServerLists {
				if server.Ip == gc.getContent {
					fmt.Printf("%d	%s\n", index, server)
					break
				}
			}
		} else {
			for index, server := range gc.configuration.ServerLists {
				fmt.Printf("%d	%s\n", index, server)
			}
		}
	}
}

func (gc GetController) searchAndPrint(contents []string) {
	if gc.getContent != "" {
		for index, content := range contents {
			if content == gc.getContent {
				fmt.Printf("%d	%s\n", index, content)
				break
			}
		}
	} else {
		for index, content := range contents {
			fmt.Printf("%d	%s\n", index, content)
		}
	}
}

func NewGetController(getType string, getContent string,
	configuration *config.MainConfig) *GetController {
	return &GetController{
		getType:       getType,
		getContent:    getContent,
		configuration: configuration,
	}
}
