package teststorage

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/physical"
	physFile "github.com/hashicorp/vault/sdk/physical/file"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/go-testing-interface"
)

func MakeInmemBackend(t testing.T, logger hclog.Logger) *vault.PhysicalBackendBundle {
	inm, err := inmem.NewTransactionalInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	inmha, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	return &vault.PhysicalBackendBundle{
		Backend:   inm,
		HABackend: inmha.(physical.HABackend),
	}
}

func MakeInmemNonTransactionalBackend(t testing.T, logger hclog.Logger) *vault.PhysicalBackendBundle {
	inm, err := inmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	inmha, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	return &vault.PhysicalBackendBundle{
		Backend:   inm,
		HABackend: inmha.(physical.HABackend),
	}
}

func MakeFileBackend(t testing.T, logger hclog.Logger) *vault.PhysicalBackendBundle {
	path, err := ioutil.TempDir("", "vault-integ-file-")
	if err != nil {
		t.Fatal(err)
	}
	fileConf := map[string]string{
		"path": path,
	}
	fileBackend, err := physFile.NewTransactionalFileBackend(fileConf, logger)
	if err != nil {
		t.Fatal(err)
	}

	inmha, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	return &vault.PhysicalBackendBundle{
		Backend:   fileBackend,
		HABackend: inmha.(physical.HABackend),
		Cleanup: func() {
			err := os.RemoveAll(path)
			if err != nil {
				t.Fatal(err)
			}
		},
	}
}

func MakeRaftBackend(t testing.T, coreIdx int, logger hclog.Logger) *vault.PhysicalBackendBundle {
	return MakeRaftBackendWithConf(t, coreIdx, logger, nil)
}

func MakeRaftWithAutopilotBackend(t testing.T, coreIdx int, logger hclog.Logger) *vault.PhysicalBackendBundle {
	extraConf := map[string]string{
		"autopilot": `[{"cleanup_dead_servers":true,"last_contact_threshold":"5s","max_trailing_logs":500,"min_quorum":3,"server_stabilization_time":"10s"}]`,
	}
	return MakeRaftBackendWithConf(t, coreIdx, logger, extraConf)
}

func MakeRaftBackendWithConf(t testing.T, coreIdx int, logger hclog.Logger, extraConf map[string]string) *vault.PhysicalBackendBundle {
	nodeID := fmt.Sprintf("core-%d", coreIdx)
	raftDir, err := ioutil.TempDir("", "vault-raft-")
	if err != nil {
		t.Fatal(err)
	}
	//t.Logf("raft dir: %s", raftDir)
	cleanupFunc := func() {
		os.RemoveAll(raftDir)
	}

	logger.Info("raft dir", "dir", raftDir)

	conf := map[string]string{
		"path":                   raftDir,
		"node_id":                nodeID,
		"performance_multiplier": "8",
	}
	for k, v := range extraConf {
		conf[k] = v
	}

	backend, err := raft.NewRaftBackend(conf, logger.Named("raft"))
	if err != nil {
		cleanupFunc()
		t.Fatal(err)
	}

	return &vault.PhysicalBackendBundle{
		Backend: backend,
		Cleanup: cleanupFunc,
	}
}

// RaftHAFactory returns a PhysicalBackendBundle with raft set as the HABackend
// and the physical.Backend provided in PhysicalBackendBundler as the storage
// backend.
func RaftHAFactory(f PhysicalBackendBundler) func(t testing.T, coreIdx int, logger hclog.Logger) *vault.PhysicalBackendBundle {
	return func(t testing.T, coreIdx int, logger hclog.Logger) *vault.PhysicalBackendBundle {
		// Call the factory func to create the storage backend
		physFactory := SharedPhysicalFactory(f)
		bundle := physFactory(t, coreIdx, logger)

		// This can happen if a shared physical backend is called on a non-0th core.
		if bundle == nil {
			bundle = new(vault.PhysicalBackendBundle)
		}

		raftDir := makeRaftDir(t)
		cleanupFunc := func() {
			os.RemoveAll(raftDir)
		}

		nodeID := fmt.Sprintf("core-%d", coreIdx)
		conf := map[string]string{
			"path":                   raftDir,
			"node_id":                nodeID,
			"performance_multiplier": "8",
		}

		// Create and set the HA Backend
		raftBackend, err := raft.NewRaftBackend(conf, logger)
		if err != nil {
			bundle.Cleanup()
			t.Fatal(err)
		}
		bundle.HABackend = raftBackend.(physical.HABackend)

		// Re-wrap the cleanup func
		bundleCleanup := bundle.Cleanup
		bundle.Cleanup = func() {
			if bundleCleanup != nil {
				bundleCleanup()
			}
			cleanupFunc()
		}

		return bundle
	}
}

type PhysicalBackendBundler func(t testing.T, logger hclog.Logger) *vault.PhysicalBackendBundle

func SharedPhysicalFactory(f PhysicalBackendBundler) func(t testing.T, coreIdx int, logger hclog.Logger) *vault.PhysicalBackendBundle {
	return func(t testing.T, coreIdx int, logger hclog.Logger) *vault.PhysicalBackendBundle {
		if coreIdx == 0 {
			return f(t, logger)
		}
		return nil
	}
}

type ClusterSetupMutator func(conf *vault.CoreConfig, opts *vault.TestClusterOptions)

func InmemBackendSetup(conf *vault.CoreConfig, opts *vault.TestClusterOptions) {
	opts.PhysicalFactory = SharedPhysicalFactory(MakeInmemBackend)
}
func InmemNonTransactionalBackendSetup(conf *vault.CoreConfig, opts *vault.TestClusterOptions) {
	opts.PhysicalFactory = SharedPhysicalFactory(MakeInmemNonTransactionalBackend)
}
func FileBackendSetup(conf *vault.CoreConfig, opts *vault.TestClusterOptions) {
	opts.PhysicalFactory = SharedPhysicalFactory(MakeFileBackend)
}

func RaftBackendSetup(conf *vault.CoreConfig, opts *vault.TestClusterOptions) {
	conf.DisablePerformanceStandby = true
	opts.KeepStandbysSealed = true
	opts.PhysicalFactory = MakeRaftBackend
	opts.SetupFunc = func(t testing.T, c *vault.TestCluster) {
		if opts.NumCores != 1 {
			testhelpers.RaftClusterJoinNodes(t, c)
			time.Sleep(15 * time.Second)
		}
	}
}

func RaftBackendWithAutopilotSetup(conf *vault.CoreConfig, opts *vault.TestClusterOptions) {
	opts.KeepStandbysSealed = true
	opts.PhysicalFactory = MakeRaftWithAutopilotBackend
	opts.SetupFunc = func(t testing.T, c *vault.TestCluster) {
		if opts.NumCores != 1 {
			testhelpers.RaftClusterJoinNodes(t, c)
			time.Sleep(15 * time.Second)
		}
	}
}

func RaftHASetup(conf *vault.CoreConfig, opts *vault.TestClusterOptions, bundler PhysicalBackendBundler) {
	opts.KeepStandbysSealed = true
	opts.PhysicalFactory = RaftHAFactory(bundler)
}

func ClusterSetup(conf *vault.CoreConfig, opts *vault.TestClusterOptions, setup ClusterSetupMutator) (*vault.CoreConfig, *vault.TestClusterOptions) {
	var localConf vault.CoreConfig
	if conf != nil {
		localConf = *conf
	}
	localOpts := vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	}
	if opts != nil {
		localOpts = *opts
	}
	if setup == nil {
		setup = InmemBackendSetup
	}
	setup(&localConf, &localOpts)
	return &localConf, &localOpts
}
