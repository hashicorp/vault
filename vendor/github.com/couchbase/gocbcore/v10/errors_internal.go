package gocbcore

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/couchbase/gocbcore/v10/memd"
)

type wrappedError struct {
	Message    string
	InnerError error
}

func (e wrappedError) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.InnerError.Error())
}

func (e wrappedError) Unwrap() error {
	return e.InnerError
}

func wrapError(err error, message string) error {
	return wrappedError{
		Message:    message,
		InnerError: err,
	}
}

// SubDocumentError provides additional contextual information to
// sub-document specific errors.  InnerError is always a KeyValueError.
type SubDocumentError struct {
	InnerError error
	Index      int
}

// Error returns the string representation of this error.
func (err SubDocumentError) Error() string {
	return fmt.Sprintf("sub-document error at index %d: %s",
		err.Index,
		err.InnerError.Error())
}

// Unwrap returns the underlying error for the operation failing.
func (err SubDocumentError) Unwrap() error {
	return err.InnerError
}

func serializeError(err error) string {
	errBytes, serErr := json.Marshal(err)
	if serErr != nil {
		logErrorf("failed to serialize error to json: %s", serErr.Error())
	}
	return string(errBytes)
}

// KeyValueError wraps key-value errors that occur within the SDK.
type KeyValueError struct {
	InnerError         error
	StatusCode         memd.StatusCode
	DocumentKey        string
	BucketName         string
	ScopeName          string
	CollectionName     string
	CollectionID       uint32
	ErrorName          string
	ErrorDescription   string
	Opaque             uint32
	Context            string
	Ref                string
	RetryReasons       []RetryReason
	RetryAttempts      uint32
	LastDispatchedTo   string
	LastDispatchedFrom string
	LastConnectionID   string

	// Internal: This should never be used and is not supported.
	Internal struct {
		ResourceUnits *ResourceUnitResult
	}
}

