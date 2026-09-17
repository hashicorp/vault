// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package database

import (
	"context"
	"fmt"
	"testing"
	"time"

	metrics "github.com/hashicorp/go-metrics/compat"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

// installTestMetricsSink installs a fresh in-memory sink as the global metrics
// sink and returns it. The global sink is process-wide, so tests using this must
// not run in parallel with each other.
func installTestMetricsSink(t *testing.T) *metrics.InmemSink {
	t.Helper()

	sink := metrics.NewInmemSink(time.Hour, 2*time.Hour)

	// An empty service name keeps emitted keys unprefixed so assertions can
	// name metrics exactly as the middleware produces them.
	config := metrics.DefaultConfig("")
	config.EnableHostname = false
	config.EnableTypePrefix = false
	config.EnableRuntimeMetrics = false

	_, err := metrics.NewGlobal(config, sink)
	require.NoError(t, err)

	return sink
}

// findMetricLabels returns the labels attached to the named metric, and reports
// whether the metric was emitted at all.
func findMetricLabels(sink *metrics.InmemSink, name string) ([]metrics.Label, bool) {
	for _, interval := range sink.Data() {
		interval.RLock()
		for _, counter := range interval.Counters {
			if counter.Name == name {
				labels := counter.Labels
				interval.RUnlock()
				return labels, true
			}
		}
		interval.RUnlock()
	}
	return nil, false
}

func labelValue(labels []metrics.Label, name string) (string, bool) {
	for _, label := range labels {
		if label.Name == name {
			return label.Value, true
		}
	}
	return "", false
}

// TestBackend_TelemetryMountPointLabels drives the real database secrets engine
// with a distinctive namespace, mount path, and connection name, and verifies
// that the namespace, mount_point, and connection_name labels are attached to
// the emitted metrics only when the operator has opted in.
func TestBackend_TelemetryMountPointLabels(t *testing.T) {
	const (
		mountNamespace = "team-a"
		mountPath      = "database-telemetry-test/"
		connectionName = "telemetry-test-connection"
	)

	testCases := []struct {
		name         string
		includeMount bool
		expectLabels bool
	}{
		{
			name:         "opted in emits mount and connection labels",
			includeMount: true,
			expectLabels: true,
		},
		{
			name:         "opted out emits unlabeled metrics",
			includeMount: false,
			expectLabels: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sink := installTestMetricsSink(t)

			cluster, sys := getCluster(t)
			defer cluster.Cleanup()

			vault.TestAddTestPlugin(t, cluster.Cores[0].Core, "postgresql-database-plugin", consts.PluginTypeDatabase, "", "TestBackend_PluginMain_PostgresMultiplexed",
				[]string{fmt.Sprintf("%s=%s", pluginutil.PluginCACertPEMEnv, cluster.CACertPEMFile)})

			config := logical.TestBackendConfig()
			config.StorageView = &logical.InmemStorage{}
			config.System = sys
			config.MountNamespace = mountNamespace
			config.MountPath = mountPath
			config.IncludeMountPointInMetrics = tc.includeMount

			b, err := Factory(context.Background(), config)
			require.NoError(t, err)
			defer b.Cleanup(context.Background())

			req := &logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "config/" + connectionName,
				Storage:   config.StorageView,
				Data: map[string]interface{}{
					"name":              connectionName,
					"plugin_name":       "postgresql-database-plugin",
					"connection_url":    "some_postgres_url",
					"username":          "postgres",
					"password":          "secret",
					"verify_connection": false,
				},
			}
			resp, err := b.HandleRequest(namespace.RootContext(nil), req)
			require.NoError(t, err)
			require.False(t, resp != nil && resp.IsError(), "connection config write failed: %#v", resp)

			// Initialize runs when the connection is created, so its metrics
			// are the ones to inspect. The metric names must be unchanged
			// regardless of whether labels were added.
			for _, metricName := range []string{"database.Initialize", "database.pgx.Initialize"} {
				labels, found := findMetricLabels(sink, metricName)
				require.True(t, found, "metric %q was not emitted", metricName)

				nsLabel, hasNS := labelValue(labels, "namespace")
				mountLabel, hasMount := labelValue(labels, "mount_point")
				connLabel, hasConn := labelValue(labels, "connection_name")

				if !tc.expectLabels {
					require.False(t, hasNS, "metric %q unexpectedly carried a namespace label", metricName)
					require.False(t, hasMount, "metric %q unexpectedly carried a mount_point label", metricName)
					require.False(t, hasConn, "metric %q unexpectedly carried a connection_name label", metricName)
					require.Empty(t, labels, "metric %q should be unlabeled by default", metricName)
					continue
				}

				require.True(t, hasNS, "metric %q missing namespace label", metricName)
				require.True(t, hasMount, "metric %q missing mount_point label", metricName)
				require.True(t, hasConn, "metric %q missing connection_name label", metricName)
				require.Equal(t, mountNamespace, nsLabel)
				require.Equal(t, mountPath, mountLabel)
				require.Equal(t, connectionName, connLabel)
			}
		})
	}
}

// TestBackend_TelemetryLabelsDisabledByDefault verifies that a backend built
// without any telemetry opt-in produces no mount identifying labels, which is
// the behavior operators upgrading from earlier versions should see.
func TestBackend_TelemetryLabelsDisabledByDefault(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.MountPath = "database/"

	b := Backend(config)
	require.False(t, b.includeMountPointInMetrics)

	labels := b.metricsLabelsForConnection("some-connection")
	require.Equal(t, databaseWrapperMetricsLabels{}, labels)
}

// TestBackend_MetricsLabelsForConnection verifies the gating logic that decides
// whether a plugin instance is given identifying labels.
func TestBackend_MetricsLabelsForConnection(t *testing.T) {
	testCases := []struct {
		name           string
		mountNamespace string
		mountPath      string
		includeMount   bool
		connectionName string
		expected       databaseWrapperMetricsLabels
	}{
		{
			name:           "opted in carries namespace, mount and connection",
			mountNamespace: "root",
			mountPath:      "database/",
			includeMount:   true,
			connectionName: "primary",
			expected: databaseWrapperMetricsLabels{
				namespace:      "root",
				mountPoint:     "database/",
				connectionName: "primary",
			},
		},
		{
			name:           "opted out carries nothing",
			mountNamespace: "root",
			mountPath:      "database/",
			includeMount:   false,
			connectionName: "primary",
			expected:       databaseWrapperMetricsLabels{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := logical.TestBackendConfig()
			config.StorageView = &logical.InmemStorage{}
			config.MountNamespace = tc.mountNamespace
			config.MountPath = tc.mountPath
			config.IncludeMountPointInMetrics = tc.includeMount

			b := Backend(config)
			require.Equal(t, tc.expected, b.metricsLabelsForConnection(tc.connectionName))
		})
	}
}
