package radix

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/mediocregopher/radix/v4/internal/proc"
	"github.com/mediocregopher/radix/v4/resp"
	"github.com/mediocregopher/radix/v4/trace"
	"github.com/tilinna/clock"
)

// poolConn is a Conn which tracks the last net.Error which was seen either
// during an Encode call or a Decode call.
type poolConn struct {
	Conn

	// A channel to which critical network errors are written. A critical
	// network error is basically any non-application level error, e.g. a
	// timeout, disconnect, etc...
	lastIOErrCh chan error

	once sync.Once
}

func newPoolConn(c Conn) *poolConn {
	return &poolConn{
		Conn:        c,
		lastIOErrCh: make(chan error, 1),
	}
}

func (pc *poolConn) EncodeDecode(ctx context.Context, m, u interface{}) error {
	err := pc.Conn.EncodeDecode(ctx, m, u)
	if err != nil && !errors.As(err, new(resp.ErrConnUsable)) {
		select {
		case pc.lastIOErrCh <- err:
		default:
		}
	}
	return err
}

func (pc *poolConn) Do(ctx context.Context, a Action) error {
	return a.Perform(ctx, pc)
}

func (pc *poolConn) Close() error {
	return pc.Conn.Close()
}

////////////////////////////////////////////////////////////////////////////////

// PoolConfig is used to create Pool instances with particular settings. All
// fields are optional, all methods are thread-safe.
type PoolConfig struct {
	// CustomPool indicates that this callback should be used in place of New
	// when New is called. All behavior of New is superceded when this is set.
	CustomPool func(ctx context.Context, network, addr string) (Client, error)

	// Dialer is used by Pool to create new Conns to the Pool's redis instance.
	Dialer Dialer

	// Size indicates the number of Conns the Pool will attempt to maintain.
	//
	// Defaults to 4.
	Size int

	// PingInterval specifies the interval at which Pool will pick a random Conn
	// and call PING on it.
	//
	// If not given then the default value is calculated to be:
	//	5*seconds / Size.
	//
	// Can be set to -1 to disable periodic pings.
	PingInterval time.Duration

	// MinReconnectInterval describes the minimum amount of time the Pool will
	// wait between creating new Conns when previous Conns in the Pool have been
	// closed due to errors.
	//
	// Failure to create new Conns will result in the time between creation
	// attempts increasing exponentially, up to MaxReconnectInterval.
	// MinReconnectInterval and MaxReconnectInterval can be set to equal values
	// to disable exponential backoff.
	//
	// MinReconnectInterval defaults to 125 * time.Millisecond.
	// MaxReconnectInterval defaults to 4 * time.Second.
	MinReconnectInterval, MaxReconnectInterval time.Duration

	// Trace contains callbacks that a Pool can use to trace itself.
	//
	// All callbacks are blocking.
	Trace trace.PoolTrace
}

type poolRander interface {
	// Context is passed in as a way for a specific test operation to affect its
	// RNG in a non-racy way which also has no performance penalty for non-test
	// operations.
	Intn(ctx context.Context, n int) int
}

type primaryPoolRander struct {
	clock clock.Clock
}

var _ poolRander = new(primaryPoolRander)

func (r *primaryPoolRander) Intn(_ context.Context, n int) int {
	// using the real rand.Intn would incur an lock, and is approx 5-10x slower
	// than doing it this way. Since we don't need "real" randomness this will
	// do fine.
	return int(r.clock.Now().UnixNano() % int64(n))
}

func poolDefaultPingInterval(poolSize int) time.Duration {
	return 5 * time.Second / time.Duration(poolSize)
}

type poolConfig struct {
	PoolConfig

	clock clock.Clock
	rand  poolRander

	// these are only used for tests
	pingSyncCh      chan struct{}
	reconnectSyncCh chan struct{}
}

func (cfg poolConfig) withDefaults() poolConfig {
	if cfg.Size == 0 {
		cfg.Size = 4
	}
	if cfg.PingInterval == -1 {
		cfg.PingInterval = 0
	} else if cfg.PingInterval == 0 {
		cfg.PingInterval = poolDefaultPingInterval(cfg.Size)
	}
	if cfg.MinReconnectInterval == 0 {
		cfg.MinReconnectInterval = 125 * time.Millisecond
	}
	if cfg.MaxReconnectInterval == 0 {
		cfg.MaxReconnectInterval = 4 * time.Second
	}
	if cfg.clock == nil {
		cfg.clock = clock.Realtime()
	}
	if cfg.rand == nil {
		cfg.rand = &primaryPoolRander{cfg.clock}
	}

	cfg.Dialer = cfg.Dialer.withDefaults()
	return cfg
}

