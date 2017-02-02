package transit

import (
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

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

			"batch_input": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Base64 encoded list of items to be rewrapped in a single batch. When this
parameter is set, if the parameters 'ciphertext', 'context' and 'nonce' are
also set, they will be ignored. JSON format for the input (which should be
bae64 encoded) goes like this:

[
  {
    "context": "c2FtcGxlY29udGV4dA==",
    "ciphertext": "vault:v1:/DupSiSbX/ATkGmKAmhqD0tvukByrx6gmps7dVI="
  },
  {
    "context": "YW5vdGhlcnNhbXBsZWNvbnRleHQ=",
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
	batchInputRaw := d.Get("batch_input").(string)
	var batchInputItems []BatchRequestItem
	var batchInput []byte
	var err error
	if len(batchInputRaw) != 0 {
		batchInput, err = base64.StdEncoding.DecodeString(batchInputRaw)
		if err != nil {
			return logical.ErrorResponse("failed to base64-decode batch input"), logical.ErrInvalidRequest
		}

		if err := jsonutil.DecodeJSON([]byte(batchInput), &batchInputItems); err != nil {
			return nil, fmt.Errorf("invalid input: %v", err)
		}

		if len(batchInputItems) == 0 {
			return logical.ErrorResponse("missing batch input to process"), logical.ErrInvalidRequest
		}
	} else {
		ciphertext := d.Get("ciphertext").(string)
		if len(ciphertext) == 0 {
			return logical.ErrorResponse("missing ciphertext to decrypt"), logical.ErrInvalidRequest
		}

		batchInputItems = make([]BatchRequestItem, 1)
		batchInputItems[0] = BatchRequestItem{
			Ciphertext: ciphertext,
		}

		// Decode the context
		contextRaw := d.Get("context").(string)
		if len(contextRaw) != 0 {
			batchInputItems[0].Context, err = base64.StdEncoding.DecodeString(contextRaw)
			if err != nil {
				return logical.ErrorResponse("failed to base64-decode context"), logical.ErrInvalidRequest
			}
		}

		// Decode the nonce
		nonceRaw := d.Get("nonce").(string)
		if len(nonceRaw) != 0 {
			batchInputItems[0].Nonce, err = base64.StdEncoding.DecodeString(nonceRaw)
			if err != nil {
				return logical.ErrorResponse("failed to base64-decode nonce"), logical.ErrInvalidRequest
			}
		}
	}

	batchResponseItems := make([]BatchResponseItem, len(batchInputItems))
	contextSet := len(batchInputItems[0].Context) != 0

	for i, item := range batchInputItems {
		if (len(item.Context) == 0 && contextSet) || (len(item.Context) != 0 && !contextSet) {
			return logical.ErrorResponse("context should be set either in all the request blocks or in none"), logical.ErrInvalidRequest
		}

		if item.Ciphertext == "" {
			batchResponseItems[i].Error = "missing ciphertext to decrypt"
			continue
		}
	}

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

	for i, item := range batchInputItems {
		if batchResponseItems[i].Error != "" {
			continue
		}

		plaintext, err := p.Decrypt(item.Context, item.Nonce, item.Ciphertext)
		if err != nil {
			switch err.(type) {
			case errutil.UserError:
				batchResponseItems[i].Error = err.Error()
				continue
			default:
				return nil, err
			}
		}

		ciphertext, err := p.Encrypt(item.Context, item.Nonce, plaintext)
		if err != nil {
			switch err.(type) {
			case errutil.UserError:
				batchResponseItems[i].Error = err.Error()
				continue
			case errutil.InternalError:
				return nil, err
			default:
				return nil, err
			}
		}

		if ciphertext == "" {
			return nil, fmt.Errorf("empty ciphertext returned for input item %d", i)
		}

		batchResponseItems[i].Ciphertext = ciphertext
	}

	resp := &logical.Response{}
	if len(batchInputRaw) != 0 {
		batchResponseJSON, err := jsonutil.EncodeJSON(batchResponseItems)
		if err != nil {
			return nil, fmt.Errorf("failed to JSON encode batch response")
		}
		resp.Data = map[string]interface{}{
			"batch_results": string(batchResponseJSON),
		}
	} else {
		if batchResponseItems[0].Error != "" {
			return logical.ErrorResponse(batchResponseItems[0].Error), logical.ErrInvalidRequest
		}
		resp.Data = map[string]interface{}{
			"ciphertext": batchResponseItems[0].Ciphertext,
		}
	}

	return resp, nil
}

const pathRewrapHelpSyn = `Rewrap ciphertext`

const pathRewrapHelpDesc = `
After key rotation, this function can be used to rewrap the given ciphertext or
a batch of given ciphertext blocks with the latest version of the named key.
If the given ciphertext is already using the latest version of the key, this
function is a no-op.
`
