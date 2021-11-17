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

type timeoutError struct {
	InnerError         error    `json:"-"`
	OperationID        string   `json:"s,omitempty"`
	Opaque             string   `json:"i,omitempty"`
	TimeObserved       uint64   `json:"t,omitempty"`
	RetryReasons       []string `json:"rr,omitempty"`
	RetryAttempts      uint32   `json:"ra,omitempty"`
	LastDispatchedTo   string   `json:"r,omitempty"`
	LastDispatchedFrom string   `json:"l,omitempty"`
	LastConnectionID   string   `json:"c,omitempty"`
}

// MarshalJSON implements the Marshaler interface.
func (err *TimeoutError) MarshalJSON() ([]byte, error) {
	var retries []string
	for _, rr := range err.RetryReasons {
		retries = append(retries, rr.Description())
	}

	toMarshal := timeoutError{
		InnerError:         err.InnerError,
		OperationID:        err.OperationID,
		Opaque:             err.Opaque,
		TimeObserved:       uint64(err.TimeObserved / time.Microsecond),
		RetryReasons:       retries,
		RetryAttempts:      err.RetryAttempts,
		LastDispatchedTo:   err.LastDispatchedTo,
		LastDispatchedFrom: err.LastDispatchedFrom,
		LastConnectionID:   err.LastConnectionID,
	}

	return json.Marshal(toMarshal)
}

// UnmarshalJSON implements the Unmarshaler interface.
func (err *TimeoutError) UnmarshalJSON(data []byte) error {
	var tErr *timeoutError
	if err := json.Unmarshal(data, &tErr); err != nil {
		return err
	}

	duration := time.Duration(tErr.TimeObserved) * time.Microsecond

	// Note that we cannot reasonably unmarshal the retry reasons
	err.OperationID = tErr.OperationID
	err.Opaque = tErr.Opaque
	err.TimeObserved = duration
	err.RetryAttempts = tErr.RetryAttempts
	err.LastDispatchedTo = tErr.LastDispatchedTo
	err.LastDispatchedFrom = tErr.LastDispatchedFrom
	err.LastConnectionID = tErr.LastConnectionID

	return nil
}

func (err TimeoutError) Error() string {
	if err.InnerError == nil {
		return serializeWrappedError(err)
	}
	return err.InnerError.Error() + " | " + serializeWrappedError(err)
}

// Unwrap returns the underlying reason for the error
func (err TimeoutError) Unwrap() error {
	return err.InnerError
}
