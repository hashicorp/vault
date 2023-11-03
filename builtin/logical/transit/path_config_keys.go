// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const keysConfigPath = "config/keys"

type keysConfig struct {
	DisableUpsert bool `json:"disable_upsert"`
}

var defaultKeysConfig = keysConfig{
	DisableUpsert: false,
}

func (b *backend) pathConfigKeys() *framework.Path {
	return &framework.Path{
		Pattern: "config/keys",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixTransit,
		},

		Fields: map[string]*framework.FieldSchema{
			"disable_upsert": {
				Type: framework.TypeBool,
				Description: `Whether to allow automatic upserting (creation) of
keys on the encrypt endpoint.`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigKeysWrite,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "keys",
				},
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigKeysRead,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "keys-configuration",
				},
			},
		},

		HelpSynopsis:    pathConfigKeysHelpSyn,
		HelpDescription: pathConfigKeysHelpDesc,
	}
}

func (b *backend) readConfigKeys(ctx context.Context, req *logical.Request) (*keysConfig, error) {
	entry, err := req.Storage.Get(ctx, keysConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch keys configuration: %w", err)
	}

	var cfg keysConfig
	if entry == nil {
		cfg = defaultKeysConfig
		return &cfg, nil
	}

	if err := entry.DecodeJSON(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode keys configuration: %w", err)
	}

	return &cfg, nil
}

func (b *backend) writeConfigKeys(ctx context.Context, req *logical.Request, cfg *keysConfig) error {
	entry, err := logical.StorageEntryJSON(keysConfigPath, cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal keys configuration: %w", err)
	}

	return req.Storage.Put(ctx, entry)
}

func respondConfigKeys(cfg *keysConfig) *logical.Response {
	return &logical.Response{
		Data: map[string]interface{}{
			"disable_upsert": cfg.DisableUpsert,
		},
	}
}

func (b *backend) pathConfigKeysWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	upsert := d.Get("disable_upsert").(bool)

	cfg, err := b.readConfigKeys(ctx, req)
	if err != nil {
		return nil, err
	}

	modified := false

	if cfg.DisableUpsert != upsert {
		cfg.DisableUpsert = upsert
		modified = true
	}

	if modified {
		if err := b.writeConfigKeys(ctx, req, cfg); err != nil {
			return nil, err
		}
	}

	return respondConfigKeys(cfg), nil
}

func (b *backend) pathConfigKeysRead(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	cfg, err := b.readConfigKeys(ctx, req)
	if err != nil {
		return nil, err
	}

	return respondConfigKeys(cfg), nil
}

const pathConfigKeysHelpSyn = `Configuration common across all keys`

const pathConfigKeysHelpDesc = `
This path is used to configure common functionality across all keys. Currently,
this supports limiting the ability to automatically create new keys when an
unknown key is used for encryption (upsert).
`
