// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// mockRollback returns a mock rollback manager
func mockRollback(t *testing.T) (*RollbackManager, *NoopBackend) {
	backend := new(NoopBackend)
	mounts := new(MountTable)
	router := NewRouter()
	core, _, _ := TestCoreUnsealed(t)

	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")

	mounts.Entries = []*MountEntry{
		{
			Path:        "foo",
			NamespaceID: namespace.RootNamespaceID,
			namespace:   namespace.RootNamespace,
		},
	}
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	if err := router.Mount(backend, "foo", &MountEntry{UUID: meUUID, Accessor: "noopaccessor", NamespaceID: namespace.RootNamespaceID, namespace: namespace.RootNamespace}, view); err != nil {
		t.Fatalf("err: %s", err)
	}

	mountsFunc := func() []*MountEntry {
		return mounts.Entries
	}

	logger := logging.NewVaultLogger(log.Trace)

	rb := NewRollbackManager(context.Background(), logger, mountsFunc, router, core)
	rb.period = 10 * time.Millisecond
	return rb, backend
}

func TestRollbackManager(t *testing.T) {
	m, backend := mockRollback(t)
	if len(backend.Paths) > 0 {
		t.Fatalf("bad: %#v", backend)
	}

	m.Start()
	time.Sleep(50 * time.Millisecond)
	m.Stop()

	count := len(backend.Paths)
	if count == 0 {
		t.Fatalf("bad: %#v", backend)
	}
	if backend.Paths[0] != "" {
		t.Fatalf("bad: %#v", backend)
	}

	time.Sleep(50 * time.Millisecond)

	if count != len(backend.Paths) {
		t.Fatalf("should stop requests: %#v", backend)
	}
}

// TestRollbackManager_ManyWorkers adds 10 backends that require a rollback
// operation, with 20 workers. The test verifies that the 10
// work items will run in parallel
func TestRollbackManager_ManyWorkers(t *testing.T) {
	core := TestCoreWithConfig(t, &CoreConfig{NumRollbackWorkers: 20, RollbackPeriod: time.Millisecond * 10})
	view := NewBarrierView(core.barrier, "logical/")

	ran := make(chan string)
	release := make(chan struct{})
	core, _, _ = testCoreUnsealed(t, core)

	// create 10 backends
	// when a rollback happens, each backend will try to write to an unbuffered
	// channel, then wait to be released
	for i := 0; i < 10; i++ {
		b := &NoopBackend{}
		b.RequestHandler = func(ctx context.Context, request *logical.Request) (*logical.Response, error) {
			if request.Operation == logical.RollbackOperation {
				ran <- request.Path
				<-release
			}
			return nil, nil
		}
		b.Root = []string{fmt.Sprintf("foo/%d", i)}
		meUUID, err := uuid.GenerateUUID()
		require.NoError(t, err)
		mountEntry := &MountEntry{
			Table:       mountTableType,
			UUID:        meUUID,
			Accessor:    fmt.Sprintf("accessor-%d", i),
			NamespaceID: namespace.RootNamespaceID,
			namespace:   namespace.RootNamespace,
			Path:        fmt.Sprintf("logical/foo/%d", i),
		}
		func() {
			core.mountsLock.Lock()
			defer core.mountsLock.Unlock()
			newTable := core.mounts.shallowClone()
			newTable.Entries = append(newTable.Entries, mountEntry)
			core.mounts = newTable
			err = core.router.Mount(b, "logical", mountEntry, view)
			require.NoError(t, core.persistMounts(context.Background(), newTable, &mountEntry.Local))
		}()
	}

	timeout, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	got := make(map[string]bool)
	hasMore := true
	for hasMore {
		// we're not bounding the number of workers, so we would expect to see
		// all 10 writes to the channel from each of the backends. Once that
		// happens, close the release channel so that the functions can exit
		select {
		case <-timeout.Done():
			require.Fail(t, "test timed out")
		case i := <-ran:
			got[i] = true
			if len(got) == 10 {
				close(release)
				hasMore = false
			}
		}
	}
	done := make(chan struct{})

	// start a goroutine to consume the remaining items from the queued work
	go func() {
		for {
			select {
			case <-ran:
			case <-done:
				return
			}
		}
	}()
	// stop the rollback worker, which will wait for all inflight rollbacks to
	// complete
	core.rollback.Stop()
	close(done)
}

