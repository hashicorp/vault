package transit

import (
	"encoding/base64"
	"fmt"
	"sync"

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
	p, lock, err := b.lm.GetPolicyShared(req.Storage, name)
	if lock != nil {
		defer lock.RUnlock()
	}
	if err != nil {
		return false, err
	}
	return p != nil, nil
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
	var err error
	if len(contextRaw) != 0 {
		context, err = base64.StdEncoding.DecodeString(contextRaw)
		if err != nil {
			return logical.ErrorResponse("failed to decode context as base64"), logical.ErrInvalidRequest
		}
	}

	// Get the policy
	var p *Policy
	var lock *sync.RWMutex
	var upserted bool
	if req.Operation == logical.CreateOperation {
		p, lock, upserted, err = b.lm.GetPolicyUpsert(req.Storage, name, len(context) != 0)
	} else {
		p, lock, err = b.lm.GetPolicyShared(req.Storage, name)
	}
	if lock != nil {
		defer lock.RUnlock()
	}
	if err != nil {
		return nil, err
	}
	if p == nil {
		return logical.ErrorResponse("policy not found"), logical.ErrInvalidRequest
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

	if req.Operation == logical.CreateOperation && !upserted {
		resp.AddWarning("Attempted creation of the key during the encrypt operation, but it was created beforehand")
	}

	return resp, nil
}

const pathEncryptHelpSyn = `Encrypt a plaintext value using a named key`

const pathEncryptHelpDesc = `
This path uses the named key from the request path to encrypt a user
provided plaintext. The plaintext must be base64 encoded.
`
