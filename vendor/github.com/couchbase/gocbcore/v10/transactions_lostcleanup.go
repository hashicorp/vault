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
	"math"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/couchbase/gocbcore/v10/memd"
	"github.com/google/uuid"
)

var clientRecordKey = []byte("_txn:client-record")

type jsonClientRecord struct {
	HeartbeatMS string `json:"heartbeat_ms,omitempty"`
	ExpiresMS   int    `json:"expires_ms,omitempty"`
	NumATRs     int    `json:"num_atrs,omitempty"`
}

type jsonClientOverride struct {
	Enabled      bool  `json:"enabled,omitempty"`
	ExpiresNanos int64 `json:"expires,omitempty"`
}

type jsonClientRecords struct {
	Clients  map[string]jsonClientRecord `json:"clients"`
	Override *jsonClientOverride         `json:"override,omitempty"`
}

type jsonHLC struct {
	NowSecs string `json:"now"`
}

// TransactionClientRecordDetails is the result of processing a client record.
// Internal: This should never be used and is not supported.
type TransactionClientRecordDetails struct {
	NumActiveClients     int
	IndexOfThisClient    int
	ClientIsNew          bool
	ExpiredClientIDs     []string
	NumExistingClients   int
	NumExpiredClients    int
	OverrideEnabled      bool
	OverrideActive       bool
	OverrideExpiresCas   int64
	CasNowNanos          int64
	AtrsHandledByClient  []string
	CheckAtrEveryNMillis int
	ClientUUID           string
}

// TransactionProcessATRStats is the stats recorded when running a ProcessATR request.
// Internal: This should never be used and is not supported.
type TransactionProcessATRStats struct {
	NumEntries        int
	NumEntriesExpired int
}

// LostTransactionCleaner is responsible for cleaning up lost transactions.
// Internal: This should never be used and is not supported.
type LostTransactionCleaner interface {
	ProcessClient(agent *Agent, oboUser string, collection, scope, uuid string, cb func(*TransactionClientRecordDetails, error))
	ProcessATR(agent *Agent, oboUser string, collection, scope, atrID string, cb func([]TransactionsCleanupAttempt, TransactionProcessATRStats, error))
	RemoveClientFromAllLocations(uuid string) error
	Close()
	GetAndResetResourceUnits() *TransactionResourceUnitResult
}

type lostTransactionCleaner interface {
	AddATRLocation(location TransactionLostATRLocation)
	ATRLocations() []TransactionLostATRLocation
	Close()
	GetAndResetResourceUnits() *TransactionResourceUnitResult
}

type noopLostTransactionCleaner struct {
}

func (ltc *noopLostTransactionCleaner) AddATRLocation(location TransactionLostATRLocation) {
}

func (ltc *noopLostTransactionCleaner) ATRLocations() []TransactionLostATRLocation {
	return nil
}

func (ltc *noopLostTransactionCleaner) Close() {
}

func (ltc *noopLostTransactionCleaner) GetAndResetResourceUnits() *TransactionResourceUnitResult {
	return nil
}

type stdLostTransactionCleaner struct {
	uuid                string
	cleanupHooks        TransactionCleanUpHooks
	clientRecordHooks   TransactionClientRecordHooks
	numAtrs             int
	cleanupWindow       time.Duration
	cleaner             TransactionsCleaner
	keyValueTimeout     time.Duration
	bucketAgentProvider TransactionsBucketAgentProviderFn
	locations           map[TransactionLostATRLocation]chan struct{}
	locationsLock       sync.Mutex
	newLocationCh       chan lostATRLocationWithShutdown
	stop                chan struct{}
	atrLocationFinder   TransactionsLostCleanupATRLocationProviderFn
	processWaitGroup    sync.WaitGroup

	numResourceUnitOps uint32
	readUnits          uint32
	writeUnits         uint32
}

type lostATRLocationWithShutdown struct {
	location TransactionLostATRLocation
	shutdown chan struct{}
}

// NewLostTransactionCleaner returns new lost transaction cleaner.
// Internal: This should never be used and is not supported.
func NewLostTransactionCleaner(config *TransactionsConfig) LostTransactionCleaner {
	return newStdLostTransactionCleaner(config)
}

