package logical

import "context"

// EventType represents a topic, and is a wrapper around eventlogger.EventType.
type EventType string

// EventSender sends events to the common event bus.
type EventSender interface {
	Send(ctx context.Context, eventType EventType, event any) error
}
