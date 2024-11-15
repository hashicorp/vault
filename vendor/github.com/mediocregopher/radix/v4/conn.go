package radix

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mediocregopher/radix/v4/internal/proc"
	"github.com/mediocregopher/radix/v4/resp"
	"github.com/mediocregopher/radix/v4/resp/resp3"
)

// Conn is a Client wrapping a single network connection which synchronously
// reads/writes data using redis's RESP protocol.
//
// A Conn can be used directly as a Client, but in general you probably want to
// use a Pool instead.
type Conn interface {
	// The Do method merely calls the Action's Perform method with the Conn as
	// the argument.
	Client

	// EncodeDecode will encode marshal onto the connection, then decode a
	// response into unmarshalInto (see resp3.Marshal and resp3.Unmarshal,
	// respectively). If either parameter is nil then that step is skipped.
	//
	// If EncodeDecode is called concurrently on the same Conn then the order of
	// decode steps will match the order of encode steps.
	//
	// NOTE If ctx is canceled then marshaling, and possibly unmarshaling, might
	// still occur in the background even though EncodeDecode has returned.
	EncodeDecode(ctx context.Context, marshal, unmarshalInto interface{}) error
}

////////////////////////////////////////////////////////////////////////////////

type wrappedNetConn struct {
	net.Conn
	prevBytesRead, totalBytesRead int
	addr                          net.Addr
}

var _ net.Conn = new(wrappedNetConn)

func (c *wrappedNetConn) Read(b []byte) (int, error) {
	n, err := c.Conn.Read(b)
	c.totalBytesRead += n

	prevBytesRead := c.prevBytesRead
	c.prevBytesRead = n

	if err == nil || !errors.Is(err, os.ErrDeadlineExceeded) || (prevBytesRead == 0 && n == 0) {
		return n, err
	}

	// a timeout was reached, but there were bytes read before it was reached.
	// In that case we pretend no timeout was reached. If there's more data to
	// be pulled then we can do so, but if there's not more data to be pulled
	// then a timeout will be returned the subsequent call to Read.
	err = c.Conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))

	return n, err
}

func (c *wrappedNetConn) resetBytesRead() {
	c.prevBytesRead = 0
	c.totalBytesRead = 0
}

func (c *wrappedNetConn) RemoteAddr() net.Addr {
	return c.addr
}

////////////////////////////////////////////////////////////////////////////////

type connMarshalerUnmarshaler struct {
	ctx                    context.Context
	marshal, unmarshalInto interface{}
	errCh                  chan error
	readSeq                uint64
}

type conn struct {
	proc *proc.Proc

	conn         *wrappedNetConn
	rOpts, wOpts *resp.Opts
	br           resp.BufferedReader
	bw           resp.BufferedWriter
	rCh, wCh     chan connMarshalerUnmarshaler

	// the readSeq is used to track the current EncodeDecode call which the
	// reader go-routine is operating on. readSeq is used to assign a uint64 to
	// each connMarshalerUnmarshaler, and currReadSeq will be updated to reflect
	// the readSeq that the reader is currently operating on.
	//
	// this mechanism is only used during error cases, primarily when the
	// context deadline is exceeded and the EncodeDecode call wants to pre-empt
	// the currently active resp3.Unmarshal call.
	readSeqL             sync.Mutex
	readSeq, currReadSeq uint64

	// errChPool is a buffered channel used as a makeshift pool of chan errors,
	// so we don't have to make a new one on every EncodeDecode call.
	errChPool chan chan error
}

var _ Conn = new(conn)

func (c *conn) Close() error {
	return c.proc.PrefixedClose(c.conn.Close, nil)
}

func isRespErr(err error) bool {
	return errors.As(err, &resp3.SimpleError{}) || errors.As(err, &resp3.BlobError{})
}

