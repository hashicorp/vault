// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/logger"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
)

// Connection pool state constants.
const (
	poolPaused int = iota
	poolReady
	poolClosed
)

// ErrPoolNotPaused is returned when attempting to mark a connection pool "ready" that is not
// currently "paused".
var ErrPoolNotPaused = PoolError("only a paused pool can be marked ready")

// ErrPoolClosed is returned when attempting to check out a connection from a closed pool.
var ErrPoolClosed = PoolError("attempted to check out a connection from closed connection pool")

// ErrConnectionClosed is returned from an attempt to use an already closed connection.
var ErrConnectionClosed = ConnectionError{ConnectionID: "<closed>", message: "connection is closed"}

// ErrWrongPool is return when a connection is returned to a pool it doesn't belong to.
var ErrWrongPool = PoolError("connection does not belong to this pool")

// PoolError is an error returned from a Pool method.
type PoolError string

func (pe PoolError) Error() string { return string(pe) }

// poolClearedError is an error returned when the connection pool is cleared or currently paused. It
// is a retryable error.
type poolClearedError struct {
	err     error
	address address.Address
}

func (pce poolClearedError) Error() string {
	return fmt.Sprintf(
		"connection pool for %v was cleared because another operation failed with: %v",
		pce.address,
		pce.err)
}

// Retryable returns true. All poolClearedErrors are retryable.
func (poolClearedError) Retryable() bool { return true }

// Assert that poolClearedError is a driver.RetryablePoolError.
var _ driver.RetryablePoolError = poolClearedError{}

// poolConfig contains all aspects of the pool that can be configured
type poolConfig struct {
	Address          address.Address
	MinPoolSize      uint64
	MaxPoolSize      uint64
	MaxConnecting    uint64
	MaxIdleTime      time.Duration
	MaintainInterval time.Duration
	LoadBalanced     bool
	PoolMonitor      *event.PoolMonitor
	Logger           *logger.Logger
	handshakeErrFn   func(error, uint64, *primitive.ObjectID)
}

type pool struct {
	// The following integer fields must be accessed using the atomic package
	// and should be at the beginning of the struct.
	// - atomic bug: https://pkg.go.dev/sync/atomic#pkg-note-BUG
	// - suggested layout: https://go101.org/article/memory-layout.html

	nextID                       uint64 // nextID is the next pool ID for a new connection.
	pinnedCursorConnections      uint64
	pinnedTransactionConnections uint64

	address       address.Address
	minSize       uint64
	maxSize       uint64
	maxConnecting uint64
	loadBalanced  bool
	monitor       *event.PoolMonitor
	logger        *logger.Logger

	// handshakeErrFn is used to handle any errors that happen during connection establishment and
	// handshaking.
	handshakeErrFn func(error, uint64, *primitive.ObjectID)

	connOpts   []ConnectionOption
	generation *poolGenerationMap

	maintainInterval time.Duration   // maintainInterval is the maintain() loop interval.
	maintainReady    chan struct{}   // maintainReady is a signal channel that starts the maintain() loop when ready() is called.
	backgroundDone   *sync.WaitGroup // backgroundDone waits for all background goroutines to return.

	stateMu      sync.RWMutex // stateMu guards state, lastClearErr
	state        int          // state is the current state of the connection pool.
	lastClearErr error        // lastClearErr is the last error that caused the pool to be cleared.

	// createConnectionsCond is the condition variable that controls when the createConnections()
	// loop runs or waits. Its lock guards cancelBackgroundCtx, conns, and newConnWait. Any changes
	// to the state of the guarded values must be made while holding the lock to prevent undefined
	// behavior in the createConnections() waiting logic.
	createConnectionsCond *sync.Cond
	cancelBackgroundCtx   context.CancelFunc     // cancelBackgroundCtx is called to signal background goroutines to stop.
	conns                 map[uint64]*connection // conns holds all currently open connections.
	newConnWait           wantConnQueue          // newConnWait holds all wantConn requests for new connections.

	idleMu       sync.Mutex    // idleMu guards idleConns, idleConnWait
	idleConns    []*connection // idleConns holds all idle connections.
	idleConnWait wantConnQueue // idleConnWait holds all wantConn requests for idle connections.
}

// getState returns the current state of the pool. Callers must not hold the stateMu lock.
func (p *pool) getState() int {
	p.stateMu.RLock()
	defer p.stateMu.RUnlock()

	return p.state
}

func mustLogPoolMessage(pool *pool) bool {
	return pool.logger != nil && pool.logger.LevelComponentEnabled(
		logger.LevelDebug, logger.ComponentConnection)
}

func logPoolMessage(pool *pool, msg string, keysAndValues ...interface{}) {
	host, port, err := net.SplitHostPort(pool.address.String())
	if err != nil {
		host = pool.address.String()
		port = ""
	}

	pool.logger.Print(logger.LevelDebug,
		logger.ComponentConnection,
		msg,
		logger.SerializeConnection(logger.Connection{
			Message:    msg,
			ServerHost: host,
			ServerPort: port,
		}, keysAndValues...)...)

}

type reason struct {
	loggerConn string
	event      string
}

// connectionPerished checks if a given connection is perished and should be removed from the pool.
func connectionPerished(conn *connection) (reason, bool) {
	switch {
	case conn.closed() || !conn.isAlive():
		// A connection would only be closed if it encountered a network error
		// during an operation and closed itself. If a connection is not alive
		// (e.g. the connection was closed by the server-side), it's also
		// considered a network error.
		return reason{
			loggerConn: logger.ReasonConnClosedError,
			event:      event.ReasonError,
		}, true
	case conn.idleTimeoutExpired():
		return reason{
			loggerConn: logger.ReasonConnClosedIdle,
			event:      event.ReasonIdle,
		}, true
	case conn.pool.stale(conn):
		return reason{
			loggerConn: logger.ReasonConnClosedStale,
			event:      event.ReasonStale,
		}, true
	}

	return reason{}, false
}

