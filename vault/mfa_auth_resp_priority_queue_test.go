// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testCases returns cached responses with TimeOfStorage offsets from now. Some
// tests rely on the ordering of the items it returns.
func testCases() (tc []*MFACachedAuthResponse) {
	// create a slice of items with times offest by these seconds
	for _, m := range []time.Duration{
		5,
		183600,  // 51 hours
		15,      // 15 seconds
		45,      // 45 seconds
		900,     // 15 minutes
		360,     // 6 minutes
		7200,    // 2 hours
		183600,  // 51 hours
		7201,    // 2 hours, 1 second
		115200,  // 32 hours
		1209600, // 2 weeks
	} {
		n := time.Now()
		ft := n.Add(time.Second * m)
		uid, err := uuid.GenerateUUID()
		if err != nil {
			continue
		}
		tc = append(tc, &MFACachedAuthResponse{
			TimeOfStorage: ft,
			RequestID:     uid,
		})
	}
	return
}

// TestLoginMFAPriorityQueue_PushPopByKey verifies that entries can be pushed
// and popped by request ID, that duplicate and empty keys are handled, and that
// popping an unknown key returns nothing.
func TestLoginMFAPriorityQueue_PushPopByKey(t *testing.T) {
	pq := NewLoginMFAPriorityQueue()

	if pq.Len() != 0 {
		t.Fatalf("expected new queue to have zero size, got (%d)", pq.Len())
	}

	tc := testCases()
	tcl := len(tc)
	for _, i := range tc {
		if err := pq.Push(i); err != nil {
			t.Fatal(err)
		}
	}

	if pq.Len() != tcl {
		t.Fatalf("error adding items, expected (%d) items, got (%d)", tcl, pq.Len())
	}

	item, err := pq.PopByKey(tc[0].RequestID)
	if err != nil {
		t.Fatalf("error popping item: %s", err)
	}
	if tc[0].TimeOfStorage != item.TimeOfStorage {
		t.Fatalf("expected tc[0] and popped item to match, got (%v) and (%v)", tc[0].TimeOfStorage, item.TimeOfStorage)
	}

	// push item with duplicate key
	dErr := pq.Push(tc[1])
	if dErr != queue.ErrDuplicateItem {
		t.Fatal(err)
	}
	// push item with no key
	tc[2].RequestID = ""
	kErr := pq.Push(tc[2])
	if kErr != nil && kErr.Error() != "error adding item: Item Key is required" {
		t.Fatal(kErr)
	}

	// check nil,nil error for not found
	i, err := pq.PopByKey("empty")
	if err != nil && i != nil {
		t.Fatalf("expected nil error for PopByKey of non-existing key, got: %s", err)
	}
}

// TestLoginMFARemoveStaleEntries verifies that RemoveExpiredMfaAuthResponse
// removes every entry whose age exceeds the TTL at the given cutoff time, and
// leaves the entries that have not yet expired in the queue. The test calls it
// repeatedly, because a single call removes at most
// maxExpiredMFAAuthResponsesPerCleanup entries.
func TestLoginMFARemoveStaleEntries(t *testing.T) {
	pq := NewLoginMFAPriorityQueue()

	tc := testCases()
	for _, i := range tc {
		if err := pq.Push(i); err != nil {
			t.Fatal(err)
		}
	}

	cutoffTime := time.Now().Add(371 * time.Second)
	timeout := time.Now().Add(5 * time.Second)
	for {
		if time.Now().After(timeout) {
			break
		}
		pq.RemoveExpiredMfaAuthResponse(defaultMFAAuthResponseTTL, cutoffTime)
	}

	if pq.Len() != 8 {
		t.Fatalf("failed to remove %d stale entries", pq.Len())
	}
}

// TestMFACachedAuthResponse_Lifetime verifies the lifetime boundary. An entry
// is valid up to, but not including, TimeOfStorage+ttl, and an unset
// TimeOfStorage is treated as past its lifetime so a missed assignment fails
// closed.
func TestMFACachedAuthResponse_Lifetime(t *testing.T) {
	t.Parallel()

	stored := time.Now()
	resp := &MFACachedAuthResponse{TimeOfStorage: stored}
	ttl := defaultMFAAuthResponseTTL

	require.False(t, resp.isExpired(ttl, stored.Add(ttl-time.Nanosecond)))
	require.True(t, resp.isExpired(ttl, stored.Add(ttl)))
	require.True(t, (&MFACachedAuthResponse{}).isExpired(ttl, time.Now()))
}

