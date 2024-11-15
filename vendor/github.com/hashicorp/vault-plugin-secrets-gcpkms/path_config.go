// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpkms

import (
	"context"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// pathConfig defines the gcpkms/config base path on the backend.
func (b *backend) pathConfig() *framework.Path {
	return &framework.Path{
		Pattern: "config",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloudKMS,
		},

		HelpSynopsis: "Configure the GCP KMS secrets engine",
		HelpDescription: "Configure the GCP KMS secrets engine with credentials " +
			"or manage the requested scope(s).",

		Fields: map[string]*framework.FieldSchema{
			"credentials": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
The credentials to use for authenticating to Google Cloud. Leave this blank to
use the Default Application Credentials or instance metadata authentication.
`,
			},

			"scopes": &framework.FieldSchema{
				Type: framework.TypeCommaStringSlice,
				Description: `
The list of full-URL scopes to request when authenticating. By default, this
requests https://www.googleapis.com/auth/cloudkms.
`,
			},
		},

		ExistenceCheck: b.pathConfigExists,

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback: withFieldValidator(b.pathConfigWrite),
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "configure",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: withFieldValidator(b.pathConfigWrite),
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "configure",
				},
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: withFieldValidator(b.pathConfigRead),
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "read",
					OperationSuffix: "configuration",
				},
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: withFieldValidator(b.pathConfigDelete),
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "delete",
					OperationSuffix: "configuration",
				},
			},
		},
	}
}

// pathConfigExists checks if the configuration exists.
func (b *backend) pathConfigExists(ctx context.Context, req *logical.Request, _ *framework.FieldData) (bool, error) {
	entry, err := req.Storage.Get(ctx, "config")
	if err != nil {
		return false, errwrap.Wrapf("failed to get configuration from storage: {{err}}", err)
	}
	if entry == nil || len(entry.Value) == 0 {
		return false, nil
	}
	return true, nil
}

// pathConfigRead corresponds to READ gcpkms/config and is used to
// read the current configuration.
func (b *backend) pathConfigRead(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	c, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"scopes": c.Scopes,
		},
	}, nil
}

// pathConfigWrite corresponds to both CREATE and UPDATE gcpkms/config and is
// used to create or update the current configuration.
func (b *backend) pathConfigWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the current configuration, if it exists
	c, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	// Update the configuration
	changed, err := c.Update(d)
	if err != nil {
		return nil, logical.CodedError(400, err.Error())
	}

	// Only do the following if the config is different
	if changed {
		// Generate a new storage entry
		entry, err := logical.StorageEntryJSON("config", c)
		if err != nil {
			return nil, errwrap.Wrapf("failed to generate JSON configuration: {{err}}", err)
		}

		// Save the storage entry
		if err := req.Storage.Put(ctx, entry); err != nil {
			return nil, errwrap.Wrapf("failed to persist configuration to storage: {{err}}", err)
		}

		// Invalidate existing client so it reads the new configuration
		b.ResetClient()
	}

	return nil, nil
}

// pathConfigDelete corresponds to DELETE gcpkms/config and is used to delete
// all the configuration.
func (b *backend) pathConfigDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(ctx, "config"); err != nil {
		return nil, errwrap.Wrapf("failed to delete from storage: {{err}}", err)
	}

	// Invalidate existing client so it reads the new configuration
	b.ResetClient()

	return nil, nil
}
