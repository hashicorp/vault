package ssh

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const SecretOneTimeKeyType = "secret_one_type_key_type"

func secretSSHKey(b *backend) *framework.Secret {
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
		DefaultDuration:    1 * time.Hour,
		DefaultGracePeriod: 10 * time.Minute,
		Renew:              b.secretSSHKeyRenew,
		Revoke:             b.secretSSHKeyRevoke,
	}
}

func (b *backend) secretSSHKeyRenew(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
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

func (b *backend) secretSSHKeyRevoke(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
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
	hostKeyEntry, err := req.Storage.Get(fmt.Sprintf("keys/%s", hostKeyName))
	if err != nil {
		return nil, fmt.Errorf("key '%s' not found error:%s", hostKeyName, err)
	}
	var hostKey sshHostKey
	if err := hostKeyEntry.DecodeJSON(&hostKey); err != nil {
		return nil, fmt.Errorf("error reading the host key: %s", err)
	}

	//write host key to file and use it as argument to scp command
	hostKeyFileName := fmt.Sprintf("./vault_ssh_%s_%s_shared.pem", username, ip)
	err = ioutil.WriteFile(hostKeyFileName, []byte(hostKey.Key), 0400)

	//write dynamicPublicKey to file and use it as argument to scp command
	dynamicPublicKeyFileName := fmt.Sprintf("vault_ssh_%s_%s_otk.pem.pub", username, ip)
	err = ioutil.WriteFile(dynamicPublicKeyFileName, []byte(dynamicPublicKey), 0400)

	//transfer the dynamic public key to target machine and use it to remove the entry from authorized_keys file
	err = uploadFileScp(dynamicPublicKeyFileName, username, ip, hostKey.Key)
	if err != nil {
		return nil, fmt.Errorf("public key transfer failed: %s", err)
	}

	//connect to target machine
	session, err := createSSHPublicKeysSession(username, ip, hostKey.Key)
	if err != nil {
		return nil, fmt.Errorf("unable to create SSH Session using public keys: %s", err)
	}
	if session == nil {
		return nil, fmt.Errorf("invalid session object")
	}

	authKeysFileName := "/home/" + username + "/.ssh/authorized_keys"
	tempKeysFileName := "/home/" + username + "/temp_authorized_keys"

	//commands to be run on target machine
	grepCmd := fmt.Sprintf("grep -vFf %s %s > %s", dynamicPublicKeyFileName, authKeysFileName, tempKeysFileName)
	catCmdRemoveDuplicate := fmt.Sprintf("cat %s > %s", tempKeysFileName, authKeysFileName)
	removeCmd := fmt.Sprintf("rm -f %s %s", tempKeysFileName, dynamicPublicKeyFileName)

	remoteCmd := fmt.Sprintf("%s;%s;%s", grepCmd, catCmdRemoveDuplicate, removeCmd)

	//run the commands in target machine
	if err := session.Run(remoteCmd); err != nil {
		return nil, err
	}

	return nil, nil
}
