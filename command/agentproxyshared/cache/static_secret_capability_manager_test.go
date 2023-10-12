// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"testing"

	"github.com/hashicorp/vault/api"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/vault/helper/testhelpers/minimal"
)

// TestGetCapabilitiesRootToken tests the getCapabilities method with the root
// token, expecting to get "root" capabilities on valid paths
func TestGetCapabilitiesRootToken(t *testing.T) {
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	capabilitiesToCheck := []string{"/auth/token/create", "/sys/health"}
	capabilities, err := getCapabilities(capabilitiesToCheck, client)
	require.NoError(t, err)

	expectedCapabilities := map[string][]string{
		"/auth/token/create": {"root"},
		"/sys/health":        {"root"},
	}
	require.Equal(t, expectedCapabilities, capabilities)
}

// TestGetCapabilitiesLowPrivilegeToken tests the getCapabilities method with
// a low privilege token, expecting to get deny or non-root capabilities
func TestGetCapabilitiesLowPrivilegeToken(t *testing.T) {
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

	capabilitiesToCheck := []string{"/auth/token/create", "/sys/capabilities-self", "/auth/token/lookup-self"}
	capabilities, err := getCapabilities(capabilitiesToCheck, client)
	require.NoError(t, err)

	expectedCapabilities := map[string][]string{
		"/auth/token/create":      {"deny"},
		"/sys/capabilities-self":  {"update"},
		"/auth/token/lookup-self": {"read"},
	}
	require.Equal(t, expectedCapabilities, capabilities)
}

// TestGetCapabilitiesBadClientToken tests that getCapabilities
// returns an error if the client token is bad.
func TestGetCapabilitiesBadClientToken(t *testing.T) {
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client
	client.SetToken("")

	capabilitiesToCheck := []string{"/auth/token/create", "/sys/capabilities-self", "/auth/token/lookup-self"}
	_, err := getCapabilities(capabilitiesToCheck, client)
	require.Error(t, err)
}

// TestGetCapabilitiesEmptyPaths tests the getCapabilities will error on an empty
// set of paths to check
func TestGetCapabilitiesEmptyPaths(t *testing.T) {
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	var capabilitiesToCheck []string
	_, err := getCapabilities(capabilitiesToCheck, client)
	require.Error(t, err)
}

// TestReconcileCapabilities tests that reconcileCapabilities will
// correctly previously remove readable paths that we don't have read access to.
func TestReconcileCapabilities(t *testing.T) {
	paths := []string{"/auth/token/create", "/sys/capabilities-self", "/auth/token/lookup-self"}
	capabilities := map[string][]string{
		"/auth/token/create":      {"deny"},
		"/sys/capabilities-self":  {"update"},
		"/auth/token/lookup-self": {"read"},
	}

	updatedCapabilities := reconcileCapabilities(paths, capabilities)
	expectedUpdatedCapabilities := map[string]struct{}{
		"/auth/token/lookup-self": {},
	}
	require.Equal(t, expectedUpdatedCapabilities, updatedCapabilities)
}

// TestReconcileCapabilitiesNoOp tests that reconcileCapabilities will
// correctly not remove capabilities when they all remain readable.
func TestReconcileCapabilitiesNoOp(t *testing.T) {
	paths := []string{"/foo/bar", "/bar/baz", "/baz/foo"}
	capabilities := map[string][]string{
		"/foo/bar": {"read"},
		"/bar/baz": {"root"},
		"/baz/foo": {"read"},
	}

	updatedCapabilities := reconcileCapabilities(paths, capabilities)
	expectedUpdatedCapabilities := map[string]struct{}{
		"/foo/bar": {},
		"/bar/baz": {},
		"/baz/foo": {},
	}
	require.Equal(t, expectedUpdatedCapabilities, updatedCapabilities)
}
