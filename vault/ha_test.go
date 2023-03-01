package vault

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/stretchr/testify/require"
)

// TestGrabLockOrStopped is a non-deterministic test to detect deadlocks in the
// grabLockOrStopped function. This test starts a bunch of workers which
// continually lock/unlock and rlock/runlock the same RWMutex. Each worker also
// starts a goroutine which closes the stop channel 1/2 the time, which races
// with acquisition of the lock.
func TestGrabLockOrStop(t *testing.T) {
	// Stop the test early if we deadlock.
	const (
		workers      = 100
		testDuration = time.Second
		testTimeout  = 10 * testDuration
	)
	done := make(chan struct{})
	defer close(done)
	var lockCount int64
	go func() {
		select {
		case <-done:
		case <-time.After(testTimeout):
			panic(fmt.Sprintf("deadlock after %d lock count",
				atomic.LoadInt64(&lockCount)))
		}
	}()

	// lock is locked/unlocked and rlocked/runlocked concurrently.
	var lock sync.RWMutex
	start := time.Now()

	// workerWg is used to wait until all workers exit.
	var workerWg sync.WaitGroup
	workerWg.Add(workers)

	// Start a bunch of worker goroutines.
	for g := 0; g < workers; g++ {
		g := g
		go func() {
			defer workerWg.Done()
			for time.Now().Sub(start) < testDuration {
				stop := make(chan struct{})

				// closerWg waits until the closer goroutine exits before we do
				// another iteration. This makes sure goroutines don't pile up.
				var closerWg sync.WaitGroup
				closerWg.Add(1)
				go func() {
					defer closerWg.Done()
					// Close the stop channel half the time.
					if rand.Int()%2 == 0 {
						close(stop)
					}
				}()

				// Half the goroutines lock/unlock and the other half rlock/runlock.
				if g%2 == 0 {
					if !grabLockOrStop(lock.Lock, lock.Unlock, stop) {
						lock.Unlock()
					}
				} else {
					if !grabLockOrStop(lock.RLock, lock.RUnlock, stop) {
						lock.RUnlock()
					}
				}

				closerWg.Wait()

				// This lets us know how many lock/unlock and rlock/runlock have
				// happened if there's a deadlock.
				atomic.AddInt64(&lockCount, 1)
			}
		}()
	}
	workerWg.Wait()
}

type testBackend struct {
	physical.Backend
	shouldFail *uint32
}

func (b *testBackend) SetShouldFail(shouldFail bool) {
	if shouldFail {
		atomic.StoreUint32(b.shouldFail, 1)
	} else {
		atomic.StoreUint32(b.shouldFail, 0)
	}
}

func (b *testBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	if key == rootKeyPath && atomic.LoadUint32(b.shouldFail) == 1 {
		return nil, fmt.Errorf("failed to get root key")
	}
	return b.Backend.Get(ctx, key)
}

// TestFailDuringKeysUpgrade checks that Vault does not seal itself when an
// error happens with the storage backend during the post-unseal opeartions but
// does report it properly in PostUnsealFailed().
func TestFailDuringKeysUpgrade(t *testing.T) {
	ctx := namespace.RootContext(context.Background())

	// Create the first core and initialize it
	logger = logging.NewVaultLogger(log.Trace).Named(t.Name())

	inm, err := inmem.NewInmemHA(nil, logger)
	require.NoError(t, err)

	inmha, err := inmem.NewInmemHA(nil, logger)
	require.NoError(t, err)

	redirectOriginal := "http://127.0.0.1:8200"
	core, err := NewCore(&CoreConfig{
		Physical:     inm,
		HAPhysical:   inmha.(physical.HABackend),
		RedirectAddr: redirectOriginal,
		DisableMlock: true,
		Logger:       logger.Named("core"),
	})
	require.NoError(t, err)
	defer core.Shutdown()
	keys, root := TestCoreInit(t, core)
	for _, key := range keys {
		_, err := TestCoreUnseal(core, TestKeyCopy(key))
		require.NoError(t, err)
	}

	// Verify unsealed
	require.False(t, core.Sealed())

	// Wait for core to become active
	TestWaitActive(t, core)

	// Check the leader is local
	isLeader, advertise, _, err := core.Leader()
	require.NoError(t, err)
	require.True(t, isLeader)
	require.Equal(t, redirectOriginal, advertise)

	// Create the second core with a faulty backend and initialize it
	redirectOriginal2 := "http://127.0.0.1:8500"
	backend := &testBackend{inm, new(uint32)}
	backend.SetShouldFail(true)
	core2, err := NewCore(&CoreConfig{
		Physical:     backend,
		HAPhysical:   inmha.(physical.HABackend),
		RedirectAddr: redirectOriginal2,
		DisableMlock: true,
		Logger:       logger.Named("core2"),
	})
	defer core2.Shutdown()
	require.NoError(t, err)
	for _, key := range keys {
		_, err := TestCoreUnseal(core2, TestKeyCopy(key))
		require.NoError(t, err)
	}

	// Verify unsealed
	require.False(t, core2.Sealed())

	// Core2 should be in standby
	standby, err := core2.Standby()
	require.NoError(t, err)
	require.True(t, standby)

	// Check the leader is not local
	isLeader, advertise, _, err = core2.Leader()
	require.NoError(t, err)
	require.False(t, isLeader)
	require.Equal(t, redirectOriginal, advertise)

	require.False(t, core2.Sealed())

	// Now we will step-down the leader, the second core should take the lead
	req := &logical.Request{
		ClientToken: root,
		Path:        "sys/step-down",
	}

	// Create an identifier for the request
	req.ID, err = uuid.GenerateUUID()
	require.NoError(t, err)

	// Step down core
	err = core.StepDown(ctx, req)
	require.NoError(t, err)

	// Give time to fail to switch leaders twice
	time.Sleep(15 * time.Second)

	// Check that core is still leader and core2 has not been sealed
	isLeader, _, _, err = core.Leader()
	require.NoError(t, err)
	require.True(t, isLeader)
	require.False(t, core2.Sealed())

	// core2 failed to unseal so it will return a different status code when
	// querying sys/health
	require.True(t, core2.PostUnsealFailed())

	backend.SetShouldFail(false)
	// Now we try again but this time core2 should be able to take the lead
	// properly
	req.ID, err = uuid.GenerateUUID()
	require.NoError(t, err)

	// Step down core
	err = core.StepDown(ctx, req)
	require.NoError(t, err)

	// Give time to fail to switch leaders twice
	time.Sleep(15 * time.Second)

	// Check that core2 is now leader and core has not been sealed
	isLeader, _, _, err = core2.Leader()
	require.NoError(t, err)
	require.True(t, isLeader)
	require.False(t, core2.Sealed())

	// Now the post-unseal operations succeeded
	require.False(t, core2.PostUnsealFailed())
}
