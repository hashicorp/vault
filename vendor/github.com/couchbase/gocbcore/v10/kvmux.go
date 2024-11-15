package gocbcore

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/couchbase/gocbcore/v10/memd"
)

type bucketCapabilityVerifier interface {
	HasBucketCapabilityStatus(cap BucketCapability, status CapabilityStatus) bool
}

type dispatcher interface {
	DispatchDirect(req *memdQRequest) (PendingOp, error)
	RequeueDirect(req *memdQRequest, isRetry bool)
	DispatchDirectToAddress(req *memdQRequest, address string) (PendingOp, error)
	CollectionsEnabled() bool
	SupportsCollections() bool
	SetPostCompleteErrorHandler(handler postCompleteErrorHandler)
	PipelineSnapshot() (*pipelineSnapshot, error)
}

type clientProvider interface {
	GetByConnID(connID string) (*memdClient, error)
}

type kvMux struct {
	muxPtr unsafe.Pointer

	bucketName         string
	collectionsEnabled bool
	queueSize          int
	poolSize           int
	cfgMgr             *configManagementComponent
	errMapMgr          *errMapComponent

	tracer *tracerComponent
	dialer *memdClientDialerComponent

	postCompleteErrHandler postCompleteErrorHandler

	// muxStateWriteLock is necessary for functions which update the muxPtr, due to the scenario where ForceReconnect and
	// OnNewRouteConfig could race. ForceReconnect must succeed and cannot fail because OnNewRouteConfig has updated
	// the mux state whilst force is attempting to update it. We could also end up in a situation where a full reconnect
	// is occurring at the same time as a pipeline takeover and scenarios like that, including missing a config update because
	// ForceReconnect has won the race.
	// There is no need for read side locks as we are locking around an atomic and it is only the write sides that present
	// a potential issue.
	muxStateWriteLock sync.Mutex

	shutdownSig   chan struct{}
	clientCloseWg sync.WaitGroup

	noTLSSeedNode bool

	hasSeenConfigCh chan struct{}
}

type kvMuxProps struct {
	CollectionsEnabled bool
	QueueSize          int
	PoolSize           int
	NoTLSSeedNode      bool
}

func newKVMux(props kvMuxProps, cfgMgr *configManagementComponent, errMapMgr *errMapComponent, tracer *tracerComponent,
	dialer *memdClientDialerComponent, muxState *kvMuxState) *kvMux {
	mux := &kvMux{
		queueSize:          props.QueueSize,
		poolSize:           props.PoolSize,
		collectionsEnabled: props.CollectionsEnabled,
		cfgMgr:             cfgMgr,
		errMapMgr:          errMapMgr,
		tracer:             tracer,
		dialer:             dialer,
		shutdownSig:        make(chan struct{}),
		noTLSSeedNode:      props.NoTLSSeedNode,
		muxPtr:             unsafe.Pointer(muxState),
		hasSeenConfigCh:    make(chan struct{}),
		bucketName:         muxState.expectedBucketName,
	}

	cfgMgr.AddConfigWatcher(mux)

	return mux
}

func (mux *kvMux) getState() *kvMuxState {
	muxPtr := atomic.LoadPointer(&mux.muxPtr)
	if muxPtr == nil {
		return nil
	}

	return (*kvMuxState)(muxPtr)
}

func (mux *kvMux) updateState(old, new *kvMuxState) bool {
	if new == nil {
		logErrorf("Attempted to update to nil kvMuxState")
		return false
	}

	if old != nil {
		return atomic.CompareAndSwapPointer(&mux.muxPtr, unsafe.Pointer(old), unsafe.Pointer(new))
	}

	if atomic.SwapPointer(&mux.muxPtr, unsafe.Pointer(new)) != nil {
		logErrorf("Updated from nil attempted on initialized kvMuxState")
		return false
	}

	return true
}

func (mux *kvMux) clear() *kvMuxState {
	mux.muxStateWriteLock.Lock()
	val := atomic.SwapPointer(&mux.muxPtr, nil)
	mux.muxStateWriteLock.Unlock()
	return (*kvMuxState)(val)
}

