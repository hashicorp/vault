// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/internal/csot"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/ocsp"
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

// Connection state constants.
const (
	connDisconnected int64 = iota
	connConnected
	connInitialized
)

var globalConnectionID uint64 = 1

var (
	defaultMaxMessageSize        uint32 = 48000000
	errResponseTooLarge                 = errors.New("length of read message too large")
	errLoadBalancedStateMismatch        = errors.New("driver attempted to initialize in load balancing mode, but the server does not support this mode")
)

func nextConnectionID() uint64 { return atomic.AddUint64(&globalConnectionID, 1) }

type connection struct {
	// state must be accessed using the atomic package and should be at the beginning of the struct.
	// - atomic bug: https://pkg.go.dev/sync/atomic#pkg-note-BUG
	// - suggested layout: https://go101.org/article/memory-layout.html
	state int64

	id                   string
	nc                   net.Conn // When nil, the connection is closed.
	addr                 address.Address
	idleTimeout          time.Duration
	idleStart            atomic.Value // Stores a time.Time
	readTimeout          time.Duration
	writeTimeout         time.Duration
	desc                 description.Server
	helloRTT             time.Duration
	compressor           wiremessage.CompressorID
	zliblevel            int
	zstdLevel            int
	connectDone          chan struct{}
	config               *connectionConfig
	cancelConnectContext context.CancelFunc
	connectContextMade   chan struct{}
	canStream            bool
	currentlyStreaming   bool
	connectContextMutex  sync.Mutex
	cancellationListener cancellationListener
	serverConnectionID   *int64 // the server's ID for this client's connection

	// pool related fields
	pool *pool

	// TODO(GODRIVER-2824): change driverConnectionID type to int64.
	driverConnectionID uint64
	generation         uint64

	// awaitRemainingBytes indicates the size of server response that was not completely
	// read before returning the connection to the pool.
	awaitRemainingBytes *int32

	// oidcTokenGenID is the monotonic generation ID for OIDC tokens, used to invalidate
	// accessTokens in the OIDC authenticator cache.
	oidcTokenGenID uint64
}

// newConnection handles the creation of a connection. It does not connect the connection.
func newConnection(addr address.Address, opts ...ConnectionOption) *connection {
	cfg := newConnectionConfig(opts...)

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
		cancellationListener: newCancellListener(),
	}
	// Connections to non-load balanced deployments should eagerly set the generation numbers so errors encountered
	// at any point during connection establishment can be processed without the connection being considered stale.
	if !c.config.loadBalanced {
		c.setGenerationNumber()
	}
	atomic.StoreInt64(&c.state, connInitialized)

	return c
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
	if err := clientHandshake(ctx, client); err != nil {
		return nil, err
	}

	// Only do OCSP verification if TLS verification is requested.
	if !config.InsecureSkipVerify {
		if ocspErr := ocsp.Verify(ctx, client.ConnectionState(), ocspOpts); ocspErr != nil {
			return nil, ocspErr
		}
	}
	return client, nil
}

