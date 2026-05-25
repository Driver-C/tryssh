package control

import (
	"fmt"
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/utils"
)

// GetController handles retrieval and display of configuration entries.
type GetController struct {
	getType       string
	getContent    string
	configuration *config.MainConfig
}

// ExecuteGet displays the configured entry or all entries of the given type.
func (gc *GetController) ExecuteGet() {
	switch gc.getType {
	case TypeUsers:
		fmt.Println("INDEX\tUSER")
		gc.searchAndPrint(gc.configuration.Main.Users, false)
	case TypePorts:
		fmt.Println("INDEX\tPORT")
		gc.searchAndPrint(gc.configuration.Main.Ports, false)
	case TypePasswords:
		fmt.Println("INDEX\tPASSWORD")
		gc.searchAndPrint(gc.configuration.Main.Passwords, true)
	case TypeKeys:
		fmt.Println("INDEX\tKEY")
		gc.searchAndPrint(gc.configuration.Main.Keys, false)
	case TypeCaches:
		fmt.Println("INDEX\tCACHE")
		if gc.getContent != "" {
			for index, server := range gc.configuration.ServerLists {
				if server.IP == gc.getContent {
					fmt.Printf("%d\t%s\n", index, server)
				}
			}
		} else {
			for index, server := range gc.configuration.ServerLists {
				fmt.Printf("%d\t%s\n", index, server)
			}
		}
	}
}

func (gc *GetController) searchAndPrint(contents []string, maskValues bool) {
	maskFn := func(s string) string { return s }
	if maskValues {
		maskFn = utils.MaskSecret
	}
	if gc.getContent != "" {
		for index, content := range contents {
			if content == gc.getContent {
				fmt.Printf("%d\t%s\n", index, maskFn(content))
			}
		}
	} else {
		for index, content := range contents {
			fmt.Printf("%d\t%s\n", index, maskFn(content))
		}
	}
}

// NewGetController creates a new GetController for the specified type and content filter.
func NewGetController(getType string, getContent string,
	configuration *config.MainConfig) *GetController {
	return &GetController{
		getType:       getType,
		getContent:    getContent,
		configuration: configuration,
	}
}
