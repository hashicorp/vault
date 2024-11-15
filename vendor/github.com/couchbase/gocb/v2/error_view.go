package gocb

import (
	"encoding/json"
	gocbcore "github.com/couchbase/gocbcore/v10"
)

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
	ErrorText          string          `json:"error_text,omitempty"`
	HTTPStatusCode     int             `json:"http_status_code,omitempty"`
}

// MarshalJSON implements the Marshaler interface.
func (e ViewError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		InnerError         string          `json:"msg,omitempty"`
		DesignDocumentName string          `json:"design_document_name,omitempty"`
		ViewName           string          `json:"view_name,omitempty"`
		Errors             []ViewErrorDesc `json:"errors,omitempty"`
		Endpoint           string          `json:"endpoint,omitempty"`
		RetryReasons       []RetryReason   `json:"retry_reasons,omitempty"`
		RetryAttempts      uint32          `json:"retry_attempts,omitempty"`
		HTTPStatusCode     int             `json:"http_status_code,omitempty"`
	}{
		InnerError:         e.InnerError.Error(),
		DesignDocumentName: e.DesignDocumentName,
		ViewName:           e.ViewName,
		Errors:             e.Errors,
		Endpoint:           e.Endpoint,
		RetryReasons:       e.RetryReasons,
		RetryAttempts:      e.RetryAttempts,
		HTTPStatusCode:     e.HTTPStatusCode,
	})
}

// Error returns the string representation of this error.
func (e ViewError) Error() string {
	errBytes, serErr := json.Marshal(struct {
		InnerError         error           `json:"-"`
		DesignDocumentName string          `json:"design_document_name,omitempty"`
		ViewName           string          `json:"view_name,omitempty"`
		Errors             []ViewErrorDesc `json:"errors,omitempty"`
		Endpoint           string          `json:"endpoint,omitempty"`
		RetryReasons       []RetryReason   `json:"retry_reasons,omitempty"`
		RetryAttempts      uint32          `json:"retry_attempts,omitempty"`
		ErrorText          string          `json:"error_text,omitempty"`
		HTTPStatusCode     int             `json:"http_status_code,omitempty"`
	}{
		InnerError:         e.InnerError,
		DesignDocumentName: e.DesignDocumentName,
		ViewName:           e.ViewName,
		Errors:             e.Errors,
		Endpoint:           e.Endpoint,
		RetryReasons:       e.RetryReasons,
		RetryAttempts:      e.RetryAttempts,
		ErrorText:          e.ErrorText,
		HTTPStatusCode:     e.HTTPStatusCode,
	})
	if serErr != nil {
		logErrorf("failed to serialize error to json: %s", serErr.Error())
	}

	return e.InnerError.Error() + " | " + string(errBytes)
}

// Unwrap returns the underlying cause for this error.
func (e ViewError) Unwrap() error {
	return e.InnerError
}
