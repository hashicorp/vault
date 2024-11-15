package gocb

import (
	"context"
	"sync"
	"time"

	gocbcore "github.com/couchbase/gocbcore/v10"
)

// Cas represents the specific state of a document on the cluster.
type Cas gocbcore.Cas

// InsertOptions are options that can be applied to an Insert operation.
type InsertOptions struct {
	Expiry          time.Duration
	PersistTo       uint
	ReplicateTo     uint
	DurabilityLevel DurabilityLevel
	Transcoder      Transcoder
	Timeout         time.Duration
	RetryStrategy   RetryStrategy
	ParentSpan      RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context

	// Internal: This should never be used and is not supported.
	Internal struct {
		User string
	}
}

// Insert creates a new document in the Collection.
func (c *Collection) Insert(id string, val interface{}, opts *InsertOptions) (mutOut *MutationResult, errOut error) {
	return autoOpControl(c.kvController(), func(agent kvProvider) (*MutationResult, error) {
		if opts == nil {
			opts = &InsertOptions{}
		}

		return agent.Insert(c, id, val, opts)
	})
}

// UpsertOptions are options that can be applied to an Upsert operation.
type UpsertOptions struct {
	Expiry          time.Duration
	PersistTo       uint
	ReplicateTo     uint
	DurabilityLevel DurabilityLevel
	Transcoder      Transcoder
	Timeout         time.Duration
	RetryStrategy   RetryStrategy
	ParentSpan      RequestSpan
	PreserveExpiry  bool

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context

	// Internal: This should never be used and is not supported.
	Internal struct {
		User string
	}
}

// Upsert creates a new document in the Collection if it does not exist, if it does exist then it updates it.
func (c *Collection) Upsert(id string, val interface{}, opts *UpsertOptions) (mutOut *MutationResult, errOut error) {
	return autoOpControl(c.kvController(), func(agent kvProvider) (*MutationResult, error) {
		if opts == nil {
			opts = &UpsertOptions{}
		}

		return agent.Upsert(c, id, val, opts)
	})
}

// ReplaceOptions are the options available to a Replace operation.
type ReplaceOptions struct {
	Expiry          time.Duration
	Cas             Cas
	PersistTo       uint
	ReplicateTo     uint
	DurabilityLevel DurabilityLevel
	Transcoder      Transcoder
	Timeout         time.Duration
	RetryStrategy   RetryStrategy
	ParentSpan      RequestSpan
	PreserveExpiry  bool

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context

	// Internal: This should never be used and is not supported.
	Internal struct {
		User string
	}
}

// Replace updates a document in the collection.
func (c *Collection) Replace(id string, val interface{}, opts *ReplaceOptions) (mutOut *MutationResult, errOut error) {
	return autoOpControl(c.kvController(), func(agent kvProvider) (*MutationResult, error) {
		if opts == nil {
			opts = &ReplaceOptions{}
		}

		if opts.Expiry > 0 && opts.PreserveExpiry {
			return nil, makeInvalidArgumentsError("cannot use expiry and preserve ttl together for replace")
		}

		return agent.Replace(c, id, val, opts)
	})
}

// GetOptions are the options available to a Get operation.
type GetOptions struct {
	WithExpiry bool
	// Project causes the Get operation to only fetch the fields indicated
	// by the paths. The result of the operation is then treated as a
	// standard GetResult.
	Project       []string
	Transcoder    Transcoder
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context

	// Internal: This should never be used and is not supported.
	Internal struct {
		User string
	}
}

// Get performs a fetch operation against the collection. This can take 3 paths, a standard full document
// fetch, a subdocument full document fetch also fetching document expiry (when WithExpiry is set),
// or a subdocument fetch (when Project is used).
func (c *Collection) Get(id string, opts *GetOptions) (docOut *GetResult, errOut error) {
	return autoOpControl(c.kvController(), func(agent kvProvider) (*GetResult, error) {
		if opts == nil {
			opts = &GetOptions{}
		}

		return agent.Get(c, id, opts)
	})
}

// ExistsOptions are the options available to the Exists command.
type ExistsOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context

	// Internal: This should never be used and is not supported.
	Internal struct {
		User string
	}
}

