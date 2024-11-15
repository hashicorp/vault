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
	"strconv"
	"time"

	"github.com/google/uuid"
)

type addCleanupRequest func(req *TransactionsCleanupRequest) bool
type addLostCleanupLocation func(bucket, scope, collection string)

// Transaction represents a single active transaction, it can be used to
// stage mutations and finally commit them.
type Transaction struct {
	parent *TransactionsManager

	expiryTime              time.Time
	startTime               time.Time
	keyValueTimeout         time.Duration
	durabilityLevel         TransactionDurabilityLevel
	enableParallelUnstaging bool
	enableNonFatalGets      bool
	enableExplicitATRs      bool
	enableMutationCaching   bool
	atrLocation             TransactionATRLocation
	bucketAgentProvider     TransactionsBucketAgentProviderFn

	transactionID string
	attempt       *transactionAttempt
	hooks         TransactionHooks

	addCleanupRequest      addCleanupRequest
	addLostCleanupLocation addLostCleanupLocation

	recordResourceUnit resourceUnitCallback

	logger *internalTransactionLogWrapper
}

// ID returns the transaction ID of this transaction.
func (t *Transaction) ID() string {
	return t.transactionID
}

// Attempt returns meta-data about the current attempt to complete the transaction.
func (t *Transaction) Attempt() TransactionAttempt {
	if t.attempt == nil {
		return TransactionAttempt{}
	}

	return t.attempt.State()
}

// NewAttempt begins a new attempt with this transaction.
func (t *Transaction) NewAttempt() error {
	attemptUUID := uuid.New().String()

	t.attempt = &transactionAttempt{
		expiryTime:              t.expiryTime,
		txnStartTime:            t.startTime,
		keyValueTimeout:         t.keyValueTimeout,
		durabilityLevel:         t.durabilityLevel,
		transactionID:           t.transactionID,
		enableNonFatalGets:      t.enableNonFatalGets,
		enableParallelUnstaging: t.enableParallelUnstaging,
		enableMutationCaching:   t.enableMutationCaching,
		enableExplicitATRs:      t.enableExplicitATRs,
		atrLocation:             t.atrLocation,
		bucketAgentProvider:     t.bucketAgentProvider,

		id:                attemptUUID,
		state:             TransactionAttemptStateNothingWritten,
		stagedMutations:   nil,
		atrAgent:          nil,
		atrScopeName:      "",
		atrCollectionName: "",
		atrKey:            nil,
		hooks:             t.hooks,

		addCleanupRequest:      t.addCleanupRequest,
		addLostCleanupLocation: t.addLostCleanupLocation,
		logger:                 t.logger,

		recordResourceUnit: t.recordResourceUnit,
	}

	return nil
}

func (t *Transaction) resumeAttempt(txnData *jsonSerializedAttempt) error {
	if txnData.ID.Attempt == "" {
		return errors.New("invalid txn data - no attempt id")
	}

	attemptUUID := txnData.ID.Attempt

	var txnState TransactionAttemptState
	var atrAgent *Agent
	var atrOboUser string
	var atrScope, atrCollection string
	var atrKey []byte
	if txnData.ATR.ID != "" {
		// ATR references the specific ATR for this transaction.

		if txnData.ATR.Bucket == "" {
			return errors.New("invalid atr data - no bucket")
		}

		foundAtrAgent, foundAtrOboUser, err := t.parent.config.BucketAgentProvider(txnData.ATR.Bucket)
		if err != nil {
			return err
		}

		txnState = TransactionAttemptStatePending
		atrAgent = foundAtrAgent
		atrOboUser = foundAtrOboUser
		atrScope = txnData.ATR.Scope
		atrCollection = txnData.ATR.Collection
		atrKey = []byte(txnData.ATR.ID)
	} else {
		// No ATR information means its pending with no custom.

		txnState = TransactionAttemptStateNothingWritten
		atrAgent = nil
		atrOboUser = ""
		atrScope = ""
		atrCollection = ""
		atrKey = nil
	}

	stagedMutations := make([]*transactionStagedMutation, len(txnData.Mutations))
	for mutationIdx, mutationData := range txnData.Mutations {
		if mutationData.Bucket == "" {
			return errors.New("invalid staged mutation - no bucket")
		}
		if mutationData.ID == "" {
			return errors.New("invalid staged mutation - no key")
		}
		if mutationData.Cas == "" {
			return errors.New("invalid staged mutation - no cas")
		}
		if mutationData.Type == "" {
			return errors.New("invalid staged mutation - no type")
		}

		agent, oboUser, err := t.parent.config.BucketAgentProvider(mutationData.Bucket)
		if err != nil {
			return err
		}

		cas, err := strconv.ParseUint(mutationData.Cas, 10, 64)
		if err != nil {
			return err
		}

		opType, err := transactionStagedMutationTypeFromString(mutationData.Type)
		if err != nil {
			return err
		}

		stagedMutations[mutationIdx] = &transactionStagedMutation{
			OpType:         opType,
			Agent:          agent,
			OboUser:        oboUser,
			ScopeName:      mutationData.Scope,
			CollectionName: mutationData.Collection,
			Key:            []byte(mutationData.ID),
			Cas:            Cas(cas),
			Staged:         nil,
		}
	}

	t.attempt = &transactionAttempt{
		expiryTime:              t.expiryTime,
		txnStartTime:            t.startTime,
		keyValueTimeout:         t.keyValueTimeout,
		durabilityLevel:         t.durabilityLevel,
		transactionID:           t.transactionID,
		enableNonFatalGets:      t.enableNonFatalGets,
		enableParallelUnstaging: t.enableParallelUnstaging,
		enableMutationCaching:   t.enableMutationCaching,
		enableExplicitATRs:      t.enableExplicitATRs,
		atrLocation:             t.atrLocation,
		bucketAgentProvider:     t.bucketAgentProvider,

		id:                attemptUUID,
		state:             txnState,
		stagedMutations:   stagedMutations,
		atrAgent:          atrAgent,
		atrOboUser:        atrOboUser,
		atrScopeName:      atrScope,
		atrCollectionName: atrCollection,
		atrKey:            atrKey,
		hooks:             t.hooks,

		addCleanupRequest:      t.addCleanupRequest,
		addLostCleanupLocation: t.addLostCleanupLocation,

		logger: t.logger,

		recordResourceUnit: t.recordResourceUnit,
	}

	return nil
}

