// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package syslog

import (
	"context"
	"testing"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// TestBackend_formatterConfig ensures that all the configuration values are parsed correctly.
func TestBackend_formatterConfig(t *testing.T) {
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
			expectedErrMsg: "audit.NewFormatterConfig: error applying options: audit.(format).validate: 'squiggly' is not a valid format: invalid parameter",
		},
		"invalid-hmac-accessor": {
			config: map[string]string{
				"format":        audit.JSONFormat.String(),
				"hmac_accessor": "maybe",
			},
			want:           audit.FormatterConfig{},
			wantErr:        true,
			expectedErrMsg: "syslog.formatterConfig: unable to parse 'hmac_accessor': strconv.ParseBool: parsing \"maybe\": invalid syntax",
		},
		"invalid-log-raw": {
			config: map[string]string{
				"format":        audit.JSONFormat.String(),
				"hmac_accessor": "true",
				"log_raw":       "maybe",
			},
			want:           audit.FormatterConfig{},
			wantErr:        true,
			expectedErrMsg: "syslog.formatterConfig: unable to parse 'log_raw': strconv.ParseBool: parsing \"maybe\": invalid syntax",
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
			expectedErrMsg: "syslog.formatterConfig: unable to parse 'elide_list_responses': strconv.ParseBool: parsing \"maybe\": invalid syntax",
		},
	}
	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := formatterConfig(tc.config)
			if tc.wantErr {
				require.Error(t, err)
				require.EqualError(t, err, tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.want, got)
		})
	}
}

// TestBackend_configureFilterNode ensures that configureFilterNode handles various
// filter values as expected. Empty (including whitespace) strings should return
// no error but skip configuration of the node.
func TestBackend_configureFilterNode(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		filter           string
		shouldSkipNode   bool
		wantErr          bool
		expectedErrorMsg string
	}{
		"happy": {
			filter: "foo == bar",
		},
		"empty": {
			filter:         "",
			shouldSkipNode: true,
		},
		"spacey": {
			filter:         "    ",
			shouldSkipNode: true,
		},
		"bad": {
			filter:           "___qwerty",
			wantErr:          true,
			expectedErrorMsg: "syslog.(Backend).configureFilterNode: error creating filter node: audit.NewEntryFilter: cannot create new audit filter",
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

			err := b.configureFilterNode(tc.filter)

			switch {
			case tc.wantErr:
				require.Error(t, err)
				require.ErrorContains(t, err, tc.expectedErrorMsg)
				require.Len(t, b.nodeIDList, 0)
				require.Len(t, b.nodeMap, 0)
			case tc.shouldSkipNode:
				require.NoError(t, err)
				require.Len(t, b.nodeIDList, 0)
				require.Len(t, b.nodeMap, 0)
			default:
				require.NoError(t, err)
				require.Len(t, b.nodeIDList, 1)
				require.Len(t, b.nodeMap, 1)
				id := b.nodeIDList[0]
				node := b.nodeMap[id]
				require.Equal(t, eventlogger.NodeTypeFilter, node.Type())
			}
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

	formatConfig, err := audit.NewFormatterConfig()
	require.NoError(t, err)

	err = b.configureFormatterNode(formatConfig)

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
			expectedErrMsg: "syslog.(Backend).configureSinkNode: name is required: invalid parameter",
		},
		"name-whitespace": {
			name:           "   ",
			wantErr:        true,
			expectedErrMsg: "syslog.(Backend).configureSinkNode: name is required: invalid parameter",
		},
		"format-empty": {
			name:           "foo",
			format:         "",
			wantErr:        true,
			expectedErrMsg: "syslog.(Backend).configureSinkNode: format is required: invalid parameter",
		},
		"format-whitespace": {
			name:           "foo",
			format:         "   ",
			wantErr:        true,
			expectedErrMsg: "syslog.(Backend).configureSinkNode: format is required: invalid parameter",
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
				sw, ok := node.(*audit.SinkWrapper)
				require.True(t, ok)
				require.Equal(t, tc.expectedName, sw.Name)
			}
		})
	}
}

// TestBackend_configureFilterFormatterSink ensures that configuring all three
// types of nodes on a Backend works as expected, i.e. we have all three nodes
// at the end and nothing gets overwritten. The order of calls influences the
// slice of IDs on the Backend.
func TestBackend_configureFilterFormatterSink(t *testing.T) {
	t.Parallel()

	b := &Backend{
		nodeIDList: []eventlogger.NodeID{},
		nodeMap:    map[eventlogger.NodeID]eventlogger.Node{},
	}

	formatConfig, err := audit.NewFormatterConfig()
	require.NoError(t, err)

	err = b.configureFilterNode("foo == bar")
	require.NoError(t, err)

	err = b.configureFormatterNode(formatConfig)
	require.NoError(t, err)

	err = b.configureSinkNode("foo", "json")
	require.NoError(t, err)

	require.Len(t, b.nodeIDList, 3)
	require.Len(t, b.nodeMap, 3)

	id := b.nodeIDList[0]
	node := b.nodeMap[id]
	require.Equal(t, eventlogger.NodeTypeFilter, node.Type())

	id = b.nodeIDList[1]
	node = b.nodeMap[id]
	require.Equal(t, eventlogger.NodeTypeFormatter, node.Type())

	id = b.nodeIDList[2]
	node = b.nodeMap[id]
	require.Equal(t, eventlogger.NodeTypeSink, node.Type())
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
			expectedErrorMessage: "syslog.Factory: nil salt config",
		},
		"nil-salt-view": {
			backendConfig: &audit.BackendConfig{
				SaltConfig: &salt.Config{},
			},
			isErrorExpected:      true,
			expectedErrorMessage: "syslog.Factory: nil salt view",
		},
		"non-fallback-device-with-filter": {
			backendConfig: &audit.BackendConfig{
				MountPath:  "discard",
				SaltConfig: &salt.Config{},
				SaltView:   &logical.InmemStorage{},
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
				Config: map[string]string{
					"fallback": "true",
					"filter":   "mount_type == kv",
				},
			},
			isErrorExpected:      true,
			expectedErrorMessage: "syslog.Factory: cannot configure a fallback device with a filter: invalid parameter",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			be, err := Factory(ctx, tc.backendConfig, true, nil)

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

			be, err := Factory(ctx, tc.backendConfig, true, nil)
			require.NoError(t, err)
			require.NotNil(t, be)
			require.Equal(t, tc.isFallbackExpected, be.IsFallback())
		})
	}
}
