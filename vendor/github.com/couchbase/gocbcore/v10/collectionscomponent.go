package gocbcore

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/couchbase/gocbcore/v10/memd"
)

func (cidMgr *collectionsComponent) createKey(scopeName, collectionName string) string {
	return fmt.Sprintf("%s.%s", scopeName, collectionName)
}

type collectionsComponent struct {
	idMap                map[string]*collectionIDCache
	mapLock              sync.Mutex
	dispatcher           dispatcher
	maxQueueSize         int
	tracer               *tracerComponent
	defaultRetryStrategy RetryStrategy
	cfgMgr               configManager

	// pendingOpQueue is used when collections are enabled but we've not yet seen a cluster config to confirm
	// whether or not collections are supported.
	pendingOpQueue *memdOpQueue
	configSeen     uint32
}

type collectionIDProps struct {
	MaxQueueSize         int
	DefaultRetryStrategy RetryStrategy
}

func newCollectionIDManager(props collectionIDProps, dispatcher dispatcher, tracer *tracerComponent,
	cfgMgr configManager) *collectionsComponent {
	cidMgr := &collectionsComponent{
		dispatcher:           dispatcher,
		idMap:                make(map[string]*collectionIDCache),
		maxQueueSize:         props.MaxQueueSize,
		tracer:               tracer,
		defaultRetryStrategy: props.DefaultRetryStrategy,
		cfgMgr:               cfgMgr,
		pendingOpQueue:       newMemdOpQueue(),
	}

	cfgMgr.AddConfigWatcher(cidMgr)
	dispatcher.SetPostCompleteErrorHandler(cidMgr.handleOpRoutingResp)

	return cidMgr
}

func (cidMgr *collectionsComponent) OnNewRouteConfig(cfg *routeConfig) {
	if !atomic.CompareAndSwapUint32(&cidMgr.configSeen, 0, 1) {
		return
	}

	colsSupported := cfg.ContainsBucketCapability("collections")
	cidMgr.cfgMgr.RemoveConfigWatcher(cidMgr)
	cidMgr.pendingOpQueue.Close()
	cidMgr.pendingOpQueue.Drain(func(request *memdQRequest) {
		// Anything in this queue is here because collections were present so if we definitely don't support collections
		// then fail them.
		if !colsSupported {
			request.tryCallback(nil, errCollectionsUnsupported)
			return
		}
		cidMgr.requeue(request)
	})
}

func (cidMgr *collectionsComponent) handleCollectionUnknown(req *memdQRequest) bool {
	if !canRetryOnCollectionUnknown(req) {
		return false
	}

	shouldRetry, retryTime := retryOrchMaybeRetry(req, KVCollectionOutdatedRetryReason)
	if shouldRetry {
		go func() {
			time.Sleep(time.Until(retryTime))
			if isDefaultCollection(req.ScopeName, req.CollectionName) {
				// If the request is against the default collection then there's no point trying to
				// refresh the cid, instead we just retry the operation.
				cidMgr.dispatcher.RequeueDirect(req, true)
				return
			}

			cidMgr.requeue(req)
		}()
	}

	return shouldRetry
}

func (cidMgr *collectionsComponent) handleOpRoutingResp(resp *memdQResponse, req *memdQRequest, err error) (bool, error) {
	if errors.Is(err, ErrCollectionNotFound) || errors.Is(err, ErrScopeNotFound) {
		if cidMgr.handleCollectionUnknown(req) {
			return true, nil
		}
	}

	return false, err
}

