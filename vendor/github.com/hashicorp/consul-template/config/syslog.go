package config

import (
	"fmt"

	"github.com/hashicorp/consul-template/version"
)

const (
	// DefaultSyslogFacility is the default facility to log to.
	DefaultSyslogFacility = "LOCAL0"
)

var (
	// DefaultSyslogName is the default app name in syslog.
	DefaultSyslogName = version.Name
)

// SyslogConfig is the configuration for syslog.
type SyslogConfig struct {
	Enabled  *bool   `mapstructure:"enabled"`
	Facility *string `mapstructure:"facility"`
	Name	 *string `mapstructure:"name"`
}

// DefaultSyslogConfig returns a configuration that is populated with the
// default values.
func DefaultSyslogConfig() *SyslogConfig {
	return &SyslogConfig{}
}

// Copy returns a deep copy of this configuration.
func (c *SyslogConfig) Copy() *SyslogConfig {
	if c == nil {
		return nil
	}

	var o SyslogConfig
	o.Enabled = c.Enabled
	o.Facility = c.Facility
	o.Name = c.Name
	return &o
}

// Merge combines all values in this configuration with the values in the other
// configuration, with values in the other configuration taking precedence.
// Maps and slices are merged, most other values are overwritten. Complex
// structs define their own merge functionality.
func (c *SyslogConfig) Merge(o *SyslogConfig) *SyslogConfig {
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

	if o.Enabled != nil {
		r.Enabled = o.Enabled
	}

	if o.Facility != nil {
		r.Facility = o.Facility
	}

	if o.Name != nil {
		r.Name = o.Name
	}

	return r
}

// Finalize ensures there no nil pointers.
func (c *SyslogConfig) Finalize() {
	if c.Enabled == nil {
		c.Enabled = Bool(StringPresent(c.Facility) || StringPresent(c.Name))
	}

	if c.Facility == nil {
		c.Facility = String(DefaultSyslogFacility)
	}

	if c.Name == nil {
		c.Name = String(DefaultSyslogName)
	}
}

// GoString defines the printable version of this struct.
func (c *SyslogConfig) GoString() string {
	if c == nil {
		return "(*SyslogConfig)(nil)"
	}

	return fmt.Sprintf("&SyslogConfig{"+
		"Enabled:%s, "+
		"Facility:%s"+
		"Name:%s"+
		"}",
		BoolGoString(c.Enabled),
		StringGoString(c.Facility),
		StringGoString(c.Name),
	)
}
