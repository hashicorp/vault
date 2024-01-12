// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ssh

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const keysStoragePrefix = "keys/"

func pathCleanupKeys(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "tidy/dynamic-keys",
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixSSH,
			OperationVerb:   "tidy",
			OperationSuffix: "dynamic-host-keys",
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: b.handleCleanupKeys,
		},
		HelpSynopsis:    `This endpoint removes the stored host keys used for the removed Dynamic Key feature, if present.`,
		HelpDescription: `For more information, refer to the API documentation.`,
	}
}

func (b *backend) handleCleanupKeys(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	names, err := req.Storage.List(ctx, keysStoragePrefix)
	if err != nil {
		return nil, fmt.Errorf("unable to list keys for removal: %w", err)
	}

	for index, name := range names {
		keyPath := keysStoragePrefix + name
		if err := req.Storage.Delete(ctx, keyPath); err != nil {
			return nil, fmt.Errorf("unable to delete key %v of %v: %w", index+1, len(names), err)
		}
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"message": fmt.Sprintf("Removed %v of %v host keys.", len(names), len(names)),
		},
	}, nil
}
