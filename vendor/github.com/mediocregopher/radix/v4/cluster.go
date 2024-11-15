package radix

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"errors"

	"github.com/mediocregopher/radix/v4/internal/proc"
	"github.com/mediocregopher/radix/v4/resp"
	"github.com/mediocregopher/radix/v4/resp/resp3"
	"github.com/mediocregopher/radix/v4/trace"
)

// dedupe is used to deduplicate a function invocation, so if multiple
// go-routines call it at the same time only the first will actually run it, and
// the others will block until that one is done.
type dedupe struct {
	l sync.Mutex
	s *sync.Once
}

func newDedupe() *dedupe {
	return &dedupe{s: new(sync.Once)}
}

func (d *dedupe) do(fn func()) {
	d.l.Lock()
	s := d.s
	d.l.Unlock()

	s.Do(func() {
		fn()
		d.l.Lock()
		d.s = new(sync.Once)
		d.l.Unlock()
	})
}

////////////////////////////////////////////////////////////////////////////////

// ClusterConfig is used to create Cluster instances with particular settings.
// All fields are optional, all methods are thread-safe.
type ClusterConfig struct {
	// PoolConfig is used by Cluster to create Clients for redis instances in
	// the cluster set.
	//
	// If PoolConfig.CustomPool and PoolConfig.Dialer.CustomConn are unset
	// then all Conns created by Cluster will have the READONLY command
	// performed on them upon creation. For Conns to primary instances this will
	// have no effect, but for secondaries this will allow DoSecondary to
	// function properly.
	//
	// If PoolConfig.CustomPool or PoolConfig.Dialer.CustomConn are set then
	// READONLY must be called by whichever is set in order for DoSecondary to
	// work.
	PoolConfig PoolConfig

	// SyncEvery tells the Cluster to synchronize itself with the cluster's
	// topology at the given interval. On every synchronization Cluster will ask
	// the cluster for its topology and make/destroy its Clients as necessary.
	//
	// Defaults to 5 * time.Second. Set to -1 to disable.
	SyncEvery time.Duration

	// OnDownDelayActionsBy tells the Cluster to delay all commands by the given
	// duration while the cluster is seen to be in the CLUSTERDOWN state. This
	// allows fewer Actions to be affected by brief outages, e.g. during a
	// failover.
	//
	// Calls to Sync will not be delayed regardless of this option.
	//
	// Defaults to 100 * time.Millisecond. Set to -1 to disable.
	OnDownDelayActionsBy time.Duration

	// Trace contains callbacks that a Cluster can use to trace itself.
	//
	// All callbacks are blocking.
	Trace trace.ClusterTrace
}

func (cfg ClusterConfig) withDefaults() ClusterConfig {
	if cfg.SyncEvery < 0 {
		cfg.SyncEvery = 0
	} else if cfg.SyncEvery == 0 {
		cfg.SyncEvery = 5 * time.Second
	}
	if cfg.OnDownDelayActionsBy < 0 {
		cfg.OnDownDelayActionsBy = 0
	} else if cfg.OnDownDelayActionsBy == 0 {
		cfg.OnDownDelayActionsBy = 100 * time.Millisecond
	}

	if cfg.PoolConfig.CustomPool == nil &&
		cfg.PoolConfig.Dialer.CustomConn == nil {
		dialer := cfg.PoolConfig.Dialer
		cfg.PoolConfig.Dialer.CustomConn = func(ctx context.Context, network, addr string) (Conn, error) {
			conn, err := dialer.Dial(ctx, network, addr)
			if err != nil {
				return nil, err
			} else if err := conn.Do(ctx, Cmd(nil, "READONLY")); err != nil {
				conn.Close()
				return nil, err
			}
			return conn, nil
		}
	}

	return cfg
}

