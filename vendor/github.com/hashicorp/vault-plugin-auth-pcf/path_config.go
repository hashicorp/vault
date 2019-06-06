package pcf

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/hashicorp/vault-plugin-auth-pcf/models"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const configStorageKey = "config"

func (b *backend) pathConfig() *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"certificates": {
				Required:    true,
				Type:        framework.TypeStringSlice,
				Description: "The PEM-format CA certificates.",
			},
			"pcf_api_addr": {
				Required:     true,
				Type:         framework.TypeString,
				DisplayName:  "PCF API Address",
				DisplayValue: "https://api.10.244.0.34.xip.io",
				Description:  "PCF’s API address.",
			},
			"pcf_username": {
				Required:     true,
				Type:         framework.TypeString,
				DisplayName:  "PCF API Username",
				DisplayValue: "admin",
				Description:  "The username for PCF’s API.",
			},
			"pcf_password": {
				Required:         true,
				Type:             framework.TypeString,
				DisplayName:      "PCF API Password",
				DisplayValue:     "admin",
				Description:      "The password for PCF’s API.",
				DisplaySensitive: true,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.operationConfigCreateUpdate,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.operationConfigCreateUpdate,
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.operationConfigRead,
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.operationConfigDelete,
			},
		},
		HelpSynopsis:    pathConfigSyn,
		HelpDescription: pathConfigDesc,
	}
}

func (b *backend) operationConfigCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		// They're creating a config.
		certificates := data.Get("certificates").([]string)
		if len(certificates) == 0 {
			return logical.ErrorResponse("'certificates' is required"), nil
		}
		pcfApiAddr := data.Get("pcf_api_addr").(string)
		if pcfApiAddr == "" {
			return logical.ErrorResponse("'pcf_api_addr' is required"), nil
		}
		pcfUsername := data.Get("pcf_username").(string)
		if pcfUsername == "" {
			return logical.ErrorResponse("'pcf_username' is required"), nil
		}
		pcfPassword := data.Get("pcf_password").(string)
		if pcfPassword == "" {
			return logical.ErrorResponse("'pcf_password' is required"), nil
		}
		config, err = models.NewConfiguration(certificates, pcfApiAddr, pcfUsername, pcfPassword)
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
	} else {
		// They're updating a config. Only update the fields that have been sent in the call.
		if raw, ok := data.GetOk("certificates"); ok {
			switch v := raw.(type) {
			case []interface{}:
				certificates := make([]string, len(v))
				for _, certificateIfc := range v {
					certificate, ok := certificateIfc.(string)
					if !ok {
						continue
					}
					certificates = append(certificates, certificate)
				}
				config.Certificates = certificates
			case string:
				config.Certificates = []string{v}
			}
		}
		if raw, ok := data.GetOk("pcf_api_addr"); ok {
			config.PCFAPIAddr = raw.(string)
		}
		if raw, ok := data.GetOk("pcf_username"); ok {
			config.PCFUsername = raw.(string)
		}
		if raw, ok := data.GetOk("pcf_password"); ok {
			config.PCFPassword = raw.(string)
		}
	}

	// To give early and explicit feedback, make sure the config works by executing a test call
	// and checking that the API version is supported. If they don't have API v2 running, we would
	// probably expect a timeout of some sort below because it's first called in the NewClient
	// method.
	client, err := cfclient.NewClient(&cfclient.Config{
		ApiAddress: config.PCFAPIAddr,
		Username:   config.PCFUsername,
		Password:   config.PCFPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to establish an initial connection to the PCF API: %s", err)
	}
	info, err := client.GetInfo()
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(info.APIVersion, "2.") {
		return nil, fmt.Errorf("the PCF auth plugin only supports version 2.X.X of the PCF API")
	}

	entry, err := logical.StorageEntryJSON(configStorageKey, config)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *backend) operationConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"certificates": config.Certificates,
			"pcf_api_addr": config.PCFAPIAddr,
			"pcf_username": config.PCFUsername,
		},
	}, nil
}

func (b *backend) operationConfigDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(ctx, configStorageKey); err != nil {
		return nil, err
	}
	return nil, nil
}

// storedConfig may return nil without error if the user doesn't currently have a config.
func config(ctx context.Context, storage logical.Storage) (*models.Configuration, error) {
	entry, err := storage.Get(ctx, configStorageKey)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	configMap := make(map[string]interface{})
	if err := entry.DecodeJSON(&configMap); err != nil {
		return nil, err
	}
	var certificates []string
	certificatesIfc := configMap["certificates"].([]interface{})
	for _, certificateIfc := range certificatesIfc {
		certificates = append(certificates, certificateIfc.(string))
	}
	config, err := models.NewConfiguration(
		certificates,
		configMap["pcf_api_addr"].(string),
		configMap["pcf_username"].(string),
		configMap["pcf_password"].(string),
	)
	if err != nil {
		return nil, err
	}
	return config, nil
}

const pathConfigSyn = `
Provide Vault with the CA certificate used to issue all client certificates.
`

const pathConfigDesc = `
When a login is attempted using a PCF client certificate, Vault will verify
that the client certificate was issued by the CA certificate configured here.
Only those passing this check will be able to gain authorization.
`