// MarshalJSON implements the Marshaler interface.
func (e KeyValueError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		InnerError         string          `json:"msg,omitempty"`
		StatusCode         memd.StatusCode `json:"status_code,omitempty"`
		DocumentKey        string          `json:"document_key,omitempty"`
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
		InnerError:         e.InnerError.Error(),
		StatusCode:         e.StatusCode,
		DocumentKey:        e.DocumentKey,
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

// Error returns the string representation of this error.
func (e KeyValueError) Error() string {
	errBytes, serErr := json.Marshal(struct {
		InnerError         error           `json:"-"`
		StatusCode         memd.StatusCode `json:"status_code,omitempty"`
		DocumentKey        string          `json:"document_key,omitempty"`
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
		DocumentKey:        e.DocumentKey,
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

// ViewQueryErrorDesc represents specific view error data.
type ViewQueryErrorDesc struct {
	SourceNode string
	Message    string
}

// ViewError represents an error returned from a view query.
type ViewError struct {
	InnerError         error
	DesignDocumentName string
	ViewName           string
	Errors             []ViewQueryErrorDesc
	Endpoint           string
	RetryReasons       []RetryReason
	RetryAttempts      uint32
	// Uncommitted: This API may change in the future.
	ErrorText string
	// Uncommitted: This API may change in the future.
	HTTPResponseCode int
}

// MarshalJSON implements the Marshaler interface.
func (e ViewError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		InnerError         string               `json:"msg,omitempty"`
		DesignDocumentName string               `json:"design_document_name,omitempty"`
		ViewName           string               `json:"view_name,omitempty"`
		Errors             []ViewQueryErrorDesc `json:"errors,omitempty"`
		Endpoint           string               `json:"endpoint,omitempty"`
		RetryReasons       []RetryReason        `json:"retry_reasons,omitempty"`
		RetryAttempts      uint32               `json:"retry_attempts,omitempty"`
		HTTPResponseCode   int                  `json:"status_code,omitempty"`
	}{
		InnerError:         e.InnerError.Error(),
		DesignDocumentName: e.DesignDocumentName,
		ViewName:           e.ViewName,
		Errors:             e.Errors,
		Endpoint:           e.Endpoint,
		RetryReasons:       e.RetryReasons,
		RetryAttempts:      e.RetryAttempts,
		HTTPResponseCode:   e.HTTPResponseCode,
	})
}

// Error returns the string representation of this error.
func (e ViewError) Error() string {
	errBytes, serErr := json.Marshal(struct {
		InnerError         error                `json:"-"`
		DesignDocumentName string               `json:"design_document_name,omitempty"`
		ViewName           string               `json:"view_name,omitempty"`
		Errors             []ViewQueryErrorDesc `json:"errors,omitempty"`
		Endpoint           string               `json:"endpoint,omitempty"`
		RetryReasons       []RetryReason        `json:"retry_reasons,omitempty"`
		RetryAttempts      uint32               `json:"retry_attempts,omitempty"`
		ErrorText          string               `json:"error_text,omitempty"`
		HTTPResponseCode   int                  `json:"status_code,omitempty"`
	}{
		InnerError:         e.InnerError,
		DesignDocumentName: e.DesignDocumentName,
		ViewName:           e.ViewName,
		Errors:             e.Errors,
		Endpoint:           e.Endpoint,
		RetryReasons:       e.RetryReasons,
		RetryAttempts:      e.RetryAttempts,
		ErrorText:          e.ErrorText,
		HTTPResponseCode:   e.HTTPResponseCode,
	})
	if serErr != nil {
		logErrorf("failed to serialize error to json: %s", serErr.Error())
	}

	return e.InnerError.Error() + " | " + string(errBytes)
}

// Unwrap returns the underlying reason for the error
func (e ViewError) Unwrap() error {
	return e.InnerError
}

// N1QLErrorDesc represents specific n1ql error data.
type N1QLErrorDesc struct {
	Code    uint32
	Message string
	Retry   bool
	Reason  map[string]interface{}
}

// MarshalJSON implements the Marshaler interface.
func (e N1QLErrorDesc) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Code    uint32                 `json:"code"`
		Message string                 `json:"message"`
		Retry   bool                   `json:"retry,omitempty"`
		Reason  map[string]interface{} `json:"reason,omitempty"`
	}{
		Code:    e.Code,
		Message: e.Message,
		Retry:   e.Retry,
		Reason:  e.Reason,
	})
}

// N1QLError represents an error returned from a n1ql query.
type N1QLError struct {
	InnerError      error
	Statement       string
	ClientContextID string
	Errors          []N1QLErrorDesc
	Endpoint        string
	RetryReasons    []RetryReason
	RetryAttempts   uint32
	// Uncommitted: This API may change in the future.
	ErrorText string
	// Uncommitted: This API may change in the future.
	HTTPResponseCode int
}

// MarshalJSON implements the Marshaler interface.
func (e N1QLError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		InnerError       string          `json:"msg,omitempty"`
		Statement        string          `json:"statement,omitempty"`
		ClientContextID  string          `json:"client_context_id,omitempty"`
		Errors           []N1QLErrorDesc `json:"errors,omitempty"`
		Endpoint         string          `json:"endpoint,omitempty"`
		RetryReasons     []RetryReason   `json:"retry_reasons,omitempty"`
		RetryAttempts    uint32          `json:"retry_attempts,omitempty"`
		HTTPResponseCode int             `json:"status_code,omitempty"`
	}{
		InnerError:       e.InnerError.Error(),
		Statement:        e.Statement,
		ClientContextID:  e.ClientContextID,
		Errors:           e.Errors,
		Endpoint:         e.Endpoint,
		RetryReasons:     e.RetryReasons,
		RetryAttempts:    e.RetryAttempts,
		HTTPResponseCode: e.HTTPResponseCode,
	})
}

