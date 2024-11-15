package gocb

import (
	"context"
	"errors"
	"time"

	"github.com/couchbase/gocbcore/v10"
)

type kvProviderCore struct {
	agent            kvProviderCoreProvider
	snapshotProvider kvProviderConfigSnapshotProvider

	tracer               RequestTracer
	meter                *meterWrapper
	preferredServerGroup string
}

var _ kvProvider = &kvProviderCore{}

func (p *kvProviderCore) Scan(c *Collection, scanType ScanType, opts *ScanOptions) (*ScanResult, error) {
	opm, err := p.newRangeScanOpManager(c, scanType, p.agent, opts.ParentSpan, opts.ConsistentWith,
		opts.IDsOnly)
	if err != nil {
		return nil, err
	}

	opm.SetTranscoder(opts.Transcoder)
	opm.SetImpersonate(opts.Internal.User)
	opm.SetTimeout(opts.Timeout)
	opm.SetItemLimit(opts.BatchItemLimit)
	opm.SetByteLimit(opts.BatchByteLimit)
	opm.SetMaxConcurrency(opts.Concurrency)

	config, err := p.snapshotProvider.WaitForConfigSnapshot(opts.Context, time.Now().Add(opm.Timeout()))
	if err != nil {
		opm.Finish()
		return nil, maybeEnhanceKVErr(err, c.bucketName(), c.ScopeName(), c.Name(), "scan")
	}

	numVbuckets, err := config.NumVbuckets()
	if err != nil {
		opm.Finish()
		return nil, err
	}

	if numVbuckets == 0 {
		opm.Finish()
		return nil, makeInvalidArgumentsError("can only use Scan with couchbase buckets")
	}

	opm.SetNumVbuckets(numVbuckets)

	serverToVbucketMap := make(map[int][]uint16)
	numServers, err := config.NumServers()
	if err != nil {
		opm.Finish()
		return nil, err
	}
	for serverIndex := 0; serverIndex < numServers; serverIndex++ {
		vbuckets, err := config.VbucketsOnServer(serverIndex)
		if err != nil {
			opm.Finish()
			return nil, err
		}
		serverToVbucketMap[serverIndex] = vbuckets
	}
	opm.SetServerToVbucketMap(serverToVbucketMap)

	cid, err := p.getCollectionID(opts.Context, c, opm.TraceSpan(), opm.Timeout(), opm.Impersonate())
	if err != nil {
		opm.Finish()
		return nil, maybeEnhanceKVErr(err, c.bucketName(), c.ScopeName(), c.Name(), "scan")
	}

	opm.SetCollectionID(cid)

	if err := opm.CheckReadyForOp(); err != nil {
		opm.Finish()
		return nil, err
	}

	return opm.Scan(opts.Context)
}

func (p *kvProviderCore) getCollectionID(ctx context.Context, c *Collection, parentSpan RequestSpan, timeout time.Duration,
	impersonate string) (uint32, error) {
	if c.isDefault() {
		return 0, nil
	}

	opm := newKvOpManagerCore(c, "get_collection_id", parentSpan, p)
	defer opm.Finish(false)

	opm.SetTimeout(timeout)
	opm.SetImpersonate(impersonate)
	opm.SetContext(ctx)

	if err := opm.CheckReadyForOp(); err != nil {
		return 0, err
	}

	var errOut error
	var cidOut uint32
	err := opm.Wait(p.agent.GetCollectionID(c.ScopeName(), c.name(), gocbcore.GetCollectionIDOptions{
		TraceContext: opm.TraceSpanContext(),
		Deadline:     opm.Deadline(),
		// User:         opm.Impersonate(),
	}, func(res *gocbcore.GetCollectionIDResult, err error) {
		if err != nil {
			errOut = opm.EnhanceErr(err)
			opm.Reject()
			return
		}

		cidOut = res.CollectionID

		opm.Resolve(nil)
	}))
	if err != nil {
		errOut = err
	}

	return cidOut, errOut
}

