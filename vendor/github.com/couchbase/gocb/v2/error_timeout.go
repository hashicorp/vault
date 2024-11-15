package gocb

import (
	"encoding/json"
	"time"
)

// TimeoutError wraps timeout errors that occur within the SDK.
// UNCOMMITTED: This API may change in the future.
type TimeoutError struct {
	InnerError         error
	OperationID        string
	Opaque             string
	TimeObserved       time.Duration
	RetryReasons       []RetryReason
	RetryAttempts      uint32
	LastDispatchedTo   string
	LastDispatchedFrom string
	LastConnectionID   string
}

// MarshalJSON implements the Marshaler interface.
func (e TimeoutError) MarshalJSON() ([]byte, error) {
	var innerError string
	if e.InnerError != nil {
		innerError = e.InnerError.Error()
	}
	var retries []string
	for _, rr := range e.RetryReasons {
		retries = append(retries, rr.Description())
	}
	return json.Marshal(struct {
		InnerError         string   `json:"msg,omitempty"`
		OperationID        string   `json:"operation_id,omitempty"`
		Opaque             string   `json:"opaque,omitempty"`
		TimeObserved       uint64   `json:"time_observed,omitempty"`
		RetryReasons       []string `json:"retry_reasons,omitempty"`
		RetryAttempts      uint32   `json:"retry_attempts,omitempty"`
		LastDispatchedTo   string   `json:"last_dispatched_to,omitempty"`
		LastDispatchedFrom string   `json:"last_dispatched_from,omitempty"`
		LastConnectionID   string   `json:"last_connection_id,omitempty"`
	}{
		InnerError:         innerError,
		OperationID:        e.OperationID,
		Opaque:             e.Opaque,
		TimeObserved:       uint64(e.TimeObserved / time.Microsecond),
		RetryReasons:       retries,
		RetryAttempts:      e.RetryAttempts,
		LastDispatchedTo:   e.LastDispatchedTo,
		LastDispatchedFrom: e.LastDispatchedFrom,
		LastConnectionID:   e.LastConnectionID,
	})
}

// Error returns the string representation of this error.
func (e TimeoutError) Error() string {
	errBytes, serErr := json.Marshal(struct {
		InnerError         error         `json:"-"`
		OperationID        string        `json:"operation_id,omitempty"`
		Opaque             string        `json:"opaque,omitempty"`
		TimeObserved       uint64        `json:"time_observed,omitempty"`
		RetryReasons       []RetryReason `json:"retry_reasons,omitempty"`
		RetryAttempts      uint32        `json:"retry_attempts,omitempty"`
		LastDispatchedTo   string        `json:"last_dispatched_to,omitempty"`
		LastDispatchedFrom string        `json:"last_dispatched_from,omitempty"`
		LastConnectionID   string        `json:"last_connection_id,omitempty"`
	}{
		InnerError:         e.InnerError,
		OperationID:        e.OperationID,
		Opaque:             e.Opaque,
		TimeObserved:       uint64(e.TimeObserved / time.Microsecond),
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
func (e TimeoutError) Unwrap() error {
	return e.InnerError
}
