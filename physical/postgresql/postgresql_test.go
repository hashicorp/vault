package postgresql

import (
	"fmt"
	"os"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/physical"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
)

func TestPostgreSQLBackend(t *testing.T) {
	logger := logging.NewVaultLogger(log.Debug)

	// Use docker as pg backend if no url is provided via environment variables
	var cleanup func()
	connURL := os.Getenv("PGURL")
	if connURL == "" {
		cleanup, connURL = prepareTestContainer(t, logger)
		defer cleanup()
	}

	table := os.Getenv("PGTABLE")
	if table == "" {
		table = "vault_kv_store"
	}

	logger.Info(fmt.Sprintf("Connection URL: %v", connURL))

	b, err := NewPostgreSQLBackend(map[string]string{
		"connection_url": connURL,
		"table":          table,
	}, logger)
	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}
	pg := b.(*PostgreSQLBackend)

	//Read postgres version to test basic connects works
	var pgversion string
	if err = pg.client.QueryRow("SELECT current_setting('server_version_num')").Scan(&pgversion); err != nil {
		t.Fatalf("Failed to check for Postgres version: %v", err)
	}
	logger.Info(fmt.Sprintf("Postgres Version: %v", pgversion))

	//Setup tables and indexes if not exists.
	createTableSQL := fmt.Sprintf(
		"  CREATE TABLE IF NOT EXISTS %v ( "+
			"  parent_path TEXT COLLATE \"C\" NOT NULL, "+
			"  path        TEXT COLLATE \"C\", "+
			"  key         TEXT COLLATE \"C\", "+
			"  value       BYTEA, "+
			"  CONSTRAINT pkey PRIMARY KEY (path, key) "+
			" ); ", table)

	_, err = pg.client.Exec(createTableSQL)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	createIndexSQL := fmt.Sprintf(" CREATE INDEX IF NOT EXISTS parent_path_idx ON %v (parent_path); ", table)

	_, err = pg.client.Exec(createIndexSQL)
	if err != nil {
		t.Fatalf("Failed to create index: %v", err)
	}

	defer func() {
		pg := b.(*PostgreSQLBackend)
		_, err := pg.client.Exec(fmt.Sprintf(" TRUNCATE TABLE %v ", pg.table))
		if err != nil {
			t.Fatalf("Failed to truncate table: %v", err)
		}
	}()

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}

func prepareTestContainer(t *testing.T, logger log.Logger) (cleanup func(), retConnString string) {
	// If environment variable is set, use this connectionstring without starting docker container
	if os.Getenv("PGURL") != "" {
		return func() {}, os.Getenv("PGURL")
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}
	//using 11.1 which is currently latest, use hard version for stabillity of tests
	resource, err := pool.Run("postgres", "11.1", []string{})
	if err != nil {
		t.Fatalf("Could not start docker Postgres: %s", err)
	}

	retConnString = fmt.Sprintf("postgres://postgres@localhost:%v/postgres?sslmode=disable", resource.GetPort("5432/tcp"))

	cleanup = func() {
		err := pool.Purge(resource)
		if err != nil {
			t.Fatalf("Failed to cleanup docker Postgres: %s", err)
		}
	}

	// Provide a test function to the pool to test if docker instance service is up.
	// We try to setup a pg backend as test for successful connect
	// exponential backoff-retry, because the dockerinstance may not be able to accept
	// connections yet, test by trying to setup a postgres backend, max-timeout is 60s
	if err := pool.Retry(func() error {
		var err error
		_, err = NewPostgreSQLBackend(map[string]string{
			"connection_url": retConnString,
		}, logger)
		return err

	}); err != nil {
		cleanup()
		t.Fatalf("Could not connect to docker: %s", err)
	}

	return cleanup, retConnString
}
