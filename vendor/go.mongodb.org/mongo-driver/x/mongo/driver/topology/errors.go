package topology

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo/description"
)

// ConnectionError represents a connection error.
type ConnectionError struct {
	ConnectionID string
	Wrapped      error

	// init will be set to true if this error occured during connection initialization or
	// during a connection handshake.
	init    bool
	message string
}

// Error implements the error interface.
func (e ConnectionError) Error() string {
	message := e.message
	if e.init {
		fullMsg := "error occured during connection handshake"
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
	Wrapped                      error
	PinnedCursorConnections      uint64
	PinnedTransactionConnections uint64
	maxPoolSize                  uint64
}

// Error implements the error interface.
func (w WaitQueueTimeoutError) Error() string {
	errorMsg := "timed out while checking out a connection from connection pool"
	if w.Wrapped != nil {
		errorMsg = fmt.Sprintf("%s: %s", errorMsg, w.Wrapped.Error())
	}

	errorMsg = fmt.Sprintf("%s; maxPoolSize: %d, connections in use by cursors: %d, connections in use by transactions: %d",
		errorMsg, w.maxPoolSize, w.PinnedCursorConnections, w.PinnedTransactionConnections)
	return fmt.Sprintf("%s, connections in use by other operations: %d", errorMsg,
		w.maxPoolSize-(w.PinnedCursorConnections+w.PinnedTransactionConnections))
}

// Unwrap returns the underlying error.
func (w WaitQueueTimeoutError) Unwrap() error {
	return w.Wrapped
}
