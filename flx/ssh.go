package flx

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"time"
	"golang.org/x/crypto/ssh"
)

func NewSshClient(user, ip, path string) (*ssh.Client, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("ssh key file read failed", err)
	}
	signer, err := ssh.ParsePrivateKey(bytes)
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 5,
		User:            user,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}


	addr := fmt.Sprintf("%s:%d", ip, 22)
	c, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func runCommand(client *ssh.Client, command string) (stdout string, err error) {
	session, err := client.NewSession()
	if err != nil {
		//log.Print(err)
		return
	}
	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf
	err = session.Run(command)
	if err != nil {
		//log.Print(err)
		return
	}
	stdout = string(buf.Bytes())

	return
}