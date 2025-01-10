// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/tink-crypto/tink-go/v2/kwp/subtle"
)

const EncryptedKeyBytes = 512

func (b *backend) pathImport() *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("name") + "/import",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixTransit,
			OperationVerb:   "import",
			OperationSuffix: "key",
		},

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "The name of the key",
			},
			"type": {
				Type:    framework.TypeString,
				Default: "aes256-gcm96",
				Description: `The type of key being imported. Currently, "aes128-gcm96" (symmetric), "aes256-gcm96" (symmetric), "ecdsa-p256"
(asymmetric), "ecdsa-p384" (asymmetric), "ecdsa-p521" (asymmetric), "ed25519" (asymmetric), "rsa-2048" (asymmetric), "rsa-3072"
(asymmetric), "rsa-4096" (asymmetric), "ml-dsa-44 (asymmetric)", "ml-dsa-65 (asymmetric)", "ml-dsa-87 (asymmetric)", "hmac", "aes128-cmac", "aes256-cmac" are supported.  Defaults to "aes256-gcm96".
`,
			},
			"hash_function": {
				Type:    framework.TypeString,
				Default: "SHA256",
				Description: `The hash function used as a random oracle in the OAEP wrapping of the user-generated,
ephemeral AES key. Can be one of "SHA1", "SHA224", "SHA256" (default), "SHA384", or "SHA512"`,
			},
			"ciphertext": {
				Type: framework.TypeString,
				Description: `The base64-encoded ciphertext of the keys. The AES key should be encrypted using OAEP 
with the wrapping key and then concatenated with the import key, wrapped by the AES key.`,
			},
			"public_key": {
				Type:        framework.TypeString,
				Description: `The plaintext PEM public key to be imported. If "ciphertext" is set, this field is ignored.`,
			},
			"allow_rotation": {
				Type:        framework.TypeBool,
				Description: "True if the imported key may be rotated within Vault; false otherwise.",
			},
			"derived": {
				Type: framework.TypeBool,
				Description: `Enables key derivation mode. This
allows for per-transaction unique
keys for encryption operations.`,
			},

			"exportable": {
				Type: framework.TypeBool,
				Description: `Enables keys to be exportable.
This allows for all the valid keys
in the key ring to be exported.`,
			},

			"allow_plaintext_backup": {
				Type: framework.TypeBool,
				Description: `Enables taking a backup of the named
key in plaintext format. Once set,
this cannot be disabled.`,
			},

			"context": {
				Type: framework.TypeString,
				Description: `Base64 encoded context for key derivation.
When reading a key with key derivation enabled,
if the key type supports public keys, this will
return the public key for the given context.`,
			},
			"auto_rotate_period": {
				Type:    framework.TypeDurationSecond,
				Default: 0,
				Description: `Amount of time the key should live before
being automatically rotated. A value of 0
(default) disables automatic rotation for the
key.`,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathImportWrite,
		},
		HelpSynopsis:    pathImportWriteSyn,
		HelpDescription: pathImportWriteDesc,
	}
}

func (b *backend) pathImportVersion() *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("name") + "/import_version",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixTransit,
			OperationVerb:   "import",
			OperationSuffix: "key-version",
		},

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "The name of the key",
			},
			"ciphertext": {
				Type: framework.TypeString,
				Description: `The base64-encoded ciphertext of the keys. The AES key should be encrypted using OAEP 
with the wrapping key and then concatenated with the import key, wrapped by the AES key.`,
			},
			"public_key": {
				Type:        framework.TypeString,
				Description: `The plaintext public key to be imported. If "ciphertext" is set, this field is ignored.`,
			},
			"hash_function": {
				Type:    framework.TypeString,
				Default: "SHA256",
				Description: `The hash function used as a random oracle in the OAEP wrapping of the user-generated,
ephemeral AES key. Can be one of "SHA1", "SHA224", "SHA256" (default), "SHA384", or "SHA512"`,
			},
			"version": {
				Type: framework.TypeInt,
				Description: `Key version to be updated, if left empty, a new version will be created unless
a private key is specified and the 'Latest' key is missing a private key.`,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathImportVersionWrite,
		},
		HelpSynopsis:    pathImportVersionWriteSyn,
		HelpDescription: pathImportVersionWriteDesc,
	}
}

