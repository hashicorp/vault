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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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
	// trailers are the trailers returned in the response, if any.
	trailers metadata.MD
	// additionalInformation optionally contains any additional information
	// about the error.
	additionalInformation string
}

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

// toSpannerError converts general Go error to *spanner.Error.
func toSpannerError(err error) error {
	return toSpannerErrorWithMetadata(err, nil)
}

// toSpannerErrorWithMetadata converts general Go error and grpc trailers to
// *spanner.Error.
//
// Note: modifies original error if trailers aren't nil.
func toSpannerErrorWithMetadata(err error, trailers metadata.MD) error {
	if err == nil {
		return nil
	}
	var se *Error
	if errorAs(err, &se) {
		if trailers != nil {
			se.trailers = metadata.Join(se.trailers, trailers)
		}
		return se
	}
	switch {
	case err == context.DeadlineExceeded || err == context.Canceled:
		return &Error{status.FromContextError(err).Code(), status.FromContextError(err).Err(), err.Error(), trailers, ""}
	case status.Code(err) == codes.Unknown:
		return &Error{codes.Unknown, err, err.Error(), trailers, ""}
	default:
		return &Error{status.Convert(err).Code(), err, status.Convert(err).Message(), trailers, ""}
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

// errTrailers extracts the grpc trailers if present from a Go error.
func errTrailers(err error) metadata.MD {
	var se *Error
	if !errorAs(err, &se) {
		return nil
	}
	return se.trailers
}
