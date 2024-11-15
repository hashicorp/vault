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
	"sync"
	"sync/atomic"
	"time"

	"github.com/couchbase/gocbcore/v10/memd"
)

// TransactionsCleanupRequest represents a complete transaction attempt that requires cleanup.
// Internal: This should never be used and is not supported.
type TransactionsCleanupRequest struct {
	AttemptID         string
	AtrID             []byte
	AtrCollectionName string
	AtrScopeName      string
	AtrBucketName     string
	Inserts           []TransactionsDocRecord
	Replaces          []TransactionsDocRecord
	Removes           []TransactionsDocRecord
	State             TransactionAttemptState
	ForwardCompat     map[string][]TransactionForwardCompatibilityEntry
	DurabilityLevel   TransactionDurabilityLevel
	Age               time.Duration
}

func (cr *TransactionsCleanupRequest) String() string {
	if isLogRedactionLevelFull() || isLogRedactionLevelPartial() {
		return cr.redacted().(string)
	}

	return fmt.Sprintf(
		"bucket: %s, collection: %s, scope: %s, atr: %s, attempt: %s, state: %s, age: %s",
		cr.AtrBucketName,
		cr.AtrCollectionName,
		cr.AtrScopeName,
		cr.AtrID,
		cr.AttemptID,
		cr.State,
		cr.Age,
	)
}

func (cr *TransactionsCleanupRequest) redacted() interface{} {
	return fmt.Sprintf(
		"bucket: %s, collection: %s, scope: %s, atr: %s, attempt: %s, state: %s, age: %s",
		redactMetaData(cr.AtrBucketName),
		redactMetaData(cr.AtrCollectionName),
		redactMetaData(cr.AtrScopeName),
		cr.AtrID,
		cr.AttemptID,
		cr.State,
		cr.Age,
	)
}

// TransactionsDocRecord represents an individual document operation requiring cleanup.
// Internal: This should never be used and is not supported.
type TransactionsDocRecord struct {
	CollectionName string
	ScopeName      string
	BucketName     string
	ID             []byte
}

// TransactionsCleanupAttempt represents the result of running cleanup for a transaction attempt.
// Internal: This should never be used and is not supported.
type TransactionsCleanupAttempt struct {
	Success           bool
	IsReqular         bool
	AttemptID         string
	AtrID             []byte
	AtrCollectionName string
	AtrScopeName      string
	AtrBucketName     string
	Request           *TransactionsCleanupRequest
}

func (ca TransactionsCleanupAttempt) String() string {
	return fmt.Sprintf("bucket: %s, collection: %s, scope: %s, atr: %s, attempt: %s", ca.AtrBucketName, ca.AtrCollectionName,
		ca.AtrScopeName, ca.AtrID, ca.AttemptID)
}

// TransactionsCleaner is responsible for performing cleanup of completed transactions.
// Internal: This should never be used and is not supported.
type TransactionsCleaner interface {
	AddRequest(req *TransactionsCleanupRequest) bool
	PopRequest() *TransactionsCleanupRequest
	ForceCleanupQueue(cb func([]TransactionsCleanupAttempt))
	QueueLength() int32
	CleanupAttempt(atrAgent *Agent, atrOboUser string, req *TransactionsCleanupRequest, regular bool, cb func(attempt TransactionsCleanupAttempt))
	Close()
	GetAndResetResourceUnits() *TransactionResourceUnitResult
}

// NewTransactionsCleaner returns a TransactionsCleaner implementation.
// Internal: This should never be used and is not supported.
func NewTransactionsCleaner(config *TransactionsConfig) TransactionsCleaner {
	return newStdCleaner(config)
}

type noopTransactionsCleaner struct {
}

func (nc *noopTransactionsCleaner) AddRequest(req *TransactionsCleanupRequest) bool {
	return true
}
func (nc *noopTransactionsCleaner) PopRequest() *TransactionsCleanupRequest {
	return nil
}