func (p *kvProviderCore) Insert(c *Collection, id string, val interface{}, opts *InsertOptions) (*MutationResult, error) {
	opm := newKvOpManagerCore(c, "insert", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetTranscoder(opts.Transcoder)
	opm.SetValue(val)
	opm.SetDuraOptions(opts.PersistTo, opts.ReplicateTo, opts.DurabilityLevel)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetImpersonate(opts.Internal.User)
	opm.SetContext(opts.Context)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	var errOut error
	var mutOut *MutationResult
	err := opm.Wait(p.agent.Add(gocbcore.AddOptions{
		Key:                    opm.DocumentID(),
		Value:                  opm.ValueBytes(),
		Flags:                  opm.ValueFlags(),
		Expiry:                 durationToExpiry(opts.Expiry),
		CollectionName:         opm.CollectionName(),
		ScopeName:              opm.ScopeName(),
		DurabilityLevel:        opm.DurabilityLevel(),
		DurabilityLevelTimeout: opm.DurabilityTimeout(),
		RetryStrategy:          opm.RetryStrategy(),
		TraceContext:           opm.TraceSpanContext(),
		Deadline:               opm.Deadline(),
		User:                   opm.Impersonate(),
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

	return mutOut, errOut
}

func (p *kvProviderCore) Upsert(c *Collection, id string, val interface{}, opts *UpsertOptions) (*MutationResult, error) {
	opm := newKvOpManagerCore(c, "upsert", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetTranscoder(opts.Transcoder)
	opm.SetValue(val)
	opm.SetDuraOptions(opts.PersistTo, opts.ReplicateTo, opts.DurabilityLevel)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetImpersonate(opts.Internal.User)
	opm.SetContext(opts.Context)
	opm.SetPreserveExpiry(opts.PreserveExpiry)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	var errOut error
	var mutOut *MutationResult
	err := opm.Wait(p.agent.Set(gocbcore.SetOptions{
		Key:                    opm.DocumentID(),
		Value:                  opm.ValueBytes(),
		Flags:                  opm.ValueFlags(),
		Expiry:                 durationToExpiry(opts.Expiry),
		CollectionName:         opm.CollectionName(),
		ScopeName:              opm.ScopeName(),
		DurabilityLevel:        opm.DurabilityLevel(),
		DurabilityLevelTimeout: opm.DurabilityTimeout(),
		RetryStrategy:          opm.RetryStrategy(),
		TraceContext:           opm.TraceSpanContext(),
		Deadline:               opm.Deadline(),
		User:                   opm.Impersonate(),
		PreserveExpiry:         opm.PreserveExpiry(),
	}, func(sr *gocbcore.StoreResult, err error) {
		if err != nil {
			errOut = opm.EnhanceErr(err)
			opm.Reject()
			return
		}

		mutOut = &MutationResult{}
		mutOut.cas = Cas(sr.Cas)
		mutOut.mt = opm.EnhanceMt(sr.MutationToken)

		opm.Resolve(mutOut.mt)
	}))
	if err != nil {
		errOut = err
	}
	return mutOut, errOut
}

func (p *kvProviderCore) Replace(c *Collection, id string, val interface{}, opts *ReplaceOptions) (*MutationResult, error) {
	opm := newKvOpManagerCore(c, "replace", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetTranscoder(opts.Transcoder)
	opm.SetValue(val)
	opm.SetDuraOptions(opts.PersistTo, opts.ReplicateTo, opts.DurabilityLevel)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetImpersonate(opts.Internal.User)
	opm.SetContext(opts.Context)
	opm.SetPreserveExpiry(opts.PreserveExpiry)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	var errOut error
	var mutOut *MutationResult
	err := opm.Wait(p.agent.Replace(gocbcore.ReplaceOptions{
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
		TraceContext:           opm.TraceSpanContext(),
		Deadline:               opm.Deadline(),
		User:                   opm.Impersonate(),
		PreserveExpiry:         opm.PreserveExpiry(),
	}, func(sr *gocbcore.StoreResult, err error) {
		if err != nil {
			errOut = opm.EnhanceErr(err)
			opm.Reject()
			return
		}

		mutOut = &MutationResult{}
		mutOut.cas = Cas(sr.Cas)
		mutOut.mt = opm.EnhanceMt(sr.MutationToken)

		opm.Resolve(mutOut.mt)

	}))
	if err != nil {
		errOut = err
	}

	return mutOut, errOut
}

func (p *kvProviderCore) Get(c *Collection, id string, opts *GetOptions) (*GetResult, error) {
	if opts == nil {
		opts = &GetOptions{}
	}

	if len(opts.Project) == 0 && !opts.WithExpiry {
		return p.getDirect(c, id, opts)
	}

	return p.getProjected(c, id, opts)
}

func (p *kvProviderCore) getDirect(c *Collection, id string, opts *GetOptions) (*GetResult, error) {
	if opts == nil {
		opts = &GetOptions{}
	}

	opm := newKvOpManagerCore(c, "get", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetTranscoder(opts.Transcoder)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetImpersonate(opts.Internal.User)
	opm.SetContext(opts.Context)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	var docOut *GetResult
	var errOut error
	err := opm.Wait(p.agent.Get(gocbcore.GetOptions{
		Key:            opm.DocumentID(),
		CollectionName: opm.CollectionName(),
		ScopeName:      opm.ScopeName(),
		RetryStrategy:  opm.RetryStrategy(),
		TraceContext:   opm.TraceSpanContext(),
		Deadline:       opm.Deadline(),
		User:           opm.Impersonate(),
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

	return docOut, errOut
}

func (p *kvProviderCore) getProjected(c *Collection, id string, opts *GetOptions) (docOut *GetResult, errOut error) {
	if opts == nil {
		opts = &GetOptions{}
	}

	opm := newKvOpManagerCore(c, "get", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetTranscoder(opts.Transcoder)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetImpersonate(opts.Internal.User)
	opm.SetContext(opts.Context)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	var withFlags bool
	numProjects := len(opts.Project)
	if opts.WithExpiry {
		if numProjects == 0 {
			// This must be a full get with expiry
			withFlags = true
		}
		numProjects = 1 + numProjects
	}

	projections := opts.Project
	if numProjects > 16 {
		projections = nil
	}

	var ops []LookupInSpec

	if opts.WithExpiry {
		ops = append(ops, GetSpec("$document.exptime", &GetSpecOptions{IsXattr: true}))

		if withFlags {
			// We also need to fetch the flags, we need them for transcoding and they aren't included in a lookupin
			// response. We only need these when doing a full get with expiry.
			ops = append(ops, GetSpec("$document.flags", &GetSpecOptions{IsXattr: true}))
		}
	}

	if len(projections) == 0 {
		ops = append(ops, GetSpec("", nil))
	} else {
		for _, path := range projections {
			ops = append(ops, GetSpec(path, nil))
		}
	}

	result, err := p.LookupIn(c, id, ops, &LookupInOptions{
		ParentSpan: opm.TraceSpan(),
		noMetrics:  true,
		Context:    opts.Context,
	})
	if err != nil {
		return nil, err
	}

	doc := &GetResult{}
	if opts.WithExpiry {
		// if expiration was requested then extract and remove it from the results
		var expires int64
		err = result.ContentAt(0, &expires)
		if err != nil {
			return nil, err
		}

		var expiryTime time.Time
		if expires > 0 {
			expiryTime = time.Unix(expires, 0)
		}
		doc.expiryTime = &expiryTime

		ops = ops[1:]
		result.contents = result.contents[1:]

		if withFlags {
			var flags uint32
			err = result.ContentAt(0, &flags)
			if err != nil {
				return nil, err
			}

			doc.flags = flags

			ops = ops[1:]
			result.contents = result.contents[1:]
		}
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

func (p *kvProviderCore) GetAndTouch(c *Collection, id string, expiry time.Duration, opts *GetAndTouchOptions) (*GetResult, error) {
	opm := newKvOpManagerCore(c, "get_and_touch", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetTranscoder(opts.Transcoder)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetImpersonate(opts.Internal.User)
	opm.SetContext(opts.Context)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	var getOut *GetResult
	var errOut error
	err := opm.Wait(p.agent.GetAndTouch(gocbcore.GetAndTouchOptions{
		Key:            opm.DocumentID(),
		Expiry:         durationToExpiry(expiry),
		CollectionName: opm.CollectionName(),
		ScopeName:      opm.ScopeName(),
		RetryStrategy:  opm.RetryStrategy(),
		TraceContext:   opm.TraceSpanContext(),
		Deadline:       opm.Deadline(),
		User:           opm.Impersonate(),
	}, func(res *gocbcore.GetAndTouchResult, err error) {
		if err != nil {
			errOut = opm.EnhanceErr(err)
			opm.Reject()
			return
		}

		if res != nil {
			getOut = &GetResult{
				Result: Result{
					cas: Cas(res.Cas),
				},
				transcoder: opm.Transcoder(),
				contents:   res.Value,
				flags:      res.Flags,
			}
		}

		opm.Resolve(nil)
	}))
	if err != nil {
		errOut = err
	}

	return getOut, errOut

}

func (p *kvProviderCore) GetAndLock(c *Collection, id string, lockTime time.Duration, opts *GetAndLockOptions) (*GetResult, error) {
	opm := newKvOpManagerCore(c, "get_and_lock", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetTranscoder(opts.Transcoder)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetImpersonate(opts.Internal.User)
	opm.SetContext(opts.Context)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	var getOut *GetResult
	var errOut error
	err := opm.Wait(p.agent.GetAndLock(gocbcore.GetAndLockOptions{
		Key:            opm.DocumentID(),
		LockTime:       uint32(lockTime / time.Second),
		CollectionName: opm.CollectionName(),
		ScopeName:      opm.ScopeName(),
		RetryStrategy:  opm.RetryStrategy(),
		TraceContext:   opm.TraceSpanContext(),
		Deadline:       opm.Deadline(),
		User:           opm.Impersonate(),
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

			getOut = doc
		}

		opm.Resolve(nil)
	}))
	if err != nil {
		errOut = err
	}

	return getOut, errOut
}

func (p *kvProviderCore) Exists(c *Collection, id string, opts *ExistsOptions) (*ExistsResult, error) {
	opm := newKvOpManagerCore(c, "exists", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetImpersonate(opts.Internal.User)
	opm.SetContext(opts.Context)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	var docExists *ExistsResult
	var errOut error
	err := opm.Wait(p.agent.GetMeta(gocbcore.GetMetaOptions{
		Key:            opm.DocumentID(),
		CollectionName: opm.CollectionName(),
		ScopeName:      opm.ScopeName(),
		RetryStrategy:  opm.RetryStrategy(),
		TraceContext:   opm.TraceSpanContext(),
		Deadline:       opm.Deadline(),
		User:           opm.Impersonate(),
	}, func(res *gocbcore.GetMetaResult, err error) {
		if errors.Is(err, ErrDocumentNotFound) {
			docExists = &ExistsResult{
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
			docExists = &ExistsResult{
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

	return docExists, errOut
}

func (p *kvProviderCore) Remove(c *Collection, id string, opts *RemoveOptions) (*MutationResult, error) {
	opm := newKvOpManagerCore(c, "remove", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetDuraOptions(opts.PersistTo, opts.ReplicateTo, opts.DurabilityLevel)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetImpersonate(opts.Internal.User)
	opm.SetContext(opts.Context)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	var errOut error
	var mutOut *MutationResult
	err := opm.Wait(p.agent.Delete(gocbcore.DeleteOptions{
		Key:                    opm.DocumentID(),
		Cas:                    gocbcore.Cas(opts.Cas),
		CollectionName:         opm.CollectionName(),
		ScopeName:              opm.ScopeName(),
		DurabilityLevel:        opm.DurabilityLevel(),
		DurabilityLevelTimeout: opm.DurabilityTimeout(),
		RetryStrategy:          opm.RetryStrategy(),
		TraceContext:           opm.TraceSpanContext(),
		Deadline:               opm.Deadline(),
		User:                   opm.Impersonate(),
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

	return mutOut, errOut
}

func (p *kvProviderCore) Unlock(c *Collection, id string, cas Cas, opts *UnlockOptions) error {
	opm := newKvOpManagerCore(c, "unlock", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetImpersonate(opts.Internal.User)
	opm.SetContext(opts.Context)

	if err := opm.CheckReadyForOp(); err != nil {
		return err
	}

	var errOut error
	err := opm.Wait(p.agent.Unlock(gocbcore.UnlockOptions{
		Key:            opm.DocumentID(),
		Cas:            gocbcore.Cas(cas),
		CollectionName: opm.CollectionName(),
		ScopeName:      opm.ScopeName(),
		RetryStrategy:  opm.RetryStrategy(),
		TraceContext:   opm.TraceSpanContext(),
		Deadline:       opm.Deadline(),
		User:           opm.Impersonate(),
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
	return errOut
}

func (p *kvProviderCore) Touch(c *Collection, id string, expiry time.Duration, opts *TouchOptions) (*MutationResult, error) {
	opm := newKvOpManagerCore(c, "touch", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetImpersonate(opts.Internal.User)
	opm.SetContext(opts.Context)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	var errOut error
	var mutOut *MutationResult
	err := opm.Wait(p.agent.Touch(gocbcore.TouchOptions{
		Key:            opm.DocumentID(),
		Expiry:         durationToExpiry(expiry),
		CollectionName: opm.CollectionName(),
		ScopeName:      opm.ScopeName(),
		RetryStrategy:  opm.RetryStrategy(),
		TraceContext:   opm.TraceSpanContext(),
		Deadline:       opm.Deadline(),
		User:           opm.Impersonate(),
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

	return mutOut, errOut

}

func (p *kvProviderCore) GetAllReplicas(c *Collection, id string, opts *GetAllReplicaOptions) (*GetAllReplicasResult, error) {
	if opts == nil {
		opts = &GetAllReplicaOptions{}
	}

	var tracectx RequestSpanContext
	if opts.ParentSpan != nil {
		tracectx = opts.ParentSpan.Context()
	}

	ctx := opts.Context
	if ctx == nil {
		ctx = context.Background()
	}

	span := p.StartKvOpTrace(c, "get_all_replicas", tracectx, false)

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

	snapshot, err := p.snapshotProvider.WaitForConfigSnapshot(ctx, deadline)
	if err != nil {
		return nil, err
	}

	var servers []int
	if opts.ReadPreference == ReadPreferenceSelectedServerGroup {
		serverGroups, err := snapshot.KeyToServersByServerGroup([]byte(id))
		if err != nil {
			return nil, err
		}

		for group, srvIdx := range serverGroups {
			if group == p.preferredServerGroup {
				servers = append(servers, srvIdx...)
			}
		}
	} else {
		numReplicas, err := snapshot.NumReplicas()
		if err != nil {
			return nil, err
		}

		numServers := numReplicas + 1
		for i := 0; i < numServers; i++ {
			servers = append(servers, i)
		}
	}

	outCh := make(chan interface{}, len(servers))
	cancelCh := make(chan struct{})

	var recorder ValueRecorder
	if !opts.noMetrics {
		recorder, err = p.meter.ValueRecorder(meterValueServiceKV, "get_all_replicas")
		if err != nil {
			logDebugf("Failed to create value recorder: %v", err)
		}
	}

	coreRes := &coreReplicasResult{
		totalRequests:       uint32(len(servers)),
		resCh:               outCh,
		cancelCh:            cancelCh,
		span:                span,
		childReqsCompleteCh: make(chan struct{}),
		valueRecorder:       recorder,
		startedTime:         time.Now(),
	}
	repRes := &GetAllReplicasResult{
		res: coreRes,
	}

	if len(servers) == 0 {
		// This can happen when the selected server group does not exist, or has not been set
		close(coreRes.resCh)
		close(coreRes.childReqsCompleteCh)
		close(coreRes.cancelCh)
		coreRes.span.End()

		return repRes, nil
	}

	// Loop all the servers and populate the result object
	for _, replicaIdx := range servers {
		go func(replicaIdx int) {
			// This timeout value will cause the getOneReplica operation to timeout after our deadline has expired,
			// as the deadline has already begun. getOneReplica timing out before our deadline would cause inconsistent
			// behaviour.
			res, err := p.getOneReplica(context.Background(), span, id, replicaIdx, transcoder, retryStrategy, cancelCh,
				timeout, opts.Internal.User, c)
			if err != nil {
				coreRes.addFailed()
				logDebugf("Failed to fetch replica from replica %d: %s", replicaIdx, err)
			} else {
				coreRes.addResult(res)
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
		case <-cancelCh:
		// If the cancel channel closes, we are done
		case <-ctx.Done():
			err := repRes.Close()
			if err != nil {
				logDebugf("failed to close GetAllReplicas response: %s", err)
			}
		}
	}()

	return repRes, nil
}

func (p *kvProviderCore) GetAnyReplica(c *Collection, id string, opts *GetAnyReplicaOptions) (*GetReplicaResult, error) {
	if opts == nil {
		opts = &GetAnyReplicaOptions{}
	}

	start := time.Now()
	defer p.meter.ValueRecord("kv", "get_any_replica", start)

	var tracectx RequestSpanContext
	if opts.ParentSpan != nil {
		tracectx = opts.ParentSpan.Context()
	}

	span := p.StartKvOpTrace(c, "get_any_replica", tracectx, false)
	defer span.End()

	repRes, err := p.GetAllReplicas(c, id, &GetAllReplicaOptions{
		Timeout:        opts.Timeout,
		Transcoder:     opts.Transcoder,
		RetryStrategy:  opts.RetryStrategy,
		Internal:       opts.Internal,
		ParentSpan:     span,
		noMetrics:      true,
		Context:        opts.Context,
		ReadPreference: opts.ReadPreference,
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

func (p *kvProviderCore) getOneReplica(
	ctx context.Context,
	span RequestSpan,
	id string,
	replicaIdx int,
	transcoder Transcoder,
	retryStrategy RetryStrategy,
	cancelCh chan struct{},
	timeout time.Duration,
	user string,
	c *Collection,
) (*GetReplicaResult, error) {
	opm := newKvOpManagerCore(c, "get_replica", span, p)
	defer opm.Finish(true)

	opm.SetDocumentID(id)
	opm.SetTranscoder(transcoder)
	opm.SetRetryStrategy(retryStrategy)
	opm.SetTimeout(timeout)
	opm.SetCancelCh(cancelCh)
	opm.SetImpersonate(user)
	opm.SetContext(ctx)

	if replicaIdx == 0 {
		var docOut *GetReplicaResult
		var errOut error
		err := opm.Wait(p.agent.Get(gocbcore.GetOptions{
			Key:            opm.DocumentID(),
			CollectionName: opm.CollectionName(),
			ScopeName:      opm.ScopeName(),
			RetryStrategy:  opm.RetryStrategy(),
			TraceContext:   opm.TraceSpanContext(),
			Deadline:       opm.Deadline(),
			User:           opm.Impersonate(),
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
		return docOut, errOut
	}

	var docOut *GetReplicaResult
	var errOut error
	err := opm.Wait(p.agent.GetOneReplica(gocbcore.GetOneReplicaOptions{
		Key:            opm.DocumentID(),
		ReplicaIdx:     replicaIdx,
		CollectionName: opm.CollectionName(),
		ScopeName:      opm.ScopeName(),
		RetryStrategy:  opm.RetryStrategy(),
		TraceContext:   opm.TraceSpanContext(),
		Deadline:       opm.Deadline(),
		User:           opm.Impersonate(),
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
	return docOut, errOut
}

func (p *kvProviderCore) Prepend(c *Collection, id string, val []byte, opts *PrependOptions) (*MutationResult, error) {
	opm := newKvOpManagerCore(c, "prepend", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetDuraOptions(opts.PersistTo, opts.ReplicateTo, opts.DurabilityLevel)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetImpersonate(opts.Internal.User)
	opm.SetContext(opts.Context)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	var errOut error
	var mutOut *MutationResult
	err := opm.Wait(p.agent.Prepend(gocbcore.AdjoinOptions{
		Key:                    opm.DocumentID(),
		Value:                  val,
		CollectionName:         opm.CollectionName(),
		ScopeName:              opm.ScopeName(),
		DurabilityLevel:        opm.DurabilityLevel(),
		DurabilityLevelTimeout: opm.DurabilityTimeout(),
		Cas:                    gocbcore.Cas(opts.Cas),
		RetryStrategy:          opm.RetryStrategy(),
		TraceContext:           opm.TraceSpanContext(),
		Deadline:               opm.Deadline(),
		User:                   opm.Impersonate(),
	}, func(res *gocbcore.AdjoinResult, err error) {
		if err != nil {
			var coreErr *gocbcore.KeyValueError
			// We do not have a ErrNotStored error and this error should be translated
			// to ErrDocumentNotFound. At time of writing this would be a breaking change
			// in gocbcore so we do it here instead.
			if errors.Is(err, gocbcore.ErrNotStored) && errors.As(err, &coreErr) {
				coreErr.InnerError = ErrDocumentNotFound
				err = coreErr
			}
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

	return mutOut, errOut
}

func (p *kvProviderCore) Append(c *Collection, id string, val []byte, opts *AppendOptions) (*MutationResult, error) {
	opm := newKvOpManagerCore(c, "append", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetDuraOptions(opts.PersistTo, opts.ReplicateTo, opts.DurabilityLevel)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetImpersonate(opts.Internal.User)
	opm.SetContext(opts.Context)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	var errOut error
	var mutOut *MutationResult
	err := opm.Wait(p.agent.Append(gocbcore.AdjoinOptions{
		Key:                    opm.DocumentID(),
		Value:                  val,
		CollectionName:         opm.CollectionName(),
		ScopeName:              opm.ScopeName(),
		DurabilityLevel:        opm.DurabilityLevel(),
		DurabilityLevelTimeout: opm.DurabilityTimeout(),
		Cas:                    gocbcore.Cas(opts.Cas),
		RetryStrategy:          opm.RetryStrategy(),
		TraceContext:           opm.TraceSpanContext(),
		Deadline:               opm.Deadline(),
		User:                   opm.Impersonate(),
	}, func(res *gocbcore.AdjoinResult, err error) {
		if err != nil {
			var coreErr *gocbcore.KeyValueError
			// We do not have a ErrNotStored error and this error should be translated
			// to ErrDocumentNotFound. At time of writing this would be a breaking change
			// in gocbcore so we do it here instead.
			if errors.Is(err, gocbcore.ErrNotStored) && errors.As(err, &coreErr) {
				coreErr.InnerError = ErrDocumentNotFound
				err = coreErr
			}
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

	return mutOut, errOut
}

func (p *kvProviderCore) Increment(c *Collection, id string, opts *IncrementOptions) (*CounterResult, error) {
	if opts.Cas > 0 {
		return nil, makeInvalidArgumentsError("cas is not supported by the server for the Increment operation")
	}
	opm := newKvOpManagerCore(c, "increment", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetDuraOptions(opts.PersistTo, opts.ReplicateTo, opts.DurabilityLevel)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetImpersonate(opts.Internal.User)
	opm.SetContext(opts.Context)

	realInitial := uint64(0xFFFFFFFFFFFFFFFF)
	if opts.Initial >= 0 {
		realInitial = uint64(opts.Initial)
	}

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	var errOut error
	var countOut *CounterResult
	err := opm.Wait(p.agent.Increment(gocbcore.CounterOptions{
		Key:                    opm.DocumentID(),
		Delta:                  opts.Delta,
		Initial:                realInitial,
		Expiry:                 durationToExpiry(opts.Expiry),
		CollectionName:         opm.CollectionName(),
		ScopeName:              opm.ScopeName(),
		DurabilityLevel:        opm.DurabilityLevel(),
		DurabilityLevelTimeout: opm.DurabilityTimeout(),
		RetryStrategy:          opm.RetryStrategy(),
		TraceContext:           opm.TraceSpanContext(),
		Deadline:               opm.Deadline(),
		User:                   opm.Impersonate(),
	}, func(res *gocbcore.CounterResult, err error) {
		if err != nil {
			errOut = opm.EnhanceErr(err)
			opm.Reject()
			return
		}

		countOut = &CounterResult{}
		countOut.cas = Cas(res.Cas)
		countOut.mt = opm.EnhanceMt(res.MutationToken)
		countOut.content = res.Value

		opm.Resolve(countOut.mt)
	}))
	if err != nil {
		errOut = err
	}

	return countOut, errOut

}

func (p *kvProviderCore) Decrement(c *Collection, id string, opts *DecrementOptions) (*CounterResult, error) {
	if opts.Cas > 0 {
		return nil, makeInvalidArgumentsError("cas is not supported by the server for the Decrement operation")
	}
	opm := newKvOpManagerCore(c, "decrement", opts.ParentSpan, p)
	defer opm.Finish(false)

	opm.SetDocumentID(id)
	opm.SetDuraOptions(opts.PersistTo, opts.ReplicateTo, opts.DurabilityLevel)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetTimeout(opts.Timeout)
	opm.SetImpersonate(opts.Internal.User)
	opm.SetContext(opts.Context)

	realInitial := uint64(0xFFFFFFFFFFFFFFFF)
	if opts.Initial >= 0 {
		realInitial = uint64(opts.Initial)
	}

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	var errOut error
	var countOut *CounterResult

	err := opm.Wait(p.agent.Decrement(gocbcore.CounterOptions{
		Key:                    opm.DocumentID(),
		Delta:                  opts.Delta,
		Initial:                realInitial,
		Expiry:                 durationToExpiry(opts.Expiry),
		CollectionName:         opm.CollectionName(),
		ScopeName:              opm.ScopeName(),
		DurabilityLevel:        opm.DurabilityLevel(),
		DurabilityLevelTimeout: opm.DurabilityTimeout(),
		RetryStrategy:          opm.RetryStrategy(),
		TraceContext:           opm.TraceSpanContext(),
		Deadline:               opm.Deadline(),
		User:                   opm.Impersonate(),
	}, func(res *gocbcore.CounterResult, err error) {
		if err != nil {
			errOut = opm.EnhanceErr(err)
			opm.Reject()
			return
		}

		countOut = &CounterResult{}
		countOut.cas = Cas(res.Cas)
		countOut.mt = opm.EnhanceMt(res.MutationToken)
		countOut.content = res.Value

		opm.Resolve(countOut.mt)
	}))
	if err != nil {
		errOut = err
	}
	return countOut, errOut

}

func (p *kvProviderCore) StartKvOpTrace(c *Collection, operationName string, tracectx RequestSpanContext, noAttributes bool) RequestSpan {
	return c.startKvOpTrace(operationName, tracectx, p.tracer, noAttributes)
}
