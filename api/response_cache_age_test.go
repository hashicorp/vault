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

// TestSysRenew_PopulatesAgeFromCachedResponse asserts that a lease renewed
// through a caching proxy reports how long its response sat in that cache.
func TestSysRenew_PopulatesAgeFromCachedResponse(t *testing.T) {
	t.Parallel()

	const (
		leaseDurationSeconds = 30
		ageSeconds           = 25
	)

	server, _ := fakeCachedSecretServer(t, leaseDurationSeconds, ageSeconds)
	defer server.Close()

	client, err := NewClient(&Config{Address: server.URL})
	require.NoError(t, err)

	secret, err := client.Sys().Renew(testCachedLeaseID, leaseDurationSeconds)
	require.NoError(t, err)

	require.Equal(t, leaseDurationSeconds, secret.LeaseDuration,
		"lease_duration should be reported as the server sent it")
	require.Equal(t, ageSeconds*time.Second, secret.Age,
		"Age header should be surfaced to the caller")
}

// TestSysRenew_UncachedResponseHasNoAge asserts that a response served directly
// by Vault, with no Age header, reports no age.
func TestSysRenew_UncachedResponseHasNoAge(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"lease_id":%q,"lease_duration":30,"renewable":true}`, testCachedLeaseID)
	}))
	defer server.Close()

	client, err := NewClient(&Config{Address: server.URL})
	require.NoError(t, err)

	secret, err := client.Sys().Renew(testCachedLeaseID, 30)
	require.NoError(t, err)

	require.Zero(t, secret.Age, "a response served directly by Vault has no age")
}
