package launcher

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/utils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

// ---------------------------------------------------------------------------
// Mock SSHDialer
// ---------------------------------------------------------------------------

type mockSSHDialer struct {
	client *ssh.Client
	err    error
}

func (m *mockSSHDialer) Dial(_ string, _ string, _ *ssh.ClientConfig) (*ssh.Client, error) {
	return m.client, m.err
}

// ---------------------------------------------------------------------------
// In-process SSH test server
// ---------------------------------------------------------------------------

// testServer holds an in-process SSH server for testing.
type testServer struct {
	listener net.Listener
	config   *ssh.ServerConfig
	hostSigner ssh.Signer
}

// newTestServer creates and starts a local SSH server that accepts password auth
// with the given user/password. It returns the server, the host signer (for known_hosts),
// and the address it is listening on.
func newTestServer(t *testing.T, user, password string) (*testServer, ssh.Signer) {
	t.Helper()
	_, priv, err := generateEd25519KeyPair()
	assert.NoError(t, err)
	signer, err := ssh.NewSignerFromKey(priv)
	assert.NoError(t, err)

	cfg := &ssh.ServerConfig{
		PasswordCallback: func(conn ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if conn.User() == user && string(pass) == password {
				return nil, nil
			}
			return nil, errors.New("auth rejected")
		},
	}
	cfg.AddHostKey(signer)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)

	ts := &testServer{listener: listener, config: cfg, hostSigner: signer}

	go ts.serve()

	return ts, signer
}

func (ts *testServer) serve() {
	for {
		conn, err := ts.listener.Accept()
		if err != nil {
			return
		}
		go ts.handleConn(conn)
	}
}

func (ts *testServer) handleConn(conn net.Conn) {
	_, chans, reqs, err := ssh.NewServerConn(conn, ts.config)
	if err != nil {
		return
	}
	go ssh.DiscardRequests(reqs)

	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			continue
		}
		go func(in <-chan *ssh.Request) {
			for req := range in {
				switch req.Type {
				case "exec":
					req.Reply(true, nil)
					channel.Write([]byte("ok\n"))
					channel.SendRequest("exit-status", false, ssh.Marshal(struct{ Status uint32 }{0}))
					channel.Close()
					return
				case "shell":
					req.Reply(true, nil)
					// Keep channel open briefly then close
					go func() {
						io.WriteString(channel, "$ ")
						// Don't close immediately; keep open for the test
					}()
				case "pty-req":
					req.Reply(true, nil)
				default:
					req.Reply(false, nil)
				}
			}
		}(requests)
	}
}

func (ts *testServer) addr() string {
	return ts.listener.Addr().String()
}

