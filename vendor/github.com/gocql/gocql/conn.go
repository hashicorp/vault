/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/*
 * Content before git sha 34fdeebefcbf183ed7f916f931aa0586fdaa1b40
 * Copyright (c) 2012, The Gocql authors,
 * provided under the BSD-3-Clause License.
 * See the NOTICE file distributed with this work for additional information.
 */

package gocql

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gocql/gocql/internal/lru"
	"github.com/gocql/gocql/internal/streams"
)

var (
	defaultApprovedAuthenticators = []string{
		"org.apache.cassandra.auth.PasswordAuthenticator",
		"com.instaclustr.cassandra.auth.SharedSecretAuthenticator",
		"com.datastax.bdp.cassandra.auth.DseAuthenticator",
		"io.aiven.cassandra.auth.AivenAuthenticator",
		"com.ericsson.bss.cassandra.ecaudit.auth.AuditPasswordAuthenticator",
		"com.amazon.helenus.auth.HelenusAuthenticator",
		"com.ericsson.bss.cassandra.ecaudit.auth.AuditAuthenticator",
		"com.scylladb.auth.SaslauthdAuthenticator",
		"com.scylladb.auth.TransitionalAuthenticator",
		"com.instaclustr.cassandra.auth.InstaclustrPasswordAuthenticator",
	}
)

// approve the authenticator with the list of allowed authenticators or default list if approvedAuthenticators is empty.
func approve(authenticator string, approvedAuthenticators []string) bool {
	if len(approvedAuthenticators) == 0 {
		approvedAuthenticators = defaultApprovedAuthenticators
	}
	for _, s := range approvedAuthenticators {
		if authenticator == s {
			return true
		}
	}
	return false
}

// JoinHostPort is a utility to return an address string that can be used
// by `gocql.Conn` to form a connection with a host.
func JoinHostPort(addr string, port int) string {
	addr = strings.TrimSpace(addr)
	if _, _, err := net.SplitHostPort(addr); err != nil {
		addr = net.JoinHostPort(addr, strconv.Itoa(port))
	}
	return addr
}

type Authenticator interface {
	Challenge(req []byte) (resp []byte, auth Authenticator, err error)
	Success(data []byte) error
}

type PasswordAuthenticator struct {
	Username              string
	Password              string
	AllowedAuthenticators []string
}

func (p PasswordAuthenticator) Challenge(req []byte) ([]byte, Authenticator, error) {
	if !approve(string(req), p.AllowedAuthenticators) {
		return nil, nil, fmt.Errorf("unexpected authenticator %q", req)
	}
	resp := make([]byte, 2+len(p.Username)+len(p.Password))
	resp[0] = 0
	copy(resp[1:], p.Username)
	resp[len(p.Username)+1] = 0
	copy(resp[2+len(p.Username):], p.Password)
	return resp, nil, nil
}

func (p PasswordAuthenticator) Success(data []byte) error {
	return nil
}

// SslOptions configures TLS use.
//
// Warning: Due to historical reasons, the SslOptions is insecure by default, so you need to set EnableHostVerification
// to true if no Config is set. Most users should set SslOptions.Config to a *tls.Config.
// SslOptions and Config.InsecureSkipVerify interact as follows:
//
//	Config.InsecureSkipVerify | EnableHostVerification | Result
//	Config is nil             | false                  | do not verify host
//	Config is nil             | true                   | verify host
//	false                     | false                  | verify host
//	true                      | false                  | do not verify host
//	false                     | true                   | verify host
//	true                      | true                   | verify host
type SslOptions struct {
	*tls.Config

	// CertPath and KeyPath are optional depending on server
	// config, but both fields must be omitted to avoid using a
	// client certificate
	CertPath string
	KeyPath  string
	CaPath   string //optional depending on server config
	// If you want to verify the hostname and server cert (like a wildcard for cass cluster) then you should turn this
	// on.
	// This option is basically the inverse of tls.Config.InsecureSkipVerify.
	// See InsecureSkipVerify in http://golang.org/pkg/crypto/tls/ for more info.
	//
	// See SslOptions documentation to see how EnableHostVerification interacts with the provided tls.Config.
	EnableHostVerification bool
}

type ConnConfig struct {
	ProtoVersion   int
	CQLVersion     string
	Timeout        time.Duration
	WriteTimeout   time.Duration
	ConnectTimeout time.Duration
	Dialer         Dialer
	HostDialer     HostDialer
	Compressor     Compressor
	Authenticator  Authenticator
	AuthProvider   func(h *HostInfo) (Authenticator, error)
	Keepalive      time.Duration
	Logger         StdLogger

	tlsConfig       *tls.Config
	disableCoalesce bool
}

func (c *ConnConfig) logger() StdLogger {
	if c.Logger == nil {
		return Logger
	}
	return c.Logger
}

type ConnErrorHandler interface {
	HandleError(conn *Conn, err error, closed bool)
}

type connErrorHandlerFn func(conn *Conn, err error, closed bool)

func (fn connErrorHandlerFn) HandleError(conn *Conn, err error, closed bool) {
	fn(conn, err, closed)
}

// If not zero, how many timeouts we will allow to occur before the connection is closed
// and restarted. This is to prevent a single query timeout from killing a connection
// which may be serving more queries just fine.
// Default is 0, should not be changed concurrently with queries.
//
// Deprecated.
var TimeoutLimit int64 = 0

// Conn is a single connection to a Cassandra node. It can be used to execute
// queries, but users are usually advised to use a more reliable, higher
// level API.
type Conn struct {
	conn net.Conn
	r    *bufio.Reader
	w    contextWriter

	timeout        time.Duration
	writeTimeout   time.Duration
	cfg            *ConnConfig
	frameObserver  FrameHeaderObserver
	streamObserver StreamObserver

	headerBuf [maxFrameHeaderSize]byte

	streams *streams.IDGenerator
	mu      sync.Mutex
	// calls stores a map from stream ID to callReq.
	// This map is protected by mu.
	// calls should not be used when closed is true, calls is set to nil when closed=true.
	calls map[int]*callReq

	errorHandler ConnErrorHandler
	compressor   Compressor
	auth         Authenticator
	addr         string

	version         uint8
	currentKeyspace string
	host            *HostInfo
	isSchemaV2      bool

	session *Session

	// true if connection close process for the connection started.
	// closed is protected by mu.
	closed bool
	ctx    context.Context
	cancel context.CancelFunc

	timeouts int64

	logger StdLogger
}

