package gocbcore

import (
	"encoding/binary"
	"errors"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/couchbase/gocbcore/v10/memd"

	"github.com/golang/snappy"
)

func isCompressibleOp(command memd.CmdCode) bool {
	switch command {
	case memd.CmdSet:
		fallthrough
	case memd.CmdAdd:
		fallthrough
	case memd.CmdReplace:
		fallthrough
	case memd.CmdAppend:
		fallthrough
	case memd.CmdPrepend:
		return true
	}
	return false
}

type postCompleteErrorHandler func(resp *memdQResponse, req *memdQRequest, err error) (bool, error)

type memdClient struct {
	lastActivity          int64
	dcpAckSize            int
	dcpFlowRecv           int
	closeNotify           chan bool
	connID                string
	closed                bool
	closedError           error
	conn                  memdConn
	opList                *memdOpMap
	features              []memd.HelloFeature
	lock                  sync.Mutex
	streamEndNotSupported bool
	breaker               circuitBreaker
	postErrHandler        postCompleteErrorHandler
	tracer                *tracerComponent
	zombieLogger          *zombieLoggerComponent

	dcpQueueSize         int
	compressionMinSize   int
	compressionMinRatio  float64
	disableDecompression bool

	cancelBootstrapSig <-chan struct{}
}

type dcpBuffer struct {
	resp       *memdQResponse
	packetLen  int
	isInternal bool
}

type memdClientProps struct {
	ClientID string

	DCPQueueSize         int
	CompressionMinSize   int
	CompressionMinRatio  float64
	DisableDecompression bool
}

func newMemdClient(props memdClientProps, conn memdConn, breakerCfg CircuitBreakerConfig, postErrHandler postCompleteErrorHandler,
	tracer *tracerComponent, zombieLogger *zombieLoggerComponent) *memdClient {
	client := memdClient{
		closeNotify:    make(chan bool),
		connID:         props.ClientID + "/" + formatCbUID(randomCbUID()),
		postErrHandler: postErrHandler,
		tracer:         tracer,
		zombieLogger:   zombieLogger,
		conn:           conn,
		opList:         newMemdOpMap(),

		dcpQueueSize:         props.DCPQueueSize,
		compressionMinRatio:  props.CompressionMinRatio,
		compressionMinSize:   props.CompressionMinSize,
		disableDecompression: props.DisableDecompression,
	}

	if breakerCfg.Enabled {
		client.breaker = newLazyCircuitBreaker(breakerCfg, client.sendCanary)
	} else {
		client.breaker = newNoopCircuitBreaker()
	}

	client.run()
	return &client
}

func (client *memdClient) SupportsFeature(feature memd.HelloFeature) bool {
	return checkSupportsFeature(client.features, feature)
}

func (client *memdClient) EnableDcpBufferAck(bufferAckSize int) {
	client.dcpAckSize = bufferAckSize
}

func (client *memdClient) maybeSendDcpBufferAck(packetLen int) {
	client.dcpFlowRecv += packetLen
	if client.dcpFlowRecv < client.dcpAckSize {
		return
	}

	ackAmt := client.dcpFlowRecv

	extrasBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(extrasBuf, uint32(ackAmt))

	err := client.conn.WritePacket(&memd.Packet{
		Magic:   memd.CmdMagicReq,
		Command: memd.CmdDcpBufferAck,
		Extras:  extrasBuf,
	})
	if err != nil {
		logWarnf("Failed to dispatch DCP buffer ack: %s", err)
	}

	client.dcpFlowRecv -= ackAmt
}

func (client *memdClient) Address() string {
	return client.conn.RemoteAddr()
}

func (client *memdClient) CloseNotify() chan bool {
	return client.closeNotify
}

func (client *memdClient) takeRequestOwnership(req *memdQRequest) error {
	client.lock.Lock()
	defer client.lock.Unlock()

	if client.closed {
		logDebugf("Attempted to put dispatched op in drained opmap")
		if client.closedError != nil {
			return client.closedError
		}
		return errMemdClientClosed
	}

	if !atomic.CompareAndSwapPointer(&req.waitingIn, nil, unsafe.Pointer(client)) {
		logDebugf("Attempted to put dispatched op in new opmap")
		return errRequestAlreadyDispatched
	}

	if req.isCancelled() {
		atomic.CompareAndSwapPointer(&req.waitingIn, unsafe.Pointer(client), nil)
		return errRequestCanceled
	}

	connInfo := memdQRequestConnInfo{
		lastDispatchedTo:   client.Address(),
		lastDispatchedFrom: client.conn.LocalAddr(),
		lastConnectionID:   client.connID,
	}
	req.SetConnectionInfo(connInfo)

	client.opList.Add(req)
	return nil
}