func (b *backend) pathImportWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	derived := d.Get("derived").(bool)
	keyType := d.Get("type").(string)
	exportable := d.Get("exportable").(bool)
	allowPlaintextBackup := d.Get("allow_plaintext_backup").(bool)
	autoRotatePeriod := time.Second * time.Duration(d.Get("auto_rotate_period").(int))
	allowRotation := d.Get("allow_rotation").(bool)

	// Ensure the caller didn't supply "convergent_encryption" as a field, since it's not supported on import.
	if _, ok := d.Raw["convergent_encryption"]; ok {
		return nil, errors.New("import cannot be used on keys with convergent encryption enabled")
	}

	if autoRotatePeriod > 0 && !allowRotation {
		return nil, errors.New("allow_rotation must be set to true if auto-rotation is enabled")
	}

	// Ensure that at least on `key` field has been set
	isCiphertextSet, err := checkKeyFieldsSet(d)
	if err != nil {
		return nil, err
	}

	polReq := keysutil.PolicyRequest{
		Storage:                  req.Storage,
		Name:                     name,
		Derived:                  derived,
		Exportable:               exportable,
		AllowPlaintextBackup:     allowPlaintextBackup,
		AutoRotatePeriod:         autoRotatePeriod,
		AllowImportedKeyRotation: allowRotation,
		IsPrivateKey:             isCiphertextSet,
	}

	switch strings.ToLower(keyType) {
	case "aes128-gcm96":
		polReq.KeyType = keysutil.KeyType_AES128_GCM96
	case "aes256-gcm96":
		polReq.KeyType = keysutil.KeyType_AES256_GCM96
	case "chacha20-poly1305":
		polReq.KeyType = keysutil.KeyType_ChaCha20_Poly1305
	case "ecdsa-p256":
		polReq.KeyType = keysutil.KeyType_ECDSA_P256
	case "ecdsa-p384":
		polReq.KeyType = keysutil.KeyType_ECDSA_P384
	case "ecdsa-p521":
		polReq.KeyType = keysutil.KeyType_ECDSA_P521
	case "ed25519":
		polReq.KeyType = keysutil.KeyType_ED25519
	case "rsa-2048":
		polReq.KeyType = keysutil.KeyType_RSA2048
	case "rsa-3072":
		polReq.KeyType = keysutil.KeyType_RSA3072
	case "rsa-4096":
		polReq.KeyType = keysutil.KeyType_RSA4096
	case "hmac":
		polReq.KeyType = keysutil.KeyType_HMAC
	case "aes128-cmac":
		polReq.KeyType = keysutil.KeyType_AES128_CMAC
	case "aes256-cmac":
		polReq.KeyType = keysutil.KeyType_AES256_CMAC
	default:
		return logical.ErrorResponse(fmt.Sprintf("unknown key type: %v", keyType)), logical.ErrInvalidRequest
	}

	if polReq.KeyType.CMACSupported() && !constants.IsEnterprise {
		return logical.ErrorResponse(ErrCmacEntOnly.Error()), logical.ErrInvalidRequest
	}

	p, _, err := b.GetPolicy(ctx, polReq, b.GetRandomReader())
	if err != nil {
		return nil, err
	}

	if p != nil {
		if b.System().CachingDisabled() {
			p.Unlock()
		}
		return nil, errors.New("the import path cannot be used with an existing key; use import-version to rotate an existing imported key")
	}

	key, resp, err := b.extractKeyFromFields(ctx, req, d, polReq.KeyType, isCiphertextSet)
	if err != nil {
		return resp, err
	}

	err = b.lm.ImportPolicy(ctx, polReq, key, b.GetRandomReader())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathImportVersionWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	isCiphertextSet, err := checkKeyFieldsSet(d)
	if err != nil {
		return nil, err
	}

	polReq := keysutil.PolicyRequest{
		Storage:      req.Storage,
		Name:         name,
		Upsert:       false,
		IsPrivateKey: isCiphertextSet,
	}
	p, _, err := b.GetPolicy(ctx, polReq, b.GetRandomReader())
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, fmt.Errorf("no key found with name %s; to import a new key, use the import/ endpoint", name)
	}
	if !p.Imported {
		return nil, errors.New("the import_version endpoint can only be used with an imported key")
	}
	if p.ConvergentEncryption {
		return nil, errors.New("import_version cannot be used on keys with convergent encryption enabled")
	}

	if !b.System().CachingDisabled() {
		p.Lock(true)
	}
	defer p.Unlock()

	key, resp, err := b.extractKeyFromFields(ctx, req, d, p.Type, isCiphertextSet)
	if err != nil {
		return resp, err
	}

	// Get param version if set else import a new version.
	if version, ok := d.GetOk("version"); ok {
		versionToUpdate := version.(int)

		// Check if given version can be updated given input
		err = p.KeyVersionCanBeUpdated(versionToUpdate, isCiphertextSet)
		if err == nil {
			err = p.ImportPrivateKeyForVersion(ctx, req.Storage, versionToUpdate, key)
		}
	} else {
		err = p.ImportPublicOrPrivate(ctx, req.Storage, key, isCiphertextSet, b.GetRandomReader())
	}

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) decryptImportedKey(ctx context.Context, storage logical.Storage, ciphertext []byte, hashFn hash.Hash) ([]byte, error) {
	// Bounds check the ciphertext to avoid panics
	if len(ciphertext) <= EncryptedKeyBytes {
		return nil, errors.New("provided ciphertext is too short")
	}

	wrappedEphKey := ciphertext[:EncryptedKeyBytes]
	wrappedImportKey := ciphertext[EncryptedKeyBytes:]

	wrappingKey, err := b.getWrappingKey(ctx, storage)
	if err != nil {
		return nil, err
	}
	if wrappingKey == nil {
		return nil, fmt.Errorf("error importing key: wrapping key was nil")
	}

	privWrappingKey := wrappingKey.Keys[strconv.Itoa(wrappingKey.LatestVersion)].RSAKey
	ephKey, err := rsa.DecryptOAEP(hashFn, b.GetRandomReader(), privWrappingKey, wrappedEphKey, []byte{})
	if err != nil {
		return nil, err
	}

	// Zero out the ephemeral AES key just to be extra cautious. Note that this
	// isn't a guarantee against memory analysis! See the documentation for the
	// `vault.memzero` utility function for more information.
	defer func() {
		for i := range ephKey {
			ephKey[i] = 0
		}
	}()

	// Ensure the ephemeral AES key is 256-bit
	if len(ephKey) != 32 {
		return nil, errors.New("expected ephemeral AES key to be 256-bit")
	}

	kwp, err := subtle.NewKWP(ephKey)
	if err != nil {
		return nil, err
	}

	importKey, err := kwp.Unwrap(wrappedImportKey)
	if err != nil {
		return nil, err
	}

	return importKey, nil
}

