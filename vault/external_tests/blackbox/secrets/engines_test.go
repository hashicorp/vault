// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package secrets

import (
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// TestSecretsEngineCreate tests creation/setup of various secrets engines
// This test covers secrets engines that work in cloud environments (HCP/Docker)
func TestSecretsEngineCreate(t *testing.T) {
	v := blackbox.New(t)

	// Verify we have a healthy cluster first
	v.AssertClusterHealthy()

	t.Run("KVSecrets", func(t *testing.T) {
		testKVSecretsCreate(t, v)
	})

	t.Run("PKISecrets", func(t *testing.T) {
		testPKISecretsCreate(t, v)
	})

	t.Run("SSHSecrets", func(t *testing.T) {
		testSSHSecretsCreate(t, v)
	})

	t.Run("IdentitySecrets", func(t *testing.T) {
		testIdentitySecretsCreate(t, v)
	})

	t.Run("TransitSecrets", func(t *testing.T) {
		testTransitSecretsCreate(t, v)
	})
}

// TestSecretsEngineRead tests read operations for various secrets engines
// This test covers secrets engines that work in cloud environments (HCP/Docker)
func TestSecretsEngineRead(t *testing.T) {
	v := blackbox.New(t)

	// Verify we have a healthy cluster first
	v.AssertClusterHealthy()

	t.Run("KVSecrets", func(t *testing.T) {
		testKVSecretsRead(t, v)
	})

	t.Run("PKISecrets", func(t *testing.T) {
		testPKISecretsRead(t, v)
	})

	t.Run("SSHSecrets", func(t *testing.T) {
		testSSHSecretsRead(t, v)
	})

	t.Run("IdentitySecrets", func(t *testing.T) {
		testIdentitySecretsRead(t, v)
	})

	t.Run("TransitSecrets", func(t *testing.T) {
		testTransitSecretsRead(t, v)
	})
}

// TestSecretsEngineDelete tests delete operations for various secrets engines
// This test covers secrets engines that work in cloud environments (HCP/Docker)
func TestSecretsEngineDelete(t *testing.T) {
	v := blackbox.New(t)

	// Verify we have a healthy cluster first
	v.AssertClusterHealthy()

	t.Run("KVSecrets", func(t *testing.T) {
		testKVSecretsDelete(t, v)
	})

	t.Run("PKISecrets", func(t *testing.T) {
		testPKISecretsDelete(t, v)
	})

	t.Run("SSHSecrets", func(t *testing.T) {
		testSSHSecretsDelete(t, v)
	})

	t.Run("IdentitySecrets", func(t *testing.T) {
		testIdentitySecretsDelete(t, v)
	})

	t.Run("TransitSecrets", func(t *testing.T) {
		testTransitSecretsDelete(t, v)
	})
}
