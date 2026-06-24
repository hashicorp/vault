// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

// Session holds the test context and Vault client
type Session struct {
	t         *testing.T
	NoCleanup bool
	Client    *api.Client
	Namespace string
}

func (s *Session) T() *testing.T {
	return s.t
}

type SessionOpts func(s *Session)

func WithNoCleanup() SessionOpts {
	return func(s *Session) {
		s.NoCleanup = true
	}
}

func New(t *testing.T, opts ...SessionOpts) *Session {
	t.Helper()

	addr := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")

	// detect the parent namespace, e.g. "admin" in HVD
	parentNS := os.Getenv("VAULT_NAMESPACE")

	if addr == "" || token == "" {
		t.Fatal("VAULT_ADDR and VAULT_TOKEN are required")
	}

	config := api.DefaultConfig()
	config.Address = addr
	config.Timeout = 120 * time.Second // Increase timeout for LDAP operations that verify service accounts

	privClient, err := api.NewClient(config)
	require.NoError(t, err)
	privClient.SetToken(token)

	// Auto-detect protocol: if we get HTTP/HTTPS mismatch, retry with correct protocol
	if strings.HasPrefix(addr, "http://") {
		// Try HTTP first, but be ready to switch to HTTPS
		testResp, testErr := privClient.Sys().Health()
		if testErr != nil && (strings.Contains(testErr.Error(), "Client sent an HTTP request to an HTTPS server") ||
			strings.Contains(testErr.Error(), "server gave HTTP response to HTTPS client")) {
			// Server is using HTTPS, create completely fresh config
			httpsAddr := strings.Replace(addr, "http://", "https://", 1)

			httpsConfig := api.DefaultConfig()
			httpsConfig.Address = httpsAddr
			httpsConfig.Timeout = 120 * time.Second

			// Disable TLS verification for test environments
			tlsConfig := &api.TLSConfig{
				Insecure: true,
			}
			if err := httpsConfig.ConfigureTLS(tlsConfig); err != nil {
				require.NoError(t, err, "Failed to configure TLS")
			}

			privClient, err = api.NewClient(httpsConfig)
			require.NoError(t, err)
			privClient.SetToken(token)
			t.Logf("Auto-detected HTTPS protocol, switched from %s to %s (TLS verification disabled)", addr, httpsAddr)
		} else if testErr != nil && testResp == nil {
			// Some other error, fail normally
			require.NoError(t, testErr, "Failed to connect to Vault at %s", addr)
		}
	}

	// Use timestamp to ensure uniqueness across test retries
	nsName := fmt.Sprintf("bbsdk-%d-%s", time.Now().UnixNano(), randomString(8))
	nsURLPath := fmt.Sprintf("sys/namespaces/%s", nsName)

	// Try to create the namespace, but if it already exists (from a previous failed run),
	// delete it first and retry
	_, err = privClient.Logical().Write(nsURLPath, nil)
	if err != nil {
		// Check if namespace already exists
		if resp, readErr := privClient.Logical().Read(nsURLPath); readErr == nil && resp != nil {
			t.Logf("RETRY DETECTED: Namespace %s already exists from previous failed test run, cleaning up and retrying", nsName)
			_, delErr := privClient.Logical().Delete(nsURLPath)
			if delErr != nil {
				t.Fatalf("RETRY CLEANUP FAILED: Could not delete existing namespace %s: %v (original creation error: %v)", nsName, delErr, err)
			}
			// Retry creation after deletion
			_, err = privClient.Logical().Write(nsURLPath, nil)
			if err != nil {
				t.Fatalf("RETRY FAILED: Could not create namespace %s after cleanup: %v", nsName, err)
			}
			t.Logf("RETRY SUCCESS: Namespace %s created after cleanup", nsName)
		} else {
			// This is the initial failure, not a retry
			require.NoError(t, err, "INITIAL TEST FAILURE: Failed to create namespace %s. If you see this error followed by a retry with a different error, the cluster state was not properly reset.", nsName)
		}
	}

	// session client should get the full namespace of parent + test
	fullNSPath := nsName
	if parentNS != "" {
		fullNSPath = path.Join(parentNS, nsName)
	}

	sessionConfig := privClient.CloneConfig()
	sessionClient, err := api.NewClient(sessionConfig)
	require.NoError(t, err)
	sessionClient.SetToken(token)
	sessionClient.SetNamespace(fullNSPath)

	session := &Session{
		t:         t,
		Client:    sessionClient,
		Namespace: nsName,
	}

	for opt := range slices.Values(opts) {
		opt(session)
	}

	t.Cleanup(func() {
		if session.NoCleanup {
			t.Logf("WARN: NoCleanup has been set, not cleaning up namespace")
			return
		}
		privClient.SetClientTimeout(time.Second)
		session.Eventually(func() error {
			_, err = privClient.Logical().Delete(nsURLPath)
			return err
		})
		t.Logf("Cleaned up namespace %s", nsName)
	})

	// make sure the namespace has been created
	session.Eventually(func() error {
		// this runs inside the new namespace, so if it succeeds, we're good
		_, err := sessionClient.Auth().Token().LookupSelf()
		return err
	})

	return session
}