func newStdLostTransactionCleaner(config *TransactionsConfig) *stdLostTransactionCleaner {
	return &stdLostTransactionCleaner{
		uuid:                uuid.New().String(),
		numAtrs:             config.Internal.NumATRs,
		cleanupWindow:       config.CleanupWindow,
		cleanupHooks:        config.Internal.CleanUpHooks,
		clientRecordHooks:   config.Internal.ClientRecordHooks,
		cleaner:             NewTransactionsCleaner(config),
		keyValueTimeout:     config.KeyValueTimeout,
		bucketAgentProvider: config.BucketAgentProvider,
		locations:           make(map[TransactionLostATRLocation]chan struct{}),
		newLocationCh:       make(chan lostATRLocationWithShutdown, 20), // Buffer of 20 should be plenty
		stop:                make(chan struct{}),
		atrLocationFinder:   config.LostCleanupATRLocationProvider,
	}
}

func startLostTransactionCleaner(config *TransactionsConfig) *stdLostTransactionCleaner {
	t := newStdLostTransactionCleaner(config)

	if config.BucketAgentProvider != nil {
		go t.start()
	}

	return t
}

func (ltc *stdLostTransactionCleaner) start() {
	logDebugf("Lost transactions %s starting", ltc.uuid)
	ltc.fetchExtraCleanupLocations()

	for {
		select {
		case <-ltc.stop:
			return
		case location := <-ltc.newLocationCh:
			logDebugf("Starting new location handler for %s, location %s", ltc.uuid, location.location)
			agent, oboUser, err := ltc.bucketAgentProvider(location.location.BucketName)
			if err != nil {
				logDebugf("Failed to fetch agent for %s, location: %s:, err: %v",
					ltc.uuid, location.location, err)
				// We should probably do something here...
				continue
			}
			go ltc.perLocation(agent, oboUser, location)
		}
	}
}

func (ltc *stdLostTransactionCleaner) GetAndResetResourceUnits() *TransactionResourceUnitResult {
	baseUnits := ltc.cleaner.GetAndResetResourceUnits()
	numOps := atomic.SwapUint32(&ltc.numResourceUnitOps, 0)
	if numOps == 0 && baseUnits == nil {
		return nil
	}

	readUnits := atomic.SwapUint32(&ltc.readUnits, 0)
	writeUnits := atomic.SwapUint32(&ltc.writeUnits, 0)

	if baseUnits == nil {
		return &TransactionResourceUnitResult{
			NumOps:     numOps,
			ReadUnits:  readUnits,
			WriteUnits: writeUnits,
		}
	} else if numOps == 0 {
		return baseUnits
	}
	return &TransactionResourceUnitResult{
		NumOps:     numOps + baseUnits.NumOps,
		ReadUnits:  readUnits + baseUnits.ReadUnits,
		WriteUnits: writeUnits + baseUnits.WriteUnits,
	}
}

func (ltc *stdLostTransactionCleaner) ATRLocations() []TransactionLostATRLocation {
	ltc.locationsLock.Lock()
	defer ltc.locationsLock.Unlock()
	if ltc.locations == nil {
		return nil
	}
	var locations []TransactionLostATRLocation
	for location := range ltc.locations {
		locations = append(locations, location)
	}
	return locations
}

func (ltc *stdLostTransactionCleaner) AddATRLocation(location TransactionLostATRLocation) {
	ltc.locationsLock.Lock()
	if ltc.locations == nil {
		ltc.locationsLock.Unlock()
		return
	}
	if _, ok := ltc.locations[location]; ok {
		ltc.locationsLock.Unlock()
		return
	}
	ch := make(chan struct{})
	ltc.locations[location] = ch
	ltc.locationsLock.Unlock()
	logDebugf("Adding location %s to lost cleanup for %s", location, ltc.uuid)
	ltc.newLocationCh <- lostATRLocationWithShutdown{
		location: location,
		shutdown: ch,
	}
}

func (ltc *stdLostTransactionCleaner) Close() {
	logDebugf("Lost transactions %s stopping", ltc.uuid)
	close(ltc.stop)
	err := ltc.RemoveClientFromAllLocations(ltc.uuid)
	if err != nil {
		logDebugf("Failed to remove client from all buckets: %v", err)
	}
}

func (ltc *stdLostTransactionCleaner) RemoveClientFromAllLocations(uuid string) error {
	ltc.locationsLock.Lock()
	if ltc.locations == nil {
		ltc.locationsLock.Unlock()
		return nil
	}
	locations := ltc.locations
	ltc.locations = nil
	ltc.locationsLock.Unlock()
	logDebugf("Removing %s from all locations", ltc.uuid)
	if ltc.atrLocationFinder != nil {
		bs, err := ltc.atrLocationFinder()
		if err != nil {
			logDebugf("Failed to get atr locations for %s: %v", ltc.uuid, err)
			return err
		}

		for _, b := range bs {
			if _, ok := locations[b]; !ok {
				locations[b] = make(chan struct{})
			}
		}
	}

	return ltc.removeClient(uuid, locations)
}

