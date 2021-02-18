// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"context"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/x/mongo/driver/address"
	"golang.org/x/sync/semaphore"
)

// ErrPoolConnected is returned from an attempt to connect an already connected pool
var ErrPoolConnected = PoolError("attempted to Connect to an already connected pool")

// ErrPoolDisconnected is returned from an attempt to Close an already disconnected
// or disconnecting pool.
var ErrPoolDisconnected = PoolError("attempted to check out a connection from closed connection pool")

// ErrConnectionClosed is returned from an attempt to use an already closed connection.
var ErrConnectionClosed = ConnectionError{ConnectionID: "<closed>", message: "connection is closed"}

// ErrWrongPool is return when a connection is returned to a pool it doesn't belong to.
var ErrWrongPool = PoolError("connection does not belong to this pool")

// ErrWaitQueueTimeout is returned when the request to get a connection from the pool timesout when on the wait queue
var ErrWaitQueueTimeout = PoolError("timed out while checking out a connection from connection pool")

// PoolError is an error returned from a Pool method.
type PoolError string

// maintainInterval is the interval at which the background routine to close stale connections will be run.
var maintainInterval = time.Minute

func (pe PoolError) Error() string { return string(pe) }

// poolConfig contains all aspects of the pool that can be configured
type poolConfig struct {
	Address     address.Address
	MinPoolSize uint64
	MaxPoolSize uint64 // MaxPoolSize is not used because handling the max number of connections in the pool is handled in server. This is only used for command monitoring
	MaxIdleTime time.Duration
	PoolMonitor *event.PoolMonitor
}

// checkOutResult is all the values that can be returned from a checkOut
type checkOutResult struct {
	c      *connection
	err    error
	reason string
}

// pool is a wrapper of resource pool that follows the CMAP spec for connection pools
type pool struct {
	address    address.Address
	opts       []ConnectionOption
	conns      *resourcePool // pool for non-checked out connections
	generation uint64        // must be accessed using atomic package
	monitor    *event.PoolMonitor

	connected int32 // Must be accessed using the sync/atomic package.
	nextid    uint64
	opened    map[uint64]*connection // opened holds all of the currently open connections.
	sem       *semaphore.Weighted
	sync.Mutex
}

// connectionExpiredFunc checks if a given connection is stale and should be removed from the resource pool
func connectionExpiredFunc(v interface{}) bool {
	if v == nil {
		return true
	}

	c, ok := v.(*connection)
	if !ok {
		return true
	}

	switch {
	case atomic.LoadInt32(&c.pool.connected) != connected:
		c.expireReason = event.ReasonPoolClosed
	case c.closed():
		// A connection would only be closed if it encountered a network error during an operation and closed itself.
		c.expireReason = event.ReasonConnectionErrored
	case c.idleTimeoutExpired():
		c.expireReason = event.ReasonIdle
	case c.pool.stale(c):
		c.expireReason = event.ReasonStale
	default:
		return false
	}

	return true
}

// connectionCloseFunc closes a given connection. If ctx is nil, the closing will occur in the background
func connectionCloseFunc(v interface{}) {
	c, ok := v.(*connection)
	if !ok || v == nil {
		return
	}

	// The resource pool will only close connections if they're expired or the pool is being disconnected and
	// resourcePool.Close() is called. For the former case, c.expireReason will be set. In the latter, it will not, so
	// we use ReasonPoolClosed.
	reason := c.expireReason
	if c.expireReason == "" {
		reason = event.ReasonPoolClosed
	}

	_ = c.pool.removeConnection(c, reason)
	go func() {
		_ = c.pool.closeConnection(c)
	}()
}

// connectionInitFunc returns an init function for the resource pool that will make new connections for this pool
func (p *pool) connectionInitFunc() interface{} {
	c, _, err := p.makeNewConnection()
	if err != nil {
		return nil
	}

	go c.connect(context.Background())

	return c
}

