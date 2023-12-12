// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL

package syslog

import (
	"testing"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/vault/audit"
	"github.com/stretchr/testify/require"
)

// TestBackend_formatterConfig ensures that all the configuration values are parsed correctly.
func TestBackend_formatterConfig(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		config          map[string]string
		want            audit.FormatterConfig
		wantErr         bool
		expectedMessage string
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
			want:            audit.FormatterConfig{},
			wantErr:         true,
			expectedMessage: "audit.NewFormatterConfig: error applying options: audit.(format).validate: 'squiggly' is not a valid format: invalid parameter",
		},
		"invalid-hmac-accessor": {
			config: map[string]string{
				"format":        audit.JSONFormat.String(),
				"hmac_accessor": "maybe",
			},
			want:            audit.FormatterConfig{},
			wantErr:         true,
			expectedMessage: "syslog.formatterConfig: unable to parse 'hmac_accessor': strconv.ParseBool: parsing \"maybe\": invalid syntax",
		},
		"invalid-log-raw": {
			config: map[string]string{
				"format":        audit.JSONFormat.String(),
				"hmac_accessor": "true",
				"log_raw":       "maybe",
			},
			want:            audit.FormatterConfig{},
			wantErr:         true,
			expectedMessage: "syslog.formatterConfig: unable to parse 'log_raw': strconv.ParseBool: parsing \"maybe\": invalid syntax",
		},
		"invalid-elide-bool": {
			config: map[string]string{
				"format":               audit.JSONFormat.String(),
				"hmac_accessor":        "true",
				"log_raw":              "true",
				"elide_list_responses": "maybe",
			},
			want:            audit.FormatterConfig{},
			wantErr:         true,
			expectedMessage: "syslog.formatterConfig: unable to parse 'elide_list_responses': strconv.ParseBool: parsing \"maybe\": invalid syntax",
		},
	}
	for name, testCase := range tests {
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := formatterConfig(testCase.config)
			if testCase.wantErr {
				require.Error(t, err)
				require.EqualError(t, err, testCase.expectedMessage)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, testCase.want, got)
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
	for name, testCase := range tests {
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			b := &Backend{
				nodeIDList: []eventlogger.NodeID{},
				nodeMap:    map[eventlogger.NodeID]eventlogger.Node{},
			}

			err := b.configureFilterNode(testCase.filter)

			switch {
			case testCase.wantErr:
				require.Error(t, err)
				require.ErrorContains(t, err, testCase.expectedErrorMsg)
				require.Len(t, b.nodeIDList, 0)
				require.Len(t, b.nodeMap, 0)
			case testCase.shouldSkipNode:
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
