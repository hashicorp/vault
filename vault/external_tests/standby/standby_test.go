// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package standby

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	"github.com/hashicorp/vault/sdk/helper/testcluster"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/cluster"
)

// Test_Echo_Duration_Skew tests that the sys/health and sys/ha-status endpoints
// report reasonable values for echo duration and clock skew.
func Test_Echo_Duration_Skew(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name        string
		perfstandby bool
	}{
		{"standby", false},
		{"perfstandby", true},
	}
	for i := range cases {
		perfstandby := cases[i].perfstandby
		if perfstandby && !constants.IsEnterprise {
			continue
		}
		t.Run(cases[i].name, func(t *testing.T) {
			t.Parallel()
			conf, opts := teststorage.ClusterSetup(nil, nil, nil)
			name := strings.Replace(t.Name(), "/", "_", -1)
			logger := corehelpers.NewTestLogger(t)
			layers, err := cluster.NewInmemLayerCluster(name, 3, logger)
			if err != nil {
				t.Fatal(err)
			}
			opts.ClusterLayers = layers
			opts.Logger = logger
			conf.DisablePerformanceStandby = !perfstandby
			cluster := vault.NewTestCluster(t, conf, opts)
			defer cluster.Cleanup()

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			leaderIdx, err := testcluster.WaitForActiveNodeAndStandbys(ctx, cluster)
			if err != nil {
				t.Fatal(err)
			}
			leader := cluster.Nodes()[leaderIdx]

			// The delay applies in both directions, hence a 0.25s delay implies a 0.5s roundtrip delay
			layers.SetReaderDelay(time.Second / 4)

			check := func(echoDuration int64, clockSkew int64) error {
				if echoDuration < time.Second.Milliseconds()/2 {
					return fmt.Errorf("echo duration must exceed 0.5s, got: %dms", echoDuration)
				}
				// Because we're using the same clock for all nodes, any clock skew will
				// be negative, as it's based on the delta of server time across both nodes,
				// but it doesn't factor in the round-trip time of the echo request.
				if clockSkew == 0 || -clockSkew < time.Second.Milliseconds()/2 {
					return fmt.Errorf("clock skew must be nonzero and exceed -0.5s, got: %dms", clockSkew)
				}

				return nil
			}

			// We need to wait for at least 2 heartbeats to happen (2s intervals)
			corehelpers.RetryUntil(t, 5*time.Second, func() error {
				haStatus, err := leader.APIClient().Sys().HAStatus()
				if err != nil {
					t.Fatal(err)
				}
				if len(haStatus.Nodes) < 3 {
					return fmt.Errorf("expected 3 nodes, got %d", len(haStatus.Nodes))
				}
				for _, node := range haStatus.Nodes {
					if node.ActiveNode {
						continue
					}

					if err := check(node.EchoDurationMillis, node.ClockSkewMillis); err != nil {
						return fmt.Errorf("ha-status node %s: %w", node.Hostname, err)
					}
				}

				for i, node := range cluster.Nodes() {
					if i == leaderIdx {
						continue
					}

					h, err := node.APIClient().Sys().Health()
					if err != nil {
						t.Fatal(err)
					}

					if err := check(h.EchoDurationMillis, h.ClockSkewMillis); err != nil {
						return fmt.Errorf("health node %s: %w", node.APIClient().Address(), err)
					}
				}
				return nil
			})
		})
	}
}
