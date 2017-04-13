package mssql

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/plugins/helper/database/connutil"
	dockertest "gopkg.in/ory-am/dockertest.v3"
)

var (
	testMSQLImagePull sync.Once
)

func prepareMSSQLTestContainer(t *testing.T) (cleanup func(), retURL string) {
	if os.Getenv("MSSQL_URL") != "" {
		return func() {}, os.Getenv("MSSQL_URL")
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	resource, err := pool.Run("microsoft/mssql-server-linux", "latest", []string{"ACCEPT_EULA=Y", "SA_PASSWORD=yourStrong(!)Password"})
	if err != nil {
		t.Fatalf("Could not start local MSSQL docker container: %s", err)
	}

	cleanup = func() {
		err := pool.Purge(resource)
		if err != nil {
			t.Fatalf("Failed to cleanup local container: %s", err)
		}
	}

	retURL = fmt.Sprintf("sqlserver://sa:yourStrong(!)Password@localhost:%s", resource.GetPort("1433/tcp"))

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		var err error
		var db *sql.DB
		db, err = sql.Open("mssql", retURL)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		t.Fatalf("Could not connect to MSSQL docker container: %s", err)
	}

	return
}

func TestMSSQL_Initialize(t *testing.T) {
	cleanup, connURL := prepareMSSQLTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := New()

	err := db.Initialize(connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	connProducer := db.ConnectionProducer.(*connutil.SQLConnectionProducer)
	if !connProducer.Initialized {
		t.Fatal("Database should be initalized")
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestMSSQL_CreateUser(t *testing.T) {
	cleanup, connURL := prepareMSSQLTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := New()
	err := db.Initialize(connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Test with no configured Creation Statememt
	_, _, err = db.CreateUser(dbplugin.Statements{}, "test", time.Now().Add(time.Minute))
	if err == nil {
		t.Fatal("Expected error when no creation statement is provided")
	}

	statements := dbplugin.Statements{
		CreationStatements: testMSSQLRole,
	}

	username, password, err := db.CreateUser(statements, "test", time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
}

func TestMSSQL_RevokeUser(t *testing.T) {
	cleanup, connURL := prepareMSSQLTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := New()
	err := db.Initialize(connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		CreationStatements: testMSSQLRole,
	}

	username, password, err := db.CreateUser(statements, "test", time.Now().Add(2*time.Second))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test default revoke statememts
	err = db.RevokeUser(statements, username)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err == nil {
		t.Fatal("Credentials were not revoked")
	}
}

func testCredsExist(t testing.TB, connURL, username, password string) error {
	// Log in with the new creds
	connURL = strings.Replace(connURL, "sa:yourStrong(!)Password", fmt.Sprintf("%s:%s", username, password), 1)
	db, err := sql.Open("mssql", connURL)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}

const testMSSQLRole = `
CREATE LOGIN [{{name}}] WITH PASSWORD = '{{password}}';
CREATE USER [{{name}}] FOR LOGIN [{{name}}];
GRANT SELECT, INSERT, UPDATE, DELETE ON SCHEMA::dbo TO [{{name}}];`
