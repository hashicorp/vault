// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cf

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/hashicorp/vault-plugin-auth-cf/models"
)

const configStorageKey = "config"

func (b *backend) pathConfig() *framework.Path {
	return &framework.Path{
		Pattern: "config",
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixCloudFoundry,
		},
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
			"cf_api_mutual_tls_certificate": {
				Type: framework.TypeString,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "CF API Mutual TLS Certificate",
					Value: `-----BEGIN CERTIFICATE----- ... -----END CERTIFICATE-----`,
				},
				Description: "The PEM-format certificates that are presented for mutual TLS with the CloudFoundry API. If not set, mutual TLS is not used",
			},
			"cf_api_mutual_tls_key": {
				Type: framework.TypeString,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "CF API Mutual TLS Key",
					Value: `-----BEGIN RSA PRIVATE KEY----- ... -----END RSA PRIVATE KEY-----`,
				},
				Description: "The PEM-format private key that are used for mutual TLS with the CloudFoundry API. If not set, mutual TLS is not used",
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
			"cf_client_id": {
				Type: framework.TypeString,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "CF API Client ID",
					Value: "client",
				},
				Description: "The client id for CF’s API.",
			},
			"cf_client_secret": {
				Type: framework.TypeString,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:      "CF API Client Secret",
					Sensitive: true,
				},
				Description: "The client secret for CF’s API.",
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
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.operationConfigWrite,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "configure",
				},
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.operationConfigRead,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "configuration",
				},
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.operationConfigDelete,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "configuration",
				},
			},
		},
		HelpSynopsis:    pathConfigSyn,
		HelpDescription: pathConfigDesc,
	}
}

func (b *backend) operationConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	config, err := getConfig(ctx, req.Storage)
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

		var cfUsername string
		cfUsernameIfc, ok := data.GetFirst("cf_username", "pcf_username")
		if ok {
			cfUsername = cfUsernameIfc.(string)
		}

		var cfPassword string
		cfPasswordIfc, ok := data.GetFirst("cf_password", "pcf_password")
		if ok {
			cfPassword = cfPasswordIfc.(string)
		}

		var cfClientId string
		cfClientIdIfc, ok := data.GetOk("cf_client_id")
		if ok {
			cfClientId = cfClientIdIfc.(string)
		}

		var cfClientSecret string
		cfClientSecretIfc, ok := data.GetOk("cf_client_secret")
		if ok {
			cfClientSecret = cfClientSecretIfc.(string)
		}

		// Before continuing, make sure that we have a pair of cf_username & cf_password,
		// pcf_username & pcf_password or cf_client_id & cf_client_secret
		// if none exist, then we should fail right away.
		if cfUsername == "" && cfClientId == "" {
			return logical.ErrorResponse("'cf_username' or 'cf_client_id' is required"), nil
		}

		if cfPassword == "" && cfClientSecret == "" {
			return logical.ErrorResponse("'cf_password' or 'cf_client_secret' is required"), nil
		}

		var cfApiCertificates []string
		cfApiCertificatesIfc, ok := data.GetFirst("cf_api_trusted_certificates", "pcf_api_trusted_certificates")
		if ok {
			cfApiCertificates = cfApiCertificatesIfc.([]string)
		}

		cfMTLSCertificate, ok := data.Get("cf_api_mutual_tls_certificate").(string)
		cfMTLSKey, ok := data.Get("cf_api_mutual_tls_key").(string)

		if (cfMTLSCertificate == "" && cfMTLSKey != "") ||
			(cfMTLSCertificate != "" && cfMTLSKey == "") {
			return logical.ErrorResponse("both 'cf_api_mutual_tls_certificate' and 'cf_api_mutual_tls_key' must be set if one is set"), nil
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
			CFMutualTLSCertificate: cfMTLSCertificate,
			CFMutualTLSKey:         cfMTLSKey,
			CFAPIAddr:              cfApiAddr,
			CFUsername:             cfUsername,
			CFPassword:             cfPassword,
			CFClientID:             cfClientId,
			CFClientSecret:         cfClientSecret,
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
		if raw, ok := data.GetOk("cf_api_mutual_tls_certificate"); ok {
			config.CFMutualTLSCertificate = raw.(string)
		}
		if raw, ok := data.GetOk("cf_api_mutual_tls_key"); ok {
			config.CFMutualTLSKey = raw.(string)
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
		if raw, ok := data.GetOk("cf_client_id"); ok {
			config.CFClientID = raw.(string)
		}
		if raw, ok := data.GetOk("cf_client_secret"); ok {
			config.CFClientSecret = raw.(string)
		}
	}

	if err := storeConfig(ctx, req.Storage, config); err != nil {
		return nil, err
	}

	// read the config back from storage to ensure that the client is updated with
	// the storage configuration
	config, err = getConfig(ctx, req.Storage)
	if err != nil {
		// should never get here
		return nil, err
	}

	if _, err := b.updateCFClient(ctx, config); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	return nil, nil
}

func (b *backend) operationConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	config, err := getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}
	resp := &logical.Response{
		Data: map[string]interface{}{
			"version":                       config.Version,
			"identity_ca_certificates":      config.IdentityCACertificates,
			"cf_api_trusted_certificates":   config.CFAPICertificates,
			"cf_api_mutual_tls_certificate": config.CFMutualTLSCertificate,
			"cf_api_addr":                   config.CFAPIAddr,
			"cf_username":                   config.CFUsername,
			"cf_client_id":                  config.CFClientID,
			"login_max_seconds_not_before":  config.LoginMaxSecNotBefore / time.Second,
			"login_max_seconds_not_after":   config.LoginMaxSecNotAfter / time.Second,
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
	b.mu.Lock()
	defer b.mu.Unlock()

	if err := req.Storage.Delete(ctx, configStorageKey); err != nil {
		return nil, err
	}
	return nil, nil
}

// getConfig returns the current configuration from storage.
func getConfig(ctx context.Context, storage logical.Storage) (*models.Configuration, error) {
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