func (client *memdClient) CancelRequest(req *memdQRequest, err error) bool {
	client.lock.Lock()
	defer client.lock.Unlock()

	if client.closed {
		logDebugf("Attempted to remove op from drained opmap")
		return false
	}

	removed := client.opList.Remove(req)
	if removed {
		atomic.CompareAndSwapPointer(&req.waitingIn, unsafe.Pointer(client), nil)
	}

	if client.breaker.CompletionCallback(err) {
		client.breaker.MarkSuccessful()
	} else {
		client.breaker.MarkFailure()
	}

	return removed
}

func (client *memdClient) SendRequest(req *memdQRequest) error {
	if !client.breaker.AllowsRequest() {
		logSchedf("Circuit breaker interrupting request. %s to %s OP=0x%x. Opaque=%d", client.conn.LocalAddr(), client.Address(), req.Command, req.Opaque)

		req.cancelWithCallback(errCircuitBreakerOpen)

		return nil
	}

	return client.internalSendRequest(req)
}

func (client *memdClient) internalSendRequest(req *memdQRequest) error {
	if err := client.takeRequestOwnership(req); err != nil {
		return err
	}

	packet := &req.Packet
	if client.SupportsFeature(memd.FeatureSnappy) {
		isCompressed := (packet.Datatype & uint8(memd.DatatypeFlagCompressed)) != 0
		packetSize := len(packet.Value)
		if !isCompressed && packetSize > client.compressionMinSize && isCompressibleOp(packet.Command) {
			compressedValue := snappy.Encode(nil, packet.Value)
			if float64(len(compressedValue))/float64(packetSize) <= client.compressionMinRatio {
				newPacket := *packet
				newPacket.Value = compressedValue
				newPacket.Datatype = newPacket.Datatype | uint8(memd.DatatypeFlagCompressed)
				packet = &newPacket
			}
		}
	}

	logSchedf("Writing request. %s to %s OP=0x%x. Opaque=%d", client.conn.LocalAddr(), client.Address(), req.Command, req.Opaque)

	client.tracer.StartNetTrace(req)

	err := client.conn.WritePacket(packet)
	if err != nil {
		logDebugf("memdClient write failure: %v", err)
		return err
	}

	return nil
}

func (client *memdClient) resolveRequest(resp *memdQResponse) {
	defer memd.ReleasePacket(resp.Packet)

	logSchedf("Handling response data. OP=0x%x. Opaque=%d. Status:%d", resp.Command, resp.Opaque, resp.Status)

	client.lock.Lock()
	// Find the request that goes with this response, don't check if the client is
	// closed so that we can handle orphaned responses.
	req := client.opList.FindAndMaybeRemove(resp.Opaque, resp.Status != memd.StatusSuccess)
	client.lock.Unlock()

	if req == nil {
		// There is no known request that goes with this response.  Ignore it.
		logDebugf("Received response with no corresponding request.")
		if client.zombieLogger != nil {
			client.zombieLogger.RecordZombieResponse(resp, client.connID, client.LocalAddress(), client.Address())
		}
		return
	}

	if !req.Persistent || resp.Status != memd.StatusSuccess {
		atomic.CompareAndSwapPointer(&req.waitingIn, unsafe.Pointer(client), nil)
	}

	req.processingLock.Lock()

	if !req.Persistent {
		stopNetTrace(req, resp, client.conn.LocalAddr(), client.conn.RemoteAddr())
	}

	isCompressed := (resp.Datatype & uint8(memd.DatatypeFlagCompressed)) != 0
	if isCompressed && !client.disableDecompression {
		newValue, err := snappy.Decode(nil, resp.Value)
		if err != nil {
			req.processingLock.Unlock()
			logDebugf("Failed to decompress value from the server for key `%s`.", req.Key)
			return
		}

		resp.Value = newValue
		resp.Datatype = resp.Datatype & ^uint8(memd.DatatypeFlagCompressed)
	}

	// Give the agent an opportunity to intercept the response first
	var err error
	if resp.Magic == memd.CmdMagicRes && resp.Status != memd.StatusSuccess {
		err = getKvStatusCodeError(resp.Status)
	}

	if client.breaker.CompletionCallback(err) {
		client.breaker.MarkSuccessful()
	} else {
		client.breaker.MarkFailure()
	}

	if !req.Persistent {
		stopCmdTrace(req)
	}

	req.processingLock.Unlock()

	if err != nil {
		shortCircuited, routeErr := client.postErrHandler(resp, req, err)
		if shortCircuited {
			logSchedf("Routing callback intercepted response")
			return
		}
		err = routeErr
	}

	// Call the requests callback handler...
	logSchedf("Dispatching response callback. OP=0x%x. Opaque=%d", resp.Command, resp.Opaque)
	req.tryCallback(resp, err)
}

