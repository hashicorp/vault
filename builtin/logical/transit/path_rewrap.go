// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

var ErrNonceNotAllowed = errors.New("provided nonce not allowed for this key")

type RewrapBatchRequestItem struct {
	// Context for key derivation. This is required for derived keys.
	Context string `json:"context" structs:"context" mapstructure:"context"`

	// DecodedContext is the base64 decoded version of Context
	DecodedContext []byte

	// Ciphertext for decryption
	Ciphertext string `json:"ciphertext" structs:"ciphertext" mapstructure:"ciphertext"`

	// Nonce to be used when v1 convergent encryption is used
	Nonce string `json:"nonce" structs:"nonce" mapstructure:"nonce"`

	// The key version to be used for encryption
	KeyVersion int `json:"key_version" structs:"key_version" mapstructure:"key_version"`

	// DecodedNonce is the base64 decoded version of Nonce
	DecodedNonce []byte

	// Associated Data for AEAD ciphers
	AssociatedData string `json:"associated_data" struct:"associated_data" mapstructure:"associated_data"`

	// Reference is an arbitrary caller supplied string value that will be placed on the
	// batch response to ease correlation between inputs and outputs
	Reference string `json:"reference" structs:"reference" mapstructure:"reference"`

	// EncryptPaddingScheme specifies the RSA padding scheme for encryption
	EncryptPaddingScheme string `json:"encrypt_padding_scheme" structs:"encrypt_padding_scheme" mapstructure:"encrypt_padding_scheme"`

	// DecryptPaddingScheme specifies the RSA padding scheme for decryption
	DecryptPaddingScheme string `json:"decrypt_padding_scheme" structs:"decrypt_padding_scheme" mapstructure:"decrypt_padding_scheme"`
}

