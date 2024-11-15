// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpkms

import (
	"context"
	"encoding/base64"
	"fmt"
	"path"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

func (b *backend) pathEncrypt() *framework.Path {
	return &framework.Path{
		Pattern: "encrypt/" + framework.GenericNameRegex("key"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloudKMS,
			OperationVerb:   "encrypt",
		},

		HelpSynopsis: "Encrypt a plaintext value using a named key",
		HelpDescription: `
Use the named encryption key to encrypt an arbitrary plaintext string. The
response will be the base64-encoded encrypted value (ciphertext).
`,

		Fields: map[string]*framework.FieldSchema{
			"key": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Name of the key in Vault to use for encryption. This key must already exist in
Vault and must map back to a Google Cloud KMS key.
`,
			},

			"additional_authenticated_data": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Optional base64-encoded data that, if specified, must also be provided to
decrypt this payload.
`,
			},

			"key_version": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `
Integer version of the crypto key version to use for encryption. If unspecified,
this defaults to the latest active crypto key version.
`,
			},

			"plaintext": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Plaintext value to be encrypted. This can be a string or binary, but the size
is limited. See the Google Cloud KMS documentation for information on size
limitations by key types.
`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: withFieldValidator(b.pathEncryptWrite),
		},
	}
}

// pathEncryptWrite corresponds to PUT/POST gcpkms/encrypt/:key and is
// used to encrypt the plaintext string using the named key.
func (b *backend) pathEncryptWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	key := d.Get("key").(string)
	aad := d.Get("additional_authenticated_data").(string)
	plaintext := d.Get("plaintext").(string)
	keyVersion := d.Get("key_version").(int)

	k, err := b.Key(ctx, req.Storage, key)
	if err != nil {
		if err == ErrKeyNotFound {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		return nil, err
	}

	cryptoKey := k.CryptoKeyID
	if keyVersion > 0 {
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

		cryptoKey = fmt.Sprintf("%s/cryptoKeyVersions/%d", cryptoKey, keyVersion)
	}

	kmsClient, closer, err := b.KMSClient(req.Storage)
	if err != nil {
		return nil, err
	}
	defer closer()

	resp, err := kmsClient.Encrypt(ctx, &kmspb.EncryptRequest{
		Name:                        cryptoKey,
		Plaintext:                   []byte(plaintext),
		AdditionalAuthenticatedData: []byte(aad),
	})
	if err != nil {
		return nil, errwrap.Wrapf("failed to encrypt plaintext: {{err}}", err)
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"key_version": path.Base(resp.Name),
			"ciphertext":  base64.StdEncoding.EncodeToString(resp.Ciphertext),
		},
	}, nil
}
