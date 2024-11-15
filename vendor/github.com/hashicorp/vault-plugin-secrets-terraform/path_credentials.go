// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfc

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathCredentials(b *tfBackend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("name"),
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixTerraformCloud,
			OperationVerb:   "generate",
		},
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeLowerCaseString,
				Description: "Name of the role",
				Required:    true,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathCredentialsRead,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "credentials",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathCredentialsRead,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "credentials2",
				},
			},
		},

		HelpSynopsis:    pathCredentialsHelpSyn,
		HelpDescription: pathCredentialsHelpDesc,
	}
}

func (b *tfBackend) terraformToken() *framework.Secret {
	return &framework.Secret{
		Type: terraformTokenType,
		Fields: map[string]*framework.FieldSchema{
			"token": {
				Type:        framework.TypeString,
				Description: "Terraform Token",
			},
		},
		Revoke: b.terraformTokenRevoke,
		Renew:  b.terraformTokenRenew,
	}
}

func (b *tfBackend) pathCredentialsRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roleName := d.Get("name").(string)

	roleEntry, err := b.getRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving role: %w", err)
	}

	if roleEntry == nil {
		return nil, errors.New("error retrieving role: role is nil")
	}

	if roleEntry.UserID != "" {
		return b.createUserCreds(ctx, req, roleEntry)
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"token_id":     roleEntry.TokenID,
			"token":        roleEntry.Token,
			"organization": roleEntry.Organization,
			"team_id":      roleEntry.TeamID,
			"role":         roleEntry.Name,
		},
	}
	return resp, nil
}

func (b *tfBackend) createUserCreds(ctx context.Context, req *logical.Request, role *terraformRoleEntry) (*logical.Response, error) {
	token, err := b.createToken(ctx, req.Storage, role)
	if err != nil {
		return nil, err
	}

	resp := b.Secret(terraformTokenType).Response(map[string]interface{}{
		"token":    token.Token,
		"token_id": token.ID,
	}, map[string]interface{}{
		"token_id": token.ID,
		"role":     role.Name,
	})

	if role.TTL > 0 {
		resp.Secret.TTL = role.TTL
	}

	if role.MaxTTL > 0 {
		resp.Secret.MaxTTL = role.MaxTTL
	}

	return resp, nil
}

func (b *tfBackend) createToken(ctx context.Context, s logical.Storage, roleEntry *terraformRoleEntry) (*terraformToken, error) {
	client, err := b.getClient(ctx, s)
	if err != nil {
		return nil, err
	}

	var token *terraformToken

	switch {
	case isOrgToken(roleEntry.Organization, roleEntry.TeamID):
		token, err = createOrgToken(ctx, client, roleEntry.Organization)
	case isTeamToken(roleEntry.TeamID):
		token, err = createTeamToken(ctx, client, roleEntry.TeamID)
	default:
		token, err = createUserToken(ctx, client, roleEntry.UserID)
	}

	if err != nil {
		return nil, fmt.Errorf("error creating Terraform token: %w", err)
	}

	if token == nil {
		return nil, errors.New("error creating Terraform token")
	}

	return token, nil
}

const pathCredentialsHelpSyn = `
Generate a Terraform Cloud or Enterprise API token from a specific Vault role.
`

const pathCredentialsHelpDesc = `
This path generates Terraform Cloud or Enterprise API Organization, Team, or
User Tokens based on a particular role. A role can only represent a single type
of Token; Organization, Team, or User, and so can only contain one parameter for
organization, team_id, or user_id.

If the role has the team ID configured, this path generates a team token.

If this role only has the organization configured, this path generates an
organization token.

If this role has a user ID configured, this path generates a user token.
`