func (ltc *stdLostTransactionCleaner) updateResourceUnits(units *ResourceUnitResult) {
	if units == nil {
		return
	}

	atomic.AddUint32(&ltc.numResourceUnitOps, 1)
	atomic.AddUint32(&ltc.readUnits, uint32(units.ReadUnits))
	atomic.AddUint32(&ltc.writeUnits, uint32(units.WriteUnits))
}

func (ltc *stdLostTransactionCleaner) updateResourceUnitsError(err error) {
	if err == nil {
		return
	}

	var kerr *KeyValueError
	if errors.As(err, &kerr) {
		ltc.updateResourceUnits(kerr.Internal.ResourceUnits)
	}
}

func (ltc *stdLostTransactionCleaner) removeClient(uuid string, locations map[TransactionLostATRLocation]chan struct{}) error {
	ltc.processWaitGroup.Wait()

	var err error
	var wg sync.WaitGroup
	for l := range locations {
		wg.Add(1)
		func(location TransactionLostATRLocation) {
			// There's a possible race between here and the client record being updated/created.
			// If that happens then it'll be expired and removed by another client anyway
			deadline := time.Now().Add(500 * time.Millisecond)

			ltc.unregisterClientRecord(location, uuid, deadline, func(unregErr error) {
				if unregErr != nil {
					logDebugf("Failed to unregister %s from cleanup record on from location %s: %v", uuid, location, unregErr)
					err = unregErr
				}
				logInfof("Unregistered %s from cleanup record for location %v", uuid, location)
				wg.Done()
			})
		}(l)
	}
	wg.Wait()

	return err
}

func (ltc *stdLostTransactionCleaner) unregisterClientRecord(location TransactionLostATRLocation, uuid string, deadline time.Time, cb func(error)) {
	logDebugf("Unregistering client %s for %s, location = %s", uuid, ltc.uuid, location)
	agent, oboUser, err := ltc.bucketAgentProvider(location.BucketName)
	if err != nil {
		logDebugf("Failed to get agent for %s, location = %s, client = %s: %v", ltc.uuid, location, uuid, err)
		select {
		case <-time.After(time.Until(deadline)):
			logDebugf("Timed out fetching agent for %s, location = %s, client = %s", ltc.uuid, location, uuid)
			cb(ErrTimeout)
			return
		case <-time.After(10 * time.Millisecond):
		}
		ltc.unregisterClientRecord(location, uuid, deadline, cb)
		return
	}

	ltc.clientRecordHooks.BeforeRemoveClient(func(err error) {
		if err != nil {
			if errors.Is(err, ErrDocumentNotFound) || errors.Is(err, ErrPathNotFound) {
				cb(nil)
				return
			}

			select {
			case <-time.After(time.Until(deadline)):
				cb(ErrTimeout)
				return
			case <-time.After(10 * time.Millisecond):
			}
			ltc.unregisterClientRecord(location, uuid, deadline, cb)
			return
		}

		var opDeadline time.Time
		if ltc.keyValueTimeout > 0 {
			opDeadline = time.Now().Add(ltc.keyValueTimeout)
		}

		_, err = agent.MutateIn(MutateInOptions{
			Key: clientRecordKey,
			Ops: []SubDocOp{
				{
					Op:    memd.SubDocOpDelete,
					Flags: memd.SubdocFlagXattrPath,
					Path:  "records.clients." + uuid,
				},
			},
			Deadline:       opDeadline,
			CollectionName: location.CollectionName,
			ScopeName:      location.ScopeName,
			User:           oboUser,
		}, func(result *MutateInResult, err error) {
			if err != nil {
				ltc.updateResourceUnitsError(err)
				if errors.Is(err, ErrDocumentNotFound) || errors.Is(err, ErrPathNotFound) {
					logDebugf("Client %s not found in client record for %s, location = %s: %v", uuid, ltc.uuid, location, err)
					cb(nil)
					return
				}
				logDebugf("Failed to remove client %s for %s, location = %s: %v", uuid, ltc.uuid, location, err)

				go func() {
					select {
					case <-time.After(time.Until(deadline)):
						logDebugf("Timed out removing client %s from client record for %s, location = %s", uuid, ltc.uuid, location)
						cb(ErrTimeout)
						return
					case <-time.After(10 * time.Millisecond):
					}
					ltc.unregisterClientRecord(location, uuid, deadline, cb)
				}()
				return
			}

			ltc.updateResourceUnits(result.Internal.ResourceUnits)

			cb(nil)
		})
		if err != nil {
			logDebugf("Failed to schedule remove client %s for %s, location = %s: %v", uuid, ltc.uuid, location, err)
			select {
			case <-time.After(time.Until(deadline)):
				logDebugf("Timed out scheduling client removal %s from client record for %s, location = %s", uuid,
					ltc.uuid, location)
				cb(ErrTimeout)
				return
			case <-time.After(10 * time.Millisecond):
			}
			ltc.unregisterClientRecord(location, uuid, deadline, cb)
		}
	})
}

