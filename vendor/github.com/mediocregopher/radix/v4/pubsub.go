package radix

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/mediocregopher/radix/v4/resp"
	"github.com/mediocregopher/radix/v4/resp/resp3"
	"github.com/mediocregopher/radix/v4/trace"
	"github.com/tilinna/clock"
)

// PubSubMessage describes a message being published to a redis pubsub channel.
type PubSubMessage struct {
	Type    string // "message" or "pmessage"
	Pattern string // will be set if Type is "pmessage"
	Channel string
	Message []byte
}

// MarshalRESP implements the Marshaler interface.
func (m PubSubMessage) MarshalRESP(w io.Writer, o *resp.Opts) error {
	var err error
	marshal := func(m resp.Marshaler) {
		if err == nil {
			err = m.MarshalRESP(w, o)
		}
	}

	if m.Type == "message" {
		marshal(resp3.ArrayHeader{NumElems: 3})
		marshal(resp3.BlobString{S: m.Type})
	} else if m.Type == "pmessage" {
		marshal(resp3.ArrayHeader{NumElems: 4})
		marshal(resp3.BlobString{S: m.Type})
		marshal(resp3.BlobString{S: m.Pattern})
	} else {
		return errors.New("unknown message Type")
	}
	marshal(resp3.BlobString{S: m.Channel})
	marshal(resp3.BlobStringBytes{B: m.Message})
	return err
}

var errNotPubSubMessage = errors.New("message is not a PubSubMessage")

// UnmarshalRESP implements the Unmarshaler interface.
func (m *PubSubMessage) UnmarshalRESP(br resp.BufferedReader, o *resp.Opts) error {
	// This method will fully consume the message on the wire, regardless of if
	// it is a PubSubMessage or not. If it is not then errNotPubSubMessage is
	// returned.

	// When in subscribe mode redis only allows (P)(UN)SUBSCRIBE commands, which
	// all return arrays, and PING, which returns an array when in subscribe
	// mode. HOWEVER, when all channels have been unsubscribed from then the
	// connection will be taken _out_ of subscribe mode. This is theoretically
	// fine, since the driver will still only allow the 5 commands, except PING
	// will return a simple string when in the non-subscribed state. So this
	// needs to check for that.
	//
	// NOTE: This is technically only true for connections using RESP2, not for
	// connections using RESP3, but for backwards compatibility with RESP2 we
	// assume the same limitations for RESP3.
	if ok, _ := resp3.NextMessageIs(br, resp3.SimpleStringPrefix); ok {
		// if it's a simple string, discard it (it's probably PONG) and error
		if err := resp3.Unmarshal(br, nil, o); err != nil {
			return err
		}
		return resp.ErrConnUsable{Err: errNotPubSubMessage}
	}

	var numElems int
	if ok, _ := resp3.NextMessageIs(br, resp3.PushHeaderPrefix); ok {
		var ph resp3.PushHeader
		if err := ph.UnmarshalRESP(br, o); err != nil {
			return err
		}
		numElems = ph.NumElems
	} else {
		var ah resp3.ArrayHeader
		if err := ah.UnmarshalRESP(br, o); err != nil {
			return err
		}
		numElems = ah.NumElems
	}

	if numElems < 2 {
		return errors.New("message has too few elements")
	}

	var msgType resp3.BlobStringBytes
	if err := msgType.UnmarshalRESP(br, o); err != nil {
		return err
	}

	switch string(msgType.B) {
	case "message":
		m.Type = "message"
		if numElems != 3 {
			return errors.New("message has wrong number of elements")
		}
	case "pmessage":
		m.Type = "pmessage"
		if numElems != 4 {
			return errors.New("message has wrong number of elements")
		}

		var pattern resp3.BlobString
		if err := pattern.UnmarshalRESP(br, o); err != nil {
			return err
		}
		m.Pattern = pattern.S
	default:
		// if it's not a PubSubMessage then discard the rest of the array
		for i := 1; i < numElems; i++ {
			if err := resp3.Unmarshal(br, nil, o); err != nil {
				return err
			}
		}
		return resp.ErrConnUsable{Err: errNotPubSubMessage}
	}

	var channel resp3.BlobString
	if err := channel.UnmarshalRESP(br, o); err != nil {
		return err
	}
	m.Channel = channel.S

	var msg resp3.BlobStringBytes
	if err := msg.UnmarshalRESP(br, o); err != nil {
		return err
	}
	m.Message = msg.B

	return nil
}

