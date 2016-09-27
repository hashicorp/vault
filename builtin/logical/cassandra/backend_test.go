package cassandra

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/mitchellh/mapstructure"
)

func TestBackend_basic(t *testing.T) {
	b := Backend()

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck:       func() { testAccPreCheck(t) },
		Backend:        b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepRole(t),
			testAccStepReadCreds(t, "test"),
		},
	})
}

func TestBackend_roleCrud(t *testing.T) {
	b := Backend()

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck:       func() { testAccPreCheck(t) },
		Backend:        b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepRole(t),
			testAccStepReadRole(t, "test", testRole),
			testAccStepDeleteRole(t, "test"),
			testAccStepReadRole(t, "test", ""),
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("CASSANDRA_HOST"); v == "" {
		t.Fatal("CASSANDRA_HOST must be set for acceptance tests")
	}
}

func testAccStepConfig(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/connection",
		Data: map[string]interface{}{
			"hosts":    os.Getenv("CASSANDRA_HOST"),
			"username": "cassandra",
			"password": "cassandra",
		},
	}
}

func testAccStepRole(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/test",
		Data: map[string]interface{}{
			"creation_cql": testRole,
			"consistency":  "All",
		},
	}
}

func testAccStepDeleteRole(t *testing.T, n string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "roles/" + n,
	}
}

func testAccStepReadCreds(t *testing.T, name string) logicaltest.TestStep {
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

			return nil
		},
	}
}

func testAccStepReadRole(t *testing.T, name string, cql string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "roles/" + name,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				if cql == "" {
					return nil
				}

				return fmt.Errorf("response is nil")
			}

			var d struct {
				CreationCQL string `mapstructure:"creation_cql"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			if d.CreationCQL != cql {
				return fmt.Errorf("bad: %#v\n%#v\n%#v\n", resp, cql, d.CreationCQL)
			}

			return nil
		},
	}
}

const testRole = `CREATE USER '{{username}}' WITH PASSWORD '{{password}}' NOSUPERUSER;
GRANT ALL PERMISSIONS ON ALL KEYSPACES TO {{username}};`
