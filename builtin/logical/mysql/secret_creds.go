package mysql

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/vault/helper/strutil"
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

			"rolename": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Rolename",
			},
		},

		Renew:  b.secretCredsRenew,
		Revoke: b.secretCredsRevoke,
	}
}

func (b *backend) secretCredsRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the lease information
	lease, err := b.Lease(req.Storage)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		lease = &configLease{}
	}

	f := framework.LeaseExtend(lease.Lease, lease.LeaseMax, b.System())
	return f(req, d)
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

	// Get the role
	// pathParts := strings.Split(req.Path, "/")
	log.Println("InternalData")
	log.Printf("%+v", req.Secret.InternalData)
	rolenameRaw, ok := req.Secret.InternalData["rolename"]
	if !ok {
		return nil, fmt.Errorf("secret is missing rollname internal data")
	}
	rolename, ok := rolenameRaw.(string)

	role, err := b.Role(req.Storage, rolename)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown role: %s", rolename)), nil
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Check for an empty revokeSQL string
	// set it to a default query if the string is empty
	if role.RevokeSQL == "" {
		role.RevokeSQL = "REVOKE ALL PRIVILEGES, GRANT OPTION FROM '" + username + "'@'%'; DROP USER '" + username + "'@'%'"
	}

	for _, query := range strutil.ParseArbitraryStringSlice(role.RevokeSQL, ";") {
		query = strings.TrimSpace(query)
		if len(query) == 0 {
			continue
		}

		// convert {{name}} to username value because we can't use
		// prepared statements for REVOKE and DROP commands
		// Generates the following error
		// 1295: This command is not supported in the prepared statement protocol yet
		// Reference https://mariadb.com/kb/en/mariadb/prepare-statement/
		query = strings.Replace(query, "{{name}}", username, -1)

		// Revoke all permissions for the user. This is done before the
		// drop, because MySQL explicitly documents that open user connections
		// will not be closed. By revoking all grants, at least we ensure
		// that the open connection is useless.
		// Drop this user. This only affects the next connection, which is
		// why we do the revoke initially.
		_, err = tx.Exec(query)
		if err != nil {
			return nil, err
		}

	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return nil, nil
}
