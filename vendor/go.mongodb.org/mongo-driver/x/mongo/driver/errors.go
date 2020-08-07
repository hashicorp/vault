package driver

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

var (
	retryableCodes          = []int32{11600, 11602, 10107, 13435, 13436, 189, 91, 7, 6, 89, 9001}
	nodeIsRecoveringCodes   = []int32{11600, 11602, 13436, 189, 91}
	notMasterCodes          = []int32{10107, 13435}
	nodeIsShuttingDownCodes = []int32{11600, 91}
)

var (
	// UnknownTransactionCommitResult is an error label for unknown transaction commit results.
	UnknownTransactionCommitResult = "UnknownTransactionCommitResult"
	// TransientTransactionError is an error label for transient errors with transactions.
	TransientTransactionError = "TransientTransactionError"
	// NetworkError is an error label for network errors.
	NetworkError = "NetworkError"
	// ErrCursorNotFound is the cursor not found error for legacy find operations.
	ErrCursorNotFound = errors.New("cursor not found")
	// ErrUnacknowledgedWrite is returned from functions that have an unacknowledged
	// write concern.
	ErrUnacknowledgedWrite = errors.New("unacknowledged write")
	// ErrUnsupportedStorageEngine is returned when a retryable write is attempted against a server
	// that uses a storage engine that does not support retryable writes
	ErrUnsupportedStorageEngine = errors.New("this MongoDB deployment does not support retryable writes. Please add retryWrites=false to your connection string")
)

// QueryFailureError is an error representing a command failure as a document.
type QueryFailureError struct {
	Message  string
	Response bsoncore.Document
}

// Error implements the error interface.
func (e QueryFailureError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Response)
}

// ResponseError is an error parsing the response to a command.
type ResponseError struct {
	Message string
	Wrapped error
}

// NewCommandResponseError creates a CommandResponseError.
func NewCommandResponseError(msg string, err error) ResponseError {
	return ResponseError{Message: msg, Wrapped: err}
}

// Error implements the error interface.
func (e ResponseError) Error() string {
	if e.Wrapped != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Wrapped)
	}
	return fmt.Sprintf("%s", e.Message)
}

// WriteCommandError is an error for a write command.
type WriteCommandError struct {
	WriteConcernError *WriteConcernError
	WriteErrors       WriteErrors
}

// UnsupportedStorageEngine returns whether or not the WriteCommandError comes from a retryable write being attempted
// against a server that has a storage engine where they are not supported
func (wce WriteCommandError) UnsupportedStorageEngine() bool {
	for _, writeError := range wce.WriteErrors {
		if writeError.Code == 20 && strings.HasPrefix(strings.ToLower(writeError.Message), "transaction numbers") {
			return true
		}
	}
	return false
}

func (wce WriteCommandError) Error() string {
	var buf bytes.Buffer
	fmt.Fprint(&buf, "write command error: [")
	fmt.Fprintf(&buf, "{%s}, ", wce.WriteErrors)
	fmt.Fprintf(&buf, "{%s}]", wce.WriteConcernError)
	return buf.String()
}

// Retryable returns true if the error is retryable
func (wce WriteCommandError) Retryable() bool {
	if wce.WriteConcernError == nil {
		return false
	}
	return (*wce.WriteConcernError).Retryable()
}

// WriteConcernError is a write concern failure that occurred as a result of a
// write operation.
type WriteConcernError struct {
	Name    string
	Code    int64
	Message string
	Details bsoncore.Document
}

func (wce WriteConcernError) Error() string {
	if wce.Name != "" {
		return fmt.Sprintf("(%v) %v", wce.Name, wce.Message)
	}
	return wce.Message
}

// Retryable returns true if the error is retryable
func (wce WriteConcernError) Retryable() bool {
	for _, code := range retryableCodes {
		if wce.Code == int64(code) {
			return true
		}
	}
	if strings.Contains(wce.Message, "not master") || strings.Contains(wce.Message, "node is recovering") {
		return true
	}

	return false
}

// NodeIsRecovering returns true if this error is a node is recovering error.
func (wce WriteConcernError) NodeIsRecovering() bool {
	for _, code := range nodeIsRecoveringCodes {
		if wce.Code == int64(code) {
			return true
		}
	}
	return strings.Contains(wce.Message, "node is recovering")
}

// NodeIsShuttingDown returns true if this error is a node is shutting down error.
func (wce WriteConcernError) NodeIsShuttingDown() bool {
	for _, code := range nodeIsShuttingDownCodes {
		if wce.Code == int64(code) {
			return true
		}
	}
	return strings.Contains(wce.Message, "node is shutting down")
}

// NotMaster returns true if this error is a not master error.
func (wce WriteConcernError) NotMaster() bool {
	for _, code := range notMasterCodes {
		if wce.Code == int64(code) {
			return true
		}
	}
	return strings.Contains(wce.Message, "not master")
}

// WriteError is a non-write concern failure that occurred as a result of a write
// operation.
type WriteError struct {
	Index   int64
	Code    int64
	Message string
}

func (we WriteError) Error() string { return we.Message }

// WriteErrors is a group of non-write concern failures that occurred as a result
// of a write operation.
type WriteErrors []WriteError

func (we WriteErrors) Error() string {
	var buf bytes.Buffer
	fmt.Fprint(&buf, "write errors: [")
	for idx, err := range we {
		if idx != 0 {
			fmt.Fprintf(&buf, ", ")
		}
		fmt.Fprintf(&buf, "{%s}", err)
	}
	fmt.Fprint(&buf, "]")
	return buf.String()
}

