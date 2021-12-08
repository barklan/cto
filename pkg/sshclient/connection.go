package sshclient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/barklan/cto/pkg/storage"
	ssh "golang.org/x/crypto/ssh"
)

var port = "22"

func executeCmd(command, hostname string, config *ssh.ClientConfig) (bytes.Buffer, error) {
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", hostname, port), config)
	if err != nil {
		return bytes.Buffer{}, err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return bytes.Buffer{}, err
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run(command)

	return stdoutBuf, nil
}

func PublicKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}

func ConnectAndExecute(
	data *storage.Data,
	projectName string,
	cmd string,
) (bytes.Buffer, error) {

	// TODO this may cause runtime panic
	baseFileName := data.Config.P[projectName].Backups.DB.SSHKeyFilename
	var keyFilename string

	if _, ok := os.LookupEnv("CTO_LOCAL_ENV"); ok {
		keyFilename = fmt.Sprintf("environment/%s", baseFileName)
	} else {
		keyFilename = fmt.Sprintf("/app/config/%s", baseFileName)
	}

	publicKey, err := PublicKeyFile(keyFilename)
	if err != nil {
		return bytes.Buffer{}, err
	}

	// TODO this may cause runtime panic
	user := data.Config.P[projectName].Backups.DB.SSHUser
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			// ssh.Password(p),
			publicKey,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// TODO this may cause runtime panic
	hostname := data.Config.P[projectName].Backups.DB.SSHHostname
	output, err := executeCmd(cmd, hostname, config)
	if err != nil {
		return bytes.Buffer{}, err
	}

	return output, nil
}
