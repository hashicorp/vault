// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/logical"
)

type DataKeyParams struct {
	keyVersion int
	bits       int
	context    []byte
	nonce      []byte
	factories  []any
}

func (b *backend) pathDatakey() *framework.Path {
	return &framework.Path{
		Pattern: "datakey/" + framework.GenericNameRegex("plaintext") + "/" + framework.GenericNameRegex("name"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixTransit,
			OperationVerb:   "generate",
			OperationSuffix: "data-key",
		},

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "The backend key used for encrypting the data key",
			},

			"plaintext": {
				Type: framework.TypeString,
				Description: `"plaintext" will return the key in both plaintext and
ciphertext; "wrapped" will return the ciphertext only.`,
			},

			"padding_scheme": {
				Type: framework.TypeString,
				Description: `The padding scheme to use for decrypt. Currently only applies to RSA key types.
Options are 'oaep' or 'pkcs1v15'. Defaults to 'oaep'`,
			},

			"context": {
				Type:        framework.TypeString,
				Description: "Context for key derivation. Required for derived keys.",
			},

			"nonce": {
				Type:        framework.TypeString,
				Description: "Nonce for when convergent encryption v1 is used (only in Vault 0.6.1)",
			},

			"bits": {
				Type: framework.TypeInt,
				Description: `Number of bits for the key; currently 128, 256,
and 512 bits are supported. Defaults to 256.`,
				Default: 256,
			},

			"key_version": {
				Type: framework.TypeInt,
				Description: `The version of the Vault key to use for
encryption of the data key. Must be 0 (for latest)
or a value greater than or equal to the
min_encryption_version configured on the key.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathDatakeyWrite,
		},

		HelpSynopsis:    pathDatakeyHelpSyn,
		HelpDescription: pathDatakeyHelpDesc,
	}
}

func (b *backend) pathDatakeyWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	params, err := getDataKeyParams(d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	name := d.Get("name").(string)

	plaintext := d.Get("plaintext").(string)
	plaintextAllowed := false
	switch plaintext {
	case "plaintext":
		plaintextAllowed = true
	case "wrapped":
	default:
		return logical.ErrorResponse("Invalid path, must be 'plaintext' or 'wrapped'"), logical.ErrInvalidRequest
	}

	// Get the policy
	p, _, err := b.GetPolicy(ctx, keysutil.PolicyRequest{
		Storage: req.Storage,
		Name:    name,
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

	params.factories = make([]any, 0)
	if ps, ok := d.GetOk("padding_scheme"); ok {
		paddingScheme, err := parsePaddingSchemeArg(p.Type, ps)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("padding_scheme argument invalid: %s", err.Error())), logical.ErrInvalidRequest
		}
		params.factories = append(params.factories, paddingScheme)

	}

	keyVersion := params.keyVersion
	if params.keyVersion == 0 {
		keyVersion = p.LatestVersion
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"key_version": keyVersion,
		},
	}

	if len(params.nonce) > 0 && !nonceAllowed(p) {
		return nil, ErrNonceNotAllowed
	}

	if constants.IsFIPS() && shouldWarnAboutNonceUsage(p, params.nonce) {
		resp.AddWarning("A provided nonce value was used within FIPS mode, this violates FIPS 140 compliance.")
	}

	ciphertext, plaintext, err := b.generateDataKey(ctx, p, params)
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

	resp.Data["ciphertext"] = ciphertext

	if plaintextAllowed {
		resp.Data["plaintext"] = plaintext
	}

	return resp, nil
}

func (b *backend) generateDataKey(ctx context.Context, p *keysutil.Policy, params *DataKeyParams) (string, string, error) {
	factories := params.factories
	if p.Type == keysutil.KeyType_MANAGED_KEY {
		managedKeySystemView, ok := b.System().(logical.ManagedKeySystemView)
		if !ok {
			return "", "", errors.New("unsupported system view")
		}

		factories = append(params.factories, ManagedKeyFactory{
			managedKeyParams: keysutil.ManagedKeyParameters{
				ManagedKeySystemView: managedKeySystemView,
				BackendUUID:          b.backendUUID,
				Context:              ctx,
			},
		})
	}

	newKey := make([]byte, params.bits/8)

	_, err := rand.Read(newKey)
	if err != nil {
		return "", "", err
	}

	opts := keysutil.EncryptionOptions{
		KeyVersion: params.keyVersion,
		Context:    params.context,
		Nonce:      params.nonce,
	}

	ciphertext, err := p.EncryptWithOptions(opts, base64.StdEncoding.EncodeToString(newKey), factories...)
	if err != nil {
		return "", "", err
	}

	if ciphertext == "" {
		return "", "", fmt.Errorf("empty ciphertext returned")
	}

	return ciphertext, base64.StdEncoding.EncodeToString(newKey), nil
}

func getDataKeyParams(d *framework.FieldData) (*DataKeyParams, error) {
	params := &DataKeyParams{}

	var err error
	params.keyVersion = d.Get("key_version").(int)

	// Decode the context if any
	if contextRaw, ok := d.GetOk("context"); ok {
		if context, ok := contextRaw.(string); ok && len(context) != 0 {
			params.context, err = base64.StdEncoding.DecodeString(context)
			if err != nil {
				return nil, errors.New("failed to base64-decode context")
			}
		}
	}

	// Decode the nonce if any
	if nonceRaw, ok := d.GetOk("nonce"); ok {
		if nonce, ok := nonceRaw.(string); ok && len(nonce) != 0 {
			params.nonce, err = base64.StdEncoding.DecodeString(nonce)
			if err != nil {
				return nil, errors.New("failed to base64-decode nonce")
			}
		}
	}

	params.bits = d.Get("bits").(int)
	if params.bits != 128 && params.bits != 256 && params.bits != 512 {
		return nil, errors.New("invalid bit length")
	}

	return params, nil
}

const pathDatakeyHelpSyn = `Generate a data key`

const pathDatakeyHelpDesc = `
This path can be used to generate a data key: a random
key of a certain length that can be used for encryption
and decryption, protected by the named backend key. 128, 256,
or 512 bits can be specified; if not specified, the default
is 256 bits. Call with the the "wrapped" path to prevent the
(base64-encoded) plaintext key from being returned along with
the encrypted key, the "plaintext" path returns both.
`
