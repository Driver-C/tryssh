package launcher

import "golang.org/x/crypto/ssh"

type Connector interface {
	Launch() bool
	CreateConnection() (sshClient *ssh.Client, err error)
	CloseConnection(sshClient *ssh.Client)
	TryToConnect() (err error)
}