func (nc *noopTransactionsCleaner) ForceCleanupQueue(cb func([]TransactionsCleanupAttempt)) {
	cb([]TransactionsCleanupAttempt{})
}

func (nc *noopTransactionsCleaner) QueueLength() int32 {
	return 0
}

func (nc *noopTransactionsCleaner) GetAndResetResourceUnits() *TransactionResourceUnitResult {
	return nil
}

func (nc *noopTransactionsCleaner) CleanupAttempt(atrAgent *Agent, atrOboUser string, req *TransactionsCleanupRequest, regular bool, cb func(attempt TransactionsCleanupAttempt)) {
	cb(TransactionsCleanupAttempt{})
}

func (nc *noopTransactionsCleaner) Close() {}

type stdTransactionsCleaner struct {
	hooks               TransactionCleanUpHooks
	qSize               uint32
	q                   chan *TransactionsCleanupRequest
	stop                chan struct{}
	bucketAgentProvider TransactionsBucketAgentProviderFn
	keyValueTimeout     time.Duration
	durabilityLevel     TransactionDurabilityLevel

	numResourceUnitOps uint32
	readUnits          uint32
	writeUnits         uint32
}

func newStdCleaner(config *TransactionsConfig) *stdTransactionsCleaner {
	return &stdTransactionsCleaner{
		hooks:               config.Internal.CleanUpHooks,
		qSize:               config.CleanupQueueSize,
		stop:                make(chan struct{}),
		bucketAgentProvider: config.BucketAgentProvider,
		q:                   make(chan *TransactionsCleanupRequest, config.CleanupQueueSize),
		keyValueTimeout:     config.KeyValueTimeout,
		durabilityLevel:     config.DurabilityLevel,
	}
}

func startCleanupThread(config *TransactionsConfig) *stdTransactionsCleaner {
	cleaner := newStdCleaner(config)

	// No point in running this if we can't get agents.
	if config.BucketAgentProvider != nil {
		go cleaner.processQ()
	}

	return cleaner
}

func (c *stdTransactionsCleaner) AddRequest(req *TransactionsCleanupRequest) bool {
	select {
	case c.q <- req:
		// success!
	default:
		logDebugf("Not queueing request for: %s, limit size reached",
			req.String())
	}

	return true
}

func (c *stdTransactionsCleaner) PopRequest() *TransactionsCleanupRequest {
	select {
	case req := <-c.q:
		return req
	default:
		return nil
	}
}

func (c *stdTransactionsCleaner) stealAllRequests() []*TransactionsCleanupRequest {
	reqs := make([]*TransactionsCleanupRequest, 0, len(c.q))
	for {
		select {
		case req := <-c.q:
			reqs = append(reqs, req)
		default:
			return reqs
		}
	}
}

func (c *stdTransactionsCleaner) updateResourceUnits(units *ResourceUnitResult) {
	if units == nil {
		return
	}

	atomic.AddUint32(&c.numResourceUnitOps, 1)
	atomic.AddUint32(&c.readUnits, uint32(units.ReadUnits))
	atomic.AddUint32(&c.writeUnits, uint32(units.WriteUnits))
}

func (c *stdTransactionsCleaner) updateResourceUnitsError(err error) {
	if err == nil {
		return
	}

	var kerr *KeyValueError
	if errors.As(err, &kerr) {
		c.updateResourceUnits(kerr.Internal.ResourceUnits)
	}
}

func (c *stdTransactionsCleaner) GetAndResetResourceUnits() *TransactionResourceUnitResult {
	numOps := atomic.SwapUint32(&c.numResourceUnitOps, 0)
	if numOps == 0 {
		return nil
	}

	return &TransactionResourceUnitResult{
		NumOps:     numOps,
		ReadUnits:  atomic.SwapUint32(&c.readUnits, 0),
		WriteUnits: atomic.SwapUint32(&c.writeUnits, 0),
	}
}

