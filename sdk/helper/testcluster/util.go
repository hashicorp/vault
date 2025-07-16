// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testcluster

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/xor"
)

// Note that OSS standbys will not accept seal requests.  And ent perf standbys
// may fail it as well if they haven't yet been able to get "elected" as perf standbys.
func SealNode(ctx context.Context, cluster VaultCluster, nodeIdx int) error {
	if nodeIdx >= len(cluster.Nodes()) {
		return fmt.Errorf("invalid nodeIdx %d for cluster", nodeIdx)
	}
	node := cluster.Nodes()[nodeIdx]
	client := node.APIClient()

	err := client.Sys().SealWithContext(ctx)
	if err != nil {
		return err
	}

	return NodeSealed(ctx, cluster, nodeIdx)
}

func SealAllNodes(ctx context.Context, cluster VaultCluster) error {
	for i := range cluster.Nodes() {
		if err := SealNode(ctx, cluster, i); err != nil {
			return err
		}
	}
	return nil
}

func UnsealNode(ctx context.Context, cluster VaultCluster, nodeIdx int) error {
	if nodeIdx >= len(cluster.Nodes()) {
		return fmt.Errorf("invalid nodeIdx %d for cluster", nodeIdx)
	}
	node := cluster.Nodes()[nodeIdx]
	client := node.APIClient()

	for _, key := range cluster.GetBarrierOrRecoveryKeys() {
		_, err := client.Sys().UnsealWithContext(ctx, hex.EncodeToString(key))
		if err != nil {
			return err
		}
	}

	return NodeHealthy(ctx, cluster, nodeIdx)
}

func UnsealAllNodes(ctx context.Context, cluster VaultCluster) error {
	for i := range cluster.Nodes() {
		if err := UnsealNode(ctx, cluster, i); err != nil {
			return err
		}
	}
	return nil
}

