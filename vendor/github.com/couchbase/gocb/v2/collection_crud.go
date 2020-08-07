package gocb

import (
	"errors"
	"sync"
	"time"

	gocbcore "github.com/couchbase/gocbcore/v9"
)

type kvProvider interface {
	Add(opts gocbcore.AddOptions, cb gocbcore.StoreCallback) (gocbcore.PendingOp, error)
	Set(opts gocbcore.SetOptions, cb gocbcore.StoreCallback) (gocbcore.PendingOp, error)
	Replace(opts gocbcore.ReplaceOptions, cb gocbcore.StoreCallback) (gocbcore.PendingOp, error)
	Get(opts gocbcore.GetOptions, cb gocbcore.GetCallback) (gocbcore.PendingOp, error)
	GetOneReplica(opts gocbcore.GetOneReplicaOptions, cb gocbcore.GetReplicaCallback) (gocbcore.PendingOp, error)
	Observe(opts gocbcore.ObserveOptions, cb gocbcore.ObserveCallback) (gocbcore.PendingOp, error)
	ObserveVb(opts gocbcore.ObserveVbOptions, cb gocbcore.ObserveVbCallback) (gocbcore.PendingOp, error)
	GetMeta(opts gocbcore.GetMetaOptions, cb gocbcore.GetMetaCallback) (gocbcore.PendingOp, error)
	Delete(opts gocbcore.DeleteOptions, cb gocbcore.DeleteCallback) (gocbcore.PendingOp, error)
	LookupIn(opts gocbcore.LookupInOptions, cb gocbcore.LookupInCallback) (gocbcore.PendingOp, error)
	MutateIn(opts gocbcore.MutateInOptions, cb gocbcore.MutateInCallback) (gocbcore.PendingOp, error)
	GetAndTouch(opts gocbcore.GetAndTouchOptions, cb gocbcore.GetAndTouchCallback) (gocbcore.PendingOp, error)
	GetAndLock(opts gocbcore.GetAndLockOptions, cb gocbcore.GetAndLockCallback) (gocbcore.PendingOp, error)
	Unlock(opts gocbcore.UnlockOptions, cb gocbcore.UnlockCallback) (gocbcore.PendingOp, error)
	Touch(opts gocbcore.TouchOptions, cb gocbcore.TouchCallback) (gocbcore.PendingOp, error)
	Increment(opts gocbcore.CounterOptions, cb gocbcore.CounterCallback) (gocbcore.PendingOp, error)
	Decrement(opts gocbcore.CounterOptions, cb gocbcore.CounterCallback) (gocbcore.PendingOp, error)
	Append(opts gocbcore.AdjoinOptions, cb gocbcore.AdjoinCallback) (gocbcore.PendingOp, error)
	Prepend(opts gocbcore.AdjoinOptions, cb gocbcore.AdjoinCallback) (gocbcore.PendingOp, error)
	ConfigSnapshot() (*gocbcore.ConfigSnapshot, error)
}

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
}

// Insert creates a new document in the Collection.
func (c *Collection) Insert(id string, val interface{}, opts *InsertOptions) (mutOut *MutationResult, errOut error) {
	if opts == nil {
		opts = &InsertOptions{}
	}

	opm := c.newKvOpManager("Insert", nil)
	defer opm.Finish()

	opm.SetDocumentID(id)
	opm.SetTranscoder(opts.Transcoder)
	opm.SetValue(val)
	opm.SetDuraOptions(opts.PersistTo, opts.ReplicateTo, opts.DurabilityLevel)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	agent, err := c.getKvProvider()
	if err != nil {
		return nil, err
	}
	err = opm.Wait(agent.Add(gocbcore.AddOptions{
		Key:                    opm.DocumentID(),
		Value:                  opm.ValueBytes(),
		Flags:                  opm.ValueFlags(),
		Expiry:                 durationToExpiry(opts.Expiry),
		CollectionName:         opm.CollectionName(),
		ScopeName:              opm.ScopeName(),
		DurabilityLevel:        opm.DurabilityLevel(),
		DurabilityLevelTimeout: opm.DurabilityTimeout(),
		RetryStrategy:          opm.RetryStrategy(),
		TraceContext:           opm.TraceSpan(),
		Deadline:               opm.Deadline(),
	}, func(res *gocbcore.StoreResult, err error) {
		if err != nil {
			errOut = opm.EnhanceErr(err)
			opm.Reject()
			return
		}

		mutOut = &MutationResult{}
		mutOut.cas = Cas(res.Cas)
		mutOut.mt = opm.EnhanceMt(res.MutationToken)

		opm.Resolve(mutOut.mt)
	}))
	if err != nil {
		errOut = err
	}
	return
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
}