func (ts *testServer) close() {
	ts.listener.Close()
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// tempKnownHosts creates a temporary known_hosts file with the given content lines.
func tempKnownHosts(t *testing.T, lines []string) (path string) {
	t.Helper()
	dir := t.TempDir()
	path = filepath.Join(dir, "known_hosts")
	var content string
	if len(lines) > 0 {
		content = strings.Join(lines, "\n") + "\n"
	}
	err := os.WriteFile(path, []byte(content), 0600)
	assert.NoError(t, err)
	return
}

// tempKeyFile writes a valid Ed25519 private key to a temp file and returns its path.
func tempKeyFile(t *testing.T) (keyPath string, signer ssh.Signer) {
	t.Helper()
	_, priv, err := generateEd25519KeyPair()
	assert.NoError(t, err)
	signer, err = ssh.NewSignerFromKey(priv)
	assert.NoError(t, err)

	dir := t.TempDir()
	keyPath = filepath.Join(dir, "id_ed25519")
	err = os.WriteFile(keyPath, marshalPrivateKey(priv), 0600)
	assert.NoError(t, err)
	return
}

// newTestSSHConnector creates a basic SSHConnector with sensible defaults for tests.
func newTestSSHConnector() *SSHConnector {
	return &SSHConnector{
		IP:         "127.0.0.1",
		Port:       "22",
		User:       "testuser",
		Password:   "testpass",
		SSHTimeout: 5 * time.Second,
	}
}

// newConnectorForServer creates a SSHConnector configured to connect to the test server.
func newConnectorForServer(ts *testServer, user, password string) *SSHConnector {
	addr := ts.addr()
	host, port, _ := net.SplitHostPort(addr)
	return &SSHConnector{
		IP:         host,
		Port:       port,
		User:       user,
		Password:   password,
		SSHTimeout: 5 * time.Second,
	}
}

// ---------------------------------------------------------------------------
// Tests for SSHConnector.LoadConfig
// ---------------------------------------------------------------------------

func TestLoadConfig_PasswordOnly(t *testing.T) {
	sc := newTestSSHConnector()
	sc.KnownHosts = tempKnownHosts(t, nil)

	cfg := sc.BuildSSHConfig()
	assert.NotNil(t, cfg)
	assert.Equal(t, "testuser", cfg.User)
	assert.Len(t, cfg.Auth, 1)
	assert.NotNil(t, cfg.HostKeyCallback)
	assert.Equal(t, 5*time.Second, cfg.Timeout)
}

func TestLoadConfig_WithKey(t *testing.T) {
	keyPath, _ := tempKeyFile(t)
	sc := newTestSSHConnector()
	sc.Key = keyPath
	sc.KnownHosts = tempKnownHosts(t, nil)

	cfg := sc.BuildSSHConfig()
	assert.NotNil(t, cfg)
	assert.Len(t, cfg.Auth, 2)
}

func TestLoadConfig_WithInvalidKeyPath(t *testing.T) {
	sc := newTestSSHConnector()
	sc.Key = "/nonexistent/key/path"
	sc.KnownHosts = tempKnownHosts(t, nil)

	cfg := sc.BuildSSHConfig()
	assert.NotNil(t, cfg)
	assert.Len(t, cfg.Auth, 1)
}

func TestLoadConfig_EmptyKey(t *testing.T) {
	sc := newTestSSHConnector()
	sc.Key = ""
	sc.KnownHosts = tempKnownHosts(t, nil)

	cfg := sc.BuildSSHConfig()
	assert.NotNil(t, cfg)
	assert.Len(t, cfg.Auth, 1)
}

func TestLoadConfig_InvalidKeyContent(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "bad_key")
	err := os.WriteFile(keyPath, []byte("not a valid key"), 0600)
	assert.NoError(t, err)

	sc := newTestSSHConnector()
	sc.Key = keyPath
	sc.KnownHosts = tempKnownHosts(t, nil)

	cfg := sc.BuildSSHConfig()
	assert.NotNil(t, cfg)
	assert.Len(t, cfg.Auth, 1)
}

// ---------------------------------------------------------------------------
// Tests for SSHConnector.CreateConnection
// ---------------------------------------------------------------------------

func TestCreateConnection_Success(t *testing.T) {
	sc := newTestSSHConnector()
	sc.KnownHosts = tempKnownHosts(t, nil)

	dialer := &mockSSHDialer{client: nil, err: nil}
	sc.Dialer = dialer

	client, err := sc.CreateConnection()
	assert.NoError(t, err)
	assert.Nil(t, client)
}

func TestCreateConnection_DialError(t *testing.T) {
	sc := newTestSSHConnector()
	sc.KnownHosts = tempKnownHosts(t, nil)

	dialErr := errors.New("connection refused")
	sc.Dialer = &mockSSHDialer{client: nil, err: dialErr}

	client, err := sc.CreateConnection()
	assert.Error(t, err)
	assert.Equal(t, dialErr, err)
	assert.Nil(t, client)
}

