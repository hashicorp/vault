// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pluginidentityutil

import (
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
)

// PluginIdentityTokenParams contains a set of common parameters that plugins
// can use for setting plugin identity token behavior
type PluginIdentityTokenParams struct {
	// IdentityTokenKey is the named key used to sign tokens
	IdentityTokenKey string `json:"identity_token_key"`
	// IdentityTokenTTLSeconds is the duration that tokens will be valid for
	IdentityTokenTTLSeconds int `json:"identity_token_ttl_seconds"`
	// IdentityTokenAudience identifies the recipient of the token
	IdentityTokenAudience string `json:"identity_token_audience"`
}

// AddPluginIdentityTokenFields adds plugin identity token fields to the given
// field schema map.
func AddPluginIdentityTokenFields(m map[string]*framework.FieldSchema) {
	fields := map[string]*framework.FieldSchema{
		"identity_token_audience": {
			Type:        framework.TypeString,
			Description: "",
			Default:     "",
			DisplayAttrs: &framework.DisplayAttributes{
				Name: "Audience of plugin identity tokens",
			},
		},
		"identity_token_key": {
			Type:        framework.TypeString,
			Description: "",
			Default:     "",
			DisplayAttrs: &framework.DisplayAttributes{
				Name: "Key used to sign plugin identity tokens",
			},
		},
		"identity_token_ttl": {
			Type:        framework.TypeDurationSecond,
			Description: "",
			DisplayAttrs: &framework.DisplayAttributes{
				Name: "Time-to-live of plugin identity tokens",
			},
			Default: 3600,
		},
	}

	for name, schema := range fields {
		if _, ok := m[name]; ok {
			panic(fmt.Sprintf("adding field %q would overwrite existing field", name))
		}
		m[name] = schema
	}
}
