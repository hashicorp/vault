package gocb

import gocbcore "github.com/couchbase/gocbcore/v9"

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
}

// Error returns the string representation of this error.
func (e AnalyticsError) Error() string {
	return e.InnerError.Error() + " | " + serializeWrappedError(e)
}

// Unwrap returns the underlying cause for this error.
func (e AnalyticsError) Unwrap() error {
	return e.InnerError
}