// Exists checks if a document exists for the given id.
func (c *Collection) Exists(id string, opts *ExistsOptions) (docOut *ExistsResult, errOut error) {
	return autoOpControl(c.kvController(), func(agent kvProvider) (*ExistsResult, error) {
		if opts == nil {
			opts = &ExistsOptions{}
		}

		return agent.Exists(c, id, opts)
	})
}

// GetAllReplicaOptions are the options available to the GetAllReplicas command.
type GetAllReplicaOptions struct {
	Transcoder    Transcoder
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// UNCOMMITTED: This API may change in the future.
	ReadPreference ReadPreference

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context

	// Internal: This should never be used and is not supported.
	Internal struct {
		User string
	}

	noMetrics bool
}

// GetAllReplicasResult represents the results of a GetAllReplicas operation.
type GetAllReplicasResult struct {
	res replicasResult
}

// Next fetches the next replica result.
func (r *GetAllReplicasResult) Next() *GetReplicaResult {
	res := r.res.Next()
	if res == nil {
		return nil
	}
	return res.(*GetReplicaResult)
}

// Close cancels all remaining get replica requests.
func (r *GetAllReplicasResult) Close() error {
	return r.res.Close()
}

type replicasResult interface {
	Next() interface{}
	Close() error
}

type coreReplicasResult struct {
	lock                sync.Mutex
	totalRequests       uint32
	successResults      uint32
	totalResults        uint32
	resCh               chan interface{}
	cancelCh            chan struct{}
	span                RequestSpan
	childReqsCompleteCh chan struct{}
	valueRecorder       ValueRecorder
	startedTime         time.Time
}

func (r *coreReplicasResult) addFailed() {
	r.lock.Lock()

	r.totalResults++
	if r.totalResults == r.totalRequests {
		close(r.childReqsCompleteCh)
		r.lock.Unlock()
		return
	}

	r.lock.Unlock()
}

func (r *coreReplicasResult) addResult(res interface{}) {
	// We use a lock here because the alternative means that there is a race
	// between the channel writes from multiple results and the channels being
	// closed.  IE: T1-Incr, T2-Incr, T2-Send, T2-Close, T1-Send[PANIC]
	r.lock.Lock()

	r.successResults++
	resultCount := r.successResults

	if resultCount <= r.totalRequests {
		r.resCh <- res
	}

	if resultCount == r.totalRequests {
		close(r.cancelCh)
		close(r.resCh)

		r.span.End()
		if r.valueRecorder != nil {
			r.valueRecorder.RecordValue(uint64(time.Since(r.startedTime).Microseconds()))
		}
	}

	r.totalResults++
	if r.totalResults == r.totalRequests {
		close(r.childReqsCompleteCh)
	}

	r.lock.Unlock()
}

// Next fetches the next replica result.
func (r *coreReplicasResult) Next() interface{} {
	return <-r.resCh
}

// Close cancels all remaining get replica requests.
func (r *coreReplicasResult) Close() error {
	// See addResult discussion on lock usage.
	r.lock.Lock()

	// Note that this number increment must be high enough to be clear that
	// the result set was closed, but low enough that it won't overflow if
	// additional result objects are processed after the close.
	prevResultCount := r.successResults
	r.successResults += 100000

	// We only have to close everything if the addResult method didn't already
	// close them due to already having completed every request
	var weClosed bool
	if prevResultCount < r.totalRequests {
		close(r.cancelCh)
		close(r.resCh)

		weClosed = true
	}

	r.lock.Unlock()

	if weClosed {
		// We need to wait for the child requests spans to be completed.
		<-r.childReqsCompleteCh
		r.span.End()
	}

	return nil
}

// GetAllReplicas returns the value of a particular document from all replica servers. This will return an iterable
// which streams results one at a time.
func (c *Collection) GetAllReplicas(id string, opts *GetAllReplicaOptions) (*GetAllReplicasResult, error) {
	return autoOpControl(c.kvController(), func(agent kvProvider) (*GetAllReplicasResult, error) {
		if opts == nil {
			opts = &GetAllReplicaOptions{}
		}

		return agent.GetAllReplicas(c, id, opts)
	})
}