// Error returns the string representation of this error.
func (e N1QLError) Error() string {
	errBytes, serErr := json.Marshal(struct {
		InnerError       error           `json:"-"`
		Statement        string          `json:"statement,omitempty"`
		ClientContextID  string          `json:"client_context_id,omitempty"`
		Errors           []N1QLErrorDesc `json:"errors,omitempty"`
		Endpoint         string          `json:"endpoint,omitempty"`
		RetryReasons     []RetryReason   `json:"retry_reasons,omitempty"`
		RetryAttempts    uint32          `json:"retry_attempts,omitempty"`
		ErrorText        string          `json:"error_text,omitempty"`
		HTTPResponseCode int             `json:"status_code,omitempty"`
	}{
		InnerError:       e.InnerError,
		Statement:        e.Statement,
		ClientContextID:  e.ClientContextID,
		Errors:           e.Errors,
		Endpoint:         e.Endpoint,
		RetryReasons:     e.RetryReasons,
		RetryAttempts:    e.RetryAttempts,
		ErrorText:        e.ErrorText,
		HTTPResponseCode: e.HTTPResponseCode,
	})
	if serErr != nil {
		logErrorf("failed to serialize error to json: %s", serErr.Error())
	}

	return e.InnerError.Error() + " | " + string(errBytes)
}

// Unwrap returns the underlying reason for the error
func (e N1QLError) Unwrap() error {
	return e.InnerError
}

// AnalyticsErrorDesc represents specific analytics error data.
type AnalyticsErrorDesc struct {
	Code    uint32
	Message string
}

// AnalyticsError represents an error returned from an analytics query.
type AnalyticsError struct {
	InnerError      error
	Statement       string
	ClientContextID string
	Errors          []AnalyticsErrorDesc
	Endpoint        string
	RetryReasons    []RetryReason
	RetryAttempts   uint32
	// Uncommitted: This API may change in the future.
	ErrorText string
	// Uncommitted: This API may change in the future.
	HTTPResponseCode int
}

// MarshalJSON implements the Marshaler interface.
func (e AnalyticsError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		InnerError       string               `json:"msg,omitempty"`
		Statement        string               `json:"statement,omitempty"`
		ClientContextID  string               `json:"client_context_id,omitempty"`
		Errors           []AnalyticsErrorDesc `json:"errors,omitempty"`
		Endpoint         string               `json:"endpoint,omitempty"`
		RetryReasons     []RetryReason        `json:"retry_reasons,omitempty"`
		RetryAttempts    uint32               `json:"retry_attempts,omitempty"`
		HTTPResponseCode int                  `json:"status_code,omitempty"`
	}{
		InnerError:       e.InnerError.Error(),
		Statement:        e.Statement,
		ClientContextID:  e.ClientContextID,
		Errors:           e.Errors,
		Endpoint:         e.Endpoint,
		RetryReasons:     e.RetryReasons,
		RetryAttempts:    e.RetryAttempts,
		HTTPResponseCode: e.HTTPResponseCode,
	})
}

// Error returns the string representation of this error.
func (e AnalyticsError) Error() string {
	errBytes, serErr := json.Marshal(struct {
		InnerError       error                `json:"-"`
		Statement        string               `json:"statement,omitempty"`
		ClientContextID  string               `json:"client_context_id,omitempty"`
		Errors           []AnalyticsErrorDesc `json:"errors,omitempty"`
		Endpoint         string               `json:"endpoint,omitempty"`
		RetryReasons     []RetryReason        `json:"retry_reasons,omitempty"`
		RetryAttempts    uint32               `json:"retry_attempts,omitempty"`
		ErrorText        string               `json:"error_text,omitempty"`
		HTTPResponseCode int                  `json:"status_code,omitempty"`
	}{
		InnerError:       e.InnerError,
		Statement:        e.Statement,
		ClientContextID:  e.ClientContextID,
		Errors:           e.Errors,
		Endpoint:         e.Endpoint,
		RetryReasons:     e.RetryReasons,
		RetryAttempts:    e.RetryAttempts,
		ErrorText:        e.ErrorText,
		HTTPResponseCode: e.HTTPResponseCode,
	})
	if serErr != nil {
		logErrorf("failed to serialize error to json: %s", serErr.Error())
	}

	return e.InnerError.Error() + " | " + string(errBytes)
}

// Unwrap returns the underlying reason for the error
func (e AnalyticsError) Unwrap() error {
	return e.InnerError
}

