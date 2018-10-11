package transit

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/helper/keysutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathTrim() *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("name") + "/trim",
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the key",
			},
			"min_version": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `The minimum version for the key ring. All versions before this version will be
permanently removed. This should always be greater than both
'min_decryption_version' and 'min_encryption_version'. This is not allowed to
be set when either 'min_encryption_version' or 'min_decryption_version' is set
to zero.`,
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

		p, _, err := b.lm.GetPolicy(ctx, keysutil.PolicyRequest{
			Storage: req.Storage,
			Name:    name,
		})
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

		minVersionRaw, ok := d.GetOk("min_version")
		if !ok {
			return logical.ErrorResponse("missing min_version"), nil
		}

		// Ensure that cache doesn't get corrupted in error cases
		originalMinVersion := p.MinVersion
		defer func() {
			if retErr != nil || (resp != nil && resp.IsError()) {
				p.MinVersion = originalMinVersion
			}
		}()

		p.MinVersion = minVersionRaw.(int)

		switch {
		case p.MinVersion < originalMinVersion:
			return logical.ErrorResponse("minimum version cannot be decremented"), nil
		case p.MinEncryptionVersion == 0:
			return logical.ErrorResponse("minimum version cannot be set when minimum encryption version is not set"), nil
		case p.MinDecryptionVersion == 0:
			return logical.ErrorResponse("minimum version cannot be set when minimum decryption version is not set"), nil
		case p.MinVersion > p.MinEncryptionVersion:
			return logical.ErrorResponse("minimum version cannot be greater than minmum encryption version"), nil
		case p.MinVersion > p.MinDecryptionVersion:
			return logical.ErrorResponse("minimum version cannot be greater than minimum decryption version"), nil
		case p.MinVersion < 0:
			return logical.ErrorResponse("minimum version cannot be negative"), nil
		case p.MinVersion == 0:
			return logical.ErrorResponse("minimum version should be positive"), nil
		case p.MinVersion < p.MinVersion:
			return logical.ErrorResponse(fmt.Sprintf("minimum version cannot be less than the already set value of %d", p.MinVersion)), nil
		}

		return nil, p.Persist(ctx, req.Storage)
	}
}

const pathTrimHelpSyn = `Trim key versions in the named key`

const pathTrimHelpDesc = `
This path is used to trim key versions of a named key. Trimming only happens
from the lower end of version numbers.
`
