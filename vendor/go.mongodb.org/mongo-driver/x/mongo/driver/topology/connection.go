// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/ocsp"
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

var globalConnectionID uint64 = 1

var (
	defaultMaxMessageSize        uint32 = 48000000
	errResponseTooLarge          error  = errors.New("length of read message too large")
	errLoadBalancedStateMismatch        = errors.New("driver attempted to initialize in load balancing mode, but the server does not support this mode")
)

func nextConnectionID() uint64 { return atomic.AddUint64(&globalConnectionID, 1) }

type connection struct {
	// connected must be accessed using the atomic package and should be at the beginning of the struct.
	// - atomic bug: https://pkg.go.dev/sync/atomic#pkg-note-BUG
	// - suggested layout: https://go101.org/article/memory-layout.html
	connected int64

	id                   string
	nc                   net.Conn // When nil, the connection is closed.
	addr                 address.Address
	idleTimeout          time.Duration
	idleDeadline         atomic.Value // Stores a time.Time
	readTimeout          time.Duration
	writeTimeout         time.Duration
	descMu               sync.RWMutex // Guards desc. TODO: Remove with or after GODRIVER-2038.
	desc                 description.Server
	isMasterRTT          time.Duration
	compressor           wiremessage.CompressorID
	zliblevel            int
	zstdLevel            int
	connectDone          chan struct{}
	connectErr           error
	config               *connectionConfig
	cancelConnectContext context.CancelFunc
	connectContextMade   chan struct{}
	canStream            bool
	currentlyStreaming   bool
	connectContextMutex  sync.Mutex
	cancellationListener cancellationListener

	// pool related fields
	pool         *pool
	poolID       uint64
	generation   uint64
	expireReason string
	poolMonitor  *event.PoolMonitor
}

// newConnection handles the creation of a connection. It does not connect the connection.
func newConnection(addr address.Address, opts ...ConnectionOption) (*connection, error) {
	cfg, err := newConnectionConfig(opts...)
	if err != nil {
		return nil, err
	}

	id := fmt.Sprintf("%s[-%d]", addr, nextConnectionID())

	c := &connection{
		id:                   id,
		addr:                 addr,
		idleTimeout:          cfg.idleTimeout,
		readTimeout:          cfg.readTimeout,
		writeTimeout:         cfg.writeTimeout,
		connectDone:          make(chan struct{}),
		config:               cfg,
		connectContextMade:   make(chan struct{}),
		cancellationListener: internal.NewCancellationListener(),
		poolMonitor:          cfg.poolMonitor,
	}
	// Connections to non-load balanced deployments should eagerly set the generation numbers so errors encountered
	// at any point during connection establishment can be processed without the connection being considered stale.
	if !c.config.loadBalanced {
		c.setGenerationNumber()
	}
	atomic.StoreInt64(&c.connected, initialized)

	return c, nil
}

func (c *connection) processInitializationError(opCtx context.Context, err error) {
	atomic.StoreInt64(&c.connected, disconnected)
	if c.nc != nil {
		_ = c.nc.Close()
	}

	c.connectErr = ConnectionError{Wrapped: err, init: true}
	if c.config.errorHandlingCallback != nil {
		c.config.errorHandlingCallback(opCtx, c.connectErr, c.generation, c.desc.ServiceID)
	}
}

// setGenerationNumber sets the connection's generation number if a callback has been provided to do so in connection
// configuration.
func (c *connection) setGenerationNumber() {
	if c.config.getGenerationFn != nil {
		c.generation = c.config.getGenerationFn(c.desc.ServiceID)
	}
}

// hasGenerationNumber returns true if the connection has set its generation number. If so, this indicates that the
// generationNumberFn provided via the connection options has been called exactly once.
func (c *connection) hasGenerationNumber() bool {
	if !c.config.loadBalanced {
		// The generation is known for all non-LB clusters once the connection object has been created.
		return true
	}

	// For LB clusters, we set the generation after the initial handshake, so we know it's set if the connection
	// description has been updated to reflect that it's behind an LB.
	return c.desc.LoadBalanced()
}

