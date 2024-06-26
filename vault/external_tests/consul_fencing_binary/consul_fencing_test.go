// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package consul_fencing

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/consul"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/sdk/helper/testcluster"
	"github.com/hashicorp/vault/sdk/helper/testcluster/docker"
	"github.com/stretchr/testify/require"
)

// TestConsulFencing_PartitionedLeaderCantWrite attempts to create an active
// node split-brain when using Consul storage to ensure the old leader doesn't
// continue to write data potentially corrupting storage. It is naturally
// non-deterministic since it relies heavily on timing of the different
// container processes, however it pretty reliably failed before the fencing fix
// (and Consul lock improvements) and should _never_ fail now we correctly fence
// writes.
func TestConsulFencing_PartitionedLeaderCantWrite(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	consulStorage := consul.NewClusterStorage()

	// Create  cluster logger that will write cluster logs to a file in CI.
	logger := corehelpers.NewTestLogger(t)
	logger.SetLevel(hclog.Trace)

	clusterOpts := docker.DefaultOptions(t)
	// We can use an enterprise image here because we are swapping out the binary anyway.
	clusterOpts.ImageRepo = "hashicorp/vault-enterprise"
	clusterOpts.ClusterOptions.Logger = logger

	clusterOpts.Storage = consulStorage

	logger.Info("==> starting cluster")
	c, err := docker.NewDockerCluster(ctx, clusterOpts)
	require.NoError(t, err)
	logger.Info("  ✅ done.", "root_token", c.GetRootToken(),
		"consul_token", consulStorage.Config().Token)

	logger.Info("==> waiting for leader")
	leaderIdx, err := testcluster.WaitForActiveNode(ctx, c)
	require.NoError(t, err)

	leader := c.Nodes()[leaderIdx]
	leaderClient := leader.APIClient()

	notLeader := c.Nodes()[1] // Assumes it's usually zero and correct below
	if leaderIdx == 1 {
		notLeader = c.Nodes()[0]
	}

	// Mount a KV v2 backend
	logger.Info("==> mounting KV")
	err = leaderClient.Sys().Mount("/test", &api.MountInput{
		Type: "kv-v2",
	})
	require.NoError(t, err)

	// We need to wait for the KV mount to be ready on all servers - it runs an
	// "Upgrade" process on mount and will error until that's done. In practice
	// that's so fast that in CE I've never seen it fail yet, but in Enterprise
	// there is the added complexity of the upgrade running on the primary while
	// standby nodes wait for it to complete. So in Enterprise the first PATCH
	// sent to a standby often comes in before the standby has "seen" that the
	// primary has completed the upgrade and fails causing the whole test to
	// error. We can prevent that by taking a short loop here to wait for a read
	// to succeed on the standby/non-leader (since KV-v2 does upgrade check on
	// read too).
	waitForKVv2Upgrade(t, ctx, notLeader.APIClient(), "/test")

	// Start two background workers that will cause writes to Consul in the
	// background. KV v2 relies on a single active node for correctness.
	// Specifically its patch operation does a read-modify-write under a
	// key-specific lock which is correct for concurrent writes to one process,
	// but which by nature of our storage API is not going to be atomic if another
	// active node is also writing the same KV. It's made worse because the cache
	// backend means the active node will not actually read from Consul after the
	// initial read and will be modifying its own in-memory version and writing
	// that back. So we should be able to detect multiple active nodes writing
	// reliably with the following setup:
	//  1. Two separate "client" goroutines each connected to different Vault
	//     servers.
	//  2. Both write to the same kv-v2 key, each one modifies only its own set
	//     of subkeys c1-X or c2-X.
	//  3. Each request adds the next sequential X value for that client. We use a
	//     Patch operation so we don't need to read the version or use CAS. On an
	//     error each client will retry the same key until it gets a success.
	//  4. In a correct system with a single active node despite a partition, we
	//     expect a complete set of consecutive X values for both clients (i.e.
	//     no writes lost). If an old leader is still allowed to write to Consul
	//     then it will continue to patch against its own last-known version from
	//     cache and so will overwrite any concurrent updates from the other
	//     client and we should see that as lost updates at the end.
	var wg sync.WaitGroup
	errCh := make(chan error, 10)
	var writeCount uint64

	// Initialise the key once
	kv := leaderClient.KVv2("/test")
	_, err = kv.Put(ctx, "data", map[string]interface{}{
		"c0-00000000": 1, // value don't matter here only keys in this set.
		"c1-00000000": 1,
	})
	require.NoError(t, err)

	const interval = 500 * time.Millisecond
	const timeout = 3 * time.Second
	runWriter := func(i int, targetServer testcluster.VaultClusterNode, ctr *uint64) {
		wg.Add(1)
		defer wg.Done()
		kv := targetServer.APIClient().KVv2("/test")

		for {
			key := fmt.Sprintf("c%d-%08d", i, atomic.LoadUint64(ctr))

			// Use a short timeout. If we don't then the one goroutine writing
			// to the partitioned active node can get stuck here until the 60
			// second request timeout kicks in without issuing another request.
			// However, this timeout being too short can cause issues too.
			// Having it set to 500 milliseconds caused the test to
			// intermittently fail in CI before.
			reqCtx, cancel := context.WithTimeout(ctx, timeout)
			logger.Debug("sending patch", "client", i, "key", key)
			_, err = kv.Patch(reqCtx, "data", map[string]interface{}{
				key: 1,
			})
			cancel()
			// Deliver errors to test, don't block if we get too many before context
			// is cancelled otherwise client 0 can end up blocked as it has so many
			// errors during the partition it doesn't actually start writing again
			// ever and so the test never sees split-brain writes.
			if err != nil {
				select {
				case <-ctx.Done():
					return
				case errCh <- fmt.Errorf("client %d error: %w", i, err):
				default:
					// errCh is blocked, carry on anyway
				}
			} else {
				// Only increment our set counter here now we've had an ack that the
				// update was successful.
				atomic.AddUint64(ctr, 1)
				atomic.AddUint64(&writeCount, 1)
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(interval):
			}
		}
	}

	logger.Info("==> starting writers")
	client0Ctr, client1Ctr := uint64(1), uint64(1)
	go runWriter(0, leader, &client0Ctr)
	go runWriter(1, notLeader, &client1Ctr)

	// Wait for some writes to have started
	var writesBeforePartition uint64
	logger.Info("==> waiting for writes")
	for {
		time.Sleep(1 * time.Millisecond)
		writesBeforePartition = atomic.LoadUint64(&writeCount)
		if writesBeforePartition >= 5 {
			break
		}
		// Also check for any write errors
		select {
		case err := <-errCh:
			require.NoError(t, err)
		default:
		}
		require.NoError(t, ctx.Err())
	}

	val, err := kv.Get(ctx, "data")
	require.NoError(t, err)

	logger.Info("==> partitioning leader")
	// Now partition the leader from everything else (including Consul)
	err = leader.(*docker.DockerClusterNode).PartitionFromCluster(ctx)
	require.NoError(t, err)

	// Reload this incase more writes occurred before the partition completed.
	writesBeforePartition = atomic.LoadUint64(&writeCount)

	// Wait for some more writes to have happened (the client writing to leader
	// will probably have sent one and hung waiting for a response but the other
	// one should eventually start committing again when new active node is
	// elected).

	logger.Info("==> waiting for writes to new leader")
	for {
		time.Sleep(1 * time.Millisecond)
		writesAfterPartition := atomic.LoadUint64(&writeCount)
		if (writesAfterPartition - writesBeforePartition) >= 20 {
			break
		}
		// Also check for any write errors or timeouts
		select {
		case err := <-errCh:
			// Don't fail here because we expect writes to the old leader to fail
			// eventually or if they need a new connection etc.
			logger.Info("failed write", "write_count", writesAfterPartition, "err", err)
		default:
		}
		require.NoError(t, ctx.Err(), "context error while waiting for writes to new leader")
	}

	// Heal partition
	logger.Info("==> healing partition")
	err = leader.(*docker.DockerClusterNode).UnpartitionFromCluster(ctx)
	require.NoError(t, err)

	// Wait for old leader to rejoin as a standby and get healthy.
	logger.Info("==> wait for old leader to rejoin")

	require.NoError(t, waitUntilNotLeader(ctx, leaderClient, logger))

	// Stop the writers and wait for them to shut down nicely
	logger.Info("==> stopping writers")
	cancel()
	wg.Wait()

	// Now verify that all Consul data is consistent with only one leader writing.
	// Use a new context since we just cancelled the general one
	reqCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	val, err = kv.Get(reqCtx, "data")
	require.NoError(t, err)

	// Ensure we have every consecutive key for both client
	sets := [][]int{make([]int, 0, 128), make([]int, 0, 128)}
	for k := range val.Data {
		var cNum, x int
		_, err := fmt.Sscanf(k, "c%d-%08d", &cNum, &x)
		require.NoError(t, err)
		sets[cNum] = append(sets[cNum], x)
	}

	// Sort both sets
	sort.Ints(sets[0])
	sort.Ints(sets[1])

	// Ensure they are both complete by creating an expected set and comparing to
	// get nice output to debug. Note that make set is an exclusive bound since we
	// don't know it the current counter value write completed or not yet so we'll
	// only create sets up to one less than that value that we know for sure
	// should be present.
	c0Writes := int(atomic.LoadUint64(&client0Ctr))
	c1Writes := int(atomic.LoadUint64(&client1Ctr))
	expect0 := makeSet(c0Writes)
	expect1 := makeSet(c1Writes)

	// Trim the sets to only the writes we know completed since that's all the
	// expected arrays contain. But only if they are longer so we don't change the
	// slice length if they are shorter than the expected number.
	if len(sets[0]) > c0Writes {
		sets[0] = sets[0][0:c0Writes]
	}
	if len(sets[1]) > c1Writes {
		sets[1] = sets[1][0:c1Writes]
	}
	require.Equal(t, expect0, sets[0], "Client 0 writes lost")
	require.Equal(t, expect1, sets[1], "Client 1 writes lost")
}

