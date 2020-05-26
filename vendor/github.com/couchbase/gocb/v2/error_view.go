package gocb

import gocbcore "github.com/couchbase/gocbcore/v9"

// ViewErrorDesc represents a specific error returned from the views service.
type ViewErrorDesc struct {
	SourceNode string
	Message    string
}

func translateCoreViewErrorDesc(descs []gocbcore.ViewQueryErrorDesc) []ViewErrorDesc {
	descsOut := make([]ViewErrorDesc, len(descs))
	for descIdx, desc := range descs {
		descsOut[descIdx] = ViewErrorDesc{
			SourceNode: desc.SourceNode,
			Message:    desc.Message,
		}
	}
	return descsOut
}

// ViewError is the error type of all view query errors.
// UNCOMMITTED: This API may change in the future.
type ViewError struct {
	InnerError         error           `json:"-"`
	DesignDocumentName string          `json:"design_document_name,omitempty"`
	ViewName           string          `json:"view_name,omitempty"`
	Errors             []ViewErrorDesc `json:"errors,omitempty"`
	Endpoint           string          `json:"endpoint,omitempty"`
	RetryReasons       []RetryReason   `json:"retry_reasons,omitempty"`
	RetryAttempts      uint32          `json:"retry_attempts,omitempty"`
}

// Error returns the string representation of this error.
func (e ViewError) Error() string {
	return e.InnerError.Error() + " | " + serializeWrappedError(e)
}

// Unwrap returns the underlying cause for this error.
func (e ViewError) Unwrap() error {
	return e.InnerError
}
