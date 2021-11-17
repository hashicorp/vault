package config

import "fmt"

const (
	// DefaultSSLVerify is the default value for SSL verification.
	DefaultSSLVerify = true
)

// SSLConfig is the configuration for SSL.
type SSLConfig struct {
	CaCert     *string `mapstructure:"ca_cert"`
	CaPath     *string `mapstructure:"ca_path"`
	Cert       *string `mapstructure:"cert"`
	Enabled    *bool   `mapstructure:"enabled"`
	Key        *string `mapstructure:"key"`
	ServerName *string `mapstructure:"server_name"`
	Verify     *bool   `mapstructure:"verify"`
}

// DefaultSSLConfig returns a configuration that is populated with the
// default values.
func DefaultSSLConfig() *SSLConfig {
	return &SSLConfig{}
}

// Copy returns a deep copy of this configuration.
func (c *SSLConfig) Copy() *SSLConfig {
	if c == nil {
		return nil
	}

	var o SSLConfig
	o.CaCert = c.CaCert
	o.CaPath = c.CaPath
	o.Cert = c.Cert
	o.Enabled = c.Enabled
	o.Key = c.Key
	o.ServerName = c.ServerName
	o.Verify = c.Verify
	return &o
}

// Merge combines all values in this configuration with the values in the other
// configuration, with values in the other configuration taking precedence.
// Maps and slices are merged, most other values are overwritten. Complex
// structs define their own merge functionality.
func (c *SSLConfig) Merge(o *SSLConfig) *SSLConfig {
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

	if o.Cert != nil {
		r.Cert = o.Cert
	}

	if o.CaCert != nil {
		r.CaCert = o.CaCert
	}

	if o.CaPath != nil {
		r.CaPath = o.CaPath
	}

	if o.Enabled != nil {
		r.Enabled = o.Enabled
	}

	if o.Key != nil {
		r.Key = o.Key
	}

	if o.ServerName != nil {
		r.ServerName = o.ServerName
	}

	if o.Verify != nil {
		r.Verify = o.Verify
	}

	return r
}

// Finalize ensures there no nil pointers.
func (c *SSLConfig) Finalize() {
	if c.Enabled == nil {
		c.Enabled = Bool(false ||
			StringPresent(c.Cert) ||
			StringPresent(c.CaCert) ||
			StringPresent(c.CaPath) ||
			StringPresent(c.Key) ||
			StringPresent(c.ServerName) ||
			BoolPresent(c.Verify))
	}

	if c.Cert == nil {
		c.Cert = String("")
	}

	if c.CaCert == nil {
		c.CaCert = String("")
	}

	if c.CaPath == nil {
		c.CaPath = String("")
	}

	if c.Key == nil {
		c.Key = String("")
	}

	if c.ServerName == nil {
		c.ServerName = String("")
	}

	if c.Verify == nil {
		c.Verify = Bool(DefaultSSLVerify)
	}
}

// GoString defines the printable version of this struct.
func (c *SSLConfig) GoString() string {
	if c == nil {
		return "(*SSLConfig)(nil)"
	}

	return fmt.Sprintf("&SSLConfig{"+
		"CaCert:%s, "+
		"CaPath:%s, "+
		"Cert:%s, "+
		"Enabled:%s, "+
		"Key:%s, "+
		"ServerName:%s, "+
		"Verify:%s"+
		"}",
		StringGoString(c.CaCert),
		StringGoString(c.CaPath),
		StringGoString(c.Cert),
		BoolGoString(c.Enabled),
		StringGoString(c.Key),
		StringGoString(c.ServerName),
		BoolGoString(c.Verify),
	)
}
