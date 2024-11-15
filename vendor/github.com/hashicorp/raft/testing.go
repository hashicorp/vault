// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package raft

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-msgpack/v2/codec"
)

var userSnapshotErrorsOnNoData = true

// Return configurations optimized for in-memory
func inmemConfig(t testing.TB) *Config {
	conf := DefaultConfig()
	conf.HeartbeatTimeout = 50 * time.Millisecond
	conf.ElectionTimeout = 50 * time.Millisecond
	conf.LeaderLeaseTimeout = 50 * time.Millisecond
	conf.CommitTimeout = 5 * time.Millisecond
	conf.Logger = newTestLogger(t)
	return conf
}

// MockFSM is an implementation of the FSM interface, and just stores
// the logs sequentially.
//
// NOTE: This is exposed for middleware testing purposes and is not a stable API
type MockFSM struct {
	sync.Mutex
	logs           [][]byte
	configurations []Configuration
}

// NOTE: This is exposed for middleware testing purposes and is not a stable API
type MockFSMConfigStore struct {
	FSM
}

// NOTE: This is exposed for middleware testing purposes and is not a stable API
type WrappingFSM interface {
	Underlying() FSM
}

func getMockFSM(fsm FSM) *MockFSM {
	switch f := fsm.(type) {
	case *MockFSM:
		return f
	case *MockFSMConfigStore:
		return f.FSM.(*MockFSM)
	case WrappingFSM:
		return getMockFSM(f.Underlying())
	}

	return nil
}

// NOTE: This is exposed for middleware testing purposes and is not a stable API
type MockSnapshot struct {
	logs     [][]byte
	maxIndex int
}

var _ ConfigurationStore = (*MockFSMConfigStore)(nil)

// NOTE: This is exposed for middleware testing purposes and is not a stable API
func (m *MockFSM) Apply(log *Log) interface{} {
	m.Lock()
	defer m.Unlock()
	m.logs = append(m.logs, log.Data)
	return len(m.logs)
}

// NOTE: This is exposed for middleware testing purposes and is not a stable API
func (m *MockFSM) Snapshot() (FSMSnapshot, error) {
	m.Lock()
	defer m.Unlock()
	return &MockSnapshot{m.logs, len(m.logs)}, nil
}

// NOTE: This is exposed for middleware testing purposes and is not a stable API
func (m *MockFSM) Restore(inp io.ReadCloser) error {
	m.Lock()
	defer m.Unlock()
	defer inp.Close()
	hd := codec.MsgpackHandle{}
	dec := codec.NewDecoder(inp, &hd)

	m.logs = nil
	return dec.Decode(&m.logs)
}

// NOTE: This is exposed for middleware testing purposes and is not a stable API
func (m *MockFSM) Logs() [][]byte {
	m.Lock()
	defer m.Unlock()
	return m.logs
}

// NOTE: This is exposed for middleware testing purposes and is not a stable API
func (m *MockFSMConfigStore) StoreConfiguration(index uint64, config Configuration) {
	mm := m.FSM.(*MockFSM)
	mm.Lock()
	defer mm.Unlock()
	mm.configurations = append(mm.configurations, config)
}

// NOTE: This is exposed for middleware testing purposes and is not a stable API
func (m *MockSnapshot) Persist(sink SnapshotSink) error {
	hd := codec.MsgpackHandle{}
	enc := codec.NewEncoder(sink, &hd)
	if err := enc.Encode(m.logs[:m.maxIndex]); err != nil {
		sink.Cancel()
		return err
	}
	sink.Close()
	return nil
}

// NOTE: This is exposed for middleware testing purposes and is not a stable API
func (m *MockSnapshot) Release() {
}

// MockMonotonicLogStore is a LogStore wrapper for testing the
// MonotonicLogStore interface.
type MockMonotonicLogStore struct {
	s LogStore
}

// IsMonotonic implements the MonotonicLogStore interface.
func (m *MockMonotonicLogStore) IsMonotonic() bool {
	return true
}

// FirstIndex implements the LogStore interface.
func (m *MockMonotonicLogStore) FirstIndex() (uint64, error) {
	return m.s.FirstIndex()
}

// LastIndex implements the LogStore interface.
func (m *MockMonotonicLogStore) LastIndex() (uint64, error) {
	return m.s.LastIndex()
}

// GetLog implements the LogStore interface.
func (m *MockMonotonicLogStore) GetLog(index uint64, log *Log) error {
	return m.s.GetLog(index, log)
}

