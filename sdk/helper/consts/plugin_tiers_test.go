// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consts

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParsePluginTier(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		pluginTier string
		want       PluginTier
		wantErr    bool
	}{
		{
			name:       "unknown",
			pluginTier: "unknown",
			want:       PluginTierUnknown,
		},
		{
			name:       "empty unknown",
			pluginTier: "",
			want:       PluginTierUnknown,
		},
		{
			name:       "community",
			pluginTier: "community",
			want:       PluginTierCommunity,
		},
		{
			name:       "partner",
			pluginTier: "partner",
			want:       PluginTierPartner,
		},
		{
			name:       "official",
			pluginTier: "official",
			want:       PluginTierOfficial,
		},
		{
			name:       "unsupported",
			pluginTier: "foo",
			want:       PluginTierUnknown,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePluginTier(tt.pluginTier)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParsePluginTier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Fatalf("ParsePluginTier() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPluginTier_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		p    PluginTier
		want []byte
	}{
		{
			name: "unknown",
			p:    PluginTierUnknown,
			want: []byte(`"unknown"`),
		},
		{
			name: "community",
			p:    PluginTierCommunity,
			want: []byte(`"community"`),
		},
		{
			name: "partner",
			p:    PluginTierPartner,
			want: []byte(`"partner"`),
		},
		{
			name: "offical",
			p:    PluginTierOfficial,
			want: []byte(`"official"`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON() error = %v, want nil", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPluginTier_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		wantTier PluginTier
		data     []byte
		wantErr  bool
	}{
		{
			name:     "unknown",
			wantTier: PluginTierUnknown,
			data:     []byte(`"unknown"`),
		},
		{
			name:     "community",
			wantTier: PluginTierCommunity,
			data:     []byte(`"community"`),
		},
		{
			name:     "partner",
			wantTier: PluginTierPartner,
			data:     []byte(`"partner"`),
		},
		{
			name:     "offical",
			wantTier: PluginTierOfficial,
			data:     []byte(`"official"`),
		},
		{
			name:     "unsupported",
			wantTier: PluginTierUnknown,
			data:     []byte(`"foo"`),
			wantErr:  true,
		},
		// The following test cases ensures new plugin tiers are added at the end of the enum list
		// Inserting a new plugin tier in the middle of the enum list will fail
		// New plugin tiers should be added at the end of the test case list
		{
			name:     "0-unknown",
			wantTier: PluginTierUnknown,
			data:     []byte(`0`),
		},
		{
			name:     "1-community",
			wantTier: PluginTierCommunity,
			data:     []byte(`1`),
		},
		{
			name:     "2-partner",
			wantTier: PluginTierPartner,
			data:     []byte(`2`),
		},
		{
			name:     "3-official",
			wantTier: PluginTierOfficial,
			data:     []byte(`3`),
		},
		{
			name:    "tier number unsupported",
			data:    []byte(`2345678`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tier PluginTier
			err := tier.UnmarshalJSON(tt.data)
			if (err != nil) != tt.wantErr {
				t.Fatalf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tier != tt.wantTier {
				t.Fatalf("UnmarshalJSON() got = %v, want %v", tier, tt.wantTier)
			}
		})
	}
}

// TestPluginTierJSONRoundTrip tests that PluginTier can be marshaled and unmarshaled
// to/from JSON in a round trip.
func TestPluginTierJSONRoundTrip(t *testing.T) {
	type testTier struct {
		PluginTier PluginTier `json:"plugin_tier"`
	}

	for _, tier := range PluginTiers {
		t.Run(tier.String(), func(t *testing.T) {
			original := testTier{
				PluginTier: tier,
			}
			asBytes, err := json.Marshal(original)
			if err != nil {
				t.Fatal(err)
			}

			var roundTripped testTier
			err = json.Unmarshal(asBytes, &roundTripped)
			if err != nil {
				t.Fatal(err)
			}

			if original != roundTripped {
				t.Fatalf("expected %v, got %v", original, roundTripped)
			}
		})
	}
}

func TestUnknownTierExcludedWithOmitEmpty(t *testing.T) {
	type testTierOmitEmpty struct {
		Type PluginTier `json:"tier,omitempty"`
	}
	bytes, err := json.Marshal(testTierOmitEmpty{})
	if err != nil {
		t.Fatal(err)
	}
	m := map[string]any{}
	json.Unmarshal(bytes, &m)
	if _, exists := m["tier"]; exists {
		t.Fatal("tier should not be present")
	}
}
