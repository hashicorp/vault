package ssh

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh"
)

func exec_command(cmdString string) error {
	cmd := exec.Command("/bin/bash", "-c", cmdString)
	if _, err := cmd.Output(); err != nil {
		return err
	}
	return nil
}

func installSshOtkInTarget(session *ssh.Session) error {
	log.Printf("Vishal: ssh.installSshOtkInTarget\n")

	grepCmd := "grep -vFf " + "vault_ssh_otk.pem.pub" + " " + "~/.ssh/authorized_keys" + " > " + "./temp_authorized_keys" + ";"
	catCmdRemoveDuplicate := "cat " + "./temp_authorized_keys" + " > " + "~/.ssh/authorized_keys" + ";"
	catCmdAppendNew := "cat " + "./vault_ssh_otk.pem.pub" + " >> " + "~/.ssh/authorized_keys" + ";"
	rmCmd := "rm -f " + "./temp_authorized_keys" + " " + "./vault_ssh_otk.pem.pub" + ";"
	remoteCmdString := strings.Join([]string{
		grepCmd,
		catCmdRemoveDuplicate,
		catCmdAppendNew,
		rmCmd,
	}, "")

	if err := session.Run(remoteCmdString); err != nil {
		return err
	}
	return nil
}
func createSSHPublicKeysSession(username string, ipAddr string) *ssh.Session {
	pemBytes, err := ioutil.ReadFile("vault_ssh_shared.pem")
	if err != nil {
		fmt.Errorf("Reading shared key failed: " + err.Error())
	}

	signer, err := ssh.ParsePrivateKey(pemBytes)
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

	session, err := client.NewSession()
	if err != nil {
		fmt.Errorf("NewSession failed: " + err.Error())
	}
	return session
}
