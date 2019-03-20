package database

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
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
			return nil, fmt.Errorf("could not find role with name: %q", req.Secret.InternalData["role"])
		}

		role, err := b.Role(ctx, req.Storage, roleNameRaw.(string))
		if err != nil {
			return nil, err
		}
		if role == nil {
			return nil, fmt.Errorf("error during renew: could not find role with name %q", req.Secret.InternalData["role"])
		}

		// Get the Database object
		db, err := b.GetConnection(ctx, req.Storage, role.DBName)
		if err != nil {
			return nil, err
		}

		db.RLock()
		defer db.RUnlock()

		// Make sure we increase the VALID UNTIL endpoint for this user.
		ttl, _, err := framework.CalculateTTL(b.System(), req.Secret.Increment, role.DefaultTTL, 0, role.MaxTTL, 0, req.Secret.IssueTime)
		if err != nil {
			return nil, err
		}
		if ttl > 0 {
			expireTime := time.Now().Add(ttl)
			// Adding a small buffer since the TTL will be calculated again after this call
			// to ensure the database credential does not expire before the lease
			expireTime = expireTime.Add(5 * time.Second)
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

		var dbName string
		var statements dbplugin.Statements

		role, err := b.Role(ctx, req.Storage, roleNameRaw.(string))
		if err != nil {
			return nil, err
		}
		if role != nil {
			dbName = role.DBName
			statements = role.Statements
		} else {
			if dbNameRaw, ok := req.Secret.InternalData["db_name"]; !ok {
				return nil, fmt.Errorf("error during revoke: could not find role with name %q or embedded revocation db name data", req.Secret.InternalData["role"])
			} else {
				dbName = dbNameRaw.(string)
			}
			if statementsRaw, ok := req.Secret.InternalData["revocation_statements"]; !ok {
				return nil, fmt.Errorf("error during revoke: could not find role with name %q or embedded revocation statement data", req.Secret.InternalData["role"])
			} else {
				// If we don't actually have any statements, because none were
				// set in the role, we'll end up with an empty one and the
				// default for the db type will be attempted
				if statementsRaw != nil {
					statementsSlice, ok := statementsRaw.([]interface{})
					if !ok {
						return nil, fmt.Errorf("error during revoke: could not find role with name %q and embedded reovcation data could not be read", req.Secret.InternalData["role"])
					} else {
						for _, v := range statementsSlice {
							statements.Revocation = append(statements.Revocation, v.(string))
						}
					}
				}
			}
		}

		// Get our connection
		db, err := b.GetConnection(ctx, req.Storage, dbName)
		if err != nil {
			return nil, err
		}

		db.RLock()
		defer db.RUnlock()

		if err := db.RevokeUser(ctx, statements, username); err != nil {
			b.CloseIfShutdown(db, err)
			return nil, err
		}
		return resp, nil
	}
}
