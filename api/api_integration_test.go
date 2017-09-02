package api_test

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/builtin/logical/database"
	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"

	auditFile "github.com/hashicorp/vault/builtin/audit/file"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	vaulthttp "github.com/hashicorp/vault/http"
	logxi "github.com/mgutz/logxi/v1"
	dockertest "gopkg.in/ory-am/dockertest.v3"
)

// testVaultServer creates a test vault cluster and returns a configured API
// client and closer function.
func testVaultServer(t testing.TB) (*api.Client, func()) {
	t.Helper()

	client, _, closer := testVaultServerUnseal(t)
	return client, closer
}

// testVaultServerUnseal creates a test vault cluster and returns a configured
// API client, list of unseal keys (as strings), and a closer function.
func testVaultServerUnseal(t testing.TB) (*api.Client, []string, func()) {
	t.Helper()

	return testVaultServerCoreConfig(t, &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       logxi.NullLog,
		CredentialBackends: map[string]logical.Factory{
			"userpass": credUserpass.Factory,
		},
		AuditBackends: map[string]audit.Factory{
			"file": auditFile.Factory,
		},
		LogicalBackends: map[string]logical.Factory{
			"database":       database.Factory,
			"generic-leased": vault.LeasedPassthroughBackendFactory,
			"pki":            pki.Factory,
			"transit":        transit.Factory,
		},
	})
}

// testVaultServerCoreConfig creates a new vault cluster with the given core
// configuration. This is a lower-level test helper.
func testVaultServerCoreConfig(t testing.TB, coreConfig *vault.CoreConfig) (*api.Client, []string, func()) {
	t.Helper()

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()

	// Make it easy to get access to the active
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	// Get the client already setup for us!
	client := cluster.Cores[0].Client
	client.SetToken(cluster.RootToken)

	// Convert the unseal keys to base64 encoded, since these are how the user
	// will get them.
	unsealKeys := make([]string, len(cluster.BarrierKeys))
	for i := range unsealKeys {
		unsealKeys[i] = base64.StdEncoding.EncodeToString(cluster.BarrierKeys[i])
	}

	return client, unsealKeys, func() { defer cluster.Cleanup() }
}

// testVaultServerBad creates an http server that returns a 500 on each request
// to simulate failures.
func testVaultServerBad(t testing.TB) (*api.Client, func()) {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	server := &http.Server{
		Addr: "127.0.0.1:0",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
		}),
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			t.Fatal(err)
		}
	}()

	client, err := api.NewClient(&api.Config{
		Address: "http://" + listener.Addr().String(),
	})
	if err != nil {
		t.Fatal(err)
	}

	return client, func() {
		ctx, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()

		server.Shutdown(ctx)
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
