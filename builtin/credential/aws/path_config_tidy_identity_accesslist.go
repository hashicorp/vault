// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package awsauth

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	identityAccessListConfigStorage = "config/tidy/identity-whitelist"
)

func (b *backend) pathConfigTidyIdentityAccessList() *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("%s$", "config/tidy/identity-accesslist"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAWS,
		},

		Fields: map[string]*framework.FieldSchema{
			"safety_buffer": {
				Type:    framework.TypeDurationSecond,
				Default: 259200, // 72h
				Description: `The amount of extra time that must have passed beyond the identity's
expiration, before it is removed from the backend storage.`,
			},
			"disable_periodic_tidy": {
				Type:        framework.TypeBool,
				Default:     false,
				Description: "If set to 'true', disables the periodic tidying of the 'identity-accesslist/<instance_id>' entries.",
			},
		},

		ExistenceCheck: b.pathConfigTidyIdentityAccessListExistenceCheck,

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathConfigTidyIdentityAccessListCreateUpdate,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "identity-access-list-tidy-operation",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigTidyIdentityAccessListCreateUpdate,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "identity-access-list-tidy-operation",
				},
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigTidyIdentityAccessListRead,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "identity-access-list-tidy-settings",
				},
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathConfigTidyIdentityAccessListDelete,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "identity-access-list-tidy-settings",
				},
			},
		},

		HelpSynopsis:    pathConfigTidyIdentityAccessListHelpSyn,
		HelpDescription: pathConfigTidyIdentityAccessListHelpDesc,
	}
}

func (b *backend) pathConfigTidyIdentityAccessListExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := b.lockedConfigTidyIdentities(ctx, req.Storage)
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

func (b *backend) lockedConfigTidyIdentities(ctx context.Context, s logical.Storage) (*tidyWhitelistIdentityConfig, error) {
	b.configMutex.RLock()
	defer b.configMutex.RUnlock()

	return b.nonLockedConfigTidyIdentities(ctx, s)
}

func (b *backend) nonLockedConfigTidyIdentities(ctx context.Context, s logical.Storage) (*tidyWhitelistIdentityConfig, error) {
	entry, err := s.Get(ctx, identityAccessListConfigStorage)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result tidyWhitelistIdentityConfig
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (b *backend) pathConfigTidyIdentityAccessListCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	configEntry, err := b.nonLockedConfigTidyIdentities(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if configEntry == nil {
		configEntry = &tidyWhitelistIdentityConfig{}
	}

	safetyBufferInt, ok := data.GetOk("safety_buffer")
	if ok {
		configEntry.SafetyBuffer = safetyBufferInt.(int)
	} else if req.Operation == logical.CreateOperation {
		configEntry.SafetyBuffer = data.Get("safety_buffer").(int)
	}

	disablePeriodicTidyBool, ok := data.GetOk("disable_periodic_tidy")
	if ok {
		configEntry.DisablePeriodicTidy = disablePeriodicTidyBool.(bool)
	} else if req.Operation == logical.CreateOperation {
		configEntry.DisablePeriodicTidy = data.Get("disable_periodic_tidy").(bool)
	}

	entry, err := logical.StorageEntryJSON(identityAccessListConfigStorage, configEntry)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathConfigTidyIdentityAccessListRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	clientConfig, err := b.lockedConfigTidyIdentities(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if clientConfig == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"safety_buffer":         clientConfig.SafetyBuffer,
			"disable_periodic_tidy": clientConfig.DisablePeriodicTidy,
		},
	}, nil
}

func (b *backend) pathConfigTidyIdentityAccessListDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	return nil, req.Storage.Delete(ctx, identityAccessListConfigStorage)
}

type tidyWhitelistIdentityConfig struct {
	SafetyBuffer        int  `json:"safety_buffer"`
	DisablePeriodicTidy bool `json:"disable_periodic_tidy"`
}

const pathConfigTidyIdentityAccessListHelpSyn = `
Configures the periodic tidying operation of the access list identity entries.
`

const pathConfigTidyIdentityAccessListHelpDesc = `
By default, the expired entries in the access list will be attempted to be removed
periodically. This operation will look for expired items in the list and purges them.
However, there is a safety buffer duration (defaults to 72h), purges the entries
only if they have been persisting this duration, past its expiration time.
`
