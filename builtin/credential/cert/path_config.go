package cert

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"disable_binding": {
				Type:        framework.TypeBool,
				Default:     false,
				Description: `If set, during renewal, skips the matching of presented client identity with the client identity used during login. Defaults to false.`,
			},
			"enable_identity_alias_metadata": {
				Type:        framework.TypeBool,
				Default:     false,
				Description: `If set, metadata of the certificate including the metadata corresponding to allowed_metadata_extensions will be stored in the alias. Defaults to false.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathConfigWrite,
			logical.ReadOperation:   b.pathConfigRead,
		},
	}
}

func (b *backend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	disableBinding := data.Get("disable_binding").(bool)
	enableIdentityAliasMetadata := data.Get("enable_identity_alias_metadata").(bool)

	entry, err := logical.StorageEntryJSON("config", config{
		DisableBinding:              disableBinding,
		EnableIdentityAliasMetadata: enableIdentityAliasMetadata,
	})
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *backend) pathConfigRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	cfg, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"disable_binding":                cfg.DisableBinding,
		"enable_identity_alias_metadata": cfg.EnableIdentityAliasMetadata,
	}

	return &logical.Response{
		Data: data,
	}, nil
}

// Config returns the configuration for this backend.
func (b *backend) Config(ctx context.Context, s logical.Storage) (*config, error) {
	entry, err := s.Get(ctx, "config")
	if err != nil {
		return nil, err
	}

	// Returning a default configuration if an entry is not found
	var result config
	if entry != nil {
		if err := entry.DecodeJSON(&result); err != nil {
			return nil, fmt.Errorf("error reading configuration: %w", err)
		}
	}
	return &result, nil
}

type config struct {
	DisableBinding              bool `json:"disable_binding"`
	EnableIdentityAliasMetadata bool `json:"enable_identity_alias_metadata"`
}
