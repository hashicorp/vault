// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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

// TestProcessManual_NilData tests ProcessManual when nil data is supplied.
func TestProcessManual_NilData(t *testing.T) {
	t.Parallel()

	var ids []eventlogger.NodeID
	nodes := make(map[eventlogger.NodeID]eventlogger.Node)

	// Formatter node
	formatterId, formatterNode := newFormatterNode(t)
	ids = append(ids, formatterId)
	nodes[formatterId] = formatterNode

	// Sink node
	sinkId, sinkNode := newSinkNode(t)
	ids = append(ids, sinkId)
	nodes[sinkId] = sinkNode

	err := ProcessManual(namespace.RootContext(context.Background()), nil, ids, nodes)
	require.Error(t, err)
	require.EqualError(t, err, "data cannot be nil")
}

// TestProcessManual_BadIDs tests ProcessManual when different bad values are
// supplied for the ID parameter.
func TestProcessManual_BadIDs(t *testing.T) {
	tests := map[string]struct {
		IDs                  []eventlogger.NodeID
		ExpectedErrorMessage string
	}{
		"nil": {
			IDs:                  nil,
			ExpectedErrorMessage: "minimum of 2 ids are required",
		},
		"one": {
			IDs:                  []eventlogger.NodeID{"1"},
			ExpectedErrorMessage: "minimum of 2 ids are required",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			nodes := make(map[eventlogger.NodeID]eventlogger.Node)

			// Formatter node
			formatterId, formatterNode := newFormatterNode(t)
			nodes[formatterId] = formatterNode

			// Sink node
			sinkId, sinkNode := newSinkNode(t)
			nodes[sinkId] = sinkNode

			// Data
			requestId, err := uuid.GenerateUUID()
			require.NoError(t, err)
			data := newData(requestId)

			err = ProcessManual(namespace.RootContext(context.Background()), data, tc.IDs, nodes)
			require.Error(t, err)
			require.EqualError(t, err, tc.ExpectedErrorMessage)
		})
	}
}

// TestProcessManual_NoNodes tests ProcessManual when no nodes are supplied.
func TestProcessManual_NoNodes(t *testing.T) {
	t.Parallel()

	var ids []eventlogger.NodeID
	nodes := make(map[eventlogger.NodeID]eventlogger.Node)

	// Formatter node
	formatterId, _ := newFormatterNode(t)
	ids = append(ids, formatterId)

	// Sink node
	sinkId, _ := newSinkNode(t)
	ids = append(ids, sinkId)

	// Data
	requestId, err := uuid.GenerateUUID()
	require.NoError(t, err)
	data := newData(requestId)

	err = ProcessManual(namespace.RootContext(context.Background()), data, ids, nodes)
	require.Error(t, err)
	require.EqualError(t, err, "nodes are required")
}

// TestProcessManual_IdNodeMismatch tests ProcessManual when IDs don't match with
// the nodes in the supplied map.
func TestProcessManual_IdNodeMismatch(t *testing.T) {
	t.Parallel()

	var ids []eventlogger.NodeID
	nodes := make(map[eventlogger.NodeID]eventlogger.Node)

	// Formatter node
	formatterId, formatterNode := newFormatterNode(t)
	ids = append(ids, formatterId)
	nodes[formatterId] = formatterNode

	// Sink node
	sinkId, _ := newSinkNode(t)
	ids = append(ids, sinkId)

	// Data
	requestId, err := uuid.GenerateUUID()
	require.NoError(t, err)
	data := newData(requestId)

	err = ProcessManual(namespace.RootContext(context.Background()), data, ids, nodes)
	require.Error(t, err)
	require.ErrorContains(t, err, "node not found: ")
}

// TestProcessManual_NotEnoughNodes tests ProcessManual when there is only one
// node provided.
func TestProcessManual_NotEnoughNodes(t *testing.T) {
	t.Parallel()

	var ids []eventlogger.NodeID
	nodes := make(map[eventlogger.NodeID]eventlogger.Node)

	// Formatter node
	formatterId, formatterNode := newFormatterNode(t)
	ids = append(ids, formatterId)
	nodes[formatterId] = formatterNode

	// Data
	requestId, err := uuid.GenerateUUID()
	require.NoError(t, err)
	data := newData(requestId)

	err = ProcessManual(namespace.RootContext(context.Background()), data, ids, nodes)
	require.Error(t, err)
	require.EqualError(t, err, "minimum of 2 ids are required")
}

