package database

import (
	"errors"
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const SecretCredsType = "creds"

func secretCreds(b *databaseBackend) *framework.Secret {
	return &framework.Secret{
		Type: SecretCredsType,
		Fields: map[string]*framework.FieldSchema{
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Username",
			},

			"password": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Password",
			},
		},

		Renew:  b.secretCredsRenew,
		Revoke: b.secretCredsRevoke,
	}
}

func (b *databaseBackend) secretCredsRenew(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	dbName := d.Get("name").(string)

	// Get the username from the internal data
	usernameRaw, ok := req.Secret.InternalData["username"]
	if !ok {
		return nil, fmt.Errorf("secret is missing username internal data")
	}
	username, ok := usernameRaw.(string)

	// Get our connection
	db, ok := b.connections[dbName]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Could not find connection with name %s", dbName))
	}

	// Get the lease information
	lease, err := b.Lease(req.Storage)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		lease = &configLease{}
	}

	f := framework.LeaseExtend(lease.Lease, lease.LeaseMax, b.System())
	resp, err := f(req, d)
	if err != nil {
		return nil, err
	}

	// Make sure we increase the VALID UNTIL endpoint for this user.
	if expireTime := resp.Secret.ExpirationTime(); !expireTime.IsZero() {
		expiration := expireTime.Format("2006-01-02 15:04:05-0700")

		err := db.RenewUser(username, expiration)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func (b *databaseBackend) secretCredsRevoke(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the username from the internal data
	usernameRaw, ok := req.Secret.InternalData["username"]
	if !ok {
		return nil, fmt.Errorf("secret is missing username internal data")
	}
	username, ok := usernameRaw.(string)

	var revocationSQL string
	var resp *logical.Response

	roleNameRaw, ok := req.Secret.InternalData["role"]
	if !ok {
		return nil, fmt.Errorf("Could not find role with name: %s", req.Secret.InternalData["role"])
	}

	role, err := b.Role(req.Storage, roleNameRaw.(string))
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, fmt.Errorf("Could not find role with name: %s", req.Secret.InternalData["role"])
	}

	/* TODO: think about how to handle this case.
	if !ok {
		role, err := b.Role(req.Storage, roleNameRaw.(string))
		if err != nil {
			return nil, err
		}
		if role == nil {
			if resp == nil {
				resp = &logical.Response{}
			}
			resp.AddWarning(fmt.Sprintf("Role %q cannot be found. Using default revocation SQL.", roleNameRaw.(string)))
		} else {
			revocationSQL = role.RevocationStatement
		}
	}*/

	// Grab the read lock
	b.RLock()
	defer b.RUnlock()

	// Get our connection
	db, ok := b.connections[role.DBName]
	if !ok {
		return nil, fmt.Errorf("Could not find database with name: %s", role.DBName)
	}

	// TODO: Maybe move this down into db package?
	switch revocationSQL {

	// This is the default revocation logic. If revocation SQL is provided it
	// is simply executed as-is.
	case "":
		err := db.DefaultRevokeUser(username)
		if err != nil {
			return nil, err
		}

	// We have revocation SQL, execute directly, within a transaction
	default:
		err := db.CustomRevokeUser(username, revocationSQL)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}
