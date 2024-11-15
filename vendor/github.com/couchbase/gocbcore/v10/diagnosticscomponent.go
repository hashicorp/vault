package gocbcore

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"

	"github.com/couchbase/gocbcore/v10/memd"
)

type diagnosticsComponent struct {
	kvMux               *kvMux
	httpMux             *httpMux
	httpComponent       *httpComponent
	bucket              string
	defaultRetry        RetryStrategy
	pollerErrorProvider pollerErrorProvider

	// preConfigBootstrapError must only be used for checking for bootstrap errors when a config has not yet been seen.
	preConfigBootstrapError     error
	preConfigBootstrapErrorLock sync.Mutex
}

func newDiagnosticsComponent(kvMux *kvMux, httpMux *httpMux, httpComponent *httpComponent, bucket string,
	defaultRetry RetryStrategy, pollerErrorProvider pollerErrorProvider) *diagnosticsComponent {
	return &diagnosticsComponent{
		kvMux:               kvMux,
		httpMux:             httpMux,
		bucket:              bucket,
		httpComponent:       httpComponent,
		defaultRetry:        defaultRetry,
		pollerErrorProvider: pollerErrorProvider,
	}
}

func (dc *diagnosticsComponent) onBootstrapFail(err error) {
	// It doesn't really matter if we overwrite this error.
	dc.preConfigBootstrapErrorLock.Lock()
	dc.preConfigBootstrapError = err
	dc.preConfigBootstrapErrorLock.Unlock()
}

func (dc *diagnosticsComponent) pingKV(ctx context.Context, interval time.Duration, deadline time.Time,
	retryStrat RetryStrategy, user string, op *pingOp) {

	var userFrame *memd.UserImpersonationFrame
	if len(user) > 0 {
		userFrame = &memd.UserImpersonationFrame{
			User: []byte(user),
		}
	}

	if !deadline.IsZero() {
		// We have to setup a new child context with its own deadline because services have their own timeout values.
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, deadline)
		defer cancel()
	}

	for {
		iter, err := dc.kvMux.PipelineSnapshot()
		if err != nil {
			if errors.Is(err, ErrShutdown) {
				op.lock.Lock()
				op.results[MemdService] = append(op.results[MemdService], EndpointPingResult{
					Error: errShutdown,
					Scope: op.bucketName,
					ID:    uuid.New().String(),
					State: PingStateError,
				})
				op.handledOneLocked(0)
				op.lock.Unlock()
				return
			}

			logErrorf("failed to get pipeline snapshot")

			select {
			case <-ctx.Done():
				ctxErr := ctx.Err()
				var cancelReason error
				if errors.Is(ctxErr, context.Canceled) {
					cancelReason = ctxErr
				} else {
					cancelReason = errUnambiguousTimeout
				}

				op.lock.Lock()
				op.results[MemdService] = append(op.results[MemdService], EndpointPingResult{
					Error: cancelReason,
					Scope: op.bucketName,
					ID:    uuid.New().String(),
					State: PingStateTimeout,
				})
				op.handledOneLocked(0)
				op.lock.Unlock()
				return
			case <-time.After(interval):
				continue
			}
		}

		if iter.RevID() > -1 {
			var wg sync.WaitGroup
			iter.Iterate(0, func(p *memdPipeline) bool {
				wg.Add(1)
				go func(pipeline *memdPipeline) {
					serverAddress := pipeline.Address()

					startTime := time.Now()
					handler := func(resp *memdQResponse, req *memdQRequest, err error) {
						pingLatency := time.Since(startTime)

						state := PingStateOK
						if err != nil {
							if errors.Is(err, ErrTimeout) {
								state = PingStateTimeout
							} else {
								state = PingStateError
							}
						}

						op.lock.Lock()
						op.results[MemdService] = append(op.results[MemdService], EndpointPingResult{
							Endpoint: serverAddress,
							Error:    err,
							Latency:  pingLatency,
							Scope:    op.bucketName,
							ID:       fmt.Sprintf("%p", pipeline),
							State:    state,
						})
						op.lock.Unlock()
						wg.Done()
					}

					req := &memdQRequest{
						Packet: memd.Packet{
							Magic:                  memd.CmdMagicReq,
							Command:                memd.CmdNoop,
							Datatype:               0,
							Cas:                    0,
							Key:                    nil,
							Value:                  nil,
							UserImpersonationFrame: userFrame,
						},
						Callback:      handler,
						RetryStrategy: retryStrat,
					}

					curOp, err := dc.kvMux.DispatchDirectToAddress(req, pipeline.Address())
					if err != nil {
						op.lock.Lock()
						op.results[MemdService] = append(op.results[MemdService], EndpointPingResult{
							Endpoint: redactSystemData(serverAddress),
							Error:    err,
							Latency:  0,
							Scope:    op.bucketName,
						})
						op.lock.Unlock()
						wg.Done()
						return
					}

					if !deadline.IsZero() {
						start := time.Now()
						req.SetTimer(time.AfterFunc(deadline.Sub(start), func() {
							connInfo := req.ConnectionInfo()
							count, reasons := req.Retries()
							req.cancelWithCallback(&TimeoutError{
								InnerError:         errUnambiguousTimeout,
								OperationID:        "PingKV",
								Opaque:             req.Identifier(),
								TimeObserved:       time.Since(start),
								RetryReasons:       reasons,
								RetryAttempts:      count,
								LastDispatchedTo:   connInfo.lastDispatchedTo,
								LastDispatchedFrom: connInfo.lastDispatchedFrom,
								LastConnectionID:   connInfo.lastConnectionID,
							})
						}))
					}

					op.lock.Lock()
					op.subops = append(op.subops, pingSubOp{
						endpoint: serverAddress,
						op:       curOp,
					})
					op.lock.Unlock()
				}(p)

				// We iterate through all pipelines
				return false
			})

			wg.Wait()
			op.lock.Lock()
			op.handledOneLocked(iter.RevID())
			op.lock.Unlock()
			return
		}

		select {
		case <-ctx.Done():
			ctxErr := ctx.Err()
			var cancelReason error
			if errors.Is(ctxErr, context.Canceled) {
				cancelReason = ctxErr
			} else {
				cancelReason = errUnambiguousTimeout
			}

			op.lock.Lock()
			op.results[MemdService] = append(op.results[MemdService], EndpointPingResult{
				Error: cancelReason,
				Scope: op.bucketName,
				ID:    uuid.New().String(),
				State: PingStateTimeout,
			})
			op.handledOneLocked(iter.RevID())
			op.lock.Unlock()
			return
		case <-time.After(interval):
		}
	}
}