// TestRollbackManager_WorkerPool adds 10 backends that require a rollback
// operation, with 5 workers. The test verifies that the 5 work items can occur
// concurrently, and that the remainder of the work is queued and run when
// workers are available
func TestRollbackManager_WorkerPool(t *testing.T) {
	core := TestCoreWithConfig(t, &CoreConfig{NumRollbackWorkers: 5, RollbackPeriod: time.Millisecond * 10})
	view := NewBarrierView(core.barrier, "logical/")

	ran := make(chan string)
	release := make(chan struct{})
	core, _, _ = testCoreUnsealed(t, core)

	// create 10 backends
	// when a rollback happens, each backend will try to write to an unbuffered
	// channel, then wait to be released
	for i := 0; i < 10; i++ {
		b := &NoopBackend{}
		b.RequestHandler = func(ctx context.Context, request *logical.Request) (*logical.Response, error) {
			if request.Operation == logical.RollbackOperation {
				ran <- request.Path
				<-release
			}
			return nil, nil
		}
		b.Root = []string{fmt.Sprintf("foo/%d", i)}
		meUUID, err := uuid.GenerateUUID()
		require.NoError(t, err)
		mountEntry := &MountEntry{
			Table:       mountTableType,
			UUID:        meUUID,
			Accessor:    fmt.Sprintf("accessor-%d", i),
			NamespaceID: namespace.RootNamespaceID,
			namespace:   namespace.RootNamespace,
			Path:        fmt.Sprintf("logical/foo/%d", i),
		}
		func() {
			core.mountsLock.Lock()
			defer core.mountsLock.Unlock()
			newTable := core.mounts.shallowClone()
			newTable.Entries = append(newTable.Entries, mountEntry)
			core.mounts = newTable
			err = core.router.Mount(b, "logical", mountEntry, view)
			require.NoError(t, core.persistMounts(context.Background(), newTable, &mountEntry.Local))
		}()
	}

	core.mountsLock.RLock()
	numMounts := len(core.mounts.Entries)
	core.mountsLock.RUnlock()

	timeout, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	got := make(map[string]bool)
	gotLock := sync.RWMutex{}
	hasMore := true
	for hasMore {
		// we're using 5 workers, so we would expect to see 5 writes to the
		// channel. Once that happens, close the release channel so that the
		// functions can exit and new rollback operations can run
		select {
		case <-timeout.Done():
			require.Fail(t, "test timed out")
		case i := <-ran:
			gotLock.Lock()
			got[i] = true
			numGot := len(got)
			gotLock.Unlock()
			if numGot == 5 {
				close(release)
				hasMore = false
			}
		}
	}
	done := make(chan struct{})

	// start a goroutine to consume the remaining items from the queued work
	go func() {
		for {
			select {
			case i := <-ran:
				gotLock.Lock()
				got[i] = true
				gotLock.Unlock()
			case <-done:
				return
			}
		}
	}()

	// wait for every mount to be rolled back at least once
	numMountsDone := 0
	for numMountsDone < numMounts {
		<-core.rollback.rollbacksDoneCh
		numMountsDone++
	}

	// stop the rollback worker, which will wait for all inflight rollbacks to
	// complete
	core.rollback.Stop()
	close(done)

	// we should have received at least 1 rollback for every backend
	gotLock.RLock()
	defer gotLock.RUnlock()
	require.GreaterOrEqual(t, len(got), 10)
}

// TestRollbackManager_numRollbackWorkers verifies that the number of rollback
// workers is parsed from the configuration, but can be overridden by an
// environment variable. This test cannot be run in parallel because of the
// environment variable
func TestRollbackManager_numRollbackWorkers(t *testing.T) {
	testCases := []struct {
		name          string
		configWorkers int
		setEnvVar     bool
		envVar        string
		wantWorkers   int
	}{
		{
			name:          "default in config",
			configWorkers: RollbackDefaultNumWorkers,
			wantWorkers:   RollbackDefaultNumWorkers,
		},
		{
			name:          "invalid envvar",
			configWorkers: RollbackDefaultNumWorkers,
			wantWorkers:   RollbackDefaultNumWorkers,
			setEnvVar:     true,
			envVar:        "invalid",
		},
		{
			name:          "envvar overrides config",
			configWorkers: RollbackDefaultNumWorkers,
			wantWorkers:   20,
			setEnvVar:     true,
			envVar:        "20",
		},
		{
			name:          "envvar negative",
			configWorkers: RollbackDefaultNumWorkers,
			wantWorkers:   RollbackDefaultNumWorkers,
			setEnvVar:     true,
			envVar:        "-1",
		},
		{
			name:          "envvar zero",
			configWorkers: RollbackDefaultNumWorkers,
			wantWorkers:   RollbackDefaultNumWorkers,
			setEnvVar:     true,
			envVar:        "0",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnvVar {
				t.Setenv(RollbackWorkersEnvVar, tc.envVar)
			}
			core := &Core{numRollbackWorkers: tc.configWorkers}
			r := &RollbackManager{logger: logger.Named("test"), core: core}
			require.Equal(t, tc.wantWorkers, r.numRollbackWorkers())
		})
	}
}

func TestRollbackManager_Join(t *testing.T) {
	m, backend := mockRollback(t)
	if len(backend.Paths) > 0 {
		t.Fatalf("bad: %#v", backend)
	}

	m.Start()
	defer m.Stop()

	wg := &sync.WaitGroup{}
	wg.Add(3)

	errCh := make(chan error, 3)
	go func() {
		defer wg.Done()
		err := m.Rollback(namespace.RootContext(nil), "foo")
		if err != nil {
			errCh <- err
		}
	}()

	go func() {
		defer wg.Done()
		err := m.Rollback(namespace.RootContext(nil), "foo")
		if err != nil {
			errCh <- err
		}
	}()

	go func() {
		defer wg.Done()
		err := m.Rollback(namespace.RootContext(nil), "foo")
		if err != nil {
			errCh <- err
		}
	}()
	wg.Wait()
	close(errCh)
	err := <-errCh
	if err != nil {
		t.Fatalf("Error on rollback:%v", err)
	}
}
