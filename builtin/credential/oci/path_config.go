// Copyright © 2019, Oracle and/or its affiliates.
package ociauth

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/framework"
	"strings"
)

const HOME_TENANCY_ID_CONFIG_NAME = "homeTenancyId"
const BASE_CONFIG_PATH = "config/"

var allowedConfigNamesForCreate = map[string]string{
	HOME_TENANCY_ID_CONFIG_NAME: HOME_TENANCY_ID_CONFIG_NAME,
}

var allowedConfigNamesForUpdate = map[string]string{
	HOME_TENANCY_ID_CONFIG_NAME: HOME_TENANCY_ID_CONFIG_NAME,
}

var allowedConfigNamesForDelete = map[string]string{
	HOME_TENANCY_ID_CONFIG_NAME: HOME_TENANCY_ID_CONFIG_NAME,
}

func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: BASE_CONFIG_PATH + framework.GenericNameRegex("configName"),
		Fields: map[string]*framework.FieldSchema{
			"configName": {
				Type:        framework.TypeString,
				Description: "Name of the config.",
			},
			"configValue": {
				Type:        framework.TypeString,
				Description: "Value of the config.",
			},
		},

		ExistenceCheck: b.pathConfigExistenceCheck,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathConfigCreate,
			logical.UpdateOperation: b.pathConfigUpdate,
			logical.DeleteOperation: b.pathConfigDelete,
			logical.ReadOperation:   b.pathConfigRead,
		},

		HelpSynopsis:    pathConfigSyn,
		HelpDescription: pathConfigDesc,
	}
}

func pathListConfigs(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: BASE_CONFIG_PATH + "?",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathConfigList,
		},

		HelpSynopsis:    pathListConfigsHelpSyn,
		HelpDescription: pathListConfigsHelpDesc,
	}
}

// Establishes dichotomy of request operation between CreateOperation and UpdateOperation.
// Returning 'true' forces an UpdateOperation, CreateOperation otherwise.
func (b *backend) pathConfigExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := b.lockedOCIConfig(ctx, req.Storage, data.Get("configName").(string))
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

// lockedOCIConfig returns the properties set on the given config. This method
// acquires the read lock before reading the config from the storage.
func (b *backend) lockedOCIConfig(ctx context.Context, s logical.Storage, configName string) (*OCIConfigEntry, error) {
	if strings.TrimSpace(configName) == "" {
		return nil, fmt.Errorf("missing configName")
	}

	b.configMutex.RLock()
	configEntry, err := b.nonLockedOCIConfig(ctx, s, configName)
	// we manually unlock rather than defer the unlock because we might need to grab
	// a read/write lock in the upgrade path
	b.configMutex.RUnlock()
	if err != nil {
		return nil, err
	}
	if configEntry == nil {
		return nil, nil
	}
	return configEntry, nil
}

// lockedSetOCIConfig creates or updates a config in the storage. This method
// acquires the write lock before creating or updating the config at the storage.
func (b *backend) lockedSetOCIConfig(ctx context.Context, s logical.Storage, configName string, configEntry *OCIConfigEntry) error {
	if strings.TrimSpace(configName) == "" {
		return fmt.Errorf("missing configName")
	}

	if configEntry == nil {
		return fmt.Errorf("config is not found")
	}

	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	return b.nonLockedSetOCIConfig(ctx, s, configName, configEntry)
}

// nonLockedSetOCIConfig creates or updates a config in the storage. This method
// does not acquire the write lock before writing the config to the storage. If
// locking is desired, use lockedSetOCIConfig instead.
func (b *backend) nonLockedSetOCIConfig(ctx context.Context, s logical.Storage, configName string,
	configEntry *OCIConfigEntry) error {
	if configName == "" {
		return fmt.Errorf("missing configName")
	}

	if configEntry == nil {
		return fmt.Errorf("config is not found")
	}

	entry, err := logical.StorageEntryJSON(BASE_CONFIG_PATH+configName, configEntry)
	if err != nil {
		return err
	}

	if err := s.Put(ctx, entry); err != nil {
		return err
	}

	return nil
}