func (client *memdClient) run() {
	var (
		// A queue for DCP commands so we can execute them out-of-band from packet receiving.  This
		// is integral to allow the higher level application to back-pressure against the DCP packet
		// processing without interfeering with the SDKs control commands (like config fetch).
		dcpBufferQ = make(chan *dcpBuffer, client.dcpQueueSize)

		// When a kill request comes in, we need to immediately stop processing all requests.  This
		// includes immediately stopping the DCP queue rather than waiting for the application to
		// flush that queue.  This means that we lose packets that were read but not processed, but
		// this is not fundementally different to if we had just not read them at all.  As a side
		// effect of this, we need to use a separate kill signal on top of closing the queue.
		isShuttingDown = uint32(0)

		// After we signal that DCP processing should stop, we need a notification so we know when
		// it has been completed, we do this to prevent leaving the goroutine around, and we need to
		// ensure that the application has finished with the last packet it received before we stop.
		dcpProcDoneCh = make(chan struct{})
	)

	go func() {
		defer close(dcpProcDoneCh)

		for {
			q, stillOpen := <-dcpBufferQ
			if !stillOpen || atomic.LoadUint32(&isShuttingDown) != 0 {
				return
			}

			logSchedf("Resolving response OP=0x%x. Opaque=%d", q.resp.Command, q.resp.Opaque)
			client.resolveRequest(q.resp)

			// See below for information on MB-26363 for why this is here.
			if !q.isInternal && client.dcpAckSize > 0 {
				client.maybeSendDcpBufferAck(q.packetLen)
			}
		}
	}()

	go func() {
		for {
			packet, n, err := client.conn.ReadPacket()
			if err != nil {
				client.lock.Lock()
				if !client.closed {
					logWarnf("memdClient read failure on conn `%v` : %v", client.connID, err)
				}
				client.lock.Unlock()
				break
			}

			resp := &memdQResponse{
				sourceAddr:   client.conn.RemoteAddr(),
				sourceConnID: client.connID,
				Packet:       packet,
			}

			atomic.StoreInt64(&client.lastActivity, time.Now().UnixNano())

			// We handle DCP no-op's directly here so we can reply immediately.
			if resp.Packet.Command == memd.CmdDcpNoop {
				err := client.conn.WritePacket(&memd.Packet{
					Magic:   memd.CmdMagicRes,
					Command: memd.CmdDcpNoop,
					Opaque:  resp.Opaque,
				})
				if err != nil {
					logWarnf("Failed to dispatch DCP noop reply: %s", err)
				}
				continue
			}

			// This is a fix for a bug in the server DCP implementation (MB-26363).  This
			// bug causes the server to fail to send a stream-end notification.  The server
			// does however synchronously stop the stream, and thus we can assume no more
			// packets will be received following the close response.
			if resp.Magic == memd.CmdMagicRes && resp.Command == memd.CmdDcpCloseStream && client.streamEndNotSupported {
				closeReq := client.opList.Find(resp.Opaque)
				if closeReq != nil {
					vbID := closeReq.Vbucket
					streamReq := client.opList.FindOpenStream(vbID)
					if streamReq != nil {
						endExtras := make([]byte, 4)
						binary.BigEndian.PutUint32(endExtras, uint32(memd.StreamEndClosed))
						endResp := &memdQResponse{
							Packet: &memd.Packet{
								Magic:   memd.CmdMagicReq,
								Command: memd.CmdDcpStreamEnd,
								Vbucket: vbID,
								Opaque:  streamReq.Opaque,
								Extras:  endExtras,
							},
						}
						dcpBufferQ <- &dcpBuffer{
							resp:       endResp,
							packetLen:  n,
							isInternal: true,
						}
					}
				}
			}

			switch resp.Packet.Command {
			case memd.CmdDcpDeletion, memd.CmdDcpExpiration, memd.CmdDcpMutation, memd.CmdDcpSnapshotMarker,
				memd.CmdDcpEvent, memd.CmdDcpOsoSnapshot, memd.CmdDcpSeqNoAdvanced, memd.CmdDcpStreamEnd:
				dcpBufferQ <- &dcpBuffer{
					resp:      resp,
					packetLen: n,
				}
			default:
				logSchedf("Resolving response OP=0x%x. Opaque=%d", resp.Command, resp.Opaque)
				client.resolveRequest(resp)
			}
		}

		var closedError error
		client.lock.Lock()
		if !client.closed {
			client.closed = true
			client.lock.Unlock()

			err := client.conn.Close()
			if err != nil {
				// Lets log a warning, as this is non-fatal
				logWarnf("Failed to shut down client connection (%s)", err)
			}
		} else {
			closedError = client.closedError
			client.lock.Unlock()
		}

		// We first mark that we are shutting down to stop the DCP processor from running any
		// additional packets up to the application.  We then close the buffer channel to wake
		// the processor if its asleep (queue was empty).  We then wait to ensure it is finished
		// with whatever packet was being processed.
		atomic.StoreUint32(&isShuttingDown, 1)
		close(dcpBufferQ)
		<-dcpProcDoneCh

		if closedError == nil {
			closedError = io.EOF
		}

		client.opList.Drain(func(req *memdQRequest) {
			if !atomic.CompareAndSwapPointer(&req.waitingIn, unsafe.Pointer(client), nil) {
				logWarnf("Encountered an unowned request in a client opMap")
			}

			shortCircuited, routeErr := client.postErrHandler(nil, req, closedError)
			if shortCircuited {
				return
			}

			req.tryCallback(nil, routeErr)
		})

		close(client.closeNotify)
	}()
}