func NodeSealed(ctx context.Context, cluster VaultCluster, nodeIdx int) error {
	if nodeIdx >= len(cluster.Nodes()) {
		return fmt.Errorf("invalid nodeIdx %d for cluster", nodeIdx)
	}
	node := cluster.Nodes()[nodeIdx]
	client := node.APIClient()

	var health *api.HealthResponse
	var err error
	for ctx.Err() == nil {
		health, err = client.Sys().HealthWithContext(ctx)
		switch {
		case err != nil:
		case !health.Sealed:
			err = fmt.Errorf("unsealed: %#v", health)
		default:
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("node %d is not sealed: %v", nodeIdx, err)
}

func WaitForNCoresSealed(ctx context.Context, cluster VaultCluster, n int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errs := make(chan error)
	for i := range cluster.Nodes() {
		go func(i int) {
			var err error
			for ctx.Err() == nil {
				err = NodeSealed(ctx, cluster, i)
				if err == nil {
					errs <- nil
					return
				}
				time.Sleep(100 * time.Millisecond)
			}
			if err == nil {
				err = ctx.Err()
			}
			errs <- err
		}(i)
	}

	var merr *multierror.Error
	var sealed int
	for range cluster.Nodes() {
		err := <-errs
		if err != nil {
			merr = multierror.Append(merr, err)
		} else {
			sealed++
			if sealed == n {
				return nil
			}
		}
	}

	return fmt.Errorf("%d cores were not sealed, errs: %v", n, merr.ErrorOrNil())
}

func NodeHealthy(ctx context.Context, cluster VaultCluster, nodeIdx int) error {
	if nodeIdx >= len(cluster.Nodes()) {
		return fmt.Errorf("invalid nodeIdx %d for cluster", nodeIdx)
	}
	node := cluster.Nodes()[nodeIdx]
	client := node.APIClient()

	var health *api.HealthResponse
	var err error
	for ctx.Err() == nil {
		health, err = client.Sys().HealthWithContext(ctx)
		switch {
		case err != nil:
		case health == nil:
			err = fmt.Errorf("nil response to health check")
		case health.Sealed:
			err = fmt.Errorf("sealed: %#v", health)
		default:
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("node %d is unhealthy: %v", nodeIdx, err)
}

func LeaderNode(ctx context.Context, cluster VaultCluster) (int, error) {
	// Be robust to multiple nodes thinking they are active. This is possible in
	// certain network partition situations where the old leader has not
	// discovered it's lost leadership yet. In tests this is only likely to come
	// up when we are specifically provoking it, but it's possible it could happen
	// at any point if leadership flaps of connectivity suffers transient errors
	// etc. so be robust against it. The best solution would be to have some sort
	// of epoch like the raft term that is guaranteed to be monotonically
	// increasing through elections, however we don't have that abstraction for
	// all HABackends in general. The best we have is the ActiveTime. In a
	// distributed systems text book this would be bad to rely on due to clock
	// sync issues etc. but for our tests it's likely fine because even if we are
	// running separate Vault containers, they are all using the same hardware
	// clock in the system.
	leaderActiveTimes := make(map[int]time.Time)
	for i, node := range cluster.Nodes() {
		client := node.APIClient()
		ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
		resp, err := client.Sys().LeaderWithContext(ctx)
		cancel()
		if err != nil || resp == nil || !resp.IsSelf {
			continue
		}
		leaderActiveTimes[i] = resp.ActiveTime
	}
	if len(leaderActiveTimes) == 0 {
		return -1, fmt.Errorf("no leader found")
	}
	// At least one node thinks it is active. If multiple, pick the one with the
	// most recent ActiveTime. Note if there is only one then this just returns
	// it.
	var newestLeaderIdx int
	var newestActiveTime time.Time
	for i, at := range leaderActiveTimes {
		if at.After(newestActiveTime) {
			newestActiveTime = at
			newestLeaderIdx = i
		}
	}
	return newestLeaderIdx, nil
}

func WaitForActiveNode(ctx context.Context, cluster VaultCluster) (int, error) {
	for ctx.Err() == nil {
		if idx, _ := LeaderNode(ctx, cluster); idx != -1 {
			return idx, nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return -1, ctx.Err()
}

func WaitForStandbyNode(ctx context.Context, cluster VaultCluster, nodeIdx int) error {
	if nodeIdx >= len(cluster.Nodes()) {
		return fmt.Errorf("invalid nodeIdx %d for cluster", nodeIdx)
	}
	node := cluster.Nodes()[nodeIdx]
	client := node.APIClient()

	var err error
	for ctx.Err() == nil {
		var resp *api.LeaderResponse

		resp, err = client.Sys().LeaderWithContext(ctx)
		switch {
		case err != nil:
		case resp.IsSelf:
			return fmt.Errorf("waiting for standby but node is leader")
		case resp.LeaderAddress == "":
			err = fmt.Errorf("node doesn't know leader address")
		default:
			return nil
		}

		time.Sleep(100 * time.Millisecond)
	}
	if err == nil {
		err = ctx.Err()
	}
	return err
}

func WaitForActiveNodeAndStandbys(ctx context.Context, cluster VaultCluster) (int, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	leaderIdx, err := WaitForActiveNode(ctx, cluster)
	if err != nil {
		return 0, err
	}

	if len(cluster.Nodes()) == 1 {
		return 0, nil
	}

	errs := make(chan error)
	for i := range cluster.Nodes() {
		if i == leaderIdx {
			continue
		}
		go func(i int) {
			errs <- WaitForStandbyNode(ctx, cluster, i)
		}(i)
	}

	var merr *multierror.Error
	expectedStandbys := len(cluster.Nodes()) - 1
	for i := 0; i < expectedStandbys; i++ {
		merr = multierror.Append(merr, <-errs)
	}

	return leaderIdx, merr.ErrorOrNil()
}

func WaitForActiveNodeAndPerfStandbys(ctx context.Context, cluster VaultCluster) error {
	logger := cluster.NamedLogger("WaitForActiveNodeAndPerfStandbys")
	// This WaitForActiveNode was added because after a Raft cluster is sealed
	// and then unsealed, when it comes up it may have a different leader than
	// Core0, making this helper fail.
	// A sleep before calling WaitForActiveNodeAndPerfStandbys seems to sort
	// things out, but so apparently does this.  We should be able to eliminate
	// this call to WaitForActiveNode by reworking the logic in this method.
	leaderIdx, err := WaitForActiveNode(ctx, cluster)
	if err != nil {
		return fmt.Errorf("did not find leader: %w", err)
	}

	if len(cluster.Nodes()) == 1 {
		return nil
	}

	expectedStandbys := len(cluster.Nodes()) - 1

	mountPoint, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	leaderClient := cluster.Nodes()[leaderIdx].APIClient()

	for ctx.Err() == nil {
		err = leaderClient.Sys().MountWithContext(ctx, mountPoint, &api.MountInput{
			Type:  "kv",
			Local: true,
		})
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		return fmt.Errorf("unable to mount KV engine: %w", err)
	}
	path := mountPoint + "/waitforactivenodeandperfstandbys"
	var standbys, actives int64
	errchan := make(chan error, len(cluster.Nodes()))
	for i := range cluster.Nodes() {
		go func(coreNo int) {
			node := cluster.Nodes()[coreNo]
			client := node.APIClient()
			val := 1
			var err error
			defer func() {
				errchan <- err
			}()

			var lastWAL uint64
			for ctx.Err() == nil {
				_, err = leaderClient.Logical().WriteWithContext(ctx, path, map[string]interface{}{
					"bar": val,
				})
				val++
				time.Sleep(250 * time.Millisecond)
				if err != nil {
					continue
				}
				var leader *api.LeaderResponse
				leader, err = client.Sys().LeaderWithContext(ctx)
				if err != nil {
					logger.Trace("waiting for core", "core", coreNo, "err", err)
					continue
				}
				switch {
				case leader.IsSelf:
					logger.Trace("waiting for core", "core", coreNo, "isLeader", true)
					atomic.AddInt64(&actives, 1)
					return
				case leader.PerfStandby && leader.PerfStandbyLastRemoteWAL > 0:
					switch {
					case lastWAL == 0:
						lastWAL = leader.PerfStandbyLastRemoteWAL
						logger.Trace("waiting for core", "core", coreNo, "lastRemoteWAL", leader.PerfStandbyLastRemoteWAL, "lastWAL", lastWAL)
					case lastWAL < leader.PerfStandbyLastRemoteWAL:
						logger.Trace("waiting for core", "core", coreNo, "lastRemoteWAL", leader.PerfStandbyLastRemoteWAL, "lastWAL", lastWAL)
						atomic.AddInt64(&standbys, 1)
						return
					}
				default:
					logger.Trace("waiting for core", "core", coreNo,
						"ha_enabled", leader.HAEnabled,
						"is_self", leader.IsSelf,
						"perf_standby", leader.PerfStandby,
						"perf_standby_remote_wal", leader.PerfStandbyLastRemoteWAL)
				}
			}
		}(i)
	}

	errs := make([]error, 0, len(cluster.Nodes()))
	for range cluster.Nodes() {
		errs = append(errs, <-errchan)
	}
	if actives != 1 || int(standbys) != expectedStandbys {
		return fmt.Errorf("expected 1 active core and %d standbys, got %d active and %d standbys, errs: %v",
			expectedStandbys, actives, standbys, errs)
	}

	for ctx.Err() == nil {
		err = leaderClient.Sys().UnmountWithContext(ctx, mountPoint)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		return fmt.Errorf("unable to unmount KV engine: %w", err)
	}
	return nil
}

func Clients(vc VaultCluster) []*api.Client {
	var ret []*api.Client
	for _, n := range vc.Nodes() {
		ret = append(ret, n.APIClient())
	}
	return ret
}

//go:generate enumer -type=GenerateRootKind -trimprefix=GenerateRoot
type GenerateRootKind int

const (
	GenerateRootRegular GenerateRootKind = iota
	GenerateRootDR
	GenerateRecovery
)

func GenerateRoot(cluster VaultCluster, kind GenerateRootKind) (string, error) {
	// If recovery keys supported, use those to perform root token generation instead
	keys := cluster.GetBarrierOrRecoveryKeys()

	client := cluster.Nodes()[0].APIClient()

	var err error
	var status *api.GenerateRootStatusResponse
	switch kind {
	case GenerateRootRegular:
		status, err = client.Sys().GenerateRootInit("", "")
	case GenerateRootDR:
		status, err = client.Sys().GenerateDROperationTokenInit("", "")
	case GenerateRecovery:
		status, err = client.Sys().GenerateRecoveryOperationTokenInit("", "")
	}
	if err != nil {
		return "", err
	}

	if status.Required > len(keys) {
		return "", fmt.Errorf("need more keys than have, need %d have %d", status.Required, len(keys))
	}

	otp := status.OTP

	for i, key := range keys {
		if i >= status.Required {
			break
		}

		strKey := base64.StdEncoding.EncodeToString(key)
		switch kind {
		case GenerateRootRegular:
			status, err = client.Sys().GenerateRootUpdate(strKey, status.Nonce)
		case GenerateRootDR:
			status, err = client.Sys().GenerateDROperationTokenUpdate(strKey, status.Nonce)
		case GenerateRecovery:
			status, err = client.Sys().GenerateRecoveryOperationTokenUpdate(strKey, status.Nonce)
		}
		if err != nil {
			return "", err
		}
	}
	if !status.Complete {
		return "", fmt.Errorf("generate root operation did not end successfully")
	}

	tokenBytes, err := base64.RawStdEncoding.DecodeString(status.EncodedToken)
	if err != nil {
		return "", err
	}
	tokenBytes, err = xor.XORBytes(tokenBytes, []byte(otp))
	if err != nil {
		return "", err
	}
	return string(tokenBytes), nil
}
