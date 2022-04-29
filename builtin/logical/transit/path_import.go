package transit

import (
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/google/tink/go/kwp/subtle"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/logical"
	"strconv"
)

const EncryptedKeyBytes = 512

func (b *backend) pathImport() *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("name") + "/import",
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The name of the key",
			},
			"type": {
				Type: framework.TypeString,
				Description: `
The type of key being imported. Currently, "aes128-gcm96" (symmetric), "aes256-gcm96" (symmetric), "ecdsa-p256"
(asymmetric), "ecdsa-p384" (asymmetric), "ecdsa-p521" (asymmetric), "ed25519" (asymmetric), "rsa-2048" (asymmetric), "rsa-3072"
(asymmetric), "rsa-4096" (asymmetric) are supported.  Defaults to "aes256-gcm96".
`,
			},
			"ciphertext": {
				Type: framework.TypeString,
				Description: `The base64-encoded ciphertext of the keys. The AES key should be encrypted using OAEP 
with the wrapping key and then concatenated with the import key, wrapped by the AES key.`,
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

			"convergent_encryption": {
				Type: framework.TypeBool,
				Description: `Whether to support convergent encryption.
This is only supported when using a key with
key derivation enabled and will require all
requests to carry both a context and 96-bit
(12-byte) nonce. The given nonce will be used
in place of a randomly generated nonce. As a
result, when the same context and nonce are
supplied, the same ciphertext is generated. It
is *very important* when using this mode that
you ensure that all nonces are unique for a
given context. Failing to do so will severely
impact the ciphertext's security.`,
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

func (b *backend) pathImportWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	derived := d.Get("derived").(bool)
	convergent := d.Get("convergent_encryption").(bool)
	keyType := d.Get("type").(string)
	exportable := d.Get("exportable").(bool)
	allowPlaintextBackup := d.Get("allow_plaintext_backup").(bool)
	autoRotatePeriod := time.Second * time.Duration(d.Get("auto_rotate_period").(int))
	ciphertextString := d.Get("ciphertext").(string)
	allowRotation := d.Get("allow_rotation").(bool)

	if autoRotatePeriod > 0 && !allowRotation {
		return nil, errors.New("allow_rotation must be set to true if auto-rotation is enabled")
	}

	polReq := keysutil.PolicyRequest{
		Storage:                  req.Storage,
		Name:                     name,
		Derived:                  derived,
		Convergent:               convergent,
		Exportable:               exportable,
		AllowPlaintextBackup:     allowPlaintextBackup,
		AutoRotatePeriod:         autoRotatePeriod,
		AllowImportedKeyRotation: allowRotation,
	}

	switch keyType {
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
	default:
		return logical.ErrorResponse(fmt.Sprintf("unknown key type: %v", keyType)), logical.ErrInvalidRequest
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

	ciphertext, err := base64.RawURLEncoding.DecodeString(ciphertextString)
	if err != nil {
		return nil, err
	}

	key, err := b.decryptImportedKey(ctx, req.Storage, ciphertext)
	if err != nil {
		return nil, err
	}

	err = b.lm.ImportPolicy(ctx, polReq, key, b.GetRandomReader())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) decryptImportedKey(ctx context.Context, storage logical.Storage, ciphertext []byte) ([]byte, error) {
	wrappedAESKey := ciphertext[:EncryptedKeyBytes]
	wrappedImportKey := ciphertext[EncryptedKeyBytes:]

	wrappingKey, err := getWrappingKey(ctx, storage)
	if err != nil {
		return nil, err
	}
	if wrappingKey == nil {
		return nil, fmt.Errorf("error importing key: wrapping key was nil")
	}

	rsaKey := wrappingKey.Keys[strconv.Itoa(wrappingKey.LatestVersion)].RSAKey
	aesKey, err := rsa.DecryptOAEP(sha256.New(), b.GetRandomReader(), rsaKey, wrappedAESKey, []byte{})
	if err != nil {
		return nil, err
	}

	kwp, err := subtle.NewKWP(aesKey)
	if err != nil {
		return nil, err
	}

	importKey, err := kwp.Unwrap(wrappedImportKey)
	if err != nil {
		return nil, err
	}

	return importKey, nil
}

const pathImportWriteSyn = "Imports an externally-generated key into transit"
const pathImportWriteDesc = ""
