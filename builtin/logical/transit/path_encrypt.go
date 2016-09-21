package transit

import (
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/hashicorp/vault/helper/errutil"
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

			"nonce": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Nonce for when convergent encryption is used",
			},

			"type": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: "aes256-gcm96",
				Description: `When performing an upsert operation, the type of key
to create. Currently, "aes256-gcm96" (symmetric) is the
only type supported. Defaults to "aes256-gcm96".`,
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

	valueRaw, ok := d.GetOk("plaintext")
	if !ok {
		return logical.ErrorResponse("missing plaintext to encrypt"), logical.ErrInvalidRequest
	}
	value := valueRaw.(string)

	// Decode the context if any
	contextRaw := d.Get("context").(string)
	var context []byte
	var err error
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
	var p *policy
	var lock *sync.RWMutex
	var upserted bool
	if req.Operation == logical.CreateOperation {
		convergent := d.Get("convergent_encryption").(bool)
		if convergent && len(context) == 0 {
			return logical.ErrorResponse("convergent encryption requires derivation to be enabled, so context is required"), nil
		}

		polReq := policyRequest{
			storage:    req.Storage,
			name:       name,
			derived:    len(context) != 0,
			convergent: convergent,
		}

		keyType := d.Get("type").(string)
		switch keyType {
		case "aes256-gcm96":
			polReq.keyType = keyType_AES256_GCM96
		case "ecdsa-p256":
			return logical.ErrorResponse(fmt.Sprintf("key type %v not supported for this operation", keyType)), logical.ErrInvalidRequest
		default:
			return logical.ErrorResponse(fmt.Sprintf("unknown key type %v", keyType)), logical.ErrInvalidRequest
		}

		p, lock, upserted, err = b.lm.GetPolicyUpsert(polReq)

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

	ciphertext, err := p.Encrypt(context, nonce, value)
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
