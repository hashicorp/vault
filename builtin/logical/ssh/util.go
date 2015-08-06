package ssh

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"

	commssh "github.com/mitchellh/packer/communicator/ssh"
	"golang.org/x/crypto/ssh"
)

// Transfers the file  to the target machine by establishing an SSH
// session with the target. Uses the public key authentication method
// and hence the parameter 'key' takes in the private key. The fileName
// parameter takes an absolute path.
func uploadPublicKeyScp(publicKey, publicKeyFileName, username, ip, port, key string) error {
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
		fmt.Fprintln(w, "C0644", len(publicKey), publicKeyFileName)
		io.Copy(w, strings.NewReader(publicKey))
		fmt.Fprint(w, "\x00")
		w.Close()
	}()
	session.Run(fmt.Sprintf("scp -vt %s", publicKeyFileName))
	return nil
}

// Creates a SSH session object which can be used to run commands
// in the target machine. The session will use public key authentication
// method with port 22.
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
// The private key will be of pem format and the public key will be
// of OpenSSH format.
func generateRSAKeys(keyBits int) (publicKeyRsa string, privateKeyRsa string, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, keyBits)
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

// Concatenates the public present in that target machine's home
// folder to ~/.ssh/authorized_keys file
func installPublicKeyInTarget(adminUser, publicKeyFileName, username, ip, port, hostkey string) error {
	session, err := createSSHPublicKeysSession(adminUser, ip, port, hostkey)
	if err != nil {
		return fmt.Errorf("unable to create SSH Session using public keys: %s", err)
	}
	if session == nil {
		return fmt.Errorf("invalid session object")
	}
	defer session.Close()

	authKeysFileName := fmt.Sprintf("/home/%s/.ssh/authorized_keys", username)
	scriptFileName := fmt.Sprintf("%s.sh", publicKeyFileName)

	// Give execute permissions to install script, run and delete it.
	chmodCmd := fmt.Sprintf("chmod +x %s", scriptFileName)
	scriptCmd := fmt.Sprintf("./%s install %s %s", scriptFileName, publicKeyFileName, authKeysFileName)
	rmCmd := fmt.Sprintf("rm -f %s", scriptFileName)
	targetCmd := fmt.Sprintf("%s;%s;%s", chmodCmd, scriptCmd, rmCmd)

	session.Run(targetCmd)
	return nil
}

// Removes the installed public key from the authorized_keys file
// in target machine
func uninstallPublicKeyInTarget(adminUser, publicKeyFileName, username, ip, port, hostKey string) error {
	session, err := createSSHPublicKeysSession(adminUser, ip, port, hostKey)
	if err != nil {
		return fmt.Errorf("unable to create SSH Session using public keys: %s", err)
	}
	if session == nil {
		return fmt.Errorf("invalid session object")
	}
	defer session.Close()

	authKeysFileName := fmt.Sprintf("/home/%s/.ssh/authorized_keys", username)
	scriptFileName := fmt.Sprintf("%s.sh", publicKeyFileName)

	// Give execute permissions to install script, run and delete it.
	chmodCmd := fmt.Sprintf("chmod +x %s", scriptFileName)
	scriptCmd := fmt.Sprintf("./%s uninstall %s %s", scriptFileName, publicKeyFileName, authKeysFileName)
	rmCmd := fmt.Sprintf("rm -f %s", scriptFileName)
	targetCmd := fmt.Sprintf("%s;%s;%s", chmodCmd, scriptCmd, rmCmd)

	session.Run(targetCmd)

	return nil
}

// Takes an IP address and role name and checks if the IP is part
// of CIDR blocks belonging to the role.
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

// Returns true if the IP supplied by the user is part of the comma
// separated CIDR blocks
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

func scpUpload(username, ip, port, hostkey, fileName, fileContent string) error {
	signer, err := ssh.ParsePrivateKey([]byte(hostkey))
	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	connfunc := func() (net.Conn, error) {
		c, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", ip, port), 15*time.Second)
		if err != nil {
			return nil, err
		}

		if tcpConn, ok := c.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(5 * time.Second)
		}

		return c, nil
	}
	config := &commssh.Config{
		SSHConfig:    clientConfig,
		Connection:   connfunc,
		Pty:          false,
		DisableAgent: true,
	}
	comm, err := commssh.New(fmt.Sprintf("%s:%s", ip, port), config)
	if err != nil {
		return fmt.Errorf("error connecting to target: %s", err)
	}
	comm.Upload(fileName, bytes.NewBufferString(fileContent), nil)
	return nil
}
