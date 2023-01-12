package logical

import (
	"context"

	"github.com/hashicorp/go-uuid"
)

// NewEvent returns an event with a new, random EID.
func NewEvent() *EventData {
	eid, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}
	return &EventData{
		Eid: eid,
	}
}

// EventType represents a topic, and is a wrapper around eventlogger.EventType.
type EventType string

// EventSender sends events to the common event bus.
type EventSender interface {
	Send(ctx context.Context, eventType EventType, event *EventData) error
}
