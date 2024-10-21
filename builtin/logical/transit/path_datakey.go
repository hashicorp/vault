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
	name := d.Get("name").(string)
	ver := d.Get("key_version").(int)

	plaintext := d.Get("plaintext").(string)
	plaintextAllowed := false
	switch plaintext {
	case "plaintext":
		plaintextAllowed = true
	case "wrapped":
	default:
		return logical.ErrorResponse("Invalid path, must be 'plaintext' or 'wrapped'"), logical.ErrInvalidRequest
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

	newKey := make([]byte, 32)
	bits := d.Get("bits").(int)
	switch bits {
	case 512:
		newKey = make([]byte, 64)
	case 256:
	case 128:
		newKey = make([]byte, 16)
	default:
		return logical.ErrorResponse("invalid bit length"), logical.ErrInvalidRequest
	}
	_, err = rand.Read(newKey)
	if err != nil {
		return nil, err
	}

	factories := make([]any, 0)
	if ps, ok := d.GetOk("padding_scheme"); ok {
		paddingScheme, err := parsePaddingSchemeArg(p.Type, ps)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("padding_scheme argument invalid: %s", err.Error())), logical.ErrInvalidRequest
		}
		factories = append(factories, paddingScheme)

	}
	if p.Type == keysutil.KeyType_MANAGED_KEY {
		managedKeySystemView, ok := b.System().(logical.ManagedKeySystemView)
		if !ok {
			return nil, errors.New("unsupported system view")
		}

		factories = append(factories, ManagedKeyFactory{
			managedKeyParams: keysutil.ManagedKeyParameters{
				ManagedKeySystemView: managedKeySystemView,
				BackendUUID:          b.backendUUID,
				Context:              ctx,
			},
		})
	}

	ciphertext, err := p.EncryptWithFactory(ver, context, nonce, base64.StdEncoding.EncodeToString(newKey), factories...)
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

	keyVersion := ver
	if keyVersion == 0 {
		keyVersion = p.LatestVersion
	}

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"ciphertext":  ciphertext,
			"key_version": keyVersion,
		},
	}

	if len(nonce) > 0 && !nonceAllowed(p) {
		return nil, ErrNonceNotAllowed
	}

	if constants.IsFIPS() && shouldWarnAboutNonceUsage(p, nonce) {
		resp.AddWarning("A provided nonce value was used within FIPS mode, this violates FIPS 140 compliance.")
	}

	if plaintextAllowed {
		resp.Data["plaintext"] = base64.StdEncoding.EncodeToString(newKey)
	}

	return resp, nil
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