func (c *conn) doDiscard(unmarshalInto interface{}) error {

	c.readSeqL.Lock()

	if err := c.conn.SetReadDeadline(time.Time{}); err != nil {
		c.readSeqL.Unlock()
		return fmt.Errorf("unsetting read deadline: %w", err)
	}

	// setting currReadSeq lets us discard the next message without any
	// possibility of being interrupted by an EncodeDecode routine.
	c.currReadSeq = 0

	c.readSeqL.Unlock()

	err := resp3.Unmarshal(c.br, unmarshalInto, c.rOpts)

	if isRespErr(err) {
		err = nil
	}

	return err
}

func (c *conn) prepareForRead(mu connMarshalerUnmarshaler) error {

	c.readSeqL.Lock()
	defer c.readSeqL.Unlock()

	if err := mu.ctx.Err(); err != nil {
		return resp.ErrConnUsable{
			Err: fmt.Errorf("checking context before read: %w", err),
		}
	}

	if err := c.conn.SetReadDeadline(time.Time{}); err != nil {
		return fmt.Errorf("unsetting read deadline: %w", err)
	}

	c.currReadSeq = mu.readSeq
	return nil
}

func (c *conn) doRead(mu connMarshalerUnmarshaler) error {

	if err := c.prepareForRead(mu); err != nil {
		return err
	}

	c.conn.resetBytesRead()
	buffered := c.br.Buffered()

	err := resp3.Unmarshal(c.br, mu.unmarshalInto, c.rOpts)

	// simplify things for the caller by translating network timeouts
	// into the context error, since that's actually what happened.
	var canceled bool
	if canceled = errors.Is(err, os.ErrDeadlineExceeded); canceled {

		err = mu.ctx.Err()

		// there might be some crazy edge case where this can happen, I'm not
		// sure... go contexts are not pleasant to work with.
		if err == nil {
			err = context.Canceled
		}

		// if Unmarshal returned a context error but data was also read it means
		// that a message was only partially read off the wire. The Conn is
		// unusable at this point, close it and bail.
		if buffered != c.br.Buffered() || c.conn.totalBytesRead > 0 {
			return fmt.Errorf("after partial read off Conn: %w", err)
		}

		return resp.ErrConnUsable{Err: err}
	}

	return err
}

func (c *conn) reader(ctx context.Context) {

	doneCh := ctx.Done()
	for {
		select {

		case <-doneCh:
			return

		case mu := <-c.rCh:

			err := c.doRead(mu)
			mu.errCh <- err

			if err == nil || isRespErr(err) {
				continue

			} else if !errors.As(err, new(resp.ErrConnUsable)) {
				go c.Close()
				return

				// if the EncodeDecode did not involve a value being marshaled
				// onto the wire, then we assume that we are not in a
				// command/response mode, but rather are waiting for
				// asynchronous writes (ie pubsub mode). In this case we don't
				// discard the next command coming down the wire.
			} else if mu.marshal != nil {

				if err := c.doDiscard(mu.unmarshalInto); err != nil {
					go c.Close()
					return
				}
			}
		}
	}
}

func (c *conn) getErrCh() chan error {
	select {
	case errCh := <-c.errChPool:
		return errCh
	default:
		return make(chan error, 1)
	}
}

func (c *conn) putErrCh(errCh chan error) {
	select {
	case c.errChPool <- errCh:
	default:
	}
}

func (c *conn) EncodeDecode(ctx context.Context, m, u interface{}) (err error) {

	mu := connMarshalerUnmarshaler{
		ctx:           ctx,
		marshal:       m,
		unmarshalInto: u,
		errCh:         c.getErrCh(),
		readSeq:       atomic.AddUint64(&c.readSeq, 1),
	}
	doneCh := ctx.Done()
	closedCh := c.proc.ClosedCh()

	select {
	case <-doneCh:
		return fmt.Errorf("writing EncodeDecode to Conn channel: %w", ctx.Err())
	case <-closedCh:
		return proc.ErrClosed
	case c.wCh <- mu:
	}

	select {
	case <-doneCh:

		// To ensure that we don't miss messages which might come in _just_ as
		// the deadline is exceeded, we only return the ctx.Err after hearing
		// back from the reader. We call SetReadDeadline with a past value (but
		// only if the reader is currently working on our message!) to unblock
		// it, if it's blocked. If the reader returns nil then the message was
		// successfully read despite the context being canceled.
		c.readSeqL.Lock()
		if c.currReadSeq == mu.readSeq {
			if err := c.conn.SetReadDeadline(time.Unix(0, 0)); err != nil {
				c.readSeqL.Unlock()
				return fmt.Errorf("canceling Conn read from EncodeDecode: %w", err)
			}
		}
		c.readSeqL.Unlock()

		select {
		case err := <-mu.errCh:
			c.putErrCh(mu.errCh)
			if err != nil {
				err = fmt.Errorf("waiting for response from Conn: %w", err)
			}
			return err
		case <-closedCh:
			return proc.ErrClosed
		}

	case <-closedCh:
		return proc.ErrClosed
	case err := <-mu.errCh:
		// it's important that we only put the error channel back in the pool if
		// it's actually been used, otherwise it might still end up with
		// something written to it.
		c.putErrCh(mu.errCh)
		if err != nil {
			err = fmt.Errorf("response returned from Conn: %w", err)
		}
		return err
	}
}

