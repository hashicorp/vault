package ad

import (
	"context"

	"github.com/fatih/structs"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/activedirectory"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	// This length is arbitrarily chosen but should work for
	// most Active Directory minimum and maximum length settings.
	// A bit of tongue-in-cheek since programmers love their base-2 exponents.
	defaultPasswordLength = 64

	// The number of minutes in 32 days.
	defaultPasswordTTLs = 24 * 60 * 32
)

type configHandler struct {
	logger hclog.Logger
}

func (h *configHandler) Path() *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"username": {
				Type:        framework.TypeString,
				Default:     "",
				Description: "Username with sufficient permissions in Active Directory to administer passwords.",
			},

			"password": {
				Type:        framework.TypeString,
				Default:     "",
				Description: "Password for username with sufficient permissions in Active Directory to administer passwords.",
			},

			"url": {
				Type:        framework.TypeString,
				Default:     "ldap://127.0.0.1",
				Description: "LDAP URL to connect to (default: ldap://127.0.0.1). Multiple URLs can be specified by concatenating them with commas; they will be tried in-order.",
			},

			"certificate": {
				Type:        framework.TypeString,
				Default:     "",
				Description: "CA certificate to use when verifying LDAP server certificate, must be x509 PEM encoded (optional)",
			},

			"dn": {
				Type:        framework.TypeString,
				Default:     "",
				Description: "The root distinguished name to bind to when managing service accounts",
			},

			"insecure_tls": {
				Type:        framework.TypeBool,
				Default:     false,
				Description: "Skip LDAP server SSL Certificate verification - VERY insecure (optional)",
			},

			"starttls": {
				Type:        framework.TypeBool,
				Default:     true,
				Description: "Issue a StartTLS command after establishing unencrypted connection (optional)",
			},

			"tls_min_version": {
				Type:        framework.TypeString,
				Default:     "tls12",
				Description: "Minimum TLS version to use. Accepted values are 'tls10', 'tls11' or 'tls12'. Defaults to 'tls12'",
			},

			"tls_max_version": {
				Type:        framework.TypeString,
				Default:     "tls12",
				Description: "Maximum TLS version to use. Accepted values are 'tls10', 'tls11' or 'tls12'. Defaults to 'tls12'",
			},

			"default_password_ttl": {
				Type:        framework.TypeInt,
				Default:     defaultPasswordTTLs,
				Description: "In minutes, the default password time-to-live.",
			},

			"max_password_ttl": {
				Type:        framework.TypeInt,
				Default:     defaultPasswordTTLs,
				Description: "In minutes, the maximum password time-to-live.",
			},

			"password_length": {
				Type:        framework.TypeInt,
				Default:     defaultPasswordLength,
				Description: "The desired length of passwords that Vault generates.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: h.DeleteOperation,
			logical.ReadOperation:   h.ReadOperation,
			logical.UpdateOperation: h.UpdateOperation,
		},
	}
}

// Config is a convenience method for other operations to use in retrieving the current config.
func (h *configHandler) Config(ctx context.Context, req *logical.Request) (*engineConfig, error) {

	entry, err := req.Storage.Get(ctx, "config")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return &engineConfig{}, nil
	}

	engineConf := &engineConfig{}
	if err := entry.DecodeJSON(engineConf); err != nil {
		return nil, err
	}

	return engineConf, nil
}

func (h *configHandler) DeleteOperation(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(ctx, "config"); err != nil {
		return nil, err
	}
	return nil, nil
}

func (h *configHandler) ReadOperation(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {

	engineConf, err := h.Config(ctx, req)
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Data: engineConf.Map(),
	}
	resp.AddWarning("Read access to this endpoint should be controlled via ACLs as it will return the configuration information as-is, including any passwords.")
	return resp, nil
}

func (h *configHandler) UpdateOperation(ctx context.Context, req *logical.Request, fieldData *framework.FieldData) (*logical.Response, error) {

	activeDirectoryConf, err := activedirectory.NewConfiguration(h.logger, fieldData)
	if err != nil {
		return logical.ErrorResponse(err.Error()), err
	}
	passwordConf := newPasswordConfig(fieldData)
	engineConf := &engineConfig{passwordConf, activeDirectoryConf}

	entry, err := logical.StorageEntryJSON("config", engineConf)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Data: engineConf.Map(),
	}
	resp.AddWarning("Write access to this endpoint should be controlled via ACLs as it will return the configuration information as-is, including any passwords.")
	return resp, nil
}

type engineConfig struct {
	PasswordConfig *passwordConfig
	ADConfig       *activedirectory.Configuration
}

func (c *engineConfig) IsSet() bool {
	if c.PasswordConfig != nil {
		return true
	}
	if c.ADConfig != nil {
		return true
	}
	return false
}

func (c *engineConfig) Map() map[string]interface{} {

	combined := make(map[string]interface{})
	if !c.IsSet() {
		return combined
	}

	for k, v := range structs.New(c.PasswordConfig).Map() {
		combined[k] = v
	}
	for k, v := range structs.New(c.ADConfig).Map() {
		combined[k] = v
	}
	return combined
}

func newPasswordConfig(fieldData *framework.FieldData) *passwordConfig {
	return &passwordConfig{
		DefaultPasswordTTL: fieldData.Get("default_password_ttl").(int),
		MaxPasswordTTL:     fieldData.Get("max_password_ttl").(int),
		PasswordLength:     fieldData.Get("password_length").(int),
	}
}

type passwordConfig struct {
	DefaultPasswordTTL int `json:"default_password_ttl" structs:"default_password_ttl" mapstructure:"default_password_ttl"`
	MaxPasswordTTL     int `json:"max_password_ttl" structs:"max_password_ttl" mapstructure:"max_password_ttl"`
	PasswordLength     int `json:"password_length" structs:"password_length" mapstructure:"password_length"`
}
