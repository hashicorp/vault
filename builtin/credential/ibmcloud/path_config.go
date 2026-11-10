// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package ibmcloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/automatedrotationutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const configStoragePath = "config"

// ibmCloudConfig represents the configuration for IBM Cloud authentication
type ibmCloudConfig struct {
	// IAM endpoint URL
	IAMEndpoint string `json:"iam_endpoint"`

	// IAM client timeout in seconds
	IAMClientTimeout int `json:"iam_client_timeout"`

	// Admin API key for account verification and root credential rotation
	// This key is used to obtain an admin token for:
	// 1. Verifying service ID and user account access
	// 2. Rotating root credentials via config/rotate-root
	APIKey string `json:"api_key"`

	// Service ID that owns the API key (required when api_key is set)
	ServiceID string `json:"service_id"`

	// IBM Cloud account ID for authentication verification (required when api_key is set)
	// Used to verify that authenticating users/service IDs belong to this account
	AccountID string `json:"account_id"`

	// AutomatedRotationParams holds the schedule/period/window/disable fields
	// for automated root credential rotation.
	automatedrotationutil.AutomatedRotationParams
}

func (b *backend) pathConfigIBMCloud() *framework.Path {
	p := &framework.Path{
		Pattern: "config$",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixIBMCloud,
			OperationSuffix: "configuration",
		},

		Fields: map[string]*framework.FieldSchema{
			"iam_endpoint": {
				Type:        framework.TypeString,
				Description: "IBM Cloud IAM endpoint URL. Defaults to https://iam.cloud.ibm.com",
				Default:     defaultIAMEndpoint,
			},
			"iam_client_timeout": {
				Type:        framework.TypeInt,
				Description: "Timeout for IAM API calls in seconds. Defaults to 30",
				Default:     30,
			},
			"api_key": {
				Type:        framework.TypeString,
				Description: "Admin API key for account verification and root credential rotation. Required for account and access group verification.",
				DisplayAttrs: &framework.DisplayAttributes{
					Sensitive: true,
				},
			},
			"service_id": {
				Type:        framework.TypeString,
				Description: "Service ID that owns the API key. Optional; only required for root rotation (config/rotate-root).",
			},
			"account_id": {
				Type:        framework.TypeString,
				Description: "IBM Cloud account ID for authentication verification. Required when api_key is set.",
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigRead,
				Summary:  "Return the configured IBM Cloud settings.",
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "read",
				},
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"iam_endpoint": {
								Type:        framework.TypeString,
								Description: "IBM Cloud IAM endpoint URL.",
							},
							"iam_client_timeout": {
								Type:        framework.TypeInt,
								Description: "Timeout for IAM API calls in seconds.",
							},
							"service_id": {
								Type:        framework.TypeString,
								Description: "Service ID that owns the API key.",
							},
							"account_id": {
								Type:        framework.TypeString,
								Description: "IBM Cloud account ID for authentication verification.",
							},
						},
					}},
				},
			},
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathConfigWrite,
				Summary:  "Configure IBM Cloud IAM settings and credentials used for authentication.",
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "configure",
				},
				ForwardPerformanceSecondary: true,
				ForwardPerformanceStandby:   true,
				Responses: map[int][]framework.Response{
					http.StatusNoContent: {{Description: "No Content"}},
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigWrite,
				Summary:  "Configure IBM Cloud IAM settings and credentials used for authentication.",
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "configure",
				},
				ForwardPerformanceSecondary: true,
				ForwardPerformanceStandby:   true,
				Responses: map[int][]framework.Response{
					http.StatusNoContent: {{Description: "No Content"}},
				},
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathConfigDelete,
				Summary:  "Delete the IBM Cloud configuration.",
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "delete",
				},
				Responses: map[int][]framework.Response{
					http.StatusNoContent: {{Description: "No Content"}},
				},
			},
		},

		ExistenceCheck: b.pathConfigExistenceCheck,

		HelpSynopsis:    pathConfigHelpSyn,
		HelpDescription: pathConfigHelpDesc,
	}
	automatedrotationutil.AddAutomatedRotationFields(p.Fields)
	return p
}

// pathConfigExistenceCheck checks if configuration exists
func (b *backend) pathConfigExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return false, err
	}

	return config != nil, nil
}

func (b *backend) pathConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return nil, nil
	}

	respData := map[string]interface{}{
		"iam_endpoint":       config.IAMEndpoint,
		"iam_client_timeout": config.IAMClientTimeout,
		"service_id":         config.ServiceID,
		"account_id":         config.AccountID,
		// Never return the API key
	}
	config.AutomatedRotationParams.PopulateSetAutomatedRotationData(respData)

	return &logical.Response{Data: respData}, nil
}

