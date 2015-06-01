package consul

import (
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	SecretTokenType      = "token"
	DefaultLeaseDuration = 1 * time.Hour
	DefaultGracePeriod   = 10 * time.Minute
)

func secretToken() *framework.Secret {
	return &framework.Secret{
		Type: SecretTokenType,
		Fields: map[string]*framework.FieldSchema{
			"token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Request token",
			},
		},

		DefaultDuration:    DefaultLeaseDuration,
		DefaultGracePeriod: DefaultGracePeriod,

		Renew:  framework.LeaseExtend(1*time.Hour, 0),
		Revoke: secretTokenRevoke,
	}
}

func secretTokenRevoke(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	c, err := client(req.Storage)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	_, err = c.ACL().Destroy(d.Get("token").(string), nil)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	return nil, nil
}