// Used only for tests
func (c *stdTransactionsCleaner) ForceCleanupQueue(cb func([]TransactionsCleanupAttempt)) {
	reqs := c.stealAllRequests()
	if len(reqs) == 0 {
		cb(nil)
		return
	}

	results := make([]TransactionsCleanupAttempt, 0, len(reqs))
	var l sync.Mutex
	handler := func(attempt TransactionsCleanupAttempt) {
		l.Lock()
		defer l.Unlock()
		results = append(results, attempt)
		if len(results) == len(reqs) {
			cb(results)
		}
	}

	for _, req := range reqs {
		agent, oboUser, err := c.bucketAgentProvider(req.AtrBucketName)
		if err != nil {
			handler(TransactionsCleanupAttempt{
				Success:           false,
				IsReqular:         false,
				AttemptID:         req.AttemptID,
				AtrID:             req.AtrID,
				AtrCollectionName: req.AtrCollectionName,
				AtrScopeName:      req.AtrScopeName,
				AtrBucketName:     req.AtrBucketName,
				Request:           req,
			})
			continue
		}

		c.CleanupAttempt(agent, oboUser, req, true, func(attempt TransactionsCleanupAttempt) {
			handler(attempt)
		})
	}
}

// Used only for tests
func (c *stdTransactionsCleaner) QueueLength() int32 {
	return int32(len(c.q))
}

// Used only for tests
func (c *stdTransactionsCleaner) Close() {
	close(c.stop)
}

func (c *stdTransactionsCleaner) processQ() {
	logDebugf("Starting cleanup for %p", c)
	for {
		select {
		case req := <-c.q:
			agent, oboUser, err := c.bucketAgentProvider(req.AtrBucketName)
			if err != nil {
				logDebugf("Failed to get agent for request: %s, err: %v", req.String(), err)
				return
			}

			logSchedf("Running cleanup for request: %s", req.String())
			waitCh := make(chan struct{}, 1)
			c.CleanupAttempt(agent, oboUser, req, true, func(attempt TransactionsCleanupAttempt) {
				if !attempt.Success {
					logDebugf("Cleanup attempt failed for entry: %s",
						attempt.String())
				}

				waitCh <- struct{}{}
			})
			<-waitCh

		case <-c.stop:
			return
		}
	}
}

func (c *stdTransactionsCleaner) checkForwardCompatability(
	stage forwardCompatStage,
	fc map[string][]TransactionForwardCompatibilityEntry,
	cb func(error),
) {
	isCompat, _, _, err := checkForwardCompatability(stage, fc)
	if err != nil {
		cb(err)
		return
	}

	if !isCompat {
		cb(ErrForwardCompatibilityFailure)
		return
	}

	cb(nil)
}

