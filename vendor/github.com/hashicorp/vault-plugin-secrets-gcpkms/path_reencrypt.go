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

func (b *backend) pathReencrypt() *framework.Path {
	return &framework.Path{
		Pattern: "reencrypt/" + framework.GenericNameRegex("key"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloudKMS,
			OperationVerb:   "reencrypt",
		},

		HelpSynopsis: "Re-encrypt existing ciphertext data to a new version",
		HelpDescription: `
Use the named encryption key to re-encrypt the underlying cryptokey to the latest
version for this ciphertext without disclosing the original plaintext value to
the requestor.
`,

		Fields: map[string]*framework.FieldSchema{
			"key": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Name of the key to use for encryption. This key must already exist in Vault and
Google Cloud KMS.
`,
			},

			"additional_authenticated_data": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Optional data that, if specified, must also be provided during decryption.
`,
			},

			"ciphertext": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Ciphertext to be re-encrypted to the latest key version. This must be ciphertext
that Vault previously generated for this named key.
`,
			},

			"key_version": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `
Integer version of the crypto key version to use for the new encryption. If
unspecified, this defaults to the latest active crypto key version.
`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: withFieldValidator(b.pathReencryptWrite),
		},
	}
}

// pathReencryptWrite corresponds to PUT/POST gcpkms/reencrypt/:key and is
// used to re-encrypt the given ciphertext to the latest cryptokey version.
func (b *backend) pathReencryptWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	key := d.Get("key").(string)
	aad := d.Get("additional_authenticated_data").(string)
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

	// We gave the user back base64-encoded ciphertext in the /encrypt payload
	ciphertext, err := base64.StdEncoding.DecodeString(d.Get("ciphertext").(string))
	if err != nil {
		return nil, errwrap.Wrapf("failed to base64 decode ciphtertext: {{err}}", err)
	}

	kmsClient, closer, err := b.KMSClient(req.Storage)
	if err != nil {
		return nil, err
	}
	defer closer()

	decResp, err := kmsClient.Decrypt(ctx, &kmspb.DecryptRequest{
		Name:                        k.CryptoKeyID, // KMS chooses the version
		Ciphertext:                  ciphertext,
		AdditionalAuthenticatedData: []byte(aad),
	})
	if err != nil {
		return nil, errwrap.Wrapf("failed to decrypt ciphertext: {{err}}", err)
	}

	encResp, err := kmsClient.Encrypt(ctx, &kmspb.EncryptRequest{
		Name:                        cryptoKey, // User-specified version
		Plaintext:                   decResp.Plaintext,
		AdditionalAuthenticatedData: []byte(aad),
	})
	if err != nil {
		return nil, errwrap.Wrapf("failed to encrypt new plaintext: {{err}}", err)
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"key_version": path.Base(encResp.Name),
			"ciphertext":  base64.StdEncoding.EncodeToString(encResp.Ciphertext),
		},
	}, nil
}
