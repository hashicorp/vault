package command

import (
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"

	vaulthttp "github.com/hashicorp/vault/http"
	logxi "github.com/mgutz/logxi/v1"
)

var testVaultServerDefaultBackends = map[string]logical.Factory{
	"transit": transit.Factory,
	"pki":     pki.Factory,
}

func testVaultServer(t testing.TB) (*api.Client, func()) {
	return testVaultServerBackends(t, testVaultServerDefaultBackends)
}

func testVaultServerBackends(t testing.TB, backends map[string]logical.Factory) (*api.Client, func()) {
	coreConfig := &vault.CoreConfig{
		DisableMlock:    true,
		DisableCache:    true,
		Logger:          logxi.NullLog,
		LogicalBackends: backends,
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()

	// make it easy to get access to the active
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	client := cluster.Cores[0].Client
	client.SetToken(cluster.RootToken)

	// Sanity check
	secret, err := client.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatal(err)
	}
	if secret == nil || secret.Data["id"].(string) != cluster.RootToken {
		t.Fatalf("token mismatch: %#v vs %q", secret, cluster.RootToken)
	}
	return client, func() { defer cluster.Cleanup() }
}

func testClient(t *testing.T, addr string, token string) *api.Client {
	config := api.DefaultConfig()
	config.Address = addr
	client, err := api.NewClient(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	client.SetToken(token)

	return client
}
