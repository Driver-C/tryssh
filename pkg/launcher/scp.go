package launcher

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Driver-C/tryssh/pkg/utils"
	"github.com/cheggaaa/pb/v3"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// ScpLauncher handles SCP file transfer operations over SSH.
type ScpLauncher struct {
	SSHConnector
	Src       string
	Dest      string
	Recursive bool
}

// Launch performs the SCP file transfer and returns true on success.
func (c *ScpLauncher) Launch() bool {
	sftpClient, sshClient, err := c.createScpClient()
	if err != nil || sftpClient == nil {
		return false
	}
	defer c.closeScpClient(sftpClient, sshClient)

	// Determine direction before tilde expansion so host prefix is intact.
	isDownload := hasHostPrefix(c.Src, c.IP)
	isUpload := hasHostPrefix(c.Dest, c.IP)

	c.replaceHomeDirPrefix(sftpClient, isDownload)

	switch {
	case isDownload:
		return c.downloadWildcards(c.Dest, splitRemotePath(c.Src), sftpClient, c.Recursive)
	case isUpload:
		return c.uploadWildcards(c.Src, splitRemotePath(c.Dest), sftpClient, c.Recursive)
	default:
		utils.Errorln("Cannot determine upload or download direction: no IP found in source or destination")
		return false
	}
}

// hasHostPrefix checks whether s starts with "host:" or "[host]:" for the given host.
func hasHostPrefix(s, host string) bool {
	if strings.HasPrefix(s, host+":") {
		return true
	}
	return strings.HasPrefix(s, "["+host+"]:")
}

// splitRemotePath splits "host:path" or "[host]:path" and returns only the path part.
func splitRemotePath(s string) string {
	// Handle [ipv6]:path format
	if strings.HasPrefix(s, "[") {
		closeBracket := strings.Index(s, "]")
		if closeBracket >= 0 && closeBracket+1 < len(s) && s[closeBracket+1] == ':' {
			return s[closeBracket+2:]
		}
		return s
	}
	// Handle host:path format
	parts := strings.SplitN(s, ":", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return s
}

// replaceHomeDirPrefix replaces a leading "~/" prefix with the remote home directory.
// Only the remote side (download source or upload destination) gets expanded.
// Paths are in "host:path" format, so the host prefix must be stripped before checking for ~.
func (c *ScpLauncher) replaceHomeDirPrefix(sftpClient *sftp.Client, isDownload bool) {
	remoteHomeDir, err := sftpClient.Getwd()
	if err != nil {
		utils.Errorf("Failed to get remote home directory: %v", err)
		return
	}
	if isDownload {
		// Download: remote = Src — expand ~ in the path portion
		c.Src = expandTildeInRemotePath(c.Src, remoteHomeDir)
	} else {
		// Upload: remote = Dest — expand ~ in the path portion
		c.Dest = expandTildeInRemotePath(c.Dest, remoteHomeDir)
	}
}

// expandTildeInRemotePath replaces ~/ in the path portion of a "host:path" string.
func expandTildeInRemotePath(s, homeDir string) string {
	tilde := "~/"
	// Handle [host]:path format
	if strings.HasPrefix(s, "[") {
		closeBracket := strings.Index(s, "]")
		if closeBracket < 0 {
			return s
		}
		hostPrefix := s[:closeBracket+1]
		pathPart := s[closeBracket+1:]
		pathPart = strings.TrimPrefix(pathPart, ":")
		if strings.HasPrefix(pathPart, tilde) {
			return hostPrefix + ":" + homeDir + pathPart[1:]
		}
		return s
	}
	// Handle host:path format
	parts := strings.SplitN(s, ":", 2)
	if len(parts) == 2 {
		if strings.HasPrefix(parts[1], tilde) {
			return parts[0] + ":" + homeDir + parts[1][1:]
		}
	}
	// No host prefix — treat entire string as path
	if strings.HasPrefix(s, tilde) {
		return homeDir + s[1:]
	}
	return s
}

// NewScpLaunchersByCombinations creates ScpLauncher instances from a channel of credential combinations.
func NewScpLaunchersByCombinations(combinations chan []interface{}, src string, dest string,
	recursive bool, sshTimeout time.Duration) (launchers []*ScpLauncher) {
	for com := range combinations {
		ip, _ := com[0].(string)
		port, _ := com[1].(string)
		user, _ := com[2].(string)
		password, _ := com[3].(string)
		key, _ := com[4].(string)
		launchers = append(launchers, &ScpLauncher{
			SSHConnector: SSHConnector{
				IP:         ip,
				Port:       port,
				User:       user,
				Password:   password,
				Key:        key,
				SSHTimeout: sshTimeout,
			},
			Src:       src,
			Dest:      dest,
			Recursive: recursive,
		})
	}
	return
}

func (c *ScpLauncher) createScpClient() (*sftp.Client, *ssh.Client, error) {
	sshClient, errSSH := c.CreateConnection()
	if errSSH != nil {
		return nil, nil, errSSH
	}
	sftpClient, errScp := sftp.NewClient(sshClient, sftp.UseConcurrentWrites(true),
		sftp.UseConcurrentReads(true))
	if errScp != nil {
		_ = sshClient.Close()
		return nil, nil, fmt.Errorf("SFTP client creation failed: %w", errScp)
	}
	return sftpClient, sshClient, nil
}

func (c *ScpLauncher) closeScpClient(sftpClient *sftp.Client, sshClient *ssh.Client) {
	if err := sftpClient.Close(); err != nil {
		utils.Errorln(err.Error())
	}
	if err := sshClient.Close(); err != nil {
		utils.Errorln(err.Error())
	}
}

// uploadWildcards expands local glob patterns and uploads matching files.
func (c *ScpLauncher) uploadWildcards(local, remote string, sftpClient *sftp.Client, recursive bool) bool {
	matches, err := filepath.Glob(local)
	if err != nil {
		utils.Errorf("Invalid glob pattern %q: %v", local, err)
		return false
	}
	if len(matches) == 0 {
		utils.Errorf("No files match pattern %q", local)
		return false
	}

	allOk := true
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			utils.Errorf("Cannot stat %q: %v", match, err)
			allOk = false
			continue
		}
		if info.IsDir() {
			if !recursive {
				utils.Warnf("Skipping directory %q (use -r for recursive)", match)
				continue
			}
			if !c.uploadDir(match, remote, sftpClient) {
				allOk = false
			}
		} else if !c.upload(match, remote, sftpClient) {
			allOk = false
		}
	}
	return allOk
}

