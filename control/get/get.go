package get

import (
	"fmt"
	"github.com/Driver-C/tryssh/config"
)

const (
	typeUsers     = "users"
	typePorts     = "ports"
	typePasswords = "passwords"
	typeCaches    = "caches"
)

type Controller struct {
	getType       string
	getContent    string
	configuration *config.MainConfig
}

func (gc Controller) ExecuteGet() {
	switch gc.getType {
	case typeUsers:
		fmt.Println("INDEX	USER")
		gc.searchAndPrint(gc.configuration.Main.Users)
	case typePorts:
		fmt.Println("INDEX	PORT")
		gc.searchAndPrint(gc.configuration.Main.Ports)
	case typePasswords:
		fmt.Println("INDEX	PASSWORD")
		gc.searchAndPrint(gc.configuration.Main.Passwords)
	case typeCaches:
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

func (gc Controller) searchAndPrint(contents []string) {
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
	configuration *config.MainConfig) *Controller {
	return &Controller{
		getType:       getType,
		getContent:    getContent,
		configuration: configuration,
	}
}
