package control

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Driver-C/tryssh/pkg/launcher"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

var errDummy = errors.New("connection failed")

type mockConnector struct {
	tryErr   error
	launchOk bool
	tryCalls int32
}

func (m *mockConnector) Launch() bool {
	return m.launchOk
}

func (m *mockConnector) CreateConnection() (*ssh.Client, error) {
	return nil, nil
}

func (m *mockConnector) CloseConnection(_ *ssh.Client) {}

func (m *mockConnector) TryToConnect() error {
	atomic.AddInt32(&m.tryCalls, 1)
	return m.tryErr
}

func TestConcurrencyTryToConnect_AllSucceed(t *testing.T) {
	connectors := make([]launcher.Connector, 5)
	for i := range connectors {
		connectors[i] = &mockConnector{tryErr: nil, launchOk: true}
	}

	hit := ConcurrencyTryToConnect(3, connectors)
	// At least one should succeed; due to concurrency, more may be added before cancel
	assert.GreaterOrEqual(t, len(hit), 1)
}

func TestConcurrencyTryToConnect_AllFail(t *testing.T) {
	connectors := make([]launcher.Connector, 5)
	for i := range connectors {
		connectors[i] = &mockConnector{tryErr: errDummy, launchOk: false}
	}

	hit := ConcurrencyTryToConnect(3, connectors)
	assert.Equal(t, 0, len(hit))

	// All connectors should have been tried
	for _, c := range connectors {
		mc := c.(*mockConnector)
		assert.Equal(t, int32(1), atomic.LoadInt32(&mc.tryCalls))
	}
}

func TestConcurrencyTryToConnect_MixedResults(t *testing.T) {
	connectors := make([]launcher.Connector, 6)
	for i := range connectors {
		connectors[i] = &mockConnector{tryErr: errDummy, launchOk: false}
	}
	connectors[5] = &mockConnector{tryErr: nil, launchOk: true}

	hit := ConcurrencyTryToConnect(3, connectors)
	assert.GreaterOrEqual(t, len(hit), 1)
}

func TestConcurrencyTryToConnect_SingleConnector(t *testing.T) {
	connectors := []launcher.Connector{
		&mockConnector{tryErr: nil, launchOk: true},
	}

	hit := ConcurrencyTryToConnect(5, connectors)
	assert.Equal(t, 1, len(hit))
}

func TestConcurrencyTryToConnect_EmptyConnectors(t *testing.T) {
	connectors := make([]launcher.Connector, 0)

	hit := ConcurrencyTryToConnect(3, connectors)
	assert.Equal(t, 0, len(hit))
}

func TestConcurrencyTryToConnect_ConcurrencyLimit(t *testing.T) {
	connectors := make([]launcher.Connector, 2)
	for i := range connectors {
		connectors[i] = &mockConnector{tryErr: nil, launchOk: true}
	}

	hit := ConcurrencyTryToConnect(100, connectors)
	assert.GreaterOrEqual(t, len(hit), 1)

	totalCalls := int32(0)
	for _, c := range connectors {
		mc := c.(*mockConnector)
		totalCalls += atomic.LoadInt32(&mc.tryCalls)
	}
	assert.True(t, totalCalls >= 1)
}

func TestConcurrencyTryToConnect_ConcurrencyOne(t *testing.T) {
	connectors := make([]launcher.Connector, 4)
	for i := range connectors {
		connectors[i] = &mockConnector{tryErr: errDummy, launchOk: false}
	}

	hit := ConcurrencyTryToConnect(1, connectors)
	assert.Equal(t, 0, len(hit))

	// All should have been tried sequentially
	for _, c := range connectors {
		mc := c.(*mockConnector)
		assert.Equal(t, int32(1), atomic.LoadInt32(&mc.tryCalls))
	}
}

func TestConcurrencyTryToConnect_FirstWins(t *testing.T) {
	connectors := make([]launcher.Connector, 10)
	connectors[0] = &mockConnector{tryErr: nil, launchOk: true}
	for i := 1; i < 10; i++ {
		connectors[i] = &mockConnector{tryErr: nil, launchOk: true}
	}

	hit := ConcurrencyTryToConnect(3, connectors)
	// At least one succeeds; multiple may succeed before cancel propagates
	assert.GreaterOrEqual(t, len(hit), 1)
}

func TestConcurrencyTryToConnect_LargeSet(t *testing.T) {
	connectors := make([]launcher.Connector, 100)
	for i := range connectors {
		connectors[i] = &mockConnector{tryErr: errDummy, launchOk: false}
	}
	// The 50th succeeds
	connectors[49] = &mockConnector{tryErr: nil, launchOk: true}

	start := time.Now()
	hit := ConcurrencyTryToConnect(10, connectors)
	elapsed := time.Since(start)

	assert.GreaterOrEqual(t, len(hit), 1)
	// Should complete quickly due to concurrency
	assert.True(t, elapsed < 5*time.Second, "should finish quickly")
}