type pool struct {
	proc          *proc.Proc
	cfg           poolConfig
	network, addr string
	conns         *poolConnColl
	notifyCh      chan struct{}
	reconnectCh   chan struct{}
}

var _ Client = new(pool)

func (cfg poolConfig) new(ctx context.Context, network, addr string) (*pool, error) {
	cfg = cfg.withDefaults()
	p := &pool{
		proc:        proc.New(),
		cfg:         cfg,
		network:     network,
		addr:        addr,
		conns:       newPoolConnColl(cfg.Size),
		notifyCh:    make(chan struct{}, cfg.Size),
		reconnectCh: make(chan struct{}, cfg.Size),
	}

	// make one Conn synchronously to ensure there's actually a redis instance
	// present. The rest will be created asynchronously.
	pc, err := p.newConn(ctx, trace.PoolConnCreatedReasonInitialization)
	if err != nil {
		return nil, err
	}
	p.putConn(pc)

	p.proc.Run(p.runReconnect())
	p.proc.Run(func(ctx context.Context) {
		startTime := p.cfg.clock.Now()
		for i := 0; i < p.cfg.Size-1; i++ {
			pc, err := p.newConn(ctx, trace.PoolConnCreatedReasonInitialization)
			if err != nil {
				p.reconnectCh <- struct{}{}
				// if there was an error connecting to the instance than it
				// might need a little breathing room, redis can sometimes get
				// sad if too many connections are created simultaneously.
				p.cfg.clock.Sleep(100 * time.Millisecond)
				continue
			} else if !p.putConn(pc) {
				// Close was called.
				break
			}
		}
		p.traceInitCompleted(p.cfg.clock.Since(startTime))

	})

	if p.cfg.PingInterval > 0 {
		p.atIntervalDo(p.cfg.PingInterval, func(ctx context.Context) {
			ctx, cancel := p.cfg.clock.TimeoutContext(ctx, 2*time.Second)
			defer cancel()
			_ = p.Do(ctx, Cmd(nil, "PING"))
			if p.cfg.pingSyncCh != nil {
				<-p.cfg.pingSyncCh
			}
		})
	}
	return p, nil
}

// New creates and returns a pool instance using the PoolConfig.
func (cfg PoolConfig) New(ctx context.Context, network, addr string) (Client, error) {
	if cfg.CustomPool != nil {
		return cfg.CustomPool(ctx, network, addr)
	}
	p, err := poolConfig{PoolConfig: cfg}.new(ctx, network, addr)
	return p, err
}

func (p *pool) traceInitCompleted(elapsedTime time.Duration) {
	if p.cfg.Trace.InitCompleted != nil {
		p.cfg.Trace.InitCompleted(trace.PoolInitCompleted{
			PoolCommon:  p.traceCommon(),
			ElapsedTime: elapsedTime,
		})
	}
}

func (p *pool) traceCommon() trace.PoolCommon {
	return trace.PoolCommon{
		Network: p.network, Addr: p.addr,
		PoolSize: p.cfg.Size,
	}
}

func (p *pool) traceConnCreated(
	ctx context.Context,
	connectTime time.Duration,
	reason trace.PoolConnCreatedReason,
	err error,
) {
	if p.cfg.Trace.ConnCreated != nil {
		p.cfg.Trace.ConnCreated(trace.PoolConnCreated{
			PoolCommon:  p.traceCommon(),
			Context:     ctx,
			Reason:      reason,
			ConnectTime: connectTime,
			Err:         err,
		})
	}
}

func (p *pool) traceConnClosed(reason trace.PoolConnClosedReason) {
	if p.cfg.Trace.ConnClosed != nil {
		p.cfg.Trace.ConnClosed(trace.PoolConnClosed{
			PoolCommon: p.traceCommon(),
			Reason:     reason,
		})
	}
}

func (p *pool) newConn(ctx context.Context, reason trace.PoolConnCreatedReason) (*poolConn, error) {
	start := p.cfg.clock.Now()
	c, err := p.cfg.Dialer.Dial(ctx, p.network, p.addr)
	elapsed := p.cfg.clock.Since(start)
	p.traceConnCreated(ctx, elapsed, reason, err)
	if err != nil {
		return nil, err
	}
	pc := newPoolConn(c)
	return pc, nil
}

