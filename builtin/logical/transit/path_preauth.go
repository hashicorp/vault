// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"

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
	if req.ClientTokenSource == logical.NoClientToken {
		return nil, logical.NewDelegatedAuthenticationError(d.Get("accessor").(string), d.Get("path").(string))
	}
	return &logical.Response{}, nil
}
