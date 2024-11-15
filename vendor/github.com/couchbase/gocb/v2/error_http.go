package gocb

import (
	"encoding/json"
	"errors"
	"github.com/couchbase/gocbcore/v10"
	"io"
)

// HTTPError is the error type of management HTTP errors.
// UNCOMMITTED: This API may change in the future.
type HTTPError struct {
	InnerError    error         `json:"-"`
	UniqueID      string        `json:"unique_id,omitempty"`
	Endpoint      string        `json:"endpoint,omitempty"`
	RetryReasons  []RetryReason `json:"retry_reasons,omitempty"`
	RetryAttempts uint32        `json:"retry_attempts,omitempty"`
	ErrorText     string        `json:"error_text,omitempty"`
	StatusCode    uint32        `json:"status_code,omitempty"`
}

// MarshalJSON implements the Marshaler interface.
func (e HTTPError) MarshalJSON() ([]byte, error) {
	var innerError string
	if e.InnerError != nil {
		innerError = e.InnerError.Error()
	}
	return json.Marshal(struct {
		InnerError    string        `json:"msg,omitempty"`
		UniqueID      string        `json:"unique_id,omitempty"`
		Endpoint      string        `json:"endpoint,omitempty"`
		RetryReasons  []RetryReason `json:"retry_reasons,omitempty"`
		RetryAttempts uint32        `json:"retry_attempts,omitempty"`
		ErrorText     string        `json:"error_text,omitempty"`
		StatusCode    uint32        `json:"status_code,omitempty"`
	}{
		InnerError:    innerError,
		UniqueID:      e.UniqueID,
		Endpoint:      e.Endpoint,
		RetryReasons:  e.RetryReasons,
		RetryAttempts: e.RetryAttempts,
		ErrorText:     e.ErrorText,
		StatusCode:    e.StatusCode,
	})
}

// Error returns the string representation of this error.
func (e HTTPError) Error() string {
	errBytes, serErr := json.Marshal(struct {
		InnerError    error         `json:"-"`
		UniqueID      string        `json:"unique_id,omitempty"`
		Endpoint      string        `json:"endpoint,omitempty"`
		RetryReasons  []RetryReason `json:"retry_reasons,omitempty"`
		RetryAttempts uint32        `json:"retry_attempts,omitempty"`
		ErrorText     string        `json:"error_text,omitempty"`
		StatusCode    uint32        `json:"status_code,omitempty"`
	}{
		InnerError:    e.InnerError,
		UniqueID:      e.UniqueID,
		Endpoint:      e.Endpoint,
		RetryReasons:  e.RetryReasons,
		RetryAttempts: e.RetryAttempts,
		ErrorText:     e.ErrorText,
		StatusCode:    e.StatusCode,
	})
	if serErr != nil {
		logErrorf("failed to serialize error to json: %s", serErr.Error())
	}

	return e.InnerError.Error() + " | " + string(errBytes)
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

	err := &HTTPError{
		InnerError: baseErr,
	}

	if req != nil {
		err.UniqueID = req.UniqueID
	}

	if resp != nil {
		err.Endpoint = resp.Endpoint
		err.StatusCode = uint32(resp.StatusCode)
	}

	return err
}

func makeGenericMgmtError(baseErr error, req *mgmtRequest, resp *mgmtResponse, errText string) error {
	if baseErr == nil {
		logErrorf("makeGenericMgmtError got an empty error")
		baseErr = errors.New("unknown error")
	}

	err := &HTTPError{
		InnerError: baseErr,
		ErrorText:  errText,
	}

	if req != nil {
		err.UniqueID = req.UniqueID
	}

	if resp != nil {
		err.Endpoint = resp.Endpoint
		err.StatusCode = resp.StatusCode
	}

	return err
}

func makeMgmtBadStatusError(message string, req *mgmtRequest, resp *mgmtResponse) error {
	var errText string
	if resp != nil {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			logDebugf("failed to read http body: %s", err)
			return nil
		}

		errText = string(b)
	}
	return makeGenericMgmtError(errors.New(message), req, resp, errText)
}
