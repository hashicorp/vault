// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cache_test

import (
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/internalshared/agentproxyshared/cache"
	"github.com/hashicorp/vault/internalshared/agentproxyshared/cache/cachememdb"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/stretchr/testify/require"
)

// testNewStaticSecretCapabilityManager returns a new StaticSecretCapabilityManager
// for use in tests. behaviour is the token capability refresh behaviour
// ("optimistic", "pessimistic", or "" for the default).
func testNewStaticSecretCapabilityManager(t *testing.T, client *api.Client, behaviour string) *cache.StaticSecretCapabilityManager {
	t.Helper()

	lc := testNewLeaseCache(t, []*cache.SendResponse{})

	updater, err := cache.NewStaticSecretCapabilityManager(&cache.StaticSecretCapabilityManagerConfig{
		LeaseCache: lc,
		Logger:     logging.NewVaultLogger(hclog.Trace).Named("cache.capabilitiesmanager"),
		Client:     client,
		StaticSecretTokenCapabilityRefreshInterval:  250 * time.Millisecond,
		StaticSecretTokenCapabilityRefreshBehaviour: behaviour,
	})
	if err != nil {
		t.Fatal(err)
	}

	return updater
}

// TestGetCapabilitiesRootToken tests the getCapabilities method with the root
// token, expecting to get "root" capabilities on valid paths
func TestGetCapabilitiesRootToken(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	capabilitiesToCheck := []string{"auth/token/create", "sys/health"}
	capabilities, err := cache.GetCapabilitiesForTest(capabilitiesToCheck, client)
	require.NoError(t, err)

	expectedCapabilities := map[string][]string{
		"auth/token/create": {"root"},
		"sys/health":        {"root"},
	}
	require.Equal(t, expectedCapabilities, capabilities)
}

// TestGetCapabilitiesLowPrivilegeToken tests the getCapabilities method with
// a low privilege token, expecting to get deny or non-root capabilities
func TestGetCapabilitiesLowPrivilegeToken(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	renewable := true
	// Set the token's policies to 'default' and nothing else
	tokenCreateRequest := &api.TokenCreateRequest{
		Policies:  []string{"default"},
		TTL:       "30m",
		Renewable: &renewable,
	}

	secret, err := client.Auth().Token().CreateOrphan(tokenCreateRequest)
	require.NoError(t, err)
	token := secret.Auth.ClientToken

	client.SetToken(token)

	capabilitiesToCheck := []string{"auth/token/create", "sys/capabilities-self", "auth/token/lookup-self"}
	capabilities, err := cache.GetCapabilitiesForTest(capabilitiesToCheck, client)
	require.NoError(t, err)

	expectedCapabilities := map[string][]string{
		"auth/token/create":      {"deny"},
		"sys/capabilities-self":  {"update"},
		"auth/token/lookup-self": {"read"},
	}
	require.Equal(t, expectedCapabilities, capabilities)
}

// TestGetCapabilitiesBadClientToken tests that getCapabilities
// returns an empty set of capabilities if the token is bad (and it gets a 403)
func TestGetCapabilitiesBadClientToken(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client
	client.SetToken("")

	capabilitiesToCheck := []string{"auth/token/create", "sys/capabilities-self", "auth/token/lookup-self"}
	capabilities, err := cache.GetCapabilitiesForTest(capabilitiesToCheck, client)
	require.Nil(t, err)
	require.Equal(t, map[string][]string{}, capabilities)
}

// TestGetCapabilitiesEmptyPaths tests the getCapabilities will error on an empty
// set of paths to check
func TestGetCapabilitiesEmptyPaths(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	var capabilitiesToCheck []string
	_, err := cache.GetCapabilitiesForTest(capabilitiesToCheck, client)
	require.Error(t, err)
}