func (cidMgr *collectionsComponent) GetCollectionManifest(opts GetCollectionManifestOptions, cb GetCollectionManifestCallback) (PendingOp, error) {
	tracer := cidMgr.tracer.StartTelemeteryHandler(metricValueServiceAnalyticsValue, "GetCollectionManifest", opts.TraceContext)

	handler := func(resp *memdQResponse, req *memdQRequest, err error) {
		if err != nil {
			cb(nil, err)
			tracer.Finish()
			return
		}

		res := GetCollectionManifestResult{
			Manifest: resp.Value,
		}
		res.Internal.ResourceUnits = req.ResourceUnits()

		tracer.Finish()
		cb(&res, nil)
	}

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = cidMgr.defaultRetryStrategy
	}

	var userFrame *memd.UserImpersonationFrame
	if len(opts.User) > 0 {
		userFrame = &memd.UserImpersonationFrame{
			User: []byte(opts.User),
		}
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:                  memd.CmdMagicReq,
			Command:                memd.CmdCollectionsGetManifest,
			UserImpersonationFrame: userFrame,
		},
		Callback:         handler,
		RetryStrategy:    opts.RetryStrategy,
		RootTraceContext: opts.TraceContext,
	}

	op, err := cidMgr.dispatcher.DispatchDirect(req)
	if err != nil {
		tracer.Finish()
		return nil, err
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.SetTimer(time.AfterFunc(opts.Deadline.Sub(start), func() {
			connInfo := req.ConnectionInfo()
			count, reasons := req.Retries()
			req.cancelWithCallbackAndFinishTracer(&TimeoutError{
				InnerError:         errUnambiguousTimeout,
				OperationID:        "GetCollectionManifest",
				Opaque:             req.Identifier(),
				TimeObserved:       time.Since(start),
				RetryReasons:       reasons,
				RetryAttempts:      count,
				LastDispatchedTo:   connInfo.lastDispatchedTo,
				LastDispatchedFrom: connInfo.lastDispatchedFrom,
				LastConnectionID:   connInfo.lastConnectionID,
			}, tracer)
		}))
	}

	return op, nil
}

func (cidMgr *collectionsComponent) GetAllCollectionManifests(opts GetAllCollectionManifestsOptions, cb GetAllCollectionManifestsCallback) (PendingOp, error) {
	tracer := cidMgr.tracer.StartTelemeteryHandler(metricValueServiceAnalyticsValue, "GetAllCollectionManifests", opts.TraceContext)

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = cidMgr.defaultRetryStrategy
	}

	iter, err := cidMgr.dispatcher.PipelineSnapshot()
	if err != nil {
		tracer.Finish()
		return nil, err
	}

	manifests := make(map[string]SingleServerManifestResult)
	manifestsLock := sync.Mutex{}

	op := &multiPendingOp{
		isIdempotent: true,
	}

	opCompleteLocked := func() {
		completed := op.IncrementCompletedOps()
		if iter.NumPipelines()-int(completed) == 0 {
			tracer.Finish()
			cb(&GetAllCollectionManifestsResult{Manifests: manifests}, nil)
		}
	}

	var setTimer func(request *memdQRequest)
	if opts.Deadline.IsZero() {
		setTimer = func(_ *memdQRequest) {}
	} else {
		start := time.Now()
		timeout := opts.Deadline.Sub(start)

		setTimer = func(req *memdQRequest) {
			req.SetTimer(time.AfterFunc(timeout, func() {
				connInfo := req.ConnectionInfo()
				count, reasons := req.Retries()
				req.cancelWithCallbackAndFinishTracer(&TimeoutError{
					InnerError:         errUnambiguousTimeout,
					OperationID:        "GetAllCollectionManifests",
					Opaque:             req.Identifier(),
					TimeObserved:       time.Since(start),
					RetryReasons:       reasons,
					RetryAttempts:      count,
					LastDispatchedTo:   connInfo.lastDispatchedTo,
					LastDispatchedFrom: connInfo.lastDispatchedFrom,
					LastConnectionID:   connInfo.lastConnectionID,
				}, tracer)
			}))
		}
	}
	var userFrame *memd.UserImpersonationFrame
	if len(opts.User) > 0 {
		userFrame = &memd.UserImpersonationFrame{
			User: []byte(opts.User),
		}
	}

	iter.Iterate(0, func(pipeline *memdPipeline) bool {
		handler := func(resp *memdQResponse, req *memdQRequest, err error) {
			manifestsLock.Lock()

			res := SingleServerManifestResult{
				Error: err,
			}

			if resp != nil {
				res.Manifest = resp.Value
			}

			manifests[pipeline.address] = res
			opCompleteLocked()

			manifestsLock.Unlock()
		}

		req := &memdQRequest{
			Packet: memd.Packet{
				Magic:                  memd.CmdMagicReq,
				Command:                memd.CmdCollectionsGetManifest,
				UserImpersonationFrame: userFrame,
			},
			Callback:         handler,
			RetryStrategy:    opts.RetryStrategy,
			RootTraceContext: opts.TraceContext,
		}

		curOp, err := cidMgr.dispatcher.DispatchDirectToAddress(req, pipeline.Address())
		if err == nil {
			setTimer(req)

			op.ops = append(op.ops, curOp)
			return false
		}

		manifestsLock.Lock()
		manifests[pipeline.address] = SingleServerManifestResult{Error: err}
		opCompleteLocked()
		manifestsLock.Unlock()

		return false
	})

	return op, nil
}