func (dc *diagnosticsComponent) pingHTTP(ctx context.Context, service ServiceType,
	interval time.Duration, deadline time.Time, retryStrat RetryStrategy, op *pingOp, ignoreMissingServices bool) {

	if !deadline.IsZero() {
		// We have to setup a new child context with its own deadline because services have their own timeout values.
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, deadline)
		defer cancel()
	}

	muxer := dc.httpMux

	var path string
	switch service {
	case N1qlService:
		path = "/admin/ping"
	case CbasService:
		path = "/admin/ping"
	case FtsService:
		path = "/api/ping"
	case CapiService:
		path = "/"
	}

	for {
		clientMux := muxer.Get()
		if clientMux == nil {
			op.lock.Lock()
			op.results[service] = append(op.results[service], EndpointPingResult{
				Error: errShutdown,
				Scope: op.bucketName,
				ID:    uuid.New().String(),
				State: PingStateError,
			})
			op.handledOneLocked(0)
			op.lock.Unlock()
			return
		}

		if clientMux.revID > -1 {
			var epList []routeEndpoint
			switch service {
			case N1qlService:
				epList = clientMux.n1qlEpList
			case CbasService:
				epList = clientMux.cbasEpList
			case FtsService:
				epList = clientMux.ftsEpList
			case MgmtService:
				epList = clientMux.mgmtEpList
			case CapiService:
				epList = clientMux.capiEpList
			}

			if len(epList) == 0 {
				op.lock.Lock()
				if !ignoreMissingServices {
					op.results[service] = append(op.results[service], EndpointPingResult{
						Error: errServiceNotAvailable,
						Scope: op.bucketName,
						ID:    uuid.New().String(),
					})
				}
				op.handledOneLocked(clientMux.revID)
				op.lock.Unlock()
				return
			}

			var wg sync.WaitGroup
			for _, ep := range epList {
				wg.Add(1)
				go func(ep string) {
					defer wg.Done()
					req := &httpRequest{
						Service:       service,
						Method:        "GET",
						Path:          path,
						Endpoint:      ep,
						IsIdempotent:  true,
						RetryStrategy: retryStrat,
						Context:       ctx,
						UniqueID:      uuid.New().String(),
					}
					start := time.Now()
					resp, err := dc.httpComponent.DoInternalHTTPRequest(req, false)
					pingLatency := time.Since(start)
					state := PingStateOK
					if err != nil {
						if errors.Is(err, ErrTimeout) {
							state = PingStateTimeout
						} else {
							state = PingStateError
						}
					} else {
						defer resp.Body.Close()
						if resp.StatusCode > 200 {
							state = PingStateError
							b, pErr := ioutil.ReadAll(resp.Body)
							if pErr != nil {
								logDebugf("Failed to read response body for ping: %v", pErr)
							}

							err = errors.New(string(b))
						}
					}
					op.lock.Lock()
					op.results[service] = append(op.results[service], EndpointPingResult{
						Endpoint: ep,
						Error:    err,
						Latency:  pingLatency,
						Scope:    op.bucketName,
						ID:       uuid.New().String(),
						State:    state,
					})
					op.lock.Unlock()
				}(ep.Address)
			}

			wg.Wait()
			op.lock.Lock()
			op.handledOneLocked(clientMux.revID)
			op.lock.Unlock()
			return
		}

		select {
		case <-ctx.Done():
			ctxErr := ctx.Err()
			var cancelReason error
			if errors.Is(ctxErr, context.Canceled) {
				cancelReason = ctxErr
			} else {
				cancelReason = errUnambiguousTimeout
			}

			op.lock.Lock()
			op.results[service] = append(op.results[service], EndpointPingResult{
				Error: cancelReason,
				Scope: op.bucketName,
				ID:    uuid.New().String(),
				State: PingStateTimeout,
			})
			op.handledOneLocked(clientMux.revID)
			op.lock.Unlock()
			return
		case <-time.After(interval):
		}
	}
}

