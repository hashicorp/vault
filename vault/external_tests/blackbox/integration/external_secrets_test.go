// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package integration

import (
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// TestSecretsEngineExternalCreate tests creation/setup of secrets engines that require external infrastructure
// These tests are excluded from cloud environments (HCP/Docker) which don't have access to AWS, LDAP servers, etc.
func TestSecretsEngineExternalCreate(t *testing.T) {
	v := blackbox.New(t)

	// Verify we have a healthy cluster first
	v.AssertClusterHealthy()

	t.Run("AWSSecrets", func(t *testing.T) {
		testAWSSecretsCreate(t, v)
	})

	t.Run("LDAPSecrets", func(t *testing.T) {
		testLDAPSecretsCreate(t, v)
	})

	t.Run("KMIPSecrets", func(t *testing.T) {
		testKMIPSecretsCreate(t, v)
	})
}

// TestSecretsEngineExternalRead tests read operations for secrets engines that require external infrastructure
func TestSecretsEngineExternalRead(t *testing.T) {
	v := blackbox.New(t)

	// Verify we have a healthy cluster first
	v.AssertClusterHealthy()

	t.Run("AWSSecrets", func(t *testing.T) {
		testAWSSecretsRead(t, v)
	})

	t.Run("LDAPSecrets", func(t *testing.T) {
		testLDAPSecretsRead(t, v)
	})

	t.Run("KMIPSecrets", func(t *testing.T) {
		testKMIPSecretsRead(t, v)
	})
}

// TestSecretsEngineExternalDelete tests delete operations for secrets engines that require external infrastructure
func TestSecretsEngineExternalDelete(t *testing.T) {
	v := blackbox.New(t)

	// Verify we have a healthy cluster first
	v.AssertClusterHealthy()

	t.Run("LDAPSecrets", func(t *testing.T) {
		testLDAPSecretsDelete(t, v)
	})
}