// newPool creates a new pool. It will use the provided options when creating connections.
func newPool(config poolConfig, connOpts ...ConnectionOption) *pool {
	if config.MaxIdleTime != time.Duration(0) {
		connOpts = append(connOpts, WithIdleTimeout(func(_ time.Duration) time.Duration { return config.MaxIdleTime }))
	}

	var maxConnecting uint64 = 2
	if config.MaxConnecting > 0 {
		maxConnecting = config.MaxConnecting
	}

	maintainInterval := 10 * time.Second
	if config.MaintainInterval != 0 {
		maintainInterval = config.MaintainInterval
	}

	pool := &pool{
		address:               config.Address,
		minSize:               config.MinPoolSize,
		maxSize:               config.MaxPoolSize,
		maxConnecting:         maxConnecting,
		loadBalanced:          config.LoadBalanced,
		monitor:               config.PoolMonitor,
		logger:                config.Logger,
		handshakeErrFn:        config.handshakeErrFn,
		connOpts:              connOpts,
		generation:            newPoolGenerationMap(),
		state:                 poolPaused,
		maintainInterval:      maintainInterval,
		maintainReady:         make(chan struct{}, 1),
		backgroundDone:        &sync.WaitGroup{},
		createConnectionsCond: sync.NewCond(&sync.Mutex{}),
		conns:                 make(map[uint64]*connection, config.MaxPoolSize),
		idleConns:             make([]*connection, 0, config.MaxPoolSize),
	}
	// minSize must not exceed maxSize if maxSize is not 0
	if pool.maxSize != 0 && pool.minSize > pool.maxSize {
		pool.minSize = pool.maxSize
	}
	pool.connOpts = append(pool.connOpts, withGenerationNumberFn(func(_ generationNumberFn) generationNumberFn { return pool.getGenerationForNewConnection }))

	pool.generation.connect()

	// Create a Context with cancellation that's used to signal the createConnections() and
	// maintain() background goroutines to stop. Also create a "backgroundDone" WaitGroup that is
	// used to wait for the background goroutines to return.
	var ctx context.Context
	ctx, pool.cancelBackgroundCtx = context.WithCancel(context.Background())

	for i := 0; i < int(pool.maxConnecting); i++ {
		pool.backgroundDone.Add(1)
		go pool.createConnections(ctx, pool.backgroundDone)
	}

	// If maintainInterval is not positive, don't start the maintain() goroutine. Expect that
	// negative values are only used in testing; this config value is not user-configurable.
	if maintainInterval > 0 {
		pool.backgroundDone.Add(1)
		go pool.maintain(ctx, pool.backgroundDone)
	}

	if mustLogPoolMessage(pool) {
		keysAndValues := logger.KeyValues{
			logger.KeyMaxIdleTimeMS, config.MaxIdleTime.Milliseconds(),
			logger.KeyMinPoolSize, config.MinPoolSize,
			logger.KeyMaxPoolSize, config.MaxPoolSize,
			logger.KeyMaxConnecting, config.MaxConnecting,
		}

		logPoolMessage(pool, logger.ConnectionPoolCreated, keysAndValues...)
	}

	if pool.monitor != nil {
		pool.monitor.Event(&event.PoolEvent{
			Type: event.PoolCreated,
			PoolOptions: &event.MonitorPoolOptions{
				MaxPoolSize: config.MaxPoolSize,
				MinPoolSize: config.MinPoolSize,
			},
			Address: pool.address.String(),
		})
	}

	return pool
}

// stale checks if a given connection's generation is below the generation of the pool
func (p *pool) stale(conn *connection) bool {
	return conn == nil || p.generation.stale(conn.desc.ServiceID, conn.generation)
}

// ready puts the pool into the "ready" state and starts the background connection creation and
// monitoring goroutines. ready must be called before connections can be checked out. An unused,
// connected pool must be closed or it will leak goroutines and will not be garbage collected.
func (p *pool) ready() error {
	// While holding the stateMu lock, set the pool to "ready" if it is currently "paused".
	p.stateMu.Lock()
	if p.state == poolReady {
		p.stateMu.Unlock()
		return nil
	}
	if p.state != poolPaused {
		p.stateMu.Unlock()
		return ErrPoolNotPaused
	}
	p.lastClearErr = nil
	p.state = poolReady
	p.stateMu.Unlock()

	if mustLogPoolMessage(p) {
		logPoolMessage(p, logger.ConnectionPoolReady)
	}

	// Send event.PoolReady before resuming the maintain() goroutine to guarantee that the
	// "pool ready" event is always sent before maintain() starts creating connections.
	if p.monitor != nil {
		p.monitor.Event(&event.PoolEvent{
			Type:    event.PoolReady,
			Address: p.address.String(),
		})
	}

	// Signal maintain() to wake up immediately when marking the pool "ready".
	select {
	case p.maintainReady <- struct{}{}:
	default:
	}

	return nil
}

