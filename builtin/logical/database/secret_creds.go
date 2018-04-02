package database

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const SecretCredsType = "creds"

func secretCreds(b *databaseBackend) *framework.Secret {
	return &framework.Secret{
		Type:   SecretCredsType,
		Fields: map[string]*framework.FieldSchema{},

		Renew:  b.secretCredsRenew(),
		Revoke: b.secretCredsRevoke(),
	}
}

func (b *databaseBackend) secretCredsRenew() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		// Get the username from the internal data
		usernameRaw, ok := req.Secret.InternalData["username"]
		if !ok {
			return nil, fmt.Errorf("secret is missing username internal data")
		}
		username, ok := usernameRaw.(string)

		roleNameRaw, ok := req.Secret.InternalData["role"]
		if !ok {
			return nil, fmt.Errorf("could not find role with name: %s", req.Secret.InternalData["role"])
		}

		role, err := b.Role(ctx, req.Storage, roleNameRaw.(string))
		if err != nil {
			return nil, err
		}
		if role == nil {
			return nil, fmt.Errorf("error during renew: could not find role with name %s", req.Secret.InternalData["role"])
		}

		// Get the Database object
		db, err := b.GetConnection(ctx, req.Storage, role.DBName)
		if err != nil {
			return nil, err
		}

		db.RLock()
		defer db.RUnlock()

		// Make sure we increase the VALID UNTIL endpoint for this user.  This value is estimated and does not
		// take into account any backend specific values.  These value will be calculated by core and will only
		// reduce the TTL based on any running max ttl.  Since vault still manages the lease, it will still get
		// revokes at the lesser time.
		if req.Secret.EstimatedTTL > 0 {
			ttl := req.Secret.EstimatedTTL
			if role.DefaultTTL > 0 && role.DefaultTTL < ttl {
				ttl = role.DefaultTTL
			}
			if role.MaxTTL > 0 && role.MaxTTL < ttl {
				ttl = role.MaxTTL
			}
			expireTime := time.Now().Add(ttl)
			err := db.RenewUser(ctx, role.Statements, username, expireTime)
			if err != nil {
				b.CloseIfShutdown(db, err)
				return nil, err
			}
		}
		resp := &logical.Response{Secret: req.Secret}
		resp.Secret.TTL = role.DefaultTTL
		resp.Secret.MaxTTL = role.MaxTTL
		return resp, nil
	}
}

func (b *databaseBackend) secretCredsRevoke() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		// Get the username from the internal data
		usernameRaw, ok := req.Secret.InternalData["username"]
		if !ok {
			return nil, fmt.Errorf("secret is missing username internal data")
		}
		username, ok := usernameRaw.(string)

		var resp *logical.Response

		roleNameRaw, ok := req.Secret.InternalData["role"]
		if !ok {
			return nil, fmt.Errorf("no role name was provided")
		}

		role, err := b.Role(ctx, req.Storage, roleNameRaw.(string))
		if err != nil {
			return nil, err
		}
		if role == nil {
			return nil, fmt.Errorf("error during revoke: could not find role with name %s", req.Secret.InternalData["role"])
		}

		// Get our connection
		db, err := b.GetConnection(ctx, req.Storage, role.DBName)
		if err != nil {
			return nil, err
		}

		db.RLock()
		defer db.RUnlock()

		if err := db.RevokeUser(ctx, role.Statements, username); err != nil {
			b.CloseIfShutdown(db, err)
			return nil, err
		}
		return resp, nil
	}
}
