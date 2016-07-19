package mongodb

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"gopkg.in/mgo.v2"
)

const SecretCredsType = "creds"

func secretCreds(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretCredsType,
		Fields: map[string]*framework.FieldSchema{
			"username": {
				Type:        framework.TypeString,
				Description: "Username",
			},

			"password": {
				Type:        framework.TypeString,
				Description: "Password",
			},
		},

		Renew:  b.secretCredsRenew,
		Revoke: b.secretCredsRevoke,
	}
}

func (b *backend) secretCredsRenew(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the lease information
	leaseConfig, err := b.LeaseConfig(req.Storage)
	if err != nil {
		return nil, err
	}
	if leaseConfig == nil {
		leaseConfig = &configLease{}
	}

	f := framework.LeaseExtend(leaseConfig.TTL, leaseConfig.MaxTTL, b.System())
	return f(req, d)
}

func (b *backend) secretCredsRevoke(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the username from the internal data
	usernameRaw, ok := req.Secret.InternalData["username"]
	if !ok {
		return nil, fmt.Errorf("secret is missing username internal data")
	}
	username, ok := usernameRaw.(string)
	if !ok {
		return nil, fmt.Errorf("username internal data is not a string")
	}

	// Get the db from the internal data
	dbRaw, ok := req.Secret.InternalData["db"]
	if !ok {
		return nil, fmt.Errorf("secret is missing db internal data")
	}
	db, ok := dbRaw.(string)
	if !ok {
		return nil, fmt.Errorf("db internal data is not a string")
	}

	// Get our connection
	session, err := b.Session(req.Storage)
	if err != nil {
		return nil, err
	}

	// Drop the user
	err = session.DB(db).RemoveUser(username)
	if err != nil && err != mgo.ErrNotFound {
		return nil, err
	}

	return nil, nil
}
