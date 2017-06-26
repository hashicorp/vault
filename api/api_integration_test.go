package api_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
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
	handlers := []http.Handler{
		http.NewServeMux(),
		http.NewServeMux(),
		http.NewServeMux(),
	}

	coreConfig := &vault.CoreConfig{
		DisableMlock:    true,
		DisableCache:    true,
		Logger:          logxi.NullLog,
		LogicalBackends: backends,
	}

	// Chicken-and-egg: Handler needs a core. So we create handlers first, then
	// add routes chained to a Handler-created handler.
	cores := vault.TestCluster(t, handlers, coreConfig, true)
	for i, core := range cores {
		handlers[i].(*http.ServeMux).Handle("/", vaulthttp.Handler(core.Core))
	}

	// make it easy to get access to the active
	core := cores[0].Core
	vault.TestWaitActive(t, core)

	rootToken := cores[0].Root
	address := fmt.Sprintf("https://127.0.0.1:%d", cores[1].Listeners[0].Address.Port)

	config := api.DefaultConfig()
	config.Address = address
	config.HttpClient = cleanhttp.DefaultClient()
	config.HttpClient.Transport.(*http.Transport).TLSClientConfig = cores[0].TLSConfig
	client, err := api.NewClient(config)
	if err != nil {
		t.Fatalf("error creating vault cluster: %s", err)
	}
	client.SetToken(rootToken)

	// Sanity check
	secret, err := client.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatal(err)
	}
	if secret == nil || secret.Data["id"].(string) != rootToken {
		t.Fatalf("token mismatch: %q vs %q", secret, rootToken)
	}

	return client, func() {
		for _, core := range cores {
			defer core.CloseListeners()
		}
	}
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
