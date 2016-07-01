package mongodb

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/mitchellh/mapstructure"
	"strings"
)

func TestBackend_basic(t *testing.T) {
	b, _ := Factory(logical.TestBackendConfig())

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck:       func() { testAccPreCheck(t) },
		Backend:        b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepRole(t),
			testAccStepReadCreds(t, "web"),
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
			testAccStepReadRole(t, "web", testDb, testMongoDBRoles),
			testAccStepDeleteRole(t, "web"),
			testAccStepReadRole(t, "web", "", ""),
		},
	})
}

func TestBackend_leaseWriteRead(t *testing.T) {
	b := Backend()

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck:       func() { testAccPreCheck(t) },
		Backend:        b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepWriteLease(t),
			testAccStepReadLease(t),
		},
	})

}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("MONGODB_URI"); v == "" {
		t.Fatal("MONGODB_URI must be set for acceptance tests")
	}
}

func testAccStepConfig(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/connection",
		Data: map[string]interface{}{
			"uri": os.Getenv("MONGODB_URI"),
		},
	}
}

func testAccStepRole(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/web",
		Data: map[string]interface{}{
			"db": testDb,
			"roles": testMongoDBRoles,
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
				DB       string `mapstructure:"db"`
				Username string `mapstructure:"username"`
				Password string `mapstructure:"password"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			if d.DB == "" {
				return fmt.Errorf("bad: %#v", resp)
			}
			if d.Username == "" {
				return fmt.Errorf("bad: %#v", resp)
			}
			if !strings.HasPrefix(d.Username, "vault-root-") {
				return fmt.Errorf("bad: %#v", resp)
			}
			if d.Password == "" {
				return fmt.Errorf("bad: %#v", resp)
			}

			log.Printf("[WARN] Generated credentials: %v", d)

			return nil
		},
	}
}

func testAccStepReadRole(t *testing.T, name, db, mongoDBRoles string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "roles/" + name,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				if db == "" && mongoDBRoles == "" {
					return nil
				}

				return fmt.Errorf("bad: %#v", resp)
			}

			var d struct {
				DB           string `mapstructure:"db"`
				MongoDBRoles string `mapstructure:"roles"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			if d.DB != db {
				return fmt.Errorf("bad: %#v", resp)
			}
			if d.MongoDBRoles != mongoDBRoles {
				return fmt.Errorf("bad: %#v", resp)
			}

			return nil
		},
	}
}

func testAccStepWriteLease(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/lease",
		Data: map[string]interface{}{
			"ttl":     "1h5m",
			"max_ttl": "24h",
		},
	}
}

func testAccStepReadLease(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "config/lease",
		Check: func(resp *logical.Response) error {
			if resp.Data["ttl"] != "1h5m0s" || resp.Data["max_ttl"] != "24h0m0s" {
				return fmt.Errorf("bad: %#v", resp)
			}

			return nil
		},
	}
}

const testDb = "foo"
const testMongoDBRoles = `["readWrite",{"db":"bar","role":"read"}]`
