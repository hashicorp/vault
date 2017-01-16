package postgresql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	dockertest "gopkg.in/ory-am/dockertest.v2"
)

var (
	testImagePull sync.Once
)

func prepareTestContainer(t *testing.T, s logical.Storage, b logical.Backend) (cid dockertest.ContainerID, retURL string) {
	if os.Getenv("PG_URL") != "" {
		return "", os.Getenv("PG_URL")
	}

	// Without this the checks for whether the container has started seem to
	// never actually pass. There's really no reason to expose the test
	// containers, so don't.
	dockertest.BindDockerToLocalhost = "yep"

	testImagePull.Do(func() {
		dockertest.Pull("postgres")
	})

	cid, connErr := dockertest.ConnectToPostgreSQL(60, 500*time.Millisecond, func(connURL string) bool {
		// This will cause a validation to run
		resp, err := b.HandleRequest(&logical.Request{
			Storage:   s,
			Operation: logical.UpdateOperation,
			Path:      "config/connection",
			Data: map[string]interface{}{
				"connection_url": connURL,
			},
		})
		if err != nil || (resp != nil && resp.IsError()) {
			// It's likely not up and running yet, so return false and try again
			return false
		}
		if resp == nil {
			t.Fatal("expected warning")
		}

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

func TestBackend_config_connection(t *testing.T) {
	var resp *logical.Response
	var err error
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	configData := map[string]interface{}{
		"connection_url":       "sample_connection_url",
		"value":                "",
		"max_open_connections": 9,
		"max_idle_connections": 7,
		"verify_connection":    false,
	}

	configReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/connection",
		Storage:   config.StorageView,
		Data:      configData,
	}
	resp, err = b.HandleRequest(configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	configReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	delete(configData, "verify_connection")
	if !reflect.DeepEqual(configData, resp.Data) {
		t.Fatalf("bad: expected:%#v\nactual:%#v\n", configData, resp.Data)
	}
}

func TestBackend_basic(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	cid, connURL := prepareTestContainer(t, config.StorageView, b)
	if cid != "" {
		defer cleanupTestContainer(t, cid)
	}
	connData := map[string]interface{}{
		"connection_url": connURL,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, connData, false),
			testAccStepCreateRole(t, "web", testRole, false),
			testAccStepReadCreds(t, b, config.StorageView, "web", connURL),
		},
	})
}

func TestBackend_roleCrud(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	cid, connURL := prepareTestContainer(t, config.StorageView, b)
	if cid != "" {
		defer cleanupTestContainer(t, cid)
	}
	connData := map[string]interface{}{
		"connection_url": connURL,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, connData, false),
			testAccStepCreateRole(t, "web", testRole, false),
			testAccStepReadRole(t, "web", testRole),
			testAccStepDeleteRole(t, "web"),
			testAccStepReadRole(t, "web", ""),
		},
	})
}

func TestBackend_BlockStatements(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	cid, connURL := prepareTestContainer(t, config.StorageView, b)
	if cid != "" {
		defer cleanupTestContainer(t, cid)
	}
	connData := map[string]interface{}{
		"connection_url": connURL,
	}

	jsonBlockStatement, err := json.Marshal(testBlockStatementRoleSlice)
	if err != nil {
		t.Fatal(err)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, connData, false),
			// This will also validate the query
			testAccStepCreateRole(t, "web-block", testBlockStatementRole, true),
			testAccStepCreateRole(t, "web-block", string(jsonBlockStatement), false),
		},
	})
}

func TestBackend_roleReadOnly(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	cid, connURL := prepareTestContainer(t, config.StorageView, b)
	if cid != "" {
		defer cleanupTestContainer(t, cid)
	}
	connData := map[string]interface{}{
		"connection_url": connURL,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, connData, false),
			testAccStepCreateRole(t, "web", testRole, false),
			testAccStepCreateRole(t, "web-readonly", testReadOnlyRole, false),
			testAccStepReadRole(t, "web-readonly", testReadOnlyRole),
			testAccStepCreateTable(t, b, config.StorageView, "web", connURL),
			testAccStepReadCreds(t, b, config.StorageView, "web-readonly", connURL),
			testAccStepDropTable(t, b, config.StorageView, "web", connURL),
			testAccStepDeleteRole(t, "web-readonly"),
			testAccStepDeleteRole(t, "web"),
			testAccStepReadRole(t, "web-readonly", ""),
		},
	})
}

