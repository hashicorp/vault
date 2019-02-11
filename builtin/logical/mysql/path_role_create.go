package mysql

import (
	"context"
	"fmt"
	"strings"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/dbtxn"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	_ "github.com/lib/pq"
)

func pathRoleCreate(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathRoleCreateRead,
		},

		HelpSynopsis:    pathRoleCreateReadHelpSyn,
		HelpDescription: pathRoleCreateReadHelpDesc,
	}
}

func (b *backend) pathRoleCreateRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)

	// Get the role
	role, err := b.Role(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown role: %s", name)), nil
	}

	// Determine if we have a lease
	lease, err := b.Lease(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		lease = &configLease{}
	}

	// Generate our username and password. The username will be a
	// concatenation of:
	//
	// - the role name, truncated to role.rolenameLength (default 4)
	// - the token display name, truncated to role.displaynameLength (default 4)
	// - a UUID
	//
	// the entire concatenated string is then truncated to role.usernameLength,
	// which by default is 16 due to limitations in older but still-prevalent
	// versions of MySQL.
	roleName := name
	if len(roleName) > role.RolenameLength {
		roleName = roleName[:role.RolenameLength]
	}
	displayName := req.DisplayName
	if len(displayName) > role.DisplaynameLength {
		displayName = displayName[:role.DisplaynameLength]
	}
	userUUID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	username := fmt.Sprintf("%s-%s-%s", roleName, displayName, userUUID)
	if len(username) > role.UsernameLength {
		username = username[:role.UsernameLength]
	}
	password, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	// Get our handle
	db, err := b.DB(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Execute each query
	for _, query := range strutil.ParseArbitraryStringSlice(role.SQL, ";") {
		query = strings.TrimSpace(query)
		if len(query) == 0 {
			continue
		}

		m := map[string]string{
			"name":     username,
			"password": password,
		}
		if err := dbtxn.ExecuteTxQuery(ctx, tx, m, query); err != nil {
			return nil, err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Return the secret
	resp := b.Secret(SecretCredsType).Response(map[string]interface{}{
		"username": username,
		"password": password,
	}, map[string]interface{}{
		"username": username,
		"role":     name,
	})

	resp.Secret.TTL = lease.Lease
	resp.Secret.MaxTTL = lease.LeaseMax

	return resp, nil
}

const pathRoleCreateReadHelpSyn = `
Request database credentials for a certain role.
`

const pathRoleCreateReadHelpDesc = `
This path reads database credentials for a certain role. The
database credentials will be generated on demand and will be automatically
revoked when the lease is up.
`