func (c *stdTransactionsCleaner) CleanupAttempt(atrAgent *Agent, atrOboUser string, req *TransactionsCleanupRequest, regular bool, cb func(attempt TransactionsCleanupAttempt)) {
	beforeCb := func(stage string, attempt TransactionsCleanupAttempt) {
		if attempt.Success {
			cb(attempt)
			return
		}

		logWarnf("Cleanup attempt %v with %p failed at %s check", req, c, stage)

		if req.Age > 2*time.Hour {
			logWarnf("Cleanup request is %s old which could indicate a serious error - please raise with support.", req.Age)
		}

		cb(attempt)
	}
	logSchedf("Cleaning up attempt %s with %p", req.AttemptID, c)
	c.checkForwardCompatability(forwardCompatStageGetsCleanupEntry, req.ForwardCompat, func(err error) {
		if err != nil {
			beforeCb("forward compatability", TransactionsCleanupAttempt{
				Success:           false,
				IsReqular:         regular,
				AttemptID:         req.AttemptID,
				AtrID:             req.AtrID,
				AtrCollectionName: req.AtrCollectionName,
				AtrScopeName:      req.AtrScopeName,
				AtrBucketName:     req.AtrBucketName,
				Request:           req,
			})
			return
		}

		c.cleanupDocs(req, func(err error) {
			if err != nil {
				beforeCb("cleanup docs", TransactionsCleanupAttempt{
					Success:           false,
					IsReqular:         regular,
					AttemptID:         req.AttemptID,
					AtrID:             req.AtrID,
					AtrCollectionName: req.AtrCollectionName,
					AtrScopeName:      req.AtrScopeName,
					AtrBucketName:     req.AtrBucketName,
					Request:           req,
				})
				return
			}

			c.cleanupATR(atrAgent, atrOboUser, req, func(err error) {
				success := true
				if err != nil {
					success = false
				}

				beforeCb("cleanup atr", TransactionsCleanupAttempt{
					Success:           success,
					IsReqular:         regular,
					AttemptID:         req.AttemptID,
					AtrID:             req.AtrID,
					AtrCollectionName: req.AtrCollectionName,
					AtrScopeName:      req.AtrScopeName,
					AtrBucketName:     req.AtrBucketName,
					Request:           req,
				})
			})
		})
	})
}

func (c *stdTransactionsCleaner) cleanupATR(agent *Agent, oboUser string, req *TransactionsCleanupRequest, cb func(error)) {
	c.hooks.BeforeATRRemove(req.AtrID, func(err error) {
		if err != nil {
			if errors.Is(err, ErrPathNotFound) {
				cb(nil)
				return
			}
			cb(err)
			return
		}

		var specs []SubDocOp
		if req.State == TransactionAttemptStatePending {
			specs = append(specs, SubDocOp{
				Op:    memd.SubDocOpDictAdd,
				Value: []byte{110, 117, 108, 108},
				Path:  "attempts." + req.AttemptID + ".p",
				Flags: memd.SubdocFlagXattrPath,
			})
		}

		specs = append(specs, SubDocOp{
			Op:    memd.SubDocOpDelete,
			Path:  "attempts." + req.AttemptID,
			Flags: memd.SubdocFlagXattrPath,
		})

		if req.DurabilityLevel == TransactionDurabilityLevelUnknown {
			req.DurabilityLevel = c.durabilityLevel
		}
		deadline, duraTimeout := transactionsMutationTimeouts(c.keyValueTimeout, req.DurabilityLevel)

		_, err = agent.MutateIn(MutateInOptions{
			Key:                    req.AtrID,
			ScopeName:              req.AtrScopeName,
			CollectionName:         req.AtrCollectionName,
			Ops:                    specs,
			Deadline:               deadline,
			DurabilityLevel:        transactionsDurabilityLevelToMemd(req.DurabilityLevel),
			DurabilityLevelTimeout: duraTimeout,
			User:                   oboUser,
		}, func(result *MutateInResult, err error) {
			if err != nil {
				c.updateResourceUnitsError(err)
				if errors.Is(err, ErrPathNotFound) {
					cb(nil)
					return
				}

				logDebugf("Failed to cleanup ATR for request: %s, err: %v", req.String(), err)
				cb(err)
				return
			}

			c.updateResourceUnits(result.Internal.ResourceUnits)
			cb(nil)
		})
		if err != nil {
			cb(err)
			return
		}

	})
}

