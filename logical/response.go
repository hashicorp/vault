package logical

import (
	"fmt"
	"time"
)

// Response is a struct that stores the response of a request.
// It is used to abstract the details of the higher level request protocol.
type Response struct {
	// IsSecret is used to indicate this is secret material instead of policy or configuration.
	// Non-secrets never have a VaultID or renewable properties.
	IsSecret bool

	// The lease settings if applicable.
	Lease *Lease

	// Response data is an opaque map that must have string keys.
	Data map[string]interface{}
}

// Lease is used to provide more information about the lease
type Lease struct {
	VaultID      string        // VaultID is the unique identifier used for renewal and revocation
	Renewable    bool          // Is the VaultID renewable
	Revokable    bool          // Is the secret revokable. Must support 'Revoke' operation.
	Duration     time.Duration // Current lease duration
	MaxDuration  time.Duration // Maximum lease duration
	MaxIncrement time.Duration // Maximum increment to lease duration
}

// Validate is used to sanity check a lease
func (l *Lease) Validate() error {
	if l.Duration <= 0 {
		return fmt.Errorf("lease duration must be greater than zero")
	}
	if l.MaxDuration <= 0 {
		return fmt.Errorf("maximum lease duration must be greater than zero")
	}
	if l.Duration > l.MaxDuration {
		return fmt.Errorf("lease duration cannot be greater than maximum lease duration")
	}
	if l.MaxIncrement < 0 {
		return fmt.Errorf("maximum lease increment cannot be negative")
	}
	return nil
}

// HelpResponse is used to format a help response
func HelpResponse(text string, seeAlso []string) *Response {
	return &Response{
		IsSecret: false,
		Data: map[string]interface{}{
			"help":     text,
			"see_also": seeAlso,
		},
	}
}

// ErrorResponse is used to format an error response
func ErrorResponse(text string) *Response {
	return &Response{
		IsSecret: false,
		Data: map[string]interface{}{
			"error": text,
		},
	}
}

// ListResponse is used to format a response to a list operation.
func ListResponse(keys []string) *Response {
	return &Response{
		Data: map[string]interface{}{
			"keys": keys,
		},
	}
}
