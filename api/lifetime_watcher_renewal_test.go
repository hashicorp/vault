// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// renewCallRecorder is a renewFunc that records when it was called and returns
// the configured result.
type renewCallRecorder struct {
	mu     sync.Mutex
	calls  []time.Time
	result func() (*Secret, error)
}

func (r *renewCallRecorder) renew(_ string, _ int) (*Secret, error) {
	r.mu.Lock()
	r.calls = append(r.calls, time.Now())
	r.mu.Unlock()
	return r.result()
}

func (r *renewCallRecorder) snapshot() []time.Time {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]time.Time(nil), r.calls...)
}

// TestLifetimeWatcher_ErrorBackoffNeverZero verifies that when renewals keep
// failing (for example, because they are rate limited), the lifetime watcher
// always waits between attempts instead of retrying in a tight loop.
//
// The back-off intervals and grace period are randomized, so a single watcher
// only sometimes reaches the problematic state. Several watchers are run in
// parallel to make the test reliable.
func TestLifetimeWatcher_ErrorBackoffNeverZero(t *testing.T) {
	t.Parallel()

	const (
		watchers             = 20
		leaseDurationSeconds = 6
		initialRetryInterval = time.Second
		// The exponential back-off uses a randomization factor of 0.5, so
		// no retry should ever happen sooner than half the initial interval.
		minExpectedGap = initialRetryInterval / 2
	)

	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)

	recorders := make([]*renewCallRecorder, watchers)
	errs := make([]error, watchers)
	var wg sync.WaitGroup
	for i := range recorders {
		rec := &renewCallRecorder{result: func() (*Secret, error) {
			return nil, errors.New("rate limit quota exceeded")
		}}
		recorders[i] = rec

		v, err := client.NewLifetimeWatcher(&LifetimeWatcherInput{
			Secret:    &Secret{LeaseDuration: leaseDurationSeconds},
			Increment: leaseDurationSeconds,
		})
		require.NoError(t, err)
		t.Cleanup(v.Stop)

		wg.Add(1)
		go func() {
			defer wg.Done()
			errs[i] = v.doRenewWithOptions(false, false, leaseDurationSeconds, "myleaseID", rec.renew, initialRetryInterval)
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(3 * leaseDurationSeconds * time.Second):
		t.Fatal("lifetime watchers did not finish")
	}

	for i, rec := range recorders {
		// Running out of time to renew is not fatal; callers are expected to
		// re-authenticate, just as when the lease enters its grace period.
		require.NoErrorf(t, errs[i], "watcher %d", i)

		calls := rec.snapshot()
		for j := 1; j < len(calls); j++ {
			if gap := calls[j].Sub(calls[j-1]); gap < minExpectedGap {
				t.Errorf("watcher %d: renewal attempt %d came %s after the previous one (total attempts: %d); expected at least %s",
					i, j+1, gap, len(calls), minExpectedGap)
				break
			}
		}
	}
}

// TestLifetimeWatcher_RetriesBeforeGracePeriod verifies that when a renewal
// fails and the next back-off would land in the grace period, the watcher makes
// one final retry just before the grace period instead of giving up straight
// away, and gives up if that retry fails too.
func TestLifetimeWatcher_RetriesBeforeGracePeriod(t *testing.T) {
	t.Parallel()

	const (
		leaseDurationSeconds = 3
		lease                = leaseDurationSeconds * time.Second
		// Longer than the lease, so the normal back-off always lands in the
		// grace period (10-20% of the lease) after the first failure.
		initialRetryInterval = 10 * time.Second
	)

	cases := []struct {
		name string
		// Renewal attempts, starting at 1, that succeed.
		succeedFrom   int64
		expectRenewal bool
	}{
		{
			name:          "final_retry_succeeds",
			succeedFrom:   2,
			expectRenewal: true,
		},
		{
			name:          "final_retry_fails",
			succeedFrom:   0,
			expectRenewal: false,
		},
	}

	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var attempt atomic.Int64
			rec := &renewCallRecorder{result: func() (*Secret, error) {
				if n := attempt.Add(1); tc.succeedFrom != 0 && n >= tc.succeedFrom {
					return &Secret{LeaseDuration: leaseDurationSeconds, Renewable: true}, nil
				}
				return nil, errors.New("rate limit quota exceeded")
			}}

			v, err := client.NewLifetimeWatcher(&LifetimeWatcherInput{
				Secret: &Secret{LeaseDuration: leaseDurationSeconds, Renewable: true},
			})
			require.NoError(t, err)
			defer v.Stop()

			doneCh := make(chan error, 1)
			go func() {
				doneCh <- v.doRenewWithOptions(false, false, leaseDurationSeconds, "myleaseID", rec.renew, initialRetryInterval)
			}()

			if tc.expectRenewal {
				select {
				case <-v.RenewCh():
				case err := <-doneCh:
					t.Fatalf("watcher gave up without retrying (err: %v), attempts: %d", err, len(rec.snapshot()))
				case <-time.After(lease):
					t.Fatal("lease not renewed before it expired")
				}
			} else {
				select {
				case err := <-doneCh:
					require.NoError(t, err)
				case <-time.After(lease + time.Second):
					t.Fatal("watcher did not give up")
				}
			}

			// The first attempt fails and the final retry is the second; if it
			// fails too the watcher gives up, and if it succeeds the next
			// renewal isn't due for roughly two thirds of the lease.
			calls := rec.snapshot()
			require.Len(t, calls, 2, "expected the first attempt plus exactly one final retry")
			if gap := calls[1].Sub(calls[0]); gap < minFinalRenewalRetryWait || gap >= lease {
				t.Fatalf("final retry came %s after the first attempt, expected between %s and %s",
					gap, minFinalRenewalRetryWait, lease)
			}
		})
	}
}