// connect establishes a connection to a Cassandra node using session's connection config.
func (s *Session) connect(ctx context.Context, host *HostInfo, errorHandler ConnErrorHandler) (*Conn, error) {
	return s.dial(ctx, host, s.connCfg, errorHandler)
}

// dial establishes a connection to a Cassandra node and notifies the session's connectObserver.
func (s *Session) dial(ctx context.Context, host *HostInfo, connConfig *ConnConfig, errorHandler ConnErrorHandler) (*Conn, error) {
	var obs ObservedConnect
	if s.connectObserver != nil {
		obs.Host = host
		obs.Start = time.Now()
	}

	conn, err := s.dialWithoutObserver(ctx, host, connConfig, errorHandler)

	if s.connectObserver != nil {
		obs.End = time.Now()
		obs.Err = err
		s.connectObserver.ObserveConnect(obs)
	}

	return conn, err
}

// dialWithoutObserver establishes connection to a Cassandra node.
//
// dialWithoutObserver does not notify the connection observer, so you most probably want to call dial() instead.
func (s *Session) dialWithoutObserver(ctx context.Context, host *HostInfo, cfg *ConnConfig, errorHandler ConnErrorHandler) (*Conn, error) {
	dialedHost, err := cfg.HostDialer.DialHost(ctx, host)
	if err != nil {
		return nil, err
	}

	writeTimeout := cfg.Timeout
	if cfg.WriteTimeout > 0 {
		writeTimeout = cfg.WriteTimeout
	}

	ctx, cancel := context.WithCancel(ctx)
	c := &Conn{
		conn:          dialedHost.Conn,
		r:             bufio.NewReader(dialedHost.Conn),
		cfg:           cfg,
		calls:         make(map[int]*callReq),
		version:       uint8(cfg.ProtoVersion),
		addr:          dialedHost.Conn.RemoteAddr().String(),
		errorHandler:  errorHandler,
		compressor:    cfg.Compressor,
		session:       s,
		streams:       streams.New(cfg.ProtoVersion),
		host:          host,
		isSchemaV2:    true, // Try using "system.peers_v2" until proven otherwise
		frameObserver: s.frameObserver,
		w: &deadlineContextWriter{
			w:         dialedHost.Conn,
			timeout:   writeTimeout,
			semaphore: make(chan struct{}, 1),
			quit:      make(chan struct{}),
		},
		ctx:            ctx,
		cancel:         cancel,
		logger:         cfg.logger(),
		streamObserver: s.streamObserver,
		writeTimeout:   writeTimeout,
	}

	if err := c.init(ctx, dialedHost); err != nil {
		cancel()
		c.Close()
		return nil, err
	}

	return c, nil
}

func (c *Conn) init(ctx context.Context, dialedHost *DialedHost) error {
	if c.session.cfg.AuthProvider != nil {
		var err error
		c.auth, err = c.cfg.AuthProvider(c.host)
		if err != nil {
			return err
		}
	} else {
		c.auth = c.cfg.Authenticator
	}

	startup := &startupCoordinator{
		frameTicker: make(chan struct{}),
		conn:        c,
	}

	c.timeout = c.cfg.ConnectTimeout
	if err := startup.setupConn(ctx); err != nil {
		return err
	}

	c.timeout = c.cfg.Timeout

	// dont coalesce startup frames
	if c.session.cfg.WriteCoalesceWaitTime > 0 && !c.cfg.disableCoalesce && !dialedHost.DisableCoalesce {
		c.w = newWriteCoalescer(c.conn, c.writeTimeout, c.session.cfg.WriteCoalesceWaitTime, ctx.Done())
	}

	go c.serve(ctx)
	go c.heartBeat(ctx)

	return nil
}

func (c *Conn) Write(p []byte) (n int, err error) {
	return c.w.writeContext(context.Background(), p)
}

func (c *Conn) Read(p []byte) (n int, err error) {
	const maxAttempts = 5

	for i := 0; i < maxAttempts; i++ {
		var nn int
		if c.timeout > 0 {
			c.conn.SetReadDeadline(time.Now().Add(c.timeout))
		}

		nn, err = io.ReadFull(c.r, p[n:])
		n += nn
		if err == nil {
			break
		}

		if verr, ok := err.(net.Error); !ok || !verr.Temporary() {
			break
		}
	}

	return
}

type startupCoordinator struct {
	conn        *Conn
	frameTicker chan struct{}
}

func (s *startupCoordinator) setupConn(ctx context.Context) error {
	var cancel context.CancelFunc
	if s.conn.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, s.conn.timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	startupErr := make(chan error)
	go func() {
		for range s.frameTicker {
			err := s.conn.recv(ctx)
			if err != nil {
				select {
				case startupErr <- err:
				case <-ctx.Done():
				}

				return
			}
		}
	}()

	go func() {
		defer close(s.frameTicker)
		err := s.options(ctx)
		select {
		case startupErr <- err:
		case <-ctx.Done():
		}
	}()

	select {
	case err := <-startupErr:
		if err != nil {
			return err
		}
	case <-ctx.Done():
		return errors.New("gocql: no response to connection startup within timeout")
	}

	return nil
}