func (client *memdClient) LocalAddress() string {
	return client.conn.LocalAddr()
}

func (client *memdClient) Close(err error) error {
	client.lock.Lock()
	client.closed = true
	client.closedError = err
	client.lock.Unlock()

	return client.conn.Close()
}

func (client *memdClient) sendCanary() {
	errChan := make(chan error)
	handler := func(resp *memdQResponse, req *memdQRequest, err error) {
		errChan <- err
	}

	req := &memdQRequest{
		Packet: memd.Packet{
			Magic:    memd.CmdMagicReq,
			Command:  memd.CmdNoop,
			Datatype: 0,
			Cas:      0,
			Key:      nil,
			Value:    nil,
		},
		Callback:      handler,
		RetryStrategy: newFailFastRetryStrategy(),
	}

	logDebugf("Sending NOOP request for %p/%s", client, client.Address())
	err := client.internalSendRequest(req)
	if err != nil {
		client.breaker.MarkFailure()
	}

	timer := AcquireTimer(client.breaker.CanaryTimeout())
	select {
	case <-timer.C:
		if !req.internalCancel(errRequestCanceled) {
			err := <-errChan
			if err == nil {
				logDebugf("NOOP request successful for %p/%s", client, client.Address())
				client.breaker.MarkSuccessful()
			} else {
				logDebugf("NOOP request failed for %p/%s", client, client.Address())
				client.breaker.MarkFailure()
			}
		}
		client.breaker.MarkFailure()
	case err := <-errChan:
		if err == nil {
			client.breaker.MarkSuccessful()
		} else {
			client.breaker.MarkFailure()
		}
	}
}

func (client *memdClient) helloFeatures(props helloProps) []memd.HelloFeature {
	var features []memd.HelloFeature

	// Send the TLS flag, which has unknown effects.
	features = append(features, memd.FeatureTLS)

	// Indicate that we understand XATTRs
	features = append(features, memd.FeatureXattr)

	// Indicates that we understand select buckets.
	features = append(features, memd.FeatureSelectBucket)

	// If the user wants to use KV Error maps, lets enable them
	if props.XErrorFeatureEnabled {
		features = append(features, memd.FeatureXerror)
	}

	// Indicate that we understand JSON
	if props.JSONFeatureEnabled {
		features = append(features, memd.FeatureJSON)
	}

	// Indicate that we understand Point in Time
	if props.PITRFeatureEnabled {
		features = append(features, memd.FeaturePITR)
	}

	// If the user wants to use mutation tokens, lets enable them
	if props.MutationTokensEnabled {
		features = append(features, memd.FeatureSeqNo)
	}

	// If the user wants on-the-wire compression, lets try to enable it
	if props.CompressionEnabled {
		features = append(features, memd.FeatureSnappy)
	}

	if props.DurationsEnabled {
		features = append(features, memd.FeatureDurations)
	}

	if props.CollectionsEnabled {
		features = append(features, memd.FeatureCollections)
	}

	if props.OutOfOrderEnabled {
		features = append(features, memd.FeatureUnorderedExec)
	}

	// These flags are informational so don't actually enable anything
	features = append(features, memd.FeatureAltRequests)
	features = append(features, memd.FeatureCreateAsDeleted)
	features = append(features, memd.FeatureReplaceBodyWithXattr)
	features = append(features, memd.FeaturePreserveExpiry)

	if props.SyncReplicationEnabled {
		features = append(features, memd.FeatureSyncReplication)
	}

	return features
}

