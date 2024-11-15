package gocb

import (
	"encoding/json"
	gocbcore "github.com/couchbase/gocbcore/v10"
)

// AnalyticsErrorDesc represents a specific error returned from the analytics service.
type AnalyticsErrorDesc struct {
	Code    uint32
	Message string
}

func translateCoreAnalyticsErrorDesc(descs []gocbcore.AnalyticsErrorDesc) []AnalyticsErrorDesc {
	descsOut := make([]AnalyticsErrorDesc, len(descs))
	for descIdx, desc := range descs {
		descsOut[descIdx] = AnalyticsErrorDesc{
			Code:    desc.Code,
			Message: desc.Message,
		}
	}
	return descsOut
}

// AnalyticsError is the error type of all analytics query errors.
// UNCOMMITTED: This API may change in the future.
type AnalyticsError struct {
	InnerError      error                `json:"-"`
	Statement       string               `json:"statement,omitempty"`
	ClientContextID string               `json:"client_context_id,omitempty"`
	Errors          []AnalyticsErrorDesc `json:"errors,omitempty"`
	Endpoint        string               `json:"endpoint,omitempty"`
	RetryReasons    []RetryReason        `json:"retry_reasons,omitempty"`
	RetryAttempts   uint32               `json:"retry_attempts,omitempty"`
	ErrorText       string               `json:"error_text,omitempty"`
	HTTPStatusCode  int                  `json:"http_status_code,omitempty"`
}

// MarshalJSON implements the Marshaler interface.

func (e AnalyticsError) MarshalJSON() ([]byte, error) {
	var innerError string
	if e.InnerError != nil {
		innerError = e.InnerError.Error()
	}
	return json.Marshal(struct {
		InnerError      string               `json:"msg,omitempty"`
		Statement       string               `json:"statement,omitempty"`
		ClientContextID string               `json:"client_context_id,omitempty"`
		Errors          []AnalyticsErrorDesc `json:"errors,omitempty"`
		Endpoint        string               `json:"endpoint,omitempty"`
		RetryReasons    []RetryReason        `json:"retry_reasons,omitempty"`
		RetryAttempts   uint32               `json:"retry_attempts,omitempty"`
		HTTPStatusCode  int                  `json:"http_status_code,omitempty"`
	}{
		InnerError:      innerError,
		Statement:       e.Statement,
		ClientContextID: e.ClientContextID,
		Errors:          e.Errors,
		Endpoint:        e.Endpoint,
		RetryReasons:    e.RetryReasons,
		RetryAttempts:   e.RetryAttempts,
		HTTPStatusCode:  e.HTTPStatusCode,
	})
}

// Error returns the string representation of this error.
func (e AnalyticsError) Error() string {
	errBytes, serErr := json.Marshal(struct {
		InnerError      error                `json:"-"`
		Statement       string               `json:"statement,omitempty"`
		ClientContextID string               `json:"client_context_id,omitempty"`
		Errors          []AnalyticsErrorDesc `json:"errors,omitempty"`
		Endpoint        string               `json:"endpoint,omitempty"`
		RetryReasons    []RetryReason        `json:"retry_reasons,omitempty"`
		RetryAttempts   uint32               `json:"retry_attempts,omitempty"`
		ErrorText       string               `json:"error_text,omitempty"`
		HTTPStatusCode  int                  `json:"http_status_code,omitempty"`
	}{
		InnerError:      e.InnerError,
		Statement:       e.Statement,
		ClientContextID: e.ClientContextID,
		Errors:          e.Errors,
		Endpoint:        e.Endpoint,
		RetryReasons:    e.RetryReasons,
		RetryAttempts:   e.RetryAttempts,
		ErrorText:       e.ErrorText,
		HTTPStatusCode:  e.HTTPStatusCode,
	})
	if serErr != nil {
		logErrorf("failed to serialize error to json: %s", serErr.Error())
	}

	return e.InnerError.Error() + " | " + string(errBytes)
}

// Unwrap returns the underlying cause for this error.
func (e AnalyticsError) Unwrap() error {
	return e.InnerError
}