// Error is a command execution error from the database.
type Error struct {
	Code    int32
	Message string
	Labels  []string
	Name    string
	Wrapped error
}

// UnsupportedStorageEngine returns whether e came as a result of an unsupported storage engine
func (e Error) UnsupportedStorageEngine() bool {
	return e.Code == 20 && strings.HasPrefix(strings.ToLower(e.Message), "transaction numbers")
}

// Error implements the error interface.
func (e Error) Error() string {
	if e.Name != "" {
		return fmt.Sprintf("(%v) %v", e.Name, e.Message)
	}
	return e.Message
}

// HasErrorLabel returns true if the error contains the specified label.
func (e Error) HasErrorLabel(label string) bool {
	if e.Labels != nil {
		for _, l := range e.Labels {
			if l == label {
				return true
			}
		}
	}
	return false
}

// Retryable returns true if the error is retryable
func (e Error) Retryable() bool {
	for _, label := range e.Labels {
		if label == NetworkError {
			return true
		}
	}
	for _, code := range retryableCodes {
		if e.Code == code {
			return true
		}
	}
	if strings.Contains(e.Message, "not master") || strings.Contains(e.Message, "node is recovering") {
		return true
	}

	return false
}

// NetworkError returns true if the error is a network error.
func (e Error) NetworkError() bool {
	for _, label := range e.Labels {
		if label == NetworkError {
			return true
		}
	}
	return false
}

// NodeIsRecovering returns true if this error is a node is recovering error.
func (e Error) NodeIsRecovering() bool {
	for _, code := range nodeIsRecoveringCodes {
		if e.Code == code {
			return true
		}
	}
	return strings.Contains(e.Message, "node is recovering")
}

// NodeIsShuttingDown returns true if this error is a node is shutting down error.
func (e Error) NodeIsShuttingDown() bool {
	for _, code := range nodeIsShuttingDownCodes {
		if e.Code == code {
			return true
		}
	}
	return strings.Contains(e.Message, "node is shutting down")
}

// NotMaster returns true if this error is a not master error.
func (e Error) NotMaster() bool {
	for _, code := range notMasterCodes {
		if e.Code == code {
			return true
		}
	}
	return strings.Contains(e.Message, "not master")
}

// NamespaceNotFound returns true if this errors is a NamespaceNotFound error.
func (e Error) NamespaceNotFound() bool {
	return e.Code == 26 || e.Message == "ns not found"
}

// helper method to extract an error from a reader if there is one; first returned item is the
// error if it exists, the second holds parsing errors
func extractError(rdr bsoncore.Document) error {
	var errmsg, codeName string
	var code int32
	var labels []string
	var ok bool
	var wcError WriteCommandError
	elems, err := rdr.Elements()
	if err != nil {
		return err
	}

	for _, elem := range elems {
		switch elem.Key() {
		case "ok":
			switch elem.Value().Type {
			case bson.TypeInt32:
				if elem.Value().Int32() == 1 {
					ok = true
				}
			case bson.TypeInt64:
				if elem.Value().Int64() == 1 {
					ok = true
				}
			case bson.TypeDouble:
				if elem.Value().Double() == 1 {
					ok = true
				}
			}
		case "errmsg":
			if str, okay := elem.Value().StringValueOK(); okay {
				errmsg = str
			}
		case "codeName":
			if str, okay := elem.Value().StringValueOK(); okay {
				codeName = str
			}
		case "code":
			if c, okay := elem.Value().Int32OK(); okay {
				code = c
			}
		case "errorLabels":
			if arr, okay := elem.Value().ArrayOK(); okay {
				elems, err := arr.Elements()
				if err != nil {
					continue
				}
				for _, elem := range elems {
					if str, ok := elem.Value().StringValueOK(); ok {
						labels = append(labels, str)
					}
				}

			}
		case "writeErrors":
			arr, exists := elem.Value().ArrayOK()
			if !exists {
				break
			}
			vals, err := arr.Values()
			if err != nil {
				continue
			}
			for _, val := range vals {
				var we WriteError
				doc, exists := val.DocumentOK()
				if !exists {
					continue
				}
				if index, exists := doc.Lookup("index").AsInt64OK(); exists {
					we.Index = index
				}
				if code, exists := doc.Lookup("code").AsInt64OK(); exists {
					we.Code = code
				}
				if msg, exists := doc.Lookup("errmsg").StringValueOK(); exists {
					we.Message = msg
				}
				wcError.WriteErrors = append(wcError.WriteErrors, we)
			}
		case "writeConcernError":
			doc, exists := elem.Value().DocumentOK()
			if !exists {
				break
			}
			wcError.WriteConcernError = new(WriteConcernError)
			if code, exists := doc.Lookup("code").AsInt64OK(); exists {
				wcError.WriteConcernError.Code = code
			}
			if name, exists := doc.Lookup("codeName").StringValueOK(); exists {
				wcError.WriteConcernError.Name = name
			}
			if msg, exists := doc.Lookup("errmsg").StringValueOK(); exists {
				wcError.WriteConcernError.Message = msg
			}
			if info, exists := doc.Lookup("errInfo").DocumentOK(); exists {
				wcError.WriteConcernError.Details = make([]byte, len(info))
				copy(wcError.WriteConcernError.Details, info)
			}
		}
	}

	if !ok {
		if errmsg == "" {
			errmsg = "command failed"
		}

		return Error{
			Code:    code,
			Message: errmsg,
			Name:    codeName,
			Labels:  labels,
		}
	}

	if len(wcError.WriteErrors) > 0 || wcError.WriteConcernError != nil {
		return wcError
	}

	return nil
}
