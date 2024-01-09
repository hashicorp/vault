// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package consul

import (
	"sync"
	realtesting "testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers/consul"
	physConsul "github.com/hashicorp/vault/physical/consul"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/go-testing-interface"
)

func MakeConsulBackend(t testing.T, logger hclog.Logger) *vault.PhysicalBackendBundle {
	cleanup, config := consul.PrepareTestContainer(t.(*realtesting.T), "", false, true)

	consulConf := map[string]string{
		"address":      config.Address(),
		"token":        config.Token,
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

func ConsulBackendSetup(conf *vault.CoreConfig, opts *vault.TestClusterOptions) {
	m := &consulContainerManager{}
	opts.PhysicalFactory = m.Backend
}

// consulContainerManager exposes Backend which matches the PhysicalFactory func
// type. When called, it will ensure that a separate Consul container is started
// for each distinct vault cluster that calls it and ensures that each Vault
// core gets a separate Consul backend instance since that contains state
// related to lock sessions. The whole test framework doesn't have a concept of
// "cluster names" outside of the prefix attached to the logger and other
// backend factories, mostly via SharedPhysicalFactory currently implicitly rely
// on being called in a sequence of core 0, 1, 2,... on one cluster and then
// core 0, 1, 2... on the next and so on. Refactoring lots of things to make
// first-class cluster identifiers a thing seems like a heavy lift given that we
// already rely on sequence of calls everywhere else anyway so we do the same
// here - each time the Backend method is called with coreIdx == 0 we create a
// whole new Consul and assume subsequent non 0 index cores are in the same
// cluster.
type consulContainerManager struct {
	mu      sync.Mutex
	current *consulContainerBackendFactory
}

func (m *consulContainerManager) Backend(t testing.T, coreIdx int,
	logger hclog.Logger, conf map[string]interface{},
) *vault.PhysicalBackendBundle {
	m.mu.Lock()
	if coreIdx == 0 || m.current == nil {
		// Create a new consul container factory
		m.current = &consulContainerBackendFactory{}
	}
	f := m.current
	m.mu.Unlock()

	return f.Backend(t, coreIdx, logger, conf)
}

type consulContainerBackendFactory struct {
	mu        sync.Mutex
	refCount  int
	cleanupFn func()
	config    map[string]string
}

func (f *consulContainerBackendFactory) Backend(t testing.T, coreIdx int,
	logger hclog.Logger, conf map[string]interface{},
) *vault.PhysicalBackendBundle {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.refCount == 0 {
		f.startContainerLocked(t)
		logger.Debug("started consul container", "clusterID", conf["cluster_id"],
			"address", f.config["address"])
	}

	f.refCount++
	consulBackend, err := physConsul.NewConsulBackend(f.config, logger.Named("consul"))
	if err != nil {
		t.Fatal(err)
	}
	return &vault.PhysicalBackendBundle{
		Backend: consulBackend,
		Cleanup: f.cleanup,
	}
}

func (f *consulContainerBackendFactory) startContainerLocked(t testing.T) {
	cleanup, config := consul.PrepareTestContainer(t.(*realtesting.T), "", false, true)
	f.config = map[string]string{
		"address":      config.Address(),
		"token":        config.Token,
		"max_parallel": "32",
	}
	f.cleanupFn = cleanup
}

func (f *consulContainerBackendFactory) cleanup() {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.refCount < 1 || f.cleanupFn == nil {
		return
	}
	f.refCount--
	if f.refCount == 0 {
		f.cleanupFn()
		f.cleanupFn = nil
	}
}