// TestPopMFAResponseAuthByID_OlderEntriesQueued verifies that redemption does
// not return an expired entry while older expired entries are still queued
// ahead of it. The expired entry is removed on the first attempt, and a live
// entry in the same queue is still returned.
func TestPopMFAResponseAuthByID_OlderEntriesQueued(t *testing.T) {
	t.Parallel()

	c := &Core{mfaResponseAuthQueue: NewLoginMFAPriorityQueue()}
	now := time.Now()

	for i := 0; i < 32; i++ {
		require.NoError(t, c.SaveMFAResponseAuth(&MFACachedAuthResponse{
			RequestID:     fmt.Sprintf("queued-%d", i),
			TimeOfStorage: now.Add(-time.Hour - time.Duration(i)*time.Second),
		}))
	}
	require.NoError(t, c.SaveMFAResponseAuth(&MFACachedAuthResponse{
		RequestID:     "target",
		TimeOfStorage: now.Add(-6 * time.Minute),
	}))
	require.NoError(t, c.SaveMFAResponseAuth(&MFACachedAuthResponse{
		RequestID:     "live",
		TimeOfStorage: now,
	}))

	got, err := c.PopMFAResponseAuthByID("target")
	require.NoError(t, err)
	require.Nil(t, got)

	got, err = c.PopMFAResponseAuthByID("target")
	require.NoError(t, err)
	require.Nil(t, got, "entry must be gone after the first attempt")

	got, err = c.PopMFAResponseAuthByID("live")
	require.NoError(t, err)
	require.NotNil(t, got)
}

// TestPopMFAResponseAuthByID_LifetimeBoundary verifies redemption honors the
// lifetime on either side of the limit.
func TestPopMFAResponseAuthByID_LifetimeBoundary(t *testing.T) {
	t.Parallel()

	c := &Core{mfaResponseAuthQueue: NewLoginMFAPriorityQueue()}
	now := time.Now()
	require.NoError(t, c.SaveMFAResponseAuth(&MFACachedAuthResponse{
		RequestID:     "just-inside",
		TimeOfStorage: now.Add(-(defaultMFAAuthResponseTTL - 30*time.Second)),
	}))
	require.NoError(t, c.SaveMFAResponseAuth(&MFACachedAuthResponse{
		RequestID:     "just-outside",
		TimeOfStorage: now.Add(-(defaultMFAAuthResponseTTL + time.Second)),
	}))

	got, err := c.PopMFAResponseAuthByID("just-inside")
	require.NoError(t, err)
	require.NotNil(t, got)

	got, err = c.PopMFAResponseAuthByID("just-outside")
	require.NoError(t, err)
	require.Nil(t, got)
}

// TestPopMFAResponseAuthByID_NilQueue verifies that redemption and save do not
// panic after teardownLoginMFA has cleared the queue.
func TestPopMFAResponseAuthByID_NilQueue(t *testing.T) {
	t.Parallel()

	c := &Core{}
	got, err := c.PopMFAResponseAuthByID("anything")
	require.NoError(t, err)
	require.Nil(t, got)
	require.Error(t, c.SaveMFAResponseAuth(&MFACachedAuthResponse{RequestID: "x"}))
}