// close closes the pool, closes all connections associated with the pool, and stops all background
// goroutines. All subsequent checkOut requests will return an error. An unused, ready pool must be
// closed or it will leak goroutines and will not be garbage collected.
func (p *pool) close(ctx context.Context) {
	p.stateMu.Lock()
	if p.state == poolClosed {
		p.stateMu.Unlock()
		return
	}
	p.state = poolClosed
	p.stateMu.Unlock()

	// Call cancelBackgroundCtx() to exit the maintain() and createConnections() background
	// goroutines. Broadcast to the createConnectionsCond to wake up all createConnections()
	// goroutines. We must hold the createConnectionsCond lock here because we're changing the
	// condition by cancelling the "background goroutine" Context, even tho cancelling the Context
	// is also synchronized by a lock. Otherwise, we run into an intermittent bug that prevents the
	// createConnections() goroutines from exiting.
	p.createConnectionsCond.L.Lock()
	p.cancelBackgroundCtx()
	p.createConnectionsCond.Broadcast()
	p.createConnectionsCond.L.Unlock()

	// Wait for all background goroutines to exit.
	p.backgroundDone.Wait()

	p.generation.disconnect()

	if ctx == nil {
		ctx = context.Background()
	}

	// If we have a deadline then we interpret it as a request to gracefully shutdown. We wait until
	// either all the connections have been checked back into the pool (i.e. total open connections
	// equals idle connections) or until the Context deadline is reached.
	if _, ok := ctx.Deadline(); ok {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

	graceful:
		for {
			if p.totalConnectionCount() == p.availableConnectionCount() {
				break graceful
			}

			select {
			case <-ticker.C:
			case <-ctx.Done():
				break graceful
			default:
			}
		}
	}

	// Empty the idle connections stack and try to deliver ErrPoolClosed to any waiting wantConns
	// from idleConnWait while holding the idleMu lock.
	p.idleMu.Lock()
	for _, conn := range p.idleConns {
		_ = p.removeConnection(conn, reason{
			loggerConn: logger.ReasonConnClosedPoolClosed,
			event:      event.ReasonPoolClosed,
		}, nil)
		_ = p.closeConnection(conn) // We don't care about errors while closing the connection.
	}
	p.idleConns = p.idleConns[:0]
	for {
		w := p.idleConnWait.popFront()
		if w == nil {
			break
		}
		w.tryDeliver(nil, ErrPoolClosed)
	}
	p.idleMu.Unlock()

	// Collect all conns from the pool and try to deliver ErrPoolClosed to any waiting wantConns
	// from newConnWait while holding the createConnectionsCond lock. We can't call removeConnection
	// on the connections while holding any locks, so do that after we release the lock.
	p.createConnectionsCond.L.Lock()
	conns := make([]*connection, 0, len(p.conns))
	for _, conn := range p.conns {
		conns = append(conns, conn)
	}
	for {
		w := p.newConnWait.popFront()
		if w == nil {
			break
		}
		w.tryDeliver(nil, ErrPoolClosed)
	}
	p.createConnectionsCond.L.Unlock()

	if mustLogPoolMessage(p) {
		logPoolMessage(p, logger.ConnectionPoolClosed)
	}

	if p.monitor != nil {
		p.monitor.Event(&event.PoolEvent{
			Type:    event.PoolClosedEvent,
			Address: p.address.String(),
		})
	}

	// Now that we're not holding any locks, remove all of the connections we collected from the
	// pool.
	for _, conn := range conns {
		_ = p.removeConnection(conn, reason{
			loggerConn: logger.ReasonConnClosedPoolClosed,
			event:      event.ReasonPoolClosed,
		}, nil)
		_ = p.closeConnection(conn) // We don't care about errors while closing the connection.
	}
}

func (p *pool) pinConnectionToCursor() {
	atomic.AddUint64(&p.pinnedCursorConnections, 1)
}

func (p *pool) unpinConnectionFromCursor() {
	// See https://golang.org/pkg/sync/atomic/#AddUint64 for an explanation of the ^uint64(0) syntax.
	atomic.AddUint64(&p.pinnedCursorConnections, ^uint64(0))
}

func (p *pool) pinConnectionToTransaction() {
	atomic.AddUint64(&p.pinnedTransactionConnections, 1)
}

func (p *pool) unpinConnectionFromTransaction() {
	// See https://golang.org/pkg/sync/atomic/#AddUint64 for an explanation of the ^uint64(0) syntax.
	atomic.AddUint64(&p.pinnedTransactionConnections, ^uint64(0))
}

