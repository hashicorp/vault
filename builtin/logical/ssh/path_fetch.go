// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ssh

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathFetchPublicKey(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `public_key`,

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixSSH,
			OperationSuffix: "public-key",
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathFetchPublicKey,
		},

		HelpSynopsis:    `Retrieve the public key.`,
		HelpDescription: `This allows the public key of the SSH CA certificate that this backend has been configured with to be fetched. This is a raw response endpoint without JSON encoding; use -format=raw or an external tool (e.g., curl) to fetch this value.`,
	}
}

func (b *backend) pathFetchPublicKey(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	publicKeyEntry, err := caKey(ctx, req.Storage, caPublicKey)
	if err != nil {
		return nil, err
	}
	if publicKeyEntry == nil || publicKeyEntry.Key == "" {
		return nil, nil
	}

	response := &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPContentType: "text/plain",
			logical.HTTPRawBody:     []byte(publicKeyEntry.Key),
			logical.HTTPStatusCode:  200,
		},
	}

	return response, nil
}