func (dc *diagnosticsComponent) Ping(opts PingOptions, cb PingCallback) (PendingOp, error) {
	bucketName := ""
	if dc.bucket != "" {
		bucketName = redactMetaData(dc.bucket)
	}

	ignoreMissingServices := false
	serviceTypes := opts.ServiceTypes
	if len(serviceTypes) == 0 {
		// We're defaulting to pinging what we can so don't ping anything that isn't in the cluster config
		ignoreMissingServices = true
		serviceTypes = []ServiceType{MemdService, CapiService, N1qlService, FtsService, CbasService, MgmtService}
	}

	ignoreMissingServices = ignoreMissingServices || opts.ignoreMissingServices

	ctx, cancelFunc := context.WithCancel(context.Background())

	op := &pingOp{
		callback:   cb,
		remaining:  int32(len(serviceTypes)),
		results:    make(map[ServiceType][]EndpointPingResult),
		bucketName: bucketName,
		httpCancel: cancelFunc,
	}

	retryStrat := newFailFastRetryStrategy()

	// interval is how long to wait between checking if we've seen a cluster config
	interval := 10 * time.Millisecond

	for _, serviceType := range serviceTypes {
		switch serviceType {
		case MemdService:
			go dc.pingKV(ctx, interval, opts.KVDeadline, retryStrat, opts.User, op)
		case CapiService:
			go dc.pingHTTP(ctx, CapiService, interval, opts.CapiDeadline, retryStrat, op, ignoreMissingServices)
		case N1qlService:
			go dc.pingHTTP(ctx, N1qlService, interval, opts.N1QLDeadline, retryStrat, op, ignoreMissingServices)
		case FtsService:
			go dc.pingHTTP(ctx, FtsService, interval, opts.FtsDeadline, retryStrat, op, ignoreMissingServices)
		case CbasService:
			go dc.pingHTTP(ctx, CbasService, interval, opts.CbasDeadline, retryStrat, op, ignoreMissingServices)
		case MgmtService:
			go dc.pingHTTP(ctx, MgmtService, interval, opts.MgmtDeadline, retryStrat, op, ignoreMissingServices)
		}
	}

	return op, nil
}

