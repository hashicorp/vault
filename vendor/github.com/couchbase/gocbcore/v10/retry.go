package gocbcore

import (
	"encoding/json"
	"math"
	"time"

	"github.com/couchbase/gocbcore/v10/memd"
)

// RetryRequest is a request that can possibly be retried.
type RetryRequest interface {
	RetryAttempts() uint32
	Identifier() string
	Idempotent() bool
	RetryReasons() []RetryReason

	retryStrategy() RetryStrategy
	recordRetryAttempt(reason RetryReason)
}

// RetryReason represents the reason for an operation possibly being retried.
type RetryReason interface {
	AllowsNonIdempotentRetry() bool
	AlwaysRetry() bool
	Description() string
}

type retryReason struct {
	allowsNonIdempotentRetry bool
	alwaysRetry              bool
	description              string
}

func (rr retryReason) AllowsNonIdempotentRetry() bool {
	return rr.allowsNonIdempotentRetry
}

func (rr retryReason) AlwaysRetry() bool {
	return rr.alwaysRetry
}

func (rr retryReason) Description() string {
	return rr.description
}

func (rr retryReason) String() string {
	return rr.description
}

func (rr retryReason) MarshalJSON() ([]byte, error) {
	return json.Marshal(rr.description)
}

var (
	// UnknownRetryReason indicates that the operation failed for an unknown reason.
	UnknownRetryReason = retryReason{allowsNonIdempotentRetry: false, alwaysRetry: false, description: "UNKNOWN"}

	// SocketNotAvailableRetryReason indicates that the operation failed because the underlying socket was not available.
	SocketNotAvailableRetryReason = retryReason{allowsNonIdempotentRetry: true, alwaysRetry: false, description: "SOCKET_NOT_AVAILABLE"}

	// ServiceNotAvailableRetryReason indicates that the operation failed because the requested service was not available.
	ServiceNotAvailableRetryReason = retryReason{allowsNonIdempotentRetry: true, alwaysRetry: false, description: "SERVICE_NOT_AVAILABLE"}

	// NodeNotAvailableRetryReason indicates that the operation failed because the requested node was not available.
	NodeNotAvailableRetryReason = retryReason{allowsNonIdempotentRetry: true, alwaysRetry: false, description: "NODE_NOT_AVAILABLE"}

	// KVNotMyVBucketRetryReason indicates that the operation failed because it was sent to the wrong node for the vbucket.
	KVNotMyVBucketRetryReason = retryReason{allowsNonIdempotentRetry: true, alwaysRetry: true, description: "KV_NOT_MY_VBUCKET"}

	// KVCollectionOutdatedRetryReason indicates that the operation failed because the collection ID on the request is outdated.
	KVCollectionOutdatedRetryReason = retryReason{allowsNonIdempotentRetry: true, alwaysRetry: true, description: "KV_COLLECTION_OUTDATED"}

	// KVErrMapRetryReason indicates that the operation failed for an unsupported reason but the KV error map indicated
	// that the operation can be retried.
	KVErrMapRetryReason = retryReason{allowsNonIdempotentRetry: true, alwaysRetry: false, description: "KV_ERROR_MAP_RETRY_INDICATED"}

	// KVLockedRetryReason indicates that the operation failed because the document was locked.
	KVLockedRetryReason = retryReason{allowsNonIdempotentRetry: true, alwaysRetry: false, description: "KV_LOCKED"}

	// KVTemporaryFailureRetryReason indicates that the operation failed because of a temporary failure.
	KVTemporaryFailureRetryReason = retryReason{allowsNonIdempotentRetry: true, alwaysRetry: false, description: "KV_TEMPORARY_FAILURE"}

	// KVSyncWriteInProgressRetryReason indicates that the operation failed because a sync write is in progress.
	KVSyncWriteInProgressRetryReason = retryReason{allowsNonIdempotentRetry: true, alwaysRetry: false, description: "KV_SYNC_WRITE_IN_PROGRESS"}

	// KVSyncWriteRecommitInProgressRetryReason indicates that the operation failed because a sync write recommit is in progress.
	KVSyncWriteRecommitInProgressRetryReason = retryReason{allowsNonIdempotentRetry: true, alwaysRetry: false, description: "KV_SYNC_WRITE_RE_COMMIT_IN_PROGRESS"}

	// ServiceResponseCodeIndicatedRetryReason indicates that the operation failed and the service responded stating that
	// the request should be retried.
	ServiceResponseCodeIndicatedRetryReason = retryReason{allowsNonIdempotentRetry: true, alwaysRetry: false, description: "SERVICE_RESPONSE_CODE_INDICATED"}

	// SocketCloseInFlightRetryReason indicates that the operation failed because the socket was closed whilst the operation
	// was in flight.
	SocketCloseInFlightRetryReason = retryReason{allowsNonIdempotentRetry: false, alwaysRetry: false, description: "SOCKET_CLOSED_WHILE_IN_FLIGHT"}

	// PipelineOverloadedRetryReason indicates that the operation failed because the pipeline queue was full.
	PipelineOverloadedRetryReason = retryReason{allowsNonIdempotentRetry: true, alwaysRetry: true, description: "PIPELINE_OVERLOADED"}

	// CircuitBreakerOpenRetryReason indicates that the operation failed because the circuit breaker for the underlying socket was open.
	CircuitBreakerOpenRetryReason = retryReason{allowsNonIdempotentRetry: true, alwaysRetry: false, description: "CIRCUIT_BREAKER_OPEN"}

	// QueryIndexNotFoundRetryReason indicates that the operation failed to to a missing query index
	QueryIndexNotFoundRetryReason = retryReason{allowsNonIdempotentRetry: true, alwaysRetry: false, description: "QUERY_INDEX_NOT_FOUND"}

	// QueryPreparedStatementFailureRetryReason indicates that the operation failed due to a prepared statement failure
	QueryPreparedStatementFailureRetryReason = retryReason{allowsNonIdempotentRetry: true, alwaysRetry: false, description: "QUERY_PREPARED_STATEMENT_FAILURE"}

	// QueryErrorRetryable indicates that the operation is retryable as indicated by the query engine.
	QueryErrorRetryable = retryReason{allowsNonIdempotentRetry: true, alwaysRetry: false, description: "QUERY_ERROR_RETRYABLE"}

	// AnalyticsTemporaryFailureRetryReason indicates that an analytics operation failed due to a temporary failure
	AnalyticsTemporaryFailureRetryReason = retryReason{allowsNonIdempotentRetry: true, alwaysRetry: false, description: "ANALYTICS_TEMPORARY_FAILURE"}

	// SearchTooManyRequestsRetryReason indicates that a search operation failed due to too many requests
	SearchTooManyRequestsRetryReason = retryReason{allowsNonIdempotentRetry: true, alwaysRetry: false, description: "SEARCH_TOO_MANY_REQUESTS"}

	// NotReadyRetryReason indicates that the WaitUntilReady operation is not ready
	NotReadyRetryReason = retryReason{allowsNonIdempotentRetry: true, alwaysRetry: true, description: "NOT_READY"}

	// NoPipelineSnapshotRetryReason indicates that there was no pipeline snapshot available
	NoPipelineSnapshotRetryReason = retryReason{allowsNonIdempotentRetry: true, alwaysRetry: false, description: "NO_PIPELINE_SNAPSHOT"}

	// BucketNotReadyReason indicates that the user has priviledges to access the bucket but the bucket doesn't exist
	// or is in warm up.
	// Uncommitted: This API may change in the future.
	BucketNotReadyReason = retryReason{allowsNonIdempotentRetry: true, alwaysRetry: false, description: "BUCKET_NOT_FOUND"}

	// ConnectionErrorRetryReason indicates that there were errors reported by underlying connections.
	// Check server ports and cluster encryption setting.
	ConnectionErrorRetryReason = retryReason{allowsNonIdempotentRetry: true, alwaysRetry: false, description: "CONNECTION_ERROR"}

	// MemdWriteFailure indicates that the operation failed because the write failed on the connection.
	MemdWriteFailure = retryReason{allowsNonIdempotentRetry: true, alwaysRetry: true, description: "MEMD_WRITE_FAILURE"}

	// CredentialsFetchFailedRetryReason indicates that the operation failed because the AuthProvider return an error for credentials.
	// Uncommitted: This API may change in the future.
	CredentialsFetchFailedRetryReason = retryReason{allowsNonIdempotentRetry: true, alwaysRetry: true, description: "CREDENTIALS_FETCH_FAILED"}
)

