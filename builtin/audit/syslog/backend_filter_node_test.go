// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package syslog

import (
	"testing"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/audit"
	"github.com/stretchr/testify/require"
)

// TestBackend_configureFilterNode ensures that configureFilterNode handles various
// filter values as expected. Empty (including whitespace) strings should return
// no error but skip configuration of the node.
// NOTE: Audit filtering is an Enterprise feature and behaves differently in the
// community edition of Vault.
func TestBackend_configureFilterNode(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		filter string
	}{
		"happy": {
			filter: "operation == update",
		},
		"empty": {
			filter: "",
		},
		"spacey": {
			filter: "    ",
		},
		"bad": {
			filter: "___qwerty",
		},
		"unsupported-field": {
			filter: "foo == bar",
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
			require.NoError(t, err)
			require.Len(t, b.nodeIDList, 0)
			require.Len(t, b.nodeMap, 0)
		})
	}
}

// TestBackend_configureFilterFormatterSink ensures that configuring all three
// types of nodes on a Backend works as expected, i.e. we have only formatter and sink
// nodes at the end and nothing gets overwritten. The order of calls influences the
// slice of IDs on the Backend.
// NOTE: Audit filtering is an Enterprise feature and behaves differently in the
// community edition of Vault.
func TestBackend_configureFilterFormatterSink(t *testing.T) {
	t.Parallel()

	b := &Backend{
		nodeIDList: []eventlogger.NodeID{},
		nodeMap:    map[eventlogger.NodeID]eventlogger.Node{},
	}

	formatConfig, err := audit.NewFormatterConfig()
	require.NoError(t, err)

	err = b.configureFilterNode("path == bar")
	require.NoError(t, err)

	err = b.configureFormatterNode("juan", formatConfig, hclog.NewNullLogger())
	require.NoError(t, err)

	err = b.configureSinkNode("foo", "json")
	require.NoError(t, err)

	require.Len(t, b.nodeIDList, 2)
	require.Len(t, b.nodeMap, 2)

	id := b.nodeIDList[0]
	node := b.nodeMap[id]
	require.Equal(t, eventlogger.NodeTypeFormatter, node.Type())

	id = b.nodeIDList[1]
	node = b.nodeMap[id]
	require.Equal(t, eventlogger.NodeTypeSink, node.Type())
}