////////////////////////////////////////////////////////////////////////////////

// PubSubConn wraps an existing Conn to support redis' pubsub system. Unlike
// Conn, a PubSubConn's methods are _not_ thread-safe.
type PubSubConn interface {

	// Subscribe subscribes the PubSubConn to the given set of channels.
	Subscribe(ctx context.Context, channels ...string) error

	// Unsubscribe unsubscribes the PubSubConn from the given set of channels.
	Unsubscribe(ctx context.Context, channels ...string) error

	// PSubscribe is like Subscribe, but it subscribes to a set of patterns and
	// not individual channels.
	PSubscribe(ctx context.Context, patterns ...string) error

	// PUnsubscribe is like Unsubscribe, but it unsubscribes from a set of
	// patterns and not individual channels.
	PUnsubscribe(ctx context.Context, patterns ...string) error

	// Ping performs a simple Ping command on the PubSubConn, returning an error
	// if it failed for some reason.
	//
	// Ping will be periodically called by Next in the default PubSubConn
	// implementations.
	Ping(ctx context.Context) error

	// Next blocks until a message is published to the PubSubConn or an error is
	// encountered. If the context is canceled then the resulting error is
	// returned immediately.
	Next(ctx context.Context) (PubSubMessage, error)

	// Close closes the PubSubConn and cleans up all resources it holds.
	Close() error
}

// PubSubConfig is used to create a PubSubConn with particular settings. All
// fields are optional, all methods are thread-safe.
type PubSubConfig struct {

	// PingInterval is the interval at which PING will be called on the
	// PubSubConn in the background.
	//
	// Defaults to 5 * time.Second. Can be set to -1 to disable periodic pings.
	PingInterval time.Duration
}

type pubSubConfig struct {
	PubSubConfig
	clock clock.Clock

	testEventCh chan string
}

const pubSubDefaultPingInterval = 5 * time.Second

func (cfg pubSubConfig) withDefaults() pubSubConfig {
	if cfg.PingInterval == -1 {
		cfg.PingInterval = 0
	} else if cfg.PingInterval == 0 {
		cfg.PingInterval = pubSubDefaultPingInterval
	}

	if cfg.clock == nil {
		cfg.clock = clock.Realtime()
	}

	return cfg
}

type pubSubConn struct {
	cfg  pubSubConfig
	conn Conn

	subs, psubs map[string]bool
	pingTicker  *clock.Ticker
}

func (cfg pubSubConfig) new(conn Conn) PubSubConn {
	c := &pubSubConn{
		cfg:   cfg.withDefaults(),
		conn:  conn,
		subs:  map[string]bool{},
		psubs: map[string]bool{},
	}

	if c.cfg.PingInterval > 0 {
		c.pingTicker = c.cfg.clock.NewTicker(c.cfg.PingInterval)
	}

	return c
}

// New returns a PubSubConn instance using the given PubSubConfig.
func (cfg PubSubConfig) New(conn Conn) PubSubConn {
	return pubSubConfig{PubSubConfig: cfg}.new(conn)
}

func (c *pubSubConn) Close() error {
	if c.pingTicker != nil {
		c.pingTicker.Stop()
	}

	return c.conn.Close()
}

func (c *pubSubConn) cmd(cmd string, args ...string) resp.Marshaler {
	return Cmd(nil, cmd, args...).(resp.Marshaler)
}

func (c *pubSubConn) Subscribe(ctx context.Context, channels ...string) error {
	if err := c.conn.EncodeDecode(ctx, c.cmd("SUBSCRIBE", channels...), nil); err != nil {
		return err
	}

	for _, ch := range channels {
		c.subs[ch] = true
	}
	return nil
}

func (c *pubSubConn) Unsubscribe(ctx context.Context, channels ...string) error {
	if err := c.conn.EncodeDecode(ctx, c.cmd("UNSUBSCRIBE", channels...), nil); err != nil {
		return err
	}

	if len(channels) == 0 {
		c.subs = map[string]bool{}
	} else {
		for _, ch := range channels {
			delete(c.subs, ch)
		}
	}

	return nil
}