func (mux *kvMux) OnNewRouteConfig(cfg *routeConfig) {
	mux.muxStateWriteLock.Lock()
	defer mux.muxStateWriteLock.Unlock()
	oldMuxState := mux.getState()
	if oldMuxState == nil {
		// We can get here if we're shutting down and a NMVB comes in from an in flight request.
		logWarnf("Received new config whilst shutting down kvmux")
		return
	}
	newMuxState := mux.newKVMuxState(cfg, oldMuxState.tlsConfig, oldMuxState.authMechanisms, oldMuxState.auth)

	// Attempt to atomically update the routing data
	if !mux.updateState(oldMuxState, newMuxState) {
		logWarnf("Someone preempted the config update, skipping update")
		return
	}

	if oldMuxState.RevID() == -1 && newMuxState.RevID() > -1 {
		if cfg.name != "" && mux.collectionsEnabled && !newMuxState.collectionsSupported {
			logDebugf("Collections disabled as unsupported")
		}

		close(mux.hasSeenConfigCh)
	}

	if !mux.collectionsEnabled {
		// If collections just aren't enabled then we never need to refresh the connections because collections
		// have come online.
		mux.pipelineTakeover(oldMuxState, newMuxState)
	} else if oldMuxState.RevID() == -1 || oldMuxState.collectionsSupported == newMuxState.collectionsSupported {
		// Get the new muxer to takeover the pipelines from the older one
		mux.pipelineTakeover(oldMuxState, newMuxState)
	} else {
		// Collections support has changed so we need to reconnect all connections in order to support the new
		// state.
		mux.reconnectPipelines(oldMuxState, newMuxState, true)
	}

	mux.requeueRequests(oldMuxState)
}

func (mux *kvMux) SetPostCompleteErrorHandler(handler postCompleteErrorHandler) {
	mux.postCompleteErrHandler = handler
}

func (mux *kvMux) ConfigRev() (int64, error) {
	clientMux := mux.getState()
	if clientMux == nil {
		return 0, errShutdown
	}
	return clientMux.RevID(), nil
}

func (mux *kvMux) ConfigUUID() string {
	clientMux := mux.getState()
	if clientMux == nil {
		return ""
	}
	return clientMux.UUID()
}

func (mux *kvMux) KeyToVbucket(key []byte) (uint16, error) {
	clientMux := mux.getState()
	if clientMux == nil || clientMux.VBMap() == nil {
		return 0, errShutdown
	}

	return clientMux.VBMap().VbucketByKey(key), nil
}

func (mux *kvMux) NumReplicas() int {
	clientMux := mux.getState()
	if clientMux == nil {
		return 0
	}

	if clientMux.VBMap() == nil {
		return 0
	}

	return clientMux.VBMap().NumReplicas()
}

func (mux *kvMux) BucketType() bucketType {
	clientMux := mux.getState()
	if clientMux == nil {
		return bktTypeInvalid
	}

	return clientMux.BucketType()
}

func (mux *kvMux) SupportsGCCCP() bool {
	clientMux := mux.getState()
	if clientMux == nil {
		return false
	}

	return clientMux.BucketType() == bktTypeNone
}

func (mux *kvMux) NumPipelines() int {
	clientMux := mux.getState()
	if clientMux == nil {
		return 0
	}

	return clientMux.NumPipelines()
}

// CollectionsEnaled returns whether or not the kv mux was created with collections enabled.
func (mux *kvMux) CollectionsEnabled() bool {
	return mux.collectionsEnabled
}

func (mux *kvMux) IsSecure() bool {
	return mux.getState().tlsConfig != nil
}

// SupportsCollections returns whether or not collections are enabled AND supported by the server.
func (mux *kvMux) SupportsCollections() bool {
	if !mux.collectionsEnabled {
		return false
	}

	clientMux := mux.getState()
	if clientMux == nil {
		return false
	}

	return clientMux.collectionsSupported
}

func (mux *kvMux) HasBucketCapabilityStatus(cap BucketCapability, status CapabilityStatus) bool {
	clientMux := mux.getState()
	if clientMux == nil {
		return status == CapabilityStatusUnknown
	}

	return clientMux.HasBucketCapabilityStatus(cap, status)
}

