// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// Common test data structures
var (
	// Standard KV test data
	StandardKVData = map[string]any{
		"api_key":     "abc123",
		"is_active":   true,
		"retry_count": 3,
	}

	// Alternative KV test data
	AltKVData = map[string]any{
		"username": "testuser",
		"password": "testpass123",
		"enabled":  true,
	}

	// Standard ops policy for KV access
	StandardOpsPolicy = `
		path "secret/data/*" { capabilities = ["create", "read", "update"] }
		path "secret/delete/*" { capabilities = ["update"] }
		path "secret/undelete/*" { capabilities = ["update"] }
		path "auth/userpass/login/*" { capabilities = ["create", "read"] }
	`

	// Read-only policy for limited access testing
	ReadOnlyPolicy = `
		path "secret/data/allowed/*" { capabilities = ["read"] }
		path "secret/data/denied/*" { capabilities = ["deny"] }
	`
)

// SetupKVEngine enables a KV v2 secrets engine at the given mount point and waits for it to be ready
func SetupKVEngine(v *blackbox.Session, mountPath string) {
	v.MustEnableSecretsEngine(mountPath, &api.MountInput{Type: "kv-v2"})

	// Wait for KV engine to finish upgrading (important for HCP environments)
	WaitForKVEngineReady(v, mountPath)
}

// WaitForKVEngineReady waits for a KV v2 engine to complete its upgrade process
func WaitForKVEngineReady(v *blackbox.Session, mountPath string) {
	maxRetries := 30
	retryDelay := 2 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Try to write a small test value to see if the engine is ready
		testPath := mountPath + "/data/__test_ready__"
		testData := map[string]any{"ready": "test"}

		_, err := v.Client.Logical().Write(testPath, map[string]any{"data": testData})
		if err != nil {
			if attempt < maxRetries {
				// Check if this is the upgrade error we're waiting for
				if strings.Contains(err.Error(), "Waiting for the primary to upgrade") {
					time.Sleep(retryDelay)
					continue
				}
				// Some other error - might still be initializing
				time.Sleep(retryDelay)
				continue
			}
			// Final attempt failed
			v.Client.Logical().Write(testPath, map[string]any{"data": testData}) // Let it fail with proper error handling
		} else {
			// Success! Clean up the test data
			v.Client.Logical().Delete(testPath)
			return
		}
	}
}

// SetupUserpassAuth enables userpass auth and creates a user with the given policy
func SetupUserpassAuth(v *blackbox.Session, username, password, policyName, policyContent string) *blackbox.Session {
	// Enable userpass auth
	v.MustEnableAuth("userpass", &api.EnableAuthOptions{Type: "userpass"})

	// Create policy if content is provided
	if policyContent != "" {
		v.MustWritePolicy(policyName, policyContent)
	}

	// Create user
	v.MustWrite("auth/userpass/users/"+username, map[string]any{
		"password": password,
		"policies": policyName,
	})

	// Try to login and return session (may fail in managed environments)
	userClient, err := v.TryLoginUserpass(username, password)
	if err != nil {
		return nil // Login not available in managed environment
	}
	return userClient
}

// SetupStandardKVUserpass is a convenience function that sets up KV engine + userpass auth with ops policy
func SetupStandardKVUserpass(v *blackbox.Session, kvMount, username, password string) *blackbox.Session {
	// Setup KV engine
	SetupKVEngine(v, kvMount)

	// Setup userpass with standard ops policy
	return SetupUserpassAuth(v, username, password, "ops-policy", StandardOpsPolicy)
}

// AssertKVData verifies standard KV data structure
func AssertKVData(t *testing.T, v *blackbox.Session, secret *api.Secret, data map[string]any) {
	t.Helper()
	assertions := v.AssertSecret(secret).KV2()

	for key, expectedValue := range data {
		assertions.HasKey(key, expectedValue)
	}
}

// CreateTestToken creates a token with the given options
func CreateTestToken(v *blackbox.Session, policies []string, ttl string) string {
	return v.MustCreateToken(blackbox.TokenOptions{
		Policies:    policies,
		TTL:         ttl,
		NoParent:    true,
		DisplayName: "test-token",
	})
}
