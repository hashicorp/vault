package postgresql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"github.com/ory-am/dockertest"
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
			testAccStepRole(t),
			testAccStepReadCreds(t, b, "web", connURL),
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
			testAccStepRole(t),
			testAccStepReadRole(t, "web", testRole),
			testAccStepDeleteRole(t, "web"),
			testAccStepReadRole(t, "web", ""),
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

func testAccStepRole(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/web",
		Data: map[string]interface{}{
			"sql": testRole,
		},
	}
}

func testAccStepDeleteRole(t *testing.T, n string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "roles/" + n,
	}
}

func testAccStepReadCreds(t *testing.T, b logical.Backend, name string, connURL string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "creds/" + name,
		Check: func(resp *logical.Response) error {
			var d struct {
				Username string `mapstructure:"username"`
				Password string `mapstructure:"password"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			log.Printf("[WARN] Generated credentials: %v", d)

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
				stmt, err := db.Prepare(fmt.Sprintf(
					"SELECT DISTINCT schemaname FROM pg_tables WHERE has_table_privilege('%s', 'information_schema.role_column_grants', 'select');",
					d.Username))
				if err != nil {
					return -1
				}
				defer stmt.Close()

				rows, err := stmt.Query()
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

			userRows := returnedRows()
			if userRows != 2 {
				t.Fatalf("did not get expected number of rows, got %d", userRows)
			}

			resp, err = b.HandleRequest(&logical.Request{
				Operation: logical.RevokeOperation,
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

			userRows = returnedRows()
			// User shouldn't exist so returnedRows() should encounter an error and exit with -1
			if userRows != -1 {
				t.Fatalf("did not get expected number of rows, got %d", userRows)
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
