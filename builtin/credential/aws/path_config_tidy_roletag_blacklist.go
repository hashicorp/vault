package awsauth

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	roletagBlacklistConfigPath = "config/tidy/roletag-blacklist"
)

func pathConfigTidyRoletagBlacklist(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("%s$", roletagBlacklistConfigPath),
		Fields: map[string]*framework.FieldSchema{
			"safety_buffer": &framework.FieldSchema{
				Type:    framework.TypeDurationSecond,
				Default: 15552000, //180d
				Description: `The amount of extra time that must have passed beyond the roletag
expiration, before it is removed from the backend storage.
Defaults to 4320h (180 days).`,
			},

			"disable_periodic_tidy": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Default:     false,
				Description: "If set to 'true', disables the periodic tidying of blacklisted entries.",
			},
		},

		ExistenceCheck: b.pathConfigTidyRoletagBlacklistExistenceCheck,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathConfigTidyRoletagBlacklistCreateUpdate,
			logical.UpdateOperation: b.pathConfigTidyRoletagBlacklistCreateUpdate,
			logical.ReadOperation:   b.pathConfigTidyRoletagBlacklistRead,
			logical.DeleteOperation: b.pathConfigTidyRoletagBlacklistDelete,
		},

		HelpSynopsis:    pathConfigTidyRoletagBlacklistHelpSyn,
		HelpDescription: pathConfigTidyRoletagBlacklistHelpDesc,
	}
}

func (b *backend) pathConfigTidyRoletagBlacklistExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := b.lockedConfigTidyRoleTags(ctx, req.Storage)
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

func (b *backend) lockedConfigTidyRoleTags(ctx context.Context, s logical.Storage) (*tidyBlacklistRoleTagConfig, error) {
	b.configMutex.RLock()
	defer b.configMutex.RUnlock()

	return b.nonLockedConfigTidyRoleTags(ctx, s)
}

func (b *backend) nonLockedConfigTidyRoleTags(ctx context.Context, s logical.Storage) (*tidyBlacklistRoleTagConfig, error) {
	entry, err := s.Get(ctx, roletagBlacklistConfigPath)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result tidyBlacklistRoleTagConfig
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathConfigTidyRoletagBlacklistCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	configEntry, err := b.nonLockedConfigTidyRoleTags(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if configEntry == nil {
		configEntry = &tidyBlacklistRoleTagConfig{}
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

	entry, err := logical.StorageEntryJSON(roletagBlacklistConfigPath, configEntry)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathConfigTidyRoletagBlacklistRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
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

func (b *backend) pathConfigTidyRoletagBlacklistDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	return nil, req.Storage.Delete(ctx, roletagBlacklistConfigPath)
}

type tidyBlacklistRoleTagConfig struct {
	SafetyBuffer        int  `json:"safety_buffer"`
	DisablePeriodicTidy bool `json:"disable_periodic_tidy"`
}

const pathConfigTidyRoletagBlacklistHelpSyn = `
Configures the periodic tidying operation of the blacklisted role tag entries.
`
const pathConfigTidyRoletagBlacklistHelpDesc = `
By default, the expired entries in the blacklist will be attempted to be removed
periodically. This operation will look for expired items in the list and purges them.
However, there is a safety buffer duration (defaults to 72h), purges the entries
only if they have been persisting this duration, past its expiration time.
`
