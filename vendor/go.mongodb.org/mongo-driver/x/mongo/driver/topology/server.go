// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
)

const minHeartbeatInterval = 500 * time.Millisecond

var (
	// ErrServerClosed occurs when an attempt to Get a connection is made after
	// the server has been closed.
	ErrServerClosed = errors.New("server is closed")
	// ErrServerConnected occurs when at attempt to Connect is made after a server
	// has already been connected.
	ErrServerConnected = errors.New("server is connected")

	errCheckCancelled = errors.New("server check cancelled")
	emptyDescription  = description.NewDefaultServer("")
)

// SelectedServer represents a specific server that was selected during server selection.
// It contains the kind of the topology it was selected from.
type SelectedServer struct {
	*Server

	Kind description.TopologyKind
}

// Description returns a description of the server as of the last heartbeat.
func (ss *SelectedServer) Description() description.SelectedServer {
	sdesc := ss.Server.Description()
	return description.SelectedServer{
		Server: sdesc,
		Kind:   ss.Kind,
	}
}

// These constants represent the connection states of a server.
const (
	disconnected int64 = iota
	disconnecting
	connected
	connecting
	initialized
)

func connectionStateString(state int64) string {
	switch state {
	case 0:
		return "Disconnected"
	case 1:
		return "Disconnecting"
	case 2:
		return "Connected"
	case 3:
		return "Connecting"
	case 4:
		return "Initialized"
	}

	return ""
}

// Server is a single server within a topology.
type Server struct {
	cfg             *serverConfig
	address         address.Address
	connectionstate int64

	// connection related fields
	pool *pool

	// goroutine management fields
	done          chan struct{}
	checkNow      chan struct{}
	disconnecting chan struct{}
	closewg       sync.WaitGroup

	// description related fields
	desc                   atomic.Value // holds a description.Server
	updateTopologyCallback atomic.Value
	topologyID             primitive.ObjectID

	// subscriber related fields
	subLock             sync.Mutex
	subscribers         map[uint64]chan description.Server
	currentSubscriberID uint64
	subscriptionsClosed bool

	// heartbeat and cancellation related fields
	// globalCtx should be created in NewServer and cancelled in Disconnect to signal that the server is shutting down.
	// heartbeatCtx should be used for individual heartbeats and should be a child of globalCtx so that it will be
	// cancelled automatically during shutdown.
	heartbeatLock      sync.Mutex
	conn               *connection
	globalCtx          context.Context
	globalCtxCancel    context.CancelFunc
	heartbeatCtx       context.Context
	heartbeatCtxCancel context.CancelFunc

	processErrorLock sync.Mutex
	rttMonitor       *rttMonitor
}

// updateTopologyCallback is a callback used to create a server that should be called when the parent Topology instance
// should be updated based on a new server description. The callback must return the server description that should be
// stored by the server.
type updateTopologyCallback func(description.Server) description.Server

// ConnectServer creates a new Server and then initializes it using the
// Connect method.
func ConnectServer(addr address.Address, updateCallback updateTopologyCallback, topologyID primitive.ObjectID, opts ...ServerOption) (*Server, error) {
	srvr, err := NewServer(addr, topologyID, opts...)
	if err != nil {
		return nil, err
	}
	err = srvr.Connect(updateCallback)
	if err != nil {
		return nil, err
	}
	return srvr, nil
}