// MaybeRetryRequest will possibly retry a request according to the strategy belonging to the request.
// It will use the reason to determine whether or not the failure reason is one that can be retried.
func (agent *Agent) MaybeRetryRequest(req RetryRequest, reason RetryReason) (bool, time.Time) {
	return retryOrchMaybeRetry(req, reason)
}

// RetryAction is used by a RetryStrategy to calculate the duration to wait before retrying an operation.
// Returning a value of 0 indicates to not retry.
type RetryAction interface {
	Duration() time.Duration
}

// NoRetryRetryAction represents an action that indicates to not retry.
type NoRetryRetryAction struct {
}

// Duration is the length of time to wait before retrying an operation.
func (ra *NoRetryRetryAction) Duration() time.Duration {
	return 0
}

// WithDurationRetryAction represents an action that indicates to retry with a given duration.
type WithDurationRetryAction struct {
	WithDuration time.Duration
}

// Duration is the length of time to wait before retrying an operation.
func (ra *WithDurationRetryAction) Duration() time.Duration {
	return ra.WithDuration
}

// RetryStrategy is to determine if an operation should be retried, and if so how long to wait before retrying.
type RetryStrategy interface {
	RetryAfter(req RetryRequest, reason RetryReason) RetryAction
}

// retryOrchMaybeRetry will possibly retry an operation according to the strategy belonging to the request.
// It will use the reason to determine whether or not the failure reason is one that can be retried.
func retryOrchMaybeRetry(req RetryRequest, reason RetryReason) (bool, time.Time) {
	if reason.AlwaysRetry() {
		duration := ControlledBackoff(req.RetryAttempts())
		logDebugf("Will retry request. Backoff=%s, OperationID=%s. Reason=%s", duration, req.Identifier(), reason)

		req.recordRetryAttempt(reason)

		return true, time.Now().Add(duration)
	}

	retryStrategy := req.retryStrategy()
	if retryStrategy == nil {
		return false, time.Time{}
	}

	action := retryStrategy.RetryAfter(req, reason)
	if action == nil {
		logDebugf("Won't retry request.  OperationID=%s. Reason=%s", req.Identifier(), reason)
		return false, time.Time{}
	}

	duration := action.Duration()
	if duration == 0 {
		logDebugf("Won't retry request.  OperationID=%s. Reason=%s", req.Identifier(), reason)
		return false, time.Time{}
	}

	logDebugf("Will retry request. Backoff=%s, OperationID=%s. Reason=%s", duration, req.Identifier(), reason)
	req.recordRetryAttempt(reason)

	return true, time.Now().Add(duration)
}

