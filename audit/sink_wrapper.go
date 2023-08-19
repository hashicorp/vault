package audit

import (
	"context"

	metrics "github.com/armon/go-metrics"

	"github.com/hashicorp/eventlogger"
)

// SinkWrapper is a wrapper for any kind of Sink Node that processes events
// containing an auditEvent payload.
type SinkWrapper struct {
	Name string
	Sink eventlogger.Node
}

// Process simply wraps the Process method of this SinkWrapper's sink field by
// taking a measurement of the time elapsed since the provided Event was created
// once this method returns.
func (s *SinkWrapper) Process(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	defer func() {
		auditEvent, ok := e.Payload.(*auditEvent)
		if ok {
			metrics.MeasureSince([]string{"audit", s.Name, auditEvent.Subtype.MetricTag()}, e.CreatedAt)
		}
	}()

	return s.Sink.Process(ctx, e)
}

// Reopen simply wraps the Reopen method of this SinkWrapper's sink field
// without doing any additional work.
func (s *SinkWrapper) Reopen() error {
	return s.Sink.Reopen()
}

// Type simply wraps the Type method of this SinkWrapper's sink field without
// doing any additional work.
func (s *SinkWrapper) Type() eventlogger.NodeType {
	return s.Sink.Type()
}
