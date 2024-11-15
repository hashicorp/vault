package radix

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/mediocregopher/radix/v4/internal/proc"
	"github.com/mediocregopher/radix/v4/trace"
)

// SentinelConfig is used to create Sentinel instances with particular settings.
// All fields are optional, all methods are thread-safe.
type SentinelConfig struct {
	// PoolConfig is used by Sentinel to create Clients for redis instances in
	// the replica set.
	PoolConfig PoolConfig

	// SentinelDialer is the Dialer instance used to create Conns to sentinels.
	SentinelDialer Dialer

	// Trace contains callbacks that a Sentinel can use to trace itself.
	//
	// All callbacks are blocking.
	Trace trace.SentinelTrace
}

// Sentinel is a MultiClient which contains all information needed to interact
// with a redis replica set managed by redis sentinel, including a set of pools
// to each of its instances. All methods on Sentinel are thread-safe.
type Sentinel struct {
	proc      *proc.Proc
	cfg       SentinelConfig
	initAddrs []string
	name      string

	// these fields are protected by proc's lock
	primAddr      string
	clients       map[string]Client
	sentinelAddrs map[string]bool // the known sentinel addresses

	// We use a persistent PubSubConn here, so we don't need to do much after
	// initialization.
	pconn PubSubConn

	// only used by tests to ensure certain actions have happened before
	// continuing on during the test
	testEventCh chan string
}

var _ MultiClient = new(Sentinel)

// New creates and returns a *Sentinel instance using the SentinelConfig.
func (cfg SentinelConfig) New(ctx context.Context, primaryName string, sentinelAddrs []string) (*Sentinel, error) {
	addrs := map[string]bool{}
	for _, addr := range sentinelAddrs {
		addrs[addr] = true
	}

	sc := &Sentinel{
		proc:          proc.New(),
		cfg:           cfg,
		initAddrs:     sentinelAddrs,
		name:          primaryName,
		clients:       map[string]Client{},
		sentinelAddrs: addrs,
		testEventCh:   make(chan string, 1),
	}

	_, sc.cfg.SentinelDialer = parseRedisURL(sentinelAddrs[0], sc.cfg.SentinelDialer)

	// first thing is to retrieve the state and create a pool using the first
	// connectable connection. This connection is only used during
	// initialization, it gets closed right after
	{
		conn, err := sc.dialSentinel(ctx)
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		if err := sc.ensureSentinelAddrs(ctx, conn); err != nil {
			return nil, err
		} else if err := sc.ensureClients(ctx, conn); err != nil {
			return nil, err
		}
	}

	pconnCfg := PersistentPubSubConnConfig{
		Dialer: Dialer{
			CustomConn: func(ctx context.Context, _, _ string) (Conn, error) {
				return sc.dialSentinel(ctx)
			},
		},
	}

	// because we're using persistent these can't _really_ fail
	var err error
	sc.pconn, err = pconnCfg.New(ctx, nil)
	if err != nil {
		sc.Close()
		return nil, err
	}

	_ = sc.pconn.Subscribe(ctx, "switch-master")
	sc.proc.Run(sc.spin)
	return sc, nil
}

func (sc *Sentinel) err(err error) {
	if sc.cfg.Trace.InternalError != nil {
		sc.cfg.Trace.InternalError(trace.SentinelInternalError{
			Err: err,
		})
	}
}

func (sc *Sentinel) testEvent(event string) {
	select {
	case sc.testEventCh <- event:
	default:
	}
}

func (sc *Sentinel) dialSentinel(ctx context.Context) (conn Conn, err error) {
	err = sc.proc.WithRLock(func() error {
		for addr := range sc.sentinelAddrs {
			if conn, err = sc.cfg.SentinelDialer.Dial(ctx, "tcp", addr); err == nil {
				return nil
			}
		}

		// try the initAddrs as a last ditch, but don't return their error if
		// this doesn't work
		for _, addr := range sc.initAddrs {
			var initErr error
			if conn, initErr = sc.cfg.SentinelDialer.Dial(ctx, "tcp", addr); initErr == nil {
				return nil
			}
		}
		return err
	})
	return
}

// Do implements the method for the Client interface. It will perform the given
// Action on the current primary.
func (sc *Sentinel) Do(ctx context.Context, a Action) error {
	var client Client
	if err := sc.proc.WithRLock(func() error {
		client = sc.clients[sc.primAddr]
		return nil
	}); err != nil {
		return err
	}
	return client.Do(ctx, a)
}

