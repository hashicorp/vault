package gocb

import (
	"github.com/couchbase/gocbcore/v10"
	"time"
)

// TransactionDocRecord represents an individual document operation requiring cleanup.
// Internal: This should never be used and is not supported.
type TransactionDocRecord struct {
	CollectionName string
	ScopeName      string
	BucketName     string
	ID             string
}

// TransactionCleanupAttempt represents the result of running cleanup for a transaction transactionAttempt.
// Internal: This should never be used and is not supported.
type TransactionCleanupAttempt struct {
	Success           bool
	IsReqular         bool
	AttemptID         string
	AtrID             string
	AtrCollectionName string
	AtrScopeName      string
	AtrBucketName     string
	Request           *TransactionCleanupRequest
}

// TransactionCleanupRequest represents a complete transaction transactionAttempt that requires cleanup.
// Internal: This should never be used and is not supported.
type TransactionCleanupRequest struct {
	AttemptID         string
	AtrID             string
	AtrCollectionName string
	AtrScopeName      string
	AtrBucketName     string
	Inserts           []TransactionDocRecord
	Replaces          []TransactionDocRecord
	Removes           []TransactionDocRecord
	State             TransactionAttemptState
	ForwardCompat     map[string][]TransactionsForwardCompatibilityEntry
}

// TransactionsForwardCompatibilityEntry represents a forward compatibility entry.
// Internal: This should never be used and is not supported.
type TransactionsForwardCompatibilityEntry struct {
	ProtocolVersion   string `json:"p,omitempty"`
	ProtocolExtension string `json:"e,omitempty"`
	Behaviour         string `json:"b,omitempty"`
	RetryInterval     int    `json:"ra,omitempty"`
}

