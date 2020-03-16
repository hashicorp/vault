package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers/docker"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/helper/dbtxn"
	"github.com/lib/pq"
	"github.com/ory/dockertest"
)

func preparePostgresTestContainer(t *testing.T) (cleanup func(), retURL string) {
	if os.Getenv("PG_URL") != "" {
		return func() {}, os.Getenv("PG_URL")
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	resource, err := pool.Run("postgres", "latest", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=database"})
	if err != nil {
		t.Fatalf("Could not start local PostgreSQL docker container: %s", err)
	}

	cleanup = func() {
		docker.CleanupResource(t, pool, resource)
	}

	retURL = fmt.Sprintf("postgres://postgres:secret@localhost:%s/database?sslmode=disable", resource.GetPort("5432/tcp"))

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		var err error
		var db *sql.DB
		db, err = sql.Open("postgres", retURL)
		if err != nil {
			return err
		}
		defer db.Close()
		return db.Ping()
	}); err != nil {
		cleanup()
		t.Fatalf("Could not connect to PostgreSQL docker container: %s", err)
	}

	return
}

func TestPostgreSQL_Initialize(t *testing.T) {
	cleanup, connURL := preparePostgresTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url":       connURL,
		"max_open_connections": 5,
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !db.Initialized {
		t.Fatal("Database should be initialized")
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Test decoding a string value for max_open_connections
	connectionDetails = map[string]interface{}{
		"connection_url":       connURL,
		"max_open_connections": "5",
	}

	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

}

func TestPostgreSQL_CreateUser_missingArgs(t *testing.T) {
	db := new()

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	username, password, err := db.CreateUser(context.Background(), dbplugin.Statements{}, usernameConfig, time.Now().Add(time.Minute))
	if err == nil {
		t.Fatalf("expected err, got nil")
	}
	if username != "" {
		t.Fatalf("expected empty username, got [%s]", username)
	}
	if password != "" {
		t.Fatalf("expected empty password, got [%s]", password)
	}
}

func TestPostgreSQL_CreateUser(t *testing.T) {
	type testCase struct {
		createStmts          []string
		shouldTestCredsExist bool
	}

	tests := map[string]testCase{
		"admin name": {
			createStmts: []string{`
				CREATE ROLE "{{name}}" WITH
				  LOGIN
				  PASSWORD '{{password}}'
				  VALID UNTIL '{{expiration}}';
				GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "{{name}}";`,
			},
			shouldTestCredsExist: true,
		},
		"admin username": {
			createStmts: []string{`
				CREATE ROLE "{{username}}" WITH
				  LOGIN
				  PASSWORD '{{password}}'
				  VALID UNTIL '{{expiration}}';
				GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "{{username}}";`,
			},
			shouldTestCredsExist: true,
		},
		"read only name": {
			createStmts: []string{`
				CREATE ROLE "{{name}}" WITH
				  LOGIN
				  PASSWORD '{{password}}'
				  VALID UNTIL '{{expiration}}';
				GRANT SELECT ON ALL TABLES IN SCHEMA public TO "{{name}}";
				GRANT SELECT ON ALL SEQUENCES IN SCHEMA public TO "{{name}}";`,
			},
			shouldTestCredsExist: true,
		},
		"read only username": {
			createStmts: []string{`
				CREATE ROLE "{{username}}" WITH
				  LOGIN
				  PASSWORD '{{password}}'
				  VALID UNTIL '{{expiration}}';
				GRANT SELECT ON ALL TABLES IN SCHEMA public TO "{{username}}";
				GRANT SELECT ON ALL SEQUENCES IN SCHEMA public TO "{{username}}";`,
			},
			shouldTestCredsExist: true,
		},
		"reproduce https://github.com/hashicorp/vault/issues/6098": {
			createStmts: []string{
				// NOTE: "rolname" in the following line is not a typo.
				"DO $$ BEGIN IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname='my_role') THEN CREATE ROLE my_role; END IF; END $$",
			},
			// This test statement doesn't generate creds.
			shouldTestCredsExist: false,
		},
		"reproduce issue with template": {
			createStmts: []string{
				`DO $$ BEGIN IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname='my_role') THEN CREATE ROLE "{{username}}"; END IF; END $$`,
			},
			// This test statement doesn't generate creds.
			shouldTestCredsExist: false,
		},
	}

	// Shared test container for speed - there should not be any overlap between the tests
	cleanup, connURL := preparePostgresTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			usernameConfig := dbplugin.UsernameConfig{
				DisplayName: "test",
				RoleName:    "test",
			}

			statements := dbplugin.Statements{
				Creation: test.createStmts,
			}

			// Give a timeout just in case the test decides to be problematic
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			username, password, err := db.CreateUser(ctx, statements, usernameConfig, time.Now().Add(time.Minute))
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if !test.shouldTestCredsExist {
				// We're done here.
				return
			}

			if err = testCredsExist(t, connURL, username, password); err != nil {
				t.Fatalf("Could not connect with new credentials: %s", err)
			}

			// Ensure that the role doesn't expire immediately
			time.Sleep(2 * time.Second)

			if err = testCredsExist(t, connURL, username, password); err != nil {
				t.Fatalf("Could not connect with new credentials: %s", err)
			}
		})
	}
}

