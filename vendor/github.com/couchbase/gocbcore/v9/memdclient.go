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

	"github.com/couchbase/gocbcore/v9/memd"

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
	conn                  memdConn
	opList                memdOpMap
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

func (client *memdClient) takeRequestOwnership(req *memdQRequest) bool {
	client.lock.Lock()
	defer client.lock.Unlock()

	if client.closed {
		logDebugf("Attempted to put dispatched op in drained opmap")
		return false
	}

	if !atomic.CompareAndSwapPointer(&req.waitingIn, nil, unsafe.Pointer(client)) {
		logDebugf("Attempted to put dispatched op in new opmap")
		return false
	}

	if req.isCancelled() {
		atomic.CompareAndSwapPointer(&req.waitingIn, unsafe.Pointer(client), nil)
		return false
	}

	connInfo := memdQRequestConnInfo{
		lastDispatchedTo:   client.Address(),
		lastDispatchedFrom: client.conn.LocalAddr(),
		lastConnectionID:   client.connID,
	}
	req.SetConnectionInfo(connInfo)

	client.opList.Add(req)
	return true
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
	addSuccess := client.takeRequestOwnership(req)
	if !addSuccess {
		return errRequestCanceled
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
	opIndex := resp.Opaque

	logSchedf("Handling response data. OP=0x%x. Opaque=%d. Status:%d", resp.Command, resp.Opaque, resp.Status)

	client.lock.Lock()
	// Find the request that goes with this response, don't check if the client is
	// closed so that we can handle orphaned responses.
	req := client.opList.FindAndMaybeRemove(opIndex, resp.Status != memd.StatusSuccess)
	client.lock.Unlock()

	if req == nil {
		// There is no known request that goes with this response.  Ignore it.
		logDebugf("Received response with no corresponding request.")
		if client.zombieLogger != nil {
			client.zombieLogger.RecordZombieResponse(resp, client.connID, client.Address())
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
	dcpBufferQ := make(chan *dcpBuffer, client.dcpQueueSize)
	dcpKillSwitch := make(chan bool)
	dcpKillNotify := make(chan bool)
	go func() {
		procDcpItem := func(q *dcpBuffer, more bool) bool {
			if !more {
				dcpKillNotify <- true
				return false
			}

			logSchedf("Resolving response OP=0x%x. Opaque=%d", q.resp.Command, q.resp.Opaque)
			client.resolveRequest(q.resp)

			// See below for information on MB-26363 for why this is here.
			if !q.isInternal && client.dcpAckSize > 0 {
				client.maybeSendDcpBufferAck(q.packetLen)
			}

			return true
		}

		for {
			select {
			case q, more := <-dcpBufferQ:
				if !procDcpItem(q, more) {
					return
				}
			case <-dcpKillSwitch:
				close(dcpBufferQ)
			}
		}
	}()

	go func() {
		for {
			packet, n, err := client.conn.ReadPacket()
			if err != nil {
				if !client.closed {
					logErrorf("memdClient read failure: %v", err)
				}
				break
			}

			resp := &memdQResponse{
				sourceAddr:   client.conn.RemoteAddr(),
				sourceConnID: client.connID,
				Packet:       *packet,
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
							Packet: memd.Packet{
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
			case memd.CmdDcpDeletion:
				fallthrough
			case memd.CmdDcpExpiration:
				fallthrough
			case memd.CmdDcpMutation:
				fallthrough
			case memd.CmdDcpSnapshotMarker:
				fallthrough
			case memd.CmdDcpEvent:
				fallthrough
			case memd.CmdDcpOsoSnapshot:
				fallthrough
			case memd.CmdDcpSeqNoAdvanced:
				fallthrough
			case memd.CmdDcpStreamEnd:
				dcpBufferQ <- &dcpBuffer{
					resp:      resp,
					packetLen: n,
				}
				continue
			default:
				logSchedf("Resolving response OP=0x%x. Opaque=%d", resp.Command, resp.Opaque)
				client.resolveRequest(resp)
			}
		}

		client.lock.Lock()
		if client.closed {
			client.lock.Unlock()
		} else {
			client.closed = true
			client.lock.Unlock()

			err := client.conn.Close()
			if err != nil {
				// Lets log an error, as this is non-fatal
				logErrorf("Failed to shut down client connection (%s)", err)
			}
		}

		dcpKillSwitch <- true
		<-dcpKillNotify

		client.opList.Drain(func(req *memdQRequest) {
			if !atomic.CompareAndSwapPointer(&req.waitingIn, unsafe.Pointer(client), nil) {
				logWarnf("Encountered an unowned request in a client opMap")
			}

			shortCircuited, routeErr := client.postErrHandler(nil, req, io.EOF)
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

func (client *memdClient) Close() error {
	client.lock.Lock()
	client.closed = true
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
	logDebugf("Fetching cluster client data")

	bucket := settings.Bucket
	features := client.helloFeatures(settings.HelloProps)
	clientInfoStr := clientInfoString(client.connID, settings.UserAgent)
	authMechanisms := settings.AuthMechanisms
	client.cancelBootstrapSig = cancelSig

	helloCh, err := client.ExecHello(clientInfoStr, features, deadline)
	if err != nil {
		logDebugf("Failed to execute HELLO (%v)", err)
		return err
	}

	errMapCh, err := client.ExecGetErrorMap(1, deadline)
	if err != nil {
		// GetErrorMap isn't integral to bootstrap succeeding
		logDebugf("Failed to execute Get error map (%v)", err)
	}

	var listMechsCh chan SaslListMechsCompleted
	firstAuthMethod := settings.AuthHandler(client, deadline, authMechanisms[0])
	// If the auth method is nil then we don't actually need to do any auth so no need to Get the mechanisms.
	if firstAuthMethod != nil {
		listMechsCh = make(chan SaslListMechsCompleted, 1)
		err = client.SaslListMechs(deadline, func(mechs []AuthMechanism, err error) {
			if err != nil {
				logDebugf("Failed to fetch list auth mechs (%v)", err)
			}
			listMechsCh <- SaslListMechsCompleted{
				Err:   err,
				Mechs: mechs,
			}
		})
		if err != nil {
			logDebugf("Failed to execute list auth mechs (%v)", err)
		}
	}

	var completedAuthCh chan BytesAndError
	var continueAuthCh chan bool
	if firstAuthMethod != nil {
		completedAuthCh, continueAuthCh, err = firstAuthMethod()
		if err != nil {
			logDebugf("Failed to execute auth (%v)", err)
			return err
		}
	}

	var selectCh chan BytesAndError
	if continueAuthCh == nil {
		if bucket != "" {
			selectCh, err = client.ExecSelectBucket([]byte(bucket), deadline)
			if err != nil {
				logDebugf("Failed to execute select bucket (%v)", err)
				return err
			}
		}
	} else {
		selectCh = client.continueAfterAuth(bucket, continueAuthCh, deadline)
	}

	helloResp := <-helloCh
	if helloResp.Err != nil {
		logDebugf("Failed to hello with server (%v)", helloResp.Err)
		return helloResp.Err
	}

	errMapResp := <-errMapCh
	if errMapResp.Err == nil {
		settings.ErrMapManager.StoreErrorMap(errMapResp.Bytes)
	} else {
		logDebugf("Failed to fetch kv error map (%s)", errMapResp.Err)
	}

	var serverAuthMechanisms []AuthMechanism
	if listMechsCh != nil {
		listMechsResp := <-listMechsCh
		if listMechsResp.Err == nil {
			serverAuthMechanisms = listMechsResp.Mechs
			logDebugf("Server supported auth mechanisms: %v", serverAuthMechanisms)
		} else {
			logDebugf("Failed to fetch auth mechs from server (%v)", listMechsResp.Err)
		}
	}

	// If completedAuthCh isn't nil then we have attempted to do auth so we need to wait on the result of that.
	if completedAuthCh != nil {
		authResp := <-completedAuthCh
		if authResp.Err != nil {
			logDebugf("Failed to perform auth against server (%v)", authResp.Err)
			if errors.Is(authResp.Err, ErrAuthenticationFailure) {
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
				logDebugf("Unsupported authentication mechanism, will attempt to find next supported mechanism")
			}

			for {
				var found bool
				var mech AuthMechanism
				found, mech, authMechanisms = findNextAuthMechanism(authMechanisms, serverAuthMechanisms)
				if !found {
					logDebugf("Failed to authenticate, all options exhausted")
					return authResp.Err
				}

				logDebugf("Retrying authentication with found supported mechanism: %s", mech)
				nextAuthFunc := settings.AuthHandler(client, deadline, mech)
				if nextAuthFunc == nil {
					// This can't really happen but just in case it somehow does.
					logDebugf("Failed to authenticate, no available credentials")
					return authResp.Err
				}
				completedAuthCh, continueAuthCh, err = nextAuthFunc()
				if err != nil {
					logDebugf("Failed to execute auth (%v)", err)
					return err
				}
				if continueAuthCh == nil {
					if bucket != "" {
						selectCh, err = client.ExecSelectBucket([]byte(bucket), deadline)
						if err != nil {
							logDebugf("Failed to execute select bucket (%v)", err)
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

				logDebugf("Failed to perform auth against server (%v)", authResp.Err)
				if errors.Is(authResp.Err, ErrAuthenticationFailure) {
					return authResp.Err
				}
			}
		}
		logDebugf("Authenticated successfully")
	}

	if selectCh != nil {
		selectResp := <-selectCh
		if selectResp.Err != nil {
			logDebugf("Failed to perform select bucket against server (%v)", selectResp.Err)
			return selectResp.Err
		}
	}

	client.features = helloResp.SrvFeatures

	logDebugf("Client Features: %+v", features)
	logDebugf("Server Features: %+v", client.features)

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
			logDebugf("Bootstrap cancellation request received")
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
			logDebugf("Failed to execute select bucket (%v)", err)
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
