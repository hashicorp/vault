package ssh

import (
	"context"
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/mitchellh/mapstructure"
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

func (b *backend) secretDynamicKeyRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return &logical.Response{Secret: req.Secret}, nil
}

func (b *backend) secretDynamicKeyRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	type sec struct {
		AdminUser        string `mapstructure:"admin_user"`
		Username         string `mapstructure:"username"`
		IP               string `mapstructure:"ip"`
		HostKeyName      string `mapstructure:"host_key_name"`
		DynamicPublicKey string `mapstructure:"dynamic_public_key"`
		InstallScript    string `mapstructure:"install_script"`
		Port             int    `mapstructure:"port"`
	}

	intSec := &sec{}
	err := mapstructure.Decode(req.Secret.InternalData, intSec)
	if err != nil {
		return nil, errwrap.Wrapf("secret internal data could not be decoded: {{err}}", err)
	}

	// Fetch the host key using the key name
	hostKey, err := b.getKey(ctx, req.Storage, intSec.HostKeyName)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("key %q not found error: {{err}}", intSec.HostKeyName), err)
	}
	if hostKey == nil {
		return nil, fmt.Errorf("key %q not found", intSec.HostKeyName)
	}

	// Remove the public key from authorized_keys file in target machine
	// The last param 'false' indicates that the key should be uninstalled.
	err = b.installPublicKeyInTarget(ctx, intSec.AdminUser, intSec.Username, intSec.IP, intSec.Port, hostKey.Key, intSec.DynamicPublicKey, intSec.InstallScript, false)
	if err != nil {
		return nil, fmt.Errorf("error removing public key from authorized_keys file in target")
	}
	return nil, nil
}