func randomString(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

// SkipIfVersionBelow skips the test if the Vault version is below the specified constraint.
// The version is read from the VAULT_VERSION environment variable.
// Example usage: s.SkipIfVersionBelow("2.0.0")
func (s *Session) SkipIfVersionBelow(minVersion string) {
	s.t.Helper()

	vaultVersion := os.Getenv("VAULT_VERSION")
	if vaultVersion == "" {
		s.t.Skip("VAULT_VERSION environment variable not set, skipping version check")
		return
	}

	// Parse the current Vault version
	currentVer, err := version.NewVersion(vaultVersion)
	if err != nil {
		s.t.Fatalf("Failed to parse VAULT_VERSION '%s': %v", vaultVersion, err)
	}

	// Parse the minimum required version
	minVer, err := version.NewVersion(minVersion)
	if err != nil {
		s.t.Fatalf("Invalid minimum version constraint '%s': %v", minVersion, err)
	}

	// Skip if current version is less than minimum required
	if currentVer.LessThan(minVer) {
		s.t.Skipf("Vault version %s is below required version %s", currentVer.String(), minVer.String())
	}
}

// SkipIfVersionAbove skips the test if the Vault version is above the specified constraint.
// The version is read from the VAULT_VERSION environment variable.
// Example usage: s.SkipIfVersionAbove("1.15.0")
func (s *Session) SkipIfVersionAbove(maxVersion string) {
	s.t.Helper()

	vaultVersion := os.Getenv("VAULT_VERSION")
	if vaultVersion == "" {
		s.t.Skip("VAULT_VERSION environment variable not set, skipping version check")
		return
	}

	// Parse the current Vault version
	currentVer, err := version.NewVersion(vaultVersion)
	if err != nil {
		s.t.Skipf("Failed to parse VAULT_VERSION '%s': %v", vaultVersion, err)
		return
	}

	// Parse the maximum allowed version
	maxVer, err := version.NewVersion(maxVersion)
	if err != nil {
		s.t.Fatalf("Invalid maximum version constraint '%s': %v", maxVersion, err)
	}

	// Skip if current version is greater than maximum version
	if currentVer.GreaterThan(maxVer) {
		s.t.Skipf("Test requires Vault version <= %s, but current version is %s", maxVersion, vaultVersion)
	}
}

// SkipIfVersionNotInRange skips the test if the Vault version is not within the specified range.
// The version is read from the VAULT_VERSION environment variable.
// Example usage: s.SkipIfVersionNotInRange("1.15.0", "2.0.0")
func (s *Session) SkipIfVersionNotInRange(minVersion, maxVersion string) {
	s.t.Helper()

	vaultVersion := os.Getenv("VAULT_VERSION")
	if vaultVersion == "" {
		s.t.Skip("VAULT_VERSION environment variable not set, skipping version check")
		return
	}

	// Parse the current Vault version
	currentVer, err := version.NewVersion(vaultVersion)
	if err != nil {
		s.t.Skipf("Failed to parse VAULT_VERSION '%s': %v", vaultVersion, err)
		return
	}

	// Parse the minimum required version
	minVer, err := version.NewVersion(minVersion)
	if err != nil {
		s.t.Fatalf("Invalid minimum version constraint '%s': %v", minVersion, err)
	}

	// Parse the maximum allowed version
	maxVer, err := version.NewVersion(maxVersion)
	if err != nil {
		s.t.Fatalf("Invalid maximum version constraint '%s': %v", maxVersion, err)
	}

	// Skip if current version is outside the range
	if currentVer.LessThan(minVer) || currentVer.GreaterThan(maxVer) {
		s.t.Skipf("Test requires Vault version between %s and %s, but current version is %s", minVersion, maxVersion, vaultVersion)
	}
}
