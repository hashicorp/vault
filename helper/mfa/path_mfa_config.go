package mfa

import (
	"context"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathMFAConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `mfa_config`,
		Fields: map[string]*framework.FieldSchema{
			"type": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Enables MFA with given backend (available: duo)",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathMFAConfigWrite,
			logical.ReadOperation:   b.pathMFAConfigRead,
		},

		HelpSynopsis:    pathMFAConfigHelpSyn,
		HelpDescription: pathMFAConfigHelpDesc,
	}
}

func (b *backend) MFAConfig(ctx context.Context, req *logical.Request) (*MFAConfig, error) {
	entry, err := req.Storage.Get(ctx, "mfa_config")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	var result MFAConfig
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (b *backend) pathMFAConfigWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := logical.StorageEntryJSON("mfa_config", MFAConfig{
		Type: d.Get("type").(string),
	})
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathMFAConfigRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	config, err := b.MFAConfig(ctx, req)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"type": config.Type,
		},
	}, nil
}

type MFAConfig struct {
	Type string `json:"type"`
}

const pathMFAConfigHelpSyn = `
Configure multi factor backend. 
`

const pathMFAConfigHelpDesc = `
This endpoint allows you to turn on multi-factor authentication with a given backend.
Currently only Duo is supported.
`
