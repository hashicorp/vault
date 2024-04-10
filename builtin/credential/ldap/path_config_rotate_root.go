// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-2.0

package ldap

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathConfigRotateRoot(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/rotate-root",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixLDAP,
			OperationVerb:   "rotate",
			OperationSuffix: "root-credentials",
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                    b.pathConfigRotateRootUpdate,
				ForwardPerformanceSecondary: true,
				ForwardPerformanceStandby:   true,
			},
		},

		HelpSynopsis:    pathConfigRotateRootHelpSyn,
		HelpDescription: pathConfigRotateRootHelpDesc,
	}
}

func (b *backend) pathConfigRotateRootUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// check for things that are logical.Response error cases first, since rotatePassword will only return an error
	b.mu.RLock()

	cfg, err := b.Config(ctx, req)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return logical.ErrorResponse("attempted to rotate root on an undefined config"), nil
	}

	u, p := cfg.BindDN, cfg.BindPassword
	if u == "" || p == "" {
		return logical.ErrorResponse("auth is not using authenticated search, no root to rotate"), nil
	}

	b.mu.RUnlock()

	err = b.RotateCredential(ctx, req)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

const pathConfigRotateRootHelpSyn = `
Request to rotate the LDAP credentials used by Vault
`

const pathConfigRotateRootHelpDesc = `
This path attempts to rotate the LDAP bindpass used by Vault for this mount.
`
