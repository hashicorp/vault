// Copyright 2021 Couchbase
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gocbcore

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrNoAttempt indicates no attempt was started before an operation was performed.
	ErrNoAttempt = errors.New("attempt was not started")

	// ErrOther indicates an non-specific error has occured.
	ErrOther = errors.New("other error")

	// ErrTransient indicates a transient error occured which may succeed at a later point in time.
	ErrTransient = errors.New("transient error")

	// ErrWriteWriteConflict indicates that another transaction conflicted with this one.
	ErrWriteWriteConflict = errors.New("write write conflict")

	// ErrHard indicates that an unrecoverable error occured.
	ErrHard = errors.New("hard")

	// ErrAmbiguous indicates that a failure occured but the outcome was not known.
	ErrAmbiguous = errors.New("ambiguous error")

	// ErrAtrFull indicates that the ATR record was too full to accept a new mutation.
	ErrAtrFull = errors.New("atr full")

	// ErrAttemptExpired indicates an attempt expired.
	ErrAttemptExpired = errors.New("attempt expired")

	// ErrAtrNotFound indicates that an expected ATR document was missing.
	ErrAtrNotFound = errors.New("atr not found")

	// ErrAtrEntryNotFound indicates that an expected ATR entry was missing.
	ErrAtrEntryNotFound = errors.New("atr entry not found")

	// ErrDocAlreadyInTransaction indicates that a document is already in a transaction.
	ErrDocAlreadyInTransaction = errors.New("doc already in transaction")

	// ErrIllegalState is used for when a transaction enters an illegal State.
	ErrIllegalState = errors.New("illegal State")

	// ErrTransactionAbortedExternally indicates the transaction was aborted externally.
	ErrTransactionAbortedExternally = errors.New("transaction aborted externally")

	// ErrPreviousOperationFailed indicates a previous operation in the transaction failed.
	ErrPreviousOperationFailed = errors.New("previous operation failed")

	// ErrForwardCompatibilityFailure indicates an operation failed due to involving a document in another transaction
	// which contains features this transaction does not support.
	ErrForwardCompatibilityFailure = errors.New("forward compatibility error")
)

type classifiedError struct {
	Source error
	Class  TransactionErrorClass
}

func (ce classifiedError) Wrap(errType error) *classifiedError {
	return &classifiedError{
		Source: &basicRetypedError{
			ErrType: errType,
			Source:  ce.Source,
		},
		Class: ce.Class,
	}
}

// TransactionOperationFailedError is used when a transaction operation fails.
// Internal: This should never be used and is not supported.
type TransactionOperationFailedError struct {
	shouldNotRetry    bool
	shouldNotRollback bool
	errorCause        error
	shouldRaise       TransactionErrorReason
	errorClass        TransactionErrorClass
}

// MarshalJSON will marshal this error for the wire.
func (tfe TransactionOperationFailedError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Retry    bool            `json:"retry"`
		Rollback bool            `json:"rollback"`
		Raise    string          `json:"raise"`
		Cause    json.RawMessage `json:"cause"`
	}{
		Retry:    !tfe.shouldNotRetry,
		Rollback: !tfe.shouldNotRollback,
		Raise:    tfe.shouldRaise.String(),
		Cause:    marshalErrorToJSON(tfe.errorCause),
	})
}

func (tfe TransactionOperationFailedError) Error() string {
	errStr := "transaction operation failed"
	errStr += " | " + fmt.Sprintf(
		"shouldRetry:%v, shouldRollback:%v, shouldRaise:%d, class:%d",
		!tfe.shouldNotRetry,
		!tfe.shouldNotRollback,
		tfe.shouldRaise,
		tfe.errorClass)
	if tfe.errorCause != nil {
		errStr += " | " + tfe.errorCause.Error()
	}
	return errStr
}

// Retry signals whether a new attempt should be made at rollback.
func (tfe TransactionOperationFailedError) Retry() bool {
	return !tfe.shouldNotRetry
}

// Rollback signals whether the attempt should be auto-rolled back.
func (tfe TransactionOperationFailedError) Rollback() bool {
	return !tfe.shouldNotRollback
}

// ToRaise signals which error type should be raised to the application.
func (tfe TransactionOperationFailedError) ToRaise() TransactionErrorReason {
	return tfe.shouldRaise
}

// ErrorClass returns the class of error which caused this error.
func (tfe TransactionOperationFailedError) ErrorClass() TransactionErrorClass {
	return tfe.errorClass
}

