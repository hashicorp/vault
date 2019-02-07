package mssql

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/mitchellh/mapstructure"
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
		"connection_string":    "sample_connection_string",
		"max_open_connections": 7,
		"verify_connection":    false,
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

	delete(configData, "verify_connection")
	delete(configData, "connection_string")
	if !reflect.DeepEqual(configData, resp.Data) {
		t.Fatalf("bad: expected:%#v\nactual:%#v\n", configData, resp.Data)
	}
}

func TestBackend_basic(t *testing.T) {
	if os.Getenv(logicaltest.TestEnvVar) == "" || os.Getenv("MSSQL_URL") == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", logicaltest.TestEnvVar))
	}
	connURL := os.Getenv("MSSQL_URL")

	b, _ := Factory(context.Background(), logical.TestBackendConfig())

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck:       testAccPreCheckFunc(t, connURL),
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, connURL),
			testAccStepRole(t),
			testAccStepReadCreds(t, "web"),
		},
	})
}

func TestBackend_roleCrud(t *testing.T) {
	if os.Getenv(logicaltest.TestEnvVar) == "" || os.Getenv("MSSQL_URL") == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", logicaltest.TestEnvVar))
	}
	connURL := os.Getenv("MSSQL_URL")

	b := Backend()

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck:       testAccPreCheckFunc(t, connURL),
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, connURL),
			testAccStepRole(t),
			testAccStepReadRole(t, "web", testRoleSQL),
			testAccStepDeleteRole(t, "web"),
			testAccStepReadRole(t, "web", ""),
		},
	})
}

func TestBackend_leaseWriteRead(t *testing.T) {
	if os.Getenv(logicaltest.TestEnvVar) == "" || os.Getenv("MSSQL_URL") == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", logicaltest.TestEnvVar))
	}
	connURL := os.Getenv("MSSQL_URL")

	b := Backend()

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck:       testAccPreCheckFunc(t, connURL),
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, connURL),
			testAccStepWriteLease(t),
			testAccStepReadLease(t),
		},
	})

}

func testAccPreCheckFunc(t *testing.T, connectionURL string) func() {
	return func() {
		if connectionURL == "" {
			t.Fatal("connection URL must be set for acceptance tests")
		}
	}
}

func testAccStepConfig(t *testing.T, connURL string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/connection",
		Data: map[string]interface{}{
			"connection_string": connURL,
		},
	}
}

func testAccStepRole(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/web",
		Data: map[string]interface{}{
			"sql": testRoleSQL,
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

func testAccStepReadRole(t *testing.T, name, sql string) logicaltest.TestStep {
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

const testRoleSQL = `
CREATE LOGIN [{{name}}] WITH PASSWORD = '{{password}}';
CREATE USER [{{name}}] FOR LOGIN [{{name}}];
GRANT SELECT ON SCHEMA::dbo TO [{{name}}]
`
