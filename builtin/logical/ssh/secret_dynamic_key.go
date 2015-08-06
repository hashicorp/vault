package ssh

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/helper/uuid"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const SecretDynamicKeyType = "secret_dynamic_key_type"

func secretDynamicKey(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretDynamicKeyType,
		Fields: map[string]*framework.FieldSchema{
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Username in host",
			},
			"ip": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "IP address of host",
			},
		},
		DefaultDuration:    10 * time.Minute,
		DefaultGracePeriod: 2 * time.Minute,
		Renew:              b.secretDynamicKeyRenew,
		Revoke:             b.secretDynamicKeyRevoke,
	}
}

func (b *backend) secretDynamicKeyRenew(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
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

func (b *backend) secretDynamicKeyRevoke(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	adminUserRaw, ok := req.Secret.InternalData["admin_user"]
	if !ok {
		return nil, fmt.Errorf("secret is missing internal data")
	}
	adminUser, ok := adminUserRaw.(string)
	if !ok {
		return nil, fmt.Errorf("secret is missing internal data")
	}

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

	installScriptRaw, ok := req.Secret.InternalData["install_script"]
	if !ok {
		return nil, fmt.Errorf("secret is missing internal data")
	}
	installScript := installScriptRaw.(string)
	if !ok {
		return nil, fmt.Errorf("secret is missing internal data")
	}

	portRaw, ok := req.Secret.InternalData["port"]
	if !ok {
		return nil, fmt.Errorf("secret is missing internal data")
	}
	port := portRaw.(string)

	// Fetch the host key using the key name
	hostKeyEntry, err := req.Storage.Get(fmt.Sprintf("keys/%s", hostKeyName))
	if err != nil {
		return nil, fmt.Errorf("key '%s' not found error:%s", hostKeyName, err)
	}
	var hostKey sshHostKey
	if err := hostKeyEntry.DecodeJSON(&hostKey); err != nil {
		return nil, fmt.Errorf("error reading the host key: %s", err)
	}

	// Transfer the dynamic public key to target machine and use it to remove the entry from authorized_keys file
	dynamicPublicKeyFileName := uuid.GenerateUUID()
	err = scpUpload(adminUser, ip, port, hostKey.Key, dynamicPublicKeyFileName, dynamicPublicKey)
	if err != nil {
		return nil, fmt.Errorf("error uploading pubic key: %s", err)
	}

	scriptFileName := fmt.Sprintf("%s.sh", dynamicPublicKeyFileName)
	err = scpUpload(adminUser, ip, port, hostKey.Key, scriptFileName, installScript)
	if err != nil {
		return nil, fmt.Errorf("error uploading script file: %s", err)
	}

	// Remove the public key from authorized_keys file in target machine
	err = uninstallPublicKeyInTarget(adminUser, dynamicPublicKeyFileName, username, ip, port, hostKey.Key)
	if err != nil {
		return nil, fmt.Errorf("error removing public key from authorized_keys file in target")
	}
	return nil, nil
}