// connect handles the I/O for a connection. It will dial, configure TLS, and perform initialization
// handshakes. All errors returned by connect are considered "before the handshake completes" and
// must be handled by calling the appropriate SDAM handshake error handler.
func (c *connection) connect(ctx context.Context) (err error) {
	if !atomic.CompareAndSwapInt64(&c.state, connInitialized, connConnected) {
		return nil
	}

	defer close(c.connectDone)

	// If connect returns an error, set the connection status as disconnected and close the
	// underlying net.Conn if it was created.
	defer func() {
		if err != nil {
			atomic.StoreInt64(&c.state, connDisconnected)

			if c.nc != nil {
				_ = c.nc.Close()
			}
		}
	}()

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
	tempNc, err := c.config.dialer.DialContext(dialCtx, c.addr.Network(), c.addr.String())
	if err != nil {
		return ConnectionError{Wrapped: err, init: true}
	}
	c.nc = tempNc

	if c.config.tlsConfig != nil {
		tlsConfig := c.config.tlsConfig.Clone()

		// store the result of configureTLS in a separate variable than c.nc to avoid overwriting c.nc with nil in
		// error cases.
		ocspOpts := &ocsp.VerifyOptions{
			Cache:                   c.config.ocspCache,
			DisableEndpointChecking: c.config.disableOCSPEndpointCheck,
			HTTPClient:              c.config.httpClient,
		}
		tlsNc, err := configureTLS(dialCtx, c.config.tlsConnectionSource, c.nc, c.addr, tlsConfig, ocspOpts)
		if err != nil {
			return ConnectionError{Wrapped: err, init: true}
		}
		c.nc = tlsNc
	}

	// running hello and authentication is handled by a handshaker on the configuration instance.
	handshaker := c.config.handshaker
	if handshaker == nil {
		return nil
	}

	var handshakeInfo driver.HandshakeInformation
	handshakeStartTime := time.Now()
	handshakeConn := initConnection{c}
	handshakeInfo, err = handshaker.GetHandshakeInformation(handshakeCtx, c.addr, handshakeConn)
	if err == nil {
		// We only need to retain the Description field as the connection's description. The authentication-related
		// fields in handshakeInfo are tracked by the handshaker if necessary.
		c.desc = handshakeInfo.Description
		c.serverConnectionID = handshakeInfo.ServerConnectionID
		c.helloRTT = time.Since(handshakeStartTime)

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
		return ConnectionError{Wrapped: err, init: true}
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
	return nil
}

func (c *connection) wait() {
	if c.connectDone != nil {
		<-c.connectDone
	}
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

func (c *connection) cancellationListenerCallback() {
	_ = c.close()
}

func transformNetworkError(ctx context.Context, originalError error, contextDeadlineUsed bool) error {
	if originalError == nil {
		return nil
	}

	// If there was an error and the context was cancelled, we assume it happened due to the cancellation.
	if errors.Is(ctx.Err(), context.Canceled) {
		return ctx.Err()
	}

	// If there was a timeout error and the context deadline was used, we convert the error into
	// context.DeadlineExceeded.
	if !contextDeadlineUsed {
		return originalError
	}
	if netErr, ok := originalError.(net.Error); ok && netErr.Timeout() {
		return fmt.Errorf("%w: %s", context.DeadlineExceeded, originalError.Error())
	}

	return originalError
}

func (c *connection) writeWireMessage(ctx context.Context, wm []byte) error {
	var err error
	if atomic.LoadInt64(&c.state) != connConnected {
		return ConnectionError{
			ConnectionID: c.id,
			message:      "connection is closed",
		}
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
func (c *connection) readWireMessage(ctx context.Context) ([]byte, error) {
	if atomic.LoadInt64(&c.state) != connConnected {
		return nil, ConnectionError{
			ConnectionID: c.id,
			message:      "connection is closed",
		}
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

	dst, errMsg, err := c.read(ctx)
	if err != nil {
		if c.awaitRemainingBytes == nil {
			// If the connection was not marked as awaiting response, use the
			// pre-CSOT behavior and close the connection because we don't know
			// if there are other bytes left to read.
			c.close()
		}
		message := errMsg
		if errors.Is(err, io.EOF) {
			message = "socket was unexpectedly closed"
		}
		return nil, ConnectionError{
			ConnectionID: c.id,
			Wrapped:      transformNetworkError(ctx, err, contextDeadlineUsed),
			message:      message,
		}
	}

	return dst, nil
}

func (c *connection) parseWmSizeBytes(wmSizeBytes [4]byte) (int32, error) {
	// read the length as an int32
	size := int32(binary.LittleEndian.Uint32(wmSizeBytes[:]))

	if size < 4 {
		return 0, fmt.Errorf("malformed message length: %d", size)
	}
	// In the case of a hello response where MaxMessageSize has not yet been set, use the hard-coded
	// defaultMaxMessageSize instead.
	maxMessageSize := c.desc.MaxMessageSize
	if maxMessageSize == 0 {
		maxMessageSize = defaultMaxMessageSize
	}
	if uint32(size) > maxMessageSize {
		return 0, errResponseTooLarge
	}

	return size, nil
}

func (c *connection) read(ctx context.Context) (bytesRead []byte, errMsg string, err error) {
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

	isCSOTTimeout := func(err error) bool {
		// If the error was a timeout error and CSOT is enabled, instead of
		// closing the connection mark it as awaiting response so the pool
		// can read the response before making it available to other
		// operations.
		nerr := net.Error(nil)
		return errors.As(err, &nerr) && nerr.Timeout() && csot.IsTimeoutContext(ctx)
	}

	// We use an array here because it only costs 4 bytes on the stack and means we'll only need to
	// reslice dst once instead of twice.
	var sizeBuf [4]byte

	// We do a ReadFull into an array here instead of doing an opportunistic ReadAtLeast into dst
	// because there might be more than one wire message waiting to be read, for example when
	// reading messages from an exhaust cursor.
	n, err := io.ReadFull(c.nc, sizeBuf[:])
	if err != nil {
		if l := int32(n); l == 0 && isCSOTTimeout(err) {
			c.awaitRemainingBytes = &l
		}
		return nil, "incomplete read of message header", err
	}
	size, err := c.parseWmSizeBytes(sizeBuf)
	if err != nil {
		return nil, err.Error(), err
	}

	dst := make([]byte, size)
	copy(dst, sizeBuf[:])

	n, err = io.ReadFull(c.nc, dst[4:])
	if err != nil {
		remainingBytes := size - 4 - int32(n)
		if remainingBytes > 0 && isCSOTTimeout(err) {
			c.awaitRemainingBytes = &remainingBytes
		}
		return dst, "incomplete read of full message", err
	}

	return dst, "", nil
}

func (c *connection) close() error {
	// Overwrite the connection state as the first step so only the first close call will execute.
	if !atomic.CompareAndSwapInt64(&c.state, connConnected, connDisconnected) {
		return nil
	}

	var err error
	if c.nc != nil {
		err = c.nc.Close()
	}

	return err
}

// closed returns true if the connection has been closed by the driver.
func (c *connection) closed() bool {
	return atomic.LoadInt64(&c.state) == connDisconnected
}

// isAlive returns true if the connection is alive and ready to be used for an
// operation.
//
// Note that the liveness check can be slow (at least 1ms), so isAlive only
// checks the liveness of the connection if it's been idle for at least 10
// seconds. For frequently in-use connections, a network error during an
// operation will be the first indication of a dead connection.
func (c *connection) isAlive() bool {
	if c.nc == nil {
		return false
	}

	// If the connection has been idle for less than 10 seconds, skip the
	// liveness check.
	//
	// The 10-seconds idle bypass is based on the liveness check implementation
	// in the Python Driver. That implementation uses 1 second as the idle
	// threshold, but we chose to be more conservative in the Go Driver because
	// this is new behavior with unknown side-effects. See
	// https://github.com/mongodb/mongo-python-driver/blob/e6b95f65953e01e435004af069a6976473eaf841/pymongo/synchronous/pool.py#L983-L985
	idleStart, ok := c.idleStart.Load().(time.Time)
	if !ok || idleStart.Add(10*time.Second).After(time.Now()) {
		return true
	}

	// Set a 1ms read deadline and attempt to read 1 byte from the connection.
	// Expect it to block for 1ms then return a deadline exceeded error. If it
	// returns any other error, the connection is not usable, so return false.
	// If it doesn't return an error and actually reads data, the connection is
	// also not usable, so return false.
	//
	// Note that we don't need to un-set the read deadline because the "read"
	// and "write" methods always reset the deadlines.
	err := c.nc.SetReadDeadline(time.Now().Add(1 * time.Millisecond))
	if err != nil {
		return false
	}
	var b [1]byte
	_, err = c.nc.Read(b[:])
	return errors.Is(err, os.ErrDeadlineExceeded)
}

func (c *connection) idleTimeoutExpired() bool {
	if c.idleTimeout == 0 {
		return false
	}

	idleStart, ok := c.idleStart.Load().(time.Time)
	return ok && idleStart.Add(c.idleTimeout).Before(time.Now())
}

func (c *connection) bumpIdleStart() {
	if c.idleTimeout > 0 {
		c.idleStart.Store(time.Now())
	}
}

func (c *connection) setCanStream(canStream bool) {
	c.canStream = canStream
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

// DriverConnectionID returns the driver connection ID.
// TODO(GODRIVER-2824): change return type to int64.
func (c *connection) DriverConnectionID() uint64 {
	return c.driverConnectionID
}

func (c *connection) ID() string {
	return c.id
}

func (c *connection) ServerConnectionID() *int64 {
	return c.serverConnectionID
}

func (c *connection) OIDCTokenGenID() uint64 {
	return c.oidcTokenGenID
}

func (c *connection) SetOIDCTokenGenID(genID uint64) {
	c.oidcTokenGenID = genID
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
func (c initConnection) ReadWireMessage(ctx context.Context) ([]byte, error) {
	return c.readWireMessage(ctx)
}
func (c initConnection) SetStreaming(streaming bool) {
	c.setStreaming(streaming)
}
func (c initConnection) CurrentlyStreaming() bool {
	return c.getCurrentlyStreaming()
}
func (c initConnection) SupportsStreaming() bool {
	return c.canStream
}

// Connection implements the driver.Connection interface to allow reading and writing wire
// messages and the driver.Expirable interface to allow expiring. It wraps an underlying
// topology.connection to make it more goroutine-safe and nil-safe.
type Connection struct {
	connection    *connection
	refCount      int
	cleanupPoolFn func()

	oidcTokenGenID uint64

	// cleanupServerFn resets the server state when a connection is returned to the connection pool
	// via Close() or expired via Expire().
	cleanupServerFn func()

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
	return c.connection.writeWireMessage(ctx, wm)
}

// ReadWireMessage handles reading a wire message from the underlying connection. The dst parameter
// will be overwritten with the new wire message.
func (c *Connection) ReadWireMessage(ctx context.Context) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.connection == nil {
		return nil, ErrConnectionClosed
	}
	return c.connection.readWireMessage(ctx)
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
	return c.connection.desc
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

	_ = c.connection.close()
	return c.cleanupReferences()
}

func (c *Connection) cleanupReferences() error {
	err := c.connection.pool.checkIn(c.connection)
	if c.cleanupPoolFn != nil {
		c.cleanupPoolFn()
		c.cleanupPoolFn = nil
	}
	if c.cleanupServerFn != nil {
		c.cleanupServerFn()
		c.cleanupServerFn = nil
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
	return c.connection.id
}

// ServerConnectionID returns the server connection ID of this connection.
func (c *Connection) ServerConnectionID() *int64 {
	if c.connection == nil {
		return nil
	}
	return c.connection.serverConnectionID
}

// Stale returns if the connection is stale.
func (c *Connection) Stale() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connection.pool.stale(c.connection)
}

// Address returns the address of this connection.
func (c *Connection) Address() address.Address {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.connection == nil {
		return address.Address("0.0.0.0")
	}
	return c.connection.addr
}

// LocalAddress returns the local address of the connection
func (c *Connection) LocalAddress() address.Address {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.connection == nil || c.connection.nc == nil {
		return address.Address("0.0.0.0")
	}
	return address.Address(c.connection.nc.LocalAddr().String())
}

// PinToCursor updates this connection to reflect that it is pinned to a cursor.
func (c *Connection) PinToCursor() error {
	return c.pin("cursor", c.connection.pool.pinConnectionToCursor, c.connection.pool.unpinConnectionFromCursor)
}

// PinToTransaction updates this connection to reflect that it is pinned to a transaction.
func (c *Connection) PinToTransaction() error {
	return c.pin("transaction", c.connection.pool.pinConnectionToTransaction, c.connection.pool.unpinConnectionFromTransaction)
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

// DriverConnectionID returns the driver connection ID.
// TODO(GODRIVER-2824): change return type to int64.
func (c *Connection) DriverConnectionID() uint64 {
	return c.connection.DriverConnectionID()
}

// OIDCTokenGenID returns the OIDC token generation ID.
func (c *Connection) OIDCTokenGenID() uint64 {
	return c.oidcTokenGenID
}

// SetOIDCTokenGenID sets the OIDC token generation ID.
func (c *Connection) SetOIDCTokenGenID(genID uint64) {
	c.oidcTokenGenID = genID
}

// TODO: Naming?

// cancellListener listens for context cancellation and notifies listeners via a
// callback function.
type cancellListener struct {
	aborted bool
	done    chan struct{}
}

// newCancellListener constructs a cancellListener.
func newCancellListener() *cancellListener {
	return &cancellListener{
		done: make(chan struct{}),
	}
}

// Listen blocks until the provided context is cancelled or listening is aborted
// via the StopListening function. If this detects that the context has been
// cancelled (i.e. errors.Is(ctx.Err(), context.Canceled), the provided callback is
// called to abort in-progress work. Even if the context expires, this function
// will block until StopListening is called.
func (c *cancellListener) Listen(ctx context.Context, abortFn func()) {
	c.aborted = false

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.Canceled) {
			c.aborted = true
			abortFn()
		}

		<-c.done
	case <-c.done:
	}
}

// StopListening stops the in-progress Listen call. This blocks if there is no
// in-progress Listen call. This function will return true if the provided abort
// callback was called when listening for cancellation on the previous context.
func (c *cancellListener) StopListening() bool {
	c.done <- struct{}{}
	return c.aborted
}
