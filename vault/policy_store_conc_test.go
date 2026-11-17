// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"container/list"
	"context"
	"maps"
	"math/rand"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/go-secure-stdlib/base62"
	"github.com/hashicorp/vault/internalshared/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/petermattis/goid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ policyLocker = &lockOwnerTracker{}

// lockOwnerTracker wraps a policyLocker and records which goroutine holds
// each key.
type lockOwnerTracker struct {
	ownerGoidByKey sync.Map
	underlying     policyLocker
}

func (l *lockOwnerTracker) Lock(key string) {
	l.underlying.Lock(key)
	l.ownerGoidByKey.Store(key, goid.Get())
}

func (l *lockOwnerTracker) Unlock(key string) {
	l.ownerGoidByKey.Delete(key)
	l.underlying.Unlock(key)
}

var _ logical.Storage = &storageOperationListener{}

// storageOperationListener wraps a logical.Storage and calls beforeOp before
// each Get, Put and Delete.
type storageOperationListener struct {
	underlying logical.Storage
	beforeOp   func(key string)
}

func (s *storageOperationListener) List(ctx context.Context, prefix string) ([]string, error) {
	return s.underlying.List(ctx, prefix)
}

func (s *storageOperationListener) Get(ctx context.Context, key string) (*logical.StorageEntry, error) {
	s.beforeOp(key)
	return s.underlying.Get(ctx, key)
}

func (s *storageOperationListener) Put(ctx context.Context, entry *logical.StorageEntry) error {
	s.beforeOp(entry.Key)
	return s.underlying.Put(ctx, entry)
}

func (s *storageOperationListener) Delete(ctx context.Context, key string) error {
	s.beforeOp(key)
	return s.underlying.Delete(ctx, key)
}

// mockConcurrencyCheckedCore creates a test cluster whose policy store
// asserts that every storage operation on a policy happens while the calling
// goroutine holds that policy's lock.
func mockConcurrencyCheckedCore(tb testing.TB) *PolicyStore {
	ps := NewTestCluster(tb, &CoreConfig{}, nil).Cores[0].Core.policyStore

	ownerTracker := &lockOwnerTracker{underlying: ps.policyLocks}
	ps.policyLocks = ownerTracker

	beforeStorageOp := func(key string) {
		cacheKey := ps.cacheKey(namespace.RootNamespace, key)
		owner, ok := ownerTracker.ownerGoidByKey.Load(cacheKey)
		assert.True(tb, ok, "storage operation on %q without holding its lock", key)
		assert.Equal(tb, goid.Get(), owner, "storage operation on %q while another goroutine holds its lock", key)
	}
	for _, bv := range []*BarrierView{ps.aclView, ps.rgpView, ps.egpView} {
		bv.storage = logical.NewStorageView(&storageOperationListener{
			underlying: bv.storage,
			beforeOp:   beforeStorageOp,
		}, "")
	}
	return ps
}

// policyState represents the state of a policy in the policy store. Possible values:
// - if the policy does not exist: { exists: false, rawContents: "" }
// - if the policy exists:         { exists: true,  rawContents: "<contents of the policy>" }
type policyState struct {
	exists      bool
	rawContents string
}

// policyStateSet represents the set of possible states of a policy in the policy store.
type policyStateSet struct {
	sync.Mutex
	// ordered list of [*policyState]s writes that have finished
	// states that are known to be invalid are pruned
	committed list.List
	// set of [*policyState]s that are currently being written
	ongoingWrites map[*policyState]struct{}
	// set of read operations on this policy
	ongoingReads map[*readHandle]struct{}
}

func newPolicyStateSet() *policyStateSet {
	set := &policyStateSet{
		ongoingWrites: map[*policyState]struct{}{},
		ongoingReads:  map[*readHandle]struct{}{},
	}
	// Add initial empty policy state
	set.committed.PushFront(&policyState{})
	return set
}

// policyStateEstimator tracks the possible states of policies in the policy store, using the following rules:
//   - Before initiating a write, the new policy state is added to the "ongoing" set.
//   - After the write has completed, the new state is moved from the "ongoing" set to the "committed" list.
//     Also, any states that were in the "committed" list before the write was initiated can be considered
//     overwritten so we discard them.
//   - At any point, the possible states of a policy are the union of the "ongoing" set and the "committed" list.
//   - Reads on the policy store can take some time to complete, so they also need to subscribe to any new writes that
//     happen during the read.
type policyStateEstimator struct {
	statesByKey map[string]*policyStateSet
}

// beginWrite should be called before starting a write operation.
func (p *policyStateEstimator) beginWrite(key string, newState policyState) *writeHandle {
	ss := p.statesByKey[key]
	ss.Lock()
	defer ss.Unlock()

	// Add new state to the ongoing set
	ss.ongoingWrites[&newState] = struct{}{}

	// Add new state to any ongoing read operations
	for rh := range ss.ongoingReads {
		rh.possibleStates = append(rh.possibleStates, &newState)
	}

	return &writeHandle{
		stateSet: ss,
		newState: &newState,
		// Save the list of committed states before the write was initiated
		committedBefore: ss.committed.Back(),
	}
}

// writeHandle is a handle to a write operation that is in progress.
type writeHandle struct {
	stateSet        *policyStateSet
	newState        *policyState
	committedBefore *list.Element
}

// commit should be called after the write operation has completed.
func (h *writeHandle) commit() {
	h.stateSet.Lock()
	defer h.stateSet.Unlock()

	// Remove all states that were committed before the write was initiated
	for cur := h.committedBefore; cur != nil; {
		prev := cur.Prev()
		h.stateSet.committed.Remove(cur)
		cur = prev
	}

	// Delete the new state from the ongoing set
	delete(h.stateSet.ongoingWrites, h.newState)

	// Add the new state to the committed list
	h.stateSet.committed.PushBack(h.newState)
}

