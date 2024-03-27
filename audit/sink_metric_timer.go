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

var _ eventlogger.Node = (*SinkMetricTimer)(nil)

// SinkMetricTimer is a wrapper for any kind of eventlogger.NodeTypeSink node that
// processes events containing an AuditEvent payload.
// It decorates the implemented eventlogger.Node Process method in order to emit
// timing metrics for the duration between the creation time of the event and the
// time the node completes processing.
type SinkMetricTimer struct {
	Name string
	Sink eventlogger.Node
}

// NewSinkMetricTimer should be used to create the SinkMetricTimer.
// It expects that an eventlogger.NodeTypeSink should be supplied as the sink.
func NewSinkMetricTimer(name string, sink eventlogger.Node) (*SinkMetricTimer, error) {
	const op = "audit.NewSinkMetricTimer"

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("%s: name is required: %w", op, ErrInvalidParameter)
	}

	if sink == nil || reflect.ValueOf(sink).IsNil() {
		return nil, fmt.Errorf("%s: sink node is required: %w", op, ErrInvalidParameter)
	}

	if sink.Type() != eventlogger.NodeTypeSink {
		return nil, fmt.Errorf("%s: sink node must be of type 'sink': %w", op, ErrInvalidParameter)
	}

	return &SinkMetricTimer{
		Name: name,
		Sink: sink,
	}, nil
}

// Process wraps the Process method of underlying sink (eventlogger.Node).
// Additionally, when the supplied eventlogger.Event has an AuditEvent as its payload,
// it measures the elapsed time between the creation of the eventlogger.Event and
// the completion of processing, emitting this as a metric.
// Examples:
// 'vault.audit.{DEVICE}.log_request'
// 'vault.audit.{DEVICE}.log_response'
func (s *SinkMetricTimer) Process(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	defer func() {
		auditEvent, ok := e.Payload.(*AuditEvent)
		if ok {
			metrics.MeasureSince([]string{"audit", s.Name, auditEvent.Subtype.MetricTag()}, e.CreatedAt)
		}
	}()

	return s.Sink.Process(ctx, e)
}

// Reopen wraps the Reopen method of this underlying sink (eventlogger.Node).
func (s *SinkMetricTimer) Reopen() error {
	return s.Sink.Reopen()
}

// Type wraps the Type method of this underlying sink (eventlogger.Node).
func (s *SinkMetricTimer) Type() eventlogger.NodeType {
	return s.Sink.Type()
}
