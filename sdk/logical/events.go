package logical

import "context"

// EventSender sends events to the common event bus.
type EventSender interface {
	Send(context.Context, string, any) error
}
