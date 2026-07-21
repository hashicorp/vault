// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	testCachedLeaseID              = "database/creds/readonly/abcd1234"
	testCachedLeaseDurationSeconds = 6
	testCachedLeaseAgeSeconds      = 5
)

// fakeCachedSecretServer stands in for a caching proxy serving a hit: it replays
// a lease that was issued ageSeconds ago, with the body still reporting the
// duration the lease had when it was issued.
//
// Renewals are timed from server start, so a test can tell whether one was
// attempted after the lease it renews had already expired.
func fakeCachedSecretServer(t *testing.T, leaseDurationSeconds int, ageSeconds int) (*httptest.Server, func() []time.Duration) {
	t.Helper()

	var (
		mu       sync.Mutex
		renewals []time.Duration
	)
	start := time.Now()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		mu.Lock()
		renewals = append(renewals, time.Since(start))
		mu.Unlock()

		w.Header().Set("Age", fmt.Sprintf("%d", ageSeconds))
		w.Header().Set("X-Cache", "HIT")
		w.Header().Set("Date", time.Now().Format(http.TimeFormat))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"lease_id":%q,"lease_duration":%d,"renewable":true}`,
			testCachedLeaseID, leaseDurationSeconds)
	}))

	return server, func() []time.Duration {
		mu.Lock()
		defer mu.Unlock()

		return slices.Clone(renewals)
	}
}

// testWatchCachedLease starts a watcher over a lease whose renewals are served by
// the given caching proxy.
func testWatchCachedLease(t *testing.T, server *httptest.Server) *LifetimeWatcher {
	t.Helper()

	client, err := NewClient(&Config{Address: server.URL})
	require.NoError(t, err)

	watcher, err := client.NewLifetimeWatcher(&LifetimeWatcherInput{
		Secret: &Secret{
			LeaseID:       testCachedLeaseID,
			LeaseDuration: testCachedLeaseDurationSeconds,
			Renewable:     true,
		},
	})
	require.NoError(t, err)

	go watcher.Start()

	return watcher
}

// TestLifetimeWatcher_SignalsReReadOfAgedLease asserts that the watcher gives up
// and reports that the secret must be re-read, rather than renewing forever.
//
// A caching proxy reports how stale its response is in the HTTP Age header, but
// the lease_duration in the body still reads from issuance. The watcher ignores
// Age, so the remaining lifetime it computes never falls and the threshold at
// which it gives up is never reached.
//
// See https://github.com/hashicorp/vault/issues/19227.
func TestLifetimeWatcher_SignalsReReadOfAgedLease(t *testing.T) {
	t.Parallel()

	server, _ := fakeCachedSecretServer(t, testCachedLeaseDurationSeconds, testCachedLeaseAgeSeconds)
	defer server.Close()

	watcher := testWatchCachedLease(t, server)
	defer watcher.Stop()

	select {
	case err := <-watcher.DoneCh():
		require.NoError(t, err)
	case <-time.After(3 * time.Second):
		t.Fatalf("watcher never signalled a re-read for a %ds lease already %ds old",
			testCachedLeaseDurationSeconds, testCachedLeaseAgeSeconds)
	}
}

// TestLifetimeWatcher_DoesNotRenewExpiredLease asserts that the watcher stops
// rather than going on renewing a lease that has already run out.
//
// The lease has only its duration less its age left when it arrives, but a
// watcher that reads lease_duration alone sleeps for two thirds of the full
// duration before renewing again, by which time the credential is long dead.
// Each such renewal is a request Vault and its secrets backend must serve for a
// lease that no longer exists.
func TestLifetimeWatcher_DoesNotRenewExpiredLease(t *testing.T) {
	t.Parallel()

	remaining := time.Duration(testCachedLeaseDurationSeconds-testCachedLeaseAgeSeconds) * time.Second

	server, renewals := fakeCachedSecretServer(t, testCachedLeaseDurationSeconds, testCachedLeaseAgeSeconds)
	defer server.Close()

	watcher := testWatchCachedLease(t, server)
	defer watcher.Stop()

	// Wait long enough to catch a second renewal, which a watcher that ignores
	// Age schedules for roughly two thirds of the reported lease duration.
	select {
	case <-watcher.DoneCh():
	case <-time.After(5 * time.Second):
	}

	for _, at := range renewals() {
		if at > remaining {
			t.Fatalf("renewed %s after the lease expired (renewal at %s, only %s of life left on arrival)",
				at-remaining, at, remaining)
		}
	}
}