// GetAnyReplicaOptions are the options available to the GetAnyReplica command.
type GetAnyReplicaOptions struct {
	Transcoder    Transcoder
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// UNCOMMITTED: This API may change in the future.
	ReadPreference ReadPreference

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context

	// Internal: This should never be used and is not supported.
	Internal struct {
		User string
	}
}

// GetAnyReplica returns the value of a particular document from a replica server.
func (c *Collection) GetAnyReplica(id string, opts *GetAnyReplicaOptions) (*GetReplicaResult, error) {
	return autoOpControl(c.kvController(), func(agent kvProvider) (*GetReplicaResult, error) {
		if opts == nil {
			opts = &GetAnyReplicaOptions{}
		}

		return agent.GetAnyReplica(c, id, opts)
	})
}

// RemoveOptions are the options available to the Remove command.
type RemoveOptions struct {
	Cas             Cas
	PersistTo       uint
	ReplicateTo     uint
	DurabilityLevel DurabilityLevel
	Timeout         time.Duration
	RetryStrategy   RetryStrategy
	ParentSpan      RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context

	// Internal: This should never be used and is not supported.
	Internal struct {
		User string
	}
}

// Remove removes a document from the collection.
func (c *Collection) Remove(id string, opts *RemoveOptions) (mutOut *MutationResult, errOut error) {
	return autoOpControl(c.kvController(), func(agent kvProvider) (*MutationResult, error) {
		if opts == nil {
			opts = &RemoveOptions{}
		}

		return agent.Remove(c, id, opts)
	})
}

// GetAndTouchOptions are the options available to the GetAndTouch operation.
type GetAndTouchOptions struct {
	Transcoder    Transcoder
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context

	// Internal: This should never be used and is not supported.
	Internal struct {
		User string
	}
}

// GetAndTouch retrieves a document and simultaneously updates its expiry time.
func (c *Collection) GetAndTouch(id string, expiry time.Duration, opts *GetAndTouchOptions) (docOut *GetResult, errOut error) {
	return autoOpControl(c.kvController(), func(agent kvProvider) (*GetResult, error) {
		if opts == nil {
			opts = &GetAndTouchOptions{}
		}

		return agent.GetAndTouch(c, id, expiry, opts)
	})
}

// GetAndLockOptions are the options available to the GetAndLock operation.
type GetAndLockOptions struct {
	Transcoder    Transcoder
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context

	// Internal: This should never be used and is not supported.
	Internal struct {
		User string
	}
}

// GetAndLock locks a document for a period of time, providing exclusive RW access to it.
// A lockTime value of over 30 seconds will be treated as 30 seconds. The resolution used to send this value to
// the server is seconds and is calculated using uint32(lockTime/time.Second).
func (c *Collection) GetAndLock(id string, lockTime time.Duration, opts *GetAndLockOptions) (docOut *GetResult, errOut error) {
	return autoOpControl(c.kvController(), func(agent kvProvider) (*GetResult, error) {
		if opts == nil {
			opts = &GetAndLockOptions{}
		}

		return agent.GetAndLock(c, id, lockTime, opts)
	})
}

// UnlockOptions are the options available to the GetAndLock operation.
type UnlockOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context

	// Internal: This should never be used and is not supported.
	Internal struct {
		User string
	}
}

// Unlock unlocks a document which was locked with GetAndLock.
func (c *Collection) Unlock(id string, cas Cas, opts *UnlockOptions) (errOut error) {
	return autoOpControlErrorOnly(c.kvController(), func(agent kvProvider) error {
		if opts == nil {
			opts = &UnlockOptions{}
		}

		return agent.Unlock(c, id, cas, opts)
	})
}

// TouchOptions are the options available to the Touch operation.
type TouchOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context

	// Internal: This should never be used and is not supported.
	Internal struct {
		User string
	}
}

// Touch touches a document, specifying a new expiry time for it.
func (c *Collection) Touch(id string, expiry time.Duration, opts *TouchOptions) (mutOut *MutationResult, errOut error) {
	return autoOpControl(c.kvController(), func(agent kvProvider) (*MutationResult, error) {
		if opts == nil {
			opts = &TouchOptions{}
		}

		return agent.Touch(c, id, expiry, opts)
	})
}

// Binary creates and returns a BinaryCollection object.
func (c *Collection) Binary() *BinaryCollection {
	return &BinaryCollection{collection: c}
}
