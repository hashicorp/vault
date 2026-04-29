// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package integration

import (
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
	helpers "github.com/hashicorp/vault/vault/external_tests/blackbox"
)

// TestAuthTokenIntegration tests token operations requiring elevated permissions
// These tests are excluded from cloud environments (HCP/Docker) which don't have
// the necessary permissions to create orphaned tokens.
func TestAuthTokenIntegration(t *testing.T) {
	v := blackbox.New(t)

	// Verify we have a healthy cluster first
	v.AssertClusterHealthy()

	t.Run("OrphanedWithPolicy", func(t *testing.T) {
		testTokenOrphanedWithPolicy(t, v)
	})
}

// testTokenOrphanedWithPolicy verifies token creation with policy assignment,
// validates token authentication, and tests policy enforcement by attempting
// both allowed and denied operations on KV secrets.
func testTokenOrphanedWithPolicy(t *testing.T, v *blackbox.Session) {
	// Use common utility to create token with read-only policy
	policyName := "read-secret-only"
	v.MustWritePolicy(policyName, helpers.ReadOnlyPolicy)

	token := helpers.CreateTestToken(v, []string{policyName}, "15m")
	t.Logf("Generated Token: %s...", token[:5])

	v.AssertTokenIsValid(token, policyName)

	// Setup KV engine and seed test data
	helpers.SetupKVEngine(v, "secret")
	v.MustWriteKV2("secret", "allowed/test", map[string]any{"val": "allowed"})
	v.MustWriteKV2("secret", "denied/test", map[string]any{"val": "denied"})

	// Test token access
	userClient := v.NewClientFromToken(token)
	secret := userClient.MustReadRequired("secret/data/allowed/test")
	userClient.AssertSecret(secret).KV2().HasKey("val", "allowed")
	userClient.AssertReadFails("secret/data/denied/test")
}