// DoSecondary implements the method for the Client interface. It will perform
// the given Action on a random secondary, or the primary if no secondary is
// available.
//
// For DoSecondary to work, replicas must be configured with replica-read-only
// enabled, otherwise calls to DoSecondary may by rejected by the replica.
func (sc *Sentinel) DoSecondary(ctx context.Context, a Action) error {
	c, err := sc.client(ctx, "")
	if err != nil {
		return err
	}
	return c.Do(ctx, a)
}

// Clients implements the method for the MultiClient interface. The returned map
// will only ever have one key/value pair.
func (sc *Sentinel) Clients() (map[string]ReplicaSet, error) {
	m := map[string]ReplicaSet{}
	err := sc.proc.WithRLock(func() error {
		var rs ReplicaSet
		for addr, client := range sc.clients {
			if addr == sc.primAddr {
				rs.Primary = client
			} else {
				rs.Secondaries = append(rs.Secondaries, client)
			}
		}
		m[sc.primAddr] = rs
		return nil
	})
	return m, err
}

// SentinelAddrs returns the addresses of all known sentinels.
func (sc *Sentinel) SentinelAddrs() ([]string, error) {
	var sentAddrs []string
	err := sc.proc.WithRLock(func() error {
		sentAddrs = make([]string, 0, len(sc.sentinelAddrs))
		for addr := range sc.sentinelAddrs {
			sentAddrs = append(sentAddrs, addr)
		}
		return nil
	})
	return sentAddrs, err
}

func (sc *Sentinel) client(ctx context.Context, addr string) (Client, error) {
	var client Client
	err := sc.proc.WithRLock(func() error {
		if addr == "" {
			for addr, client = range sc.clients {
				if addr != sc.primAddr {
					break
				}
			}

			if client == nil {
				client = sc.clients[sc.primAddr]
			}
		} else {
			client = sc.clients[addr]
		}
		return nil
	})
	if err != nil {
		return nil, err
	} else if client != nil {
		return client, nil
	} else if addr == "" {
		return nil, errors.New("no Clients available")
	}

	// if client was nil but ok was true it means the address is a secondary but
	// a Client for it has never been created. Create one now and store it into
	// clients.
	newClient, err := sc.cfg.PoolConfig.New(ctx, "tcp", addr)
	if err != nil {
		return nil, err
	}

	// two routines might be requesting the same addr at the same time, and
	// both create the client. The second one needs to make sure it closes its
	// own pool when it sees the other got there first.
	err = sc.proc.WithLock(func() error {
		if client = sc.clients[addr]; client == nil {
			sc.clients[addr] = newClient
		}
		return nil
	})

	if client != nil || err != nil {
		newClient.Close()
		return client, err
	}

	return newClient, nil
}

// Close implements the method for the Client interface.
func (sc *Sentinel) Close() error {
	return sc.proc.Close(func() error {
		for _, client := range sc.clients {
			if client != nil {
				client.Close()
			}
		}
		return nil
	})
}

// cmd should be the command called which generated m.
func sentinelMtoAddr(m map[string]string, cmd string) (string, error) {
	if m["ip"] == "" || m["port"] == "" {
		return "", fmt.Errorf("malformed %q response: %#v", cmd, m)
	}
	return net.JoinHostPort(m["ip"], m["port"]), nil
}

// given a connection to a sentinel, ensures that the Clients currently being
// held agrees with what the sentinel thinks they should be.
func (sc *Sentinel) ensureClients(ctx context.Context, conn Conn) error {
	var primM map[string]string
	var secMM []map[string]string
	p := NewPipeline()
	p.Append(Cmd(&primM, "SENTINEL", "MASTER", sc.name))
	p.Append(Cmd(&secMM, "SENTINEL", "SLAVES", sc.name))
	if err := conn.Do(ctx, p); err != nil {
		return err
	}

	newPrimAddr, err := sentinelMtoAddr(primM, "SENTINEL MASTER")
	if err != nil {
		return err
	}

	newClients := map[string]Client{newPrimAddr: nil}
	for _, secM := range secMM {
		// means it's down with flag "s_down,slave"
		if secM["flags"] != "slave" {
			continue
		}
		newSecAddr, err := sentinelMtoAddr(secM, "SENTINEL SLAVES")
		if err != nil {
			return err
		}
		newClients[newSecAddr] = nil
	}

	// ensure all current clients exist
	newTraceNodes := map[string]trace.SentinelNodeInfo{}
	for addr := range newClients {
		client, err := sc.client(ctx, addr)
		if err != nil {
			return fmt.Errorf("error creating client for %q: %w", addr, err)
		}
		newClients[addr] = client
		newTraceNodes[addr] = trace.SentinelNodeInfo{
			Addr:      addr,
			IsPrimary: addr == newPrimAddr,
		}
	}

	var toClose []Client
	prevTraceNodes := map[string]trace.SentinelNodeInfo{}
	err = sc.proc.WithLock(func() error {

		// for each actual Client instance in sc.client, either move it over to
		// newClients (if the address is shared) or make sure it is closed
		for addr, client := range sc.clients {
			prevTraceNodes[addr] = trace.SentinelNodeInfo{
				Addr:      addr,
				IsPrimary: addr == sc.primAddr,
			}

			if _, ok := newClients[addr]; ok {
				newClients[addr] = client
			} else {
				toClose = append(toClose, client)
			}
		}

		sc.primAddr = newPrimAddr
		sc.clients = newClients

		return nil
	})
	if err != nil {
		return err
	}

	for _, client := range toClose {
		client.Close()
	}
	sc.traceTopoChanged(prevTraceNodes, newTraceNodes)
	return nil
}