// TestLoginMFAPriorityQueue_CleanupRemovesAllStale verifies that when there are
// fewer than maxExpiredMFAAuthResponsesPerCleanup expired entries at the head
// of the queue, a single cleanup call removes all of them and leaves the live
// entries in place.
func TestLoginMFAPriorityQueue_CleanupRemovesAllStale(t *testing.T) {
	t.Parallel()

	pq := NewLoginMFAPriorityQueue()
	now := time.Now()
	for i := 0; i < 50; i++ {
		require.NoError(t, pq.Push(&MFACachedAuthResponse{
			RequestID:     fmt.Sprintf("stale-%d", i),
			TimeOfStorage: now.Add(-time.Hour - time.Duration(i)*time.Second),
		}))
	}
	for i := 0; i < 3; i++ {
		require.NoError(t, pq.Push(&MFACachedAuthResponse{
			RequestID:     fmt.Sprintf("live-%d", i),
			TimeOfStorage: now.Add(time.Duration(i) * time.Second),
		}))
	}

	require.NoError(t, pq.RemoveExpiredMfaAuthResponse(defaultMFAAuthResponseTTL, now))
	require.Equal(t, 3, pq.Len())

	got, err := pq.PopByKey("live-0")
	require.NoError(t, err)
	require.NotNil(t, got)
}

// TestLoginMFAPriorityQueue_CleanupKeepsLiveEntryAvailable verifies that
// cleanup never makes a live request ID briefly unavailable.
func TestLoginMFAPriorityQueue_CleanupKeepsLiveEntryAvailable(t *testing.T) {
	t.Parallel()

	c := &Core{mfaResponseAuthQueue: NewLoginMFAPriorityQueue()}
	require.NoError(t, c.SaveMFAResponseAuth(&MFACachedAuthResponse{
		RequestID:     "live",
		TimeOfStorage: time.Now(),
	}))

	stop := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
				assert.NoError(t, c.mfaResponseAuthQueue.RemoveExpiredMfaAuthResponse(defaultMFAAuthResponseTTL, time.Now()))
			}
		}
	}()

	// The entry is the queue head, so every cleanup call pops and re-pushes it.
	for i := 0; i < 5000; i++ {
		got, err := c.PopMFAResponseAuthByID("live")
		if !assert.NoError(t, err) || !assert.NotNil(t, got, "live entry went missing on iteration %d", i) {
			break
		}
		require.NoError(t, c.SaveMFAResponseAuth(got))
	}

	close(stop)
	wg.Wait()
}

// TestPopMFAResponseAuthByID_ConcurrentWithCleanup races redemption against
// cleanup. Each live ID must be returned exactly once across several
// concurrent redeemers, and stale IDs must never be returned.
func TestPopMFAResponseAuthByID_ConcurrentWithCleanup(t *testing.T) {
	t.Parallel()

	c := &Core{mfaResponseAuthQueue: NewLoginMFAPriorityQueue()}
	now := time.Now()
	const n = 200
	for i := 0; i < n; i++ {
		require.NoError(t, c.SaveMFAResponseAuth(&MFACachedAuthResponse{
			RequestID:     fmt.Sprintf("stale-%d", i),
			TimeOfStorage: now.Add(-time.Hour),
		}))
		require.NoError(t, c.SaveMFAResponseAuth(&MFACachedAuthResponse{
			RequestID:     fmt.Sprintf("live-%d", i),
			TimeOfStorage: now,
		}))
	}

	var wg sync.WaitGroup
	var liveHits, staleHits atomic.Int64
	for i := 0; i < n; i++ {
		for _, id := range []string{fmt.Sprintf("live-%d", i), fmt.Sprintf("stale-%d", i)} {
			for r := 0; r < 3; r++ { // three redeemers per ID
				wg.Add(1)
				go func(id string) {
					defer wg.Done()
					got, err := c.PopMFAResponseAuthByID(id)
					assert.NoError(t, err)
					if got == nil {
						return
					}
					if strings.HasPrefix(id, "live-") {
						liveHits.Add(1)
					} else {
						staleHits.Add(1)
					}
				}(id)
			}
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.NoError(t, c.mfaResponseAuthQueue.RemoveExpiredMfaAuthResponse(defaultMFAAuthResponseTTL, time.Now()))
		}()
	}
	wg.Wait()

	require.EqualValues(t, n, liveHits.Load())
	require.Zero(t, staleHits.Load())
}