// downloadWildcards expands remote glob patterns and downloads matching files.
func (c *ScpLauncher) downloadWildcards(local, remote string, sftpClient *sftp.Client, recursive bool) bool {
	matches, err := sftpClient.Glob(remote)
	if err != nil {
		utils.Errorf("Invalid remote glob pattern %q: %v", remote, err)
		return false
	}
	if len(matches) == 0 {
		utils.Errorf("No remote files match pattern %q", remote)
		return false
	}

	allOk := true
	for _, match := range matches {
		info, err := sftpClient.Stat(match)
		if err != nil {
			utils.Errorf("Cannot stat remote %q: %v", match, err)
			allOk = false
			continue
		}
		if info.IsDir() {
			if !recursive {
				utils.Warnf("Skipping remote directory %q (use -r for recursive)", match)
				continue
			}
			if !c.downloadDir(local, match, sftpClient) {
				allOk = false
			}
		} else if !c.download(local, match, sftpClient) {
			allOk = false
		}
	}
	return allOk
}

func (c *ScpLauncher) upload(local, remote string, sftpClient *sftp.Client) bool {
	localFileName := filepath.Base(local)

	var targetPath string
	if strings.HasSuffix(remote, "/") {
		targetPath = sftp.Join(remote, localFileName)
	} else {
		targetPath = remote
	}
	prefix := local + " "

	localFile, err := os.Open(local) //nolint:gosec // G304: local path from user input
	if err != nil {
		utils.Errorln(err.Error())
		return false
	}
	defer func(localFile *os.File) {
		if closeErr := localFile.Close(); closeErr != nil {
			utils.Errorln(closeErr.Error())
		}
	}(localFile)

	remoteFile, err := sftpClient.Create(targetPath)
	if err != nil {
		utils.Errorln(err.Error())
		return false
	}
	defer func(remoteFile *sftp.File) {
		if closeErr := remoteFile.Close(); closeErr != nil {
			utils.Errorln(closeErr.Error())
		}
	}(remoteFile)

	localFileInfo, err := localFile.Stat()
	if err != nil {
		utils.Errorln("Get local file stat failed: ", err)
		return false
	}
	localFileSize := localFileInfo.Size()
	localFilePerm := localFileInfo.Mode().Perm()
	if chmodErr := remoteFile.Chmod(localFilePerm); chmodErr != nil {
		utils.Errorln("Sync file permission failed: ", chmodErr)
		return false
	}
	progressBar := pb.New64(localFileSize)
	progressBar.Set("prefix", prefix)
	barReader := progressBar.NewProxyReader(localFile)
	localReader := io.LimitReader(barReader, localFileSize)
	progressBar.Start()
	if _, copyErr := io.Copy(remoteFile, localReader); copyErr != nil {
		utils.Errorln(copyErr.Error())
		progressBar.Finish()
		return false
	}
	progressBar.Finish()
	return true
}

