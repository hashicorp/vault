// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pluginidentityutil

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
)

// PluginIdentityTokenParams contains a set of common parameters that plugins
// can use for setting plugin identity token behavior.
type PluginIdentityTokenParams struct {
	// IdentityTokenTTL is the duration that tokens will be valid for
	IdentityTokenTTL time.Duration `json:"identity_token_ttl"`
	// IdentityTokenAudience identifies the recipient of the token
	IdentityTokenAudience string `json:"identity_token_audience"`
}

// ParsePluginIdentityTokenFields provides common field parsing to embedding structs.
func (p *PluginIdentityTokenParams) ParsePluginIdentityTokenFields(d *framework.FieldData) error {
	if tokenTTLRaw, ok := d.GetOk("identity_token_ttl"); ok {
		p.IdentityTokenTTL = time.Duration(tokenTTLRaw.(int)) * time.Second
	}

	if tokenAudienceRaw, ok := d.GetOk("identity_token_audience"); ok {
		p.IdentityTokenAudience = tokenAudienceRaw.(string)
	}

	return nil
}

// PopulatePluginIdentityTokenData adds PluginIdentityTokenParams info into the given map.
func (p *PluginIdentityTokenParams) PopulatePluginIdentityTokenData(m map[string]interface{}) {
	m["identity_token_ttl"] = int64(p.IdentityTokenTTL.Seconds())
	m["identity_token_audience"] = p.IdentityTokenAudience
}

// AddPluginIdentityTokenFields adds plugin identity token fields to the given
// field schema map.
func AddPluginIdentityTokenFields(m map[string]*framework.FieldSchema) {
	fields := map[string]*framework.FieldSchema{
		"identity_token_audience": {
			Type:        framework.TypeString,
			Description: "Audience of plugin identity tokens",
			Default:     "",
		},
		"identity_token_ttl": {
			Type:        framework.TypeDurationSecond,
			Description: "Time-to-live of plugin identity tokens",
			Default:     3600,
		},
	}

	for name, schema := range fields {
		if _, ok := m[name]; ok {
			panic(fmt.Sprintf("adding field %q would overwrite existing field", name))
		}
		m[name] = schema
	}
}
