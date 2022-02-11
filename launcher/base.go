package launcher

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"time"
)

const (
	sshProtocol      string = "tcp"
	sshClientTimeout        = 1000 * time.Millisecond
	terminalTerm            = "xterm"
)

// Connector 连接器接口
type Connector interface {
	Launch() bool
	CreateConnection() (sshClient *ssh.Client, err error)
	CloseConnection(sshClient *ssh.Client)
	TryToConnect() (err error)
}

// SshConnector ssh连接器
type SshConnector struct {
	Ip       string
	Port     string
	User     string
	Password string
}

// Launch 执行
func (sc *SshConnector) Launch() bool {
	return false
}

// LoadConfig 构造连接配置
func (sc *SshConnector) LoadConfig() (config *ssh.ClientConfig) {
	config = &ssh.ClientConfig{
		User: sc.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(sc.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		ClientVersion:   "",
		Timeout:         sshClientTimeout,
	}
	return
}

// CreateConnection 创建连接
func (sc *SshConnector) CreateConnection() (sshClient *ssh.Client, err error) {
	addr := sc.Ip + ":" + sc.Port
	config := sc.LoadConfig()

	sshClient, err = ssh.Dial(sshProtocol, addr, config)
	if err != nil {
		log.Warnf("Unable to connect: %s@%s, Password:%s Cause: %s",
			sc.User, addr, sc.Password, err.Error())
	}
	return
}

// CloseConnection 关闭连接
func (sc *SshConnector) CloseConnection(sshClient *ssh.Client) {
	err := sshClient.Close()
	if err != nil {
		log.Errorln("Unable to close connection: ", err.Error())
	}
}

// TryToConnect 尝试创建连接
func (sc *SshConnector) TryToConnect() (err error) {
	sshClient, err := sc.CreateConnection()
	if err != nil {
		return
	}
	defer sc.CloseConnection(sshClient)
	return
}
