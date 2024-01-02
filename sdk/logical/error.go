// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"context"
	"errors"
)

var (
	// ErrUnsupportedOperation is returned if the operation is not supported
	// by the logical backend.
	ErrUnsupportedOperation = errors.New("unsupported operation")

	// ErrUnsupportedPath is returned if the path is not supported
	// by the logical backend.
	ErrUnsupportedPath = errors.New("unsupported path")

	// ErrInvalidRequest is returned if the request is invalid
	ErrInvalidRequest = errors.New("invalid request")

	// ErrPermissionDenied is returned if the client is not authorized
	ErrPermissionDenied = errors.New("permission denied")

	// ErrInvalidCredentials is returned when the provided credentials are incorrect
	// This is used internally for user lockout purposes. This is not seen externally.
	// The status code returned does not change because of this error
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrMultiAuthzPending is returned if the the request needs more
	// authorizations
	ErrMultiAuthzPending = errors.New("request needs further approval")

	// ErrUpstreamRateLimited is returned when Vault receives a rate limited
	// response from an upstream
	ErrUpstreamRateLimited = errors.New("upstream rate limited")

	// ErrPerfStandbyForward is returned when Vault is in a state such that a
	// perf standby cannot satisfy a request
	ErrPerfStandbyPleaseForward = errors.New("please forward to the active node")

	// ErrLeaseCountQuotaExceeded is returned when a request is rejected due to a lease
	// count quota being exceeded.
	ErrLeaseCountQuotaExceeded = errors.New("lease count quota exceeded")

	// ErrRateLimitQuotaExceeded is returned when a request is rejected due to a
	// rate limit quota being exceeded.
	ErrRateLimitQuotaExceeded = errors.New("rate limit quota exceeded")

	// ErrUnrecoverable is returned when a request fails due to something that
	// is likely to require manual intervention. This is a generic form of an
	// unrecoverable error.
	// e.g.: misconfigured or disconnected storage backend.
	ErrUnrecoverable = errors.New("unrecoverable error")

	// ErrMissingRequiredState is returned when a request can't be satisfied
	// with the data in the local node's storage, based on the provided
	// X-Vault-Index request header.
	ErrMissingRequiredState = errors.New("required index state not present")

	// Error indicating that the requested path used to serve a purpose in older
	// versions, but the functionality has now been removed
	ErrPathFunctionalityRemoved = errors.New("functionality on this path has been removed")

	// ErrNotFound is an error used to indicate that a particular resource was
	// not found.
	ErrNotFound = errors.New("not found")
)

type DelegatedAuthErrorHandler func(ctx context.Context, initiatingRequest, authRequest *Request, authResponse *Response, err error) (*Response, error)

var _ error = &RequestDelegatedAuthError{}

// RequestDelegatedAuthError Special error indicating the backend wants to delegate authentication elsewhere
type RequestDelegatedAuthError struct {
	mountAccessor string
	path          string
	data          map[string]interface{}
	errHandler    DelegatedAuthErrorHandler
}

func NewDelegatedAuthenticationRequest(mountAccessor, path string, data map[string]interface{}, errHandler DelegatedAuthErrorHandler) *RequestDelegatedAuthError {
	return &RequestDelegatedAuthError{
		mountAccessor: mountAccessor,
		path:          path,
		data:          data,
		errHandler:    errHandler,
	}
}

func (d *RequestDelegatedAuthError) Error() string {
	return "authentication delegation requested"
}

func (d *RequestDelegatedAuthError) MountAccessor() string {
	return d.mountAccessor
}

func (d *RequestDelegatedAuthError) Path() string {
	return d.path
}

func (d *RequestDelegatedAuthError) Data() map[string]interface{} {
	return d.data
}

func (d *RequestDelegatedAuthError) AuthErrorHandler() DelegatedAuthErrorHandler {
	return d.errHandler
}

type HTTPCodedError interface {
	Error() string
	Code() int
}

func CodedError(status int, msg string) HTTPCodedError {
	return &codedError{
		Status:  status,
		Message: msg,
	}
}

var _ HTTPCodedError = (*codedError)(nil)

type codedError struct {
	Status  int
	Message string
}

func (e *codedError) Error() string {
	return e.Message
}

func (e *codedError) Code() int {
	return e.Status
}

// Struct to identify user input errors.  This is helpful in responding the
// appropriate status codes to clients from the HTTP endpoints.
type StatusBadRequest struct {
	Err string
}

// Implementing error interface
func (s *StatusBadRequest) Error() string {
	return s.Err
}

// This is a new type declared to not cause potential compatibility problems if
// the logic around the CodedError changes; in particular for logical request
// paths it is basically ignored, and changing that behavior might cause
// unforeseen issues.
type ReplicationCodedError struct {
	Msg  string
	Code int
}

func (r *ReplicationCodedError) Error() string {
	return r.Msg
}

type KeyNotFoundError struct {
	Err error
}

func (e *KeyNotFoundError) WrappedErrors() []error {
	return []error{e.Err}
}

func (e *KeyNotFoundError) Error() string {
	return e.Err.Error()
}
