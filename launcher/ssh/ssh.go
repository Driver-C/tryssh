package ssh

import (
	"github.com/Driver-C/tryssh/launcher"
	"github.com/Driver-C/tryssh/utils"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"time"
)

type Launcher struct {
	launcher.SshConnector
}

func (h *Launcher) Launch() bool {
	return h.dialServer()
}

func NewSshLaunchersByCombinations(combinations chan []interface{},
	sshTimeout time.Duration) (launchers []*Launcher) {
	for com := range combinations {
		launchers = append(launchers, &Launcher{launcher.SshConnector{
			Ip:         com[0].(string),
			Port:       com[1].(string),
			User:       com[2].(string),
			Password:   com[3].(string),
			SshTimeout: sshTimeout,
		}})
	}
	return
}

func (h *Launcher) dialServer() (res bool) {
	res = false
	sshClient, err := h.CreateConnection()
	if err == nil {
		utils.Logger.Infoln("[ LOGIN SUCCESSFUL ]\n")
		utils.Logger.Infoln("User:", sshClient.User())
		utils.Logger.Infoln("Password:", h.Password)
		utils.Logger.Infoln("Port:", h.Port)
		utils.Logger.Infoln("Ssh Server Version:", string(sshClient.ServerVersion()))
		utils.Logger.Infof("Ssh Client Version: %s\n\n", string(sshClient.ClientVersion()))
		res = true
		h.createTerminal(sshClient)
	} else {
		return
	}
	defer h.CloseConnection(sshClient)
	return
}

func (h *Launcher) createTerminal(conn *ssh.Client) {
	session, err := conn.NewSession()
	if err != nil {
		utils.Logger.Fatalln(err.Error())
	}
	defer func(conn *ssh.Client) {
		if err := session.Close(); err != nil {
			if err.Error() != "EOF" {
				utils.Logger.Fatalln(err.Error())
			}
		}
	}(conn)

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
		ssh.VSTATUS:       1,
	}
	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		utils.Logger.Fatalln(err.Error())
	}
	defer func(fd int, oldState *terminal.State) {
		if err := terminal.Restore(fd, oldState); err != nil {
			utils.Logger.Fatalln(err.Error())
		}
	}(fd, oldState)

	termWidth, termHeight, err := terminal.GetSize(fd)
	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	err = session.RequestPty(launcher.TerminalTerm, termHeight, termWidth, modes)
	if err != nil {
		utils.Logger.Fatalln(err.Error())
	}

	err = session.Shell()
	if err != nil {
		utils.Logger.Fatalln(err.Error())
	}

	err = session.Wait()
	if err != nil {
		utils.Logger.Warnln(err.Error())
	}
}
