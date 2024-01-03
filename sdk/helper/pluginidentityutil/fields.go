// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pluginidentityutil

import (
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
)

func AddPluginIdentityTokenFields(m map[string]*framework.FieldSchema) {
	f := PluginIdentityTokenFields()
	for k, v := range f {
		if _, ok := m[k]; ok {
			panic(fmt.Sprintf("adding field %q would overwrite existing field", k))
		}
		m[k] = v
	}
}

func PluginIdentityTokenFields() map[string]*framework.FieldSchema {
	return map[string]*framework.FieldSchema{
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
}
