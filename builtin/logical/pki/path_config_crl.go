package pki

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// CRLConfig holds basic CRL configuration information
type crlConfig struct {
	Expiry  string `json:"expiry" mapstructure:"expiry"`
	Disable bool   `json:"disable"`
}

func pathConfigCRL(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/crl",
		Fields: map[string]*framework.FieldSchema{
			"expiry": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The amount of time the generated CRL should be
valid; defaults to 72 hours`,
				Default: "72h",
			},
			"disable": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Description: `If set to true, disables generating the CRL entirely.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathCRLRead,
			logical.UpdateOperation: b.pathCRLWrite,
		},

		HelpSynopsis:    pathConfigCRLHelpSyn,
		HelpDescription: pathConfigCRLHelpDesc,
	}
}

func (b *backend) CRL(ctx context.Context, s logical.Storage) (*crlConfig, error) {
	entry, err := s.Get(ctx, "config/crl")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result crlConfig
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathCRLRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.CRL(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"expiry":  config.Expiry,
			"disable": config.Disable,
		},
	}, nil
}

func (b *backend) pathCRLWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	config, err := b.CRL(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		config = &crlConfig{}
	}

	if expiryRaw, ok := d.GetOk("expiry"); ok {
		expiry := expiryRaw.(string)
		_, err := time.ParseDuration(expiry)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("given expiry could not be decoded: %s", err)), nil
		}
		config.Expiry = expiry
	}

	var oldDisable bool
	if disableRaw, ok := d.GetOk("disable"); ok {
		oldDisable = config.Disable
		config.Disable = disableRaw.(bool)
	}

	entry, err := logical.StorageEntryJSON("config/crl", config)
	if err != nil {
		return nil, err
	}
	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	if oldDisable != config.Disable {
		// It wasn't disabled but now it is, rotate
		crlErr := buildCRL(ctx, b, req, true)
		switch crlErr.(type) {
		case errutil.UserError:
			return logical.ErrorResponse(fmt.Sprintf("Error during CRL building: %s", crlErr)), nil
		case errutil.InternalError:
			return nil, errwrap.Wrapf("error encountered during CRL building: {{err}}", crlErr)
		}
	}

	return nil, nil
}

const pathConfigCRLHelpSyn = `
Configure the CRL expiration.
`

const pathConfigCRLHelpDesc = `
This endpoint allows configuration of the CRL lifetime.
`