type helloProps struct {
	MutationTokensEnabled  bool
	CollectionsEnabled     bool
	CompressionEnabled     bool
	DurationsEnabled       bool
	OutOfOrderEnabled      bool
	JSONFeatureEnabled     bool
	XErrorFeatureEnabled   bool
	SyncReplicationEnabled bool
	PITRFeatureEnabled     bool
}

type bootstrapProps struct {
	Bucket         string
	UserAgent      string
	AuthMechanisms []AuthMechanism
	AuthHandler    authFuncHandler
	ErrMapManager  *errMapComponent
	HelloProps     helloProps
}

type memdInitFunc func(*memdClient, time.Time) error

func (client *memdClient) Bootstrap(cancelSig <-chan struct{}, settings bootstrapProps, deadline time.Time, cb memdInitFunc) error {
	logDebugf("Memdclient `%s/%p` Fetching cluster client data", client.Address(), client)

	bucket := settings.Bucket
	features := client.helloFeatures(settings.HelloProps)
	clientInfoStr := clientInfoString(client.connID, settings.UserAgent)
	authMechanisms := settings.AuthMechanisms
	client.cancelBootstrapSig = cancelSig

	helloCh, err := client.ExecHello(clientInfoStr, features, deadline)
	if err != nil {
		logDebugf("Memdclient `%s/%p` Failed to execute HELLO (%v)", client.Address(), client, err)
		return err
	}

	errMapCh, err := client.ExecGetErrorMap(1, deadline)
	if err != nil {
		// GetErrorMap isn't integral to bootstrap succeeding
		logDebugf("Memdclient `%s/%p`Failed to execute Get error map (%v)", client.Address(), client, err)
	}

	var listMechsCh chan SaslListMechsCompleted
	firstAuthMethod := settings.AuthHandler(client, deadline, authMechanisms[0])
	// If the auth method is nil then we don't actually need to do any auth so no need to Get the mechanisms.
	if firstAuthMethod != nil {
		listMechsCh = make(chan SaslListMechsCompleted, 1)
		err = client.SaslListMechs(deadline, func(mechs []AuthMechanism, err error) {
			if err != nil {
				logDebugf("Memdclient `%s/%p` Failed to fetch list auth mechs (%v)", client.Address(), client, err)
			}
			listMechsCh <- SaslListMechsCompleted{
				Err:   err,
				Mechs: mechs,
			}
		})
		if err != nil {
			logDebugf("Memdclient `%s/%p` Failed to execute list auth mechs (%v)", client.Address(), client, err)
		}
	}

	var completedAuthCh chan BytesAndError
	var continueAuthCh chan bool
	if firstAuthMethod != nil {
		completedAuthCh, continueAuthCh, err = firstAuthMethod()
		if err != nil {
			logDebugf("Memdclient `%s/%p` Failed to execute auth (%v)", client.Address(), client, err)
			return err
		}
	}

	var selectCh chan BytesAndError
	if continueAuthCh == nil {
		if bucket != "" {
			selectCh, err = client.ExecSelectBucket([]byte(bucket), deadline)
			if err != nil {
				logDebugf("Memdclient `%s/%p` Failed to execute select bucket (%v)", client.Address(), client, err)
				return err
			}
		}
	} else {
		selectCh = client.continueAfterAuth(bucket, continueAuthCh, deadline)
	}

	helloResp := <-helloCh
	if helloResp.Err != nil {
		logDebugf("Memdclient `%s/%p` Failed to hello with server (%v)", client.Address(), client, helloResp.Err)
		return helloResp.Err
	}

	if errMapCh != nil {
		errMapResp := <-errMapCh
		if errMapResp.Err == nil {
			settings.ErrMapManager.StoreErrorMap(errMapResp.Bytes)
		} else {
			logDebugf("Memdclient `%s/%p` Failed to fetch kv error map (%s)", client.Address(), client, errMapResp.Err)
		}
	}

	var serverAuthMechanisms []AuthMechanism
	if listMechsCh != nil {
		listMechsResp := <-listMechsCh
		if listMechsResp.Err == nil {
			serverAuthMechanisms = listMechsResp.Mechs
			logDebugf("Memdclient `%s/%p` Server supported auth mechanisms: %v", client.Address(), client, serverAuthMechanisms)
		} else {
			logDebugf("Memdclient `%s/%p` Failed to fetch auth mechs from server (%v)", client.Address(), client, listMechsResp.Err)
		}
	}

	// If completedAuthCh isn't nil then we have attempted to do auth so we need to wait on the result of that.
	if completedAuthCh != nil {
		authResp := <-completedAuthCh
		if authResp.Err != nil {
			logDebugf("Memdclient `%s/%p` Failed to perform auth against server (%v)", client.Address(), client, authResp.Err)
			if errors.Is(authResp.Err, ErrRequestCanceled) {
				// There's no point in us trying different mechanisms if something has cancelled bootstrapping.
				return authResp.Err
			} else if errors.Is(authResp.Err, ErrAuthenticationFailure) {
				// If there's only one auth mechanism then we can just fail.
				if len(authMechanisms) == 1 {
					return authResp.Err
				}
				// If the server supports the mechanism we've tried then this auth error can't be due to an unsupported
				// mechanism.
				for _, mech := range serverAuthMechanisms {
					if mech == authMechanisms[0] {
						return authResp.Err
					}
				}

				// If we've got here then the auth mechanism we tried is unsupported so let's keep trying with the next
				// supported mechanism.
				logInfof("Memdclient `%p` Unsupported authentication mechanism, will attempt to find next supported mechanism", client)
			}

			for {
				var found bool
				var mech AuthMechanism
				found, mech, authMechanisms = findNextAuthMechanism(authMechanisms, serverAuthMechanisms)
				if !found {
					logDebugf("Memdclient `%s/%p` Failed to authenticate, all options exhausted", client.Address(), client)
					return authResp.Err
				}

				logDebugf("Memdclient `%s/%p` Retrying authentication with found supported mechanism: %s", client.Address(), client, mech)
				nextAuthFunc := settings.AuthHandler(client, deadline, mech)
				if nextAuthFunc == nil {
					// This can't really happen but just in case it somehow does.
					logInfof("Memdclient `%p` Failed to authenticate, no available credentials", client)
					return authResp.Err
				}
				completedAuthCh, continueAuthCh, err = nextAuthFunc()
				if err != nil {
					logDebugf("Memdclient `%s/%p` Failed to execute auth (%v)", client.Address(), client, err)
					return err
				}
				if continueAuthCh == nil {
					if bucket != "" {
						selectCh, err = client.ExecSelectBucket([]byte(bucket), deadline)
						if err != nil {
							logDebugf("Memdclient `%s/%p` Failed to execute select bucket (%v)", client.Address(), client, err)
							return err
						}
					}
				} else {
					selectCh = client.continueAfterAuth(bucket, continueAuthCh, deadline)
				}
				authResp = <-completedAuthCh
				if authResp.Err == nil {
					break
				}

				logDebugf("Memdclient `%s/%p` Failed to perform auth against server (%v)", client.Address(), client, authResp.Err)
				if errors.Is(authResp.Err, ErrAuthenticationFailure) || errors.Is(err, ErrRequestCanceled) {
					return authResp.Err
				}
			}
		}
		logDebugf("Memdclient `%s/%p` Authenticated successfully", client.Address(), client)
	}

	if selectCh != nil {
		selectResp := <-selectCh
		if selectResp.Err != nil {
			logDebugf("Memdclient `%s/%p` Failed to perform select bucket against server (%v)", client.Address(), client, selectResp.Err)
			return selectResp.Err
		}
	}

	client.features = helloResp.SrvFeatures

	logDebugf("Memdclient `%s/%p` Client Features: %+v", client.Address(), client, features)
	logDebugf("Memdclient `%s/%p` Server Features: %+v", client.Address(), client, client.features)

	for _, feature := range client.features {
		client.conn.EnableFeature(feature)
	}

	err = cb(client, deadline)
	if err != nil {
		return err
	}

	return nil
}