func (ltc *stdLostTransactionCleaner) perLocation(agent *Agent, oboUser string, location lostATRLocationWithShutdown) {
	logSchedf("Running cleanup %s on %s", ltc.uuid, location.location)
	ltc.processWaitGroup.Add(1)
	ltc.process(agent, oboUser, location.location.CollectionName, location.location.ScopeName, func(err error) {
		ltc.processWaitGroup.Done()
		if err != nil {
			logDebugf("Cleanup failed for %s on %s", ltc.uuid, location.location)
			// See comment in process for explanation of why we have a goroutine here.
			go func() {
				if errors.Is(err, ErrCollectionNotFound) || errors.Is(err, ErrScopeNotFound) {
					logDebugf("Removing %s.%s.%s from lost cleanup %s due to collection no longer existing",
						location.location.BucketName,
						location.location.ScopeName,
						location.location.CollectionName,
						ltc.uuid,
					)
					close(location.shutdown) // This is unlikely to do anything as we're only listening here but best be safe.
					ltc.locationsLock.Lock()
					if ltc.locations == nil {
						ltc.locationsLock.Unlock()
						return
					}
					delete(ltc.locations, location.location)
					ltc.locationsLock.Unlock()
					return
				}
				select {
				case <-ltc.stop:
					return
				case <-location.shutdown:
					return
				case <-time.After(1 * time.Second):
					ltc.perLocation(agent, oboUser, location)
					return
				}
			}()
			return
		}

		select {
		case <-ltc.stop:
			return
		case <-location.shutdown:
			return
		default:
		}
		ltc.perLocation(agent, oboUser, location)
	})
}

func (ltc *stdLostTransactionCleaner) process(agent *Agent, oboUser string, collection, scope string, cb func(error)) {
	ltc.ProcessClient(agent, oboUser, collection, scope, ltc.uuid, func(recordDetails *TransactionClientRecordDetails, err error) {
		if err != nil {
			logDebugf("Failed to process client %s on %s.%s.%s", ltc.uuid, agent.BucketName(), scope, collection)
			var coreErr *TimeoutError
			if errors.As(err, &coreErr) {
				for _, reason := range coreErr.RetryReasons {
					if reason == KVCollectionOutdatedRetryReason {
						// We translate from outdated to not found here because at the point in time when we tried to
						// use the collection it could not be found.
						cb(ErrCollectionNotFound)
						return
					}
				}
			}
			cb(err)
			return
		}

		logDebugf("%s will check %d atrs, check every %d ms", ltc.uuid, len(recordDetails.AtrsHandledByClient),
			recordDetails.CheckAtrEveryNMillis)

		// We need this goroutine so we can release the scope of the callback. We're still in the callback from the
		// LookupIn here so we're blocking the gocbcore read loop for the node, any further requests against that node
		// will never complete and timeout.
		go func() {
			d := time.Duration(recordDetails.CheckAtrEveryNMillis) * time.Millisecond
			for _, atr := range recordDetails.AtrsHandledByClient {
				select {
				case <-ltc.stop:
					cb(nil)
					return
				case <-time.After(d):
				}

				waitCh := make(chan error, 1)
				ltc.ProcessATR(agent, oboUser, collection, scope, atr, func(attempts []TransactionsCleanupAttempt,
					stats TransactionProcessATRStats, err error) {
					waitCh <- err
				})
				err := <-waitCh
				var coreErr *TimeoutError
				if errors.As(err, &coreErr) {
					for _, reason := range coreErr.RetryReasons {
						if reason == KVCollectionOutdatedRetryReason {
							cb(ErrCollectionNotFound)
							return
						}
					}
				}
			}

			cb(nil)
		}()
	})
}

