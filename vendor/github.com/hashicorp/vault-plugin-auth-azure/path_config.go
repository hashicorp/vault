// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azureauth

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/pluginidentityutil"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathConfig(b *azureAuthBackend) *framework.Path {
	p := &framework.Path{
		Pattern: "config",
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAzure,
		},
		Fields: map[string]*framework.FieldSchema{
			"tenant_id": {
				Type:        framework.TypeString,
				Description: `The tenant id for the Azure Active Directory. This is sometimes referred to as Directory ID in AD. This value can also be provided with the AZURE_TENANT_ID environment variable.`,
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Tenant ID",
				},
			},
			"resource": {
				Type:        framework.TypeString,
				Description: `The resource URL for the vault application in Azure Active Directory. This value can also be provided with the AZURE_AD_RESOURCE environment variable.`,
			},
			"environment": {
				Type:        framework.TypeString,
				Description: `The Azure environment name. If not provided, AzurePublicCloud is used. This value can also be provided with the AZURE_ENVIRONMENT environment variable.`,
			},
			"client_id": {
				Type:        framework.TypeString,
				Description: `The OAuth2 client id to connection to Azure. This value can also be provided with the AZURE_CLIENT_ID environment variable.`,
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Client ID",
				},
			},
			"client_secret": {
				Type:        framework.TypeString,
				Description: `The OAuth2 client secret to connection to Azure. This value can also be provided with the AZURE_CLIENT_SECRET environment variable.`,
			},
			"root_password_ttl": {
				Type:        framework.TypeDurationSecond,
				Default:     defaultRootPasswordTTL,
				Description: "The TTL of the root password in Azure. This can be either a number of seconds or a time formatted duration (ex: 24h, 48ds)",
				Required:    false,
			},
			"max_retries": {
				Type:        framework.TypeInt,
				Default:     defaultMaxRetries,
				Description: "The maximum number of attempts a failed operation will be retried before producing an error.",
				Required:    false,
			},
			"max_retry_delay": {
				Type:        framework.TypeSignedDurationSecond,
				Default:     defaultMaxRetryDelay,
				Description: "The maximum delay allowed before retrying an operation.",
				Required:    false,
			},
			"retry_delay": {
				Type:        framework.TypeSignedDurationSecond,
				Default:     defaultRetryDelay,
				Description: "The initial amount of delay to use before retrying an operation, increasing exponentially.",
				Required:    false,
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
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathConfigWrite,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "auth",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigWrite,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "auth",
				},
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathConfigDelete,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "delete",
					OperationSuffix: "auth-configuration",
				},
			},
		},
		ExistenceCheck: b.pathConfigExistenceCheck,

		HelpSynopsis:    confHelpSyn,
		HelpDescription: confHelpDesc,
	}

	pluginidentityutil.AddPluginIdentityTokenFields(p.Fields)

	return p
}

type azureConfig struct {
	pluginidentityutil.PluginIdentityTokenParams

	TenantID                      string        `json:"tenant_id"`
	Resource                      string        `json:"resource"`
	Environment                   string        `json:"environment"`
	ClientID                      string        `json:"client_id"`
	ClientSecret                  string        `json:"client_secret"`
	ClientSecretKeyID             string        `json:"client_secret_key_id"`
	NewClientSecret               string        `json:"new_client_secret"`
	NewClientSecretCreated        time.Time     `json:"new_client_secret_created"`
	NewClientSecretExpirationDate time.Time     `json:"new_client_secret_expiration_date"`
	NewClientSecretKeyID          string        `json:"new_client_secret_key_id"`
	RootPasswordTTL               time.Duration `json:"root_password_ttl"`
	RootPasswordExpirationDate    time.Time     `json:"root_password_expiration_date"`
	MaxRetries                    int32         `json:"max_retries"`
	MaxRetryDelay                 time.Duration `json:"max_retry_delay"`
	RetryDelay                    time.Duration `json:"retry_delay"`
}

func (b *azureAuthBackend) config(ctx context.Context, s logical.Storage) (*azureConfig, error) {
	entry, err := s.Get(ctx, "config")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	config := new(azureConfig)
	if err := entry.DecodeJSON(config); err != nil {
		return nil, err
	}
	return config, nil
}

func (b *azureAuthBackend) pathConfigExistenceCheck(ctx context.Context, req *logical.Request, _ *framework.FieldData) (bool, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return false, err
	}
	return config != nil, nil
}

