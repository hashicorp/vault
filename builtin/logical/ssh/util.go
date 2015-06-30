package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

func exec_command(cmdString string) error {
	cmd := exec.Command("/bin/bash", "-c", cmdString)
	if _, err := cmd.Output(); err != nil {
		return err
	}
	return nil
}

func uploadFileScp(fileName, username, ip, key string) error {
	nameBase := filepath.Base(fileName)
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("Unable to open file")
	}
	stat, err := file.Stat()
	if os.IsNotExist(err) {
		return fmt.Errorf("File does not exist")
	}
	session, err := createSSHPublicKeysSession(username, ip, key)
	if err != nil {
		return fmt.Errorf("Unable to create SSH Session using public keys: %s", err)
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
		return fmt.Errorf("Failed to run: %s", err)
	}
	return nil
}

func createSSHPublicKeysSession(username string, ipAddr string, hostKey string) (*ssh.Session, error) {
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
		return nil, fmt.Errorf("Dial Failed: %s", err)
	}
	if client == nil {
		return nil, fmt.Errorf("Invalid client object: %s", err)
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Creating new client session failed: %s", err)
	}
	return session, nil
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