func (s *startupCoordinator) write(ctx context.Context, frame frameBuilder) (frame, error) {
	select {
	case s.frameTicker <- struct{}{}:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	framer, err := s.conn.exec(ctx, frame, nil)
	if err != nil {
		return nil, err
	}

	return framer.parseFrame()
}

func (s *startupCoordinator) options(ctx context.Context) error {
	frame, err := s.write(ctx, &writeOptionsFrame{})
	if err != nil {
		return err
	}

	supported, ok := frame.(*supportedFrame)
	if !ok {
		return NewErrProtocol("Unknown type of response to startup frame: %T", frame)
	}

	return s.startup(ctx, supported.supported)
}

func (s *startupCoordinator) startup(ctx context.Context, supported map[string][]string) error {
	m := map[string]string{
		"CQL_VERSION":    s.conn.cfg.CQLVersion,
		"DRIVER_NAME":    driverName,
		"DRIVER_VERSION": driverVersion,
	}

	if s.conn.compressor != nil {
		comp := supported["COMPRESSION"]
		name := s.conn.compressor.Name()
		for _, compressor := range comp {
			if compressor == name {
				m["COMPRESSION"] = compressor
				break
			}
		}

		if _, ok := m["COMPRESSION"]; !ok {
			s.conn.compressor = nil
		}
	}

	frame, err := s.write(ctx, &writeStartupFrame{opts: m})
	if err != nil {
		return err
	}

	switch v := frame.(type) {
	case error:
		return v
	case *readyFrame:
		return nil
	case *authenticateFrame:
		return s.authenticateHandshake(ctx, v)
	default:
		return NewErrProtocol("Unknown type of response to startup frame: %s", v)
	}
}

func (s *startupCoordinator) authenticateHandshake(ctx context.Context, authFrame *authenticateFrame) error {
	if s.conn.auth == nil {
		return fmt.Errorf("authentication required (using %q)", authFrame.class)
	}

	resp, challenger, err := s.conn.auth.Challenge([]byte(authFrame.class))
	if err != nil {
		return err
	}

	req := &writeAuthResponseFrame{data: resp}
	for {
		frame, err := s.write(ctx, req)
		if err != nil {
			return err
		}

		switch v := frame.(type) {
		case error:
			return v
		case *authSuccessFrame:
			if challenger != nil {
				return challenger.Success(v.data)
			}
			return nil
		case *authChallengeFrame:
			resp, challenger, err = challenger.Challenge(v.data)
			if err != nil {
				return err
			}

			req = &writeAuthResponseFrame{
				data: resp,
			}
		default:
			return fmt.Errorf("unknown frame response during authentication: %v", v)
		}
	}
}

func (c *Conn) closeWithError(err error) {
	if c == nil {
		return
	}

	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return
	}
	c.closed = true

	var callsToClose map[int]*callReq

	// We should attempt to deliver the error back to the caller if it
	// exists. However, don't block c.mu while we are delivering the
	// error to outstanding calls.
	if err != nil {
		callsToClose = c.calls
		// It is safe to change c.calls to nil. Nobody should use it after c.closed is set to true.
		c.calls = nil
	}
	c.mu.Unlock()

	for _, req := range callsToClose {
		// we need to send the error to all waiting queries.
		select {
		case req.resp <- callResp{err: err}:
		case <-req.timeout:
		}
		if req.streamObserverContext != nil {
			req.streamObserverEndOnce.Do(func() {
				req.streamObserverContext.StreamAbandoned(ObservedStream{
					Host: c.host,
				})
			})
		}
	}

	// if error was nil then unblock the quit channel
	c.cancel()
	cerr := c.close()

	if err != nil {
		c.errorHandler.HandleError(c, err, true)
	} else if cerr != nil {
		// TODO(zariel): is it a good idea to do this?
		c.errorHandler.HandleError(c, cerr, true)
	}
}

func (c *Conn) close() error {
	return c.conn.Close()
}

func (c *Conn) Close() {
	c.closeWithError(nil)
}

// Serve starts the stream multiplexer for this connection, which is required
// to execute any queries. This method runs as long as the connection is
// open and is therefore usually called in a separate goroutine.
func (c *Conn) serve(ctx context.Context) {
	var err error
	for err == nil {
		err = c.recv(ctx)
	}

	c.closeWithError(err)
}

func (c *Conn) discardFrame(head frameHeader) error {
	_, err := io.CopyN(ioutil.Discard, c, int64(head.length))
	if err != nil {
		return err
	}
	return nil
}

type protocolError struct {
	frame frame
}

func (p *protocolError) Error() string {
	if err, ok := p.frame.(error); ok {
		return err.Error()
	}
	return fmt.Sprintf("gocql: received unexpected frame on stream %d: %v", p.frame.Header().stream, p.frame)
}

func (c *Conn) heartBeat(ctx context.Context) {
	sleepTime := 1 * time.Second
	timer := time.NewTimer(sleepTime)
	defer timer.Stop()

	var failures int

	for {
		if failures > 5 {
			c.closeWithError(fmt.Errorf("gocql: heartbeat failed"))
			return
		}

		timer.Reset(sleepTime)

		select {
		case <-ctx.Done():
			return
		case <-timer.C:
		}

		framer, err := c.exec(context.Background(), &writeOptionsFrame{}, nil)
		if err != nil {
			failures++
			continue
		}

		resp, err := framer.parseFrame()
		if err != nil {
			// invalid frame
			failures++
			continue
		}

		switch resp.(type) {
		case *supportedFrame:
			// Everything ok
			sleepTime = 5 * time.Second
			failures = 0
		case error:
			// TODO: should we do something here?
		default:
			panic(fmt.Sprintf("gocql: unknown frame in response to options: %T", resp))
		}
	}
}