// checkOut checks out a connection from the pool. If an idle connection is not available, the
// checkOut enters a queue waiting for either the next idle or new connection. If the pool is not
// ready, checkOut returns an error.
// Based partially on https://cs.opensource.google/go/go/+/refs/tags/go1.16.6:src/net/http/transport.go;l=1324
func (p *pool) checkOut(ctx context.Context) (conn *connection, err error) {
	if mustLogPoolMessage(p) {
		logPoolMessage(p, logger.ConnectionCheckoutStarted)
	}

	// TODO(CSOT): If a Timeout was specified at any level, respect the Timeout is server selection, connection
	// TODO checkout.
	if p.monitor != nil {
		p.monitor.Event(&event.PoolEvent{
			Type:    event.GetStarted,
			Address: p.address.String(),
		})
	}

	start := time.Now()
	// Check the pool state while holding a stateMu read lock. If the pool state is not "ready",
	// return an error. Do all of this while holding the stateMu read lock to prevent a state change between
	// checking the state and entering the wait queue. Not holding the stateMu read lock here may
	// allow a checkOut() to enter the wait queue after clear() pauses the pool and clears the wait
	// queue, resulting in createConnections() doing work while the pool is "paused".
	p.stateMu.RLock()
	switch p.state {
	case poolClosed:
		p.stateMu.RUnlock()

		duration := time.Since(start)
		if mustLogPoolMessage(p) {
			keysAndValues := logger.KeyValues{
				logger.KeyDurationMS, duration.Milliseconds(),
				logger.KeyReason, logger.ReasonConnCheckoutFailedPoolClosed,
			}

			logPoolMessage(p, logger.ConnectionCheckoutFailed, keysAndValues...)
		}

		if p.monitor != nil {
			p.monitor.Event(&event.PoolEvent{
				Type:     event.GetFailed,
				Address:  p.address.String(),
				Duration: duration,
				Reason:   event.ReasonPoolClosed,
			})
		}
		return nil, ErrPoolClosed
	case poolPaused:
		err := poolClearedError{err: p.lastClearErr, address: p.address}
		p.stateMu.RUnlock()

		duration := time.Since(start)
		if mustLogPoolMessage(p) {
			keysAndValues := logger.KeyValues{
				logger.KeyDurationMS, duration.Milliseconds(),
				logger.KeyReason, logger.ReasonConnCheckoutFailedError,
			}

			logPoolMessage(p, logger.ConnectionCheckoutFailed, keysAndValues...)
		}

		if p.monitor != nil {
			p.monitor.Event(&event.PoolEvent{
				Type:     event.GetFailed,
				Address:  p.address.String(),
				Duration: duration,
				Reason:   event.ReasonConnectionErrored,
				Error:    err,
			})
		}
		return nil, err
	}

	if ctx == nil {
		ctx = context.Background()
	}

	// Create a wantConn, which we will use to request an existing idle or new connection. Always
	// cancel the wantConn if checkOut() returned an error to make sure any delivered connections
	// are returned to the pool (e.g. if a connection was delivered immediately after the Context
	// timed out).
	w := newWantConn()
	defer func() {
		if err != nil {
			w.cancel(p, err)
		}
	}()

	// Get in the queue for an idle connection. If getOrQueueForIdleConn returns true, it was able to
	// immediately deliver an idle connection to the wantConn, so we can return the connection or
	// error from the wantConn without waiting for "ready".
	if delivered := p.getOrQueueForIdleConn(w); delivered {
		// If delivered = true, we didn't enter the wait queue and will return either a connection
		// or an error, so unlock the stateMu lock here.
		p.stateMu.RUnlock()

		duration := time.Since(start)
		if w.err != nil {
			if mustLogPoolMessage(p) {
				keysAndValues := logger.KeyValues{
					logger.KeyDurationMS, duration.Milliseconds(),
					logger.KeyReason, logger.ReasonConnCheckoutFailedError,
				}

				logPoolMessage(p, logger.ConnectionCheckoutFailed, keysAndValues...)
			}

			if p.monitor != nil {
				p.monitor.Event(&event.PoolEvent{
					Type:     event.GetFailed,
					Address:  p.address.String(),
					Duration: duration,
					Reason:   event.ReasonConnectionErrored,
					Error:    w.err,
				})
			}
			return nil, w.err
		}

		duration = time.Since(start)
		if mustLogPoolMessage(p) {
			keysAndValues := logger.KeyValues{
				logger.KeyDriverConnectionID, w.conn.driverConnectionID,
				logger.KeyDurationMS, duration.Milliseconds(),
			}

			logPoolMessage(p, logger.ConnectionCheckedOut, keysAndValues...)
		}

		if p.monitor != nil {
			p.monitor.Event(&event.PoolEvent{
				Type:         event.GetSucceeded,
				Address:      p.address.String(),
				ConnectionID: w.conn.driverConnectionID,
				Duration:     duration,
			})
		}

		return w.conn, nil
	}

	// If we didn't get an immediately available idle connection, also get in the queue for a new
	// connection while we're waiting for an idle connection.
	p.queueForNewConn(w)
	p.stateMu.RUnlock()

	// Wait for either the wantConn to be ready or for the Context to time out.
	waitQueueStart := time.Now()
	select {
	case <-w.ready:
		if w.err != nil {
			duration := time.Since(start)
			if mustLogPoolMessage(p) {
				keysAndValues := logger.KeyValues{
					logger.KeyDurationMS, duration.Milliseconds(),
					logger.KeyReason, logger.ReasonConnCheckoutFailedError,
					logger.KeyError, w.err.Error(),
				}

				logPoolMessage(p, logger.ConnectionCheckoutFailed, keysAndValues...)
			}

			if p.monitor != nil {
				p.monitor.Event(&event.PoolEvent{
					Type:     event.GetFailed,
					Address:  p.address.String(),
					Duration: duration,
					Reason:   event.ReasonConnectionErrored,
					Error:    w.err,
				})
			}

			return nil, w.err
		}

		duration := time.Since(start)
		if mustLogPoolMessage(p) {
			keysAndValues := logger.KeyValues{
				logger.KeyDriverConnectionID, w.conn.driverConnectionID,
				logger.KeyDurationMS, duration.Milliseconds(),
			}

			logPoolMessage(p, logger.ConnectionCheckedOut, keysAndValues...)
		}

		if p.monitor != nil {
			p.monitor.Event(&event.PoolEvent{
				Type:         event.GetSucceeded,
				Address:      p.address.String(),
				ConnectionID: w.conn.driverConnectionID,
				Duration:     duration,
			})
		}
		return w.conn, nil
	case <-ctx.Done():
		waitQueueDuration := time.Since(waitQueueStart)

		duration := time.Since(start)
		if mustLogPoolMessage(p) {
			keysAndValues := logger.KeyValues{
				logger.KeyDurationMS, duration.Milliseconds(),
				logger.KeyReason, logger.ReasonConnCheckoutFailedTimout,
			}

			logPoolMessage(p, logger.ConnectionCheckoutFailed, keysAndValues...)
		}

		if p.monitor != nil {
			p.monitor.Event(&event.PoolEvent{
				Type:     event.GetFailed,
				Address:  p.address.String(),
				Duration: duration,
				Reason:   event.ReasonTimedOut,
				Error:    ctx.Err(),
			})
		}

		err := WaitQueueTimeoutError{
			Wrapped:              ctx.Err(),
			maxPoolSize:          p.maxSize,
			totalConnections:     p.totalConnectionCount(),
			availableConnections: p.availableConnectionCount(),
			waitDuration:         waitQueueDuration,
		}
		if p.loadBalanced {
			err.pinnedConnections = &pinnedConnections{
				cursorConnections:      atomic.LoadUint64(&p.pinnedCursorConnections),
				transactionConnections: atomic.LoadUint64(&p.pinnedTransactionConnections),
			}
		}
		return nil, err
	}
}

// closeConnection closes a connection.
func (p *pool) closeConnection(conn *connection) error {
	if conn.pool != p {
		return ErrWrongPool
	}

	if atomic.LoadInt64(&conn.state) == connConnected {
		conn.closeConnectContext()
		conn.wait() // Make sure that the connection has finished connecting.
	}

	err := conn.close()
	if err != nil {
		return ConnectionError{ConnectionID: conn.id, Wrapped: err, message: "failed to close net.Conn"}
	}

	return nil
}

func (p *pool) getGenerationForNewConnection(serviceID *primitive.ObjectID) uint64 {
	return p.generation.addConnection(serviceID)
}

