package control

import (
	"bufio"
	"fmt"
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/launcher"
	"github.com/Driver-C/tryssh/pkg/utils"
	"os"
	"strings"
	"sync"
	"time"
)

type PruneController struct {
	configuration *config.MainConfig
	auto          bool
	sshTimeout    time.Duration
	concurrency   int
}

func (pc *PruneController) PruneCaches() {
	newServerList := make([]config.ServerListConfig, 0)
	if pc.auto {
		newServerList = pc.concurrencyDeleteCache()
	} else {
		for _, server := range pc.configuration.ServerLists {
			lan := &launcher.SshLauncher{SshConnector: *launcher.GetSshConnectorFromConfig(&server)}
			// Set timeout
			lan.SshTimeout = pc.sshTimeout
			// Determine if connection is possible
			if err := lan.TryToConnect(); err != nil {
				if !pc.interactiveDeleteCache(server) {
					newServerList = append(newServerList, server)
				}
			} else {
				utils.Logger.Infof("Cache %v is still available.", server)
				newServerList = append(newServerList, server)
			}
		}
	}
	pc.configuration.ServerLists = newServerList
	if config.UpdateConfig(pc.configuration) {
		utils.Logger.Infoln("Update config successful.")
	} else {
		utils.Logger.Errorln("Update config failed.")
	}
}

func (pc *PruneController) interactiveDeleteCache(server config.ServerListConfig) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Are you sure you want to delete this cache? "+
			"Please enter \"yes\" to confirm deletion, or \"no\" to cancel. %s\n"+
			"(yes/no): ", server)
		stdin, _ := reader.ReadString('\n')
		// Delete space
		stdin = strings.TrimSpace(stdin)
		switch stdin {
		case "yes":
			utils.Logger.Infof("The cache %v has been marked for deletion.", server)
			return true
		case "no":
			utils.Logger.Infof("Cache %v skipped.", server)
			return false
		default:
			utils.Logger.Errorln("Input error:", stdin)
		}
	}
}

func (pc *PruneController) concurrencyDeleteCache() []config.ServerListConfig {
	newServerList := make([]config.ServerListConfig, 0)
	serversChan := make(chan *config.ServerListConfig)
	var mutex sync.Mutex
	var wg sync.WaitGroup

	go func(serversChan chan<- *config.ServerListConfig) {
		for _, server := range pc.configuration.ServerLists {
			newServer := server
			serversChan <- &newServer
		}
		close(serversChan)
	}(serversChan)

	for i := 0; i < pc.concurrency; i++ {
		wg.Add(1)
		go func(serversChan <-chan *config.ServerListConfig, wg *sync.WaitGroup) {
			defer wg.Done()
			for {
				serverP, ok := <-serversChan
				if !ok {
					break
				}
				lan := &launcher.SshLauncher{SshConnector: *launcher.GetSshConnectorFromConfig(serverP)}
				lan.SshTimeout = pc.sshTimeout
				if err := lan.TryToConnect(); err == nil {
					utils.Logger.Infof("Cache %v is still available.", *serverP)
					mutex.Lock()
					newServerList = append(newServerList, *serverP)
					mutex.Unlock()
				} else {
					utils.Logger.Infof("The cache %v has been marked for deletion.", *serverP)
				}
			}
		}(serversChan, &wg)
	}
	wg.Wait()
	return newServerList
}

func NewPruneController(configuration *config.MainConfig, auto bool, timeout time.Duration,
	concurrency int) *PruneController {
	return &PruneController{
		configuration: configuration,
		auto:          auto,
		sshTimeout:    timeout,
		concurrency:   concurrency,
	}
}
