package cassandra

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/helper/uuid"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathCredsCreate(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `creds/(?P<name>\w+)`,
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathCredsCreateRead,
		},

		HelpSynopsis:    pathCredsCreateReadHelpSyn,
		HelpDescription: pathCredsCreateReadHelpDesc,
	}
}

func (b *backend) pathCredsCreateRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)

	// Get the role
	role, err := getRole(req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("Unknown role: %s", name)), nil
	}

	displayName := req.DisplayName
	username := fmt.Sprintf("vault-%s-%s-%s-%d", name, displayName, uuid.GenerateUUID(), time.Now().Unix())
	password := uuid.GenerateUUID()

	// Get our connection
	session, err := b.DB(req.Storage)
	if err != nil {
		return nil, err
	}

	// Execute each query
	for _, query := range splitSQL(role.CreationCQL) {
		err = session.Query(substQuery(query, map[string]string{
			"username": username,
			"password": password,
		})).Exec()
		if err != nil {
			for _, query := range splitSQL(role.RollbackCQL) {
				session.Query(substQuery(query, map[string]string{
					"username": username,
					"password": password,
				})).Exec()
			}
			return nil, err
		}
	}

	// Return the secret
	resp := b.Secret(SecretCredsType).Response(map[string]interface{}{
		"username": username,
		"password": password,
	}, map[string]interface{}{
		"username": username,
		"role":     name,
	})
	resp.Secret.Lease = role.Lease
	resp.Secret.LeaseGracePeriod = role.LeaseGracePeriod

	return resp, nil
}

const pathCredsCreateReadHelpSyn = `
Request database credentials for a certain role.
`

const pathCredsCreateReadHelpDesc = `
This path creates database credentials for a certain role. The
database credentials will be generated on demand and will be automatically
revoked when the lease is up.
`
