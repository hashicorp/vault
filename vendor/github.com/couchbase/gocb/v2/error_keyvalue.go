package gocb

import (
	"encoding/json"
	"github.com/couchbase/gocbcore/v10/memd"
)

// KeyValueError wraps key-value errors that occur within the SDK.
// UNCOMMITTED: This API may change in the future.
type KeyValueError struct {
	InnerError         error           `json:"-"`
	StatusCode         memd.StatusCode `json:"status_code,omitempty"`
	DocumentID         string          `json:"document_id,omitempty"`
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

// MarshalJSON implements the Marshaler interface.
func (e KeyValueError) MarshalJSON() ([]byte, error) {
	var innerError string
	if e.InnerError != nil {
		innerError = e.InnerError.Error()
	}
	return json.Marshal(struct {
		InnerError         string          `json:"msg,omitempty"`
		StatusCode         memd.StatusCode `json:"status_code,omitempty"`
		DocumentID         string          `json:"document_id,omitempty"`
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
	}{
		InnerError:         innerError,
		StatusCode:         e.StatusCode,
		DocumentID:         e.DocumentID,
		BucketName:         e.BucketName,
		ScopeName:          e.ScopeName,
		CollectionName:     e.CollectionName,
		CollectionID:       e.CollectionID,
		ErrorName:          e.ErrorName,
		ErrorDescription:   e.ErrorDescription,
		Opaque:             e.Opaque,
		Context:            e.Context,
		Ref:                e.Ref,
		RetryReasons:       e.RetryReasons,
		RetryAttempts:      e.RetryAttempts,
		LastDispatchedTo:   e.LastDispatchedTo,
		LastDispatchedFrom: e.LastDispatchedFrom,
		LastConnectionID:   e.LastConnectionID,
	})
}

// Error returns the string representation of a kv error.
func (e KeyValueError) Error() string {
	errBytes, serErr := json.Marshal(struct {
		InnerError         error           `json:"-"`
		StatusCode         memd.StatusCode `json:"status_code,omitempty"`
		DocumentID         string          `json:"document_id,omitempty"`
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
	}{
		InnerError:         e.InnerError,
		StatusCode:         e.StatusCode,
		DocumentID:         e.DocumentID,
		BucketName:         e.BucketName,
		ScopeName:          e.ScopeName,
		CollectionName:     e.CollectionName,
		CollectionID:       e.CollectionID,
		ErrorName:          e.ErrorName,
		ErrorDescription:   e.ErrorDescription,
		Opaque:             e.Opaque,
		Context:            e.Context,
		Ref:                e.Ref,
		RetryReasons:       e.RetryReasons,
		RetryAttempts:      e.RetryAttempts,
		LastDispatchedTo:   e.LastDispatchedTo,
		LastDispatchedFrom: e.LastDispatchedFrom,
		LastConnectionID:   e.LastConnectionID,
	})
	if serErr != nil {
		logErrorf("failed to serialize error to json: %s", serErr.Error())
	}

	return e.InnerError.Error() + " | " + string(errBytes)
}

// Unwrap returns the underlying reason for the error
func (e KeyValueError) Unwrap() error {
	return e.InnerError
}
