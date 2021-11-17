package gocb

import (
	"encoding/json"
	gocbcore "github.com/couchbase/gocbcore/v10"
)

// QueryErrorDesc represents a specific error returned from the query service.
type QueryErrorDesc struct {
	Code    uint32
	Message string
}

func translateCoreQueryErrorDesc(descs []gocbcore.N1QLErrorDesc) []QueryErrorDesc {
	descsOut := make([]QueryErrorDesc, len(descs))
	for descIdx, desc := range descs {
		descsOut[descIdx] = QueryErrorDesc{
			Code:    desc.Code,
			Message: desc.Message,
		}
	}
	return descsOut
}

// QueryError is the error type of all query errors.
// UNCOMMITTED: This API may change in the future.
type QueryError struct {
	InnerError      error            `json:"-"`
	Statement       string           `json:"statement,omitempty"`
	ClientContextID string           `json:"client_context_id,omitempty"`
	Errors          []QueryErrorDesc `json:"errors,omitempty"`
	Endpoint        string           `json:"endpoint,omitempty"`
	RetryReasons    []RetryReason    `json:"retry_reasons,omitempty"`
	RetryAttempts   uint32           `json:"retry_attempts,omitempty"`
}

// MarshalJSON implements the Marshaler interface.

func (e QueryError) MarshalJSON() ([]byte, error) {
	var innerError string
	if e.InnerError != nil {
		innerError = e.InnerError.Error()
	}
	return json.Marshal(struct {
		InnerError      string           `json:"msg,omitempty"`
		Statement       string           `json:"statement,omitempty"`
		ClientContextID string           `json:"client_context_id,omitempty"`
		Errors          []QueryErrorDesc `json:"errors,omitempty"`
		Endpoint        string           `json:"endpoint,omitempty"`
		RetryReasons    []RetryReason    `json:"retry_reasons,omitempty"`
		RetryAttempts   uint32           `json:"retry_attempts,omitempty"`
	}{
		InnerError:      innerError,
		Statement:       e.Statement,
		ClientContextID: e.ClientContextID,
		Errors:          e.Errors,
		Endpoint:        e.Endpoint,
		RetryReasons:    e.RetryReasons,
		RetryAttempts:   e.RetryAttempts,
	})
}

// Error returns the string representation of this error.
func (e QueryError) Error() string {
	errBytes, serErr := json.Marshal(struct {
		InnerError      error            `json:"-"`
		Statement       string           `json:"statement,omitempty"`
		ClientContextID string           `json:"client_context_id,omitempty"`
		Errors          []QueryErrorDesc `json:"errors,omitempty"`
		Endpoint        string           `json:"endpoint,omitempty"`
		RetryReasons    []RetryReason    `json:"retry_reasons,omitempty"`
		RetryAttempts   uint32           `json:"retry_attempts,omitempty"`
	}{
		InnerError:      e.InnerError,
		Statement:       e.Statement,
		ClientContextID: e.ClientContextID,
		Errors:          e.Errors,
		Endpoint:        e.Endpoint,
		RetryReasons:    e.RetryReasons,
		RetryAttempts:   e.RetryAttempts,
	})
	if serErr != nil {
		logErrorf("failed to serialize error to json: %s", serErr.Error())
	}

	return e.InnerError.Error() + " | " + string(errBytes)
}

// Unwrap returns the underlying cause for this error.
func (e QueryError) Unwrap() error {
	return e.InnerError
}
