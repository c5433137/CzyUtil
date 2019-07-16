package main
/*
	测试demo：开启sftp 发送、获取文件demo
*/
import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"os"
	"path"
	"time"
)

func main() {
	sendfile()
	fmt.Println("-----------")
	recvfile()
}

func sendfile() {
	var (
		err        error
		sftpClient *sftp.Client
	)

	// 这里换成实际的 SSH 连接的 用户名，密码，主机名或IP，SSH端口
	sftpClient, err = connect("root", "****", "0.0.0.0", 22)
	if err != nil {
		fmt.Println("connect err  = ",err)
		return
	}
	defer sftpClient.Close()

	// 用来测试的本地文件路径 和 远程机器上的文件夹
	var localFilePath = "./test.txt"
	var remoteDir = "/home/sftp_root/"
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer srcFile.Close()

	var remoteFileName = path.Base(localFilePath)
	fmt.Println("remoteFileName=",remoteFileName)
	dstFile, err := sftpClient.Create(path.Join(remoteDir, remoteFileName))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dstFile.Close()

	buf := make([]byte, 1024)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf)
	}

	fmt.Println("copy file to remote server finished!")
}

func recvfile() {
	var (
		err        error
		sftpClient *sftp.Client
	)

	// 这里换成实际的 SSH 连接的 用户名，密码，主机名或IP，SSH端口
	sftpClient, err = connect("root", "****", "0.0.0.0", 22)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sftpClient.Close()

	// 用来测试的远程文件路径 和 本地文件夹
	var remoteFilePath = "/home/sftp_root/test.txt"
	var localDir = "/home/sftp_root"


	srcFile, err := sftpClient.Open(remoteFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer srcFile.Close()

	var localFileName = path.Base(remoteFilePath)
	dstFile, err := os.Create(path.Join(localDir, localFileName))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dstFile.Close()

	dstFile.Chmod(os.ModePerm)

	if _, err = srcFile.WriteTo(dstFile); err != nil {
		fmt.Println(err)
	}

	fmt.Println("copy file from remote server finished!")
}


func connect(user, password, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,

		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}
