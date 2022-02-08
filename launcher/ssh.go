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
	sshClientTimeout        = 10000 * time.Millisecond
)

// SshLauncher ssh终端连接器
type SshLauncher struct {
	target.SshTarget
}

// Launch 执行
func (h *SshLauncher) Launch() bool {
	return h.DialServer()
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

// DialServer 连接服务器
func (h *SshLauncher) DialServer() (res bool) {
	res = false
	addr := h.Ip + ":" + h.Port
	config := h.LoadConfig()

	conn, err := ssh.Dial(sshProtocol, addr, config)
	if err != nil {
		log.Warnf("Unable to connect: %s@%s, Password:%s Cause: %s",
			h.User, addr, h.Password, err.Error())
		return
	} else {
		log.Infoln("[ LOGIN SUCCESSFUL ]\n")
		log.Infoln("User:", conn.User())
		log.Infoln("Password:", h.Password)
		log.Infoln("Ssh Server Version:", string(conn.ServerVersion()))
		log.Infof("Ssh Client Version: %s\n\n", string(conn.ClientVersion()))
		res = true
		h.CreateTerminal(conn)
	}
	defer func(conn *ssh.Client) {
		err := conn.Close()
		if err != nil {
			log.Errorln("Unable to close connection: ", err.Error())
		}
	}(conn)
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