func (c *Conn) recv(ctx context.Context) error {
	// not safe for concurrent reads

	// read a full header, ignore timeouts, as this is being ran in a loop
	// TODO: TCP level deadlines? or just query level deadlines?
	if c.timeout > 0 {
		c.conn.SetReadDeadline(time.Time{})
	}

	headStartTime := time.Now()
	// were just reading headers over and over and copy bodies
	head, err := readHeader(c.r, c.headerBuf[:])
	headEndTime := time.Now()
	if err != nil {
		return err
	}

	if c.frameObserver != nil {
		c.frameObserver.ObserveFrameHeader(context.Background(), ObservedFrameHeader{
			Version: protoVersion(head.version),
			Flags:   head.flags,
			Stream:  int16(head.stream),
			Opcode:  frameOp(head.op),
			Length:  int32(head.length),
			Start:   headStartTime,
			End:     headEndTime,
			Host:    c.host,
		})
	}

	if head.stream > c.streams.NumStreams {
		return fmt.Errorf("gocql: frame header stream is beyond call expected bounds: %d", head.stream)
	} else if head.stream == -1 {
		// TODO: handle cassandra event frames, we shouldnt get any currently
		framer := newFramer(c.compressor, c.version)
		if err := framer.readFrame(c, &head); err != nil {
			return err
		}
		go c.session.handleEvent(framer)
		return nil
	} else if head.stream <= 0 {
		// reserved stream that we dont use, probably due to a protocol error
		// or a bug in Cassandra, this should be an error, parse it and return.
		framer := newFramer(c.compressor, c.version)
		if err := framer.readFrame(c, &head); err != nil {
			return err
		}

		frame, err := framer.parseFrame()
		if err != nil {
			return err
		}

		return &protocolError{
			frame: frame,
		}
	}

	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return ErrConnectionClosed
	}
	call, ok := c.calls[head.stream]
	delete(c.calls, head.stream)
	c.mu.Unlock()
	if call == nil || !ok {
		c.logger.Printf("gocql: received response for stream which has no handler: header=%v\n", head)
		return c.discardFrame(head)
	} else if head.stream != call.streamID {
		panic(fmt.Sprintf("call has incorrect streamID: got %d expected %d", call.streamID, head.stream))
	}

	framer := newFramer(c.compressor, c.version)

	err = framer.readFrame(c, &head)
	if err != nil {
		// only net errors should cause the connection to be closed. Though
		// cassandra returning corrupt frames will be returned here as well.
		if _, ok := err.(net.Error); ok {
			return err
		}
	}

	// we either, return a response to the caller, the caller timedout, or the
	// connection has closed. Either way we should never block indefinatly here
	select {
	case call.resp <- callResp{framer: framer, err: err}:
	case <-call.timeout:
		c.releaseStream(call)
	case <-ctx.Done():
	}

	return nil
}

func (c *Conn) releaseStream(call *callReq) {
	if call.timer != nil {
		call.timer.Stop()
	}

	c.streams.Clear(call.streamID)

	if call.streamObserverContext != nil {
		call.streamObserverEndOnce.Do(func() {
			call.streamObserverContext.StreamFinished(ObservedStream{
				Host: c.host,
			})
		})
	}
}

func (c *Conn) handleTimeout() {
	if TimeoutLimit > 0 && atomic.AddInt64(&c.timeouts, 1) > TimeoutLimit {
		c.closeWithError(ErrTooManyTimeouts)
	}
}

type callReq struct {
	// resp will receive the frame that was sent as a response to this stream.
	resp     chan callResp
	timeout  chan struct{} // indicates to recv() that a call has timed out
	streamID int           // current stream in use

	timer *time.Timer

	// streamObserverContext is notified about events regarding this stream
	streamObserverContext StreamObserverContext

	// streamObserverEndOnce ensures that either StreamAbandoned or StreamFinished is called,
	// but not both.
	streamObserverEndOnce sync.Once
}

type callResp struct {
	// framer is the response frame.
	// May be nil if err is not nil.
	framer *framer
	// err is error encountered, if any.
	err error
}

// contextWriter is like io.Writer, but takes context as well.
type contextWriter interface {
	// writeContext writes p to the connection.
	//
	// If ctx is canceled before we start writing p (e.g. during waiting while another write is currently in progress),
	// p is not written and ctx.Err() is returned. Context is ignored after we start writing p (i.e. we don't interrupt
	// blocked writes that are in progress) so that we always either write the full frame or not write it at all.
	//
	// It returns the number of bytes written from p (0 <= n <= len(p)) and any error that caused the write to stop
	// early. writeContext must return a non-nil error if it returns n < len(p). writeContext must not modify the
	// data in p, even temporarily.
	writeContext(ctx context.Context, p []byte) (n int, err error)
}

type deadlineWriter interface {
	SetWriteDeadline(time.Time) error
	io.Writer
}

type deadlineContextWriter struct {
	w       deadlineWriter
	timeout time.Duration
	// semaphore protects critical section for SetWriteDeadline/Write.
	// It is a channel with capacity 1.
	semaphore chan struct{}

	// quit closed once the connection is closed.
	quit chan struct{}
}

// writeContext implements contextWriter.
func (c *deadlineContextWriter) writeContext(ctx context.Context, p []byte) (int, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case <-c.quit:
		return 0, ErrConnectionClosed
	case c.semaphore <- struct{}{}:
		// acquired
	}

	defer func() {
		// release
		<-c.semaphore
	}()

	if c.timeout > 0 {
		err := c.w.SetWriteDeadline(time.Now().Add(c.timeout))
		if err != nil {
			return 0, err
		}
	}
	return c.w.Write(p)
}

func newWriteCoalescer(conn deadlineWriter, writeTimeout, coalesceDuration time.Duration,
	quit <-chan struct{}) *writeCoalescer {
	wc := &writeCoalescer{
		writeCh: make(chan writeRequest),
		c:       conn,
		quit:    quit,
		timeout: writeTimeout,
	}
	go wc.writeFlusher(coalesceDuration)
	return wc
}

type writeCoalescer struct {
	c deadlineWriter

	mu sync.Mutex

	quit    <-chan struct{}
	writeCh chan writeRequest

	timeout time.Duration

	testEnqueuedHook func()
	testFlushedHook  func()
}

type writeRequest struct {
	// resultChan is a channel (with buffer size 1) where to send results of the write.
	resultChan chan<- writeResult
	// data to write.
	data []byte
}

type writeResult struct {
	n   int
	err error
}

// writeContext implements contextWriter.
func (w *writeCoalescer) writeContext(ctx context.Context, p []byte) (int, error) {
	resultChan := make(chan writeResult, 1)
	wr := writeRequest{
		resultChan: resultChan,
		data:       p,
	}

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case <-w.quit:
		return 0, io.EOF // TODO: better error here?
	case w.writeCh <- wr:
		// enqueued for writing
	}

	if w.testEnqueuedHook != nil {
		w.testEnqueuedHook()
	}

	result := <-resultChan
	return result.n, result.err
}

