// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package syslog

import (
	"context"
	"testing"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/internal/observability/event"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// TestBackend_newFormatterConfig ensures that all the configuration values are parsed correctly.
func TestBackend_newFormatterConfig(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		config         map[string]string
		want           audit.FormatterConfig
		wantErr        bool
		expectedErrMsg string
	}{
		"happy-path-json": {
			config: map[string]string{
				"format":               audit.JSONFormat.String(),
				"hmac_accessor":        "true",
				"log_raw":              "true",
				"elide_list_responses": "true",
			},
			want: audit.FormatterConfig{
				Raw:                true,
				HMACAccessor:       true,
				ElideListResponses: true,
				RequiredFormat:     "json",
			}, wantErr: false,
		},
		"happy-path-jsonx": {
			config: map[string]string{
				"format":               audit.JSONxFormat.String(),
				"hmac_accessor":        "true",
				"log_raw":              "true",
				"elide_list_responses": "true",
			},
			want: audit.FormatterConfig{
				Raw:                true,
				HMACAccessor:       true,
				ElideListResponses: true,
				RequiredFormat:     "jsonx",
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
			want:           audit.FormatterConfig{},
			wantErr:        true,
			expectedErrMsg: "error applying options: invalid format \"squiggly\": invalid parameter",
		},
		"invalid-hmac-accessor": {
			config: map[string]string{
				"format":        audit.JSONFormat.String(),
				"hmac_accessor": "maybe",
			},
			want:           audit.FormatterConfig{},
			wantErr:        true,
			expectedErrMsg: "unable to parse 'hmac_accessor': invalid configuration: strconv.ParseBool: parsing \"maybe\": invalid syntax",
		},
		"invalid-log-raw": {
			config: map[string]string{
				"format":        audit.JSONFormat.String(),
				"hmac_accessor": "true",
				"log_raw":       "maybe",
			},
			want:           audit.FormatterConfig{},
			wantErr:        true,
			expectedErrMsg: "unable to parse 'log_raw: invalid configuration: strconv.ParseBool: parsing \"maybe\": invalid syntax",
		},
		"invalid-elide-bool": {
			config: map[string]string{
				"format":               audit.JSONFormat.String(),
				"hmac_accessor":        "true",
				"log_raw":              "true",
				"elide_list_responses": "maybe",
			},
			want:           audit.FormatterConfig{},
			wantErr:        true,
			expectedErrMsg: "unable to parse 'elide_list_responses': invalid configuration: strconv.ParseBool: parsing \"maybe\": invalid syntax",
		},
		"prefix": {
			config: map[string]string{
				"format": audit.JSONFormat.String(),
				"prefix": "foo",
			},
			want: audit.FormatterConfig{
				RequiredFormat: audit.JSONFormat,
				Prefix:         "foo",
				HMACAccessor:   true,
			},
		},
	}
	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := newFormatterConfig(&corehelpers.NoopHeaderFormatter{}, tc.config)
			if tc.wantErr {
				require.Error(t, err)
				require.EqualError(t, err, tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.want.RequiredFormat, got.RequiredFormat)
			require.Equal(t, tc.want.Raw, got.Raw)
			require.Equal(t, tc.want.ElideListResponses, got.ElideListResponses)
			require.Equal(t, tc.want.HMACAccessor, got.HMACAccessor)
			require.Equal(t, tc.want.OmitTime, got.OmitTime)
			require.Equal(t, tc.want.Prefix, got.Prefix)
		})
	}
}

// TestBackend_configureFormatterNode ensures that configureFormatterNode
// populates the nodeIDList and nodeMap on Backend when given valid formatConfig.
func TestBackend_configureFormatterNode(t *testing.T) {
	t.Parallel()

	b := &Backend{
		nodeIDList: []eventlogger.NodeID{},
		nodeMap:    map[eventlogger.NodeID]eventlogger.Node{},
	}

	formatConfig, err := audit.NewFormatterConfig(&corehelpers.NoopHeaderFormatter{})
	require.NoError(t, err)

	err = b.configureFormatterNode("juan", formatConfig, hclog.NewNullLogger())

	require.NoError(t, err)
	require.Len(t, b.nodeIDList, 1)
	require.Len(t, b.nodeMap, 1)
	id := b.nodeIDList[0]
	node := b.nodeMap[id]
	require.Equal(t, eventlogger.NodeTypeFormatter, node.Type())
}

