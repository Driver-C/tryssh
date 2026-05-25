package launcher

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

var errConnectionRefused = errors.New("connection refused")

// ---------------------------------------------------------------------------
// Tests for NewSSHLaunchersByCombinations
// ---------------------------------------------------------------------------

func TestNewSSHLaunchersByCombinations(t *testing.T) {
	combinations := make(chan []interface{}, 3)
	combinations <- []interface{}{"192.168.1.1", "22", "user1", "pass1", "/key1"}
	combinations <- []interface{}{"10.0.0.1", "2222", "user2", "pass2", "/key2"}
	combinations <- []interface{}{"172.16.0.1", "8022", "user3", "pass3", ""}
	close(combinations)

	timeout := 10 * time.Second
	launchers := NewSSHLaunchersByCombinations(combinations, timeout)

	assert.Len(t, launchers, 3)

	assert.Equal(t, "192.168.1.1", launchers[0].IP)
	assert.Equal(t, "22", launchers[0].Port)
	assert.Equal(t, "user1", launchers[0].User)
	assert.Equal(t, "pass1", launchers[0].Password)
	assert.Equal(t, "/key1", launchers[0].Key)
	assert.Equal(t, timeout, launchers[0].SSHTimeout)

	assert.Equal(t, "10.0.0.1", launchers[1].IP)
	assert.Equal(t, "2222", launchers[1].Port)

	assert.Equal(t, "172.16.0.1", launchers[2].IP)
	assert.Equal(t, "", launchers[2].Key)
}

func TestNewSSHLaunchersByCombinations_Empty(t *testing.T) {
	combinations := make(chan []interface{})
	close(combinations)

	launchers := NewSSHLaunchersByCombinations(combinations, 5*time.Second)
	assert.Empty(t, launchers)
}

// ---------------------------------------------------------------------------
// Tests for SSHLauncher.Launch with mock dialer
// ---------------------------------------------------------------------------

func TestSSHLauncher_Launch_ConnectionFails(t *testing.T) {
	launcher := &SSHLauncher{
		SSHConnector: SSHConnector{
			IP:         "127.0.0.1",
			Port:       "22",
			User:       "testuser",
			Password:   "testpass",
			SSHTimeout: 5 * time.Second,
			Dialer:     &mockSSHDialer{client: nil, err: errConnectionRefused},
			KnownHosts: tempKnownHosts(t, nil),
		},
	}

	result := launcher.Launch()
	assert.False(t, result)
}

// ---------------------------------------------------------------------------
// Tests for SSHLauncher.Launch with real server (dialServer success)
// ---------------------------------------------------------------------------

func TestSSHLauncher_Launch_WithRealServer(t *testing.T) {
	// This test verifies the dialServer success path.  Since createTerminal
	// requires os.Stdin to be a terminal, it will error out, but the connection
	// succeeds and the function returns true.
	//
	// We use a mock dialer that returns a real *ssh.Client connected to our
	// test server.

	ts, _ := newTestServer(t, "testuser", "testpass")
	defer ts.close()

	// Create a real SSH client by connecting to the test server
	addr := ts.addr()
	client, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User:            "testuser",
		Auth:            []ssh.AuthMethod{ssh.Password("testpass")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	})
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// Use a custom dialer that returns this pre-connected client
	host, port, _ := net.SplitHostPort(addr)
	launcher := &SSHLauncher{
		SSHConnector: SSHConnector{
			IP:         host,
			Port:       port,
			User:       "testuser",
			Password:   "testpass",
			SSHTimeout: 5 * time.Second,
			Dialer:     &mockSSHDialer{client: client, err: nil},
			KnownHosts: tempKnownHosts(t, nil),
		},
	}

	// Launch will call dialServer -> CreateConnection -> success -> createTerminal
	// createTerminal will fail because os.Stdin is not a terminal, but
	// dialServer returns true before that error matters (it logs but returns true).
	result := launcher.Launch()
	assert.True(t, result)
}
