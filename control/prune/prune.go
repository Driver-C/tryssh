package prune

import (
	"bufio"
	"fmt"
	"github.com/Driver-C/tryssh/config"
	"github.com/Driver-C/tryssh/launcher"
	"github.com/Driver-C/tryssh/launcher/ssh"
	"github.com/Driver-C/tryssh/utils"
	"os"
	"strings"
	"sync"
	"time"
)

type Controller struct {
	configuration *config.MainConfig
	auto          bool
	sshTimeout    time.Duration
	concurrency   int
}

func (pc *Controller) PruneCaches() {
	newServerList := make([]config.ServerListConfig, 0)
	if pc.auto {
		newServerList = pc.concurrencyDeleteCache()
	} else {
		for _, server := range pc.configuration.ServerLists {
			lan := &ssh.Launcher{SshConnector: *launcher.GetSshConnectorFromConfig(&server)}
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

func (pc *Controller) interactiveDeleteCache(server config.ServerListConfig) bool {
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

func (pc *Controller) concurrencyDeleteCache() []config.ServerListConfig {
	newServerList := make([]config.ServerListConfig, 0)
	launchersChan := make(chan *ssh.Launcher)
	var mutex sync.Mutex
	var wg sync.WaitGroup

	go func(launchersChan chan<- *ssh.Launcher) {
		for _, server := range pc.configuration.ServerLists {
			lan := &ssh.Launcher{SshConnector: *launcher.GetSshConnectorFromConfig(&server)}
			launchersChan <- lan
		}
		close(launchersChan)
	}(launchersChan)

	for i := 0; i < pc.concurrency; i++ {
		wg.Add(1)
		go func(launchersChan <-chan *ssh.Launcher, wg *sync.WaitGroup) {
			defer wg.Done()
			for {
				launcherP, ok := <-launchersChan
				if !ok {
					break
				}
				server := launcher.GetConfigFromSshConnector(&launcherP.SshConnector)
				launcherP.SshTimeout = pc.sshTimeout
				if err := launcherP.TryToConnect(); err == nil {
					utils.Logger.Infof("Cache %v is still available.", server)
					mutex.Lock()
					newServerList = append(newServerList, *server)
					mutex.Unlock()
				} else {
					utils.Logger.Infof("The cache %v has been marked for deletion.", server)
				}
			}
		}(launchersChan, &wg)
	}
	wg.Wait()
	return newServerList
}

func NewPruneController(configuration *config.MainConfig, auto bool, timeout time.Duration,
	concurrency int) *Controller {
	return &Controller{
		configuration: configuration,
		auto:          auto,
		sshTimeout:    timeout,
		concurrency:   concurrency,
	}
}