// Upsert creates a new document in the Collection if it does not exist, if it does exist then it updates it.
func (c *Collection) Upsert(id string, val interface{}, opts *UpsertOptions) (mutOut *MutationResult, errOut error) {
	if opts == nil {
		opts = &UpsertOptions{}
	}

	opm := c.newKvOpManager("Upsert", nil)
	defer opm.Finish()

	opm.SetDocumentID(id)
	opm.SetTranscoder(opts.Transcoder)
	opm.SetValue(val)
	opm.SetDuraOptions(opts.PersistTo, opts.ReplicateTo, opts.DurabilityLevel)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	agent, err := c.getKvProvider()
	if err != nil {
		return nil, err
	}
	err = opm.Wait(agent.Set(gocbcore.SetOptions{
		Key:                    opm.DocumentID(),
		Value:                  opm.ValueBytes(),
		Flags:                  opm.ValueFlags(),
		Expiry:                 durationToExpiry(opts.Expiry),
		CollectionName:         opm.CollectionName(),
		ScopeName:              opm.ScopeName(),
		DurabilityLevel:        opm.DurabilityLevel(),
		DurabilityLevelTimeout: opm.DurabilityTimeout(),
		RetryStrategy:          opm.RetryStrategy(),
		TraceContext:           opm.TraceSpan(),
		Deadline:               opm.Deadline(),
	}, func(res *gocbcore.StoreResult, err error) {
		if err != nil {
			errOut = opm.EnhanceErr(err)
			opm.Reject()
			return
		}

		mutOut = &MutationResult{}
		mutOut.cas = Cas(res.Cas)
		mutOut.mt = opm.EnhanceMt(res.MutationToken)

		opm.Resolve(mutOut.mt)
	}))
	if err != nil {
		errOut = err
	}
	return
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
}

// Replace updates a document in the collection.
func (c *Collection) Replace(id string, val interface{}, opts *ReplaceOptions) (mutOut *MutationResult, errOut error) {
	if opts == nil {
		opts = &ReplaceOptions{}
	}

	opm := c.newKvOpManager("Replace", nil)
	defer opm.Finish()

	opm.SetDocumentID(id)
	opm.SetTranscoder(opts.Transcoder)
	opm.SetValue(val)
	opm.SetDuraOptions(opts.PersistTo, opts.ReplicateTo, opts.DurabilityLevel)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	agent, err := c.getKvProvider()
	if err != nil {
		return nil, err
	}
	err = opm.Wait(agent.Replace(gocbcore.ReplaceOptions{
		Key:                    opm.DocumentID(),
		Value:                  opm.ValueBytes(),
		Flags:                  opm.ValueFlags(),
		Expiry:                 durationToExpiry(opts.Expiry),
		Cas:                    gocbcore.Cas(opts.Cas),
		CollectionName:         opm.CollectionName(),
		ScopeName:              opm.ScopeName(),
		DurabilityLevel:        opm.DurabilityLevel(),
		DurabilityLevelTimeout: opm.DurabilityTimeout(),
		RetryStrategy:          opm.RetryStrategy(),
		TraceContext:           opm.TraceSpan(),
		Deadline:               opm.Deadline(),
	}, func(res *gocbcore.StoreResult, err error) {
		if err != nil {
			errOut = opm.EnhanceErr(err)
			opm.Reject()
			return
		}

		mutOut = &MutationResult{}
		mutOut.cas = Cas(res.Cas)
		mutOut.mt = opm.EnhanceMt(res.MutationToken)

		opm.Resolve(mutOut.mt)
	}))
	if err != nil {
		errOut = err
	}
	return
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
}