// NewServer creates a new server. The mongodb server at the address will be monitored
// on an internal monitoring goroutine.
func NewServer(addr address.Address, topologyID primitive.ObjectID, opts ...ServerOption) (*Server, error) {
	cfg, err := newServerConfig(opts...)
	if err != nil {
		return nil, err
	}

	globalCtx, globalCtxCancel := context.WithCancel(context.Background())
	s := &Server{
		cfg:     cfg,
		address: addr,

		done:          make(chan struct{}),
		checkNow:      make(chan struct{}, 1),
		disconnecting: make(chan struct{}),

		topologyID: topologyID,

		subscribers:     make(map[uint64]chan description.Server),
		globalCtx:       globalCtx,
		globalCtxCancel: globalCtxCancel,
	}
	s.desc.Store(description.NewDefaultServer(addr))
	rttCfg := &rttConfig{
		interval:           cfg.heartbeatInterval,
		createConnectionFn: s.createConnection,
		createOperationFn:  s.createBaseOperation,
	}
	s.rttMonitor = newRttMonitor(rttCfg)

	pc := poolConfig{
		Address:     addr,
		MinPoolSize: cfg.minConns,
		MaxPoolSize: cfg.maxConns,
		MaxIdleTime: cfg.connectionPoolMaxIdleTime,
		PoolMonitor: cfg.poolMonitor,
	}

	connectionOpts := make([]ConnectionOption, len(cfg.connectionOpts))
	copy(connectionOpts, cfg.connectionOpts)
	connectionOpts = append(connectionOpts, withErrorHandlingCallback(s.ProcessHandshakeError))
	s.pool, err = newPool(pc, connectionOpts...)
	if err != nil {
		return nil, err
	}

	s.publishServerOpeningEvent(s.address)

	return s, nil
}

// Connect initializes the Server by starting background monitoring goroutines.
// This method must be called before a Server can be used.
func (s *Server) Connect(updateCallback updateTopologyCallback) error {
	if !atomic.CompareAndSwapInt64(&s.connectionstate, disconnected, connected) {
		return ErrServerConnected
	}

	desc := description.NewDefaultServer(s.address)
	if s.cfg.loadBalanced {
		// LBs automatically start off with kind LoadBalancer because there is no monitoring routine for state changes.
		desc.Kind = description.LoadBalancer
	}
	s.desc.Store(desc)
	s.updateTopologyCallback.Store(updateCallback)

	if !s.cfg.monitoringDisabled && !s.cfg.loadBalanced {
		s.rttMonitor.connect()
		s.closewg.Add(1)
		go s.update()
	}
	return s.pool.connect()
}

// Disconnect closes sockets to the server referenced by this Server.
// Subscriptions to this Server will be closed. Disconnect will shutdown
// any monitoring goroutines, closeConnection the idle connection pool, and will
// wait until all the in use connections have been returned to the connection
// pool and are closed before returning. If the context expires via
// cancellation, deadline, or timeout before the in use connections have been
// returned, the in use connections will be closed, resulting in the failure of
// any in flight read or write operations. If this method returns with no
// errors, all connections associated with this Server have been closed.
func (s *Server) Disconnect(ctx context.Context) error {
	if !atomic.CompareAndSwapInt64(&s.connectionstate, connected, disconnecting) {
		return ErrServerClosed
	}

	s.updateTopologyCallback.Store((updateTopologyCallback)(nil))

	// Cancel the global context so any new contexts created from it will be automatically cancelled. Close the done
	// channel so the update() routine will know that it can stop. Cancel any in-progress monitoring checks at the end.
	// The done channel is closed before cancelling the check so the update routine() will immediately detect that it
	// can stop rather than trying to create new connections until the read from done succeeds.
	s.globalCtxCancel()
	close(s.done)
	s.cancelCheck()

	s.rttMonitor.disconnect()
	err := s.pool.disconnect(ctx)
	if err != nil {
		return err
	}

	s.closewg.Wait()
	atomic.StoreInt64(&s.connectionstate, disconnected)

	return nil
}

// Connection gets a connection to the server.
func (s *Server) Connection(ctx context.Context) (driver.Connection, error) {

	if s.pool.monitor != nil {
		s.pool.monitor.Event(&event.PoolEvent{
			Type:    "ConnectionCheckOutStarted",
			Address: s.pool.address.String(),
		})
	}

	if atomic.LoadInt64(&s.connectionstate) != connected {
		return nil, ErrServerClosed
	}

	connImpl, err := s.pool.get(ctx)
	if err != nil {
		// The error has already been handled by connection.connect, which calls Server.ProcessHandshakeError.
		return nil, err
	}

	return &Connection{connection: connImpl}, nil
}

