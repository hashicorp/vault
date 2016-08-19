package pki

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// CRLConfig holds basic CRL configuration information
type crlConfig struct {
	Expiry string `json:"expiry" mapstructure:"expiry" structs:"expiry"`
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
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathCRLRead,
			logical.UpdateOperation: b.pathCRLWrite,
		},

		HelpSynopsis:    pathConfigCRLHelpSyn,
		HelpDescription: pathConfigCRLHelpDesc,
	}
}

func (b *backend) CRL(s logical.Storage) (*crlConfig, error) {
	entry, err := s.Get("config/crl")
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

func (b *backend) pathCRLRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.CRL(req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"expiry": config.Expiry,
		},
	}, nil
}

func (b *backend) pathCRLWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	expiry := d.Get("expiry").(string)

	_, err := time.ParseDuration(expiry)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("Given expiry could not be decoded: %s", err)), nil
	}

	config := &crlConfig{
		Expiry: expiry,
	}

	entry, err := logical.StorageEntryJSON("config/crl", config)
	if err != nil {
		return nil, err
	}
	err = req.Storage.Put(entry)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

const pathConfigCRLHelpSyn = `
Configure the CRL expiration.
`

const pathConfigCRLHelpDesc = `
This endpoint allows configuration of the CRL lifetime.
`