// SearchError represents an error returned from a search query.
type SearchError struct {
	InnerError       error
	IndexName        string
	Query            interface{}
	ErrorText        string
	HTTPResponseCode int
	Endpoint         string
	RetryReasons     []RetryReason
	RetryAttempts    uint32
}

// MarshalJSON implements the Marshaler interface.
func (e SearchError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		InnerError       string        `json:"msg,omitempty"`
		IndexName        string        `json:"index_name,omitempty"`
		Query            interface{}   `json:"query,omitempty"`
		ErrorText        string        `json:"error_text"`
		HTTPResponseCode int           `json:"status_code,omitempty"`
		Endpoint         string        `json:"endpoint,omitempty"`
		RetryReasons     []RetryReason `json:"retry_reasons,omitempty"`
		RetryAttempts    uint32        `json:"retry_attempts,omitempty"`
	}{
		InnerError:       e.InnerError.Error(),
		IndexName:        e.IndexName,
		Query:            e.Query,
		ErrorText:        e.ErrorText,
		HTTPResponseCode: e.HTTPResponseCode,
		Endpoint:         e.Endpoint,
		RetryReasons:     e.RetryReasons,
		RetryAttempts:    e.RetryAttempts,
	})
}

// Error returns the string representation of this error.
func (e SearchError) Error() string {
	errBytes, serErr := json.Marshal(struct {
		InnerError       error         `json:"-"`
		IndexName        string        `json:"index_name,omitempty"`
		Query            interface{}   `json:"query,omitempty"`
		ErrorText        string        `json:"error_text"`
		HTTPResponseCode int           `json:"status_code,omitempty"`
		Endpoint         string        `json:"endpoint,omitempty"`
		RetryReasons     []RetryReason `json:"retry_reasons,omitempty"`
		RetryAttempts    uint32        `json:"retry_attempts,omitempty"`
	}{
		InnerError:       e.InnerError,
		IndexName:        e.IndexName,
		Query:            e.Query,
		ErrorText:        e.ErrorText,
		HTTPResponseCode: e.HTTPResponseCode,
		Endpoint:         e.Endpoint,
		RetryReasons:     e.RetryReasons,
		RetryAttempts:    e.RetryAttempts,
	})
	if serErr != nil {
		logErrorf("failed to serialize error to json: %s", serErr.Error())
	}

	return e.InnerError.Error() + " | " + string(errBytes)
}

// Unwrap returns the underlying reason for the error
func (e SearchError) Unwrap() error {
	return e.InnerError
}

// HTTPError represents an error returned from an HTTP request.
type HTTPError struct {
	InnerError    error
	UniqueID      string
	Endpoint      string
	RetryReasons  []RetryReason
	RetryAttempts uint32
}

// MarshalJSON implements the Marshaler interface.
func (e HTTPError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		InnerError    string        `json:"msg,omitempty"`
		UniqueID      string        `json:"unique_id,omitempty"`
		Endpoint      string        `json:"endpoint,omitempty"`
		RetryReasons  []RetryReason `json:"retry_reasons,omitempty"`
		RetryAttempts uint32        `json:"retry_attempts,omitempty"`
	}{
		InnerError:    e.InnerError.Error(),
		UniqueID:      e.UniqueID,
		Endpoint:      e.Endpoint,
		RetryReasons:  e.RetryReasons,
		RetryAttempts: e.RetryAttempts,
	})
}

// Error returns the string representation of this error.
func (e HTTPError) Error() string {
	errBytes, serErr := json.Marshal(struct {
		InnerError    error         `json:"-"`
		UniqueID      string        `json:"unique_id,omitempty"`
		Endpoint      string        `json:"endpoint,omitempty"`
		RetryReasons  []RetryReason `json:"retry_reasons,omitempty"`
		RetryAttempts uint32        `json:"retry_attempts,omitempty"`
	}{
		InnerError:    e.InnerError,
		UniqueID:      e.UniqueID,
		Endpoint:      e.Endpoint,
		RetryReasons:  e.RetryReasons,
		RetryAttempts: e.RetryAttempts,
	})
	if serErr != nil {
		logErrorf("failed to serialize error to json: %s", serErr.Error())
	}

	return e.InnerError.Error() + " | " + string(errBytes)
}

