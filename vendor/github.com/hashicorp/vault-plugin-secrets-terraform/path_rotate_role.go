// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfc

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathRotateRole(b *tfBackend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "rotate-role/" + framework.GenericNameRegex("name"),

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixTerraformCloud,
				OperationVerb:   "rotate",
				OperationSuffix: "role",
			},

			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the team or organization role",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback:                    b.pathRotateRole,
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: true,
				},
			},

			HelpSynopsis:    pathRotateRoleHelpSyn,
			HelpDescription: pathRotateRoleHelpDesc,
		},
	}
}

func (b *tfBackend) pathRotateRole(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse("missing role name"), nil
	}

	roleEntry, err := b.getRole(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}

	if roleEntry == nil {
		return logical.ErrorResponse("missing role entry"), nil
	}

	if roleEntry.UserID != "" {
		return logical.ErrorResponse("cannot rotate credentials for user roles"), nil
	}

	token, err := b.createToken(ctx, req.Storage, roleEntry)
	if err != nil {
		return nil, err
	}

	roleEntry.Token = token.Token

	if err := setRole(ctx, req.Storage, name, roleEntry); err != nil {
		return nil, err
	}

	return nil, nil
}

const pathRotateRoleHelpSyn = `
Request to rotate the credentials for a team or organization.
`

const pathRotateRoleHelpDesc = `
This path attempts to rotate the credentials for the given team or organization role. 
This endpoint returns an error if attempting to rotate a user role.
`