// Get performs a fetch operation against the collection. This can take 3 paths, a standard full document
// fetch, a subdocument full document fetch also fetching document expiry (when WithExpiry is set),
// or a subdocument fetch (when Project is used).
func (c *Collection) Get(id string, opts *GetOptions) (docOut *GetResult, errOut error) {
	if opts == nil {
		opts = &GetOptions{}
	}

	if len(opts.Project) == 0 && !opts.WithExpiry {
		return c.getDirect(id, opts)
	}

	return c.getProjected(id, opts)
}

func (c *Collection) getDirect(id string, opts *GetOptions) (docOut *GetResult, errOut error) {
	if opts == nil {
		opts = &GetOptions{}
	}

	opm := c.newKvOpManager("Get", nil)
	defer opm.Finish()

	opm.SetDocumentID(id)
	opm.SetTranscoder(opts.Transcoder)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	agent, err := c.getKvProvider()
	if err != nil {
		return nil, err
	}
	err = opm.Wait(agent.Get(gocbcore.GetOptions{
		Key:            opm.DocumentID(),
		CollectionName: opm.CollectionName(),
		ScopeName:      opm.ScopeName(),
		RetryStrategy:  opm.RetryStrategy(),
		TraceContext:   opm.TraceSpan(),
		Deadline:       opm.Deadline(),
	}, func(res *gocbcore.GetResult, err error) {
		if err != nil {
			errOut = opm.EnhanceErr(err)
			opm.Reject()
			return
		}

		doc := &GetResult{
			Result: Result{
				cas: Cas(res.Cas),
			},
			transcoder: opm.Transcoder(),
			contents:   res.Value,
			flags:      res.Flags,
		}

		docOut = doc

		opm.Resolve(nil)
	}))
	if err != nil {
		errOut = err
	}
	return
}

func (c *Collection) getProjected(id string, opts *GetOptions) (docOut *GetResult, errOut error) {
	if opts == nil {
		opts = &GetOptions{}
	}

	opm := c.newKvOpManager("Get", nil)
	defer opm.Finish()

	opm.SetDocumentID(id)
	opm.SetTranscoder(opts.Transcoder)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)

	if opts.Transcoder != nil {
		return nil, errors.New("Cannot specify custom transcoder for projected gets")
	}

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	numProjects := len(opts.Project)
	if opts.WithExpiry {
		numProjects = 1 + numProjects
	}

	projections := opts.Project
	if numProjects > 16 {
		projections = nil
	}

	var ops []LookupInSpec

	if opts.WithExpiry {
		ops = append(ops, GetSpec("$document.exptime", &GetSpecOptions{IsXattr: true}))
	}

	if len(projections) == 0 {
		ops = append(ops, GetSpec("", nil))
	} else {
		for _, path := range projections {
			ops = append(ops, GetSpec(path, nil))
		}
	}

	result, err := c.internalLookupIn(opm, ops, false)
	if err != nil {
		return nil, err
	}

	doc := &GetResult{}
	if opts.WithExpiry {
		// if expiration was requested then extract and remove it from the results
		err = result.ContentAt(0, &doc.expiry)
		if err != nil {
			return nil, err
		}
		ops = ops[1:]
		result.contents = result.contents[1:]
	}

	doc.transcoder = opm.Transcoder()
	doc.cas = result.cas
	if projections == nil {
		err = doc.fromFullProjection(ops, result, opts.Project)
		if err != nil {
			return nil, err
		}
	} else {
		err = doc.fromSubDoc(ops, result)
		if err != nil {
			return nil, err
		}
	}

	return doc, nil
}

// ExistsOptions are the options available to the Exists command.
type ExistsOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

// Exists checks if a document exists for the given id.
func (c *Collection) Exists(id string, opts *ExistsOptions) (docOut *ExistsResult, errOut error) {
	if opts == nil {
		opts = &ExistsOptions{}
	}

	opm := c.newKvOpManager("Exists", nil)
	defer opm.Finish()

	opm.SetDocumentID(id)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	agent, err := c.getKvProvider()
	if err != nil {
		return nil, err
	}
	err = opm.Wait(agent.GetMeta(gocbcore.GetMetaOptions{
		Key:            opm.DocumentID(),
		CollectionName: opm.CollectionName(),
		ScopeName:      opm.ScopeName(),
		RetryStrategy:  opm.RetryStrategy(),
		TraceContext:   opm.TraceSpan,
		Deadline:       opm.Deadline(),
	}, func(res *gocbcore.GetMetaResult, err error) {
		if errors.Is(err, ErrDocumentNotFound) {
			docOut = &ExistsResult{
				Result: Result{
					cas: Cas(0),
				},
				docExists: false,
			}
			opm.Resolve(nil)
			return
		}

		if err != nil {
			errOut = opm.EnhanceErr(err)
			opm.Reject()
			return
		}

		if res != nil {
			docOut = &ExistsResult{
				Result: Result{
					cas: Cas(res.Cas),
				},
				docExists: res.Deleted == 0,
			}
		}

		opm.Resolve(nil)
	}))
	if err != nil {
		errOut = err
	}
	return
}