func (mux *kvMux) BucketCapabilityStatus(cap BucketCapability) CapabilityStatus {
	clientMux := mux.getState()
	if clientMux == nil || clientMux.RevID() == -1 {
		return CapabilityStatusUnknown
	}

	return clientMux.BucketCapabilityStatus(cap)
}

func (mux *kvMux) RouteRequest(req *memdQRequest) (*memdPipeline, error) {
	clientMux := mux.getState()
	if clientMux == nil {
		return nil, errShutdown
	}

	// We haven't seen a valid config yet so put this in the dead pipeline so
	// it'll get requeued once we do get a config.
	if clientMux.RevID() == -1 {
		return clientMux.deadPipe, nil
	}

	var srvIdx int
	repIdx := req.ReplicaIdx

	// Route to specific server
	if repIdx < 0 {
		srvIdx = -repIdx - 1
	} else {
		var err error

		bktType := clientMux.BucketType()
		if bktType == bktTypeCouchbase {
			if req.Key != nil {
				req.Vbucket = clientMux.VBMap().VbucketByKey(req.Key)
			}

			srvIdx, err = clientMux.VBMap().NodeByVbucket(req.Vbucket, uint32(repIdx))
			if err != nil {
				return nil, err
			}

		} else if bktType == bktTypeMemcached {
			if repIdx > 0 {
				// Error. Memcached buckets don't understand replicas!
				return nil, errInvalidReplica
			}

			if len(req.Key) == 0 {
				// Non-broadcast keyless Memcached bucket request
				return nil, errInvalidArgument
			}

			srvIdx, err = clientMux.KetamaMap().NodeByKey(req.Key)
			if err != nil {
				return nil, err
			}
		} else if bktType == bktTypeNone {
			// This means that we're using GCCCP and not connected to a bucket
			return nil, errGCCCPInUse
		}
	}

	pipeline := clientMux.GetPipeline(srvIdx)
	if req.ServerGroup != "" && pipeline.serverGroup != req.ServerGroup {
		return nil, ErrServerGroupMismatch
	}

	return clientMux.GetPipeline(srvIdx), nil
}

func (mux *kvMux) DispatchDirect(req *memdQRequest) (PendingOp, error) {
	mux.tracer.StartCmdTrace(req)
	req.dispatchTime = time.Now()

	for {
		pipeline, err := mux.RouteRequest(req)
		if err != nil {
			return nil, err
		}

		err = pipeline.SendRequest(req)
		if err == errPipelineClosed {
			continue
		} else if err != nil {
			if err == errPipelineFull {
				err = errOverload
			}

			shortCircuit, routeErr := mux.handleOpRoutingResp(nil, req, err)
			if shortCircuit {
				return req, nil
			}

			return nil, routeErr
		}

		break
	}

	return req, nil
}

func (mux *kvMux) RequeueDirect(req *memdQRequest, isRetry bool) {
	mux.requeueDirect(nil, req, isRetry)
}

func (mux *kvMux) requeueDirect(pipeline *memdPipeline, req *memdQRequest, isRetry bool) {
	mux.tracer.StartCmdTrace(req)

	handleError := func(err error) {
		// We only want to log an error on retries if the error isn't cancelled.
		if !isRetry || (isRetry && !errors.Is(err, ErrRequestCanceled)) {
			logErrorf("Reschedule failed, failing request, Opaque=%d, Opcode=0x%x, (%s)", req.Opaque, req.Command, err)
		}

		req.tryCallback(nil, err)
	}

	logDebugf("Request being requeued, Opaque=%d, Opcode=0x%x", req.Opaque, req.Command)

	if pipeline == nil {
		var err error
		pipeline, err = mux.RouteRequest(req)
		if err != nil {
			handleError(err)
			return
		}
	}

	for {
		err := pipeline.RequeueRequest(req)
		if err == nil {
			return
		}

		if !errors.Is(err, errPipelineClosed) {
			handleError(err)
			return
		}

		pipeline, err = mux.RouteRequest(req)
		if err != nil {
			handleError(err)
			return
		}
	}
}

