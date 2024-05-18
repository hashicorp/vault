// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/internal/observability/event"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// TestEntryFilter_NewEntryFilter tests that we can create entryFilter types correctly.
func TestEntryFilter_NewEntryFilter(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		Filter               string
		IsErrorExpected      bool
		ExpectedErrorMessage string
	}{
		"empty-filter": {
			Filter:               "",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "cannot create new audit filter with empty filter expression: invalid configuration",
		},
		"spacey-filter": {
			Filter:               "    ",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "cannot create new audit filter with empty filter expression: invalid configuration",
		},
		"bad-filter": {
			Filter:               "____",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "cannot create new audit filter",
		},
		"unsupported-field-filter": {
			Filter:               "foo == bar",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "filter references an unsupported field: foo == bar",
		},
		"good-filter-operation": {
			Filter:          "operation == create",
			IsErrorExpected: false,
		},
		"good-filter-mount_type": {
			Filter:          "mount_type == kv",
			IsErrorExpected: false,
		},
		"good-filter-mount_point": {
			Filter:          "mount_point == \"/auth/userpass\"",
			IsErrorExpected: false,
		},
		"good-filter-namespace": {
			Filter:          "namespace == juan",
			IsErrorExpected: false,
		},
		"good-filter-path": {
			Filter:          "path == foo",
			IsErrorExpected: false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			f, err := newEntryFilter(tc.Filter)
			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.ErrorContains(t, err, tc.ExpectedErrorMessage)
				require.Nil(t, f)
			default:
				require.NoError(t, err)
				require.NotNil(t, f)
			}
		})
	}
}

// TestEntryFilter_Reopen ensures we can reopen the filter node.
func TestEntryFilter_Reopen(t *testing.T) {
	t.Parallel()

	f := &entryFilter{}
	res := f.Reopen()
	require.Nil(t, res)
}

// TestEntryFilter_Type ensures we always return the right type for this node.
func TestEntryFilter_Type(t *testing.T) {
	t.Parallel()

	f := &entryFilter{}
	require.Equal(t, eventlogger.NodeTypeFilter, f.Type())
}

// TestEntryFilter_Process_ContextDone ensures that we stop processing the event
// if the context was cancelled.
func TestEntryFilter_Process_ContextDone(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())

	// Explicitly cancel the context
	cancel()

	l, err := newEntryFilter("operation == foo")
	require.NoError(t, err)

	// Fake audit event
	a, err := NewEvent(RequestType)
	require.NoError(t, err)

	// Fake event logger event
	e := &eventlogger.Event{
		Type:      event.AuditType.AsEventType(),
		CreatedAt: time.Now(),
		Formatted: make(map[string][]byte),
		Payload:   a,
	}

	e2, err := l.Process(ctx, e)

	require.Error(t, err)
	require.ErrorContains(t, err, "context canceled")

	// Ensure that the pipeline won't continue.
	require.Nil(t, e2)
}

// TestEntryFilter_Process_NilEvent ensures we receive the right error when the
// event we are trying to process is nil.
func TestEntryFilter_Process_NilEvent(t *testing.T) {
	t.Parallel()

	l, err := newEntryFilter("operation == foo")
	require.NoError(t, err)
	e, err := l.Process(context.Background(), nil)
	require.Error(t, err)
	require.EqualError(t, err, "event is nil: invalid internal parameter")

	// Ensure that the pipeline won't continue.
	require.Nil(t, e)
}

// TestEntryFilter_Process_BadPayload ensures we receive the correct error when
// attempting to process an event with a payload that cannot be parsed back to
// an audit event.
func TestEntryFilter_Process_BadPayload(t *testing.T) {
	t.Parallel()

	l, err := newEntryFilter("operation == foo")
	require.NoError(t, err)

	e := &eventlogger.Event{
		Type:      event.AuditType.AsEventType(),
		CreatedAt: time.Now(),
		Formatted: make(map[string][]byte),
		Payload:   nil,
	}

	e2, err := l.Process(context.Background(), e)
	require.Error(t, err)
	require.EqualError(t, err, "cannot parse event payload: invalid internal parameter")

	// Ensure that the pipeline won't continue.
	require.Nil(t, e2)
}

// TestEntryFilter_Process_NoAuditDataInPayload ensure we stop processing a pipeline
// when the data in the audit event is nil.
func TestEntryFilter_Process_NoAuditDataInPayload(t *testing.T) {
	t.Parallel()

	l, err := newEntryFilter("operation == foo")
	require.NoError(t, err)

	a, err := NewEvent(RequestType)
	require.NoError(t, err)

	// Ensure audit data is nil
	a.Data = nil

	e := &eventlogger.Event{
		Type:      event.AuditType.AsEventType(),
		CreatedAt: time.Now(),
		Formatted: make(map[string][]byte),
		Payload:   a,
	}

	e2, err := l.Process(context.Background(), e)

	// Make sure we get the 'nil, nil' response to stop processing this pipeline.
	require.NoError(t, err)
	require.Nil(t, e2)
}

// TestEntryFilter_Process_FilterSuccess tests that when a filter matches we
// receive no error and the event is not nil so it continues in the pipeline.
func TestEntryFilter_Process_FilterSuccess(t *testing.T) {
	t.Parallel()

	l, err := newEntryFilter("mount_type == juan")
	require.NoError(t, err)

	a, err := NewEvent(RequestType)
	require.NoError(t, err)

	a.Data = &logical.LogInput{
		Request: &logical.Request{
			Operation: logical.CreateOperation,
			MountType: "juan",
		},
	}

	e := &eventlogger.Event{
		Type:      event.AuditType.AsEventType(),
		CreatedAt: time.Now(),
		Formatted: make(map[string][]byte),
		Payload:   a,
	}

	ctx := namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)

	e2, err := l.Process(ctx, e)

	require.NoError(t, err)
	require.NotNil(t, e2)
}

// TestEntryFilter_Process_FilterFail tests that when a filter fails to match we
// receive no error, but also the event is nil so that the pipeline completes.
func TestEntryFilter_Process_FilterFail(t *testing.T) {
	t.Parallel()

	l, err := newEntryFilter("mount_type == john and operation == create and namespace == root")
	require.NoError(t, err)

	a, err := NewEvent(RequestType)
	require.NoError(t, err)

	a.Data = &logical.LogInput{
		Request: &logical.Request{
			Operation: logical.CreateOperation,
			MountType: "juan",
		},
	}

	e := &eventlogger.Event{
		Type:      event.AuditType.AsEventType(),
		CreatedAt: time.Now(),
		Formatted: make(map[string][]byte),
		Payload:   a,
	}

	ctx := namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)

	e2, err := l.Process(ctx, e)

	require.NoError(t, err)
	require.Nil(t, e2)
}
