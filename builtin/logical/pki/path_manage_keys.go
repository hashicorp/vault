// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"bytes"
	"context"
	"encoding/pem"
	"net/http"
	"strings"

	"github.com/hashicorp/vault/builtin/logical/pki/managed_key"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathGenerateKey(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "keys/generate/(internal|exported|kms)",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
			OperationVerb:   "generate",
			OperationSuffix: "internal-key|exported-key|kms-key",
		},

		Fields: map[string]*framework.FieldSchema{
			keyNameParam: {
				Type:        framework.TypeString,
				Description: "Optional name to be used for this key",
			},
			keyTypeParam: {
				Type:    framework.TypeString,
				Default: "rsa",
				Description: `The type of key to use; defaults to RSA. "rsa"
"ec" and "ed25519" are the only valid values.`,
				AllowedValues: []interface{}{"rsa", "ec", "ed25519"},
				DisplayAttrs: &framework.DisplayAttributes{
					Value: "rsa",
				},
			},
			keyBitsParam: {
				Type:    framework.TypeInt,
				Default: 0,
				Description: `The number of bits to use. Allowed values are
0 (universal default); with rsa key_type: 2048 (default), 3072, 4096 or 8192;
with ec key_type: 224, 256 (default), 384, or 521; ignored with ed25519.`,
			},
			"managed_key_name": {
				Type: framework.TypeString,
				Description: `The name of the managed key to use when the exported
type is kms. When kms type is the key type, this field or managed_key_id
is required. Ignored for other types.`,
			},
			"managed_key_id": {
				Type: framework.TypeString,
				Description: `The name of the managed key to use when the exported
type is kms. When kms type is the key type, this field or managed_key_name
is required. Ignored for other types.`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathGenerateKeyHandler,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"key_id": {
								Type:        framework.TypeString,
								Description: `ID assigned to this key.`,
								Required:    true,
							},
							"key_name": {
								Type:        framework.TypeString,
								Description: `Name assigned to this key.`,
								Required:    true,
							},
							"key_type": {
								Type: framework.TypeString,
								Description: `The type of key to use; defaults to RSA. "rsa"
								"ec" and "ed25519" are the only valid values.`,
								Required: true,
							},
							"private_key": {
								Type:        framework.TypeString,
								Description: `The private key string`,
								Required:    false,
							},
						},
					}},
				},

				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
		},

		HelpSynopsis:    pathGenerateKeyHelpSyn,
		HelpDescription: pathGenerateKeyHelpDesc,
	}
}

const (
	pathGenerateKeyHelpSyn  = `Generate a new private key used for signing.`
	pathGenerateKeyHelpDesc = `This endpoint will generate a new key pair of the specified type (internal, exported, or kms).`
)

func (b *backend) pathGenerateKeyHandler(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Since we're planning on updating issuers here, grab the lock so we've
	// got a consistent view.
	b.issuersLock.Lock()
	defer b.issuersLock.Unlock()

	if b.UseLegacyBundleCaStorage() {
		return logical.ErrorResponse("Can not generate keys until migration has completed"), nil
	}

	sc := b.makeStorageContext(ctx, req.Storage)
	keyName, err := getKeyName(sc, data)
	if err != nil { // Fail Immediately if Key Name is in Use, etc...
		return logical.ErrorResponse(err.Error()), nil
	}

	exportPrivateKey := false
	var keyBundle certutil.KeyBundle
	var actualPrivateKeyType certutil.PrivateKeyType
	switch {
	case strings.HasSuffix(req.Path, "/exported"):
		exportPrivateKey = true
		fallthrough
	case strings.HasSuffix(req.Path, "/internal"):
		keyType := data.Get(keyTypeParam).(string)
		keyBits := data.Get(keyBitsParam).(int)

		keyBits, _, err := certutil.ValidateDefaultOrValueKeyTypeSignatureLength(keyType, keyBits, 0)
		if err != nil {
			return logical.ErrorResponse("Validation for key_type, key_bits failed: %s", err.Error()), nil
		}

		// Internal key generation, stored in storage
		keyBundle, err = certutil.CreateKeyBundle(keyType, keyBits, b.GetRandomReader())
		if err != nil {
			return nil, err
		}

		actualPrivateKeyType = keyBundle.PrivateKeyType
	case strings.HasSuffix(req.Path, "/kms"):
		keyId, err := getManagedKeyId(data)
		if err != nil {
			return nil, err
		}

		keyBundle, actualPrivateKeyType, err = managed_key.CreateKmsKeyBundle(ctx, b, keyId)
		if err != nil {
			return nil, err
		}
	default:
		return logical.ErrorResponse("Unknown type of key to generate"), nil
	}

	privateKeyPemString, err := keyBundle.ToPrivateKeyPemString()
	if err != nil {
		return nil, err
	}

	key, _, err := sc.importKey(privateKeyPemString, keyName, keyBundle.PrivateKeyType)
	if err != nil {
		return nil, err
	}
	responseData := map[string]interface{}{
		keyIdParam:   key.ID,
		keyNameParam: key.Name,
		keyTypeParam: string(actualPrivateKeyType),
	}
	if exportPrivateKey {
		responseData["private_key"] = privateKeyPemString
	}
	return &logical.Response{
		Data: responseData,
	}, nil
}

