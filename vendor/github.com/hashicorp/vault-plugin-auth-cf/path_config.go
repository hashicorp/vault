package cf

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault-plugin-auth-cf/models"
	"github.com/hashicorp/vault-plugin-auth-cf/util"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const configStorageKey = "config"

func (b *backend) pathConfig() *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"identity_ca_certificates": {
				Type: framework.TypeStringSlice,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "Identity CA Certificates",
					Value: `-----BEGIN CERTIFICATE----- ... -----END CERTIFICATE-----`,
				},
				Description: "The PEM-format CA certificates that are required to have issued the instance certificates presented for logging in.",
			},
			"cf_api_trusted_certificates": {
				Type: framework.TypeStringSlice,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "CF API Trusted IdentityCACertificates",
					Value: `-----BEGIN CERTIFICATE----- ... -----END CERTIFICATE-----`,
				},
				Description: "The PEM-format CA certificates that are acceptable for the CF API to present.",
			},
			"cf_api_addr": {
				Type: framework.TypeString,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "CF API Address",
					Value: "https://api.10.244.0.34.xip.io",
				},
				Description: "CF’s API address.",
			},
			"cf_username": {
				Type: framework.TypeString,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "CF API Username",
					Value: "admin",
				},
				Description: "The username for CF’s API.",
			},
			"cf_password": {
				Type: framework.TypeString,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:      "CF API Password",
					Sensitive: true,
				},
				Description: "The password for CF’s API.",
			},
			// These fields were in the original release, but are being deprecated because Cloud Foundry is moving
			// away from using "PCF" to refer to themselves.
			"pcf_api_trusted_certificates": {
				Deprecated: true,
				Type:       framework.TypeStringSlice,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "CF API Trusted IdentityCACertificates",
					Value: `-----BEGIN CERTIFICATE----- ... -----END CERTIFICATE-----`,
				},
				Description: `Deprecated. Please use "cf_api_trusted_certificates".`,
			},
			"pcf_api_addr": {
				Deprecated: true,
				Type:       framework.TypeString,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "CF API Address",
					Value: "https://api.10.244.0.34.xip.io",
				},
				Description: `Deprecated. Please use "cf_api_addr".`,
			},
			"pcf_username": {
				Deprecated: true,
				Type:       framework.TypeString,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "CF API Username",
					Value: "admin",
				},
				Description: `Deprecated. Please use "cf_username".`,
			},
			"pcf_password": {
				Deprecated: true,
				Type:       framework.TypeString,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:      "CF API Password",
					Sensitive: true,
				},
				Description: `Deprecated. Please use "cf_password".`,
			},
			"login_max_seconds_not_before": {
				Type: framework.TypeDurationSecond,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "Login Max Seconds Old",
					Value: "300",
				},
				Description: `Duration in seconds for the maximum acceptable age of a "signing_time". Useful for clock drift. 
Set low to reduce the opportunity for replay attacks.`,
				Default: 300,
			},
			"login_max_seconds_not_after": {
				Type: framework.TypeInt,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "Login Max Seconds Ahead",
					Value: "60",
				},
				Description: `Duration in seconds for the maximum acceptable length in the future a "signing_time" can be. Useful for clock drift.
Set low to reduce the opportunity for replay attacks.`,
				Default: 60,
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
		// All new configs will be created as config version 1.
		identityCACerts := data.Get("identity_ca_certificates").([]string)
		if len(identityCACerts) == 0 {
			return logical.ErrorResponse("'identity_ca_certificates' is required"), nil
		}

		cfApiAddrIfc, ok := data.GetFirst("cf_api_addr", "pcf_api_addr")
		if !ok {
			return logical.ErrorResponse("'cf_api_addr' is required"), nil
		}
		cfApiAddr := cfApiAddrIfc.(string)

		cfUsernameIfc, ok := data.GetFirst("cf_username", "pcf_username")
		if !ok {
			return logical.ErrorResponse("'cf_username' is required"), nil
		}
		cfUsername := cfUsernameIfc.(string)

		cfPasswordIfc, ok := data.GetFirst("cf_password", "pcf_password")
		if !ok {
			return logical.ErrorResponse("'cf_password' is required"), nil
		}
		cfPassword := cfPasswordIfc.(string)

		var cfApiCertificates []string
		cfApiCertificatesIfc, ok := data.GetFirst("cf_api_trusted_certificates", "pcf_api_trusted_certificates")
		if ok {
			cfApiCertificates = cfApiCertificatesIfc.([]string)
		}

		// Default this to 5 minutes.
		loginMaxSecNotBefore := 300 * time.Second
		if raw, ok := data.GetOk("login_max_seconds_not_before"); ok {
			loginMaxSecNotBefore = time.Duration(raw.(int)) * time.Second
		}

		// Default this to 1 minute.
		loginMaxSecNotAfter := 60 * time.Second
		if raw, ok := data.GetOk("login_max_seconds_not_after"); ok {
			loginMaxSecNotAfter = time.Duration(raw.(int)) * time.Second
		}
		config = &models.Configuration{
			Version:                1,
			IdentityCACertificates: identityCACerts,
			CFAPICertificates:      cfApiCertificates,
			CFAPIAddr:              cfApiAddr,
			CFUsername:             cfUsername,
			CFPassword:             cfPassword,
			LoginMaxSecNotBefore:   loginMaxSecNotBefore,
			LoginMaxSecNotAfter:    loginMaxSecNotAfter,
		}
	} else {
		// They're updating a config. Only update the fields that have been sent in the call.
		// The stored config will have already handled any version upgrades necessary on read,
		// so here we only need to deal with setting up the present version of the config.
		if raw, ok := data.GetOk("identity_ca_certificates"); ok {
			config.IdentityCACertificates = raw.([]string)
		}
		if raw, ok := data.GetFirst("cf_api_trusted_certificates", "pcf_api_trusted_certificates"); ok {
			config.CFAPICertificates = raw.([]string)
		}
		if raw, ok := data.GetFirst("cf_api_addr", "pcf_api_addr"); ok {
			config.CFAPIAddr = raw.(string)
		}
		if raw, ok := data.GetFirst("cf_username", "pcf_username"); ok {
			config.CFUsername = raw.(string)
		}
		if raw, ok := data.GetFirst("cf_password", "pcf_password"); ok {
			config.CFPassword = raw.(string)
		}
		if raw, ok := data.GetOk("login_max_seconds_not_before"); ok {
			config.LoginMaxSecNotBefore = time.Duration(raw.(int)) * time.Second
		}
		if raw, ok := data.GetOk("login_max_seconds_not_after"); ok {
			config.LoginMaxSecNotAfter = time.Duration(raw.(int)) * time.Second
		}
	}

	// To give early and explicit feedback, make sure the config works by executing a test call
	// and checking that the API version is supported. If they don't have API v2 running, we would
	// probably expect a timeout of some sort below because it's first called in the NewCFClient
	// method.
	client, err := util.NewCFClient(config)
	if err != nil {
		return nil, fmt.Errorf("unable to establish an initial connection to the CF API: %s", err)
	}
	info, err := client.GetInfo()
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(info.APIVersion, "2.") {
		return nil, fmt.Errorf("the CF auth plugin only supports version 2.X.X of the CF API")
	}

	if err := storeConfig(ctx, req.Storage, config); err != nil {
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
	resp := &logical.Response{
		Data: map[string]interface{}{
			"version":                      config.Version,
			"identity_ca_certificates":     config.IdentityCACertificates,
			"cf_api_trusted_certificates":  config.CFAPICertificates,
			"cf_api_addr":                  config.CFAPIAddr,
			"cf_username":                  config.CFUsername,
			"login_max_seconds_not_before": config.LoginMaxSecNotBefore / time.Second,
			"login_max_seconds_not_after":  config.LoginMaxSecNotAfter / time.Second,
		},
	}
	// Populate any deprecated values and warn about them. These should just be stripped when we go to
	// version 2 of the config.
	if len(config.PCFAPICertificates) > 0 {
		resp.Data["pcf_api_trusted_certificates"] = config.PCFAPICertificates
		resp.AddWarning(deprecationText("cf_api_trusted_certificates", "pcf_api_trusted_certificates"))
	}
	if config.PCFAPIAddr != "" {
		resp.Data["pcf_api_addr"] = config.PCFAPIAddr
		resp.AddWarning(deprecationText("cf_api_addr", "pcf_api_addr"))
	}
	if config.PCFUsername != "" {
		resp.Data["pcf_username"] = config.PCFUsername
		resp.AddWarning(deprecationText("cf_username", "pcf_username"))
	}
	return resp, nil
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
	config := &models.Configuration{}
	if err := entry.DecodeJSON(config); err != nil {
		return nil, err
	}

	// Perform config version migrations if needed.
	if config.Version == 0 {
		if config.CFAPIAddr == "" && config.PCFAPIAddr != "" {
			config.CFAPIAddr = config.PCFAPIAddr
		}
		if len(config.CFAPICertificates) == 0 && len(config.PCFAPICertificates) > 0 {
			config.CFAPICertificates = config.PCFAPICertificates
		}
		if config.CFUsername == "" && config.PCFUsername != "" {
			config.CFUsername = config.PCFUsername
		}
		if config.CFPassword == "" && config.PCFPassword != "" {
			config.CFPassword = config.PCFPassword
		}
		config.Version = 1
		if err := storeConfig(ctx, storage, config); err != nil {
			return nil, err
		}
	}
	return config, nil
}

func storeConfig(ctx context.Context, storage logical.Storage, conf *models.Configuration) error {
	entry, err := logical.StorageEntryJSON(configStorageKey, conf)
	if err != nil {
		return err
	}
	return storage.Put(ctx, entry)
}

func deprecationText(newParam, oldParam string) string {
	return fmt.Sprintf("Use %q instead. If this and %q are both specified, only %q will be used.", newParam, oldParam, newParam)
}

const pathConfigSyn = `
Provide Vault with the CA certificate used to issue all client certificates.
`

const pathConfigDesc = `
When a login is attempted using a CF client certificate, Vault will verify
that the client certificate was issued by the CA certificate configured here.
Only those passing this check will be able to gain authorization.
`