func (w *writeCoalescer) writeFlusher(interval time.Duration) {
	timer := time.NewTimer(interval)
	defer timer.Stop()

	if !timer.Stop() {
		<-timer.C
	}

	w.writeFlusherImpl(timer.C, func() { timer.Reset(interval) })
}

func (w *writeCoalescer) writeFlusherImpl(timerC <-chan time.Time, resetTimer func()) {
	running := false

	var buffers net.Buffers
	var resultChans []chan<- writeResult

	for {
		select {
		case req := <-w.writeCh:
			buffers = append(buffers, req.data)
			resultChans = append(resultChans, req.resultChan)
			if !running {
				// Start timer on first write.
				resetTimer()
				running = true
			}
		case <-w.quit:
			result := writeResult{
				n:   0,
				err: io.EOF, // TODO: better error here?
			}
			// Unblock whoever was waiting.
			for _, resultChan := range resultChans {
				// resultChan has capacity 1, so it does not block.
				resultChan <- result
			}
			return
		case <-timerC:
			running = false
			w.flush(resultChans, buffers)
			buffers = nil
			resultChans = nil
			if w.testFlushedHook != nil {
				w.testFlushedHook()
			}
		}
	}
}

func (w *writeCoalescer) flush(resultChans []chan<- writeResult, buffers net.Buffers) {
	// Flush everything we have so far.
	if w.timeout > 0 {
		err := w.c.SetWriteDeadline(time.Now().Add(w.timeout))
		if err != nil {
			for i := range resultChans {
				resultChans[i] <- writeResult{
					n:   0,
					err: err,
				}
			}
			return
		}
	}
	// Copy buffers because WriteTo modifies buffers in-place.
	buffers2 := make(net.Buffers, len(buffers))
	copy(buffers2, buffers)
	n, err := buffers2.WriteTo(w.c)
	// Writes of bytes before n succeeded, writes of bytes starting from n failed with err.
	// Use n as remaining byte counter.
	for i := range buffers {
		if int64(len(buffers[i])) <= n {
			// this buffer was fully written.
			resultChans[i] <- writeResult{
				n:   len(buffers[i]),
				err: nil,
			}
			n -= int64(len(buffers[i]))
		} else {
			// this buffer was not (fully) written.
			resultChans[i] <- writeResult{
				n:   int(n),
				err: err,
			}
			n = 0
		}
	}
}

// addCall attempts to add a call to c.calls.
// It fails with error if the connection already started closing or if a call for the given stream
// already exists.
func (c *Conn) addCall(call *callReq) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return ErrConnectionClosed
	}
	existingCall := c.calls[call.streamID]
	if existingCall != nil {
		return fmt.Errorf("attempting to use stream already in use: %d -> %d", call.streamID,
			existingCall.streamID)
	}
	c.calls[call.streamID] = call
	return nil
}

func (c *Conn) exec(ctx context.Context, req frameBuilder, tracer Tracer) (*framer, error) {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return nil, ctxErr
	}

	// TODO: move tracer onto conn
	stream, ok := c.streams.GetStream()
	if !ok {
		return nil, ErrNoStreams
	}

	// resp is basically a waiting semaphore protecting the framer
	framer := newFramer(c.compressor, c.version)

	call := &callReq{
		timeout:  make(chan struct{}),
		streamID: stream,
		resp:     make(chan callResp),
	}

	if c.streamObserver != nil {
		call.streamObserverContext = c.streamObserver.StreamContext(ctx)
	}

	if err := c.addCall(call); err != nil {
		return nil, err
	}

	// After this point, we need to either read from call.resp or close(call.timeout)
	// since closeWithError can try to write a connection close error to call.resp.
	// If we don't close(call.timeout) or read from call.resp, closeWithError can deadlock.

	if tracer != nil {
		framer.trace()
	}

	if call.streamObserverContext != nil {
		call.streamObserverContext.StreamStarted(ObservedStream{
			Host: c.host,
		})
	}

	err := req.buildFrame(framer, stream)
	if err != nil {
		// closeWithError will block waiting for this stream to either receive a response
		// or for us to timeout.
		close(call.timeout)
		// We failed to serialize the frame into a buffer.
		// This should not affect the connection as we didn't write anything. We just free the current call.
		c.mu.Lock()
		if !c.closed {
			delete(c.calls, call.streamID)
		}
		c.mu.Unlock()
		// We need to release the stream after we remove the call from c.calls, otherwise the existingCall != nil
		// check above could fail.
		c.releaseStream(call)
		return nil, err
	}

	n, err := c.w.writeContext(ctx, framer.buf)
	if err != nil {
		// closeWithError will block waiting for this stream to either receive a response
		// or for us to timeout, close the timeout chan here. Im not entirely sure
		// but we should not get a response after an error on the write side.
		close(call.timeout)
		if (errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) && n == 0 {
			// We have not started to write this frame.
			// Release the stream as no response can come from the server on the stream.
			c.mu.Lock()
			if !c.closed {
				delete(c.calls, call.streamID)
			}
			c.mu.Unlock()
			// We need to release the stream after we remove the call from c.calls, otherwise the existingCall != nil
			// check above could fail.
			c.releaseStream(call)
		} else {
			// I think this is the correct thing to do, im not entirely sure. It is not
			// ideal as readers might still get some data, but they probably wont.
			// Here we need to be careful as the stream is not available and if all
			// writes just timeout or fail then the pool might use this connection to
			// send a frame on, with all the streams used up and not returned.
			c.closeWithError(err)
		}
		return nil, err
	}

	var timeoutCh <-chan time.Time
	if c.timeout > 0 {
		if call.timer == nil {
			call.timer = time.NewTimer(0)
			<-call.timer.C
		} else {
			if !call.timer.Stop() {
				select {
				case <-call.timer.C:
				default:
				}
			}
		}

		call.timer.Reset(c.timeout)
		timeoutCh = call.timer.C
	}

	var ctxDone <-chan struct{}
	if ctx != nil {
		ctxDone = ctx.Done()
	}

	select {
	case resp := <-call.resp:
		close(call.timeout)
		if resp.err != nil {
			if !c.Closed() {
				// if the connection is closed then we cant release the stream,
				// this is because the request is still outstanding and we have
				// been handed another error from another stream which caused the
				// connection to close.
				c.releaseStream(call)
			}
			return nil, resp.err
		}
		// dont release the stream if detect a timeout as another request can reuse
		// that stream and get a response for the old request, which we have no
		// easy way of detecting.
		//
		// Ensure that the stream is not released if there are potentially outstanding
		// requests on the stream to prevent nil pointer dereferences in recv().
		defer c.releaseStream(call)

		if v := resp.framer.header.version.version(); v != c.version {
			return nil, NewErrProtocol("unexpected protocol version in response: got %d expected %d", v, c.version)
		}

		return resp.framer, nil
	case <-timeoutCh:
		close(call.timeout)
		c.handleTimeout()
		return nil, ErrTimeoutNoResponse
	case <-ctxDone:
		close(call.timeout)
		return nil, ctx.Err()
	case <-c.ctx.Done():
		close(call.timeout)
		return nil, ErrConnectionClosed
	}
}