func (c *pubSubConn) PSubscribe(ctx context.Context, patterns ...string) error {
	if err := c.conn.EncodeDecode(ctx, c.cmd("PSUBSCRIBE", patterns...), nil); err != nil {
		return err
	}

	for _, p := range patterns {
		c.psubs[p] = true
	}
	return nil
}

func (c *pubSubConn) PUnsubscribe(ctx context.Context, patterns ...string) error {
	if err := c.conn.EncodeDecode(ctx, c.cmd("PUNSUBSCRIBE", patterns...), nil); err != nil {
		return err
	}

	if len(patterns) == 0 {
		c.psubs = map[string]bool{}
	} else {
		for _, p := range patterns {
			delete(c.psubs, p)
		}
	}

	return nil
}

func (c *pubSubConn) Ping(ctx context.Context) error {
	return c.conn.EncodeDecode(ctx, c.cmd("PING"), nil)
}

func (c *pubSubConn) testEvent(event string) {
	if c.cfg.testEventCh != nil {
		c.cfg.testEventCh <- event
		<-c.cfg.testEventCh
	}
}

var nullCtxCancel = context.CancelFunc(func() {})

const pubSubCtxWrapTimeout = 1 * time.Second

func (c *pubSubConn) wrapNextCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	if c.pingTicker == nil {
		return ctx, nullCtxCancel
	}

	return c.cfg.clock.TimeoutContext(ctx, pubSubCtxWrapTimeout)
}

func (c *pubSubConn) Next(ctx context.Context) (PubSubMessage, error) {
	for {
		c.testEvent("next-top")

		if c.pingTicker != nil {
			select {
			case <-c.pingTicker.C:
				if err := c.Ping(ctx); err != nil {
					return PubSubMessage{}, fmt.Errorf("calling PING internally: %w", err)
				}
				c.testEvent("pinged")
			default:
			}
		}

		// ctx has the potential to be wrapped so that it will have a 1 second
		// deadline, so that we can loop back up to check the pingTicker now and
		// then.
		innerCtx, cancel := c.wrapNextCtx(ctx)
		c.testEvent("wrapped-ctx")

		var msg PubSubMessage
		err := c.conn.EncodeDecode(innerCtx, nil, &msg)
		cancel()
		c.testEvent("decode-returned")

		if errors.Is(err, errNotPubSubMessage) {
			continue
		} else if ctxErr := ctx.Err(); ctxErr != nil {
			return msg, ctxErr
		} else if errors.Is(err, context.DeadlineExceeded) {
			continue
		} else if err != nil {
			return msg, err
		}

		if msg.Pattern != "" {
			if !c.psubs[msg.Pattern] {
				c.testEvent("skipped-pattern")
				continue
			}
		} else if !c.subs[msg.Channel] {
			c.testEvent("skipped-channel")
			continue
		}

		return msg, nil
	}
}

////////////////////////////////////////////////////////////////////////////////

// PersistentPubSubConnConfig is used to create a persistent PubSubConn with
// particular settings. All fields are optional, all methods are thread-safe.
type PersistentPubSubConnConfig struct {
	// Dialer is used to create new Conns.
	Dialer Dialer

	// PubSubConfig is used to create PubSubConns from the Conns created by
	// Dialer.
	PubSubConfig PubSubConfig

	// Trace contains callbacks that a persistent PubSubConn can use to trace
	// itself.
	//
	// All callbacks are blocking.
	Trace trace.PersistentPubSubTrace
}

type persistentPubSubConn struct {
	cfg  PersistentPubSubConnConfig
	dial func(context.Context) (Conn, error)

	subs, psubs map[string]bool
	conn        PubSubConn
}