// connect handles the I/O for a connection. It will dial, configure TLS, and perform
// initialization handshakes.
func (c *connection) connect(ctx context.Context) {
	if !atomic.CompareAndSwapInt64(&c.connected, initialized, connected) {
		return
	}
	defer close(c.connectDone)

	// Create separate contexts for dialing a connection and doing the MongoDB/auth handshakes.
	//
	// handshakeCtx is simply a cancellable version of ctx because there's no default timeout that needs to be applied
	// to the full handshake. The cancellation allows consumers to bail out early when dialing a connection if it's no
	// longer required. This is done in lock because it accesses the shared cancelConnectContext field.
	//
	// dialCtx is equal to handshakeCtx if connectTimeoutMS=0. Otherwise, it is derived from handshakeCtx so the
	// cancellation still applies but with an added timeout to ensure the connectTimeoutMS option is applied to socket
	// establishment and the TLS handshake as a whole. This is created outside of the connectContextMutex lock to avoid
	// holding the lock longer than necessary.
	c.connectContextMutex.Lock()
	var handshakeCtx context.Context
	handshakeCtx, c.cancelConnectContext = context.WithCancel(ctx)
	c.connectContextMutex.Unlock()

	dialCtx := handshakeCtx
	var dialCancel context.CancelFunc
	if c.config.connectTimeout != 0 {
		dialCtx, dialCancel = context.WithTimeout(handshakeCtx, c.config.connectTimeout)
		defer dialCancel()
	}

	defer func() {
		var cancelFn context.CancelFunc

		c.connectContextMutex.Lock()
		cancelFn = c.cancelConnectContext
		c.cancelConnectContext = nil
		c.connectContextMutex.Unlock()

		if cancelFn != nil {
			cancelFn()
		}
	}()

	close(c.connectContextMade)

	// Assign the result of DialContext to a temporary net.Conn to ensure that c.nc is not set in an error case.
	var err error
	var tempNc net.Conn
	tempNc, err = c.config.dialer.DialContext(dialCtx, c.addr.Network(), c.addr.String())
	if err != nil {
		c.processInitializationError(ctx, err)
		return
	}
	c.nc = tempNc

	if c.config.tlsConfig != nil {
		tlsConfig := c.config.tlsConfig.Clone()

		// store the result of configureTLS in a separate variable than c.nc to avoid overwriting c.nc with nil in
		// error cases.
		ocspOpts := &ocsp.VerifyOptions{
			Cache:                   c.config.ocspCache,
			DisableEndpointChecking: c.config.disableOCSPEndpointCheck,
		}
		tlsNc, err := configureTLS(dialCtx, c.config.tlsConnectionSource, c.nc, c.addr, tlsConfig, ocspOpts)
		if err != nil {
			c.processInitializationError(ctx, err)
			return
		}
		c.nc = tlsNc
	}

	c.bumpIdleDeadline()

	// running isMaster and authentication is handled by a handshaker on the configuration instance.
	handshaker := c.config.handshaker
	if handshaker == nil {
		if c.poolMonitor != nil {
			c.poolMonitor.Event(&event.PoolEvent{
				Type:         event.ConnectionReady,
				Address:      c.addr.String(),
				ConnectionID: c.poolID,
			})
		}
		return
	}

	var handshakeInfo driver.HandshakeInformation
	handshakeStartTime := time.Now()
	handshakeConn := initConnection{c}
	handshakeInfo, err = handshaker.GetHandshakeInformation(handshakeCtx, c.addr, handshakeConn)
	if err == nil {
		// We only need to retain the Description field as the connection's description. The authentication-related
		// fields in handshakeInfo are tracked by the handshaker if necessary.
		c.descMu.Lock()
		c.desc = handshakeInfo.Description
		c.descMu.Unlock()
		c.isMasterRTT = time.Since(handshakeStartTime)

		// If the application has indicated that the cluster is load balanced, ensure the server has included serviceId
		// in its handshake response to signal that it knows it's behind an LB as well.
		if c.config.loadBalanced && c.desc.ServiceID == nil {
			err = errLoadBalancedStateMismatch
		}
	}
	if err == nil {
		// For load-balanced connections, the generation number depends on the service ID, which isn't known until the
		// initial MongoDB handshake is done. To account for this, we don't attempt to set the connection's generation
		// number unless GetHandshakeInformation succeeds.
		if c.config.loadBalanced {
			c.setGenerationNumber()
		}

		// If we successfully finished the first part of the handshake and verified LB state, continue with the rest of
		// the handshake.
		err = handshaker.FinishHandshake(handshakeCtx, handshakeConn)
	}

	// We have a failed handshake here
	if err != nil {
		c.processInitializationError(ctx, err)
		return
	}

	if len(c.desc.Compression) > 0 {
	clientMethodLoop:
		for _, method := range c.config.compressors {
			for _, serverMethod := range c.desc.Compression {
				if method != serverMethod {
					continue
				}

				switch strings.ToLower(method) {
				case "snappy":
					c.compressor = wiremessage.CompressorSnappy
				case "zlib":
					c.compressor = wiremessage.CompressorZLib
					c.zliblevel = wiremessage.DefaultZlibLevel
					if c.config.zlibLevel != nil {
						c.zliblevel = *c.config.zlibLevel
					}
				case "zstd":
					c.compressor = wiremessage.CompressorZstd
					c.zstdLevel = wiremessage.DefaultZstdLevel
					if c.config.zstdLevel != nil {
						c.zstdLevel = *c.config.zstdLevel
					}
				}
				break clientMethodLoop
			}
		}
	}
	if c.poolMonitor != nil {
		c.poolMonitor.Event(&event.PoolEvent{
			Type:         event.ConnectionReady,
			Address:      c.addr.String(),
			ConnectionID: c.poolID,
		})
	}
}

