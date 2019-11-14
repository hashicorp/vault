package consul

import (
	"sync"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/consul/api"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers/consul"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/hashicorp/vault/vault"
)

// TestConsul_ServiceDiscovery tests whether consul ServiceDiscovery works
func TestConsul_ServiceDiscovery(t *testing.T) {

	// Prepare a docker-based consul instance
	cleanup, addr, token := consul.PrepareTestContainer(t, "1.4.0-rc1")
	defer cleanup()

	// Create a consul client
	cfg := api.DefaultConfig()
	cfg.Address = addr
	cfg.Token = token
	client, err := api.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	// transitionFrom waits patiently for the services in the Consul catalog to
	// transition from a known value, and then returns the new value.
	transitionFrom := func(t *testing.T, known map[string][]string) map[string][]string {
		t.Helper()
		// Wait for up to 10 seconds
		for i := 0; i < 10; i++ {
			services, _, err := client.Catalog().Services(nil)
			if err != nil {
				t.Fatal(err)
			}
			if diff := deep.Equal(services, known); diff != nil {
				return services
			}
			time.Sleep(time.Second)
		}
		t.Fatalf("Catalog Services never transitioned from %v", known)
		return nil
	}

	// Create a ServiceDiscovery that points to our consul instance
	logger := logging.NewVaultLogger(log.Trace)
	sd, err := NewConsulServiceDiscovery(map[string]string{
		"address": addr,
		"token":   token,
	}, logger)
	if err != nil {
		t.Fatal(err)
	}

	// Create the core
	inm, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	inmha, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	const redirectAddr = "http://127.0.0.1:8200"
	core, err := vault.NewCore(&vault.CoreConfig{
		ConfigServiceDiscovery: sd,
		Physical:               inm,
		HAPhysical:             inmha.(physical.HABackend),
		RedirectAddr:           redirectAddr,
		DisableMlock:           true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Vault should not yet be registered with Consul
	services, _, err := client.Catalog().Services(nil)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(services, map[string][]string{
		"consul": []string{},
	}); diff != nil {
		t.Fatal(diff)
	}

	// Run service discovery on the core
	wg := &sync.WaitGroup{}
	var shutdown chan struct{}
	activeFunc := func() bool {
		if isLeader, _, _, err := core.Leader(); err == nil {
			return isLeader
		}
		return false
	}
	err = sd.RunServiceDiscovery(
		wg, shutdown, redirectAddr, activeFunc, core.Sealed, core.PerfStandby)
	if err != nil {
		t.Fatal(err)
	}

	// Vault should soon be registered with Consul in standby mode
	services = transitionFrom(t, map[string][]string{
		"consul": []string{},
	})
	if diff := deep.Equal(services, map[string][]string{
		"consul": []string{},
		"vault":  []string{"standby"},
	}); diff != nil {
		t.Fatal(diff)
	}

	// Initialize and unseal the core
	keys, _ := vault.TestCoreInit(t, core)
	for _, key := range keys {
		if _, err := vault.TestCoreUnseal(core, vault.TestKeyCopy(key)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}
	if core.Sealed() {
		t.Fatal("should not be sealed")
	}

	// Wait for the core to become active
	vault.TestWaitActive(t, core)

	// Vault should soon be registered with Consul in active mode
	services = transitionFrom(t, map[string][]string{
		"consul": []string{},
		"vault":  []string{"standby"},
	})
	if diff := deep.Equal(services, map[string][]string{
		"consul": []string{},
		"vault":  []string{"active"},
	}); diff != nil {
		t.Fatal(diff)
	}
}
