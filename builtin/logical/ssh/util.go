package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"os"
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

func removeFile(fileName string) {
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("Error fetching working directory:%s", err)
		return
	}
	absFileName := wd + "/" + fileName

	if _, err := os.Stat(absFileName); err == nil {
		err := os.Remove(absFileName)
		if err != nil {
			log.Printf(fmt.Sprintf("Failed: %s", err))
			return
		} else {
			log.Printf("Successful\n")
		}
	}
}

func generateRSAKeys() (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", fmt.Errorf("error generating RSA key-pair: %s", err)
	}

	privateKeyRsa := string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}))

	sshPublicKey, err := ssh.NewPublicKey(privateKey.Public())
	if err != nil {
		return "", "", fmt.Errorf("error generating RSA key-pair: %s", err)
	}
	publicKeyRsa := "ssh-rsa " + base64.StdEncoding.EncodeToString(sshPublicKey.Marshal())

	//ioutil.WriteFile("testkey.pem", []byte(privateKeyRsa), 0600)
	//ioutil.WriteFile("testkey.pub", []byte(publicKeyRsa), 0600)

	return publicKeyRsa, privateKeyRsa, nil
}
