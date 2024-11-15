// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpkms

import (
	"strings"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/framework"
)

const (
	defaultScope = "https://www.googleapis.com/auth/cloudkms"
)

// Config is the stored configuration.
type Config struct {
	Credentials string   `json:"credentials"`
	Scopes      []string `json:"scopes"`
}

// DefaultConfig returns a config with the default values.
func DefaultConfig() *Config {
	return &Config{
		Scopes: []string{defaultScope},
	}
}

// Update updates the configuration from the given field data.
func (c *Config) Update(d *framework.FieldData) (bool, error) {
	if d == nil {
		return false, nil
	}

	changed := false

	if v, ok := d.GetOk("credentials"); ok {
		nv := strings.TrimSpace(v.(string))
		if nv != c.Credentials {
			c.Credentials = nv
			changed = true
		}
	}

	if v, ok := d.GetOk("scopes"); ok {
		nv := strutil.RemoveDuplicates(v.([]string), true)
		if !strutil.EquivalentSlices(nv, c.Scopes) {
			c.Scopes = nv
			changed = true
		}
	}

	return changed, nil
}
