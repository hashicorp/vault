package gocbcore

import (
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/couchbase/gocbcore/v9/memd"
)

func (cidMgr *collectionsComponent) createKey(scopeName, collectionName string) string {
	return fmt.Sprintf("%s.%s", scopeName, collectionName)
}

type collectionsComponent struct {
	idMap                map[string]*collectionIDCache
	mapLock              sync.Mutex
	dispatcher           dispatcher
	maxQueueSize         int
	tracer               tracerManager
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

func newCollectionIDManager(props collectionIDProps, dispatcher dispatcher, tracer tracerManager,
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
	// We cannot retry requests with no collection information.
	// This also prevents the GetCollectionID requests from being automatically retried.
	if req.CollectionName == "" && req.ScopeName == "" {
		return false
	}

	shouldRetry, retryTime := retryOrchMaybeRetry(req, KVCollectionOutdatedRetryReason)
	if shouldRetry {
		go func() {
			time.Sleep(time.Until(retryTime))
			cidMgr.requeue(req)
		}()
	}

	return shouldRetry
}

func (cidMgr *collectionsComponent) handleOpRoutingResp(resp *memdQResponse, req *memdQRequest, err error) (bool, error) {
	if errors.Is(err, ErrCollectionNotFound) {
		if cidMgr.handleCollectionUnknown(req) {
			return true, nil
		}
	}

	return false, err
}

func (cidMgr *collectionsComponent) GetCollectionManifest(opts GetCollectionManifestOptions, cb GetCollectionManifestCallback) (PendingOp, error) {
	tracer := cidMgr.tracer.CreateOpTrace("GetCollectionManifest", opts.TraceContext)

	handler := func(resp *memdQResponse, req *memdQRequest, err error) {
		if err != nil {
			cb(nil, err)
			tracer.Finish()
			return
		}

		res := GetCollectionManifestResult{
			Manifest: resp.Value,
		}

		tracer.Finish()
		cb(&res, nil)
	}

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = cidMgr.defaultRetryStrategy
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:    memd.CmdMagicReq,
			Command:  memd.CmdCollectionsGetManifest,
			Datatype: 0,
			Cas:      0,
			Extras:   nil,
			Key:      nil,
			Value:    nil,
		},
		Callback:         handler,
		RetryStrategy:    opts.RetryStrategy,
		RootTraceContext: opts.TraceContext,
	}

	return cidMgr.dispatcher.DispatchDirect(req)
}

// GetCollectionID does not trigger retries on unknown collection. This is because the request sets the scope and collection
// name in the key rather than in the corresponding fields.
func (cidMgr *collectionsComponent) GetCollectionID(scopeName string, collectionName string, opts GetCollectionIDOptions,
	cb GetCollectionIDCallback) (PendingOp, error) {
	tracer := cidMgr.tracer.CreateOpTrace("GetCollectionID", opts.TraceContext)

	handler := func(resp *memdQResponse, req *memdQRequest, err error) {
		cidCache, ok := cidMgr.get(scopeName, collectionName)
		if !ok {
			cidCache = cidMgr.newCollectionIDCache(scopeName, collectionName)
			cidMgr.add(cidCache, scopeName, collectionName)
		}

		if err != nil {
			tracer.Finish()
			cb(nil, err)
			return
		}

		manifestID := binary.BigEndian.Uint64(resp.Extras[0:])
		collectionID := binary.BigEndian.Uint32(resp.Extras[8:])

		cidCache.lock.Lock()
		cidCache.setID(collectionID)
		cidCache.lock.Unlock()

		res := GetCollectionIDResult{
			ManifestID:   manifestID,
			CollectionID: collectionID,
		}

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

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:    memd.CmdMagicReq,
			Command:  memd.CmdCollectionsGetID,
			Datatype: 0,
			Cas:      0,
			Extras:   nil,
			Key:      []byte(fmt.Sprintf("%s.%s", keyScopeName, keyCollectionName)),
			Value:    nil,
			Vbucket:  0,
		},
		ReplicaIdx:       -1,
		RetryStrategy:    opts.RetryStrategy,
		RootTraceContext: opts.TraceContext,
	}

	req.Callback = handler

	return cidMgr.dispatcher.DispatchDirect(req)
}

func (cidMgr *collectionsComponent) add(id *collectionIDCache, scopeName, collectionName string) {
	key := cidMgr.createKey(scopeName, collectionName)
	cidMgr.mapLock.Lock()
	cidMgr.idMap[key] = id
	cidMgr.mapLock.Unlock()
}

func (cidMgr *collectionsComponent) get(scopeName, collectionName string) (*collectionIDCache, bool) {
	cidMgr.mapLock.Lock()
	id, ok := cidMgr.idMap[cidMgr.createKey(scopeName, collectionName)]
	cidMgr.mapLock.Unlock()
	if !ok {
		return nil, false
	}

	return id, true
}

func (cidMgr *collectionsComponent) remove(scopeName, collectionName string) {
	logDebugf("Removing cache entry for", scopeName, collectionName)
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
	req.CollectionID = cid.id
	cid.lock.Unlock()
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
				if errors.Is(err, ErrCollectionNotFound) {
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
			logDebugf("Collection %s.%s refresh succeeded, requeuing requests", req.ScopeName, req.CollectionName)
			cid.lock.Lock()
			opQueue := cid.opQueue
			cid.opQueue = newMemdOpQueue()
			cid.lock.Unlock()

			opQueue.Close()
			opQueue.Drain(func(request *memdQRequest) {
				request.CollectionID = result.CollectionID
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
		cid.opQueue = newMemdOpQueue()

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
		logDebugf("Collection %s.%s pending, queueing request", req.ScopeName, req.CollectionName)
		cid.lock.Unlock()
		return cid.queueRequest(req)
	default:
		cid.lock.Unlock()
		return cid.sendWithCid(req)
	}
}

func (cidMgr *collectionsComponent) Dispatch(req *memdQRequest) (PendingOp, error) {
	noCollection := req.CollectionName == "" && req.ScopeName == ""
	defaultCollection := req.CollectionName == "_default" && req.ScopeName == "_default"
	collectionIDPresent := req.CollectionID > 0

	// If the user didn't enable collections then we can just not bother with any collections logic.
	if !cidMgr.dispatcher.CollectionsEnabled() {
		if !(noCollection || defaultCollection) || collectionIDPresent {
			return nil, errCollectionsUnsupported
		}
		_, err := cidMgr.dispatcher.DispatchDirect(req)
		if err != nil {
			return nil, err
		}

		return req, nil
	}

	if noCollection || defaultCollection || collectionIDPresent {
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

	cidCache, ok := cidMgr.get(req.ScopeName, req.CollectionName)
	if !ok {
		cidCache = cidMgr.newCollectionIDCache(req.ScopeName, req.CollectionName)
		cidCache.setID(unknownCid)
		cidMgr.add(cidCache, req.ScopeName, req.CollectionName)
	}
	err := cidCache.dispatch(req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (cidMgr *collectionsComponent) requeue(req *memdQRequest) {
	cidCache, ok := cidMgr.get(req.ScopeName, req.CollectionName)
	if !ok {
		cidCache = cidMgr.newCollectionIDCache(req.ScopeName, req.CollectionName)
		cidCache.setID(unknownCid)
		cidMgr.add(cidCache, req.ScopeName, req.CollectionName)
	}
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
