// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/internal/csot"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// LegacyNotPrimaryErrMsg is the error message that older MongoDB servers (see
// SERVER-50412 for versions) return when a write operation is erroneously sent
// to a non-primary node.
const LegacyNotPrimaryErrMsg = "not master"

var (
	retryableCodes          = []int32{11600, 11602, 10107, 13435, 13436, 189, 91, 7, 6, 89, 9001, 262}
	nodeIsRecoveringCodes   = []int32{11600, 11602, 13436, 189, 91}
	notPrimaryCodes         = []int32{10107, 13435, 10058}
	nodeIsShuttingDownCodes = []int32{11600, 91}

	unknownReplWriteConcernCode   = int32(79)
	unsatisfiableWriteConcernCode = int32(100)
)

var (
	// UnknownTransactionCommitResult is an error label for unknown transaction commit results.
	UnknownTransactionCommitResult = "UnknownTransactionCommitResult"
	// TransientTransactionError is an error label for transient errors with transactions.
	TransientTransactionError = "TransientTransactionError"
	// NetworkError is an error label for network errors.
	NetworkError = "NetworkError"
	// RetryableWriteError is an error label for retryable write errors.
	RetryableWriteError = "RetryableWriteError"
	// NoWritesPerformed is an error label indicated that no writes were performed for an operation.
	NoWritesPerformed = "NoWritesPerformed"
	// ErrCursorNotFound is the cursor not found error for legacy find operations.
	ErrCursorNotFound = errors.New("cursor not found")
	// ErrUnacknowledgedWrite is returned from functions that have an unacknowledged
	// write concern.
	ErrUnacknowledgedWrite = errors.New("unacknowledged write")
	// ErrUnsupportedStorageEngine is returned when a retryable write is attempted against a server
	// that uses a storage engine that does not support retryable writes
	ErrUnsupportedStorageEngine = errors.New("this MongoDB deployment does not support retryable writes. Please add retryWrites=false to your connection string")
	// ErrDeadlineWouldBeExceeded is returned when a Timeout set on an operation
	// would be exceeded if the operation were sent to the server. It wraps
	// context.DeadlineExceeded.
	ErrDeadlineWouldBeExceeded = fmt.Errorf(
		"operation not sent to server, as Timeout would be exceeded: %w",
		context.DeadlineExceeded)
	// ErrNegativeMaxTime is returned when MaxTime on an operation is a negative value.
	ErrNegativeMaxTime = errors.New("a negative value was provided for MaxTime on an operation")
)

// QueryFailureError is an error representing a command failure as a document.
type QueryFailureError struct {
	Message  string
	Response bsoncore.Document
	Wrapped  error
}

// Error implements the error interface.
func (e QueryFailureError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Response)
}

// Unwrap returns the underlying error.
func (e QueryFailureError) Unwrap() error {
	return e.Wrapped
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
	return e.Message
}

