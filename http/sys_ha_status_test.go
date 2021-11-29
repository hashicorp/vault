package http

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/hashicorp/vault/vault"
)

func TestSysHAStatus(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)
	inm, err := inmem.NewTransactionalInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	inmha, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	conf := &vault.CoreConfig{
		Physical:   inm,
		HAPhysical: inmha.(physical.HABackend),
	}
	opts := &vault.TestClusterOptions{
		HandlerFunc: Handler,
	}
	cluster := vault.NewTestCluster(t, conf, opts)
	cluster.Start()
	defer cluster.Cleanup()
	testhelpers.WaitForActiveNodeAndStandbys(t, cluster)
	// Make sure standbys have time to echo and populate the cache
	time.Sleep(6 * time.Second)

	// Use standby deliberately to make sure it forwards
	client := cluster.Cores[1].Client
	r := client.NewRequest("GET", "/v1/sys/ha-status")
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := client.RawRequestWithContext(ctx, r)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var result HaStatusResponse
	err = resp.DecodeJSON(&result)
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Nodes) != len(cluster.Cores) {
		t.Fatalf("expected %d nodes, got %d", len(cluster.Cores), len(result.Nodes))
	}
}