func (c *conn) Do(ctx context.Context, a Action) error {
	return a.Perform(ctx, c)
}

func (c *conn) Addr() net.Addr {
	return c.conn.RemoteAddr()
}

////////////////////////////////////////////////////////////////////////////////

// Dialer is used to create Conns with particular settings. All fields are
// optional, all methods are thread-safe.
type Dialer struct {
	// CustomConn indicates that this callback should be used in place of Dial
	// when Dial is called. All behavior of Dialer/Dial is superceded when this
	// is set.
	CustomConn func(ctx context.Context, network, addr string) (Conn, error)

	// AuthPass will cause Dial to perform an AUTH command once the connection
	// is created, using AuthUser (if given) and AuthPass.
	//
	// If this is set and a redis URI is passed to Dial which also has a password
	// set, this takes precedence.
	AuthUser, AuthPass string

	// SelectDB will cause Dial to perform a SELECT command once the connection
	// is created, using the given database index.
	//
	// If this is set and a redis URI is passed to Dial which also has a
	// database index set, this takes precedence.
	SelectDB string

	// Protocol can be used to automatically set the RESP protocol version.
	//
	// If Protocol is not empty the Dialer will send a HELLO command with the
	// value of Protocol as version, otherwise no HELLO command will be send.
	Protocol string

	// NetDialer is used to create the underlying network connection.
	//
	// Defaults to net.Dialer.
	NetDialer interface {
		DialContext(context.Context, string, string) (net.Conn, error)
	}

	// WriteFlushInterval indicates how often the Conn should flush writes
	// to the underlying net.Conn.
	//
	// Conn uses a bufio.Writer to write data to the underlying net.Conn, and so
	// requires Flush to be called on that bufio.Writer in order for the data to
	// be fully written. By delaying Flush calls until multiple concurrent
	// EncodeDecode calls have been made Conn can reduce system calls and
	// significantly improve performance in that case.
	//
	// All EncodeDecode calls will be delayed up to WriteFlushInterval, with one
	// exception: if more than WriteFlushInterval has elapsed since the last
	// EncodeDecode call then the next EncodeDecode will Flush immediately. This
	// allows Conns to behave well during both low and high activity periods.
	//
	// Defaults to 0, indicating Flush will be called upon each EncodeDecode
	// call without delay.
	WriteFlushInterval time.Duration

	// NewRespOpts returns a fresh instance of a *resp.Opts to be used by the
	// underlying connection. This maybe be called more than once.
	//
	// Defaults to resp.NewOpts.
	NewRespOpts func() *resp.Opts
}

func (d Dialer) withDefaults() Dialer {
	if d.NetDialer == nil {
		d.NetDialer = new(net.Dialer)
	}
	if d.NewRespOpts == nil {
		d.NewRespOpts = resp.NewOpts
	}
	return d
}