func (c *stdTransactionsCleaner) cleanupDocs(req *TransactionsCleanupRequest, cb func(error)) {
	var memdDuraLevel memd.DurabilityLevel
	if req.DurabilityLevel > TransactionDurabilityLevelUnknown {
		// We want to ensure that we don't panic here, if the durability level is unknown then we'll just not set
		// a durability level.
		memdDuraLevel = transactionsDurabilityLevelToMemd(req.DurabilityLevel)
	}
	deadline, duraTimeout := transactionsMutationTimeouts(c.keyValueTimeout, req.DurabilityLevel)

	switch req.State {
	case TransactionAttemptStateCommitted:

		waitCh := make(chan error, 1)
		c.commitInsRepDocs(req.AttemptID, req.Inserts, deadline, memdDuraLevel, duraTimeout, func(err error) {
			waitCh <- err
		})
		err := <-waitCh
		if err != nil {
			cb(err)
			return
		}

		waitCh = make(chan error, 1)
		c.commitInsRepDocs(req.AttemptID, req.Replaces, deadline, memdDuraLevel, duraTimeout, func(err error) {
			waitCh <- err
		})
		err = <-waitCh
		if err != nil {
			cb(err)
			return
		}

		waitCh = make(chan error, 1)
		c.commitRemDocs(req.AttemptID, req.Removes, deadline, memdDuraLevel, duraTimeout, func(err error) {
			waitCh <- err
		})
		err = <-waitCh
		if err != nil {
			cb(err)
			return
		}

		cb(nil)
	case TransactionAttemptStateAborted:
		waitCh := make(chan error, 1)
		c.rollbackInsDocs(req.AttemptID, req.Inserts, deadline, memdDuraLevel, duraTimeout, func(err error) {
			waitCh <- err
		})
		err := <-waitCh
		if err != nil {
			cb(err)
			return
		}

		waitCh = make(chan error, 1)
		c.rollbackRepRemDocs(req.AttemptID, req.Replaces, deadline, memdDuraLevel, duraTimeout, func(err error) {
			waitCh <- err
		})
		err = <-waitCh
		if err != nil {
			cb(err)
			return
		}

		waitCh = make(chan error, 1)
		c.rollbackRepRemDocs(req.AttemptID, req.Removes, deadline, memdDuraLevel, duraTimeout, func(err error) {
			waitCh <- err
		})
		err = <-waitCh
		if err != nil {
			cb(err)
			return
		}

		cb(nil)
	case TransactionAttemptStatePending:
		cb(nil)
	case TransactionAttemptStateCompleted:
		cb(nil)
	case TransactionAttemptStateRolledBack:
		cb(nil)
	case TransactionAttemptStateNothingWritten:
		cb(nil)
	default:
		cb(nil)
	}
}

func (c *stdTransactionsCleaner) rollbackRepRemDocs(attemptID string, docs []TransactionsDocRecord, deadline time.Time, durability memd.DurabilityLevel,
	duraTimeout time.Duration, cb func(err error)) {

	for _, doc := range docs {
		waitCh := make(chan error, 1)

		agent, oboUser, err := c.bucketAgentProvider(doc.BucketName)
		if err != nil {
			cb(err)
			return
		}

		c.perDoc(false, attemptID, doc, agent, oboUser, func(getRes *transactionGetDoc, err error) {
			if err != nil {
				waitCh <- err
				return
			}

			if getRes == nil {
				// This violates implicit contract idioms but needs must.
				waitCh <- nil
				return
			}

			c.hooks.BeforeRemoveLinks(doc.ID, func(err error) {
				if err != nil {
					waitCh <- err
					return
				}

				_, err = agent.MutateIn(MutateInOptions{
					Key:            doc.ID,
					ScopeName:      doc.ScopeName,
					CollectionName: doc.CollectionName,
					Cas:            getRes.Cas,
					Ops: []SubDocOp{
						{
							Op:    memd.SubDocOpDelete,
							Path:  "txn",
							Flags: memd.SubdocFlagXattrPath,
						},
					},
					Flags:                  memd.SubdocDocFlagAccessDeleted,
					Deadline:               deadline,
					DurabilityLevel:        durability,
					DurabilityLevelTimeout: duraTimeout,
					User:                   oboUser,
				}, func(result *MutateInResult, err error) {
					if err != nil {
						c.updateResourceUnitsError(err)
						logDebugf("Failed to rollback for bucket: %s, collection: %s, scope: %s, id: %s, err: %v",
							doc.BucketName, doc.CollectionName, doc.ScopeName, doc.ID, err)
						waitCh <- err
						return
					}

					c.updateResourceUnits(result.Internal.ResourceUnits)
					waitCh <- nil
				})
				if err != nil {
					waitCh <- err
					return
				}
			})
		})

		err = <-waitCh

		if err != nil {
			cb(err)
			return
		}
	}

	cb(nil)
}

