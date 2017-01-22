package transit

import (
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/mitchellh/mapstructure"
)

// BatchRewrapItemRequest represents an item in the batch rewrap
// request
type BatchRewrapItemRequest struct {
	// Context for key derivation. This is required for derived keys.
	Context string `json:"context" structs:"context" mapstructure:"context"`

	// Ciphertext which needs rewrap
	Ciphertext string `json:"ciphertext" structs:"ciphertext" mapstructure:"ciphertext"`

	// Nonce to be used when v1 convergent encryption is used
	Nonce string `json:"nonce" structs:"nonce" mapstructure:"nonce"`
}

// BatchRewrapItemResponse represents an item in the batch rewrap
// response
type BatchRewrapItemResponse struct {
	// Ciphertext is a rewrapped version of the ciphertext in the corresponding
	// batch request item
	Ciphertext string `json:"ciphertext" structs:"ciphertext" mapstructure:"ciphertext"`

	// Error, if set represents a failure encountered while encrypting rewrapping a
	// corresponding batch request item
	Error string `json:"error" structs:"error" mapstructure:"error"`
}

func (b *backend) pathRewrap() *framework.Path {
	return &framework.Path{
		Pattern: "rewrap/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the key",
			},

			"ciphertext": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Ciphertext value to rewrap",
			},

			"context": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Base64 encoded context for key derivation. Required for derived keys.",
			},

			"nonce": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Nonce for when convergent encryption is used",
			},

			"batch": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Base64 encoded list of items to be rewrapped in a single batch. When this
parameter is set, if the parameters 'ciphertext', 'context' and 'nonce' are
also set, they will be ignored. JSON format for the input goes like this:
[
  {
    "context": "context1",
    "ciphertext": "vault:v1:/DupSiSbX/ATkGmKAmhqD0tvukByrx6gmps7dVI="
  },
  {
    "context": "context2",
    "ciphertext": "vault:v1:XjsPWPjqPrBi1N2Ms2s1QM798YyFWnO4TR4lsFA="
  },
  ...
]`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathRewrapWrite,
		},

		HelpSynopsis:    pathRewrapHelpSyn,
		HelpDescription: pathRewrapHelpDesc,
	}
}

func (b *backend) pathRewrapWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	// Get the policy
	p, lock, err := b.lm.GetPolicyShared(req.Storage, d.Get("name").(string))
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
			return logical.ErrorResponse("failed to base64-decode batch input"), logical.ErrInvalidRequest
		}
	} else {
		ciphertext := d.Get("ciphertext").(string)
		if len(ciphertext) == 0 {
			return logical.ErrorResponse("missing ciphertext to decrypt"), logical.ErrInvalidRequest
		}

		var singleItemBatch []BatchRewrapItemRequest
		singleItemBatch = append(singleItemBatch, BatchRewrapItemRequest{
			Ciphertext: ciphertext,
			Context:    d.Get("context").(string),
			Nonce:      d.Get("nonce").(string),
		})

		batchInput, err = jsonutil.EncodeJSON(singleItemBatch)
		if err != nil {
			return nil, fmt.Errorf("failed to encode batch input")
		}
	}

	var batchInputArray []interface{}
	if err := jsonutil.DecodeJSON([]byte(batchInput), &batchInputArray); err != nil {
		return nil, fmt.Errorf("invalid input: %v", err)
	}

	var batchItems []BatchRewrapItemRequest
	var batchResponseItems []BatchRewrapItemResponse
	for _, batchItem := range batchInputArray {
		var item BatchRewrapItemRequest
		if err := mapstructure.Decode(batchItem, &item); err != nil {
			return logical.ErrorResponse(fmt.Sprintf("failed to parse the input: %v", err)), logical.ErrInvalidRequest
		}
		batchItems = append(batchItems, item)

		if item.Ciphertext == "" {
			batchResponseItems = append(batchResponseItems, BatchRewrapItemResponse{
				Error: "missing ciphertext to decrypt",
			})
			continue
		}

		var itemContext []byte
		if len(item.Context) != 0 {
			itemContext, err = base64.StdEncoding.DecodeString(item.Context)
			if err != nil {
				batchResponseItems = append(batchResponseItems, BatchRewrapItemResponse{
					Error: "failed to base64-decode context",
				})
				continue
			}
		}

		var itemNonce []byte
		if len(item.Nonce) != 0 {
			itemNonce, err = base64.StdEncoding.DecodeString(item.Nonce)
			if err != nil {
				batchResponseItems = append(batchResponseItems, BatchRewrapItemResponse{
					Error: "failed to base64-decode nonce",
				})
			}
			continue
		}

		plaintext, err := p.Decrypt(itemContext, itemNonce, item.Ciphertext)
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

		if plaintext == "" {
			batchResponseItems = append(batchResponseItems, BatchRewrapItemResponse{
				Error: "empty plaintext returned during rewrap",
			})
			continue
		}

		ciphertext, err := p.Encrypt(itemContext, itemNonce, plaintext)
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
			batchResponseItems = append(batchResponseItems, BatchRewrapItemResponse{
				Error: "empty ciphertext returned",
			})
			continue
		}

		batchResponseItems = append(batchResponseItems, BatchRewrapItemResponse{
			Ciphertext: ciphertext,
		})
	}

	if len(batchItems) != len(batchResponseItems) {
		return nil, fmt.Errorf("number of request and the number of response items do not match")
	}

	if len(batchResponseItems) == 0 {
		return nil, fmt.Errorf("number of response items cannot be zero")
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

	return resp, nil

	/*
		name := d.Get("name").(string)

		value := d.Get("ciphertext").(string)
		if len(value) == 0 {
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
		// Error if invalid policy
		if p == nil {
			return logical.ErrorResponse("policy not found"), logical.ErrInvalidRequest
		}

		plaintext, err := p.Decrypt(context, nonce, value)
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

		if plaintext == "" {
			return nil, fmt.Errorf("empty plaintext returned during rewrap")
		}

		ciphertext, err := p.Encrypt(context, nonce, plaintext)
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
		return resp, nil
	*/
}

const pathRewrapHelpSyn = `Rewrap ciphertext`

const pathRewrapHelpDesc = `
After key rotation, this function can be used to rewrap the
given ciphertext with the latest version of the named key.
If the given ciphertext is already using the latest version
of the key, this function is a no-op.
`
