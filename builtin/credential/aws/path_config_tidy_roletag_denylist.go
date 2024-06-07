// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package awsauth

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	roletagDenyListConfigStorage = "config/tidy/roletag-blacklist"
)

func (b *backend) pathConfigTidyRoletagDenyList() *framework.Path {
	return &framework.Path{
		Pattern: "config/tidy/roletag-denylist$",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAWS,
		},

		Fields: map[string]*framework.FieldSchema{
			"safety_buffer": {
				Type:    framework.TypeDurationSecond,
				Default: 15552000, // 180d
				Description: `The amount of extra time that must have passed beyond the roletag
expiration, before it is removed from the backend storage.
Defaults to 4320h (180 days).`,
			},

			"disable_periodic_tidy": {
				Type:        framework.TypeBool,
				Default:     false,
				Description: "If set to 'true', disables the periodic tidying of deny listed entries.",
			},
		},

		ExistenceCheck: b.pathConfigTidyRoletagDenyListExistenceCheck,

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathConfigTidyRoletagDenyListCreateUpdate,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "role-tag-deny-list-tidy-operation",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigTidyRoletagDenyListCreateUpdate,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "role-tag-deny-list-tidy-operation",
				},
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigTidyRoletagDenyListRead,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "role-tag-deny-list-tidy-settings",
				},
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathConfigTidyRoletagDenyListDelete,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "role-tag-deny-list-tidy-settings",
				},
			},
		},

		HelpSynopsis:    pathConfigTidyRoletagDenyListHelpSyn,
		HelpDescription: pathConfigTidyRoletagDenyListHelpDesc,
	}
}

func (b *backend) pathConfigTidyRoletagDenyListExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := b.lockedConfigTidyRoleTags(ctx, req.Storage)
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

func (b *backend) lockedConfigTidyRoleTags(ctx context.Context, s logical.Storage) (*tidyDenyListRoleTagConfig, error) {
	b.configMutex.RLock()
	defer b.configMutex.RUnlock()

	return b.nonLockedConfigTidyRoleTags(ctx, s)
}

func (b *backend) nonLockedConfigTidyRoleTags(ctx context.Context, s logical.Storage) (*tidyDenyListRoleTagConfig, error) {
	entry, err := s.Get(ctx, roletagDenyListConfigStorage)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result tidyDenyListRoleTagConfig
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathConfigTidyRoletagDenyListCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	configEntry, err := b.nonLockedConfigTidyRoleTags(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if configEntry == nil {
		configEntry = &tidyDenyListRoleTagConfig{}
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

	entry, err := logical.StorageEntryJSON(roletagDenyListConfigStorage, configEntry)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathConfigTidyRoletagDenyListRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	clientConfig, err := b.lockedConfigTidyRoleTags(ctx, req.Storage)
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

func (b *backend) pathConfigTidyRoletagDenyListDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	return nil, req.Storage.Delete(ctx, roletagDenyListConfigStorage)
}

type tidyDenyListRoleTagConfig struct {
	SafetyBuffer        int  `json:"safety_buffer"`
	DisablePeriodicTidy bool `json:"disable_periodic_tidy"`
}

const pathConfigTidyRoletagDenyListHelpSyn = `
Configures the periodic tidying operation of the deny listed role tag entries.
`

const pathConfigTidyRoletagDenyListHelpDesc = `
By default, the expired entries in the deny list will be attempted to be removed
periodically. This operation will look for expired items in the list and purges them.
However, there is a safety buffer duration (defaults to 72h), purges the entries
only if they have been persisting this duration, past its expiration time.
`