// nonLockedOCIConfig returns the properties set on the given config. This method
// does not acquire the read lock before reading the config from the storage. If
// locking is desired, use lockedOCIConfig instead.
// This method also does NOT check to see if a config upgrade is required. It is
// the responsibility of the caller to check if a config upgrade is required and,
// if so, to upgrade the config
func (b *backend) nonLockedOCIConfig(ctx context.Context, s logical.Storage, configName string) (*OCIConfigEntry, error) {
	if strings.TrimSpace(configName) == "" {
		return nil, fmt.Errorf("missing configName")
	}

	entry, err := s.Get(ctx, BASE_CONFIG_PATH+configName)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result OCIConfigEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathConfigList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.RLock()
	defer b.configMutex.RUnlock()

	configs, err := req.Storage.List(ctx, BASE_CONFIG_PATH)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(configs), nil
}

func (b *backend) pathConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	configEntry, err := b.lockedOCIConfig(ctx, req.Storage, data.Get("configName").(string))
	if err != nil {
		return nil, err
	}
	if configEntry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: configEntry.ToResponseData(),
	}, nil
}

// create a Config
func (b *backend) pathConfigCreate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	configName := data.Get("configName").(string)
	if strings.TrimSpace(configName) == "" {
		return logical.ErrorResponse("missing configName"), nil
	}

	_, ok := allowedConfigNamesForCreate[configName]
	if ok == false {
		return logical.ErrorResponse(fmt.Sprintf("%s The specified configName is not allowed to be created.", configName)), nil
	}

	configValue := data.Get("configValue").(string)
	if strings.TrimSpace(configValue) == "" {
		return logical.ErrorResponse("missing configValue"), nil
	}

	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	configEntry, err := b.nonLockedOCIConfig(ctx, req.Storage, configName)
	if err != nil {
		return nil, err
	}

	if configEntry != nil {
		return logical.ErrorResponse("The specified config already exists"), nil
	}

	configEntry = &OCIConfigEntry{
		ConfigName:  configName,
		ConfigValue: configValue,
	}

	if err := b.nonLockedSetOCIConfig(ctx, req.Storage, configName, configEntry); err != nil {
		return nil, err
	}

	var resp logical.Response

	return &resp, nil
}

// update a Config
func (b *backend) pathConfigUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	configName := data.Get("configName").(string)
	if strings.TrimSpace(configName) == "" {
		return logical.ErrorResponse("missing configName"), nil
	}

	_, ok := allowedConfigNamesForUpdate[configName]
	if ok == false {
		return logical.ErrorResponse(fmt.Sprintf("%s The specified configName is not allowed to be updated.", configName)), nil
	}

	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	configEntry, err := b.nonLockedOCIConfig(ctx, req.Storage, configName)
	if err != nil {
		return nil, err
	}

	if configEntry == nil {
		return logical.ErrorResponse("The specified config does not exist"), nil
	}

	configValue := data.Get("configValue").(string)
	if strings.TrimSpace(configValue) == "" {
		return logical.ErrorResponse("missing configValue"), nil
	}

	configEntry.ConfigValue = configValue

	if err := b.nonLockedSetOCIConfig(ctx, req.Storage, configName, configEntry); err != nil {
		return nil, err
	}

	var resp logical.Response
	return &resp, nil
}

// delete a Config
func (b *backend) pathConfigDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	configName := data.Get("configName").(string)
	if strings.TrimSpace(configName) == "" {
		return logical.ErrorResponse("missing configName"), nil
	}

	_, ok := allowedConfigNamesForDelete[configName]
	if ok == false {
		return logical.ErrorResponse(fmt.Sprintf("%s The specified configName is not allowed to be deleted.", configName)), nil
	}

	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	return nil, req.Storage.Delete(ctx, BASE_CONFIG_PATH+configName)
}

// Struct to hold the information associated with an OCI config
type OCIConfigEntry struct {
	ConfigName  string `json:"configName" `
	ConfigValue string `json:"configValue" `
}

func (r *OCIConfigEntry) ToResponseData() map[string]interface{} {
	responseData := map[string]interface{}{
		"configName":  r.ConfigName,
		"configValue": r.ConfigValue,
	}

	return responseData
}

const pathConfigSyn = `
Create a config. Allowed values are:
homeTenancyId

Example:

vault write /auth/oci/homeTenancyId configValue=myocid
`

const pathConfigDesc = `
Create a config.
`

const pathListConfigsHelpSyn = `
Lists all the configs that are registered with Vault.
`

const pathListConfigsHelpDesc = `
Configs will be listed by their respective config names.
`