// We pass uuid to this so that it's testable externally.
func (ltc *stdLostTransactionCleaner) ProcessClient(agent *Agent, oboUser string, collection, scope, uuid string, cb func(*TransactionClientRecordDetails, error)) {
	logSchedf("Processing client %s for %s.%s.%s", uuid, agent.BucketName(), scope, collection)
	ltc.clientRecordHooks.BeforeGetRecord(func(err error) {
		if err != nil {
			ec := classifyHookError(err)
			switch ec.Class {
			default:
				cb(nil, err)
				return
			case TransactionErrorClassFailDocAlreadyExists:
			case TransactionErrorClassFailCasMismatch:
			}
		}

		var deadline time.Time
		if ltc.keyValueTimeout > 0 {
			deadline = time.Now().Add(ltc.keyValueTimeout)
		}

		_, err = agent.LookupIn(LookupInOptions{
			Key: clientRecordKey,
			Ops: []SubDocOp{
				{
					Op:    memd.SubDocOpGet,
					Path:  "records",
					Flags: memd.SubdocFlagXattrPath,
				},
				{
					Op:    memd.SubDocOpGet,
					Path:  hlcMacro,
					Flags: memd.SubdocFlagXattrPath,
				},
			},
			Deadline:       deadline,
			CollectionName: collection,
			ScopeName:      scope,
			User:           oboUser,
		}, func(result *LookupInResult, err error) {
			if err != nil {
				ltc.updateResourceUnitsError(err)
				ec := classifyError(err)

				switch ec.Class {
				case TransactionErrorClassFailDocNotFound:
					ltc.createClientRecord(agent, oboUser, collection, scope, func(err error) {
						if err != nil {
							logDebugf("%s failed to create client record: %v", ltc.uuid, err)
							cb(nil, err)
							return
						}

						ltc.ProcessClient(agent, oboUser, collection, scope, uuid, cb)
					})
				default:
					cb(nil, err)
				}
				return
			}

			ltc.updateResourceUnits(result.Internal.ResourceUnits)

			recordOp := result.Ops[0]
			hlcOp := result.Ops[1]
			if recordOp.Err != nil {
				logDebugf("Failed to get records from client record for %s: %v", ltc.uuid, err)
				cb(nil, recordOp.Err)
				return
			}

			if hlcOp.Err != nil {
				logDebugf("Failed to get hlc from client record for %s: %v", ltc.uuid, err)
				cb(nil, hlcOp.Err)
				return
			}

			var records jsonClientRecords
			err = json.Unmarshal(recordOp.Value, &records)
			if err != nil {
				logDebugf("Failed to unmarshal records from client record for %s: %v", ltc.uuid, err)
				cb(nil, err)
				return
			}

			var hlc jsonHLC
			err = json.Unmarshal(hlcOp.Value, &hlc)
			if err != nil {
				logDebugf("Failed to unmarshal hlc from client record for %s: %v", ltc.uuid, err)
				cb(nil, err)
				return
			}

			nowSecs, err := parseHLCToSeconds(hlc)
			if err != nil {
				logDebugf("Failed to parse hlc from client record for %s: %v", ltc.uuid, err)
				cb(nil, err)
				return
			}
			nowMS := nowSecs * 1000 // we need it in millis

			recordDetails, err := ltc.parseClientRecords(records, uuid, nowMS)
			if err != nil {
				logDebugf("Failed to parse records from client record for %s: %v", ltc.uuid, err)
				cb(nil, err)
				return
			}

			if recordDetails.OverrideActive {
				cb(&recordDetails, nil)
				return
			}

			ltc.processClientRecord(agent, oboUser, collection, scope, uuid, recordDetails, func(err error) {
				if err != nil {
					logDebugf("%s failed to process client record %s: %v", ltc.uuid, uuid, err)
					cb(nil, err)
					return
				}

				cb(&recordDetails, nil)
			})
		})
		if err != nil {
			cb(nil, err)
			return
		}
	})
}