func TestCreateConnection_SSHKeyKeywordError(t *testing.T) {
	sc := newTestSSHConnector()
	sc.KnownHosts = tempKnownHosts(t, nil)

	dialErr := fmt.Errorf("some error containing SSH-KEY in message")
	sc.Dialer = &mockSSHDialer{client: nil, err: dialErr}

	client, err := sc.CreateConnection()
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestCreateConnection_DefaultDialer(t *testing.T) {
	sc := newTestSSHConnector()
	sc.KnownHosts = tempKnownHosts(t, nil)
	assert.Nil(t, sc.Dialer)

	client, err := sc.CreateConnection()
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestCreateConnection_WithRealServer(t *testing.T) {
	ts, _ := newTestServer(t, "testuser", "testpass")
	defer ts.close()

	sc := newConnectorForServer(ts, "testuser", "testpass")
	sc.KnownHosts = tempKnownHosts(t, nil)

	client, err := sc.CreateConnection()
	assert.NoError(t, err)
	assert.NotNil(t, client)
	if client != nil {
		client.Close()
	}
}

func TestCreateConnection_WithRealServerBadCreds(t *testing.T) {
	ts, _ := newTestServer(t, "testuser", "testpass")
	defer ts.close()

	sc := newConnectorForServer(ts, "testuser", "wrongpass")
	sc.KnownHosts = tempKnownHosts(t, nil)

	client, err := sc.CreateConnection()
	assert.Error(t, err)
	assert.Nil(t, client)
}

// ---------------------------------------------------------------------------
// Tests for SSHConnector.CloseConnection
// ---------------------------------------------------------------------------

func TestCloseConnection_Success(t *testing.T) {
	ts, _ := newTestServer(t, "testuser", "testpass")
	defer ts.close()

	sc := newConnectorForServer(ts, "testuser", "testpass")
	sc.KnownHosts = tempKnownHosts(t, nil)

	client, err := sc.CreateConnection()
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// Close should succeed
	sc.CloseConnection(client)
}

// ---------------------------------------------------------------------------
// Tests for SSHConnector.TryToConnect
// ---------------------------------------------------------------------------

func TestTryToConnect_Success(t *testing.T) {
	ts, _ := newTestServer(t, "testuser", "testpass")
	defer ts.close()

	sc := newConnectorForServer(ts, "testuser", "testpass")
	sc.KnownHosts = tempKnownHosts(t, nil)

	err := sc.TryToConnect()
	assert.NoError(t, err)
}

func TestTryToConnect_Failure(t *testing.T) {
	sc := newTestSSHConnector()
	sc.KnownHosts = tempKnownHosts(t, nil)

	dialErr := errors.New("dial failed")
	sc.Dialer = &mockSSHDialer{client: nil, err: dialErr}

	err := sc.TryToConnect()
	assert.Error(t, err)
	assert.Equal(t, dialErr, err)
}

func TestTryToConnect_ServerAuthFails(t *testing.T) {
	ts, _ := newTestServer(t, "testuser", "testpass")
	defer ts.close()

	sc := newConnectorForServer(ts, "testuser", "wrongpass")
	sc.KnownHosts = tempKnownHosts(t, nil)

	err := sc.TryToConnect()
	assert.Error(t, err)
}

// ---------------------------------------------------------------------------
// Tests for GetSSHConnectorFromConfig
// ---------------------------------------------------------------------------

func TestGetSSHConnectorFromConfig(t *testing.T) {
	conf := &config.ServerListConfig{
		IP:       "192.168.1.1",
		Port:     "2222",
		User:     "admin",
		Password: "secret",
		Key:      "/path/to/key",
	}

	sc := GetSSHConnectorFromConfig(conf)
	assert.Equal(t, "192.168.1.1", sc.IP)
	assert.Equal(t, "2222", sc.Port)
	assert.Equal(t, "admin", sc.User)
	assert.Equal(t, "secret", sc.Password)
	assert.Equal(t, "/path/to/key", sc.Key)
}

// ---------------------------------------------------------------------------
// Tests for GetConfigFromSSHConnector
// ---------------------------------------------------------------------------

func TestGetConfigFromSSHConnector(t *testing.T) {
	sc := &SSHConnector{
		IP:       "10.0.0.1",
		Port:     "22",
		User:     "root",
		Password: "pass",
		Key:      "/home/user/.ssh/id_rsa",
	}

	conf := GetConfigFromSSHConnector(sc)
	assert.Equal(t, "10.0.0.1", conf.IP)
	assert.Equal(t, "22", conf.Port)
	assert.Equal(t, "root", conf.User)
	assert.Equal(t, "pass", conf.Password)
	assert.Equal(t, "/home/user/.ssh/id_rsa", conf.Key)
}

// ---------------------------------------------------------------------------
// Tests for searchKeyFromAddress
// ---------------------------------------------------------------------------

func TestSearchKeyFromAddress_Found(t *testing.T) {
	knownHosts := tempKnownHosts(t, []string{
		"192.168.1.1 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAItestkey",
		"10.0.0.1 ssh-rsa AAAAB3NzaC1yc2EAAAAtrankey",
	})

	result := searchKeyFromAddress(knownHosts, "192.168.1.1")
	assert.Equal(t, "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAItestkey", result)
}

func TestSearchKeyFromAddress_NotFound(t *testing.T) {
	knownHosts := tempKnownHosts(t, []string{
		"192.168.1.1 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAItestkey",
	})

	result := searchKeyFromAddress(knownHosts, "10.0.0.1")
	assert.Equal(t, "", result)
}

func TestSearchKeyFromAddress_FileNotFound(t *testing.T) {
	result := searchKeyFromAddress("/nonexistent/known_hosts", "192.168.1.1")
	assert.Equal(t, "", result)
}

func TestSearchKeyFromAddress_EmptyFile(t *testing.T) {
	knownHosts := tempKnownHosts(t, nil)
	result := searchKeyFromAddress(knownHosts, "192.168.1.1")
	assert.Equal(t, "", result)
}

// ---------------------------------------------------------------------------
// Tests for keyString
// ---------------------------------------------------------------------------

func TestKeyString(t *testing.T) {
	_, priv, err := generateEd25519KeyPair()
	assert.NoError(t, err)
	signer, err := ssh.NewSignerFromKey(priv)
	assert.NoError(t, err)
	pubKey := signer.PublicKey()

	result := keyString(pubKey)
	assert.Contains(t, result, pubKey.Type())
	parts := strings.SplitN(result, " ", 2)
	assert.Len(t, parts, 2)
	assert.Equal(t, "ssh-ed25519", parts[0])
	assert.NotEmpty(t, parts[1])
}

// ---------------------------------------------------------------------------
// Tests for trustedHostKeyCallback
// ---------------------------------------------------------------------------

func TestTrustedHostKeyCallback_NewHost(t *testing.T) {
	knownHosts := tempKnownHosts(t, nil)

	_, priv, err := generateEd25519KeyPair()
	assert.NoError(t, err)
	signer, err := ssh.NewSignerFromKey(priv)
	assert.NoError(t, err)
	pubKey := signer.PublicKey()

	hostKeyMutex := &sync.Mutex{}
	cb := trustedHostKeyCallback("", "192.168.1.100", hostKeyMutex, knownHosts)
	assert.NotNil(t, cb)

	err = cb("192.168.1.100", &net.IPAddr{IP: net.ParseIP("192.168.1.100")}, pubKey)
	assert.NoError(t, err)

	content, ok := utils.ReadFile(knownHosts)
	assert.True(t, ok)
	assert.Contains(t, string(content), "192.168.1.100")
	assert.Contains(t, string(content), pubKey.Type())
}

func TestTrustedHostKeyCallback_TrustedHost(t *testing.T) {
	_, priv, err := generateEd25519KeyPair()
	assert.NoError(t, err)
	signer, err := ssh.NewSignerFromKey(priv)
	assert.NoError(t, err)
	pubKey := signer.PublicKey()

	ks := keyString(pubKey)
	knownHosts := tempKnownHosts(t, []string{
		"192.168.1.100 " + ks,
	})

	hostKeyMutex := &sync.Mutex{}
	cb := trustedHostKeyCallback(ks, "192.168.1.100", hostKeyMutex, knownHosts)
	assert.NotNil(t, cb)

	err = cb("192.168.1.100", &net.IPAddr{IP: net.ParseIP("192.168.1.100")}, pubKey)
	assert.NoError(t, err)
}

func TestTrustedHostKeyCallback_Mismatch(t *testing.T) {
	_, priv1, err := generateEd25519KeyPair()
	assert.NoError(t, err)
	signer1, err := ssh.NewSignerFromKey(priv1)
	assert.NoError(t, err)

	_, priv2, err := generateEd25519KeyPair()
	assert.NoError(t, err)
	signer2, err := ssh.NewSignerFromKey(priv2)
	assert.NoError(t, err)

	trustedKey := keyString(signer1.PublicKey())

	knownHosts := tempKnownHosts(t, []string{
		"192.168.1.100 " + trustedKey,
	})

	hostKeyMutex := &sync.Mutex{}
	cb := trustedHostKeyCallback(trustedKey, "192.168.1.100", hostKeyMutex, knownHosts)
	assert.NotNil(t, cb)

	err = cb("192.168.1.100", &net.IPAddr{IP: net.ParseIP("192.168.1.100")}, signer2.PublicKey())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), SSHKeyKeyword)
}

