// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/eventlogger"
)

var _ eventlogger.Node = (*sinkMetricTimer)(nil)

// sinkMetricTimer is a wrapper for any kind of eventlogger.NodeTypeSink node that
// processes events containing an AuditEvent payload.
// It decorates the implemented eventlogger.Node Process method in order to emit
// timing metrics for the duration between the creation time of the event and the
// time the node completes processing.
type sinkMetricTimer struct {
	name string
	sink eventlogger.Node
}

// newSinkMetricTimer should be used to create the sinkMetricTimer.
// It expects that an eventlogger.NodeTypeSink should be supplied as the sink.
func newSinkMetricTimer(name string, sink eventlogger.Node) (*sinkMetricTimer, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("name is required: %w", ErrInvalidParameter)
	}

	if sink == nil || reflect.ValueOf(sink).IsNil() {
		return nil, fmt.Errorf("sink node is required: %w", ErrInvalidParameter)
	}

	if sink.Type() != eventlogger.NodeTypeSink {
		return nil, fmt.Errorf("sink node must be of type 'sink': %w", ErrInvalidParameter)
	}

	return &sinkMetricTimer{
		name: name,
		sink: sink,
	}, nil
}

// Process wraps the Process method of underlying sink (eventlogger.Node).
// Additionally, when the supplied eventlogger.Event has an AuditEvent as its payload,
// it measures the elapsed time between the creation of the eventlogger.Event and
// the completion of processing, emitting this as a metric.
// Examples:
// 'vault.audit.{DEVICE}.log_request'
// 'vault.audit.{DEVICE}.log_response'
func (s *sinkMetricTimer) Process(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	defer func() {
		auditEvent, ok := e.Payload.(*Event)
		if ok {
			metrics.MeasureSince([]string{"audit", s.name, auditEvent.Subtype.MetricTag()}, e.CreatedAt)
		}
	}()

	return s.sink.Process(ctx, e)
}

// Reopen wraps the Reopen method of this underlying sink (eventlogger.Node).
func (s *sinkMetricTimer) Reopen() error {
	return s.sink.Reopen()
}

// Type wraps the Type method of this underlying sink (eventlogger.Node).
func (s *sinkMetricTimer) Type() eventlogger.NodeType {
	return s.sink.Type()
}
