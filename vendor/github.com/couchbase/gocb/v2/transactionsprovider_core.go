package gocb

import (
	"errors"
	"math"
	"sync"
	"time"

	"github.com/couchbase/gocbcore/v10"
)

type transactionsInternalCore struct {
	parent *transactionsProviderCore
}

type agentProvider interface {
	OpenBucket(string) error
	GetAgent(string) *gocbcore.Agent
}

type transactionsProviderCore struct {
	config     TransactionsConfig
	cluster    *Cluster
	transcoder Transcoder

	// Transactions bypasses the connection manager when getting agents to allow
	// lost cleanup to fetch agents once close has been called on the cluster.
	getAgentProvider agentProvider

	txns                *gocbcore.TransactionsManager
	hooksWrapper        transactionHooksWrapper
	cleanupHooksWrapper transactionCleanupHooksWrapper
	cleanupCollections  []gocbcore.TransactionLostATRLocation
}

func (t *transactionsProviderCore) Init(config TransactionsConfig, c *Cluster) error {
	// Note that gocbcore will handle a lot of default values for us.
	if config.QueryConfig.ScanConsistency == 0 {
		config.QueryConfig.ScanConsistency = QueryScanConsistencyRequestPlus
	}
	if config.DurabilityLevel == DurabilityLevelUnknown {
		config.DurabilityLevel = DurabilityLevelMajority
	}

	var hooksWrapper transactionHooksWrapper
	if config.Internal.Hooks == nil {
		hooksWrapper = &noopHooksWrapper{
			TransactionDefaultHooks: gocbcore.TransactionDefaultHooks{},
			hooks:                   transactionsDefaultHooks{},
		}
	} else {
		hooksWrapper = &coreTxnsHooksWrapper{
			hooks: config.Internal.Hooks,
		}
	}

	var cleanupHooksWrapper transactionCleanupHooksWrapper
	if config.Internal.CleanupHooks == nil {
		cleanupHooksWrapper = &noopCleanupHooksWrapper{
			TransactionDefaultCleanupHooks: gocbcore.TransactionDefaultCleanupHooks{},
		}
	} else {
		cleanupHooksWrapper = &coreTxnsCleanupHooksWrapper{
			CleanupHooks: config.Internal.CleanupHooks,
		}
	}

	var clientRecordHooksWrapper clientRecordHooksWrapper
	if config.Internal.ClientRecordHooks == nil {
		clientRecordHooksWrapper = &noopClientRecordHooksWrapper{
			TransactionDefaultCleanupHooks:      gocbcore.TransactionDefaultCleanupHooks{},
			TransactionDefaultClientRecordHooks: gocbcore.TransactionDefaultClientRecordHooks{},
		}
	} else {
		clientRecordHooksWrapper = &coreTxnsClientRecordHooksWrapper{
			coreTxnsCleanupHooksWrapper: coreTxnsCleanupHooksWrapper{
				CleanupHooks: config.Internal.CleanupHooks,
			},
			ClientRecordHooks: config.Internal.ClientRecordHooks,
		}
	}

	atrLocation := gocbcore.TransactionATRLocation{}
	if config.MetadataCollection != nil {
		customATRAgent, err := c.Bucket(config.MetadataCollection.BucketName).Internal().IORouter()
		if err != nil {
			return err
		}

		atrLocation.Agent = customATRAgent
		atrLocation.CollectionName = config.MetadataCollection.CollectionName
		atrLocation.ScopeName = config.MetadataCollection.ScopeName

		// We add the custom metadata collection to the cleanup collections so that lost cleanup starts watching it
		// immediately. Note that we don't do the same for the custom metadata on TransactionOptions, this is because
		// we know that that collection will be used in a transaction.
		var alreadyInCleanup bool
		for _, keySpace := range config.CleanupConfig.CleanupCollections {
			if keySpace == *config.MetadataCollection {
				alreadyInCleanup = true
				break
			}
		}

		if !alreadyInCleanup {
			config.CleanupConfig.CleanupCollections = append(config.CleanupConfig.CleanupCollections, *config.MetadataCollection)
		}
	}

	var cleanupLocs []gocbcore.TransactionLostATRLocation
	for _, keyspace := range config.CleanupConfig.CleanupCollections {
		cleanupLocs = append(cleanupLocs, gocbcore.TransactionLostATRLocation{
			BucketName:     keyspace.BucketName,
			ScopeName:      keyspace.ScopeName,
			CollectionName: keyspace.CollectionName,
		})
	}

	t.cluster = c
	t.config = config
	t.transcoder = NewJSONTranscoder()
	t.hooksWrapper = hooksWrapper
	t.cleanupHooksWrapper = cleanupHooksWrapper
	t.cleanupCollections = cleanupLocs

	corecfg := &gocbcore.TransactionsConfig{}
	corecfg.DurabilityLevel = gocbcore.TransactionDurabilityLevel(config.DurabilityLevel)
	corecfg.BucketAgentProvider = t.agentProvider
	corecfg.LostCleanupATRLocationProvider = t.atrLocationsProvider
	corecfg.CleanupClientAttempts = !config.CleanupConfig.DisableClientAttemptCleanup
	corecfg.CleanupQueueSize = config.CleanupConfig.CleanupQueueSize
	corecfg.ExpirationTime = config.Timeout
	corecfg.CleanupWindow = config.CleanupConfig.CleanupWindow
	corecfg.CleanupLostAttempts = !config.CleanupConfig.DisableLostAttemptCleanup
	corecfg.CustomATRLocation = atrLocation
	corecfg.Internal.Hooks = hooksWrapper
	corecfg.Internal.CleanUpHooks = cleanupHooksWrapper
	corecfg.Internal.ClientRecordHooks = clientRecordHooksWrapper
	corecfg.Internal.NumATRs = config.Internal.NumATRs
	corecfg.KeyValueTimeout = c.timeoutsConfig.KVTimeout

	txns, err := gocbcore.InitTransactions(corecfg)
	if err != nil {
		return err
	}

	t.txns = txns
	return nil
}

