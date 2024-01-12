// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package teststorage

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/audit"
	auditFile "github.com/hashicorp/vault/builtin/audit/file"
	auditSocket "github.com/hashicorp/vault/builtin/audit/socket"
	auditSyslog "github.com/hashicorp/vault/builtin/audit/syslog"
	logicalDb "github.com/hashicorp/vault/builtin/logical/database"
	"github.com/hashicorp/vault/builtin/plugin"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/logical"
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

func MakeLatentInmemBackend(t testing.T, logger hclog.Logger) *vault.PhysicalBackendBundle {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	jitter := r.Intn(20)
	latency := time.Duration(r.Intn(15)) * time.Millisecond

	pbb := MakeInmemBackend(t, logger)
	latencyInjector := physical.NewTransactionalLatencyInjector(pbb.Backend, latency, jitter, logger)
	pbb.Backend = latencyInjector
	return pbb
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

func MakeRaftBackend(t testing.T, coreIdx int, logger hclog.Logger, extraConf map[string]interface{}, bridge *raft.ClusterAddrBridge) *vault.PhysicalBackendBundle {
	nodeID := fmt.Sprintf("core-%d", coreIdx)
	raftDir, err := ioutil.TempDir("", "vault-raft-")
	if err != nil {
		t.Fatal(err)
	}
	// t.Logf("raft dir: %s", raftDir)
	cleanupFunc := func() {
		os.RemoveAll(raftDir)
	}

	logger.Info("raft dir", "dir", raftDir)

	backend, err := makeRaftBackend(logger, nodeID, raftDir, extraConf, bridge)
	if err != nil {
		cleanupFunc()
		t.Fatal(err)
	}

	return &vault.PhysicalBackendBundle{
		Backend: backend,
		Cleanup: cleanupFunc,
	}
}

func makeRaftBackend(logger hclog.Logger, nodeID, raftDir string, extraConf map[string]interface{}, bridge *raft.ClusterAddrBridge) (physical.Backend, error) {
	conf := map[string]string{
		"path":                         raftDir,
		"node_id":                      nodeID,
		"performance_multiplier":       "8",
		"autopilot_reconcile_interval": "300ms",
		"autopilot_update_interval":    "100ms",
	}
	for k, v := range extraConf {
		val, ok := v.(string)
		if ok {
			conf[k] = val
		}
	}

	backend, err := raft.NewRaftBackend(conf, logger.Named("raft"))
	if err != nil {
		return nil, err
	}
	if bridge != nil {
		backend.(*raft.RaftBackend).SetServerAddressProvider(bridge)
	}

	return backend, nil
}

// RaftHAFactory returns a PhysicalBackendBundle with raft set as the HABackend
// and the physical.Backend provided in PhysicalBackendBundler as the storage
// backend.
func RaftHAFactory(f PhysicalBackendBundler) func(t testing.T, coreIdx int, logger hclog.Logger, conf map[string]interface{}) *vault.PhysicalBackendBundle {
	return func(t testing.T, coreIdx int, logger hclog.Logger, conf map[string]interface{}) *vault.PhysicalBackendBundle {
		// Call the factory func to create the storage backend
		physFactory := SharedPhysicalFactory(f)
		bundle := physFactory(t, coreIdx, logger, nil)

		// This can happen if a shared physical backend is called on a non-0th core.
		if bundle == nil {
			bundle = new(vault.PhysicalBackendBundle)
		}

		raftDir := makeRaftDir(t)
		cleanupFunc := func() {
			os.RemoveAll(raftDir)
		}

		nodeID := fmt.Sprintf("core-%d", coreIdx)
		backendConf := map[string]string{
			"path":                         raftDir,
			"node_id":                      nodeID,
			"performance_multiplier":       "8",
			"autopilot_reconcile_interval": "300ms",
			"autopilot_update_interval":    "100ms",
		}

		// Create and set the HA Backend
		raftBackend, err := raft.NewRaftBackend(backendConf, logger)
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

func SharedPhysicalFactory(f PhysicalBackendBundler) func(t testing.T, coreIdx int, logger hclog.Logger, conf map[string]interface{}) *vault.PhysicalBackendBundle {
	return func(t testing.T, coreIdx int, logger hclog.Logger, conf map[string]interface{}) *vault.PhysicalBackendBundle {
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

func InmemLatentBackendSetup(conf *vault.CoreConfig, opts *vault.TestClusterOptions) {
	opts.PhysicalFactory = SharedPhysicalFactory(MakeLatentInmemBackend)
}

func InmemNonTransactionalBackendSetup(conf *vault.CoreConfig, opts *vault.TestClusterOptions) {
	opts.PhysicalFactory = SharedPhysicalFactory(MakeInmemNonTransactionalBackend)
}

func FileBackendSetup(conf *vault.CoreConfig, opts *vault.TestClusterOptions) {
	opts.PhysicalFactory = SharedPhysicalFactory(MakeFileBackend)
}

func RaftBackendSetup(conf *vault.CoreConfig, opts *vault.TestClusterOptions) {
	opts.KeepStandbysSealed = true
	var bridge *raft.ClusterAddrBridge
	if !opts.InmemClusterLayers && opts.ClusterLayers == nil {
		bridge = raft.NewClusterAddrBridge()
	}
	conf.ClusterAddrBridge = bridge
	opts.PhysicalFactory = func(t testing.T, coreIdx int, logger hclog.Logger, conf map[string]interface{}) *vault.PhysicalBackendBundle {
		return MakeRaftBackend(t, coreIdx, logger, conf, bridge)
	}
	opts.SetupFunc = func(t testing.T, c *vault.TestCluster) {
		if opts.NumCores != 1 {
			testhelpers.RaftClusterJoinNodes(t, c)
			time.Sleep(15 * time.Second)
		}
	}
}

func RaftHASetup(conf *vault.CoreConfig, opts *vault.TestClusterOptions, bundler PhysicalBackendBundler) {
	opts.InmemClusterLayers = true
	opts.PhysicalFactory = RaftHAFactory(bundler)
}

func ClusterSetup(conf *vault.CoreConfig, opts *vault.TestClusterOptions, setup ClusterSetupMutator) (*vault.CoreConfig, *vault.TestClusterOptions) {
	var localConf vault.CoreConfig
	localConf.DisableAutopilot = true
	if conf != nil {
		localConf = *conf
	}
	localOpts := vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		DefaultHandlerProperties: vault.HandlerProperties{
			ListenerConfig: &configutil.Listener{},
		},
	}
	if opts != nil {
		localOpts = *opts
	}
	if setup == nil {
		setup = InmemBackendSetup
	}
	setup(&localConf, &localOpts)
	if localConf.CredentialBackends == nil {
		localConf.CredentialBackends = map[string]logical.Factory{
			"plugin": plugin.Factory,
		}
	}
	if localConf.LogicalBackends == nil {
		localConf.LogicalBackends = map[string]logical.Factory{
			"plugin":   plugin.Factory,
			"database": logicalDb.Factory,
			// This is also available in the plugin catalog, but is here due to the need to
			// automatically mount it.
			"kv": logicalKv.Factory,
		}
	}
	if localConf.AuditBackends == nil {
		localConf.AuditBackends = map[string]audit.Factory{
			"file":   auditFile.Factory,
			"socket": auditSocket.Factory,
			"syslog": auditSyslog.Factory,
			"noop":   corehelpers.NoopAuditFactory(nil),
		}
	}

	return &localConf, &localOpts
}
