package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

/*
Executes the command represented by the input.
Multiple commands can be run by concatinating strings with ';'.
Currently, it is supported only for linux platforms and user bash shell.
*/
func exec_command(cmdString string) error {
	cmd := exec.Command("/bin/bash", "-c", cmdString)
	if _, err := cmd.Output(); err != nil {
		return err
	}
	return nil
}

/*
Transfers the file  to the target machine by establishing an SSH session with the target.
Uses the public key authentication method and hence the parameter 'key' takes in the private key.
The fileName parameter takes an absolute path.
*/
func uploadFileScp(fileName, username, ip, key string) error {
	nameBase := filepath.Base(fileName)
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	stat, err := file.Stat()
	if os.IsNotExist(err) {
		return fmt.Errorf("File does not exist")
	}
	session, err := createSSHPublicKeysSession(username, ip, key)
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("Invalid session object")
	}
	defer session.Close()
	go func() {
		w, _ := session.StdinPipe()
		fmt.Fprintln(w, "C0644", stat.Size(), nameBase)
		io.Copy(w, file)
		fmt.Fprint(w, "\x00")
		w.Close()
	}()
	if err := session.Run(fmt.Sprintf("scp -vt %s", nameBase)); err != nil {
		return err
	}
	return nil
}

/*
Creates a SSH session object which can be used to run commands in the target machine.
The session will use public key authentication method with port 22.
*/
func createSSHPublicKeysSession(username, ipAddr, hostKey string) (*ssh.Session, error) {
	if username == "" || ipAddr == "" || hostKey == "" {
		return nil, fmt.Errorf("Invalid parameters")
	}
	signer, err := ssh.ParsePrivateKey([]byte(hostKey))
	if err != nil {
		return nil, fmt.Errorf("Parsing Private Key failed: %s", err)
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	client, err := ssh.Dial("tcp", ipAddr+":22", config)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, fmt.Errorf("Invalid client object: %s", err)
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

/*
Deletes the file in the current directory.
The parameter is just the name of the file and not a path.
*/
func removeFile(fileName string) error {
	if fileName == "" {
		return fmt.Errorf("Invalid file name")
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	absFileName := wd + "/" + fileName

	if _, err := os.Stat(absFileName); err == nil {
		err := os.Remove(absFileName)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
Creates a new RSA key pair with key length of 2048.
The private key will be of pem format and the public key will be of OpenSSH format.
*/
func generateRSAKeys() (publicKeyRsa string, privateKeyRsa string, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", fmt.Errorf("error generating RSA key-pair: %s", err)
	}

	privateKeyRsa = string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}))

	sshPublicKey, err := ssh.NewPublicKey(privateKey.Public())
	if err != nil {
		return "", "", fmt.Errorf("error generating RSA key-pair: %s", err)
	}
	publicKeyRsa = "ssh-rsa " + base64.StdEncoding.EncodeToString(sshPublicKey.Marshal())
	return
}