func (sc *Sentinel) traceTopoChanged(prevTopo, newTopo map[string]trace.SentinelNodeInfo) {
	if sc.cfg.Trace.TopoChanged == nil {
		return
	}

	var added, removed, changed []trace.SentinelNodeInfo
	for addr, prevNodeInfo := range prevTopo {
		if newNodeInfo, ok := newTopo[addr]; !ok {
			removed = append(removed, prevNodeInfo)
		} else if newNodeInfo != prevNodeInfo {
			changed = append(changed, newNodeInfo)
		}
	}
	for addr, newNodeInfo := range newTopo {
		if _, ok := prevTopo[addr]; !ok {
			added = append(added, newNodeInfo)
		}
	}

	if len(added)+len(removed)+len(changed) == 0 {
		return
	}
	sc.cfg.Trace.TopoChanged(trace.SentinelTopoChanged{
		Added:   added,
		Removed: removed,
		Changed: changed,
	})
}

// annoyingly the SENTINEL SENTINELS <name> command doesn't return _this_
// sentinel instance, only the others it knows about for that primary.
func (sc *Sentinel) ensureSentinelAddrs(ctx context.Context, conn Conn) error {
	var mm []map[string]string
	err := conn.Do(ctx, Cmd(&mm, "SENTINEL", "SENTINELS", sc.name))
	if err != nil {
		return err
	}

	addrs := map[string]bool{conn.Addr().String(): true}
	for _, m := range mm {
		addrs[net.JoinHostPort(m["ip"], m["port"])] = true
	}

	return sc.proc.WithLock(func() error {
		sc.sentinelAddrs = addrs
		return nil
	})
}

func (sc *Sentinel) spin(ctx context.Context) {
	defer sc.pconn.Close()
	for {
		err := sc.innerSpin(ctx)

		// This also gets checked within innerSpin to short-circuit that, but
		// we also must check in here to short-circuit this. The error returned
		// doesn't really matter if the whole Sentinel is closing.
		select {
		case <-ctx.Done():
			return
		default:
		}

		if err != nil {
			sc.err(err)
			// sleep a second so we don't end up in a tight loop
			time.Sleep(1 * time.Second)
		}
	}
}

// makes connection to an address in sc.addrs and handles
// the sentinel until that connection goes bad.
//
// Things this handles:
//   - Listening for switch-master events (from pconn, which has reconnect logic
//     external to this package)
//   - Periodically re-ensuring that the list of sentinel addresses is up-to-date
//   - Periodically re-checking the current primary, in case the switch-master was
//     missed somehow
func (sc *Sentinel) innerSpin(ctx context.Context) error {
	conn, err := sc.dialSentinel(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	var switchMaster bool
	for {
		err := func() error {
			// putting this in an anonymous function is only slightly less ugly
			// than calling cancel in every if-error case.
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			if err := sc.ensureSentinelAddrs(ctx, conn); err != nil {
				return fmt.Errorf("retrieving addresses of sentinel instances: %w", err)
			} else if err := sc.ensureClients(ctx, conn); err != nil {
				return fmt.Errorf("creating clients based on sentinel addresses: %w", err)
			} else if err := sc.pconn.Ping(ctx); err != nil {
				return fmt.Errorf("calling PING on sentinel instance: %w", err)
			}
			return nil
		}()
		if err != nil {
			return err
		}

		// the tests want to know when the client state has been updated due to
		// a switch-master event
		if switchMaster {
			sc.testEvent("switch-master completed")
			switchMaster = false
		}

		innerCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		_, err = sc.pconn.Next(innerCtx)
		cancel()

		if err == nil {
			switchMaster = true

		} else if ctx.Err() != nil {
			return nil

		} else if innerCtx.Err() == nil {
			sc.err(fmt.Errorf("unexpected error from pubsub conn: %w", err))
		}
	}
}