// TestProcessManual_LastNodeNotSink tests ProcessManual when the last node is
// not a Sink node.
func TestProcessManual_LastNodeNotSink(t *testing.T) {
	t.Parallel()

	var ids []eventlogger.NodeID
	nodes := make(map[eventlogger.NodeID]eventlogger.Node)

	// Formatter node
	formatterId, formatterNode := newFormatterNode(t)
	ids = append(ids, formatterId)
	nodes[formatterId] = formatterNode

	// Another Formatter node
	formatterId, formatterNode = newFormatterNode(t)
	ids = append(ids, formatterId)
	nodes[formatterId] = formatterNode

	// Data
	requestId, err := uuid.GenerateUUID()
	require.NoError(t, err)
	data := newData(requestId)

	err = ProcessManual(namespace.RootContext(context.Background()), data, ids, nodes)
	require.Error(t, err)
	require.EqualError(t, err, "last node must be a filter or sink")
}

// TestProcessManualEndWithSink ensures that the manual processing of a test
// message works as expected with proper inputs, which mean processing ends with
// sink node.
func TestProcessManualEndWithSink(t *testing.T) {
	t.Parallel()

	var ids []eventlogger.NodeID
	nodes := make(map[eventlogger.NodeID]eventlogger.Node)

	// Formatter node
	formatterId, formatterNode := newFormatterNode(t)
	ids = append(ids, formatterId)
	nodes[formatterId] = formatterNode

	// Sink node
	sinkId, sinkNode := newSinkNode(t)
	ids = append(ids, sinkId)
	nodes[sinkId] = sinkNode

	// Data
	requestId, err := uuid.GenerateUUID()
	require.NoError(t, err)
	data := newData(requestId)

	err = ProcessManual(namespace.RootContext(context.Background()), data, ids, nodes)
	require.NoError(t, err)
}

// TestProcessManual_EndWithFilter ensures that the manual processing of a test
// message works as expected with proper inputs, which mean processing ends with
// sink node.
func TestProcessManual_EndWithFilter(t *testing.T) {
	t.Parallel()

	var ids []eventlogger.NodeID
	nodes := make(map[eventlogger.NodeID]eventlogger.Node)

	// Filter node
	filterId, filterNode := newFilterNode(t)
	ids = append(ids, filterId)
	nodes[filterId] = filterNode

	// Formatter node
	formatterId, formatterNode := newFormatterNode(t)
	ids = append(ids, formatterId)
	nodes[formatterId] = formatterNode

	// Sink node
	sinkId, sinkNode := newSinkNode(t)
	ids = append(ids, sinkId)
	nodes[sinkId] = sinkNode

	// Data
	requestId, err := uuid.GenerateUUID()
	require.NoError(t, err)
	data := newData(requestId)

	err = ProcessManual(namespace.RootContext(context.Background()), data, ids, nodes)
	require.NoError(t, err)
}

// newSinkNode creates a new UUID and NoopSink (sink node).
func newSinkNode(t *testing.T) (eventlogger.NodeID, *event.NoopSink) {
	t.Helper()

	sinkId, err := event.GenerateNodeID()
	require.NoError(t, err)
	sinkNode := event.NewNoopSink()

	return sinkId, sinkNode
}

// TestFilter is a trivial implementation of eventlogger.Node used as a placeholder
// for Filter nodes in tests.
type TestFilter struct{}

// Process trivially filters the event preventing it from being processed by subsequent nodes.
func (f *TestFilter) Process(_ context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	return nil, nil
}

// Reopen does nothing.
func (f *TestFilter) Reopen() error {
	return nil
}

// Type returns the eventlogger.NodeTypeFormatter type.
func (f *TestFilter) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeFilter
}

// TestFormatter is a trivial implementation of the eventlogger.Node interface
// used as a place-holder for Formatter nodes in tests.
type TestFormatter struct{}

// Process trivially formats the event by storing "test" as a byte slice under
// the test format type.
func (f *TestFormatter) Process(_ context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	e.FormattedAs("test", []byte("test"))

	return e, nil
}

// Reopen does nothing.
func (f *TestFormatter) Reopen() error {
	return nil
}

// Type returns the eventlogger.NodeTypeFormatter type.
func (f *TestFormatter) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeFormatter
}

// newFilterNode creates a new TestFormatter (filter node).
func newFilterNode(t *testing.T) (eventlogger.NodeID, *TestFilter) {
	nodeId, err := event.GenerateNodeID()
	require.NoError(t, err)
	node := &TestFilter{}

	return nodeId, node
}

// newFormatterNode creates a new TestFormatter (formatter node).
func newFormatterNode(t *testing.T) (eventlogger.NodeID, *TestFormatter) {
	nodeId, err := event.GenerateNodeID()
	require.NoError(t, err)
	node := &TestFormatter{}

	return nodeId, node
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
