package scp

import (
	"github.com/cheggaaa/pb/v3"
	"github.com/pkg/sftp"
	"io"
	"os"
	"strings"
	"time"
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

func NewScpLaunchersByCombinations(combinations chan []interface{}, src string, dest string,
	sshTimeout time.Duration) (launchers []*Launcher) {
	for com := range combinations {
		launchers = append(launchers, &Launcher{
			SshConnector: launcher.SshConnector{
				Ip:         com[0].(string),
				Port:       com[1].(string),
				User:       com[2].(string),
				Password:   com[3].(string),
				SshTimeout: sshTimeout,
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
	sftpClient, errScp := sftp.NewClient(sshClient, sftp.UseConcurrentWrites(true),
		sftp.UseConcurrentReads(true))
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

	remoteFile, err := sftpClient.Create(sftp.Join(remote, remoteFileName))
	if err != nil {
		utils.Logger.Fatalln(err.Error())
	}
	defer func(remoteFile *sftp.File) {
		err := remoteFile.Close()
		if err != nil {
			utils.Logger.Errorln(err.Error())
		}
	}(remoteFile)

	localFileInfo, _ := localFile.Stat()
	localFileSize := localFileInfo.Size()
	progressBar := pb.New64(localFileSize)
	barReader := progressBar.NewProxyReader(localFile)
	localReader := io.LimitReader(barReader, localFileSize)
	progressBar.Start()
	// Reader must be io.Reader, bytes.Reader or satisfy one of the following interfaces:
	// Len() int, Size() int64, Stat() (os.FileInfo, error).
	// Or the concurrent upload can not work.
	transSize, err := io.Copy(remoteFile, localReader)
	if err != nil {
		utils.Logger.Fatalln(err.Error())
	}
	progressBar.Finish()
	utils.Logger.Infof("File: %s uploaded successfully. Transmission size: %s.\n",
		remoteFileName, utils.ByteSizeFormat(float64(transSize)))
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

	localFile, err := os.Create(sftp.Join(local, localFileName))
	if err != nil {
		utils.Logger.Fatalln(err.Error())
	}
	defer func(localFile *os.File) {
		err := localFile.Close()
		if err != nil {
			utils.Logger.Errorln(err.Error())
		}
	}(localFile)

	remoteFileInfo, _ := remoteFile.Stat()
	remoteFileSize := remoteFileInfo.Size()
	progressBar := pb.New64(remoteFileSize)
	barWriter := progressBar.NewProxyWriter(localFile)
	progressBar.Start()
	transSize, err := io.Copy(barWriter, remoteFile)
	if err != nil {
		utils.Logger.Fatalln(err.Error())
	}
	progressBar.Finish()
	utils.Logger.Infof("File: %s downloaded successfully. Transmission size: %s.\n",
		remoteFileName, utils.ByteSizeFormat(float64(transSize)))
	return true
}
