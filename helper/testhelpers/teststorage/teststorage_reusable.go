package teststorage

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"

	"github.com/mitchellh/go-testing-interface"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/vault"
)

// ReusableStorage is a physical backend that can be re-used across
// multiple test clusters in sequence.  It is useful for testing things like
// seal migration, wherein a given physical backend must be re-used as several
// test clusters are sequentially created, tested, and discarded.
type ReusableStorage struct {

	// IsRaft specifies whether the storage is using a raft backend.
	IsRaft bool

	// Setup should be called just before a new TestCluster is created.
	Setup ClusterSetupMutator

	// Cleanup should be called after a TestCluster is no longer
	// needed -- generally in a defer, just before the call to
	// cluster.Cleanup().
	Cleanup func(t testing.T, cluster *vault.TestCluster)
}

// StorageCleanup is a function that should be called once -- at the very end
// of a given unit test, after each of the sequence of clusters have been
// created, tested, and discarded.
type StorageCleanup func()

// MakeReusableStorage makes a physical backend that can be re-used across
// multiple test clusters in sequence.
func MakeReusableStorage(t testing.T, logger hclog.Logger, bundle *vault.PhysicalBackendBundle) (ReusableStorage, StorageCleanup) {

	storage := ReusableStorage{
		IsRaft: false,

		Setup: func(conf *vault.CoreConfig, opts *vault.TestClusterOptions) {
			opts.PhysicalFactory = func(t testing.T, coreIdx int, logger hclog.Logger) *vault.PhysicalBackendBundle {
				if coreIdx == 0 {
					// We intentionally do not clone the backend's Cleanup func,
					// because we don't want it to be run until the entire test has
					// been completed.
					return &vault.PhysicalBackendBundle{
						Backend:   bundle.Backend,
						HABackend: bundle.HABackend,
					}
				}
				return nil
			}
		},

		// No-op
		Cleanup: func(t testing.T, cluster *vault.TestCluster) {
		},
	}

	cleanup := func() {
		if bundle.Cleanup != nil {
			bundle.Cleanup()
		}
	}

	return storage, cleanup
}

// MakeReusableRaftStorage makes a physical raft backend that can be re-used
// across multiple test clusters in sequence.
func MakeReusableRaftStorage(t testing.T, logger hclog.Logger) (ReusableStorage, StorageCleanup) {

	raftDirs := make([]string, vault.DefaultNumCores)
	for i := 0; i < vault.DefaultNumCores; i++ {
		raftDirs[i] = makeRaftDir(t)
	}

	storage := ReusableStorage{
		IsRaft: true,

		Setup: func(conf *vault.CoreConfig, opts *vault.TestClusterOptions) {
			conf.DisablePerformanceStandby = true
			opts.KeepStandbysSealed = true
			opts.PhysicalFactory = func(t testing.T, coreIdx int, logger hclog.Logger) *vault.PhysicalBackendBundle {
				return makeReusableRaftBackend(t, coreIdx, logger, raftDirs[coreIdx])
			}
		},

		// Close open files being used by raft.
		Cleanup: func(t testing.T, cluster *vault.TestCluster) {
			for _, core := range cluster.Cores {
				raftStorage := core.UnderlyingRawStorage.(*raft.RaftBackend)
				if err := raftStorage.Close(); err != nil {
					t.Fatal(err)
				}
			}
		},
	}

	cleanup := func() {
		for _, rd := range raftDirs {
			os.RemoveAll(rd)
		}
	}

	return storage, cleanup
}

func makeRaftDir(t testing.T) string {
	raftDir, err := ioutil.TempDir("", "vault-raft-")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("raft dir: %s", raftDir)
	return raftDir
}

func makeReusableRaftBackend(t testing.T, coreIdx int, logger hclog.Logger, raftDir string) *vault.PhysicalBackendBundle {

	nodeID := fmt.Sprintf("core-%d", coreIdx)
	conf := map[string]string{
		"path":                   raftDir,
		"node_id":                nodeID,
		"performance_multiplier": "8",
	}

	backend, err := raft.NewRaftBackend(conf, logger)
	if err != nil {
		debug.PrintStack()
		t.Fatal(err)
	}

	return &vault.PhysicalBackendBundle{
		Backend: backend,
	}
}