func (c *stdTransactionsCleaner) rollbackInsDocs(attemptID string, docs []TransactionsDocRecord, deadline time.Time, durability memd.DurabilityLevel,
	duraTimeout time.Duration, cb func(err error)) {

	for _, doc := range docs {
		waitCh := make(chan error, 1)

		agent, oboUser, err := c.bucketAgentProvider(doc.BucketName)
		if err != nil {
			cb(err)
			return
		}

		c.perDoc(false, attemptID, doc, agent, oboUser, func(getRes *transactionGetDoc, err error) {
			if err != nil {
				waitCh <- err
				return
			}

			if getRes == nil {
				// This violates implicit contract idioms but needs must.
				waitCh <- nil
				return
			}

			c.hooks.BeforeRemoveDoc(doc.ID, func(err error) {
				if err != nil {
					waitCh <- err
					return
				}

				if getRes.Deleted {
					_, err := agent.MutateIn(MutateInOptions{
						Key:            doc.ID,
						ScopeName:      doc.ScopeName,
						CollectionName: doc.CollectionName,
						Cas:            getRes.Cas,
						Ops: []SubDocOp{
							{
								Op:    memd.SubDocOpDelete,
								Path:  "txn",
								Flags: memd.SubdocFlagXattrPath,
							},
						},
						Flags:                  memd.SubdocDocFlagAccessDeleted,
						Deadline:               deadline,
						DurabilityLevel:        durability,
						DurabilityLevelTimeout: duraTimeout,
						User:                   oboUser,
					}, func(result *MutateInResult, err error) {
						if err != nil {
							c.updateResourceUnitsError(err)
							logDebugf("Failed to rollback for bucket: %s, collection: %s, scope: %s, id: %s, err: %v",
								doc.BucketName, doc.CollectionName, doc.ScopeName, doc.ID, err)
							waitCh <- err
							return
						}

						c.updateResourceUnits(result.Internal.ResourceUnits)
						waitCh <- nil
					})
					if err != nil {
						waitCh <- err
						return
					}
				} else {
					_, err := agent.Delete(DeleteOptions{
						Key:                    doc.ID,
						ScopeName:              doc.ScopeName,
						CollectionName:         doc.CollectionName,
						Cas:                    getRes.Cas,
						Deadline:               deadline,
						DurabilityLevel:        durability,
						DurabilityLevelTimeout: duraTimeout,
						User:                   oboUser,
					}, func(result *DeleteResult, err error) {
						if err != nil {
							c.updateResourceUnitsError(err)
							logDebugf("Failed to rollback for bucket: %s, collection: %s, scope: %s, id: %s, err: %v",
								doc.BucketName, doc.CollectionName, doc.ScopeName, doc.ID, err)
							waitCh <- err
							return
						}

						c.updateResourceUnits(result.Internal.ResourceUnits)
						waitCh <- nil
					})
					if err != nil {
						waitCh <- err
						return
					}
				}
			})
		})

		err = <-waitCh
		if err != nil {
			cb(err)
			return
		}
	}

	cb(nil)
}

