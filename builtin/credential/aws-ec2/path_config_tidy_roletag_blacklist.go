package awsec2

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
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

func (b *backend) pathConfigTidyRoletagBlacklistExistenceCheck(req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := b.lockedConfigTidyRoleTags(req.Storage)
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

func (b *backend) lockedConfigTidyRoleTags(s logical.Storage) (*tidyBlacklistRoleTagConfig, error) {
	b.configMutex.RLock()
	defer b.configMutex.RUnlock()

	return b.nonLockedConfigTidyRoleTags(s)
}

func (b *backend) nonLockedConfigTidyRoleTags(s logical.Storage) (*tidyBlacklistRoleTagConfig, error) {
	entry, err := s.Get(roletagBlacklistConfigPath)
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

func (b *backend) pathConfigTidyRoletagBlacklistCreateUpdate(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	configEntry, err := b.nonLockedConfigTidyRoleTags(req.Storage)
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

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathConfigTidyRoletagBlacklistRead(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	clientConfig, err := b.lockedConfigTidyRoleTags(req.Storage)
	if err != nil {
		return nil, err
	}
	if clientConfig == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: structs.New(clientConfig).Map(),
	}, nil
}

func (b *backend) pathConfigTidyRoletagBlacklistDelete(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	return nil, req.Storage.Delete(roletagBlacklistConfigPath)
}

type tidyBlacklistRoleTagConfig struct {
	SafetyBuffer        int  `json:"safety_buffer" structs:"safety_buffer" mapstructure:"safety_buffer"`
	DisablePeriodicTidy bool `json:"disable_periodic_tidy" structs:"disable_periodic_tidy" mapstructure:"disable_periodic_tidy"`
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
