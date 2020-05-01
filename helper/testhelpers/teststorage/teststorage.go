package teststorage

import (
	"fmt"
	"io/ioutil"
	"os"
	realtesting "testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/consul"
	vaulthttp "github.com/hashicorp/vault/http"
	physConsul "github.com/hashicorp/vault/physical/consul"
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

func MakeConsulBackend(t testing.T, logger hclog.Logger) *vault.PhysicalBackendBundle {
	cleanup, consulAddress, consulToken := consul.PrepareTestContainer(t.(*realtesting.T), "")
	consulConf := map[string]string{
		"address":      consulAddress,
		"token":        consulToken,
		"max_parallel": "32",
	}
	consulBackend, err := physConsul.NewConsulBackend(consulConf, logger)
	if err != nil {
		t.Fatal(err)
	}
	return &vault.PhysicalBackendBundle{
		Backend: consulBackend,
		Cleanup: cleanup,
	}
}

func MakeRaftBackend(t testing.T, coreIdx int, logger hclog.Logger) *vault.PhysicalBackendBundle {
	nodeID := fmt.Sprintf("core-%d", coreIdx)
	raftDir, err := ioutil.TempDir("", "vault-raft-")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("raft dir: %s", raftDir)
	cleanupFunc := func() {
		os.RemoveAll(raftDir)
	}

	logger.Info("raft dir", "dir", raftDir)

	conf := map[string]string{
		"path":                   raftDir,
		"node_id":                nodeID,
		"performance_multiplier": "8",
	}

	backend, err := raft.NewRaftBackend(conf, logger)
	if err != nil {
		cleanupFunc()
		t.Fatal(err)
	}

	return &vault.PhysicalBackendBundle{
		Backend: backend,
		Cleanup: cleanupFunc,
	}
}

type ClusterSetupMutator func(conf *vault.CoreConfig, opts *vault.TestClusterOptions)

func SharedPhysicalFactory(f func(t testing.T, logger hclog.Logger) *vault.PhysicalBackendBundle) func(t testing.T, coreIdx int, logger hclog.Logger) *vault.PhysicalBackendBundle {
	return func(t testing.T, coreIdx int, logger hclog.Logger) *vault.PhysicalBackendBundle {
		if coreIdx == 0 {
			return f(t, logger)
		}
		return nil
	}
}

func InmemBackendSetup(conf *vault.CoreConfig, opts *vault.TestClusterOptions) {
	opts.PhysicalFactory = SharedPhysicalFactory(MakeInmemBackend)
}
func InmemNonTransactionalBackendSetup(conf *vault.CoreConfig, opts *vault.TestClusterOptions) {
	opts.PhysicalFactory = SharedPhysicalFactory(MakeInmemNonTransactionalBackend)
}
func FileBackendSetup(conf *vault.CoreConfig, opts *vault.TestClusterOptions) {
	opts.PhysicalFactory = SharedPhysicalFactory(MakeFileBackend)
}
func ConsulBackendSetup(conf *vault.CoreConfig, opts *vault.TestClusterOptions) {
	opts.PhysicalFactory = SharedPhysicalFactory(MakeConsulBackend)
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
