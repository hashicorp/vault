// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpkms

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

func (b *backend) pathVerify() *framework.Path {
	return &framework.Path{
		Pattern: "verify/" + framework.GenericNameRegex("key"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloudKMS,
			OperationVerb:   "verify",
		},

		HelpSynopsis: "Verify a signature using a named key",
		HelpDescription: `
Use the named key to verify the given signature. The response will be the
base64-encoded encrypted value (ciphertext).
`,

		Fields: map[string]*framework.FieldSchema{
			"key": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Name of the key in Vault to use for verification. This key must already exist in
Vault and must map back to a Google Cloud KMS key.
`,
			},

			"digest": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Digest to verify. This digest must use the same SHA algorithm as the underlying
Cloud KMS key. The digest must be the base64-encoded binary value. This field is
required.
`,
			},

			"key_version": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `
Integer version of the crypto key version to use for verification. This field is
required.
`,
			},

			"signature": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Base64-encoded signature to use for verification. This field is required.
`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: withFieldValidator(b.pathVerifyWrite),
		},
	}
}

// pathVerifyWrite corresponds to PUT/POST gcpkms/sign/:key and is used to
// verify the digest using the named key.
func (b *backend) pathVerifyWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	key := d.Get("key").(string)
	digest := d.Get("digest").(string)
	signature := d.Get("signature").(string)
	keyVersion := d.Get("key_version").(int)

	if digest == "" {
		return nil, errMissingFields("digest")
	}

	if signature == "" {
		return nil, errMissingFields("signature")
	}

	if keyVersion == 0 {
		return nil, errMissingFields("key_version")
	}

	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return nil, errwrap.Wrapf("failed to base64 decode signature: {{err}}", err)
	}

	dig, err := base64.StdEncoding.DecodeString(digest)
	if err != nil {
		return nil, errwrap.Wrapf("failed to base64 decode digest: {{err}}", err)
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

	// Get the public key
	pk, err := kmsClient.GetPublicKey(ctx, &kmspb.GetPublicKeyRequest{
		Name: fmt.Sprintf("%s/cryptoKeyVersions/%d", k.CryptoKeyID, keyVersion),
	})
	if err != nil {
		return nil, errwrap.Wrapf("failed to get public key: {{err}}", err)
	}

	// Extract the PEM-encoded data block
	block, _ := pem.Decode([]byte(pk.Pem))
	if block == nil {
		return nil, fmt.Errorf("public key is not in pem format: %s", pk.Pem)
	}

	// Decode the public key
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errwrap.Wrapf("failed to parse public key: {{err}}", err)
	}

	validSig := false

	switch pk.Algorithm {
	case kmspb.CryptoKeyVersion_EC_SIGN_P256_SHA256:
		var parsedSig struct{ R, S *big.Int }
		if _, err := asn1.Unmarshal(sig, &parsedSig); err != nil {
			return nil, errwrap.Wrapf("failed to unmarshal signature: {{err}}", err)
		}
		validSig = ecdsa.Verify(pub.(*ecdsa.PublicKey), dig, parsedSig.R, parsedSig.S)
	case kmspb.CryptoKeyVersion_EC_SIGN_P384_SHA384:
		var parsedSig struct{ R, S *big.Int }
		if _, err := asn1.Unmarshal(sig, &parsedSig); err != nil {
			return nil, errwrap.Wrapf("failed to unmarshal signature: {{err}}", err)
		}
		validSig = ecdsa.Verify(pub.(*ecdsa.PublicKey), dig, parsedSig.R, parsedSig.S)
	case kmspb.CryptoKeyVersion_RSA_SIGN_PSS_2048_SHA256,
		kmspb.CryptoKeyVersion_RSA_SIGN_PSS_3072_SHA256,
		kmspb.CryptoKeyVersion_RSA_SIGN_PSS_4096_SHA256:
		err := rsa.VerifyPSS(pub.(*rsa.PublicKey), crypto.SHA256, dig, sig, &rsa.PSSOptions{})
		validSig = err == nil
	case kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_2048_SHA256,
		kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_3072_SHA256,
		kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_4096_SHA256:
		err := rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA256, dig, sig)
		validSig = err == nil
	default:
		return nil, fmt.Errorf("unknown key signing algorithm: %s", pk.Algorithm)
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"valid": validSig,
		},
	}, nil
}
