package transit

import (
	"crypto/elliptic"
	"fmt"
	"strconv"

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
				Description: `The type of key to create. Currently,
"aes256-gcm96" (symmetric) and "ecdsa-p256" (asymmetric) are
supported. Defaults to "aes256-gcm96".`,
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

func (b *backend) pathKeysList(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List("policy/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

func (b *backend) pathPolicyWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	derived := d.Get("derived").(bool)
	convergent := d.Get("convergent_encryption").(bool)
	keyType := d.Get("type").(string)

	if !derived && convergent {
		return logical.ErrorResponse("convergent encryption requires derivation to be enabled"), nil
	}

	polReq := keysutil.PolicyRequest{
		Storage:    req.Storage,
		Name:       name,
		Derived:    derived,
		Convergent: convergent,
	}
	switch keyType {
	case "aes256-gcm96":
		polReq.KeyType = keysutil.KeyType_AES256_GCM96
	case "ecdsa-p256":
		polReq.KeyType = keysutil.KeyType_ECDSA_P256
	default:
		return logical.ErrorResponse(fmt.Sprintf("unknown key type %v", keyType)), logical.ErrInvalidRequest
	}

	p, lock, upserted, err := b.lm.GetPolicyUpsert(polReq)
	if lock != nil {
		defer lock.RUnlock()
	}
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, fmt.Errorf("error generating key: returned policy was nil")
	}

	resp := &logical.Response{}
	if !upserted {
		resp.AddWarning(fmt.Sprintf("key %s already existed", name))
	}

	return nil, nil
}

func (b *backend) pathPolicyRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	p, lock, err := b.lm.GetPolicyShared(req.Storage, name)
	if lock != nil {
		defer lock.RUnlock()
	}
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}

	// Return the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"name":                   p.Name,
			"type":                   p.Type.String(),
			"derived":                p.Derived,
			"deletion_allowed":       p.DeletionAllowed,
			"min_decryption_version": p.MinDecryptionVersion,
			"latest_version":         p.LatestVersion,
		},
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

	switch p.Type {
	case keysutil.KeyType_AES256_GCM96:
		retKeys := map[string]int64{}
		for k, v := range p.Keys {
			retKeys[strconv.Itoa(k)] = v.CreationTime
		}
		resp.Data["keys"] = retKeys

	case keysutil.KeyType_ECDSA_P256:
		type ecdsaKey struct {
			Name      string `json:"name"`
			PublicKey string `json:"public_key"`
		}
		retKeys := map[string]ecdsaKey{}
		for k, v := range p.Keys {
			retKeys[strconv.Itoa(k)] = ecdsaKey{
				Name:      elliptic.P256().Params().Name,
				PublicKey: v.FormattedPublicKey,
			}
		}
		resp.Data["keys"] = retKeys
	}

	return resp, nil
}

func (b *backend) pathPolicyDelete(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	// Delete does its own locking
	err := b.lm.DeletePolicy(req.Storage, name)
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