// Cluster is a MultiClient which contains all information about a redis cluster
// needed to interact with it, including a set of pools to each of its
// instances.
//
// All methods on Cluster are thread-safe.
//
// Cluster will automatically attempt to handle MOVED/ASK errors.
type Cluster struct {
	// Atomic fields must be at the beginning of the struct since they must be
	// correctly aligned or else access may cause panics on 32-bit architectures
	// See https://golang.org/pkg/sync/atomic/#pkg-note-BUG
	lastClusterdown int64 // unix timestamp in milliseconds, atomic

	proc *proc.Proc
	cfg  ClusterConfig

	// used to deduplicate calls to sync
	syncDedupe *dedupe

	// these fields are protected by proc's lock
	pools          map[string]Client
	primTopo, topo ClusterTopo
	secondaries    map[string]map[string]ClusterNode
}

var _ MultiClient = new(Cluster)

// New initializes and returns a Cluster instance using the ClusterConfig. It
// will try every address given until it finds a usable one.  From there it uses
// CLUSTER SLOTS to discover the cluster topology and make all the necessary
// connections.
func (cfg ClusterConfig) New(ctx context.Context, clusterAddrs []string) (*Cluster, error) {
	c := &Cluster{
		proc:       proc.New(),
		cfg:        cfg.withDefaults(),
		syncDedupe: newDedupe(),
		pools:      map[string]Client{},
	}

	var err error

	// make a pool to base the cluster on
	for _, addr := range clusterAddrs {

		var client Client

		if client, err = c.newClient(ctx, addr); err != nil {
			continue
		}

		c.pools[addr] = client
		break
	}

	if len(c.pools) == 0 {
		return nil, fmt.Errorf("could not connect to any redis instances, last error was: %w", err)
	}

	if err := c.Sync(ctx); err != nil {
		for _, p := range c.pools {
			p.Close()
		}
		return nil, err
	}

	if c.cfg.SyncEvery > 0 {
		c.proc.Run(func(ctx context.Context) { c.syncEvery(ctx, c.cfg.SyncEvery) })
	}

	return c, nil
}

func (c *Cluster) newClient(ctx context.Context, addr string) (Client, error) {
	return c.cfg.PoolConfig.New(ctx, "tcp", addr)
}

func (c *Cluster) err(err error) {
	if c.cfg.Trace.InternalError != nil {
		c.cfg.Trace.InternalError(trace.ClusterInternalError{
			Err: err,
		})
	}
}

func assertKeysSlot(keys []string) error {
	var ok bool
	var prevKey string
	var slot uint16
	for _, key := range keys {
		thisSlot := ClusterSlot([]byte(key))
		if !ok {
			ok = true
		} else if slot != thisSlot {
			return fmt.Errorf("keys %q and %q do not belong to the same slot", prevKey, key)
		}
		prevKey = key
		slot = thisSlot
	}
	return nil
}

// may return nil, nil if no pool for the addr.
func (c *Cluster) rpool(addr string) (client Client, err error) {
	err = c.proc.WithRLock(func() error {
		if addr == "" {
			for _, client = range c.pools {
				return nil
			}
			return errors.New("no Clients available")
		}
		client = c.pools[addr]
		return nil
	})
	return
}

// if addr is "" returns a random pool. If addr is given but there's no pool for
// it one will be created on-the-fly.
func (c *Cluster) pool(ctx context.Context, addr string) (Client, error) {
	p, err := c.rpool(addr)
	if p != nil || err != nil {
		return p, err
	}

	// if the pool isn't available make it on-the-fly. This behavior isn't
	// _great_, but theoretically the syncEvery process should clean up any
	// extraneous pools which aren't really needed

	// it's important that the cluster pool set isn't locked while this is
	// happening, because this could block for a while
	if p, err = c.newClient(ctx, addr); err != nil {
		return nil, err
	}

	// we've made a new pool, but we need to double-check someone else didn't
	// make one at the same time and add it in first. If they did, close this
	// one and return that one
	err = c.proc.WithLock(func() error {
		if p2, ok := c.pools[addr]; ok {
			p.Close()
			p = p2
			return nil
		}
		c.pools[addr] = p
		return nil
	})
	return p, err
}

