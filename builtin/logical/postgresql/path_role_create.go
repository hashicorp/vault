package postgresql

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	_ "github.com/lib/pq"
)

func pathRoleCreate(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `creds/(?P<name>\w+)`,
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

func (b *backend) pathRoleCreateRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)

	// Get the role
	role, err := b.Role(req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown role: %s", name)), nil
	}

	// Determine if we have a lease
	lease, err := b.Lease(req.Storage)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		lease = &configLease{Lease: 1 * time.Hour}
	}

	// Get our connection
	db, err := b.DB(req.Storage)
	if err != nil {
		return nil, err
	}

	// Generate our query
	username := fmt.Sprintf(
		"vault-%s-%d-%d",
		req.DisplayName, time.Now().Unix(), rand.Int31n(10000))
	password := generateUUID()
	expiration := time.Now().UTC().
		Add(lease.Lease + time.Duration((float64(lease.Lease) * 0.1))).
		Format("2006-01-02 15:04:05")
	query := Query(role.SQL, map[string]string{
		"name":       username,
		"password":   password,
		"expiration": expiration,
	})

	// Prepare the statement and execute it
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(); err != nil {
		return nil, err
	}

	// Return the secret
	resp := b.Secret(SecretCredsType).Response(map[string]interface{}{
		"username": username,
		"password": password,
	}, map[string]interface{}{
		"username": username,
	})
	resp.Secret.Lease = lease.Lease
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
