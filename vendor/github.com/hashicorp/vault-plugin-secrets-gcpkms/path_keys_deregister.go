// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpkms

import (
	"context"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathKeysDeregister() *framework.Path {
	return &framework.Path{
		Pattern: "keys/deregister/" + framework.GenericNameRegex("key"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloudKMS,
			OperationVerb:   "deregister",
		},

		HelpSynopsis: "Deregister an existing key in Vault",
		HelpDescription: `
This endpoint deregisters an existing reference Vault has to a crypto key in
Google Cloud KMS. The underlying Google Cloud KMS key remains unchanged.
`,

		Fields: map[string]*framework.FieldSchema{
			"key": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Name of the key to deregister in Vault. If the key exists in Google Cloud KMS,
it will be left untouched.
`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: withFieldValidator(b.pathKeysDeregisterWrite),
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "key",
				},
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: withFieldValidator(b.pathKeysDeregisterWrite),
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "key2",
				},
			},
		},
	}
}

// pathKeysDeregisterWrite corresponds to POST gcpkms/keys/deregister/:key
// and deregisters a key for use in Vault. It does not delete or disable the
// underlying GCP KMS keys.
func (b *backend) pathKeysDeregisterWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	key := d.Get("key").(string)

	if err := req.Storage.Delete(ctx, "keys/"+key); err != nil {
		return nil, errwrap.Wrapf("failed to delete from storage: {{err}}", err)
	}
	return nil, nil
}