// removeConnection removes a connection from the pool and emits a "ConnectionClosed" event.
func (p *pool) removeConnection(conn *connection, reason reason, err error) error {
	if conn == nil {
		return nil
	}

	if conn.pool != p {
		return ErrWrongPool
	}

	p.createConnectionsCond.L.Lock()
	_, ok := p.conns[conn.driverConnectionID]
	if !ok {
		// If the connection has been removed from the pool already, exit without doing any
		// additional state changes.
		p.createConnectionsCond.L.Unlock()
		return nil
	}
	delete(p.conns, conn.driverConnectionID)
	// Signal the createConnectionsCond so any goroutines waiting for a new connection slot in the
	// pool will proceed.
	p.createConnectionsCond.Signal()
	p.createConnectionsCond.L.Unlock()

	// Only update the generation numbers map if the connection has retrieved its generation number.
	// Otherwise, we'd decrement the count for the generation even though it had never been
	// incremented.
	if conn.hasGenerationNumber() {
		p.generation.removeConnection(conn.desc.ServiceID)
	}

	if mustLogPoolMessage(p) {
		keysAndValues := logger.KeyValues{
			logger.KeyDriverConnectionID, conn.driverConnectionID,
			logger.KeyReason, reason.loggerConn,
		}

		if err != nil {
			keysAndValues.Add(logger.KeyError, err.Error())
		}

		logPoolMessage(p, logger.ConnectionClosed, keysAndValues...)
	}

	if p.monitor != nil {
		p.monitor.Event(&event.PoolEvent{
			Type:         event.ConnectionClosed,
			Address:      p.address.String(),
			ConnectionID: conn.driverConnectionID,
			Reason:       reason.event,
			Error:        err,
		})
	}

	return nil
}

var (
	// BGReadTimeout is the maximum amount of the to wait when trying to read
	// the server reply on a connection after an operation timed out. The
	// default is 1 second.
	//
	// Deprecated: BGReadTimeout is intended for internal use only and may be
	// removed or modified at any time.
	BGReadTimeout = 1 * time.Second

	// BGReadCallback is a callback for monitoring the behavior of the
	// background-read-on-timeout connection preserving mechanism.
	//
	// Deprecated: BGReadCallback is intended for internal use only and may be
	// removed or modified at any time.
	BGReadCallback func(addr string, start, read time.Time, errs []error, connClosed bool)
)

// bgRead sets a new read deadline on the provided connection (1 second in the
// future) and tries to read any bytes returned by the server. If successful, it
// checks the connection into the provided pool. If there are any errors, it
// closes the connection.
//
// It calls the package-global BGReadCallback function, if set, with the
// address, timings, and any errors that occurred.
func bgRead(pool *pool, conn *connection, size int32) {
	var err error
	start := time.Now()

	defer func() {
		read := time.Now()
		errs := make([]error, 0)
		connClosed := false
		if err != nil {
			errs = append(errs, err)
			connClosed = true
			err = conn.close()
			if err != nil {
				errs = append(errs, fmt.Errorf("error closing conn after reading: %w", err))
			}
		}

		// No matter what happens, always check the connection back into the
		// pool, which will either make it available for other operations or
		// remove it from the pool if it was closed.
		err = pool.checkInNoEvent(conn)
		if err != nil {
			errs = append(errs, fmt.Errorf("error checking in: %w", err))
		}

		if BGReadCallback != nil {
			BGReadCallback(conn.addr.String(), start, read, errs, connClosed)
		}
	}()

	err = conn.nc.SetReadDeadline(time.Now().Add(BGReadTimeout))
	if err != nil {
		err = fmt.Errorf("error setting a read deadline: %w", err)
		return
	}

	if size == 0 {
		var sizeBuf [4]byte
		_, err = io.ReadFull(conn.nc, sizeBuf[:])
		if err != nil {
			err = fmt.Errorf("error reading the message size: %w", err)
			return
		}
		size, err = conn.parseWmSizeBytes(sizeBuf)
		if err != nil {
			return
		}
		size -= 4
	}
	_, err = io.CopyN(io.Discard, conn.nc, int64(size))
	if err != nil {
		err = fmt.Errorf("error discarding %d byte message: %w", size, err)
	}
}

// checkIn returns an idle connection to the pool. If the connection is perished or the pool is
// closed, it is removed from the connection pool and closed.
func (p *pool) checkIn(conn *connection) error {
	if conn == nil {
		return nil
	}
	if conn.pool != p {
		return ErrWrongPool
	}

	if mustLogPoolMessage(p) {
		keysAndValues := logger.KeyValues{
			logger.KeyDriverConnectionID, conn.driverConnectionID,
		}

		logPoolMessage(p, logger.ConnectionCheckedIn, keysAndValues...)
	}

	if p.monitor != nil {
		p.monitor.Event(&event.PoolEvent{
			Type:         event.ConnectionReturned,
			ConnectionID: conn.driverConnectionID,
			Address:      conn.addr.String(),
		})
	}

	return p.checkInNoEvent(conn)
}

// checkInNoEvent returns a connection to the pool. It behaves identically to checkIn except it does
// not publish events. It is only intended for use by pool-internal functions.
func (p *pool) checkInNoEvent(conn *connection) error {
	if conn == nil {
		return nil
	}
	if conn.pool != p {
		return ErrWrongPool
	}

	// If the connection has an awaiting server response, try to read the
	// response in another goroutine before checking it back into the pool.
	//
	// Do this here because we want to publish checkIn events when the operation
	// is done with the connection, not when it's ready to be used again. That
	// means that connections in "awaiting response" state are checked in but
	// not usable, which is not covered by the current pool events. We may need
	// to add pool event information in the future to communicate that.
	if conn.awaitRemainingBytes != nil {
		size := *conn.awaitRemainingBytes
		conn.awaitRemainingBytes = nil
		go bgRead(p, conn, size)
		return nil
	}

	// Bump the connection idle start time here because we're about to make the
	// connection "available". The idle start time is used to determine how long
	// a connection has been idle and when it has reached its max idle time and
	// should be closed. A connection reaches its max idle time when it has been
	// "available" in the idle connections stack for more than the configured
	// duration (maxIdleTimeMS). Set it before we call connectionPerished(),
	// which checks the idle deadline, because a newly "available" connection
	// should never be perished due to max idle time.
	conn.bumpIdleStart()

	r, perished := connectionPerished(conn)
	if !perished && conn.pool.getState() == poolClosed {
		perished = true
		r = reason{
			loggerConn: logger.ReasonConnClosedPoolClosed,
			event:      event.ReasonPoolClosed,
		}
	}
	if perished {
		_ = p.removeConnection(conn, r, nil)
		go func() {
			_ = p.closeConnection(conn)
		}()
		return nil
	}

	p.idleMu.Lock()
	defer p.idleMu.Unlock()

	for {
		w := p.idleConnWait.popFront()
		if w == nil {
			break
		}
		if w.tryDeliver(conn, nil) {
			return nil
		}
	}

	for _, idle := range p.idleConns {
		if idle == conn {
			return fmt.Errorf("duplicate idle conn %p in idle connections stack", conn)
		}
	}

	p.idleConns = append(p.idleConns, conn)
	return nil
}