// WriteCommandError is an error for a write command.
type WriteCommandError struct {
	WriteConcernError *WriteConcernError
	WriteErrors       WriteErrors
	Labels            []string
	Raw               bsoncore.Document
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
func (wce WriteCommandError) Retryable(wireVersion *description.VersionRange) bool {
	for _, label := range wce.Labels {
		if label == RetryableWriteError {
			return true
		}
	}
	if wireVersion != nil && wireVersion.Max >= 9 {
		return false
	}

	if wce.WriteConcernError == nil {
		return false
	}
	return wce.WriteConcernError.Retryable()
}

// HasErrorLabel returns true if the error contains the specified label.
func (wce WriteCommandError) HasErrorLabel(label string) bool {
	if wce.Labels != nil {
		for _, l := range wce.Labels {
			if l == label {
				return true
			}
		}
	}
	return false
}

// WriteConcernError is a write concern failure that occurred as a result of a
// write operation.
type WriteConcernError struct {
	Name            string
	Code            int64
	Message         string
	Details         bsoncore.Document
	Labels          []string
	TopologyVersion *description.TopologyVersion
	Raw             bsoncore.Document
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

	return false
}

// NodeIsRecovering returns true if this error is a node is recovering error.
func (wce WriteConcernError) NodeIsRecovering() bool {
	for _, code := range nodeIsRecoveringCodes {
		if wce.Code == int64(code) {
			return true
		}
	}
	hasNoCode := wce.Code == 0
	return hasNoCode && strings.Contains(wce.Message, "node is recovering")
}

// NodeIsShuttingDown returns true if this error is a node is shutting down error.
func (wce WriteConcernError) NodeIsShuttingDown() bool {
	for _, code := range nodeIsShuttingDownCodes {
		if wce.Code == int64(code) {
			return true
		}
	}
	hasNoCode := wce.Code == 0
	return hasNoCode && strings.Contains(wce.Message, "node is shutting down")
}

// NotPrimary returns true if this error is a not primary error.
func (wce WriteConcernError) NotPrimary() bool {
	for _, code := range notPrimaryCodes {
		if wce.Code == int64(code) {
			return true
		}
	}
	hasNoCode := wce.Code == 0
	return hasNoCode && strings.Contains(wce.Message, LegacyNotPrimaryErrMsg)
}

// WriteError is a non-write concern failure that occurred as a result of a write
// operation.
type WriteError struct {
	Index   int64
	Code    int64
	Message string
	Details bsoncore.Document
	Raw     bsoncore.Document
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
	Code            int32
	Message         string
	Labels          []string
	Name            string
	Wrapped         error
	TopologyVersion *description.TopologyVersion
	Raw             bsoncore.Document
}

// UnsupportedStorageEngine returns whether e came as a result of an unsupported storage engine
func (e Error) UnsupportedStorageEngine() bool {
	return e.Code == 20 && strings.HasPrefix(strings.ToLower(e.Message), "transaction numbers")
}

// Error implements the error interface.
func (e Error) Error() string {
	var msg string
	if e.Name != "" {
		msg = fmt.Sprintf("(%v)", e.Name)
	}
	msg += " " + e.Message
	if e.Wrapped != nil {
		msg += ": " + e.Wrapped.Error()
	}
	return msg
}

// Unwrap returns the underlying error.
func (e Error) Unwrap() error {
	return e.Wrapped
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

// RetryableRead returns true if the error is retryable for a read operation
func (e Error) RetryableRead() bool {
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

	return false
}

// RetryableWrite returns true if the error is retryable for a write operation
func (e Error) RetryableWrite(wireVersion *description.VersionRange) bool {
	for _, label := range e.Labels {
		if label == NetworkError || label == RetryableWriteError {
			return true
		}
	}
	if wireVersion != nil && wireVersion.Max >= 9 {
		return false
	}
	for _, code := range retryableCodes {
		if e.Code == code {
			return true
		}
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
	hasNoCode := e.Code == 0
	return hasNoCode && strings.Contains(e.Message, "node is recovering")
}

// NodeIsShuttingDown returns true if this error is a node is shutting down error.
func (e Error) NodeIsShuttingDown() bool {
	for _, code := range nodeIsShuttingDownCodes {
		if e.Code == code {
			return true
		}
	}
	hasNoCode := e.Code == 0
	return hasNoCode && strings.Contains(e.Message, "node is shutting down")
}

// NotPrimary returns true if this error is a not primary error.
func (e Error) NotPrimary() bool {
	for _, code := range notPrimaryCodes {
		if e.Code == code {
			return true
		}
	}
	hasNoCode := e.Code == 0
	return hasNoCode && strings.Contains(e.Message, LegacyNotPrimaryErrMsg)
}

// NamespaceNotFound returns true if this errors is a NamespaceNotFound error.
func (e Error) NamespaceNotFound() bool {
	return e.Code == 26 || e.Message == "ns not found"
}

// ExtractErrorFromServerResponse extracts an error from a server response bsoncore.Document
// if there is one. Also used in testing for SDAM.
func ExtractErrorFromServerResponse(ctx context.Context, doc bsoncore.Document) error {
	var errmsg, codeName string
	var code int32
	var labels []string
	var ok bool
	var tv *description.TopologyVersion
	var wcError WriteCommandError
	elems, err := doc.Elements()
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
			case bson.TypeBoolean:
				if elem.Value().Boolean() {
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
				vals, err := arr.Values()
				if err != nil {
					continue
				}
				for _, val := range vals {
					if str, ok := val.StringValueOK(); ok {
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
				if info, exists := doc.Lookup("errInfo").DocumentOK(); exists {
					we.Details = make([]byte, len(info))
					copy(we.Details, info)
				}
				we.Raw = doc
				wcError.WriteErrors = append(wcError.WriteErrors, we)
			}
		case "writeConcernError":
			doc, exists := elem.Value().DocumentOK()
			if !exists {
				break
			}
			wcError.WriteConcernError = new(WriteConcernError)
			wcError.WriteConcernError.Raw = doc
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
			if errLabels, exists := doc.Lookup("errorLabels").ArrayOK(); exists {
				vals, err := errLabels.Values()
				if err != nil {
					continue
				}
				for _, val := range vals {
					if str, ok := val.StringValueOK(); ok {
						labels = append(labels, str)
					}
				}
			}
		case "topologyVersion":
			doc, ok := elem.Value().DocumentOK()
			if !ok {
				break
			}
			version, err := description.NewTopologyVersion(bson.Raw(doc))
			if err == nil {
				tv = version
			}
		}
	}

	if !ok {
		if errmsg == "" {
			errmsg = "command failed"
		}

		err := Error{
			Code:            code,
			Message:         errmsg,
			Name:            codeName,
			Labels:          labels,
			TopologyVersion: tv,
			Raw:             doc,
		}

		// If CSOT is enabled and we get a MaxTimeMSExpired error, assume that
		// the error was caused by setting "maxTimeMS" on the command based on
		// the context deadline or on "timeoutMS". In that case, make the error
		// wrap context.DeadlineExceeded so that users can always check
		//
		//  errors.Is(err, context.DeadlineExceeded)
		//
		// for either client-side or server-side timeouts.
		if csot.IsTimeoutContext(ctx) && err.Code == 50 {
			err.Wrapped = context.DeadlineExceeded
		}

		return err
	}

	if len(wcError.WriteErrors) > 0 || wcError.WriteConcernError != nil {
		wcError.Labels = labels
		if wcError.WriteConcernError != nil {
			wcError.WriteConcernError.TopologyVersion = tv
		}
		wcError.Raw = doc
		return wcError
	}

	return nil
}
