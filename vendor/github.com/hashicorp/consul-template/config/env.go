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
	// Denylist specifies a list of environment variables to explicitly
	// exclude from the list of environment variables populated to the child.
	// If both Allowlist and Denylist are provided, Denylist takes
	// precedence over the values in Allowlist.
	Denylist []string `mapstructure:"denylist"`

	// DenylistDeprecated is the backward compatible option for Denylist for
	// configuration supported by v0.25.0 and older. This should not be used
	// directly, use Denylist instead. Values from this are combined to
	// Denylist in Finalize().
	DenylistDeprecated []string `mapstructure:"blacklist" json:"-"`

	// CustomEnv specifies custom environment variables to pass to the child
	// process. These are provided programmatically, override any environment
	// variables of the same name, are ignored from allowlist/denylist, and
	// are still included even if PristineEnv is set to true.
	Custom []string `mapstructure:"custom"`

	// PristineEnv specifies if the child process should inherit the parent's
	// environment.
	Pristine *bool `mapstructure:"pristine"`

	// Allowlist specifies a list of environment variables to exclusively
	// include in the list of environment variables populated to the child.
	Allowlist []string `mapstructure:"allowlist"`

	// AllowlistDeprecated is the backward compatible option for Allowlist for
	// configuration supported by v0.25.0 and older. This should not be used
	// directly, use Allowlist instead. Values from this are combined to
	// Allowlist in Finalize().
	AllowlistDeprecated []string `mapstructure:"whitelist" json:"-"`
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

	if c.Denylist != nil {
		o.Denylist = append([]string{}, c.Denylist...)
	}

	if c.DenylistDeprecated != nil {
		o.DenylistDeprecated = append([]string{}, c.DenylistDeprecated...)
	}

	if c.Custom != nil {
		o.Custom = append([]string{}, c.Custom...)
	}

	o.Pristine = c.Pristine

	if c.Allowlist != nil {
		o.Allowlist = append([]string{}, c.Allowlist...)
	}

	if c.AllowlistDeprecated != nil {
		o.AllowlistDeprecated = append([]string{}, c.AllowlistDeprecated...)
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

	if o.Denylist != nil {
		r.Denylist = append(r.Denylist, o.Denylist...)
	}

	if o.DenylistDeprecated != nil {
		r.DenylistDeprecated = append(r.DenylistDeprecated, o.DenylistDeprecated...)
	}

	if o.Custom != nil {
		r.Custom = append(r.Custom, o.Custom...)
	}

	if o.Pristine != nil {
		r.Pristine = o.Pristine
	}

	if o.Allowlist != nil {
		r.Allowlist = append(r.Allowlist, o.Allowlist...)
	}

	if o.AllowlistDeprecated != nil {
		r.AllowlistDeprecated = append(r.AllowlistDeprecated, o.AllowlistDeprecated...)
	}

	return r
}

// Env calculates and returns the finalized environment for this exec
// configuration. It takes into account pristine, custom environment, allowlist,
// and denylist values.
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

	// Pull out any envvars that match the allowlist.
	// Combining lists on each reference may be slightly inefficient but this
	// allows for out of order method calls, not requiring the config to be
	// finalized first.
	allowlist := combineLists(c.Allowlist, c.AllowlistDeprecated)
	if len(allowlist) > 0 {
		newKeys := make([]string, 0, len(keys))
		for _, k := range keys {
			if anyGlobMatch(k, allowlist) {
				newKeys = append(newKeys, k)
			}
		}
		keys = newKeys
	}

	// Remove any envvars that match the denylist.
	// Combining lists on each reference may be slightly inefficient but this
	// allows for out of order method calls, not requiring the config to be
	// finalized first.
	denylist := combineLists(c.Denylist, c.DenylistDeprecated)
	if len(denylist) > 0 {
		newKeys := make([]string, 0, len(keys))
		for _, k := range keys {
			if !anyGlobMatch(k, denylist) {
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
	if c.Denylist == nil && c.DenylistDeprecated == nil {
		c.Denylist = []string{}
		c.DenylistDeprecated = []string{}
	} else {
		c.Denylist = combineLists(c.Denylist, c.DenylistDeprecated)
	}

	if c.Custom == nil {
		c.Custom = []string{}
	}

	if c.Pristine == nil {
		c.Pristine = Bool(false)
	}

	if c.Allowlist == nil && c.AllowlistDeprecated == nil {
		c.Allowlist = []string{}
		c.AllowlistDeprecated = []string{}
	} else {
		c.Allowlist = combineLists(c.Allowlist, c.AllowlistDeprecated)
	}
}

// GoString defines the printable version of this struct.
func (c *EnvConfig) GoString() string {
	if c == nil {
		return "(*EnvConfig)(nil)"
	}

	return fmt.Sprintf("&EnvConfig{"+
		"Denylist:%v, "+
		"Custom:%v, "+
		"Pristine:%s, "+
		"Allowlist:%v"+
		"}",
		combineLists(c.Denylist, c.DenylistDeprecated),
		c.Custom,
		BoolGoString(c.Pristine),
		combineLists(c.Allowlist, c.AllowlistDeprecated),
	)
}

// combineLists makes a new list that combines 2 lists by adding values from
// the second list without removing any duplicates from the first.
func combineLists(a, b []string) []string {
	combined := make([]string, len(a), len(a)+len(b))
	m := make(map[string]bool)
	for i, v := range a {
		m[v] = true
		combined[i] = v
	}

	for _, v := range b {
		if !m[v] {
			combined = append(combined, v)
		}
	}

	return combined
}