func TestPostgreSQL_RenewUser(t *testing.T) {
	type testCase struct {
		renewalStmts []string
	}

	tests := map[string]testCase{
		"empty renewal statements": {
			renewalStmts: nil,
		},
		"default renewal name": {
			renewalStmts: []string{defaultPostgresRenewSQL},
		},
		"default renewal username": {
			renewalStmts: []string{`
				ALTER ROLE "{{username}}" VALID UNTIL '{{expiration}}';`,
			},
		},
	}

	// Shared test container for speed - there should not be any overlap between the tests
	cleanup, connURL := preparePostgresTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := new()

	// Give a timeout just in case the test decides to be problematic
	initCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := db.Init(initCtx, connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			statements := dbplugin.Statements{
				Creation: []string{createAdminUser},
				Renewal:  test.renewalStmts,
			}

			usernameConfig := dbplugin.UsernameConfig{
				DisplayName: "test",
				RoleName:    "test",
			}

			// Give a timeout just in case the test decides to be problematic
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			username, password, err := db.CreateUser(ctx, statements, usernameConfig, time.Now().Add(2*time.Second))
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if err = testCredsExist(t, connURL, username, password); err != nil {
				t.Fatalf("Could not connect with new credentials: %s", err)
			}

			err = db.RenewUser(ctx, statements, username, time.Now().Add(time.Minute))
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			// Sleep longer than the initial expiration time
			time.Sleep(2 * time.Second)

			if err = testCredsExist(t, connURL, username, password); err != nil {
				t.Fatalf("Could not connect with new credentials: %s", err)
			}
		})
	}
}

func TestPostgreSQL_RotateRootCredentials(t *testing.T) {
	type testCase struct {
		statements []string
	}

	tests := map[string]testCase{
		"empty statements": {
			statements: nil,
		},
		"default name": {
			statements: []string{`
				ALTER ROLE "{{name}}" WITH PASSWORD '{{password}}';`,
			},
		},
		"default username": {
			statements: []string{defaultPostgresRotateRootCredentialsSQL},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cleanup, connURL := preparePostgresTestContainer(t)
			defer cleanup()

			connURL = strings.Replace(connURL, "postgres:secret", `{{username}}:{{password}}`, -1)

			connectionDetails := map[string]interface{}{
				"connection_url":       connURL,
				"max_open_connections": 5,
				"username":             "postgres",
				"password":             "secret",
			}

			db := new()

			connProducer := db.SQLConnectionProducer

			// Give a timeout just in case the test decides to be problematic
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			_, err := db.Init(ctx, connectionDetails, true)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if !connProducer.Initialized {
				t.Fatal("Database should be initialized")
			}

			newConf, err := db.RotateRootCredentials(ctx, test.statements)
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			if newConf["password"] == "secret" {
				t.Fatal("password was not updated")
			}

			err = db.Close()
			if err != nil {
				t.Fatalf("failed to close: %s", err)
			}
		})
	}
}

