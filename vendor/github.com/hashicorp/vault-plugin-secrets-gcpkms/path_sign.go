// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpkms

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

func (b *backend) pathSign() *framework.Path {
	return &framework.Path{
		Pattern: "sign/" + framework.GenericNameRegex("key"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloudKMS,
			OperationVerb:   "sign",
		},

		HelpSynopsis: "Signs a message or digest using a named key",
		HelpDescription: `
Use the named key to sign a digest string. The response will be the
base64-encoded signature.
`,

		Fields: map[string]*framework.FieldSchema{
			"key": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Name of the key in Vault to use for signing. This key must already exist in
Vault and must map back to a Google Cloud KMS key.
`,
			},

			"digest": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Digest to sign. This digest must use the same SHA algorithm as the underlying
Cloud KMS key. The digest must be the base64-encoded binary value. This field
is required.
`,
			},

			"key_version": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `
Integer version of the crypto key version to use for signing. This field is
required.
`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: withFieldValidator(b.pathSignWrite),
		},
	}
}

// pathSignWrite corresponds to PUT/POST gcpkms/sign/:key and is used to sign
// the digest using the named key.
func (b *backend) pathSignWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	key := d.Get("key").(string)
	digest := d.Get("digest").(string)
	keyVersion := d.Get("key_version").(int)

	if digest == "" {
		return nil, errMissingFields("digest")
	}

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

	ckv, err := kmsClient.GetCryptoKeyVersion(ctx, &kmspb.GetCryptoKeyVersionRequest{
		Name: fmt.Sprintf("%s/cryptoKeyVersions/%d", k.CryptoKeyID, keyVersion),
	})
	if err != nil {
		return nil, errwrap.Wrapf("failed to get underlying crypto key: {{err}}", err)
	}

	var dig *kmspb.Digest

	switch ckv.Algorithm {
	case kmspb.CryptoKeyVersion_RSA_SIGN_PSS_2048_SHA256,
		kmspb.CryptoKeyVersion_RSA_SIGN_PSS_3072_SHA256,
		kmspb.CryptoKeyVersion_RSA_SIGN_PSS_4096_SHA256,
		kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_2048_SHA256,
		kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_3072_SHA256,
		kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_4096_SHA256,
		kmspb.CryptoKeyVersion_EC_SIGN_P256_SHA256:
		d, err := base64.StdEncoding.DecodeString(digest)
		if err != nil {
			return nil, errwrap.Wrapf("failed to decode base64 digest: {{err}}", err)
		}
		dig = &kmspb.Digest{
			Digest: &kmspb.Digest_Sha256{
				Sha256: d,
			},
		}
	case kmspb.CryptoKeyVersion_EC_SIGN_P384_SHA384:
		d, err := base64.StdEncoding.DecodeString(digest)
		if err != nil {
			return nil, errwrap.Wrapf("failed to decode base64 digest: {{err}}", err)
		}
		dig = &kmspb.Digest{
			Digest: &kmspb.Digest_Sha384{
				Sha384: d,
			},
		}
	default:
		return nil, fmt.Errorf("unknown key signing algorithm: %s", ckv.Algorithm)
	}

	resp, err := kmsClient.AsymmetricSign(ctx, &kmspb.AsymmetricSignRequest{
		Name:   ckv.Name,
		Digest: dig,
	})
	if err != nil {
		return nil, errwrap.Wrapf("failed to sign digest: {{err}}", err)
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"signature": base64.StdEncoding.EncodeToString(resp.Signature),
		},
	}, nil
}