// ProcessHandshakeError implements SDAM error handling for errors that occur before a connection
// finishes handshaking. opCtx is the context passed to Server.Connection() and is used to determine
// whether or not an operation-scoped context deadline or cancellation was the cause of the
// handshake error; it is not used for timeout or cancellation of ProcessHandshakeError.
func (s *Server) ProcessHandshakeError(opCtx context.Context, err error, startingGenerationNumber uint64, serviceID *primitive.ObjectID) {
	// Ignore the error if the server is behind a load balancer but the service ID is unknown. This indicates that the
	// error happened when dialing the connection or during the MongoDB handshake, so we don't know the service ID to
	// use for clearing the pool.
	if err == nil || s.cfg.loadBalanced && serviceID == nil {
		return
	}
	// Ignore the error if the connection is stale.
	if startingGenerationNumber < s.pool.generation.getGeneration(serviceID) {
		return
	}

	wrappedConnErr := unwrapConnectionError(err)
	if wrappedConnErr == nil {
		return
	}

	isCtxTimeoutOrCanceled := func(ctx context.Context) bool {
		if ctx == nil {
			return false
		}

		if ctx.Err() != nil {
			return true
		}

		// In some networking functions, the deadline from the context is used to determine timeouts
		// instead of the ctx.Done() chan closure. In that case, there can be a race condition
		// between the networking functing returning an error and ctx.Err() returning an an error
		// (i.e. the networking function returns an error caused by the context deadline, but
		// ctx.Err() returns nil). If the operation-scoped context deadline was exceeded, assume
		// operation-scoped context timeout caused the error.
		if deadline, ok := ctx.Deadline(); ok && time.Now().After(deadline) {
			return true
		}

		return false
	}

	isErrTimeoutOrCanceled := func(err error) bool {
		for err != nil {
			// Check for errors that implement the "net.Error" interface and self-report as timeout
			// errors. Includes some "*net.OpError" errors and "context.DeadlineExceeded".
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				return true
			}
			// Check for context cancellation. Also handle the case where the cancellation error has
			// been replaced by "net.errCanceled" (which isn't exported and can't be compared
			// directly) by checking the error message.
			if err == context.Canceled || err.Error() == "operation was canceled" {
				return true
			}

			wrapper, ok := err.(interface{ Unwrap() error })
			if !ok {
				break
			}
			err = wrapper.Unwrap()
		}

		return false
	}

	// Ignore errors that indicate a client-side timeout occurred when the context passed into an
	// operation timed out or was canceled (i.e. errors caused by an operation-scoped timeout).
	// Timeouts caused by reaching connectTimeoutMS or other non-operation-scoped timeouts should
	// still clear the pool.
	// TODO(GODRIVER-2038): Remove this condition when connections are no longer created with an
	// operation-scoped timeout.
	if isCtxTimeoutOrCanceled(opCtx) && isErrTimeoutOrCanceled(wrappedConnErr) {
		return
	}

	// Since the only kind of ConnectionError we receive from pool.Get will be an initialization error, we should set
	// the description.Server appropriately. The description should not have a TopologyVersion because the staleness
	// checking logic above has already determined that this description is not stale.
	s.updateDescription(description.NewServerFromError(s.address, wrappedConnErr, nil))
	s.pool.clear(serviceID)
	s.cancelCheck()
}

// Description returns a description of the server as of the last heartbeat.
func (s *Server) Description() description.Server {
	return s.desc.Load().(description.Server)
}

// SelectedDescription returns a description.SelectedServer with a Kind of
// Single. This can be used when performing tasks like monitoring a batch
// of servers and you want to run one off commands against those servers.
func (s *Server) SelectedDescription() description.SelectedServer {
	sdesc := s.Description()
	return description.SelectedServer{
		Server: sdesc,
		Kind:   description.Single,
	}
}

// Subscribe returns a ServerSubscription which has a channel on which all
// updated server descriptions will be sent. The channel will have a buffer
// size of one, and will be pre-populated with the current description.
func (s *Server) Subscribe() (*ServerSubscription, error) {
	if atomic.LoadInt64(&s.connectionstate) != connected {
		return nil, ErrSubscribeAfterClosed
	}
	ch := make(chan description.Server, 1)
	ch <- s.desc.Load().(description.Server)

	s.subLock.Lock()
	defer s.subLock.Unlock()
	if s.subscriptionsClosed {
		return nil, ErrSubscribeAfterClosed
	}
	id := s.currentSubscriberID
	s.subscribers[id] = ch
	s.currentSubscriberID++

	ss := &ServerSubscription{
		C:  ch,
		s:  s,
		id: id,
	}

	return ss, nil
}

