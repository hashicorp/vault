package database

import (
	"context"

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

		conf, err := b.DatabaseConfig(ctx, req.Storage, name)
		if err != nil {
			return nil, err
		}

		// Close plugin and delete the entry in the connections cache.
		db, err := b.GetConnection(ctx, req.Storage, name)
		if err != nil {
			return nil, err
		}

		_, err = db.RotateRootCredentials(ctx, conf.RootCredentialsRotateStatements)
		if err != nil {
			return nil, err
		}

		// SAVE CONF

		return nil, nil
	}
}
