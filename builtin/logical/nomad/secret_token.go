package nomad

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	SecretTokenType = "token"
)

func secretToken(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretTokenType,
		Fields: map[string]*framework.FieldSchema{
			"token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Request token",
			},
		},

		Renew:  b.secretTokenRenew,
		Revoke: b.secretTokenRevoke,
	}
}

func (b *backend) secretTokenRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	lease, err := b.LeaseConfig(req.Storage)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		lease = &configLease{}
	}

	return framework.LeaseExtend(lease.TTL, lease.MaxTTL, b.System())(req, d)
}

func (b *backend) secretTokenRevoke(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	c, err := b.client(req.Storage)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, fmt.Errorf("Error connecting with Nomad")
	}

	accessorIDRaw, ok := req.Secret.InternalData["accessor_id"]
	if !ok {
		return nil, fmt.Errorf("accessor_id is missing on the lease")
	}

	_, err = c.ACLTokens().Delete(accessorIDRaw.(string), nil)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
