// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/logical"
)

type DecryptBatchResponseItem struct {
	// Plaintext for the ciphertext present in the corresponding batch
	// request item
	Plaintext string `json:"plaintext" structs:"plaintext" mapstructure:"plaintext"`

	// Error, if set represents a failure encountered while encrypting a
	// corresponding batch request item
	Error string `json:"error,omitempty" structs:"error" mapstructure:"error"`

	// Reference is an arbitrary caller supplied string value that will be placed on the
	// batch response to ease correlation between inputs and outputs
	Reference string `json:"reference" structs:"reference" mapstructure:"reference"`
}

func (b *backend) pathDecrypt() *framework.Path {
	return &framework.Path{
		Pattern: "decrypt/" + framework.GenericNameRegex("name"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixTransit,
			OperationVerb:   "decrypt",
		},

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the key",
			},

			"ciphertext": {
				Type: framework.TypeString,
				Description: `
The ciphertext to decrypt, provided as returned by encrypt.`,
			},

			"padding_scheme": {
				Type: framework.TypeString,
				Description: `The padding scheme to use for decrypt. Currently only applies to RSA key types.
Options are 'oaep' or 'pkcs1v15'. Defaults to 'oaep'`,
			},

			"context": {
				Type: framework.TypeString,
				Description: `
Base64 encoded context for key derivation. Required if key derivation is
enabled.`,
			},

			"nonce": {
				Type: framework.TypeString,
				Description: `
Base64 encoded nonce value used during encryption. Must be provided if
convergent encryption is enabled for this key and the key was generated with
Vault 0.6.1. Not required for keys created in 0.6.2+.`,
			},

			"partial_failure_response_code": {
				Type: framework.TypeInt,
				Description: `
Ordinarily, if a batch item fails to decrypt due to a bad input, but other batch items succeed, 
the HTTP response code is 400 (Bad Request).  Some applications may want to treat partial failures differently.
Providing the parameter returns the given response code integer instead of a 400 in this case.  If all values fail
HTTP 400 is still returned.`,
			},

			"associated_data": {
				Type: framework.TypeString,
				Description: `
When using an AEAD cipher mode, such as AES-GCM, this parameter allows
passing associated data (AD/AAD) into the encryption function; this data
must be passed on subsequent decryption requests but can be transited in
plaintext. On successful decryption, both the ciphertext and the associated
data are attested not to have been tampered with.
                `,
			},

			"batch_input": {
				Type: framework.TypeSlice,
				Description: `
Specifies a list of items to be decrypted in a single batch. When this
parameter is set, if the parameters 'ciphertext', 'context' and 'nonce' are
also set, they will be ignored. Any batch output will preserve the order
of the batch input.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathDecryptWrite,
		},

		HelpSynopsis:    pathDecryptHelpSyn,
		HelpDescription: pathDecryptHelpDesc,
	}
}

func (b *backend) pathDecryptWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	batchInputRaw := d.Raw["batch_input"]
	var batchInputItems []BatchRequestItem
	var err error
	if batchInputRaw != nil {
		err = decodeDecryptBatchRequestItems(batchInputRaw, &batchInputItems)
		if err != nil {
			return nil, fmt.Errorf("failed to parse batch input: %w", err)
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
			Ciphertext:     ciphertext,
			Context:        d.Get("context").(string),
			Nonce:          d.Get("nonce").(string),
			AssociatedData: d.Get("associated_data").(string),
		}
		if ps, ok := d.GetOk("padding_scheme"); ok {
			batchInputItems[0].PaddingScheme = ps.(string)
		}
	}

	batchResponseItems := make([]DecryptBatchResponseItem, len(batchInputItems))
	contextSet := len(batchInputItems[0].Context) != 0

	userErrorInBatch := false
	internalErrorInBatch := false

	for i, item := range batchInputItems {
		if (len(item.Context) == 0 && contextSet) || (len(item.Context) != 0 && !contextSet) {
			return logical.ErrorResponse("context should be set either in all the request blocks or in none"), logical.ErrInvalidRequest
		}

		if item.Ciphertext == "" {
			userErrorInBatch = true
			batchResponseItems[i].Error = "missing ciphertext to decrypt"
			continue
		}

		// Decode the context
		if len(item.Context) != 0 {
			batchInputItems[i].DecodedContext, err = base64.StdEncoding.DecodeString(item.Context)
			if err != nil {
				userErrorInBatch = true
				batchResponseItems[i].Error = err.Error()
				continue
			}
		}

		// Decode the nonce
		if len(item.Nonce) != 0 {
			batchInputItems[i].DecodedNonce, err = base64.StdEncoding.DecodeString(item.Nonce)
			if err != nil {
				userErrorInBatch = true
				batchResponseItems[i].Error = err.Error()
				continue
			}
		}
	}

	// Get the policy
	p, _, err := b.GetPolicy(ctx, keysutil.PolicyRequest{
		Storage: req.Storage,
		Name:    d.Get("name").(string),
	}, b.GetRandomReader())
	if err != nil {
		return nil, err
	}
	if p == nil {
		return logical.ErrorResponse("encryption key not found"), logical.ErrInvalidRequest
	}
	if !b.System().CachingDisabled() {
		p.Lock(false)
	}
	defer p.Unlock()

	successesInBatch := false
	for i, item := range batchInputItems {
		if batchResponseItems[i].Error != "" {
			continue
		}

		var factories []any
		if item.PaddingScheme != "" {
			paddingScheme, err := parsePaddingSchemeArg(p.Type, item.PaddingScheme)
			if err != nil {
				batchResponseItems[i].Error = fmt.Sprintf("'[%d].padding_scheme' invalid: %s", i, err.Error())
				continue
			}
			factories = append(factories, paddingScheme)
		}
		if item.AssociatedData != "" {
			if !p.Type.AssociatedDataSupported() {
				batchResponseItems[i].Error = fmt.Sprintf("'[%d].associated_data' provided for non-AEAD cipher suite %v", i, p.Type.String())
				continue
			}

			factories = append(factories, AssocDataFactory{item.AssociatedData})
		}

		if p.Type == keysutil.KeyType_MANAGED_KEY {
			managedKeySystemView, ok := b.System().(logical.ManagedKeySystemView)
			if !ok {
				batchResponseItems[i].Error = errors.New("unsupported system view").Error()
			}

			factories = append(factories, ManagedKeyFactory{
				managedKeyParams: keysutil.ManagedKeyParameters{
					ManagedKeySystemView: managedKeySystemView,
					BackendUUID:          b.backendUUID,
					Context:              ctx,
				},
			})
		}

		plaintext, err := p.DecryptWithFactory(item.DecodedContext, item.DecodedNonce, item.Ciphertext, factories...)
		if err != nil {
			switch err.(type) {
			case errutil.InternalError:
				internalErrorInBatch = true
			default:
				userErrorInBatch = true
			}
			batchResponseItems[i].Error = err.Error()
			continue
		}
		successesInBatch = true
		batchResponseItems[i].Plaintext = plaintext
	}

	resp := &logical.Response{}
	if batchInputRaw != nil {
		// Copy the references
		for i := range batchInputItems {
			batchResponseItems[i].Reference = batchInputItems[i].Reference
		}
		resp.Data = map[string]interface{}{
			"batch_results": batchResponseItems,
		}
	} else {
		if batchResponseItems[0].Error != "" {
			if internalErrorInBatch {
				return nil, errutil.InternalError{Err: batchResponseItems[0].Error}
			}

			return logical.ErrorResponse(batchResponseItems[0].Error), logical.ErrInvalidRequest
		}
		resp.Data = map[string]interface{}{
			"plaintext": batchResponseItems[0].Plaintext,
		}
	}

	return batchRequestResponse(d, resp, req, successesInBatch, userErrorInBatch, internalErrorInBatch)
}

const pathDecryptHelpSyn = `Decrypt a ciphertext value using a named key`

const pathDecryptHelpDesc = `
This path uses the named key from the request path to decrypt a user
provided ciphertext. The plaintext is returned base64 encoded.
`
