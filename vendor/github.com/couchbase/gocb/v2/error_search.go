package gocb

import (
	"encoding/json"
)

// SearchError is the error type of all search query errors.
// UNCOMMITTED: This API may change in the future.
type SearchError struct {
	InnerError     error         `json:"-"`
	Query          interface{}   `json:"query,omitempty"`
	Endpoint       string        `json:"endpoint,omitempty"`
	RetryReasons   []RetryReason `json:"retry_reasons,omitempty"`
	RetryAttempts  uint32        `json:"retry_attempts,omitempty"`
	ErrorText      string        `json:"error_text"`
	IndexName      string        `json:"index_name,omitempty"`
	HTTPStatusCode int           `json:"http_status_code,omitempty"`
}

// MarshalJSON implements the Marshaler interface.

func (e SearchError) MarshalJSON() ([]byte, error) {
	var innerError string
	if e.InnerError != nil {
		innerError = e.InnerError.Error()
	}
	return json.Marshal(struct {
		InnerError     string        `json:"msg,omitempty"`
		IndexName      string        `json:"index_name,omitempty"`
		Query          interface{}   `json:"query,omitempty"`
		ErrorText      string        `json:"error_text"`
		Endpoint       string        `json:"endpoint,omitempty"`
		RetryReasons   []RetryReason `json:"retry_reasons,omitempty"`
		RetryAttempts  uint32        `json:"retry_attempts,omitempty"`
		HTTPStatusCode int           `json:"http_status_code,omitempty"`
	}{
		InnerError:     innerError,
		IndexName:      e.IndexName,
		Query:          e.Query,
		ErrorText:      e.ErrorText,
		Endpoint:       e.Endpoint,
		RetryReasons:   e.RetryReasons,
		RetryAttempts:  e.RetryAttempts,
		HTTPStatusCode: e.HTTPStatusCode,
	})
}

// Error returns the string representation of this error.
func (e SearchError) Error() string {
	errBytes, serErr := json.Marshal(struct {
		InnerError       error         `json:"-"`
		IndexName        string        `json:"index_name,omitempty"`
		Query            interface{}   `json:"query,omitempty"`
		ErrorText        string        `json:"error_text"`
		HTTPResponseCode int           `json:"status_code,omitempty"`
		Endpoint         string        `json:"endpoint,omitempty"`
		RetryReasons     []RetryReason `json:"retry_reasons,omitempty"`
		RetryAttempts    uint32        `json:"retry_attempts,omitempty"`
		HTTPStatusCode   int           `json:"http_status_code,omitempty"`
	}{
		InnerError:     e.InnerError,
		IndexName:      e.IndexName,
		Query:          e.Query,
		ErrorText:      e.ErrorText,
		Endpoint:       e.Endpoint,
		RetryReasons:   e.RetryReasons,
		RetryAttempts:  e.RetryAttempts,
		HTTPStatusCode: e.HTTPStatusCode,
	})
	if serErr != nil {
		logErrorf("failed to serialize error to json: %s", serErr.Error())
	}

	return e.InnerError.Error() + " | " + string(errBytes)
}

// Unwrap returns the underlying cause for this error.
func (e SearchError) Unwrap() error {
	return e.InnerError
}
