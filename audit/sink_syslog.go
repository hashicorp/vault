// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package audit

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/internal/observability/event"

	gsyslog "github.com/hashicorp/go-syslog"

	"github.com/hashicorp/eventlogger"
)

// AuditSyslogSink is a sink node which handles writing audit events to syslog.
type AuditSyslogSink struct {
	format format
	logger gsyslog.Syslogger
}

// NewAuditSyslogSink should be used to create a new AuditSyslogSink.
// Accepted options: WithFacility and WithTag.
func NewAuditSyslogSink(format format, opt ...Option) (*AuditSyslogSink, error) {
	const op = "audit.NewAuditSyslogSink"

	opts, err := getOpts(opt...)
	if err != nil {
		return nil, fmt.Errorf("%s: error applying options: %w", op, err)
	}

	logger, err := gsyslog.NewLogger(gsyslog.LOG_INFO, opts.withFacility, opts.withTag)
	if err != nil {
		return nil, fmt.Errorf("%s: error creating syslogger: %w", op, err)
	}

	return &AuditSyslogSink{format: format, logger: logger}, nil
}

// Process handles writing the event to the syslog.
func (s *AuditSyslogSink) Process(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	const op = "audit.(AuditSyslogSink).Process"

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if e == nil {
		return nil, fmt.Errorf("%s: event is nil: %w", op, event.ErrInvalidParameter)
	}

	formatted, found := e.Format(s.format.String())
	if !found {
		return nil, fmt.Errorf("%s: unable to retrieve event formatted as %q", op, s.format)
	}

	_, err := s.logger.Write(formatted)
	if err != nil {
		return nil, fmt.Errorf("%s: error writing to syslog: %w", op, err)
	}

	// return nil for the event to indicate the pipeline is complete.
	return nil, nil
}

// Reopen is a no-op for a syslog sink.
func (s *AuditSyslogSink) Reopen() error {
	return nil
}

// Type describes the type of this node (sink).
func (s *AuditSyslogSink) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeSink
}
