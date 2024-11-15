// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"
)

// UnexpectedResponseError tracks the components for API errors encountered when
// requireOK and requireStatusIn's conditions are not met.
type UnexpectedResponseError struct {
	expected   []int
	statusCode int
	statusText string
	body       string
	err        error
	additional error
}

func (e UnexpectedResponseError) HasExpectedStatuses() bool { return len(e.expected) > 0 }
func (e UnexpectedResponseError) ExpectedStatuses() []int   { return e.expected }
func (e UnexpectedResponseError) HasStatusCode() bool       { return e.statusCode != 0 }
func (e UnexpectedResponseError) StatusCode() int           { return e.statusCode }
func (e UnexpectedResponseError) HasStatusText() bool       { return e.statusText != "" }
func (e UnexpectedResponseError) StatusText() string        { return e.statusText }
func (e UnexpectedResponseError) HasBody() bool             { return e.body != "" }
func (e UnexpectedResponseError) Body() string              { return e.body }
func (e UnexpectedResponseError) HasError() bool            { return e.err != nil }
func (e UnexpectedResponseError) Unwrap() error             { return e.err }
func (e UnexpectedResponseError) HasAdditional() bool       { return e.additional != nil }
func (e UnexpectedResponseError) Additional() error         { return e.additional }
func newUnexpectedResponseError(src unexpectedResponseErrorSource, opts ...unexpectedResponseErrorOption) UnexpectedResponseError {
	nErr := src()
	for _, opt := range opts {
		opt(nErr)
	}
	if nErr.statusText == "" {
		// the stdlib's http.StatusText function is a good place to start
		nErr.statusFromCode(http.StatusText)
	}

	return *nErr
}

// Use textual representation of the given integer code. Called when status text
// is not set using the WithStatusText option.
func (e UnexpectedResponseError) statusFromCode(f func(int) string) {
	e.statusText = f(e.statusCode)
	if !e.HasStatusText() {
		e.statusText = "unknown status code"
	}
}

func (e UnexpectedResponseError) Error() string {
	var eTxt strings.Builder
	eTxt.WriteString("Unexpected response code")
	if e.HasBody() || e.HasStatusCode() {
		eTxt.WriteString(": ")
	}
	if e.HasStatusCode() {
		eTxt.WriteString(fmt.Sprint(e.statusCode))
		if e.HasBody() {
			eTxt.WriteRune(' ')
		}
	}
	if e.HasBody() {
		eTxt.WriteString(fmt.Sprintf("(%s)", e.body))
	}

	if e.HasAdditional() {
		eTxt.WriteString(fmt.Sprintf(". Additionally, an error occurred while constructing this error (%s); the body might be truncated or missing.", e.additional.Error()))
	}

	return eTxt.String()
}

// UnexpectedResponseErrorOptions are functions passed to NewUnexpectedResponseError
// to customize the created error.
type unexpectedResponseErrorOption func(*UnexpectedResponseError)

// withError allows the addition of a Go error that may have been encountered
// while processing the response. For example, if there is an error constructing
// the gzip reader to process a gzip-encoded response body.
func withError(e error) unexpectedResponseErrorOption {
	return func(u *UnexpectedResponseError) { u.err = e }
}

// withBody overwrites the Body value with the provided custom value
func withBody(b string) unexpectedResponseErrorOption {
	return func(u *UnexpectedResponseError) { u.body = b }
}

// withStatusText overwrites the StatusText value the provided custom value
func withStatusText(st string) unexpectedResponseErrorOption {
	return func(u *UnexpectedResponseError) { u.statusText = st }
}

// withExpectedStatuses provides a list of statuses that the receiving function
// expected to receive. This can be used by API callers to provide more feedback
// to end-users.
func withExpectedStatuses(s []int) unexpectedResponseErrorOption {
	return func(u *UnexpectedResponseError) { u.expected = slices.Clone(s) }
}

// unexpectedResponseErrorSource provides the basis for a NewUnexpectedResponseError.
type unexpectedResponseErrorSource func() *UnexpectedResponseError

// fromHTTPResponse read an open HTTP response, drains and closes its body as
// the data for the UnexpectedResponseError.
func fromHTTPResponse(resp *http.Response) unexpectedResponseErrorSource {
	return func() *UnexpectedResponseError {
		u := new(UnexpectedResponseError)

		if resp != nil {
			// collect and close the body
			var buf bytes.Buffer
			if _, e := io.Copy(&buf, resp.Body); e != nil {
				u.additional = e
			}

			// Body has been tested as safe to close more than once
			_ = resp.Body.Close()
			body := strings.TrimSpace(buf.String())

			// make and return the error
			u.statusCode = resp.StatusCode
			u.statusText = strings.TrimSpace(strings.TrimPrefix(resp.Status, fmt.Sprint(resp.StatusCode)))
			u.body = body
		}
		return u
	}
}

// fromStatusCode attempts to resolve the status code to status text using
// the resolving function provided inside of the NewUnexpectedResponseError
// implementation.
func fromStatusCode(sc int) unexpectedResponseErrorSource {
	return func() *UnexpectedResponseError { return &UnexpectedResponseError{statusCode: sc} }
}

// doRequestWrapper is a function that wraps the client's doRequest method
// and can be used to provide error and response handling
type doRequestWrapper = func(time.Duration, *http.Response, error) (time.Duration, *http.Response, error)

// requireOK is used to wrap doRequest and check for a 200
func requireOK(d time.Duration, resp *http.Response, e error) (time.Duration, *http.Response, error) {
	f := requireStatusIn(http.StatusOK)
	return f(d, resp, e)
}

// requireStatusIn is a doRequestWrapper generator that takes expected HTTP
// response codes and validates that the received response code is among them
func requireStatusIn(statuses ...int) doRequestWrapper {
	return func(d time.Duration, resp *http.Response, e error) (time.Duration, *http.Response, error) {
		if e != nil {
			if resp != nil {
				_ = resp.Body.Close()
			}
			return d, nil, e
		}

		for _, status := range statuses {
			if resp.StatusCode == status {
				return d, resp, nil
			}
		}

		return d, nil, newUnexpectedResponseError(fromHTTPResponse(resp), withExpectedStatuses(statuses))
	}
}