func TestTrustedHostKeyCallback_NewHostButAlreadyExistsInFile(t *testing.T) {
	knownHosts := tempKnownHosts(t, []string{
		"192.168.1.100 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIsomeoldkey",
	})

	_, priv, err := generateEd25519KeyPair()
	assert.NoError(t, err)
	signer, err := ssh.NewSignerFromKey(priv)
	assert.NoError(t, err)

	hostKeyMutex := &sync.Mutex{}
	cb := trustedHostKeyCallback("", "192.168.1.100", hostKeyMutex, knownHosts)

	err = cb("192.168.1.100", &net.IPAddr{IP: net.ParseIP("192.168.1.100")}, signer.PublicKey())
	assert.NoError(t, err)
}

func TestTrustedHostKeyCallback_NewHostReadFails(t *testing.T) {
	nonexistent := filepath.Join(t.TempDir(), "missing_known_hosts")

	_, priv, err := generateEd25519KeyPair()
	assert.NoError(t, err)
	signer, err := ssh.NewSignerFromKey(priv)
	assert.NoError(t, err)

	hostKeyMutex := &sync.Mutex{}
	cb := trustedHostKeyCallback("", "192.168.1.100", hostKeyMutex, nonexistent)

	err = cb("192.168.1.100", &net.IPAddr{IP: net.ParseIP("192.168.1.100")}, signer.PublicKey())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "read known_hosts failed")
}

