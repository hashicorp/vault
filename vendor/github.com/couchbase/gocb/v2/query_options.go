package gocb

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// QueryScanConsistency indicates the level of data consistency desired for a query.
type QueryScanConsistency uint

const (
	// QueryScanConsistencyNotBounded indicates no data consistency is required.
	QueryScanConsistencyNotBounded QueryScanConsistency = iota + 1
	// QueryScanConsistencyRequestPlus indicates that request-level data consistency is required.
	QueryScanConsistencyRequestPlus
)

// QueryUseReplicaLevel specifies whether to enable replica reads for the request.
// If left unset will default to the value specified on the cluster.
type QueryUseReplicaLevel uint

const (
	// QueryUseReplicaLevelNotSet indicates to not set any replica level for the request.
	QueryUseReplicaLevelNotSet QueryUseReplicaLevel = iota
	// QueryUseReplicaLevelOff indicates to disable replica reads.
	QueryUseReplicaLevelOff
	// QueryUseReplicaLevelOn indicates to enable replica reads.
	QueryUseReplicaLevelOn
)

// QueryOptions represents the options available when executing a query.
type QueryOptions struct {
	ScanConsistency QueryScanConsistency
	ConsistentWith  *MutationState
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

	// MaxParallelism is the maximum number of index partitions, for computing aggregation in parallel.
	MaxParallelism uint32

	// ClientContextID provides a unique ID for this query which can be used matching up requests between connectionManager and
	// server. If not provided will be assigned a uuid value.
	ClientContextID      string
	PositionalParameters []interface{}
	NamedParameters      map[string]interface{}
	Metrics              bool

	// Raw provides a way to provide extra parameters in the request body for the query.
	Raw map[string]interface{}

	Adhoc         bool
	Timeout       time.Duration
	RetryStrategy RetryStrategy

	// FlexIndex tells the query engine to use a flex index (utilizing the search service).
	FlexIndex bool

	// PreserveExpiry tells the query engine to preserve expiration values set on any documents modified by this query.
	PreserveExpiry bool

	ParentSpan RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context

	// AsTransaction indicates to run this query as a transaction, providing any additional transaction specific
	// configuration.
	// UNCOMMITTED: This API may change in the future.
	AsTransaction *SingleQueryTransactionOptions

	// UseReplica tells the query engine to use replicas, if required, for the query.
	// This means that results could come from either active or replica nodes, depending on the state of the active node.
	// If any of the results came from a replica node then a warning will be populated in the query metadata.
	// If not set then this field is not sent in the query payload and the default setting on the cluster/node will be used.
	UseReplica QueryUseReplicaLevel

	// Internal: This should never be used and is not supported.
	Internal struct {
		User     string
		Endpoint string
	}
}

func (opts *QueryOptions) toMap() (map[string]interface{}, error) {
	execOpts := make(map[string]interface{})

	if opts.ScanConsistency != 0 && opts.ConsistentWith != nil {
		return nil, makeInvalidArgumentsError("ScanConsistency and ConsistentWith must be used exclusively")
	}

	if opts.ScanConsistency != 0 {
		if opts.ScanConsistency == QueryScanConsistencyNotBounded {
			execOpts["scan_consistency"] = "not_bounded"
		} else if opts.ScanConsistency == QueryScanConsistencyRequestPlus {
			execOpts["scan_consistency"] = "request_plus"
		} else {
			return nil, makeInvalidArgumentsError("Unexpected consistency option")
		}
	}

	if opts.ConsistentWith != nil {
		execOpts["scan_consistency"] = "at_plus"
		execOpts["scan_vectors"] = opts.ConsistentWith
	}

	if opts.Profile != "" {
		execOpts["profile"] = opts.Profile
	}

	if opts.Readonly {
		execOpts["readonly"] = opts.Readonly
	}

	if opts.PositionalParameters != nil {
		execOpts["args"] = opts.PositionalParameters
	}

	if opts.NamedParameters != nil {
		for key, value := range opts.NamedParameters {
			if !strings.HasPrefix(key, "$") {
				key = "$" + key
			}
			execOpts[key] = value
		}
	}

	if opts.ScanCap != 0 {
		execOpts["scan_cap"] = strconv.FormatUint(uint64(opts.ScanCap), 10)
	}

	if opts.PipelineBatch != 0 {
		execOpts["pipeline_batch"] = strconv.FormatUint(uint64(opts.PipelineBatch), 10)
	}

	if opts.PipelineCap != 0 {
		execOpts["pipeline_cap"] = strconv.FormatUint(uint64(opts.PipelineCap), 10)
	}

	if opts.ScanWait > 0 {
		execOpts["scan_wait"] = opts.ScanWait.String()
	}

	if opts.Raw != nil {
		for k, v := range opts.Raw {
			execOpts[k] = v
		}
	}

	if opts.MaxParallelism > 0 {
		execOpts["max_parallelism"] = strconv.FormatUint(uint64(opts.MaxParallelism), 10)
	}

	if !opts.Metrics {
		execOpts["metrics"] = false
	}

	if opts.ClientContextID == "" {
		execOpts["client_context_id"] = uuid.New()
	} else {
		execOpts["client_context_id"] = opts.ClientContextID
	}

	if opts.FlexIndex {
		execOpts["use_fts"] = true
	}

	if opts.PreserveExpiry {
		execOpts["preserve_expiry"] = true
	}

	if opts.UseReplica != QueryUseReplicaLevelNotSet {
		if opts.UseReplica == QueryUseReplicaLevelOff {
			execOpts["use_replica"] = "off"
		} else if opts.UseReplica == QueryUseReplicaLevelOn {
			execOpts["use_replica"] = "on"
		} else {
			return nil, makeInvalidArgumentsError("Unexpected replica level option")
		}
	}

	return execOpts, nil
}
