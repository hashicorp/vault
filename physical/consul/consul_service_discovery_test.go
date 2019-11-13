package consul

import (
	"sync"
	"testing"
	"time"

	"github.com/go-test/deep"
	consulapi "github.com/hashicorp/consul/api"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers/consul"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/hashicorp/vault/vault"
)

// TestConsul_ServiceDiscovery tests whether consul ServiceDiscovery works
func TestConsul_ServiceDiscovery(t *testing.T) {

	logger := logging.NewVaultLogger(log.Trace)

	// Prepare a docker-based consul instance
	cleanup, connURL, connToken := consul.PrepareTestContainer(t, "1.4.0-rc1")
	defer cleanup()

	// Create a consul client
	consulCfg := consulapi.DefaultConfig()
	consulCfg.Address = connURL
	consulCfg.Token = connToken
	client, err := consulapi.NewClient(consulCfg)
	if err != nil {
		t.Fatal(err)
	}

	// checkServices is a helper function that checks whether the Consul
	// catalog is in the state that we expect.
	checkServices := func(t *testing.T, expected map[string][]string) {
		t.Helper()
		services, _, err := client.Catalog().Services(nil)
		if err != nil {
			t.Fatal(err)
		}
		if diff := deep.Equal(services, expected); diff != nil {
			t.Fatal(diff)
		}
	}

	// awaitServicesTransition is a helper function that waits patiently for
	// the Consul catalog to transition from a known state to some other state.
	awaitServicesTransition := func(t *testing.T, from map[string][]string) {
		t.Helper()
		// Wait for up to 10 seconds
		for i := 0; i < 10; i++ {
			services, _, err := client.Catalog().Services(nil)
			if err != nil {
				t.Fatal(err)
			}
			if diff := deep.Equal(services, from); diff != nil {
				return
			}
			time.Sleep(time.Second)
		}
		t.Fatalf("Catalog Services never transitioned from %v", from)
	}

	// Create a ServiceDiscovery that points to our consul instance
	sd, err := NewConsulServiceDiscovery(map[string]string{
		"address": connURL,
		"token":   connToken,
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
	checkServices(t, map[string][]string{
		"consul": []string{},
	})

	// Run service discovery on the core
	wg := &sync.WaitGroup{}
	var shutdown chan struct{}
	activeFunc := func() bool {
		if isLeader, _, _, err := core.Leader(); err == nil {
			return isLeader
		}
		return false
	}
	err = sd.RunServiceDiscovery(wg, shutdown, redirectAddr, activeFunc, core.Sealed, core.PerfStandby)
	if err != nil {
		t.Fatal(err)
	}

	// Vault should now be registered with Consul in standby mode
	awaitServicesTransition(t, map[string][]string{
		"consul": []string{},
	})
	checkServices(t, map[string][]string{
		"consul": []string{},
		"vault":  []string{"standby"},
	})

	// Initialize the core
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

	// Vault should now be registered with Consul in active mode
	awaitServicesTransition(t, map[string][]string{
		"consul": []string{},
		"vault":  []string{"standby"},
	})
	checkServices(t, map[string][]string{
		"consul": []string{},
		"vault":  []string{"active"},
	})
}
