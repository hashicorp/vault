package ssh

import (
	"fmt"

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

		Renew:  b.secretDynamicKeyRenew,
		Revoke: b.secretDynamicKeyRevoke,
	}
}

func (b *backend) secretDynamicKeyRenew(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	f := framework.LeaseExtend(0, 0, b.System())
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
	port := int(portRaw.(float64))

	// Fetch the host key using the key name
	hostKey, err := b.getKey(req.Storage, hostKeyName)
	if err != nil {
		return nil, fmt.Errorf("key '%s' not found error:%s", hostKeyName, err)
	}
	if hostKey == nil {
		return nil, fmt.Errorf("key '%s' not found", hostKeyName)
	}

	// Remove the public key from authorized_keys file in target machine
	// The last param 'false' indicates that the key should be uninstalled.
	err = b.installPublicKeyInTarget(adminUser, username, ip, port, hostKey.Key, dynamicPublicKey, installScript, false)
	if err != nil {
		return nil, fmt.Errorf("error removing public key from authorized_keys file in target")
	}
	return nil, nil
}
