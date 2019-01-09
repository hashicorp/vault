package postgresql

import (
	"context"
	"fmt"
	"strings"
	"time"

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
	// Unlike some other backends we need a lease here (can't leave as 0 and
	// let core fill it in) because Postgres also expires users as a safety
	// measure, so cannot be zero
	if lease == nil {
		lease = &configLease{
			Lease: b.System().DefaultLeaseTTL(),
		}
	}

	// Generate the username, password and expiration. PG limits user to 63 characters
	displayName := req.DisplayName
	if len(displayName) > 26 {
		displayName = displayName[:26]
	}
	userUUID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	username := fmt.Sprintf("%s-%s", displayName, userUUID)
	if len(username) > 63 {
		username = username[:63]
	}
	password, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	ttl, _, err := framework.CalculateTTL(b.System(), 0, lease.Lease, 0, lease.LeaseMax, 0, time.Time{})
	if err != nil {
		return nil, err
	}
	expiration := time.Now().
		Add(ttl).
		Format("2006-01-02 15:04:05-0700")

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
	defer func() {
		tx.Rollback()
	}()

	// Execute each query
	for _, query := range strutil.ParseArbitraryStringSlice(role.SQL, ";") {
		query = strings.TrimSpace(query)
		if len(query) == 0 {
			continue
		}

		m := map[string]string{
			"name":       username,
			"password":   password,
			"expiration": expiration,
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
