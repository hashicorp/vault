// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpkms

import (
	"context"
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

func (b *backend) pathPubkey() *framework.Path {
	return &framework.Path{
		Pattern: "pubkey/" + framework.GenericNameRegex("key"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloudKMS,
			OperationVerb:   "retrieve",
			OperationSuffix: "public-key",
		},

		HelpSynopsis: "Retrieve the public key associated with the named key",
		HelpDescription: `
Retrieve the PEM-encoded Google Cloud KMS public key associated with the Vault
named key. The key will only be available if the key is asymmetric.
`,

		Fields: map[string]*framework.FieldSchema{
			"key": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Name of the key for which to get the public key. This key must already exist in
Vault and Google Cloud KMS.
`,
			},

			"key_version": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `
Integer version of the crypto key version from which to exact the public key.
This field is required.
`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: withFieldValidator(b.pathPubkeyRead),
		},
	}
}

// pathPubkeyRead corresponds to GET gcpkms/pubkey/:key and is used to read the
// public key contents of the crypto key version.
func (b *backend) pathPubkeyRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	key := d.Get("key").(string)
	keyVersion := d.Get("key_version").(int)

	if keyVersion == 0 {
		return nil, errMissingFields("key_version")
	}

	k, err := b.Key(ctx, req.Storage, key)
	if err != nil {
		if err == ErrKeyNotFound {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		return nil, err
	}

	if k.MinVersion > 0 && keyVersion < k.MinVersion {
		resp := fmt.Sprintf("requested version %d is less than minimum allowed version of %d",
			keyVersion, k.MinVersion)
		return logical.ErrorResponse(resp), logical.ErrPermissionDenied
	}

	if k.MaxVersion > 0 && keyVersion > k.MaxVersion {
		resp := fmt.Sprintf("requested version %d is greater than maximum allowed version of %d",
			keyVersion, k.MaxVersion)
		return logical.ErrorResponse(resp), logical.ErrPermissionDenied
	}

	kmsClient, closer, err := b.KMSClient(req.Storage)
	if err != nil {
		return nil, err
	}
	defer closer()

	pk, err := kmsClient.GetPublicKey(ctx, &kmspb.GetPublicKeyRequest{
		Name: fmt.Sprintf("%s/cryptoKeyVersions/%d", k.CryptoKeyID, keyVersion),
	})
	if err != nil {
		return nil, errwrap.Wrapf("failed to get public key: {{err}}", err)
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"pem":       pk.Pem,
			"algorithm": algorithmToString(pk.Algorithm),
		},
	}, nil
}