func (mux *kvMux) GetByConnID(connID string) (*memdClient, error) {
	clientMux := mux.getState()
	if clientMux == nil {
		return nil, errShutdown
	}

	for _, p := range clientMux.pipelines {
		p.clientsLock.Lock()
		for _, pipeCli := range p.clients {
			pipeCli.lock.Lock()
			if pipeCli.client.connID == connID {
				pipeCli.lock.Unlock()
				p.clientsLock.Unlock()
				return pipeCli.client, nil
			}
			pipeCli.lock.Unlock()
		}
		p.clientsLock.Unlock()
	}

	return nil, errConnectionIDInvalid

}

func (mux *kvMux) DispatchDirectToAddress(req *memdQRequest, address string) (PendingOp, error) {
	mux.tracer.StartCmdTrace(req)
	req.dispatchTime = time.Now()

	// We set the ReplicaIdx to a negative number to ensure it is not redispatched
	// and we check that it was 0 to begin with to ensure it wasn't miss-used.
	if req.ReplicaIdx != 0 {
		return nil, errInvalidReplica
	}
	req.ReplicaIdx = -999999999

	for {
		clientMux := mux.getState()
		if clientMux == nil {
			return nil, errShutdown
		}

		var pipeline *memdPipeline
		for _, p := range clientMux.pipelines {
			if p.Address() == address {
				pipeline = p
				break
			}
		}

		if pipeline == nil {
			return nil, errInvalidServer
		}

		err := pipeline.SendRequest(req)
		if err == errPipelineClosed {
			continue
		} else if err != nil {
			if err == errPipelineFull {
				err = errOverload
			}

			shortCircuit, routeErr := mux.handleOpRoutingResp(nil, req, err)
			if shortCircuit {
				return req, nil
			}

			return nil, routeErr
		}

		break
	}

	return req, nil
}

func (mux *kvMux) Close() error {
	logInfof("KV Mux closing")

	mux.cfgMgr.RemoveConfigWatcher(mux)
	clientMux := mux.clear()

	if clientMux == nil {
		return errShutdown
	}

	// Trigger any memdclients that are in graceful close to forcibly close.
	close(mux.shutdownSig)

	var muxErr error
	// Shut down the client multiplexer which will close all its queues
	// effectively causing all the clients to shut down.
	for _, pipeline := range clientMux.pipelines {
		err := pipeline.Close()
		if err != nil {
			logErrorf("failed to shut down pipeline: %s", err)
			muxErr = errCliInternalError
		}
	}

	if clientMux.deadPipe != nil {
		err := clientMux.deadPipe.Close()
		if err != nil {
			logErrorf("failed to shut down deadpipe: %s", err)
			muxErr = errCliInternalError
		}
	}

	// Drain all the pipelines and error their requests, then
	//  drain the dead queue and error those requests.
	cb := func(req *memdQRequest) {
		req.tryCallback(nil, errShutdown)
	}

	mux.drainPipelines(clientMux, cb)

	mux.clientCloseWg.Wait()

	logInfof("KV Mux closed")

	return muxErr
}

func (mux *kvMux) ForceReconnect(tlsConfig *dynTLSConfig, authMechanisms []AuthMechanism, auth AuthProvider,
	reconnectLocal bool) {
	logDebugf("Forcing reconnect of all connections")
	mux.muxStateWriteLock.Lock()
	muxState := mux.getState()
	newMuxState := mux.newKVMuxState(muxState.RouteConfig(), tlsConfig, authMechanisms, auth)

	atomic.SwapPointer(&mux.muxPtr, unsafe.Pointer(newMuxState))

	mux.reconnectPipelines(muxState, newMuxState, reconnectLocal)
	mux.muxStateWriteLock.Unlock()
}

func (mux *kvMux) PipelineSnapshot() (*pipelineSnapshot, error) {
	clientMux := mux.getState()
	if clientMux == nil {
		return nil, errShutdown
	}

	return &pipelineSnapshot{
		state: clientMux,
	}, nil
}

type waitForConfigSnapshotOp struct {
	cancelCh chan struct{}
}

func (w *waitForConfigSnapshotOp) Cancel() {
	close(w.cancelCh)
}

