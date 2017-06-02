package hana

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
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

		Renew:  b.secretCredsRenew,
		Revoke: b.secretCredsRevoke,
	}
}

// Renews lease, and resets the HANA account's valid until field on the server
func (b *backend) secretCredsRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the lease information
	leaseConfig, err := b.LeaseConfig(req.Storage)
	if err != nil {
		return nil, err
	}
	if leaseConfig == nil {
		leaseConfig = &configLease{}
	}

	// Get the username
	usernameRaw, ok := req.Secret.InternalData["username"]
	if !ok {
		return nil, fmt.Errorf("secret is missing username internal data")
	}
	username, ok := usernameRaw.(string)

	// Get our handle
	db, err := b.DB(req.Storage)
	if err != nil {
		return nil, err
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Request server's current time plus lease duration
	var validUntil string
	timeQuery := fmt.Sprintf("SELECT TO_NVARCHAR(add_seconds(CURRENT_TIMESTAMP," +
		"%f), 'YYYY-MM-DD HH24:MI:SS') FROM DUMMY", (leaseConfig.TTL).Seconds())
	err = db.QueryRow(timeQuery).Scan(&validUntil)
	if err != nil {
		return nil, err
	}

	// Renew user's valid until property field
	stmt, err := tx.Prepare("ALTER USER " + username + " VALID UNTIL " + "'" + validUntil + "'")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	f := framework.LeaseExtend(leaseConfig.TTL, leaseConfig.TTLMax, b.System())
	return f(req, d)
}

// Revoking tries to deactivate user and also drop user
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

	// Disable server login for user
	disableStmt, err := db.Prepare(fmt.Sprintf("ALTER USER %s DEACTIVATE USER NOW", username))
	if err != nil {
		return nil, err
	}
	defer disableStmt.Close()
	if _, err := disableStmt.Exec(); err != nil {
		return nil, err
	}

	// Drop user (Restrict - only drop user if they do not own any object other than
	// their own schema and other schemas created by the user)
	// This also invalidates current sessions
	dropStmt, err := db.Prepare(fmt.Sprintf("DROP USER %s RESTRICT", username))
	if err != nil {
		return nil, err
	}
	defer dropStmt.Close()
	if _, err := dropStmt.Exec(); err != nil {
		return nil, err
	}

	return nil, nil
}