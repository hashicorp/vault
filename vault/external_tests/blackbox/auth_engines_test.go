// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// TestAuthEngineCreate tests creation/setup of various auth engines
func TestAuthEngineCreate(t *testing.T) {
	v := blackbox.New(t)

	// Verify we have a healthy cluster first
	v.AssertClusterHealthy()

	t.Run("UserpassAuth", func(t *testing.T) {
		testUserpassAuthCreate(t, v)
	})

	t.Run("LDAPAuth", func(t *testing.T) {
		t.Skip("LDAP auth engine create test - implementation pending")
	})

	t.Run("OIDCAuth", func(t *testing.T) {
		t.Skip("OIDC auth engine create test - implementation pending")
	})

	t.Run("AWSAuth", func(t *testing.T) {
		t.Skip("AWS auth engine create test - implementation pending")
	})

	t.Run("KubernetesAuth", func(t *testing.T) {
		t.Skip("Kubernetes auth engine create test - implementation pending")
	})

	t.Run("AppRoleAuth", func(t *testing.T) {
		t.Skip("AppRole auth engine create test - implementation pending")
	})

	t.Run("CertAuth", func(t *testing.T) {
		t.Skip("Cert auth engine create test - implementation pending")
	})
}

// TestAuthEngineRead tests read operations for various auth engines
func TestAuthEngineRead(t *testing.T) {
	v := blackbox.New(t)

	// Verify we have a healthy cluster first
	v.AssertClusterHealthy()

	t.Run("UserpassAuth", func(t *testing.T) {
		testUserpassAuthRead(t, v)
	})

	t.Run("LDAPAuth", func(t *testing.T) {
		t.Skip("LDAP auth engine read test - implementation pending")
	})

	t.Run("OIDCAuth", func(t *testing.T) {
		t.Skip("OIDC auth engine read test - implementation pending")
	})

	t.Run("AWSAuth", func(t *testing.T) {
		t.Skip("AWS auth engine read test - implementation pending")
	})

	t.Run("KubernetesAuth", func(t *testing.T) {
		t.Skip("Kubernetes auth engine read test - implementation pending")
	})

	t.Run("AppRoleAuth", func(t *testing.T) {
		t.Skip("AppRole auth engine read test - implementation pending")
	})

	t.Run("CertAuth", func(t *testing.T) {
		t.Skip("Cert auth engine read test - implementation pending")
	})
}

// TestAuthEngineDelete tests delete operations for various auth engines
func TestAuthEngineDelete(t *testing.T) {
	v := blackbox.New(t)

	// Verify we have a healthy cluster first
	v.AssertClusterHealthy()

	t.Run("UserpassAuth", func(t *testing.T) {
		testUserpassAuthDelete(t, v)
	})

	t.Run("LDAPAuth", func(t *testing.T) {
		t.Skip("LDAP auth engine delete test - implementation pending")
	})

	t.Run("OIDCAuth", func(t *testing.T) {
		t.Skip("OIDC auth engine delete test - implementation pending")
	})

	t.Run("AWSAuth", func(t *testing.T) {
		t.Skip("AWS auth engine delete test - implementation pending")
	})

	t.Run("KubernetesAuth", func(t *testing.T) {
		t.Skip("Kubernetes auth engine delete test - implementation pending")
	})

	t.Run("AppRoleAuth", func(t *testing.T) {
		t.Skip("AppRole auth engine delete test - implementation pending")
	})

	t.Run("CertAuth", func(t *testing.T) {
		t.Skip("Cert auth engine delete test - implementation pending")
	})
}
