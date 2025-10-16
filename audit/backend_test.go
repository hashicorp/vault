// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"testing"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/stretchr/testify/require"
)

// testAuditEntryMinimal represents a minimal version of a MountEntry which would be
// required by the audit system in order to compare two entries and establish if
// they are distinct. It is declared here in testing, because 'auditEntryMinimal'
// is not exported, and 'MountEntry' exists in the 'vault' package and would cause
// an import cycle.
type testAuditEntryMinimal struct {
	path       string
	deviceType string
	options    map[string]string
}

// Type returns the device type of the entry.
func (m *testAuditEntryMinimal) Type() string {
	return m.deviceType
}

// Path returns the device mount path of the entry.
func (m *testAuditEntryMinimal) Path() string {
	return m.path
}

// Options returns the specified option given a key.
func (m *testAuditEntryMinimal) Options(key string) string {
	return m.options[key]
}

// TestBackend_newFormatterConfig ensures that all the configuration values are
// parsed correctly when trying to create a new formatterConfig via newFormatterConfig.
func TestBackend_newFormatterConfig(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		config          map[string]string
		want            formatterConfig
		wantErr         bool
		expectedMessage string
	}{
		"happy-path-json": {
			config: map[string]string{
				"format":               jsonFormat.String(),
				"hmac_accessor":        "true",
				"log_raw":              "true",
				"elide_list_responses": "true",
			},
			want: formatterConfig{
				raw:                true,
				hmacAccessor:       true,
				elideListResponses: true,
				requiredFormat:     "json",
			}, wantErr: false,
		},
		"happy-path-jsonx": {
			config: map[string]string{
				"format":               jsonxFormat.String(),
				"hmac_accessor":        "true",
				"log_raw":              "true",
				"elide_list_responses": "true",
			},
			want: formatterConfig{
				raw:                true,
				hmacAccessor:       true,
				elideListResponses: true,
				requiredFormat:     "jsonx",
			},
			wantErr: false,
		},
		"invalid-format": {
			config: map[string]string{
				"format":               " squiggly ",
				"hmac_accessor":        "true",
				"log_raw":              "true",
				"elide_list_responses": "true",
			},
			want:            formatterConfig{},
			wantErr:         true,
			expectedMessage: "unsupported \"format\": invalid configuration",
		},
		"invalid-hmac-accessor": {
			config: map[string]string{
				"format":        jsonFormat.String(),
				"hmac_accessor": "maybe",
			},
			want:            formatterConfig{},
			wantErr:         true,
			expectedMessage: "unable to parse \"hmac_accessor\": invalid configuration",
		},
		"invalid-log-raw": {
			config: map[string]string{
				"format":        jsonFormat.String(),
				"hmac_accessor": "true",
				"log_raw":       "maybe",
			},
			want:            formatterConfig{},
			wantErr:         true,
			expectedMessage: "unable to parse \"log_raw\": invalid configuration",
		},
		"invalid-elide-bool": {
			config: map[string]string{
				"format":               jsonFormat.String(),
				"hmac_accessor":        "true",
				"log_raw":              "true",
				"elide_list_responses": "maybe",
			},
			want:            formatterConfig{},
			wantErr:         true,
			expectedMessage: "unable to parse \"elide_list_responses\": invalid configuration",
		},
		"prefix": {
			config: map[string]string{
				"format": jsonFormat.String(),
				"prefix": "foo",
			},
			want: formatterConfig{
				requiredFormat: jsonFormat,
				prefix:         "foo",
				hmacAccessor:   true,
			},
		},
	}
	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := newFormatterConfig(&noopHeaderFormatter{}, tc.config)
			if tc.wantErr {
				require.Error(t, err)
				require.EqualError(t, err, tc.expectedMessage)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.want.requiredFormat, got.requiredFormat)
			require.Equal(t, tc.want.raw, got.raw)
			require.Equal(t, tc.want.elideListResponses, got.elideListResponses)
			require.Equal(t, tc.want.hmacAccessor, got.hmacAccessor)
			require.Equal(t, tc.want.omitTime, got.omitTime)
			require.Equal(t, tc.want.prefix, got.prefix)
		})
	}
}

// TestBackend_configureFormatterNode ensures that configureFormatterNode
// populates the nodeIDList and nodeMap on backend when given valid config.
func TestBackend_configureFormatterNode(t *testing.T) {
	t.Parallel()

	b, err := newBackend(&noopHeaderFormatter{}, &BackendConfig{
		MountPath: "foo",
		Logger:    hclog.NewNullLogger(),
	})
	require.NoError(t, err)

	require.Len(t, b.nodeIDList, 1)
	require.Len(t, b.nodeMap, 1)
	id := b.nodeIDList[0]
	node := b.nodeMap[id]
	require.Equal(t, eventlogger.NodeTypeFormatter, node.Type())
}

// TestBackend_hasEnterpriseAuditOptions checks that the existence of any Enterprise
// only options in the options which can be supplied to enable an audit device can
// be flagged.
func TestBackend_hasEnterpriseAuditOptions(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input    map[string]string
		expected bool
	}{
		"nil": {
			expected: false,
		},
		"empty": {
			input:    make(map[string]string),
			expected: false,
		},
		"non-ent-opts": {
			input: map[string]string{
				"log_raw": "true",
			},
			expected: false,
		},
		"ent-opt-filter": {
			input: map[string]string{
				"filter": "mount_type == kv",
			},
			expected: true,
		},
		"ent-opt-fallback": {
			input: map[string]string{
				"fallback": "true",
			},
			expected: true,
		},
		"ent-opt-filter-and-fallback": {
			input: map[string]string{
				"filter":   "mount_type == kv",
				"fallback": "true",
			},
			expected: true,
		},
		"ent-opt-exclude": {
			input: map[string]string{
				"exclude": `{
					"condition": "\"/request/mount_type\" == transit",
					"fields": [ "/request/data", "/response/data" ]
				}`,
			},
			expected: true,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.expected, hasEnterpriseAuditOptions(tc.input))
		})
	}
}

