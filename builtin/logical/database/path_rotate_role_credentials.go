package database

import (
	"context"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathRotateRoleCredentials(b *databaseBackend) *framework.Path {
	return &framework.Path{
		Pattern: "rotate-role/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the static role",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathRotateRoleCredentialsUpdate(),
		},

		HelpSynopsis:    pathCredsCreateReadHelpSyn,
		HelpDescription: pathCredsCreateReadHelpDesc,
	}
}

func (b *databaseBackend) pathRotateRoleCredentialsUpdate() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)
		if name == "" {
			return logical.ErrorResponse("empty role name attribute given"), nil
		}

		role, err := b.Role(ctx, req.Storage, data.Get("name").(string))
		if err != nil {
			return nil, err
		}
		if role == nil {
			return nil, nil
		}

		if role.StaticAccount != nil {
			// in create/update of static accounts, we only care if the operation
			// err'd , and this call does not return credentials

			//TODO wrap in WAL, rollback
			// TODO: pop and replace item from queue?
			_, err = b.createUpdateStaticAccount(ctx, req.Storage, &setPasswordInput{
				RoleName: name,
				Role:     role,
			})
			if err != nil {
				return nil, err
			}
		} else {
			return logical.ErrorResponse("cannot rotate credentials of non-static accounts"), nil
		}

		return nil, nil
	}
}

const pathRotateRoleCredentialsUpdateHelpSyn = `
Request to rotate the credentials for a static user account.
`

const pathRotateRoleCredentialsUpdateHelpDesc = `
This path attempts to rotate the credentials for the given static user account. 
`