// TestPopMFAResponseAuthByID_ConcurrentWithTeardown races redemption and save
// against the step teardownLoginMFA performs on the queue (clearing it under
// the queue lock). Neither call may panic, and each must either succeed or
// fail cleanly.
func TestPopMFAResponseAuthByID_ConcurrentWithTeardown(t *testing.T) {
	t.Parallel()

	c := &Core{mfaResponseAuthQueue: NewLoginMFAPriorityQueue()}
	for i := 0; i < 100; i++ {
		require.NoError(t, c.SaveMFAResponseAuth(&MFACachedAuthResponse{
			RequestID:     fmt.Sprintf("id-%d", i),
			TimeOfStorage: time.Now(),
		}))
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			_, err := c.PopMFAResponseAuthByID(fmt.Sprintf("id-%d", i))
			assert.NoError(t, err)
		}(i)
		go func(i int) {
			defer wg.Done()
			// An error is acceptable once the queue is gone; a panic is not.
			_ = c.SaveMFAResponseAuth(&MFACachedAuthResponse{
				RequestID:     fmt.Sprintf("new-%d", i),
				TimeOfStorage: time.Now(),
			})
		}(i)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.mfaResponseAuthQueueLock.Lock()
		c.mfaResponseAuthQueue = nil
		c.mfaResponseAuthQueueLock.Unlock()
	}()
	wg.Wait()
}

// pendingTestResp builds an MFACachedAuthResponse for the pending-limit tests.
// The entry counts against the given entity ID and the host part of addr, and
// is considered expired once stored is older than defaultMFAAuthResponseTTL.
// An empty entityID leaves CachedAuth nil, so no per-entity limit applies.
func pendingTestResp(id, entityID, addr string, stored time.Time) *MFACachedAuthResponse {
	resp := &MFACachedAuthResponse{RequestID: id, RequestConnRemoteAddr: addr, TimeOfStorage: stored}
	if entityID != "" {
		resp.CachedAuth = &logical.Auth{EntityID: entityID}
	}
	return resp
}

// TestLoginMFAPriorityQueue_PendingLimits verifies that each pending MFA cap
// (per entity, per source IP and global) accepts exactly its limit of live
// entries, rejects the next one with an *mfaPendingLimitError naming the
// scope, and still accepts an entry that falls outside the full scope. The
// per-source-IP case also checks that the IP is compared without its port, so
// changing ports does not get around the cap.
//
// The per-entity case uses the real limit of 100, because 100 pushes are
// cheap. The per-source-IP and global cases lower the limit on their own
// queue (to 10 and 20). Their real values of 1,000 and 100,000 would need that
// many pushes, and each push deep copies the entry, which made the global case
// take about 85 seconds under -race. The limit check is the same code for any
// value, so a small limit tests the same logic. The production values
// themselves are checked by the defaults case, which asserts that a new queue
// starts with 100, 1,000 and 100,000.
func TestLoginMFAPriorityQueue_PendingLimits(t *testing.T) {
	t.Parallel()

	t.Run("defaults", func(t *testing.T) {
		t.Parallel()
		pq := NewLoginMFAPriorityQueue()
		require.Equal(t, mfaPendingLimits{perEntity: 100, perSourceIP: 1000, global: 100000}, pq.limits)
	})

	tests := []struct {
		name string
		// limits replaces the default limits on the test's queue, if set.
		limits *mfaPendingLimits
		// fillN is how many entries fill the scope to its limit.
		fillN int
		// fill returns the entry for index i while filling the scope.
		fill func(i int) *MFACachedAuthResponse
		// over is the entry that must be rejected once the scope is full.
		over *MFACachedAuthResponse
		// outside is an entry in a different scope that must still be accepted.
		// It is nil for the global case, where no entry is outside the scope.
		outside   *MFACachedAuthResponse
		wantScope string
	}{
		{
			name:  "per entity uses the real limit",
			fillN: maxPendingMFAPerEntity,
			fill: func(i int) *MFACachedAuthResponse {
				return pendingTestResp(fmt.Sprintf("e-%d", i), "ent", fmt.Sprintf("10.0.%d.%d:1", i/200, i%200), time.Now())
			},
			over:      pendingTestResp("e-over", "ent", "9.9.9.9:1", time.Now()),
			outside:   pendingTestResp("e-other", "ent2", "9.9.9.9:1", time.Now()),
			wantScope: "entity",
		},
		{
			name:   "per source IP ignores the port",
			limits: &mfaPendingLimits{perEntity: maxPendingMFAPerEntity, perSourceIP: 10, global: maxPendingMFAGlobal},
			fillN:  10,
			fill: func(i int) *MFACachedAuthResponse {
				return pendingTestResp(fmt.Sprintf("s-%d", i), fmt.Sprintf("ent-%d", i), fmt.Sprintf("1.2.3.4:%d", i+1), time.Now())
			},
			over:      pendingTestResp("s-over", "x", "1.2.3.4:60000", time.Now()),
			outside:   pendingTestResp("s-other", "x", "1.2.3.5:1", time.Now()),
			wantScope: "source_ip",
		},
		{
			name:   "global",
			limits: &mfaPendingLimits{perEntity: maxPendingMFAPerEntity, perSourceIP: maxPendingMFAPerSourceIP, global: 20},
			fillN:  20,
			fill: func(i int) *MFACachedAuthResponse {
				return pendingTestResp(fmt.Sprintf("g-%d", i), fmt.Sprintf("ent-%d", i), fmt.Sprintf("10.0.0.%d:1", i), time.Now())
			},
			over:      pendingTestResp("g-over", "new", "8.8.8.8:1", time.Now()),
			wantScope: "global",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			pq := NewLoginMFAPriorityQueue()
			if tc.limits != nil {
				pq.limits = *tc.limits
			}

			for i := 0; i < tc.fillN; i++ {
				require.NoError(t, pq.Push(tc.fill(i)), "entry %d should be accepted", i+1)
			}

			var limitErr *mfaPendingLimitError
			require.ErrorAs(t, pq.Push(tc.over), &limitErr)
			require.Equal(t, tc.wantScope, limitErr.scope)

			if tc.outside != nil {
				require.NoError(t, pq.Push(tc.outside))
			}
		})
	}
}

