package api_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"

	vaulthttp "github.com/hashicorp/vault/http"
	logxi "github.com/mgutz/logxi/v1"
	dockertest "gopkg.in/ory-am/dockertest.v3"
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

// testPostgresDB creates a testing postgres database in a Docker container,
// returning the connection URL and the associated closer function.
func testPostgresDB(t testing.TB) (string, func()) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("postgresdb: failed to connect to docker: %s", err)
	}

	resource, err := pool.Run("postgres", "latest", []string{
		"POSTGRES_PASSWORD=secret",
		"POSTGRES_DB=database",
	})
	if err != nil {
		t.Fatalf("postgresdb: could not start container: %s", err)
	}

	addr := fmt.Sprintf("postgres://postgres:secret@localhost:%s/database?sslmode=disable", resource.GetPort("5432/tcp"))

	if err := pool.Retry(func() error {
		db, err := sql.Open("postgres", addr)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		t.Fatalf("postgresdb: could not connect: %s", err)
	}

	return addr, func() {
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("postgresdb: failed to cleanup container: %s", err)
		}
	}
}
