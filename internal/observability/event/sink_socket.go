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

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
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
	logger         hclog.Logger
}

// NewSocketSink should be used to create a new SocketSink.
// Accepted options: WithMaxDuration and WithSocketType.
func NewSocketSink(address string, format string, opt ...Option) (*SocketSink, error) {
	address = strings.TrimSpace(address)
	if address == "" {
		return nil, fmt.Errorf("address is required: %w", ErrInvalidParameter)
	}

	format = strings.TrimSpace(format)
	if format == "" {
		return nil, fmt.Errorf("format is required: %w", ErrInvalidParameter)
	}

	opts, err := getOpts(opt...)
	if err != nil {
		return nil, err
	}

	sink := &SocketSink{
		requiredFormat: format,
		address:        address,
		socketType:     opts.withSocketType,
		maxDuration:    opts.withMaxDuration,
		socketLock:     sync.RWMutex{},
		connection:     nil,
		logger:         opts.withLogger,
	}

	return sink, nil
}

// Process handles writing the event to the socket.
func (s *SocketSink) Process(ctx context.Context, e *eventlogger.Event) (_ *eventlogger.Event, retErr error) {
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
			s.logger.Error("socket sink error", "context", err, "error", retErr)
		}
	}()

	if e == nil {
		return nil, fmt.Errorf("event is nil: %w", ErrInvalidParameter)
	}

	formatted, found := e.Format(s.requiredFormat)
	if !found {
		return nil, fmt.Errorf("unable to retrieve event formatted as %q: %w", s.requiredFormat, ErrInvalidParameter)
	}

	// Wait for the lock, but ensure we check for a cancelled context as soon as
	// we have it, as there's no point in continuing if we're cancelled.
	s.socketLock.Lock()
	defer s.socketLock.Unlock()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
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
		err = fmt.Errorf("error writing to socket %q: %w", s.address, err)
	}

	// return nil for the event to indicate the pipeline is complete.
	return nil, err
}

// Reopen handles reopening the connection for the socket sink.
func (s *SocketSink) Reopen() error {
	s.socketLock.Lock()
	defer s.socketLock.Unlock()

	err := s.reconnect(nil)
	if err != nil {
		return fmt.Errorf("error reconnecting %q: %w", s.address, err)
	}

	return nil
}

// Type describes the type of this node (sink).
func (_ *SocketSink) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeSink
}

// connect attempts to establish a connection using the socketType and address.
// NOTE: connect is context aware and will not attempt to connect if the context is 'done'.
func (s *SocketSink) connect(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// If we're already connected, we should have disconnected first.
	if s.connection != nil {
		return nil
	}

	timeoutContext, cancel := context.WithTimeout(ctx, s.maxDuration)
	defer cancel()

	dialer := net.Dialer{}
	conn, err := dialer.DialContext(timeoutContext, s.socketType, s.address)
	if err != nil {
		return fmt.Errorf("error connecting to %q address %q: %w", s.socketType, s.address, err)
	}

	s.connection = conn

	return nil
}

// disconnect attempts to close and clear an existing connection.
func (s *SocketSink) disconnect() error {
	// If we're already disconnected, we can return early.
	if s.connection == nil {
		return nil
	}

	err := s.connection.Close()
	if err != nil {
		return fmt.Errorf("error closing connection to %q address %q: %w", s.socketType, s.address, err)
	}
	s.connection = nil

	return nil
}

// reconnect attempts to disconnect and then connect to the configured socketType and address.
func (s *SocketSink) reconnect(ctx context.Context) error {
	err := s.disconnect()
	if err != nil {
		return err
	}

	err = s.connect(ctx)
	if err != nil {
		return err
	}

	return nil
}

// write attempts to write the specified data using the established connection.
func (s *SocketSink) write(ctx context.Context, data []byte) error {
	// Ensure we're connected.
	err := s.connect(ctx)
	if err != nil {
		return err
	}

	err = s.connection.SetWriteDeadline(time.Now().Add(s.maxDuration))
	if err != nil {
		return fmt.Errorf("unable to set write deadline: %w", err)
	}

	_, err = s.connection.Write(data)
	if err != nil {
		return fmt.Errorf("unable to write to socket: %w", err)
	}

	return nil
}
