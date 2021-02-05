/*
Copyright 2017 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package spanner

import (
	"context"
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error is the structured error returned by Cloud Spanner client.
type Error struct {
	// Code is the canonical error code for describing the nature of a
	// particular error.
	//
	// Deprecated: The error code should be extracted from the wrapped error by
	// calling ErrCode(err error). This field will be removed in a future
	// release.
	Code codes.Code
	// err is the wrapped error that caused this Spanner error. The wrapped
	// error can be read with the Unwrap method.
	err error
	// Desc explains more details of the error.
	Desc string
	// additionalInformation optionally contains any additional information
	// about the error.
	additionalInformation string
}

// TransactionOutcomeUnknownError is wrapped in a Spanner error when the error
// occurred during a transaction, and the outcome of the transaction is
// unknown as a result of the error. This could be the case if a timeout or
// canceled error occurs after a Commit request has been sent, but before the
// client has received a response from the server.
type TransactionOutcomeUnknownError struct {
	// err is the wrapped error that caused this TransactionOutcomeUnknownError
	// error. The wrapped error can be read with the Unwrap method.
	err error
}

const transactionOutcomeUnknownMsg = "transaction outcome unknown"

// Error implements error.Error.
func (*TransactionOutcomeUnknownError) Error() string { return transactionOutcomeUnknownMsg }

// Unwrap returns the wrapped error (if any).
func (e *TransactionOutcomeUnknownError) Unwrap() error { return e.err }

// Error implements error.Error.
func (e *Error) Error() string {
	if e == nil {
		return fmt.Sprintf("spanner: OK")
	}
	code := ErrCode(e)
	if e.additionalInformation == "" {
		return fmt.Sprintf("spanner: code = %q, desc = %q", code, e.Desc)
	}
	return fmt.Sprintf("spanner: code = %q, desc = %q, additional information = %s", code, e.Desc, e.additionalInformation)
}

// Unwrap returns the wrapped error (if any).
func (e *Error) Unwrap() error {
	return e.err
}

// GRPCStatus returns the corresponding gRPC Status of this Spanner error.
// This allows the error to be converted to a gRPC status using
// `status.Convert(error)`.
func (e *Error) GRPCStatus() *status.Status {
	err := unwrap(e)
	for {
		// If the base error is nil, return status created from e.Code and e.Desc.
		if err == nil {
			return status.New(e.Code, e.Desc)
		}
		code := status.Code(err)
		if code != codes.Unknown {
			return status.New(code, e.Desc)
		}
		err = unwrap(err)
	}
}

// decorate decorates an existing spanner.Error with more information.
func (e *Error) decorate(info string) {
	e.Desc = fmt.Sprintf("%v, %v", info, e.Desc)
}

// spannerErrorf generates a *spanner.Error with the given description and a
// status error with the given error code as its wrapped error.
func spannerErrorf(code codes.Code, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	wrapped := status.Error(code, msg)
	return &Error{
		Code: code,
		err:  wrapped,
		Desc: msg,
	}
}

// ToSpannerError converts a general Go error to *spanner.Error. If the given
// error is already a *spanner.Error, the original error will be returned.
//
// Spanner Errors are normally created by the Spanner client library from the
// returned status of a RPC. This method can also be used to create Spanner
// errors for use in tests. The recommended way to create test errors is
// calling this method with a status error, e.g.
// ToSpannerError(status.New(codes.NotFound, "Table not found").Err())
func ToSpannerError(err error) error {
	return toSpannerErrorWithCommitInfo(err, false)
}

// toSpannerErrorWithCommitInfo converts general Go error to *spanner.Error
// with additional information if the error occurred during a Commit request.
//
// If err is already a *spanner.Error, err is returned unmodified.
func toSpannerErrorWithCommitInfo(err error, errorDuringCommit bool) error {
	if err == nil {
		return nil
	}
	var se *Error
	if errorAs(err, &se) {
		return se
	}
	switch {
	case err == context.DeadlineExceeded || err == context.Canceled:
		desc := err.Error()
		wrapped := status.FromContextError(err).Err()
		if errorDuringCommit {
			desc = fmt.Sprintf("%s, %s", desc, transactionOutcomeUnknownMsg)
			wrapped = &TransactionOutcomeUnknownError{err: wrapped}
		}
		return &Error{status.FromContextError(err).Code(), wrapped, desc, ""}
	case status.Code(err) == codes.Unknown:
		return &Error{codes.Unknown, err, err.Error(), ""}
	default:
		statusErr := status.Convert(err)
		code, desc := statusErr.Code(), statusErr.Message()
		wrapped := err
		if errorDuringCommit && (code == codes.DeadlineExceeded || code == codes.Canceled) {
			desc = fmt.Sprintf("%s, %s", desc, transactionOutcomeUnknownMsg)
			wrapped = &TransactionOutcomeUnknownError{err: wrapped}
		}
		return &Error{code, wrapped, desc, ""}
	}
}

// ErrCode extracts the canonical error code from a Go error.
func ErrCode(err error) codes.Code {
	s, ok := status.FromError(err)
	if !ok {
		return codes.Unknown
	}
	return s.Code()
}

// ErrDesc extracts the Cloud Spanner error description from a Go error.
func ErrDesc(err error) string {
	var se *Error
	if !errorAs(err, &se) {
		return err.Error()
	}
	return se.Desc
}

// extractResourceType extracts the resource type from any ResourceInfo detail
// included in the error.
func extractResourceType(err error) (string, bool) {
	var s *status.Status
	var se *Error
	if errorAs(err, &se) {
		// Unwrap statusError.
		s = status.Convert(se.Unwrap())
	} else {
		s = status.Convert(err)
	}
	if s == nil {
		return "", false
	}
	for _, detail := range s.Details() {
		if resourceInfo, ok := detail.(*errdetails.ResourceInfo); ok {
			return resourceInfo.ResourceType, true
		}
	}
	return "", false
}
