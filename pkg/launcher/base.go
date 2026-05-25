// Package launcher provides SSH/SCP connection and file transfer capabilities.
package launcher

import (
	"encoding/base64"
	"fmt"
	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/utils"
	"golang.org/x/crypto/ssh"
	"net"
	"strings"
	"sync"
	"time"
)

// SSHProtocol, TerminalTerm, and SSHKeyKeyword are constants used in SSH connection setup.
const (
	SSHProtocol  = "tcp"
	TerminalTerm        = "xterm"
	SSHKeyKeyword       = "SSH-KEY"
)

var keysMap = sync.Map{}

// hostKeyCallbackMutex protects concurrent known_hosts file access across all connections.
var hostKeyCallbackMutex sync.Mutex

// SSHDialer abstracts the SSH dial operation for testability.
type SSHDialer interface {
	Dial(network, address string, config *ssh.ClientConfig) (*ssh.Client, error)
}

type defaultSSHDialer struct{}

func (d defaultSSHDialer) Dial(network, address string, config *ssh.ClientConfig) (*ssh.Client, error) {
	return ssh.Dial(network, address, config)
}

// Connector defines the interface for launching SSH/SCP connections and testing connectivity.
type Connector interface {
	Launch() bool
	CreateConnection() (sshClient *ssh.Client, err error)
	CloseConnection(sshClient *ssh.Client)
	TryToConnect() (err error)
}

// SSHConnector holds the parameters and state needed to establish an SSH connection.
type SSHConnector struct {
	IP         string
	Port       string
	User       string
	Password   string
	Key        string
	SSHTimeout time.Duration
	Dialer     SSHDialer
	KnownHosts string
}

func (sc *SSHConnector) getDialer() SSHDialer {
	if sc.Dialer != nil {
		return sc.Dialer
	}
	return defaultSSHDialer{}
}

func (sc *SSHConnector) getKnownHosts() string {
	if sc.KnownHosts != "" {
		return sc.KnownHosts
	}
	return config.DefaultKnownHostsPath
}

// BuildSSHConfig creates an SSH client configuration with appropriate auth methods.
func (sc *SSHConnector) BuildSSHConfig() (cfg *ssh.ClientConfig) {
	var authMethods []ssh.AuthMethod
	if sc.Key != "" {
		privateKey := sc.loadKey()
		if privateKey != nil {
			signer, err := ssh.ParsePrivateKey(privateKey)
			if err == nil {
				authMethods = append(authMethods, ssh.PublicKeys(signer))
			} else {
				utils.Errorln("Failed to parse private key: ", err)
			}
		}
	}
	if sc.Password != "" {
		authMethods = append(authMethods, ssh.Password(sc.Password))
	}
	cfg = &ssh.ClientConfig{
		User:            sc.User,
		Auth:            authMethods,
		HostKeyCallback: trustedHostKeyCallback(searchKeyFromAddress(sc.getKnownHosts(), sc.IP), sc.IP, &hostKeyCallbackMutex, sc.getKnownHosts()),
		Timeout:         sc.SSHTimeout,
	}
	return
}

func (sc *SSHConnector) loadKey() []byte {
	if sc.Key == "" {
		return nil
	}
	if cached, ok := keysMap.Load(sc.Key); ok {
		return cached.([]byte)
	}
	if pk, status := utils.ReadFile(sc.Key); status {
		keysMap.Store(sc.Key, pk)
		return pk
	}
	return nil
}

// CreateConnection establishes an SSH connection using the configured parameters.
func (sc *SSHConnector) CreateConnection() (sshClient *ssh.Client, err error) {
	addr := sc.IP + ":" + sc.Port
	conf := sc.BuildSSHConfig()

	sshClient, err = sc.getDialer().Dial(SSHProtocol, addr, conf)
	if err != nil {
		if strings.Contains(err.Error(), SSHKeyKeyword) {
			utils.Errorf("Unable to connect: %s Cause: %s", addr, err.Error())
		}
	}
	return
}

