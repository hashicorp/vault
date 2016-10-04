package mysql

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const SecretCredsType = "creds"

// Default Revoke and Drop user SQL statment
// Revoking permissions for the user is done before the
// drop, because MySQL explicitly documents that open user connections
// will not be closed. By revoking all grants, at least we ensure
// that the open connection is useless.
// Dropping the user will only affect the next connection
// This is not a prepared statement because not all commands are supported
// 1295: This command is not supported in the prepared statement protocol yet
// Reference https://mariadb.com/kb/en/mariadb/prepare-statement/
const defaultRevokeSQL = `
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

			"role": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Role",
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

	// Get the role name
	// we may not always have role data in the secret InternalData
	// so don't exit if the roleNameRaw fails. Instead it is set
	// as an empty string.
	var roleName string
	roleNameRaw, ok := req.Secret.InternalData["role"]
	if !ok {
		roleName = ""
	} else {
		roleName, _ = roleNameRaw.(string)
	}
	// init default revoke sql string.
	// this will replaced by a user provided one if one exists
	// otherwise this is what will be used when lease is revoked
	revokeSQL := defaultRevokeSQL

	// init bool to track if we should responding with warning about nil role
	nonNilResponse := false

	// if we were successful in finding a role name
	// create role entry from that name
	if roleName != "" {
		role, err := b.Role(req.Storage, roleName)
		if err != nil {
			return nil, err
		}

		if role == nil {
			nonNilResponse = true
		}

		// Check for a revokeSQL string
		// if one exists use that instead of the default
		if role.RevokeSQL != "" && role != nil {
			revokeSQL = role.RevokeSQL
		}
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	for _, query := range strutil.ParseArbitraryStringSlice(revokeSQL, ";") {
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

	// Let the user know that since we had a nil role we used the default SQL revocation statment
	if nonNilResponse == true {
		// unable to get role continuing with default sql statements for revoking users
		var resp *logical.Response
		resp = &logical.Response{}
		resp.AddWarning("Role " + roleName + "cannot be found. Using default SQL for revoking user")

		// return non-nil response and nil error
		return resp, nil
	}

	return nil, nil
}