// StoreLog implements the LogStore interface.
func (m *MockMonotonicLogStore) StoreLog(log *Log) error {
	return m.s.StoreLog(log)
}

// StoreLogs implements the LogStore interface.
func (m *MockMonotonicLogStore) StoreLogs(logs []*Log) error {
	return m.s.StoreLogs(logs)
}

// DeleteRange implements the LogStore interface.
func (m *MockMonotonicLogStore) DeleteRange(min uint64, max uint64) error {
	return m.s.DeleteRange(min, max)
}

// This can be used as the destination for a logger and it'll
// map them into calls to testing.T.Log, so that you only see
// the logging for failed tests.
type testLoggerAdapter struct {
	tb     testing.TB
	prefix string
}

func (a *testLoggerAdapter) Write(d []byte) (int, error) {
	if d[len(d)-1] == '\n' {
		d = d[:len(d)-1]
	}
	if a.prefix != "" {
		l := a.prefix + ": " + string(d)
		a.tb.Log(l)
		return len(l), nil
	}

	a.tb.Log(string(d))
	return len(d), nil
}

func newTestLogger(tb testing.TB) hclog.Logger {
	return newTestLoggerWithPrefix(tb, "")
}

// newTestLoggerWithPrefix returns a Logger that can be used in tests. prefix
// will be added as the name of the logger.
//
// If tests are run with -v (verbose mode, or -json which implies verbose) the
// log output will go to stderr directly. If tests are run in regular "quiet"
// mode, logs will be sent to t.Log so that the logs only appear when a test
// fails.
//
// Be careful where this is used though - calling t.Log after the test completes
// causes a panic. This is common if you use it for a NetworkTransport for
// example and then close the transport at the end of the test because an error
// is logged after the test is complete.
func newTestLoggerWithPrefix(tb testing.TB, prefix string) hclog.Logger {
	if testing.Verbose() {
		return hclog.New(&hclog.LoggerOptions{Name: prefix, Level: hclog.Trace})
	}

	return hclog.New(&hclog.LoggerOptions{
		Name:   prefix,
		Output: &testLoggerAdapter{tb: tb, prefix: prefix},
	})
}

type cluster struct {
	dirs             []string
	stores           []*InmemStore
	fsms             []FSM
	snaps            []*FileSnapshotStore
	trans            []LoopbackTransport
	rafts            []*Raft
	t                *testing.T
	observationCh    chan Observation
	conf             *Config
	propagateTimeout time.Duration
	longstopTimeout  time.Duration
	logger           hclog.Logger
	startTime        time.Time

	failedLock sync.Mutex
	failedCh   chan struct{}
	failed     bool
}

func (c *cluster) Merge(other *cluster) {
	c.dirs = append(c.dirs, other.dirs...)
	c.stores = append(c.stores, other.stores...)
	c.fsms = append(c.fsms, other.fsms...)
	c.snaps = append(c.snaps, other.snaps...)
	c.trans = append(c.trans, other.trans...)
	c.rafts = append(c.rafts, other.rafts...)
}

func (c *cluster) RemoveServer(id ServerID) {
	for i, n := range c.rafts {
		if n.localID == id {
			c.rafts = append(c.rafts[:i], c.rafts[i+1:]...)
			return
		}
	}
}

// notifyFailed will close the failed channel which can signal the goroutine
// running the test that another goroutine has detected a failure in order to
// terminate the test.
func (c *cluster) notifyFailed() {
	c.failedLock.Lock()
	defer c.failedLock.Unlock()
	if !c.failed {
		c.failed = true
		close(c.failedCh)
	}
}

// Failf provides a logging function that fails the tests, prints the output
// with microseconds, and does not mysteriously eat the string. This can be
// safely called from goroutines but won't immediately halt the test. The
// failedCh will be closed to allow blocking functions in the main thread to
// detect the failure and react. Note that you should arrange for the main
// thread to block until all goroutines have completed in order to reliably
// fail tests using this function.
func (c *cluster) Failf(format string, args ...interface{}) {
	c.logger.Error(fmt.Sprintf(format, args...))
	c.t.Fail()
	c.notifyFailed()
}

// FailNowf provides a logging function that fails the tests, prints the output
// with microseconds, and does not mysteriously eat the string. FailNowf must be
// called from the goroutine running the test or benchmark function, not from
// other goroutines created during the test. Calling FailNowf does not stop
// those other goroutines.
func (c *cluster) FailNowf(format string, args ...interface{}) {
	c.t.Helper()
	c.t.Fatalf(format, args...)
}