// Topo returns the Cluster's topology as it currently knows it. See
// ClusterTopo's docs for more on its default order.
func (c *Cluster) Topo() ClusterTopo {
	var topo ClusterTopo
	_ = c.proc.WithRLock(func() error {
		topo = c.topo
		return nil
	})
	return topo
}

// Clients implements the method for the MultiClient interface.
func (c *Cluster) Clients() (map[string]ReplicaSet, error) {
	m := map[string]ReplicaSet{}
	err := c.proc.WithRLock(func() error {
		for _, primNode := range c.primTopo {
			primAddr := primNode.Addr
			primClient, ok := c.pools[primAddr]
			if !ok {
				return fmt.Errorf("no available Client for primary %q", primAddr)
			}
			rs := ReplicaSet{Primary: primClient}
			m[primAddr] = rs
		}

		for primAddr, secondaries := range c.secondaries {
			rs := m[primAddr]
			for secAddr := range secondaries {
				secClient, ok := c.pools[secAddr]
				if !ok {
					return fmt.Errorf("no available Client for secondary %q (secondary of %q)", secAddr, primAddr)
				}
				rs.Secondaries = append(rs.Secondaries, secClient)
			}
			m[primAddr] = rs
		}
		return nil
	})
	return m, err
}

func (c *Cluster) getTopo(ctx context.Context, p Client) (ClusterTopo, error) {
	var tt ClusterTopo
	err := p.Do(ctx, Cmd(&tt, "CLUSTER", "SLOTS"))
	if len(tt) == 0 && err == nil {
		//This will happen between when nodes starts coming up after cluster goes down and
		//Cluster swarm yet not ready using those nodes.
		err = errors.New("no cluster slots assigned")
	}
	return tt, err
}

// Sync will synchronize the Cluster with the actual cluster, making new pools
// to new instances and removing ones from instances no longer in the cluster.
// This will be called periodically automatically, but you can manually call it
// at any time as well.
func (c *Cluster) Sync(ctx context.Context) error {
	p, err := c.pool(ctx, "")
	if err != nil {
		return err
	}
	c.syncDedupe.do(func() {
		err = c.sync(ctx, p)
	})
	return err
}

func nodeInfoFromNode(node ClusterNode) trace.ClusterNodeInfo {
	return trace.ClusterNodeInfo{
		Addr:      node.Addr,
		Slots:     node.Slots,
		IsPrimary: node.SecondaryOfAddr == "",
	}
}

func (c *Cluster) traceTopoChanged(prevTopo ClusterTopo, newTopo ClusterTopo) {
	if c.cfg.Trace.TopoChanged == nil {
		return
	}

	var addedNodes []trace.ClusterNodeInfo
	var removedNodes []trace.ClusterNodeInfo
	var changedNodes []trace.ClusterNodeInfo

	prevTopoMap := prevTopo.Map()
	newTopoMap := newTopo.Map()

	for addr, newNode := range newTopoMap {
		if prevNode, ok := prevTopoMap[addr]; ok {
			// Check whether two nodes which have the same address changed its value or not
			if !reflect.DeepEqual(prevNode, newNode) {
				changedNodes = append(changedNodes, nodeInfoFromNode(newNode))
			}
			// No need to handle this address for finding removed nodes
			delete(prevTopoMap, addr)
		} else {
			// The node's address not found from prevTopo is newly added node
			addedNodes = append(addedNodes, nodeInfoFromNode(newNode))
		}
	}

	// Find removed nodes, prevTopoMap has reduced
	for addr, prevNode := range prevTopoMap {
		if _, ok := newTopoMap[addr]; !ok {
			removedNodes = append(removedNodes, nodeInfoFromNode(prevNode))
		}
	}

	// Callback when any changes detected
	if len(addedNodes) != 0 || len(removedNodes) != 0 || len(changedNodes) != 0 {
		c.cfg.Trace.TopoChanged(trace.ClusterTopoChanged{
			Added:   addedNodes,
			Removed: removedNodes,
			Changed: changedNodes,
		})
	}
}

