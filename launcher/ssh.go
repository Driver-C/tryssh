package launcher

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

// SshLauncher ssh终端连接器
type SshLauncher struct {
	SshConnector
}

// Launch 执行
func (h *SshLauncher) Launch() bool {
	return h.dialServer()
}

// NewSshLaunchersByCombinations 通过用户、密码和端口的组合生成SshLauncher对象切片
func NewSshLaunchersByCombinations(combinations chan []interface{}) (launchers []SshLauncher) {
	for com := range combinations {
		launchers = append(launchers, SshLauncher{SshConnector{
			Ip:       com[0].(string),
			Port:     com[1].(string),
			User:     com[2].(string),
			Password: com[3].(string),
		}})
	}
	return
}

// dialServer 连接服务器
func (h *SshLauncher) dialServer() (res bool) {
	res = false
	sshClient, err := h.CreateConnection()
	if err == nil {
		log.Infoln("[ LOGIN SUCCESSFUL ]\n")
		log.Infoln("User:", sshClient.User())
		log.Infoln("Password:", h.Password)
		log.Infoln("Ssh Server Version:", string(sshClient.ServerVersion()))
		log.Infof("Ssh Client Version: %s\n\n", string(sshClient.ClientVersion()))
		res = true
		h.createTerminal(sshClient)
	} else {
		return
	}
	defer h.CloseConnection(sshClient)
	return
}

// createTerminal 创建session和终端
func (h *SshLauncher) createTerminal(conn *ssh.Client) {
	// 创建session
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer func(conn *ssh.Client) {
		if err := session.Close(); err != nil {
			if err.Error() != "EOF" {
				log.Fatalln(err.Error())
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
			log.Fatalln(err.Error())
		}
	}(fd, oldState)

	termWidth, termHeight, err := terminal.GetSize(fd)
	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// 打开伪终端
	err = session.RequestPty(terminalTerm, termHeight, termWidth, modes)
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