// Close shuts down the cluster and cleans up.
func (c *cluster) Close() {
	var futures []Future
	for _, r := range c.rafts {
		futures = append(futures, r.Shutdown())
	}

	// Wait for shutdown
	limit := time.AfterFunc(c.longstopTimeout, func() {
		// We can't FailNowf here, and c.Failf won't do anything if we
		// hang, so panic.
		panic("timed out waiting for shutdown")
	})
	defer limit.Stop()

	for _, f := range futures {
		if err := f.Error(); err != nil {
			c.t.Fatalf("shutdown future err: %v", err)
		}
	}

	for _, d := range c.dirs {
		os.RemoveAll(d)
	}
}

// WaitEventChan returns a channel which will signal if an observation is made
// or a timeout occurs. It is possible to set a filter to look for specific
// observations. Setting timeout to 0 means that it will wait forever until a
// non-filtered observation is made.
func (c *cluster) WaitEventChan(ctx context.Context, filter FilterFn) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			case o, ok := <-c.observationCh:
				if !ok || filter == nil || filter(&o) {
					return
				}
			}
		}
	}()
	return ch
}

// WaitEvent waits until an observation is made, a timeout occurs, or a test
// failure is signaled. It is possible to set a filter to look for specific
// observations. Setting timeout to 0 means that it will wait forever until a
// non-filtered observation is made or a test failure is signaled.
func (c *cluster) WaitEvent(filter FilterFn, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	eventCh := c.WaitEventChan(ctx, filter)
	select {
	case <-c.failedCh:
		c.t.FailNow()
	case <-eventCh:
	}
}

// WaitForReplication blocks until every FSM in the cluster has the given
// length, or the long sanity check timeout expires.
func (c *cluster) WaitForReplication(fsmLength int) {
	limitCh := time.After(c.longstopTimeout)

CHECK:
	for {
		ctx, cancel := context.WithTimeout(context.Background(), c.conf.CommitTimeout)
		defer cancel()
		ch := c.WaitEventChan(ctx, nil)
		select {
		case <-c.failedCh:
			c.t.FailNow()

		case <-limitCh:
			c.t.Fatalf("timeout waiting for replication")

		case <-ch:
			for _, fsmRaw := range c.fsms {
				fsm := getMockFSM(fsmRaw)
				fsm.Lock()
				num := len(fsm.logs)
				fsm.Unlock()
				if num != fsmLength {
					continue CHECK
				}
			}
			return
		}
	}
}

// pollState takes a snapshot of the state of the cluster. This might not be
// stable, so use GetInState() to apply some additional checks when waiting
// for the cluster to achieve a particular state.
func (c *cluster) pollState(s RaftState) ([]*Raft, uint64) {
	var highestTerm uint64
	in := make([]*Raft, 0, 1)
	for _, r := range c.rafts {
		if r.State() == s {
			in = append(in, r)
		}
		term := r.getCurrentTerm()
		if term > highestTerm {
			highestTerm = term
		}
	}
	return in, highestTerm
}