// TestBackend_hasInvalidAuditOptions tests that depending on whether we are running
// an Enterprise or non-Enterprise version of Vault, the options supplied to enable
// an audit device may or may not be valid.
// NOTE: In the non-Enterprise version of Vault supplying audit options such as
// 'filter' or 'fallback' is not allowed.
func TestBackend_hasInvalidAuditOptions(t *testing.T) {
	tests := map[string]struct {
		input    map[string]string
		expected bool
	}{
		"non-ent-opts": {
			input: map[string]string{
				"log_raw": "true",
			},
			expected: false,
		},
		"ent-opt": {
			input: map[string]string{
				"filter": "mount_type == kv",
			},
			expected: !constants.IsEnterprise,
		},
		"ent-opt-filter": {
			input: map[string]string{
				"filter": "mount_type == kv",
			},
			expected: !constants.IsEnterprise,
		},
		"ent-opt-fallback": {
			input: map[string]string{
				"fallback": "true",
			},
			expected: !constants.IsEnterprise,
		},
		"ent-opt-filter-and-fallback": {
			input: map[string]string{
				"filter":   "mount_type == kv",
				"fallback": "true",
			},
			expected: !constants.IsEnterprise,
		},
		"ent-opt-exclude": {
			input: map[string]string{
				"exclude": `{
					"condition": "\"/request/mount_type\" == transit",
					"fields": [ "/request/data", "/response/data" ]
				}`,
			},
			expected: !constants.IsEnterprise,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.expected, HasInvalidOptions(tc.input))
		})
	}
}

// TestBackend_IsDistinct tests a variety of scenarios for deciding whether a proposed
// audit backend is distinct from an existing one (basically just comparing some
// properties of two backends).
func TestBackend_IsDistinct(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		isDistinctExpected   bool
		expectedErrorMessage string
		proposed             backendComparer
		existing             backendComparer
	}{
		"different-path-and-file-path": {
			isDistinctExpected: true,
			proposed: &testAuditEntryMinimal{
				path:       "foo1/",
				deviceType: "file",
				options:    map[string]string{"file_path": "/var/vault/audit.log"},
			},
			existing: &testAuditEntryMinimal{
				path:       "foo2/",
				deviceType: "file",
				options:    map[string]string{"file_path": "/var/vault/audit2.log"},
			},
		},
		"different-path-same-file-path": {
			isDistinctExpected:   false,
			expectedErrorMessage: "\"file_path\" \"/var/vault/audit.log\" already in use on device: \"foo2/\": invalid configuration",
			proposed: &testAuditEntryMinimal{
				path:       "foo1/",
				deviceType: "file",
				options:    map[string]string{"file_path": "/var/vault/audit.log"},
			},
			existing: &testAuditEntryMinimal{
				path:       "foo2/",
				deviceType: "file",
				options:    map[string]string{"file_path": "/var/vault/audit.log"},
			},
		},
		"same-path-different-file-path": {
			isDistinctExpected:   false,
			expectedErrorMessage: "path \"foo1/\" already in use: invalid configuration",
			proposed: &testAuditEntryMinimal{
				path:       "foo1/",
				deviceType: "file",
				options:    map[string]string{"file_path": "/var/vault/audit2.log"},
			},
			existing: &testAuditEntryMinimal{
				path:       "foo1/",
				deviceType: "file",
				options:    map[string]string{"file_path": "/var/vault/audit.log"},
			},
		},
		"bad-type": {
			isDistinctExpected:   false,
			expectedErrorMessage: "unexpected 'type' for audit: \"juan\": invalid internal parameter",
			proposed: &testAuditEntryMinimal{
				path:       "foo2/",
				deviceType: "juan",
				options:    map[string]string{"file_path": "/var/vault/audit2.log"},
			},
			existing: &testAuditEntryMinimal{
				path:       "foo1/",
				deviceType: "file",
				options:    map[string]string{"file_path": "/var/vault/audit.log"},
			},
		},
		// These tests feel odd, but it's what the code was doing before any move
		// to eventlogger:
		// https://github.com/hashicorp/vault/blob/release/1.14.x/vault/audit.go#L81-L86
		"proposed-has-prefix-of-existing": {
			isDistinctExpected:   false,
			expectedErrorMessage: "path \"foo1/\" already in use: invalid configuration",
			proposed: &testAuditEntryMinimal{
				path:       "foo1/bar/",
				deviceType: "juan",
				options:    map[string]string{"file_path": "/var/vault/audit2.log"},
			},
			existing: &testAuditEntryMinimal{
				path:       "foo1/",
				deviceType: "file",
				options:    map[string]string{"file_path": "/var/vault/audit.log"},
			},
		},
		"existing-has-prefix-of-proposed": {
			isDistinctExpected:   false,
			expectedErrorMessage: "path \"foo1/bar\" already in use: invalid configuration",
			proposed: &testAuditEntryMinimal{
				path:       "foo1/",
				deviceType: "juan",
				options:    map[string]string{"file_path": "/var/vault/audit2.log"},
			},
			existing: &testAuditEntryMinimal{
				path:       "foo1/bar",
				deviceType: "file",
				options:    map[string]string{"file_path": "/var/vault/audit.log"},
			},
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ok, err := IsDistinct(tc.proposed, tc.existing)
			require.Equal(t, tc.isDistinctExpected, ok)
			if tc.isDistinctExpected {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.EqualError(t, err, tc.expectedErrorMessage)
			}
		})
	}
}
