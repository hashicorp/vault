package topology

import "fmt"

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
	if e.Wrapped != nil {
		return fmt.Sprintf("connection(%s) %s: %s", e.ConnectionID, e.message, e.Wrapped.Error())
	}
	return fmt.Sprintf("connection(%s) %s", e.ConnectionID, e.message)
}

// Unwrap returns the underlying error.
func (e ConnectionError) Unwrap() error {
	return e.Wrapped
}