// Diagnostics returns diagnostics information about the client.
// Mainly containing a list of open connections and their current
// states.
func (dc *diagnosticsComponent) Diagnostics(opts DiagnosticsOptions) (*DiagnosticInfo, error) {
	for {
		iter, err := dc.kvMux.PipelineSnapshot()
		if err != nil {
			return nil, err
		}

		var conns []MemdConnInfo

		iter.Iterate(0, func(pipeline *memdPipeline) bool {
			pipeline.clientsLock.Lock()
			for _, pipecli := range pipeline.clients {
				localAddr := ""
				remoteAddr := ""
				var lastActivity time.Time

				pipecli.lock.Lock()
				if pipecli.client != nil {
					localAddr = pipecli.client.LocalAddress()
					remoteAddr = pipecli.client.Address()
					lastActivityUs := atomic.LoadInt64(&pipecli.client.lastActivity)
					if lastActivityUs != 0 {
						lastActivity = time.Unix(0, lastActivityUs)
					}
				}
				pipecli.lock.Unlock()

				conn := MemdConnInfo{
					LocalAddr:    localAddr,
					RemoteAddr:   remoteAddr,
					LastActivity: lastActivity,
					ID:           fmt.Sprintf("%p", pipecli),
					State:        pipecli.State(),
				}
				if dc.bucket != "" {
					conn.Scope = redactMetaData(dc.bucket)
				}
				conns = append(conns, conn)
			}
			pipeline.clientsLock.Unlock()
			return false
		})

		expected := len(conns)
		connected := 0
		for _, conn := range conns {
			if conn.State == EndpointStateConnected {
				connected++
			}
		}

		state := ClusterStateOffline
		if connected == expected {
			state = ClusterStateOnline
		} else if connected > 1 {
			state = ClusterStateDegraded
		}

		endIter, err := dc.kvMux.PipelineSnapshot()
		if err != nil {
			return nil, err
		}
		if iter.RevID() == endIter.RevID() {
			return &DiagnosticInfo{
				ConfigRev: iter.RevID(),
				MemdConns: conns,
				State:     state,
			}, nil
		}
	}
}

func (dc *diagnosticsComponent) checkKVReady(desiredState ClusterState, op *waitUntilOp) {
	for {
		iter, err := dc.kvMux.PipelineSnapshot()
		if err != nil {
			if errors.Is(err, ErrShutdown) {
				op.cancel(err)
				return
			}

			logErrorf("failed to get pipeline snapshot: %v", err)

			shouldRetry, until := retryOrchMaybeRetry(op, NoPipelineSnapshotRetryReason)
			if !shouldRetry {
				op.cancel(err)
				return
			}

			select {
			case <-op.stopCh:
				return
			case <-time.After(time.Until(until)):
				continue
			}
		}

		var connectErr error
		revID := iter.RevID()
		if revID == -1 {
			// We've not seen a config so let's see if we've been informed about any errors.
			dc.preConfigBootstrapErrorLock.Lock()
			connectErr = dc.preConfigBootstrapError
			logDebugf("Bootstrap error found before config seen: %v", connectErr)
			dc.preConfigBootstrapErrorLock.Unlock()

			// If there's no error appearing from the pipeline client then let's check the poller
			if connectErr == nil && dc.pollerErrorProvider != nil {
				pollerErr := dc.pollerErrorProvider.PollerError()

				// We don't care about timeouts, they don't tell us anything we want to know.
				if pollerErr != nil && !errors.Is(pollerErr, ErrTimeout) {
					logDebugf("Error found in poller before config seen: %v", pollerErr)
					connectErr = pollerErr
				}
			}

			if connectErr == nil {
				logDebugf("No config seen yet in kv muxer but no errors found.")
			}
		} else if revID > -1 {
			expected := iter.NumPipelines()
			connected := 0
			iter.Iterate(0, func(pipeline *memdPipeline) bool {
				pipeline.clientsLock.Lock()
				defer pipeline.clientsLock.Unlock()
				for _, cli := range pipeline.clients {
					state := cli.State()
					if state == EndpointStateConnected {
						connected++
						if desiredState == ClusterStateDegraded {
							// If we're after degraded state then we can just bail early as we've already fulfilled that.
							return true
						}

						// We only need one of the pipeline clients to be connected for this pipeline to be considered
						// online.
						break
					}

					err := cli.Error()
					if err != nil {
						logDebugf("Error found in client after config seen: %v", err)
						connectErr = err

						// If the desired state is degraded then we need to keep trying as a different client or pipeline
						// might be connected. If it's online then we can bail now as we'll never achieve that.
						if desiredState == ClusterStateOnline {
							return true
						}
					}
				}

				return false
			})

			// If there's no error appearing from the pipeline client then let's check the poller
			if connectErr == nil && dc.pollerErrorProvider != nil {
				pollerErr := dc.pollerErrorProvider.PollerError()

				// We don't care about timeouts, they don't tell us anything we want to know.
				if pollerErr != nil && !errors.Is(pollerErr, ErrTimeout) {
					logDebugf("Error found in poller after config seen: %v", pollerErr)
					connectErr = pollerErr
				}
			}

			switch desiredState {
			case ClusterStateDegraded:
				if connected > 0 {
					op.lock.Lock()
					op.handledOneLocked()
					op.lock.Unlock()

					return
				}
			case ClusterStateOnline:
				if connected == expected {
					op.lock.Lock()
					op.handledOneLocked()
					op.lock.Unlock()

					return
				}
			default:
				// How we got here no-one does know
				// But round and round we must go
			}
		}

		var until time.Time
		if connectErr == nil {
			var shouldRetry bool
			shouldRetry, until = retryOrchMaybeRetry(op, NotReadyRetryReason)
			if !shouldRetry {
				op.cancel(errCliInternalError)
				return
			}
		} else {
			var shouldRetry bool
			if errors.Is(connectErr, ErrBucketNotFound) {
				shouldRetry, until = retryOrchMaybeRetry(op, BucketNotReadyReason)
			} else {
				shouldRetry, until = retryOrchMaybeRetry(op, ConnectionErrorRetryReason)
			}
			if !shouldRetry {
				op.cancel(connectErr)
				return
			}
		}

		select {
		case <-op.stopCh:
			return
		case <-time.After(time.Until(until)):
		}
	}
}

