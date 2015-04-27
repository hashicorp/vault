package transit

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathEncrypt() *framework.Path {
	return &framework.Path{
		Pattern: `encrypt/(?P<name>\w+)`,
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the policy",
			},

			"plaintext": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Plaintext value to encrypt",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: pathEncryptWrite,
		},

		HelpSynopsis:    pathEncryptHelpSyn,
		HelpDescription: pathEncryptHelpDesc,
	}
}

func pathEncryptWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	value := d.Get("plaintext").(string)
	if len(value) == 0 {
		return logical.ErrorResponse("missing plaintext to encrypt"), logical.ErrInvalidRequest
	}

	// Decode the plaintext value
	plaintext, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return logical.ErrorResponse("failed to decode plaintext as base64"), logical.ErrInvalidRequest
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

	// Guard against a potentially invalid cipher-mode
	switch p.CipherMode {
	case "aes-gcm":
	default:
		return logical.ErrorResponse("unsupported cipher mode"), logical.ErrInvalidRequest
	}

	// Setup the cipher
	aesCipher, err := aes.NewCipher(p.Key)
	if err != nil {
		return nil, err
	}

	// Setup the GCM AEAD
	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return nil, err
	}

	// Compute random nonce
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	// Encrypt and tag with GCM
	out := gcm.Seal(nil, nonce, plaintext, nil)

	// Place the encrypted data after the nonce
	full := append(nonce, out...)

	// Convert to base64
	encoded := base64.StdEncoding.EncodeToString(full)

	// Prepend some information
	encoded = "vault:v0:" + encoded

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"ciphertext": encoded,
		},
	}
	return resp, nil
}

const pathEncryptHelpSyn = `Encrypt a plaintext value using a named key`

const pathEncryptHelpDesc = `
This path uses the named key from the request path to encrypt a user
provided plaintext. The plaintext must be base64 encoded.
`
