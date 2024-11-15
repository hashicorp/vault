// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpkms

import (
	"context"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathKeysConfigCRUD() *framework.Path {
	return &framework.Path{
		Pattern: "keys/config/" + framework.GenericNameRegex("key"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloudKMS,
			OperationVerb:   "configure",
			OperationSuffix: "key",
		},

		HelpSynopsis: "Configure the key in Vault",
		HelpDescription: `
Update the Vault's configuration of this key such as the minimum allowed key
version and other metadata.
`,

		Fields: map[string]*framework.FieldSchema{
			"key": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Name of the key in Vault.
`,
			},

			"min_version": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `
Minimum allowed crypto key version. If set to a positive value, key versions
less than the given value are not permitted to be used. If set to 0 or a
negative value, there is no minimum key version. This value only affects
encryption/re-encryption, not decryption. To restrict old values from being
decrypted, increase this value and then perform a trim operation.
`,
			},

			"max_version": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `
Maximum allowed crypto key version. If set to a positive value, key versions
greater than the given value are not permitted to be used. If set to 0 or a
negative value, there is no maximum key version.
`,
			},
		},

		ExistenceCheck: b.pathKeysExistenceCheck,

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: withFieldValidator(b.pathKeysConfigRead),
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "read",
					OperationSuffix: "key-configuration",
				},
			},
			logical.CreateOperation: &framework.PathOperation{
				Callback: withFieldValidator(b.pathKeysConfigWrite),
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "key",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: withFieldValidator(b.pathKeysConfigWrite),
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "key",
				},
			},
		},
	}
}

// pathKeysConfigRead corresponds to GET gcpkms/keys/config/:name and is used to
// show information about the key configuration in Vault.
func (b *backend) pathKeysConfigRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	key := d.Get("key").(string)

	k, err := b.Key(ctx, req.Storage, key)
	if err != nil {
		if err == ErrKeyNotFound {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		return nil, err
	}

	data := map[string]interface{}{
		"name":       k.Name,
		"crypto_key": k.CryptoKeyID,
	}

	if k.MinVersion > 0 {
		data["min_version"] = k.MinVersion
	}

	if k.MaxVersion > 0 {
		data["max_version"] = k.MaxVersion
	}

	return &logical.Response{
		Data: data,
	}, nil
}

// pathKeysConfigWrite corresponds to PUT/POST gcpkms/keys/config/:key and
// configures information about the key in Vault.
func (b *backend) pathKeysConfigWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	key := d.Get("key").(string)

	k, err := b.Key(ctx, req.Storage, key)
	if err != nil {
		if err == ErrKeyNotFound {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		return nil, err
	}

	if v, ok := d.GetOk("min_version"); ok {
		if v.(int) <= 0 {
			k.MinVersion = 0
		} else {
			k.MinVersion = v.(int)
		}
	}

	if v, ok := d.GetOk("max_version"); ok {
		if v.(int) <= 0 {
			k.MaxVersion = 0
		} else {
			k.MaxVersion = v.(int)
		}
	}

	// Save it
	entry, err := logical.StorageEntryJSON("keys/"+key, k)
	if err != nil {
		return nil, errwrap.Wrapf("failed to create storage entry: {{err}}", err)
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, errwrap.Wrapf("failed to write to storage: {{err}}", err)
	}

	return nil, nil
}