// New is like PubSubConfig.New, but instead of taking in an existing Conn to
// wrap it will create its own using the network/address returned from the given
// callback.
//
// If the Conn is ever severed then the callback will be re-called, a new Conn
// will be created, and that Conn will be reset to the previous Conn's state.
//
// This is effectively a way to have a permanent PubSubConn established which
// supports subscribing/unsubscribing but without the hassle of implementing
// reconnect/re-subscribe logic.
func (cfg PersistentPubSubConnConfig) New(
	ctx context.Context,
	cb func() (network, addr string, err error),
) (
	PubSubConn, error,
) {
	p := &persistentPubSubConn{
		cfg: cfg,
		dial: func(ctx context.Context) (Conn, error) {
			var network, addr string
			var err error
			if cb != nil {
				// hidden feature, you don't have to give a callback if the
				// Dialer is just going to ignore its results anyway.
				if network, addr, err = cb(); err != nil {
					return nil, fmt.Errorf("calling PersistentPubSub callback to determine network/address: %w", err)
				}
			}
			return cfg.Dialer.Dial(ctx, network, addr)
		},

		subs:  map[string]bool{},
		psubs: map[string]bool{},
	}

	if err := p.refresh(ctx); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *persistentPubSubConn) refresh(ctx context.Context) error {
	if p.conn != nil {
		p.conn.Close()

		// sleep a bit in between closing and re-opening the connection. If the
		// server is unavailable we don't want the client to get stuck in a
		// tight loop creating/closing connections.
		time.Sleep(250 * time.Millisecond)
	}

	conn, err := p.dial(ctx)
	if err != nil {
		return err
	}

	p.conn = p.cfg.PubSubConfig.New(conn)

	mtos := func(m map[string]bool) []string {
		strs := make([]string, 0, len(m))
		for str := range m {
			strs = append(strs, str)
		}
		return strs
	}

	if len(p.subs) > 0 {
		if err := p.conn.Subscribe(ctx, mtos(p.subs)...); err != nil {
			return fmt.Errorf("recreating subscriptions: %w", err)
		}
	}

	if len(p.psubs) > 0 {
		if err := p.conn.PSubscribe(ctx, mtos(p.psubs)...); err != nil {
			return fmt.Errorf("recreating pattern subscriptions: %w", err)
		}
	}

	return nil
}

func (p *persistentPubSubConn) traceErr(err error) {
	if p.cfg.Trace.InternalError != nil {
		p.cfg.Trace.InternalError(trace.PersistentPubSubInternalError{
			Err: err,
		})
	}
}

func (p *persistentPubSubConn) do(ctx context.Context, fn func() error) error {
	var err error
	for {
		if errors.Is(err, context.Canceled) ||
			errors.Is(err, context.DeadlineExceeded) {
			return err

		} else if err != nil {
			p.traceErr(err)
			err = p.refresh(ctx)

		} else if err = fn(); err == nil {
			return nil
		}
	}
}

func (p *persistentPubSubConn) Subscribe(ctx context.Context, channels ...string) error {
	for _, ch := range channels {
		p.subs[ch] = true
	}

	return p.do(ctx, func() error { return p.conn.Subscribe(ctx, channels...) })
}

func (p *persistentPubSubConn) Unsubscribe(ctx context.Context, channels ...string) error {

	if len(channels) == 0 {
		p.subs = map[string]bool{}
	} else {
		for _, ch := range channels {
			delete(p.subs, ch)
		}
	}

	return p.do(ctx, func() error { return p.conn.Unsubscribe(ctx, channels...) })
}

func (p *persistentPubSubConn) PSubscribe(ctx context.Context, patterns ...string) error {
	for _, pat := range patterns {
		p.psubs[pat] = true
	}

	return p.do(ctx, func() error { return p.conn.PSubscribe(ctx, patterns...) })
}

func (p *persistentPubSubConn) PUnsubscribe(ctx context.Context, patterns ...string) error {

	if len(patterns) == 0 {
		p.psubs = map[string]bool{}
	} else {
		for _, pat := range patterns {
			delete(p.psubs, pat)
		}
	}

	return p.do(ctx, func() error { return p.conn.PUnsubscribe(ctx, patterns...) })
}

func (p *persistentPubSubConn) Ping(ctx context.Context) error {
	return p.do(ctx, func() error { return p.conn.Ping(ctx) })
}

func (p *persistentPubSubConn) Next(ctx context.Context) (PubSubMessage, error) {
	var msg PubSubMessage
	var err error
	err = p.do(ctx, func() error {
		msg, err = p.conn.Next(ctx)
		return err
	})
	return msg, err
}

func (p *persistentPubSubConn) Close() error {
	return p.conn.Close()
}
