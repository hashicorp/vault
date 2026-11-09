// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package event

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"syscall"
	"time"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/go-hclog"
	gsyslog "github.com/hashicorp/go-syslog"
)

var _ eventlogger.Node = (*SyslogSink)(nil)

const (
	// syslogSocketLimit is the max bytes per Unix datagram write.
	syslogSocketLimit = 2048

	// syslogFramingOverheadFixed is the fixed overhead go-syslog adds before
	syslogFramingOverheadFixed = 33
)

// SyslogSink is a sink node which handles writing events to syslog.
type SyslogSink struct {
	requiredFormat string
	syslogger      gsyslog.Syslogger
	logger         hclog.Logger
	// chunkSize is the max payload per write: socket limit minus framing and tag.
	chunkSize int
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

	syslog := &SyslogSink{
		requiredFormat: format,
		syslogger:      logger,
		logger:         opts.withLogger,
		chunkSize:      syslogSocketLimit - syslogFramingOverheadFixed - len(opts.withTag),
	}

	return syslog, nil
}

// Process handles writing the event to the syslog.
// Large payloads are split into chunks to avoid EMSGSIZE on the syslog socket.
func (s *SyslogSink) Process(ctx context.Context, e *eventlogger.Event) (_ *eventlogger.Event, retErr error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	defer func() {
		// If the context is errored (cancelled), and we were planning to return
		// an error, let's also log (if we have a logger) in case the eventlogger's
		// status channel and errors propagated.
		if err := ctx.Err(); err != nil && retErr != nil && s.logger != nil {
			s.logger.Error("syslog sink error", "context", err, "error", retErr)
		}
	}()

	if e == nil {
		return nil, fmt.Errorf("event is nil: %w", ErrInvalidParameter)
	}

	formatted, found := e.Format(s.requiredFormat)
	if !found {
		return nil, fmt.Errorf("unable to retrieve event formatted as %q: %w", s.requiredFormat, ErrInvalidParameter)
	}

	if err := s.writeChunked(formatted); err != nil {
		return nil, fmt.Errorf("error writing to syslog: %w", err)
	}

	// return nil for the event to indicate the pipeline is complete.
	return nil, nil
}

// writeChunked writes b to syslog, splitting into chunkSize segments if needed.
// ENOBUFS is retried since it clears quickly once syslogd drains the buffer.
func (s *SyslogSink) writeChunked(b []byte) error {
	if len(b) <= s.chunkSize {
		return s.writeWithRetry(b)
	}

	total := len(b)
	chunks := (total + s.chunkSize - 1) / s.chunkSize
	for i := range chunks {
		start := i * s.chunkSize
		end := start + s.chunkSize
		if end > total {
			end = total
		}
		if err := s.writeWithRetry(b[start:end]); err != nil {
			return err
		}
	}
	return nil
}

// writeWithRetry writes b, retrying up to maxRetries times on ENOBUFS.
func (s *SyslogSink) writeWithRetry(b []byte) error {
	const maxRetries = 5
	wait := 2 * time.Millisecond

	for attempt := range maxRetries + 1 {
		_, err := s.syslogger.Write(b)
		if err == nil {
			return nil
		}
		if attempt < maxRetries && isErrNoBuf(err) {
			time.Sleep(wait)
			wait *= 2
			continue
		}
		return err
	}
	return nil
}

// isErrNoBuf reports whether err wraps ENOBUFS.
func isErrNoBuf(err error) bool {
	var errno syscall.Errno
	return errors.As(err, &errno) && errno == syscall.ENOBUFS
}

// Reopen is a no-op for a syslog sink.
func (_ *SyslogSink) Reopen() error {
	return nil
}

// Type describes the type of this node (sink).
func (_ *SyslogSink) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeSink
}
