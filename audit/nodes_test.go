// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package audit

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/helper/namespace"

	"github.com/hashicorp/go-uuid"

	"github.com/hashicorp/vault/sdk/logical"

	"github.com/hashicorp/eventlogger"

	"github.com/hashicorp/vault/internal/observability/event"
	"github.com/stretchr/testify/require"
)

// TestProcessManual ensures that the manual processing of a test message works
// as expected, covering various inputs.
func TestProcessManual(t *testing.T) {
	tests := map[string]struct {
		IsErrorExpected        bool
		IsErrorContains        bool
		ExpectedErrorMessage   string
		ShouldUseData          bool
		ShouldUseIDs           bool
		ShouldUseNodes         bool
		ShouldCreateFilterNode bool
		ShouldCreateSinkNode   bool
		ShouldMatchIdToFilter  bool
		ShouldMatchIdToSink    bool
	}{
		"nil-data": {
			IsErrorExpected:        true,
			IsErrorContains:        false,
			ExpectedErrorMessage:   "data cannot be nil",
			ShouldUseData:          false,
			ShouldUseIDs:           false,
			ShouldUseNodes:         false,
			ShouldCreateFilterNode: false,
			ShouldCreateSinkNode:   false,
			ShouldMatchIdToFilter:  false,
			ShouldMatchIdToSink:    false,
		},
		"no-ids": {
			IsErrorExpected:        true,
			IsErrorContains:        false,
			ExpectedErrorMessage:   "ids are required",
			ShouldUseData:          true,
			ShouldUseIDs:           false,
			ShouldUseNodes:         false,
			ShouldCreateFilterNode: false,
			ShouldCreateSinkNode:   false,
			ShouldMatchIdToFilter:  false,
			ShouldMatchIdToSink:    false,
		},
		"no-nodes": {
			IsErrorExpected:        true,
			IsErrorContains:        false,
			ExpectedErrorMessage:   "nodes are required",
			ShouldUseData:          true,
			ShouldUseIDs:           true,
			ShouldUseNodes:         true,
			ShouldCreateFilterNode: true,
			ShouldCreateSinkNode:   true,
			ShouldMatchIdToFilter:  false,
			ShouldMatchIdToSink:    false,
		},
		"id-node-mismatch": {
			IsErrorExpected:        true,
			IsErrorContains:        true,
			ExpectedErrorMessage:   "node not found",
			ShouldUseData:          true,
			ShouldUseIDs:           true,
			ShouldUseNodes:         true,
			ShouldCreateFilterNode: true,
			ShouldMatchIdToFilter:  true,
			ShouldCreateSinkNode:   true,
			ShouldMatchIdToSink:    false,
		},
		"last-node-not-sink": {
			IsErrorExpected:        true,
			IsErrorContains:        false,
			ExpectedErrorMessage:   "last node must be a sink",
			ShouldUseData:          true,
			ShouldUseIDs:           true,
			ShouldUseNodes:         true,
			ShouldCreateFilterNode: true,
			ShouldCreateSinkNode:   false,
			ShouldMatchIdToFilter:  true,
			ShouldMatchIdToSink:    false,
		},
		"normal-operation": {
			IsErrorExpected:        false,
			ShouldUseData:          true,
			ShouldUseIDs:           true,
			ShouldUseNodes:         true,
			ShouldCreateFilterNode: true,
			ShouldCreateSinkNode:   true,
			ShouldMatchIdToFilter:  true,
			ShouldMatchIdToSink:    true,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var ids []eventlogger.NodeID
			var nodes map[eventlogger.NodeID]eventlogger.Node
			var data *logical.LogInput

			if tc.ShouldUseNodes {
				nodes = make(map[eventlogger.NodeID]eventlogger.Node)

				if tc.ShouldCreateFilterNode {
					filterId, filterNode := newFilterNode(t)

					if tc.ShouldUseIDs {
						ids = append(ids, filterId)
					}

					if tc.ShouldMatchIdToFilter {
						nodes[filterId] = filterNode
					}
				}

				if tc.ShouldCreateSinkNode {
					sinkId, sinkNode := newSinkNode(t)
					if tc.ShouldUseIDs {
						ids = append(ids, sinkId)
					}

					if tc.ShouldMatchIdToSink {
						nodes[sinkId] = sinkNode
					}
				}
			}

			if tc.ShouldUseData {
				requestId, err := uuid.GenerateUUID()
				require.NoError(t, err)
				data = newData(requestId)
			}

			err := ProcessManual(namespace.RootContext(context.Background()), data, ids, nodes)
			if tc.IsErrorExpected {
				require.Error(t, err)
				if tc.IsErrorContains {
					require.ErrorContains(t, err, tc.ExpectedErrorMessage)
				} else {
					require.EqualError(t, err, tc.ExpectedErrorMessage)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// newFilterNode creates a new UUID and EntryFormatter (filter node).
func newFilterNode(t *testing.T) (eventlogger.NodeID, *EntryFormatter) {
	t.Helper()

	filterId, err := event.GenerateNodeID()
	require.NoError(t, err)
	cfg, err := NewFormatterConfig()
	require.NoError(t, err)
	filterNode, err := NewEntryFormatter(cfg, newStaticSalt(t))
	require.NoError(t, err)

	return filterId, filterNode
}

// newSinkNode creates a new UUID and NoopSink (sink node).
func newSinkNode(t *testing.T) (eventlogger.NodeID, *event.NoopSink) {
	t.Helper()

	sinkId, err := event.GenerateNodeID()
	require.NoError(t, err)
	sinkNode := event.NewNoopSink()

	return sinkId, sinkNode
}

// newData creates a sample logical.LogInput to be used as data for tests.
func newData(id string) *logical.LogInput {
	return &logical.LogInput{
		Type: "request",
		Auth: nil,
		Request: &logical.Request{
			ID:        id,
			Operation: "update",
			Path:      "sys/audit/test",
		},
		Response: nil,
		OuterErr: nil,
	}
}
