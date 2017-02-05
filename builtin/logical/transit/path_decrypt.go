package transit

import (
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/helper/jsonutil"
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
				Type: framework.TypeString,
				Description: `
The ciphertext to decrypt, provided as returned by encrypt.`,
			},

			"context": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Base64 encoded context for key derivation. Required if key derivation is
enabled.`,
			},

			"nonce": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Base64 encoded nonce value used during encryption. Must be provided if
convergent encryption is enabled for this key and the key was generated with
Vault 0.6.1. Not required for keys created in 0.6.2+.`,
			},

			"batch_input": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Base64 encoded list of items to be decrypted in a single batch. When this
parameter is set, if the parameters 'ciphertext', 'context' and 'nonce' are
also set, they will be ignored. JSON format for the input (which should be
base64 encoded) goes like this:

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
			logical.UpdateOperation: b.pathDecryptWrite,
		},

		HelpSynopsis:    pathDecryptHelpSyn,
		HelpDescription: pathDecryptHelpDesc,
	}
}

func (b *backend) pathDecryptWrite(
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
		batchResponseItems[i].Plaintext = plaintext
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
			"plaintext": batchResponseItems[0].Plaintext,
		}
	}

	return resp, nil
}

const pathDecryptHelpSyn = `Decrypt a ciphertext value using a named key`

const pathDecryptHelpDesc = `
This path uses the named key from the request path to decrypt a user
provided ciphertext. The plaintext is returned base64 encoded.
`
