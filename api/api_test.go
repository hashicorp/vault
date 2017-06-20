package api

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"sync/atomic"
	"testing"
	"time"

	_ "github.com/lib/pq"
	dockertest "gopkg.in/ory-am/dockertest.v3"

	"golang.org/x/net/http2"
)

// testHTTPServer creates a test HTTP server that handles requests until
// the listener returned is closed.
func testHTTPServer(
	t *testing.T, handler http.Handler) (*Config, net.Listener) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	server := &http.Server{Handler: handler}
	if err := http2.ConfigureServer(server, nil); err != nil {
		t.Fatal(err)
	}
	go server.Serve(ln)

	config := DefaultConfig()
	config.Address = fmt.Sprintf("http://%s", ln.Addr())

	return config, ln
}

// nextPort is the next port to use for the API server.
var nextPort int32 = 28200

// restVaultServer runs an instance of the Vault server in development mode.
// This requires that the vault binary is installed and in the $PATH.
func testVaultServer(t *testing.T) (*Client, func()) {
	bin, err := exec.LookPath("vault")
	if err != nil || bin == "" {
		t.Fatal("vault binary not found")
	}

	// Get the port number
	port := atomic.AddInt32(&nextPort, 1)

	// Construct the address
	addr := fmt.Sprintf("127.0.0.1:%d", port)

	// Start the server
	cmd := exec.Command(
		bin, "server", "-dev",
		"-dev-listen-address", addr,
		"-dev-root-token-id", "root",
	)
	if err := cmd.Start(); err != nil {
		t.Fatalf("err: %s", err)
	}

	for i := 0; i < 10; i++ {
		conn, err := net.DialTimeout("tcp", addr, time.Second)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		conn.Close()

		config := DefaultConfig()
		config.Address = fmt.Sprintf("http://%s", addr)
		client, err := NewClient(config)
		if err != nil {
			t.Fatal(err)
		}
		client.SetToken("root")

		return client, func() {
			cmd.Process.Signal(os.Interrupt)
			cmd.Process.Wait()
		}
	}

	t.Fatalf("timeout waiting for vault server")
	return nil, nil
}

func testPostgresDatabase(t *testing.T) (string, func()) {
	if os.Getenv("PG_URL") != "" {
		return os.Getenv("PG_URL"), func() {}
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	resource, err := pool.Run("postgres", "latest", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=database"})
	if err != nil {
		t.Fatalf("Could not start local PostgreSQL docker container: %s", err)
	}

	cleanup := func() {
		err := pool.Purge(resource)
		if err != nil {
			t.Fatalf("Failed to cleanup local container: %s", err)
		}
	}

	pgURL := fmt.Sprintf("postgres://postgres:secret@localhost:%s/database?sslmode=disable", resource.GetPort("5432/tcp"))

	// exponential backoff-retry
	if err := pool.Retry(func() error {
		db, err := sql.Open("postgres", pgURL)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		t.Fatalf("Could not connect to PostgreSQL docker container: %s", err)
	}

	return pgURL, cleanup
}
