package audit

import (
	"errors"
	"fmt"
)

var (
	ErrFilterParameter    = errors.New("filter parameter")
	ErrFallbackParameter  = errors.New("fallback parameter")
	ErrInvalidParameter   = errors.New("invalid parameter")
	ErrEnterpriseOnly     = errors.New("enterprise-only")
	ErrConfiguration      = errors.New("configuration error")
	ErrConflict           = errors.New("audit conflict")
	ErrUnknown            = errors.New("unknown error")
	ErrPersistence        = errors.New("persistence error")
	ErrBrokerRegistration = errors.New("registration error")
)

// Error represents an error within the audit subsystem.
type Error struct {
	op  string
	msg string
	err error
	// detail should be used to store any extra error information that may
	// have been returned from other upstream function calls.
	detail  error
	wrapped *Error
}

// NewError should be used to create instances of an Error.
func NewError(operation string, message string, err error) *Error {
	return &Error{
		op:  operation,
		msg: message,
		err: err,
	}
}

// Wrap should be used to wrap an upstream Error.
// The original Error is returned after wrapping the upstream error.
func (e *Error) Wrap(e2 *Error) *Error {
	e.wrapped = e2

	return e
}

// Detail is used to set a detailed Go error which came from upstream.
func (e *Error) Detail(err error) *Error {
	e.detail = err

	return e
}

// Error satisfies the Go Error interface.
func (e *Error) Error() string {
	return e.Internal().Error()
}

// String satisfies the Go Stringer interface
func (e *Error) String() string {
	return e.Error()
}

// Internal should be used when an error is required for internal systems such as logs.
func (e *Error) Internal() error {
	err := fmt.Errorf("%s: %s: %w", e.op, e.msg, e.err)

	if e.detail != nil {
		err = fmt.Errorf("%w: %w", err, e.detail)
	}

	if e.wrapped != nil {
		err = fmt.Errorf("%w: %w", err, e.wrapped.Internal())
	}

	return err
}

// External should be used when an error is required for external systems, such
// as API responses.
func (e *Error) External() error {
	err := fmt.Errorf("%s: %w", e.msg, e.err)

	if e.wrapped != nil {
		err = fmt.Errorf("%w: %w", err, e.wrapped.External())
	}

	return err
}
