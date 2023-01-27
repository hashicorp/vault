package ssh

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const SecretDynamicKeyType = "secret_dynamic_key_type"

func secretDynamicKey(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretDynamicKeyType,
		Fields: map[string]*framework.FieldSchema{
			"username": {
				Type:        framework.TypeString,
				Description: "Username in host",
			},
			"ip": {
				Type:        framework.TypeString,
				Description: "IP address of host",
			},
		},

		Renew:  b.secretDynamicKeyRenew,
		Revoke: b.secretDynamicKeyRevoke,
	}
}

func (b *backend) secretDynamicKeyRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return nil, nil
}

func (b *backend) secretDynamicKeyRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return nil, nil
}
