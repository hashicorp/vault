//go:build isolated
// +build isolated

// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package verify

// TestDefaultLCQ verifies that the default lease count quota (max_leases) is set
// to the expected value after an autopilot upgrade.
//
// This test converts the enos module vault_verify_default_lcq:
//   enos/modules/vault_verify_default_lcq/scripts/smoke-verify-default-lcq.sh
//
// Shell logic translated:
//  1. [[ -z "$DEFAULT_LCQ" ]] && exit 0          → t.Skip when env var is empty
//  2. GET /v1/sys/quotas/lease-count/default       → v.Req(..., WithClientRootNamespace())
//  3. jq '.data.max_leases' == DEFAULT_LCQ         → numeric comparison with require
//  4. while timeout loop with RETRY_INTERVAL sleep → v.EventuallyWithTimeout
//
// Required env vars (injected by the enos autopilot scenario / CI pipeline):
//   DEFAULT_LCQ      — expected max_leases (e.g. "300000"); if empty → skip
//   TIMEOUT_SECONDS  — polling deadline in seconds (default: 60)
//   RETRY_INTERVAL   — documented intent; EventuallyWithTimeout polls at 200ms fixed interval

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
	"github.com/stretchr/testify/require"
)

func TestDefaultLCQ(t *testing.T) {
	t.Parallel()

	// Shell: [[ -z "$DEFAULT_LCQ" ]] && exit 0
	// The enos autopilot scenario sets DEFAULT_LCQ="300000" only when upgrading from >= 1.16.0.
	// An empty value means the initial version was < 1.16.0 and no default LCQ is expected.
	defaultLCQStr := os.Getenv("DEFAULT_LCQ")
	if defaultLCQStr == "" {
		t.Skip("DEFAULT_LCQ not set — initial version is < 1.16.0, no default LCQ expected")
	}

	expectedMaxLeases, err := strconv.Atoi(defaultLCQStr)
	require.NoError(t, err, "DEFAULT_LCQ must be a valid integer, got: %q", defaultLCQStr)

	// Shell: TIMEOUT_SECONDS (module default: 60)
	timeoutSeconds := 60
	if v := os.Getenv("TIMEOUT_SECONDS"); v != "" {
		timeoutSeconds, err = strconv.Atoi(v)
		require.NoError(t, err, "TIMEOUT_SECONDS must be a valid integer, got: %q", v)
	}
	timeout := time.Duration(timeoutSeconds) * time.Second

	// Shell: RETRY_INTERVAL (module default: 2) — noted for documentation purposes.
	// EventuallyWithTimeout polls at a fixed 200ms interval regardless (see session_metrics.go:39).

	v := blackbox.New(t)

	// Shell: while timeout loop → v.EventuallyWithTimeout
	// sys/quotas/lease-count/default is a root-namespace-only path.
	// v.Client is scoped to the test's isolated child namespace; use WithClientRootNamespace()
	// to clear it before the request.
	var secret *api.Secret
	v.EventuallyWithTimeout(func() error {
		return v.Req(func(c *api.Client) error {
			var readErr error
			secret, readErr = c.Logical().Read("sys/quotas/lease-count/default")
			if readErr != nil {
				return readErr
			}
			if secret == nil {
				return fmt.Errorf("nil response from sys/quotas/lease-count/default")
			}
			return nil
		}, blackbox.WithClientRootNamespace())
	}, timeout)

	// Shell: jq '.data.max_leases // empty' — assert key exists.
	require.NotNil(t, secret)
	raw, ok := secret.Data["max_leases"]
	require.True(t, ok, "max_leases key not found in sys/quotas/lease-count/default response")

	// The Vault API may return max_leases as json.Number, float64, or int depending on
	// the client decode path. Normalise to int for comparison.
	var actualMaxLeases int
	switch typedVal := raw.(type) {
	case json.Number:
		i, parseErr := typedVal.Int64()
		require.NoError(t, parseErr, "failed to parse max_leases json.Number as int64")
		actualMaxLeases = int(i)
	case float64:
		actualMaxLeases = int(typedVal)
	case int:
		actualMaxLeases = typedVal
	default:
		t.Fatalf("unexpected type for max_leases: %T (value: %v)", raw, raw)
	}

	// Shell: if [[ "$max_leases" == "$DEFAULT_LCQ" ]]; then exit 0
	require.Equal(t, expectedMaxLeases, actualMaxLeases,
		"expected Default LCQ (max_leases) to be %d but got %d", expectedMaxLeases, actualMaxLeases)

	t.Logf("✓ Default LCQ (max_leases) is %d as expected", actualMaxLeases)
}
