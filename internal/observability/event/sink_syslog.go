// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package event

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/eventlogger"
	gsyslog "github.com/hashicorp/go-syslog"
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
	format = strings.TrimSpace(format)
	if format == "" {
		return nil, fmt.Errorf("format is required: %w", ErrInvalidParameter)
	}

	opts, err := getOpts(opt...)
	if err != nil {
		return nil, err
	}

	logger, err := gsyslog.NewLogger(gsyslog.LOG_INFO, opts.withFacility, opts.withTag)
	if err != nil {
		return nil, fmt.Errorf("error creating syslogger: %w", err)
	}

	return &SyslogSink{requiredFormat: format, logger: logger}, nil
}

// Process handles writing the event to the syslog.
func (s *SyslogSink) Process(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
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

	_, err := s.logger.Write(formatted)
	if err != nil {
		return nil, fmt.Errorf("error writing to syslog: %w", err)
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