func (c *Collection) getOneReplica(
	span requestSpanContext,
	id string,
	replicaIdx int,
	transcoder Transcoder,
	retryStrategy RetryStrategy,
	cancelCh chan struct{},
	timeout time.Duration,
) (docOut *GetReplicaResult, errOut error) {
	opm := c.newKvOpManager("getOneReplica", span)
	defer opm.Finish()

	opm.SetDocumentID(id)
	opm.SetTranscoder(transcoder)
	opm.SetRetryStrategy(retryStrategy)
	opm.SetTimeout(timeout)
	opm.SetCancelCh(cancelCh)

	agent, err := c.getKvProvider()
	if err != nil {
		return nil, err
	}
	if replicaIdx == 0 {
		err = opm.Wait(agent.Get(gocbcore.GetOptions{
			Key:            opm.DocumentID(),
			CollectionName: opm.CollectionName(),
			ScopeName:      opm.ScopeName(),
			RetryStrategy:  opm.RetryStrategy(),
			TraceContext:   opm.TraceSpan(),
			Deadline:       opm.Deadline(),
		}, func(res *gocbcore.GetResult, err error) {
			if err != nil {
				errOut = opm.EnhanceErr(err)
				opm.Reject()
				return
			}

			docOut = &GetReplicaResult{}
			docOut.cas = Cas(res.Cas)
			docOut.transcoder = opm.Transcoder()
			docOut.contents = res.Value
			docOut.flags = res.Flags
			docOut.isReplica = false

			opm.Resolve(nil)
		}))
		if err != nil {
			errOut = err
		}
		return
	}

	err = opm.Wait(agent.GetOneReplica(gocbcore.GetOneReplicaOptions{
		Key:            opm.DocumentID(),
		ReplicaIdx:     replicaIdx,
		CollectionName: opm.CollectionName(),
		ScopeName:      opm.ScopeName(),
		RetryStrategy:  opm.RetryStrategy(),
		TraceContext:   opm.TraceSpan(),
		Deadline:       opm.Deadline(),
	}, func(res *gocbcore.GetReplicaResult, err error) {
		if err != nil {
			errOut = opm.EnhanceErr(err)
			opm.Reject()
			return
		}

		docOut = &GetReplicaResult{}
		docOut.cas = Cas(res.Cas)
		docOut.transcoder = opm.Transcoder()
		docOut.contents = res.Value
		docOut.flags = res.Flags
		docOut.isReplica = true

		opm.Resolve(nil)
	}))
	if err != nil {
		errOut = err
	}
	return
}

