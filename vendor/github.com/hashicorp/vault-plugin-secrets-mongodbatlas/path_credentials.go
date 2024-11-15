// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-secure-stdlib/base62"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

var displayNameRegex = regexp.MustCompile("[^a-zA-Z0-9+=,.@_-]")

func (b *Backend) pathCredentials() *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("name"),
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixMongoDBAtlas,
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

func (b *Backend) pathCredentialsRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	role := d.Get("name").(string)

	cred, err := b.credentialRead(ctx, req.Storage, role)
	if err != nil {
		return nil, errwrap.Wrapf("error retrieving credential: {{err}}", err)
	}

	if cred == nil {
		return nil, errors.New("error retrieving credential: credential is nil")
	}

	return b.programmaticAPIKeyCreate(ctx, req.Storage, role, cred)
}

type walEntry struct {
	Role                 string
	ProjectID            string `mapstructure:"project_id"`
	OrganizationID       string `mapstructure:"organization_id"`
	ProgrammaticAPIKeyID string `mapstructure:"programmatic_api_key_id"`
}

func genAPIKeyDescription(displayName string) (string, error) {
	midString := displayNameRegex.ReplaceAllString(displayName, "_")

	id, err := base62.Random(20)
	if err != nil {
		return "", err
	}
	ret := fmt.Sprintf("vault-%s-%s", midString, id)
	return ret, nil
}

const pathCredentialsHelpSyn = `
Generate MongoDB Atlas Programmatic API from a specific Vault role.
`

const pathCredentialsHelpDesc = `
This path reads generates MongoDB Atlas Programmatic API Keys for
a particular role. Atlas Programmatic API Keys will be
generated on demand and will be automatically revoked when
the lease is up.
`
