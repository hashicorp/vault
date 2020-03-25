package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// EnvConfig is an embeddable struct for things that accept environment
// variable filtering. You should not use this directly and it is only public
// for mapstructure's decoding.
type EnvConfig struct {
	// BlacklistEnv specifies a list of environment variables to explicitly
	// exclude from the list of environment variables populated to the child.
	// If both WhitelistEnv and BlacklistEnv are provided, BlacklistEnv takes
	// precedence over the values in WhitelistEnv.
	Blacklist []string `mapstructure:"blacklist"`

	// CustomEnv specifies custom environment variables to pass to the child
	// process. These are provided programmatically, override any environment
	// variables of the same name, are ignored from whitelist/blacklist, and
	// are still included even if PristineEnv is set to true.
	Custom []string `mapstructure:"custom"`

	// PristineEnv specifies if the child process should inherit the parent's
	// environment.
	Pristine *bool `mapstructure:"pristine"`

	// WhitelistEnv specifies a list of environment variables to exclusively
	// include in the list of environment variables populated to the child.
	Whitelist []string `mapstructure:"whitelist"`
}

// DefaultEnvConfig returns a configuration that is populated with the
// default values.
func DefaultEnvConfig() *EnvConfig {
	return &EnvConfig{}
}

// Copy returns a deep copy of this configuration.
func (c *EnvConfig) Copy() *EnvConfig {
	if c == nil {
		return nil
	}

	var o EnvConfig

	if c.Blacklist != nil {
		o.Blacklist = append([]string{}, c.Blacklist...)
	}

	if c.Custom != nil {
		o.Custom = append([]string{}, c.Custom...)
	}

	o.Pristine = c.Pristine

	if c.Whitelist != nil {
		o.Whitelist = append([]string{}, c.Whitelist...)
	}

	return &o
}

// Merge combines all values in this configuration with the values in the other
// configuration, with values in the other configuration taking precedence.
// Maps and slices are merged, most other values are overwritten. Complex
// structs define their own merge functionality.
func (c *EnvConfig) Merge(o *EnvConfig) *EnvConfig {
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

	if o.Blacklist != nil {
		r.Blacklist = append(r.Blacklist, o.Blacklist...)
	}

	if o.Custom != nil {
		r.Custom = append(r.Custom, o.Custom...)
	}

	if o.Pristine != nil {
		r.Pristine = o.Pristine
	}

	if o.Whitelist != nil {
		r.Whitelist = append(r.Whitelist, o.Whitelist...)
	}

	return r
}

// Env calculates and returns the finalized environment for this exec
// configuration. It takes into account pristine, custom environment, whitelist,
// and blacklist values.
func (c *EnvConfig) Env() []string {
	// In pristine mode, just return the custom environment. If the user did not
	// specify a custom environment, just return the empty slice to force an
	// empty environment. We cannot return nil here because the later call to
	// os/exec will think we want to inherit the parent.
	if BoolVal(c.Pristine) {
		if len(c.Custom) > 0 {
			return c.Custom
		}
		return []string{}
	}

	// Pull all the key-value pairs out of the environment
	environ := os.Environ()
	keys := make([]string, len(environ))
	env := make(map[string]string, len(environ))
	for i, v := range environ {
		list := strings.SplitN(v, "=", 2)
		keys[i] = list[0]
		env[list[0]] = list[1]
	}

	// anyGlobMatch is a helper function which checks if any of the given globs
	// match the string.
	anyGlobMatch := func(s string, patterns []string) bool {
		for _, pattern := range patterns {
			if matched, _ := filepath.Match(pattern, s); matched {
				return true
			}
		}
		return false
	}

	// Pull out any envvars that match the whitelist.
	if len(c.Whitelist) > 0 {
		newKeys := make([]string, 0, len(keys))
		for _, k := range keys {
			if anyGlobMatch(k, c.Whitelist) {
				newKeys = append(newKeys, k)
			}
		}
		keys = newKeys
	}

	// Remove any envvars that match the blacklist.
	if len(c.Blacklist) > 0 {
		newKeys := make([]string, 0, len(keys))
		for _, k := range keys {
			if !anyGlobMatch(k, c.Blacklist) {
				newKeys = append(newKeys, k)
			}
		}
		keys = newKeys
	}

	// Build the final list using only the filtered keys.
	finalEnv := make([]string, 0, len(keys)+len(c.Custom))
	for _, k := range keys {
		finalEnv = append(finalEnv, k+"="+env[k])
	}

	// Append remaining custom environment.
	finalEnv = append(finalEnv, c.Custom...)

	return finalEnv
}

// Finalize ensures there no nil pointers.
func (c *EnvConfig) Finalize() {
	if c.Blacklist == nil {
		c.Blacklist = []string{}
	}

	if c.Custom == nil {
		c.Custom = []string{}
	}

	if c.Pristine == nil {
		c.Pristine = Bool(false)
	}

	if c.Whitelist == nil {
		c.Whitelist = []string{}
	}
}

// GoString defines the printable version of this struct.
func (c *EnvConfig) GoString() string {
	if c == nil {
		return "(*EnvConfig)(nil)"
	}

	return fmt.Sprintf("&EnvConfig{"+
		"Blacklist:%v, "+
		"Custom:%v, "+
		"Pristine:%s, "+
		"Whitelist:%v"+
		"}",
		c.Blacklist,
		c.Custom,
		BoolGoString(c.Pristine),
		c.Whitelist,
	)
}