// RequestImmediateCheck will cause the server to send a heartbeat immediately
// instead of waiting for the heartbeat timeout.
func (s *Server) RequestImmediateCheck() {
	select {
	case s.checkNow <- struct{}{}:
	default:
	}
}

// getWriteConcernErrorForProcessing extracts a driver.WriteConcernError from the provided error. This function returns
// (error, true) if the error is a WriteConcernError and the falls under the requirements for SDAM error
// handling and (nil, false) otherwise.
func getWriteConcernErrorForProcessing(err error) (*driver.WriteConcernError, bool) {
	writeCmdErr, ok := err.(driver.WriteCommandError)
	if !ok {
		return nil, false
	}

	wcerr := writeCmdErr.WriteConcernError
	if wcerr != nil && (wcerr.NodeIsRecovering() || wcerr.NotMaster()) {
		return wcerr, true
	}
	return nil, false
}

// ProcessError handles SDAM error handling and implements driver.ErrorProcessor.
func (s *Server) ProcessError(err error, conn driver.Connection) driver.ProcessErrorResult {
	// ignore nil error
	if err == nil {
		return driver.NoChange
	}

	s.processErrorLock.Lock()
	defer s.processErrorLock.Unlock()

	// ignore stale error
	if conn.Stale() {
		return driver.NoChange
	}
	// Invalidate server description if not master or node recovering error occurs.
	// These errors can be reported as a command error or a write concern error.
	desc := conn.Description()
	if cerr, ok := err.(driver.Error); ok && (cerr.NodeIsRecovering() || cerr.NotMaster()) {
		// ignore stale error
		if desc.TopologyVersion.CompareToIncoming(cerr.TopologyVersion) >= 0 {
			return driver.NoChange
		}

		// updates description to unknown
		s.updateDescription(description.NewServerFromError(s.address, err, cerr.TopologyVersion))
		s.RequestImmediateCheck()

		res := driver.ServerMarkedUnknown
		// If the node is shutting down or is older than 4.2, we synchronously clear the pool
		if cerr.NodeIsShuttingDown() || desc.WireVersion == nil || desc.WireVersion.Max < 8 {
			res = driver.ConnectionPoolCleared
			s.pool.clear(desc.ServiceID)
		}
		return res
	}
	if wcerr, ok := getWriteConcernErrorForProcessing(err); ok {
		// ignore stale error
		if desc.TopologyVersion.CompareToIncoming(wcerr.TopologyVersion) >= 0 {
			return driver.NoChange
		}

		// updates description to unknown
		s.updateDescription(description.NewServerFromError(s.address, err, wcerr.TopologyVersion))
		s.RequestImmediateCheck()

		res := driver.ServerMarkedUnknown
		// If the node is shutting down or is older than 4.2, we synchronously clear the pool
		if wcerr.NodeIsShuttingDown() || desc.WireVersion == nil || desc.WireVersion.Max < 8 {
			res = driver.ConnectionPoolCleared
			s.pool.clear(desc.ServiceID)
		}
		return res
	}

	wrappedConnErr := unwrapConnectionError(err)
	if wrappedConnErr == nil {
		return driver.NoChange
	}

	// Ignore transient timeout errors.
	if netErr, ok := wrappedConnErr.(net.Error); ok && netErr.Timeout() {
		return driver.NoChange
	}
	if wrappedConnErr == context.Canceled || wrappedConnErr == context.DeadlineExceeded {
		return driver.NoChange
	}

	// For a non-timeout network error, we clear the pool, set the description to Unknown, and cancel the in-progress
	// monitoring check. The check is cancelled last to avoid a post-cancellation reconnect racing with
	// updateDescription.
	s.updateDescription(description.NewServerFromError(s.address, err, nil))
	s.pool.clear(desc.ServiceID)
	s.cancelCheck()
	return driver.ConnectionPoolCleared
}