func (t *transactionsProviderCore) Run(logicFn AttemptFunc, perConfig *TransactionOptions, singleQueryMode bool) (*TransactionResult, error) {
	if perConfig == nil {
		perConfig = &TransactionOptions{
			DurabilityLevel: t.config.DurabilityLevel,
			Timeout:         t.config.Timeout,
		}
	}

	scanConsistency := t.config.QueryConfig.ScanConsistency

	// Gocbcore looks at whether the location agent is nil to verify whether CustomATRLocation has been set.
	atrLocation := gocbcore.TransactionATRLocation{}
	if perConfig.MetadataCollection != nil {
		customATRAgent, err := perConfig.MetadataCollection.bucket.Internal().IORouter()
		if err != nil {
			return nil, err
		}

		atrLocation.Agent = customATRAgent
		atrLocation.CollectionName = perConfig.MetadataCollection.Name()
		atrLocation.ScopeName = perConfig.MetadataCollection.ScopeName()
	}

	logger := newTransactionLogger()

	// TODO: fill in the rest of this config
	config := &gocbcore.TransactionOptions{
		DurabilityLevel:   gocbcore.TransactionDurabilityLevel(perConfig.DurabilityLevel),
		ExpirationTime:    perConfig.Timeout,
		CustomATRLocation: atrLocation,
		TransactionLogger: logger,
	}

	hooksWrapper := t.hooksWrapper
	if perConfig.Internal.Hooks != nil {
		hooksWrapper = &coreTxnsHooksWrapper{
			hooks: perConfig.Internal.Hooks,
		}
		config.Internal.Hooks = hooksWrapper
	}

	txn, err := t.txns.BeginTransaction(config)
	if err != nil {
		return nil, err
	}

	logger.setTxnID(txn.ID())

	retries := 0
	backoffCalc := func() time.Duration {
		var max float64 = 100000000 // 100 Milliseconds
		var min float64 = 1000000   // 1 Millisecond
		retries++
		backoff := min * (math.Pow(2, float64(retries)))

		if backoff > max {
			backoff = max
		}
		if backoff < min {
			backoff = min
		}

		return time.Duration(backoff)
	}

	for {
		err = txn.NewAttempt()
		if err != nil {
			return nil, err
		}

		attemptID := txn.Attempt().ID
		logDebugf("New transaction attempt starting for %s, %s", txn.ID(), attemptID)
		logger.logInfof(attemptID, "New transaction attempt starting")

		attempt := TransactionAttemptContext{
			txn:            txn,
			transcoder:     t.transcoder,
			hooks:          hooksWrapper.Hooks(),
			cluster:        t.cluster,
			queryStateLock: new(sync.Mutex),
			queryConfig: TransactionQueryOptions{
				ScanConsistency: scanConsistency,
			},
			logger:    logger,
			attemptID: attemptID,
		}

		if hooksWrapper != nil {
			hooksWrapper.SetAttemptContext(attempt)
		}

		lambdaErr := logicFn(&attempt)

		if !singleQueryMode && lambdaErr != nil {
			logger.logInfof(attemptID, "Lambda returned error and not single query mode")
			var txnErr *TransactionOperationFailedError
			if !errors.As(lambdaErr, &txnErr) {
				// We wrap non-TOF errors in a TOF.
				lambdaErr = operationFailed(transactionQueryOperationFailedDef{
					ShouldNotRetry:    true,
					ShouldNotRollback: false,
					Reason:            gocbcore.TransactionErrorReasonTransactionFailed,
					ErrorCause:        lambdaErr,
					ShouldNotCommit:   true,
				}, &attempt)
			}
		}

		finalErr := lambdaErr
		if !singleQueryMode {
			if attempt.canCommit() {
				finalErr = attempt.commit()
			}
			if attempt.shouldRollback() {
				rollbackErr := attempt.rollback()
				if rollbackErr != nil {
					logWarnf("rollback after error failed: %s", rollbackErr)
				}
			}
		}
		toRaise := attempt.finalErrorToRaise()

		if attempt.shouldRetry() && toRaise != gocbcore.TransactionErrorReasonSuccess {
			logDebugf("retrying lambda after backoff")
			sleep := backoffCalc()
			logger.logInfof(attemptID, "Will retry lambda after %s", sleep)
			time.Sleep(sleep)
			continue
		}

		// We don't want the TOF to be the cause in the final error we return so we unwrap it.
		var finalErrCause error
		if finalErr != nil {
			var txnErr *TransactionOperationFailedError
			if errors.As(finalErr, &txnErr) {
				finalErrCause = txnErr.InternalUnwrap()
			} else {
				finalErrCause = finalErr
			}
		}

		switch toRaise {
		case gocbcore.TransactionErrorReasonSuccess:
			if singleQueryMode && finalErr != nil {
				return nil, finalErr
			}

			unstagingComplete := attempt.attempt().State == TransactionAttemptStateCompleted

			return &TransactionResult{
				TransactionID:     txn.ID(),
				UnstagingComplete: unstagingComplete,
				Logs:              logger.Logs(),
			}, nil
		case gocbcore.TransactionErrorReasonTransactionFailed:
			return nil, &TransactionFailedError{
				cause: finalErrCause,
				result: &TransactionResult{
					TransactionID:     txn.ID(),
					UnstagingComplete: false,
					Logs:              logger.Logs(),
				},
			}
		case gocbcore.TransactionErrorReasonTransactionExpired:
			// If we expired during gocbcore auto-rollback then we return failed with the error cause rather
			// than expired. This occurs when we commit itself errors and gocbcore auto rolls back the transaction.
			if attempt.attempt().PreExpiryAutoRollback {
				return nil, &TransactionFailedError{
					cause: finalErrCause,
					result: &TransactionResult{
						TransactionID:     txn.ID(),
						UnstagingComplete: false,
						Logs:              logger.Logs(),
					},
				}
			}
			return nil, &TransactionExpiredError{
				result: &TransactionResult{
					TransactionID:     txn.ID(),
					UnstagingComplete: false,
					Logs:              logger.Logs(),
				},
			}
		case gocbcore.TransactionErrorReasonTransactionCommitAmbiguous:
			return nil, &TransactionCommitAmbiguousError{
				cause: finalErrCause,
				result: &TransactionResult{
					TransactionID:     txn.ID(),
					UnstagingComplete: false,
					Logs:              logger.Logs(),
				},
			}
		case gocbcore.TransactionErrorReasonTransactionFailedPostCommit:
			return &TransactionResult{
				TransactionID:     txn.ID(),
				UnstagingComplete: false,
				Logs:              logger.Logs(),
			}, nil
		default:
			return nil, errors.New("invalid final transaction state")
		}
	}
}

