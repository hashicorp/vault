package transit

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathDatakey() *framework.Path {
	return &framework.Path{
		Pattern: "datakey/" + framework.GenericNameRegex("plaintext") + "/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The backend key used for encrypting the data key",
			},

			"plaintext": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `"plaintext" will return the key in both plaintext and
ciphertext; "wrapped" will return the ciphertext only.`,
			},

			"context": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Context for key derivation. Required for derived keys.",
			},

			"nonce": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Nonce for when convergent encryption v1 is used (only in Vault 0.6.1)",
			},

			"bits": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `Number of bits for the key; currently 128, 256,
and 512 bits are supported. Defaults to 256.`,
				Default: 256,
			},

			"key_version": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `The version of the Vault key to use for
encryption of the data key. Must be 0 (for latest)
or a value greater than or equal to the
min_encryption_version configured on the key.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathDatakeyWrite,
		},

		HelpSynopsis:    pathDatakeyHelpSyn,
		HelpDescription: pathDatakeyHelpDesc,
	}
}

func (b *backend) pathDatakeyWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	ver := d.Get("key_version").(int)

	plaintext := d.Get("plaintext").(string)
	plaintextAllowed := false
	switch plaintext {
	case "plaintext":
		plaintextAllowed = true
	case "wrapped":
	default:
		return logical.ErrorResponse("Invalid path, must be 'plaintext' or 'wrapped'"), logical.ErrInvalidRequest
	}

	var err error

	// Decode the context if any
	contextRaw := d.Get("context").(string)
	var context []byte
	if len(contextRaw) != 0 {
		context, err = base64.StdEncoding.DecodeString(contextRaw)
		if err != nil {
			return logical.ErrorResponse("failed to base64-decode context"), logical.ErrInvalidRequest
		}
	}

	// Decode the nonce if any
	nonceRaw := d.Get("nonce").(string)
	var nonce []byte
	if len(nonceRaw) != 0 {
		nonce, err = base64.StdEncoding.DecodeString(nonceRaw)
		if err != nil {
			return logical.ErrorResponse("failed to base64-decode nonce"), logical.ErrInvalidRequest
		}
	}

	// Get the policy
	p, lock, err := b.lm.GetPolicyShared(req.Storage, name)
	if lock != nil {
		defer lock.RUnlock()
	}
	if err != nil {
		return nil, err
	}
	if p == nil {
		return logical.ErrorResponse("encryption key not found"), logical.ErrInvalidRequest
	}

	newKey := make([]byte, 32)
	bits := d.Get("bits").(int)
	switch bits {
	case 512:
		newKey = make([]byte, 64)
	case 256:
	case 128:
		newKey = make([]byte, 16)
	default:
		return logical.ErrorResponse("invalid bit length"), logical.ErrInvalidRequest
	}
	_, err = rand.Read(newKey)
	if err != nil {
		return nil, err
	}

	ciphertext, err := p.Encrypt(ver, context, nonce, base64.StdEncoding.EncodeToString(newKey))
	if err != nil {
		switch err.(type) {
		case errutil.UserError:
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		case errutil.InternalError:
			return nil, err
		default:
			return nil, err
		}
	}

	if ciphertext == "" {
		return nil, fmt.Errorf("empty ciphertext returned")
	}

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"ciphertext": ciphertext,
		},
	}

	if plaintextAllowed {
		resp.Data["plaintext"] = base64.StdEncoding.EncodeToString(newKey)
	}

	return resp, nil
}

const pathDatakeyHelpSyn = `Generate a data key`

const pathDatakeyHelpDesc = `
This path can be used to generate a data key: a random
key of a certain length that can be used for encryption
and decryption, protected by the named backend key. 128, 256,
or 512 bits can be specified; if not specified, the default
is 256 bits. Call with the the "wrapped" path to prevent the
(base64-encoded) plaintext key from being returned along with
the encrypted key, the "plaintext" path returns both.
`
