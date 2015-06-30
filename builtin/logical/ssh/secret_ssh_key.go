package ssh

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const SecretOneTimeKeyType = "secret_one_type_key_type"

func secretSshKey(b *backend) *framework.Secret {
	log.Printf("Vishal: ssh.secretPrivateKey\n")
	return &framework.Secret{
		Type: SecretOneTimeKeyType,
		Fields: map[string]*framework.FieldSchema{
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Username in host",
			},
			"ip": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "ip address of host",
			},
		},
		DefaultDuration:    10 * time.Second, //TODO: change this
		DefaultGracePeriod: 10 * time.Second, //TODO: change this
		Renew:              b.secretSshKeyRenew,
		Revoke:             b.secretSshKeyRevoke,
	}
}

func (b *backend) secretSshKeyRenew(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	log.Printf("Vishal: ssh.secretPrivateKeyRenew\n")
	lease, err := b.Lease(req.Storage)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		lease = &configLease{Lease: 1 * time.Hour}
	}
	f := framework.LeaseExtend(lease.Lease, lease.LeaseMax, false)
	return f(req, d)
}

func (b *backend) secretSshKeyRevoke(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	log.Printf("Vishal: ssh.secretPrivateKeyRevoke req: %#v\n", req)
	//fetch the values from secret
	usernameRaw, ok := req.Secret.InternalData["username"]
	if !ok {
		return nil, fmt.Errorf("secret is missing internal data")
	}
	username, ok := usernameRaw.(string)
	if !ok {
		return nil, fmt.Errorf("secret is missing internal data")
	}
	ipRaw, ok := req.Secret.InternalData["ip"]
	if !ok {
		return nil, fmt.Errorf("secret is missing internal data")
	}
	ip, ok := ipRaw.(string)
	if !ok {
		return nil, fmt.Errorf("secret is missing internal data")
	}
	hostKeyNameRaw, ok := req.Secret.InternalData["host_key_name"]
	if !ok {
		return nil, fmt.Errorf("secret is missing internal data")
	}
	hostKeyName := hostKeyNameRaw.(string)
	if !ok {
		return nil, fmt.Errorf("secret is missing internal data")
	}
	dynamicPublicKeyRaw, ok := req.Secret.InternalData["dynamic_public_key"]
	if !ok {
		return nil, fmt.Errorf("secret is missing internal data")
	}
	dynamicPublicKey := dynamicPublicKeyRaw.(string)
	if !ok {
		return nil, fmt.Errorf("secret is missing internal data")
	}
	log.Printf("Vishal: username:%s ip:%s keyName:%s\n", username, ip, hostKeyName)

	//fetch the host key using the key name
	hostKeyPath := "keys/" + hostKeyName
	hostKeyEntry, err := req.Storage.Get(hostKeyPath)
	if err != nil {
		return nil, fmt.Errorf("Key '%s' not found error:%s", hostKeyName, err)
	}
	var hostKey sshHostKey
	if err := hostKeyEntry.DecodeJSON(&hostKey); err != nil {
		return nil, fmt.Errorf("Error reading the host key: %s", err)
	}

	//write host key to file and use it as argument to scp command
	hostKeyFileName := "./vault_ssh_" + username + "_" + ip + "_shared.pem"
	err = ioutil.WriteFile(hostKeyFileName, []byte(hostKey.Key), 0400)

	//write dynamicPublicKey to file and use it as argument to scp command
	otkPrivateKeyFileName := "vault_ssh_" + username + "_" + ip + "_otk.pem"
	otkPublicKeyFileName := otkPrivateKeyFileName + ".pub"
	err = ioutil.WriteFile(otkPublicKeyFileName, []byte(dynamicPublicKey), 0400)

	//transfer the dynamic public key to target machine and use it to remove the entry from authorized_keys file
	scpCmd := "scp -i " + hostKeyFileName + " " + otkPublicKeyFileName + " " + username + "@" + ip + ":~;"
	err = exec_command(scpCmd)
	if err != nil {
		fmt.Errorf("Running command scp failed " + err.Error())
	}

	authKeysFileName := "~/.ssh/authorized_keys"
	tempKeysFileName := "~/temp_authorized_keys"

	//commands to be run on target machine
	grepCmd := "grep -vFf " + otkPublicKeyFileName + " " + authKeysFileName + " > " + tempKeysFileName + ";"
	catCmdRemoveDuplicate := "cat " + tempKeysFileName + " > " + authKeysFileName + ";"
	rmCmd := "rm -f " + tempKeysFileName + " " + otkPublicKeyFileName + ";"
	remoteCmdString := strings.Join([]string{
		grepCmd,
		catCmdRemoveDuplicate,
		rmCmd,
	}, "")

	//connect to target machine
	session, err := createSSHPublicKeysSession(username, ip, hostKey.Key)
	if err != nil {
		return nil, fmt.Errorf("Unable to create SSH Session using public keys: %s", err)
	}
	if session == nil {
		return nil, fmt.Errorf("Invalid session object")
	}

	var buf bytes.Buffer
	session.Stdout = &buf
	if err := session.Run(remoteCmdString); err != nil {
		return nil, err
	}
	session.Close()
	fmt.Println(buf.String())
	return nil, nil
}