// ObservedStream observes a single request/response stream.
type ObservedStream struct {
	// Host of the connection used to send the stream.
	Host *HostInfo
}

// StreamObserver is notified about request/response pairs.
// Streams are created for executing queries/batches or
// internal requests to the database and might live longer than
// execution of the query - the stream is still tracked until
// response arrives so that stream IDs are not reused.
type StreamObserver interface {
	// StreamContext is called before creating a new stream.
	// ctx is context passed to Session.Query / Session.Batch,
	// but might also be an internal context (for example
	// for internal requests that use control connection).
	// StreamContext might return nil if it is not interested
	// in the details of this stream.
	// StreamContext is called before the stream is created
	// and the returned StreamObserverContext might be discarded
	// without any methods called on the StreamObserverContext if
	// creation of the stream fails.
	// Note that if you don't need to track per-stream data,
	// you can always return the same StreamObserverContext.
	StreamContext(ctx context.Context) StreamObserverContext
}

// StreamObserverContext is notified about state of a stream.
// A stream is started every time a request is written to the server
// and is finished when a response is received.
// It is abandoned when the underlying network connection is closed
// before receiving a response.
type StreamObserverContext interface {
	// StreamStarted is called when the stream is started.
	// This happens just before a request is written to the wire.
	StreamStarted(observedStream ObservedStream)

	// StreamAbandoned is called when we stop waiting for response.
	// This happens when the underlying network connection is closed.
	// StreamFinished won't be called if StreamAbandoned is.
	StreamAbandoned(observedStream ObservedStream)

	// StreamFinished is called when we receive a response for the stream.
	StreamFinished(observedStream ObservedStream)
}

type preparedStatment struct {
	id       []byte
	request  preparedMetadata
	response resultMetadata
}

type inflightPrepare struct {
	done chan struct{}
	err  error

	preparedStatment *preparedStatment
}

func (c *Conn) prepareStatement(ctx context.Context, stmt string, tracer Tracer) (*preparedStatment, error) {
	stmtCacheKey := c.session.stmtsLRU.keyFor(c.host.HostID(), c.currentKeyspace, stmt)
	flight, ok := c.session.stmtsLRU.execIfMissing(stmtCacheKey, func(lru *lru.Cache) *inflightPrepare {
		flight := &inflightPrepare{
			done: make(chan struct{}),
		}
		lru.Add(stmtCacheKey, flight)
		return flight
	})

	if !ok {
		go func() {
			defer close(flight.done)

			prep := &writePrepareFrame{
				statement: stmt,
			}
			if c.version > protoVersion4 {
				prep.keyspace = c.currentKeyspace
			}

			// we won the race to do the load, if our context is canceled we shouldnt
			// stop the load as other callers are waiting for it but this caller should get
			// their context cancelled error.
			framer, err := c.exec(c.ctx, prep, tracer)
			if err != nil {
				flight.err = err
				c.session.stmtsLRU.remove(stmtCacheKey)
				return
			}

			frame, err := framer.parseFrame()
			if err != nil {
				flight.err = err
				c.session.stmtsLRU.remove(stmtCacheKey)
				return
			}

			// TODO(zariel): tidy this up, simplify handling of frame parsing so its not duplicated
			// everytime we need to parse a frame.
			if len(framer.traceID) > 0 && tracer != nil {
				tracer.Trace(framer.traceID)
			}

			switch x := frame.(type) {
			case *resultPreparedFrame:
				flight.preparedStatment = &preparedStatment{
					// defensively copy as we will recycle the underlying buffer after we
					// return.
					id: copyBytes(x.preparedID),
					// the type info's should _not_ have a reference to the framers read buffer,
					// therefore we can just copy them directly.
					request:  x.reqMeta,
					response: x.respMeta,
				}
			case error:
				flight.err = x
			default:
				flight.err = NewErrProtocol("Unknown type in response to prepare frame: %s", x)
			}

			if flight.err != nil {
				c.session.stmtsLRU.remove(stmtCacheKey)
			}
		}()
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-flight.done:
		return flight.preparedStatment, flight.err
	}
}

func marshalQueryValue(typ TypeInfo, value interface{}, dst *queryValues) error {
	if named, ok := value.(*namedValue); ok {
		dst.name = named.name
		value = named.value
	}

	if _, ok := value.(unsetColumn); !ok {
		val, err := Marshal(typ, value)
		if err != nil {
			return err
		}

		dst.value = val
	} else {
		dst.isUnset = true
	}

	return nil
}

