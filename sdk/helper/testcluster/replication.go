package testcluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
)

// PassiveWaitForActiveNodeAndPerfStandbys should be used instead of
// WaitForActiveNodeAndPerfStandbys when you don't want to do any writes
// as a side-effect. This returns perfStandby nodes in the cluster and
// an error.
func PassiveWaitForActiveNodeAndPerfStandbys(ctx context.Context, pri VaultCluster) (VaultClusterNode, []VaultClusterNode, error) {
	leaderNode, standbys, err := GetActiveAndStandbys(ctx, pri)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive standby nodes, %w", err)
	}

	for i, node := range standbys {
		client := node.APIClient()
		// Make sure we get perf standby nodes
		if err = EnsureCoreIsPerfStandby(ctx, client); err != nil {
			return nil, nil, fmt.Errorf("standby node %d is not a perfStandby, %w", i, err)
		}
	}

	return leaderNode, standbys, nil
}

func EnsureCoreIsPerfStandby(ctx context.Context, client *api.Client) error {
	var err error
	var health *api.HealthResponse
	for ctx.Err() == nil {
		health, err = client.Sys().HealthWithContext(ctx)
		if err == nil && health.PerformanceStandby {
			return nil
		}
		time.Sleep(time.Millisecond * 500)
	}
	if err == nil {
		err = ctx.Err()
	}
	return err
}

func GetActiveAndStandbys(ctx context.Context, cluster VaultCluster) (VaultClusterNode, []VaultClusterNode, error) {
	var leaderIndex int
	var err error
	if leaderIndex, err = WaitForActiveNode(ctx, cluster); err != nil {
		return nil, nil, err
	}

	var leaderNode VaultClusterNode
	var nodes []VaultClusterNode
	for i, node := range cluster.Nodes() {
		if i == leaderIndex {
			leaderNode = node
			continue
		}
		nodes = append(nodes, node)
	}

	return leaderNode, nodes, nil
}
