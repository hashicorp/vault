package transit

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathDecrypt() *framework.Path {
	return &framework.Path{
		Pattern: `decrypt/(?P<name>\w+)`,
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the policy",
			},

			"ciphertext": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Ciphertext value to decrypt",
			},

			"context": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Context for key derivation. Required for derived keys.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: pathDecryptWrite,
		},

		HelpSynopsis:    pathDecryptHelpSyn,
		HelpDescription: pathDecryptHelpDesc,
	}
}

func pathDecryptWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	value := d.Get("ciphertext").(string)
	if len(value) == 0 {
		return logical.ErrorResponse("missing ciphertext to decrypt"), logical.ErrInvalidRequest
	}

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

	// Derive the key that should be used
	key, err := p.DeriveKey(context)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	// Guard against a potentially invalid cipher-mode
	switch p.CipherMode {
	case "aes-gcm":
	default:
		return logical.ErrorResponse("unsupported cipher mode"), logical.ErrInvalidRequest
	}

	// Verify the prefix
	if !strings.HasPrefix(value, "vault:v0:") {
		return logical.ErrorResponse("invalid ciphertext"), logical.ErrInvalidRequest
	}

	// Decode the base64
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(value, "vault:v0:"))
	if err != nil {
		return logical.ErrorResponse("invalid ciphertext"), logical.ErrInvalidRequest
	}

	// Setup the cipher
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Setup the GCM AEAD
	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return nil, err
	}

	// Extract the nonce and ciphertext
	nonce := decoded[:gcm.NonceSize()]
	ciphertext := decoded[gcm.NonceSize():]

	// Verify and Decrypt
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return logical.ErrorResponse("invalid ciphertext"), logical.ErrInvalidRequest
	}

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"plaintext": base64.StdEncoding.EncodeToString(plain),
		},
	}
	return resp, nil
}

const pathDecryptHelpSyn = `Decrypt a ciphertext value using a named key`

const pathDecryptHelpDesc = `
This path uses the named key from the request path to decrypt a user
provided ciphertext. The plaintext is returned base64 encoded.
`
