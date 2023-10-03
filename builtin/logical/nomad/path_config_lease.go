// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package nomad

import (
	"context"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const leaseConfigKey = "config/lease"

func pathConfigLease(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/lease",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixNomad,
		},

		Fields: map[string]*framework.FieldSchema{
			"ttl": {
				Type:        framework.TypeDurationSecond,
				Description: "Duration before which the issued token needs renewal",
			},
			"max_ttl": {
				Type:        framework.TypeDurationSecond,
				Description: `Duration after which the issued token should not be allowed to be renewed`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathLeaseRead,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "read",
					OperationSuffix: "lease-configuration",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathLeaseUpdate,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "lease",
				},
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathLeaseDelete,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "delete",
					OperationSuffix: "lease-configuration",
				},
			},
		},

		HelpSynopsis:    pathConfigLeaseHelpSyn,
		HelpDescription: pathConfigLeaseHelpDesc,
	}
}

// Sets the lease configuration parameters
func (b *backend) pathLeaseUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := logical.StorageEntryJSON("config/lease", &configLease{
		TTL:    time.Second * time.Duration(d.Get("ttl").(int)),
		MaxTTL: time.Second * time.Duration(d.Get("max_ttl").(int)),
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathLeaseDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(ctx, leaseConfigKey); err != nil {
		return nil, err
	}

	return nil, nil
}

// Returns the lease configuration parameters
func (b *backend) pathLeaseRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	lease, err := b.LeaseConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"ttl":     int64(lease.TTL.Seconds()),
			"max_ttl": int64(lease.MaxTTL.Seconds()),
		},
	}, nil
}

// Lease returns the lease information
func (b *backend) LeaseConfig(ctx context.Context, s logical.Storage) (*configLease, error) {
	entry, err := s.Get(ctx, leaseConfigKey)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result configLease
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Lease configuration information for the secrets issued by this backend
type configLease struct {
	TTL    time.Duration `json:"ttl" mapstructure:"ttl"`
	MaxTTL time.Duration `json:"max_ttl" mapstructure:"max_ttl"`
}

var pathConfigLeaseHelpSyn = "Configure the lease parameters for generated tokens"

var pathConfigLeaseHelpDesc = `
Sets the ttl and max_ttl values for the secrets to be issued by this backend.
Both ttl and max_ttl takes in an integer number of seconds as input as well as
inputs like "1h".
`