// TransactionGetOptions provides options for a Get operation.
type TransactionGetOptions struct {
	Agent          *Agent
	OboUser        string
	ScopeName      string
	CollectionName string
	Key            []byte

	// NoRYOW will disable the RYOW logic used to enable transactions
	// to naturally read any mutations they have performed.
	// VOLATILE: This parameter is subject to change.
	NoRYOW bool
}

// TransactionMutableItemMetaATR represents the ATR for meta.
type TransactionMutableItemMetaATR struct {
	BucketName     string `json:"bkt"`
	ScopeName      string `json:"scp"`
	CollectionName string `json:"coll"`
	DocID          string `json:"key"`
}

// TransactionMutableItemMeta represents all the meta-data for a fetched
// item.  Most of this is used for later mutation operations.
type TransactionMutableItemMeta struct {
	TransactionID string                                            `json:"txn"`
	AttemptID     string                                            `json:"atmpt"`
	ATR           TransactionMutableItemMetaATR                     `json:"atr"`
	ForwardCompat map[string][]TransactionForwardCompatibilityEntry `json:"fc,omitempty"`
}

// TransactionGetResult represents the result of a Get or GetOptional operation.
type TransactionGetResult struct {
	agent          *Agent
	oboUser        string
	scopeName      string
	collectionName string
	key            []byte

	Meta  *TransactionMutableItemMeta
	Value []byte
	Cas   Cas
}

// TransactionGetCallback describes a callback for a completed Get or GetOptional operation.
type TransactionGetCallback func(*TransactionGetResult, error)

// Get will attempt to fetch a document, and fail the transaction if it does not exist.
func (t *Transaction) Get(opts TransactionGetOptions, cb TransactionGetCallback) error {
	if t.attempt == nil {
		return ErrNoAttempt
	}

	return t.attempt.Get(opts, cb)
}

// TransactionInsertOptions provides options for a Insert operation.
type TransactionInsertOptions struct {
	Agent          *Agent
	OboUser        string
	ScopeName      string
	CollectionName string
	Key            []byte
	Value          json.RawMessage
}

// TransactionStoreCallback describes a callback for a completed Replace operation.
type TransactionStoreCallback func(*TransactionGetResult, error)

// Insert will attempt to insert a document.
func (t *Transaction) Insert(opts TransactionInsertOptions, cb TransactionStoreCallback) error {
	if t.attempt == nil {
		return ErrNoAttempt
	}

	return t.attempt.Insert(opts, cb)
}

// TransactionReplaceOptions provides options for a Replace operation.
type TransactionReplaceOptions struct {
	Document *TransactionGetResult
	Value    json.RawMessage
}

// Replace will attempt to replace an existing document.
func (t *Transaction) Replace(opts TransactionReplaceOptions, cb TransactionStoreCallback) error {
	if t.attempt == nil {
		return ErrNoAttempt
	}

	return t.attempt.Replace(opts, cb)
}

// TransactionRemoveOptions provides options for a Remove operation.
type TransactionRemoveOptions struct {
	Document *TransactionGetResult
}

// Remove will attempt to remove a previously fetched document.
func (t *Transaction) Remove(opts TransactionRemoveOptions, cb TransactionStoreCallback) error {
	if t.attempt == nil {
		return ErrNoAttempt
	}

	return t.attempt.Remove(opts, cb)
}

// TransactionCommitCallback describes a callback for a completed commit operation.
type TransactionCommitCallback func(error)