func (dc *diagnosticsComponent) checkHTTPReady(ctx context.Context, service ServiceType,
	desiredState ClusterState, forceWait bool, op *waitUntilOp) {
	retryStrat := &failFastRetryStrategy{}
	muxer := dc.httpMux

	var path string
	switch service {
	case N1qlService:
		path = "/admin/ping"
	case CbasService:
		path = "/admin/ping"
	case FtsService:
		path = "/api/ping"
	case CapiService:
		path = "/"
	case MgmtService:
		path = ""
	}

	for {
		clientMux := muxer.Get()
		if clientMux == nil {
			op.cancel(errShutdown)
			return
		}
		var connectErr error
		if clientMux.revID == -1 {
			// We've not seen a config so let's see if we've been informed about any errors.
			dc.preConfigBootstrapErrorLock.Lock()
			connectErr = dc.preConfigBootstrapError
			logDebugf("Bootstrap error found before config seen: %v", connectErr)
			dc.preConfigBootstrapErrorLock.Unlock()

			// If there's no error appearing from the pipeline client then let's check the poller
			if connectErr == nil && dc.pollerErrorProvider != nil {
				pollerErr := dc.pollerErrorProvider.PollerError()

				// We don't care about timeouts, they don't tell us anything we want to know.
				if pollerErr != nil && !errors.Is(pollerErr, ErrTimeout) {
					logDebugf("Error found in poller before config seen: %v", pollerErr)
					connectErr = pollerErr
				}
			}

			if connectErr == nil {
				logDebugf("No config seen yet in http muxer but no errors found.")
			}
		} else {
			var epList []routeEndpoint
			switch service {
			case N1qlService:
				epList = clientMux.n1qlEpList
			case CbasService:
				epList = clientMux.cbasEpList
			case FtsService:
				epList = clientMux.ftsEpList
			case CapiService:
				epList = clientMux.capiEpList
			case MgmtService:
				epList = clientMux.mgmtEpList
			}

			connected := uint32(0)
			func() {
				ctx, cancel := context.WithCancel(ctx)
				defer cancel()

				var wg sync.WaitGroup
				for _, ep := range epList {
					wg.Add(1)
					go func(ep string) {
						defer wg.Done()
						req := &httpRequest{
							Service:       service,
							Method:        "GET",
							Path:          path,
							RetryStrategy: retryStrat,
							Endpoint:      ep,
							IsIdempotent:  true,
							Context:       ctx,
							UniqueID:      uuid.New().String(),
						}
						resp, err := dc.httpComponent.DoInternalHTTPRequest(req, false)
						if err != nil {
							if errors.Is(err, context.Canceled) {
								return
							}

							logDebugf("Error returned for HTTP request for service %d: %v", service, err)

							if desiredState == ClusterStateOnline {
								// Cancel this run entirely, we can't satisfy the requirements
								cancel()
							}
							return
						}
						err = resp.Body.Close()
						if err != nil {
							logDebugf("Failed to close response body: %s", err)
						}
						if resp.StatusCode != 200 {
							logDebugf("Non-200 status code returned for HTTP request for service %d: %d", service, resp.StatusCode)
							if desiredState == ClusterStateOnline {
								// Cancel this run entirely, we can't satisfy the requirements
								cancel()
							}
							return
						}
						atomic.AddUint32(&connected, 1)
						if desiredState == ClusterStateDegraded {
							// Cancel this run entirely, we've successfully satisfied the requirements
							cancel()
						}
					}(ep.Address)
				}

				wg.Wait()
			}()

			switch desiredState {
			case ClusterStateDegraded:
				if !forceWait && len(epList) == 0 {
					op.lock.Lock()
					op.handledOneLocked()
					op.lock.Unlock()

					return
				}
				// If there are no entries in the epList then the service is not online and so cannot be ready.
				if len(epList) > 0 && atomic.LoadUint32(&connected) > 0 {
					op.lock.Lock()
					op.handledOneLocked()
					op.lock.Unlock()

					return
				}
			case ClusterStateOnline:
				if !forceWait && len(epList) == 0 {
					op.lock.Lock()
					op.handledOneLocked()
					op.lock.Unlock()

					return
				}
				// If there are no entries in the epList then the service is not online and so cannot be ready.
				if len(epList) > 0 && atomic.LoadUint32(&connected) == uint32(len(epList)) {
					op.lock.Lock()
					op.handledOneLocked()
					op.lock.Unlock()

					return
				}
			default:
				// How we got here no-one does know
				// But round and round we must go
			}
		}

		var until time.Time
		if connectErr == nil {
			var shouldRetry bool
			shouldRetry, until = retryOrchMaybeRetry(op, NotReadyRetryReason)
			if !shouldRetry {
				op.cancel(errCliInternalError)
				return
			}
		} else {
			var shouldRetry bool
			if errors.Is(connectErr, ErrBucketNotFound) {
				shouldRetry, until = retryOrchMaybeRetry(op, BucketNotReadyReason)
			} else {
				shouldRetry, until = retryOrchMaybeRetry(op, ConnectionErrorRetryReason)
			}
			if !shouldRetry {
				op.cancel(connectErr)
				return
			}
		}

		select {
		case <-op.stopCh:
			return
		case <-time.After(time.Until(until)):
		}
	}
}

