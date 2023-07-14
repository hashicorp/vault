// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"

	"github.com/hashicorp/eventlogger"
)

// AuditSocketSink is a sink node which handles writing audit events to socket.
type AuditSocketSink struct {
	format        auditFormat
	address       string
	socketType    string        // socketType (default 'tcp')
	writeDuration time.Duration // default 2s
	socketLock    sync.RWMutex
	connection    net.Conn
}

// NewAuditSocketSink should be used to create a new AuditSocketSink.
// Accepted options: WithDuration and WithSocketType.
func NewAuditSocketSink(format auditFormat, address string, opt ...Option) (*AuditSocketSink, error) {
	const op = "event.NewAuditSocketSink"

	opts, err := getOpts(opt...)
	if err != nil {
		return nil, fmt.Errorf("%s: error applying options: %w", op, err)
	}

	sink := &AuditSocketSink{
		format:        format,
		address:       address,
		socketType:    opts.withSocketType,    // default tcp
		writeDuration: opts.withWriteDuration, // default 2s
		socketLock:    sync.RWMutex{},
		connection:    nil,
	}

	return sink, nil
}

// Process handles writing the event to the socket.
func (s *AuditSocketSink) Process(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	const op = "event.(AuditSocketSink).Process"

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if e == nil {
		return nil, fmt.Errorf("%s: event is nil: %w", op, ErrInvalidParameter)
	}

	formatted, found := e.Format(s.format.String())
	if !found {
		return nil, fmt.Errorf("%s: unable to retrieve event formatted as %q", op, s.format)
	}

	buffer := bytes.NewBuffer(formatted)

	// Try writing and return early if successful.
	err := s.write(ctx, buffer.Bytes())
	if err == nil {
		return nil, nil
	}

	// We will try to reconnect and retry a single write.
	reconErr := s.reconnect(ctx)
	switch {
	case reconErr != nil:
		// Add the reconnection error to the existing error.
		err = multierror.Append(err, reconErr)
	default:
		err = s.write(ctx, buffer.Bytes())
	}

	// Format the error nicely if we need to return one.
	if err != nil {
		err = fmt.Errorf("%s: error writing to socket: %w", op, err)
	}

	return nil, err
}

// Reopen handles reopening the connection for the socket sink.
func (s *AuditSocketSink) Reopen() error {
	const op = "event.(AuditSocketSink).Reopen"

	s.socketLock.Lock()
	defer s.socketLock.Unlock()

	err := s.reconnect(nil)

	return err
}

// Type describes the type of this node (sink).
func (s *AuditSocketSink) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeSink
}

// connect attempts to establish a connection using the socketType and address.
func (s *AuditSocketSink) connect(ctx context.Context) error {
	const op = "event.(AuditSocketSink).connect"

	// If we're already connected, we should have disconnected first.
	if s.connection != nil {
		return nil
	}

	timeoutContext, cancel := context.WithTimeout(ctx, s.writeDuration)
	defer cancel()

	dialer := net.Dialer{}
	conn, err := dialer.DialContext(timeoutContext, s.socketType, s.address)
	if err != nil {
		return fmt.Errorf("%s: error connecting to %q address %q: %w", op, s.socketType, s.address, err)
	}

	s.connection = conn

	return nil
}

// disconnect attempts to close and clear an existing connection.
func (s *AuditSocketSink) disconnect() error {
	const op = "event.(AuditSocketSink).disconnect"

	// If we're already disconnected, we can return early.
	if s.connection == nil {
		return nil
	}

	err := s.connection.Close()
	if err != nil {
		return fmt.Errorf("%s: error closing connection: %w", op, err)
	}
	s.connection = nil

	return nil
}

// reconnect attempts to disconnect and then connect to the configured socketType and address.
func (s *AuditSocketSink) reconnect(ctx context.Context) error {
	const op = "event.(AuditSocketSink).reconnect"

	err := s.disconnect()
	if err != nil {
		return fmt.Errorf("%s: error disconnecting: %w", op, err)
	}

	err = s.connect(ctx)
	if err != nil {
		return fmt.Errorf("%s: error connecting: %w", op, err)
	}

	return nil
}

// write attempts to write the specified data using the established connection.
func (s *AuditSocketSink) write(ctx context.Context, data []byte) error {
	const op = "event.(AuditSocketSink).write"

	// Ensure we're connected.
	err := s.connect(ctx)
	if err != nil {
		return fmt.Errorf("%s: connection error: %w", op, err)
	}

	err = s.connection.SetWriteDeadline(time.Now().Add(s.writeDuration))
	if err != nil {
		return fmt.Errorf("%s: unable to set write deadline: %w", op, err)
	}

	_, err = s.connection.Write(data)
	if err != nil {
		return fmt.Errorf("%s: unable to write to socket: %w", op, err)
	}

	return nil
}
