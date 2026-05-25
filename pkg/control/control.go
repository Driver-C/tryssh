package control

import (
	"context"
	"github.com/Driver-C/tryssh/pkg/launcher"
	"github.com/cheggaaa/pb/v3"
	"sync"
	"time"
)

// Resource type constants used to identify the kind of configuration entry.
const (
	TypeUsers                 = "users"
	TypePorts                 = "ports"
	TypePasswords             = "passwords"
	TypeCaches                = "caches"
	TypeKeys                  = "keys"
	sshClientTimeoutWhenLogin = 5 * time.Second
)

// ConcurrencyTryToConnect attempts to connect using the given connectors concurrently,
// returning the ones that succeed.
func ConcurrencyTryToConnect(concurrency int, connectors []launcher.Connector) []launcher.Connector {
	if len(connectors) == 0 {
		return nil
	}
	if concurrency < 1 {
		concurrency = 1
	}
	if concurrency > len(connectors) {
		concurrency = len(connectors)
	}

	hitConnectors := make([]launcher.Connector, 0)
	bar := pb.StartNew(len(connectors))
	bar.Set("prefix", "Attempting:")

	connectorsChan := make(chan launcher.Connector)
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// Producer
	go func() {
		defer close(connectorsChan)
		for _, connector := range connectors {
			select {
			case <-ctx.Done():
				return
			case connectorsChan <- connector:
			}
		}
	}()

	// Consumer
	var mu sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case connector, ok := <-connectorsChan:
					if !ok {
						return
					}
					if err := connector.TryToConnect(); err == nil {
						mu.Lock()
						hitConnectors = append(hitConnectors, connector)
						mu.Unlock()
						cancelFunc()
					}
					bar.Increment()
				}
			}
		}()
	}
	wg.Wait()
	bar.Finish()
	return hitConnectors
}
