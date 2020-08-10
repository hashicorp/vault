// Copyright (c) 2012 The gocql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
	approvedAuthenticators = [...]string{
		"org.apache.cassandra.auth.PasswordAuthenticator",
		"com.instaclustr.cassandra.auth.SharedSecretAuthenticator",
		"com.datastax.bdp.cassandra.auth.DseAuthenticator",
		"io.aiven.cassandra.auth.AivenAuthenticator",
		"com.ericsson.bss.cassandra.ecaudit.auth.AuditPasswordAuthenticator",
		"com.amazon.helenus.auth.HelenusAuthenticator",
		"com.ericsson.bss.cassandra.ecaudit.auth.AuditAuthenticator",
	}
)

func approve(authenticator string) bool {
	for _, s := range approvedAuthenticators {
		if authenticator == s {
			return true
		}
	}
	return false
}

//JoinHostPort is a utility to return a address string that can be used
//gocql.Conn to form a connection with a host.
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
	Username string
	Password string
}

func (p PasswordAuthenticator) Challenge(req []byte) ([]byte, Authenticator, error) {
	if !approve(string(req)) {
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

type SslOptions struct {
	*tls.Config

	// CertPath and KeyPath are optional depending on server
	// config, but both fields must be omitted to avoid using a
	// client certificate
	CertPath string
	KeyPath  string
	CaPath   string //optional depending on server config
	// If you want to verify the hostname and server cert (like a wildcard for cass cluster) then you should turn this on
	// This option is basically the inverse of InSecureSkipVerify
	// See InSecureSkipVerify in http://golang.org/pkg/crypto/tls/ for more info
	EnableHostVerification bool
}

type ConnConfig struct {
	ProtoVersion   int
	CQLVersion     string
	Timeout        time.Duration
	ConnectTimeout time.Duration
	Dialer         Dialer
	Compressor     Compressor
	Authenticator  Authenticator
	AuthProvider   func(h *HostInfo) (Authenticator, error)
	Keepalive      time.Duration

	tlsConfig       *tls.Config
	disableCoalesce bool
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
// depreciated
var TimeoutLimit int64 = 0

// Conn is a single connection to a Cassandra node. It can be used to execute
// queries, but users are usually advised to use a more reliable, higher
// level API.
type Conn struct {
	conn net.Conn
	r    *bufio.Reader
	w    io.Writer

	timeout       time.Duration
	cfg           *ConnConfig
	frameObserver FrameHeaderObserver

	headerBuf [maxFrameHeaderSize]byte

	streams *streams.IDGenerator
	mu      sync.Mutex
	calls   map[int]*callReq

	errorHandler ConnErrorHandler
	compressor   Compressor
	auth         Authenticator
	addr         string

	version         uint8
	currentKeyspace string
	host            *HostInfo

	session *Session

	closed int32
	ctx    context.Context
	cancel context.CancelFunc

	timeouts int64
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
	ip := host.ConnectAddress()
	port := host.port

	// TODO(zariel): remove these
	if !validIpAddr(ip) {
		panic(fmt.Sprintf("host missing connect ip address: %v", ip))
	} else if port == 0 {
		panic(fmt.Sprintf("host missing port: %v", port))
	}

	dialer := cfg.Dialer
	if dialer == nil {
		d := &net.Dialer{
			Timeout: cfg.ConnectTimeout,
		}
		if cfg.Keepalive > 0 {
			d.KeepAlive = cfg.Keepalive
		}
		dialer = d
	}

	conn, err := dialer.DialContext(ctx, "tcp", host.HostnameAndPort())
	if err != nil {
		return nil, err
	}
	if cfg.tlsConfig != nil {
		// the TLS config is safe to be reused by connections but it must not
		// be modified after being used.
		tconn := tls.Client(conn, cfg.tlsConfig)
		if err := tconn.Handshake(); err != nil {
			conn.Close()
			return nil, err
		}
		conn = tconn
	}

	ctx, cancel := context.WithCancel(ctx)
	c := &Conn{
		conn:          conn,
		r:             bufio.NewReader(conn),
		cfg:           cfg,
		calls:         make(map[int]*callReq),
		version:       uint8(cfg.ProtoVersion),
		addr:          conn.RemoteAddr().String(),
		errorHandler:  errorHandler,
		compressor:    cfg.Compressor,
		session:       s,
		streams:       streams.New(cfg.ProtoVersion),
		host:          host,
		frameObserver: s.frameObserver,
		w: &deadlineWriter{
			w:       conn,
			timeout: cfg.Timeout,
		},
		ctx:    ctx,
		cancel: cancel,
	}

	if err := c.init(ctx); err != nil {
		cancel()
		c.Close()
		return nil, err
	}

	return c, nil
}

func (c *Conn) init(ctx context.Context) error {
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
	if c.session.cfg.WriteCoalesceWaitTime > 0 && !c.cfg.disableCoalesce {
		c.w = newWriteCoalescer(c.conn, c.timeout, c.session.cfg.WriteCoalesceWaitTime, ctx.Done())
	}

	go c.serve(ctx)
	go c.heartBeat(ctx)

	return nil
}

func (c *Conn) Write(p []byte) (n int, err error) {
	return c.w.Write(p)
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

func (s *startupCoordinator) write(ctx context.Context, frame frameWriter) (frame, error) {
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
		"CQL_VERSION": s.conn.cfg.CQLVersion,
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
	if !atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		return
	}

	// we should attempt to deliver the error back to the caller if it
	// exists
	if err != nil {
		c.mu.Lock()
		for _, req := range c.calls {
			// we need to send the error to all waiting queries, put the state
			// of this conn into not active so that it can not execute any queries.
			select {
			case req.resp <- err:
			case <-req.timeout:
			}
		}
		c.mu.Unlock()
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
		framer := newFramer(c, c, c.compressor, c.version)
		if err := framer.readFrame(&head); err != nil {
			return err
		}
		go c.session.handleEvent(framer)
		return nil
	} else if head.stream <= 0 {
		// reserved stream that we dont use, probably due to a protocol error
		// or a bug in Cassandra, this should be an error, parse it and return.
		framer := newFramer(c, c, c.compressor, c.version)
		if err := framer.readFrame(&head); err != nil {
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
	call, ok := c.calls[head.stream]
	delete(c.calls, head.stream)
	c.mu.Unlock()
	if call == nil || call.framer == nil || !ok {
		Logger.Printf("gocql: received response for stream which has no handler: header=%v\n", head)
		return c.discardFrame(head)
	} else if head.stream != call.streamID {
		panic(fmt.Sprintf("call has incorrect streamID: got %d expected %d", call.streamID, head.stream))
	}

	err = call.framer.readFrame(&head)
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
	case call.resp <- err:
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
}

func (c *Conn) handleTimeout() {
	if TimeoutLimit > 0 && atomic.AddInt64(&c.timeouts, 1) > TimeoutLimit {
		c.closeWithError(ErrTooManyTimeouts)
	}
}

type callReq struct {
	// could use a waitgroup but this allows us to do timeouts on the read/send
	resp     chan error
	framer   *framer
	timeout  chan struct{} // indicates to recv() that a call has timedout
	streamID int           // current stream in use

	timer *time.Timer
}

type deadlineWriter struct {
	w interface {
		SetWriteDeadline(time.Time) error
		io.Writer
	}
	timeout time.Duration
}

func (c *deadlineWriter) Write(p []byte) (int, error) {
	if c.timeout > 0 {
		c.w.SetWriteDeadline(time.Now().Add(c.timeout))
	}
	return c.w.Write(p)
}

func newWriteCoalescer(conn net.Conn, timeout time.Duration, d time.Duration, quit <-chan struct{}) *writeCoalescer {
	wc := &writeCoalescer{
		writeCh: make(chan struct{}), // TODO: could this be sync?
		cond:    sync.NewCond(&sync.Mutex{}),
		c:       conn,
		quit:    quit,
		timeout: timeout,
	}
	go wc.writeFlusher(d)
	return wc
}

type writeCoalescer struct {
	c net.Conn

	quit    <-chan struct{}
	writeCh chan struct{}
	running bool

	// cond waits for the buffer to be flushed
	cond    *sync.Cond
	buffers net.Buffers
	timeout time.Duration

	// result of the write
	err error
}

func (w *writeCoalescer) flushLocked() {
	w.running = false
	if len(w.buffers) == 0 {
		return
	}

	if w.timeout > 0 {
		w.c.SetWriteDeadline(time.Now().Add(w.timeout))
	}

	// Given we are going to do a fanout n is useless and according to
	// the docs WriteTo should return 0 and err or bytes written and
	// no error.
	_, w.err = w.buffers.WriteTo(w.c)
	if w.err != nil {
		w.buffers = nil
	}
	w.cond.Broadcast()
}

func (w *writeCoalescer) flush() {
	w.cond.L.Lock()
	w.flushLocked()
	w.cond.L.Unlock()
}

func (w *writeCoalescer) stop() {
	w.cond.L.Lock()
	defer w.cond.L.Unlock()

	w.flushLocked()
	// nil the channel out sends block forever on it
	// instead of closing which causes a send on closed channel
	// panic.
	w.writeCh = nil
}

func (w *writeCoalescer) Write(p []byte) (int, error) {
	w.cond.L.Lock()

	if !w.running {
		select {
		case w.writeCh <- struct{}{}:
			w.running = true
		case <-w.quit:
			w.cond.L.Unlock()
			return 0, io.EOF // TODO: better error here?
		}
	}

	w.buffers = append(w.buffers, p)
	for len(w.buffers) != 0 {
		w.cond.Wait()
	}

	err := w.err
	w.cond.L.Unlock()

	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (w *writeCoalescer) writeFlusher(interval time.Duration) {
	timer := time.NewTimer(interval)
	defer timer.Stop()
	defer w.stop()

	if !timer.Stop() {
		<-timer.C
	}

	for {
		// wait for a write to start the flush loop
		select {
		case <-w.writeCh:
		case <-w.quit:
			return
		}

		timer.Reset(interval)

		select {
		case <-w.quit:
			return
		case <-timer.C:
		}

		w.flush()
	}
}

func (c *Conn) exec(ctx context.Context, req frameWriter, tracer Tracer) (*framer, error) {
	// TODO: move tracer onto conn
	stream, ok := c.streams.GetStream()
	if !ok {
		return nil, ErrNoStreams
	}

	// resp is basically a waiting semaphore protecting the framer
	framer := newFramer(c, c, c.compressor, c.version)

	call := &callReq{
		framer:   framer,
		timeout:  make(chan struct{}),
		streamID: stream,
		resp:     make(chan error),
	}

	c.mu.Lock()
	existingCall := c.calls[stream]
	if existingCall == nil {
		c.calls[stream] = call
	}
	c.mu.Unlock()

	if existingCall != nil {
		return nil, fmt.Errorf("attempting to use stream already in use: %d -> %d", stream, existingCall.streamID)
	}

	if tracer != nil {
		framer.trace()
	}

	err := req.writeFrame(framer, stream)
	if err != nil {
		// closeWithError will block waiting for this stream to either receive a response
		// or for us to timeout, close the timeout chan here. Im not entirely sure
		// but we should not get a response after an error on the write side.
		close(call.timeout)
		// I think this is the correct thing to do, im not entirely sure. It is not
		// ideal as readers might still get some data, but they probably wont.
		// Here we need to be careful as the stream is not available and if all
		// writes just timeout or fail then the pool might use this connection to
		// send a frame on, with all the streams used up and not returned.
		c.closeWithError(err)
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
	case err := <-call.resp:
		close(call.timeout)
		if err != nil {
			if !c.Closed() {
				// if the connection is closed then we cant release the stream,
				// this is because the request is still outstanding and we have
				// been handed another error from another stream which caused the
				// connection to close.
				c.releaseStream(call)
			}
			return nil, err
		}
	case <-timeoutCh:
		close(call.timeout)
		c.handleTimeout()
		return nil, ErrTimeoutNoResponse
	case <-ctxDone:
		close(call.timeout)
		return nil, ctx.Err()
	case <-c.ctx.Done():
		return nil, ErrConnectionClosed
	}

	// dont release the stream if detect a timeout as another request can reuse
	// that stream and get a response for the old request, which we have no
	// easy way of detecting.
	//
	// Ensure that the stream is not released if there are potentially outstanding
	// requests on the stream to prevent nil pointer dereferences in recv().
	defer c.releaseStream(call)

	if v := framer.header.version.version(); v != c.version {
		return nil, NewErrProtocol("unexpected protocol version in response: got %d expected %d", v, c.version)
	}

	return framer, nil
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
	stmtCacheKey := c.session.stmtsLRU.keyFor(c.addr, c.currentKeyspace, stmt)
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
		frame frameWriter
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
			iter.next = &nextIter{
				qry: qry,
				pos: int((1 - qry.prefetch) * float64(x.numRows)),
			}

			iter.next.qry.pageState = copyBytes(x.meta.pagingState)
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
			Logger.Println(err)
		}
		// dont return an error from this, might be a good idea to give a warning
		// though. The impact of this returning an error would be that the cluster
		// is not consistent with regards to its schema.
		return iter
	case *RequestErrUnprepared:
		stmtCacheKey := c.session.stmtsLRU.keyFor(c.addr, c.currentKeyspace, qry.stmt)
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
	return atomic.LoadInt32(&c.closed) == 1
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
			info, err := c.prepareStatement(batch.Context(), entry.Stmt, nil)
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

	// TODO: should batch support tracing?
	framer, err := c.exec(batch.Context(), req, nil)
	if err != nil {
		return &Iter{err: err}
	}

	resp, err := framer.parseFrame()
	if err != nil {
		return &Iter{err: err, framer: framer}
	}

	switch x := resp.(type) {
	case *resultVoidFrame:
		return &Iter{}
	case *RequestErrUnprepared:
		stmt, found := stmts[string(x.StatementId)]
		if found {
			key := c.session.stmtsLRU.keyFor(c.addr, c.currentKeyspace, stmt)
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
	q := c.session.Query(statement, values...).Consistency(One)
	q.trace = nil
	q.skipPrepare = true
	q.disableSkipMetadata = true
	return c.executeQuery(ctx, q)
}

func (c *Conn) awaitSchemaAgreement(ctx context.Context) (err error) {
	const (
		peerSchemas  = "SELECT * FROM system.peers"
		localSchemas = "SELECT schema_version FROM system.local WHERE key='local'"
	)

	var versions map[string]struct{}
	var schemaVersion string

	endDeadline := time.Now().Add(c.session.cfg.MaxWaitSchemaAgreement)
	for time.Now().Before(endDeadline) {
		iter := c.query(ctx, peerSchemas)

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
				Logger.Printf("invalid peer or peer with empty schema_version: peer=%q", host)
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

func (c *Conn) localHostInfo(ctx context.Context) (*HostInfo, error) {
	row, err := c.query(ctx, "SELECT * FROM system.local WHERE key='local'").rowMap()
	if err != nil {
		return nil, err
	}

	port := c.conn.RemoteAddr().(*net.TCPAddr).Port

	// TODO(zariel): avoid doing this here
	host, err := c.session.hostInfoFromMap(row, &HostInfo{connectAddress: c.host.connectAddress, port: port})
	if err != nil {
		return nil, err
	}

	return c.session.ring.addOrUpdate(host), nil
}

var (
	ErrQueryArgLength    = errors.New("gocql: query argument length mismatch")
	ErrTimeoutNoResponse = errors.New("gocql: no response received from cassandra within timeout period")
	ErrTooManyTimeouts   = errors.New("gocql: too many query timeouts on the connection")
	ErrConnectionClosed  = errors.New("gocql: connection closed waiting for response")
	ErrNoStreams         = errors.New("gocql: no streams available on connection")
)