// GetInState polls the state of the cluster and attempts to identify when it has
// settled into the given state.
func (c *cluster) GetInState(s RaftState) []*Raft {
	c.logger.Info("starting stability test", "raft-state", s)
	limitCh := time.After(c.longstopTimeout)

	// An election should complete after 2 * max(HeartbeatTimeout, ElectionTimeout)
	// because of the randomised timer expiring in 1 x interval ... 2 x interval.
	// We add a bit for propagation delay. If the election fails (e.g. because
	// two elections start at once), we will have got something through our
	// observer channel indicating a different state (i.e. one of the nodes
	// will have moved to candidate state) which will reset the timer.
	//
	// Because of an implementation peculiarity, it can actually be 3 x timeout.
	timeout := c.conf.HeartbeatTimeout
	if timeout < c.conf.ElectionTimeout {
		timeout = c.conf.ElectionTimeout
	}
	timeout = 2*timeout + c.conf.CommitTimeout
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// Wait until we have a stable instate slice. Each time we see an
	// observation a state has changed, recheck it and if it has changed,
	// restart the timer.
	pollStartTime := time.Now()
	for {
		_, highestTerm := c.pollState(s)
		inStateTime := time.Now()

		// Sometimes this routine is called very early on before the
		// rafts have started up. We then timeout even though no one has
		// even started an election. So if the highest term in use is
		// zero, we know there are no raft processes that have yet issued
		// a RequestVote, and we set a long time out. This is fixed when
		// we hear the first RequestVote, at which point we reset the
		// timer.
		if highestTerm == 0 {
			timer.Reset(c.longstopTimeout)
		} else {
			timer.Reset(timeout)
		}

		// Filter will wake up whenever we observe a RequestVote.
		filter := func(ob *Observation) bool {
			switch ob.Data.(type) {
			case RaftState:
				return true
			case RequestVoteRequest:
				return true
			default:
				return false
			}
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		eventCh := c.WaitEventChan(ctx, filter)
		select {
		case <-c.failedCh:
			c.t.FailNow()

		case <-limitCh:
			c.t.Fatalf("timeout waiting for stable %s state", s)

		case <-eventCh:
			c.logger.Debug("resetting stability timeout")

		case t, ok := <-timer.C:
			if !ok {
				c.t.Fatalf("timer channel errored")
			}

			inState, highestTerm := c.pollState(s)
			c.logger.Info(fmt.Sprintf("stable state for %s reached at %s (%d nodes), highestTerm is %d, %s from start of poll, %s from cluster start. Timeout at %s, %s after stability",
				s, inStateTime, len(inState), highestTerm, inStateTime.Sub(pollStartTime), inStateTime.Sub(c.startTime), t, t.Sub(inStateTime)))
			return inState
		}
	}
}

// Leader waits for the cluster to elect a leader and stay in a stable state.
func (c *cluster) Leader() *Raft {
	c.t.Helper()
	leaders := c.GetInState(Leader)
	if len(leaders) != 1 {
		c.t.Fatalf("expected one leader: %v", leaders)
	}
	return leaders[0]
}

// Followers waits for the cluster to have N-1 followers and stay in a stable
// state.
func (c *cluster) Followers() []*Raft {
	expFollowers := len(c.rafts) - 1
	return c.WaitForFollowers(expFollowers)
}

// WaitForFollowers waits for the cluster to have a given number of followers and stay in a stable
// state.
func (c *cluster) WaitForFollowers(expFollowers int) []*Raft {
	followers := c.GetInState(Follower)
	if len(followers) != expFollowers {
		c.t.Fatalf("timeout waiting for %d followers (followers are %v)", expFollowers, followers)
	}
	return followers
}

// FullyConnect connects all the transports together.
func (c *cluster) FullyConnect() {
	c.logger.Debug("fully connecting")
	for i, t1 := range c.trans {
		for j, t2 := range c.trans {
			if i != j {
				t1.Connect(t2.LocalAddr(), t2)
				t2.Connect(t1.LocalAddr(), t1)
			}
		}
	}
}

// Disconnect disconnects all transports from the given address.
func (c *cluster) Disconnect(a ServerAddress) {
	c.logger.Debug("disconnecting", "address", a)
	for _, t := range c.trans {
		if t.LocalAddr() == a {
			t.DisconnectAll()
		} else {
			t.Disconnect(a)
		}
	}
}

// Partition keeps the given list of addresses connected but isolates them
// from the other members of the cluster.
func (c *cluster) Partition(far []ServerAddress) {
	c.logger.Debug("partitioning", "addresses", far)

	// Gather the set of nodes on the "near" side of the partition (we
	// will call the supplied list of nodes the "far" side).
	near := make(map[ServerAddress]struct{})
OUTER:
	for _, t := range c.trans {
		l := t.LocalAddr()
		for _, a := range far {
			if l == a {
				continue OUTER
			}
		}
		near[l] = struct{}{}
	}

	// Now fixup all the connections. The near side will be separated from
	// the far side, and vice-versa.
	for _, t := range c.trans {
		l := t.LocalAddr()
		if _, ok := near[l]; ok {
			for _, a := range far {
				t.Disconnect(a)
			}
		} else {
			for a := range near {
				t.Disconnect(a)
			}
		}
	}
}

// IndexOf returns the index of the given raft instance.
func (c *cluster) IndexOf(r *Raft) int {
	for i, n := range c.rafts {
		if n == r {
			return i
		}
	}
	return -1
}

// EnsureLeader checks that ALL the nodes think the leader is the given expected
// leader.
func (c *cluster) EnsureLeader(t *testing.T, expect ServerAddress) {
	// We assume c.Leader() has been called already; now check all the rafts
	// think the leader is correct
	fail := false
	for _, r := range c.rafts {
		leaderAddr, _ := r.LeaderWithID()

		if leaderAddr != expect {
			if leaderAddr == "" {
				leaderAddr = "[none]"
			}
			if expect == "" {
				c.logger.Error("peer sees incorrect leader", "peer", r, "leader", leaderAddr, "expected-leader", "[none]")
			} else {
				c.logger.Error("peer sees incorrect leader", "peer", r, "leader", leaderAddr, "expected-leader", expect)
			}
			fail = true
		}
	}
	if fail {
		t.Fatalf("at least one peer has the wrong notion of leader")
	}
}

// EnsureSame makes sure all the FSMs have the same contents.
func (c *cluster) EnsureSame(t *testing.T) {
	limit := time.Now().Add(c.longstopTimeout)
	first := getMockFSM(c.fsms[0])

CHECK:
	first.Lock()
	for i, fsmRaw := range c.fsms {
		fsm := getMockFSM(fsmRaw)
		if i == 0 {
			continue
		}
		fsm.Lock()

		if len(first.logs) != len(fsm.logs) {
			fsm.Unlock()
			if time.Now().After(limit) {
				t.Fatalf("FSM log length mismatch: %d %d",
					len(first.logs), len(fsm.logs))
			} else {
				goto WAIT
			}
		}

		for idx := 0; idx < len(first.logs); idx++ {
			if bytes.Compare(first.logs[idx], fsm.logs[idx]) != 0 {
				fsm.Unlock()
				if time.Now().After(limit) {
					t.Fatalf("FSM log mismatch at index %d", idx)
				} else {
					goto WAIT
				}
			}
		}
		if len(first.configurations) != len(fsm.configurations) {
			fsm.Unlock()
			if time.Now().After(limit) {
				t.Fatalf("FSM configuration length mismatch: %d %d",
					len(first.logs), len(fsm.logs))
			} else {
				goto WAIT
			}
		}

		for idx := 0; idx < len(first.configurations); idx++ {
			if !reflect.DeepEqual(first.configurations[idx], fsm.configurations[idx]) {
				fsm.Unlock()
				if time.Now().After(limit) {
					t.Fatalf("FSM configuration mismatch at index %d: %v, %v", idx, first.configurations[idx], fsm.configurations[idx])
				} else {
					goto WAIT
				}
			}
		}
		fsm.Unlock()
	}

	first.Unlock()
	return

WAIT:
	first.Unlock()
	c.WaitEvent(nil, c.conf.CommitTimeout)
	goto CHECK
}

// getConfiguration returns the configuration of the given Raft instance, or
// fails the test if there's an error
func (c *cluster) getConfiguration(r *Raft) Configuration {
	future := r.GetConfiguration()
	if err := future.Error(); err != nil {
		c.t.Fatalf("failed to get configuration: %v", err)
		return Configuration{}
	}

	return future.Configuration()
}

// EnsureSamePeers makes sure all the rafts have the same set of peers.
func (c *cluster) EnsureSamePeers(t *testing.T) {
	limit := time.Now().Add(c.longstopTimeout)
	peerSet := c.getConfiguration(c.rafts[0])

CHECK:
	for i, raft := range c.rafts {
		if i == 0 {
			continue
		}

		otherSet := c.getConfiguration(raft)
		if !reflect.DeepEqual(peerSet, otherSet) {
			if time.Now().After(limit) {
				t.Fatalf("peer mismatch: %+v %+v", peerSet, otherSet)
			} else {
				goto WAIT
			}
		}
	}
	return

WAIT:
	c.WaitEvent(nil, c.conf.CommitTimeout)
	goto CHECK
}

// NOTE: This is exposed for middleware testing purposes and is not a stable API
type MakeClusterOpts struct {
	Peers           int
	Bootstrap       bool
	Conf            *Config
	ConfigStoreFSM  bool
	MakeFSMFunc     func() FSM
	LongstopTimeout time.Duration
	MonotonicLogs   bool
}

// makeCluster will return a cluster with the given config and number of peers.
// If bootstrap is true, the servers will know about each other before starting,
// otherwise their transports will be wired up but they won't yet have configured
// each other.
func makeCluster(t *testing.T, opts *MakeClusterOpts) *cluster {
	if opts.Conf == nil {
		opts.Conf = inmemConfig(t)
	}

	c := &cluster{
		observationCh: make(chan Observation, 1024),
		conf:          opts.Conf,
		// Propagation takes a maximum of 2 heartbeat timeouts (time to
		// get a new heartbeat that would cause a commit) plus a bit.
		propagateTimeout: opts.Conf.HeartbeatTimeout*2 + opts.Conf.CommitTimeout,
		longstopTimeout:  5 * time.Second,
		logger:           newTestLoggerWithPrefix(t, "cluster"),
		failedCh:         make(chan struct{}),
	}
	if opts.LongstopTimeout > 0 {
		c.longstopTimeout = opts.LongstopTimeout
	}

	c.t = t
	var configuration Configuration

	// Setup the stores and transports
	for i := 0; i < opts.Peers; i++ {
		dir, err := os.MkdirTemp("", "raft")
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		store := NewInmemStore()
		c.dirs = append(c.dirs, dir)
		c.stores = append(c.stores, store)
		if opts.ConfigStoreFSM {
			c.fsms = append(c.fsms, &MockFSMConfigStore{
				FSM: &MockFSM{},
			})
		} else {
			var fsm FSM
			if opts.MakeFSMFunc != nil {
				fsm = opts.MakeFSMFunc()
			} else {
				fsm = &MockFSM{}
			}
			c.fsms = append(c.fsms, fsm)
		}

		dir2, snap := FileSnapTest(t)
		c.dirs = append(c.dirs, dir2)
		c.snaps = append(c.snaps, snap)

		addr, trans := NewInmemTransport("")
		c.trans = append(c.trans, trans)
		localID := ServerID(fmt.Sprintf("server-%s", addr))
		if opts.Conf.ProtocolVersion < 3 {
			localID = ServerID(addr)
		}
		configuration.Servers = append(configuration.Servers, Server{
			Suffrage: Voter,
			ID:       localID,
			Address:  addr,
		})
	}

	// Wire the transports together
	c.FullyConnect()

	// Create all the rafts
	c.startTime = time.Now()
	for i := 0; i < opts.Peers; i++ {
		var logs LogStore
		logs = c.stores[i]
		store := c.stores[i]
		snap := c.snaps[i]
		trans := c.trans[i]

		if opts.MonotonicLogs {
			logs = &MockMonotonicLogStore{s: logs}
		}

		peerConf := opts.Conf
		peerConf.LocalID = configuration.Servers[i].ID
		peerConf.Logger = newTestLoggerWithPrefix(t, string(configuration.Servers[i].ID))

		if opts.Bootstrap {
			err := BootstrapCluster(peerConf, logs, store, snap, trans, configuration)
			if err != nil {
				t.Fatalf("BootstrapCluster failed: %v", err)
			}
		}

		raft, err := NewRaft(peerConf, c.fsms[i], logs, store, snap, trans)
		if err != nil {
			t.Fatalf("NewRaft failed: %v", err)
		}

		raft.RegisterObserver(NewObserver(c.observationCh, false, nil))
		if err != nil {
			t.Fatalf("RegisterObserver failed: %v", err)
		}
		c.rafts = append(c.rafts, raft)
	}

	return c
}

// NOTE: This is exposed for middleware testing purposes and is not a stable API
func MakeCluster(n int, t *testing.T, conf *Config) *cluster {
	return makeCluster(t, &MakeClusterOpts{
		Peers:     n,
		Bootstrap: true,
		Conf:      conf,
	})
}

// NOTE: This is exposed for middleware testing purposes and is not a stable API
func MakeClusterNoBootstrap(n int, t *testing.T, conf *Config) *cluster {
	return makeCluster(t, &MakeClusterOpts{
		Peers: n,
		Conf:  conf,
	})
}

// NOTE: This is exposed for middleware testing purposes and is not a stable API
func MakeClusterCustom(t *testing.T, opts *MakeClusterOpts) *cluster {
	return makeCluster(t, opts)
}

// NOTE: This is exposed for middleware testing purposes and is not a stable API
func FileSnapTest(t *testing.T) (string, *FileSnapshotStore) {
	// Create a test dir
	dir, err := os.MkdirTemp("", "raft")
	if err != nil {
		t.Fatalf("err: %v ", err)
	}

	snap, err := NewFileSnapshotStoreWithLogger(dir, 3, newTestLogger(t))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	snap.noSync = true
	return dir, snap
}
