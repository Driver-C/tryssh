package launcher

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Driver-C/tryssh/pkg/utils"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

// SSHLauncher handles interactive SSH terminal sessions.
type SSHLauncher struct {
	SSHConnector
}

// Launch starts an interactive SSH session and returns true on success.
func (h *SSHLauncher) Launch() bool {
	return h.dialServer()
}

// NewSSHLaunchersByCombinations creates SSHLauncher instances from a channel of credential combinations.
func NewSSHLaunchersByCombinations(combinations chan []interface{},
	sshTimeout time.Duration) (launchers []*SSHLauncher) {
	for com := range combinations {
		ip, _ := com[0].(string)
		port, _ := com[1].(string)
		user, _ := com[2].(string)
		password, _ := com[3].(string)
		key, _ := com[4].(string)
		launchers = append(launchers, &SSHLauncher{SSHConnector{
			IP:         ip,
			Port:       port,
			User:       user,
			Password:   password,
			Key:        key,
			SSHTimeout: sshTimeout,
		}})
	}
	return
}

func (h *SSHLauncher) dialServer() bool {
	sshClient, err := h.CreateConnection()
	if err != nil {
		return false
	}
	defer h.CloseConnection(sshClient)

	utils.Infoln("[ LOGIN SUCCESSFUL ]")
	utils.Infoln("User:", sshClient.User())
	utils.Infoln("Port:", h.Port)
	if err := h.createTerminal(sshClient); err != nil {
		utils.Errorf("Terminal session failed: %v", err)
	}
	return true
}

func (h *SSHLauncher) createTerminal(conn *ssh.Client) error {
	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := session.Close(); closeErr != nil && closeErr.Error() != "EOF" {
			utils.Errorln(closeErr.Error())
		}
	}()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer func() {
		if restoreErr := term.Restore(fd, oldState); restoreErr != nil {
			utils.Errorln(restoreErr.Error())
		}
	}()

	termWidth, termHeight, _ := term.GetSize(fd)
	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if ptyErr := session.RequestPty(TerminalTerm, termHeight, termWidth, modes); ptyErr != nil {
		return ptyErr
	}

	// Handle terminal resize via SIGWINCH
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGWINCH)
	go func() {
		for range sigChan {
			w, h, _ := term.GetSize(fd)
			if w > 0 && h > 0 {
				_ = session.WindowChange(h, w)
			}
		}
	}()
	defer signal.Stop(sigChan)

	if shellErr := session.Shell(); shellErr != nil {
		return shellErr
	}

	if waitErr := session.Wait(); waitErr != nil {
		utils.Warnln(waitErr.Error())
	}
	return nil
}
