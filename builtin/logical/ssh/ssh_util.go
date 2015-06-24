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

func installSshOtkInTarget(session *ssh.Session, username string, ipAddr string) error {
	log.Printf("Vishal: ssh.installSshOtkInTarget\n")

	//TODO: Input validation for the commands below
	otkPrivateKeyFileName := "vault_ssh_" + username + "_" + ipAddr + "_otk.pem"
	otkPublicKeyFileName := otkPrivateKeyFileName + ".pub"
	authKeysFileName := "~/.ssh/authorized_keys"
	tempKeysFileName := "./temp_authorized_keys"

	grepCmd := "grep -vFf " + otkPublicKeyFileName + " " + authKeysFileName + " > " + tempKeysFileName + ";"
	catCmdRemoveDuplicate := "cat " + tempKeysFileName + " > " + authKeysFileName + ";"
	catCmdAppendNew := "cat " + otkPublicKeyFileName + " >> " + authKeysFileName + ";"
	rmCmd := "rm -f " + tempKeysFileName + " " + otkPublicKeyFileName + ";"
	log.Printf("Vishal: grepCmd:%#v\n catCmdRemoveDuplicate:%#v\n catCmdAppendNew:%#v\n rmCmd: %#v\n", grepCmd, catCmdRemoveDuplicate, catCmdAppendNew, rmCmd)
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
	hostKeyFileName := "./vault_ssh_" + username + "_" + ipAddr + "_shared.pem"
	pemBytes, err := ioutil.ReadFile(hostKeyFileName)
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
	if client == nil {
		fmt.Errorf("SSH Dial to target failed: ", err.Error())
	}

	session, err := client.NewSession()
	if err != nil {
		fmt.Errorf("NewSession failed: " + err.Error())
	}
	return session
}