// update handles performing heartbeats and updating any subscribers of the
// newest description.Server retrieved.
func (s *Server) update() {
	defer s.closewg.Done()
	heartbeatTicker := time.NewTicker(s.cfg.heartbeatInterval)
	rateLimiter := time.NewTicker(minHeartbeatInterval)
	defer heartbeatTicker.Stop()
	defer rateLimiter.Stop()
	checkNow := s.checkNow
	done := s.done

	var doneOnce bool
	defer func() {
		if r := recover(); r != nil {
			if doneOnce {
				return
			}
			// We keep this goroutine alive attempting to read from the done channel.
			<-done
		}
	}()

	closeServer := func() {
		doneOnce = true
		s.subLock.Lock()
		for id, c := range s.subscribers {
			close(c)
			delete(s.subscribers, id)
		}
		s.subscriptionsClosed = true
		s.subLock.Unlock()

		// We don't need to take s.heartbeatLock here because closeServer is called synchronously when the select checks
		// below detect that the server is being closed, so we can be sure that the connection isn't being used.
		if s.conn != nil {
			_ = s.conn.close()
		}
	}

	waitUntilNextCheck := func() {
		// Wait until heartbeatFrequency elapses, an application operation requests an immediate check, or the server
		// is disconnecting.
		select {
		case <-heartbeatTicker.C:
		case <-checkNow:
		case <-done:
			// Return because the next update iteration will check the done channel again and clean up.
			return
		}

		// Ensure we only return if minHeartbeatFrequency has elapsed or the server is disconnecting.
		select {
		case <-rateLimiter.C:
		case <-done:
			return
		}
	}

	for {
		// Check if the server is disconnecting. Even if waitForNextCheck has already read from the done channel, we
		// can safely read from it again because Disconnect closes the channel.
		select {
		case <-done:
			closeServer()
			return
		default:
		}

		previousDescription := s.Description()

		// Perform the next check.
		desc, err := s.check()
		if err == errCheckCancelled {
			if atomic.LoadInt64(&s.connectionstate) != connected {
				continue
			}

			// If the server is not disconnecting, the check was cancelled by an application operation after an error.
			// Wait before running the next check.
			waitUntilNextCheck()
			continue
		}

		s.updateDescription(desc)
		if desc.LastError != nil {
			// Clear the pool once the description has been updated to Unknown. Pass in a nil service ID to clear
			// because the monitoring routine only runs for non-load balanced deployments in which servers don't return
			// IDs.
			s.pool.clear(nil)
		}

		// If the server supports streaming or we're already streaming, we want to move to streaming the next response
		// without waiting. If the server has transitioned to Unknown from a network error, we want to do another
		// check without waiting in case it was a transient error and the server isn't actually down.
		serverSupportsStreaming := desc.Kind != description.Unknown && desc.TopologyVersion != nil
		connectionIsStreaming := s.conn != nil && s.conn.getCurrentlyStreaming()
		transitionedFromNetworkError := desc.LastError != nil && unwrapConnectionError(desc.LastError) != nil &&
			previousDescription.Kind != description.Unknown

		if serverSupportsStreaming || connectionIsStreaming || transitionedFromNetworkError {
			continue
		}

		// The server either does not support the streamable protocol or is not in a healthy state, so we wait until
		// the next check.
		waitUntilNextCheck()
	}
}

// updateDescription handles updating the description on the Server, notifying
// subscribers, and potentially draining the connection pool. The initial
// parameter is used to determine if this is the first description from the
// server.
func (s *Server) updateDescription(desc description.Server) {
	if s.cfg.loadBalanced {
		// In load balanced mode, there are no updates from the monitoring routine. For errors encountered in pooled
		// connections, the server should not be marked Unknown to ensure that the LB remains selectable.
		return
	}

	defer func() {
		//  ¯\_(ツ)_/¯
		_ = recover()
	}()

	// Use the updateTopologyCallback to update the parent Topology and get the description that should be stored.
	callback, ok := s.updateTopologyCallback.Load().(updateTopologyCallback)
	if ok && callback != nil {
		desc = callback(desc)
	}
	s.desc.Store(desc)

	s.subLock.Lock()
	for _, c := range s.subscribers {
		select {
		// drain the channel if it isn't empty
		case <-c:
		default:
		}
		c <- desc
	}
	s.subLock.Unlock()
}

