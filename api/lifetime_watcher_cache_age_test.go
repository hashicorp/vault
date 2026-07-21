// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
func fakeCachedSecretServer(t *testing.T, leaseDurationSeconds int, ageSeconds int) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Age", fmt.Sprintf("%d", ageSeconds))
		w.Header().Set("X-Cache", "HIT")
		w.Header().Set("Date", time.Now().Format(http.TimeFormat))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"lease_id":%q,"lease_duration":%d,"renewable":true}`,
			testCachedLeaseID, leaseDurationSeconds)
	}))
}

// testWatchCachedLease starts a watcher over a lease whose renewals are served
// by the given caching proxy.
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

	server := fakeCachedSecretServer(t, testCachedLeaseDurationSeconds, testCachedLeaseAgeSeconds)
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
