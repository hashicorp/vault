package gocb

import (
	gocbcore "github.com/couchbase/gocbcore/v9"
	"github.com/pkg/errors"
)

// HTTPError is the error type of management HTTP errors.
// UNCOMMITTED: This API may change in the future.
type HTTPError struct {
	InnerError    error         `json:"-"`
	UniqueID      string        `json:"unique_id,omitempty"`
	Endpoint      string        `json:"endpoint,omitempty"`
	RetryReasons  []RetryReason `json:"retry_reasons,omitempty"`
	RetryAttempts uint32        `json:"retry_attempts,omitempty"`
}

// Error returns the string representation of this error.
func (e HTTPError) Error() string {
	return e.InnerError.Error() + " | " + serializeWrappedError(e)
}

// Unwrap returns the underlying cause for this error.
func (e HTTPError) Unwrap() error {
	return e.InnerError
}

func makeGenericHTTPError(baseErr error, req *gocbcore.HTTPRequest, resp *gocbcore.HTTPResponse) error {
	if baseErr == nil {
		logErrorf("makeGenericHTTPError got an empty error")
		baseErr = errors.New("unknown error")
	}

	err := HTTPError{
		InnerError: baseErr,
	}

	if req != nil {
		err.UniqueID = req.UniqueID
	}

	if resp != nil {
		err.Endpoint = resp.Endpoint
	}

	return err
}

func makeGenericMgmtError(baseErr error, req *mgmtRequest, resp *mgmtResponse) error {
	if baseErr == nil {
		logErrorf("makeGenericMgmtError got an empty error")
		baseErr = errors.New("unknown error")
	}

	err := HTTPError{
		InnerError: baseErr,
	}

	if req != nil {
		err.UniqueID = req.UniqueID
	}

	if resp != nil {
		err.Endpoint = resp.Endpoint
	}

	return err
}

func makeMgmtBadStatusError(message string, req *mgmtRequest, resp *mgmtResponse) error {
	return makeGenericMgmtError(errors.New(message), req, resp)
}
