package testcluster

import (
	"context"
	"fmt"
	"time"
)

func WaitForActiveNode(ctx context.Context, cluster VaultCluster) (int, error) {
	for ctx.Err() == nil {
		if idx, _ := LeaderNode(ctx, cluster); idx != -1 {
			return idx, nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return -1, ctx.Err()
}

func LeaderNode(ctx context.Context, cluster VaultCluster) (int, error) {
	for i, node := range cluster.Nodes() {
		client := node.APIClient()
		ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		resp, err := client.Sys().LeaderWithContext(ctx)
		cancel()
		if err != nil || resp == nil || !resp.IsSelf {
			continue
		}
		return i, nil
	}
	return -1, fmt.Errorf("no leader found")
}
