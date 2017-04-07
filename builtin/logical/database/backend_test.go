package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/helper/builtinplugins"
	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	dockertest "gopkg.in/ory-am/dockertest.v3"
)

var (
	testImagePull sync.Once
)

func preparePostgresTestContainer(t *testing.T, s logical.Storage, b logical.Backend) (cleanup func(), retURL string) {
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
		err := pool.Purge(resource)
		if err != nil {
			t.Fatalf("Failed to cleanup local container: %s", err)
		}
	}

	retURL = fmt.Sprintf("postgres://postgres:secret@localhost:%s/database?sslmode=disable", resource.GetPort("5432/tcp"))

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		// This will cause a validation to run
		resp, err := b.HandleRequest(&logical.Request{
			Storage:   s,
			Operation: logical.UpdateOperation,
			Path:      "config/postgresql",
			Data: map[string]interface{}{
				"plugin_name":    "postgresql-database-plugin",
				"connection_url": retURL,
			},
		})
		if err != nil || (resp != nil && resp.IsError()) {
			// It's likely not up and running yet, so return error and try again
			return fmt.Errorf("err:%s resp:%#v\n", err, resp)
		}
		if resp == nil {
			t.Fatal("expected warning")
		}

		return nil
	}); err != nil {
		t.Fatalf("Could not connect to PostgreSQL docker container: %s", err)
	}

	return
}

func getCore(t *testing.T) (*vault.Core, net.Listener, logical.SystemView, string) {
	core, _, token, ln := vault.TestCoreUnsealedWithListener(t)
	http.TestServerWithListener(t, ln, "", core)
	sys := vault.TestDynamicSystemView(core)
	vault.TestAddTestPlugin(t, core, "postgresql-database-plugin", fmt.Sprintf("%s -test.run=TestBackend_PluginMain", os.Args[0]))

	return core, ln, sys, token
}

func TestBackend_PluginMain(t *testing.T) {
	if os.Getenv(pluginutil.PluginUnwrapTokenEnv) == "" {
		return
	}

	f, _ := builtinplugins.BuiltinPlugins.Get("postgresql-database-plugin")
	f()
}

func TestBackend_config_connection(t *testing.T) {
	var resp *logical.Response
	var err error
	_, ln, sys, _ := getCore(t)
	defer ln.Close()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup()

	configData := map[string]interface{}{
		"connection_url":    "sample_connection_url",
		"plugin_name":       "postgresql-database-plugin",
		"verify_connection": false,
	}

	configReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      configData,
	}
	resp, err = b.HandleRequest(configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	expected := map[string]interface{}{
		"plugin_name":        "postgresql-database-plugin",
		"connection_details": configData,
	}
	configReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	delete(resp.Data["connection_details"].(map[string]interface{}), "name")
	if !reflect.DeepEqual(expected, resp.Data) {
		t.Fatalf("bad: expected:%#v\nactual:%#v\n", expected, resp.Data)
	}
}

func TestBackend_basic(t *testing.T) {
	_, ln, sys, _ := getCore(t)
	defer ln.Close()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys

	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup()

	cleanup, connURL := preparePostgresTestContainer(t, config.StorageView, b)
	defer cleanup()

	// Configure a connection
	data := map[string]interface{}{
		"connection_url": connURL,
		"plugin_name":    "postgresql-database-plugin",
	}
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err := b.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Create a role
	data = map[string]interface{}{
		"db_name":             "plugin-test",
		"creation_statements": testRole,
		"default_ttl":         "5m",
		"max_ttl":             "10m",
	}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/plugin-role-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Get creds
	data = map[string]interface{}{}
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/plugin-role-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	credsResp, err := b.HandleRequest(req)
	if err != nil || (credsResp != nil && credsResp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, credsResp)
	}

	if testCredsByCount(t, credsResp, connURL) != 2 {
		t.Fatalf("Got wrong number of creds")
	}

	// Revoke creds
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.RevokeOperation,
		Storage:   config.StorageView,
		Secret: &logical.Secret{
			InternalData: map[string]interface{}{
				"secret_type": "creds",
				"username":    credsResp.Data["username"],
				"role":        "plugin-role-test",
			},
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	if testCredsByCount(t, credsResp, connURL) != -1 {
		t.Fatalf("Got wrong number of creds")
	}

}

