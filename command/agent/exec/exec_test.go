package exec

import (
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"

	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
)

func testVaultServer(t *testing.T) (*api.Client, func()) {
	t.Helper()

	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       log.NewNullLogger(),
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()

	cores := cluster.Cores

	vault.TestWaitActive(t, cores[0].Core)

	client := cores[0].Client

	// enable kv-v2 backend
	if err := client.Sys().Mount("kv/", &api.MountInput{
		Type: "kv-v2",
	}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)

	return client, cluster.Cleanup
}

func TestServer_Run_Golden(t *testing.T) {
	_, closer := testVaultServer(t)
	defer closer()

}
