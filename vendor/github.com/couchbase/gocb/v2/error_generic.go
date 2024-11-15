package gocb

import (
	"encoding/json"
)

// GenericError wraps errors that come from the SDK, and can be returned from any service.
// Errors returned when protostellar is used are of this type.
//
// # UNCOMMITTED
//
// This API is UNCOMMITTED and may change in the future.
type GenericError struct {
	InnerError error                  `json:"-"`
	Context    map[string]interface{} `json:"context,omitempty"`
}

// MarshalJSON implements the Marshaler interface.
func (e GenericError) MarshalJSON() ([]byte, error) {
	var innerError string
	if e.InnerError != nil {
		innerError = e.InnerError.Error()
	}
	return json.Marshal(struct {
		InnerError string                 `json:"msg,omitempty"`
		Context    map[string]interface{} `json:"context,omitempty"`
	}{
		InnerError: innerError,
		Context:    e.Context,
	})
}

// Error returns the string representation of a kv error.
func (e GenericError) Error() string {
	errBytes, serErr := json.Marshal(struct {
		InnerError error                  `json:"-"`
		Context    map[string]interface{} `json:"context,omitempty"`
	}{
		InnerError: e.InnerError,
		Context:    e.Context,
	})
	if serErr != nil {
		logErrorf("failed to serialize error to json: %s", serErr.Error())
	}

	return e.InnerError.Error() + " | " + string(errBytes)
}

// Unwrap returns the underlying reason for the error
func (e GenericError) Unwrap() error {
	return e.InnerError
}

func makeGenericError(baseErr error, context map[string]interface{}) *GenericError {
	return &GenericError{
		InnerError: baseErr,
		Context:    context,
	}
}
