// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pluginidentityutil

import (
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/stretchr/testify/assert"
)

const (
	fieldIDTokenTTL      = "identity_token_ttl"
	fieldIDTokenAudience = "identity_token_audience"
)

func identityTokenFieldData(raw map[string]interface{}) *framework.FieldData {
	return &framework.FieldData{
		Raw: raw,
		Schema: map[string]*framework.FieldSchema{
			fieldIDTokenTTL: {
				Type: framework.TypeDurationSecond,
			},
			fieldIDTokenAudience: {
				Type: framework.TypeString,
			},
		},
	}
}

func TestParsePluginIdentityTokenFields(t *testing.T) {
	testcases := []struct {
		name    string
		d       *framework.FieldData
		wantErr bool
		want    map[string]interface{}
	}{
		{
			name: "all input",
			d: identityTokenFieldData(map[string]interface{}{
				fieldIDTokenTTL:      10,
				fieldIDTokenAudience: "test-aud",
			}),
			want: map[string]interface{}{
				fieldIDTokenTTL:      time.Duration(10) * time.Second,
				fieldIDTokenAudience: "test-aud",
			},
		},
		{
			name: "empty ttl",
			d: identityTokenFieldData(map[string]interface{}{
				fieldIDTokenAudience: "test-aud",
			}),
			want: map[string]interface{}{
				fieldIDTokenTTL:      time.Duration(0),
				fieldIDTokenAudience: "test-aud",
			},
		},
		{
			name: "empty audience",
			d: identityTokenFieldData(map[string]interface{}{
				fieldIDTokenTTL: 10,
			}),
			want: map[string]interface{}{
				fieldIDTokenTTL:      time.Duration(10) * time.Second,
				fieldIDTokenAudience: "",
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			p := new(PluginIdentityTokenParams)
			err := p.ParsePluginIdentityTokenFields(tt.d)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			got := map[string]interface{}{
				fieldIDTokenTTL:      p.IdentityTokenTTL,
				fieldIDTokenAudience: p.IdentityTokenAudience,
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPopulatePluginIdentityTokenData(t *testing.T) {
	testcases := []struct {
		name string
		p    *PluginIdentityTokenParams
		want map[string]interface{}
	}{
		{
			name: "basic",
			p: &PluginIdentityTokenParams{
				IdentityTokenAudience: "test-aud",
				IdentityTokenTTL:      time.Duration(10) * time.Second,
			},
			want: map[string]interface{}{
				fieldIDTokenTTL:      int64(10),
				fieldIDTokenAudience: "test-aud",
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			got := make(map[string]interface{})
			tt.p.PopulatePluginIdentityTokenData(got)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAddPluginIdentityTokenFields(t *testing.T) {
	testcases := []struct {
		name  string
		input map[string]*framework.FieldSchema
		want  map[string]*framework.FieldSchema
	}{
		{
			name:  "basic",
			input: map[string]*framework.FieldSchema{},
			want: map[string]*framework.FieldSchema{
				fieldIDTokenAudience: {
					Type:        framework.TypeString,
					Description: "Audience of plugin identity tokens",
					Default:     "",
				},
				fieldIDTokenTTL: {
					Type:        framework.TypeDurationSecond,
					Description: "Time-to-live of plugin identity tokens",
					Default:     3600,
				},
			},
		},
		{
			name: "additional-fields",
			input: map[string]*framework.FieldSchema{
				"test": {
					Type:        framework.TypeString,
					Description: "Test description",
					Default:     "default",
				},
			},
			want: map[string]*framework.FieldSchema{
				fieldIDTokenAudience: {
					Type:        framework.TypeString,
					Description: "Audience of plugin identity tokens",
					Default:     "",
				},
				fieldIDTokenTTL: {
					Type:        framework.TypeDurationSecond,
					Description: "Time-to-live of plugin identity tokens",
					Default:     3600,
				},
				"test": {
					Type:        framework.TypeString,
					Description: "Test description",
					Default:     "default",
				},
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input
			AddPluginIdentityTokenFields(got)
			assert.Equal(t, tt.want, got)
		})
	}
}