func parseRedisURL(urlStr string, d Dialer) (string, Dialer) {
	// do a quick check before we bust out url.Parse, in case that is very
	// unperformant
	if !strings.HasPrefix(urlStr, "redis://") {
		return urlStr, d
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr, d
	}

	q := u.Query()

	if d.AuthUser == "" {
		d.AuthUser = q.Get("username")
		if n := u.User.Username(); n != "" {
			d.AuthUser = n
		}
	}

	if d.AuthPass == "" {
		d.AuthPass = q.Get("password")
		if p, ok := u.User.Password(); ok {
			d.AuthPass = p
		}
	}

	if d.SelectDB == "" {
		d.SelectDB = q.Get("db")
		if u.Path != "" && u.Path != "/" {
			d.SelectDB = u.Path[1:]
		}
	}

	return u.Host, d
}

// Dial creates a Conn using the Dialer configuration.
//
// In place of a host:port address, Dial also accepts a URI, as per:
// 	https://www.iana.org/assignments/uri-schemes/prov/redis
// If the URI has an AUTH password or db specified Dial will attempt to perform
// the AUTH and/or SELECT as well.
func (d Dialer) Dial(ctx context.Context, network, addr string) (Conn, error) {
	if d.CustomConn != nil {
		return d.CustomConn(ctx, network, addr)
	}

	d = d.withDefaults()
	addr, d = parseRedisURL(addr, d)

	netConn, err := d.NetDialer.DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}

	// If the netConn is a net.TCPConn (or some wrapper for it) and so can have
	// keepalive enabled, do so with a sane (though slightly aggressive)
	// default.
	{
		type keepaliveConn interface {
			SetKeepAlive(bool) error
			SetKeepAlivePeriod(time.Duration) error
		}

		if kaConn, ok := netConn.(keepaliveConn); ok {
			if err = kaConn.SetKeepAlive(true); err != nil {
				netConn.Close()
				return nil, err
			} else if err = kaConn.SetKeepAlivePeriod(10 * time.Second); err != nil {
				netConn.Close()
				return nil, err
			}
		}
	}

	wrappedNetConn := &wrappedNetConn{
		Conn: netConn,

		// wrap the conn so that it will return exactly what was used for
		// dialing when Addr is called. If Conn's normal RemoteAddr is used then
		// it returns the fully resolved form of the host.
		addr: rawAddr{network: network, addr: addr},
	}

	conn := &conn{
		proc:      proc.New(),
		conn:      wrappedNetConn,
		rOpts:     d.NewRespOpts(),
		wOpts:     d.NewRespOpts(),
		rCh:       make(chan connMarshalerUnmarshaler, 128),
		wCh:       make(chan connMarshalerUnmarshaler, 128),
		errChPool: make(chan chan error, 16),
	}

	conn.br = conn.rOpts.GetBufferedReader(wrappedNetConn)
	conn.bw = conn.wOpts.GetBufferedWriter(wrappedNetConn)

	conn.proc.Run(conn.reader)
	conn.proc.Run(func(ctx context.Context) {
		conn.writer(ctx, d.WriteFlushInterval)
	})

	if d.Protocol != "" {
		args := []string{d.Protocol}
		if d.AuthUser != "" {
			args = append(args, "AUTH", d.AuthUser, d.AuthPass)
		} else if d.AuthPass != "" {
			args = append(args, "AUTH", "default", d.AuthPass)
		}
		if err := conn.Do(ctx, Cmd(nil, "HELLO", args...)); err != nil {
			conn.Close()
			return nil, err
		}
	} else if d.AuthUser != "" {
		if err := conn.Do(ctx, Cmd(nil, "AUTH", d.AuthUser, d.AuthPass)); err != nil {
			conn.Close()
			return nil, err
		}
	} else if d.AuthPass != "" {
		if err := conn.Do(ctx, Cmd(nil, "AUTH", d.AuthPass)); err != nil {
			conn.Close()
			return nil, err
		}
	}

	if d.SelectDB != "" {
		if err := conn.Do(ctx, Cmd(nil, "SELECT", d.SelectDB)); err != nil {
			conn.Close()
			return nil, err
		}
	}

	return conn, nil
}

// Dial is a shortcut for calling Dial on a zero-value Dialer.
func Dial(ctx context.Context, network, addr string) (Conn, error) {
	return (Dialer{}).Dial(ctx, network, addr)
}
