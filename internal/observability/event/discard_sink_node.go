package event

import (
	"context"

	"github.com/hashicorp/eventlogger"
)

// Make sure that the DiscardSinkNode type satisfies the eventlogger.Node
// interface.
var _ eventlogger.Node = (*DiscardSinkNode)(nil)

// DiscardSinkNode is a structure that implements the eventlogger.Node interface
// as a Sink node that simply discards the event.
type DiscardSinkNode struct {
}

// Process an eventlogger.Event by simply discarding it.
func (n *DiscardSinkNode) Process(ctx context.Context, event *eventlogger.Event) (*eventlogger.Event, error) {
	// Return nil, nil to indicate the pipeline is complete.
	return nil, nil
}

// Reopen is a no-op for the DiscardSinkNode type.
func (n *DiscardSinkNode) Reopen() error {
	return nil
}

// Type returns the eventlogger.NodeTypeSink constant.
func (n *DiscardSinkNode) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeSink
}