// failFastRetryStrategy represents a strategy that will never retry.
type failFastRetryStrategy struct {
}

// newFailFastRetryStrategy returns a new FailFastRetryStrategy.
func newFailFastRetryStrategy() *failFastRetryStrategy {
	return &failFastRetryStrategy{}
}

// RetryAfter calculates and returns a RetryAction describing how long to wait before retrying an operation.
func (rs *failFastRetryStrategy) RetryAfter(req RetryRequest, reason RetryReason) RetryAction {
	return &NoRetryRetryAction{}
}

// BackoffCalculator is used by retry strategies to calculate backoff durations.
type BackoffCalculator func(retryAttempts uint32) time.Duration

// BestEffortRetryStrategy represents a strategy that will keep retrying until it succeeds (or the caller times out
// the request).
type BestEffortRetryStrategy struct {
	backoffCalculator BackoffCalculator
}

// NewBestEffortRetryStrategy returns a new BestEffortRetryStrategy which will use the supplied calculator function
// to calculate retry durations. If calculator is nil then ControlledBackoff will be used.
func NewBestEffortRetryStrategy(calculator BackoffCalculator) *BestEffortRetryStrategy {
	if calculator == nil {
		calculator = ControlledBackoff
	}

	return &BestEffortRetryStrategy{backoffCalculator: calculator}
}

// RetryAfter calculates and returns a RetryAction describing how long to wait before retrying an operation.
func (rs *BestEffortRetryStrategy) RetryAfter(req RetryRequest, reason RetryReason) RetryAction {
	if req.Idempotent() || reason.AllowsNonIdempotentRetry() {
		return &WithDurationRetryAction{WithDuration: rs.backoffCalculator(req.RetryAttempts())}
	}

	return &NoRetryRetryAction{}
}

// ExponentialBackoff calculates a backoff time duration from the retry attempts on a given request.
func ExponentialBackoff(min, max time.Duration, backoffFactor float64) BackoffCalculator {
	var minBackoff float64 = 1000000   // 1 Millisecond
	var maxBackoff float64 = 500000000 // 500 Milliseconds
	var factor float64 = 2

	if min > 0 {
		minBackoff = float64(min)
	}
	if max > 0 {
		maxBackoff = float64(max)
	}
	if backoffFactor > 0 {
		factor = backoffFactor
	}

	return func(retryAttempts uint32) time.Duration {
		backoff := minBackoff * (math.Pow(factor, float64(retryAttempts)))

		if backoff > maxBackoff {
			backoff = maxBackoff
		}
		if backoff < minBackoff {
			backoff = minBackoff
		}

		return time.Duration(backoff)
	}
}

// ControlledBackoff calculates a backoff time duration from the retry attempts on a given request.
func ControlledBackoff(retryAttempts uint32) time.Duration {
	switch retryAttempts {
	case 0:
		return 1 * time.Millisecond
	case 1:
		return 10 * time.Millisecond
	case 2:
		return 50 * time.Millisecond
	case 3:
		return 100 * time.Millisecond
	case 4:
		return 500 * time.Millisecond
	default:
		return 1000 * time.Millisecond
	}
}

var idempotentOps = map[memd.CmdCode]bool{
	memd.CmdGet:                    true,
	memd.CmdGetReplica:             true,
	memd.CmdGetMeta:                true,
	memd.CmdSubDocGet:              true,
	memd.CmdSubDocExists:           true,
	memd.CmdSubDocGetCount:         true,
	memd.CmdNoop:                   true,
	memd.CmdStat:                   true,
	memd.CmdGetRandom:              true,
	memd.CmdCollectionsGetID:       true,
	memd.CmdCollectionsGetManifest: true,
	memd.CmdGetClusterConfig:       true,
	memd.CmdObserve:                true,
	memd.CmdObserveSeqNo:           true,
}