// TransactionsClientRecordDetails is the result of processing a client record.
// Internal: This should never be used and is not supported.
type TransactionsClientRecordDetails struct {
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

// TransactionsProcessATRStats is the stats recorded when running a ProcessATR request.
// Internal: This should never be used and is not supported.
type TransactionsProcessATRStats struct {
	NumEntries        int
	NumEntriesExpired int
}

// TransactionsCleaner is responsible for performing cleanup of completed transactions.
// Internal: This should never be used and is not supported.
type TransactionsCleaner interface {
	AddRequest(req *TransactionCleanupRequest) bool
	PopRequest() *TransactionCleanupRequest
	ForceCleanupQueue() []TransactionCleanupAttempt
	QueueLength() int32
	CleanupAttempt(bucket *Bucket, isRegular bool, req *TransactionCleanupRequest) TransactionCleanupAttempt
	Close()
}

// NewTransactionsCleaner returns a TransactionsCleaner implementation.
// Internal: This should never be used and is not supported.
func NewTransactionsCleaner(bucketProvider TransactionsBucketProviderFn, config *TransactionsConfig) TransactionsCleaner {
	cleanupHooksWrapper := &coreTxnsCleanupHooksWrapper{
		CleanupHooks: config.Internal.CleanupHooks,
	}

	corecfg := &gocbcore.TransactionsConfig{}
	corecfg.DurabilityLevel = gocbcore.TransactionDurabilityLevel(config.DurabilityLevel)
	corecfg.Internal.Hooks = nil
	corecfg.CleanupQueueSize = config.CleanupConfig.CleanupQueueSize
	corecfg.BucketAgentProvider = func(bucketName string) (*gocbcore.Agent, string, error) {
		bucket, user, err := bucketProvider(bucketName)
		if err != nil {
			return nil, "", err
		}

		agent, err := bucket.Internal().IORouter()
		if err != nil {
			return nil, "", err
		}

		return agent, user, nil
	}
	corecfg.Internal.CleanUpHooks = cleanupHooksWrapper
	corecfg.Internal.NumATRs = config.Internal.NumATRs
	corecfg.KeyValueTimeout = 2500 * time.Millisecond

	return &coreTransactionsCleanerWrapper{
		wrapped: gocbcore.NewTransactionsCleaner(corecfg),
	}
}

type coreTransactionsCleanerWrapper struct {
	wrapped gocbcore.TransactionsCleaner
}

func (ccw *coreTransactionsCleanerWrapper) AddRequest(req *TransactionCleanupRequest) bool {
	return ccw.wrapped.AddRequest(cleanupRequestToCore(req))
}

func (ccw *coreTransactionsCleanerWrapper) PopRequest() *TransactionCleanupRequest {
	return cleanupRequestFromCore(ccw.wrapped.PopRequest())
}

func (ccw *coreTransactionsCleanerWrapper) ForceCleanupQueue() []TransactionCleanupAttempt {
	waitCh := make(chan []TransactionCleanupAttempt, 1)
	ccw.wrapped.ForceCleanupQueue(func(coreAttempts []gocbcore.TransactionsCleanupAttempt) {
		var attempts []TransactionCleanupAttempt
		for _, attempt := range coreAttempts {
			attempts = append(attempts, cleanupAttemptFromCore(attempt))
		}
		waitCh <- attempts
	})
	return <-waitCh
}

func (ccw *coreTransactionsCleanerWrapper) QueueLength() int32 {
	return ccw.wrapped.QueueLength()
}

func (ccw *coreTransactionsCleanerWrapper) CleanupAttempt(bucket *Bucket, isRegular bool, req *TransactionCleanupRequest) TransactionCleanupAttempt {
	waitCh := make(chan TransactionCleanupAttempt, 1)
	a, err := bucket.Internal().IORouter()
	if err != nil {
		return TransactionCleanupAttempt{
			Success:           false,
			IsReqular:         isRegular,
			AttemptID:         req.AttemptID,
			AtrID:             req.AtrID,
			AtrCollectionName: req.AtrCollectionName,
			AtrScopeName:      req.AtrScopeName,
			AtrBucketName:     req.AtrBucketName,
			Request:           req,
		}
	}
	ccw.wrapped.CleanupAttempt(a, "", cleanupRequestToCore(req), isRegular, func(attempt gocbcore.TransactionsCleanupAttempt) {
		waitCh <- cleanupAttemptFromCore(attempt)
	})
	return <-waitCh
}

func (ccw *coreTransactionsCleanerWrapper) Close() {
	ccw.wrapped.Close()
}

// LostTransactionsCleaner is responsible for performing cleanup of lost transactions.
// Internal: This should never be used and is not supported.
type LostTransactionsCleaner interface {
	ProcessATR(bucket *Bucket, collection, scope, atrID string) ([]TransactionCleanupAttempt, TransactionsProcessATRStats)
	ProcessClient(bucket *Bucket, collection, scope, clientUUID string) (*TransactionsClientRecordDetails, error)
	RemoveClient(uuid string) error
	Close()
}

type coreLostTransactionsCleanerWrapper struct {
	wrapped gocbcore.LostTransactionCleaner
}

func (clcw *coreLostTransactionsCleanerWrapper) ProcessATR(bucket *Bucket, collection, scope, atrID string) ([]TransactionCleanupAttempt, TransactionsProcessATRStats) {
	a, err := bucket.Internal().IORouter()
	if err != nil {
		return nil, TransactionsProcessATRStats{}
	}

	var ourAttempts []TransactionCleanupAttempt
	var ourStats TransactionsProcessATRStats
	waitCh := make(chan struct{}, 1)
	clcw.wrapped.ProcessATR(a, "", collection, scope, atrID, func(attempts []gocbcore.TransactionsCleanupAttempt, stats gocbcore.TransactionProcessATRStats, _ error) {
		for _, a := range attempts {
			ourAttempts = append(ourAttempts, cleanupAttemptFromCore(a))
		}
		ourStats = TransactionsProcessATRStats(stats)

		waitCh <- struct{}{}
	})

	<-waitCh
	return ourAttempts, ourStats
}

func (clcw *coreLostTransactionsCleanerWrapper) ProcessClient(bucket *Bucket, collection, scope, clientUUID string) (*TransactionsClientRecordDetails, error) {
	type result struct {
		recordDetails *TransactionsClientRecordDetails
		err           error
	}
	waitCh := make(chan result, 1)
	a, err := bucket.Internal().IORouter()
	if err != nil {
		return nil, err
	}

	clcw.wrapped.ProcessClient(a, "", collection, scope, clientUUID, func(details *gocbcore.TransactionClientRecordDetails, err error) {
		if err != nil {
			waitCh <- result{
				err: err,
			}
			return
		}
		waitCh <- result{
			recordDetails: &TransactionsClientRecordDetails{
				NumActiveClients:     details.NumActiveClients,
				IndexOfThisClient:    details.IndexOfThisClient,
				ClientIsNew:          details.ClientIsNew,
				ExpiredClientIDs:     details.ExpiredClientIDs,
				NumExistingClients:   details.NumExistingClients,
				NumExpiredClients:    details.NumExpiredClients,
				OverrideEnabled:      details.OverrideEnabled,
				OverrideActive:       details.OverrideActive,
				OverrideExpiresCas:   details.OverrideExpiresCas,
				CasNowNanos:          details.CasNowNanos,
				AtrsHandledByClient:  details.AtrsHandledByClient,
				CheckAtrEveryNMillis: details.CheckAtrEveryNMillis,
				ClientUUID:           details.ClientUUID,
			},
		}
	})

	res := <-waitCh
	return res.recordDetails, res.err
}

func (clcw *coreLostTransactionsCleanerWrapper) RemoveClient(uuid string) error {
	return clcw.wrapped.RemoveClientFromAllLocations(uuid)
}

func (clcw *coreLostTransactionsCleanerWrapper) Close() {
	clcw.wrapped.Close()
}

// TransactionsBucketProviderFn is a function used to provide a bucket for
// a particular bucket by name.
// Internal: This should never be used and is not supported.
type TransactionsBucketProviderFn func(bucket string) (*Bucket, string, error)

type TransactionsLostCleanupKeyspaceProviderFn func() ([]TransactionKeyspace, error)

// NewLostTransactionsCleanup returns a LostTransactionsCleaner implementation.
// Internal: This should never be used and is not supported.
func NewLostTransactionsCleanup(bucketProvider TransactionsBucketProviderFn, locationProvider TransactionsLostCleanupKeyspaceProviderFn,
	config *TransactionsConfig) LostTransactionsCleaner {
	cleanupHooksWrapper := &coreTxnsClientRecordHooksWrapper{
		coreTxnsCleanupHooksWrapper: coreTxnsCleanupHooksWrapper{
			CleanupHooks: config.Internal.CleanupHooks,
		},
		ClientRecordHooks: config.Internal.ClientRecordHooks,
	}

	corecfg := &gocbcore.TransactionsConfig{}
	corecfg.DurabilityLevel = gocbcore.TransactionDurabilityLevel(config.DurabilityLevel)
	corecfg.Internal.Hooks = nil
	corecfg.CleanupQueueSize = config.CleanupConfig.CleanupQueueSize
	corecfg.BucketAgentProvider = func(bucketName string) (*gocbcore.Agent, string, error) {
		bucket, user, err := bucketProvider(bucketName)
		if err != nil {
			return nil, "", err
		}

		agent, err := bucket.Internal().IORouter()
		if err != nil {
			return nil, "", err
		}

		return agent, user, nil
	}
	corecfg.LostCleanupATRLocationProvider = func() ([]gocbcore.TransactionLostATRLocation, error) {
		locations, err := locationProvider()
		if err != nil {
			return nil, err
		}

		atrLocs := make([]gocbcore.TransactionLostATRLocation, len(locations))
		for i, loc := range locations {
			atrLocs[i] = gocbcore.TransactionLostATRLocation{
				BucketName:     loc.BucketName,
				CollectionName: loc.CollectionName,
				ScopeName:      loc.ScopeName,
			}
		}

		return atrLocs, nil
	}
	corecfg.Internal.CleanUpHooks = cleanupHooksWrapper
	corecfg.Internal.ClientRecordHooks = cleanupHooksWrapper
	corecfg.Internal.NumATRs = config.Internal.NumATRs
	corecfg.KeyValueTimeout = 2500 * time.Millisecond

	return &coreLostTransactionsCleanerWrapper{
		wrapped: gocbcore.NewLostTransactionCleaner(corecfg),
	}
}

func cleanupAttemptFromCore(attempt gocbcore.TransactionsCleanupAttempt) TransactionCleanupAttempt {
	var req *TransactionCleanupRequest
	if attempt.Request != nil {
		req = &TransactionCleanupRequest{
			AttemptID:         attempt.Request.AttemptID,
			AtrID:             string(attempt.Request.AtrID),
			AtrCollectionName: attempt.Request.AtrCollectionName,
			AtrScopeName:      attempt.Request.AtrScopeName,
			AtrBucketName:     attempt.Request.AtrBucketName,
			Inserts:           docRecordsFromCore(attempt.Request.Inserts),
			Replaces:          docRecordsFromCore(attempt.Request.Replaces),
			Removes:           docRecordsFromCore(attempt.Request.Removes),
			State:             TransactionAttemptState(attempt.Request.State),
		}
	}
	return TransactionCleanupAttempt{
		Success:           attempt.Success,
		IsReqular:         attempt.IsReqular,
		AttemptID:         attempt.AttemptID,
		AtrID:             string(attempt.AtrID),
		AtrCollectionName: attempt.AtrCollectionName,
		AtrScopeName:      attempt.AtrScopeName,
		AtrBucketName:     attempt.AtrBucketName,
		Request:           req,
	}
}

func docRecordsFromCore(drs []gocbcore.TransactionsDocRecord) []TransactionDocRecord {
	var recs []TransactionDocRecord
	for _, i := range drs {
		recs = append(recs, TransactionDocRecord{
			CollectionName: i.CollectionName,
			ScopeName:      i.ScopeName,
			BucketName:     i.BucketName,
			ID:             string(i.ID),
		})
	}

	return recs
}

func cleanupRequestFromCore(request *gocbcore.TransactionsCleanupRequest) *TransactionCleanupRequest {
	forwardCompat := make(map[string][]TransactionsForwardCompatibilityEntry)
	for k, entries := range request.ForwardCompat {
		if _, ok := forwardCompat[k]; !ok {
			forwardCompat[k] = make([]TransactionsForwardCompatibilityEntry, len(entries))
		}

		for i, entry := range entries {
			forwardCompat[k][i] = TransactionsForwardCompatibilityEntry(entry)
		}
	}

	return &TransactionCleanupRequest{
		AttemptID:         request.AttemptID,
		AtrID:             string(request.AtrID),
		AtrCollectionName: request.AtrCollectionName,
		AtrScopeName:      request.AtrScopeName,
		AtrBucketName:     request.AtrBucketName,
		Inserts:           docRecordsFromCore(request.Inserts),
		Replaces:          docRecordsFromCore(request.Replaces),
		Removes:           docRecordsFromCore(request.Removes),
		State:             TransactionAttemptState(request.State),
		ForwardCompat:     forwardCompat,
	}
}

func cleanupRequestToCore(request *TransactionCleanupRequest) *gocbcore.TransactionsCleanupRequest {
	forwardCompat := make(map[string][]gocbcore.TransactionForwardCompatibilityEntry)
	for k, entries := range request.ForwardCompat {
		if _, ok := forwardCompat[k]; !ok {
			forwardCompat[k] = make([]gocbcore.TransactionForwardCompatibilityEntry, len(entries))
		}

		for i, entry := range entries {
			forwardCompat[k][i] = gocbcore.TransactionForwardCompatibilityEntry(entry)
		}
	}

	return &gocbcore.TransactionsCleanupRequest{
		AttemptID:         request.AttemptID,
		AtrID:             []byte(request.AtrID),
		AtrCollectionName: request.AtrCollectionName,
		AtrScopeName:      request.AtrScopeName,
		AtrBucketName:     request.AtrBucketName,
		Inserts:           docRecordsToCore(request.Inserts),
		Replaces:          docRecordsToCore(request.Replaces),
		Removes:           docRecordsToCore(request.Removes),
		State:             gocbcore.TransactionAttemptState(request.State),
		ForwardCompat:     forwardCompat,
	}
}

func docRecordsToCore(drs []TransactionDocRecord) []gocbcore.TransactionsDocRecord {
	var recs []gocbcore.TransactionsDocRecord
	for _, i := range drs {
		recs = append(recs, gocbcore.TransactionsDocRecord{
			CollectionName: i.CollectionName,
			ScopeName:      i.ScopeName,
			BucketName:     i.BucketName,
			ID:             []byte(i.ID),
		})
	}

	return recs
}