func TestBackend_roleCrud(t *testing.T) {
	_, ln, sys, _ := getCore(t)
	defer ln.Close()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys

	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup()

	cleanup, connURL := preparePostgresTestContainer(t, config.StorageView, b)
	defer cleanup()

	// Configure a connection
	data := map[string]interface{}{
		"connection_url": connURL,
		"plugin_name":    "postgresql-database-plugin",
	}
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err := b.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Create a role
	data = map[string]interface{}{
		"db_name":               "plugin-test",
		"creation_statements":   testRole,
		"revocation_statements": defaultRevocationSQL,
		"default_ttl":           "5m",
		"max_ttl":               "10m",
	}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/plugin-role-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Read the role
	data = map[string]interface{}{}
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "roles/plugin-role-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	expected := dbplugin.Statements{
		CreationStatements:   testRole,
		RevocationStatements: defaultRevocationSQL,
	}

	var actual dbplugin.Statements
	if err := mapstructure.Decode(resp.Data, &actual); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Statements did not match, exepected %#v, got %#v", expected, actual)
	}

	// Delete the role
	data = map[string]interface{}{}
	req = &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      "roles/plugin-role-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Read the role
	data = map[string]interface{}{}
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "roles/plugin-role-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Should be empty
	if resp != nil {
		t.Fatal("Expected response to be nil")
	}
}

func TestBackend_roleReadOnly(t *testing.T) {
	_, ln, sys, _ := getCore(t)
	defer ln.Close()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys

	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup()

	cleanup, connURL := preparePostgresTestContainer(t, config.StorageView, b)
	defer cleanup()

	// Configure a connection
	data := map[string]interface{}{
		"connection_url": connURL,
		"plugin_name":    "postgresql-database-plugin",
	}
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err := b.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Create a role
	data = map[string]interface{}{
		"db_name":             "plugin-test",
		"creation_statements": testRole,
		"default_ttl":         "5m",
		"max_ttl":             "10m",
	}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/plugin-role-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Create a readonly role
	data = map[string]interface{}{
		"db_name":             "plugin-test",
		"creation_statements": testReadOnlyRole,
		"default_ttl":         "5m",
		"max_ttl":             "10m",
	}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/plugin-readonly-role-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Get creds
	data = map[string]interface{}{}
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/plugin-role-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	credsResp, err := b.HandleRequest(req)
	if err != nil || (credsResp != nil && credsResp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, credsResp)
	}

	if i := testCredsByCount(t, credsResp, connURL); i != 2 {
		t.Fatalf("Got wrong number of creds got %d, expected 2", i)
	}

	// Get readonly creds
	data = map[string]interface{}{}
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/plugin-readonly-role-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	readOnlyCredsResp, err := b.HandleRequest(req)
	if err != nil || (credsResp != nil && credsResp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, readOnlyCredsResp)
	}

	if i := testCredsByCount(t, readOnlyCredsResp, connURL); i != 2 {
		t.Fatalf("Got wrong number of creds got %d, expected 2", i)
	}

	if err := testCreateTable(t, readOnlyCredsResp, connURL); err == nil {
		t.Fatal("Read only creds should return error on table creation")
	}

	if err := testCreateTable(t, credsResp, connURL); err != nil {
		t.Fatalf("Error on table creation: %s", err)
	}
}

func testCredsByCount(t *testing.T, resp *logical.Response, connURL string) int {
	var d struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	}
	if err := mapstructure.Decode(resp.Data, &d); err != nil {
		t.Fatal(err)
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

	return returnedRows()
}

func testCreateTable(t *testing.T, resp *logical.Response, connURL string) error {
	var d struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	}
	if err := mapstructure.Decode(resp.Data, &d); err != nil {
		t.Fatal(err)
	}

	connURL = strings.Replace(connURL, "postgres:secret", fmt.Sprintf("%s:%s", d.Username, d.Password), 1)

	fmt.Println(connURL)
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

	r, err := db.Exec("CREATE TABLE test1 (id SERIAL PRIMARY KEY);")
	if err != nil {
		return err
	}

	if i, _ := r.RowsAffected(); i != 1 {
		return errors.New("Did not create db")
	}

	return nil
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
REVOKE ALL ON SCHEMA public FROM "{{name}}";
GRANT SELECT ON ALL TABLES IN SCHEMA public TO "{{name}}";
GRANT SELECT ON ALL SEQUENCES IN SCHEMA public TO "{{name}}";
`

const defaultRevocationSQL = `
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM {{name}};
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM {{name}};
REVOKE USAGE ON SCHEMA public FROM {{name}};

DROP ROLE IF EXISTS {{name}};
`