func (t *transactionsProviderCore) Internal() transactionsInternal {
	return &transactionsInternalCore{parent: t}
}

func (t *transactionsInternalCore) ForceCleanupQueue() []TransactionCleanupAttempt {
	waitCh := make(chan []gocbcore.TransactionsCleanupAttempt, 1)
	t.parent.txns.Internal().ForceCleanupQueue(func(attempts []gocbcore.TransactionsCleanupAttempt) {
		waitCh <- attempts
	})
	coreAttempts := <-waitCh

	var attempts []TransactionCleanupAttempt
	for _, attempt := range coreAttempts {
		attempts = append(attempts, cleanupAttemptFromCore(attempt))
	}

	return attempts
}

func (t *transactionsInternalCore) CleanupQueueLength() int32 {
	return t.parent.txns.Internal().CleanupQueueLength()
}

func (t *transactionsInternalCore) ClientCleanupEnabled() bool {
	return t.parent.txns.Config().CleanupClientAttempts
}

func (t *transactionsInternalCore) CleanupLocations() []gocbcore.TransactionLostATRLocation {
	return t.parent.txns.Internal().CleanupLocations()
}

func (t *transactionsProviderCore) agentProvider(bucketName string) (*gocbcore.Agent, string, error) {
	err := t.getAgentProvider.OpenBucket(bucketName)
	if err != nil {
		return nil, "", err
	}

	agent := t.getAgentProvider.GetAgent(bucketName)

	return agent, "", nil
}

func (t *transactionsProviderCore) atrLocationsProvider() ([]gocbcore.TransactionLostATRLocation, error) {
	return t.cleanupCollections, nil
}

// Close will shut down this Transactions object, shutting down all
// background tasks associated with it.
func (t *transactionsProviderCore) close() error {
	return t.txns.Close()
}
