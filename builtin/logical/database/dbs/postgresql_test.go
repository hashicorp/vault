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
	testPostgresImagePull sync.Once
)

func preparePostgresTestContainer(t *testing.T) (cid dockertest.ContainerID, retURL string) {
	if os.Getenv("PG_URL") != "" {
		return "", os.Getenv("PG_URL")
	}

	// Without this the checks for whether the container has started seem to
	// never actually pass. There's really no reason to expose the test
	// containers, so don't.
	dockertest.BindDockerToLocalhost = "yep"

	testPostgresImagePull.Do(func() {
		dockertest.Pull("postgres")
	})

	cid, connErr := dockertest.ConnectToPostgreSQL(60, 500*time.Millisecond, func(connURL string) bool {
		// This will cause a validation to run
		connProducer := &sqlConnectionProducer{}
		connProducer.ConnectionURL = connURL
		connProducer.config = &DatabaseConfig{
			DatabaseType: postgreSQLTypeName,
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

func cleanupTestContainer(t *testing.T, cid dockertest.ContainerID) {
	err := cid.KillRemove()
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostgreSQL_Initialize(t *testing.T) {
	cid, connURL := preparePostgresTestContainer(t)
	if cid != "" {
		defer cleanupTestContainer(t, cid)
	}

	conf := &DatabaseConfig{
		DatabaseType: postgreSQLTypeName,
		ConnectionDetails: map[string]interface{}{
			"connection_url": connURL,
		},
	}

	dbRaw, err := BuiltinFactory(conf, nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Deconsturct the middleware chain to get the underlying postgres object
	dbMetrics := dbRaw.(*databaseMetricsMiddleware)
	db := dbMetrics.next.(*PostgreSQL)

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

func TestPostgreSQL_CreateUser(t *testing.T) {
	cid, connURL := preparePostgresTestContainer(t)
	if cid != "" {
		defer cleanupTestContainer(t, cid)
	}

	conf := &DatabaseConfig{
		DatabaseType: postgreSQLTypeName,
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
		CreationStatements: testPostgresRole,
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
	statements.CreationStatements = testPostgresReadOnlyRole
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

func TestPostgreSQL_RenewUser(t *testing.T) {
	cid, connURL := preparePostgresTestContainer(t)
	if cid != "" {
		defer cleanupTestContainer(t, cid)
	}

	conf := &DatabaseConfig{
		DatabaseType: postgreSQLTypeName,
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
		CreationStatements: testPostgresRole,
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

func TestPostgreSQL_RevokeUser(t *testing.T) {
	cid, connURL := preparePostgresTestContainer(t)
	if cid != "" {
		defer cleanupTestContainer(t, cid)
	}

	conf := &DatabaseConfig{
		DatabaseType: postgreSQLTypeName,
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
		CreationStatements: testPostgresRole,
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

	err = db.CreateUser(statements, username, password, expiration)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Test custom revoke statements
	statements.RevocationStatements = defaultPostgresRevocationSQL
	err = db.RevokeUser(statements, username)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

}

const testPostgresRole = `
CREATE ROLE "{{name}}" WITH
  LOGIN
  PASSWORD '{{password}}'
  VALID UNTIL '{{expiration}}';
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "{{name}}";
`

const testPostgresReadOnlyRole = `
CREATE ROLE "{{name}}" WITH
  LOGIN
  PASSWORD '{{password}}'
  VALID UNTIL '{{expiration}}';
GRANT SELECT ON ALL TABLES IN SCHEMA public TO "{{name}}";
GRANT SELECT ON ALL SEQUENCES IN SCHEMA public TO "{{name}}";
`

const testPostgresBlockStatementRole = `
DO $$
BEGIN
   IF NOT EXISTS (SELECT * FROM pg_catalog.pg_roles WHERE rolname='foo-role') THEN
      CREATE ROLE "foo-role";
      CREATE SCHEMA IF NOT EXISTS foo AUTHORIZATION "foo-role";
      ALTER ROLE "foo-role" SET search_path = foo;
      GRANT TEMPORARY ON DATABASE "postgres" TO "foo-role";
      GRANT ALL PRIVILEGES ON SCHEMA foo TO "foo-role";
      GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA foo TO "foo-role";
      GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA foo TO "foo-role";
      GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA foo TO "foo-role";
   END IF;
END
$$

CREATE ROLE "{{name}}" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}';
GRANT "foo-role" TO "{{name}}";
ALTER ROLE "{{name}}" SET search_path = foo;
GRANT CONNECT ON DATABASE "postgres" TO "{{name}}";
`

var testPostgresBlockStatementRoleSlice = []string{
	`
DO $$
BEGIN
   IF NOT EXISTS (SELECT * FROM pg_catalog.pg_roles WHERE rolname='foo-role') THEN
      CREATE ROLE "foo-role";
      CREATE SCHEMA IF NOT EXISTS foo AUTHORIZATION "foo-role";
      ALTER ROLE "foo-role" SET search_path = foo;
      GRANT TEMPORARY ON DATABASE "postgres" TO "foo-role";
      GRANT ALL PRIVILEGES ON SCHEMA foo TO "foo-role";
      GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA foo TO "foo-role";
      GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA foo TO "foo-role";
      GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA foo TO "foo-role";
   END IF;
END
$$
`,
	`CREATE ROLE "{{name}}" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}';`,
	`GRANT "foo-role" TO "{{name}}";`,
	`ALTER ROLE "{{name}}" SET search_path = foo;`,
	`GRANT CONNECT ON DATABASE "postgres" TO "{{name}}";`,
}

const defaultPostgresRevocationSQL = `
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM "{{name}}";
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM "{{name}}";
REVOKE USAGE ON SCHEMA public FROM "{{name}}";

DROP ROLE IF EXISTS "{{name}}";
`
