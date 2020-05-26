package gocbcore

import (
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/couchbase/gocbcore/v9/memd"
)

func (cidMgr *collectionsComponent) createKey(scopeName, collectionName string) string {
	return fmt.Sprintf("%s.%s", scopeName, collectionName)
}

type collectionsComponent struct {
	idMap                map[string]*collectionIDCache
	mapLock              sync.Mutex
	mux                  *kvMux
	maxQueueSize         int
	tracer               *tracerComponent
	defaultRetryStrategy RetryStrategy
}

type collectionIDProps struct {
	MaxQueueSize         int
	DefaultRetryStrategy RetryStrategy
}

func newCollectionIDManager(props collectionIDProps, mux *kvMux, tracerCmpt *tracerComponent) *collectionsComponent {
	cidMgr := &collectionsComponent{
		mux:                  mux,
		idMap:                make(map[string]*collectionIDCache),
		maxQueueSize:         props.MaxQueueSize,
		tracer:               tracerCmpt,
		defaultRetryStrategy: props.DefaultRetryStrategy,
	}

	mux.SetPostCompleteErrorHandler(cidMgr.handleOpRoutingResp)

	return cidMgr
}

func (cidMgr *collectionsComponent) handleCollectionUnknown(req *memdQRequest) bool {
	// We cannot retry requests with no collection information
	if req.CollectionName == "" && req.ScopeName == "" {
		return false
	}

	shouldRetry, retryTime := retryOrchMaybeRetry(req, KVCollectionOutdatedRetryReason)
	if shouldRetry {
		go func() {
			time.Sleep(retryTime.Sub(time.Now()))
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

	return cidMgr.mux.DispatchDirect(req)
}

// GetCollectionID does not trigger retries on unknown collection. This is because the request sets the scope and collection
// name in the key rather than in the corresponding fields.
func (cidMgr *collectionsComponent) GetCollectionID(scopeName string, collectionName string, opts GetCollectionIDOptions, cb GetCollectionIDCallback) (PendingOp, error) {
	tracer := cidMgr.tracer.CreateOpTrace("GetCollectionID", opts.TraceContext)

	handler := func(resp *memdQResponse, req *memdQRequest, err error) {
		cidCache, ok := cidMgr.get(scopeName, collectionName)
		if !ok {
			cidCache = cidMgr.newCollectionIDCache()
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
		cidCache.id = collectionID
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

	return cidMgr.mux.DispatchDirect(req)
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
	cidMgr.mapLock.Lock()
	delete(cidMgr.idMap, cidMgr.createKey(scopeName, collectionName))
	cidMgr.mapLock.Unlock()
}

func (cidMgr *collectionsComponent) newCollectionIDCache() *collectionIDCache {
	return &collectionIDCache{
		mux:          cidMgr.mux,
		maxQueueSize: cidMgr.maxQueueSize,
		parent:       cidMgr,
	}
}

type collectionIDCache struct {
	opQueue        *memdOpQueue
	id             uint32
	collectionName string
	scopeName      string
	parent         *collectionsComponent
	mux            *kvMux
	lock           sync.Mutex
	err            error
	maxQueueSize   int
}

func (cid *collectionIDCache) sendWithCid(req *memdQRequest) error {
	cid.lock.Lock()
	req.CollectionID = cid.id
	cid.lock.Unlock()
	_, err := cid.mux.DispatchDirect(req)
	if err != nil {
		return err
	}

	return nil
}

func (cid *collectionIDCache) queueRequest(req *memdQRequest) error {
	return cid.opQueue.Push(req, cid.maxQueueSize)
}

func (cid *collectionIDCache) refreshCid(req *memdQRequest) error {
	err := cid.opQueue.Push(req, cid.maxQueueSize)
	if err != nil {
		return err
	}

	_, err = cid.parent.GetCollectionID(req.ScopeName, req.CollectionName, GetCollectionIDOptions{TraceContext: req.RootTraceContext},
		func(result *GetCollectionIDResult, err error) {
			if err != nil {
				if errors.Is(err, ErrCollectionNotFound) {
					// If this is a collection unknown error then it'll get retried so we just leave the
					// cid in pending state. Either the collection will eventually come online or this request will
					// timeout.
					if cid.opQueue.Remove(req) {
						cid.lock.Lock()
						cid.id = unknownCid
						cid.lock.Unlock()
						if cid.parent.handleCollectionUnknown(req) {
							return
						}
					}
				}
				cid.opQueue.Close()
				cid.opQueue.Drain(func(request *memdQRequest) {
					request.tryCallback(nil, err)
				})
				cid.opQueue = nil
				return
			}

			cid.opQueue.Close()
			cid.opQueue.Drain(func(request *memdQRequest) {
				request.CollectionID = result.CollectionID
				cid.mux.RequeueDirect(request, false)
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
		cid.id = pendingCid
		cid.opQueue = newMemdOpQueue()
		cid.lock.Unlock()
		return cid.refreshCid(req)
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

	if !cidMgr.mux.SupportsCollections() {
		if !(noCollection || defaultCollection) || collectionIDPresent {
			return nil, errCollectionsUnsupported
		}
		_, err := cidMgr.mux.DispatchDirect(req)
		if err != nil {
			return nil, err
		}

		return req, nil
	}

	if noCollection || defaultCollection || collectionIDPresent {
		return cidMgr.mux.DispatchDirect(req)
	}

	cidCache, ok := cidMgr.get(req.ScopeName, req.CollectionName)
	if !ok {
		cidCache = cidMgr.newCollectionIDCache()
		cidCache.id = unknownCid
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
		cidCache = cidMgr.newCollectionIDCache()
		cidCache.id = unknownCid
		cidMgr.add(cidCache, req.ScopeName, req.CollectionName)
	}
	cidCache.lock.Lock()
	if cidCache.id != unknownCid && cidCache.id != pendingCid {
		cidCache.id = unknownCid
	}
	cidCache.lock.Unlock()

	err := cidCache.dispatch(req)
	if err != nil {
		req.tryCallback(nil, err)
	}
}

func (cidMgr *collectionsComponent) HasDurabilityLevelStatus(status durabilityLevelStatus) bool {
	return cidMgr.mux.HasDurabilityLevelStatus(status)
}

func (cidMgr *collectionsComponent) BucketType() bucketType {
	return cidMgr.mux.BucketType()
}

func (cidMgr *collectionsComponent) KeyToVbucket(key []byte) (uint16, error) {
	return cidMgr.mux.KeyToVbucket(key)
}