func TestTrustedHostKeyCallback_NewHostUpdateFails(t *testing.T) {
	dir := t.TempDir()
	knownHosts := filepath.Join(dir, "known_hosts")
	err := os.WriteFile(knownHosts, []byte(""), 0600)
	assert.NoError(t, err)
	err = os.Chmod(dir, 0500)
	assert.NoError(t, err)
	defer os.Chmod(dir, 0755)

	_, priv, err2 := generateEd25519KeyPair()
	assert.NoError(t, err2)
	signer, err2 := ssh.NewSignerFromKey(priv)
	assert.NoError(t, err2)

	hostKeyMutex := &sync.Mutex{}
	cb := trustedHostKeyCallback("", "10.0.0.99", hostKeyMutex, knownHosts)

	_ = cb("10.0.0.99", &net.IPAddr{IP: net.ParseIP("10.0.0.99")}, signer.PublicKey())
}

// ---------------------------------------------------------------------------
// Tests for loadKey
// ---------------------------------------------------------------------------

func TestLoadKey_EmptyKey(t *testing.T) {
	sc := &SSHConnector{Key: ""}
	result := sc.loadKey()
	assert.Nil(t, result)
}

func TestLoadKey_ValidKey(t *testing.T) {
	keyPath, _ := tempKeyFile(t)
	keysMap.Delete(keyPath)

	sc := &SSHConnector{Key: keyPath}
	result := sc.loadKey()
	assert.NotNil(t, result)
}

func TestLoadKey_InvalidPath(t *testing.T) {
	sc := &SSHConnector{Key: "/nonexistent/key/file"}
	result := sc.loadKey()
	assert.Nil(t, result)
}

func TestLoadKey_Caching(t *testing.T) {
	keyPath, _ := tempKeyFile(t)
	keysMap.Delete(keyPath)

	sc := &SSHConnector{Key: keyPath}

	first := sc.loadKey()
	assert.NotNil(t, first)

	os.Remove(keyPath)

	second := sc.loadKey()
	assert.NotNil(t, second)
	assert.Equal(t, first, second)

	keysMap.Delete(keyPath)
}