// newPool creates a new pool that will hold size number of idle connections. It will use the
// provided options when creating connections.
func newPool(config poolConfig, connOpts ...ConnectionOption) (*pool, error) {
	opts := connOpts
	if config.MaxIdleTime != time.Duration(0) {
		opts = append(opts, WithIdleTimeout(func(_ time.Duration) time.Duration { return config.MaxIdleTime }))
	}

	var maxConns = config.MaxPoolSize
	if maxConns == 0 {
		maxConns = math.MaxInt64
	}

	pool := &pool{
		address:   config.Address,
		monitor:   config.PoolMonitor,
		connected: disconnected,
		opened:    make(map[uint64]*connection),
		opts:      opts,
		sem:       semaphore.NewWeighted(int64(maxConns)),
	}

	// we do not pass in config.MaxPoolSize because we manage the max size at this level rather than the resource pool level
	rpc := resourcePoolConfig{
		MaxSize:          maxConns,
		MinSize:          config.MinPoolSize,
		MaintainInterval: maintainInterval,
		ExpiredFn:        connectionExpiredFunc,
		CloseFn:          connectionCloseFunc,
		InitFn:           pool.connectionInitFunc,
	}

	if pool.monitor != nil {
		pool.monitor.Event(&event.PoolEvent{
			Type: event.PoolCreated,
			PoolOptions: &event.MonitorPoolOptions{
				MaxPoolSize:        rpc.MaxSize,
				MinPoolSize:        rpc.MinSize,
				WaitQueueTimeoutMS: uint64(config.MaxIdleTime) / uint64(time.Millisecond),
			},
			Address: pool.address.String(),
		})
	}

	rp, err := newResourcePool(rpc)
	if err != nil {
		return nil, err
	}
	pool.conns = rp

	return pool, nil
}

// stale checks if a given connection's generation is below the generation of the pool
func (p *pool) stale(c *connection) bool {
	return c == nil || c.generation < atomic.LoadUint64(&p.generation)
}

// connect puts the pool into the connected state, allowing it to be used and will allow items to begin being processed from the wait queue
func (p *pool) connect() error {
	if !atomic.CompareAndSwapInt32(&p.connected, disconnected, connected) {
		return ErrPoolConnected
	}
	p.conns.initialize()
	return nil
}

// disconnect disconnects the pool and closes all connections including those both in and out of the pool
func (p *pool) disconnect(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&p.connected, connected, disconnecting) {
		return ErrPoolDisconnected
	}

	if ctx == nil {
		ctx = context.Background()
	}

	p.conns.Close()
	atomic.AddUint64(&p.generation, 1)

	var err error
	if dl, ok := ctx.Deadline(); ok {
		// If we have a deadline then we interpret it as a request to gracefully shutdown. We wait
		// until either all the connections have landed back in the pool (and have been closed) or
		// until the timer is done.
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		timer := time.NewTimer(time.Now().Sub(dl))
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
			case <-ctx.Done():
			case <-ticker.C: // Can we replace this with an actual signal channel? We will know when p.inflight hits zero from the close method.
				p.Lock()
				if len(p.opened) > 0 {
					p.Unlock()
					continue
				}
				p.Unlock()
			}
			break
		}
	}

	// We copy the remaining connections into a slice, then iterate it to close them. This allows us
	// to use a single function to actually clean up and close connections at the expense of a
	// double iteration in the worse case.
	p.Lock()
	toClose := make([]*connection, 0, len(p.opened))
	for _, pc := range p.opened {
		toClose = append(toClose, pc)
	}
	p.Unlock()
	for _, pc := range toClose {
		_ = p.removeConnection(pc, event.ReasonPoolClosed)
		_ = p.closeConnection(pc) // We don't care about errors while closing the connection.
	}
	atomic.StoreInt32(&p.connected, disconnected)
	p.conns.clearTotal()

	if p.monitor != nil {
		p.monitor.Event(&event.PoolEvent{
			Type:    event.PoolClosedEvent,
			Address: p.address.String(),
		})
	}

	return err
}

