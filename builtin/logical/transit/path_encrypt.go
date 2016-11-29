package transit

import (
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/keysutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/mitchellh/mapstructure"
)

const (
	MB int = 1048 * 1048
)

type BatchEncryptionItemRequest struct {
	Context   string `json:"context" structs:"context" mapstructure:"context"`
	Plaintext string `json:"plaintext" structs:"plaintext" mapstructure:"plaintext"`
	Nonce     string `json:"nonce" structs:"nonce" mapstructure:"nonce"`
}

type BatchEncryptionItemResponse struct {
	CipherText string `json:"ciphertext" structs:"ciphertext" mapstructure:"ciphertext"`
	Error      string `json:"error" structs:"error" mapstructure:"error"`
}

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
				Description: "Base64 encoded context for key derivation. Required for derived keys.",
			},

			"nonce": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Base64 encoded nonce for when convergent encryption is used",
			},

			"type": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: "aes256-gcm96",
				Description: `
This parameter is required when encryption key is expected to be created.
When performing an upsert operation, the type of key to create. Currently,
"aes256-gcm96" (symmetric) is the only type supported. Defaults to
"aes256-gcm96".`,
			},

			"convergent_encryption": &framework.FieldSchema{
				Type: framework.TypeBool,
				Description: `
This parameter will only be used when a key is expected to be created.  Whether
to support convergent encryption. This is only supported when using a key with
key derivation enabled and will require all requests to carry both a context
and 96-bit (12-byte) nonce. The given nonce will be used in place of a randomly
generated nonce. As a result, when the same context and nonce are supplied, the
same ciphertext is generated. It is *very important* when using this mode that
you ensure that all nonces are unique for a given context.  Failing to do so
will severely impact the ciphertext's security.`,
			},

			"derived": &framework.FieldSchema{
				Type: framework.TypeBool,
				Description: `
This parameter will only be used when a key is expected to be created.  Enables
key derivation mode. This allows for per-transaction unique keys for encryption
operations.`,
			},

			"batch": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: "",
				Description: `
Base64 encoded list of items to be encrypted in a single batch. The size of the
input is limited to 4MB. When this parameter is set, the parameters
'plaintext', 'context' and 'nonce' will be ignored. JSON format for the input
goes like this:

[
  {
    "context": "context1",
    "plaintext": "dGhlIHF1aWNrIGJyb3duIGZveA=="
  },
  {
    "context": "context2",
    "plaintext": "dGhlIHF1aWNrIGJyb3duIGZveA=="
  },
  ...
]`,
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
	var err error

	convergent := d.Get("convergent_encryption").(bool)

	// Get the policy
	var p *keysutil.Policy
	var lock *sync.RWMutex
	var upserted bool
	if req.Operation == logical.CreateOperation {
		polReq := keysutil.PolicyRequest{
			Storage:    req.Storage,
			Name:       name,
			Derived:    d.Get("derived").(bool),
			Convergent: convergent,
		}

		keyType := d.Get("type").(string)
		switch keyType {
		case "aes256-gcm96":
			polReq.KeyType = keysutil.KeyType_AES256_GCM96
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

	batchInputRaw := d.Get("batch").(string)
	var batchInput []byte
	if len(batchInputRaw) != 0 {
		batchInput, err = base64.StdEncoding.DecodeString(batchInputRaw)
		if err != nil {
			return logical.ErrorResponse("failed to base64-decode batch"), logical.ErrInvalidRequest
		}
	}

	if len(batchInput) != 0 {
		// The size of the batch input is limited to 4MB
		if len(batchInput) > 4*MB {
			return logical.ErrorResponse("Input for batch encryption should be less than 4MB"), logical.ErrInvalidRequest
		}

		var batchInputArray []interface{}
		if err := jsonutil.DecodeJSON([]byte(batchInput), &batchInputArray); err != nil {
			return nil, err
		}

		var batchItems []BatchEncryptionItemRequest
		var batchResponseItems []BatchEncryptionItemResponse

		// Process batch request items. If encryption of any request
		// item fails, respectively mark the error in the response
		// collection and continue to process other items.
		for _, batchItem := range batchInputArray {
			var item BatchEncryptionItemRequest
			if err := mapstructure.Decode(batchItem, &item); err != nil {
				batchResponseItems = append(batchResponseItems, BatchEncryptionItemResponse{
					Error: fmt.Sprintf("failed to decode the request: %v", err),
				})
				continue
			}
			batchItems = append(batchItems, item)

			if item.Plaintext == "" {
				batchResponseItems = append(batchResponseItems, BatchEncryptionItemResponse{
					Error: "missing plaintext to encrypt",
				})
				continue
			}

			// Decode the context
			var itemContext []byte
			if len(item.Context) != 0 {
				itemContext, err = base64.StdEncoding.DecodeString(item.Context)
				if err != nil {
					batchResponseItems = append(batchResponseItems, BatchEncryptionItemResponse{
						Error: "failed to base64-decode context",
					})
					continue
				}
			}

			// Decode the nonce
			var itemNonce []byte
			if len(item.Nonce) != 0 {
				itemNonce, err = base64.StdEncoding.DecodeString(item.Nonce)
				if err != nil {
					batchResponseItems = append(batchResponseItems, BatchEncryptionItemResponse{
						Error: "failed to base64-decode nonce",
					})
					continue
				}
			}

			ciphertext, err := p.Encrypt(itemContext, itemNonce, item.Plaintext)
			if err != nil {
				batchResponseItems = append(batchResponseItems, BatchEncryptionItemResponse{
					Error: fmt.Sprintf("encryption failed: %s", err.Error()),
				})
				continue
			}

			if ciphertext == "" {
				batchResponseItems = append(batchResponseItems, BatchEncryptionItemResponse{
					Error: "empty ciphertext returned",
				})
				continue
			}

			batchResponseItems = append(batchResponseItems, BatchEncryptionItemResponse{
				CipherText: ciphertext,
			})
		}

		if len(batchItems) != len(batchResponseItems) {
			return nil, fmt.Errorf("number of request and the response items does not match")
		}

		batchResponseJSON, err := jsonutil.EncodeJSON(batchResponseItems)
		if err != nil {
			return nil, fmt.Errorf("failed to JSON encode batch response")
		}

		// Generate the response
		resp := &logical.Response{
			Data: map[string]interface{}{
				"data": batchResponseJSON,
			},
		}

		if req.Operation == logical.CreateOperation && !upserted {
			resp.AddWarning("Attempted creation of the key during the encrypt operation, but it was created beforehand")
		}
		return resp, nil
	}

	valueRaw, ok := d.GetOk("plaintext")
	if !ok {
		return logical.ErrorResponse("missing plaintext to encrypt"), logical.ErrInvalidRequest
	}
	value := valueRaw.(string)

	// Decode the context if any
	contextRaw := d.Get("context").(string)
	var context []byte
	if len(contextRaw) != 0 {
		context, err = base64.StdEncoding.DecodeString(contextRaw)
		if err != nil {
			return logical.ErrorResponse("failed to base64-decode context"), logical.ErrInvalidRequest
		}
	}

	if convergent && len(context) == 0 {
		return logical.ErrorResponse("convergent encryption requires derivation to be enabled, so context is required"), nil
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
