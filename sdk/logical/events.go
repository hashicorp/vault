package logical

import (
	"context"

	"github.com/hashicorp/go-uuid"
)

// ID is an alias to GetId() for CloudEvents compatibility.
func (x *EventData) ID() string {
	return x.GetId()
}

// NewEvent returns an event with a new, random EID.
func NewEvent() (*EventData, error) {
	id, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	return &EventData{
		Id: id,
	}, nil
}

// EventType represents a topic, and is a wrapper around eventlogger.EventType.
type EventType string

// EventSender sends events to the common event bus.
type EventSender interface {
	Send(ctx context.Context, eventType EventType, event *EventData) error
}
