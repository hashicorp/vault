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
)

// BatchEncryptionItemRequest represents an item in the batch encryption
// request
type BatchEncryptionItemRequest struct {
	// Context for key derivation. This is required for derived keys.
	Context string `json:"context" structs:"context" mapstructure:"context"`

	// DecodedContext, for internal use, which is the base64 decoded version of
	// the Context field
	DecodedContext []byte `json:"decoded_context" structs:"decoded_context" mapstructure:"decoded_context"`

	// Plaintext for encryption
	Plaintext string `json:"plaintext" structs:"plaintext" mapstructure:"plaintext"`

	// Nonce to be used when v1 convergent encryption is used
	Nonce string `json:"nonce" structs:"nonce" mapstructure:"nonce"`

	// DecodedNonce, for internal use, which is the base64 decoded version of
	// the Nonce field
	DecodedNonce []byte `json:"decoded_nonce" structs:"decoded_nonce" mapstructure:"decoded_nonce"`
}

// BatchEncryptionItemResponse represents an item in the batch encryption
// response
type BatchEncryptionItemResponse struct {
	// Ciphertext for the plaintext present in the corresponding batch
	// request item
	Ciphertext string `json:"ciphertext" structs:"ciphertext" mapstructure:"ciphertext"`

	// Error, if set represents a failure encountered while encrypting a
	// corresponding batch request item
	Error string `json:"error" structs:"error" mapstructure:"error"`
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
				Description: "Base64 encoded plaintext value to be encrypted",
			},

			"context": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Base64 encoded context for key derivation. Required if key derivation is enabled",
			},

			"nonce": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Base64 encoded nonce value. Must be provided if convergent encryption is
enabled for this key and the key was generated with Vault 0.6.1. Not required
for keys created in 0.6.2+. The value must be exactly 96 bits (12 bytes) long
and the user must ensure that for any given context (and thus, any given
encryption key) this nonce value is **never reused**.
`,
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

			"batch": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Base64 encoded list of items to be encrypted in a single batch. When this
parameter is set, if the parameters 'plaintext', 'context' and 'nonce' are also
set, they will be ignored. JSON format for the input goes like this:

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

	batchInputRaw := d.Get("batch").(string)
	var batchInput []byte
	var batchInputItems []BatchEncryptionItemRequest
	if len(batchInputRaw) != 0 {
		batchInput, err = base64.StdEncoding.DecodeString(batchInputRaw)
		if err != nil {
			return logical.ErrorResponse("failed to base64-decode batch input"), logical.ErrInvalidRequest
		}

		if err := jsonutil.DecodeJSON([]byte(batchInput), &batchInputItems); err != nil {
			return nil, fmt.Errorf("invalid input: %v", err)
		}

	} else {
		valueRaw, ok := d.GetOk("plaintext")
		if !ok {
			return logical.ErrorResponse("missing plaintext to encrypt"), logical.ErrInvalidRequest
		}

		batchInputItems = make([]BatchEncryptionItemRequest, 1)
		batchInputItems[0] = BatchEncryptionItemRequest{
			Plaintext: valueRaw.(string),
			Context:   d.Get("context").(string),
			Nonce:     d.Get("nonce").(string),
		}
	}

	contextSet := batchInputItems[0].Context != ""
	if len(batchInputItems) == 0 {
		return logical.ErrorResponse("missing input to process"), logical.ErrInvalidRequest
	}

	batchResponseItems := make([]BatchEncryptionItemResponse, len(batchInputItems))

	// Before processing the batch request items, get the policy. If the
	// policy is supposed to be upserted, then determine if 'derived' is to
	// be set or not, based on the presence of 'context' field in all the
	// input items.
	for i, item := range batchInputItems {
		if (item.Context == "" && contextSet) || (item.Context != "" && !contextSet) {
			return logical.ErrorResponse("context should be set either in all the request blocks or in none"), logical.ErrInvalidRequest
		}

		_, err := base64.StdEncoding.DecodeString(item.Plaintext)
		if err != nil {
			batchResponseItems[i] = BatchEncryptionItemResponse{
				Error: "failed to base64-decode plaintext",
			}
			continue
		}

		// Decode the context
		if len(item.Context) != 0 {
			batchInputItems[i].DecodedContext, err = base64.StdEncoding.DecodeString(item.Context)
			if err != nil {
				batchResponseItems[i] = BatchEncryptionItemResponse{
					Error: "failed to base64-decode context",
				}
				continue
			}
		}

		// Decode the nonce
		if len(item.Nonce) != 0 {
			batchInputItems[i].DecodedNonce, err = base64.StdEncoding.DecodeString(item.Nonce)
			if err != nil {
				batchResponseItems[i] = BatchEncryptionItemResponse{
					Error: "failed to base64-decode nonce",
				}
				continue
			}
		}
	}

	// Get the policy
	var p *keysutil.Policy
	var lock *sync.RWMutex
	var upserted bool
	if req.Operation == logical.CreateOperation {
		convergent := d.Get("convergent_encryption").(bool)
		if convergent && !contextSet {
			return logical.ErrorResponse("convergent encryption requires derivation to be enabled, so context is required"), nil
		}

		polReq := keysutil.PolicyRequest{
			Storage:    req.Storage,
			Name:       name,
			Derived:    contextSet,
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

	// Process batch request items. If encryption of any request
	// item fails, respectively mark the error in the response
	// collection and continue to process other items.
	for i, item := range batchInputItems {
		if batchResponseItems[i].Error != "" {
			continue
		}

		ciphertext, err := p.Encrypt(item.DecodedContext, item.DecodedNonce, item.Plaintext)
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
			batchResponseItems[i] = BatchEncryptionItemResponse{
				Error: "empty ciphertext returned",
			}
			continue
		}

		batchResponseItems[i] = BatchEncryptionItemResponse{
			Ciphertext: ciphertext,
		}
	}

	resp := &logical.Response{}
	if len(batchInputRaw) != 0 {
		batchResponseJSON, err := jsonutil.EncodeJSON(batchResponseItems)
		if err != nil {
			return nil, fmt.Errorf("failed to JSON encode batch response")
		}
		resp.Data = map[string]interface{}{
			"data": string(batchResponseJSON),
		}
	} else {
		if batchResponseItems[0].Error != "" {
			return nil, fmt.Errorf(batchResponseItems[0].Error)
		}
		resp.Data = map[string]interface{}{
			"ciphertext": batchResponseItems[0].Ciphertext,
		}
	}

	if req.Operation == logical.CreateOperation && !upserted {
		resp.AddWarning("Attempted creation of the key during the encrypt operation, but it was created beforehand")
	}
	return resp, nil
}

const pathEncryptHelpSyn = `Encrypt a plaintext value or a batch of plaintext
blocks using a named key`

const pathEncryptHelpDesc = `
This path uses the named key from the request path to encrypt a user provided
plaintext or a batch of plaintext blocks. The plaintext must be base64 encoded.
`