func (mux *kvMux) WaitForConfigSnapshot(deadline time.Time, cb WaitForConfigSnapshotCallback) (PendingOp, error) {
	// No point in doing anything if we're shutdown.
	clientMux := mux.getState()
	if clientMux == nil {
		return nil, errShutdown
	}

	op := &waitForConfigSnapshotOp{
		cancelCh: make(chan struct{}),
	}

	start := time.Now()
	go func() {
		select {
		case <-mux.shutdownSig:
			cb(nil, errShutdown)
		case <-op.cancelCh:
			cb(nil, errRequestCanceled)
		case <-time.After(time.Until(deadline)):
			cb(nil, &TimeoutError{
				InnerError:   errUnambiguousTimeout,
				OperationID:  "WaitForConfigSnapshot",
				TimeObserved: time.Since(start),
			})
		case <-mux.hasSeenConfigCh:
			// Just in case.
			clientMux := mux.getState()
			if clientMux == nil {
				cb(nil, errShutdown)
				return
			}

			cb(&WaitForConfigSnapshotResult{
				Snapshot: &ConfigSnapshot{
					state: clientMux,
				},
			}, nil)
		}
	}()

	return op, nil
}

func (mux *kvMux) ConfigSnapshot() (*ConfigSnapshot, error) {
	clientMux := mux.getState()
	if clientMux == nil {
		return nil, errShutdown
	}

	return &ConfigSnapshot{
		state: clientMux,
	}, nil
}

func (mux *kvMux) handleOpRoutingResp(resp *memdQResponse, req *memdQRequest, originalErr error) (bool, error) {
	// If there is no error, we should return immediately
	if originalErr == nil {
		return false, nil
	}

	// If this operation has been cancelled, we just fail immediately.
	if errors.Is(originalErr, ErrRequestCanceled) || errors.Is(originalErr, ErrTimeout) {
		return false, originalErr
	}

	err := translateMemdError(originalErr, req)

	if err == originalErr {
		if errors.Is(err, io.EOF) && !mux.closed() {
			// The connection has gone away.
			if req.Command == memd.CmdGetClusterConfig {
				return false, err
			}

			// If the request is idempotent or not written yet then we should retry.
			if req.Idempotent() || req.ConnectionInfo().lastDispatchedTo == "" {
				if mux.waitAndRetryOperation(req, SocketNotAvailableRetryReason) {
					return true, nil
				}
			} else {
				// If the request has been dispatched then the retry reason is in flight.
				// For not causing a breaking change reasons we use socket not available for all idempotent
				// requests.
				if mux.waitAndRetryOperation(req, SocketCloseInFlightRetryReason) {
					return true, nil
				}
			}
		} else if errors.Is(err, ErrMemdClientClosed) && !mux.closed() {
			if req.Command == memd.CmdGetClusterConfig {
				return false, err
			}

			// The request can't have been dispatched yet.
			if mux.waitAndRetryOperation(req, SocketNotAvailableRetryReason) {
				return true, nil
			}
		} else if errors.Is(err, io.ErrShortWrite) {
			// This is a special case where the write has failed on the underlying connection and not all the bytes
			// were written to the network.
			if mux.waitAndRetryOperation(req, MemdWriteFailure) {
				return true, nil
			}
		} else if errors.Is(err, ErrMemdConfigOnly) {
			logWarnf("Received config-only status, will attempt to refresh config map and retry operation")
			if mux.handleConfigOnly(resp, req) {
				return true, nil
			}
		} else if resp != nil && resp.Magic == memd.CmdMagicRes {
			// We don't know anything about this error so send it to the error map
			shouldRetry := mux.errMapMgr.ShouldRetry(resp.Status)
			if shouldRetry {
				if mux.waitAndRetryOperation(req, KVErrMapRetryReason) {
					return true, nil
				}
			}
		}
	} else {
		// Handle potentially retrying the operation
		if errors.Is(err, ErrNotMyVBucket) {
			if mux.handleNotMyVbucket(resp, req) {
				return true, nil
			}
		} else if errors.Is(err, ErrDocumentLocked) {
			if mux.waitAndRetryOperation(req, KVLockedRetryReason) {
				return true, nil
			}
		} else if errors.Is(err, ErrTemporaryFailure) {
			if mux.waitAndRetryOperation(req, KVTemporaryFailureRetryReason) {
				return true, nil
			}
		} else if errors.Is(err, ErrDurableWriteInProgress) {
			if mux.waitAndRetryOperation(req, KVSyncWriteInProgressRetryReason) {
				return true, nil
			}
		} else if errors.Is(err, ErrDurableWriteReCommitInProgress) {
			if mux.waitAndRetryOperation(req, KVSyncWriteRecommitInProgressRetryReason) {
				return true, nil
			}
		}
		// If an error isn't in this list then we know what this error is but we don't support retries for it.
	}

	err = mux.errMapMgr.EnhanceKvError(err, resp, req)

	if mux.postCompleteErrHandler == nil {
		return false, err
	}

	return mux.postCompleteErrHandler(resp, req, err)
}