// InternalUnwrap returns the underlying error for this error.
func (tfe TransactionOperationFailedError) InternalUnwrap() error {
	return tfe.errorCause
}

type aggregateError []error

func (agge aggregateError) MarshalJSON() ([]byte, error) {
	suberrs := make([]json.RawMessage, len(agge))
	for i, err := range agge {
		suberrs[i] = marshalErrorToJSON(err)
	}
	return json.Marshal(suberrs)
}

func (agge aggregateError) Error() string {
	errStrs := []string{}
	for _, err := range agge {
		errStrs = append(errStrs, err.Error())
	}
	return "[" + strings.Join(errStrs, ", ") + "]"
}

func (agge aggregateError) Is(err error) bool {
	for _, aerr := range agge {
		if errors.Is(aerr, err) {
			return true
		}
	}
	return false
}

type writeWriteConflictError struct {
	BucketName     string
	ScopeName      string
	CollectionName string
	DocumentKey    []byte
	Source         error
}

func (wwce writeWriteConflictError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Msg            string          `json:"msg"`
		Cause          json.RawMessage `json:"cause"`
		BucketName     string          `json:"bucket"`
		ScopeName      string          `json:"scope"`
		CollectionName string          `json:"collection"`
		DocumentKey    string          `json:"document_key"`
	}{
		Msg:            "write write conflict",
		Cause:          marshalErrorToJSON(wwce.Source),
		BucketName:     wwce.BucketName,
		ScopeName:      wwce.ScopeName,
		CollectionName: wwce.CollectionName,
		DocumentKey:    string(wwce.DocumentKey),
	})
}

func (wwce writeWriteConflictError) Error() string {
	errStr := "write write conflict"
	errStr += " | " + fmt.Sprintf(
		"bucket:%s, scope:%s, collection:%s, key:%s",
		wwce.BucketName,
		wwce.ScopeName,
		wwce.CollectionName,
		wwce.DocumentKey)
	if wwce.Source != nil {
		errStr += " | " + wwce.Source.Error()
	}
	return errStr
}

func (wwce writeWriteConflictError) Is(err error) bool {
	if err == ErrWriteWriteConflict {
		return true
	}
	return errors.Is(wwce.Source, err)
}

func (wwce writeWriteConflictError) Unwrap() error {
	return wwce.Source
}

type basicRetypedError struct {
	ErrType error
	Source  error
}

func (bre basicRetypedError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Msg   string          `json:"msg"`
		Cause json.RawMessage `json:"cause"`
	}{
		Msg:   bre.ErrType.Error(),
		Cause: marshalErrorToJSON(bre.Source),
	})
}

func (bre basicRetypedError) Error() string {
	errStr := bre.ErrType.Error()
	if bre.Source != nil {
		errStr += " | " + bre.Source.Error()
	}
	return errStr
}

func (bre basicRetypedError) Is(err error) bool {
	if errors.Is(bre.ErrType, err) {
		return true
	}
	return errors.Is(bre.Source, err)
}

func (bre basicRetypedError) Unwrap() error {
	return bre.Source
}

type forwardCompatError struct {
	BucketName     string
	ScopeName      string
	CollectionName string
	DocumentKey    []byte
}

func (fce forwardCompatError) Error() string {
	errStr := ErrForwardCompatibilityFailure.Error()
	errStr += " | " + fmt.Sprintf(
		"bucket:%s, scope:%s, collection:%s, key:%s",
		fce.BucketName,
		fce.ScopeName,
		fce.CollectionName,
		fce.DocumentKey)
	return errStr
}

func (fce forwardCompatError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		BucketName     string `json:"bucket,omitempty"`
		ScopeName      string `json:"scope,omitempty"`
		CollectionName string `json:"collection,omitempty"`
		DocumentKey    string `json:"document_key,omitempty"`
		Message        string `json:"msg"`
	}{
		BucketName:     fce.BucketName,
		ScopeName:      fce.ScopeName,
		CollectionName: fce.CollectionName,
		DocumentKey:    string(fce.DocumentKey),
		Message:        ErrForwardCompatibilityFailure.Error(),
	})
}

func (fce forwardCompatError) Unwrap() error {
	return ErrForwardCompatibilityFailure
}

func marshalErrorToJSON(err error) json.RawMessage {
	if marshaler, ok := err.(json.Marshaler); ok {
		if data, err := marshaler.MarshalJSON(); err == nil {
			return data
		}
	}

	data, err := json.Marshal(err.Error())
	if err != nil {
		logWarnf("Failed to marshal error: %v", err)
	}
	return data
}
