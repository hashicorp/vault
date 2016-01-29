package transit

import (
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/vault/helper/certutil"
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
	lp, err := b.policies.getPolicy(req, name)
	if err != nil {
		return nil, err
	}

	// Error if invalid policy
	if lp == nil {
		return logical.ErrorResponse("policy not found"), logical.ErrInvalidRequest
	}

	lp.RLock()
	defer lp.RUnlock()

	// Verify if wasn't deleted before we grabbed the lock
	if lp.policy == nil {
		return nil, fmt.Errorf("no existing policy named %s could be found", name)
	}

	plaintext, err := lp.policy.Decrypt(context, ciphertext)
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

	if plaintext == "" {
		return nil, fmt.Errorf("empty plaintext returned")
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