func (c *connection) wait() error {
	if c.connectDone != nil {
		<-c.connectDone
	}
	return c.connectErr
}

func (c *connection) closeConnectContext() {
	<-c.connectContextMade
	var cancelFn context.CancelFunc

	c.connectContextMutex.Lock()
	cancelFn = c.cancelConnectContext
	c.cancelConnectContext = nil
	c.connectContextMutex.Unlock()

	if cancelFn != nil {
		cancelFn()
	}
}

func transformNetworkError(ctx context.Context, originalError error, contextDeadlineUsed bool) error {
	if originalError == nil {
		return nil
	}

	// If there was an error and the context was cancelled, we assume it happened due to the cancellation.
	if ctx.Err() == context.Canceled {
		return context.Canceled
	}

	// If there was a timeout error and the context deadline was used, we convert the error into
	// context.DeadlineExceeded.
	if !contextDeadlineUsed {
		return originalError
	}
	if netErr, ok := originalError.(net.Error); ok && netErr.Timeout() {
		return context.DeadlineExceeded
	}

	return originalError
}

func (c *connection) cancellationListenerCallback() {
	_ = c.close()
}

func (c *connection) writeWireMessage(ctx context.Context, wm []byte) error {
	var err error
	if atomic.LoadInt64(&c.connected) != connected {
		return ConnectionError{ConnectionID: c.id, message: "connection is closed"}
	}
	select {
	case <-ctx.Done():
		return ConnectionError{ConnectionID: c.id, Wrapped: ctx.Err(), message: "failed to write"}
	default:
	}

	var deadline time.Time
	if c.writeTimeout != 0 {
		deadline = time.Now().Add(c.writeTimeout)
	}

	var contextDeadlineUsed bool
	if dl, ok := ctx.Deadline(); ok && (deadline.IsZero() || dl.Before(deadline)) {
		contextDeadlineUsed = true
		deadline = dl
	}

	if err := c.nc.SetWriteDeadline(deadline); err != nil {
		return ConnectionError{ConnectionID: c.id, Wrapped: err, message: "failed to set write deadline"}
	}

	err = c.write(ctx, wm)
	if err != nil {
		c.close()
		return ConnectionError{
			ConnectionID: c.id,
			Wrapped:      transformNetworkError(ctx, err, contextDeadlineUsed),
			message:      "unable to write wire message to network",
		}
	}

	c.bumpIdleDeadline()
	return nil
}

func (c *connection) write(ctx context.Context, wm []byte) (err error) {
	go c.cancellationListener.Listen(ctx, c.cancellationListenerCallback)
	defer func() {
		// There is a race condition between Write and StopListening. If the context is cancelled after c.nc.Write
		// succeeds, the cancellation listener could fire and close the connection. In this case, the connection has
		// been invalidated but the error is nil. To account for this, overwrite the error to context.Cancelled if
		// the abortedForCancellation flag was set.

		if aborted := c.cancellationListener.StopListening(); aborted && err == nil {
			err = context.Canceled
		}
	}()

	_, err = c.nc.Write(wm)
	return err
}

