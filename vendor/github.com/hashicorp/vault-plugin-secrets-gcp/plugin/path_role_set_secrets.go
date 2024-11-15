// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpsecrets

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func fieldSchemaRoleSetServiceAccountKey() map[string]*framework.FieldSchema {
	return map[string]*framework.FieldSchema{
		"roleset": {
			Type:        framework.TypeString,
			Description: "Required. Name of the role set.",
		},
		"key_algorithm": {
			Type:        framework.TypeString,
			Description: fmt.Sprintf(`Private key algorithm for service account key - defaults to %s"`, keyAlgorithmRSA2k),
			Default:     keyAlgorithmRSA2k,
			Query:       true,
		},
		"key_type": {
			Type:        framework.TypeString,
			Description: fmt.Sprintf(`Private key type for service account key - defaults to %s"`, privateKeyTypeJson),
			Default:     privateKeyTypeJson,
			Query:       true,
		},
		"ttl": {
			Type:        framework.TypeDurationSecond,
			Description: "Lifetime of the service account key",
			Query:       true,
		},
	}
}

func fieldSchemaRoleSetAccessToken() map[string]*framework.FieldSchema {
	return map[string]*framework.FieldSchema{
		"roleset": {
			Type:        framework.TypeString,
			Description: "Required. Name of the role set.",
		},
	}

}

func pathRoleSetSecretServiceAccountKey(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("roleset/%s/key", framework.GenericNameRegex("roleset")),
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloud,
			OperationVerb:   "generate",
		},
		Fields:         fieldSchemaRoleSetServiceAccountKey(),
		ExistenceCheck: b.pathRoleSetExistenceCheck("roleset"),
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathRoleSetSecretKey,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "roleset-key2",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathRoleSetSecretKey,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "roleset-key",
				},
			},
		},
		HelpSynopsis:    pathServiceAccountKeySyn,
		HelpDescription: pathServiceAccountKeyDesc,
	}
}

func deprecatedPathRoleSetSecretServiceAccountKey(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("key/%s", framework.GenericNameRegex("roleset")),
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloud,
			OperationVerb:   "generate",
		},
		Deprecated:     true,
		Fields:         fieldSchemaRoleSetServiceAccountKey(),
		ExistenceCheck: b.pathRoleSetExistenceCheck("roleset"),
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathRoleSetSecretKey,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "roleset-key4",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathRoleSetSecretKey,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "roleset-key3",
				},
			},
		},
		HelpSynopsis:    pathServiceAccountKeySyn,
		HelpDescription: pathServiceAccountKeyDesc,
	}
}

func pathRoleSetSecretAccessToken(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("roleset/%s/token", framework.GenericNameRegex("roleset")),
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloud,
			OperationVerb:   "generate",
		},
		Fields:         fieldSchemaRoleSetAccessToken(),
		ExistenceCheck: b.pathRoleSetExistenceCheck("roleset"),
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathRoleSetSecretAccessToken,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "roleset-access-token2",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathRoleSetSecretAccessToken,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "roleset-access-token",
				},
			},
		},
		HelpSynopsis:    pathTokenHelpSyn,
		HelpDescription: pathTokenHelpDesc,
	}
}

func deprecatedPathRoleSetSecretAccessToken(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("token/%s", framework.GenericNameRegex("roleset")),
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloud,
			OperationVerb:   "generate",
		},
		Deprecated:     true,
		Fields:         fieldSchemaRoleSetAccessToken(),
		ExistenceCheck: b.pathRoleSetExistenceCheck("roleset"),
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathRoleSetSecretAccessToken,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "roleset-access-token4",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathRoleSetSecretAccessToken,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "roleset-access-token3",
				},
			},
		},
		HelpSynopsis:    pathTokenHelpSyn,
		HelpDescription: pathTokenHelpDesc,
	}
}

func (b *backend) pathRoleSetSecretKey(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	rsName := d.Get("roleset").(string)
	keyType := d.Get("key_type").(string)
	keyAlg := d.Get("key_algorithm").(string)
	ttl := d.Get("ttl").(int)

	rs, err := getRoleSet(rsName, ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return logical.ErrorResponse("role set %q does not exists", rsName), nil
	}

	if rs.SecretType != SecretTypeKey {
		return logical.ErrorResponse("role set %q cannot generate service account keys (has secret type %s)", rsName, rs.SecretType), nil
	}

	params := secretKeyParams{
		keyType:      keyType,
		keyAlgorithm: keyAlg,
		ttl:          ttl,
		extraInternalData: map[string]interface{}{
			"role_set":          rs.Name,
			"role_set_bindings": rs.bindingHash(),
		},
	}

	return b.createServiceAccountKeySecret(ctx, req.Storage, rs.AccountId, params)
}

func (b *backend) pathRoleSetSecretAccessToken(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	rsName := d.Get("roleset").(string)

	rs, err := getRoleSet(rsName, ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return logical.ErrorResponse("role set '%s' does not exists", rsName), nil
	}

	if rs.SecretType != SecretTypeAccessToken {
		return logical.ErrorResponse("role set '%s' cannot generate access tokens (has secret type %s)", rsName, rs.SecretType), nil
	}

	return b.secretAccessTokenResponse(ctx, req.Storage, rs.TokenGen)
}
