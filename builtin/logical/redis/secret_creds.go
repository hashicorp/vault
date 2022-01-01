package redis

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func secretCreds(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: "creds",
		Fields: map[string]*framework.FieldSchema{
			"username": {
				Type: framework.TypeString,
			},
			"password": {
				Type: framework.TypeString,
			},
		},

		Renew:  b.secretCredsRenew,
		Revoke: b.secretCredsRevoke,
	}
}

func (b *backend) secretCredsRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return &logical.Response{Secret: req.Secret}, nil
}

func (b *backend) secretCredsRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	client, err := b.Client(ctx, req.Storage)
	if err != nil {
		return logical.ErrorResponse("failed to get Redis client: %s", err), nil
	}

	_, err = client.Do(ctx, "ACL", "DELUSER", req.Secret.InternalData["username"]).Result()
	if err != nil {
		return logical.ErrorResponse("failed to delete user: %s", err), nil
	}

	return nil, nil
}
