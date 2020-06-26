package config

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
)

const (
	// XXX Change use to api.EnvVaultSkipVerify once we've updated vendored
	// vault to version 1.1.0 or newer.
	EnvVaultSkipVerify = "VAULT_SKIP_VERIFY"

	// DefaultVaultRenewToken is the default value for if the Vault token should
	// be renewed.
	DefaultVaultRenewToken = true

	// DefaultVaultUnwrapToken is the default value for if the Vault token should
	// be unwrapped.
	DefaultVaultUnwrapToken = false

	// DefaultVaultRetryBase is the default value for the base time to use for
	// exponential backoff.
	DefaultVaultRetryBase = 250 * time.Millisecond

	// DefaultVaultRetryMaxAttempts is the default maximum number of attempts to
	// retry before quitting.
	DefaultVaultRetryMaxAttempts = 5
)

// VaultConfig is the configuration for connecting to a vault server.
type VaultConfig struct {
	// Address is the URI to the Vault server.
	Address *string `mapstructure:"address"`

	// Enabled controls whether the Vault integration is active.
	Enabled *bool `mapstructure:"enabled"`

	// Namespace is the Vault namespace to use for reading/writing secrets. This can
	// also be set via the VAULT_NAMESPACE environment variable.
	Namespace *string `mapstructure:"namespace"`

	// RenewToken renews the Vault token.
	RenewToken *bool `mapstructure:"renew_token"`

	// Retry is the configuration for specifying how to behave on failure.
	Retry *RetryConfig `mapstructure:"retry"`

	// SSL indicates we should use a secure connection while talking to Vault.
	SSL *SSLConfig `mapstructure:"ssl"`

	// Token is the Vault token to communicate with for requests. It may be
	// a wrapped token or a real token. This can also be set via the VAULT_TOKEN
	// environment variable, or via the VaultAgentTokenFile.
	Token *string `mapstructure:"token" json:"-"`

	// VaultAgentTokenFile is the path of file that contains a Vault Agent token.
	// If vault_agent_token_file is specified:
	//   - Consul Template will not try to renew the Vault token.
	//   - Consul Template will periodically stat the file and update the token if it has
	// changed.
	VaultAgentTokenFile *string `mapstructure:"vault_agent_token_file" json:"-"`

	// Transport configures the low-level network connection details.
	Transport *TransportConfig `mapstructure:"transport"`

	// UnwrapToken unwraps the provided Vault token as a wrapped token.
	UnwrapToken *bool `mapstructure:"unwrap_token"`
}

// DefaultVaultConfig returns a configuration that is populated with the
// default values.
func DefaultVaultConfig() *VaultConfig {
	v := &VaultConfig{
		Retry:     DefaultRetryConfig(),
		SSL:       DefaultSSLConfig(),
		Transport: DefaultTransportConfig(),
	}

	// Force SSL when communicating with Vault.
	v.SSL.Enabled = Bool(true)

	return v
}

// Copy returns a deep copy of this configuration.
func (c *VaultConfig) Copy() *VaultConfig {
	if c == nil {
		return nil
	}

	var o VaultConfig
	o.Address = c.Address

	o.Enabled = c.Enabled

	o.Namespace = c.Namespace

	o.RenewToken = c.RenewToken

	if c.Retry != nil {
		o.Retry = c.Retry.Copy()
	}

	if c.SSL != nil {
		o.SSL = c.SSL.Copy()
	}

	o.Token = c.Token

	o.VaultAgentTokenFile = c.VaultAgentTokenFile

	if c.Transport != nil {
		o.Transport = c.Transport.Copy()
	}

	o.UnwrapToken = c.UnwrapToken

	return &o
}

// Merge combines all values in this configuration with the values in the other
// configuration, with values in the other configuration taking precedence.
// Maps and slices are merged, most other values are overwritten. Complex
// structs define their own merge functionality.
func (c *VaultConfig) Merge(o *VaultConfig) *VaultConfig {
	if c == nil {
		if o == nil {
			return nil
		}
		return o.Copy()
	}

	if o == nil {
		return c.Copy()
	}

	r := c.Copy()

	if o.Address != nil {
		r.Address = o.Address
	}

	if o.Enabled != nil {
		r.Enabled = o.Enabled
	}

	if o.Namespace != nil {
		r.Namespace = o.Namespace
	}

	if o.RenewToken != nil {
		r.RenewToken = o.RenewToken
	}

	if o.Retry != nil {
		r.Retry = r.Retry.Merge(o.Retry)
	}

	if o.SSL != nil {
		r.SSL = r.SSL.Merge(o.SSL)
	}

	if o.Token != nil {
		r.Token = o.Token
	}

	if o.VaultAgentTokenFile != nil {
		r.VaultAgentTokenFile = o.VaultAgentTokenFile
	}

	if o.Transport != nil {
		r.Transport = r.Transport.Merge(o.Transport)
	}

	if o.UnwrapToken != nil {
		r.UnwrapToken = o.UnwrapToken
	}

	return r
}