// GetCollectionID does not trigger retries on unknown collection. This is because the request sets the scope and collection
// name in the key rather than in the corresponding fields.
func (cidMgr *collectionsComponent) GetCollectionID(scopeName string, collectionName string, opts GetCollectionIDOptions,
	cb GetCollectionIDCallback) (PendingOp, error) {
	tracer := cidMgr.tracer.StartTelemeteryHandler(metricValueServiceAnalyticsValue, "GetCollectionID", opts.TraceContext)

	handler := func(resp *memdQResponse, req *memdQRequest, err error) {
		if err != nil {
			tracer.Finish()
			cb(nil, err)
			return
		}

		manifestID := binary.BigEndian.Uint64(resp.Extras[0:])
		collectionID := binary.BigEndian.Uint32(resp.Extras[8:])

		cidMgr.upsert(scopeName, collectionName, collectionID)

		res := GetCollectionIDResult{
			ManifestID:   manifestID,
			CollectionID: collectionID,
		}
		res.Internal.ResourceUnits = req.ResourceUnits()

		tracer.Finish()
		cb(&res, nil)
	}

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = cidMgr.defaultRetryStrategy
	}

	keyScopeName := scopeName
	if keyScopeName == "" {
		keyScopeName = "_default"
	}
	keyCollectionName := collectionName
	if keyCollectionName == "" {
		keyCollectionName = "_default"
	}

	var userFrame *memd.UserImpersonationFrame
	if len(opts.User) > 0 {
		userFrame = &memd.UserImpersonationFrame{
			User: []byte(opts.User),
		}
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:                  memd.CmdMagicReq,
			Command:                memd.CmdCollectionsGetID,
			Datatype:               0,
			Cas:                    0,
			Extras:                 nil,
			Key:                    nil,
			Value:                  []byte(fmt.Sprintf("%s.%s", keyScopeName, keyCollectionName)),
			Vbucket:                0,
			UserImpersonationFrame: userFrame,
		},
		ReplicaIdx:       -1,
		RetryStrategy:    opts.RetryStrategy,
		RootTraceContext: opts.TraceContext,
	}

	req.Callback = handler

	op, err := cidMgr.dispatcher.DispatchDirect(req)
	if err != nil {
		tracer.Finish()
		return nil, err
	}

	if !opts.Deadline.IsZero() {
		start := time.Now()
		req.SetTimer(time.AfterFunc(opts.Deadline.Sub(start), func() {
			connInfo := req.ConnectionInfo()
			count, reasons := req.Retries()
			req.cancelWithCallbackAndFinishTracer(&TimeoutError{
				InnerError:         errUnambiguousTimeout,
				OperationID:        "GetCollectionID",
				Opaque:             req.Identifier(),
				TimeObserved:       time.Since(start),
				RetryReasons:       reasons,
				RetryAttempts:      count,
				LastDispatchedTo:   connInfo.lastDispatchedTo,
				LastDispatchedFrom: connInfo.lastDispatchedFrom,
				LastConnectionID:   connInfo.lastConnectionID,
			}, tracer)
		}))
	}

	return op, nil
}

func (cidMgr *collectionsComponent) upsert(scopeName, collectionName string, value uint32) *collectionIDCache {
	cidMgr.mapLock.Lock()
	id, ok := cidMgr.idMap[cidMgr.createKey(scopeName, collectionName)]
	if !ok {
		id = cidMgr.newCollectionIDCache(scopeName, collectionName)
		key := cidMgr.createKey(scopeName, collectionName)
		cidMgr.idMap[key] = id
	}
	id.lock.Lock()
	id.setID(value)
	id.lock.Unlock()
	cidMgr.mapLock.Unlock()

	return id
}

func (cidMgr *collectionsComponent) getAndMaybeInsert(scopeName, collectionName string, value uint32) *collectionIDCache {
	cidMgr.mapLock.Lock()
	id, ok := cidMgr.idMap[cidMgr.createKey(scopeName, collectionName)]
	if !ok {
		id = cidMgr.newCollectionIDCache(scopeName, collectionName)
		id.lock.Lock()
		id.setID(value)
		id.lock.Unlock()

		key := cidMgr.createKey(scopeName, collectionName)
		cidMgr.idMap[key] = id
	}
	cidMgr.mapLock.Unlock()

	return id
}

