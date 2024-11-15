// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpkms

import (
	"context"
	"path"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

const (
	primaryVersionWarning = "The crypto key version was rotated successfully, " +
		"but it can take up to 2 hours for the new crypto key version to become " +
		"the primary. In practice, it is usually much shorter. To ensure you are " +
		"using this latest key version, specify the key_version attribute during " +
		"operations or issue a read operation and verify the key version has " +
		"propagated."
)

func (b *backend) pathKeysRotate() *framework.Path {
	return &framework.Path{
		Pattern: "keys/rotate/" + framework.GenericNameRegex("key"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloudKMS,
			OperationVerb:   "rotate",
			OperationSuffix: "key",
		},

		HelpSynopsis: "Rotate a crypto key to a new primary version",
		HelpDescription: `
This endpoint creates a new crypto key version for the corresponding Google
Cloud KMS key and updates the new crypto key to be the primary key for future
encryptions.

It can take up to 2 hours for a new crypto key version to become the primary,
so be sure to issue a read operation if you require new data to be encrypted
with this key.
`,

		Fields: map[string]*framework.FieldSchema{
			"key": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Name of the key to rotate. This key must already be registered with Vault and
point to a valid Google Cloud KMS crypto key.
`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: withFieldValidator(b.pathKeysRotateWrite),
		},
	}
}

// pathKeysRotateWrite corresponds to PUT/POST gcpkms/keys/rotate/:name and is
// used to create a new underlying GCP KMS crypto key version and set that
// version to the primary for future encryption.
func (b *backend) pathKeysRotateWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	key := d.Get("key").(string)

	kmsClient, closer, err := b.KMSClient(req.Storage)
	if err != nil {
		return nil, err
	}
	defer closer()

	entry, err := b.Key(ctx, req.Storage, key)
	if err != nil {
		if err == ErrKeyNotFound {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		return nil, err
	}

	// Create a new cyrpto key version
	resp, err := kmsClient.CreateCryptoKeyVersion(ctx, &kmspb.CreateCryptoKeyVersionRequest{
		Parent: entry.CryptoKeyID,
		CryptoKeyVersion: &kmspb.CryptoKeyVersion{
			State: kmspb.CryptoKeyVersion_ENABLED,
		},
	})
	if err != nil {
		return nil, errwrap.Wrapf("failed to create new crypto key version: {{err}}", err)
	}

	// The API expects JUST the version, not the full resource ID
	cryptoKeyVersion := path.Base(resp.Name)

	// Set the new version as primary, only valid for symmetric keys
	if resp.Algorithm == kmspb.CryptoKeyVersion_GOOGLE_SYMMETRIC_ENCRYPTION {
		if _, err := kmsClient.UpdateCryptoKeyPrimaryVersion(ctx, &kmspb.UpdateCryptoKeyPrimaryVersionRequest{
			Name:               entry.CryptoKeyID,
			CryptoKeyVersionId: cryptoKeyVersion,
		}); err != nil {
			return nil, errwrap.Wrapf("failed to update crypto key primary version: {{err}}", err)
		}
	}

	return &logical.Response{
		Warnings: []string{primaryVersionWarning},
		Data: map[string]interface{}{
			"key_version": cryptoKeyVersion,
		},
	}, nil
}
