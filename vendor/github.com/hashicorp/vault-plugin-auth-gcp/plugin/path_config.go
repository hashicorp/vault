// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpauth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/authmetadata"
	"github.com/hashicorp/vault/sdk/helper/pluginidentityutil"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
)

var (
	// The default gce_alias is "instance_id". The default fields
	// below are selected because they're unlikely to change often
	// for a particular instance ID.
	gceAuthMetadataFields = &authmetadata.Fields{
		FieldName: "gce_metadata",
		Default: []string{
			"instance_creation_timestamp",
			"instance_id",
			"instance_name",
			"project_id",
			"project_number",
			"role",
			"service_account_id",
			"service_account_email",
			"zone",
		},
		AvailableToAdd: []string{},
	}

	// The default iam_alias is "unique_id". The default fields
	// below are selected because they're unlikely to change often
	// for a particular instance ID.
	iamAuthMetadataFields = &authmetadata.Fields{
		FieldName: "iam_metadata",
		Default: []string{
			"project_id",
			"role",
			"service_account_id",
			"service_account_email",
		},
		AvailableToAdd: []string{},
	}
)

func pathConfig(b *GcpAuthBackend) *framework.Path {
	p := &framework.Path{
		Pattern: "config",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloud,
		},

		Fields: map[string]*framework.FieldSchema{
			"credentials": {
				Type: framework.TypeString,
				Description: `
Google credentials JSON that Vault will use to verify users against GCP APIs.
If not specified, will use application default credentials`,
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Credentials",
				},
			},
			"iam_alias": {
				Type:        framework.TypeString,
				Default:     defaultIAMAlias,
				Description: "Indicates what value to use when generating an alias for IAM authentications.",
			},
			iamAuthMetadataFields.FieldName: authmetadata.FieldSchema(iamAuthMetadataFields),
			"gce_alias": {
				Type:        framework.TypeString,
				Default:     defaultGCEAlias,
				Description: "Indicates what value to use when generating an alias for GCE authentications.",
			},
			gceAuthMetadataFields.FieldName: authmetadata.FieldSchema(gceAuthMetadataFields),
			"custom_endpoint": {
				Type:        framework.TypeKVPairs,
				Description: `Specifies overrides for various Google API Service Endpoints used in requests.`,
			},
			// Deprecated
			"google_certs_endpoint": {
				Type: framework.TypeString,
				Description: `
Deprecated. This field does nothing and be removed in a future release`,
				Deprecated: true,
			},
			"service_account_email": {
				Type:        framework.TypeString,
				Description: `Email ID for the Service Account to impersonate for Workload Identity Federation.`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigRead,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "read",
					OperationSuffix: "auth-configuration",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigWrite,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "auth",
				},
			},
		},

		HelpSynopsis: `Configure credentials used to query the GCP IAM API to verify authenticating service accounts`,
		HelpDescription: `
The GCP IAM auth backend makes queries to the GCP IAM auth backend to verify a service account
attempting login. It verifies the service account exists and retrieves a public key to verify
signed JWT requests passed in on login. The credentials should have the following permissions:

iam AUTH:
* iam.serviceAccountKeys.get
`,
	}

	pluginidentityutil.AddPluginIdentityTokenFields(p.Fields)

	return p
}

func (b *GcpAuthBackend) pathConfigWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if err := validateFields(req, d); err != nil {
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	c, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if err := c.Update(d); err != nil {
		return nil, logical.CodedError(http.StatusBadRequest, err.Error())
	}

	// generate token to check if WIF is enabled on this edition of Vault
	if c.IdentityTokenAudience != "" {
		_, err := b.System().GenerateIdentityToken(ctx, &pluginutil.IdentityTokenRequest{
			Audience: c.IdentityTokenAudience,
		})
		if err != nil {
			if errors.Is(err, pluginidentityutil.ErrPluginWorkloadIdentityUnsupported) {
				return logical.ErrorResponse(err.Error()), nil
			}
			return nil, err
		}
	}

	// Create/update the storage entry
	entry, err := logical.StorageEntryJSON("config", c)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JSON configuration: %w", err)
	}

	// Save the storage entry
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, fmt.Errorf("failed to persist configuration to storage: %w", err)
	}

	// Invalidate existing client so it reads the new configuration
	b.ClearCaches()

	return nil, nil
}

func (b *GcpAuthBackend) pathConfigRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if err := validateFields(req, d); err != nil {
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		gceAuthMetadataFields.FieldName: config.GCEAuthMetadata.AuthMetadata(),
		iamAuthMetadataFields.FieldName: config.IAMAuthMetadata.AuthMetadata(),
	}

	if config.Credentials != nil {
		if v := config.Credentials.ClientEmail; v != "" {
			resp["client_email"] = v
		}
		if v := config.Credentials.ClientId; v != "" {
			resp["client_id"] = v
		}
		if v := config.Credentials.PrivateKeyId; v != "" {
			resp["private_key_id"] = v
		}
		if v := config.Credentials.ProjectId; v != "" {
			resp["project_id"] = v
		}
	}

	if v := config.IAMAliasType; v != "" {
		resp["iam_alias"] = v
	}
	if v := config.GCEAliasType; v != "" {
		resp["gce_alias"] = v
	}

	endpoints := make(map[string]string)
	if v := config.APICustomEndpoint; v != "" {
		endpoints["api"] = v
	}
	if v := config.IAMCustomEndpoint; v != "" {
		endpoints["iam"] = v
	}
	if v := config.CRMCustomEndpoint; v != "" {
		endpoints["crm"] = v
	}
	if v := config.ComputeCustomEndpoint; v != "" {
		endpoints["compute"] = v
	}
	if len(endpoints) > 0 {
		resp["custom_endpoint"] = endpoints
	}

	if v := config.ServiceAccountEmail; v != "" {
		resp["service_account_email"] = v
	}

	config.PopulatePluginIdentityTokenData(resp)

	return &logical.Response{
		Data: resp,
	}, nil
}

// config reads the backend's gcpConfig from storage.
// This assumes the caller has already obtained the backend's config lock.
func (b *GcpAuthBackend) config(ctx context.Context, s logical.Storage) (*gcpConfig, error) {
	config := &gcpConfig{
		GCEAuthMetadata: authmetadata.NewHandler(gceAuthMetadataFields),
		IAMAuthMetadata: authmetadata.NewHandler(iamAuthMetadataFields),
	}
	entry, err := s.Get(ctx, "config")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return config, nil
	}
	if err := entry.DecodeJSON(config); err != nil {
		return nil, err
	}
	return config, nil
}
