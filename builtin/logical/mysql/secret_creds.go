package mysql

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const SecretCredsType = "creds"

// defaultRevocationSQL is a default SQL statement for revoking a user. Revoking
// permissions for the user is done before the drop, because MySQL explicitly
// documents that open user connections will not be closed. By revoking all
// grants, at least we ensure that the open connection is useless. Dropping the
// user will only affect the next connection.
const defaultRevocationSQL = `
REVOKE ALL PRIVILEGES, GRANT OPTION FROM '{{name}}'@'%'; 
DROP USER '{{name}}'@'%'
`

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
	var resp *logical.Response

	// Get the username from the internal data
	usernameRaw, ok := req.Secret.InternalData["username"]
	if !ok {
		return nil, fmt.Errorf("secret is missing username internal data")
	}
	username, ok := usernameRaw.(string)
	if !ok {
		return nil, fmt.Errorf("usernameRaw is not a string")
	}

	// Get our connection
	db, err := b.DB(req.Storage)
	if err != nil {
		return nil, err
	}

	roleName := ""
	roleNameRaw, ok := req.Secret.InternalData["role"]
	if ok {
		roleName = roleNameRaw.(string)
	}

	var role *roleEntry
	if roleName != "" {
		role, err = b.Role(req.Storage, roleName)
		if err != nil {
			return nil, err
		}
	}

	// Use a default SQL statement for revocation if one cannot be fetched from the role
	revocationSQL := defaultRevocationSQL

	if role != nil && role.RevocationSQL != "" {
		revocationSQL = role.RevocationSQL
	} else {
		if resp == nil {
			resp = &logical.Response{}
		}
		resp.AddWarning(fmt.Sprintf("Role %q cannot be found. Using default SQL for revoking user.", roleName))
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	for _, query := range strutil.ParseArbitraryStringSlice(revocationSQL, ";") {
		query = strings.TrimSpace(query)
		if len(query) == 0 {
			continue
		}

		// This is not a prepared statement because not all commands are supported
		// 1295: This command is not supported in the prepared statement protocol yet
		// Reference https://mariadb.com/kb/en/mariadb/prepare-statement/
		query = strings.Replace(query, "{{name}}", username, -1)
		_, err = tx.Exec(query)
		if err != nil {
			return nil, err
		}

	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return resp, nil
}