func TestBackend_roleReadOnly_revocationSQL(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	cid, connURL := prepareTestContainer(t, config.StorageView, b)
	if cid != "" {
		defer cleanupTestContainer(t, cid)
	}
	connData := map[string]interface{}{
		"connection_url": connURL,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, connData, false),
			testAccStepCreateRoleWithRevocationSQL(t, "web", testRole, defaultRevocationSQL, false),
			testAccStepCreateRoleWithRevocationSQL(t, "web-readonly", testReadOnlyRole, defaultRevocationSQL, false),
			testAccStepReadRole(t, "web-readonly", testReadOnlyRole),
			testAccStepCreateTable(t, b, config.StorageView, "web", connURL),
			testAccStepReadCreds(t, b, config.StorageView, "web-readonly", connURL),
			testAccStepDropTable(t, b, config.StorageView, "web", connURL),
			testAccStepDeleteRole(t, "web-readonly"),
			testAccStepDeleteRole(t, "web"),
			testAccStepReadRole(t, "web-readonly", ""),
		},
	})
}

func testAccStepConfig(t *testing.T, d map[string]interface{}, expectError bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/connection",
		Data:      d,
		ErrorOk:   true,
		Check: func(resp *logical.Response) error {
			if expectError {
				if resp.Data == nil {
					return fmt.Errorf("data is nil")
				}
				var e struct {
					Error string `mapstructure:"error"`
				}
				if err := mapstructure.Decode(resp.Data, &e); err != nil {
					return err
				}
				if len(e.Error) == 0 {
					return fmt.Errorf("expected error, but write succeeded.")
				}
				return nil
			} else if resp != nil && resp.IsError() {
				return fmt.Errorf("got an error response: %v", resp.Error())
			}
			return nil
		},
	}
}

func testAccStepCreateRole(t *testing.T, name string, sql string, expectFail bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      path.Join("roles", name),
		Data: map[string]interface{}{
			"sql": sql,
		},
		ErrorOk: expectFail,
	}
}

func testAccStepCreateRoleWithRevocationSQL(t *testing.T, name, sql, revocationSQL string, expectFail bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      path.Join("roles", name),
		Data: map[string]interface{}{
			"sql":            sql,
			"revocation_sql": revocationSQL,
		},
		ErrorOk: expectFail,
	}
}

func testAccStepDeleteRole(t *testing.T, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      path.Join("roles", name),
	}
}

