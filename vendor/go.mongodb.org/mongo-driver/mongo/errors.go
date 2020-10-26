// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"bytes"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
)

// ErrUnacknowledgedWrite is returned by operations that have an unacknowledged write concern.
var ErrUnacknowledgedWrite = errors.New("unacknowledged write")

// ErrClientDisconnected is returned when disconnected Client is used to run an operation.
var ErrClientDisconnected = errors.New("client is disconnected")

// ErrNilDocument is returned when a nil document is passed to a CRUD method.
var ErrNilDocument = errors.New("document is nil")

// ErrNilValue is returned when a nil value is passed to a CRUD method.
var ErrNilValue = errors.New("value is nil")

// ErrEmptySlice is returned when an empty slice is passed to a CRUD method that requires a non-empty slice.
var ErrEmptySlice = errors.New("must provide at least one element in input slice")

func replaceErrors(err error) error {
	if err == topology.ErrTopologyClosed {
		return ErrClientDisconnected
	}
	if de, ok := err.(driver.Error); ok {
		return CommandError{
			Code:    de.Code,
			Message: de.Message,
			Labels:  de.Labels,
			Name:    de.Name,
			Wrapped: de.Wrapped,
		}
	}
	if qe, ok := err.(driver.QueryFailureError); ok {
		// qe.Message is "command failure"
		ce := CommandError{
			Name:    qe.Message,
			Wrapped: qe.Wrapped,
		}

		dollarErr, err := qe.Response.LookupErr("$err")
		if err == nil {
			ce.Message, _ = dollarErr.StringValueOK()
		}
		code, err := qe.Response.LookupErr("code")
		if err == nil {
			ce.Code, _ = code.Int32OK()
		}

		return ce
	}
	if me, ok := err.(mongocrypt.Error); ok {
		return MongocryptError{Code: me.Code, Message: me.Message}
	}

	return err
}

// MongocryptError represents an libmongocrypt error during client-side encryption.
type MongocryptError struct {
	Code    int32
	Message string
}

// Error implements the error interface.
func (m MongocryptError) Error() string {
	return fmt.Sprintf("mongocrypt error %d: %v", m.Code, m.Message)
}

// EncryptionKeyVaultError represents an error while communicating with the key vault collection during client-side
// encryption.
type EncryptionKeyVaultError struct {
	Wrapped error
}

// Error implements the error interface.
func (ekve EncryptionKeyVaultError) Error() string {
	return fmt.Sprintf("key vault communication error: %v", ekve.Wrapped)
}

// Unwrap returns the underlying error.
func (ekve EncryptionKeyVaultError) Unwrap() error {
	return ekve.Wrapped
}

// MongocryptdError represents an error while communicating with mongocryptd during client-side encryption.
type MongocryptdError struct {
	Wrapped error
}

// Error implements the error interface.
func (e MongocryptdError) Error() string {
	return fmt.Sprintf("mongocryptd communication error: %v", e.Wrapped)
}

// Unwrap returns the underlying error.
func (e MongocryptdError) Unwrap() error {
	return e.Wrapped
}

// CommandError represents a server error during execution of a command. This can be returned by any operation.
type CommandError struct {
	Code    int32
	Message string
	Labels  []string // Categories to which the error belongs
	Name    string   // A human-readable name corresponding to the error code
	Wrapped error    // The underlying error, if one exists.
}

// Error implements the error interface.
func (e CommandError) Error() string {
	if e.Name != "" {
		return fmt.Sprintf("(%v) %v", e.Name, e.Message)
	}
	return e.Message
}

// Unwrap returns the underlying error.
func (e CommandError) Unwrap() error {
	return e.Wrapped
}

// HasErrorLabel returns true if the error contains the specified label.
func (e CommandError) HasErrorLabel(label string) bool {
	if e.Labels != nil {
		for _, l := range e.Labels {
			if l == label {
				return true
			}
		}
	}
	return false
}

// IsMaxTimeMSExpiredError returns true if the error is a MaxTimeMSExpired error.
func (e CommandError) IsMaxTimeMSExpiredError() bool {
	return e.Code == 50 || e.Name == "MaxTimeMSExpired"
}

// WriteError is an error that occurred during execution of a write operation. This error type is only returned as part
// of a WriteException or BulkWriteException.
type WriteError struct {
	// The index of the write in the slice passed to an InsertMany or BulkWrite operation that caused this error.
	Index int

	Code    int
	Message string
}

func (we WriteError) Error() string { return we.Message }

// WriteErrors is a group of write errors that occurred during execution of a write operation.
type WriteErrors []WriteError

// Error implements the error interface.
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

