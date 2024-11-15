// Copyright © 2019, Oracle and/or its affiliates.
package ociauth

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"strings"
)

// These constants store the configuration keys
const (
	HomeTenancyIdConfigName = "home_tenancy_id"
)

func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixOCI,
		},

		Fields: map[string]*framework.FieldSchema{
			HomeTenancyIdConfigName: {
				Type:        framework.TypeString,
				Description: "The tenancy id of the account.",
			},
		},

		ExistenceCheck: b.pathConfigExistenceCheck,

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathConfigCreateUpdate,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "configure",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigCreateUpdate,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "configure",
				},
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathConfigDelete,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "configuration",
				},
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigRead,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "configuration",
				},
			},
		},

		HelpSynopsis:    pathConfigSyn,
		HelpDescription: pathConfigDesc,
	}
}

// Establishes dichotomy of request operation between CreateOperation and UpdateOperation.
// Returning 'true' forces an UpdateOperation, CreateOperation otherwise.
func (b *backend) pathConfigExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := b.getOCIConfig(ctx, req.Storage)
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

// setOCIConfig creates or updates a config in the storage.
func (b *backend) setOCIConfig(ctx context.Context, s logical.Storage, configEntry *OCIConfigEntry) error {
	if configEntry == nil {
		return fmt.Errorf("config is not found")
	}

	entry, err := logical.StorageEntryJSON("config", configEntry)
	if err != nil {
		return err
	}

	if err := s.Put(ctx, entry); err != nil {
		return err
	}

	return nil
}

// getOCIConfig returns the properties set on the given config.
// This method also does NOT check to see if a config upgrade is required. It is
// the responsibility of the caller to check if a config upgrade is required and,
// if so, to upgrade the config
func (b *backend) getOCIConfig(ctx context.Context, s logical.Storage) (*OCIConfigEntry, error) {
	entry, err := s.Get(ctx, "config")
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

func (b *backend) pathConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	configEntry, err := b.getOCIConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if configEntry == nil {
		return nil, nil
	}

	responseData := map[string]interface{}{
		HomeTenancyIdConfigName: configEntry.HomeTenancyId,
	}

	return &logical.Response{
		Data: responseData,
	}, nil
}

// Create a Config
func (b *backend) pathConfigCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	homeTenancyId := data.Get(HomeTenancyIdConfigName).(string)
	if strings.TrimSpace(homeTenancyId) == "" {
		return logical.ErrorResponse("Missing homeTenancyId"), nil
	}

	configEntry, err := b.getOCIConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if configEntry == nil && req.Operation == logical.UpdateOperation {
		return logical.ErrorResponse("The specified config does not exist"), nil
	}

	configEntry = &OCIConfigEntry{
		HomeTenancyId: homeTenancyId,
	}

	if err := b.setOCIConfig(ctx, req.Storage, configEntry); err != nil {
		return nil, err
	}

	var resp logical.Response

	return &resp, nil
}

// Delete a Config
func (b *backend) pathConfigDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return nil, req.Storage.Delete(ctx, "config")
}

// Struct to hold the information associated with an OCI config
type OCIConfigEntry struct {
	HomeTenancyId string `json:"home_tenancy_id" `
}

const pathConfigSyn = `
Manages the configuration for the Vault Auth Plugin.
`

const pathConfigDesc = `
The home_tenancy_id configuration is the Tenant OCID of your OCI Account. Only login requests from entities present in this tenant are accepted.

Example:

vault write /auth/oci/config home_tenancy_id=myocid
`