// clear calls clearImpl internally with a false interruptAllConnections value.
func (p *pool) clear(err error, serviceID *primitive.ObjectID) {
	p.clearImpl(err, serviceID, false)
}

// clearAll does same as the "clear" method but interrupts all connections.
func (p *pool) clearAll(err error, serviceID *primitive.ObjectID) {
	p.clearImpl(err, serviceID, true)
}

// interruptConnections interrupts the input connections.
func (p *pool) interruptConnections(conns []*connection) {
	for _, conn := range conns {
		_ = p.removeConnection(conn, reason{
			loggerConn: logger.ReasonConnClosedStale,
			event:      event.ReasonStale,
		}, nil)
		go func(c *connection) {
			_ = p.closeConnection(c)
		}(conn)
	}
}

// clear marks all connections as stale by incrementing the generation number, stops all background
// goroutines, removes all requests from idleConnWait and newConnWait, and sets the pool state to
// "paused". If serviceID is nil, clear marks all connections as stale. If serviceID is not nil,
// clear marks only connections associated with the given serviceID stale (for use in load balancer
// mode).
// If interruptAllConnections is true, this function calls interruptConnections to interrupt all
// non-idle connections.
func (p *pool) clearImpl(err error, serviceID *primitive.ObjectID, interruptAllConnections bool) {
	if p.getState() == poolClosed {
		return
	}

	p.generation.clear(serviceID)

	// If serviceID is nil (i.e. not in load balancer mode), transition the pool to a paused state
	// by stopping all background goroutines, clearing the wait queues, and setting the pool state
	// to "paused".
	sendEvent := true
	if serviceID == nil {
		// While holding the stateMu lock, set the pool state to "paused" if it's currently "ready",
		// and set lastClearErr to the error that caused the pool to be cleared. If the pool is
		// already paused, don't send another "ConnectionPoolCleared" event.
		p.stateMu.Lock()
		if p.state == poolPaused {
			sendEvent = false
		}
		if p.state == poolReady {
			p.state = poolPaused
		}
		p.lastClearErr = err
		p.stateMu.Unlock()
	}

	if mustLogPoolMessage(p) {
		keysAndValues := logger.KeyValues{
			logger.KeyServiceID, serviceID,
		}

		logPoolMessage(p, logger.ConnectionPoolCleared, keysAndValues...)
	}

	if sendEvent && p.monitor != nil {
		event := &event.PoolEvent{
			Type:         event.PoolCleared,
			Address:      p.address.String(),
			ServiceID:    serviceID,
			Interruption: interruptAllConnections,
			Error:        err,
		}
		p.monitor.Event(event)
	}

	p.removePerishedConns()
	if interruptAllConnections {
		p.createConnectionsCond.L.Lock()
		p.idleMu.Lock()

		idleConns := make(map[*connection]bool, len(p.idleConns))
		for _, idle := range p.idleConns {
			idleConns[idle] = true
		}

		conns := make([]*connection, 0, len(p.conns))
		for _, conn := range p.conns {
			if _, ok := idleConns[conn]; !ok && p.stale(conn) {
				conns = append(conns, conn)
			}
		}

		p.idleMu.Unlock()
		p.createConnectionsCond.L.Unlock()

		p.interruptConnections(conns)
	}

	if serviceID == nil {
		pcErr := poolClearedError{err: err, address: p.address}

		// Clear the idle connections wait queue.
		p.idleMu.Lock()
		for {
			w := p.idleConnWait.popFront()
			if w == nil {
				break
			}
			w.tryDeliver(nil, pcErr)
		}
		p.idleMu.Unlock()

		// Clear the new connections wait queue. This effectively pauses the createConnections()
		// background goroutine because newConnWait is empty and checkOut() won't insert any more
		// wantConns into newConnWait until the pool is marked "ready" again.
		p.createConnectionsCond.L.Lock()
		for {
			w := p.newConnWait.popFront()
			if w == nil {
				break
			}
			w.tryDeliver(nil, pcErr)
		}
		p.createConnectionsCond.L.Unlock()
	}
}

// getOrQueueForIdleConn attempts to deliver an idle connection to the given wantConn. If there is
// an idle connection in the idle connections stack, it pops an idle connection, delivers it to the
// wantConn, and returns true. If there are no idle connections in the idle connections stack, it
// adds the wantConn to the idleConnWait queue and returns false.
func (p *pool) getOrQueueForIdleConn(w *wantConn) bool {
	p.idleMu.Lock()
	defer p.idleMu.Unlock()

	// Try to deliver an idle connection from the idleConns stack first.
	for len(p.idleConns) > 0 {
		conn := p.idleConns[len(p.idleConns)-1]
		p.idleConns = p.idleConns[:len(p.idleConns)-1]

		if conn == nil {
			continue
		}

		if reason, perished := connectionPerished(conn); perished {
			_ = conn.pool.removeConnection(conn, reason, nil)
			go func() {
				_ = conn.pool.closeConnection(conn)
			}()
			continue
		}

		if !w.tryDeliver(conn, nil) {
			// If we couldn't deliver the conn to w, put it back in the idleConns stack.
			p.idleConns = append(p.idleConns, conn)
		}

		// If we got here, we tried to deliver an idle conn to w. No matter if tryDeliver() returned
		// true or false, w is no longer waiting and doesn't need to be added to any wait queues, so
		// return delivered = true.
		return true
	}

	p.idleConnWait.cleanFront()
	p.idleConnWait.pushBack(w)
	return false
}