func (ltc *stdLostTransactionCleaner) ProcessATR(agent *Agent, oboUser string, collection, scope, atrID string, cb func([]TransactionsCleanupAttempt, TransactionProcessATRStats, error)) {
	ltc.getATR(agent, oboUser, collection, scope, atrID, func(attempts map[string]jsonAtrAttempt, hlc int64, err error) {
		if err != nil {
			// We want to be careful to not flood the logs with atr not found messages.
			if !errors.Is(err, ErrDocumentNotFound) {
				logSchedf("%s failed to get atr %s on %s.%s.%s", ltc.uuid, atrID, agent.BucketName(), scope, collection)
			}
			cb(nil, TransactionProcessATRStats{}, err)
			return
		}

		if len(attempts) == 0 {
			cb([]TransactionsCleanupAttempt{}, TransactionProcessATRStats{}, nil)
			return
		}

		logSchedf("%s processing %d entries for atr %s on %s.%s.%s", ltc.uuid, len(attempts), atrID, agent.BucketName(), scope, collection)

		stats := TransactionProcessATRStats{
			NumEntries: len(attempts),
		}

		// See the explanation in process, same idea.
		go func() {
			var results []TransactionsCleanupAttempt
			for key, attempt := range attempts {
				select {
				case <-ltc.stop:
					return
				default:
				}
				parsedCAS, err := parseCASToMilliseconds(attempt.PendingCAS)
				if err != nil {
					logDebugf("%s failed to parse CAS value %s for attempt %s on atr %s: %v", ltc.uuid,
						attempt.PendingCAS, key, atrID, err)
					cb(nil, TransactionProcessATRStats{}, err)
					return
				}
				var inserts []TransactionsDocRecord
				var replaces []TransactionsDocRecord
				var removes []TransactionsDocRecord
				for _, staged := range attempt.Inserts {
					inserts = append(inserts, TransactionsDocRecord{
						CollectionName: staged.CollectionName,
						ScopeName:      staged.ScopeName,
						BucketName:     staged.BucketName,
						ID:             []byte(staged.DocID),
					})
				}
				for _, staged := range attempt.Replaces {
					replaces = append(replaces, TransactionsDocRecord{
						CollectionName: staged.CollectionName,
						ScopeName:      staged.ScopeName,
						BucketName:     staged.BucketName,
						ID:             []byte(staged.DocID),
					})
				}
				for _, staged := range attempt.Removes {
					removes = append(removes, TransactionsDocRecord{
						CollectionName: staged.CollectionName,
						ScopeName:      staged.ScopeName,
						BucketName:     staged.BucketName,
						ID:             []byte(staged.DocID),
					})
				}

				var st TransactionAttemptState
				switch jsonAtrState(attempt.State) {
				case jsonAtrStateCommitted:
					st = TransactionAttemptStateCommitted
				case jsonAtrStateCompleted:
					st = TransactionAttemptStateCompleted
				case jsonAtrStatePending:
					st = TransactionAttemptStatePending
				case jsonAtrStateAborted:
					st = TransactionAttemptStateAborted
				case jsonAtrStateRolledBack:
					st = TransactionAttemptStateRolledBack
				default:
					continue
				}

				if int64(attempt.ExpiryTime)+parsedCAS < hlc {
					logDebugf("%s detected expired attempt %s on atr %s", ltc.uuid, key, atrID)
					req := &TransactionsCleanupRequest{
						AttemptID:         key,
						AtrID:             []byte(atrID),
						AtrCollectionName: collection,
						AtrScopeName:      scope,
						AtrBucketName:     agent.BucketName(),
						Inserts:           inserts,
						Replaces:          replaces,
						Removes:           removes,
						State:             st,
						ForwardCompat:     jsonForwardCompatToForwardCompat(attempt.ForwardCompat),
						DurabilityLevel:   transactionsDurabilityLevelFromShorthand(attempt.DurabilityLevel),
						Age:               time.Duration(hlc - parsedCAS),
					}

					waitCh := make(chan TransactionsCleanupAttempt, 1)
					ltc.cleaner.CleanupAttempt(agent, oboUser, req, false, func(attempt TransactionsCleanupAttempt) {
						waitCh <- attempt
					})
					attempt := <-waitCh
					results = append(results, attempt)
					stats.NumEntriesExpired++
				}
			}
			cb(results, stats, nil)
		}()
	})
}

func (ltc *stdLostTransactionCleaner) getATR(agent *Agent, oboUser string, collection, scope, atrID string,
	cb func(map[string]jsonAtrAttempt, int64, error)) {
	ltc.cleanupHooks.BeforeATRGet([]byte(atrID), func(err error) {
		if err != nil {
			cb(nil, 0, err)
			return
		}

		var deadline time.Time
		if ltc.keyValueTimeout > 0 {
			deadline = time.Now().Add(ltc.keyValueTimeout)
		}

		_, err = agent.LookupIn(LookupInOptions{
			Key: []byte(atrID),
			Ops: []SubDocOp{
				{
					Op:    memd.SubDocOpGet,
					Path:  "attempts",
					Flags: memd.SubdocFlagXattrPath,
				},
				{
					Op:    memd.SubDocOpGet,
					Path:  hlcMacro,
					Flags: memd.SubdocFlagXattrPath,
				},
			},
			Deadline:       deadline,
			CollectionName: collection,
			ScopeName:      scope,
			User:           oboUser,
		}, func(result *LookupInResult, err error) {
			if err != nil {
				ltc.updateResourceUnitsError(err)
				cb(nil, 0, err)
				return
			}

			ltc.updateResourceUnits(result.Internal.ResourceUnits)

			if result.Ops[0].Err != nil {
				cb(nil, 0, result.Ops[0].Err)
				return
			}

			if result.Ops[1].Err != nil {
				cb(nil, 0, result.Ops[1].Err)
				return
			}

			var attempts map[string]jsonAtrAttempt
			err = json.Unmarshal(result.Ops[0].Value, &attempts)
			if err != nil {
				cb(nil, 0, err)
				return
			}

			var hlc jsonHLC
			err = json.Unmarshal(result.Ops[1].Value, &hlc)
			if err != nil {
				cb(nil, 0, err)
				return
			}

			nowSecs, err := parseHLCToSeconds(hlc)
			if err != nil {
				cb(nil, 0, err)
				return
			}
			nowMS := nowSecs * 1000 // we need it in millis

			cb(attempts, nowMS, err)
		})
		if err != nil {
			cb(nil, 0, err)
			return
		}
	})
}

