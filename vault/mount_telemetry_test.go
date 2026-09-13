// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// TestNewLogicalBackend_MountTelemetryConfig verifies that a mounted secrets
// engine is told its own mount path, and whether the operator opted in to
// per-mount telemetry labels via add_mount_point_database_metrics.
//
// This covers the wiring between the telemetry configuration and the backend,
// which is what allows the database secrets engine to label its metrics.
func TestNewLogicalBackend_MountTelemetryConfig(t *testing.T) {
	testCases := []struct {
		name         string
		includeMount bool
	}{
		{
			name:         "opted in",
			includeMount: true,
		},
		{
			name:         "opted out",
			includeMount: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// captured records the BackendConfig the core handed to the factory.
			var captured *logical.BackendConfig

			captureFactory := func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
				captured = config
				return LeasedPassthroughBackendFactory(ctx, config)
			}

			core, _, _ := TestCoreUnsealedWithConfig(t, &CoreConfig{
				LogicalBackends: map[string]logical.Factory{
					"capture": captureFactory,
				},
			})

			// The flag is read when the backend is constructed, so it must be
			// set before mounting.
			core.metricSink.TelemetryConsts.DatabaseMetricsIncludeMountPoint = tc.includeMount

			const mountPath = "telemetry-test/"
			err := core.mount(namespace.RootContext(nil), &MountEntry{
				Table: mountTableType,
				Path:  mountPath,
				Type:  "capture",
			})
			require.NoError(t, err)

			require.NotNil(t, captured, "the backend factory was never called")
			require.Equal(t, mountPath, captured.MountPath)
			require.Equal(t, "root", captured.MountNamespace)
			require.Equal(t, tc.includeMount, captured.IncludeMountPointInMetrics)
		})
	}
}