// TestLoginMFAPriorityQueue_PendingLimitsCounters verifies that the counters
// shrink back to empty after entries are popped or expire, that expired entries
// do not count toward a cap even before background cleanup runs, and that a
// re-push of a popped entry at the cap is rejected.
func TestLoginMFAPriorityQueue_PendingLimitsCounters(t *testing.T) {
	t.Parallel()
	now := time.Now()

	pq := NewLoginMFAPriorityQueue()
	expired := now.Add(-2 * defaultMFAAuthResponseTTL)
	for i := 0; i < maxPendingMFAPerEntity; i++ {
		require.NoError(t, pq.Push(pendingTestResp(fmt.Sprintf("old-%d", i), "ent", "1.1.1.1:1", expired)))
	}
	// Expired entries must not block a new one, even without cleanup.
	require.NoError(t, pq.Push(pendingTestResp("new", "ent", "1.1.1.1:2", now)))
	require.Equal(t, 1, pq.byEntity["ent"])
	require.Equal(t, 1, pq.bySource["1.1.1.1"])

	got, err := pq.PopByKey("new")
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Empty(t, pq.byEntity)
	require.Empty(t, pq.bySource)

	// Cleanup leaves live entries' counts unchanged and removes expired counts.
	require.NoError(t, pq.Push(pendingTestResp("live", "ent", "2.2.2.2:1", now)))
	require.NoError(t, pq.Push(pendingTestResp("dead", "ent", "2.2.2.2:1", expired)))
	require.NoError(t, pq.RemoveExpiredMfaAuthResponse(defaultMFAAuthResponseTTL, time.Now()))
	require.Equal(t, 1, pq.byEntity["ent"])
	require.Equal(t, 1, pq.bySource["2.2.2.2"])

	// Re-push at the cap follows the cap.
	pq = NewLoginMFAPriorityQueue()
	for i := 0; i < maxPendingMFAPerEntity; i++ {
		require.NoError(t, pq.Push(pendingTestResp(fmt.Sprintf("r-%d", i), "ent", "3.3.3.3:1", now)))
	}
	popped, err := pq.PopByKey("r-0")
	require.NoError(t, err)
	require.NoError(t, pq.Push(pendingTestResp("filler", "ent", "3.3.3.3:1", now)))
	require.Error(t, pq.Push(popped))
	require.Equal(t, maxPendingMFAPerEntity, pq.byEntity["ent"])
}