func (mux *kvMux) closed() bool {
	return mux.getState() == nil
}

func (mux *kvMux) waitAndRetryOperation(req *memdQRequest, reason RetryReason) bool {
	shouldRetry, retryTime := retryOrchMaybeRetry(req, reason)
	if shouldRetry {
		go func() {
			time.Sleep(time.Until(retryTime))
			mux.RequeueDirect(req, true)
		}()
		return true
	}

	return false
}

func (mux *kvMux) parseNotMyVbucketValue(value []byte, sourceAddr string) *cfgBucket {
	// Grab just the hostname from the source address
	sourceHost, err := hostFromHostPort(sourceAddr)
	if err != nil {
		logErrorf("NMV response source address was invalid, skipping config update")
		return nil
	}
	// Try to parse the value as a bucket configuration
	logDebugf("Got NMV Block: %v", string(value))
	bk, err := parseConfig(value, sourceHost)
	if err != nil {
		return nil
	}

	return bk
}

func (mux *kvMux) handleNotMyVbucket(resp *memdQResponse, req *memdQRequest) bool {
	// For range scan continue we never want to retry, the range scan is now invalid.
	isRetryableReq := req.Command != memd.CmdRangeScanContinue

	if len(resp.Value) == 0 {
		logDebugf("NMV response containing no new config")
		if !isRetryableReq {
			return false
		}
	} else {
		bk := mux.parseNotMyVbucketValue(resp.Value, resp.sourceAddr)
		if bk == nil {
			if !isRetryableReq {
				return false
			}
		} else {
			// We need to push this upstream which will then internal update the state with a new config.
			mux.cfgMgr.OnNewConfig(bk)

			if !isRetryableReq {
				return false
			}

			originalVBID := req.Vbucket
			pipeline, err := mux.RouteRequest(req)
			if err == nil {
				// If the address or vbucket has changed then just redispatch directly.
				if pipeline.Address() != resp.sourceAddr || originalVBID != req.Vbucket {
					mux.requeueDirect(pipeline, req, true)
					return true
				}
			}
		}
	}

	// Redirect it!  This may actually come back to this server, but I won't tell
	//   if you don't ;)
	return mux.waitAndRetryOperation(req, KVNotMyVBucketRetryReason)
}

func (mux *kvMux) handleConfigOnly(resp *memdQResponse, req *memdQRequest) bool {
	snapshot, err := mux.PipelineSnapshot()
	if err != nil {
		logInfof("Failed to get pipeline snapshot: %s", err)
		// Not much we can do here, attempt a retry.
		mux.RequeueDirect(req, true)
		return true
	}

	go func() {
		// Don't block the client read loop whilst we apply the config and redispatch.
		// For a start if the node this status has originated from is now not in the config then
		// calling RefreshConfig will end up blocking because we're holding the client read thread open
		// whilst also trying to shutdown the client.
		mux.cfgMgr.RefreshConfig(snapshot)
		mux.RequeueDirect(req, true)
	}()
	return true
}

func (mux *kvMux) drainPipelines(clientMux *kvMuxState, cb func(req *memdQRequest)) {
	for _, pipeline := range clientMux.pipelines {
		logDebugf("Draining queue. Address=`%s`. Num Clients=%d. Server Group=`%s`. Op Queue={%s}",
			pipeline.Address(),
			len(pipeline.Clients()),
			pipeline.ServerGroup(),
			strings.ReplaceAll(pipeline.queue.debugString(), "\n", ", "))
		pipeline.Drain(cb)
	}
	if clientMux.deadPipe != nil {
		clientMux.deadPipe.Drain(cb)
	}
}