// Unwrap returns the underlying reason for the error
func (e HTTPError) Unwrap() error {
	return e.InnerError
}

// TimeoutError wraps timeout errors that occur within the SDK.
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

	// Internal: This should never be used and is not supported.
	Internal struct {
		ResourceUnits *ResourceUnitResult
	}
}

func makeTimeoutError(start time.Time, op string, innerErr error, req *memdQRequest) *TimeoutError {
	connInfo := req.ConnectionInfo()
	count, reasons := req.Retries()
	err := &TimeoutError{
		InnerError:         innerErr,
		OperationID:        op,
		Opaque:             req.Identifier(),
		TimeObserved:       time.Since(start),
		RetryReasons:       reasons,
		RetryAttempts:      count,
		LastDispatchedTo:   connInfo.lastDispatchedTo,
		LastDispatchedFrom: connInfo.lastDispatchedFrom,
		LastConnectionID:   connInfo.lastConnectionID,
	}
	err.Internal.ResourceUnits = req.ResourceUnits()

	return err
}

type timeoutError struct {
	InnerError         error         `json:"-,omitempty"`
	OperationID        string        `json:"s,omitempty"`
	Opaque             string        `json:"i,omitempty"`
	TimeObserved       uint64        `json:"t,omitempty"`
	RetryReasons       []RetryReason `json:"rr,omitempty"`
	RetryAttempts      uint32        `json:"ra,omitempty"`
	LastDispatchedTo   string        `json:"r,omitempty"`
	LastDispatchedFrom string        `json:"l,omitempty"`
	LastConnectionID   string        `json:"c,omitempty"`
}

// MarshalJSON implements the Marshaler interface.
func (err *TimeoutError) MarshalJSON() ([]byte, error) {
	toMarshal := timeoutError{
		InnerError:         err.InnerError,
		OperationID:        err.OperationID,
		Opaque:             err.Opaque,
		TimeObserved:       uint64(err.TimeObserved / time.Microsecond),
		RetryReasons:       err.RetryReasons,
		RetryAttempts:      err.RetryAttempts,
		LastDispatchedTo:   err.LastDispatchedTo,
		LastDispatchedFrom: err.LastDispatchedFrom,
		LastConnectionID:   err.LastConnectionID,
	}

	return json.Marshal(toMarshal)
}

// UnmarshalJSON implements the Unmarshaler interface.
func (err *TimeoutError) UnmarshalJSON(data []byte) error {
	var tErr timeoutError
	if jErr := json.Unmarshal(data, &tErr); jErr != nil {
		return jErr
	}

	duration := time.Duration(tErr.TimeObserved) * time.Microsecond

	err.InnerError = tErr.InnerError
	err.OperationID = tErr.OperationID
	err.Opaque = tErr.Opaque
	err.TimeObserved = duration
	err.RetryReasons = tErr.RetryReasons
	err.RetryAttempts = tErr.RetryAttempts
	err.LastDispatchedTo = tErr.LastDispatchedTo
	err.LastDispatchedFrom = tErr.LastDispatchedFrom
	err.LastConnectionID = tErr.LastConnectionID

	return nil
}

func (err TimeoutError) Error() string {
	return err.InnerError.Error() + " | " + serializeError(err)
}

// Unwrap returns the underlying reason for the error
func (err TimeoutError) Unwrap() error {
	return err.InnerError
}

type DCPRollbackError struct {
	InnerError error
	SeqNo      SeqNo
}

// MarshalJSON implements the Marshaler interface.
func (e DCPRollbackError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		InnerError string `json:"msg,omitempty"`
		SeqNo      uint64 `json:"seq_no,omitempty"`
	}{
		InnerError: e.InnerError.Error(),
		SeqNo:      uint64(e.SeqNo),
	})
}

// Error returns the string representation of this error.
func (e DCPRollbackError) Error() string {
	errBytes, serErr := json.Marshal(struct {
		InnerError error  `json:"-"`
		SeqNo      uint64 `json:"seq_no,omitempty"`
	}{
		InnerError: e.InnerError,
		SeqNo:      uint64(e.SeqNo),
	})
	if serErr != nil {
		logErrorf("failed to serialize error to json: %s", serErr.Error())
	}

	return e.InnerError.Error() + " | " + string(errBytes)
}

