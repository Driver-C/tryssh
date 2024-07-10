package launcher

import (
	"encoding/base64"
	"fmt"
	"github.com/Driver-C/tryssh/config"
	"github.com/Driver-C/tryssh/utils"
	"golang.org/x/crypto/ssh"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	sshProtocol   string = "tcp"
	TerminalTerm         = "xterm"
	SSHKeyKeyword        = "SSH-KEY"
)

var keysMap = sync.Map{}

type Connector interface {
	Launch() bool
	CreateConnection() (sshClient *ssh.Client, err error)
	CloseConnection(sshClient *ssh.Client)
	TryToConnect() (err error)
}

type SshConnector struct {
	Ip           string
	Port         string
	User         string
	Password     string
	Key          string
	SshTimeout   time.Duration
	HostKeyMutex *sync.Mutex
}

func (sc *SshConnector) Launch() bool {
	return false
}

func (sc *SshConnector) LoadConfig() (config *ssh.ClientConfig) {
	// If no mutex is passed in, initialize one
	if sc.HostKeyMutex == nil {
		sc.HostKeyMutex = new(sync.Mutex)
	}

	var authMethods []ssh.AuthMethod
	var privateKey []byte
	if sc.Key != "" {
		if _, ok := keysMap.Load(sc.Key); !ok {
			if pk, status := utils.ReadFile(sc.Key); status {
				keysMap.Store(sc.Key, pk)
				privateKey = pk
			}
		} else {
			pk, _ := keysMap.Load(sc.Key)
			privateKey = pk.([]byte)
		}
		signer, err := ssh.ParsePrivateKey(privateKey)
		if err == nil {
			authMethods = append(authMethods, ssh.PublicKeys(signer))
		} else {
			utils.Logger.Errorln("Failed to parse private key: %v", err)
		}
	}
	authMethods = append(authMethods, ssh.Password(sc.Password))
	config = &ssh.ClientConfig{
		User:            sc.User,
		Auth:            authMethods,
		HostKeyCallback: trustedHostKeyCallback(searchKeyFromAddress(sc.Ip), sc.Ip, sc.HostKeyMutex),
		Timeout:         sc.SshTimeout,
	}
	return
}

func (sc *SshConnector) CreateConnection() (sshClient *ssh.Client, err error) {
	addr := sc.Ip + ":" + sc.Port
	conf := sc.LoadConfig()

	sshClient, err = ssh.Dial(sshProtocol, addr, conf)
	if err != nil {
		if strings.Contains(err.Error(), SSHKeyKeyword) {
			// If it's a public key verification issue, just exit
			utils.Logger.Fatalf("Unable to connect: %s@%s, Password:%s Cause: %s\n",
				sc.User, addr, sc.Password, err.Error())
		} else {
			utils.Logger.Warnf("Unable to connect: %s@%s, Password:%s Cause: %s\n",
				sc.User, addr, sc.Password, err.Error())
		}
	}
	return
}

func (sc *SshConnector) CloseConnection(sshClient *ssh.Client) {
	err := sshClient.Close()
	if err != nil {
		utils.Logger.Errorln("Unable to close connection: ", err.Error())
	}
}

func (sc *SshConnector) TryToConnect() (err error) {
	sshClient, err := sc.CreateConnection()
	if err != nil {
		return
	}
	defer sc.CloseConnection(sshClient)
	return
}

// GetSshConnectorFromConfig Get SshConnector by ServerListConfig
func GetSshConnectorFromConfig(conf *config.ServerListConfig) *SshConnector {
	return &SshConnector{
		Ip:       conf.Ip,
		Port:     conf.Port,
		User:     conf.User,
		Password: conf.Password,
		Key:      conf.Key,
	}
}

// GetConfigFromSshConnector Get ServerListConfig by SshConnector
func GetConfigFromSshConnector(tgt *SshConnector) *config.ServerListConfig {
	return &config.ServerListConfig{
		Ip:       tgt.Ip,
		Port:     tgt.Port,
		User:     tgt.User,
		Password: tgt.Password,
		Key:      tgt.Key,
	}
}

func searchKeyFromAddress(address string) string {
	knownHostsContent, status := utils.ReadFile(config.KnownHostsPath)
	if !status {
		utils.Logger.Fatalln("Read known_hosts failed")
	}
	knownHostsLines := strings.Split(string(knownHostsContent), "\n")
	for _, line := range knownHostsLines {
		if strings.Split(line, " ")[0] == address {
			return strings.Join(strings.Split(line, " ")[1:], " ")
		}
	}
	return ""
}

func keyString(k ssh.PublicKey) string {
	return k.Type() + " " + base64.StdEncoding.EncodeToString(k.Marshal())
}

func trustedHostKeyCallback(trustedKey string, address string, hostKeyMutex *sync.Mutex) ssh.HostKeyCallback {
	if trustedKey == "" {
		return func(_ string, _ net.Addr, k ssh.PublicKey) error {
			hostKeyMutex.Lock()
			defer hostKeyMutex.Unlock()
			// Re search for key to avoid duplicate operations
			if searchKeyFromAddress(address) != "" {
				return nil
			}
			newHostKeyInfo := address + " " + keyString(k) + "\n"
			if knownHostsContent, status := utils.ReadFile(config.KnownHostsPath); status {
				knownHostsContent = append(knownHostsContent, []byte(newHostKeyInfo)...)
				if utils.UpdateFile(config.KnownHostsPath, knownHostsContent, 0600) {
					utils.Logger.Infoln("First login, automatically add key to known_hosts")
					return nil
				} else {
					return fmt.Errorf("update known_hosts failed")
				}
			} else {
				return fmt.Errorf("read known_hosts failed")
			}
		}
	}

	return func(_ string, _ net.Addr, k ssh.PublicKey) error {
		ks := keyString(k)
		if trustedKey != ks {
			return fmt.Errorf("\n*[%s]* ssh-key verification: expected %q but got %q\n"+
				"*[%s]* Server [%s] may have been impersonated. "+
				"If you can confirm that the public key change of server [%s] "+
				"is normal, please delete the entry for server [%s] in ~/.tryssh/known_hosts "+
				"and try logging in again.",
				SSHKeyKeyword, trustedKey, ks, SSHKeyKeyword, address, address, address)
		}
		return nil
	}
}
