// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package expiration

import (
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

// TestSysRenew verifies that a lease on a leased secret can be renewed through
// both the legacy sys/renew and the sys/leases/renew endpoints without the
// lease ID changing.
func TestSysRenew(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": vault.LeasedPassthroughBackendFactory,
		},
	})
	client := cluster.Cores[0].Client

	err := client.Sys().Mount("secret", &api.MountInput{Type: "kv"})
	require.NoError(t, err)

	_, err = client.Logical().Write("secret/foo", map[string]interface{}{
		"data":  "bar",
		"lease": "1h",
	})
	require.NoError(t, err)

	secret, err := client.Logical().Read("secret/foo")
	require.NoError(t, err)
	require.NotNil(t, secret)
	leaseID := secret.LeaseID
	require.NotEmpty(t, leaseID)

	for _, path := range []string{"sys/renew/", "sys/leases/renew/"} {
		renewed, err := client.Logical().Write(path+leaseID, nil)
		require.NoError(t, err, path)
		require.NotNil(t, renewed, path)
		require.Equal(t, leaseID, renewed.LeaseID, path)
	}
}

// TestSysRevokeAndRevokePrefix verifies that sys/revoke revokes a single lease
// and sys/revoke-prefix revokes only the leases under the given prefix, and
// that revoking leases that don't exist succeeds, as revocation is idempotent.
func TestSysRevokeAndRevokePrefix(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": vault.LeasedPassthroughBackendFactory,
		},
	})
	client := cluster.Cores[0].Client

	err := client.Sys().Mount("secret", &api.MountInput{Type: "kv"})
	require.NoError(t, err)

	for _, name := range []string{"foo", "bar"} {
		_, err = client.Logical().Write("secret/"+name, map[string]interface{}{
			"data":  name,
			"lease": "1h",
		})
		require.NoError(t, err)
	}

	readLease := func(name string) string {
		t.Helper()
		secret, err := client.Logical().Read("secret/" + name)
		require.NoError(t, err)
		require.NotNil(t, secret)
		require.NotEmpty(t, secret.LeaseID)
		return secret.LeaseID
	}
	leaseExists := func(leaseID string) bool {
		t.Helper()
		secret, err := client.Logical().Write("sys/leases/lookup", map[string]interface{}{"lease_id": leaseID})
		return err == nil && secret != nil
	}

	// Revoke a single lease and verify only that lease is gone.
	revokedLease := readLease("foo")
	otherLease := readLease("foo")
	require.NotEqual(t, revokedLease, otherLease)
	require.True(t, leaseExists(revokedLease))
	require.True(t, leaseExists(otherLease))

	_, err = client.Logical().Write("sys/revoke/"+revokedLease, nil)
	require.NoError(t, err)
	require.Eventually(t, func() bool { return !leaseExists(revokedLease) }, 5*time.Second, 50*time.Millisecond)
	require.True(t, leaseExists(otherLease))

	// Revoke by prefix and verify leases outside the prefix are untouched.
	fooLease := readLease("foo")
	barLease := readLease("bar")
	require.True(t, leaseExists(fooLease))
	require.True(t, leaseExists(barLease))

	_, err = client.Logical().Write("sys/revoke-prefix/secret/foo", nil)
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		return !leaseExists(fooLease) && !leaseExists(otherLease)
	}, 5*time.Second, 50*time.Millisecond)
	require.True(t, leaseExists(barLease))

	// Revoking nonexistent leases is not an error.
	_, err = client.Logical().Write("sys/revoke/secret/foo/1234", nil)
	require.NoError(t, err)

	_, err = client.Logical().Write("sys/revoke-prefix/secret/foo/1234", nil)
	require.NoError(t, err)
}