// readWireMessage reads a wiremessage from the connection. The dst parameter will be overwritten.
func (c *connection) readWireMessage(ctx context.Context, dst []byte) ([]byte, error) {
	if atomic.LoadInt64(&c.connected) != connected {
		return dst, ConnectionError{ConnectionID: c.id, message: "connection is closed"}
	}

	select {
	case <-ctx.Done():
		// We closeConnection the connection because we don't know if there is an unread message on the wire.
		c.close()
		return nil, ConnectionError{ConnectionID: c.id, Wrapped: ctx.Err(), message: "failed to read"}
	default:
	}

	var deadline time.Time
	if c.readTimeout != 0 {
		deadline = time.Now().Add(c.readTimeout)
	}

	var contextDeadlineUsed bool
	if dl, ok := ctx.Deadline(); ok && (deadline.IsZero() || dl.Before(deadline)) {
		contextDeadlineUsed = true
		deadline = dl
	}

	if err := c.nc.SetReadDeadline(deadline); err != nil {
		return nil, ConnectionError{ConnectionID: c.id, Wrapped: err, message: "failed to set read deadline"}
	}

	dst, errMsg, err := c.read(ctx, dst)
	if err != nil {
		// We closeConnection the connection because we don't know if there are other bytes left to read.
		c.close()
		message := errMsg
		if err == io.EOF {
			message = "socket was unexpectedly closed"
		}
		return nil, ConnectionError{
			ConnectionID: c.id,
			Wrapped:      transformNetworkError(ctx, err, contextDeadlineUsed),
			message:      message,
		}
	}

	c.bumpIdleDeadline()
	return dst, nil
}

func (c *connection) read(ctx context.Context, dst []byte) (bytesRead []byte, errMsg string, err error) {
	go c.cancellationListener.Listen(ctx, c.cancellationListenerCallback)
	defer func() {
		// If the context is cancelled after we finish reading the server response, the cancellation listener could fire
		// even though the socket reads succeed. To account for this, we overwrite err to be context.Canceled if the
		// abortedForCancellation flag is set.

		if aborted := c.cancellationListener.StopListening(); aborted && err == nil {
			errMsg = "unable to read server response"
			err = context.Canceled
		}
	}()

	// We use an array here because it only costs 4 bytes on the stack and means we'll only need to
	// reslice dst once instead of twice.
	var sizeBuf [4]byte

	// We do a ReadFull into an array here instead of doing an opportunistic ReadAtLeast into dst
	// because there might be more than one wire message waiting to be read, for example when
	// reading messages from an exhaust cursor.
	_, err = io.ReadFull(c.nc, sizeBuf[:])
	if err != nil {
		return nil, "incomplete read of message header", err
	}

	// read the length as an int32
	size := (int32(sizeBuf[0])) | (int32(sizeBuf[1]) << 8) | (int32(sizeBuf[2]) << 16) | (int32(sizeBuf[3]) << 24)

	// In the case of an isMaster response where MaxMessageSize has not yet been set, use the hard-coded
	// defaultMaxMessageSize instead.
	maxMessageSize := c.desc.MaxMessageSize
	if maxMessageSize == 0 {
		maxMessageSize = defaultMaxMessageSize
	}
	if uint32(size) > maxMessageSize {
		return nil, errResponseTooLarge.Error(), errResponseTooLarge
	}

	if int(size) > cap(dst) {
		// Since we can't grow this slice without allocating, just allocate an entirely new slice.
		dst = make([]byte, 0, size)
	}
	// We need to ensure we don't accidentally read into a subsequent wire message, so we set the
	// size to read exactly this wire message.
	dst = dst[:size]
	copy(dst, sizeBuf[:])

	_, err = io.ReadFull(c.nc, dst[4:])
	if err != nil {
		return nil, "incomplete read of full message", err
	}

	return dst, "", nil
}

func (c *connection) close() error {
	// Overwrite the connection state as the first step so only the first close call will execute.
	if !atomic.CompareAndSwapInt64(&c.connected, connected, disconnected) {
		return nil
	}

	var err error
	if c.nc != nil {
		err = c.nc.Close()
	}

	return err
}

func (c *connection) closed() bool {
	return atomic.LoadInt64(&c.connected) == disconnected
}

func (c *connection) idleTimeoutExpired() bool {
	now := time.Now()
	if c.idleTimeout > 0 {
		idleDeadline, ok := c.idleDeadline.Load().(time.Time)
		if ok && now.After(idleDeadline) {
			return true
		}
	}

	return false
}

func (c *connection) bumpIdleDeadline() {
	if c.idleTimeout > 0 {
		c.idleDeadline.Store(time.Now().Add(c.idleTimeout))
	}
}

func (c *connection) setCanStream(canStream bool) {
	c.canStream = canStream
}

