// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import "fmt"

// NomadConfig is the configuration for connecting to a Nomad agent.
type NomadConfig struct {
	// Address is the URI to the Nomad agent.
	Address *string `mapstructure:"address"`

	// Enabled controls whether the Nomad integration is active.
	Enabled *bool `mapstructure:"enabled"`

	// Namespace is the Nomad namespace to use. This can also be set via
	// the NOMAD_NAMESPACE environment variable.
	Namespace *string `mapstructure:"namespace"`

	// SSL indicates we should use a secure connection while talking
	// to Nomad.
	SSL *SSLConfig `mapstructure:"ssl"`

	// Token is the Nomad ACL token to use with API requests. This can also
	// be set via the NOMAD_TOKEN environment variable.
	Token *string `mapstructure:"token" json:"-"`

	// AuthUsername and AuthPassword are the HTTP Basic Auth username and
	// password to use when authenticating with the Nomad API.
	AuthUsername *string `mapstructure:"auth_username"`
	AuthPassword *string `mapstructure:"auth_password"`

	// Transport configures the low-level network connection details.
	Transport *TransportConfig `mapstructure:"transport"`

	// Retry is the configuration for specifying how to behave on failure.
	Retry *RetryConfig `mapstructure:"retry"`
}

func DefaultNomadConfig() *NomadConfig {
	return &NomadConfig{
		SSL:       DefaultSSLConfig(),
		Transport: DefaultTransportConfig(),
	}
}

// Copy returns a deep copy of this configuration.
func (n *NomadConfig) Copy() *NomadConfig {
	if n == nil {
		return nil
	}

	var o NomadConfig

	o.Address = n.Address
	o.Enabled = n.Enabled
	o.Namespace = n.Namespace
	o.SSL = n.SSL.Copy()
	o.Token = n.Token
	o.AuthUsername = n.AuthUsername
	o.AuthPassword = n.AuthPassword
	o.Transport = n.Transport.Copy()
	o.Retry = n.Retry.Copy()

	return &o
}

// Merge combines all values in this configuration with the values in the other
// configuration, with values in the other configuration taking precedence.
// Maps and slices are merged, most other values are overwritten. Complex
// structs define their own merge functionality.
func (n *NomadConfig) Merge(o *NomadConfig) *NomadConfig {
	if n == nil {
		if o == nil {
			return nil
		}
		return o.Copy()
	}

	if o == nil {
		return n.Copy()
	}

	r := n.Copy()

	if o.Address != nil {
		r.Address = o.Address
	}

	if o.Enabled != nil {
		r.Enabled = o.Enabled
	}

	if o.Namespace != nil {
		r.Namespace = o.Namespace
	}

	if o.SSL != nil {
		r.SSL = r.SSL.Merge(o.SSL)
	}

	if o.Token != nil {
		r.Token = o.Token
	}

	if o.AuthUsername != nil {
		r.AuthUsername = o.AuthUsername
	}

	if o.AuthPassword != nil {
		r.AuthPassword = o.AuthPassword
	}

	if o.Transport != nil {
		r.Transport = r.Transport.Merge(o.Transport)
	}

	if o.Retry != nil {
		r.Retry = r.Retry.Merge(o.Retry)
	}

	return r
}

// Finalize ensures there no nil pointers.
func (n *NomadConfig) Finalize() {
	if n.Address == nil {
		n.Address = stringFromEnv([]string{"NOMAD_ADDR"}, "")
	}

	if n.Enabled == nil {
		// Enable if there's an address or custom dialer
		customDialer := n.Transport != nil && n.Transport.CustomDialer != nil
		addressPresent := n.Address != nil && *n.Address != ""
		n.Enabled = Bool(addressPresent || customDialer)
	}

	if n.Namespace == nil {
		n.Namespace = stringFromEnv([]string{"NOMAD_NAMESPACE"}, "")
	}

	if n.SSL == nil {
		n.SSL = DefaultSSLConfig()
	}
	n.SSL.Finalize()

	if n.Token == nil {
		n.Token = stringFromEnv([]string{"NOMAD_TOKEN"}, "")
	}

	if n.AuthUsername == nil {
		n.AuthUsername = String("")
	}

	if n.AuthPassword == nil {
		n.AuthPassword = String("")
	}

	if n.Transport == nil {
		n.Transport = DefaultTransportConfig()
	}
	n.Transport.Finalize()

	if n.Retry == nil {
		n.Retry = DefaultRetryConfig()
	}
	n.Retry.Finalize()
}

func (n *NomadConfig) GoString() string {
	if n == nil {
		return "(*NomadConfig)(nil)"
	}

	return fmt.Sprintf("&NomadConfig{"+
		"Address:%s, "+
		"Enabled:%s, "+
		"Namespace:%s, "+
		"SSL:%#v, "+
		"Token:%s, "+
		"AuthUsername:%s, "+
		"AuthPassword:%s, "+
		"Transport:%#v, "+
		"Retry:%#v, "+
		"}",
		StringGoString(n.Address),
		BoolGoString(n.Enabled),
		StringGoString(n.Namespace),
		n.SSL,
		StringGoString(n.Token),
		StringGoString(n.AuthUsername),
		StringGoString(n.AuthPassword),
		n.Transport,
		n.Retry,
	)
}