func (p *pool) queueForNewConn(w *wantConn) {
	p.createConnectionsCond.L.Lock()
	defer p.createConnectionsCond.L.Unlock()

	p.newConnWait.cleanFront()
	p.newConnWait.pushBack(w)
	p.createConnectionsCond.Signal()
}

func (p *pool) totalConnectionCount() int {
	p.createConnectionsCond.L.Lock()
	defer p.createConnectionsCond.L.Unlock()

	return len(p.conns)
}

func (p *pool) availableConnectionCount() int {
	p.idleMu.Lock()
	defer p.idleMu.Unlock()

	return len(p.idleConns)
}

// createConnections creates connections for wantConn requests on the newConnWait queue.
func (p *pool) createConnections(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	// condition returns true if the createConnections() loop should continue and false if it should
	// wait. Note that the condition also listens for Context cancellation, which also causes the
	// loop to continue, allowing for a subsequent check to return from createConnections().
	condition := func() bool {
		checkOutWaiting := p.newConnWait.len() > 0
		poolHasSpace := p.maxSize == 0 || uint64(len(p.conns)) < p.maxSize
		cancelled := ctx.Err() != nil
		return (checkOutWaiting && poolHasSpace) || cancelled
	}

	// wait waits for there to be an available wantConn and for the pool to have space for a new
	// connection. When the condition becomes true, it creates a new connection and returns the
	// waiting wantConn and new connection. If the Context is cancelled or there are any
	// errors, wait returns with "ok = false".
	wait := func() (*wantConn, *connection, bool) {
		p.createConnectionsCond.L.Lock()
		defer p.createConnectionsCond.L.Unlock()

		for !condition() {
			p.createConnectionsCond.Wait()
		}

		if ctx.Err() != nil {
			return nil, nil, false
		}

		p.newConnWait.cleanFront()
		w := p.newConnWait.popFront()
		if w == nil {
			return nil, nil, false
		}

		conn := newConnection(p.address, p.connOpts...)
		conn.pool = p
		conn.driverConnectionID = atomic.AddUint64(&p.nextID, 1)
		p.conns[conn.driverConnectionID] = conn

		return w, conn, true
	}

	for ctx.Err() == nil {
		w, conn, ok := wait()
		if !ok {
			continue
		}

		if mustLogPoolMessage(p) {
			keysAndValues := logger.KeyValues{
				logger.KeyDriverConnectionID, conn.driverConnectionID,
			}

			logPoolMessage(p, logger.ConnectionCreated, keysAndValues...)
		}

		if p.monitor != nil {
			p.monitor.Event(&event.PoolEvent{
				Type:         event.ConnectionCreated,
				Address:      p.address.String(),
				ConnectionID: conn.driverConnectionID,
			})
		}

		start := time.Now()
		// Pass the createConnections context to connect to allow pool close to cancel connection
		// establishment so shutdown doesn't block indefinitely if connectTimeout=0.
		err := conn.connect(ctx)
		if err != nil {
			w.tryDeliver(nil, err)

			// If there's an error connecting the new connection, call the handshake error handler
			// that implements the SDAM handshake error handling logic. This must be called after
			// delivering the connection error to the waiting wantConn. If it's called before, the
			// handshake error handler may clear the connection pool, leading to a different error
			// message being delivered to the same waiting wantConn in idleConnWait when the wait
			// queues are cleared.
			if p.handshakeErrFn != nil {
				p.handshakeErrFn(err, conn.generation, conn.desc.ServiceID)
			}

			_ = p.removeConnection(conn, reason{
				loggerConn: logger.ReasonConnClosedError,
				event:      event.ReasonError,
			}, err)

			_ = p.closeConnection(conn)

			continue
		}

		duration := time.Since(start)
		if mustLogPoolMessage(p) {
			keysAndValues := logger.KeyValues{
				logger.KeyDriverConnectionID, conn.driverConnectionID,
				logger.KeyDurationMS, duration.Milliseconds(),
			}

			logPoolMessage(p, logger.ConnectionReady, keysAndValues...)
		}

		if p.monitor != nil {
			p.monitor.Event(&event.PoolEvent{
				Type:         event.ConnectionReady,
				Address:      p.address.String(),
				ConnectionID: conn.driverConnectionID,
				Duration:     duration,
			})
		}

		if w.tryDeliver(conn, nil) {
			continue
		}

		_ = p.checkInNoEvent(conn)
	}
}

func (p *pool) maintain(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(p.maintainInterval)
	defer ticker.Stop()

	// remove removes the *wantConn at index i from the slice and returns the new slice. The order
	// of the slice is not maintained.
	remove := func(arr []*wantConn, i int) []*wantConn {
		end := len(arr) - 1
		arr[i], arr[end] = arr[end], arr[i]
		return arr[:end]
	}

	// removeNotWaiting removes any wantConns that are no longer waiting from given slice of
	// wantConns. That allows maintain() to use the size of its wantConns slice as an indication of
	// how many new connection requests are outstanding and subtract that from the number of
	// connections to ask for when maintaining minPoolSize.
	removeNotWaiting := func(arr []*wantConn) []*wantConn {
		for i := len(arr) - 1; i >= 0; i-- {
			w := arr[i]
			if !w.waiting() {
				arr = remove(arr, i)
			}
		}

		return arr
	}

	wantConns := make([]*wantConn, 0, p.minSize)
	defer func() {
		for _, w := range wantConns {
			w.tryDeliver(nil, ErrPoolClosed)
		}
	}()

	for {
		select {
		case <-ticker.C:
		case <-p.maintainReady:
		case <-ctx.Done():
			return
		}

		// Only maintain the pool while it's in the "ready" state. If the pool state is not "ready",
		// wait for the next tick or "ready" signal. Do all of this while holding the stateMu read
		// lock to prevent a state change between checking the state and entering the wait queue.
		// Not holding the stateMu read lock here may allow maintain() to request wantConns after
		// clear() pauses the pool and clears the wait queue, resulting in createConnections()
		// doing work while the pool is "paused".
		p.stateMu.RLock()
		if p.state != poolReady {
			p.stateMu.RUnlock()
			continue
		}

		p.removePerishedConns()

		// Remove any wantConns that are no longer waiting.
		wantConns = removeNotWaiting(wantConns)

		// Figure out how many more wantConns we need to satisfy minPoolSize. Assume that the
		// outstanding wantConns (i.e. the ones that weren't removed from the slice) will all return
		// connections when they're ready, so only add wantConns to make up the difference. Limit
		// the number of connections requested to max 10 at a time to prevent overshooting
		// minPoolSize in case other checkOut() calls are requesting new connections, too.
		total := p.totalConnectionCount()
		n := int(p.minSize) - total - len(wantConns)
		if n > 10 {
			n = 10
		}

		for i := 0; i < n; i++ {
			w := newWantConn()
			p.queueForNewConn(w)
			wantConns = append(wantConns, w)

			// Start a goroutine for each new wantConn, waiting for it to be ready.
			go func() {
				<-w.ready
				if w.conn != nil {
					_ = p.checkInNoEvent(w.conn)
				}
			}()
		}
		p.stateMu.RUnlock()
	}
}

