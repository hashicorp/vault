package openldap

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/vault-plugin-secrets-openldap/client"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/ldaputil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	configPath            = "config"
	defaultPasswordLength = 64
	defaultSchema         = "openldap"
	defaultTLSVersion     = "tls12"
	defaultCtxTimeout     = 1 * time.Minute
)

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
			HelpSynopsis: "Configure the OpenLDAP secret engine plugin.",
			HelpDescription: "This path configures the OpenLDAP secret engine plugin. See the documentation for the " +
				"plugin specified for a full list of accepted connection details.",
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
	fields["schema"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Default:     defaultSchema,
		Description: "The desired OpenLDAP schema used when modifying user account passwords.",
	}
	fields["password_policy"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Description: "Password policy to use to generate passwords",
	}

	// Deprecated
	fields["length"] = &framework.FieldSchema{
		Type:        framework.TypeInt,
		Description: "The desired length of passwords that Vault generates.",
		Deprecated:  true,
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

	rawPassLength, hasPassLen := fieldData.GetOk("length")
	if rawPassLength == nil {
		rawPassLength = 0 // Don't set to the default but keep this as the zero value so we know it hasn't been set
	}
	passLength := rawPassLength.(int)
	url := fieldData.Get("url").(string)

	if url == "" {
		return nil, errors.New("url is required")
	}

	schema := fieldData.Get("schema").(string)
	if schema == "" {
		return nil, errors.New("schema is required")
	}

	if !client.ValidSchema(schema) {
		return nil, fmt.Errorf("the configured schema %s is not valid.  Supported schemas: %s",
			schema, client.SupportedSchemas())
	}

	passPolicy := fieldData.Get("password_policy").(string)

	if passPolicy != "" && hasPassLen {
		// If both a password policy and a password length are set, we can't figure out what to do
		return nil, fmt.Errorf("cannot set both 'password_policy' and 'length'")
	}

	config := config{
		LDAP: &client.Config{
			ConfigEntry: ldapConf,
			Schema:      schema,
		},
		PasswordPolicy: passPolicy,
		PasswordLength: passLength,
	}

	err = writeConfig(ctx, req.Storage, config)
	if err != nil {
		return nil, err
	}

	// Respond with a 204.
	return nil, nil
}

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

func writeConfig(ctx context.Context, storage logical.Storage, config config) (err error) {
	entry, err := logical.StorageEntryJSON(configPath, config)
	if err != nil {
		return err
	}
	err = storage.Put(ctx, entry)
	if err != nil {
		return err
	}
	return nil
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
	if config.PasswordLength > 0 {
		configMap["length"] = config.PasswordLength
	}
	if config.PasswordPolicy != "" {
		configMap["password_policy"] = config.PasswordPolicy
	}

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
	PasswordPolicy string `json:"password_policy,omitempty"`

	// Deprecated
	PasswordLength int `json:"length,omitempty"`
}