func waitForKVv2Upgrade(t *testing.T, ctx context.Context, client *api.Client, path string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	kv := client.KVv2(path)
	attempts := 0
	wait := 20 * time.Millisecond
	for {
		// Attempt to perform a write on the KVv2 mount. It will fail until the
		// backend is done upgrading. If this is a performance standby in Ent then
		// it will not complete until the primary is done upgrading AND the standby
		// has noticed that!
		_, err := kv.Put(ctx, "test-upgrade-done", map[string]interface{}{
			"ok": 1,
		})
		if err == nil {
			return
		}
		t.Logf("waitForKVv2Upgrade: write failed: %s", err)
		select {
		case <-ctx.Done():
			t.Fatalf("context cancelled waiting for KVv2 (%s) upgrade to complete: %s",
				path, ctx.Err())
			return
		case <-time.After(wait):
		}
		attempts++
		// We don't quite want exponential backoff because it really should be fast,
		// but just reduce log spam on failures if it's taking a while.
		if attempts > 4 {
			wait = 250 * time.Millisecond
		}
		if attempts > 10 {
			wait = time.Second
		}
	}
}

func makeSet(n int) []int {
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = i
	}
	return a
}

func waitUntilNotLeader(ctx context.Context, oldLeaderClient *api.Client, logger hclog.Logger) error {
	for {
		// Wait for the original leader to acknowledge it's not active any more.
		resp, err := oldLeaderClient.Sys().LeaderWithContext(ctx)
		// Tolerate errors as the old leader is in a difficult place right now.
		if err == nil {
			if !resp.IsSelf {
				// We are not leader!
				return nil
			}
			logger.Info("old leader not ready yet", "IsSelf", resp.IsSelf)
		} else {
			logger.Info("failed to read old leader status", "err", err)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
			// Loop again
		}
	}
}