func (b *azureAuthBackend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		config = new(azureConfig)
	}

	tenantID, ok := data.GetOk("tenant_id")
	if ok {
		config.TenantID = tenantID.(string)
	}

	resource, ok := data.GetOk("resource")
	if ok {
		config.Resource = resource.(string)
	}

	environment, ok := data.GetOk("environment")
	if ok {
		config.Environment = environment.(string)
	}

	clientID, ok := data.GetOk("client_id")
	if ok {
		config.ClientID = clientID.(string)
	}

	clientSecret, ok := data.GetOk("client_secret")
	if ok {
		config.ClientSecret = clientSecret.(string)
	}

	config.RootPasswordTTL = defaultRootPasswordTTL
	rootExpirationRaw, ok := data.GetOk("root_password_ttl")
	if ok {
		config.RootPasswordTTL = time.Second * time.Duration(rootExpirationRaw.(int))
	}

	config.MaxRetries = defaultMaxRetries
	maxRetriesRaw, ok := data.GetOk("max_retries")
	if ok {
		config.MaxRetries = int32(maxRetriesRaw.(int))
	}

	config.MaxRetryDelay = defaultMaxRetryDelay
	maxRetryDelayRaw, ok := data.GetOk("max_retry_delay")
	if ok {
		config.MaxRetryDelay = time.Second * time.Duration(maxRetryDelayRaw.(int))
	}

	config.RetryDelay = defaultRetryDelay
	retryDelayRaw, ok := data.GetOk("retry_delay")
	if ok {
		config.RetryDelay = time.Second * time.Duration(retryDelayRaw.(int))
	}

	if err := config.ParsePluginIdentityTokenFields(data); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	if config.IdentityTokenAudience != "" && config.ClientSecret != "" {
		return logical.ErrorResponse("only one of 'client_secret' or 'identity_token_audience' can be set"), nil
	}

	// generate token to check if WIF is enabled on this edition of Vault
	if config.IdentityTokenAudience != "" {
		_, err := b.System().GenerateIdentityToken(ctx, &pluginutil.IdentityTokenRequest{
			Audience: config.IdentityTokenAudience,
		})
		if err != nil {
			if errors.Is(err, pluginidentityutil.ErrPluginWorkloadIdentityUnsupported) {
				return logical.ErrorResponse(err.Error()), nil
			}
			return nil, err
		}
	}

	// Create a settings object to validate all required settings
	// are available
	if _, err := b.getAzureSettings(ctx, config); err != nil {
		return nil, err
	}

	entry, err := logical.StorageEntryJSON("config", config)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	// Reset backend
	b.reset()

	return nil, nil
}

func (b *azureAuthBackend) pathConfigRead(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"tenant_id":         config.TenantID,
			"resource":          config.Resource,
			"environment":       config.Environment,
			"client_id":         config.ClientID,
			"root_password_ttl": int(config.RootPasswordTTL.Seconds()),
			"retry_delay":       config.RetryDelay,
			"max_retry_delay":   config.MaxRetryDelay,
			"max_retries":       config.MaxRetries,
		},
	}
	config.PopulatePluginIdentityTokenData(resp.Data)

	if !config.RootPasswordExpirationDate.IsZero() {
		resp.Data["root_password_expiration_date"] = config.RootPasswordExpirationDate
	}

	return resp, nil
}

func (b *azureAuthBackend) pathConfigDelete(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, "config")

	if err == nil {
		b.reset()
	}

	return nil, err
}

func (b *azureAuthBackend) saveConfig(ctx context.Context, config *azureConfig, s logical.Storage) error {
	entry, err := logical.StorageEntryJSON(configStoragePath, config)
	if err != nil {
		return err
	}

	err = s.Put(ctx, entry)
	if err != nil {
		return err
	}

	// reset the backend since the client and provider will have been
	// built using old versions of this data
	b.reset()

	return nil
}

const (
	// The default password expiration duration is 6 months in
	// the Azure UI, so we're setting it to 6 months (in hours)
	// as the default.
	defaultRootPasswordTTL = 4380 * time.Hour
	defaultRetryDelay      = 4 * time.Second
	defaultMaxRetries      = int32(3)
	defaultMaxRetryDelay   = 60 * time.Second
	configStoragePath      = "config"
	confHelpSyn            = `Configures the Azure authentication backend.`
	confHelpDesc           = `
The Azure authentication backend validates the login JWTs using the
configured credentials.  In order to validate machine information, the
OAuth2 client id and secret are used to query the Azure API.  The OAuth2
credentials require Microsoft.Compute/virtualMachines/read permission on
the resource requesting credentials.
`
)