func (p *pool) atIntervalDo(d time.Duration, do func(context.Context)) {
	ticker := p.cfg.clock.NewTicker(d)
	p.proc.Run(func(ctx context.Context) {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				do(ctx)
			case <-ctx.Done():
				return
			}
		}
	})
}

func (p *pool) runReconnect() func(context.Context) {
	wait := p.cfg.MinReconnectInterval
	waitCh := p.cfg.clock.After(wait)

	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-waitCh:
			}

			select {
			case <-ctx.Done():
				return
			case <-p.reconnectCh:
			}

			ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
			pc, err := p.newConn(ctx, trace.PoolConnCreatedReasonReconnect)
			cancel()
			if err != nil {
				// the user can find out about the error via tracing.
				wait *= 2
				if wait > p.cfg.MaxReconnectInterval {
					wait = p.cfg.MaxReconnectInterval
				}
				p.reconnectCh <- struct{}{}
			} else {
				wait = p.cfg.MinReconnectInterval
				p.putConn(pc)
			}

			waitCh = p.cfg.clock.After(wait)

			if p.cfg.reconnectSyncCh != nil {
				<-p.cfg.reconnectSyncCh
			}
		}
	}
}

func (p *pool) discardConn(pc *poolConn, reason trace.PoolConnClosedReason) {

	// ensure that the discard logic for the conn only occurs once, specifically
	// buffering a message on reconnectCh.
	var ok bool
	pc.once.Do(func() { ok = true })
	if !ok {
		return
	}

	err := p.proc.WithLock(func() error {
		p.conns.remove(pc)
		return nil
	})
	if err != nil {
		return
	}

	pc.Close()
	p.traceConnClosed(reason)
	p.reconnectCh <- struct{}{}
}

func (p *pool) getConn(ctx context.Context) (*poolConn, error) {
	for {
		var pc *poolConn
		err := p.proc.WithLock(func() error {
			pc = p.conns.popBack() // will be nil if conns is empty
			return nil
		})
		if err != nil || pc != nil {
			return pc, err
		}

		select {
		case <-p.notifyCh:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// returns true if the connection was put back, false if it was closed and
// discarded.
func (p *pool) putConn(pc *poolConn) bool {
	err := p.proc.WithLock(func() error {
		p.conns.pushFront(pc)
		select {
		case p.notifyCh <- struct{}{}:
		default:
		}
		return nil
	})
	return err == nil
}

func (p *pool) useSharedConn(ctx context.Context, a Action) error {
	for {
		var pc *poolConn
		err := p.proc.WithRLock(func() error {
			l := p.conns.len()
			if l == 0 {
				return nil
			}
			i := p.cfg.rand.Intn(ctx, l)
			pc = p.conns.get(i)
			return pc.Do(ctx, a)
		})

		if pc != nil {
			select {
			case <-pc.lastIOErrCh:
				p.discardConn(pc, trace.PoolConnClosedReasonError)
			default:
			}
		}

		if pc != nil || err != nil {
			return err
		}

		select {
		case <-p.notifyCh:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Do implements the Do method of the Client interface by retrieving a Conn out
// of the pool, calling Perform on the given Action with it, and returning the
// Conn to the pool.
func (p *pool) Do(ctx context.Context, a Action) error {
	if a.Properties().CanShareConn {
		return p.useSharedConn(ctx, a)
	}

	pc, err := p.getConn(ctx)
	if err != nil {
		return err
	}

	err = pc.Do(ctx, a)
	if err != nil && !isRespErr(err) {
		// Non-shared conns are used for commands which might block. Therefore
		// any non-application errors result in closing the connection, because
		// it might still have some blocking command holding it up, and we don't
		// want to have other connections be blocked by it.
		p.discardConn(pc, trace.PoolConnClosedReasonError)
	} else {
		p.putConn(pc)
	}
	return err
}

// Addr implements the method for the Client interface.
func (p *pool) Addr() net.Addr {
	return rawAddr{network: p.network, addr: p.addr}
}

// Close implements the method for the Client interface.
func (p *pool) Close() error {
	return p.proc.Close(func() error {
		for {
			pc := p.conns.popBack()
			if pc == nil {
				break
			}
			pc.Close()
			p.traceConnClosed(trace.PoolConnClosedReasonPoolClosed)
		}
		return nil
	})
}
