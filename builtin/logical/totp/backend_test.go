package totp

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/mitchellh/mapstructure"
	otplib "github.com/pquerna/otp"
	totplib "github.com/pquerna/otp/totp"
)

var (
	masterKey string
)

func TestBackend_basic(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
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
