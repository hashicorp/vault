package mongodb

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"testing"

	logicaltest "github.com/hashicorp/vault/helper/testhelpers/logical"
	"github.com/hashicorp/vault/helper/testhelpers/mongodb"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

var (
	testImagePull sync.Once
)

func TestBackend_config_connection(t *testing.T) {
	var resp *logical.Response
	var err error
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	configData := map[string]interface{}{
		"uri":               "sample_connection_uri",
		"verify_connection": false,
	}

	configReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/connection",
		Storage:   config.StorageView,
		Data:      configData,
	}
	resp, err = b.HandleRequest(context.Background(), configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	configReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}
}

func TestBackend_basic(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, connURI := mongodb.PrepareTestContainer(t, "latest")
	defer cleanup()
	connData := map[string]interface{}{
		"uri": connURI,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(connData, false),
			testAccStepRole(),
			testAccStepReadCreds("web"),
		},
	})
}

func TestBackend_roleCrud(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, connURI := mongodb.PrepareTestContainer(t, "latest")
	defer cleanup()
	connData := map[string]interface{}{
		"uri": connURI,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(connData, false),
			testAccStepRole(),
			testAccStepReadRole("web", testDb, testMongoDBRoles),
			testAccStepDeleteRole("web"),
			testAccStepReadRole("web", "", ""),
		},
	})
}

func TestBackend_leaseWriteRead(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, connURI := mongodb.PrepareTestContainer(t, "latest")
	defer cleanup()
	connData := map[string]interface{}{
		"uri": connURI,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(connData, false),
			testAccStepWriteLease(),
			testAccStepReadLease(),
		},
	})

}

func testAccStepConfig(d map[string]interface{}, expectError bool) logicaltest.TestStep {
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
					return fmt.Errorf("expected error, but write succeeded")
				}
				return nil
			} else if resp != nil && resp.IsError() {
				return fmt.Errorf("got an error response: %v", resp.Error())
			}
			return nil
		},
	}
}

func testAccStepRole() logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/web",
		Data: map[string]interface{}{
			"db":    testDb,
			"roles": testMongoDBRoles,
		},
	}
}

func testAccStepDeleteRole(n string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "roles/" + n,
	}
}

func testAccStepReadCreds(name string) logicaltest.TestStep {
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

func testAccStepReadRole(name, db, mongoDBRoles string) logicaltest.TestStep {
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

func testAccStepWriteLease() logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/lease",
		Data: map[string]interface{}{
			"ttl":     "1h5m",
			"max_ttl": "24h",
		},
	}
}

func testAccStepReadLease() logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "config/lease",
		Check: func(resp *logical.Response) error {
			if resp.Data["ttl"].(float64) != 3900 || resp.Data["max_ttl"].(float64) != 86400 {
				return fmt.Errorf("bad: %#v", resp)
			}

			return nil
		},
	}
}

const testDb = "foo"
const testMongoDBRoles = `["readWrite",{"role":"read","db":"bar"}]`
