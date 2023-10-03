// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	paths "path"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathPreauthTest() *framework.Path {
	return &framework.Path{
		Pattern: "preauth-test",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixTransit,
		},
		Fields: map[string]*framework.FieldSchema{
			"accessor": {
				Type:     framework.TypeString,
				Required: false,
			},
			"path": {
				Type:     framework.TypeString,
				Required: false,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.handlePreauthTest,
				Summary:  "Returns the size of the active cache",
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "cache-configuration",
				},
			},
		},
	}
}

func (b *backend) handlePreauthTest(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if req.ClientTokenSource != logical.ClientTokenFromInternalAuth {
		da := logical.NewDelegatedAuthenticationError(d.Get("accessor").(string), paths.Join(d.Get("path").(string), req.Data["username"].(string)), nil)
		delete(req.Data, "username")
		return nil, da
	} else {
		delete(req.Data, "password")
	}
	return &logical.Response{}, nil
}
