// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathTrim() *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("name") + "/trim",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixTransit,
			OperationVerb:   "trim",
			OperationSuffix: "key",
		},

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the key",
			},
			"min_available_version": {
				Type: framework.TypeInt,
				Description: `
The minimum available version for the key ring. All versions before this
version will be permanently deleted. This value can at most be equal to the
lesser of 'min_decryption_version' and 'min_encryption_version'. This is not
allowed to be set when either 'min_encryption_version' or
'min_decryption_version' is set to zero.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathTrimUpdate(),
		},

		HelpSynopsis:    pathTrimHelpSyn,
		HelpDescription: pathTrimHelpDesc,
	}
}

func (b *backend) pathTrimUpdate() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (resp *logical.Response, retErr error) {
		name := d.Get("name").(string)

		p, _, err := b.GetPolicy(ctx, keysutil.PolicyRequest{
			Storage: req.Storage,
			Name:    name,
		}, b.GetRandomReader())
		if err != nil {
			return nil, err
		}
		if p == nil {
			return logical.ErrorResponse("invalid key name"), logical.ErrInvalidRequest
		}
		if !b.System().CachingDisabled() {
			p.Lock(true)
		}
		defer p.Unlock()

		minAvailableVersionRaw, ok, err := d.GetOkErr("min_available_version")
		if err != nil {
			return nil, err
		}
		if !ok {
			return logical.ErrorResponse("missing min_available_version"), nil
		}
		minAvailableVersion := minAvailableVersionRaw.(int)

		originalMinAvailableVersion := p.MinAvailableVersion

		switch {
		case minAvailableVersion < originalMinAvailableVersion:
			return logical.ErrorResponse("minimum available version cannot be decremented"), nil
		case p.MinEncryptionVersion == 0:
			return logical.ErrorResponse("minimum available version cannot be set when minimum encryption version is not set"), nil
		case p.MinDecryptionVersion == 0:
			return logical.ErrorResponse("minimum available version cannot be set when minimum decryption version is not set"), nil
		case minAvailableVersion > p.MinEncryptionVersion:
			return logical.ErrorResponse("minimum available version cannot be greater than minmum encryption version"), nil
		case minAvailableVersion > p.MinDecryptionVersion:
			return logical.ErrorResponse("minimum available version cannot be greater than minimum decryption version"), nil
		case minAvailableVersion < 0:
			return logical.ErrorResponse("minimum available version cannot be negative"), nil
		case minAvailableVersion == 0:
			return logical.ErrorResponse("minimum available version should be positive"), nil
		}

		// Ensure that cache doesn't get corrupted in error cases
		p.MinAvailableVersion = minAvailableVersion
		if err := p.Persist(ctx, req.Storage); err != nil {
			p.MinAvailableVersion = originalMinAvailableVersion
			return nil, err
		}

		b.TryRecordObservationWithRequest(ctx, req, ObservationTypeTransitKeyTrim, b.keyPolicyObservationMetadata(p))

		return b.formatKeyPolicy(p, nil)
	}
}

const pathTrimHelpSyn = `Trim key versions of a named key`

const pathTrimHelpDesc = `
This path is used to trim key versions of a named key. Trimming only happens
from the lower end of version numbers.
`