func (c initConnection) supportsStreaming() bool {
	return c.canStream
}

func (c *connection) setStreaming(streaming bool) {
	c.currentlyStreaming = streaming
}

func (c *connection) getCurrentlyStreaming() bool {
	return c.currentlyStreaming
}

func (c *connection) setSocketTimeout(timeout time.Duration) {
	c.readTimeout = timeout
	c.writeTimeout = timeout
}

func (c *connection) ID() string {
	return c.id
}

// initConnection is an adapter used during connection initialization. It has the minimum
// functionality necessary to implement the driver.Connection interface, which is required to pass a
// *connection to a Handshaker.
type initConnection struct{ *connection }

var _ driver.Connection = initConnection{}
var _ driver.StreamerConnection = initConnection{}

func (c initConnection) Description() description.Server {
	if c.connection == nil {
		return description.Server{}
	}
	return c.connection.desc
}
func (c initConnection) Close() error             { return nil }
func (c initConnection) ID() string               { return c.id }
func (c initConnection) Address() address.Address { return c.addr }
func (c initConnection) Stale() bool              { return false }
func (c initConnection) LocalAddress() address.Address {
	if c.connection == nil || c.nc == nil {
		return address.Address("0.0.0.0")
	}
	return address.Address(c.nc.LocalAddr().String())
}
func (c initConnection) WriteWireMessage(ctx context.Context, wm []byte) error {
	return c.writeWireMessage(ctx, wm)
}
func (c initConnection) ReadWireMessage(ctx context.Context, dst []byte) ([]byte, error) {
	return c.readWireMessage(ctx, dst)
}
func (c initConnection) SetStreaming(streaming bool) {
	c.setStreaming(streaming)
}
func (c initConnection) CurrentlyStreaming() bool {
	return c.getCurrentlyStreaming()
}
func (c initConnection) SupportsStreaming() bool {
	return c.supportsStreaming()
}

// Connection implements the driver.Connection interface to allow reading and writing wire
// messages and the driver.Expirable interface to allow expiring.
type Connection struct {
	*connection
	refCount      int
	cleanupPoolFn func()

	mu sync.RWMutex
}

var _ driver.Connection = (*Connection)(nil)
var _ driver.Expirable = (*Connection)(nil)
var _ driver.PinnedConnection = (*Connection)(nil)

// WriteWireMessage handles writing a wire message to the underlying connection.
func (c *Connection) WriteWireMessage(ctx context.Context, wm []byte) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.connection == nil {
		return ErrConnectionClosed
	}
	return c.writeWireMessage(ctx, wm)
}

// ReadWireMessage handles reading a wire message from the underlying connection. The dst parameter
// will be overwritten with the new wire message.
func (c *Connection) ReadWireMessage(ctx context.Context, dst []byte) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.connection == nil {
		return dst, ErrConnectionClosed
	}
	return c.readWireMessage(ctx, dst)
}

// CompressWireMessage handles compressing the provided wire message using the underlying
// connection's compressor. The dst parameter will be overwritten with the new wire message. If
// there is no compressor set on the underlying connection, then no compression will be performed.
func (c *Connection) CompressWireMessage(src, dst []byte) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.connection == nil {
		return dst, ErrConnectionClosed
	}
	if c.connection.compressor == wiremessage.CompressorNoOp {
		return append(dst, src...), nil
	}
	_, reqid, respto, origcode, rem, ok := wiremessage.ReadHeader(src)
	if !ok {
		return dst, errors.New("wiremessage is too short to compress, less than 16 bytes")
	}
	idx, dst := wiremessage.AppendHeaderStart(dst, reqid, respto, wiremessage.OpCompressed)
	dst = wiremessage.AppendCompressedOriginalOpCode(dst, origcode)
	dst = wiremessage.AppendCompressedUncompressedSize(dst, int32(len(rem)))
	dst = wiremessage.AppendCompressedCompressorID(dst, c.connection.compressor)
	opts := driver.CompressionOpts{
		Compressor: c.connection.compressor,
		ZlibLevel:  c.connection.zliblevel,
		ZstdLevel:  c.connection.zstdLevel,
	}
	compressed, err := driver.CompressPayload(rem, opts)
	if err != nil {
		return nil, err
	}
	dst = wiremessage.AppendCompressedCompressedMessage(dst, compressed)
	return bsoncore.UpdateLength(dst, idx, int32(len(dst[idx:]))), nil
}