func (c *Conn) executeQuery(ctx context.Context, qry *Query) *Iter {
	params := queryParams{
		consistency: qry.cons,
	}

	// frame checks that it is not 0
	params.serialConsistency = qry.serialCons
	params.defaultTimestamp = qry.defaultTimestamp
	params.defaultTimestampValue = qry.defaultTimestampValue

	if len(qry.pageState) > 0 {
		params.pagingState = qry.pageState
	}
	if qry.pageSize > 0 {
		params.pageSize = qry.pageSize
	}
	if c.version > protoVersion4 {
		params.keyspace = c.currentKeyspace
	}

	var (
		frame frameBuilder
		info  *preparedStatment
	)

	if !qry.skipPrepare && qry.shouldPrepare() {
		// Prepare all DML queries. Other queries can not be prepared.
		var err error
		info, err = c.prepareStatement(ctx, qry.stmt, qry.trace)
		if err != nil {
			return &Iter{err: err}
		}

		values := qry.values
		if qry.binding != nil {
			values, err = qry.binding(&QueryInfo{
				Id:          info.id,
				Args:        info.request.columns,
				Rval:        info.response.columns,
				PKeyColumns: info.request.pkeyColumns,
			})

			if err != nil {
				return &Iter{err: err}
			}
		}

		if len(values) != info.request.actualColCount {
			return &Iter{err: fmt.Errorf("gocql: expected %d values send got %d", info.request.actualColCount, len(values))}
		}

		params.values = make([]queryValues, len(values))
		for i := 0; i < len(values); i++ {
			v := &params.values[i]
			value := values[i]
			typ := info.request.columns[i].TypeInfo
			if err := marshalQueryValue(typ, value, v); err != nil {
				return &Iter{err: err}
			}
		}

		params.skipMeta = !(c.session.cfg.DisableSkipMetadata || qry.disableSkipMetadata)

		frame = &writeExecuteFrame{
			preparedID:    info.id,
			params:        params,
			customPayload: qry.customPayload,
		}

		// Set "keyspace" and "table" property in the query if it is present in preparedMetadata
		qry.routingInfo.mu.Lock()
		qry.routingInfo.keyspace = info.request.keyspace
		qry.routingInfo.table = info.request.table
		qry.routingInfo.mu.Unlock()
	} else {
		frame = &writeQueryFrame{
			statement:     qry.stmt,
			params:        params,
			customPayload: qry.customPayload,
		}
	}

	framer, err := c.exec(ctx, frame, qry.trace)
	if err != nil {
		return &Iter{err: err}
	}

	resp, err := framer.parseFrame()
	if err != nil {
		return &Iter{err: err}
	}

	if len(framer.traceID) > 0 && qry.trace != nil {
		qry.trace.Trace(framer.traceID)
	}

	switch x := resp.(type) {
	case *resultVoidFrame:
		return &Iter{framer: framer}
	case *resultRowsFrame:
		iter := &Iter{
			meta:    x.meta,
			framer:  framer,
			numRows: x.numRows,
		}

		if params.skipMeta {
			if info != nil {
				iter.meta = info.response
				iter.meta.pagingState = copyBytes(x.meta.pagingState)
			} else {
				return &Iter{framer: framer, err: errors.New("gocql: did not receive metadata but prepared info is nil")}
			}
		} else {
			iter.meta = x.meta
		}

		if x.meta.morePages() && !qry.disableAutoPage {
			newQry := new(Query)
			*newQry = *qry
			newQry.pageState = copyBytes(x.meta.pagingState)
			newQry.metrics = &queryMetrics{m: make(map[string]*hostMetrics)}

			iter.next = &nextIter{
				qry: newQry,
				pos: int((1 - qry.prefetch) * float64(x.numRows)),
			}

			if iter.next.pos < 1 {
				iter.next.pos = 1
			}
		}

		return iter
	case *resultKeyspaceFrame:
		return &Iter{framer: framer}
	case *schemaChangeKeyspace, *schemaChangeTable, *schemaChangeFunction, *schemaChangeAggregate, *schemaChangeType:
		iter := &Iter{framer: framer}
		if err := c.awaitSchemaAgreement(ctx); err != nil {
			// TODO: should have this behind a flag
			c.logger.Println(err)
		}
		// dont return an error from this, might be a good idea to give a warning
		// though. The impact of this returning an error would be that the cluster
		// is not consistent with regards to its schema.
		return iter
	case *RequestErrUnprepared:
		stmtCacheKey := c.session.stmtsLRU.keyFor(c.host.HostID(), c.currentKeyspace, qry.stmt)
		c.session.stmtsLRU.evictPreparedID(stmtCacheKey, x.StatementId)
		return c.executeQuery(ctx, qry)
	case error:
		return &Iter{err: x, framer: framer}
	default:
		return &Iter{
			err:    NewErrProtocol("Unknown type in response to execute query (%T): %s", x, x),
			framer: framer,
		}
	}
}

func (c *Conn) Pick(qry *Query) *Conn {
	if c.Closed() {
		return nil
	}
	return c
}

func (c *Conn) Closed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closed
}

func (c *Conn) Address() string {
	return c.addr
}

func (c *Conn) AvailableStreams() int {
	return c.streams.Available()
}

func (c *Conn) UseKeyspace(keyspace string) error {
	q := &writeQueryFrame{statement: `USE "` + keyspace + `"`}
	q.params.consistency = c.session.cons

	framer, err := c.exec(c.ctx, q, nil)
	if err != nil {
		return err
	}

	resp, err := framer.parseFrame()
	if err != nil {
		return err
	}

	switch x := resp.(type) {
	case *resultKeyspaceFrame:
	case error:
		return x
	default:
		return NewErrProtocol("unknown frame in response to USE: %v", x)
	}

	c.currentKeyspace = keyspace

	return nil
}

