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

func (b *backend) pathRewrap() *framework.Path {
	return &framework.Path{
		Pattern: "rewrap/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the key",
			},

			"ciphertext": {
				Type:        framework.TypeString,
				Description: "Ciphertext value to rewrap",
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
	var batchInputItems []BatchRequestItem
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

		batchInputItems = make([]BatchRequestItem, 1)
		batchInputItems[0] = BatchRequestItem{
			Ciphertext: ciphertext,
			Context:    d.Get("context").(string),
			Nonce:      d.Get("nonce").(string),
			KeyVersion: d.Get("key_version").(int),
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

	warnAboutNonceUsage := false
	for i, item := range batchInputItems {
		if batchResponseItems[i].Error != "" {
			continue
		}

		if item.Nonce != "" && !nonceAllowed(p) {
			batchResponseItems[i].Error = ErrNonceNotAllowed.Error()
			continue
		}

		plaintext, err := p.Decrypt(item.DecodedContext, item.DecodedNonce, item.Ciphertext)
		if err != nil {
			switch err.(type) {
			case errutil.UserError:
				batchResponseItems[i].Error = err.Error()
				continue
			default:
				p.Unlock()
				return nil, err
			}
		}

		if !warnAboutNonceUsage && shouldWarnAboutNonceUsage(p, item.DecodedNonce) {
			warnAboutNonceUsage = true
		}

		ciphertext, err := p.Encrypt(item.KeyVersion, item.DecodedContext, item.DecodedNonce, plaintext)
		if err != nil {
			switch err.(type) {
			case errutil.UserError:
				batchResponseItems[i].Error = err.Error()
				continue
			case errutil.InternalError:
				p.Unlock()
				return nil, err
			default:
				p.Unlock()
				return nil, err
			}
		}

		if ciphertext == "" {
			p.Unlock()
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
			p.Unlock()
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

	p.Unlock()
	return resp, nil
}

const pathRewrapHelpSyn = `Rewrap ciphertext`

const pathRewrapHelpDesc = `
After key rotation, this function can be used to rewrap the given ciphertext or
a batch of given ciphertext blocks with the latest version of the named key.
If the given ciphertext is already using the latest version of the key, this
function is a no-op.
`