func writeErrorsFromDriverWriteErrors(errs driver.WriteErrors) WriteErrors {
	wes := make(WriteErrors, 0, len(errs))
	for _, err := range errs {
		wes = append(wes, WriteError{Index: int(err.Index), Code: int(err.Code), Message: err.Message})
	}
	return wes
}

// WriteConcernError represents a write concern failure during execution of a write operation. This error type is only
// returned as part of a WriteException or a BulkWriteException.
type WriteConcernError struct {
	Name    string
	Code    int
	Message string
	Details bson.Raw
}

// Error implements the error interface.
func (wce WriteConcernError) Error() string {
	if wce.Name != "" {
		return fmt.Sprintf("(%v) %v", wce.Name, wce.Message)
	}
	return wce.Message
}

// WriteException is the error type returned by the InsertOne, DeleteOne, DeleteMany, UpdateOne, UpdateMany, and
// ReplaceOne operations.
type WriteException struct {
	// The write concern error that occurred, or nil if there was none.
	WriteConcernError *WriteConcernError

	// The write errors that occurred during operation execution.
	WriteErrors WriteErrors

	// The categories to which the exception belongs.
	Labels []string
}

// Error implements the error interface.
func (mwe WriteException) Error() string {
	var buf bytes.Buffer
	fmt.Fprint(&buf, "multiple write errors: [")
	fmt.Fprintf(&buf, "{%s}, ", mwe.WriteErrors)
	fmt.Fprintf(&buf, "{%s}]", mwe.WriteConcernError)
	return buf.String()
}

// HasErrorLabel returns true if the error contains the specified label.
func (mwe WriteException) HasErrorLabel(label string) bool {
	if mwe.Labels != nil {
		for _, l := range mwe.Labels {
			if l == label {
				return true
			}
		}
	}
	return false
}

func convertDriverWriteConcernError(wce *driver.WriteConcernError) *WriteConcernError {
	if wce == nil {
		return nil
	}

	return &WriteConcernError{
		Name:    wce.Name,
		Code:    int(wce.Code),
		Message: wce.Message,
		Details: bson.Raw(wce.Details),
	}
}

// BulkWriteError is an error that occurred during execution of one operation in a BulkWrite. This error type is only
// returned as part of a BulkWriteException.
type BulkWriteError struct {
	WriteError            // The WriteError that occurred.
	Request    WriteModel // The WriteModel that caused this error.
}

// Error implements the error interface.
func (bwe BulkWriteError) Error() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "{%s}", bwe.WriteError)
	return buf.String()
}

// BulkWriteException is the error type returned by BulkWrite and InsertMany operations.
type BulkWriteException struct {
	// The write concern error that occurred, or nil if there was none.
	WriteConcernError *WriteConcernError

	// The write errors that occurred during operation execution.
	WriteErrors []BulkWriteError

	// The categories to which the exception belongs.
	Labels []string
}

// Error implements the error interface.
func (bwe BulkWriteException) Error() string {
	var buf bytes.Buffer
	fmt.Fprint(&buf, "bulk write error: [")
	fmt.Fprintf(&buf, "{%s}, ", bwe.WriteErrors)
	fmt.Fprintf(&buf, "{%s}]", bwe.WriteConcernError)
	return buf.String()
}

// HasErrorLabel returns true if the error contains the specified label.
func (bwe BulkWriteException) HasErrorLabel(label string) bool {
	if bwe.Labels != nil {
		for _, l := range bwe.Labels {
			if l == label {
				return true
			}
		}
	}
	return false
}

// returnResult is used to determine if a function calling processWriteError should return
// the result or return nil. Since the processWriteError function is used by many different
// methods, both *One and *Many, we need a way to differentiate if the method should return
// the result and the error.
type returnResult int

const (
	rrNone returnResult = 1 << iota // None means do not return the result ever.
	rrOne                           // One means return the result if this was called by a *One method.
	rrMany                          // Many means return the result is this was called by a *Many method.

	rrAll returnResult = rrOne | rrMany // All means always return the result.
)

// processWriteError handles processing the result of a write operation. If the retrunResult matches
// the calling method's type, it should return the result object in addition to the error.
// This function will wrap the errors from other packages and return them as errors from this package.
//
// WriteConcernError will be returned over WriteErrors if both are present.
func processWriteError(err error) (returnResult, error) {
	switch {
	case err == driver.ErrUnacknowledgedWrite:
		return rrAll, ErrUnacknowledgedWrite
	case err != nil:
		switch tt := err.(type) {
		case driver.WriteCommandError:
			return rrMany, WriteException{
				WriteConcernError: convertDriverWriteConcernError(tt.WriteConcernError),
				WriteErrors:       writeErrorsFromDriverWriteErrors(tt.WriteErrors),
				Labels:            tt.Labels,
			}
		default:
			return rrNone, replaceErrors(err)
		}
	default:
		return rrAll, nil
	}
}
