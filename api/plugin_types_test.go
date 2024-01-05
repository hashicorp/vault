// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

// NOTE: this file was copied from
// https://github.com/hashicorp/vault/blob/main/sdk/helper/consts/plugin_types_test.go
// Any changes made should be made to both files at the same time.

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		require.NoError(t, err)

		var roundTripped testType
		require.NoError(t, json.Unmarshal(asBytes, &roundTripped))

		assert.Equal(t, original, roundTripped)
	}
}

func TestPluginTypeJSONUnmarshal(t *testing.T) {
	// Failure/unsupported cases.
	for _, tc := range map[string]string{
		"unsupported":   `{"plugin_type":"unsupported"}`,
		"random string": `{"plugin_type":"foo"}`,
		"integer":       `{"plugin_type":0}`,
		"boolean":       `{"plugin_type":true}`,
		"empty":         `{"plugin_type":""}`,
	} {
		var result testType
		t.Log(result)
		err := json.Unmarshal([]byte(tc), &result)
		assert.Error(t, err)
		t.Log(result)
	}

	// Valid cases.
	for name, tc := range map[string]struct {
		json     string
		expected PluginType
	}{
		"unknown":  {`{"plugin_type":"unknown"}`, PluginTypeUnknown},
		"auth":     {`{"plugin_type":"auth"}`, PluginTypeCredential},
		"secret":   {`{"plugin_type":"secret"}`, PluginTypeSecrets},
		"database": {`{"plugin_type":"database"}`, PluginTypeDatabase},
		"absent":   {`{}`, PluginTypeUnknown},
	} {
		t.Run(name, func(t *testing.T) {
			var result testType
			assert.NoError(t, json.Unmarshal([]byte(tc.json), &result))
			assert.Equal(t, tc.expected, result.PluginType)
		})
	}
}
