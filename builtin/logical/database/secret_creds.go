package database

import (
	"fmt"

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
	return func(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
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

		role, err := b.Role(req.Storage, roleNameRaw.(string))
		if err != nil {
			return nil, err
		}
		if role == nil {
			return nil, fmt.Errorf("error during renew: could not find role with name %s", req.Secret.InternalData["role"])
		}

		f := framework.LeaseExtend(role.DefaultTTL, role.MaxTTL, b.System())
		resp, err := f(req, data)
		if err != nil {
			return nil, err
		}

		// Grab the read lock
		b.RLock()
		var unlockFunc func() = b.RUnlock

		// Get the Database object
		db, ok := b.getDBObj(role.DBName)
		if !ok {
			// Upgrade lock
			b.RUnlock()
			b.Lock()
			unlockFunc = b.Unlock

			// Create a new DB object
			db, err = b.createDBObj(req.Storage, role.DBName)
			if err != nil {
				unlockFunc()
				return nil, fmt.Errorf("cound not retrieve db with name: %s, got error: %s", role.DBName, err)
			}
		}

		// Make sure we increase the VALID UNTIL endpoint for this user.
		if expireTime := resp.Secret.ExpirationTime(); !expireTime.IsZero() {
			err := db.RenewUser(role.Statements, username, expireTime)
			// Unlock
			unlockFunc()
			if err != nil {
				b.closeIfShutdown(role.DBName, err)
				return nil, err
			}
		}

		return resp, nil
	}
}

func (b *databaseBackend) secretCredsRevoke() framework.OperationFunc {
	return func(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
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

		role, err := b.Role(req.Storage, roleNameRaw.(string))
		if err != nil {
			return nil, err
		}
		if role == nil {
			return nil, fmt.Errorf("error during revoke: could not find role with name %s", req.Secret.InternalData["role"])
		}

		// Grab the read lock
		b.RLock()
		var unlockFunc func() = b.RUnlock

		// Get our connection
		db, ok := b.getDBObj(role.DBName)
		if !ok {
			// Upgrade lock
			b.RUnlock()
			b.Lock()
			unlockFunc = b.Unlock

			// Create a new DB object
			db, err = b.createDBObj(req.Storage, role.DBName)
			if err != nil {
				unlockFunc()
				return nil, fmt.Errorf("cound not retrieve db with name: %s, got error: %s", role.DBName, err)
			}
		}

		err = db.RevokeUser(role.Statements, username)
		// Unlock
		unlockFunc()
		if err != nil {
			b.closeIfShutdown(role.DBName, err)
			return nil, err
		}

		return resp, nil
	}
}
