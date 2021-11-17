package azuresecrets

import (
	"context"
	"errors"
	"time"

	"github.com/Azure/go-autorest/autorest/azure"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	configStoragePath = "config"
	// The default password expiration duration is 6 months in
	// the Azure UI, so we're setting it to 6 months (in hours)
	// as the default.
	defaultRootPasswordTTL = 4380
)

// azureConfig contains values to configure Azure clients and
// defaults for roles. The zero value is useful and results in
// environments variable and system defaults being used.
type azureConfig struct {
	SubscriptionID                string        `json:"subscription_id"`
	TenantID                      string        `json:"tenant_id"`
	ClientID                      string        `json:"client_id"`
	ClientSecret                  string        `json:"client_secret"`
	ClientSecretKeyID             string        `json:"client_secret_key_id"`
	NewClientSecret               string        `json:"new_client_secret"`
	NewClientSecretCreated        time.Time     `json:"new_client_secret_created"`
	NewClientSecretExpirationDate time.Time     `json:"new_client_secret_expiration_date"`
	NewClientSecretKeyID          string        `json:"new_client_secret_key_id"`
	Environment                   string        `json:"environment"`
	PasswordPolicy                string        `json:"password_policy"`
	UseMsGraphAPI                 bool          `json:"use_microsoft_graph_api"`
	RootPasswordTTL               time.Duration `json:"root_password_ttl"`
	RootPasswordExpirationDate    time.Time     `json:"root_password_expiration_date"`
}

func pathConfig(b *azureSecretBackend) *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"subscription_id": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The subscription id for the Azure Active Directory.
				This value can also be provided with the AZURE_SUBSCRIPTION_ID environment variable.`,
			},
			"tenant_id": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The tenant id for the Azure Active Directory. This value can also
				be provided with the AZURE_TENANT_ID environment variable.`,
			},
			"environment": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The Azure environment name. If not provided, AzurePublicCloud is used.
				This value can also be provided with the AZURE_ENVIRONMENT environment variable.`,
			},
			"client_id": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The OAuth2 client id to connect to Azure.
				This value can also be provided with the AZURE_CLIENT_ID environment variable.`,
			},
			"client_secret": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The OAuth2 client secret to connect to Azure.
				This value can also be provided with the AZURE_CLIENT_SECRET environment variable.`,
			},
			"password_policy": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the password policy to use to generate passwords for dynamic credentials.",
			},
			"use_microsoft_graph_api": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Description: "Enable usage of the Microsoft Graph API over the deprecated Azure AD Graph API.",
			},
			"root_password_ttl": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Default:     defaultRootPasswordTTL,
				Description: "The TTL of the root password in Azure. This can be either a number of seconds or a time formatted duration (ex: 24h, 48ds)",
				Required:    false,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigRead,
			},
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathConfigWrite,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigWrite,
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathConfigDelete,
			},
		},
		ExistenceCheck:  b.pathConfigExistenceCheck,
		HelpSynopsis:    confHelpSyn,
		HelpDescription: confHelpDesc,
	}
}

func (b *azureSecretBackend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var merr *multierror.Error

	config, err := b.getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if config == nil {
		if req.Operation == logical.UpdateOperation {
			return nil, errors.New("config not found during update operation")
		}
		config = new(azureConfig)
	}

	if subscriptionID, ok := data.GetOk("subscription_id"); ok {
		config.SubscriptionID = subscriptionID.(string)
	}

	if tenantID, ok := data.GetOk("tenant_id"); ok {
		config.TenantID = tenantID.(string)
	}

	if environment, ok := data.GetOk("environment"); ok {
		e := environment.(string)
		if _, err := azure.EnvironmentFromName(e); err != nil {
			merr = multierror.Append(merr, err)
		} else {
			config.Environment = e
		}
	}

	if clientID, ok := data.GetOk("client_id"); ok {
		config.ClientID = clientID.(string)
	}

	if clientSecret, ok := data.GetOk("client_secret"); ok {
		config.ClientSecret = clientSecret.(string)
	}

	if useMsGraphApi, ok := data.GetOk("use_microsoft_graph_api"); ok {
		config.UseMsGraphAPI = useMsGraphApi.(bool)
	}

	config.PasswordPolicy = data.Get("password_policy").(string)

	config.RootPasswordTTL = defaultRootPasswordTTL * time.Hour
	rootExpirationRaw, ok := data.GetOk("root_password_ttl")
	if ok {
		config.RootPasswordTTL = time.Second * time.Duration(rootExpirationRaw.(int))
	}

	if merr.ErrorOrNil() != nil {
		return logical.ErrorResponse(merr.Error()), nil
	}

	err = b.saveConfig(ctx, config, req.Storage)
	if err != nil {
		return nil, err
	}

	resp := addAADWarning(nil, config)

	return resp, nil
}

const aadWarning = "This configuration is using the Azure Active Directory API which is being " +
	"removed soon. Please migrate to using the Microsoft Graph API using the " +
	"use_microsoft_graph_api configuration parameter."

func addAADWarning(resp *logical.Response, config *azureConfig) *logical.Response {
	if config.UseMsGraphAPI {
		return resp
	}
	if resp == nil {
		resp = &logical.Response{}
	}
	resp.AddWarning(aadWarning)
	return resp
}

func (b *azureSecretBackend) pathConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.getConfig(ctx, req.Storage)

	if err != nil {
		return nil, err
	}

	if config == nil {
		config = new(azureConfig)
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"subscription_id":         config.SubscriptionID,
			"tenant_id":               config.TenantID,
			"environment":             config.Environment,
			"client_id":               config.ClientID,
			"use_microsoft_graph_api": config.UseMsGraphAPI,
			"root_password_ttl":       int(config.RootPasswordTTL.Seconds()),
		},
	}

	if !config.RootPasswordExpirationDate.IsZero() {
		resp.Data["root_password_expiration_date"] = config.RootPasswordExpirationDate
	}

	return addAADWarning(resp, config), nil
}

func (b *azureSecretBackend) pathConfigDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, configStoragePath)

	if err == nil {
		b.reset()
	}

	return nil, err
}

func (b *azureSecretBackend) pathConfigExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	config, err := b.getConfig(ctx, req.Storage)
	if err != nil {
		return false, err
	}

	return config != nil, err
}

func (b *azureSecretBackend) getConfig(ctx context.Context, s logical.Storage) (*azureConfig, error) {
	entry, err := s.Get(ctx, configStoragePath)
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

func (b *azureSecretBackend) saveConfig(ctx context.Context, config *azureConfig, s logical.Storage) error {
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

const confHelpSyn = `Configure the Azure Secret backend.`
const confHelpDesc = `
The Azure secret backend requires credentials for managing applications and
service principals. This endpoint is used to configure those credentials as
well as default values for the backend in general.
`
