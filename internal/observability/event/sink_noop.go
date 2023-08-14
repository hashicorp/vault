// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package event

import (
	"context"

	"github.com/hashicorp/eventlogger"
)

var _ eventlogger.Node = (*NoopSink)(nil)

// NoopSink is a sink node which handles ignores everything.
type NoopSink struct{}

// NewNoopSink should be used to create a new NoopSink.
func NewNoopSink() *NoopSink {
	return &NoopSink{}
}

// Process is a no-op and always returns nil event and nil error.
func (_ *NoopSink) Process(ctx context.Context, _ *eventlogger.Event) (*eventlogger.Event, error) {
	// return nil for the event to indicate the pipeline is complete.
	return nil, nil
}

// Reopen is a no-op and always returns nil.
func (_ *NoopSink) Reopen() error {
	return nil
}

// Type describes the type of this node (sink).
func (_ *NoopSink) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeSink
}
