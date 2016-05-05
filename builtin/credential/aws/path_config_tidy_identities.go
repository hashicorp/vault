package aws

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	identityWhitelistConfigPath = "config/tidy/identity-whitelist"
)

func pathConfigTidyIdentities(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("%s$", identityWhitelistConfigPath),
		Fields: map[string]*framework.FieldSchema{
			"safety_buffer": &framework.FieldSchema{
				Type:    framework.TypeDurationSecond,
				Default: 259200, //72h
				Description: `The amount of extra time that must have passed beyond the identity's
expiration, before it is removed from the backend storage.`,
			},
			"disable_periodic_tidy": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Default:     false,
				Description: "If set to 'true', disables the periodic tidying of the 'whitelist/identity/<instance_id>' entries.",
			},
		},

		ExistenceCheck: b.pathConfigTidyIdentitiesExistenceCheck,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathConfigTidyIdentitiesCreateUpdate,
			logical.UpdateOperation: b.pathConfigTidyIdentitiesCreateUpdate,
			logical.ReadOperation:   b.pathConfigTidyIdentitiesRead,
			logical.DeleteOperation: b.pathConfigTidyIdentitiesDelete,
		},

		HelpSynopsis:    pathConfigTidyIdentitiesHelpSyn,
		HelpDescription: pathConfigTidyIdentitiesHelpDesc,
	}
}

func (b *backend) pathConfigTidyIdentitiesExistenceCheck(req *logical.Request, data *framework.FieldData) (bool, error) {
	b.configMutex.RLock()
	defer b.configMutex.RUnlock()

	entry, err := b.configTidyIdentities(req.Storage)
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

func (b *backend) configTidyIdentities(s logical.Storage) (*tidyWhitelistIdentityConfig, error) {
	entry, err := s.Get(identityWhitelistConfigPath)
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

func (b *backend) pathConfigTidyIdentitiesCreateUpdate(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	configEntry, err := b.configTidyIdentities(req.Storage)
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

	entry, err := logical.StorageEntryJSON(identityWhitelistConfigPath, configEntry)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathConfigTidyIdentitiesRead(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.RLock()
	defer b.configMutex.RUnlock()

	clientConfig, err := b.configTidyIdentities(req.Storage)
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

func (b *backend) pathConfigTidyIdentitiesDelete(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	return nil, req.Storage.Delete(identityWhitelistConfigPath)
}

type tidyWhitelistIdentityConfig struct {
	SafetyBuffer        int  `json:"safety_buffer" structs:"safety_buffer" mapstructure:"safety_buffer"`
	DisablePeriodicTidy bool `json:"disable_periodic_tidy" structs:"disable_periodic_tidy" mapstructure:"disable_periodic_tidy"`
}

const pathConfigTidyIdentitiesHelpSyn = `
Configures the periodic tidying operation of the whitelisted identity entries.
`
const pathConfigTidyIdentitiesHelpDesc = `
By default, the expired entries in the whitelist will be attempted to be removed
periodically. This operation will look for expired items in the list and purge them.
However, there is a safety buffer duration (defaults to 72h), purges the entries
only if they have been persisting this duration, past its expiration time.
`