func (ltc *stdLostTransactionCleaner) parseClientRecords(records jsonClientRecords, uuid string, hlc int64) (TransactionClientRecordDetails, error) {
	var expiredIDs []string
	var activeIDs []string
	var clientAlreadyExists bool

	for u, client := range records.Clients {
		if u == uuid {
			activeIDs = append(activeIDs, u)
			clientAlreadyExists = true
			continue
		}

		heartbeatMS, err := parseCASToMilliseconds(client.HeartbeatMS)
		if err != nil {
			return TransactionClientRecordDetails{}, err
		}
		expiredPeriod := hlc - heartbeatMS

		if expiredPeriod >= int64(client.ExpiresMS) {
			expiredIDs = append(expiredIDs, u)
		} else {
			activeIDs = append(activeIDs, u)
		}
	}

	if !clientAlreadyExists {
		activeIDs = append(activeIDs, uuid)
	}

	sort.Strings(activeIDs)

	clientIndex := 0
	for i, u := range activeIDs {
		if u == uuid {
			clientIndex = i
			break
		}
	}

	var overrideEnabled bool
	var overrideActive bool
	var overrideExpiresCas int64

	if records.Override != nil {
		overrideEnabled = records.Override.Enabled
		overrideExpiresCas = records.Override.ExpiresNanos
		hlcNanos := hlc * 1000000

		if overrideEnabled && overrideExpiresCas > hlcNanos {
			overrideActive = true
		}
	}

	numActive := len(activeIDs)
	numExpired := len(expiredIDs)

	atrsHandled := atrsToHandle(clientIndex, numActive, ltc.numAtrs)

	checkAtrEveryNS := ltc.cleanupWindow.Milliseconds() / int64(len(atrsHandled))
	checkAtrEveryNMS := int(math.Max(1, float64(checkAtrEveryNS)))

	return TransactionClientRecordDetails{
		NumActiveClients:     numActive,
		IndexOfThisClient:    clientIndex,
		ClientIsNew:          clientAlreadyExists,
		ExpiredClientIDs:     expiredIDs,
		NumExistingClients:   numActive + numExpired,
		NumExpiredClients:    numExpired,
		OverrideEnabled:      overrideEnabled,
		OverrideActive:       overrideActive,
		OverrideExpiresCas:   overrideExpiresCas,
		CasNowNanos:          hlc,
		AtrsHandledByClient:  atrsHandled,
		CheckAtrEveryNMillis: checkAtrEveryNMS,
		ClientUUID:           uuid,
	}, nil
}

