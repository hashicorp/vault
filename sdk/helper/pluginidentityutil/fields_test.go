// Copyright IBM Corp. 2016, 2025
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

func TestAddPluginIdentityTokenFieldsWithGroup(t *testing.T) {
	testcases := []struct {
		name  string
		group string
		input map[string]*framework.FieldSchema
		want  map[string]*framework.FieldSchema
	}{
		{
			name:  "basic",
			input: map[string]*framework.FieldSchema{},
			group: "Tokens",
			want: map[string]*framework.FieldSchema{
				fieldIDTokenAudience: {
					Type:        framework.TypeString,
					Description: "Audience of plugin identity tokens",
					Default:     "",
					DisplayAttrs: &framework.DisplayAttributes{
						Group: "Tokens",
					},
				},
				fieldIDTokenTTL: {
					Type:        framework.TypeDurationSecond,
					Description: "Time-to-live of plugin identity tokens",
					Default:     3600,
					DisplayAttrs: &framework.DisplayAttributes{
						Name:  "Identity token TTL",
						Group: "Tokens",
					},
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
			group: "Token Options",
			want: map[string]*framework.FieldSchema{
				fieldIDTokenAudience: {
					Type:        framework.TypeString,
					Description: "Audience of plugin identity tokens",
					Default:     "",
					DisplayAttrs: &framework.DisplayAttributes{
						Group: "Token Options",
					},
				},
				fieldIDTokenTTL: {
					Type:        framework.TypeDurationSecond,
					Description: "Time-to-live of plugin identity tokens",
					Default:     3600,
					DisplayAttrs: &framework.DisplayAttributes{
						Name:  "Identity token TTL",
						Group: "Token Options",
					},
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
			AddPluginIdentityTokenFieldsWithGroup(got, tt.group)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAddPluginIdentityTokenFields(t *testing.T) {
	testcases := []struct {
		name  string
		group string
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
					DisplayAttrs: &framework.DisplayAttributes{
						Group: "default",
					},
				},
				fieldIDTokenTTL: {
					Type:        framework.TypeDurationSecond,
					Description: "Time-to-live of plugin identity tokens",
					Default:     3600,
					DisplayAttrs: &framework.DisplayAttributes{
						Name:  "Identity token TTL",
						Group: "default",
					},
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

func TestPluginIdentityTokenParams_Equals(t *testing.T) {
	testcases := []struct {
		name     string
		p1       PluginIdentityTokenParams
		p2       PluginIdentityTokenParams
		expected bool
	}{
		{
			name: "equal-all-fields",
			p1: PluginIdentityTokenParams{
				IdentityTokenTTL:      10 * time.Second,
				IdentityTokenAudience: "test-aud",
			},
			p2: PluginIdentityTokenParams{
				IdentityTokenTTL:      10 * time.Second,
				IdentityTokenAudience: "test-aud",
			},
			expected: true,
		},
		{
			name: "equal-zero-values",
			p1: PluginIdentityTokenParams{
				IdentityTokenTTL:      0,
				IdentityTokenAudience: "",
			},
			p2: PluginIdentityTokenParams{
				IdentityTokenTTL:      0,
				IdentityTokenAudience: "",
			},
			expected: true,
		},
		{
			name: "different-ttl",
			p1: PluginIdentityTokenParams{
				IdentityTokenTTL:      10 * time.Second,
				IdentityTokenAudience: "test-aud",
			},
			p2: PluginIdentityTokenParams{
				IdentityTokenTTL:      20 * time.Second,
				IdentityTokenAudience: "test-aud",
			},
			expected: false,
		},
		{
			name: "different-audience",
			p1: PluginIdentityTokenParams{
				IdentityTokenTTL:      10 * time.Second,
				IdentityTokenAudience: "test-aud",
			},
			p2: PluginIdentityTokenParams{
				IdentityTokenTTL:      10 * time.Second,
				IdentityTokenAudience: "different-aud",
			},
			expected: false,
		},
		{
			name: "different-both",
			p1: PluginIdentityTokenParams{
				IdentityTokenTTL:      10 * time.Second,
				IdentityTokenAudience: "test-aud",
			},
			p2: PluginIdentityTokenParams{
				IdentityTokenTTL:      20 * time.Second,
				IdentityTokenAudience: "different-aud",
			},
			expected: false,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.p1.Equals(tt.p2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPluginIdentityTokenParams_Equals_Embedded(t *testing.T) {
	// Test with embedded struct
	type ConfigWithIdentityParams struct {
		Name string
		PluginIdentityTokenParams
	}

	testcases := []struct {
		name     string
		c1       ConfigWithIdentityParams
		c2       ConfigWithIdentityParams
		expected bool
	}{
		{
			name: "embedded-equal",
			c1: ConfigWithIdentityParams{
				Name: "config1",
				PluginIdentityTokenParams: PluginIdentityTokenParams{
					IdentityTokenTTL:      10 * time.Second,
					IdentityTokenAudience: "test-aud",
				},
			},
			c2: ConfigWithIdentityParams{
				Name: "config2", // Different name, but we only compare PluginIdentityTokenParams
				PluginIdentityTokenParams: PluginIdentityTokenParams{
					IdentityTokenTTL:      10 * time.Second,
					IdentityTokenAudience: "test-aud",
				},
			},
			expected: true,
		},
		{
			name: "embedded-different",
			c1: ConfigWithIdentityParams{
				Name: "config1",
				PluginIdentityTokenParams: PluginIdentityTokenParams{
					IdentityTokenTTL:      10 * time.Second,
					IdentityTokenAudience: "test-aud",
				},
			},
			c2: ConfigWithIdentityParams{
				Name: "config1",
				PluginIdentityTokenParams: PluginIdentityTokenParams{
					IdentityTokenTTL:      20 * time.Second,
					IdentityTokenAudience: "test-aud",
				},
			},
			expected: false,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			// Test comparing the embedded fields directly
			result := tt.c1.PluginIdentityTokenParams.Equals(tt.c2.PluginIdentityTokenParams)
			assert.Equal(t, tt.expected, result)

			// Test using method promotion
			result2 := tt.c1.Equals(tt.c2.PluginIdentityTokenParams)
			assert.Equal(t, tt.expected, result2)
		})
	}
}
