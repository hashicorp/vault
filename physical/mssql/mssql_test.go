package mssql

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/physical"
	"github.com/ory/dockertest"

	_ "github.com/denisenkom/go-mssqldb"
)

var conf = map[string]string{
	"database": "master",
	"server":   "localhost",
	"username": "SA",
	"password": "pa$$w0rd!",
}

func prepareMSSQLTestContainer(t *testing.T) (cleanup func(), retURL string, retPort string) {
	if os.Getenv("MSSQL_URL") != "" {
		return func() {}, os.Getenv("MSSQL_URL"), ""
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	ro := &dockertest.RunOptions{
		Repository: "mcr.microsoft.com/mssql/server",
		Tag:        "latest",
		Env:        []string{"ACCEPT_EULA=Y", "SA_PASSWORD=pa$$w0rd!"},
	}
	resource, err := pool.RunWithOptions(ro)
	if err != nil {
		t.Fatalf("Could not start local mssql docker container: %s", err)
	}

	cleanup = func() {
		err := pool.Purge(resource)
		if err != nil {
			t.Fatalf("Failed to cleanup local container: %s", err)
		}
	}

	retPort = resource.GetPort("1433/tcp")
	retURL = fmt.Sprintf("sqlserver://SA:pa$$w0rd!@localhost:%s", retPort)

	// exponential backoff-retry
	if retryErr := pool.Retry(func() error {
		db, err := sql.Open("sqlserver", retURL)
		if err != nil {
			return err
		}
		return db.Ping()

	}); retryErr != nil {
		cleanup()
		t.Fatalf("Could not connect to mssql docker container: %s", err)
	}

	return
}

func TestMSSQLBackendWithConnectionURL(t *testing.T) {
	cleanup, connURL, _ := prepareMSSQLTestContainer(t)
	defer cleanup()

	// Run vault tests
	logger := logging.NewVaultLogger(log.Debug)

	// connect with a connection url
	conf["connection_url"] = connURL
	b, err := NewMSSQLBackend(conf, logger)
	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	defer func() {
		mssql := b.(*MSSQLBackend)
		_, err := mssql.client.Exec("DROP TABLE " + mssql.dbTable)
		if err != nil {
			t.Fatalf("Failed to drop table: %v", err)
		}
	}()

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}

func TestMSSQLBackendWithoutConnectionURL(t *testing.T) {
	cleanup, _, connPort := prepareMSSQLTestContainer(t)
	defer cleanup()

	// Run vault tests
	logger := logging.NewVaultLogger(log.Debug)

	// connect without a connection url
	conf["port"] = connPort
	delete(conf, "connection_url")
	b, err := NewMSSQLBackend(conf, logger)
	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	defer func() {
		mssql := b.(*MSSQLBackend)
		_, err := mssql.client.Exec("DROP TABLE " + mssql.dbTable)
		if err != nil {
			t.Fatalf("Failed to drop table: %v", err)
		}
	}()

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}
