package database

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathRoleCreate(b *databaseBackend) *framework.Path {
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

func (b *databaseBackend) pathRoleCreateRead(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)

	// Get the role
	role, err := b.Role(req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("Unknown role: %s", name)), nil
	}

	b.Lock()
	defer b.Unlock()

	// Get the Database object
	db, err := b.getOrCreateDBObj(req.Storage, role.DBName)
	if err != nil {
		// TODO: return a resp error instead?
		return nil, fmt.Errorf("cound not retrieve db with name: %s, got error: %s", role.DBName, err)
	}

	// Generate the username, password and expiration
	username, err := db.GenerateUsername(req.DisplayName)
	if err != nil {
		return nil, err
	}

	password, err := db.GeneratePassword()
	if err != nil {
		return nil, err
	}

	expiration, err := db.GenerateExpiration(role.DefaultTTL)
	if err != nil {
		return nil, err
	}

	// Create the user
	err = db.CreateUser(role.Statements, username, password, expiration)
	if err != nil {
		return nil, err
	}

	resp := b.Secret(SecretCredsType).Response(map[string]interface{}{
		"username": username,
		"password": password,
	}, map[string]interface{}{
		"username": username,
		"role":     name,
	})
	resp.Secret.TTL = role.DefaultTTL
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
