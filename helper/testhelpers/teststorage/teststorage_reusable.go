package teststorage

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"
	"testing"

	mtesting "github.com/mitchellh/go-testing-interface"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/vault"
)

// ReusableStorage is a physical backend that can be re-used across
// multiple test clusters in sequence.
type ReusableStorage struct {
	IsRaft  bool
	Setup   ClusterSetupMutator
	Cleanup func(t *testing.T, cluster *vault.TestCluster)
}

// MakeReusableStorage makes a backend that can be re-used across
// multiple test clusters in sequence.
func MakeReusableStorage(t *testing.T, logger hclog.Logger, bundle *vault.PhysicalBackendBundle) (ReusableStorage, func()) {

	storage := ReusableStorage{
		IsRaft: false,

		Setup: func(conf *vault.CoreConfig, opts *vault.TestClusterOptions) {
			opts.PhysicalFactory = func(t mtesting.T, coreIdx int, logger hclog.Logger) *vault.PhysicalBackendBundle {
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

		Cleanup: func(t *testing.T, cluster *vault.TestCluster) {
		},
	}

	cleanup := func() {
		if bundle.Cleanup != nil {
			bundle.Cleanup()
		}
	}

	return storage, cleanup
}

// MakeReusableRaftStorage makes a raft backend that can be re-used across
// multiple test clusters in sequence.
func MakeReusableRaftStorage(t *testing.T, logger hclog.Logger) (ReusableStorage, func()) {

	raftDirs := make([]string, vault.DefaultNumCores)
	for i := 0; i < vault.DefaultNumCores; i++ {
		raftDirs[i] = makeRaftDir(t)
	}

	storage := ReusableStorage{
		IsRaft: true,

		Setup: func(conf *vault.CoreConfig, opts *vault.TestClusterOptions) {
			conf.DisablePerformanceStandby = true
			opts.KeepStandbysSealed = true
			opts.PhysicalFactory = func(t mtesting.T, coreIdx int, logger hclog.Logger) *vault.PhysicalBackendBundle {
				return makeReusableRaftBackend(t, coreIdx, logger, raftDirs[coreIdx])
			}
		},

		Cleanup: func(t *testing.T, cluster *vault.TestCluster) {
			// Close open files.
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

func makeRaftDir(t *testing.T) string {
	raftDir, err := ioutil.TempDir("", "vault-raft-")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("raft dir: %s", raftDir)
	return raftDir
}

func makeReusableRaftBackend(t mtesting.T, coreIdx int, logger hclog.Logger, raftDir string) *vault.PhysicalBackendBundle {

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

//// SetRaftAddressProviders sets a ServerAddressProvider for all the raft cores
//func SetRaftAddressProviders(t mtesting.T, cluster *vault.TestCluster) {
//
//	addressProvider := &testhelpers.TestRaftServerAddressProvider{Cluster: cluster}
//	atomic.StoreUint32(&vault.UpdateClusterAddrForTests, 1)
//
//	for _, core := range cluster.Cores {
//		core.UnderlyingRawStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
//	}
//}
//
//// JoinRaftFollower joins a follower to the cluster
//func JoinRaftFollower(t *testing.T, cluster *vault.TestCluster, leader, follower *vault.TestClusterCore) {
//
//	info := []*raft.LeaderJoinInfo{
//		&raft.LeaderJoinInfo{
//			LeaderAPIAddr: leader.Client.Address(),
//			TLSConfig:     leader.TLSConfig,
//		},
//	}
//
//	ctx := namespace.RootContext(context.Background())
//	_, err := follower.JoinRaftCluster(ctx, info, false)
//	if err != nil {
//		t.Fatal(err)
//	}
//}
