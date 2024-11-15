// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpkms

import (
	"context"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

func (b *backend) pathKeysRegister() *framework.Path {
	return &framework.Path{
		Pattern: "keys/register/" + framework.GenericNameRegex("key"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloudKMS,
			OperationVerb:   "register",
			OperationSuffix: "key",
		},

		HelpSynopsis: "Register an existing crypto key in Google Cloud KMS",
		HelpDescription: `
Registers an existing crypto key in Google Cloud KMS and make it available for
encryption and decryption in Vault.

To have Vault create a crypto key, use the create method instead. This function
is for existing crypto keys which you now want to manage via Vault.
`,

		Fields: map[string]*framework.FieldSchema{
			"key": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Name of the key to register in Vault. This will be the named used to refer to
the underlying crypto key when encrypting or decrypting data.
`,
			},

			"crypto_key": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Full resource ID of the crypto key including the project, location, key ring,
and crypto key like "projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s". This
crypto key must already exist in Google Cloud KMS unless verify is set to
"false".
`,
			},

			"verify": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: true,
				Description: `
Verify that the given Google Cloud KMS crypto key exists and is accessible
before creating the storage entry in Vault. Set this to "false" if the key will
not exist at creation time.
`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: withFieldValidator(b.pathKeysRegisterWrite),
		},
	}
}

// pathKeysRegisterWrite corresponds to PUT/POST gcpkms/keys/register/:key and
// registers an existing GCP KMS key for use in Vault.
func (b *backend) pathKeysRegisterWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	key := d.Get("key").(string)
	cryptoKey := d.Get("crypto_key").(string)
	verify := d.Get("verify").(bool)

	if verify {
		kmsClient, closer, err := b.KMSClient(req.Storage)
		if err != nil {
			return nil, err
		}
		defer closer()

		if _, err := kmsClient.GetCryptoKey(ctx, &kmspb.GetCryptoKeyRequest{
			Name: cryptoKey,
		}); err != nil {
			return nil, errwrap.Wrapf("failed to read crypto key: {{err}}", err)
		}
	}

	entry, err := logical.StorageEntryJSON("keys/"+key, &Key{
		Name:        key,
		CryptoKeyID: cryptoKey,
	})
	if err != nil {
		return nil, errwrap.Wrapf("failed to create storage entry: {{err}}", err)
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, errwrap.Wrapf("failed to write to storage: {{err}}", err)
	}

	return nil, nil
}