func (p *pool) removePerishedConns() {
	p.idleMu.Lock()
	defer p.idleMu.Unlock()

	for i := range p.idleConns {
		conn := p.idleConns[i]
		if conn == nil {
			continue
		}

		if reason, perished := connectionPerished(conn); perished {
			p.idleConns[i] = nil

			_ = p.removeConnection(conn, reason, nil)
			go func() {
				_ = p.closeConnection(conn)
			}()
		}
	}

	p.idleConns = compact(p.idleConns)
}

// compact removes any nil pointers from the slice and keeps the non-nil pointers, retaining the
// order of the non-nil pointers.
func compact(arr []*connection) []*connection {
	offset := 0
	for i := range arr {
		if arr[i] == nil {
			continue
		}
		arr[offset] = arr[i]
		offset++
	}
	return arr[:offset]
}

// A wantConn records state about a wanted connection (that is, an active call to checkOut).
// The conn may be gotten by creating a new connection or by finding an idle connection, or a
// cancellation may make the conn no longer wanted. These three options are racing against each
// other and use wantConn to coordinate and agree about the winning outcome.
// Based on https://cs.opensource.google/go/go/+/refs/tags/go1.16.6:src/net/http/transport.go;l=1174-1240
type wantConn struct {
	ready chan struct{}

	mu   sync.Mutex // Guards conn, err
	conn *connection
	err  error
}

func newWantConn() *wantConn {
	return &wantConn{
		ready: make(chan struct{}, 1),
	}
}

// waiting reports whether w is still waiting for an answer (connection or error).
func (w *wantConn) waiting() bool {
	select {
	case <-w.ready:
		return false
	default:
		return true
	}
}

// tryDeliver attempts to deliver conn, err to w and reports whether it succeeded.
func (w *wantConn) tryDeliver(conn *connection, err error) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.conn != nil || w.err != nil {
		return false
	}

	w.conn = conn
	w.err = err
	if w.conn == nil && w.err == nil {
		panic("x/mongo/driver/topology: internal error: misuse of tryDeliver")
	}

	close(w.ready)

	return true
}

// cancel marks w as no longer wanting a result (for example, due to cancellation). If a connection
// has been delivered already, cancel returns it with p.checkInNoEvent(). Note that the caller must
// not hold any locks on the pool while calling cancel.
func (w *wantConn) cancel(p *pool, err error) {
	if err == nil {
		panic("x/mongo/driver/topology: internal error: misuse of cancel")
	}

	w.mu.Lock()
	if w.conn == nil && w.err == nil {
		close(w.ready) // catch misbehavior in future delivery
	}
	conn := w.conn
	w.conn = nil
	w.err = err
	w.mu.Unlock()

	if conn != nil {
		_ = p.checkInNoEvent(conn)
	}
}

// A wantConnQueue is a queue of wantConns.
// Based on https://cs.opensource.google/go/go/+/refs/tags/go1.16.6:src/net/http/transport.go;l=1242-1306
type wantConnQueue struct {
	// This is a queue, not a deque.
	// It is split into two stages - head[headPos:] and tail.
	// popFront is trivial (headPos++) on the first stage, and
	// pushBack is trivial (append) on the second stage.
	// If the first stage is empty, popFront can swap the
	// first and second stages to remedy the situation.
	//
	// This two-stage split is analogous to the use of two lists
	// in Okasaki's purely functional queue but without the
	// overhead of reversing the list when swapping stages.
	head    []*wantConn
	headPos int
	tail    []*wantConn
}

// len returns the number of items in the queue.
func (q *wantConnQueue) len() int {
	return len(q.head) - q.headPos + len(q.tail)
}

// pushBack adds w to the back of the queue.
func (q *wantConnQueue) pushBack(w *wantConn) {
	q.tail = append(q.tail, w)
}

// popFront removes and returns the wantConn at the front of the queue.
func (q *wantConnQueue) popFront() *wantConn {
	if q.headPos >= len(q.head) {
		if len(q.tail) == 0 {
			return nil
		}
		// Pick up tail as new head, clear tail.
		q.head, q.headPos, q.tail = q.tail, 0, q.head[:0]
	}
	w := q.head[q.headPos]
	q.head[q.headPos] = nil
	q.headPos++
	return w
}

// peekFront returns the wantConn at the front of the queue without removing it.
func (q *wantConnQueue) peekFront() *wantConn {
	if q.headPos < len(q.head) {
		return q.head[q.headPos]
	}
	if len(q.tail) > 0 {
		return q.tail[0]
	}
	return nil
}

// cleanFront pops any wantConns that are no longer waiting from the head of the queue.
func (q *wantConnQueue) cleanFront() {
	for {
		w := q.peekFront()
		if w == nil || w.waiting() {
			return
		}
		q.popFront()
	}
}