func TestPostgreSQL_RevokeUser(t *testing.T) {
	type testCase struct {
		revokeStmts []string
	}

	tests := map[string]testCase{
		"empty statements": {
			revokeStmts: nil,
		},
		"explicit default name": {
			revokeStmts: []string{defaultPostgresRevocationSQL},
		},
		"explicit default username": {
			revokeStmts: []string{`
				REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM "{{username}}";
				REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM "{{username}}";
				REVOKE USAGE ON SCHEMA public FROM "{{username}}";
				
				DROP ROLE IF EXISTS "{{username}}";`,
			},
		},
	}

	// Shared test container for speed - there should not be any overlap between the tests
	cleanup, connURL := preparePostgresTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := new()

	// Give a timeout just in case the test decides to be problematic
	initCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := db.Init(initCtx, connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			statements := dbplugin.Statements{
				Creation:   []string{createAdminUser},
				Revocation: test.revokeStmts,
			}

			usernameConfig := dbplugin.UsernameConfig{
				DisplayName: "test",
				RoleName:    "test",
			}

			username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(2*time.Second))
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if err = testCredsExist(t, connURL, username, password); err != nil {
				t.Fatalf("Could not connect with new credentials: %s", err)
			}

			// Test default revoke statements
			err = db.RevokeUser(context.Background(), statements, username)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if err := testCredsExist(t, connURL, username, password); err == nil {
				t.Fatal("Credentials were not revoked")
			}
		})
	}
}

func TestPostgreSQL_SetCredentials_missingArgs(t *testing.T) {
	type testCase struct {
		statements dbplugin.Statements
		userConfig dbplugin.StaticUserConfig
	}

	tests := map[string]testCase{
		"empty rotation statements": {
			statements: dbplugin.Statements{
				Rotation: nil,
			},
			userConfig: dbplugin.StaticUserConfig{
				Username: "testuser",
				Password: "password",
			},
		},
		"empty username": {
			statements: dbplugin.Statements{
				Rotation: []string{`
					ALTER ROLE "{{name}}" WITH PASSWORD '{{password}}';`,
				},
			},
			userConfig: dbplugin.StaticUserConfig{
				Username: "",
				Password: "password",
			},
		},
		"empty password": {
			statements: dbplugin.Statements{
				Rotation: []string{`
					ALTER ROLE "{{name}}" WITH PASSWORD '{{password}}';`,
				},
			},
			userConfig: dbplugin.StaticUserConfig{
				Username: "testuser",
				Password: "",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db := new()

			username, password, err := db.SetCredentials(context.Background(), test.statements, test.userConfig)
			if err == nil {
				t.Fatalf("expected err, got nil")
			}
			if username != "" {
				t.Fatalf("expected empty username, got [%s]", username)
			}
			if password != "" {
				t.Fatalf("expected empty password, got [%s]", password)
			}
		})
	}
}

