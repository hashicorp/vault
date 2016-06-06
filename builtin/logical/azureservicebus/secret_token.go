package azureservicebus

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const SecretTokenType = "token"

func secretToken(b *backend) *framework.Secret {
	return &framework.Secret{
		Type:   SecretTokenType,
		Fields: map[string]*framework.FieldSchema{},
		Revoke: b.secretTokenRevoke,
	}
}

func (b *backend) secretTokenRevoke(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Since SAS tokens have a natural expiry time, no revokation on the part of vault is needed
	return nil, nil
}