// TestSubmitWorkUpdatesIndex tests that an index will be correctly updated if the capabilities differ.
func TestSubmitWorkUpdatesIndex(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	// Create a low permission token
	renewable := true
	// Set the token's policies to 'default' and nothing else
	tokenCreateRequest := &api.TokenCreateRequest{
		Policies:  []string{"default"},
		TTL:       "30m",
		Renewable: &renewable,
	}

	secret, err := client.Auth().Token().CreateOrphan(tokenCreateRequest)
	require.NoError(t, err)
	token := secret.Auth.ClientToken
	indexId := cache.HashStaticSecretIndexForTest(token)

	sscm := testNewStaticSecretCapabilityManager(t, client, "")
	index := &cachememdb.CapabilitiesIndex{
		ID:    indexId,
		Token: token,
		// The token will (perhaps obviously) not have
		// read access to /foo/bar, but will to /auth/token/lookup-self
		ReadablePaths: map[string]struct{}{
			"foo/bar":                {},
			"auth/token/lookup-self": {},
		},
	}
	err = sscm.LeaseCacheForTest().DBForTest().SetCapabilitiesIndex(index)
	require.Nil(t, err)

	sscm.StartRenewingCapabilities(index)

	// Wait for the job to complete at least once...
	time.Sleep(3 * time.Second)

	newIndex, err := sscm.LeaseCacheForTest().DBForTest().GetCapabilitiesIndex(cachememdb.IndexNameID, indexId)
	require.Nil(t, err)
	newIndex.IndexLock.RLock()
	require.Equal(t, map[string]struct{}{
		"auth/token/lookup-self": {},
	}, newIndex.ReadablePaths)
	newIndex.IndexLock.RUnlock()

	// Forcefully stop any remaining workers
	sscm.Stop()
}

// TestSubmitWorkUpdatesIndexWithBadToken tests that an index will be correctly updated if the token
// has expired and we cannot access the sys capabilities endpoint.
func TestSubmitWorkUpdatesIndexWithBadToken(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	token := "not real token"
	indexId := cache.HashStaticSecretIndexForTest(token)

	sscm := testNewStaticSecretCapabilityManager(t, client, "")
	index := &cachememdb.CapabilitiesIndex{
		ID:    indexId,
		Token: token,
		ReadablePaths: map[string]struct{}{
			"foo/bar":                {},
			"auth/token/lookup-self": {},
		},
	}
	err := sscm.LeaseCacheForTest().DBForTest().SetCapabilitiesIndex(index)
	require.Nil(t, err)

	sscm.StartRenewingCapabilities(index)

	// Wait for the job to complete at least once...
	time.Sleep(3 * time.Second)

	// This entry should be evicted.
	newIndex, err := sscm.LeaseCacheForTest().DBForTest().GetCapabilitiesIndex(cachememdb.IndexNameID, indexId)
	require.Equal(t, err, cachememdb.ErrCacheItemNotFound)
	require.Nil(t, newIndex)

	// Forcefully stop any remaining workers
	sscm.Stop()
}

// TestSubmitWorkSealedVaultOptimistic tests that the capability manager
// behaves as expected when
// sscm.tokenCapabilityRefreshBehaviour == TokenCapabilityRefreshBehaviourOptimistic
func TestSubmitWorkSealedVaultOptimistic(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	token := "not real token"
	indexId := cache.HashStaticSecretIndexForTest(token)

	sscm := testNewStaticSecretCapabilityManager(t, client, "")
	index := &cachememdb.CapabilitiesIndex{
		ID:    indexId,
		Token: token,
		ReadablePaths: map[string]struct{}{
			"foo/bar":                {},
			"auth/token/lookup-self": {},
		},
	}
	err := sscm.LeaseCacheForTest().DBForTest().SetCapabilitiesIndex(index)
	require.Nil(t, err)

	// Seal the cluster
	cluster.EnsureCoresSealed(t)

	sscm.StartRenewingCapabilities(index)

	// Wait for the job to complete at least once...
	time.Sleep(3 * time.Second)

	// This entry should not be evicted.
	newIndex, err := sscm.LeaseCacheForTest().DBForTest().GetCapabilitiesIndex(cachememdb.IndexNameID, indexId)
	require.NoError(t, err)
	require.NotNil(t, newIndex)

	// Forcefully stop any remaining workers
	sscm.Stop()
}