// Finalize ensures there no nil pointers.
func (c *VaultConfig) Finalize() {
	if c.Address == nil {
		c.Address = stringFromEnv([]string{
			api.EnvVaultAddress,
		}, "")
	}

	if c.Namespace == nil {
		c.Namespace = stringFromEnv([]string{"VAULT_NAMESPACE"}, "")
	}

	if c.Retry == nil {
		c.Retry = DefaultRetryConfig()
	}
	c.Retry.Finalize()

	// Vault has custom SSL settings
	if c.SSL == nil {
		c.SSL = DefaultSSLConfig()
	}
	if c.SSL.Enabled == nil {
		c.SSL.Enabled = Bool(true)
	}
	if c.SSL.CaCert == nil {
		c.SSL.CaCert = stringFromEnv([]string{api.EnvVaultCACert}, "")
	}
	if c.SSL.CaPath == nil {
		c.SSL.CaPath = stringFromEnv([]string{api.EnvVaultCAPath}, "")
	}
	if c.SSL.Cert == nil {
		c.SSL.Cert = stringFromEnv([]string{api.EnvVaultClientCert}, "")
	}
	if c.SSL.Key == nil {
		c.SSL.Key = stringFromEnv([]string{api.EnvVaultClientKey}, "")
	}
	if c.SSL.ServerName == nil {
		c.SSL.ServerName = stringFromEnv([]string{api.EnvVaultTLSServerName}, "")
	}
	if c.SSL.Verify == nil {
		c.SSL.Verify = antiboolFromEnv([]string{
			EnvVaultSkipVerify, api.EnvVaultInsecure}, true)
	}
	c.SSL.Finalize()

	// Order of precedence
	// 1. `vault_agent_token_file` configuration value
	// 2. `token` configuration value`
	// 3. `VAULT_TOKEN` environment variable
	if c.Token == nil {
		c.Token = stringFromEnv([]string{
			"VAULT_TOKEN",
		}, "")
	}

	if c.VaultAgentTokenFile == nil {
		if StringVal(c.Token) == "" {
			if homePath != "" {
				c.Token = stringFromFile([]string{
					homePath + "/.vault-token",
				}, "")
			}
		}
	} else {
		c.Token = stringFromFile([]string{*c.VaultAgentTokenFile}, "")
	}

	// must be after c.Token setting, as default depends on that.
	if c.RenewToken == nil {
		default_renew := DefaultVaultRenewToken
		if c.VaultAgentTokenFile != nil {
			default_renew = false
		} else if StringVal(c.Token) == "" {
			default_renew = false
		}
		c.RenewToken = boolFromEnv([]string{
			"VAULT_RENEW_TOKEN",
		}, default_renew)
	}

	if c.Transport == nil {
		c.Transport = DefaultTransportConfig()
	}
	c.Transport.Finalize()

	if c.UnwrapToken == nil {
		c.UnwrapToken = boolFromEnv([]string{
			"VAULT_UNWRAP_TOKEN",
		}, DefaultVaultUnwrapToken)
	}

	if c.Enabled == nil {
		c.Enabled = Bool(StringPresent(c.Address))
	}
}

// GoString defines the printable version of this struct.
func (c *VaultConfig) GoString() string {
	if c == nil {
		return "(*VaultConfig)(nil)"
	}

	return fmt.Sprintf("&VaultConfig{"+
		"Address:%s, "+
		"Enabled:%s, "+
		"Namespace:%s,"+
		"RenewToken:%s, "+
		"Retry:%#v, "+
		"SSL:%#v, "+
		"Token:%t, "+
		"VaultAgentTokenFile:%t, "+
		"Transport:%#v, "+
		"UnwrapToken:%s"+
		"}",
		StringGoString(c.Address),
		BoolGoString(c.Enabled),
		StringGoString(c.Namespace),
		BoolGoString(c.RenewToken),
		c.Retry,
		c.SSL,
		StringPresent(c.Token),
		StringPresent(c.VaultAgentTokenFile),
		c.Transport,
		BoolGoString(c.UnwrapToken),
	)
}