func pathImportKey(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "keys/import",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
			OperationVerb:   "import",
			OperationSuffix: "key",
		},

		Fields: map[string]*framework.FieldSchema{
			keyNameParam: {
				Type:        framework.TypeString,
				Description: "Optional name to be used for this key",
			},
			"pem_bundle": {
				Type:        framework.TypeString,
				Description: `PEM-format, unencrypted secret key`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathImportKeyHandler,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"key_id": {
								Type:        framework.TypeString,
								Description: `ID assigned to this key.`,
								Required:    true,
							},
							"key_name": {
								Type:        framework.TypeString,
								Description: `Name assigned to this key.`,
								Required:    true,
							},
							"key_type": {
								Type: framework.TypeString,
								Description: `The type of key to use; defaults to RSA. "rsa"
								"ec" and "ed25519" are the only valid values.`,
								Required: true,
							},
						},
					}},
				},
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
		},

		HelpSynopsis:    pathImportKeyHelpSyn,
		HelpDescription: pathImportKeyHelpDesc,
	}
}

const (
	pathImportKeyHelpSyn  = `Import the specified key.`
	pathImportKeyHelpDesc = `This endpoint allows importing a specified issuer key from a pem bundle.
If key_name is set, that will be set on the key, assuming the key did not exist previously.`
)

func (b *backend) pathImportKeyHandler(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Since we're planning on updating issuers here, grab the lock so we've
	// got a consistent view.
	b.issuersLock.Lock()
	defer b.issuersLock.Unlock()

	if b.UseLegacyBundleCaStorage() {
		return logical.ErrorResponse("Cannot import keys until migration has completed"), nil
	}

	sc := b.makeStorageContext(ctx, req.Storage)
	pemBundle := data.Get("pem_bundle").(string)
	keyName, err := getKeyName(sc, data)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	if len(pemBundle) < 64 {
		// It is almost nearly impossible to store a complete key in
		// less than 64 bytes. It is definitely impossible to do so when PEM
		// encoding has been applied. Detect this and give a better warning
		// than "provided PEM block contained no data" in this case. This is
		// because the PEM headers contain 5*4 + 6 + 4 + 2 + 2 = 34 characters
		// minimum (five dashes, "BEGIN" + space + at least one character
		// identifier, "END" + space + at least one character identifier, and
		// a pair of new lines). That would leave 30 bytes for Base64 data,
		// meaning at most a 22-byte DER key. Even with a 128-bit key, 6 bytes
		// is not sufficient for the required ASN.1 structure and OID encoding.
		//
		// However, < 64 bytes is probably a good length for a file path so
		// suggest that is the case.
		return logical.ErrorResponse("provided data for import was too short; perhaps a path was passed to the API rather than the contents of a PEM file"), nil
	}

	pemBytes := []byte(pemBundle)
	var pemBlock *pem.Block

	var keys []string
	for len(bytes.TrimSpace(pemBytes)) > 0 {
		pemBlock, pemBytes = pem.Decode(pemBytes)
		if pemBlock == nil {
			return logical.ErrorResponse("provided PEM block contained no data"), nil
		}

		pemBlockString := string(pem.EncodeToMemory(pemBlock))
		keys = append(keys, pemBlockString)
	}

	if len(keys) != 1 {
		return logical.ErrorResponse("only a single key can be present within the pem_bundle for importing"), nil
	}

	key, existed, err := importKeyFromBytes(sc, keys[0], keyName)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	resp := logical.Response{
		Data: map[string]interface{}{
			keyIdParam:   key.ID,
			keyNameParam: key.Name,
			keyTypeParam: key.PrivateKeyType,
		},
	}

	if existed {
		resp.AddWarning("Key already imported, use key/ endpoint to update name.")
	}

	return &resp, nil
}
