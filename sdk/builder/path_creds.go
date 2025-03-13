package builder

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// pathRole extends the Vault API with a `/role`
// endpoint for the backend. You can choose whether
// or not certain attributes should be displayed,
// required, and named. You can also define different
// path patterns to list all roles.
func (gb *GenericBackend[O, C, R]) pathCredentials() *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeLowerCaseString,
				Description: "Name of the role",
				Required:    true,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   gb.pathCredentialsRead,
			logical.UpdateOperation: gb.pathCredentialsRead,
		},
		HelpSynopsis:    "helllp",
		HelpDescription: "get some help man",
	}
}

// pathRolesRead makes a request to Vault storage to read a role and return response data
func (gb *GenericBackend[O, C, R]) pathCredentialsRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roleName := d.Get("name").(string)

	roleEntry, err := gb.getRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving role: %w", err)
	}

	if roleEntry == nil {
		return nil, errors.New("error retrieving role: role is nil")
	}

	client, err := gb.getClient(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	return gb.role.Secret.FetchSecretFunc(req, d, client, roleEntry, gb.Secret(gb.role.Secret.Type))
}