func (mux *kvMux) newKVMuxState(cfg *routeConfig, tlsConfig *dynTLSConfig, authMechanisms []AuthMechanism,
	auth AuthProvider) *kvMuxState {
	poolSize := 1
	if !cfg.IsGCCCPConfig() {
		poolSize = mux.poolSize
	}

	useTls := tlsConfig != nil

	var kvServerList []routeEndpoint
	if mux.noTLSSeedNode {
		// The order of the kv server list matters, so we need to maintain the same order and just replace the seed
		// node.
		if useTls {
			kvServerList = make([]routeEndpoint, len(cfg.kvServerList.SSLEndpoints))
			copy(kvServerList, cfg.kvServerList.SSLEndpoints)

			for i, ep := range cfg.kvServerList.NonSSLEndpoints {
				if ep.IsSeedNode {
					kvServerList[i] = ep
				}
			}
		} else {
			kvServerList = cfg.kvServerList.NonSSLEndpoints
		}
	} else {
		if useTls {
			kvServerList = cfg.kvServerList.SSLEndpoints
		} else {
			kvServerList = cfg.kvServerList.NonSSLEndpoints
		}
	}

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintln("KV muxer applying endpoints:"))
	buffer.WriteString(fmt.Sprintf("Bucket: %s\n", cfg.name))
	for _, ep := range kvServerList {
		buffer.WriteString(fmt.Sprintf("  - %s\n", ep.Address))
	}

	logDebugf(buffer.String())

	pipelines := make([]*memdPipeline, len(kvServerList))
	for i, hostPort := range kvServerList {
		trimmedHostPort := routeEndpoint{
			Address:     trimSchemePrefix(hostPort.Address),
			IsSeedNode:  hostPort.IsSeedNode,
			ServerGroup: hostPort.ServerGroup,
		}

		getCurClientFn := func(cancelSig <-chan struct{}) (*memdClient, error) {
			return mux.dialer.SlowDialMemdClient(cancelSig, trimmedHostPort, tlsConfig, auth, authMechanisms,
				mux.handleOpRoutingResp, mux.handleServerRequest)
		}
		pipeline := newPipeline(trimmedHostPort, poolSize, mux.queueSize, getCurClientFn)

		pipelines[i] = pipeline
	}

	return newKVMuxState(cfg, kvServerList, tlsConfig, authMechanisms, auth, mux.bucketName, pipelines,
		newDeadPipeline(mux.queueSize))
}

func (mux *kvMux) reconnectPipelines(oldMuxState *kvMuxState, newMuxState *kvMuxState, reconnectSeed bool) {
	oldPipelines := list.New()

	for _, pipeline := range oldMuxState.pipelines {
		oldPipelines.PushBack(pipeline)
	}

	for _, pipeline := range newMuxState.pipelines {
		// If we aren't reconnecting the seed node then we need to take its clients and make sure we don't
		// end up closing it down.
		if pipeline.isSeedNode && !reconnectSeed {
			oldPipeline := mux.stealPipeline(pipeline.Address(), oldPipelines)

			if oldPipeline != nil {
				pipeline.Takeover(oldPipeline)
			}
		}

		pipeline.StartClients()
	}

	for e := oldPipelines.Front(); e != nil; e = e.Next() {
		pipeline, ok := e.Value.(*memdPipeline)
		if !ok {
			logErrorf("Failed to cast old pipeline")
			continue
		}

		clients := pipeline.GracefulClose()

		for _, client := range clients {
			mux.closeMemdClient(client, errForcedReconnect)
		}
	}
}

