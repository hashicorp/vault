package transit

import (
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathEncrypt() *framework.Path {
	return &framework.Path{
		Pattern: "encrypt/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the policy",
			},

			"plaintext": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Plaintext value to encrypt",
			},

			"context": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Context for key derivation. Required for derived keys.",
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

	// Get the policy
	p, err := getPolicy(req, name)
	if err != nil {
		return nil, err
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

	// Error if invalid policy
	if p == nil {
		isDerived := len(context) != 0
		p, err = generatePolicy(req.Storage, name, isDerived)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("failed to upsert policy: %v", err)), logical.ErrInvalidRequest
		}
	}

	ciphertext, err := p.Encrypt(context, value)
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
	return resp, nil
}

const pathEncryptHelpSyn = `Encrypt a plaintext value using a named key`

const pathEncryptHelpDesc = `
This path uses the named key from the request path to encrypt a user
provided plaintext. The plaintext must be base64 encoded.
`
