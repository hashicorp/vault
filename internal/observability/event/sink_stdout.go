// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package event

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/eventlogger"
)

var _ eventlogger.Node = (*StdoutSink)(nil)

// StdoutSink is structure that implements the eventlogger.Node interface
// as a Sink node that writes the events to the standard output stream.
type StdoutSink struct {
	requiredFormat string
}

// NewStdoutSinkNode creates a new StdoutSink that will persist the events
// it processes using the specified expected format.
func NewStdoutSinkNode(format string) (*StdoutSink, error) {
	format = strings.TrimSpace(format)
	if format == "" {
		return nil, fmt.Errorf("format is required: %w", ErrInvalidParameter)
	}

	return &StdoutSink{
		requiredFormat: format,
	}, nil
}

// Process persists the provided eventlogger.Event to the standard output stream.
func (s *StdoutSink) Process(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if e == nil {
		return nil, fmt.Errorf("event is nil: %w", ErrInvalidParameter)
	}

	formatted, found := e.Format(s.requiredFormat)
	if !found {
		return nil, fmt.Errorf("unable to retrieve event formatted as %q: %w", s.requiredFormat, ErrInvalidParameter)
	}

	_, err := os.Stdout.Write(formatted)
	if err != nil {
		return nil, fmt.Errorf("error writing to stdout: %w", err)
	}

	// Return nil, nil to indicate the pipeline is complete.
	return nil, nil
}

// Reopen is a no-op for the StdoutSink type.
func (s *StdoutSink) Reopen() error {
	return nil
}

// Type returns the eventlogger.NodeTypeSink constant.
func (s *StdoutSink) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeSink
}
