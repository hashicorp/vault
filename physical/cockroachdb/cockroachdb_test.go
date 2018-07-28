package cockroachdb

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/ory/dockertest"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/physical"

	_ "github.com/lib/pq"
)

func prepareCockroachDBTestContainer(t *testing.T) (cleanup func(), retURL, tableName string) {
	tableName = os.Getenv("CR_TABLE")
	if tableName == "" {
		tableName = "vault_kv_store"
	}
	retURL = os.Getenv("CR_URL")
	if retURL != "" {
		return func() {}, retURL, tableName
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	dockerOptions := &dockertest.RunOptions{
		Repository: "cockroachdb/cockroach",
		Tag:        "release-1.0",
		Cmd:        []string{"start", "--insecure"},
	}
	resource, err := pool.RunWithOptions(dockerOptions)
	if err != nil {
		t.Fatalf("Could not start local CockroachDB docker container: %s", err)
	}

	cleanup = func() {
		err := pool.Purge(resource)
		if err != nil {
			t.Fatalf("Failed to cleanup local container: %s", err)
		}
	}

	retURL = fmt.Sprintf("postgresql://root@localhost:%s/?sslmode=disable", resource.GetPort("26257/tcp"))
	database := "database"
	tableName = database + ".vault_kv"

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		var err error
		db, err := sql.Open("postgres", retURL)
		if err != nil {
			return err
		}
		defer db.Close()
		_, err = db.Exec("CREATE DATABASE database")
		return err
	}); err != nil {
		cleanup()
		t.Fatalf("Could not connect to docker: %s", err)
	}
	return cleanup, retURL, tableName
}

func TestCockroachDBBackend(t *testing.T) {
	cleanup, connURL, table := prepareCockroachDBTestContainer(t)
	defer cleanup()

	// Run vault tests
	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewCockroachDBBackend(map[string]string{
		"connection_url": connURL,
		"table":          table,
	}, logger)

	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	defer func() {
		truncate(t, b)
	}()

	physical.ExerciseBackend(t, b)
	truncate(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
	truncate(t, b)
	physical.ExerciseTransactionalBackend(t, b)
}

func truncate(t *testing.T, b physical.Backend) {
	crdb := b.(*CockroachDBBackend)
	_, err := crdb.client.Exec("TRUNCATE TABLE " + crdb.table)
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}
}
