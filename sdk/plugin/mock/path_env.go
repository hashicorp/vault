// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mock

import (
	"context"
	"os"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// pathEnv is used to interrogate plugin env vars.
func pathEnv(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "env/" + framework.GenericNameRegex("key"),
		Fields: map[string]*framework.FieldSchema{
			"key": {
				Type:        framework.TypeString,
				Required:    true,
				Description: "The name of the environment variable to read.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathEnvRead,
		},
	}
}

func (b *backend) pathEnvRead(_ context.Context, _ *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Return the secret
	return &logical.Response{
		Data: map[string]interface{}{
			"key": os.Getenv(data.Get("key").(string)),
		},
	}, nil
}
