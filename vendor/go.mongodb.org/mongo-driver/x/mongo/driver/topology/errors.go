// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo/description"
)

// ConnectionError represents a connection error.
type ConnectionError struct {
	ConnectionID string
	Wrapped      error

	// init will be set to true if this error occurred during connection initialization or
	// during a connection handshake.
	init    bool
	message string
}

// Error implements the error interface.
func (e ConnectionError) Error() string {
	message := e.message
	if e.init {
		fullMsg := "error occurred during connection handshake"
		if message != "" {
			fullMsg = fmt.Sprintf("%s: %s", fullMsg, message)
		}
		message = fullMsg
	}
	if e.Wrapped != nil && message != "" {
		return fmt.Sprintf("connection(%s) %s: %s", e.ConnectionID, message, e.Wrapped.Error())
	}
	if e.Wrapped != nil {
		return fmt.Sprintf("connection(%s) %s", e.ConnectionID, e.Wrapped.Error())
	}
	return fmt.Sprintf("connection(%s) %s", e.ConnectionID, message)
}

// Unwrap returns the underlying error.
func (e ConnectionError) Unwrap() error {
	return e.Wrapped
}

// ServerSelectionError represents a Server Selection error.
type ServerSelectionError struct {
	Desc    description.Topology
	Wrapped error
}

// Error implements the error interface.
func (e ServerSelectionError) Error() string {
	if e.Wrapped != nil {
		return fmt.Sprintf("server selection error: %s, current topology: { %s }", e.Wrapped.Error(), e.Desc.String())
	}
	return fmt.Sprintf("server selection error: current topology: { %s }", e.Desc.String())
}

// Unwrap returns the underlying error.
func (e ServerSelectionError) Unwrap() error {
	return e.Wrapped
}

// WaitQueueTimeoutError represents a timeout when requesting a connection from the pool
type WaitQueueTimeoutError struct {
	Wrapped              error
	pinnedConnections    *pinnedConnections
	maxPoolSize          uint64
	totalConnections     int
	availableConnections int
	waitDuration         time.Duration
}

type pinnedConnections struct {
	cursorConnections      uint64
	transactionConnections uint64
}

// Error implements the error interface.
func (w WaitQueueTimeoutError) Error() string {
	errorMsg := "timed out while checking out a connection from connection pool"
	switch {
	case w.Wrapped == nil:
	case errors.Is(w.Wrapped, context.Canceled):
		errorMsg = fmt.Sprintf(
			"%s: %s",
			"canceled while checking out a connection from connection pool",
			w.Wrapped.Error(),
		)
	default:
		errorMsg = fmt.Sprintf(
			"%s: %s",
			errorMsg,
			w.Wrapped.Error(),
		)
	}

	msg := fmt.Sprintf("%s; total connections: %d, maxPoolSize: %d, ", errorMsg, w.totalConnections, w.maxPoolSize)
	if pinnedConnections := w.pinnedConnections; pinnedConnections != nil {
		openConnectionCount := uint64(w.totalConnections) -
			pinnedConnections.cursorConnections -
			pinnedConnections.transactionConnections
		msg += fmt.Sprintf("connections in use by cursors: %d, connections in use by transactions: %d, connections in use by other operations: %d, ",
			pinnedConnections.cursorConnections,
			pinnedConnections.transactionConnections,
			openConnectionCount,
		)
	}
	msg += fmt.Sprintf("idle connections: %d, wait duration: %s", w.availableConnections, w.waitDuration.String())
	return msg
}

// Unwrap returns the underlying error.
func (w WaitQueueTimeoutError) Unwrap() error {
	return w.Wrapped
}
