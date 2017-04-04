package dbs

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	log "github.com/mgutz/logxi/v1"
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
			t.Fatalf("Failed to cleanup local DynamoDB: %s", err)
		}
	}

	retURL = fmt.Sprintf("sqlserver://sa:yourStrong(!)Password@localhost:%s", resource.GetPort("1433/tcp"))

	// exponential backoff-retry, because the mssql container may not be able to accept connections yet
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

	conf := &DatabaseConfig{
		DatabaseType: msSQLTypeName,
		ConnectionDetails: map[string]interface{}{
			"connection_url": connURL,
		},
	}

	dbRaw, err := BuiltinFactory(conf, nil, &log.NullLogger{})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Deconsturct the middleware chain to get the underlying mssql object
	dbTracer := dbRaw.(*databaseTracingMiddleware)
	dbMetrics := dbTracer.next.(*databaseMetricsMiddleware)
	db := dbMetrics.next.(*MSSQL)
	connProducer := db.ConnectionProducer.(*sqlConnectionProducer)

	err = dbRaw.Initialize(conf.ConnectionDetails)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !connProducer.initalized {
		t.Fatal("Database should be initalized")
	}

	err = dbRaw.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if connProducer.db != nil {
		t.Fatal("db object should be nil")
	}
}

func TestMSSQL_CreateUser(t *testing.T) {
	cleanup, connURL := prepareMSSQLTestContainer(t)
	defer cleanup()

	conf := &DatabaseConfig{
		DatabaseType: msSQLTypeName,
		ConnectionDetails: map[string]interface{}{
			"connection_url": connURL,
		},
	}

	db, err := BuiltinFactory(conf, nil, &log.NullLogger{})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	err = db.Initialize(conf.ConnectionDetails)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	username, err := db.GenerateUsername("test")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	password, err := db.GeneratePassword()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expiration, err := db.GenerateExpiration(time.Minute)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Test with no configured Creation Statememt
	err = db.CreateUser(Statements{}, username, password, expiration)
	if err == nil {
		t.Fatal("Expected error when no creation statement is provided")
	}

	statements := Statements{
		CreationStatements: testMSSQLRole,
	}

	err = db.CreateUser(statements, username, password, expiration)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	username, err = db.GenerateUsername("test")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	password, err = db.GeneratePassword()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expiration, err = db.GenerateExpiration(time.Minute)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestMSSQL_RevokeUser(t *testing.T) {
	cleanup, connURL := prepareMSSQLTestContainer(t)
	defer cleanup()

	conf := &DatabaseConfig{
		DatabaseType: msSQLTypeName,
		ConnectionDetails: map[string]interface{}{
			"connection_url": connURL,
		},
	}

	db, err := BuiltinFactory(conf, nil, &log.NullLogger{})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	err = db.Initialize(conf.ConnectionDetails)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	username, err := db.GenerateUsername("test")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	password, err := db.GeneratePassword()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expiration, err := db.GenerateExpiration(time.Minute)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	statements := Statements{
		CreationStatements: testMSSQLRole,
	}

	err = db.CreateUser(statements, username, password, expiration)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Test default revoke statememts
	err = db.RevokeUser(statements, username)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

const testMSSQLRole = `
CREATE LOGIN [{{name}}] WITH PASSWORD = '{{password}}';
CREATE USER [{{name}}] FOR LOGIN [{{name}}];
GRANT SELECT, INSERT, UPDATE, DELETE ON SCHEMA::dbo TO [{{name}}];`
