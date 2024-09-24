// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agentproxyshared/cache/cachememdb"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/stretchr/testify/require"
)

// testNewStaticSecretCapabilityManager returns a new StaticSecretCapabilityManager
// for use in tests.
func testNewStaticSecretCapabilityManager(t *testing.T, client *api.Client) *StaticSecretCapabilityManager {
	t.Helper()

	lc := testNewLeaseCache(t, []*SendResponse{})

	updater, err := NewStaticSecretCapabilityManager(&StaticSecretCapabilityManagerConfig{
		LeaseCache: lc,
		Logger:     logging.NewVaultLogger(hclog.Trace).Named("cache.capabilitiesmanager"),
		Client:     client,
		StaticSecretTokenCapabilityRefreshInterval: 250 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}

	return updater
}

// TestNewStaticSecretCapabilityManager tests the NewStaticSecretCapabilityManager method,
// to ensure it errors out when appropriate.
func TestNewStaticSecretCapabilityManager(t *testing.T) {
	t.Parallel()

	lc := testNewLeaseCache(t, []*SendResponse{})
	logger := logging.NewVaultLogger(hclog.Trace).Named("cache.capabilitiesmanager")
	client, err := api.NewClient(api.DefaultConfig())
	require.Nil(t, err)

	// Expect an error if any of the arguments are nil:
	updater, err := NewStaticSecretCapabilityManager(&StaticSecretCapabilityManagerConfig{
		LeaseCache: nil,
		Logger:     logger,
		Client:     client,
	})
	require.Error(t, err)
	require.Nil(t, updater)

	updater, err = NewStaticSecretCapabilityManager(&StaticSecretCapabilityManagerConfig{
		LeaseCache: lc,
		Logger:     nil,
		Client:     client,
	})
	require.Error(t, err)
	require.Nil(t, updater)

	updater, err = NewStaticSecretCapabilityManager(&StaticSecretCapabilityManagerConfig{
		LeaseCache: lc,
		Logger:     logger,
		Client:     nil,
	})
	require.Error(t, err)
	require.Nil(t, updater)

	// Don't expect an error if the arguments are as expected
	updater, err = NewStaticSecretCapabilityManager(&StaticSecretCapabilityManagerConfig{
		LeaseCache: lc,
		Logger:     logging.NewVaultLogger(hclog.Trace).Named("cache.capabilitiesmanager"),
		Client:     client,
	})
	if err != nil {
		t.Fatal(err)
	}
	require.NotNil(t, updater)
	require.NotNil(t, updater.workerPool)
	require.NotNil(t, updater.staticSecretTokenCapabilityRefreshInterval)
	require.NotNil(t, updater.client)
	require.NotNil(t, updater.leaseCache)
	require.NotNil(t, updater.logger)
	require.Equal(t, DefaultStaticSecretTokenCapabilityRefreshInterval, updater.staticSecretTokenCapabilityRefreshInterval)

	// Lastly, double check that the refresh interval can be properly set
	updater, err = NewStaticSecretCapabilityManager(&StaticSecretCapabilityManagerConfig{
		LeaseCache: lc,
		Logger:     logging.NewVaultLogger(hclog.Trace).Named("cache.capabilitiesmanager"),
		Client:     client,
		StaticSecretTokenCapabilityRefreshInterval: time.Hour,
	})
	if err != nil {
		t.Fatal(err)
	}
	require.NotNil(t, updater)
	require.NotNil(t, updater.workerPool)
	require.NotNil(t, updater.staticSecretTokenCapabilityRefreshInterval)
	require.NotNil(t, updater.client)
	require.NotNil(t, updater.leaseCache)
	require.NotNil(t, updater.logger)
	require.Equal(t, time.Hour, updater.staticSecretTokenCapabilityRefreshInterval)
}

// TestGetCapabilitiesRootToken tests the getCapabilities method with the root
// token, expecting to get "root" capabilities on valid paths
func TestGetCapabilitiesRootToken(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	capabilitiesToCheck := []string{"auth/token/create", "sys/health"}
	capabilities, err := getCapabilities(capabilitiesToCheck, client)
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
	capabilities, err := getCapabilities(capabilitiesToCheck, client)
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
	capabilities, err := getCapabilities(capabilitiesToCheck, client)
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
	_, err := getCapabilities(capabilitiesToCheck, client)
	require.Error(t, err)
}

// TestReconcileCapabilities tests that reconcileCapabilities will
// correctly previously remove readable paths that we don't have read access to.
func TestReconcileCapabilities(t *testing.T) {
	t.Parallel()
	paths := []string{"auth/token/create", "sys/capabilities-self", "auth/token/lookup-self"}
	capabilities := map[string][]string{
		"auth/token/create":      {"deny"},
		"sys/capabilities-self":  {"update"},
		"auth/token/lookup-self": {"read"},
	}

	updatedCapabilities := reconcileCapabilities(paths, capabilities)
	expectedUpdatedCapabilities := map[string]struct{}{
		"auth/token/lookup-self": {},
	}
	require.Equal(t, expectedUpdatedCapabilities, updatedCapabilities)
}

// TestReconcileCapabilitiesNoOp tests that reconcileCapabilities will
// correctly not remove capabilities when they all remain readable.
func TestReconcileCapabilitiesNoOp(t *testing.T) {
	t.Parallel()
	paths := []string{"foo/bar", "bar/baz", "baz/foo"}
	capabilities := map[string][]string{
		"foo/bar": {"read"},
		"bar/baz": {"root"},
		"baz/foo": {"read"},
	}

	updatedCapabilities := reconcileCapabilities(paths, capabilities)
	expectedUpdatedCapabilities := map[string]struct{}{
		"foo/bar": {},
		"bar/baz": {},
		"baz/foo": {},
	}
	require.Equal(t, expectedUpdatedCapabilities, updatedCapabilities)
}

// TestReconcileCapabilitiesNoAdding tests that reconcileCapabilities will
// not add any capabilities that weren't present in the first argument to the function
func TestReconcileCapabilitiesNoAdding(t *testing.T) {
	t.Parallel()
	paths := []string{"auth/token/create", "sys/capabilities-self", "auth/token/lookup-self"}
	capabilities := map[string][]string{
		"auth/token/create":      {"deny"},
		"sys/capabilities-self":  {"update"},
		"auth/token/lookup-self": {"read"},
		"some/new/path":          {"read"},
	}

	updatedCapabilities := reconcileCapabilities(paths, capabilities)
	expectedUpdatedCapabilities := map[string]struct{}{
		"auth/token/lookup-self": {},
	}
	require.Equal(t, expectedUpdatedCapabilities, updatedCapabilities)
}

// TestSubmitWorkNoOp tests that we will gracefully end if the capabilities index
// does not exist in the cache
func TestSubmitWorkNoOp(t *testing.T) {
	t.Parallel()
	client, err := api.NewClient(api.DefaultConfig())
	require.Nil(t, err)
	sscm := testNewStaticSecretCapabilityManager(t, client)
	// This index will be a no-op, as this does not exist in the cache
	index := &cachememdb.CapabilitiesIndex{
		ID: "test",
	}
	sscm.StartRenewingCapabilities(index)

	// Wait for the job to complete...
	time.Sleep(1 * time.Second)
	require.Equal(t, 0, sscm.workerPool.WaitingQueueSize())
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
	indexId := hashStaticSecretIndex(token)

	sscm := testNewStaticSecretCapabilityManager(t, client)
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
	err = sscm.leaseCache.db.SetCapabilitiesIndex(index)
	require.Nil(t, err)

	sscm.StartRenewingCapabilities(index)

	// Wait for the job to complete at least once...
	time.Sleep(3 * time.Second)

	newIndex, err := sscm.leaseCache.db.GetCapabilitiesIndex(cachememdb.IndexNameID, indexId)
	require.Nil(t, err)
	newIndex.IndexLock.RLock()
	require.Equal(t, map[string]struct{}{
		"auth/token/lookup-self": {},
	}, newIndex.ReadablePaths)
	newIndex.IndexLock.RUnlock()

	// Forcefully stop any remaining workers
	sscm.workerPool.Stop()
}

// TestSubmitWorkUpdatesIndexWithBadToken tests that an index will be correctly updated if the token
// has expired and we cannot access the sys capabilities endpoint.
func TestSubmitWorkUpdatesIndexWithBadToken(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	token := "not real token"
	indexId := hashStaticSecretIndex(token)

	sscm := testNewStaticSecretCapabilityManager(t, client)
	index := &cachememdb.CapabilitiesIndex{
		ID:    indexId,
		Token: token,
		ReadablePaths: map[string]struct{}{
			"foo/bar":                {},
			"auth/token/lookup-self": {},
		},
	}
	err := sscm.leaseCache.db.SetCapabilitiesIndex(index)
	require.Nil(t, err)

	sscm.StartRenewingCapabilities(index)

	// Wait for the job to complete at least once...
	time.Sleep(3 * time.Second)

	// This entry should be evicted.
	newIndex, err := sscm.leaseCache.db.GetCapabilitiesIndex(cachememdb.IndexNameID, indexId)
	require.Equal(t, err, cachememdb.ErrCacheItemNotFound)
	require.Nil(t, newIndex)

	// Forcefully stop any remaining workers
	sscm.workerPool.Stop()
}

// TestSubmitWorkSealedVaultOptimistic tests that the capability manager
// behaves as expected when
// sscm.tokenCapabilityRefreshBehaviour == TokenCapabilityRefreshBehaviourOptimistic
func TestSubmitWorkSealedVaultOptimistic(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	token := "not real token"
	indexId := hashStaticSecretIndex(token)

	sscm := testNewStaticSecretCapabilityManager(t, client)
	index := &cachememdb.CapabilitiesIndex{
		ID:    indexId,
		Token: token,
		ReadablePaths: map[string]struct{}{
			"foo/bar":                {},
			"auth/token/lookup-self": {},
		},
	}
	err := sscm.leaseCache.db.SetCapabilitiesIndex(index)
	require.Nil(t, err)

	// Seal the cluster
	cluster.EnsureCoresSealed(t)

	sscm.StartRenewingCapabilities(index)

	// Wait for the job to complete at least once...
	time.Sleep(3 * time.Second)

	// This entry should not be evicted.
	newIndex, err := sscm.leaseCache.db.GetCapabilitiesIndex(cachememdb.IndexNameID, indexId)
	require.NoError(t, err)
	require.NotNil(t, newIndex)

	// Forcefully stop any remaining workers
	sscm.workerPool.Stop()
}

// TestSubmitWorkSealedVaultPessimistic tests that the capability manager
// behaves as expected when
// sscm.tokenCapabilityRefreshBehaviour == TokenCapabilityRefreshBehaviourPessimistic
func TestSubmitWorkSealedVaultPessimistic(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	token := "not real token"
	indexId := hashStaticSecretIndex(token)

	sscm := testNewStaticSecretCapabilityManager(t, client)
	sscm.tokenCapabilityRefreshBehaviour = TokenCapabilityRefreshBehaviourPessimistic

	index := &cachememdb.CapabilitiesIndex{
		ID:    indexId,
		Token: token,
		ReadablePaths: map[string]struct{}{
			"foo/bar":                {},
			"auth/token/lookup-self": {},
		},
	}
	err := sscm.leaseCache.db.SetCapabilitiesIndex(index)
	require.Nil(t, err)

	// Seal the cluster
	cluster.EnsureCoresSealed(t)

	sscm.StartRenewingCapabilities(index)

	// Wait for the job to complete at least once...
	time.Sleep(3 * time.Second)

	// This entry should be evicted.
	newIndex, err := sscm.leaseCache.db.GetCapabilitiesIndex(cachememdb.IndexNameID, indexId)
	require.Error(t, err)
	require.Nil(t, newIndex)

	// Forcefully stop any remaining workers
	sscm.workerPool.Stop()
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
	indexId := hashStaticSecretIndex(token)

	sscm := testNewStaticSecretCapabilityManager(t, client)
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
	err = sscm.leaseCache.db.SetCapabilitiesIndex(index)
	require.Nil(t, err)

	pathIndexId1 := hashStaticSecretIndex("foo/bar")
	pathIndex1 := &cachememdb.Index{
		ID:        pathIndexId1,
		Namespace: "root/",
		Tokens: map[string]struct{}{
			token: {},
		},
		RequestPath: "foo/bar",
		Response:    []byte{},
	}

	pathIndexId2 := hashStaticSecretIndex("auth/token/lookup-self")
	pathIndex2 := &cachememdb.Index{
		ID:        pathIndexId2,
		Namespace: "root/",
		Tokens: map[string]struct{}{
			token: {},
		},
		RequestPath: "auth/token/lookup-self",
		Response:    []byte{},
	}

	err = sscm.leaseCache.db.Set(pathIndex1)
	require.Nil(t, err)

	err = sscm.leaseCache.db.Set(pathIndex2)
	require.Nil(t, err)

	sscm.StartRenewingCapabilities(index)

	// Wait for the job to complete at least once...
	time.Sleep(1 * time.Second)

	newIndex, err := sscm.leaseCache.db.GetCapabilitiesIndex(cachememdb.IndexNameID, indexId)
	require.Nil(t, err)
	newIndex.IndexLock.RLock()
	require.Equal(t, map[string]struct{}{
		"auth/token/lookup-self": {},
	}, newIndex.ReadablePaths)
	newIndex.IndexLock.RUnlock()

	// For this, we expect the token to have been deleted
	newPathIndex1, err := sscm.leaseCache.db.Get(cachememdb.IndexNameID, pathIndexId1)
	require.Nil(t, err)
	require.Equal(t, map[string]struct{}{}, newPathIndex1.Tokens)

	// For this, we expect no change
	newPathIndex2, err := sscm.leaseCache.db.Get(cachememdb.IndexNameID, pathIndexId2)
	require.Nil(t, err)
	require.Equal(t, newPathIndex2, newPathIndex2)

	// Forcefully stop any remaining workers
	sscm.workerPool.Stop()
}