func (ltc *stdLostTransactionCleaner) processClientRecord(agent *Agent, oboUser string, collection, scope, uuid string,
	recordDetails TransactionClientRecordDetails, cb func(error)) {
	logSchedf("%s processing client record %s for %s.%s.%s", ltc.uuid, uuid, agent.BucketName(), scope, collection)
	ltc.clientRecordHooks.BeforeUpdateRecord(func(err error) {
		if err != nil {
			cb(err)
			return
		}

		prefix := "records.clients." + uuid + "."
		var marshalErr error
		fieldOp := func(fieldName string, data interface{}, op memd.SubDocOpType, flags memd.SubdocFlag) SubDocOp {
			b, err := json.Marshal(data)
			if err != nil {
				marshalErr = err
				return SubDocOp{}
			}

			return SubDocOp{
				Op:    op,
				Flags: flags,
				Path:  prefix + fieldName,
				Value: b,
			}
		}

		if marshalErr != nil {
			cb(err)
			return
		}

		ops := []SubDocOp{
			fieldOp("heartbeat_ms", "${Mutation.CAS}", memd.SubDocOpDictSet,
				memd.SubdocFlagXattrPath|memd.SubdocFlagExpandMacros|memd.SubdocFlagMkDirP),
			fieldOp("expires_ms", (ltc.cleanupWindow + 20000*time.Millisecond).Milliseconds(),
				memd.SubDocOpDictSet, memd.SubdocFlagXattrPath),
			fieldOp("num_atrs", ltc.numAtrs, memd.SubDocOpDictSet, memd.SubdocFlagXattrPath),
			{
				Op:    memd.SubDocOpSetDoc,
				Flags: memd.SubdocFlagNone,
				Value: []byte{0},
			},
		}

		numOps := 12
		if len(recordDetails.ExpiredClientIDs) < 12 {
			numOps = len(recordDetails.ExpiredClientIDs)
		}

		for i := 0; i < numOps; i++ {
			ops = append(ops, SubDocOp{
				Op:    memd.SubDocOpDelete,
				Flags: memd.SubdocFlagXattrPath,
				Path:  "records.clients." + recordDetails.ExpiredClientIDs[i],
			})
		}

		deadline := time.Time{}
		if ltc.keyValueTimeout > 0 {
			deadline = time.Now().Add(ltc.keyValueTimeout)
		}

		_, err = agent.MutateIn(MutateInOptions{
			Key:            clientRecordKey,
			Ops:            ops,
			CollectionName: collection,
			ScopeName:      scope,
			Deadline:       deadline,
			User:           oboUser,
		}, func(result *MutateInResult, err error) {
			if err != nil {
				ltc.updateResourceUnitsError(err)
				cb(err)
				return
			}

			ltc.updateResourceUnits(result.Internal.ResourceUnits)

			cb(nil)
		})
		if err != nil {
			cb(err)
			return
		}
	})
}

func (ltc *stdLostTransactionCleaner) createClientRecord(agent *Agent, oboUser string, collection, scope string, cb func(error)) {
	logDebugf("%s creating client record in %s.%s.%s", ltc.uuid, agent.BucketName(), scope, collection)
	ltc.clientRecordHooks.BeforeCreateRecord(func(err error) {
		if err != nil {
			ec := classifyHookError(err)

			switch ec.Class {
			default:
				cb(err)
				return
			case TransactionErrorClassFailDocNotFound:
			}
		}

		var deadline time.Time
		if ltc.keyValueTimeout > 0 {
			deadline = time.Now().Add(ltc.keyValueTimeout)
		}

		_, err = agent.MutateIn(MutateInOptions{
			Key: clientRecordKey,
			Ops: []SubDocOp{
				{
					Op:    memd.SubDocOpDictAdd,
					Flags: memd.SubdocFlagXattrPath,
					Path:  "records.clients",
					Value: []byte{123, 125}, // {}
				},
				{
					Op:    memd.SubDocOpSetDoc,
					Flags: memd.SubdocFlagNone,
					Path:  "",
					Value: []byte{0},
				},
			},
			Flags:          memd.SubdocDocFlagAddDoc,
			Deadline:       deadline,
			CollectionName: collection,
			ScopeName:      scope,
			User:           oboUser,
		}, func(result *MutateInResult, err error) {
			if err != nil {
				ltc.updateResourceUnitsError(err)
				ec := classifyError(err)

				switch ec.Class {
				default:
					cb(err)
					return
				case TransactionErrorClassFailDocAlreadyExists:
				case TransactionErrorClassFailCasMismatch:
				}

				cb(nil)
				return
			}

			ltc.updateResourceUnits(result.Internal.ResourceUnits)
			cb(nil)
		})
		if err != nil {
			cb(err)
			return
		}
	})
}

func (ltc *stdLostTransactionCleaner) fetchExtraCleanupLocations() {
	if ltc.atrLocationFinder != nil {
		locations, err := ltc.atrLocationFinder()
		if err != nil {
			logDebugf("%s failed to fetch extra cleanup locations: %v", ltc.uuid, err)
			return
		}

		locationMap := make(map[TransactionLostATRLocation]struct{})
		for _, location := range locations {
			ltc.AddATRLocation(location)
			locationMap[location] = struct{}{}
		}
	}

}

func atrsToHandle(index int, numActive int, numAtrs int) []string {
	allAtrs := transactionAtrIDList[:numAtrs]
	var selectedAtrs []string
	for i := index; i < len(allAtrs); i += numActive {
		selectedAtrs = append(selectedAtrs, allAtrs[i])
	}

	return selectedAtrs
}
