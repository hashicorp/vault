package transit

import (
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

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
			"allow_plaintext_backup": {
				Type: framework.TypeBool,
				Description: `Enables taking a backup of the named
key in plaintext format. Once set,
this cannot be disabled.`,
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
	keyType := d.Get("type").(string)
	ciphertextString := d.Get("ciphertext").(string)
	allowRotation := d.Get("allow_rotation").(bool)

	polReq := keysutil.PolicyRequest{
		Storage:                  req.Storage,
		Name:                     name,
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
		p.Unlock()
		return nil, fmt.Errorf("the import path cannot overwrite an existing key; use import-version to rotate an imported key")
	}

	ciphertext, err := base64.RawURLEncoding.DecodeString(ciphertextString)
	if err != nil {
		return nil, err
	}

	wrappedAESKey := ciphertext[:EncryptedKeyBytes]
	wrappedImportKey := ciphertext[EncryptedKeyBytes:]

	wrappingKey, err := getWrappingKey(ctx, req.Storage)
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

	err = b.lm.ImportPolicy(ctx, polReq, importKey, b.GetRandomReader())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

const pathImportWriteSyn = "Imports an externally-generated key into transit"
const pathImportWriteDesc = ""
