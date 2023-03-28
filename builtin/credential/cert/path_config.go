// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cert

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const maxCacheSize = 100000

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
			"ocsp_cache_size": {
				Type:        framework.TypeInt,
				Default:     100,
				Description: `The size of the in memory OCSP response cache, shared by all configured certs`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathConfigWrite,
			logical.ReadOperation:   b.pathConfigRead,
		},
	}
}

func (b *backend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if disableBindingRaw, ok := data.GetOk("disable_binding"); ok {
		config.DisableBinding = disableBindingRaw.(bool)
	}
	if enableIdentityAliasMetadataRaw, ok := data.GetOk("enable_identity_alias_metadata"); ok {
		config.EnableIdentityAliasMetadata = enableIdentityAliasMetadataRaw.(bool)
	}
	if cacheSizeRaw, ok := data.GetOk("ocsp_cache_size"); ok {
		cacheSize := cacheSizeRaw.(int)
		if cacheSize < 2 || cacheSize > maxCacheSize {
			return logical.ErrorResponse("invalid cache size, must be >= 2 and <= %d", maxCacheSize), nil
		}
		config.OcspCacheSize = cacheSize
	}
	if err := b.storeConfig(ctx, req.Storage, config); err != nil {
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
		"ocsp_cache_size":                cfg.OcspCacheSize,
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
	OcspCacheSize               int  `json:"ocsp_cache_size"`
}
