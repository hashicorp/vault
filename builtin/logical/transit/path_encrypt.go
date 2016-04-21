package transit

import (
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathEncrypt() *framework.Path {
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
			logical.CreateOperation: b.pathEncryptWrite,
			logical.UpdateOperation: b.pathEncryptWrite,
		},

		ExistenceCheck: b.pathEncryptExistenceCheck,

		HelpSynopsis:    pathEncryptHelpSyn,
		HelpDescription: pathEncryptHelpDesc,
	}
}

func (b *backend) pathEncryptExistenceCheck(
	req *logical.Request, d *framework.FieldData) (bool, error) {
	name := d.Get("name").(string)
	lp, err := b.policies.getPolicy(req.Storage, name)
	if err != nil {
		return false, err
	}

	return lp != nil, nil
}

func (b *backend) pathEncryptWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	value := d.Get("plaintext").(string)
	if len(value) == 0 {
		return logical.ErrorResponse("missing plaintext to encrypt"), logical.ErrInvalidRequest
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
	lp, err := b.policies.getPolicy(req.Storage, name)
	if err != nil {
		return nil, err
	}

	// Error or upsert if invalid policy
	if lp == nil {
		if req.Operation != logical.CreateOperation {
			return logical.ErrorResponse("policy not found"), logical.ErrInvalidRequest
		}

		// Get a write lock
		b.policies.Lock()

		isDerived := len(context) != 0

		// This also checks to make sure one hasn't been created since we grabbed the write lock
		lp, err = b.policies.generatePolicy(req.Storage, name, isDerived)
		// If the error is that the policy has been created in the interim we
		// will get the policy back, so only consider it an error if err is not
		// nil and we do not get a policy back
		if err != nil && lp != nil {
			b.policies.Unlock()
			return nil, err
		}
		b.policies.Unlock()
	}

	lp.RLock()
	defer lp.RUnlock()

	// Verify if wasn't deleted before we grabbed the lock
	if lp.Policy() == nil {
		return nil, fmt.Errorf("no existing policy named %s could be found", name)
	}

	ciphertext, err := lp.Policy().Encrypt(context, value)
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