// Unwrap returns the underlying reason for the error
func (err DCPRollbackError) Unwrap() error {
	return err.InnerError
}

// ncError is a wrapper error that provides no additional context to one of the
// publicly exposed error types.  This is to force people to correctly use the
// error handling behaviours to check the error, rather than direct compares.
type ncError struct {
	InnerError error
}

func (err ncError) Error() string {
	return err.InnerError.Error()
}

func (err ncError) Unwrap() error {
	return err.InnerError
}

func isErrorStatus(err error, code memd.StatusCode) bool {
	var kvErr *KeyValueError
	if errors.As(err, &kvErr) {
		return kvErr.StatusCode == code
	}
	return false
}

var (
	errNoCCCPHosts = errors.New("no cccp hosts available")
)

// This list contains protected versions of all the errors we throw
// to ensure no users inadvertently rely on direct comparisons.
// nolint: deadcode,varcheck,unused
var (
	errTimeout                  = ncError{ErrTimeout}
	errRequestCanceled          = ncError{ErrRequestCanceled}
	errInvalidArgument          = ncError{ErrInvalidArgument}
	errServiceNotAvailable      = ncError{ErrServiceNotAvailable}
	errInternalServerFailure    = ncError{ErrInternalServerFailure}
	errAuthenticationFailure    = ncError{ErrAuthenticationFailure}
	errTemporaryFailure         = ncError{ErrTemporaryFailure}
	errBusy                     = ncError{ErrBusy}
	errParsingFailure           = ncError{ErrParsingFailure}
	errCasMismatch              = ncError{ErrCasMismatch}
	errBucketNotFound           = ncError{ErrBucketNotFound}
	errCollectionNotFound       = ncError{ErrCollectionNotFound}
	errEncodingFailure          = ncError{ErrEncodingFailure}
	errDecodingFailure          = ncError{ErrDecodingFailure}
	errUnsupportedOperation     = ncError{ErrUnsupportedOperation}
	errAmbiguousTimeout         = ncError{ErrAmbiguousTimeout}
	errUnambiguousTimeout       = ncError{ErrUnambiguousTimeout}
	errFeatureNotAvailable      = ncError{ErrFeatureNotAvailable}
	errScopeNotFound            = ncError{ErrScopeNotFound}
	errIndexNotFound            = ncError{ErrIndexNotFound}
	errIndexExists              = ncError{ErrIndexExists}
	errGCCCPInUse               = ncError{ErrGCCCPInUse}
	errNotMyVBucket             = ncError{ErrNotMyVBucket}
	errDMLFailure               = ncError{ErrDMLFailure}
	errMemdClientClosed         = ncError{ErrMemdClientClosed}
	errRequestAlreadyDispatched = ncError{ErrRequestAlreadyDispatched}

	errDocumentNotFound                  = ncError{ErrDocumentNotFound}
	errDocumentUnretrievable             = ncError{ErrDocumentUnretrievable}
	errDocumentLocked                    = ncError{ErrDocumentLocked}
	errDocumentNotLocked                 = ncError{ErrDocumentNotLocked}
	errValueTooLarge                     = ncError{ErrValueTooLarge}
	errDocumentExists                    = ncError{ErrDocumentExists}
	errNotStored                         = ncError{ErrNotStored}
	errValueNotJSON                      = ncError{ErrValueNotJSON}
	errDurabilityLevelNotAvailable       = ncError{ErrDurabilityLevelNotAvailable}
	errDurabilityImpossible              = ncError{ErrDurabilityImpossible}
	errDurabilityAmbiguous               = ncError{ErrDurabilityAmbiguous}
	errDurableWriteInProgress            = ncError{ErrDurableWriteInProgress}
	errDurableWriteReCommitInProgress    = ncError{ErrDurableWriteReCommitInProgress}
	errMutationLost                      = ncError{ErrMutationLost}
	errPathNotFound                      = ncError{ErrPathNotFound}
	errPathMismatch                      = ncError{ErrPathMismatch}
	errPathInvalid                       = ncError{ErrPathInvalid}
	errPathTooBig                        = ncError{ErrPathTooBig}
	errPathTooDeep                       = ncError{ErrPathTooDeep}
	errValueTooDeep                      = ncError{ErrValueTooDeep}
	errValueInvalid                      = ncError{ErrValueInvalid}
	errDocumentNotJSON                   = ncError{ErrDocumentNotJSON}
	errNumberTooBig                      = ncError{ErrNumberTooBig}
	errDeltaInvalid                      = ncError{ErrDeltaInvalid}
	errPathExists                        = ncError{ErrPathExists}
	errXattrUnknownMacro                 = ncError{ErrXattrUnknownMacro}
	errXattrInvalidFlagCombo             = ncError{ErrXattrInvalidFlagCombo}
	errXattrInvalidKeyCombo              = ncError{ErrXattrInvalidKeyCombo}
	errXattrUnknownVirtualAttribute      = ncError{ErrXattrUnknownVirtualAttribute}
	errXattrCannotModifyVirtualAttribute = ncError{ErrXattrCannotModifyVirtualAttribute}
	errXattrInvalidOrder                 = ncError{ErrXattrInvalidOrder}

	errPlanningFailure          = ncError{ErrPlanningFailure}
	errIndexFailure             = ncError{ErrIndexFailure}
	errPreparedStatementFailure = ncError{ErrPreparedStatementFailure}

	errCompilationFailure = ncError{ErrCompilationFailure}
	errJobQueueFull       = ncError{ErrJobQueueFull}
	errDatasetNotFound    = ncError{ErrDatasetNotFound}
	errDataverseNotFound  = ncError{ErrDataverseNotFound}
	errDatasetExists      = ncError{ErrDatasetExists}
	errDataverseExists    = ncError{ErrDataverseExists}
	errLinkNotFound       = ncError{ErrLinkNotFound}

	errViewNotFound           = ncError{ErrViewNotFound}
	errDesignDocumentNotFound = ncError{ErrDesignDocumentNotFound}

	errNoSupportedMechanisms  = ncError{ErrNoSupportedMechanisms}
	errBadHosts               = ncError{ErrBadHosts}
	errProtocol               = ncError{ErrProtocol}
	errNoReplicas             = ncError{ErrNoReplicas}
	errCliInternalError       = ncError{ErrCliInternalError}
	errInvalidCredentials     = ncError{ErrInvalidCredentials}
	errInvalidServer          = ncError{ErrInvalidServer}
	errInvalidVBucket         = ncError{ErrInvalidVBucket}
	errInvalidReplica         = ncError{ErrInvalidReplica}
	errInvalidService         = ncError{ErrInvalidService}
	errInvalidCertificate     = ncError{ErrInvalidCertificate}
	errCollectionsUnsupported = ncError{ErrCollectionsUnsupported}
	errBucketAlreadySelected  = ncError{ErrBucketAlreadySelected}
	errShutdown               = ncError{ErrShutdown}
	errOverload               = ncError{ErrOverload}
	errStreamIDNotEnabled     = ncError{ErrStreamIDNotEnabled}
	errDCPStreamIDInvalid     = ncError{ErrDCPStreamIDInvalid}
	errForcedReconnect        = ncError{ErrForcedReconnect}

	errRateLimitedFailure  = ncError{ErrRateLimitedFailure}
	errQuotaLimitedFailure = ncError{ErrQuotaLimitedFailure}

	errRangeScanCancelled      = ncError{ErrRangeScanCancelled}
	errRangeScanMore           = ncError{ErrRangeScanMore}
	errRangeScanComplete       = ncError{ErrRangeScanComplete}
	errRangeScanVbUUIDNotEqual = ncError{ErrRangeScanVbUUIDNotEqual}

	errConnectionIDInvalid = ncError{ErrConnectionIDInvalid}

	errCircuitBreakerOpen = ncError{ErrCircuitBreakerOpen}
)