func testAccStepReadCreds(t *testing.T, b logical.Backend, s logical.Storage, name string, connURL string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      path.Join("creds", name),
		Check: func(resp *logical.Response) error {
			var d struct {
				Username string `mapstructure:"username"`
				Password string `mapstructure:"password"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			log.Printf("[TRACE] Generated credentials: %v", d)
			conn, err := pq.ParseURL(connURL)

			if err != nil {
				t.Fatal(err)
			}

			conn += " timezone=utc"

			db, err := sql.Open("postgres", conn)
			if err != nil {
				t.Fatal(err)
			}

			returnedRows := func() int {
				stmt, err := db.Prepare("SELECT DISTINCT schemaname FROM pg_tables WHERE has_table_privilege($1, 'information_schema.role_column_grants', 'select');")
				if err != nil {
					return -1
				}
				defer stmt.Close()

				rows, err := stmt.Query(d.Username)
				if err != nil {
					return -1
				}
				defer rows.Close()

				i := 0
				for rows.Next() {
					i++
				}
				return i
			}

			// minNumPermissions is the minimum number of permissions that will always be present.
			const minNumPermissions = 2

			userRows := returnedRows()
			if userRows < minNumPermissions {
				t.Fatalf("did not get expected number of rows, got %d", userRows)
			}

			resp, err = b.HandleRequest(&logical.Request{
				Operation: logical.RevokeOperation,
				Storage:   s,
				Secret: &logical.Secret{
					InternalData: map[string]interface{}{
						"secret_type": "creds",
						"username":    d.Username,
						"role":        name,
					},
				},
			})
			if err != nil {
				return err
			}
			if resp != nil {
				if resp.IsError() {
					return fmt.Errorf("Error on resp: %#v", *resp)
				}
			}

			userRows = returnedRows()
			// User shouldn't exist so returnedRows() should encounter an error and exit with -1
			if userRows != -1 {
				t.Fatalf("did not get expected number of rows, got %d", userRows)
			}

			return nil
		},
	}
}

func testAccStepCreateTable(t *testing.T, b logical.Backend, s logical.Storage, name string, connURL string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      path.Join("creds", name),
		Check: func(resp *logical.Response) error {
			var d struct {
				Username string `mapstructure:"username"`
				Password string `mapstructure:"password"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			log.Printf("[TRACE] Generated credentials: %v", d)
			conn, err := pq.ParseURL(connURL)

			if err != nil {
				t.Fatal(err)
			}

			conn += " timezone=utc"

			db, err := sql.Open("postgres", conn)
			if err != nil {
				t.Fatal(err)
			}

			_, err = db.Exec("CREATE TABLE test (id SERIAL PRIMARY KEY);")
			if err != nil {
				t.Fatal(err)
			}

			resp, err = b.HandleRequest(&logical.Request{
				Operation: logical.RevokeOperation,
				Storage:   s,
				Secret: &logical.Secret{
					InternalData: map[string]interface{}{
						"secret_type": "creds",
						"username":    d.Username,
					},
				},
			})
			if err != nil {
				return err
			}
			if resp != nil {
				if resp.IsError() {
					return fmt.Errorf("Error on resp: %#v", *resp)
				}
			}

			return nil
		},
	}
}

func testAccStepDropTable(t *testing.T, b logical.Backend, s logical.Storage, name string, connURL string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      path.Join("creds", name),
		Check: func(resp *logical.Response) error {
			var d struct {
				Username string `mapstructure:"username"`
				Password string `mapstructure:"password"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			log.Printf("[TRACE] Generated credentials: %v", d)
			conn, err := pq.ParseURL(connURL)

			if err != nil {
				t.Fatal(err)
			}

			conn += " timezone=utc"

			db, err := sql.Open("postgres", conn)
			if err != nil {
				t.Fatal(err)
			}

			_, err = db.Exec("DROP TABLE test;")
			if err != nil {
				t.Fatal(err)
			}

			resp, err = b.HandleRequest(&logical.Request{
				Operation: logical.RevokeOperation,
				Storage:   s,
				Secret: &logical.Secret{
					InternalData: map[string]interface{}{
						"secret_type": "creds",
						"username":    d.Username,
					},
				},
			})
			if err != nil {
				return err
			}
			if resp != nil {
				if resp.IsError() {
					return fmt.Errorf("Error on resp: %#v", *resp)
				}
			}

			return nil
		},
	}
}

func testAccStepReadRole(t *testing.T, name string, sql string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "roles/" + name,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				if sql == "" {
					return nil
				}

				return fmt.Errorf("bad: %#v", resp)
			}

			var d struct {
				SQL string `mapstructure:"sql"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			if d.SQL != sql {
				return fmt.Errorf("bad: %#v", resp)
			}

			return nil
		},
	}
}

const testRole = `
CREATE ROLE "{{name}}" WITH
  LOGIN
  PASSWORD '{{password}}'
  VALID UNTIL '{{expiration}}';
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "{{name}}";
`

const testReadOnlyRole = `
CREATE ROLE "{{name}}" WITH
  LOGIN
  PASSWORD '{{password}}'
  VALID UNTIL '{{expiration}}';
GRANT SELECT ON ALL TABLES IN SCHEMA public TO "{{name}}";
GRANT SELECT ON ALL SEQUENCES IN SCHEMA public TO "{{name}}";
`

const testBlockStatementRole = `
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

var testBlockStatementRoleSlice = []string{
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

const defaultRevocationSQL = `
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM {{name}};
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM {{name}};
REVOKE USAGE ON SCHEMA public FROM {{name}};

DROP ROLE IF EXISTS {{name}};
`
