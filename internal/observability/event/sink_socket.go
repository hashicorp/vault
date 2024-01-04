// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package event

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"

	"github.com/hashicorp/eventlogger"
)

var _ eventlogger.Node = (*SocketSink)(nil)

// SocketSink is a sink node which handles writing events to socket.
type SocketSink struct {
	requiredFormat string
	address        string
	socketType     string
	maxDuration    time.Duration
	socketLock     sync.RWMutex
	connection     net.Conn
}

// NewSocketSink should be used to create a new SocketSink.
// Accepted options: WithMaxDuration and WithSocketType.
func NewSocketSink(address string, format string, opt ...Option) (*SocketSink, error) {
	const op = "event.NewSocketSink"

	address = strings.TrimSpace(address)
	if address == "" {
		return nil, fmt.Errorf("%s: address is required: %w", op, ErrInvalidParameter)
	}

	format = strings.TrimSpace(format)
	if format == "" {
		return nil, fmt.Errorf("%s: format is required: %w", op, ErrInvalidParameter)
	}

	opts, err := getOpts(opt...)
	if err != nil {
		return nil, fmt.Errorf("%s: error applying options: %w", op, err)
	}

	sink := &SocketSink{
		requiredFormat: format,
		address:        address,
		socketType:     opts.withSocketType,
		maxDuration:    opts.withMaxDuration,
		socketLock:     sync.RWMutex{},
		connection:     nil,
	}

	return sink, nil
}

// Process handles writing the event to the socket.
func (s *SocketSink) Process(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	const op = "event.(SocketSink).Process"

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.socketLock.Lock()
	defer s.socketLock.Unlock()

	if e == nil {
		return nil, fmt.Errorf("%s: event is nil: %w", op, ErrInvalidParameter)
	}

	formatted, found := e.Format(s.requiredFormat)
	if !found {
		return nil, fmt.Errorf("%s: unable to retrieve event formatted as %q", op, s.requiredFormat)
	}

	// Try writing and return early if successful.
	err := s.write(ctx, formatted)
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
		err = s.write(ctx, formatted)
	}

	// Format the error nicely if we need to return one.
	if err != nil {
		err = fmt.Errorf("%s: error writing to socket: %w", op, err)
	}

	// return nil for the event to indicate the pipeline is complete.
	return nil, err
}

// Reopen handles reopening the connection for the socket sink.
func (s *SocketSink) Reopen() error {
	const op = "event.(SocketSink).Reopen"

	s.socketLock.Lock()
	defer s.socketLock.Unlock()

	err := s.reconnect(nil)
	if err != nil {
		return fmt.Errorf("%s: error reconnecting: %w", op, err)
	}

	return nil
}

// Type describes the type of this node (sink).
func (_ *SocketSink) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeSink
}

// connect attempts to establish a connection using the socketType and address.
func (s *SocketSink) connect(ctx context.Context) error {
	const op = "event.(SocketSink).connect"

	// If we're already connected, we should have disconnected first.
	if s.connection != nil {
		return nil
	}

	timeoutContext, cancel := context.WithTimeout(ctx, s.maxDuration)
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
func (s *SocketSink) disconnect() error {
	const op = "event.(SocketSink).disconnect"

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
func (s *SocketSink) reconnect(ctx context.Context) error {
	const op = "event.(SocketSink).reconnect"

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
func (s *SocketSink) write(ctx context.Context, data []byte) error {
	const op = "event.(SocketSink).write"

	// Ensure we're connected.
	err := s.connect(ctx)
	if err != nil {
		return fmt.Errorf("%s: connection error: %w", op, err)
	}

	err = s.connection.SetWriteDeadline(time.Now().Add(s.maxDuration))
	if err != nil {
		return fmt.Errorf("%s: unable to set write deadline: %w", op, err)
	}

	_, err = s.connection.Write(data)
	if err != nil {
		return fmt.Errorf("%s: unable to write to socket: %w", op, err)
	}

	return nil
}