// while this method is normally deduplicated by the Sync method's use of
// dedupe it is perfectly thread-safe on its own and can be used whenever.
func (c *Cluster) sync(ctx context.Context, p Client) error {
	tt, err := c.getTopo(ctx, p)
	if err != nil {
		return err
	}

	for _, t := range tt {
		// call pool just to ensure one exists for this addr
		if _, err := c.pool(ctx, t.Addr); err != nil {
			return fmt.Errorf("error creating client for %q: %w", t.Addr, err)
		}
	}

	c.traceTopoChanged(c.topo, tt)

	var toClose []Client
	err = c.proc.WithLock(func() error {
		c.topo = tt
		c.primTopo = tt.Primaries()

		c.secondaries = make(map[string]map[string]ClusterNode, len(c.primTopo))
		for _, node := range c.topo {
			if node.SecondaryOfAddr != "" {
				m := c.secondaries[node.SecondaryOfAddr]
				if m == nil {
					m = make(map[string]ClusterNode, len(c.topo)/len(c.primTopo))
					c.secondaries[node.SecondaryOfAddr] = m
				}
				m[node.Addr] = node
			}
		}

		tm := tt.Map()
		for addr, p := range c.pools {
			if _, ok := tm[addr]; !ok {
				toClose = append(toClose, p)
				delete(c.pools, addr)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, p := range toClose {
		p.Close()
	}
	return nil
}

func (c *Cluster) syncEvery(ctx context.Context, d time.Duration) {
	t := time.NewTicker(d)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			ctx, cancel := context.WithTimeout(ctx, d)
			err := c.Sync(ctx)
			cancel()
			if err != nil {
				c.err(fmt.Errorf("calling Sync internally: %w", err))
			}
		case <-c.proc.ClosedCh():
			return
		}
	}
}

func (c *Cluster) clientForKey(key string, random, secondary bool) (Client, string, error) {
	var addr string
	var client Client
	err := c.proc.WithRLock(func() error {
		var primAddr string
		if random {
			for _, primNode := range c.primTopo {
				primAddr = primNode.Addr
				break
			}
		} else {
			s := ClusterSlot([]byte(key))
		loop:
			for _, t := range c.primTopo {
				for _, slot := range t.Slots {
					if s >= slot[0] && s < slot[1] {
						primAddr = t.Addr
						break loop
					}
				}
			}
		}

		if primAddr == "" {
			return fmt.Errorf("could not find primary address for key %q", key)
		} else if secondary {
			for addr = range c.secondaries[primAddr] {
				break
			}
		}
		if addr == "" {
			addr = primAddr
		}

		client = c.pools[addr]
		return nil
	})
	if err != nil {
		return nil, "", err
	} else if client == nil {
		return nil, "", fmt.Errorf("no Client available for key %q (address:%q)", key, addr)
	}
	return client, addr, nil
}

type prefixAsking struct {
	marshal, unmarshalInto interface{}
}

func (a prefixAsking) MarshalRESP(w io.Writer, o *resp.Opts) error {
	if err := Cmd(nil, "ASKING").(resp.Marshaler).MarshalRESP(w, o); err != nil {
		return err
	}
	return resp3.Marshal(w, a.marshal, o)
}

func (a prefixAsking) UnmarshalRESP(br resp.BufferedReader, o *resp.Opts) error {
	if err := resp3.Unmarshal(br, nil, o); err != nil {
		return err
	}
	return resp3.Unmarshal(br, a.unmarshalInto, o)
}

type askConn struct {
	Conn
}

func (ac askConn) EncodeDecode(ctx context.Context, m, u interface{}) error {
	a := prefixAsking{m, u}
	return ac.Conn.EncodeDecode(ctx, a, a)
}

func (ac askConn) Do(ctx context.Context, a Action) error {
	return a.Perform(ctx, ac)
}

const doAttempts = 5

// Do performs an Action on a redis instance in the cluster, with the instance
// being determeined by the keys returned from the Action's Properties() method.
//
// This method handles MOVED and ASK errors automatically in most cases.
func (c *Cluster) Do(ctx context.Context, a Action) error {
	var client Client
	var err error
	var addr, key string
	keys := a.Properties().Keys
	if len(keys) == 0 {
		if client, addr, err = c.clientForKey("", true, false); err != nil {
			return err
		}
		// key will be ""
	} else if err := assertKeysSlot(keys); err != nil {
		return err
	} else {
		key = keys[0]
		if client, addr, err = c.clientForKey(key, false, false); err != nil {
			return err
		}
	}

	return c.doInner(clusterDoInnerParams{
		ctx:      ctx,
		action:   a,
		client:   client,
		addr:     addr,
		key:      key,
		attempts: doAttempts,
	})
}

// DoSecondary implements the method for the MultiClient interface. It will
// perform the Action on a random secondary for the affected keys, or the
// primary if no secondary is available.
//
// For DoSecondary to work, all connections must be created in read-only mode by
// using the READONLY command. See the PoolConfig field of ClusterConfig for
// more details.
func (c *Cluster) DoSecondary(ctx context.Context, a Action) error {
	var client Client
	var err error
	var addr, key string
	keys := a.Properties().Keys
	if len(keys) == 0 {
		if client, addr, err = c.clientForKey("", true, true); err != nil {
			return err
		}
		// key will be ""
	} else if err := assertKeysSlot(keys); err != nil {
		return err
	} else {
		key = keys[0]
		if client, addr, err = c.clientForKey(key, false, true); err != nil {
			return err
		}
	}

	return c.doInner(clusterDoInnerParams{
		ctx:      ctx,
		action:   a,
		client:   client,
		addr:     addr,
		key:      key,
		attempts: doAttempts,
	})
}

func (c *Cluster) getClusterDownSince() int64 {
	return atomic.LoadInt64(&c.lastClusterdown)
}

func (c *Cluster) setClusterDown(down bool) (changed bool) {
	// There is a race when calling this method concurrently when the cluster
	// healed after being down.
	//
	// If we have 2 goroutines, one that sends a command before the cluster
	// heals and once that sends a command after the cluster healed, both
	// goroutines will call this method, but with different values
	// (down == true and down == false).
	//
	// Since there is bi ordering between the two goroutines, it can happen
	// that the call to setClusterDown in the second goroutine runs before
	// the call in the first goroutine. In that case the state would be
	// changed from down to up by the second goroutine, as it should, only
	// for the first goroutine to set it back to down a few microseconds later.
	//
	// If this happens other commands will be needlessly delayed until
	// another goroutine sets the state to up again and we will trace two
	// unnecessary state transitions.
	//
	// We can not reliably avoid this race without more complex tracking of
	// previous states, which would be rather complex and possibly expensive.

	// Swapping values is expensive (on amd64, an uncontended swap can be 10x
	// slower than a load) and can easily become quite contended when we have
	// many goroutines trying to update the value concurrently, which would
	// slow it down even more.
	//
	// We avoid the overhead of swapping when not necessary by loading the
	// value first and checking if the value is already what we want it to be.
	//
	// Since atomic loads are fast (on amd64 an atomic load can be as fast as
	// a non-atomic load, and is perfectly scalable as long as there are no
	// writes to the same cache line), we can safely do this without adding
	// unnecessary extra latency.
	prevVal := atomic.LoadInt64(&c.lastClusterdown)

	var newVal int64
	if down {
		newVal = time.Now().UnixNano() / 1000 / 1000
		// Since the exact value is only used for delaying commands small
		// differences don't matter much and we can avoid many updates by
		// ignoring small differences (<5ms).
		if prevVal != 0 && newVal-prevVal < 5 {
			return false
		}
	} else {
		if prevVal == 0 {
			return false
		}
	}

	prevVal = atomic.SwapInt64(&c.lastClusterdown, newVal)

	changed = (prevVal == 0 && newVal != 0) || (prevVal != 0 && newVal == 0)

	if changed && c.cfg.Trace.StateChange != nil {
		c.cfg.Trace.StateChange(trace.ClusterStateChange{IsDown: down})
	}

	return changed
}

func (c *Cluster) traceRedirected(
	ctx context.Context,
	addr, key string,
	moved, ask bool,
	count int,
	final bool,
) {
	if c.cfg.Trace.Redirected != nil {
		c.cfg.Trace.Redirected(trace.ClusterRedirected{
			Context:       ctx,
			Addr:          addr,
			Key:           key,
			Moved:         moved,
			Ask:           ask,
			RedirectCount: count,
			Final:         final,
		})
	}
}

type clusterDoInnerParams struct {
	ctx       context.Context
	action    Action
	addr, key string
	client    Client
	ask       bool
	attempts  int
}

func (c *Cluster) doInner(params clusterDoInnerParams) error {
	if params.attempts <= 0 {
		return errors.New("cluster action redirected too many times")
	}

	if downSince := c.getClusterDownSince(); downSince > 0 && c.cfg.OnDownDelayActionsBy > 0 {
		// only wait when the last command was not too long, because
		// otherwise the chance is high that the cluster already healed
		elapsed := (time.Now().UnixNano() / 1000 / 1000) - downSince
		if elapsed < int64(c.cfg.OnDownDelayActionsBy/time.Millisecond) {
			time.Sleep(c.cfg.OnDownDelayActionsBy)
		}
	}

	thisA := params.action
	if params.ask {
		// the reason for doing this in such a round-about way, rather than just
		// pipelining an ASKING command with the action, is that this better
		// handles actions which can't be pipelined like EvalScript.
		thisA = WithConn(params.key, func(ctx context.Context, conn Conn) error {
			return askConn{conn}.Do(params.ctx, params.action)
		})
	}

	err := params.client.Do(params.ctx, thisA)
	if err == nil {
		c.setClusterDown(false)
		return nil
	}

	var respErr resp3.SimpleError
	if !errors.As(err, &respErr) {
		return err
	}

	msg := respErr.Error()

	clusterDown := strings.HasPrefix(msg, "CLUSTERDOWN ")
	clusterDownChanged := c.setClusterDown(clusterDown)
	if clusterDown && c.cfg.OnDownDelayActionsBy > 0 && clusterDownChanged {
		params.attempts--
		return c.doInner(params)
	}

	// if the error was a MOVED or ASK we can potentially retry
	moved := strings.HasPrefix(msg, "MOVED ")
	params.ask = strings.HasPrefix(msg, "ASK ")
	if !moved && !params.ask {
		return err
	}

	// if we get an ASK there's no need to do a sync quite yet, we can continue
	// normally. But MOVED always prompts a sync. In the section after this one
	// we figure out what address to use based on the returned error so the sync
	// isn't used _immediately_, but it still needs to happen.
	//
	// Also, even if the Action isn't a retryable Action we want a MOVED to
	// prompt a Sync
	if moved {
		if serr := c.Sync(params.ctx); serr != nil {
			return serr
		}
	}

	if !params.action.Properties().CanRetry {
		return err
	}

	msgParts := strings.Split(msg, " ")
	if len(msgParts) < 3 {
		return fmt.Errorf("malformed MOVED/ASK error %q", msg)
	}

	ogAddr := params.addr
	params.addr = msgParts[2]
	if params.client, err = c.pool(params.ctx, params.addr); err != nil {
		return err
	}

	c.traceRedirected(
		params.ctx,
		ogAddr, params.key,
		moved, params.ask,
		doAttempts-params.attempts+1,
		params.attempts <= 1,
	)
	params.attempts--
	return c.doInner(params)
}

// Close cleans up all goroutines spawned by Cluster and closes all of its
// Pools.
func (c *Cluster) Close() error {
	return c.proc.Close(func() error {
		var pErr error
		for _, p := range c.pools {
			if err := p.Close(); pErr == nil && err != nil {
				pErr = err
			}
		}
		return pErr
	})
}
