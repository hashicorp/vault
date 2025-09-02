// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

// NOTE: this file was copied from
// https://github.com/hashicorp/vault/blob/main/sdk/helper/consts/plugin_types_test.go
// Any changes made should be made to both files at the same time.

import (
	"encoding/json"
	"testing"
)

type testType struct {
	PluginType PluginType `json:"plugin_type"`
}

func TestPluginTypeJSONRoundTrip(t *testing.T) {
	for _, pluginType := range PluginTypes {
		original := testType{
			PluginType: pluginType,
		}
		asBytes, err := json.Marshal(original)
		if err != nil {
			t.Fatal(err)
		}

		var roundTripped testType
		err = json.Unmarshal(asBytes, &roundTripped)
		if err != nil {
			t.Fatal(err)
		}

		if original != roundTripped {
			t.Fatalf("expected %v, got %v", original, roundTripped)
		}
	}
}

func TestPluginTypeJSONUnmarshal(t *testing.T) {
	// Failure/unsupported cases.
	for name, tc := range map[string]string{
		"unsupported":   `{"plugin_type":"unsupported"}`,
		"random string": `{"plugin_type":"foo"}`,
		"boolean":       `{"plugin_type":true}`,
		"empty":         `{"plugin_type":""}`,
		"negative":      `{"plugin_type":-1}`,
		"out of range":  `{"plugin_type":10}`,
	} {
		t.Run(name, func(t *testing.T) {
			var result testType
			err := json.Unmarshal([]byte(tc), &result)
			if err == nil {
				t.Fatal("expected error")
			}
		})
	}

	// Valid cases.
	for name, tc := range map[string]struct {
		json     string
		expected PluginType
	}{
		"unknown":         {`{"plugin_type":"unknown"}`, PluginTypeUnknown},
		"auth":            {`{"plugin_type":"auth"}`, PluginTypeCredential},
		"secret":          {`{"plugin_type":"secret"}`, PluginTypeSecrets},
		"database":        {`{"plugin_type":"database"}`, PluginTypeDatabase},
		"absent":          {`{}`, PluginTypeUnknown},
		"integer unknown": {`{"plugin_type":0}`, PluginTypeUnknown},
		"integer auth":    {`{"plugin_type":1}`, PluginTypeCredential},
		"integer db":      {`{"plugin_type":2}`, PluginTypeDatabase},
		"integer secret":  {`{"plugin_type":3}`, PluginTypeSecrets},
	} {
		t.Run(name, func(t *testing.T) {
			var result testType
			err := json.Unmarshal([]byte(tc.json), &result)
			if err != nil {
				t.Fatal(err)
			}
			if tc.expected != result.PluginType {
				t.Fatalf("expected %v, got %v", tc.expected, result.PluginType)
			}
		})
	}
}

func TestUnknownTypeExcludedWithOmitEmpty(t *testing.T) {
	type testTypeOmitEmpty struct {
		Type PluginType `json:"type,omitempty"`
	}
	bytes, err := json.Marshal(testTypeOmitEmpty{})
	if err != nil {
		t.Fatal(err)
	}
	m := map[string]any{}
	json.Unmarshal(bytes, &m)
	if _, exists := m["type"]; exists {
		t.Fatal("type should not be present")
	}
}
