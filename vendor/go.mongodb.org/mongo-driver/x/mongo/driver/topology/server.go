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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/driverutil"
	"go.mongodb.org/mongo-driver/internal/logger"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
)

const minHeartbeatInterval = 500 * time.Millisecond
const wireVersion42 = 8 // Wire version for MongoDB 4.2

// Server state constants.
const (
	serverDisconnected int64 = iota
	serverDisconnecting
	serverConnected
)

func serverStateString(state int64) string {
	switch state {
	case serverDisconnected:
		return "Disconnected"
	case serverDisconnecting:
		return "Disconnecting"
	case serverConnected:
		return "Connected"
	}

	return ""
}

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

// Server is a single server within a topology.
type Server struct {
	// The following integer fields must be accessed using the atomic package and should be at the
	// beginning of the struct.
	// - atomic bug: https://pkg.go.dev/sync/atomic#pkg-note-BUG
	// - suggested layout: https://go101.org/article/memory-layout.html

	state          int64
	operationCount int64

	cfg     *serverConfig
	address address.Address

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
	monitorOnce      sync.Once
}

// updateTopologyCallback is a callback used to create a server that should be called when the parent Topology instance
// should be updated based on a new server description. The callback must return the server description that should be
// stored by the server.
type updateTopologyCallback func(description.Server) description.Server

// ConnectServer creates a new Server and then initializes it using the
// Connect method.
func ConnectServer(
	addr address.Address,
	updateCallback updateTopologyCallback,
	topologyID primitive.ObjectID,
	opts ...ServerOption,
) (*Server, error) {
	srvr := NewServer(addr, topologyID, opts...)
	err := srvr.Connect(updateCallback)
	if err != nil {
		return nil, err
	}
	return srvr, nil
}

// NewServer creates a new server. The mongodb server at the address will be monitored
// on an internal monitoring goroutine.
func NewServer(addr address.Address, topologyID primitive.ObjectID, opts ...ServerOption) *Server {
	cfg := newServerConfig(opts...)
	globalCtx, globalCtxCancel := context.WithCancel(context.Background())
	s := &Server{
		state: serverDisconnected,

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
		minRTTWindow:       5 * time.Minute,
		createConnectionFn: s.createConnection,
		createOperationFn:  s.createBaseOperation,
	}
	s.rttMonitor = newRTTMonitor(rttCfg)

	pc := poolConfig{
		Address:          addr,
		MinPoolSize:      cfg.minConns,
		MaxPoolSize:      cfg.maxConns,
		MaxConnecting:    cfg.maxConnecting,
		MaxIdleTime:      cfg.poolMaxIdleTime,
		MaintainInterval: cfg.poolMaintainInterval,
		LoadBalanced:     cfg.loadBalanced,
		PoolMonitor:      cfg.poolMonitor,
		Logger:           cfg.logger,
		handshakeErrFn:   s.ProcessHandshakeError,
	}

	connectionOpts := copyConnectionOpts(cfg.connectionOpts)
	s.pool = newPool(pc, connectionOpts...)
	s.publishServerOpeningEvent(s.address)

	return s
}

func mustLogServerMessage(srv *Server) bool {
	return srv.cfg.logger != nil && srv.cfg.logger.LevelComponentEnabled(
		logger.LevelDebug, logger.ComponentTopology)
}

func logServerMessage(srv *Server, msg string, keysAndValues ...interface{}) {
	serverHost, serverPort, err := net.SplitHostPort(srv.address.String())
	if err != nil {
		serverHost = srv.address.String()
		serverPort = ""
	}

	var driverConnectionID uint64
	var serverConnectionID *int64

	if srv.conn != nil {
		driverConnectionID = srv.conn.driverConnectionID
		serverConnectionID = srv.conn.serverConnectionID
	}

	srv.cfg.logger.Print(logger.LevelDebug,
		logger.ComponentTopology,
		msg,
		logger.SerializeServer(logger.Server{
			DriverConnectionID: driverConnectionID,
			TopologyID:         srv.topologyID,
			Message:            msg,
			ServerConnectionID: serverConnectionID,
			ServerHost:         serverHost,
			ServerPort:         serverPort,
		}, keysAndValues...)...)
}

