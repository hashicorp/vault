// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package system

import (
	"testing"

	"github.com/hashicorp/vault/command/server"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

// TestSystemBackend_InternalUISettings verifies that GET sys/internal/ui/settings
// reflects the ui_settings stanza and requires authentication.
func TestSystemBackend_InternalUISettings(t *testing.T) {
	t.Parallel()

	newCluster := func(t *testing.T, raw *server.Config) *vault.TestCluster {
		cluster := vault.NewTestCluster(t, &vault.CoreConfig{RawConfig: raw}, &vault.TestClusterOptions{
			HandlerFunc: vaulthttp.Handler,
			NumCores:    1,
		})
		t.Cleanup(cluster.Cleanup)
		return cluster
	}

	t.Run("defaults false when ui_settings absent", func(t *testing.T) {
		t.Parallel()
		cluster := newCluster(t, &server.Config{})
		client := cluster.Cores[0].Client

		resp, err := client.Logical().Read("sys/internal/ui/settings")
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, false, resp.Data["ui_telemetry_enabled"])
	})

	t.Run("returns configured true", func(t *testing.T) {
		t.Parallel()
		raw := &server.Config{UISettings: &server.UISettings{UITelemetry: true, UITelemetrySet: true}}
		cluster := newCluster(t, raw)
		client := cluster.Cores[0].Client

		resp, err := client.Logical().Read("sys/internal/ui/settings")
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, true, resp.Data["ui_telemetry_enabled"])
	})

	t.Run("returns configured false", func(t *testing.T) {
		t.Parallel()
		raw := &server.Config{UISettings: &server.UISettings{UITelemetry: false, UITelemetrySet: true}}
		cluster := newCluster(t, raw)
		client := cluster.Cores[0].Client

		resp, err := client.Logical().Read("sys/internal/ui/settings")
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, false, resp.Data["ui_telemetry_enabled"])
	})

	t.Run("requires authentication", func(t *testing.T) {
		t.Parallel()
		cluster := newCluster(t, &server.Config{})
		client := cluster.Cores[0].Client
		client.ClearToken()

		_, err := client.Logical().Read("sys/internal/ui/settings")
		require.Error(t, err)
	})
}
