// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/internalshared/agentproxyshared/cache/cachememdb"
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
	updater.workerPoolLock.RLock()
	require.NotNil(t, updater.workerPool)
	updater.workerPoolLock.RUnlock()
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
	updater.workerPoolLock.RLock()
	require.NotNil(t, updater.workerPool)
	updater.workerPoolLock.RUnlock()
	require.NotNil(t, updater.staticSecretTokenCapabilityRefreshInterval)
	require.NotNil(t, updater.client)
	require.NotNil(t, updater.leaseCache)
	require.NotNil(t, updater.logger)
	require.Equal(t, time.Hour, updater.staticSecretTokenCapabilityRefreshInterval)
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