// CloseConnection closes the given SSH client connection.
func (sc *SSHConnector) CloseConnection(sshClient *ssh.Client) {
	err := sshClient.Close()
	if err != nil {
		utils.Errorln("Unable to close connection: ", err.Error())
	}
}

// TryToConnect attempts to establish and then immediately close an SSH connection.
func (sc *SSHConnector) TryToConnect() (err error) {
	sshClient, err := sc.CreateConnection()
	if err != nil {
		return
	}
	defer sc.CloseConnection(sshClient)
	return
}

// GetSSHConnectorFromConfig creates an SSHConnector from a ServerListConfig entry.
func GetSSHConnectorFromConfig(conf *config.ServerListConfig) *SSHConnector {
	return &SSHConnector{
		IP:         conf.IP,
		Port:       conf.Port,
		User:       conf.User,
		Password:   conf.Password,
		Key:        conf.Key,
		SSHTimeout: 5 * time.Second,
	}
}

// GetConfigFromSSHConnector converts an SSHConnector back into a ServerListConfig.
func GetConfigFromSSHConnector(tgt *SSHConnector) *config.ServerListConfig {
	return &config.ServerListConfig{
		IP:       tgt.IP,
		Port:     tgt.Port,
		User:     tgt.User,
		Password: tgt.Password,
		Key:      tgt.Key,
	}
}

func searchKeyFromAddress(knownHostsPath, address string) string {
	knownHostsContent, status := utils.ReadFile(knownHostsPath)
	if !status {
		return ""
	}
	knownHostsLines := strings.Split(string(knownHostsContent), "\n")
	for _, line := range knownHostsLines {
		if len(line) == 0 || line[0] == '@' {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		// Match against plain address or [address]:port format
		hostPart := parts[0]
		if hostPart == address {
			return parts[1]
		}
		// Check if hostPart is a comma-separated list of hostnames
		for _, h := range strings.Split(hostPart, ",") {
			if h == address {
				return parts[1]
			}
			// Match [address]:port entries in known_hosts
			if strings.HasPrefix(h, "[") {
				if bracketClose := strings.Index(h, "]"); bracketClose > 0 {
					if h[1:bracketClose] == address {
						return parts[1]
					}
				}
			}
		}
	}
	return ""
}

func keyString(k ssh.PublicKey) string {
	return k.Type() + " " + base64.StdEncoding.EncodeToString(k.Marshal())
}

func trustedHostKeyCallback(trustedKey string, address string, hostKeyMutex *sync.Mutex, knownHostsPath string) ssh.HostKeyCallback {
	if trustedKey == "" {
		return func(_ string, _ net.Addr, k ssh.PublicKey) error {
			hostKeyMutex.Lock()
			defer hostKeyMutex.Unlock()
			if searchKeyFromAddress(knownHostsPath, address) != "" {
				return nil
			}
			newHostKeyInfo := address + " " + keyString(k) + "\n"
			if knownHostsContent, status := utils.ReadFile(knownHostsPath); status {
				knownHostsContent = append(knownHostsContent, []byte(newHostKeyInfo)...)
				if err := utils.UpdateFile(knownHostsPath, knownHostsContent, 0600); err != nil {
					return fmt.Errorf("update known_hosts failed: %w", err)
				}
				utils.Infof("First login to %s, automatically adding host key to known_hosts (TOFU)\n", address)
				return nil
			}
			return fmt.Errorf("read known_hosts failed")
		}
	}

	return func(_ string, _ net.Addr, k ssh.PublicKey) error {
		ks := keyString(k)
		if trustedKey != ks {
			return fmt.Errorf("*[%s]* ssh-key verification: expected %q but got %q "+
				"*[%s]* Server [%s] may have been impersonated. "+
				"If you can confirm that the public key change of server [%s] "+
				"is normal, please delete the entry for server [%s] in ~/.tryssh/known_hosts "+
				"and try logging in again",
				SSHKeyKeyword, trustedKey, ks, SSHKeyKeyword, address, address, address)
		}
		return nil
	}
}