func (c *stdTransactionsCleaner) commitRemDocs(attemptID string, docs []TransactionsDocRecord, deadline time.Time, durability memd.DurabilityLevel,
	duraTimeout time.Duration, cb func(err error)) {

	for _, doc := range docs {
		waitCh := make(chan error, 1)

		agent, oboUser, err := c.bucketAgentProvider(doc.BucketName)
		if err != nil {
			cb(err)
			return
		}

		c.perDoc(true, attemptID, doc, agent, oboUser, func(getRes *transactionGetDoc, err error) {
			if err != nil {
				waitCh <- err
				return
			}

			if getRes == nil {
				// This violates implicit contract idioms but needs must.
				waitCh <- nil
				return
			}

			c.hooks.BeforeRemoveDocStagedForRemoval(doc.ID, func(err error) {
				if err != nil {
					waitCh <- err
					return
				}

				if getRes.TxnMeta.Operation.Type != jsonMutationRemove {
					waitCh <- nil
					return
				}

				_, err = agent.Delete(DeleteOptions{
					Key:                    doc.ID,
					ScopeName:              doc.ScopeName,
					CollectionName:         doc.CollectionName,
					Cas:                    getRes.Cas,
					Deadline:               deadline,
					DurabilityLevel:        durability,
					DurabilityLevelTimeout: duraTimeout,
					User:                   oboUser,
				}, func(result *DeleteResult, err error) {
					if err != nil {
						c.updateResourceUnitsError(err)
						logDebugf("Failed to commit for bucket: %s, collection: %s, scope: %s, id: %s, err: %v",
							doc.BucketName, doc.CollectionName, doc.ScopeName, doc.ID, err)
						waitCh <- err
						return
					}

					c.updateResourceUnits(result.Internal.ResourceUnits)
					waitCh <- nil
				})
				if err != nil {
					waitCh <- err
					return
				}
			})
		})
		err = <-waitCh
		if err != nil {
			cb(err)
			return
		}
	}

	cb(nil)
}

func (c *stdTransactionsCleaner) commitInsRepDocs(attemptID string, docs []TransactionsDocRecord, deadline time.Time, durability memd.DurabilityLevel,
	duraTimeout time.Duration, cb func(err error)) {

	for _, doc := range docs {
		waitCh := make(chan error, 1)

		agent, oboUser, err := c.bucketAgentProvider(doc.BucketName)
		if err != nil {
			cb(err)
			return
		}

		c.perDoc(true, attemptID, doc, agent, oboUser, func(getRes *transactionGetDoc, err error) {
			if err != nil {
				waitCh <- err
				return
			}

			if getRes == nil {
				// This violates implicit contract idioms but needs must.
				waitCh <- nil
				return
			}

			c.hooks.BeforeCommitDoc(doc.ID, func(err error) {
				if err != nil {
					waitCh <- err
					return
				}

				if getRes.Deleted {
					_, err := agent.Set(SetOptions{
						Value:                  getRes.Body,
						Key:                    doc.ID,
						ScopeName:              doc.ScopeName,
						CollectionName:         doc.CollectionName,
						Deadline:               deadline,
						DurabilityLevel:        durability,
						DurabilityLevelTimeout: duraTimeout,
						User:                   oboUser,
					}, func(result *StoreResult, err error) {
						if err != nil {
							c.updateResourceUnitsError(err)
							logDebugf("Failed to commit for bucket: %s, collection: %s, scope: %s, id: %s, err: %v",
								doc.BucketName, doc.CollectionName, doc.ScopeName, doc.ID, err)
							waitCh <- err
							return
						}

						c.updateResourceUnits(result.Internal.ResourceUnits)
						waitCh <- nil
					})
					if err != nil {
						waitCh <- err
						return
					}
				} else {
					_, err := agent.MutateIn(MutateInOptions{
						Key:            doc.ID,
						ScopeName:      doc.ScopeName,
						CollectionName: doc.CollectionName,
						Cas:            getRes.Cas,
						Ops: []SubDocOp{
							{
								Op:    memd.SubDocOpDelete,
								Path:  "txn",
								Flags: memd.SubdocFlagXattrPath,
							},
							{
								Op:    memd.SubDocOpSetDoc,
								Path:  "",
								Value: getRes.Body,
							},
						},
						Deadline:               deadline,
						DurabilityLevel:        durability,
						DurabilityLevelTimeout: duraTimeout,
						User:                   oboUser,
					}, func(result *MutateInResult, err error) {
						if err != nil {
							c.updateResourceUnitsError(err)
							logDebugf("Failed to commit for bucket: %s, collection: %s, scope: %s, id: %s, err: %v",
								doc.BucketName, doc.CollectionName, doc.ScopeName, doc.ID, err)
							waitCh <- err
							return
						}

						c.updateResourceUnits(result.Internal.ResourceUnits)
						waitCh <- nil
					})
					if err != nil {
						waitCh <- err
						return
					}
				}
			})
		})
		err = <-waitCh
		if err != nil {
			cb(err)
			return
		}
	}

	cb(nil)
}