// createConnection creates a new connection instance but does not call connect on it. The caller must call connect
// before the connection can be used for network operations.
func (s *Server) createConnection() (*connection, error) {
	opts := make([]ConnectionOption, len(s.cfg.connectionOpts))
	copy(opts, s.cfg.connectionOpts)
	opts = append(opts,
		WithConnectTimeout(func(time.Duration) time.Duration { return s.cfg.heartbeatTimeout }),
		WithReadTimeout(func(time.Duration) time.Duration { return s.cfg.heartbeatTimeout }),
		WithWriteTimeout(func(time.Duration) time.Duration { return s.cfg.heartbeatTimeout }),
		// We override whatever handshaker is currently attached to the options with a basic
		// one because need to make sure we don't do auth.
		WithHandshaker(func(h Handshaker) Handshaker {
			return operation.NewIsMaster().AppName(s.cfg.appname).Compressors(s.cfg.compressionOpts).
				ServerAPI(s.cfg.serverAPI)
		}),
		// Override any monitors specified in options with nil to avoid monitoring heartbeats.
		WithMonitor(func(*event.CommandMonitor) *event.CommandMonitor { return nil }),
		withPoolMonitor(func(*event.PoolMonitor) *event.PoolMonitor { return nil }),
	)

	return newConnection(s.address, opts...)
}

func (s *Server) setupHeartbeatConnection() error {
	conn, err := s.createConnection()
	if err != nil {
		return err
	}

	// Take the lock when assigning the context and connection because they're accessed by cancelCheck.
	s.heartbeatLock.Lock()
	if s.heartbeatCtxCancel != nil {
		// Ensure the previous context is cancelled to avoid a leak.
		s.heartbeatCtxCancel()
	}
	s.heartbeatCtx, s.heartbeatCtxCancel = context.WithCancel(s.globalCtx)
	s.conn = conn
	s.heartbeatLock.Unlock()

	s.conn.connect(s.heartbeatCtx)
	return s.conn.wait()
}

// cancelCheck cancels in-progress connection dials and reads. It does not set any fields on the server.
func (s *Server) cancelCheck() {
	var conn *connection

	// Take heartbeatLock for mutual exclusion with the checks in the update function.
	s.heartbeatLock.Lock()
	if s.heartbeatCtx != nil {
		s.heartbeatCtxCancel()
	}
	conn = s.conn
	s.heartbeatLock.Unlock()

	if conn == nil {
		return
	}

	// If the connection exists, we need to wait for it to be connected conn.connect() and conn.close() cannot be called
	// concurrently. We can ignore the error from conn.wait(). If the connection wasn't successfully opened, its state
	// was set back to disconnected, so calling conn.close() will be a noop.
	conn.closeConnectContext()
	_ = conn.wait()
	_ = conn.close()
}

func (s *Server) checkWasCancelled() bool {
	return s.heartbeatCtx.Err() != nil
}

func (s *Server) createBaseOperation(conn driver.Connection) *operation.IsMaster {
	return operation.
		NewIsMaster().
		ClusterClock(s.cfg.clock).
		Deployment(driver.SingleConnectionDeployment{conn}).
		ServerAPI(s.cfg.serverAPI)
}

