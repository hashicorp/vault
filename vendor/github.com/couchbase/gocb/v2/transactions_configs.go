package gocb

import (
	"time"
)

// TransactionsCleanupConfig specifies various tunable options related to transactions cleanup.
type TransactionsCleanupConfig struct {
	// CleanupWindow specifies how often to the cleanup process runs
	// attempting to garbage collection transactions that have failed but
	// were not cleaned up by the previous client.
	CleanupWindow time.Duration

	// DisableClientAttemptCleanup controls where any transaction attempts made
	// by this client are automatically removed.
	DisableClientAttemptCleanup bool

	// DisableLostAttemptCleanup controls where a background process is created
	// to cleanup any ‘lost’ transaction attempts.
	DisableLostAttemptCleanup bool

	// CleanupQueueSize controls the maximum queue size for the cleanup thread.
	CleanupQueueSize uint32

	// CleanupCollections is a set of extra collections that should be monitored
	// by the cleanup thread.
	CleanupCollections []TransactionKeyspace
}

// TransactionsConfig specifies various tunable options related to transactions.
type TransactionsConfig struct {
	// MetadataCollection specifies a specific location to place meta-data.
	MetadataCollection *TransactionKeyspace

	// ExpirationTimout sets the maximum time that transactions created
	// by this Transactions object can run for, before expiring.
	Timeout time.Duration

	// DurabilityLevel specifies the durability level that should be used
	// for all write operations performed by this Transactions object.
	DurabilityLevel DurabilityLevel

	// QueryConfig specifies any query configuration to use in transactions.
	QueryConfig TransactionsQueryConfig

	// CleanupConfig specifies cleanup configuration to use in transactions.
	CleanupConfig TransactionsCleanupConfig

	// Internal specifies a set of options for internal use.
	// Internal: This should never be used and is not supported.
	Internal struct {
		Hooks             TransactionHooks
		CleanupHooks      TransactionCleanupHooks
		ClientRecordHooks TransactionClientRecordHooks
		NumATRs           int
	}
}

// TransactionOptions specifies options which can be overridden on a per transaction basis.
type TransactionOptions struct {
	// DurabilityLevel specifies the durability level that should be used
	// for all write operations performed by this transaction.
	DurabilityLevel DurabilityLevel

	// Timeout sets the maximum time that this transaction can run for, before expiring.
	Timeout time.Duration

	// MetadataCollection specifies a specific Collection to place meta-data.
	MetadataCollection *Collection

	// Internal specifies a set of options for internal use.
	// Internal: This should never be used and is not supported.
	Internal struct {
		Hooks TransactionHooks
	}
}

// TransactionsQueryConfig specifies various tunable query options related to transactions.
type TransactionsQueryConfig struct {
	ScanConsistency QueryScanConsistency
}

// SingleQueryTransactionOptions specifies various tunable query options related to single query transactions.
type SingleQueryTransactionOptions struct {
	DurabilityLevel DurabilityLevel

	// Internal specifies a set of options for internal use.
	// Internal: This should never be used and is not supported.
	Internal struct {
		Hooks TransactionHooks
	}
}

// TransactionKeyspace specifies a specific location where ATR entries should be
// placed when performing transactions.
type TransactionKeyspace struct {
	BucketName     string
	ScopeName      string
	CollectionName string
}

// TransactionQueryOptions specifies the set of options available when running queries as a part of a transaction.
// This is a subset of QueryOptions.
type TransactionQueryOptions struct {
	ScanConsistency QueryScanConsistency
	Profile         QueryProfileMode

	// ScanCap is the maximum buffered channel size between the indexer connectionManager and the query service for index scans.
	ScanCap uint32

	// PipelineBatch controls the number of items execution operators can batch for Fetch from the KV.
	PipelineBatch uint32

	// PipelineCap controls the maximum number of items each execution operator can buffer between various operators.
	PipelineCap uint32

	// ScanWait is how long the indexer is allowed to wait until it can satisfy ScanConsistency/ConsistentWith criteria.
	ScanWait time.Duration
	Readonly bool

	// ClientContextID provides a unique ID for this query which can be used matching up requests between connectionManager and
	// server. If not provided will be assigned a uuid value.
	ClientContextID      string
	PositionalParameters []interface{}
	NamedParameters      map[string]interface{}

	// FlexIndex tells the query engine to use a flex index (utilizing the search service).
	FlexIndex bool

	// Raw provides a way to provide extra parameters in the request body for the query.
	Raw map[string]interface{}

	Prepared bool

	Scope *Scope
}

func (qo *TransactionQueryOptions) toSDKOptions() QueryOptions {
	scanc := qo.ScanConsistency
	if scanc == 0 {
		scanc = QueryScanConsistencyRequestPlus
	}

	return QueryOptions{
		ScanConsistency:      scanc,
		Profile:              qo.Profile,
		ScanCap:              qo.ScanCap,
		PipelineBatch:        qo.PipelineBatch,
		PipelineCap:          qo.PipelineCap,
		ScanWait:             qo.ScanWait,
		Readonly:             qo.Readonly,
		ClientContextID:      qo.ClientContextID,
		PositionalParameters: qo.PositionalParameters,
		NamedParameters:      qo.NamedParameters,
		Raw:                  qo.Raw,
		Adhoc:                !qo.Prepared,
		FlexIndex:            qo.FlexIndex,
	}
}