func (c *stdTransactionsCleaner) perDoc(crc32MatchStaging bool, attemptID string, dr TransactionsDocRecord, agent *Agent, oboUser string,
	cb func(getRes *transactionGetDoc, err error)) {
	c.hooks.BeforeDocGet(dr.ID, func(err error) {
		if err != nil {
			cb(nil, err)
			return
		}

		var deadline time.Time
		if c.keyValueTimeout > 0 {
			deadline = time.Now().Add(c.keyValueTimeout)
		}

		_, err = agent.LookupIn(LookupInOptions{
			ScopeName:      dr.ScopeName,
			CollectionName: dr.CollectionName,
			Key:            dr.ID,
			Ops: []SubDocOp{
				{
					Op:    memd.SubDocOpGet,
					Path:  "$document",
					Flags: memd.SubdocFlagXattrPath,
				},
				{
					Op:    memd.SubDocOpGet,
					Path:  "txn",
					Flags: memd.SubdocFlagXattrPath,
				},
			},
			Deadline: deadline,
			Flags:    memd.SubdocDocFlagAccessDeleted,
			User:     oboUser,
		}, func(result *LookupInResult, err error) {
			if err != nil {
				c.updateResourceUnitsError(err)
				if errors.Is(err, ErrDocumentNotFound) {
					// We can consider this success.
					cb(nil, nil)
					return
				}

				logDebugf("Failed to lookup doc for bucket: %s, collection: %s, scope: %s, id: %s, err: %v",
					dr.BucketName, dr.CollectionName, dr.ScopeName, dr.ID, err)
				cb(nil, err)
				return
			}

			c.updateResourceUnits(result.Internal.ResourceUnits)

			if result.Ops[0].Err != nil {
				// This is not so good.
				cb(nil, result.Ops[0].Err)
				return
			}

			if result.Ops[1].Err != nil {
				// Txn probably committed so this is success.
				cb(nil, nil)
				return
			}

			var txnMetaVal *jsonTxnXattr
			if err := json.Unmarshal(result.Ops[1].Value, &txnMetaVal); err != nil {
				cb(nil, err)
				return
			}

			if attemptID != txnMetaVal.ID.Attempt {
				// Document involved in another txn, was probably committed, this is success.
				cb(nil, nil)
				return
			}

			var meta *transactionDocMeta
			if err := json.Unmarshal(result.Ops[0].Value, &meta); err != nil {
				cb(nil, err)
				return
			}
			if crc32MatchStaging {
				if meta.CRC32 != txnMetaVal.Operation.CRC32 {
					// This document is a part of this txn but its body has changed, we'll continue as success.
					cb(nil, nil)
					return
				}
			}

			cb(&transactionGetDoc{
				Body:    txnMetaVal.Operation.Staged,
				DocMeta: meta,
				Cas:     result.Cas,
				Deleted: result.Internal.IsDeleted,
				TxnMeta: txnMetaVal,
			}, nil)
		})
		if err != nil {
			cb(nil, err)
		}
	})
}