// Description returns the server description of the server this connection is connected to.
func (c *Connection) Description() description.Server {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.connection == nil {
		return description.Server{}
	}
	return c.desc
}

// Close returns this connection to the connection pool. This method may not closeConnection the underlying
// socket.
func (c *Connection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.connection == nil || c.refCount > 0 {
		return nil
	}

	return c.cleanupReferences()
}

// Expire closes this connection and will closeConnection the underlying socket.
func (c *Connection) Expire() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.connection == nil {
		return nil
	}

	_ = c.close()
	return c.cleanupReferences()
}

func (c *Connection) cleanupReferences() error {
	err := c.pool.put(c.connection)
	if c.cleanupPoolFn != nil {
		c.cleanupPoolFn()
		c.cleanupPoolFn = nil
	}
	c.connection = nil
	return err
}

// Alive returns if the connection is still alive.
func (c *Connection) Alive() bool {
	return c.connection != nil
}

// ID returns the ID of this connection.
func (c *Connection) ID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.connection == nil {
		return "<closed>"
	}
	return c.id
}

// Stale returns if the connection is stale.
func (c *Connection) Stale() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.pool.stale(c.connection)
}

// Address returns the address of this connection.
func (c *Connection) Address() address.Address {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.connection == nil {
		return address.Address("0.0.0.0")
	}
	return c.addr
}

// LocalAddress returns the local address of the connection
func (c *Connection) LocalAddress() address.Address {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.connection == nil || c.nc == nil {
		return address.Address("0.0.0.0")
	}
	return address.Address(c.nc.LocalAddr().String())
}

// PinToCursor updates this connection to reflect that it is pinned to a cursor.
func (c *Connection) PinToCursor() error {
	return c.pin("cursor", c.pool.pinConnectionToCursor, c.pool.unpinConnectionFromCursor)
}

// PinToTransaction updates this connection to reflect that it is pinned to a transaction.
func (c *Connection) PinToTransaction() error {
	return c.pin("transaction", c.pool.pinConnectionToTransaction, c.pool.unpinConnectionFromTransaction)
}

func (c *Connection) pin(reason string, updatePoolFn, cleanupPoolFn func()) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.connection == nil {
		return fmt.Errorf("attempted to pin a connection for a %s, but the connection has already been returned to the pool", reason)
	}

	// Only use the provided callbacks for the first reference to avoid double-counting pinned connection statistics
	// in the pool.
	if c.refCount == 0 {
		updatePoolFn()
		c.cleanupPoolFn = cleanupPoolFn
	}
	c.refCount++
	return nil
}

// UnpinFromCursor updates this connection to reflect that it is no longer pinned to a cursor.
func (c *Connection) UnpinFromCursor() error {
	return c.unpin("cursor")
}

// UnpinFromTransaction updates this connection to reflect that it is no longer pinned to a transaction.
func (c *Connection) UnpinFromTransaction() error {
	return c.unpin("transaction")
}

func (c *Connection) unpin(reason string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.connection == nil {
		// We don't error here because the resource could have been forcefully closed via Expire.
		return nil
	}
	if c.refCount == 0 {
		return fmt.Errorf("attempted to unpin a connection from a %s, but the connection is not pinned by any resources", reason)
	}

	c.refCount--
	return nil
}

var notMasterCodes = []int32{10107, 13435}
var recoveringCodes = []int32{11600, 11602, 13436, 189, 91}

func configureTLS(ctx context.Context,
	tlsConnSource tlsConnectionSource,
	nc net.Conn,
	addr address.Address,
	config *tls.Config,
	ocspOpts *ocsp.VerifyOptions,
) (net.Conn, error) {

	// Ensure config.ServerName is always set for SNI.
	if config.ServerName == "" {
		hostname := addr.String()
		colonPos := strings.LastIndex(hostname, ":")
		if colonPos == -1 {
			colonPos = len(hostname)
		}

		hostname = hostname[:colonPos]
		config.ServerName = hostname
	}

	client := tlsConnSource.Client(nc, config)
	errChan := make(chan error, 1)
	go func() {
		errChan <- client.Handshake()
	}()

	select {
	case err := <-errChan:
		if err != nil {
			return nil, err
		}

		// Only do OCSP verification if TLS verification is requested.
		if config.InsecureSkipVerify {
			break
		}

		if ocspErr := ocsp.Verify(ctx, client.ConnectionState(), ocspOpts); ocspErr != nil {
			return nil, ocspErr
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	return client, nil
}