func TestPostgresSQL_SetCredentials(t *testing.T) {
	type testCase struct {
		rotationStmts []string
	}

	tests := map[string]testCase{
		"name rotation": {
			rotationStmts: []string{`
				ALTER ROLE "{{name}}" WITH PASSWORD '{{password}}';`,
			},
		},
		"username rotation": {
			rotationStmts: []string{`
				ALTER ROLE "{{username}}" WITH PASSWORD '{{password}}';`,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Shared test container for speed - there should not be any overlap between the tests
			cleanup, connURL := preparePostgresTestContainer(t)
			defer cleanup()

			// create the database user
			dbUser := "vaultstatictest"
			initPassword := "password"
			createTestPGUser(t, connURL, dbUser, initPassword, testRoleStaticCreate)

			connectionDetails := map[string]interface{}{
				"connection_url": connURL,
			}

			db := new()

			// Give a timeout just in case the test decides to be problematic
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			_, err := db.Init(ctx, connectionDetails, true)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			statements := dbplugin.Statements{
				Rotation: test.rotationStmts,
			}

			password, err := db.GenerateCredentials(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			usernameConfig := dbplugin.StaticUserConfig{
				Username: dbUser,
				Password: password,
			}

			if err := testCredsExist(t, connURL, dbUser, initPassword); err != nil {
				t.Fatalf("Could not connect with initial credentials: %s", err)
			}

			username, password, err := db.SetCredentials(ctx, statements, usernameConfig)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if err := testCredsExist(t, connURL, username, password); err != nil {
				t.Fatalf("Could not connect with new credentials: %s", err)
			}

			if err := testCredsExist(t, connURL, username, initPassword); err == nil {
				t.Fatalf("Should not be able to connect with initial credentials")
			}
		})
	}
}

func testCredsExist(t testing.TB, connURL, username, password string) error {
	t.Helper()
	// Log in with the new creds
	connURL = strings.Replace(connURL, "postgres:secret", fmt.Sprintf("%s:%s", username, password), 1)
	db, err := sql.Open("postgres", connURL)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}

const createAdminUser = `
CREATE ROLE "{{name}}" WITH
  LOGIN
  PASSWORD '{{password}}'
  VALID UNTIL '{{expiration}}';
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "{{name}}";
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

const testRoleStaticCreate = `
CREATE ROLE "{{name}}" WITH
  LOGIN
  PASSWORD '{{password}}';
`

// This is a copy of a test helper method also found in
// builtin/logical/database/rotation_test.go , and should be moved into a shared
// helper file in the future.
func createTestPGUser(t *testing.T, connURL string, username, password, query string) {
	t.Helper()
	conn, err := pq.ParseURL(connURL)
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("postgres", conn)
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Start a transaction
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	m := map[string]string{
		"name":     username,
		"password": password,
	}
	if err := dbtxn.ExecuteTxQuery(ctx, tx, m, query); err != nil {
		t.Fatal(err)
	}
	// Commit the transaction
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}

func TestContainsMultilineStatement(t *testing.T) {
	type testCase struct {
		Input    string
		Expected bool
	}

	testCases := map[string]*testCase{
		"issue 6098 repro": {
			Input:    `DO $$ BEGIN IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname='my_role') THEN CREATE ROLE my_role; END IF; END $$`,
			Expected: true,
		},
		"multiline with template fields": {
			Input:    `DO $$ BEGIN IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname=\"{{name}}\") THEN CREATE ROLE {{name}}; END IF; END $$`,
			Expected: true,
		},
		"docs example": {
			Input: `CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}'; \
        GRANT SELECT ON ALL TABLES IN SCHEMA public TO \"{{name}}\";`,
			Expected: false,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			if containsMultilineStatement(tCase.Input) != tCase.Expected {
				t.Fatalf("%q should be %t for multiline input", tCase.Input, tCase.Expected)
			}
		})
	}
}

func TestExtractQuotedStrings(t *testing.T) {
	type testCase struct {
		Input    string
		Expected []string
	}

	testCases := map[string]*testCase{
		"no quotes": {
			Input:    `Five little monkeys jumping on the bed`,
			Expected: []string{},
		},
		"two of both quote types": {
			Input:    `"Five" little 'monkeys' "jumping on" the' 'bed`,
			Expected: []string{`"Five"`, `"jumping on"`, `'monkeys'`, `' '`},
		},
		"one single quote": {
			Input:    `Five little monkeys 'jumping on the bed`,
			Expected: []string{},
		},
		"empty string": {
			Input:    ``,
			Expected: []string{},
		},
		"templated field": {
			Input:    `DO $$ BEGIN IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname=\"{{name}}\") THEN CREATE ROLE {{name}}; END IF; END $$`,
			Expected: []string{`"{{name}}\"`},
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			results, err := extractQuotedStrings(tCase.Input)
			if err != nil {
				t.Fatal(err)
			}
			if len(results) != len(tCase.Expected) {
				t.Fatalf("%s isn't equal to %s", results, tCase.Expected)
			}
			for i := range results {
				if results[i] != tCase.Expected[i] {
					t.Fatalf(`expected %q but received %q`, tCase.Expected, results[i])
				}
			}
		})
	}
}