func (dc *diagnosticsComponent) WaitUntilReady(deadline time.Time, forceWait bool, opts WaitUntilReadyOptions,
	cb WaitUntilReadyCallback) (PendingOp, error) {
	desiredState := opts.DesiredState
	if desiredState == ClusterStateOffline {
		return nil, wrapError(errInvalidArgument, "cannot use offline as a desired state")
	}

	if desiredState == 0 {
		desiredState = ClusterStateOnline
	}

	retry := opts.RetryStrategy
	if retry == nil {
		retry = dc.defaultRetry
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	op := &waitUntilOp{
		remaining:  int32(len(opts.ServiceTypes)),
		stopCh:     make(chan struct{}),
		callback:   cb,
		httpCancel: cancelFunc,
		retryStrat: retry,
	}

	op.lock.Lock()
	start := time.Now()
	op.timer = time.AfterFunc(deadline.Sub(start), func() {
		op.cancel(&TimeoutError{
			InnerError:    errUnambiguousTimeout,
			OperationID:   "WaitUntilReady",
			TimeObserved:  time.Since(start),
			RetryReasons:  op.RetryReasons(),
			RetryAttempts: op.RetryAttempts(),
		})
	})
	op.lock.Unlock()

	for _, serviceType := range opts.ServiceTypes {
		switch serviceType {
		case MemdService:
			go dc.checkKVReady(desiredState, op)
		case CapiService:
			go dc.checkHTTPReady(ctx, CapiService, desiredState, forceWait, op)
		case N1qlService:
			go dc.checkHTTPReady(ctx, N1qlService, desiredState, forceWait, op)
		case FtsService:
			go dc.checkHTTPReady(ctx, FtsService, desiredState, forceWait, op)
		case CbasService:
			go dc.checkHTTPReady(ctx, CbasService, desiredState, forceWait, op)
		case MgmtService:
			go dc.checkHTTPReady(ctx, MgmtService, desiredState, forceWait, op)
		}
	}

	return op, nil
}
