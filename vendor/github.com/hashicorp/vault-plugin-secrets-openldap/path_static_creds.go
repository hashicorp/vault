// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package openldap

import (
	"context"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const staticCredPath = "static-cred/"

func (b *backend) pathStaticCredsCreate() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: strings.TrimSuffix(staticCredPath, "/") + genericNameWithForwardSlashRegex("name"),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixLDAP,
				OperationVerb:   "request",
				OperationSuffix: "static-role-credentials",
			},
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeLowerCaseString,
					Description: "Name of the static role.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathStaticCredsRead,
				},
			},
			HelpSynopsis:    pathStaticCredsReadHelpSyn,
			HelpDescription: pathStaticCredsReadHelpDesc,
		},
	}
}

func (b *backend) pathStaticCredsRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)

	role, err := b.staticRole(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse("unknown role: %s", name), nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"dn":                  role.StaticAccount.DN,
			"username":            role.StaticAccount.Username,
			"password":            role.StaticAccount.Password,
			"last_password":       role.StaticAccount.LastPassword,
			"ttl":                 role.StaticAccount.PasswordTTL().Seconds(),
			"rotation_period":     role.StaticAccount.RotationPeriod.Seconds(),
			"last_vault_rotation": role.StaticAccount.LastVaultRotation,
		},
	}, nil
}

const pathStaticCredsReadHelpSyn = `
Request LDAP credentials for a certain static role. These credentials are
rotated periodically.`

const pathStaticCredsReadHelpDesc = `
This path reads LDAP credentials for a certain static role. The LDAPs 
credentials are rotated periodically according to their configuration, and will
return the same password until they are rotated.
`
