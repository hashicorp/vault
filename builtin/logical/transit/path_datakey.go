package transit

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathDatakey() *framework.Path {
	return &framework.Path{
		Pattern: "datakey/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The backend key used for encrypting the data key",
			},

			"context": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Context for key derivation. Required for derived keys.",
			},

			"bits": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `Number of bits for the key; currently 128 and
256 are supported. Defaults to 256.`,
				Default: 256,
			},

			"no_plaintext": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Description: "If set, the plaintext of the key will not be returned",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: pathDatakeyWrite,
		},

		HelpSynopsis:    pathDatakeyHelpSyn,
		HelpDescription: pathDatakeyHelpDesc,
	}
}

func pathDatakeyWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	// Decode the context if any
	contextRaw := d.Get("context").(string)
	var context []byte
	if len(contextRaw) != 0 {
		var err error
		context, err = base64.StdEncoding.DecodeString(contextRaw)
		if err != nil {
			return logical.ErrorResponse("failed to decode context as base64"), logical.ErrInvalidRequest
		}
	}

	// Get the policy
	p, err := getPolicy(req, name)
	if err != nil {
		return nil, err
	}

	// Error if invalid policy
	if p == nil {
		return logical.ErrorResponse("policy not found"), logical.ErrInvalidRequest
	}

	newKey := make([]byte, 32)
	bits := d.Get("bits").(int)
	switch bits {
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

	ciphertext, err := p.Encrypt(context, base64.StdEncoding.EncodeToString(newKey))
	if err != nil {
		switch err.(type) {
		case certutil.UserError:
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		case certutil.InternalError:
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

	if !d.Get("no_plaintext").(bool) {
		resp.Data["plaintext"] = base64.StdEncoding.EncodeToString(newKey)
	}

	return resp, nil
}

const pathDatakeyHelpSyn = `Generate a data key`

const pathDatakeyHelpDesc = `
This path can be used to generate a data key: a random
key of a certain length that can be used for encryption
and decryption, protected by the named backend key. 128
or 256 bits can be specified; if not specified, the default
is 256 bits. The "no_plaintext" parameter can be used to
prevent the (base64-encoded) plaintext key from being
returned along with the encrypted key.
`
