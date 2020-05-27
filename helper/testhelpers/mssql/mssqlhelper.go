package mssqlhelper

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/vault/helper/testhelpers/docker"
	"github.com/ory/dockertest"
)

func PrepareMSSQLTestContainer(t *testing.T) (cleanup func(), retURL string) {
	if os.Getenv("MSSQL_URL") != "" {
		return func() {}, os.Getenv("MSSQL_URL")
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	resource, err := pool.Run("mcr.microsoft.com/mssql/server", "2017-latest-ubuntu", []string{"ACCEPT_EULA=Y", "SA_PASSWORD=yourStrong(!)Password"})
	if err != nil {
		t.Fatalf("Could not start local MSSQL docker container: %s", err)
	}

	cleanup = func() {
		docker.CleanupResource(t, pool, resource)
	}

	retURL = fmt.Sprintf("sqlserver://sa:yourStrong(!)Password@127.0.0.1:%s", resource.GetPort("1433/tcp"))

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		var err error
		var db *sql.DB
		db, err = sql.Open("mssql", retURL)
		if err != nil {
			return err
		}
		defer db.Close()
		return db.Ping()
	}); err != nil {
		cleanup()
		t.Fatalf("Could not connect to MSSQL docker container: %s", err)
	}

	return
}
