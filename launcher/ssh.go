package launcher

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"time"
	"tryssh/target"
)

const (
	sshProtocol      string = "tcp"
	sshClientTimeout        = 500 * time.Millisecond
)

// SshLauncher ssh终端连接器
type SshLauncher struct {
	target.SshTarget
}

// Launch 执行
func (h *SshLauncher) Launch() bool {
	return h.DialServer()
}

// NewSshLaunchersByCombinations 通过用户、密码和端口的组合生成SshLauncher对象切片
func NewSshLaunchersByCombinations(combinations chan []interface{}) (launchers []SshLauncher) {
	for com := range combinations {
		launchers = append(launchers, SshLauncher{SshTarget: target.SshTarget{
			Ip:       com[0].(string),
			Port:     com[1].(string),
			User:     com[2].(string),
			Password: com[3].(string),
		}})
	}
	return
}

// LoadConfig 构造连接配置
func (h *SshLauncher) LoadConfig() (config *ssh.ClientConfig) {
	config = &ssh.ClientConfig{
		User: h.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(h.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		ClientVersion:   "",
		Timeout:         sshClientTimeout,
	}
	return
}

// CreateConnection 创建连接
func (h *SshLauncher) CreateConnection() (sshClient *ssh.Client, err error) {
	addr := h.Ip + ":" + h.Port
	config := h.LoadConfig()

	sshClient, err = ssh.Dial(sshProtocol, addr, config)
	if err != nil {
		log.Warnf("Unable to connect: %s@%s, Password:%s Cause: %s",
			h.User, addr, h.Password, err.Error())
	}
	return
}

// CloseConnection 关闭连接
func (h *SshLauncher) CloseConnection(sshClient *ssh.Client) {
	err := sshClient.Close()
	if err != nil {
		log.Errorln("Unable to close connection: ", err.Error())
	}
}

// TryToConnect 尝试创建连接
func (h *SshLauncher) TryToConnect() (err error) {
	sshClient, err := h.CreateConnection()
	if err != nil {
		return
	}
	defer h.CloseConnection(sshClient)
	return
}

// DialServer 连接服务器
func (h *SshLauncher) DialServer() (res bool) {
	res = false
	sshClient, err := h.CreateConnection()
	if err == nil {
		log.Infoln("[ LOGIN SUCCESSFUL ]\n")
		log.Infoln("User:", sshClient.User())
		log.Infoln("Password:", h.Password)
		log.Infoln("Ssh Server Version:", string(sshClient.ServerVersion()))
		log.Infof("Ssh Client Version: %s\n\n", string(sshClient.ClientVersion()))
		res = true
		h.CreateTerminal(sshClient)
	} else {
		return
	}
	defer h.CloseConnection(sshClient)
	return
}

// CreateTerminal 创建session和终端
func (h *SshLauncher) CreateTerminal(conn *ssh.Client) {
	// 创建session
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalln("Failed to create ssh session: ", err.Error())
	}
	defer func(conn *ssh.Client) {
		if err := session.Close(); err != nil {
			if err.Error() != "EOF" {
				log.Fatalln("Failed to close ssh session: ", err.Error())
			}
		}
	}(conn)

	// 配置终端
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // 打开回显 0关 1开
		ssh.TTY_OP_ISPEED: 14400, // 输入速率 14.4k baud
		ssh.TTY_OP_OSPEED: 14400, // 输出速率 14.4k baud
		ssh.VSTATUS:       1,
	}
	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer func(fd int, oldState *terminal.State) {
		if err := terminal.Restore(fd, oldState); err != nil {
			log.Fatalln("Failed to restore terminal: ", err.Error())
		}
	}(fd, oldState)

	termWidth, termHeight, err := terminal.GetSize(fd)
	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// 打开伪终端
	err = session.RequestPty("xterm", termHeight, termWidth, modes)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// 启动一个远程shell
	err = session.Shell()
	if err != nil {
		log.Fatalln(err.Error())
	}

	// 等待远程命令结束或远程shell退出
	err = session.Wait()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
