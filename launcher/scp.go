package launcher

import (
	"github.com/pkg/sftp"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
	"tryssh/utils"
)

// ScpLauncher scp连接器
type ScpLauncher struct {
	SshConnector
	Src  string
	Dest string
}

// Launch 执行
func (c *ScpLauncher) Launch() bool {
	switch {
	case strings.Contains(c.Src, c.Ip):
		return c.download(c.Dest, strings.Split(c.Src, ":")[1])
	case strings.Contains(c.Dest, c.Ip):
		return c.upload(c.Src, strings.Split(c.Dest, ":")[1])
	}
	return false
}

// NewScpLaunchersByCombinations 通过用户、密码、端口以及文件传输源和目标的组合生成ScpLauncher对象切片
func NewScpLaunchersByCombinations(combinations chan []interface{}, src string, dest string) (launchers []ScpLauncher) {
	for com := range combinations {
		launchers = append(launchers, ScpLauncher{
			SshConnector: SshConnector{
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

// createScpClient 创建scp客户端
func (c *ScpLauncher) createScpClient() (sftpClient *sftp.Client) {
	sshClient, errSsh := c.CreateConnection()
	if errSsh != nil {
		return
	}
	sftpClient, errScp := sftp.NewClient(sshClient)
	if errScp != nil {
		log.Fatalln(errScp.Error())
	}
	return
}

// closeScpClient 关闭scp客户端
func (c *ScpLauncher) closeScpClient(sftpClient *sftp.Client) {
	err := sftpClient.Close()
	if err != nil {
		log.Errorln(err.Error())
	}
}

// upload 上传文件
func (c *ScpLauncher) upload(local, remote string) bool {
	sftpClient := c.createScpClient()
	if sftpClient == nil {
		return false
	}
	defer c.closeScpClient(sftpClient)
	localPath := strings.Split(local, "/")
	localFileName := localPath[len(localPath)-1]
	remoteFileName := localFileName
	// 打开远程文件
	remoteFile, err := sftpClient.OpenFile(sftp.Join(remote, remoteFileName), os.O_CREATE|os.O_RDWR|os.O_EXCL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer func(remoteFile *sftp.File) {
		err := remoteFile.Close()
		if err != nil {
			log.Errorln(err.Error())
		}
	}(remoteFile)
	// 打开本地文件
	localFile, err := os.Open(local)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer func(localFile *os.File) {
		err := localFile.Close()
		if err != nil {
			log.Errorln(err.Error())
		}
	}(localFile)
	// 上传
	transSize, err := io.Copy(remoteFile, localFile)
	if err != nil {
		log.Fatalln(err.Error())
	} else {
		log.Infof("File: %s uploaded successfully. Transmission size: %d Bytes.", localFileName, transSize)
	}
	return true
}

// download 下载文件
func (c *ScpLauncher) download(local, remote string) bool {
	sftpClient := c.createScpClient()
	if sftpClient == nil {
		return false
	}
	defer c.closeScpClient(sftpClient)
	remotePath := strings.Split(remote, "/")
	remoteFileName := remotePath[len(remotePath)-1]
	localFileName := remoteFileName
	// 打开本地文件
	localFile, err := os.OpenFile(sftp.Join(local, localFileName), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer func(localFile *os.File) {
		err := localFile.Close()
		if err != nil {
			log.Errorln(err.Error())
		}
	}(localFile)
	// 打开远程文件
	remoteFile, err := sftpClient.Open(remote)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer func(remoteFile *sftp.File) {
		err := remoteFile.Close()
		if err != nil {
			log.Errorln(err.Error())
		}
	}(remoteFile)
	// 下载
	transSize, err := io.Copy(localFile, remoteFile)
	if err != nil {
		log.Fatalln(err.Error())
	} else {
		log.Infof("File: %s downloaded successfully. Transmission size: %s.",
			remoteFileName, utils.ByteSizeFormat(float64(transSize)))
	}
	return true
}