// ---------------------------------------------------------------------------
// Tests for getDialer / getKnownHosts
// ---------------------------------------------------------------------------

func TestGetDialer_Custom(t *testing.T) {
	mock := &mockSSHDialer{}
	sc := &SSHConnector{Dialer: mock}
	assert.Equal(t, mock, sc.getDialer())
}

func TestGetDialer_Default(t *testing.T) {
	sc := &SSHConnector{}
	d := sc.getDialer()
	assert.NotNil(t, d)
	_, ok := d.(defaultSSHDialer)
	assert.True(t, ok)
}

func TestGetKnownHosts_Custom(t *testing.T) {
	sc := &SSHConnector{KnownHosts: "/custom/known_hosts"}
	assert.Equal(t, "/custom/known_hosts", sc.getKnownHosts())
}

func TestGetKnownHosts_Default(t *testing.T) {
	sc := &SSHConnector{}
	assert.Equal(t, config.DefaultKnownHostsPath, sc.getKnownHosts())
}

// ---------------------------------------------------------------------------
// Additional tests for searchKeyFromAddress
// ---------------------------------------------------------------------------

func TestSearchKeyFromAddress_BracketHostPort(t *testing.T) {
	knownHosts := tempKnownHosts(t, []string{
		"[192.168.1.1]:2222 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIbracketkey",
	})

	result := searchKeyFromAddress(knownHosts, "192.168.1.1")
	assert.Equal(t, "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIbracketkey", result)
}

func TestSearchKeyFromAddress_BracketHostPortNoMatch(t *testing.T) {
	knownHosts := tempKnownHosts(t, []string{
		"[192.168.1.1]:2222 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIbracketkey",
	})

	result := searchKeyFromAddress(knownHosts, "10.0.0.1")
	assert.Equal(t, "", result)
}

func TestSearchKeyFromAddress_CommaSeparatedHosts(t *testing.T) {
	knownHosts := tempKnownHosts(t, []string{
		"host1,host2,host3 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIcommakey",
	})

	result := searchKeyFromAddress(knownHosts, "host2")
	assert.Equal(t, "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIcommakey", result)
}

func TestSearchKeyFromAddress_CommaSeparatedHostsNoMatch(t *testing.T) {
	knownHosts := tempKnownHosts(t, []string{
		"host1,host2,host3 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIcommakey",
	})

	result := searchKeyFromAddress(knownHosts, "host4")
	assert.Equal(t, "", result)
}

func TestSearchKeyFromAddress_LinesStartingWithAt(t *testing.T) {
	knownHosts := tempKnownHosts(t, []string{
		"@revoked 192.168.1.1 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIrevokedkey",
		"192.168.1.1 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIvalidkey",
	})

	result := searchKeyFromAddress(knownHosts, "192.168.1.1")
	assert.Equal(t, "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIvalidkey", result)
}

func TestSearchKeyFromAddress_CommaSeparatedWithBracketHost(t *testing.T) {
	knownHosts := tempKnownHosts(t, []string{
		"host1,[::1]:22 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAImixedkey",
	})

	result := searchKeyFromAddress(knownHosts, "::1")
	assert.Equal(t, "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAImixedkey", result)
}

func TestSearchKeyFromAddress_MalformedLine(t *testing.T) {
	knownHosts := tempKnownHosts(t, []string{
		"malformedlinewithoutspaces",
		"192.168.1.1 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIvalidkey",
	})

	result := searchKeyFromAddress(knownHosts, "192.168.1.1")
	assert.Equal(t, "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIvalidkey", result)
}

// ---------------------------------------------------------------------------
// Tests for CloseConnection error path
// ---------------------------------------------------------------------------

func TestCloseConnection_AlreadyClosed(t *testing.T) {
	ts, _ := newTestServer(t, "testuser", "testpass")
	defer ts.close()

	sc := newConnectorForServer(ts, "testuser", "testpass")
	sc.KnownHosts = tempKnownHosts(t, nil)

	client, err := sc.CreateConnection()
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// Close once normally
	sc.CloseConnection(client)
	// Close again - exercises the error branch in CloseConnection
	sc.CloseConnection(client)
}
