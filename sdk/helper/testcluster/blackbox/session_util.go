// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"os"
	"time"

	"github.com/hashicorp/vault/api"
)

// Eventually retries the function 'fn' until it returns nil or timeout occurs.
func (s *Session) Eventually(fn func() error) {
	s.t.Helper()

	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	var lastErr error

	for {
		select {
		case <-timeout:
			s.t.Fatalf("Eventually failed after 5s. Last error: %v", lastErr)
		case <-ticker.C:
			lastErr = fn()
			if lastErr == nil {
				return
			}
		}
	}
}

func (s *Session) WithRootNamespace(fn func() (*api.Secret, error)) (*api.Secret, error) {
	s.t.Helper()

	oldNamespace := s.Client.Namespace()
	defer s.Client.SetNamespace(oldNamespace)
	s.Client.ClearNamespace()

	return fn()
}

// WithParentNamespace temporarily switches to the parent namespace (e.g., "admin" in HVD)
// and executes the provided function, then restores the original namespace.
func (s *Session) WithParentNamespace(fn func() (*api.Secret, error)) (*api.Secret, error) {
	s.t.Helper()

	oldNamespace := s.Client.Namespace()
	defer s.Client.SetNamespace(oldNamespace)

	// Get the parent namespace from environment (e.g., "admin" in HVD)
	parentNS := s.GetParentNamespace()
	s.Client.SetNamespace(parentNS)

	return fn()
}

// GetParentNamespace returns the namespace from VAULT_NAMESPACE environment variable.
// The blackbox test framework auto-creates a unique child namespace for each test
// (e.g., "admin/bbsdk-xxxxx") for isolation. VAULT_NAMESPACE contains the base namespace
// (e.g., "admin"), which is the parent of the test's namespace.
// Example: VAULT_NAMESPACE="admin" → test runs in "admin/bbsdk-xxxxx" → returns "admin"
// Note: This doesn't traverse the namespace hierarchy - it simply returns VAULT_NAMESPACE,
// which happens to be the parent of the test namespace.
func (s *Session) GetParentNamespace() string {
	ns := os.Getenv("VAULT_NAMESPACE")
	// If VAULT_NAMESPACE is not set, default to root namespace (empty string).
	// This handles cases where tests run in non-namespaced environments.
	if ns == "" {
		return ""
	}
	return ns
}
