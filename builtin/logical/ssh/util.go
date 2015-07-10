package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/hashicorp/vault/logical"

	"golang.org/x/crypto/ssh"
)

// Transfers the file  to the target machine by establishing an SSH session with the target.
// Uses the public key authentication method and hence the parameter 'key' takes in the private key.
// The fileName parameter takes an absolute path.
func uploadPublicKeyScp(publicKey, username, ip, port, key string) error {
	dynamicPublicKeyFileName := fmt.Sprintf("vault_ssh_%s_%s.pub", username, ip)
	session, err := createSSHPublicKeysSession(username, ip, port, key)
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("invalid session object")
	}
	defer session.Close()
	go func() {
		w, _ := session.StdinPipe()
		fmt.Fprintln(w, "C0644", len(publicKey), dynamicPublicKeyFileName)
		io.Copy(w, strings.NewReader(publicKey))
		fmt.Fprint(w, "\x00")
		w.Close()
	}()
	err = session.Run(fmt.Sprintf("scp -vt %s", dynamicPublicKeyFileName))
	return nil
}

// Creates a SSH session object which can be used to run commands in the target machine.
// The session will use public key authentication method with port 22.
func createSSHPublicKeysSession(username, ipAddr, port, hostKey string) (*ssh.Session, error) {
	if username == "" {
		return nil, fmt.Errorf("missing username")
	}
	if ipAddr == "" {
		return nil, fmt.Errorf("missing ip address")
	}
	if hostKey == "" {
		return nil, fmt.Errorf("missing host key")
	}
	signer, err := ssh.ParsePrivateKey([]byte(hostKey))
	if err != nil {
		return nil, fmt.Errorf("parsing Private Key failed: %s", err)
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", ipAddr, port), config)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, fmt.Errorf("invalid client object: %s", err)
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

// Creates a new RSA key pair with key length of 2048.
// The private key will be of pem format and the public key will be of OpenSSH format.
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

// Concatenates the public present in that target machine's home folder to ~/.ssh/authorized_keys file
func installPublicKeyInTarget(username, ip, port, hostKey string) error {
	session, err := createSSHPublicKeysSession(username, ip, port, hostKey)
	if err != nil {
		return fmt.Errorf("unable to create SSH Session using public keys: %s", err)
	}
	if session == nil {
		return fmt.Errorf("invalid session object")
	}
	defer session.Close()

	authKeysFileName := fmt.Sprintf("/home/%s/.ssh/authorized_keys", username)
	tempKeysFileName := fmt.Sprintf("/home/%s/temp_authorized_keys", username)

	// Commands to be run on target machine
	dynamicPublicKeyFileName := fmt.Sprintf("vault_ssh_%s_%s.pub", username, ip)
	grepCmd := fmt.Sprintf("grep -vFf %s %s > %s", dynamicPublicKeyFileName, authKeysFileName, tempKeysFileName)
	catCmdRemoveDuplicate := fmt.Sprintf("cat %s > %s", tempKeysFileName, authKeysFileName)
	catCmdAppendNew := fmt.Sprintf("cat %s >> %s", dynamicPublicKeyFileName, authKeysFileName)
	removeCmd := fmt.Sprintf("rm -f %s %s", tempKeysFileName, dynamicPublicKeyFileName)

	targetCmd := fmt.Sprintf("%s;%s;%s;%s", grepCmd, catCmdRemoveDuplicate, catCmdAppendNew, removeCmd)
	session.Run(targetCmd)
	return nil
}

// Removes the installed public key from the authorized_keys file in target machine
func uninstallPublicKeyInTarget(username, ip, port, hostKey string) error {
	session, err := createSSHPublicKeysSession(username, ip, port, hostKey)
	if err != nil {
		return fmt.Errorf("unable to create SSH Session using public keys: %s", err)
	}
	if session == nil {
		return fmt.Errorf("invalid session object")
	}
	defer session.Close()

	authKeysFileName := fmt.Sprintf("/home/%s/.ssh/authorized_keys", username)
	tempKeysFileName := fmt.Sprintf("/home/%s/temp_authorized_keys", username)

	// Commands to be run on target machine
	dynamicPublicKeyFileName := fmt.Sprintf("vault_ssh_%s_%s.pub", username, ip)
	grepCmd := fmt.Sprintf("grep -vFf %s %s > %s", dynamicPublicKeyFileName, authKeysFileName, tempKeysFileName)
	catCmdRemoveDuplicate := fmt.Sprintf("cat %s > %s", tempKeysFileName, authKeysFileName)
	removeCmd := fmt.Sprintf("rm -f %s %s", tempKeysFileName, dynamicPublicKeyFileName)

	remoteCmd := fmt.Sprintf("%s;%s;%s", grepCmd, catCmdRemoveDuplicate, removeCmd)
	session.Run(remoteCmd)
	return nil
}

// Takes an IP address and role name and checks if the IP is part of CIDR blocks belonging to the role.
func roleContainsIP(s logical.Storage, roleName string, ip string) (bool, error) {
	if roleName == "" {
		return false, fmt.Errorf("missing role name")
	}
	if ip == "" {
		return false, fmt.Errorf("missing ip")
	}
	roleEntry, err := s.Get(fmt.Sprintf("policy/%s", roleName))
	if err != nil {
		return false, fmt.Errorf("error retrieving role '%s'", err)
	}
	if roleEntry == nil {
		return false, fmt.Errorf("role '%s' not found", roleName)
	}
	var role sshRole
	if err := roleEntry.DecodeJSON(&role); err != nil {
		return false, fmt.Errorf("error decoding role '%s'", roleName)
	}

	if matched, err := cidrContainsIP(ip, role.CIDR); err != nil {
		return false, err
	} else {
		return matched, nil
	}
}

// Returns true if the IP supplied by the user is part of the comma separated CIDR blocks
func cidrContainsIP(ip, cidr string) (bool, error) {
	for _, item := range strings.Split(cidr, ",") {
		_, cidrIPNet, err := net.ParseCIDR(item)
		if err != nil {
			return false, fmt.Errorf("invalid cidr entry '%s'", item)
		}
		if cidrIPNet.Contains(net.ParseIP(ip)) {
			return true, nil
		}
	}
	return false, nil
}