// BytesAndError contains the raw bytes of the result of an operation, and/or the error that occurred.
type BytesAndError struct {
	Err   error
	Bytes []byte
}

func (client *memdClient) SaslAuth(k, v []byte, deadline time.Time, cb func(b []byte, err error)) error {
	err := client.doBootstrapRequest(
		&memdQRequest{
			Packet: memd.Packet{
				Magic:   memd.CmdMagicReq,
				Command: memd.CmdSASLAuth,
				Key:     k,
				Value:   v,
			},
			Callback: func(resp *memdQResponse, _ *memdQRequest, err error) {
				// Auth is special, auth continue is surfaced as an error
				var val []byte
				if resp != nil {
					val = resp.Value
				}

				cb(val, err)
			},
			RetryStrategy: newFailFastRetryStrategy(),
		},
		deadline,
	)
	if err != nil {
		return err
	}

	return nil
}

func (client *memdClient) SaslStep(k, v []byte, deadline time.Time, cb func(err error)) error {
	err := client.doBootstrapRequest(
		&memdQRequest{
			Packet: memd.Packet{
				Magic:   memd.CmdMagicReq,
				Command: memd.CmdSASLStep,
				Key:     k,
				Value:   v,
			},
			Callback: func(resp *memdQResponse, _ *memdQRequest, err error) {
				if err != nil {
					cb(err)
					return
				}

				cb(nil)
			},
			RetryStrategy: newFailFastRetryStrategy(),
		},
		deadline,
	)
	if err != nil {
		return err
	}

	return nil
}