func (b *backend) extractKeyFromFields(ctx context.Context, req *logical.Request, d *framework.FieldData, keyType keysutil.KeyType, isPrivateKey bool) ([]byte, *logical.Response, error) {
	var key []byte
	if isPrivateKey {
		hashFnStr := d.Get("hash_function").(string)
		hashFn, err := parseHashFn(hashFnStr)
		if err != nil {
			return key, logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}

		ciphertextString := d.Get("ciphertext").(string)
		ciphertext, err := base64.StdEncoding.DecodeString(ciphertextString)
		if err != nil {
			return key, nil, err
		}

		key, err = b.decryptImportedKey(ctx, req.Storage, ciphertext, hashFn)
		if err != nil {
			return key, nil, err
		}
	} else {
		publicKeyString := d.Get("public_key").(string)
		if !keyType.ImportPublicKeySupported() {
			return key, nil, errors.New("provided type does not support public_key import")
		}
		key = []byte(publicKeyString)
	}

	return key, nil, nil
}

func parseHashFn(hashFn string) (hash.Hash, error) {
	switch strings.ToUpper(hashFn) {
	case "SHA1":
		return sha1.New(), nil
	case "SHA224":
		return sha256.New224(), nil
	case "SHA256":
		return sha256.New(), nil
	case "SHA384":
		return sha512.New384(), nil
	case "SHA512":
		return sha512.New(), nil
	default:
		return nil, fmt.Errorf("unknown hash function: %s", hashFn)
	}
}

// checkKeyFieldsSet: Checks which key fields are set. If both are set, an error is returned
func checkKeyFieldsSet(d *framework.FieldData) (bool, error) {
	ciphertextSet := isFieldSet("ciphertext", d)
	publicKeySet := isFieldSet("publicKey", d)

	if ciphertextSet && publicKeySet {
		return false, errors.New("only one of the following fields, ciphertext and public_key, can be set")
	} else if ciphertextSet {
		return true, nil
	} else {
		return false, nil
	}
}

func isFieldSet(fieldName string, d *framework.FieldData) bool {
	_, fieldSet := d.Raw[fieldName]
	if !fieldSet {
		return false
	}

	return true
}

const (
	pathImportWriteSyn  = "Imports an externally-generated key into a new transit key"
	pathImportWriteDesc = "This path is used to import an externally-generated " +
		"key into Vault. The import operation creates a new key and cannot be used to " +
		"replace an existing key."
)

const pathImportVersionWriteSyn = "Imports an externally-generated key into an " +
	"existing imported key"

const pathImportVersionWriteDesc = "This path is used to import a new version of an " +
	"externally-generated key into an existing import key. The import_version endpoint " +
	"only supports importing key material into existing imported keys."