func (b *backend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if config == nil {
		config = &ibmCloudConfig{}
	}

	// Update fields if provided
	if iamEndpoint, ok := data.GetOk("iam_endpoint"); ok {
		config.IAMEndpoint = iamEndpoint.(string)
	}

	if iamClientTimeout, ok := data.GetOk("iam_client_timeout"); ok {
		config.IAMClientTimeout = iamClientTimeout.(int)
	}

	if apiKey, ok := data.GetOk("api_key"); ok {
		config.APIKey = apiKey.(string)
	}

	if serviceID, ok := data.GetOk("service_id"); ok {
		config.ServiceID = serviceID.(string)
	}

	if accountID, ok := data.GetOk("account_id"); ok {
		config.AccountID = accountID.(string)
	}

	// Validate that if api_key is set, account_id must also be set. service_id
	// is optional: it is only needed for root rotation, which checks for it.
	if config.APIKey != "" && config.AccountID == "" {
		return logical.ErrorResponse("account_id is required when api_key is set"), nil
	}

	// Validate that if account_id is set, api_key must also be set
	if config.AccountID != "" && config.APIKey == "" {
		return logical.ErrorResponse("api_key is required when account_id is set"), nil
	}

	// Set defaults if not provided
	if config.IAMEndpoint == "" {
		config.IAMEndpoint = defaultIAMEndpoint
	}
	if config.IAMClientTimeout == 0 {
		config.IAMClientTimeout = 30
	}

	// Handle automated rotation fields. HandleRotationJob calls
	// ParseAutomatedRotationFields internally and decides whether to
	// register or deregister based on the parsed values; calling it
	// unconditionally is correct and avoids a redundant pre-parse.
	rotationOpResp, err := config.HandleRotationJob(ctx, b.Backend, data, req)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	entry, err := logical.StorageEntryJSON(configStoragePath, config)
	if err != nil {
		return nil, rotationOpResp.HandleStorageErrorAfterRotationJob(req, err)
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, rotationOpResp.HandleStorageErrorAfterRotationJob(req, err)
	}

	// Invalidate the cached client
	b.invalidate(ctx, configStoragePath)

	return nil, nil
}

func (b *backend) pathConfigDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(ctx, configStoragePath); err != nil {
		return nil, err
	}

	// Invalidate the cached client
	b.invalidate(ctx, configStoragePath)

	return nil, nil
}

// config reads the backend configuration from storage
func (b *backend) config(ctx context.Context, s logical.Storage) (*ibmCloudConfig, error) {
	entry, err := s.Get(ctx, configStoragePath)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	config := &ibmCloudConfig{}
	if err := entry.DecodeJSON(config); err != nil {
		return nil, fmt.Errorf("failed to decode configuration: %w", err)
	}

	return config, nil
}

const pathConfigHelpSyn = `
Configure the IBM Cloud authentication backend.
`

const pathConfigHelpDesc = `
The IBM Cloud authentication backend requires configuration of IBM Cloud
endpoints and optional credentials for root rotation.

Configuration Parameters:

  * iam_endpoint (string, optional)
      IBM Cloud IAM endpoint URL. Defaults to https://iam.cloud.ibm.com

  * iam_client_timeout (int, optional)
      Timeout for IAM API calls in seconds. Defaults to 30

  * api_key (string, optional)
      Admin API key for account verification and root credential rotation.
      Required for proper account and access group verification during authentication.
      When set, account_id must also be provided.

  * service_id (string, optional)
      Service ID that owns the API key. Only required for root credential
      rotation (config/rotate-root); not needed for authentication.

  * account_id (string, optional)
      IBM Cloud account ID for authentication verification. Required when api_key is set,
      and api_key is required when account_id is set.
      Used to verify that authenticating users/service IDs belong to this account.

Example (minimal configuration):

  $ vault write auth/ibmcloud/config \
      iam_endpoint="https://iam.cloud.ibm.com"

Example (with admin credentials for verification):

  $ vault write auth/ibmcloud/config \
      iam_endpoint="https://iam.cloud.ibm.com" \
      api_key="<admin-api-key>" \
      account_id="abc123def456"

Example (also enabling root rotation):

  $ vault write auth/ibmcloud/config \
      iam_endpoint="https://iam.cloud.ibm.com" \
      api_key="<admin-api-key>" \
      service_id="ServiceId-xxx" \
      account_id="abc123def456"
`

// Made with Bob
