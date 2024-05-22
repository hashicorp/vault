// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package audit

import (
	"testing"

	"github.com/hashicorp/eventlogger"
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
			filter: "operation == \"update\"",
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

			b := &backend{
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