func (b *backend) pathRewrap() *framework.Path {
	return &framework.Path{
		Pattern: "rewrap/" + framework.GenericNameRegex("name"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixTransit,
			OperationVerb:   "rewrap",
		},

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the key",
			},

			"ciphertext": {
				Type:        framework.TypeString,
				Description: "Ciphertext value to rewrap",
			},

			"encrypt_padding_scheme": {
				Type: framework.TypeString,
				Description: `The padding scheme to use for rewrap's encrypt step. Currently only applies to RSA key types.
Options are 'oaep' or 'pkcs1v15'. Defaults to 'oaep'`,
			},

			"decrypt_padding_scheme": {
				Type: framework.TypeString,
				Description: `The padding scheme to use for rewrap's decrypt step. Currently only applies to RSA key types.
Options are 'oaep' or 'pkcs1v15'. Defaults to 'oaep'`,
			},

			"context": {
				Type:        framework.TypeString,
				Description: "Base64 encoded context for key derivation. Required for derived keys.",
			},

			"nonce": {
				Type:        framework.TypeString,
				Description: "Nonce for when convergent encryption is used",
			},

			"key_version": {
				Type: framework.TypeInt,
				Description: `The version of the key to use for encryption.
Must be 0 (for latest) or a value greater than or equal
to the min_encryption_version configured on the key.`,
			},

			"batch_input": {
				Type: framework.TypeSlice,
				Description: `
Specifies a list of items to be re-encrypted in a single batch. When this parameter is set,
if the parameters 'ciphertext', 'context' and 'nonce' are also set, they will be ignored.
Any batch output will preserve the order of the batch input.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathRewrapWrite,
		},

		HelpSynopsis:    pathRewrapHelpSyn,
		HelpDescription: pathRewrapHelpDesc,
	}
}

func (b *backend) pathRewrapWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	batchInputRaw := d.Raw["batch_input"]
	var batchInputItems []RewrapBatchRequestItem
	var err error
	if batchInputRaw != nil {
		err = mapstructure.Decode(batchInputRaw, &batchInputItems)
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

		batchInputItems = make([]RewrapBatchRequestItem, 1)
		batchInputItems[0] = RewrapBatchRequestItem{
			Ciphertext: ciphertext,
			Context:    d.Get("context").(string),
			Nonce:      d.Get("nonce").(string),
			KeyVersion: d.Get("key_version").(int),
		}
		if ps, ok := d.GetOk("decrypt_padding_scheme"); ok {
			batchInputItems[0].DecryptPaddingScheme = ps.(string)
		}
		if ps, ok := d.GetOk("encrypt_padding_scheme"); ok {
			batchInputItems[0].EncryptPaddingScheme = ps.(string)
		}
	}

	batchResponseItems := make([]EncryptBatchResponseItem, len(batchInputItems))
	contextSet := len(batchInputItems[0].Context) != 0

	for i, item := range batchInputItems {
		if (len(item.Context) == 0 && contextSet) || (len(item.Context) != 0 && !contextSet) {
			return logical.ErrorResponse("context should be set either in all the request blocks or in none"), logical.ErrInvalidRequest
		}

		if item.Ciphertext == "" {
			batchResponseItems[i].Error = "missing ciphertext to decrypt"
			continue
		}

		// Decode the context
		if len(item.Context) != 0 {
			batchInputItems[i].DecodedContext, err = base64.StdEncoding.DecodeString(item.Context)
			if err != nil {
				batchResponseItems[i].Error = err.Error()
				continue
			}
		}

		// Decode the nonce
		if len(item.Nonce) != 0 {
			batchInputItems[i].DecodedNonce, err = base64.StdEncoding.DecodeString(item.Nonce)
			if err != nil {
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

	warnAboutNonceUsage := false
	for i, item := range batchInputItems {
		if batchResponseItems[i].Error != "" {
			continue
		}

		var factories []any
		if item.DecryptPaddingScheme != "" {
			paddingScheme, err := parsePaddingSchemeArg(p.Type, item.DecryptPaddingScheme)
			if err != nil {
				batchResponseItems[i].Error = fmt.Sprintf("'[%d].decrypt_padding_scheme' invalid: %s", i, err.Error())
				continue
			}
			factories = append(factories, paddingScheme)
		}
		if item.Nonce != "" && !nonceAllowed(p) {
			batchResponseItems[i].Error = ErrNonceNotAllowed.Error()
			continue
		}

		plaintext, err := p.DecryptWithFactory(item.DecodedContext, item.DecodedNonce, item.Ciphertext, factories...)
		if err != nil {
			switch err.(type) {
			case errutil.UserError:
				batchResponseItems[i].Error = err.Error()
				continue
			default:
				return nil, err
			}
		}

		factories = make([]any, 0)
		if item.EncryptPaddingScheme != "" {
			paddingScheme, err := parsePaddingSchemeArg(p.Type, item.EncryptPaddingScheme)
			if err != nil {
				batchResponseItems[i].Error = fmt.Sprintf("'[%d].encrypt_padding_scheme' invalid: %s", i, err.Error())
				continue
			}
			factories = append(factories, paddingScheme)
			factories = append(factories, keysutil.PaddingScheme(item.EncryptPaddingScheme))
		}
		if !warnAboutNonceUsage && shouldWarnAboutNonceUsage(p, item.DecodedNonce) {
			warnAboutNonceUsage = true
		}

		ciphertext, err := p.EncryptWithFactory(item.KeyVersion, item.DecodedContext, item.DecodedNonce, plaintext, factories...)
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

		keyVersion := item.KeyVersion
		if keyVersion == 0 {
			keyVersion = p.LatestVersion
		}

		batchResponseItems[i].Ciphertext = ciphertext
		batchResponseItems[i].KeyVersion = keyVersion
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
			return logical.ErrorResponse(batchResponseItems[0].Error), logical.ErrInvalidRequest
		}
		resp.Data = map[string]interface{}{
			"ciphertext":  batchResponseItems[0].Ciphertext,
			"key_version": batchResponseItems[0].KeyVersion,
		}
	}

	if constants.IsFIPS() && warnAboutNonceUsage {
		resp.AddWarning("A provided nonce value was used within FIPS mode, this violates FIPS 140 compliance.")
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