// GetAllReplicaOptions are the options available to the GetAllReplicas command.
type GetAllReplicaOptions struct {
	Transcoder    Transcoder
	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

// GetAllReplicasResult represents the results of a GetAllReplicas operation.
type GetAllReplicasResult struct {
	lock          sync.Mutex
	totalRequests uint32
	totalResults  uint32
	resCh         chan *GetReplicaResult
	cancelCh      chan struct{}
}

func (r *GetAllReplicasResult) addResult(res *GetReplicaResult) {
	// We use a lock here because the alternative means that there is a race
	// between the channel writes from multiple results and the channels being
	// closed.  IE: T1-Incr, T2-Incr, T2-Send, T2-Close, T1-Send[PANIC]
	r.lock.Lock()

	r.totalResults++
	resultCount := r.totalResults

	if resultCount <= r.totalRequests {
		r.resCh <- res
	}

	if resultCount == r.totalRequests {
		close(r.cancelCh)
		close(r.resCh)
	}

	r.lock.Unlock()
}

// Next fetches the next replica result.
func (r *GetAllReplicasResult) Next() *GetReplicaResult {
	return <-r.resCh
}

// Close cancels all remaining get replica requests.
func (r *GetAllReplicasResult) Close() error {
	// See addResult discussion on lock usage.
	r.lock.Lock()

	// Note that this number increment must be high enough to be clear that
	// the result set was closed, but low enough that it won't overflow if
	// additional result objects are processed after the close.
	prevResultCount := r.totalResults
	r.totalResults += 100000

	// We only have to close everything if the addResult method didn't already
	// close them due to already having completed every request
	if prevResultCount < r.totalRequests {
		close(r.cancelCh)
		close(r.resCh)
	}

	r.lock.Unlock()

	return nil
}

// GetAllReplicas returns the value of a particular document from all replica servers. This will return an iterable
// which streams results one at a time.
func (c *Collection) GetAllReplicas(id string, opts *GetAllReplicaOptions) (docOut *GetAllReplicasResult, errOut error) {
	if opts == nil {
		opts = &GetAllReplicaOptions{}
	}

	span := c.startKvOpTrace("GetAllReplicas", nil)
	defer span.Finish()

	// Timeout needs to be adjusted here, since we use it at the bottom of this
	// function, but the remaining options are all passed downwards and get handled
	// by those functions rather than us.
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = c.timeoutsConfig.KVTimeout
	}

	deadline := time.Now().Add(timeout)
	transcoder := opts.Transcoder
	retryStrategy := opts.RetryStrategy

	agent, err := c.getKvProvider()
	if err != nil {
		return nil, err
	}

	snapshot, err := agent.ConfigSnapshot()
	if err != nil {
		return nil, err
	}

	numReplicas, err := snapshot.NumReplicas()
	if err != nil {
		return nil, err
	}

	numServers := numReplicas + 1
	outCh := make(chan *GetReplicaResult, numServers)
	cancelCh := make(chan struct{})

	repRes := &GetAllReplicasResult{
		totalRequests: uint32(numServers),
		resCh:         outCh,
		cancelCh:      cancelCh,
	}

	// Loop all the servers and populate the result object
	for replicaIdx := 0; replicaIdx < numServers; replicaIdx++ {
		go func(replicaIdx int) {
			// This timeout value will cause the getOneReplica operation to timeout after our deadline has expired,
			// as the deadline has already begun. getOneReplica timing out before our deadline would cause inconsistent
			// behaviour.
			res, err := c.getOneReplica(span, id, replicaIdx, transcoder, retryStrategy, cancelCh, timeout)
			if err != nil {
				logDebugf("Failed to fetch replica from replica %d: %s", replicaIdx, err)
			} else {
				repRes.addResult(res)
			}
		}(replicaIdx)
	}

	// Start a timer to close it after the deadline
	go func() {
		select {
		case <-time.After(time.Until(deadline)):
			// If we timeout, we should close the result
			err := repRes.Close()
			if err != nil {
				logDebugf("failed to close GetAllReplicas response: %s", err)
			}
			return
		case <-cancelCh:
			// If the cancel channel closes, we are done
			return
		}
	}()

	return repRes, nil
}

