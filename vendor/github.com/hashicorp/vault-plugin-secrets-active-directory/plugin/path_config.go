package plugin

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/vault-plugin-secrets-active-directory/plugin/util"
	"github.com/hashicorp/vault/helper/ldaputil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	configPath       = "config"
	configStorageKey = "config"

	// This length is arbitrarily chosen but should work for
	// most Active Directory minimum and maximum length settings.
	// A bit tongue-in-cheek since programmers love their base-2 exponents.
	defaultPasswordLength = 64

	defaultTLSVersion = "tls12"
)

func (b *backend) readConfig(ctx context.Context, storage logical.Storage) (*configuration, error) {
	entry, err := storage.Get(ctx, configStorageKey)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	config := &configuration{&passwordConf{}, &ldaputil.ConfigEntry{}}
	if err := entry.DecodeJSON(config); err != nil {
		return nil, err
	}
	return config, nil
}

func (b *backend) pathConfig() *framework.Path {
	return &framework.Path{
		Pattern: configPath,
		Fields:  b.configFields(),
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.configUpdateOperation,
			logical.ReadOperation:   b.configReadOperation,
			logical.DeleteOperation: b.configDeleteOperation,
		},
		HelpSynopsis:    configHelpSynopsis,
		HelpDescription: configHelpDescription,
	}
}

func (b *backend) configFields() map[string]*framework.FieldSchema {
	fields := ldaputil.ConfigFields()
	fields["ttl"] = &framework.FieldSchema{
		Type:        framework.TypeDurationSecond,
		Description: "In seconds, the default password time-to-live.",
	}
	fields["max_ttl"] = &framework.FieldSchema{
		Type:        framework.TypeDurationSecond,
		Description: "In seconds, the maximum password time-to-live.",
	}
	fields["length"] = &framework.FieldSchema{
		Type:        framework.TypeInt,
		Default:     defaultPasswordLength,
		Description: "The desired length of passwords that Vault generates.",
	}
	return fields
}

func (b *backend) configUpdateOperation(ctx context.Context, req *logical.Request, fieldData *framework.FieldData) (*logical.Response, error) {
	// Build and validate the ldap conf.
	activeDirectoryConf, err := ldaputil.NewConfigEntry(fieldData)
	if err != nil {
		return nil, err
	}
	if err := activeDirectoryConf.Validate(); err != nil {
		return nil, err
	}

	// Build the password conf.
	ttl := fieldData.Get("ttl").(int)
	maxTTL := fieldData.Get("max_ttl").(int)
	length := fieldData.Get("length").(int)

	if ttl == 0 {
		ttl = int(b.System().DefaultLeaseTTL().Seconds())
	}
	if maxTTL == 0 {
		maxTTL = int(b.System().MaxLeaseTTL().Seconds())
	}
	if ttl > maxTTL {
		return nil, errors.New("ttl must be smaller than or equal to max_ttl")
	}
	if ttl < 1 {
		return nil, errors.New("ttl must be positive")
	}
	if maxTTL < 1 {
		return nil, errors.New("max_ttl must be positive")
	}
	if length < util.MinimumPasswordLength {
		return nil, fmt.Errorf("minimum password length is %d for sufficient complexity to be secure, though Vault recommends a higher length", util.MinimumPasswordLength)
	}
	passwordConf := &passwordConf{
		TTL:    ttl,
		MaxTTL: maxTTL,
		Length: length,
	}

	config := &configuration{passwordConf, activeDirectoryConf}
	entry, err := logical.StorageEntryJSON(configStorageKey, config)
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
	config, err := b.readConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}

	// NOTE:
	// "password" is intentionally not returned by this endpoint,
	// as we lean away from returning sensitive information unless it's absolutely necessary.
	// Also, we don't return the full ADConf here because not all parameters are used by this engine.
	configMap := map[string]interface{}{
		"url":             config.ADConf.Url,
		"starttls":        config.ADConf.StartTLS,
		"insecure_tls":    config.ADConf.InsecureTLS,
		"certificate":     config.ADConf.Certificate,
		"binddn":          config.ADConf.BindDN,
		"userdn":          config.ADConf.UserDN,
		"upndomain":       config.ADConf.UPNDomain,
		"tls_min_version": config.ADConf.TLSMinVersion,
		"tls_max_version": config.ADConf.TLSMaxVersion,
	}
	for k, v := range config.PasswordConf.Map() {
		configMap[k] = v
	}

	resp := &logical.Response{
		Data: configMap,
	}
	return resp, nil
}

func (b *backend) configDeleteOperation(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(ctx, configStorageKey); err != nil {
		return nil, err
	}
	return nil, nil
}

const (
	configHelpSynopsis    = ``
	configHelpDescription = ``
)
