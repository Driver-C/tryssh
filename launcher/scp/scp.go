package scp

import (
	"github.com/cheggaaa/pb/v3"
	"github.com/pkg/sftp"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"tryssh/launcher"
	"tryssh/utils"
)

type Launcher struct {
	launcher.SshConnector
	Src       string
	Dest      string
	Recursive bool
}

func (c *Launcher) Launch() bool {
	sftpClient := c.createScpClient()
	if sftpClient == nil {
		return false
	}
	defer c.closeScpClient(sftpClient)

	// Replace ~ to the real home directory
	c.replaceHomeDirSymbol(sftpClient)

	switch {
	case strings.Contains(c.Src, c.Ip) && !c.Recursive:
		return c.download(c.Dest, strings.Split(c.Src, ":")[1], sftpClient)
	case strings.Contains(c.Src, c.Ip) && c.Recursive:
		return c.downloadDir(c.Dest, strings.Split(c.Src, ":")[1], sftpClient)
	case strings.Contains(c.Dest, c.Ip) && !c.Recursive:
		return c.upload(c.Src, strings.Split(c.Dest, ":")[1], sftpClient)
	case strings.Contains(c.Dest, c.Ip) && c.Recursive:
		return c.uploadDir(c.Src, strings.Split(c.Dest, ":")[1], sftpClient)
	}

	return false
}

func (c *Launcher) replaceHomeDirSymbol(sftpClient *sftp.Client) {
	remoteHomeDir, err := sftpClient.Getwd()
	if err != nil {
		utils.Logger.Fatalf("Failed to get home directory: %v", err)
	}
	homeDirSymbol := "~"
	c.Src = strings.Replace(c.Src, homeDirSymbol, remoteHomeDir, -1)
	c.Dest = strings.Replace(c.Dest, homeDirSymbol, remoteHomeDir, -1)
}

func NewScpLaunchersByCombinations(combinations chan []interface{}, src string, dest string,
	recursive bool, sshTimeout time.Duration) (launchers []*Launcher) {
	for com := range combinations {
		launchers = append(launchers, &Launcher{
			SshConnector: launcher.SshConnector{
				Ip:         com[0].(string),
				Port:       com[1].(string),
				User:       com[2].(string),
				Password:   com[3].(string),
				SshTimeout: sshTimeout,
			},
			Src:       src,
			Dest:      dest,
			Recursive: recursive,
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

func (c *Launcher) upload(local, remote string, sftpClient *sftp.Client) bool {
	localPathSegments := strings.Split(local, "/")
	localFileName := localPathSegments[len(localPathSegments)-1]
	// Openssh scp options rule imitation
	var remoteFileName string
	if strings.HasSuffix(remote, "/") {
		remoteFileName = localFileName
	}
	prefix := local + " "

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

	localFileInfo, err := localFile.Stat()
	if err != nil {
		utils.Logger.Errorln("Get local file stat failed: ", err)
		return false
	}
	localFileSize := localFileInfo.Size()
	localFilePerm := localFileInfo.Mode().Perm()
	// Sync file permission
	if err := remoteFile.Chmod(localFilePerm); err != nil {
		utils.Logger.Errorln("Sync file permission failed: ", err)
		return false
	}
	progressBar := pb.New64(localFileSize)
	progressBar.Set("prefix", prefix)
	barReader := progressBar.NewProxyReader(localFile)
	localReader := io.LimitReader(barReader, localFileSize)
	progressBar.Start()
	// Reader must be io.Reader, bytes.Reader or satisfy one of the following interfaces:
	// Len() int, Size() int64, Stat() (os.FileInfo, error).
	// Or the concurrent upload can not work.
	_, err = io.Copy(remoteFile, localReader)
	if err != nil {
		utils.Logger.Fatalln(err.Error())
	}
	progressBar.Finish()
	return true
}

func (c *Launcher) uploadDir(local, remote string, sftpClient *sftp.Client) bool {
	// Openssh scp options rule imitation
	if strings.HasSuffix(remote, "/") {
		remote = filepath.Join(remote, filepath.Base(local))
	}

	// Create remote root directory
	if err := sftpClient.MkdirAll(remote); err != nil {
		utils.Logger.Errorln("Unable to create remote directory: ", err)
		return false
	}
	entries, err := os.ReadDir(local)
	if err != nil {
		utils.Logger.Errorln(err.Error())
		return false
	}
	for _, entry := range entries {
		localPath := filepath.Join(local, entry.Name())
		remotePath := filepath.Join(remote, entry.Name())
		if entry.IsDir() {
			// Create remote directory
			if err := sftpClient.MkdirAll(remotePath); err != nil {
				utils.Logger.Errorln("Unable to create remote directory: ", err)
				return false
			}
			c.uploadDir(localPath, remotePath, sftpClient)
		} else {
			c.upload(localPath, remotePath, sftpClient)
		}
	}
	return true
}

func (c *Launcher) download(local, remote string, sftpClient *sftp.Client) bool {
	remotePath := strings.Split(remote, "/")
	remoteFileName := remotePath[len(remotePath)-1]
	// Openssh scp options rule imitation
	var localFileName string
	if strings.HasSuffix(local, "/") {
		localFileName = remoteFileName
	}
	prefix := remote + " "

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

	remoteFileInfo, err := remoteFile.Stat()
	if err != nil {
		utils.Logger.Errorln("Get remote file stat failed: ", err)
		return false
	}
	remoteFilePerm := remoteFileInfo.Mode().Perm()
	// Sync file permission
	if err := localFile.Chmod(remoteFilePerm); err != nil {
		utils.Logger.Errorln("Sync file permission failed: ", err)
		return false
	}
	remoteFileSize := remoteFileInfo.Size()
	progressBar := pb.New64(remoteFileSize)
	progressBar.Set("prefix", prefix)
	barWriter := progressBar.NewProxyWriter(localFile)
	progressBar.Start()
	_, err = io.Copy(barWriter, remoteFile)
	if err != nil {
		utils.Logger.Fatalln(err.Error())
	}
	progressBar.Finish()
	return true
}

func (c *Launcher) downloadDir(local, remote string, sftpClient *sftp.Client) bool {
	// Openssh scp options rule imitation
	if strings.HasSuffix(local, "/") {
		local = filepath.Join(local, filepath.Base(remote))
	}

	// Create local root directory
	if err := os.MkdirAll(local, 0755); err != nil {
		utils.Logger.Errorln("Unable to create local directory: ", err)
		return false
	}
	entries, err := sftpClient.ReadDir(remote)
	if err != nil {
		utils.Logger.Errorln(err.Error())
		return false
	}
	for _, entry := range entries {
		localPath := filepath.Join(local, entry.Name())
		remotePath := filepath.Join(remote, entry.Name())
		if entry.IsDir() {
			// Create local directory
			if err := os.MkdirAll(localPath, 0755); err != nil {
				utils.Logger.Errorln("Unable to create local directory: ", err)
				return false
			}
			c.downloadDir(localPath, remotePath, sftpClient)
		} else {
			c.download(localPath, remotePath, sftpClient)
		}
	}
	return true
}
