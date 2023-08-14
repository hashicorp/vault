// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathRotate() *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("name") + "/rotate",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixTransit,
			OperationVerb:   "rotate",
			OperationSuffix: "key",
		},

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the key",
			},
			"managed_key_name": {
				Type:        framework.TypeString,
				Description: "The name of the managed key to use for the new version of this transit key",
			},
			"managed_key_id": {
				Type:        framework.TypeString,
				Description: "The UUID of the managed key to use for the new version of this transit key",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathRotateWrite,
		},

		HelpSynopsis:    pathRotateHelpSyn,
		HelpDescription: pathRotateHelpDesc,
	}
}

func (b *backend) pathRotateWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	managedKeyName := d.Get("managed_key_name").(string)
	managedKeyId := d.Get("managed_key_id").(string)

	// Get the policy
	p, _, err := b.GetPolicy(ctx, keysutil.PolicyRequest{
		Storage: req.Storage,
		Name:    name,
	}, b.GetRandomReader())
	if err != nil {
		return nil, err
	}
	if p == nil {
		return logical.ErrorResponse("key not found"), logical.ErrInvalidRequest
	}
	if !b.System().CachingDisabled() {
		p.Lock(true)
	}
	defer p.Unlock()

	if p.Type == keysutil.KeyType_MANAGED_KEY {
		var keyId string
		keyId, err = GetManagedKeyUUID(ctx, b, managedKeyName, managedKeyId)
		if err != nil {
			p.Unlock()
			return nil, err
		}
		err = p.RotateManagedKey(ctx, req.Storage, keyId)
	} else {
		// Rotate the policy
		err = p.Rotate(ctx, req.Storage, b.GetRandomReader())
	}

	if err != nil {
		return nil, err
	}

	return b.formatKeyPolicy(p, nil)
}

const pathRotateHelpSyn = `Rotate named encryption key`

const pathRotateHelpDesc = `
This path is used to rotate the named key. After rotation,
new encryption requests using this name will use the new key,
but decryption will still be supported for older versions.
`