func (c *Conn) executeBatch(ctx context.Context, batch *Batch) *Iter {
	if c.version == protoVersion1 {
		return &Iter{err: ErrUnsupported}
	}

	n := len(batch.Entries)
	req := &writeBatchFrame{
		typ:                   batch.Type,
		statements:            make([]batchStatment, n),
		consistency:           batch.Cons,
		serialConsistency:     batch.serialCons,
		defaultTimestamp:      batch.defaultTimestamp,
		defaultTimestampValue: batch.defaultTimestampValue,
		customPayload:         batch.CustomPayload,
	}

	stmts := make(map[string]string, len(batch.Entries))

	for i := 0; i < n; i++ {
		entry := &batch.Entries[i]
		b := &req.statements[i]

		if len(entry.Args) > 0 || entry.binding != nil {
			info, err := c.prepareStatement(batch.Context(), entry.Stmt, batch.trace)
			if err != nil {
				return &Iter{err: err}
			}

			var values []interface{}
			if entry.binding == nil {
				values = entry.Args
			} else {
				values, err = entry.binding(&QueryInfo{
					Id:          info.id,
					Args:        info.request.columns,
					Rval:        info.response.columns,
					PKeyColumns: info.request.pkeyColumns,
				})
				if err != nil {
					return &Iter{err: err}
				}
			}

			if len(values) != info.request.actualColCount {
				return &Iter{err: fmt.Errorf("gocql: batch statement %d expected %d values send got %d", i, info.request.actualColCount, len(values))}
			}

			b.preparedID = info.id
			stmts[string(info.id)] = entry.Stmt

			b.values = make([]queryValues, info.request.actualColCount)

			for j := 0; j < info.request.actualColCount; j++ {
				v := &b.values[j]
				value := values[j]
				typ := info.request.columns[j].TypeInfo
				if err := marshalQueryValue(typ, value, v); err != nil {
					return &Iter{err: err}
				}
			}
		} else {
			b.statement = entry.Stmt
		}
	}

	framer, err := c.exec(batch.Context(), req, batch.trace)
	if err != nil {
		return &Iter{err: err}
	}

	resp, err := framer.parseFrame()
	if err != nil {
		return &Iter{err: err, framer: framer}
	}

	if len(framer.traceID) > 0 && batch.trace != nil {
		batch.trace.Trace(framer.traceID)
	}

	switch x := resp.(type) {
	case *resultVoidFrame:
		return &Iter{}
	case *RequestErrUnprepared:
		stmt, found := stmts[string(x.StatementId)]
		if found {
			key := c.session.stmtsLRU.keyFor(c.host.HostID(), c.currentKeyspace, stmt)
			c.session.stmtsLRU.evictPreparedID(key, x.StatementId)
		}
		return c.executeBatch(ctx, batch)
	case *resultRowsFrame:
		iter := &Iter{
			meta:    x.meta,
			framer:  framer,
			numRows: x.numRows,
		}

		return iter
	case error:
		return &Iter{err: x, framer: framer}
	default:
		return &Iter{err: NewErrProtocol("Unknown type in response to batch statement: %s", x), framer: framer}
	}
}

func (c *Conn) query(ctx context.Context, statement string, values ...interface{}) (iter *Iter) {
	q := c.session.Query(statement, values...).Consistency(One).Trace(nil)
	q.skipPrepare = true
	q.disableSkipMetadata = true
	// we want to keep the query on this connection
	q.conn = c
	return c.executeQuery(ctx, q)
}

func (c *Conn) querySystemPeers(ctx context.Context, version cassVersion) *Iter {
	const (
		peerSchema    = "SELECT * FROM system.peers"
		peerV2Schemas = "SELECT * FROM system.peers_v2"
	)

	c.mu.Lock()
	isSchemaV2 := c.isSchemaV2
	c.mu.Unlock()

	if version.AtLeast(4, 0, 0) && isSchemaV2 {
		// Try "system.peers_v2" and fallback to "system.peers" if it's not found
		iter := c.query(ctx, peerV2Schemas)

		err := iter.checkErrAndNotFound()
		if err != nil {
			if errFrame, ok := err.(errorFrame); ok && errFrame.code == ErrCodeInvalid { // system.peers_v2 not found, try system.peers
				c.mu.Lock()
				c.isSchemaV2 = false
				c.mu.Unlock()
				return c.query(ctx, peerSchema)
			} else {
				return iter
			}
		}
		return iter
	} else {
		return c.query(ctx, peerSchema)
	}
}

func (c *Conn) querySystemLocal(ctx context.Context) *Iter {
	return c.query(ctx, "SELECT * FROM system.local WHERE key='local'")
}

func (c *Conn) awaitSchemaAgreement(ctx context.Context) (err error) {
	const localSchemas = "SELECT schema_version FROM system.local WHERE key='local'"

	var versions map[string]struct{}
	var schemaVersion string

	endDeadline := time.Now().Add(c.session.cfg.MaxWaitSchemaAgreement)

	for time.Now().Before(endDeadline) {
		iter := c.querySystemPeers(ctx, c.host.version)

		versions = make(map[string]struct{})

		rows, err := iter.SliceMap()
		if err != nil {
			goto cont
		}

		for _, row := range rows {
			host, err := c.session.hostInfoFromMap(row, &HostInfo{connectAddress: c.host.ConnectAddress(), port: c.session.cfg.Port})
			if err != nil {
				goto cont
			}
			if !isValidPeer(host) || host.schemaVersion == "" {
				c.logger.Printf("invalid peer or peer with empty schema_version: peer=%q", host)
				continue
			}

			versions[host.schemaVersion] = struct{}{}
		}

		if err = iter.Close(); err != nil {
			goto cont
		}

		iter = c.query(ctx, localSchemas)
		for iter.Scan(&schemaVersion) {
			versions[schemaVersion] = struct{}{}
			schemaVersion = ""
		}

		if err = iter.Close(); err != nil {
			goto cont
		}

		if len(versions) <= 1 {
			return nil
		}

	cont:
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(200 * time.Millisecond):
		}
	}

	if err != nil {
		return err
	}

	schemas := make([]string, 0, len(versions))
	for schema := range versions {
		schemas = append(schemas, schema)
	}

	// not exported
	return fmt.Errorf("gocql: cluster schema versions not consistent: %+v", schemas)
}

var (
	ErrQueryArgLength    = errors.New("gocql: query argument length mismatch")
	ErrTimeoutNoResponse = errors.New("gocql: no response received from cassandra within timeout period")
	ErrTooManyTimeouts   = errors.New("gocql: too many query timeouts on the connection")
	ErrConnectionClosed  = errors.New("gocql: connection closed waiting for response")
	ErrNoStreams         = errors.New("gocql: no streams available on connection")
)
