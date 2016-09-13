package transit

import (
	"encoding/base64"

	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathDecrypt() *framework.Path {
	return &framework.Path{
		Pattern: "decrypt/" + framework.GenericNameRegex("name"),
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

			"nonce": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Nonce for when convergent encryption is used",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathDecryptWrite,
		},

		HelpSynopsis:    pathDecryptHelpSyn,
		HelpDescription: pathDecryptHelpDesc,
	}
}

func (b *backend) pathDecryptWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	ciphertext := d.Get("ciphertext").(string)
	if len(ciphertext) == 0 {
		return logical.ErrorResponse("missing ciphertext to decrypt"), logical.ErrInvalidRequest
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
		return logical.ErrorResponse("policy not found"), logical.ErrInvalidRequest
	}

	plaintext, err := p.Decrypt(context, nonce, ciphertext)
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

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"plaintext": plaintext,
		},
	}
	return resp, nil
}

const pathDecryptHelpSyn = `Decrypt a ciphertext value using a named key`

const pathDecryptHelpDesc = `
This path uses the named key from the request path to decrypt a user
provided ciphertext. The plaintext is returned base64 encoded.
`
