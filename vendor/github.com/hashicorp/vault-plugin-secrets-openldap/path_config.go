// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	defaultSchema         = client.SchemaOpenLDAP
	defaultTLSVersion     = "tls12"
	defaultCtxTimeout     = 1 * time.Minute
)

func (b *backend) pathConfig() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: configPath,
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixLDAP,
			},
			Fields: b.configFields(),
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.configCreateUpdateOperation,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "configure",
					},
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.configCreateUpdateOperation,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "configure",
					},
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.configReadOperation,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationSuffix: "configuration",
					},
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.configDeleteOperation,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationSuffix: "configuration",
					},
				},
			},
			ExistenceCheck: b.pathConfigExistenceCheck,
			HelpSynopsis:   "Configure the LDAP secrets engine plugin.",
			HelpDescription: "This path configures the LDAP secrets engine plugin. See the " +
				"documentation for the plugin for a full list of accepted parameters.",
		},
	}
}

func (b *backend) pathConfigExistenceCheck(ctx context.Context, req *logical.Request, _ *framework.FieldData) (bool, error) {
	entry, err := readConfig(ctx, req.Storage)
	if err != nil {
		return false, err
	}

	return entry != nil, nil
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
		Description: "The desired LDAP schema used when modifying user account passwords.",
	}
	fields["password_policy"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Description: "Password policy to use to generate passwords",
	}
	fields["skip_static_role_import_rotation"] = &framework.FieldSchema{
		Type:        framework.TypeBool,
		Description: "Whether to skip the 'import' rotation.",
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
	conf, err := readConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if conf == nil {
		conf = new(config)
		conf.LDAP = new(client.Config)
	}

	// Use the existing ldap client config if it is set
	var existing *ldaputil.ConfigEntry
	if conf.LDAP != nil && conf.LDAP.ConfigEntry != nil {
		existing = conf.LDAP.ConfigEntry
	}

	// Build and validate the ldap conf.
	ldapConf, err := ldaputil.NewConfigEntry(existing, fieldData)
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

	schema := fieldData.Get("schema").(string)
	_, schemaChanged := fieldData.Raw["schema"]

	// if update operation and schema not updated in raw payload, keep existing schema
	if existing != nil && !schemaChanged {
		schema = conf.LDAP.Schema
	}

	if schema == "" {
		return nil, errors.New("schema is required")
	}

	if !client.ValidSchema(schema) {
		return nil, fmt.Errorf("the configured schema %s is not valid. Supported schemas: %s",
			schema, client.SupportedSchemas())
	}

	// Set the userattr if given. Otherwise, set the default for creates.
	if userAttrRaw, ok := fieldData.GetOk("userattr"); ok {
		ldapConf.UserAttr = userAttrRaw.(string)
	} else if req.Operation == logical.CreateOperation {
		ldapConf.UserAttr = defaultUserAttr(schema)
	}

	passPolicy := fieldData.Get("password_policy").(string)
	_, passPolicyChanged := fieldData.Raw["password_policy"]

	// if update operation and password_policy not updated in raw payload, keep existing password_policy
	if existing != nil && !passPolicyChanged {
		passPolicy = conf.PasswordPolicy
	}

	if passPolicy != "" && hasPassLen {
		// If both a password policy and a password length are set, we can't figure out what to do
		return nil, fmt.Errorf("cannot set both 'password_policy' and 'length'")
	}

	staticSkip := fieldData.Get("skip_static_role_import_rotation").(bool)
	if _, set := fieldData.Raw["skip_static_role_import_rotation"]; existing != nil && !set {
		staticSkip = conf.SkipStaticRoleImportRotation // use existing value if not set
	}

	// Update config field values
	conf.PasswordPolicy = passPolicy
	conf.PasswordLength = passLength
	conf.SkipStaticRoleImportRotation = staticSkip
	conf.LDAP.ConfigEntry = ldapConf
	conf.LDAP.Schema = schema

	err = writeConfig(ctx, req.Storage, *conf)
	if err != nil {
		return nil, err
	}

	// Respond with a 204.
	return nil, nil
}

// defaultUserAttr returns the default user attribute for the given
// schema or an empty string if the schema is unknown.
func defaultUserAttr(schema string) string {
	switch schema {
	case client.SchemaAD:
		return "userPrincipalName"
	case client.SchemaRACF:
		return "racfid"
	case client.SchemaOpenLDAP:
		return "cn"
	default:
		return ""
	}
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
	configMap := config.LDAP.PasswordlessMap()
	if config.PasswordLength > 0 {
		configMap["length"] = config.PasswordLength
	}
	if config.PasswordPolicy != "" {
		configMap["password_policy"] = config.PasswordPolicy
	}
	configMap["skip_static_role_import_rotation"] = config.SkipStaticRoleImportRotation
	if !config.LDAP.LastBindPasswordRotation.IsZero() {
		configMap["last_bind_password_rotation"] = config.LDAP.LastBindPasswordRotation
	}
	if config.LDAP.Schema != "" {
		configMap["schema"] = config.LDAP.Schema
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
	LDAP                         *client.Config
	PasswordPolicy               string `json:"password_policy,omitempty"`
	SkipStaticRoleImportRotation bool   `json:"skip_static_role_import_rotation"`

	// Deprecated
	PasswordLength int `json:"length,omitempty"`
}