func (s *Server) check() (description.Server, error) {
	var descPtr *description.Server
	var err error
	var durationNanos int64

	// Create a new connection if this is the first check, the connection was closed after an error during the previous
	// check, or the previous check was cancelled.
	if s.conn == nil || s.conn.closed() || s.checkWasCancelled() {
		// Create a new connection and add it's handshake RTT as a sample.
		err = s.setupHeartbeatConnection()
		if err == nil {
			// Use the description from the connection handshake as the value for this check.
			s.rttMonitor.addSample(s.conn.isMasterRTT)
			descPtr = &s.conn.desc
			durationNanos = s.conn.isMasterRTT.Nanoseconds()
		}
	}

	if descPtr == nil && err == nil {
		// An existing connection is being used. Use the server description properties to execute the right heartbeat.

		// Wrap conn in a type that implements driver.StreamerConnection.
		heartbeatConn := initConnection{s.conn}
		baseOperation := s.createBaseOperation(heartbeatConn)
		previousDescription := s.Description()
		streamable := previousDescription.TopologyVersion != nil

		s.publishServerHeartbeatStartedEvent(s.conn.ID(), s.conn.getCurrentlyStreaming() || streamable)
		start := time.Now()
		switch {
		case s.conn.getCurrentlyStreaming():
			// The connection is already in a streaming state, so we stream the next response.
			err = baseOperation.StreamResponse(s.heartbeatCtx, heartbeatConn)
		case streamable:
			// The server supports the streamable protocol. Set the socket timeout to
			// connectTimeoutMS+heartbeatFrequencyMS and execute an awaitable isMaster request. Set conn.canStream so
			// the wire message will advertise streaming support to the server.

			// Calculation for maxAwaitTimeMS is taken from time.Duration.Milliseconds (added in Go 1.13).
			maxAwaitTimeMS := int64(s.cfg.heartbeatInterval) / 1e6
			// If connectTimeoutMS=0, the socket timeout should be infinite. Otherwise, it is connectTimeoutMS +
			// heartbeatFrequencyMS to account for the fact that the query will block for heartbeatFrequencyMS
			// server-side.
			socketTimeout := s.cfg.heartbeatTimeout
			if socketTimeout != 0 {
				socketTimeout += s.cfg.heartbeatInterval
			}
			s.conn.setSocketTimeout(socketTimeout)
			baseOperation = baseOperation.TopologyVersion(previousDescription.TopologyVersion).
				MaxAwaitTimeMS(maxAwaitTimeMS)
			s.conn.setCanStream(true)
			err = baseOperation.Execute(s.heartbeatCtx)
		default:
			// The server doesn't support the awaitable protocol. Set the socket timeout to connectTimeoutMS and
			// execute a regular heartbeat without any additional parameters.

			s.conn.setSocketTimeout(s.cfg.heartbeatTimeout)
			err = baseOperation.Execute(s.heartbeatCtx)
		}
		durationNanos = time.Since(start).Nanoseconds()

		if err == nil {
			tempDesc := baseOperation.Result(s.address)
			descPtr = &tempDesc
			s.publishServerHeartbeatSucceededEvent(s.conn.ID(), durationNanos, tempDesc, s.conn.getCurrentlyStreaming() || streamable)
		} else {
			// Close the connection here rather than below so we ensure we're not closing a connection that wasn't
			// successfully created.
			if s.conn != nil {
				_ = s.conn.close()
			}
			s.publishServerHeartbeatFailedEvent(s.conn.ID(), durationNanos, err, s.conn.getCurrentlyStreaming() || streamable)
		}
	}

	if descPtr != nil {
		// The check was successful. Set the average RTT and return.
		desc := *descPtr
		desc = desc.SetAverageRTT(s.rttMonitor.getRTT())
		desc.HeartbeatInterval = s.cfg.heartbeatInterval
		return desc, nil
	}

	if s.checkWasCancelled() {
		// If the previous check was cancelled, we don't want to clear the pool. Return a sentinel error so the caller
		// will know that an actual error didn't occur.
		return emptyDescription, errCheckCancelled
	}

	// An error occurred. We reset the RTT monitor for all errors and return an Unknown description. The pool must also
	// be cleared, but only after the description has already been updated, so that is handled by the caller.
	topologyVersion := extractTopologyVersion(err)
	s.rttMonitor.reset()
	return description.NewServerFromError(s.address, err, topologyVersion), nil
}

func extractTopologyVersion(err error) *description.TopologyVersion {
	if ce, ok := err.(ConnectionError); ok {
		err = ce.Wrapped
	}

	switch converted := err.(type) {
	case driver.Error:
		return converted.TopologyVersion
	case driver.WriteCommandError:
		if converted.WriteConcernError != nil {
			return converted.WriteConcernError.TopologyVersion
		}
	}

	return nil
}

