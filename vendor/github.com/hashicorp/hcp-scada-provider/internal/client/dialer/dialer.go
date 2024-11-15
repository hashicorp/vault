// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dialer

import (
	"context"
	"net"
)

// Dialer is the interface for dialing a SCADA broker.
type Dialer interface {
	// Dial takes a string identifying the broker to dial and returns a
	// net.Conn that is the connection to SCADA or an error if it couldn't be
	// dialed.
	Dial(string) (net.Conn, error)
	// DialContext has the same functionality Dial has but operates with a context.
	DialContext(context.Context, string) (net.Conn, error)
}

// TransientError marks dial errors created from scada dialers as being
// transient.
type TransientError struct {
	Err error
}

// Error implements the error interface.
func (e *TransientError) Error() string {
	return e.Err.Error()
}

// Unwrap supports unwrapping for use with errors.Is and errors.As.
func (e *TransientError) Unwrap() error {
	return e.Err
}

// NewTransientError returns new transient error.
func NewTransientError(err error) error {
	return &TransientError{Err: err}
}
