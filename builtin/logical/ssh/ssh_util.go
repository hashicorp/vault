package ssh

import (
	"fmt"
	"os/exec"

	"golang.org/x/crypto/ssh"
)

func exec_command(cmdString string) error {
	cmd := exec.Command("/bin/bash", "-c", cmdString)
	if _, err := cmd.Output(); err != nil {
		return err
	}
	return nil
}

func createSSHPublicKeysSession(username string, ipAddr string, hostKey string) *ssh.Session {
	signer, err := ssh.ParsePrivateKey([]byte(hostKey))
	if err != nil {
		fmt.Errorf("Parsing Private Key failed: " + err.Error())
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	client, err := ssh.Dial("tcp", ipAddr+":22", config)
	if err != nil {
		fmt.Errorf("Dial Failed: " + err.Error())
	}
	if client == nil {
		fmt.Errorf("SSH Dial to target failed: ", err.Error())
	}

	session, err := client.NewSession()
	if err != nil {
		fmt.Errorf("NewSession failed: " + err.Error())
	}
	return session
}
