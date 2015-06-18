package postgresql

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/lib/pq"
)

const SecretCredsType = "creds"

func secretCreds(b *backend) *framework.Secret {
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

		DefaultDuration:    1 * time.Hour,
		DefaultGracePeriod: 10 * time.Minute,

		Renew:  b.secretCredsRenew,
		Revoke: b.secretCredsRevoke,
	}
}

func (b *backend) secretCredsRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the username from the internal data
	usernameRaw, ok := req.Secret.InternalData["username"]
	if !ok {
		return nil, fmt.Errorf("secret is missing username internal data")
	}
	username, ok := usernameRaw.(string)

	// Get our connection
	db, err := b.DB(req.Storage)
	if err != nil {
		return nil, err
	}

	// Get the lease information
	lease, err := b.Lease(req.Storage)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		lease = &configLease{Lease: 1 * time.Hour}
	}

	f := framework.LeaseExtend(lease.Lease, lease.LeaseMax, false)
	resp, err := f(req, d)
	if err != nil {
		return nil, err
	}

	// Make sure we increase the VALID UNTIL endpoint for this user.
	if expireTime := resp.Secret.ExpirationTime(); !expireTime.IsZero() {
		expiration := expireTime.Add(10 * time.Minute).
			Format("2006-01-02 15:04:05")

		query := fmt.Sprintf(
			"ALTER ROLE %s VALID UNTIL '%s';",
			pq.QuoteIdentifier(username),
			expiration)
		stmt, err := db.Prepare(query)
		if err != nil {
			return nil, err
		}
		defer stmt.Close()
		if _, err := stmt.Exec(); err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func (b *backend) secretCredsRevoke(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the username from the internal data
	usernameRaw, ok := req.Secret.InternalData["username"]
	if !ok {
		return nil, fmt.Errorf("secret is missing username internal data")
	}
	username, ok := usernameRaw.(string)

	// Get our connection
	db, err := b.DB(req.Storage)
	if err != nil {
		return nil, err
	}

	// Drop this user
	stmt, err := db.Prepare(fmt.Sprintf(
		"DROP ROLE IF EXISTS %s;", pq.QuoteIdentifier(username)))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(); err != nil {
		return nil, err
	}

	return nil, nil
}
