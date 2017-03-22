package dbs

import (
	"database/sql"
	"os"
	"sync"
	"testing"
	"time"

	dockertest "gopkg.in/ory-am/dockertest.v2"
)

var (
	testMySQLImagePull sync.Once
)

func prepareMySQLTestContainer(t *testing.T) (cid dockertest.ContainerID, retURL string) {
	if os.Getenv("MYSQL_URL") != "" {
		return "", os.Getenv("MYSQL_URL")
	}

	// Without this the checks for whether the container has started seem to
	// never actually pass. There's really no reason to expose the test
	// containers, so don't.
	dockertest.BindDockerToLocalhost = "yep"

	testImagePull.Do(func() {
		dockertest.Pull("mysql")
	})

	cid, connErr := dockertest.ConnectToMySQL(60, 500*time.Millisecond, func(connURL string) bool {
		// This will cause a validation to run
		connProducer := &sqlConnectionProducer{}
		connProducer.ConnectionURL = connURL
		connProducer.config = &DatabaseConfig{
			DatabaseType: mySQLTypeName,
		}

		conn, err := connProducer.connection()
		if err != nil {
			return false
		}
		if err := conn.(*sql.DB).Ping(); err != nil {
			return false
		}

		connProducer.Close()

		retURL = connURL
		return true
	})

	if connErr != nil {
		t.Fatalf("could not connect to database: %v", connErr)
	}

	return
}

func TestMySQL_Initialize(t *testing.T) {
	cid, connURL := prepareMySQLTestContainer(t)
	if cid != "" {
		defer cleanupTestContainer(t, cid)
	}

	conf := &DatabaseConfig{
		DatabaseType: mySQLTypeName,
		ConnectionDetails: map[string]interface{}{
			"connection_url": connURL,
		},
	}

	dbRaw, err := BuiltinFactory(conf, nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Deconsturct the middleware chain to get the underlying mysql object
	dbMetrics := dbRaw.(*databaseMetricsMiddleware)
	db := dbMetrics.next.(*MySQL)

	err = dbRaw.Initialize(conf.ConnectionDetails)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	connProducer := db.ConnectionProducer.(*sqlConnectionProducer)
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

func TestMySQL_CreateUser(t *testing.T) {
	cid, connURL := prepareMySQLTestContainer(t)
	if cid != "" {
		defer cleanupTestContainer(t, cid)
	}

	conf := &DatabaseConfig{
		DatabaseType: mySQLTypeName,
		ConnectionDetails: map[string]interface{}{
			"connection_url": connURL,
		},
	}

	db, err := BuiltinFactory(conf, nil)
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
		CreationStatements: testMySQLRoleWildCard,
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
	statements.CreationStatements = testMySQLRoleHost
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
	/*	statements.CreationStatements = testBlockStatementRole
		err = db.CreateUser(statements, username, password, expiration)
		if err != nil {
			t.Fatalf("err: %s", err)
		}*/
}

func TestMySQL_RenewUser(t *testing.T) {
	cid, connURL := prepareMySQLTestContainer(t)
	if cid != "" {
		defer cleanupTestContainer(t, cid)
	}

	conf := &DatabaseConfig{
		DatabaseType: mySQLTypeName,
		ConnectionDetails: map[string]interface{}{
			"connection_url": connURL,
		},
	}

	db, err := BuiltinFactory(conf, nil)
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
		CreationStatements: testMySQLRoleWildCard,
	}

	err = db.CreateUser(statements, username, password, expiration)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expiration, err = db.GenerateExpiration(time.Minute)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	err = db.RenewUser(statements, username, expiration)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestMySQL_RevokeUser(t *testing.T) {
	cid, connURL := prepareMySQLTestContainer(t)
	if cid != "" {
		defer cleanupTestContainer(t, cid)
	}

	conf := &DatabaseConfig{
		DatabaseType: mySQLTypeName,
		ConnectionDetails: map[string]interface{}{
			"connection_url": connURL,
		},
	}

	db, err := BuiltinFactory(conf, nil)
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
		CreationStatements: testMySQLRoleWildCard,
	}

	err = db.CreateUser(statements, username, password, expiration)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expiration, err = db.GenerateExpiration(time.Minute)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Test default revoke statememts
	err = db.RevokeUser(statements, username)
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

	statements.CreationStatements = testMySQLRoleHost
	err = db.CreateUser(statements, username, password, expiration)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Test custom revoke statements
	statements.RevocationStatements = testMySQLRevocationSQL
	err = db.RevokeUser(statements, username)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

}

const testMySQLRoleWildCard = `
CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';
GRANT SELECT ON *.* TO '{{name}}'@'%';
`
const testMySQLRoleHost = `
CREATE USER '{{name}}'@'10.1.1.2' IDENTIFIED BY '{{password}}';
GRANT SELECT ON *.* TO '{{name}}'@'10.1.1.2';
`
const testMySQLRevocationSQL = `
REVOKE ALL PRIVILEGES, GRANT OPTION FROM '{{name}}'@'10.1.1.2'; 
DROP USER '{{name}}'@'10.1.1.2';
`
