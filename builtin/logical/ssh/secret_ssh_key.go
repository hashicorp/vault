package ssh

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const SecretOneTimeKeyType = "secret_one_type_key_type"

func secretSshKey(b *backend) *framework.Secret {
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
		DefaultDuration:    5 * time.Second, //TODO: change this
		DefaultGracePeriod: 1 * time.Second, //TODO: change this
		Renew:              b.secretSshKeyRenew,
		Revoke:             b.secretSshKeyRevoke,
	}
}

func (b *backend) secretSshKeyRenew(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
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
	dynamicPrivateKeyFileName := "vault_ssh_" + username + "_" + ip + "_otk.pem"
	dynamicPublicKeyFileName := dynamicPrivateKeyFileName + ".pub"
	err = ioutil.WriteFile(dynamicPublicKeyFileName, []byte(dynamicPublicKey), 0400)

	//transfer the dynamic public key to target machine and use it to remove the entry from authorized_keys file
	err = uploadFileScp(dynamicPublicKeyFileName, username, ip, hostKey.Key)
	if err != nil {
		return nil, fmt.Errorf("Public key transfer failed: %s", err)
	}

	authKeysFileName := "~/.ssh/authorized_keys"
	tempKeysFileName := "~/temp_authorized_keys"

	//commands to be run on target machine
	grepCmd := "grep -vFf " + dynamicPublicKeyFileName + " " + authKeysFileName + " > " + tempKeysFileName + ";"
	catCmdRemoveDuplicate := "cat " + tempKeysFileName + " > " + authKeysFileName + ";"
	rmCmd := "rm -f " + tempKeysFileName + " " + dynamicPublicKeyFileName + ";"
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

	//run the commands in target machine
	if err := session.Run(remoteCmdString); err != nil {
		return nil, err
	}

	return nil, nil
}
