// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

func testHealthCommand(tb testing.TB, client *api.Client) (*cli.MockUi, *HealthCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &HealthCommand{
		BaseCommand: &BaseCommand{
			UI:     ui,
			client: client,
		},
		now: func() time.Time {
			return time.Date(2026, time.September, 17, 12, 0, 0, 0, time.UTC)
		},
	}
}

func TestHealthCommand_Run(t *testing.T) {
	t.Parallel()

	var primaryAddress, secondaryAddress string
	heartbeat := "2026-09-17T11:59:55Z"

	primary := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/sys/health":
			writeHealthTestJSON(t, w, map[string]interface{}{
				"initialized":                  true,
				"sealed":                       false,
				"standby":                      false,
				"version":                      "1.21.0",
				"replication_dr_mode":          "primary",
				"replication_performance_mode": "disabled",
			})
		case "/v1/sys/replication/status":
			writeHealthTestJSON(t, w, map[string]interface{}{
				"data": map[string]interface{}{
					"dr": map[string]interface{}{
						"mode":        "primary",
						"state":       "running",
						"cluster_id":  "cluster-1",
						"last_dr_wal": 1200,
						"secondaries": []map[string]interface{}{
							{
								"api_address":       secondaryAddress,
								"connection_status": "connected",
								"last_heartbeat":    heartbeat,
							},
						},
					},
					"performance": map[string]interface{}{"mode": "disabled"},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer primary.Close()

	secondary := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/sys/health":
			writeHealthTestJSON(t, w, map[string]interface{}{
				"initialized":                  true,
				"sealed":                       false,
				"standby":                      false,
				"version":                      "1.21.0",
				"replication_dr_mode":          "secondary",
				"replication_performance_mode": "disabled",
			})
		case "/v1/sys/replication/status":
			writeHealthTestJSON(t, w, map[string]interface{}{
				"data": map[string]interface{}{
					"dr": map[string]interface{}{
						"mode":             "secondary",
						"state":            "stream-wals",
						"connection_state": "ready",
						"cluster_id":       "cluster-1",
						"last_remote_wal":  1150,
						"primaries": []map[string]interface{}{
							{
								"api_address":       primaryAddress,
								"connection_status": "connected",
								"last_heartbeat":    heartbeat,
							},
						},
					},
					"performance": map[string]interface{}{"mode": "disabled"},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer secondary.Close()

	primaryAddress = primary.URL
	secondaryAddress = secondary.URL

	client := testHealthAPIClient(t, primaryAddress)
	ui, cmd := testHealthCommand(t, client)

	code := cmd.Run(nil)
	require.Equal(t, 0, code)

	output := ui.OutputWriter.String()
	require.Contains(t, output, primaryAddress)
	require.Contains(t, output, secondaryAddress)
	require.Contains(t, output, "stream-wals")
	require.Contains(t, output, "+50")
	require.Contains(t, output, "Overall: HEALTHY")
}

func TestHealthCommand_CriticalReplication(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/sys/health":
			writeHealthTestJSON(t, w, map[string]interface{}{
				"initialized": true,
				"sealed":      false,
				"version":     "1.21.0",
			})
		case "/v1/sys/replication/status":
			writeHealthTestJSON(t, w, map[string]interface{}{
				"data": map[string]interface{}{
					"dr": map[string]interface{}{
						"mode":             "secondary",
						"state":            "idle",
						"connection_state": "transient_failure",
						"cluster_id":       "cluster-1",
						"last_remote_wal":  0,
					},
					"performance": map[string]interface{}{"mode": "disabled"},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	ui, cmd := testHealthCommand(t, testHealthAPIClient(t, server.URL))
	code := cmd.Run(nil)

	require.Equal(t, 2, code)
	output := ui.OutputWriter.String()
	require.Contains(t, output, "state=idle")
	require.Contains(t, output, "connection=transient_failure")
	require.Contains(t, output, "never streamed (remote_wal=0)")
	require.Contains(t, output, "Overall: CRITICAL")
}

func TestHealthCommand_NoReplication(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/sys/health":
			writeHealthTestJSON(t, w, map[string]interface{}{
				"initialized": true,
				"sealed":      false,
				"version":     "1.21.0",
			})
		case "/v1/sys/replication/status":
			writeHealthTestJSON(t, w, map[string]interface{}{
				"data": map[string]interface{}{
					"dr":          map[string]interface{}{"mode": "disabled"},
					"performance": map[string]interface{}{"mode": "disabled"},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	ui, cmd := testHealthCommand(t, testHealthAPIClient(t, server.URL))
	code := cmd.Run(nil)

	require.Equal(t, 0, code)
	require.Contains(t, ui.OutputWriter.String(), "Replication is not configured.")
	require.Contains(t, ui.OutputWriter.String(), "false")
	require.Contains(t, ui.OutputWriter.String(), "Overall: NOT_APPLICABLE")
}

func TestHealthCommand_TLS(t *testing.T) {
	t.Parallel()

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/sys/health":
			writeHealthTestJSON(t, w, map[string]interface{}{
				"initialized": true,
				"sealed":      false,
				"version":     "1.21.0",
			})
		case "/v1/sys/replication/status":
			writeHealthTestJSON(t, w, map[string]interface{}{
				"data": map[string]interface{}{
					"dr":          map[string]interface{}{"mode": "disabled"},
					"performance": map[string]interface{}{"mode": "disabled"},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	config := api.DefaultConfig()
	config.Address = server.URL
	config.HttpClient = server.Client()
	client, err := api.NewClient(config)
	require.NoError(t, err)

	ui, cmd := testHealthCommand(t, client)
	code := cmd.Run(nil)

	require.Equal(t, 0, code)
	output := ui.OutputWriter.String()
	require.Contains(t, output, "HTTPS")
	require.Contains(t, output, "true")
	require.NotContains(t, output, "true       n/a")
}

func TestHealthCommand_TLSExpiryAffectsOverallHealth(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.September, 17, 12, 0, 0, 0, time.UTC)
	expiresAt := now.Add(6 * 24 * time.Hour)
	report := healthReport{
		Overall: healthLevelNotApplicable,
		Nodes: []healthNode{
			{
				Address: "https://vault.example",
				TLS: healthTLS{
					Enabled:   true,
					ExpiresAt: &expiresAt,
				},
			},
		},
	}
	cmd := HealthCommand{now: func() time.Time { return now }}

	cmd.evaluateReplicationHealth(&report, nil)

	require.Equal(t, healthLevelCritical, report.Overall)
	require.NotNil(t, report.Nodes[0].TLS.ExpiresInDays)
	require.Equal(t, 6, *report.Nodes[0].TLS.ExpiresInDays)
}

func TestHealthCommand_Validation(t *testing.T) {
	t.Parallel()

	ui, cmd := testHealthCommand(t, nil)
	code := cmd.Run([]string{"extra"})

	require.Equal(t, 1, code)
	require.Contains(t, ui.ErrorWriter.String(), "Too many arguments")
}

func TestHealthCommand_NoTabs(t *testing.T) {
	t.Parallel()

	_, cmd := testHealthCommand(t, nil)
	assertNoTabs(t, cmd)
}

func testHealthAPIClient(tb testing.TB, address string) *api.Client {
	tb.Helper()

	config := api.DefaultConfig()
	config.Address = address
	client, err := api.NewClient(config)
	require.NoError(tb, err)
	return client
}

func writeHealthTestJSON(tb testing.TB, w http.ResponseWriter, value interface{}) {
	tb.Helper()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		tb.Errorf("failed to write JSON response: %v", err)
	}
}

func TestHealthNodeName(t *testing.T) {
	t.Parallel()

	require.Equal(t, "eastus2-prd-3", healthNodeName("https://eastus2-prd-3.vault.example:8200"))
	require.Equal(t, "not a URL", healthNodeName("not a URL"))
	require.Empty(t, strings.TrimSpace(healthNodeName("")))
}
