// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package event

import (
	"context"
	"fmt"
	"strings"

	gsyslog "github.com/hashicorp/go-syslog"

	"github.com/hashicorp/eventlogger"
)

var _ eventlogger.Node = (*SyslogSink)(nil)

// SyslogSink is a sink node which handles writing events to syslog.
type SyslogSink struct {
	requiredFormat string
	logger         gsyslog.Syslogger
}

// NewSyslogSink should be used to create a new SyslogSink.
// Accepted options: WithFacility and WithTag.
func NewSyslogSink(format string, opt ...Option) (*SyslogSink, error) {
	const op = "event.NewSyslogSink"

	format = strings.TrimSpace(format)
	if format == "" {
		return nil, fmt.Errorf("%s: format is required: %w", op, ErrInvalidParameter)
	}

	opts, err := getOpts(opt...)
	if err != nil {
		return nil, fmt.Errorf("%s: error applying options: %w", op, err)
	}

	logger, err := gsyslog.NewLogger(gsyslog.LOG_INFO, opts.withFacility, opts.withTag)
	if err != nil {
		return nil, fmt.Errorf("%s: error creating syslogger: %w", op, err)
	}

	return &SyslogSink{requiredFormat: format, logger: logger}, nil
}

// Process handles writing the event to the syslog.
func (s *SyslogSink) Process(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	const op = "event.(SyslogSink).Process"

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if e == nil {
		return nil, fmt.Errorf("%s: event is nil: %w", op, ErrInvalidParameter)
	}

	formatted, found := e.Format(s.requiredFormat)
	if !found {
		return nil, fmt.Errorf("%s: unable to retrieve event formatted as %q", op, s.requiredFormat)
	}

	_, err := s.logger.Write(formatted)
	if err != nil {
		return nil, fmt.Errorf("%s: error writing to syslog: %w", op, err)
	}

	// return nil for the event to indicate the pipeline is complete.
	return nil, nil
}

// Reopen is a no-op for a syslog sink.
func (_ *SyslogSink) Reopen() error {
	return nil
}

// Type describes the type of this node (sink).
func (_ *SyslogSink) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeSink
}
