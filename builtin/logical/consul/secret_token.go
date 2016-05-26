package consul

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
		Revoke: secretTokenRevoke,
	}
}

func (b *backend) secretTokenRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	return framework.LeaseExtend(0, 0, b.System())(req, d)
}

func secretTokenRevoke(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	c, err := client(req.Storage)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	tokenRaw, ok := req.Secret.InternalData["token"]
	if !ok {
		return nil, fmt.Errorf("secret is missing internal data: token")
	}

	_, err = c.ACL().Destroy(tokenRaw.(string), nil)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	return nil, nil
}