// String implements the Stringer interface.
func (s *Server) String() string {
	desc := s.Description()
	connState := atomic.LoadInt64(&s.connectionstate)
	str := fmt.Sprintf("Addr: %s, Type: %s, State: %s",
		s.address, desc.Kind, connectionStateString(connState))
	if len(desc.Tags) != 0 {
		str += fmt.Sprintf(", Tag sets: %s", desc.Tags)
	}
	if connState == connected {
		str += fmt.Sprintf(", Average RTT: %d", desc.AverageRTT)
	}
	if desc.LastError != nil {
		str += fmt.Sprintf(", Last error: %s", desc.LastError)
	}

	return str
}

// ServerSubscription represents a subscription to the description.Server updates for
// a specific server.
type ServerSubscription struct {
	C  <-chan description.Server
	s  *Server
	id uint64
}

// Unsubscribe unsubscribes this ServerSubscription from updates and closes the
// subscription channel.
func (ss *ServerSubscription) Unsubscribe() error {
	ss.s.subLock.Lock()
	defer ss.s.subLock.Unlock()
	if ss.s.subscriptionsClosed {
		return nil
	}

	ch, ok := ss.s.subscribers[ss.id]
	if !ok {
		return nil
	}

	close(ch)
	delete(ss.s.subscribers, ss.id)

	return nil
}

// publishes a ServerOpeningEvent to indicate the server is being initialized
func (s *Server) publishServerOpeningEvent(addr address.Address) {
	serverOpening := &event.ServerOpeningEvent{
		Address:    addr,
		TopologyID: s.topologyID,
	}

	if s != nil && s.cfg.serverMonitor != nil && s.cfg.serverMonitor.ServerOpening != nil {
		s.cfg.serverMonitor.ServerOpening(serverOpening)
	}
}

// publishes a ServerHeartbeatStartedEvent to indicate an ismaster command has started
func (s *Server) publishServerHeartbeatStartedEvent(connectionID string, await bool) {
	serverHeartbeatStarted := &event.ServerHeartbeatStartedEvent{
		ConnectionID: connectionID,
		Awaited:      await,
	}

	if s != nil && s.cfg.serverMonitor != nil && s.cfg.serverMonitor.ServerHeartbeatStarted != nil {
		s.cfg.serverMonitor.ServerHeartbeatStarted(serverHeartbeatStarted)
	}
}

// publishes a ServerHeartbeatSucceededEvent to indicate ismaster has succeeded
func (s *Server) publishServerHeartbeatSucceededEvent(connectionID string,
	durationNanos int64,
	desc description.Server,
	await bool) {
	serverHeartbeatSucceeded := &event.ServerHeartbeatSucceededEvent{
		DurationNanos: durationNanos,
		Reply:         desc,
		ConnectionID:  connectionID,
		Awaited:       await,
	}

	if s != nil && s.cfg.serverMonitor != nil && s.cfg.serverMonitor.ServerHeartbeatSucceeded != nil {
		s.cfg.serverMonitor.ServerHeartbeatSucceeded(serverHeartbeatSucceeded)
	}
}

// publishes a ServerHeartbeatFailedEvent to indicate ismaster has failed
func (s *Server) publishServerHeartbeatFailedEvent(connectionID string,
	durationNanos int64,
	err error,
	await bool) {
	serverHeartbeatFailed := &event.ServerHeartbeatFailedEvent{
		DurationNanos: durationNanos,
		Failure:       err,
		ConnectionID:  connectionID,
		Awaited:       await,
	}

	if s != nil && s.cfg.serverMonitor != nil && s.cfg.serverMonitor.ServerHeartbeatFailed != nil {
		s.cfg.serverMonitor.ServerHeartbeatFailed(serverHeartbeatFailed)
	}
}

// unwrapConnectionError returns the connection error wrapped by err, or nil if err does not wrap a connection error.
func unwrapConnectionError(err error) error {
	// This is essentially an implementation of errors.As to unwrap this error until we get a ConnectionError and then
	// return ConnectionError.Wrapped.

	connErr, ok := err.(ConnectionError)
	if ok {
		return connErr.Wrapped
	}

	driverErr, ok := err.(driver.Error)
	if !ok || !driverErr.NetworkError() {
		return nil
	}

	connErr, ok = driverErr.Wrapped.(ConnectionError)
	if ok {
		return connErr.Wrapped
	}

	return nil
}
