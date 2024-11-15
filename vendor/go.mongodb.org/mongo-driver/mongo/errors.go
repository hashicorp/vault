// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/internal/codecutil"
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

// ErrMapForOrderedArgument is returned when a map with multiple keys is passed to a CRUD method for an ordered parameter
type ErrMapForOrderedArgument struct {
	ParamName string
}

// Error implements the error interface.
func (e ErrMapForOrderedArgument) Error() string {
	return fmt.Sprintf("multi-key map passed in for ordered parameter %v", e.ParamName)
}

func replaceErrors(err error) error {
	// Return nil when err is nil to avoid costly reflection logic below.
	if err == nil {
		return nil
	}

	if errors.Is(err, topology.ErrTopologyClosed) {
		return ErrClientDisconnected
	}
	if de, ok := err.(driver.Error); ok {
		return CommandError{
			Code:    de.Code,
			Message: de.Message,
			Labels:  de.Labels,
			Name:    de.Name,
			Wrapped: de.Wrapped,
			Raw:     bson.Raw(de.Raw),
		}
	}
	if qe, ok := err.(driver.QueryFailureError); ok {
		// qe.Message is "command failure"
		ce := CommandError{
			Name:    qe.Message,
			Wrapped: qe.Wrapped,
			Raw:     bson.Raw(qe.Response),
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

	if errors.Is(err, codecutil.ErrNilValue) {
		return ErrNilValue
	}

	if marshalErr, ok := err.(codecutil.MarshalError); ok {
		return MarshalError{
			Value: marshalErr.Value,
			Err:   marshalErr.Err,
		}
	}

	return err
}

// IsDuplicateKeyError returns true if err is a duplicate key error.
func IsDuplicateKeyError(err error) bool {
	if se := ServerError(nil); errors.As(err, &se) {
		return se.HasErrorCode(11000) || // Duplicate key error.
			se.HasErrorCode(11001) || // Duplicate key error on update.
			// Duplicate key error in a capped collection. See SERVER-7164.
			se.HasErrorCode(12582) ||
			// Mongos insert error caused by a duplicate key error. See
			// SERVER-11493.
			se.HasErrorCodeWithMessage(16460, " E11000 ")
	}
	return false
}

// timeoutErrs is a list of error values that indicate a timeout happened.
var timeoutErrs = [...]error{
	context.DeadlineExceeded,
	driver.ErrDeadlineWouldBeExceeded,
	topology.ErrServerSelectionTimeout,
}

// IsTimeout returns true if err was caused by a timeout. For error chains,
// IsTimeout returns true if any error in the chain was caused by a timeout.
func IsTimeout(err error) bool {
	// Check if the error chain contains any of the timeout error values.
	for _, target := range timeoutErrs {
		if errors.Is(err, target) {
			return true
		}
	}

	// Check if the error chain contains any error types that can indicate
	// timeout.
	if errors.As(err, &topology.WaitQueueTimeoutError{}) {
		return true
	}
	if ce := (CommandError{}); errors.As(err, &ce) && ce.IsMaxTimeMSExpiredError() {
		return true
	}
	if we := (WriteException{}); errors.As(err, &we) && we.WriteConcernError != nil && we.WriteConcernError.IsMaxTimeMSExpiredError() {
		return true
	}
	if ne := net.Error(nil); errors.As(err, &ne) {
		return ne.Timeout()
	}
	// Check timeout error labels.
	if le := LabeledError(nil); errors.As(err, &le) {
		if le.HasErrorLabel("NetworkTimeoutError") || le.HasErrorLabel("ExceededTimeLimitError") {
			return true
		}
	}

	return false
}

// unwrap returns the inner error if err implements Unwrap(), otherwise it returns nil.
func unwrap(err error) error {
	u, ok := err.(interface {
		Unwrap() error
	})
	if !ok {
		return nil
	}
	return u.Unwrap()
}

// errorHasLabel returns true if err contains the specified label
func errorHasLabel(err error, label string) bool {
	for ; err != nil; err = unwrap(err) {
		if le, ok := err.(LabeledError); ok && le.HasErrorLabel(label) {
			return true
		}
	}
	return false
}

// IsNetworkError returns true if err is a network error
func IsNetworkError(err error) bool {
	return errorHasLabel(err, "NetworkError")
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

// LabeledError is an interface for errors with labels.
type LabeledError interface {
	error
	// HasErrorLabel returns true if the error contains the specified label.
	HasErrorLabel(string) bool
}

// ServerError is the interface implemented by errors returned from the server. Custom implementations of this
// interface should not be used in production.
type ServerError interface {
	LabeledError
	// HasErrorCode returns true if the error has the specified code.
	HasErrorCode(int) bool
	// HasErrorMessage returns true if the error contains the specified message.
	HasErrorMessage(string) bool
	// HasErrorCodeWithMessage returns true if any of the contained errors have the specified code and message.
	HasErrorCodeWithMessage(int, string) bool

	serverError()
}

var _ ServerError = CommandError{}
var _ ServerError = WriteError{}
var _ ServerError = WriteException{}
var _ ServerError = BulkWriteException{}

// CommandError represents a server error during execution of a command. This can be returned by any operation.
type CommandError struct {
	Code    int32
	Message string
	Labels  []string // Categories to which the error belongs
	Name    string   // A human-readable name corresponding to the error code
	Wrapped error    // The underlying error, if one exists.
	Raw     bson.Raw // The original server response containing the error.
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

// HasErrorCode returns true if the error has the specified code.
func (e CommandError) HasErrorCode(code int) bool {
	return int(e.Code) == code
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

// HasErrorMessage returns true if the error contains the specified message.
func (e CommandError) HasErrorMessage(message string) bool {
	return strings.Contains(e.Message, message)
}

// HasErrorCodeWithMessage returns true if the error has the specified code and Message contains the specified message.
func (e CommandError) HasErrorCodeWithMessage(code int, message string) bool {
	return int(e.Code) == code && strings.Contains(e.Message, message)
}

// IsMaxTimeMSExpiredError returns true if the error is a MaxTimeMSExpired error.
func (e CommandError) IsMaxTimeMSExpiredError() bool {
	return e.Code == 50 || e.Name == "MaxTimeMSExpired"
}

// serverError implements the ServerError interface.
func (e CommandError) serverError() {}

// WriteError is an error that occurred during execution of a write operation. This error type is only returned as part
// of a WriteException or BulkWriteException.
type WriteError struct {
	// The index of the write in the slice passed to an InsertMany or BulkWrite operation that caused this error.
	Index int

	Code    int
	Message string
	Details bson.Raw

	// The original write error from the server response.
	Raw bson.Raw
}

func (we WriteError) Error() string {
	msg := we.Message
	if len(we.Details) > 0 {
		msg = fmt.Sprintf("%s: %s", msg, we.Details.String())
	}
	return msg
}

// HasErrorCode returns true if the error has the specified code.
func (we WriteError) HasErrorCode(code int) bool {
	return we.Code == code
}

// HasErrorLabel returns true if the error contains the specified label. WriteErrors do not contain labels,
// so we always return false.
func (we WriteError) HasErrorLabel(string) bool {
	return false
}

// HasErrorMessage returns true if the error contains the specified message.
func (we WriteError) HasErrorMessage(message string) bool {
	return strings.Contains(we.Message, message)
}

// HasErrorCodeWithMessage returns true if the error has the specified code and Message contains the specified message.
func (we WriteError) HasErrorCodeWithMessage(code int, message string) bool {
	return we.Code == code && strings.Contains(we.Message, message)
}

// serverError implements the ServerError interface.
func (we WriteError) serverError() {}

// WriteErrors is a group of write errors that occurred during execution of a write operation.
type WriteErrors []WriteError

// Error implements the error interface.
func (we WriteErrors) Error() string {
	errs := make([]error, len(we))
	for i := 0; i < len(we); i++ {
		errs[i] = we[i]
	}
	// WriteErrors isn't returned from batch operations, but we can still use the same formatter.
	return "write errors: " + joinBatchErrors(errs)
}

func writeErrorsFromDriverWriteErrors(errs driver.WriteErrors) WriteErrors {
	wes := make(WriteErrors, 0, len(errs))
	for _, err := range errs {
		wes = append(wes, WriteError{
			Index:   int(err.Index),
			Code:    int(err.Code),
			Message: err.Message,
			Details: bson.Raw(err.Details),
			Raw:     bson.Raw(err.Raw),
		})
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
	Raw     bson.Raw // The original write concern error from the server response.
}

// Error implements the error interface.
func (wce WriteConcernError) Error() string {
	if wce.Name != "" {
		return fmt.Sprintf("(%v) %v", wce.Name, wce.Message)
	}
	return wce.Message
}

// IsMaxTimeMSExpiredError returns true if the error is a MaxTimeMSExpired error.
func (wce WriteConcernError) IsMaxTimeMSExpiredError() bool {
	return wce.Code == 50
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

	// The original server response containing the error.
	Raw bson.Raw
}

// Error implements the error interface.
func (mwe WriteException) Error() string {
	causes := make([]string, 0, 2)
	if mwe.WriteConcernError != nil {
		causes = append(causes, "write concern error: "+mwe.WriteConcernError.Error())
	}
	if len(mwe.WriteErrors) > 0 {
		// The WriteErrors error message already starts with "write errors:", so don't add it to the
		// error message again.
		causes = append(causes, mwe.WriteErrors.Error())
	}

	message := "write exception: "
	if len(causes) == 0 {
		return message + "no causes"
	}
	return message + strings.Join(causes, ", ")
}

// HasErrorCode returns true if the error has the specified code.
func (mwe WriteException) HasErrorCode(code int) bool {
	if mwe.WriteConcernError != nil && mwe.WriteConcernError.Code == code {
		return true
	}
	for _, we := range mwe.WriteErrors {
		if we.Code == code {
			return true
		}
	}
	return false
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

// HasErrorMessage returns true if the error contains the specified message.
func (mwe WriteException) HasErrorMessage(message string) bool {
	if mwe.WriteConcernError != nil && strings.Contains(mwe.WriteConcernError.Message, message) {
		return true
	}
	for _, we := range mwe.WriteErrors {
		if strings.Contains(we.Message, message) {
			return true
		}
	}
	return false
}

// HasErrorCodeWithMessage returns true if any of the contained errors have the specified code and message.
func (mwe WriteException) HasErrorCodeWithMessage(code int, message string) bool {
	if mwe.WriteConcernError != nil &&
		mwe.WriteConcernError.Code == code && strings.Contains(mwe.WriteConcernError.Message, message) {
		return true
	}
	for _, we := range mwe.WriteErrors {
		if we.Code == code && strings.Contains(we.Message, message) {
			return true
		}
	}
	return false
}

// serverError implements the ServerError interface.
func (mwe WriteException) serverError() {}

func convertDriverWriteConcernError(wce *driver.WriteConcernError) *WriteConcernError {
	if wce == nil {
		return nil
	}

	return &WriteConcernError{
		Name:    wce.Name,
		Code:    int(wce.Code),
		Message: wce.Message,
		Details: bson.Raw(wce.Details),
		Raw:     bson.Raw(wce.Raw),
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
	return bwe.WriteError.Error()
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
	causes := make([]string, 0, 2)
	if bwe.WriteConcernError != nil {
		causes = append(causes, "write concern error: "+bwe.WriteConcernError.Error())
	}
	if len(bwe.WriteErrors) > 0 {
		errs := make([]error, len(bwe.WriteErrors))
		for i := 0; i < len(bwe.WriteErrors); i++ {
			errs[i] = &bwe.WriteErrors[i]
		}
		causes = append(causes, "write errors: "+joinBatchErrors(errs))
	}

	message := "bulk write exception: "
	if len(causes) == 0 {
		return message + "no causes"
	}
	return "bulk write exception: " + strings.Join(causes, ", ")
}

// HasErrorCode returns true if any of the errors have the specified code.
func (bwe BulkWriteException) HasErrorCode(code int) bool {
	if bwe.WriteConcernError != nil && bwe.WriteConcernError.Code == code {
		return true
	}
	for _, we := range bwe.WriteErrors {
		if we.Code == code {
			return true
		}
	}
	return false
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

// HasErrorMessage returns true if the error contains the specified message.
func (bwe BulkWriteException) HasErrorMessage(message string) bool {
	if bwe.WriteConcernError != nil && strings.Contains(bwe.WriteConcernError.Message, message) {
		return true
	}
	for _, we := range bwe.WriteErrors {
		if strings.Contains(we.Message, message) {
			return true
		}
	}
	return false
}

// HasErrorCodeWithMessage returns true if any of the contained errors have the specified code and message.
func (bwe BulkWriteException) HasErrorCodeWithMessage(code int, message string) bool {
	if bwe.WriteConcernError != nil &&
		bwe.WriteConcernError.Code == code && strings.Contains(bwe.WriteConcernError.Message, message) {
		return true
	}
	for _, we := range bwe.WriteErrors {
		if we.Code == code && strings.Contains(we.Message, message) {
			return true
		}
	}
	return false
}

// serverError implements the ServerError interface.
func (bwe BulkWriteException) serverError() {}

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
	case errors.Is(err, driver.ErrUnacknowledgedWrite):
		return rrAll, ErrUnacknowledgedWrite
	case err != nil:
		switch tt := err.(type) {
		case driver.WriteCommandError:
			return rrMany, WriteException{
				WriteConcernError: convertDriverWriteConcernError(tt.WriteConcernError),
				WriteErrors:       writeErrorsFromDriverWriteErrors(tt.WriteErrors),
				Labels:            tt.Labels,
				Raw:               bson.Raw(tt.Raw),
			}
		default:
			return rrNone, replaceErrors(err)
		}
	default:
		return rrAll, nil
	}
}

// batchErrorsTargetLength is the target length of error messages returned by batch operation
// error types. Try to limit batch error messages to 2kb to prevent problems when printing error
// messages from large batch operations.
const batchErrorsTargetLength = 2000

// joinBatchErrors appends messages from the given errors to a comma-separated string. If the
// string exceeds 2kb, it stops appending error messages and appends the message "+N more errors..."
// to the end.
//
// Example format:
//
//	"[message 1, message 2, +8 more errors...]"
func joinBatchErrors(errs []error) string {
	var buf bytes.Buffer
	fmt.Fprint(&buf, "[")
	for idx, err := range errs {
		if idx != 0 {
			fmt.Fprint(&buf, ", ")
		}
		// If the error message has exceeded the target error message length, stop appending errors
		// to the message and append the number of remaining errors instead.
		if buf.Len() > batchErrorsTargetLength {
			fmt.Fprintf(&buf, "+%d more errors...", len(errs)-idx)
			break
		}
		fmt.Fprint(&buf, err.Error())
	}
	fmt.Fprint(&buf, "]")

	return buf.String()
}
