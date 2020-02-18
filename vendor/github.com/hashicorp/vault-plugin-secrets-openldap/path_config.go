package openldap

import (
	"context"
	"errors"

	"github.com/hashicorp/vault-plugin-secrets-openldap/client"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/ldaputil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	configPath            = "config"
	defaultPasswordLength = 64
	defaultTLSVersion     = "tls12"
)

func readConfig(ctx context.Context, storage logical.Storage) (*config, error) {
	entry, err := storage.Get(ctx, configPath)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	config := &config{}
	if err := entry.DecodeJSON(config); err != nil {
		return nil, err
	}
	return config, nil
}

func (b *backend) pathConfig() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: configPath,
			Fields:  b.configFields(),
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.configCreateUpdateOperation,
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.configCreateUpdateOperation,
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.configReadOperation,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.configDeleteOperation,
				},
			},
			HelpSynopsis:    configHelpSynopsis,
			HelpDescription: configHelpDescription,
		},
	}
}

func (b *backend) configFields() map[string]*framework.FieldSchema {
	fields := ldaputil.ConfigFields()
	fields["ttl"] = &framework.FieldSchema{
		Type:        framework.TypeDurationSecond,
		Description: "The default password time-to-live.",
	}
	fields["max_ttl"] = &framework.FieldSchema{
		Type:        framework.TypeDurationSecond,
		Description: "The maximum password time-to-live.",
	}
	fields["length"] = &framework.FieldSchema{
		Type:        framework.TypeInt,
		Default:     defaultPasswordLength,
		Description: "The desired length of passwords that Vault generates.",
	}
	return fields
}

func (b *backend) configCreateUpdateOperation(ctx context.Context, req *logical.Request, fieldData *framework.FieldData) (*logical.Response, error) {
	// Build and validate the ldap conf.
	ldapConf, err := ldaputil.NewConfigEntry(nil, fieldData)
	if err != nil {
		return nil, err
	}

	if err := ldapConf.Validate(); err != nil {
		return nil, err
	}

	length := fieldData.Get("length").(int)
	url := fieldData.Get("url").(string)

	if url == "" {
		return nil, errors.New("url is required")
	}

	config := &config{
		LDAP: &client.Config{
			ConfigEntry: ldapConf,
		},
		PasswordLength: length,
	}

	entry, err := logical.StorageEntryJSON(configPath, config)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	// Respond with a 204.
	return nil, nil
}

func (b *backend) configReadOperation(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	config, err := readConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}

	// "password" is intentionally not returned by this endpoint
	configMap := config.LDAP.Map()
	delete(configMap, "bindpass")
	configMap["length"] = config.PasswordLength

	resp := &logical.Response{
		Data: configMap,
	}
	return resp, nil
}

func (b *backend) configDeleteOperation(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(ctx, configPath); err != nil {
		return nil, err
	}
	return nil, nil
}

type config struct {
	LDAP           *client.Config
	PasswordLength int `json:"length"`
}

const configHelpSynopsis = `
Configure the OpenLDAP secret engine plugin.
`

const configHelpDescription = `
This path configures the OpenLDAP secret engine plugin. See the documentation for the plugin specified
for a full list of accepted connection details.
`
