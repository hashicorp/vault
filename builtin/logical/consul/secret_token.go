package consul

import (
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
	c, userErr, intErr := client(req.Storage)
	if intErr != nil {
		return nil, intErr
	}
	if userErr != nil {
		// Returning logical.ErrorResponse from revocation function is risky
		return nil, userErr
	}

	tokenRaw, ok := req.Secret.InternalData["token"]
	if !ok {
		// We return nil here because this is a pre-0.5.3 problem and there is
		// nothing we can do about it. We already can't revoke the lease
		// properly if it has been renewed and this is documented pre-0.5.3
		// behavior with a security bulletin about it.
		return nil, nil
	}

	_, err := c.ACL().Destroy(tokenRaw.(string), nil)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
