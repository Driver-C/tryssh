package launcher

import (
	"github.com/Driver-C/tryssh/pkg/utils"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"time"
)

type SshLauncher struct {
	SshConnector
}

func (h *SshLauncher) Launch() bool {
	return h.dialServer()
}

func NewSshLaunchersByCombinations(combinations chan []interface{},
	sshTimeout time.Duration) (launchers []*SshLauncher) {
	for com := range combinations {
		launchers = append(launchers, &SshLauncher{SshConnector{
			Ip:         com[0].(string),
			Port:       com[1].(string),
			User:       com[2].(string),
			Password:   com[3].(string),
			Key:        com[4].(string),
			SshTimeout: sshTimeout,
		}})
	}
	return
}

func (h *SshLauncher) dialServer() (res bool) {
	res = false
	sshClient, err := h.CreateConnection()
	if err == nil {
		utils.Logger.Infoln("[ LOGIN SUCCESSFUL ]\n")
		utils.Logger.Infoln("User:", sshClient.User())
		utils.Logger.Infoln("Port:", h.Port)
		res = true
		h.createTerminal(sshClient)
	} else {
		return
	}
	defer h.CloseConnection(sshClient)
	return
}

func (h *SshLauncher) createTerminal(conn *ssh.Client) {
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

	err = session.RequestPty(TerminalTerm, termHeight, termWidth, modes)
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