func (client *memdClient) ExecSelectBucket(b []byte, deadline time.Time) (chan BytesAndError, error) {
	completedCh := make(chan BytesAndError, 1)
	err := client.doBootstrapRequest(
		&memdQRequest{
			Packet: memd.Packet{
				Magic:   memd.CmdMagicReq,
				Command: memd.CmdSelectBucket,
				Key:     b,
			},
			Callback: func(resp *memdQResponse, _ *memdQRequest, err error) {
				if err != nil {
					if errors.Is(err, ErrDocumentNotFound) {
						// Bucket not found means that the user has priviledges to access the bucket but that the bucket
						// is in some way not existing right now (e.g. in warmup).
						err = errBucketNotFound
					}
					completedCh <- BytesAndError{
						Err: err,
					}
					return
				}

				completedCh <- BytesAndError{
					Bytes: resp.Value,
				}
			},
			RetryStrategy: newFailFastRetryStrategy(),
		},
		deadline,
	)
	if err != nil {
		return nil, err
	}

	return completedCh, nil
}

func (client *memdClient) ExecGetErrorMap(version uint16, deadline time.Time) (chan BytesAndError, error) {
	completedCh := make(chan BytesAndError, 1)
	valueBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(valueBuf, version)

	err := client.doBootstrapRequest(
		&memdQRequest{
			Packet: memd.Packet{
				Magic:   memd.CmdMagicReq,
				Command: memd.CmdGetErrorMap,
				Value:   valueBuf,
			},
			Callback: func(resp *memdQResponse, _ *memdQRequest, err error) {
				if err != nil {
					completedCh <- BytesAndError{
						Err: err,
					}
					return
				}

				completedCh <- BytesAndError{
					Bytes: resp.Value,
				}
			},
			RetryStrategy: newFailFastRetryStrategy(),
		},
		deadline,
	)
	if err != nil {
		return nil, err
	}

	return completedCh, nil
}

func (client *memdClient) SaslListMechs(deadline time.Time, cb func(mechs []AuthMechanism, err error)) error {
	err := client.doBootstrapRequest(
		&memdQRequest{
			Packet: memd.Packet{
				Magic:   memd.CmdMagicReq,
				Command: memd.CmdSASLListMechs,
			},
			Callback: func(resp *memdQResponse, _ *memdQRequest, err error) {
				if err != nil {
					cb(nil, err)
					return
				}

				mechs := strings.Split(string(resp.Value), " ")
				var authMechs []AuthMechanism
				for _, mech := range mechs {
					authMechs = append(authMechs, AuthMechanism(mech))
				}

				cb(authMechs, nil)
			},
			RetryStrategy: newFailFastRetryStrategy(),
		},
		deadline,
	)
	if err != nil {
		return err
	}

	return nil
}

