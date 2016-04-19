package aws

import (
	"github.com/fatih/structs"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigTidyWhitelistIdentity(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/tidy/whitelist/identity$",
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
				Description: "If set to 'true', disables the periodic tidying of the 'whitelist/identity/<instance_id>' entries and 'whitelist/identity' entries.",
			},
		},

		ExistenceCheck: b.pathConfigTidyWhitelistIdentityExistenceCheck,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathConfigTidyWhitelistIdentityCreateUpdate,
			logical.UpdateOperation: b.pathConfigTidyWhitelistIdentityCreateUpdate,
		},

		HelpSynopsis:    pathConfigTidyWhitelistIdentityHelpSyn,
		HelpDescription: pathConfigTidyWhitelistIdentityHelpDesc,
	}
}

func (b *backend) pathConfigTidyWhitelistIdentityExistenceCheck(req *logical.Request, data *framework.FieldData) (bool, error) {
	b.configMutex.RLock()
	defer b.configMutex.RUnlock()

	entry, err := configTidyWhitelistIdentity(req.Storage)
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

func configTidyWhitelistIdentity(s logical.Storage) (*tidyWhitelistIdentityConfig, error) {
	entry, err := s.Get("config/tidy/whitelist/identity")
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

func (b *backend) pathConfigTidyWhitelistIdentityCreateUpdate(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.Lock()
	defer b.configMutex.Unlock()
	configEntry, err := configTidyWhitelistIdentity(req.Storage)
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

	entry, err := logical.StorageEntryJSON("config/tidy/whitelist/identity", configEntry)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathConfigTidyWhitelistIdentityRead(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.RLock()
	defer b.configMutex.RUnlock()

	clientConfig, err := configTidyWhitelistIdentity(req.Storage)
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

func (b *backend) pathConfigTidyWhitelistIdentityDelete(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	if err := req.Storage.Delete("config/tidy/whitelist/identity"); err != nil {
		return nil, err
	}

	return nil, nil
}

type tidyWhitelistIdentityConfig struct {
	SafetyBuffer        int  `json:"safety_buffer" structs:"safety_buffer" mapstructure:"safety_buffer"`
	DisablePeriodicTidy bool `json:"disable_periodic_tidy" structs:"disable_periodic_tidy" mapstructure:"disable_periodic_tidy"`
}

const pathConfigTidyWhitelistIdentityHelpSyn = `
Configures the periodic tidying operation of the whitelisted identity entries.
`
const pathConfigTidyWhitelistIdentityHelpDesc = `
By default, the expired entries in teb whitelist will be attempted to be removed
periodically. This operation will look for expired items in the list and purge them.
However, there is a safety buffer duration (defaults to 72h), which purges the entries,
only if they have been persisting this duration, past its expiration time.
`