func (cidMgr *collectionsComponent) remove(scopeName, collectionName string) {
	logDebugf("Removing cache entry for %s.%s", scopeName, collectionName)
	cidMgr.mapLock.Lock()
	delete(cidMgr.idMap, cidMgr.createKey(scopeName, collectionName))
	cidMgr.mapLock.Unlock()
}

func (cidMgr *collectionsComponent) newCollectionIDCache(scope, collection string) *collectionIDCache {
	return &collectionIDCache{
		dispatcher:     cidMgr.dispatcher,
		maxQueueSize:   cidMgr.maxQueueSize,
		parent:         cidMgr,
		scopeName:      scope,
		collectionName: collection,
	}
}

type collectionIDCache struct {
	opQueue        *memdOpQueue
	id             uint32
	collectionName string
	scopeName      string
	parent         *collectionsComponent
	dispatcher     dispatcher
	lock           sync.Mutex
	maxQueueSize   int
}

func (cid *collectionIDCache) sendWithCid(req *memdQRequest) error {
	cid.lock.Lock()
	id := cid.id
	cid.lock.Unlock()
	if err := setRequestCid(req, id); err != nil {
		logDebugf("Failed to set collection ID on request: %v", err)
		return err
	}

	_, err := cid.dispatcher.DispatchDirect(req)
	if err != nil {
		return err
	}

	return nil
}

func (cid *collectionIDCache) queueRequest(req *memdQRequest) error {
	cid.lock.Lock()
	defer cid.lock.Unlock()
	return cid.opQueue.Push(req, cid.maxQueueSize)
}

func (cid *collectionIDCache) setID(id uint32) {
	logDebugf("Setting cache ID to %d for %s.%s", id, cid.scopeName, cid.collectionName)
	cid.id = id
}

func (cid *collectionIDCache) refreshCid(req *memdQRequest) error {
	err := cid.opQueue.Push(req, cid.maxQueueSize)
	if err != nil {
		return err
	}

	logDebugf("Refreshing collection ID for %s.%s", req.ScopeName, req.CollectionName)
	_, err = cid.parent.GetCollectionID(req.ScopeName, req.CollectionName, GetCollectionIDOptions{TraceContext: req.RootTraceContext},
		func(result *GetCollectionIDResult, err error) {
			if err != nil {
				if errors.Is(err, ErrCollectionNotFound) || errors.Is(err, ErrScopeNotFound) {
					// The collection is unknown so we need to mark the cid unknown and attempt to retry the request.
					// Retrying the request will requeue it in the cid manager so either it will pick up the unknown cid
					// and cause a refresh or another request will and this one will get queued within the cache.
					// Either the collection will eventually come online or this request will timeout.
					logDebugf("Collection %s.%s not found, attempting retry", req.ScopeName, req.CollectionName)
					cid.lock.Lock()
					cid.setID(unknownCid)
					cid.lock.Unlock()
					if cid.opQueue.Remove(req) {
						if cid.parent.handleCollectionUnknown(req) {
							return
						}
					} else {
						logDebugf("Request no longer existed in op queue, possibly cancelled?",
							req.Opaque, req.CollectionName)
					}
				} else {
					logDebugf("Collection ID refresh failed: %v", err)
				}

				// There was an error getting this collection ID so lets remove the cache from the manager and try to
				// callback on all of the queued requests.
				cid.parent.remove(req.ScopeName, req.CollectionName)
				cid.opQueue.Close()
				cid.opQueue.Drain(func(request *memdQRequest) {
					request.tryCallback(nil, err)
				})
				return
			}

			// We successfully got the cid, the GetCollectionID itself will have handled setting the ID on this cache,
			// so lets reset the op queue and requeue all of our requests.
			cid.lock.Lock()
			opQueue := cid.opQueue
			cid.opQueue = newMemdOpQueue()
			cid.lock.Unlock()

			logDebugf("Collection %s.%s refresh succeeded, requeuing %d requests", req.ScopeName, req.CollectionName, opQueue.items.Len())
			opQueue.Close()
			opQueue.Drain(func(request *memdQRequest) {
				request.AddResourceUnitsFromUnitResult(result.Internal.ResourceUnits)

				if err := setRequestCid(request, result.CollectionID); err != nil {
					logDebugf("Failed to set collection ID on request: %v", err)
					request.cancelWithCallback(err)
					return
				}
				cid.dispatcher.RequeueDirect(request, false)
			})
		},
	)

	return err
}