// Connect initializes the Server by starting background monitoring goroutines.
// This method must be called before a Server can be used.
func (s *Server) Connect(updateCallback updateTopologyCallback) error {
	if !atomic.CompareAndSwapInt64(&s.state, serverDisconnected, serverConnected) {
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
		s.closewg.Add(1)
		go s.update()
	}

	// The CMAP spec describes that pools should only be marked "ready" when the server description
	// is updated to something other than "Unknown". However, we maintain the previous Server
	// behavior here and immediately mark the pool as ready during Connect() to simplify and speed
	// up the Client startup behavior. The risk of marking a pool as ready proactively during
	// Connect() is that we could attempt to create connections to a server that was configured
	// erroneously until the first server check or checkOut() failure occurs, when the SDAM error
	// handler would transition the Server back to "Unknown" and set the pool to "paused".
	return s.pool.ready()
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
	if !atomic.CompareAndSwapInt64(&s.state, serverConnected, serverDisconnecting) {
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

	s.pool.close(ctx)

	s.closewg.Wait()
	s.rttMonitor.disconnect()
	atomic.StoreInt64(&s.state, serverDisconnected)

	return nil
}

// Connection gets a connection to the server.
func (s *Server) Connection(ctx context.Context) (driver.Connection, error) {
	if atomic.LoadInt64(&s.state) != serverConnected {
		return nil, ErrServerClosed
	}

	// Increment the operation count before calling checkOut to make sure that all connection
	// requests are included in the operation count, including those in the wait queue. If we got an
	// error instead of a connection, immediately decrement the operation count.
	atomic.AddInt64(&s.operationCount, 1)
	conn, err := s.pool.checkOut(ctx)
	if err != nil {
		atomic.AddInt64(&s.operationCount, -1)
		return nil, err
	}

	return &Connection{
		connection: conn,
		cleanupServerFn: func() {
			// Decrement the operation count whenever the caller is done with the connection. Note
			// that cleanupServerFn() is not called while the connection is pinned to a cursor or
			// transaction, so the operation count is not decremented until the cursor is closed or
			// the transaction is committed or aborted. Use an int64 instead of a uint64 to mitigate
			// the impact of any possible bugs that could cause the uint64 to underflow, which would
			// make the server much less selectable.
			atomic.AddInt64(&s.operationCount, -1)
		},
	}, nil
}

// ProcessHandshakeError implements SDAM error handling for errors that occur before a connection
// finishes handshaking.
func (s *Server) ProcessHandshakeError(err error, startingGenerationNumber uint64, serviceID *primitive.ObjectID) {
	// Ignore the error if the server is behind a load balancer but the service ID is unknown. This indicates that the
	// error happened when dialing the connection or during the MongoDB handshake, so we don't know the service ID to
	// use for clearing the pool.
	if err == nil || s.cfg.loadBalanced && serviceID == nil {
		return
	}
	// Ignore the error if the connection is stale.
	if generation, _ := s.pool.generation.getGeneration(serviceID); startingGenerationNumber < generation {
		return
	}

	// Unwrap any connection errors. If there is no wrapped connection error, then the error should
	// not result in any Server state change (e.g. a command error from the database).
	wrappedConnErr := unwrapConnectionError(err)
	if wrappedConnErr == nil {
		return
	}

	// Must hold the processErrorLock while updating the server description and clearing the pool.
	// Not holding the lock leads to possible out-of-order processing of pool.clear() and
	// pool.ready() calls from concurrent server description updates.
	s.processErrorLock.Lock()
	defer s.processErrorLock.Unlock()

	// Since the only kind of ConnectionError we receive from pool.Get will be an initialization error, we should set
	// the description.Server appropriately. The description should not have a TopologyVersion because the staleness
	// checking logic above has already determined that this description is not stale.
	s.updateDescription(description.NewServerFromError(s.address, wrappedConnErr, nil))
	s.pool.clear(err, serviceID)
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
	if atomic.LoadInt64(&s.state) != serverConnected {
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
	var writeCmdErr driver.WriteCommandError
	if !errors.As(err, &writeCmdErr) {
		return nil, false
	}

	wcerr := writeCmdErr.WriteConcernError
	if wcerr != nil && (wcerr.NodeIsRecovering() || wcerr.NotPrimary()) {
		return wcerr, true
	}
	return nil, false
}

// ProcessError handles SDAM error handling and implements driver.ErrorProcessor.
func (s *Server) ProcessError(err error, conn driver.Connection) driver.ProcessErrorResult {
	// Ignore nil errors.
	if err == nil {
		return driver.NoChange
	}

	// Ignore errors from stale connections because the error came from a previous generation of the
	// connection pool. The root cause of the error has already been handled, which is what caused
	// the pool generation to increment. Processing errors for stale connections could result in
	// handling the same error root cause multiple times (e.g. a temporary network interrupt causing
	// all connections to the same server to return errors).
	if conn.Stale() {
		return driver.NoChange
	}

	// Must hold the processErrorLock while updating the server description and clearing the pool.
	// Not holding the lock leads to possible out-of-order processing of pool.clear() and
	// pool.ready() calls from concurrent server description updates.
	s.processErrorLock.Lock()
	defer s.processErrorLock.Unlock()

	// Get the wire version and service ID from the connection description because they will never
	// change for the lifetime of a connection and can possibly be different between connections to
	// the same server.
	connDesc := conn.Description()
	wireVersion := connDesc.WireVersion
	serviceID := connDesc.ServiceID

	// Get the topology version from the Server description because the Server description is
	// updated by heartbeats and errors, so typically has a more up-to-date topology version.
	serverDesc := s.desc.Load().(description.Server)
	topologyVersion := serverDesc.TopologyVersion

	// We don't currently update the Server topology version when we create new application
	// connections, so it's possible for a connection's topology version to be newer than the
	// Server's topology version. Pick the "newest" of the two topology versions.
	// Technically a nil topology version on a new database response should be considered a new
	// topology version and replace the Server's topology version. However, we don't know if the
	// connection's topology version is based on a new or old database response, so we ignore a nil
	// topology version on the connection for now.
	//
	// TODO(GODRIVER-2841): Remove this logic once we set the Server description when we create
	// TODO application connections because then the Server's topology version will always be the
	// TODO latest known.
	if tv := connDesc.TopologyVersion; tv != nil && topologyVersion.CompareToIncoming(tv) < 0 {
		topologyVersion = tv
	}

	// Invalidate server description if not primary or node recovering error occurs.
	// These errors can be reported as a command error or a write concern error.
	if cerr, ok := err.(driver.Error); ok && (cerr.NodeIsRecovering() || cerr.NotPrimary()) {
		// Ignore errors that came from when the database was on a previous topology version.
		if topologyVersion.CompareToIncoming(cerr.TopologyVersion) >= 0 {
			return driver.NoChange
		}

		// updates description to unknown
		s.updateDescription(description.NewServerFromError(s.address, err, cerr.TopologyVersion))
		s.RequestImmediateCheck()

		res := driver.ServerMarkedUnknown
		// If the node is shutting down or is older than 4.2, we synchronously clear the pool
		if cerr.NodeIsShuttingDown() || wireVersion == nil || wireVersion.Max < wireVersion42 {
			res = driver.ConnectionPoolCleared
			s.pool.clear(err, serviceID)
		}

		return res
	}
	if wcerr, ok := getWriteConcernErrorForProcessing(err); ok {
		// Ignore errors that came from when the database was on a previous topology version.
		if topologyVersion.CompareToIncoming(wcerr.TopologyVersion) >= 0 {
			return driver.NoChange
		}

		// updates description to unknown
		s.updateDescription(description.NewServerFromError(s.address, err, wcerr.TopologyVersion))
		s.RequestImmediateCheck()

		res := driver.ServerMarkedUnknown
		// If the node is shutting down or is older than 4.2, we synchronously clear the pool
		if wcerr.NodeIsShuttingDown() || wireVersion == nil || wireVersion.Max < wireVersion42 {
			res = driver.ConnectionPoolCleared
			s.pool.clear(err, serviceID)
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
	if errors.Is(wrappedConnErr, context.Canceled) || errors.Is(wrappedConnErr, context.DeadlineExceeded) {
		return driver.NoChange
	}

	// For a non-timeout network error, we clear the pool, set the description to Unknown, and cancel the in-progress
	// monitoring check. The check is cancelled last to avoid a post-cancellation reconnect racing with
	// updateDescription.
	s.updateDescription(description.NewServerFromError(s.address, err, nil))
	s.pool.clear(err, serviceID)
	s.cancelCheck()
	return driver.ConnectionPoolCleared
}

// update handle performing heartbeats and updating any subscribers of the
// newest description.Server retrieved.
func (s *Server) update() {
	defer s.closewg.Done()
	heartbeatTicker := time.NewTicker(s.cfg.heartbeatInterval)
	rateLimiter := time.NewTicker(minHeartbeatInterval)
	defer heartbeatTicker.Stop()
	defer rateLimiter.Stop()
	checkNow := s.checkNow
	done := s.done

	defer logUnexpectedFailure(s.cfg.logger, "Encountered unexpected failure updating server")

	closeServer := func() {
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

	timeoutCnt := 0
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
		if errors.Is(err, errCheckCancelled) {
			if atomic.LoadInt64(&s.state) != serverConnected {
				continue
			}

			// If the server is not disconnecting, the check was cancelled by an application operation after an error.
			// Wait before running the next check.
			waitUntilNextCheck()
			continue
		}

		if isShortcut := func() bool {
			// Must hold the processErrorLock while updating the server description and clearing the
			// pool. Not holding the lock leads to possible out-of-order processing of pool.clear() and
			// pool.ready() calls from concurrent server description updates.
			s.processErrorLock.Lock()
			defer s.processErrorLock.Unlock()

			s.updateDescription(desc)
			// Retry after the first timeout before clearing the pool in case of a FAAS pause as
			// described in GODRIVER-2577.
			if err := unwrapConnectionError(desc.LastError); err != nil && timeoutCnt < 1 {
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					timeoutCnt++
					// We want to immediately retry on timeout error. Continue to next loop.
					return true
				}
				if err, ok := err.(net.Error); ok && err.Timeout() {
					timeoutCnt++
					// We want to immediately retry on timeout error. Continue to next loop.
					return true
				}
			}
			if err := desc.LastError; err != nil {
				// Clear the pool once the description has been updated to Unknown. Pass in a nil service ID to clear
				// because the monitoring routine only runs for non-load balanced deployments in which servers don't return
				// IDs.
				if timeoutCnt > 0 {
					s.pool.clearAll(err, nil)
				} else {
					s.pool.clear(err, nil)
				}
			}
			// We're either not handling a timeout error, or we just handled the 2nd consecutive
			// timeout error. In either case, reset the timeout count to 0 and return false to
			// continue the normal check process.
			timeoutCnt = 0
			return false
		}(); isShortcut {
			continue
		}

		// If the server supports streaming or we're already streaming, we want to move to streaming the next response
		// without waiting. If the server has transitioned to Unknown from a network error, we want to do another
		// check without waiting in case it was a transient error and the server isn't actually down.
		connectionIsStreaming := s.conn != nil && s.conn.getCurrentlyStreaming()
		transitionedFromNetworkError := desc.LastError != nil && unwrapConnectionError(desc.LastError) != nil &&
			previousDescription.Kind != description.Unknown

		if isStreamingEnabled(s) && isStreamable(s) {
			s.monitorOnce.Do(s.rttMonitor.connect)
		}

		if isStreamingEnabled(s) && (isStreamable(s) || connectionIsStreaming) || transitionedFromNetworkError {
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

	defer logUnexpectedFailure(s.cfg.logger, "Encountered unexpected failure updating server description")

	// Anytime we update the server description to something other than "unknown", set the pool to
	// "ready". Do this before updating the description so that connections can be checked out as
	// soon as the server is selectable. If the pool is already ready, this operation is a no-op.
	// Note that this behavior is roughly consistent with the current Go driver behavior (connects
	// to all servers, even non-data-bearing nodes) but deviates slightly from CMAP spec, which
	// specifies a more restricted set of server descriptions and topologies that should mark the
	// pool ready. We don't have access to the topology here, so prefer the current Go driver
	// behavior for simplicity.
	if desc.Kind != description.Unknown {
		_ = s.pool.ready()
	}

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
func (s *Server) createConnection() *connection {
	opts := copyConnectionOpts(s.cfg.connectionOpts)
	opts = append(opts,
		WithConnectTimeout(func(time.Duration) time.Duration { return s.cfg.heartbeatTimeout }),
		WithReadTimeout(func(time.Duration) time.Duration { return s.cfg.heartbeatTimeout }),
		WithWriteTimeout(func(time.Duration) time.Duration { return s.cfg.heartbeatTimeout }),
		// We override whatever handshaker is currently attached to the options with a basic
		// one because need to make sure we don't do auth.
		WithHandshaker(func(Handshaker) Handshaker {
			return operation.NewHello().AppName(s.cfg.appname).Compressors(s.cfg.compressionOpts).
				ServerAPI(s.cfg.serverAPI)
		}),
		// Override any monitors specified in options with nil to avoid monitoring heartbeats.
		WithMonitor(func(*event.CommandMonitor) *event.CommandMonitor { return nil }),
	)

	return newConnection(s.address, opts...)
}

func copyConnectionOpts(opts []ConnectionOption) []ConnectionOption {
	optsCopy := make([]ConnectionOption, len(opts))
	copy(optsCopy, opts)
	return optsCopy
}

func (s *Server) setupHeartbeatConnection() error {
	conn := s.createConnection()

	// Take the lock when assigning the context and connection because they're accessed by cancelCheck.
	s.heartbeatLock.Lock()
	if s.heartbeatCtxCancel != nil {
		// Ensure the previous context is cancelled to avoid a leak.
		s.heartbeatCtxCancel()
	}
	s.heartbeatCtx, s.heartbeatCtxCancel = context.WithCancel(s.globalCtx)
	s.conn = conn
	s.heartbeatLock.Unlock()

	return s.conn.connect(s.heartbeatCtx)
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

	// If the connection exists, we need to wait for it to be connected because conn.connect() and
	// conn.close() cannot be called concurrently. If the connection wasn't successfully opened, its
	// state was set back to disconnected, so calling conn.close() will be a no-op.
	conn.closeConnectContext()
	conn.wait()
	_ = conn.close()
}

func (s *Server) checkWasCancelled() bool {
	return s.heartbeatCtx.Err() != nil
}

func (s *Server) createBaseOperation(conn driver.Connection) *operation.Hello {
	return operation.
		NewHello().
		ClusterClock(s.cfg.clock).
		Deployment(driver.SingleConnectionDeployment{C: conn}).
		ServerAPI(s.cfg.serverAPI)
}

func isStreamingEnabled(srv *Server) bool {
	switch srv.cfg.serverMonitoringMode {
	case connstring.ServerMonitoringModeStream:
		return true
	case connstring.ServerMonitoringModePoll:
		return false
	default:
		return driverutil.GetFaasEnvName() == ""
	}
}

func isStreamable(srv *Server) bool {
	return srv.Description().Kind != description.Unknown && srv.Description().TopologyVersion != nil
}

func (s *Server) check() (description.Server, error) {
	var descPtr *description.Server
	var err error
	var duration time.Duration

	start := time.Now()

	// Create a new connection if this is the first check, the connection was closed after an error during the previous
	// check, or the previous check was cancelled.
	if s.conn == nil || s.conn.closed() || s.checkWasCancelled() {
		connID := "0"
		if s.conn != nil {
			connID = s.conn.ID()
		}
		s.publishServerHeartbeatStartedEvent(connID, false)
		// Create a new connection and add it's handshake RTT as a sample.
		err = s.setupHeartbeatConnection()
		duration = time.Since(start)
		connID = "0"
		if s.conn != nil {
			connID = s.conn.ID()
		}
		if err == nil {
			// Use the description from the connection handshake as the value for this check.
			s.rttMonitor.addSample(s.conn.helloRTT)
			descPtr = &s.conn.desc
			s.publishServerHeartbeatSucceededEvent(connID, duration, s.conn.desc, false)
		} else {
			err = unwrapConnectionError(err)
			s.publishServerHeartbeatFailedEvent(connID, duration, err, false)
		}
	} else {
		// An existing connection is being used. Use the server description properties to execute the right heartbeat.

		// Wrap conn in a type that implements driver.StreamerConnection.
		heartbeatConn := initConnection{s.conn}
		baseOperation := s.createBaseOperation(heartbeatConn)
		previousDescription := s.Description()
		streamable := isStreamingEnabled(s) && isStreamable(s)

		s.publishServerHeartbeatStartedEvent(s.conn.ID(), s.conn.getCurrentlyStreaming() || streamable)

		switch {
		case s.conn.getCurrentlyStreaming():
			// The connection is already in a streaming state, so we stream the next response.
			err = baseOperation.StreamResponse(s.heartbeatCtx, heartbeatConn)
		case streamable:
			// The server supports the streamable protocol. Set the socket timeout to
			// connectTimeoutMS+heartbeatFrequencyMS and execute an awaitable hello request. Set conn.canStream so
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

		duration = time.Since(start)

		// We need to record an RTT sample in the polling case so that if the server
		// is < 4.4, or if polling is specified by the user, then the
		// RTT-short-circuit feature of CSOT is not disabled.
		if !streamable {
			s.rttMonitor.addSample(duration)
		}

		if err == nil {
			tempDesc := baseOperation.Result(s.address)
			descPtr = &tempDesc
			s.publishServerHeartbeatSucceededEvent(s.conn.ID(), duration, tempDesc, s.conn.getCurrentlyStreaming() || streamable)
		} else {
			// Close the connection here rather than below so we ensure we're not closing a connection that wasn't
			// successfully created.
			if s.conn != nil {
				_ = s.conn.close()
			}
			s.publishServerHeartbeatFailedEvent(s.conn.ID(), duration, err, s.conn.getCurrentlyStreaming() || streamable)
		}
	}

	if descPtr != nil {
		// The check was successful. Set the average RTT and the 90th percentile RTT and return.
		desc := *descPtr
		desc = desc.SetAverageRTT(s.rttMonitor.EWMA())
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

// RTTMonitor returns this server's round-trip-time monitor.
func (s *Server) RTTMonitor() driver.RTTMonitor {
	return s.rttMonitor
}

// OperationCount returns the current number of in-progress operations for this server.
func (s *Server) OperationCount() int64 {
	return atomic.LoadInt64(&s.operationCount)
}

// String implements the Stringer interface.
func (s *Server) String() string {
	desc := s.Description()
	state := atomic.LoadInt64(&s.state)
	str := fmt.Sprintf("Addr: %s, Type: %s, State: %s",
		s.address, desc.Kind, serverStateString(state))
	if len(desc.Tags) != 0 {
		str += fmt.Sprintf(", Tag sets: %s", desc.Tags)
	}
	if state == serverConnected {
		str += fmt.Sprintf(", Average RTT: %s, Min RTT: %s", desc.AverageRTT, s.RTTMonitor().Min())
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
	if s == nil {
		return
	}

	serverOpening := &event.ServerOpeningEvent{
		Address:    addr,
		TopologyID: s.topologyID,
	}

	if s.cfg.serverMonitor != nil && s.cfg.serverMonitor.ServerOpening != nil {
		s.cfg.serverMonitor.ServerOpening(serverOpening)
	}

	if mustLogServerMessage(s) {
		logServerMessage(s, logger.TopologyServerOpening)
	}
}

// publishes a ServerHeartbeatStartedEvent to indicate a hello command has started
func (s *Server) publishServerHeartbeatStartedEvent(connectionID string, await bool) {
	serverHeartbeatStarted := &event.ServerHeartbeatStartedEvent{
		ConnectionID: connectionID,
		Awaited:      await,
	}

	if s != nil && s.cfg.serverMonitor != nil && s.cfg.serverMonitor.ServerHeartbeatStarted != nil {
		s.cfg.serverMonitor.ServerHeartbeatStarted(serverHeartbeatStarted)
	}

	if mustLogServerMessage(s) {
		logServerMessage(s, logger.TopologyServerHeartbeatStarted,
			logger.KeyAwaited, await)
	}
}

// publishes a ServerHeartbeatSucceededEvent to indicate hello has succeeded
func (s *Server) publishServerHeartbeatSucceededEvent(connectionID string,
	duration time.Duration,
	desc description.Server,
	await bool,
) {
	serverHeartbeatSucceeded := &event.ServerHeartbeatSucceededEvent{
		DurationNanos: duration.Nanoseconds(),
		Duration:      duration,
		Reply:         desc,
		ConnectionID:  connectionID,
		Awaited:       await,
	}

	if s != nil && s.cfg.serverMonitor != nil && s.cfg.serverMonitor.ServerHeartbeatSucceeded != nil {
		s.cfg.serverMonitor.ServerHeartbeatSucceeded(serverHeartbeatSucceeded)
	}

	if mustLogServerMessage(s) {
		descRaw, _ := bson.Marshal(struct {
			description.Server `bson:",inline"`
			Ok                 int32
		}{
			Server: desc,
			Ok: func() int32 {
				if desc.LastError != nil {
					return 0
				}

				return 1
			}(),
		})

		logServerMessage(s, logger.TopologyServerHeartbeatSucceeded,
			logger.KeyAwaited, await,
			logger.KeyDurationMS, duration.Milliseconds(),
			logger.KeyReply, bson.Raw(descRaw).String())
	}
}

// publishes a ServerHeartbeatFailedEvent to indicate hello has failed
func (s *Server) publishServerHeartbeatFailedEvent(connectionID string,
	duration time.Duration,
	err error,
	await bool,
) {
	serverHeartbeatFailed := &event.ServerHeartbeatFailedEvent{
		DurationNanos: duration.Nanoseconds(),
		Duration:      duration,
		Failure:       err,
		ConnectionID:  connectionID,
		Awaited:       await,
	}

	if s != nil && s.cfg.serverMonitor != nil && s.cfg.serverMonitor.ServerHeartbeatFailed != nil {
		s.cfg.serverMonitor.ServerHeartbeatFailed(serverHeartbeatFailed)
	}

	if mustLogServerMessage(s) {
		logServerMessage(s, logger.TopologyServerHeartbeatFailed,
			logger.KeyAwaited, await,
			logger.KeyDurationMS, duration.Milliseconds(),
			logger.KeyFailure, err.Error())
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