// TestBackend_configureSinkNode ensures that we can correctly configure the sink
// node on the Backend, and any incorrect parameters result in the relevant errors.
func TestBackend_configureSinkNode(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		name           string
		format         string
		wantErr        bool
		expectedErrMsg string
		expectedName   string
	}{
		"name-empty": {
			name:           "",
			wantErr:        true,
			expectedErrMsg: "name is required: invalid parameter",
		},
		"name-whitespace": {
			name:           "   ",
			wantErr:        true,
			expectedErrMsg: "name is required: invalid parameter",
		},
		"format-empty": {
			name:           "foo",
			format:         "",
			wantErr:        true,
			expectedErrMsg: "format is required: invalid parameter",
		},
		"format-whitespace": {
			name:           "foo",
			format:         "   ",
			wantErr:        true,
			expectedErrMsg: "format is required: invalid parameter",
		},
		"happy": {
			name:         "foo",
			format:       "json",
			wantErr:      false,
			expectedName: "foo",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			b := &Backend{
				nodeIDList: []eventlogger.NodeID{},
				nodeMap:    map[eventlogger.NodeID]eventlogger.Node{},
			}

			err := b.configureSinkNode(tc.name, tc.format)

			if tc.wantErr {
				require.Error(t, err)
				require.EqualError(t, err, tc.expectedErrMsg)
				require.Len(t, b.nodeIDList, 0)
				require.Len(t, b.nodeMap, 0)
			} else {
				require.NoError(t, err)
				require.Len(t, b.nodeIDList, 1)
				require.Len(t, b.nodeMap, 1)
				id := b.nodeIDList[0]
				node := b.nodeMap[id]
				require.Equal(t, eventlogger.NodeTypeSink, node.Type())
				mc, ok := node.(*event.MetricsCounter)
				require.True(t, ok)
				require.Equal(t, tc.expectedName, mc.Name)
			}
		})
	}
}

// TestBackend_Factory_Conf is used to ensure that any configuration which is
// supplied, is validated and tested.
func TestBackend_Factory_Conf(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := map[string]struct {
		backendConfig        *audit.BackendConfig
		isErrorExpected      bool
		expectedErrorMessage string
	}{
		"nil-salt-config": {
			backendConfig: &audit.BackendConfig{
				SaltConfig: nil,
			},
			isErrorExpected:      true,
			expectedErrorMessage: "nil salt config: invalid parameter",
		},
		"nil-salt-view": {
			backendConfig: &audit.BackendConfig{
				SaltConfig: &salt.Config{},
			},
			isErrorExpected:      true,
			expectedErrorMessage: "nil salt view: invalid parameter",
		},
		"non-fallback-device-with-filter": {
			backendConfig: &audit.BackendConfig{
				MountPath:  "discard",
				SaltConfig: &salt.Config{},
				SaltView:   &logical.InmemStorage{},
				Logger:     hclog.NewNullLogger(),
				Config: map[string]string{
					"fallback": "false",
					"filter":   "mount_type == kv",
				},
			},
			isErrorExpected: false,
		},
		"fallback-device-with-filter": {
			backendConfig: &audit.BackendConfig{
				MountPath:  "discard",
				SaltConfig: &salt.Config{},
				SaltView:   &logical.InmemStorage{},
				Logger:     hclog.NewNullLogger(),
				Config: map[string]string{
					"fallback": "true",
					"filter":   "mount_type == kv",
				},
			},
			isErrorExpected:      true,
			expectedErrorMessage: "cannot configure a fallback device with a filter: invalid parameter",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			be, err := Factory(ctx, tc.backendConfig, &corehelpers.NoopHeaderFormatter{})

			switch {
			case tc.isErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.expectedErrorMessage)
			default:
				require.NoError(t, err)
				require.NotNil(t, be)
			}
		})
	}
}

// TestBackend_IsFallback ensures that the 'fallback' config setting is parsed
// and set correctly, then exposed via the interface method IsFallback().
func TestBackend_IsFallback(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := map[string]struct {
		backendConfig      *audit.BackendConfig
		isFallbackExpected bool
	}{
		"fallback": {
			backendConfig: &audit.BackendConfig{
				MountPath:  "qwerty",
				SaltConfig: &salt.Config{},
				SaltView:   &logical.InmemStorage{},
				Logger:     hclog.NewNullLogger(),
				Config: map[string]string{
					"fallback": "true",
				},
			},
			isFallbackExpected: true,
		},
		"no-fallback": {
			backendConfig: &audit.BackendConfig{
				MountPath:  "qwerty",
				SaltConfig: &salt.Config{},
				SaltView:   &logical.InmemStorage{},
				Logger:     hclog.NewNullLogger(),
				Config: map[string]string{
					"fallback": "false",
				},
			},
			isFallbackExpected: false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			be, err := Factory(ctx, tc.backendConfig, &corehelpers.NoopHeaderFormatter{})
			require.NoError(t, err)
			require.NotNil(t, be)
			require.Equal(t, tc.isFallbackExpected, be.IsFallback())
		})
	}
}