func (p *policyStateEstimator) beginRead(key string) *readHandle {
	ss := p.statesByKey[key]
	ss.Lock()
	defer ss.Unlock()

	// Collect the current possible states of the policy into a read handle
	rh := &readHandle{
		key:            key,
		ss:             ss,
		possibleStates: make([]*policyState, 0, ss.committed.Len()+len(ss.ongoingWrites)),
	}
	for e := ss.committed.Front(); e != nil; e = e.Next() {
		rh.possibleStates = append(rh.possibleStates, e.Value.(*policyState))
	}
	rh.possibleStates = slices.AppendSeq(rh.possibleStates, maps.Keys(ss.ongoingWrites))

	// Add the read operation to the ongoing set to subscribe to any new writes
	ss.ongoingReads[rh] = struct{}{}

	return rh
}

type readHandle struct {
	key            string
	ss             *policyStateSet
	possibleStates []*policyState
}

func (rh *readHandle) finishAndAssert(tb testing.TB, state policyState) {
	rh.ss.Lock()
	// Remove the read operation from the ongoing set
	delete(rh.ss.ongoingReads, rh)
	rh.ss.Unlock()

	// Assert that the state is possible
	for _, s := range rh.possibleStates {
		if *s == state {
			return
		}
	}
	tb.Errorf("state %v is not possible for key %s, possible states are: %v", state, rh.key, rh.possibleStates)
}

// policyStoreConcurrencyHarness runs random concurrent set, get and delete
// operations against a concurrency checked policy store, and checks that
// every read returns a state consistent with the writes that overlapped it.
type policyStoreConcurrencyHarness struct {
	ps             *PolicyStore
	ctx            context.Context
	keys           []string
	samplePolicy   *Policy
	stateEstimator *policyStateEstimator
}

func newPolicyStoreConcurrencyHarness(tb testing.TB) *policyStoreConcurrencyHarness {
	const keyCount = 100

	// Include names that only differ by case, which share a policy
	keys := []string{"keya", "KeyA", "keyb", "keyc"}
	r := newTestRand()
	for len(keys) < keyCount {
		k, err := base62.RandomWithReader(5, r)
		require.NoError(tb, err)
		keys = append(keys, k)
	}

	samplePolicy, err := ParseACLPolicy(namespace.RootNamespace, `
path "secret/sample" {
	capabilities = ["read", "delete"]
}
`)
	require.NoError(tb, err)

	ps := mockConcurrencyCheckedCore(tb)
	stateEstimator := &policyStateEstimator{statesByKey: make(map[string]*policyStateSet, len(keys))}
	for _, key := range keys {
		stateEstimator.statesByKey[ps.sanitizeName(key)] = newPolicyStateSet()
	}

	return &policyStoreConcurrencyHarness{
		ps:             ps,
		ctx:            namespace.RootContext(context.Background()),
		keys:           keys,
		samplePolicy:   samplePolicy,
		stateEstimator: stateEstimator,
	}
}

func (h *policyStoreConcurrencyHarness) doRandomOp(tb testing.TB, r *rand.Rand) {
	key := h.keys[r.Intn(len(h.keys))]
	sanitizedKey := h.ps.sanitizeName(key)

	switch r.Intn(3) {
	case 0:
		p := h.samplePolicy.ShallowClone()
		p.Name = key
		wh := h.stateEstimator.beginWrite(sanitizedKey, policyState{exists: true, rawContents: p.Raw})
		err := h.ps.SetPolicy(h.ctx, p)
		wh.commit()
		assert.NoError(tb, err)
	case 1:
		rh := h.stateEstimator.beginRead(sanitizedKey)
		p, err := h.ps.GetPolicy(h.ctx, key, PolicyTypeACL)
		assert.NoError(tb, err)
		var state policyState
		if p != nil {
			state = policyState{exists: true, rawContents: p.Raw}
		}
		rh.finishAndAssert(tb, state)
	case 2:
		wh := h.stateEstimator.beginWrite(sanitizedKey, policyState{exists: false})
		err := h.ps.DeletePolicy(h.ctx, key, PolicyTypeACL)
		wh.commit()
		assert.NoError(tb, err)
	}
}

func (h *policyStoreConcurrencyHarness) assertNoOngoingOps(tb testing.TB) {
	for _, ss := range h.stateEstimator.statesByKey {
		assert.Empty(tb, ss.ongoingReads)
		assert.Empty(tb, ss.ongoingWrites)
	}
}

func newTestRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// TestPolicyStore_Concurrency runs random concurrent policy operations and
// verifies that storage is only touched while holding the policy's lock and
// that reads never observe a state that no overlapping write could produce.
func TestPolicyStore_Concurrency(t *testing.T) {
	t.Parallel()

	const (
		workers      = 16
		opsPerWorker = 500
	)

	h := newPolicyStoreConcurrencyHarness(t)
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := newTestRand()
			for range opsPerWorker {
				h.doRandomOp(t, r)
			}
		}()
	}
	wg.Wait()
	h.assertNoOngoingOps(t)
}

// BenchmarkPolicyStore_Concurrency runs the same checks as
// TestPolicyStore_Concurrency for as long as the benchmark runs, which
// surfaces rarer interleavings.
func BenchmarkPolicyStore_Concurrency(b *testing.B) {
	h := newPolicyStoreConcurrencyHarness(b)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		r := newTestRand()
		for pb.Next() {
			h.doRandomOp(b, r)
		}
	})
	b.StopTimer()

	h.assertNoOngoingOps(b)
}