// ExecHelloResponse contains the features and/or error from an ExecHello operation.
type ExecHelloResponse struct {
	SrvFeatures []memd.HelloFeature
	Err         error
}

func (client *memdClient) ExecHello(clientID string, features []memd.HelloFeature, deadline time.Time) (chan ExecHelloResponse, error) {
	appendFeatureCode := func(bytes []byte, feature memd.HelloFeature) []byte {
		bytes = append(bytes, 0, 0)
		binary.BigEndian.PutUint16(bytes[len(bytes)-2:], uint16(feature))
		return bytes
	}

	var featureBytes []byte
	for _, feature := range features {
		featureBytes = appendFeatureCode(featureBytes, feature)
	}

	completedCh := make(chan ExecHelloResponse, 1)
	err := client.doBootstrapRequest(
		&memdQRequest{
			Packet: memd.Packet{
				Magic:   memd.CmdMagicReq,
				Command: memd.CmdHello,
				Key:     []byte(clientID),
				Value:   featureBytes,
			},
			Callback: func(resp *memdQResponse, _ *memdQRequest, err error) {
				if err != nil {
					completedCh <- ExecHelloResponse{
						Err: err,
					}
					return
				}

				var srvFeatures []memd.HelloFeature
				for i := 0; i < len(resp.Value); i += 2 {
					feature := binary.BigEndian.Uint16(resp.Value[i:])
					srvFeatures = append(srvFeatures, memd.HelloFeature(feature))
				}

				completedCh <- ExecHelloResponse{
					SrvFeatures: srvFeatures,
				}
			},
			RetryStrategy: newFailFastRetryStrategy(),
		},
		deadline,
	)
	if err != nil {
		return nil, err
	}

	return completedCh, nil
}

func (client *memdClient) doBootstrapRequest(req *memdQRequest, deadline time.Time) error {
	origCb := req.Callback
	doneCh := make(chan struct{})
	handler := func(resp *memdQResponse, req *memdQRequest, err error) {
		close(doneCh)
		origCb(resp, req, err)
	}

	req.Callback = handler
	start := time.Now()
	req.SetTimer(time.AfterFunc(deadline.Sub(start), func() {
		connInfo := req.ConnectionInfo()
		count, reasons := req.Retries()
		req.cancelWithCallback(&TimeoutError{
			InnerError:         errAmbiguousTimeout,
			OperationID:        req.Command.Name(),
			Opaque:             req.Identifier(),
			TimeObserved:       time.Since(start),
			RetryReasons:       reasons,
			RetryAttempts:      count,
			LastDispatchedTo:   connInfo.lastDispatchedTo,
			LastDispatchedFrom: connInfo.lastDispatchedFrom,
			LastConnectionID:   connInfo.lastConnectionID,
		})
	}))

	go func() {
		select {
		case <-doneCh:
			return
		case <-client.cancelBootstrapSig:
			req.Cancel()
			<-doneCh
			return
		}
	}()

	err := client.SendRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (client *memdClient) continueAfterAuth(bucketName string, continueAuthCh chan bool, deadline time.Time) chan BytesAndError {
	if bucketName == "" {
		return nil
	}

	selectCh := make(chan BytesAndError, 1)
	go func() {
		success := <-continueAuthCh
		if !success {
			selectCh <- BytesAndError{}
			return
		}
		execCh, err := client.ExecSelectBucket([]byte(bucketName), deadline)
		if err != nil {
			logDebugf("Memdclient `%s/%p` Failed to execute select bucket (%v)", client.Address(), client, err)
			selectCh <- BytesAndError{Err: err}
			return
		}

		execResp := <-execCh
		selectCh <- execResp
	}()

	return selectCh
}

func checkSupportsFeature(srvFeatures []memd.HelloFeature, feature memd.HelloFeature) bool {
	for _, srvFeature := range srvFeatures {
		if srvFeature == feature {
			return true
		}
	}
	return false
}

func findNextAuthMechanism(authMechanisms []AuthMechanism, serverAuthMechanisms []AuthMechanism) (bool, AuthMechanism, []AuthMechanism) {
	for {
		if len(authMechanisms) <= 1 {
			break
		}
		authMechanisms = authMechanisms[1:]
		mech := authMechanisms[0]
		for _, serverMech := range serverAuthMechanisms {
			if mech == serverMech {
				return true, mech, authMechanisms
			}
		}
	}

	return false, "", authMechanisms
}
