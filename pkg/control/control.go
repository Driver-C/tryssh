package control

import (
	"context"
	"github.com/Driver-C/tryssh/pkg/launcher"
	"github.com/cheggaaa/pb/v3"
	"sync"
	"time"
)

const (
	TypeUsers                 = "users"
	TypePorts                 = "ports"
	TypePasswords             = "passwords"
	TypeCaches                = "caches"
	TypeKeys                  = "keys"
	sshClientTimeoutWhenLogin = 5 * time.Second
)

func ConcurrencyTryToConnect(concurrency int, connectors []launcher.Connector) []launcher.Connector {
	hitConnectors := make([]launcher.Connector, 0)
	mutex := new(sync.Mutex)
	bar := pb.StartNew(len(connectors))
	bar.Set("prefix", "Attempting:")
	// If the number of connectors is less than the set concurrency, change the concurrency to the number of connectors
	if concurrency > len(connectors) {
		concurrency = len(connectors)
	}
	connectorsChan := make(chan launcher.Connector)
	ctx, cancelFunc := context.WithCancel(context.Background())
	// Producer
	go func(ctx context.Context, connectorsChan chan<- launcher.Connector, connectors []launcher.Connector) {
		for _, connector := range connectors {
			select {
			case <-ctx.Done():
				break
			default:
				connectorsChan <- connector
			}
		}
		close(connectorsChan)
	}(ctx, connectorsChan, connectors)
	// Consumer
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(ctx context.Context, cancelFunc context.CancelFunc,
			connectorsChan <-chan launcher.Connector, cwg *sync.WaitGroup, mutex *sync.Mutex) {
			defer cwg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case connector, ok := <-connectorsChan:
					if !ok {
						return
					}
					if err := connector.TryToConnect(); err == nil {
						mutex.Lock()
						hitConnectors = append(hitConnectors, connector)
						mutex.Unlock()
						bar.Finish()
						cancelFunc()
					}
					bar.Increment()
				}
			}
		}(ctx, cancelFunc, connectorsChan, &wg, mutex)
	}
	wg.Wait()
	bar.Finish()
	cancelFunc()
	return hitConnectors
}