func (mux *kvMux) requeueRequests(oldMuxState *kvMuxState) {
	// Gather all the requests from all the old pipelines and then
	//  sort and redispatch them (which will use the new pipelines)
	var requestList []*memdQRequest
	mux.drainPipelines(oldMuxState, func(req *memdQRequest) {
		requestList = append(requestList, req)
	})

	sort.Sort(memdQRequestSorter(requestList))

	for _, req := range requestList {
		req.processingLock.Lock()
		stopCmdTraceLocked(req)
		req.processingLock.Unlock()

		// If the command is a get cluster config then we cancel it rather than requeuing.
		// Get cluster config is explicitly sent a specific pipeline so we do not want to requeue.
		// This may seem like it'll cause the poller to take longer to fetch a config but that's
		// OK because we can only have here by something fetching a new config anyway.
		if req.Command == memd.CmdGetClusterConfig {
			req.tryCallback(nil, ErrRequestCanceled)
			continue
		}
		mux.RequeueDirect(req, false)
	}
}

// closeMemdClient will gracefully close the memdclient, spinning up a goroutine to watch for when the client
// shuts down. The error provided is the error sent to any callback handlers for persistent operations which are
// currently live in the client.
func (mux *kvMux) closeMemdClient(client *memdClient, err error) {
	mux.clientCloseWg.Add(1)
	client.GracefulClose(err)
	go func(client *memdClient) {
		select {
		case <-client.CloseNotify():
			logDebugf("Memdclient %s/%p completed graceful shutdown", client.Address(), client)
		case <-mux.shutdownSig:
			logDebugf("Memdclient %s/%p being forcibly shutdown", client.Address(), client)
			// Force the client to close even if there are requests in flight.
			err := client.Close()
			if err != nil {
				logErrorf("failed to shutdown memdclient: %s", err)
			}
			<-client.CloseNotify()
			logDebugf("Memdclient %s/%p completed shutdown", client.Address(), client)
		}
		mux.clientCloseWg.Done()
	}(client)
}

func (mux *kvMux) stealPipeline(address string, oldPipelines *list.List) *memdPipeline {
	for e := oldPipelines.Front(); e != nil; e = e.Next() {
		pipeline, ok := e.Value.(*memdPipeline)
		if !ok {
			logErrorf("Failed to cast old pipeline")
			continue
		}

		if pipeline.Address() == address {
			oldPipelines.Remove(e)
			return pipeline
		}
	}

	return nil
}

func (mux *kvMux) pipelineTakeover(oldMux, newMux *kvMuxState) {
	oldPipelines := list.New()

	// Gather all our old pipelines up for takeover and what not
	if oldMux != nil {
		for _, pipeline := range oldMux.pipelines {
			oldPipelines.PushBack(pipeline)
		}
	}

	// Initialize new pipelines (possibly with a takeover)
	for _, pipeline := range newMux.pipelines {
		oldPipeline := mux.stealPipeline(pipeline.Address(), oldPipelines)
		if oldPipeline != nil {
			pipeline.Takeover(oldPipeline)
		}

		pipeline.StartClients()
	}

	// Shut down any pipelines that were not taken over
	for e := oldPipelines.Front(); e != nil; e = e.Next() {
		pipeline, ok := e.Value.(*memdPipeline)
		if !ok {
			logErrorf("Failed to cast old pipeline")
			continue
		}

		clients := pipeline.GracefulClose()
		for _, client := range clients {
			mux.closeMemdClient(client, nil)
		}
	}

	if oldMux != nil && oldMux.deadPipe != nil {
		err := oldMux.deadPipe.Close()
		if err != nil {
			logErrorf("Failed to properly close abandoned dead pipe (%s)", err)
		}
	}
}

func (mux *kvMux) handleServerRequest(pak *memd.Packet) {
	if pak.Command == memd.CmdSet {
		// We copy out the extras before handling the packet in its own goroutine.
		// If we don't do this then the memdclient is going to free the packet and by the
		// time that we access extras they'll be nil.
		extras := make([]byte, len(pak.Extras))
		copy(extras, pak.Extras)
		go func() {
			snapshot, err := mux.PipelineSnapshot()
			if err != nil {
				logInfof("Failed to get pipeline snapshot: %s", err)
				return
			}
			mux.cfgMgr.OnNewConfigChangeNotifBrief(snapshot, extras)
		}()
		return
	}

	logWarnf("Received an unknown command type for a server request: OP=0x%x", pak.Command)
}
