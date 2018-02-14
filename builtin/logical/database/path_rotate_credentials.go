package database

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathRotateCredentials(b *databaseBackend) *framework.Path {
	return &framework.Path{
		Pattern: "rotate-root/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of this database connection",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathRotateCredentialsUpdate(),
		},

		HelpSynopsis:    pathCredsCreateReadHelpSyn,
		HelpDescription: pathCredsCreateReadHelpDesc,
	}
}

func (b *databaseBackend) pathRotateCredentialsUpdate() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)
		if name == "" {
			return logical.ErrorResponse(respErrEmptyName), nil
		}

		b.Lock()
		defer b.Unlock()

		config, err := b.DatabaseConfig(ctx, req.Storage, name)
		if err != nil {
			return nil, err
		}

		// Close plugin and delete the entry in the connections cache.
		db, err := b.GetConnection(ctx, req.Storage, name)
		if err != nil {
			return nil, err
		}

		connectionDetails, err := db.RotateRootCredentials(ctx, config.RootCredentialsRotateStatements)
		if err != nil {
			return nil, err
		}

		config.ConnectionDetails = connectionDetails
		entry, err := logical.StorageEntryJSON(fmt.Sprintf("config/%s", name), config)
		if err != nil {
			return nil, err
		}
		if err := req.Storage.Put(ctx, entry); err != nil {
			return nil, err
		}

		if err := b.clearConnection(name); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

const pathRotateCredentialsUpdateHelpSyn = `
Request to rotate the root credentials for a certain database connection.
`

const pathRotateCredentialsUpdateHelpDesc = `
This path attempts to rotate the root credentials for the given database. 
`