// makeNewConnection creates a new connection instance and emits a ConnectionCreatedEvent. The caller must call
// connection.connect on the returned instance before using it for operations. This function ensures that a
// ConnectionClosed event is published if there is an error after the ConnectionCreated event has been published. The
// caller must not hold the pool lock when calling this function.
func (p *pool) makeNewConnection() (*connection, string, error) {
	c, err := newConnection(p.address, p.opts...)
	if err != nil {
		return nil, event.ReasonConnectionErrored, err
	}

	c.pool = p
	c.poolID = atomic.AddUint64(&p.nextid, 1)
	c.generation = atomic.LoadUint64(&p.generation)

	if p.monitor != nil {
		p.monitor.Event(&event.PoolEvent{
			Type:         event.ConnectionCreated,
			Address:      p.address.String(),
			ConnectionID: c.poolID,
		})
	}

	if atomic.LoadInt32(&p.connected) != connected {
		// Manually publish a ConnectionClosed event here because the connection reference hasn't been stored and we
		// need to ensure each ConnectionCreated event has a corresponding ConnectionClosed event.
		if p.monitor != nil {
			p.monitor.Event(&event.PoolEvent{
				Type:         event.ConnectionClosed,
				Address:      p.address.String(),
				ConnectionID: c.poolID,
				Reason:       event.ReasonPoolClosed,
			})
		}
		_ = p.closeConnection(c) // The pool is disconnected or disconnecting, ignore the error from closing the connection.
		return nil, event.ReasonPoolClosed, ErrPoolDisconnected
	}

	p.Lock()
	p.opened[c.poolID] = c
	p.Unlock()

	return c, "", nil

}

func (p *pool) getGeneration() uint64 {
	return atomic.LoadUint64(&p.generation)
}

// Checkout returns a connection from the pool
func (p *pool) get(ctx context.Context) (*connection, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if atomic.LoadInt32(&p.connected) != connected {
		if p.monitor != nil {
			p.monitor.Event(&event.PoolEvent{
				Type:    event.GetFailed,
				Address: p.address.String(),
				Reason:  event.ReasonPoolClosed,
			})
		}
		return nil, ErrPoolDisconnected
	}

	err := p.sem.Acquire(ctx, 1)
	if err != nil {
		if p.monitor != nil {
			p.monitor.Event(&event.PoolEvent{
				Type:    event.GetFailed,
				Address: p.address.String(),
				Reason:  event.ReasonTimedOut,
			})
		}
		return nil, ErrWaitQueueTimeout
	}

	// This loop is so that we don't end up with more than maxPoolSize connections if p.conns.Maintain runs between
	// calling p.conns.Get() and making the new connection
	for {
		if atomic.LoadInt32(&p.connected) != connected {
			if p.monitor != nil {
				p.monitor.Event(&event.PoolEvent{
					Type:    event.GetFailed,
					Address: p.address.String(),
					Reason:  event.ReasonPoolClosed,
				})
			}
			p.sem.Release(1)
			return nil, ErrPoolDisconnected
		}

		connVal := p.conns.Get()
		if c, ok := connVal.(*connection); ok && connVal != nil {
			// call connect if not connected
			if atomic.LoadInt32(&c.connected) == initialized {
				c.connect(ctx)
			}

			err := c.wait()
			if err != nil {
				// Call removeConnection to remove the connection reference and emit a ConnectionClosed event.
				_ = p.removeConnection(c, event.ReasonConnectionErrored)
				p.conns.decrementTotal()
				p.sem.Release(1)

				if p.monitor != nil {
					p.monitor.Event(&event.PoolEvent{
						Type:    event.GetFailed,
						Address: p.address.String(),
						Reason:  event.ReasonConnectionErrored,
					})
				}
				return nil, err
			}

			if p.monitor != nil {
				p.monitor.Event(&event.PoolEvent{
					Type:         event.GetSucceeded,
					Address:      p.address.String(),
					ConnectionID: c.poolID,
				})
			}
			return c, nil
		}

		select {
		case <-ctx.Done():
			if p.monitor != nil {
				p.monitor.Event(&event.PoolEvent{
					Type:    event.GetFailed,
					Address: p.address.String(),
					Reason:  event.ReasonTimedOut,
				})
			}
			p.sem.Release(1)
			return nil, ctx.Err()
		default:
			// The pool is empty, so we try to make a new connection. If incrementTotal fails, the resource pool has
			// more resources than we previously thought, so we try to get a resource again.
			made := p.conns.incrementTotal()
			if !made {
				continue
			}
			c, reason, err := p.makeNewConnection()

			if err != nil {
				if p.monitor != nil {
					// We only publish a GetFailed event because makeNewConnection has already published
					// ConnectionClosed if needed.
					p.monitor.Event(&event.PoolEvent{
						Type:    event.GetFailed,
						Address: p.address.String(),
						Reason:  reason,
					})
				}
				p.conns.decrementTotal()
				p.sem.Release(1)
				return nil, err
			}

			c.connect(ctx)
			// wait for conn to be connected
			err = c.wait()
			if err != nil {
				// Call removeConnection to remove the connection reference and fire a ConnectionClosedEvent.
				_ = p.removeConnection(c, event.ReasonConnectionErrored)
				p.conns.decrementTotal()
				p.sem.Release(1)

				if p.monitor != nil {
					p.monitor.Event(&event.PoolEvent{
						Type:    event.GetFailed,
						Address: p.address.String(),
						Reason:  reason,
					})
				}
				return nil, err
			}

			if p.monitor != nil {
				p.monitor.Event(&event.PoolEvent{
					Type:         event.GetSucceeded,
					Address:      p.address.String(),
					ConnectionID: c.poolID,
				})
			}
			return c, nil
		}
	}
}