func (cid *collectionIDCache) dispatch(req *memdQRequest) error {
	cid.lock.Lock()
	// if the cid is unknown then mark the request pending and refresh cid first
	// if it's pending then queue the request
	// otherwise send the request
	switch cid.id {
	case unknownCid:
		logDebugf("Collection %s.%s unknown, refreshing id", req.ScopeName, req.CollectionName)
		cid.setID(pendingCid)
		newOpQueue := newMemdOpQueue()
		if cid.opQueue != nil {
			// Drain the old queue into the new one so that we move over outstanding requests.
			cid.opQueue.Close()
			cid.opQueue.Drain(func(request *memdQRequest) {
				err := newOpQueue.Push(request, 0)
				if err != nil {
					request.tryCallback(nil, err)
				}
			})
		}
		cid.opQueue = newOpQueue

		// We attempt to send the refresh inside of the lock, that way we haven't released the lock and allowed an op
		// to get queued if we need to move the status back to unknown. Without doing this it's possible for one or
		// more op(s) to sneak into the queue and then no more requests come in and those sit in the queue until they
		// timeout because nothing is triggering the cid refresh.
		err := cid.refreshCid(req)
		if err != nil {
			// We've failed to send the cid refresh so we need to set it back to unknown otherwise it'll never
			// get updated.
			cid.setID(unknownCid)
			cid.lock.Unlock()
			return err
		}
		cid.lock.Unlock()
		return nil
	case pendingCid:
		logDebugf("Collection %s.%s pending, queueing request OP=0x%x", req.ScopeName, req.CollectionName, req.Command)
		cid.lock.Unlock()
		return cid.queueRequest(req)
	default:
		cid.lock.Unlock()
		return cid.sendWithCid(req)
	}
}

func (cidMgr *collectionsComponent) Dispatch(req *memdQRequest) (PendingOp, error) {
	isDefaultCollectionName := isDefaultCollection(req.ScopeName, req.CollectionName)
	collectionIDPresent := req.CollectionID > 0

	// If the user didn't enable collections then we can just not bother with any collections logic.
	if !cidMgr.dispatcher.CollectionsEnabled() {
		if !isDefaultCollectionName || collectionIDPresent {
			return nil, errCollectionsUnsupported
		}
		_, err := cidMgr.dispatcher.DispatchDirect(req)
		if err != nil {
			return nil, err
		}

		return req, nil
	}

	if isDefaultCollectionName || collectionIDPresent {
		return cidMgr.dispatcher.DispatchDirect(req)
	}

	if atomic.LoadUint32(&cidMgr.configSeen) == 0 {
		logDebugf("Collections are enabled but we've not yet seen a config so queueing request")
		err := cidMgr.pendingOpQueue.Push(req, cidMgr.maxQueueSize)
		if err != nil {
			return nil, err
		}

		return req, nil
	}

	if !cidMgr.dispatcher.SupportsCollections() {
		return nil, errCollectionsUnsupported
	}

	cidCache := cidMgr.getAndMaybeInsert(req.ScopeName, req.CollectionName, unknownCid)
	err := cidCache.dispatch(req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (cidMgr *collectionsComponent) requeue(req *memdQRequest) {
	cidCache := cidMgr.getAndMaybeInsert(req.ScopeName, req.CollectionName, unknownCid)
	cidCache.lock.Lock()
	if cidCache.id != unknownCid && cidCache.id != pendingCid {
		cidCache.setID(unknownCid)
	}
	cidCache.lock.Unlock()

	err := cidCache.dispatch(req)
	if err != nil {
		req.tryCallback(nil, err)
	}
}

func setRequestCid(req *memdQRequest, cid uint32) error {
	if req.Command == memd.CmdRangeScanCreate {
		var createReq *rangeScanCreateRequest
		if err := json.Unmarshal(req.Value, &createReq); err != nil {
			return err
		}
		createReq.Collection = strconv.FormatUint(uint64(cid), 16)
		value, err := json.Marshal(createReq)
		if err != nil {
			return err
		}

		req.Value = value
		return nil
	}
	req.CollectionID = cid
	return nil
}

func canRetryOnCollectionUnknown(req *memdQRequest) bool {
	switch req.Command {
	case memd.CmdCollectionsGetID:
		return false
	case memd.CmdRangeScanContinue:
		return false
	case memd.CmdRangeScanCancel:
		return false
	case memd.CmdDcpStreamReq:
		return false
	}

	return true
}

func isDefaultCollection(scopeName, collectionName string) bool {
	return (collectionName == "" || collectionName == "_default") && (scopeName == "" || scopeName == "_default")
}
