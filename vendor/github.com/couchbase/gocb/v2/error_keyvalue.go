package gocb

import "github.com/couchbase/gocbcore/v9/memd"

// KeyValueError wraps key-value errors that occur within the SDK.
// UNCOMMITTED: This API may change in the future.
type KeyValueError struct {
	InnerError         error           `json:"-"`
	StatusCode         memd.StatusCode `json:"status_code,omitempty"`
	BucketName         string          `json:"bucket,omitempty"`
	ScopeName          string          `json:"scope,omitempty"`
	CollectionName     string          `json:"collection,omitempty"`
	CollectionID       uint32          `json:"collection_id,omitempty"`
	ErrorName          string          `json:"error_name,omitempty"`
	ErrorDescription   string          `json:"error_description,omitempty"`
	Opaque             uint32          `json:"opaque,omitempty"`
	Context            string          `json:"context,omitempty"`
	Ref                string          `json:"ref,omitempty"`
	RetryReasons       []RetryReason   `json:"retry_reasons,omitempty"`
	RetryAttempts      uint32          `json:"retry_attempts,omitempty"`
	LastDispatchedTo   string          `json:"last_dispatched_to,omitempty"`
	LastDispatchedFrom string          `json:"last_dispatched_from,omitempty"`
	LastConnectionID   string          `json:"last_connection_id,omitempty"`
}

// Error returns the string representation of a kv error.
func (e KeyValueError) Error() string {
	return e.InnerError.Error() + " | " + serializeWrappedError(e)
}

// Unwrap returns the underlying reason for the error
func (e KeyValueError) Unwrap() error {
	return e.InnerError
}
