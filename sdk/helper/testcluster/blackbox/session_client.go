// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"os"
	"time"

	"github.com/hashicorp/vault/api"
)

type ClientOpt func(*api.Client)

func WithClientRootNamespace() ClientOpt {
	return func(c *api.Client) {
		c.ClearNamespace()
	}
}

func WithClientParentNamespace() ClientOpt {
	return func(c *api.Client) {
		c.SetNamespace(getParentNamespace())
	}
}

func WithClientTimeout(d time.Duration) ClientOpt {
	return func(c *api.Client) {
		c.SetClientTimeout(d)
	}
}

func (s *Session) Req(fn func(*api.Client) error, opts ...ClientOpt) error {
	s.t.Helper()

	c, err := api.NewClient(s.Client.CloneConfig())
	if err != nil {
		return err
	}

	for _, opt := range opts {
		opt(c)
	}

	return fn(c)
}

// GetParentNamespace returns the namespace from VAULT_NAMESPACE environment variable.
// The blackbox test framework auto-creates a unique child namespace for each test
// (e.g., "admin/bbsdk-xxxxx") for isolation. VAULT_NAMESPACE contains the base namespace
// (e.g., "admin"), which is the parent of the test's namespace.
// Example: VAULT_NAMESPACE="admin" → test runs in "admin/bbsdk-xxxxx" → returns "admin"
// Note: This doesn't traverse the namespace hierarchy - it simply returns VAULT_NAMESPACE,
// which happens to be the parent of the test namespace.
func (s *Session) GetParentNamespace() string {
	return getParentNamespace()
}

func getParentNamespace() string {
	ns := os.Getenv("VAULT_NAMESPACE")
	// If VAULT_NAMESPACE is not set, default to root namespace (empty string).
	// This handles cases where tests run in non-namespaced environments.
	if ns == "" {
		return ""
	}
	return ns
}