// closeConnection closes a connection, not the pool itself. This method will actually closeConnection the connection,
// making it unusable, to instead return the connection to the pool, use put.
func (p *pool) closeConnection(c *connection) error {
	if c.pool != p {
		return ErrWrongPool
	}

	if atomic.LoadInt32(&c.connected) == connected {
		c.closeConnectContext()
		_ = c.wait() // Make sure that the connection has finished connecting
	}

	if !atomic.CompareAndSwapInt32(&c.connected, connected, disconnected) {
		return nil // We're closing an already closed connection
	}

	if c.nc != nil {
		err := c.nc.Close()
		if err != nil {
			return ConnectionError{ConnectionID: c.id, Wrapped: err, message: "failed to close net.Conn"}
		}
	}

	return nil
}

// removeConnection removes a connection from the pool.
func (p *pool) removeConnection(c *connection, reason string) error {
	if c.pool != p {
		return ErrWrongPool
	}

	var publishEvent bool
	p.Lock()
	if _, ok := p.opened[c.poolID]; ok {
		publishEvent = true
		delete(p.opened, c.poolID)
	}
	p.Unlock()

	if publishEvent && p.monitor != nil {
		c.pool.monitor.Event(&event.PoolEvent{
			Type:         event.ConnectionClosed,
			Address:      c.pool.address.String(),
			ConnectionID: c.poolID,
			Reason:       reason,
		})
	}
	return nil
}

// put returns a connection to this pool. If the pool is connected, the connection is not
// stale, and there is space in the cache, the connection is returned to the cache. This
// assumes that the connection has already been counted in p.conns.totalSize.
func (p *pool) put(c *connection) error {
	defer p.sem.Release(1)
	if p.monitor != nil {
		var cid uint64
		var addr string
		if c != nil {
			cid = c.poolID
			addr = c.addr.String()
		}
		p.monitor.Event(&event.PoolEvent{
			Type:         event.ConnectionReturned,
			ConnectionID: cid,
			Address:      addr,
		})
	}

	if c == nil {
		return nil
	}

	if c.pool != p {
		return ErrWrongPool
	}

	_ = p.conns.Put(c)

	return nil
}

// clear clears the pool by incrementing the generation
func (p *pool) clear() {
	if p.monitor != nil {
		p.monitor.Event(&event.PoolEvent{
			Type:    event.PoolCleared,
			Address: p.address.String(),
		})
	}
	atomic.AddUint64(&p.generation, 1)
}
