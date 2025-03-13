package builder

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// return the secret
func (gb *GenericBackend[O, C, R]) secret(secret *Secret[R, C]) *framework.Secret {
	return &framework.Secret{
		Type:   secret.Type,
		Fields: secret.Fields,
		Revoke: gb.tokenRevoke,
		Renew:  gb.tokenRenew,
	}
}

// tokenRevoke removes the token from the Vault storage API and calls the client to revoke the token
func (gb *GenericBackend[O, C, R]) tokenRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	client, role, err := gb.clientAndRole(ctx, req, d)
	if err != nil {
		return nil, err
	}

	return gb.role.Secret.RevokeFunc(req, d, client, role)
}

// tokenRenew calls the client to create a new token and stores it in the Vault storage API
func (gb *GenericBackend[O, C, R]) tokenRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	client, role, err := gb.clientAndRole(ctx, req, d)
	if err != nil {
		return nil, err
	}

	return gb.role.Secret.RenewFunc(req, d, client, role)
}