// GetAnyReplicaOptions are the options available to the GetAnyReplica command.
type GetAnyReplicaOptions struct {
	Transcoder    Transcoder
	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

// GetAnyReplica returns the value of a particular document from a replica server.
func (c *Collection) GetAnyReplica(id string, opts *GetAnyReplicaOptions) (docOut *GetReplicaResult, errOut error) {
	if opts == nil {
		opts = &GetAnyReplicaOptions{}
	}

	span := c.startKvOpTrace("GetAnyReplica", nil)
	defer span.Finish()

	repRes, err := c.GetAllReplicas(id, &GetAllReplicaOptions{
		Timeout:       opts.Timeout,
		Transcoder:    opts.Transcoder,
		RetryStrategy: opts.RetryStrategy,
	})
	if err != nil {
		return nil, err
	}

	// Try to fetch at least one result
	res := repRes.Next()
	if res == nil {
		return nil, &KeyValueError{
			InnerError:     ErrDocumentUnretrievable,
			BucketName:     c.bucketName(),
			ScopeName:      c.scope,
			CollectionName: c.collectionName,
		}
	}

	// Close the results channel since we don't care about any of the
	// remaining result objects at this point.
	err = repRes.Close()
	if err != nil {
		logDebugf("failed to close GetAnyReplica response: %s", err)
	}

	return res, nil
}

// RemoveOptions are the options available to the Remove command.
type RemoveOptions struct {
	Cas             Cas
	PersistTo       uint
	ReplicateTo     uint
	DurabilityLevel DurabilityLevel
	Timeout         time.Duration
	RetryStrategy   RetryStrategy
}

// Remove removes a document from the collection.
func (c *Collection) Remove(id string, opts *RemoveOptions) (mutOut *MutationResult, errOut error) {
	if opts == nil {
		opts = &RemoveOptions{}
	}

	opm := c.newKvOpManager("Remove", nil)
	defer opm.Finish()

	opm.SetDocumentID(id)
	opm.SetDuraOptions(opts.PersistTo, opts.ReplicateTo, opts.DurabilityLevel)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	agent, err := c.getKvProvider()
	if err != nil {
		return nil, err
	}
	err = opm.Wait(agent.Delete(gocbcore.DeleteOptions{
		Key:                    opm.DocumentID(),
		Cas:                    gocbcore.Cas(opts.Cas),
		CollectionName:         opm.CollectionName(),
		ScopeName:              opm.ScopeName(),
		DurabilityLevel:        opm.DurabilityLevel(),
		DurabilityLevelTimeout: opm.DurabilityTimeout(),
		RetryStrategy:          opm.RetryStrategy(),
		TraceContext:           opm.TraceSpan(),
		Deadline:               opm.Deadline(),
	}, func(res *gocbcore.DeleteResult, err error) {
		if err != nil {
			errOut = opm.EnhanceErr(err)
			opm.Reject()
			return
		}

		mutOut = &MutationResult{}
		mutOut.cas = Cas(res.Cas)
		mutOut.mt = opm.EnhanceMt(res.MutationToken)

		opm.Resolve(mutOut.mt)
	}))
	if err != nil {
		errOut = err
	}
	return
}

// GetAndTouchOptions are the options available to the GetAndTouch operation.
type GetAndTouchOptions struct {
	Transcoder    Transcoder
	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

// GetAndTouch retrieves a document and simultaneously updates its expiry time.
func (c *Collection) GetAndTouch(id string, expiry time.Duration, opts *GetAndTouchOptions) (docOut *GetResult, errOut error) {
	if opts == nil {
		opts = &GetAndTouchOptions{}
	}

	opm := c.newKvOpManager("GetAndTouch", nil)
	defer opm.Finish()

	opm.SetDocumentID(id)
	opm.SetTranscoder(opts.Transcoder)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	agent, err := c.getKvProvider()
	if err != nil {
		return nil, err
	}
	err = opm.Wait(agent.GetAndTouch(gocbcore.GetAndTouchOptions{
		Key:            opm.DocumentID(),
		Expiry:         durationToExpiry(expiry),
		CollectionName: opm.CollectionName(),
		ScopeName:      opm.ScopeName(),
		RetryStrategy:  opm.RetryStrategy(),
		TraceContext:   opm.TraceSpan(),
		Deadline:       opm.Deadline(),
	}, func(res *gocbcore.GetAndTouchResult, err error) {
		if err != nil {
			errOut = opm.EnhanceErr(err)
			opm.Reject()
			return
		}

		if res != nil {
			doc := &GetResult{
				Result: Result{
					cas: Cas(res.Cas),
				},
				transcoder: opm.Transcoder(),
				contents:   res.Value,
				flags:      res.Flags,
			}

			docOut = doc
		}

		opm.Resolve(nil)
	}))
	if err != nil {
		errOut = err
	}
	return
}

// GetAndLockOptions are the options available to the GetAndLock operation.
type GetAndLockOptions struct {
	Transcoder    Transcoder
	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

// GetAndLock locks a document for a period of time, providing exclusive RW access to it.
// A lockTime value of over 30 seconds will be treated as 30 seconds. The resolution used to send this value to
// the server is seconds and is calculated using uint32(lockTime/time.Second).
func (c *Collection) GetAndLock(id string, lockTime time.Duration, opts *GetAndLockOptions) (docOut *GetResult, errOut error) {
	if opts == nil {
		opts = &GetAndLockOptions{}
	}

	opm := c.newKvOpManager("GetAndLock", nil)
	defer opm.Finish()

	opm.SetDocumentID(id)
	opm.SetTranscoder(opts.Transcoder)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	agent, err := c.getKvProvider()
	if err != nil {
		return nil, err
	}
	err = opm.Wait(agent.GetAndLock(gocbcore.GetAndLockOptions{
		Key:            opm.DocumentID(),
		LockTime:       uint32(lockTime / time.Second),
		CollectionName: opm.CollectionName(),
		ScopeName:      opm.ScopeName(),
		RetryStrategy:  opm.RetryStrategy(),
		TraceContext:   opm.TraceSpan(),
		Deadline:       opm.Deadline(),
	}, func(res *gocbcore.GetAndLockResult, err error) {
		if err != nil {
			errOut = opm.EnhanceErr(err)
			opm.Reject()
			return
		}

		if res != nil {
			doc := &GetResult{
				Result: Result{
					cas: Cas(res.Cas),
				},
				transcoder: opm.Transcoder(),
				contents:   res.Value,
				flags:      res.Flags,
			}

			docOut = doc
		}

		opm.Resolve(nil)
	}))
	if err != nil {
		errOut = err
	}
	return
}

// UnlockOptions are the options available to the GetAndLock operation.
type UnlockOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

// Unlock unlocks a document which was locked with GetAndLock.
func (c *Collection) Unlock(id string, cas Cas, opts *UnlockOptions) (errOut error) {
	if opts == nil {
		opts = &UnlockOptions{}
	}

	opm := c.newKvOpManager("Unlock", nil)
	defer opm.Finish()

	opm.SetDocumentID(id)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)

	if err := opm.CheckReadyForOp(); err != nil {
		return err
	}

	agent, err := c.getKvProvider()
	if err != nil {
		return err
	}
	err = opm.Wait(agent.Unlock(gocbcore.UnlockOptions{
		Key:            opm.DocumentID(),
		Cas:            gocbcore.Cas(cas),
		CollectionName: opm.CollectionName(),
		ScopeName:      opm.ScopeName(),
		RetryStrategy:  opm.RetryStrategy(),
		TraceContext:   opm.TraceSpan(),
		Deadline:       opm.Deadline(),
	}, func(res *gocbcore.UnlockResult, err error) {
		if err != nil {
			errOut = opm.EnhanceErr(err)
			opm.Reject()
			return
		}

		mt := opm.EnhanceMt(res.MutationToken)
		opm.Resolve(mt)
	}))
	if err != nil {
		errOut = err
	}
	return
}

// TouchOptions are the options available to the Touch operation.
type TouchOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
}

// Touch touches a document, specifying a new expiry time for it.
func (c *Collection) Touch(id string, expiry time.Duration, opts *TouchOptions) (mutOut *MutationResult, errOut error) {
	if opts == nil {
		opts = &TouchOptions{}
	}

	opm := c.newKvOpManager("Touch", nil)
	defer opm.Finish()

	opm.SetDocumentID(id)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	agent, err := c.getKvProvider()
	if err != nil {
		return nil, err
	}
	err = opm.Wait(agent.Touch(gocbcore.TouchOptions{
		Key:            opm.DocumentID(),
		Expiry:         durationToExpiry(expiry),
		CollectionName: opm.CollectionName(),
		ScopeName:      opm.ScopeName(),
		RetryStrategy:  opm.RetryStrategy(),
		TraceContext:   opm.TraceSpan(),
		Deadline:       opm.Deadline(),
	}, func(res *gocbcore.TouchResult, err error) {
		if err != nil {
			errOut = opm.EnhanceErr(err)
			opm.Reject()
			return
		}

		mutOut = &MutationResult{}
		mutOut.cas = Cas(res.Cas)
		mutOut.mt = opm.EnhanceMt(res.MutationToken)

		opm.Resolve(mutOut.mt)
	}))
	if err != nil {
		errOut = err
	}
	return
}

// Binary creates and returns a BinaryCollection object.
func (c *Collection) Binary() *BinaryCollection {
	return &BinaryCollection{collection: c}
}
