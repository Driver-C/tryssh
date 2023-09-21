package launcher

import (
	"golang.org/x/crypto/ssh"
	"time"
	"tryssh/utils"
)

const (
	sshProtocol      string = "tcp"
	sshClientTimeout        = 2 * time.Second
	TerminalTerm            = "xterm"
)

type Connector interface {
	Launch() bool
	CreateConnection() (sshClient *ssh.Client, err error)
	CloseConnection(sshClient *ssh.Client)
	TryToConnect() (err error)
}

type SshConnector struct {
	Ip         string
	Port       string
	User       string
	Password   string
	SshTimeout time.Duration
}

func (sc *SshConnector) Launch() bool {
	return false
}

func (sc *SshConnector) LoadConfig() (config *ssh.ClientConfig) {
	// If no timeout has been set, set 5 second
	if sc.SshTimeout == 0 {
		sc.SshTimeout = 5 * time.Second
	}
	config = &ssh.ClientConfig{
		User: sc.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(sc.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		ClientVersion:   "",
		Timeout:         sc.SshTimeout,
	}
	return
}

func (sc *SshConnector) CreateConnection() (sshClient *ssh.Client, err error) {
	addr := sc.Ip + ":" + sc.Port
	config := sc.LoadConfig()

	sshClient, err = ssh.Dial(sshProtocol, addr, config)
	if err != nil {
		utils.Logger.Warnf("Unable to connect: %s@%s, Password:%s Cause: %s\n",
			sc.User, addr, sc.Password, err.Error())
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