// TestSubmitWorkSealedVaultPessimistic tests that the capability manager
// behaves as expected when
// sscm.tokenCapabilityRefreshBehaviour == TokenCapabilityRefreshBehaviourPessimistic
func TestSubmitWorkSealedVaultPessimistic(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	token := "not real token"
	indexId := cache.HashStaticSecretIndexForTest(token)

	sscm := testNewStaticSecretCapabilityManager(t, client, "pessimistic")

	index := &cachememdb.CapabilitiesIndex{
		ID:    indexId,
		Token: token,
		ReadablePaths: map[string]struct{}{
			"foo/bar":                {},
			"auth/token/lookup-self": {},
		},
	}
	err := sscm.LeaseCacheForTest().DBForTest().SetCapabilitiesIndex(index)
	require.Nil(t, err)

	// Seal the cluster
	cluster.EnsureCoresSealed(t)

	sscm.StartRenewingCapabilities(index)

	// Wait for the job to complete at least once...
	time.Sleep(3 * time.Second)

	// This entry should be evicted.
	newIndex, err := sscm.LeaseCacheForTest().DBForTest().GetCapabilitiesIndex(cachememdb.IndexNameID, indexId)
	require.Error(t, err)
	require.Nil(t, newIndex)

	// Forcefully stop any remaining workers
	sscm.Stop()
}

// TestSubmitWorkUpdatesAllIndexes tests that an index will be correctly updated if the capabilities differ, as
// well as the indexes related to the paths that are being checked for.
func TestSubmitWorkUpdatesAllIndexes(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	// Create a low permission token
	renewable := true
	// Set the token's policies to 'default' and nothing else
	tokenCreateRequest := &api.TokenCreateRequest{
		Policies:  []string{"default"},
		TTL:       "30m",
		Renewable: &renewable,
	}

	secret, err := client.Auth().Token().CreateOrphan(tokenCreateRequest)
	require.NoError(t, err)
	token := secret.Auth.ClientToken
	indexId := cache.HashStaticSecretIndexForTest(token)

	sscm := testNewStaticSecretCapabilityManager(t, client, "")
	index := &cachememdb.CapabilitiesIndex{
		ID:    indexId,
		Token: token,
		// The token will (perhaps obviously) not have
		// read access to /foo/bar, but will to /auth/token/lookup-self
		ReadablePaths: map[string]struct{}{
			"foo/bar":                {},
			"auth/token/lookup-self": {},
		},
	}
	err = sscm.LeaseCacheForTest().DBForTest().SetCapabilitiesIndex(index)
	require.Nil(t, err)

	pathIndexId1 := cache.HashStaticSecretIndexForTest("foo/bar")
	pathIndex1 := &cachememdb.Index{
		ID:        pathIndexId1,
		Namespace: "root/",
		Tokens: map[string]struct{}{
			token: {},
		},
		RequestPath: "foo/bar",
		Response:    []byte{},
	}

	pathIndexId2 := cache.HashStaticSecretIndexForTest("auth/token/lookup-self")
	pathIndex2 := &cachememdb.Index{
		ID:        pathIndexId2,
		Namespace: "root/",
		Tokens: map[string]struct{}{
			token: {},
		},
		RequestPath: "auth/token/lookup-self",
		Response:    []byte{},
	}

	err = sscm.LeaseCacheForTest().DBForTest().Set(pathIndex1)
	require.Nil(t, err)

	err = sscm.LeaseCacheForTest().DBForTest().Set(pathIndex2)
	require.Nil(t, err)

	sscm.StartRenewingCapabilities(index)

	// Wait for the job to complete at least once...
	time.Sleep(1 * time.Second)

	newIndex, err := sscm.LeaseCacheForTest().DBForTest().GetCapabilitiesIndex(cachememdb.IndexNameID, indexId)
	require.Nil(t, err)
	newIndex.IndexLock.RLock()
	require.Equal(t, map[string]struct{}{
		"auth/token/lookup-self": {},
	}, newIndex.ReadablePaths)
	newIndex.IndexLock.RUnlock()

	// For this, we expect the token to have been deleted
	newPathIndex1, err := sscm.LeaseCacheForTest().DBForTest().Get(cachememdb.IndexNameID, pathIndexId1)
	require.Nil(t, err)
	require.Equal(t, map[string]struct{}{}, newPathIndex1.Tokens)

	// For this, we expect no change
	newPathIndex2, err := sscm.LeaseCacheForTest().DBForTest().Get(cachememdb.IndexNameID, pathIndexId2)
	require.Nil(t, err)
	require.Equal(t, newPathIndex2, newPathIndex2)

	// Forcefully stop any remaining workers
	sscm.Stop()
}