func (c *ScpLauncher) uploadDir(local, remote string, sftpClient *sftp.Client) bool {
	if strings.HasSuffix(remote, "/") {
		remote = sftp.Join(remote, filepath.Base(local))
	}

	if mkdirErr := sftpClient.MkdirAll(remote); mkdirErr != nil {
		utils.Errorln("Unable to create remote directory: ", mkdirErr)
		return false
	}
	entries, err := os.ReadDir(local)
	if err != nil {
		utils.Errorln(err.Error())
		return false
	}
	for _, entry := range entries {
		localPath := filepath.Join(local, entry.Name())
		remotePath := sftp.Join(remote, entry.Name())
		if entry.IsDir() {
			if !c.uploadDir(localPath, remotePath, sftpClient) {
				return false
			}
		} else {
			if !c.upload(localPath, remotePath, sftpClient) {
				return false
			}
		}
	}
	return true
}

func (c *ScpLauncher) download(local, remote string, sftpClient *sftp.Client) bool {
	remoteFileName := filepath.Base(remote)

	var targetPath string
	if strings.HasSuffix(local, "/") {
		targetPath = filepath.Join(local, remoteFileName)
	} else {
		targetPath = local
	}
	prefix := remote + " "

	remoteFile, err := sftpClient.Open(remote)
	if err != nil {
		utils.Errorln(err.Error())
		return false
	}
	defer func(remoteFile *sftp.File) {
		if closeErr := remoteFile.Close(); closeErr != nil {
			utils.Errorln(closeErr.Error())
		}
	}(remoteFile)

	remoteFileInfo, err := remoteFile.Stat()
	if err != nil {
		utils.Errorln("Get remote file stat failed: ", err)
		return false
	}
	remoteFileSize := remoteFileInfo.Size()
	remoteFilePerm := remoteFileInfo.Mode().Perm()

	// Write to a temporary file first to avoid truncating the target on failure.
	tmpFile, err := os.CreateTemp(filepath.Dir(targetPath), ".tryssh-dl-*")
	if err != nil {
		utils.Errorln("Failed to create temp file: ", err)
		return false
	}
	tmpPath := tmpFile.Name()

	success := false
	defer func() {
		_ = tmpFile.Close()
		if !success {
			_ = os.Remove(tmpPath)
		}
	}()

	if chmodErr := tmpFile.Chmod(remoteFilePerm); chmodErr != nil {
		utils.Errorln("Sync file permission failed: ", chmodErr)
		return false
	}

	progressBar := pb.New64(remoteFileSize)
	progressBar.Set("prefix", prefix)
	barWriter := progressBar.NewProxyWriter(tmpFile)
	progressBar.Start()
	if _, copyErr := io.Copy(barWriter, io.LimitReader(remoteFile, remoteFileSize)); copyErr != nil {
		utils.Errorln(copyErr.Error())
		progressBar.Finish()
		return false
	}
	progressBar.Finish()

	if renameErr := os.Rename(tmpPath, targetPath); renameErr != nil {
		utils.Errorln("Failed to rename temp file: ", renameErr)
		return false
	}
	success = true
	return true
}

func (c *ScpLauncher) downloadDir(local, remote string, sftpClient *sftp.Client) bool {
	if strings.HasSuffix(local, "/") {
		local = filepath.Join(local, filepath.Base(remote))
	}

	if mkdirErr := os.MkdirAll(local, 0700); mkdirErr != nil {
		utils.Errorln("Unable to create local directory: ", mkdirErr)
		return false
	}
	entries, err := sftpClient.ReadDir(remote)
	if err != nil {
		utils.Errorln(err.Error())
		return false
	}
	for _, entry := range entries {
		localPath := filepath.Join(local, entry.Name())
		remotePath := sftp.Join(remote, entry.Name())
		if entry.IsDir() {
			if !c.downloadDir(localPath, remotePath, sftpClient) {
				return false
			}
		} else {
			if !c.download(localPath, remotePath, sftpClient) {
				return false
			}
		}
	}
	return true
}