// Commit will attempt to commit the transaction, rolling it back and cancelling
// it if it is not capable of doing so.
func (t *Transaction) Commit(cb TransactionCommitCallback) error {
	if t.attempt == nil {
		return ErrNoAttempt
	}

	return t.attempt.Commit(cb)
}

// TransactionRollbackCallback describes a callback for a completed rollback operation.
type TransactionRollbackCallback func(error)

// Rollback will attempt to rollback the transaction.
func (t *Transaction) Rollback(cb TransactionRollbackCallback) error {
	if t.attempt == nil {
		return ErrNoAttempt
	}

	return t.attempt.Rollback(cb)
}

// HasExpired indicates whether this attempt has expired.
func (t *Transaction) HasExpired() bool {
	if t.attempt == nil {
		return false
	}

	return t.attempt.HasExpired()
}

// CanCommit indicates whether this attempt can still be committed.
func (t *Transaction) CanCommit() bool {
	if t.attempt == nil {
		return false
	}

	return t.attempt.CanCommit()
}

// ShouldRollback indicates if this attempt should be rolled back.
func (t *Transaction) ShouldRollback() bool {
	if t.attempt == nil {
		return false
	}

	return t.attempt.ShouldRollback()
}

// ShouldRetry indicates if this attempt thinks we can retry.
func (t *Transaction) ShouldRetry() bool {
	if t.attempt == nil {
		return false
	}

	return t.attempt.ShouldRetry()
}

// FinalErrorToRaise returns the TransactionErrorReason corresponding to the final state of the transaction.
func (t *Transaction) FinalErrorToRaise() TransactionErrorReason {
	if t.attempt == nil {
		return 0
	}

	return t.attempt.FinalErrorToRaise()
}

func (t *Transaction) TimeRemaining() time.Duration {
	if t.attempt == nil {
		return 0
	}

	return t.attempt.TimeRemaining()
}

// SerializeAttempt will serialize the current transaction attempt, allowing it
// to be resumed later, potentially under a different transactions client.  It
// is no longer safe to use this attempt once this has occurred, a new attempt
// must be started to use this object following this call.
func (t *Transaction) SerializeAttempt(cb func([]byte, error)) error {
	return t.attempt.Serialize(cb)
}

// GetMutations returns a list of all the current mutations that have been performed
// under this transaction.
func (t *Transaction) GetMutations() []TransactionStagedMutation {
	if t.attempt == nil {
		return nil
	}

	return t.attempt.GetMutations()
}

// GetATRLocation returns the ATR location for the current attempt, either by
// identifying where it was placed, or where it will be based on custom atr
// configurations.
func (t *Transaction) GetATRLocation() TransactionATRLocation {
	if t.attempt != nil {
		return t.attempt.GetATRLocation()
	}

	return t.atrLocation
}

// SetATRLocation forces the ATR location for the current attempt to a specific
// location.  Note that this cannot be called if it has already been set.  This
// is currently only safe to call before any mutations have occurred.
func (t *Transaction) SetATRLocation(location TransactionATRLocation) error {
	if t.attempt == nil {
		return errors.New("cannot set ATR location without an active attempt")
	}

	return t.attempt.SetATRLocation(location)
}

// Config returns the configured parameters for this transaction.
// Note that the Expiration time is adjusted based on the time left.
// Note also that after a transaction is resumed, the custom atr location
// may no longer reflect the originally configured value.
func (t *Transaction) Config() TransactionOptions {
	return TransactionOptions{
		CustomATRLocation: t.atrLocation,
		ExpirationTime:    t.TimeRemaining(),
		DurabilityLevel:   t.durabilityLevel,
		KeyValueTimeout:   t.keyValueTimeout,
	}
}

// TransactionUpdateStateOptions are the settings available to UpdateState.
// This function must only be called once the transaction has entered query mode.
// Internal: This should never be used and is not supported.
type TransactionUpdateStateOptions struct {
	ShouldNotCommit   bool
	ShouldNotRollback bool
	ShouldNotRetry    bool
	State             TransactionAttemptState
	Reason            TransactionErrorReason
}

func (tuso TransactionUpdateStateOptions) String() string {
	return fmt.Sprintf("Should not commit: %t, should not rollback: %t, should not retry: %t, state: %s, reason: %s",
		tuso.ShouldNotCommit, tuso.ShouldNotRollback, tuso.ShouldNotRetry, tuso.State, tuso.Reason)
}

// UpdateState will update the internal state of the current attempt.
// Internal: This should never be used and is not supported.
func (t *Transaction) UpdateState(opts TransactionUpdateStateOptions) {
	if t.attempt == nil {
		return
	}

	t.attempt.UpdateState(opts)
}

// Logger returns the logger used by this transaction.
// Uncommitted: This API may change in the future.
func (t *Transaction) Logger() TransactionLogger {
	return t.logger.wrapped
}
