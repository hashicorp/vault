package transit

import (
	"context"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/crypto/ed25519"

	"github.com/fatih/structs"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/keysutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathListKeys() *framework.Path {
	return &framework.Path{
		Pattern: "keys/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathKeysList,
		},

		HelpSynopsis:    pathPolicyHelpSyn,
		HelpDescription: pathPolicyHelpDesc,
	}
}

func (b *backend) pathKeys() *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the key",
			},

			"type": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: "aes256-gcm96",
				Description: `
The type of key to create. Currently, "aes256-gcm96" (symmetric), "ecdsa-p256"
(asymmetric), 'ed25519' (asymmetric), 'rsa-2048' (asymmetric), 'rsa-4096'
(asymmetric) are supported.  Defaults to "aes256-gcm96".
`,
			},

			"derived": &framework.FieldSchema{
				Type: framework.TypeBool,
				Description: `Enables key derivation mode. This
allows for per-transaction unique
keys for encryption operations.`,
			},

			"convergent_encryption": &framework.FieldSchema{
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

			"exportable": &framework.FieldSchema{
				Type: framework.TypeBool,
				Description: `Enables keys to be exportable.
This allows for all the valid keys
in the key ring to be exported.`,
			},

			"allow_plaintext_backup": &framework.FieldSchema{
				Type: framework.TypeBool,
				Description: `Enables taking a backup of the named
key in plaintext format. Once set,
this cannot be disabled.`,
			},

			"context": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `Base64 encoded context for key derivation.
When reading a key with key derivation enabled,
if the key type supports public keys, this will
return the public key for the given context.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathPolicyWrite,
			logical.DeleteOperation: b.pathPolicyDelete,
			logical.ReadOperation:   b.pathPolicyRead,
		},

		HelpSynopsis:    pathPolicyHelpSyn,
		HelpDescription: pathPolicyHelpDesc,
	}
}

func (b *backend) pathKeysList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "policy/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

func (b *backend) pathPolicyWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	derived := d.Get("derived").(bool)
	convergent := d.Get("convergent_encryption").(bool)
	keyType := d.Get("type").(string)
	exportable := d.Get("exportable").(bool)
	allowPlaintextBackup := d.Get("allow_plaintext_backup").(bool)

	if !derived && convergent {
		return logical.ErrorResponse("convergent encryption requires derivation to be enabled"), nil
	}

	polReq := keysutil.PolicyRequest{
		Upsert:               true,
		Storage:              req.Storage,
		Name:                 name,
		Derived:              derived,
		Convergent:           convergent,
		Exportable:           exportable,
		AllowPlaintextBackup: allowPlaintextBackup,
	}
	switch keyType {
	case "aes256-gcm96":
		polReq.KeyType = keysutil.KeyType_AES256_GCM96
	case "chacha20-poly1305":
		polReq.KeyType = keysutil.KeyType_ChaCha20_Poly1305
	case "ecdsa-p256":
		polReq.KeyType = keysutil.KeyType_ECDSA_P256
	case "ed25519":
		polReq.KeyType = keysutil.KeyType_ED25519
	case "rsa-2048":
		polReq.KeyType = keysutil.KeyType_RSA2048
	case "rsa-4096":
		polReq.KeyType = keysutil.KeyType_RSA4096
	default:
		return logical.ErrorResponse(fmt.Sprintf("unknown key type %v", keyType)), logical.ErrInvalidRequest
	}

	p, upserted, err := b.lm.GetPolicy(ctx, polReq)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, fmt.Errorf("error generating key: returned policy was nil")
	}
	if b.System().CachingDisabled() {
		p.Unlock()
	}

	resp := &logical.Response{}
	if !upserted {
		resp.AddWarning(fmt.Sprintf("key %s already existed", name))
	}

	return nil, nil
}

// Built-in helper type for returning asymmetric keys
type asymKey struct {
	Name         string    `json:"name" structs:"name" mapstructure:"name"`
	PublicKey    string    `json:"public_key" structs:"public_key" mapstructure:"public_key"`
	CreationTime time.Time `json:"creation_time" structs:"creation_time" mapstructure:"creation_time"`
}

func (b *backend) pathPolicyRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	p, _, err := b.lm.GetPolicy(ctx, keysutil.PolicyRequest{
		Storage: req.Storage,
		Name:    name,
	})
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}
	if !b.System().CachingDisabled() {
		p.Lock(false)
	}
	defer p.Unlock()

	// Return the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"name":                   p.Name,
			"type":                   p.Type.String(),
			"derived":                p.Derived,
			"deletion_allowed":       p.DeletionAllowed,
			"min_available_version":  p.MinAvailableVersion,
			"min_decryption_version": p.MinDecryptionVersion,
			"min_encryption_version": p.MinEncryptionVersion,
			"latest_version":         p.LatestVersion,
			"exportable":             p.Exportable,
			"allow_plaintext_backup": p.AllowPlaintextBackup,
			"supports_encryption":    p.Type.EncryptionSupported(),
			"supports_decryption":    p.Type.DecryptionSupported(),
			"supports_signing":       p.Type.SigningSupported(),
			"supports_derivation":    p.Type.DerivationSupported(),
		},
	}

	if p.BackupInfo != nil {
		resp.Data["backup_info"] = map[string]interface{}{
			"time":    p.BackupInfo.Time,
			"version": p.BackupInfo.Version,
		}
	}
	if p.RestoreInfo != nil {
		resp.Data["restore_info"] = map[string]interface{}{
			"time":    p.RestoreInfo.Time,
			"version": p.RestoreInfo.Version,
		}
	}

	if p.Derived {
		switch p.KDF {
		case keysutil.Kdf_hmac_sha256_counter:
			resp.Data["kdf"] = "hmac-sha256-counter"
			resp.Data["kdf_mode"] = "hmac-sha256-counter"
		case keysutil.Kdf_hkdf_sha256:
			resp.Data["kdf"] = "hkdf_sha256"
		}
		resp.Data["convergent_encryption"] = p.ConvergentEncryption
		if p.ConvergentEncryption {
			resp.Data["convergent_encryption_version"] = p.ConvergentVersion
		}
	}

	contextRaw := d.Get("context").(string)
	var context []byte
	if len(contextRaw) != 0 {
		context, err = base64.StdEncoding.DecodeString(contextRaw)
		if err != nil {
			return logical.ErrorResponse("failed to base64-decode context"), logical.ErrInvalidRequest
		}
	}

	switch p.Type {
	case keysutil.KeyType_AES256_GCM96, keysutil.KeyType_ChaCha20_Poly1305:
		retKeys := map[string]int64{}
		for k, v := range p.Keys {
			retKeys[k] = v.DeprecatedCreationTime
		}
		resp.Data["keys"] = retKeys

	case keysutil.KeyType_ECDSA_P256, keysutil.KeyType_ED25519, keysutil.KeyType_RSA2048, keysutil.KeyType_RSA4096:
		retKeys := map[string]map[string]interface{}{}
		for k, v := range p.Keys {
			key := asymKey{
				PublicKey:    v.FormattedPublicKey,
				CreationTime: v.CreationTime,
			}
			if key.CreationTime.IsZero() {
				key.CreationTime = time.Unix(v.DeprecatedCreationTime, 0)
			}

			switch p.Type {
			case keysutil.KeyType_ECDSA_P256:
				key.Name = elliptic.P256().Params().Name
			case keysutil.KeyType_ED25519:
				if p.Derived {
					if len(context) == 0 {
						key.PublicKey = ""
					} else {
						ver, err := strconv.Atoi(k)
						if err != nil {
							return nil, errwrap.Wrapf(fmt.Sprintf("invalid version %q: {{err}}", k), err)
						}
						derived, err := p.DeriveKey(context, ver, 32)
						if err != nil {
							return nil, fmt.Errorf("failed to derive key to return public component")
						}
						pubKey := ed25519.PrivateKey(derived).Public().(ed25519.PublicKey)
						key.PublicKey = base64.StdEncoding.EncodeToString(pubKey)
					}
				}
				key.Name = "ed25519"
			case keysutil.KeyType_RSA2048, keysutil.KeyType_RSA4096:
				key.Name = "rsa-2048"
				if p.Type == keysutil.KeyType_RSA4096 {
					key.Name = "rsa-4096"
				}

				// Encode the RSA public key in PEM format to return over the
				// API
				derBytes, err := x509.MarshalPKIXPublicKey(v.RSAKey.Public())
				if err != nil {
					return nil, errwrap.Wrapf("error marshaling RSA public key: {{err}}", err)
				}
				pemBlock := &pem.Block{
					Type:  "PUBLIC KEY",
					Bytes: derBytes,
				}
				pemBytes := pem.EncodeToMemory(pemBlock)
				if pemBytes == nil || len(pemBytes) == 0 {
					return nil, fmt.Errorf("failed to PEM-encode RSA public key")
				}
				key.PublicKey = string(pemBytes)
			}

			retKeys[k] = structs.New(key).Map()
		}
		resp.Data["keys"] = retKeys
	}

	return resp, nil
}

func (b *backend) pathPolicyDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	// Delete does its own locking
	err := b.lm.DeletePolicy(ctx, req.Storage, name)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("error deleting policy %s: %s", name, err)), err
	}

	return nil, nil
}

const pathPolicyHelpSyn = `Managed named encryption keys`

const pathPolicyHelpDesc = `
This path is used to manage the named keys that are available.
Doing a write with no value against a new named key will create
it using a randomly generated key.
`
