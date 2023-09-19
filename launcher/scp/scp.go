package scp

import (
	"github.com/pkg/sftp"
	"io"
	"os"
	"strings"
	"tryssh/launcher"
	"tryssh/utils"
)

type Launcher struct {
	launcher.SshConnector
	Src  string
	Dest string
}

func (c *Launcher) Launch() bool {
	switch {
	case strings.Contains(c.Src, c.Ip):
		return c.download(c.Dest, strings.Split(c.Src, ":")[1])
	case strings.Contains(c.Dest, c.Ip):
		return c.upload(c.Src, strings.Split(c.Dest, ":")[1])
	}
	return false
}

func NewScpLaunchersByCombinations(combinations chan []interface{}, src string, dest string) (launchers []Launcher) {
	for com := range combinations {
		launchers = append(launchers, Launcher{
			SshConnector: launcher.SshConnector{
				Ip:       com[0].(string),
				Port:     com[1].(string),
				User:     com[2].(string),
				Password: com[3].(string),
			},
			Src:  src,
			Dest: dest,
		})
	}
	return
}

func (c *Launcher) createScpClient() (sftpClient *sftp.Client) {
	sshClient, errSsh := c.CreateConnection()
	if errSsh != nil {
		return
	}
	sftpClient, errScp := sftp.NewClient(sshClient)
	if errScp != nil {
		utils.Logger.Fatalln(errScp.Error())
	}
	return
}

func (c *Launcher) closeScpClient(sftpClient *sftp.Client) {
	err := sftpClient.Close()
	if err != nil {
		utils.Logger.Errorln(err.Error())
	}
}

func (c *Launcher) upload(local, remote string) bool {
	sftpClient := c.createScpClient()
	if sftpClient == nil {
		return false
	}
	defer c.closeScpClient(sftpClient)
	localPath := strings.Split(local, "/")
	localFileName := localPath[len(localPath)-1]
	remoteFileName := localFileName

	localFile, err := os.Open(local)
	if err != nil {
		utils.Logger.Fatalln(err.Error())
	}
	defer func(localFile *os.File) {
		err := localFile.Close()
		if err != nil {
			utils.Logger.Errorln(err.Error())
		}
	}(localFile)

	remoteFile, err := sftpClient.OpenFile(sftp.Join(remote, remoteFileName), os.O_CREATE|os.O_RDWR|os.O_EXCL)
	if err != nil {
		utils.Logger.Fatalln(err.Error())
	}
	defer func(remoteFile *sftp.File) {
		err := remoteFile.Close()
		if err != nil {
			utils.Logger.Errorln(err.Error())
		}
	}(remoteFile)

	transSize, err := io.Copy(remoteFile, localFile)
	if err != nil {
		utils.Logger.Fatalln(err.Error())
	} else {
		utils.Logger.Infof("File: %s uploaded successfully. Transmission size: %s.\n",
			remoteFileName, utils.ByteSizeFormat(float64(transSize)))
	}
	return true
}

func (c *Launcher) download(local, remote string) bool {
	sftpClient := c.createScpClient()
	if sftpClient == nil {
		return false
	}
	defer c.closeScpClient(sftpClient)
	remotePath := strings.Split(remote, "/")
	remoteFileName := remotePath[len(remotePath)-1]
	localFileName := remoteFileName

	remoteFile, err := sftpClient.Open(remote)
	if err != nil {
		utils.Logger.Fatalln(err.Error())
	}
	defer func(remoteFile *sftp.File) {
		err := remoteFile.Close()
		if err != nil {
			utils.Logger.Errorln(err.Error())
		}
	}(remoteFile)

	localFile, err := os.OpenFile(sftp.Join(local, localFileName), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		utils.Logger.Fatalln(err.Error())
	}
	defer func(localFile *os.File) {
		err := localFile.Close()
		if err != nil {
			utils.Logger.Errorln(err.Error())
		}
	}(localFile)

	transSize, err := io.Copy(localFile, remoteFile)
	if err != nil {
		utils.Logger.Fatalln(err.Error())
	} else {
		utils.Logger.Infof("File: %s downloaded successfully. Transmission size: %s.\n",
			remoteFileName, utils.ByteSizeFormat(float64(transSize)))
	}
	return true
}
