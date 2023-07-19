package event

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/eventlogger"
)

var _ eventlogger.Node = (*StdoutSinkNode)(nil)

// StdoutSinkNode is structure that implements the eventlogger.Node interface
// as a Sink node that writes the events to the standard output stream.
type StdoutSinkNode struct {
	expectedFormat string
}

// NewStdoutSinkNode creates a new StdoutSinkNode that will persist the events
// it processes using the specified expected format.
func NewStdoutSinkNode(expectedFormat string) *StdoutSinkNode {
	return &StdoutSinkNode{
		expectedFormat: expectedFormat,
	}
}

// Process persists the provided eventlogger.Event to the standard output stream.
func (n *StdoutSinkNode) Process(ctx context.Context, event *eventlogger.Event) (*eventlogger.Event, error) {
	const op = "event.(StdoutSinkNode).Process"

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if event == nil {
		return nil, fmt.Errorf("%s: event is nil: %w", op, ErrInvalidParameter)
	}

	formattedBytes, found := event.Format(n.expectedFormat)
	if !found {
		return nil, fmt.Errorf("%s: unable to retrieve event formatted as %q", op, n.expectedFormat)
	}

	_, err := os.Stdout.Write(formattedBytes)
	if err != nil {
		return nil, fmt.Errorf("%s: error writing to stdout: %w", op, err)
	}

	// Return nil, nil to indicate the pipeline is complete.
	return nil, nil
}

// Reopen is a no-op for the StdoutSinkNode type.
func (n *StdoutSinkNode) Reopen() error {
	return nil
}

// Type returns the eventlogger.NodeTypeSink constant.
func (n *StdoutSinkNode) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeSink
}
