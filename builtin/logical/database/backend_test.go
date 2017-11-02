package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/helper/pluginutil"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/plugins/database/postgresql"
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

func getCluster(t *testing.T) (*vault.TestCluster, logical.SystemView) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"database": Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	cores := cluster.Cores

	os.Setenv(pluginutil.PluginCACertPEMEnv, cluster.CACertPEMFile)

	sys := vault.TestDynamicSystemView(cores[0].Core)
	vault.TestAddTestPlugin(t, cores[0].Core, "postgresql-database-plugin", "TestBackend_PluginMain")

	return cluster, sys
}

func TestBackend_PluginMain(t *testing.T) {
	if os.Getenv(pluginutil.PluginUnwrapTokenEnv) == "" {
		return
	}

	caPEM := os.Getenv(pluginutil.PluginCACertPEMEnv)
	if caPEM == "" {
		t.Fatal("CA cert not passed in")
	}

	args := []string{"--ca-cert=" + caPEM}

	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(args)

	postgresql.Run(apiClientMeta.GetTLSConfig())
}

func TestBackend_config_connection(t *testing.T) {
	var resp *logical.Response
	var err error

	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

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
		"allowed_roles":     []string{"*"},
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
		"plugin_name": "postgresql-database-plugin",
		"connection_details": map[string]interface{}{
			"connection_url": "sample_connection_url",
		},
		"allowed_roles": []string{"*"},
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

	configReq.Operation = logical.ListOperation
	configReq.Data = nil
	configReq.Path = "config/"
	resp, err = b.HandleRequest(configReq)
	if err != nil {
		t.Fatal(err)
	}
	keys := resp.Data["keys"].([]string)
	key := keys[0]
	if key != "plugin-test" {
		t.Fatalf("bad key: %q", key)
	}
}

func TestBackend_basic(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

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
		"allowed_roles":  []string{"plugin-role-test"},
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

	if !testCredsExist(t, credsResp, connURL) {
		t.Fatalf("Creds should exist")
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

	if testCredsExist(t, credsResp, connURL) {
		t.Fatalf("Creds should not exist")
	}

}

func TestBackend_connectionCrud(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

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
		"connection_url":    "test",
		"plugin_name":       "postgresql-database-plugin",
		"verify_connection": false,
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

	// Update the connection
	data = map[string]interface{}{
		"connection_url": connURL,
		"plugin_name":    "postgresql-database-plugin",
		"allowed_roles":  []string{"plugin-role-test"},
	}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Read connection
	expected := map[string]interface{}{
		"plugin_name": "postgresql-database-plugin",
		"connection_details": map[string]interface{}{
			"connection_url": connURL,
		},
		"allowed_roles": []string{"plugin-role-test"},
	}
	req.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	delete(resp.Data["connection_details"].(map[string]interface{}), "name")
	if !reflect.DeepEqual(expected, resp.Data) {
		t.Fatalf("bad: expected:%#v\nactual:%#v\n", expected, resp.Data)
	}

	// Reset Connection
	data = map[string]interface{}{}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "reset/plugin-test",
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

	if !testCredsExist(t, credsResp, connURL) {
		t.Fatalf("Creds should exist")
	}

	// Delete Connection
	data = map[string]interface{}{}
	req = &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Read connection
	req.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Should be empty
	if resp != nil {
		t.Fatal("Expected response to be nil")
	}
}

func TestBackend_roleCrud(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

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
func TestBackend_allowedRoles(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

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

	// Create a denied and an allowed role
	data = map[string]interface{}{
		"db_name":             "plugin-test",
		"creation_statements": testRole,
		"default_ttl":         "5m",
		"max_ttl":             "10m",
	}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/denied",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	data = map[string]interface{}{
		"db_name":             "plugin-test",
		"creation_statements": testRole,
		"default_ttl":         "5m",
		"max_ttl":             "10m",
	}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/allowed",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Get creds from denied role, should fail
	data = map[string]interface{}{}
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/denied",
		Storage:   config.StorageView,
		Data:      data,
	}
	credsResp, err := b.HandleRequest(req)
	if err != logical.ErrPermissionDenied {
		t.Fatalf("expected error to be:%s got:%#v\n", logical.ErrPermissionDenied, err)
	}

	// update connection with glob allowed roles connection
	data = map[string]interface{}{
		"connection_url": connURL,
		"plugin_name":    "postgresql-database-plugin",
		"allowed_roles":  "allow*",
	}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Get creds, should work.
	data = map[string]interface{}{}
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/allowed",
		Storage:   config.StorageView,
		Data:      data,
	}
	credsResp, err = b.HandleRequest(req)
	if err != nil || (credsResp != nil && credsResp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, credsResp)
	}

	if !testCredsExist(t, credsResp, connURL) {
		t.Fatalf("Creds should exist")
	}

	// update connection with * allowed roles connection
	data = map[string]interface{}{
		"connection_url": connURL,
		"plugin_name":    "postgresql-database-plugin",
		"allowed_roles":  "*",
	}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Get creds, should work.
	data = map[string]interface{}{}
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/allowed",
		Storage:   config.StorageView,
		Data:      data,
	}
	credsResp, err = b.HandleRequest(req)
	if err != nil || (credsResp != nil && credsResp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, credsResp)
	}

	if !testCredsExist(t, credsResp, connURL) {
		t.Fatalf("Creds should exist")
	}

	// update connection with allowed roles
	data = map[string]interface{}{
		"connection_url": connURL,
		"plugin_name":    "postgresql-database-plugin",
		"allowed_roles":  "allow, allowed",
	}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Get creds from denied role, should fail
	data = map[string]interface{}{}
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/denied",
		Storage:   config.StorageView,
		Data:      data,
	}
	credsResp, err = b.HandleRequest(req)
	if err != logical.ErrPermissionDenied {
		t.Fatalf("expected error to be:%s got:%#v\n", logical.ErrPermissionDenied, err)
	}

	// Get creds from allowed role, should work.
	data = map[string]interface{}{}
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/allowed",
		Storage:   config.StorageView,
		Data:      data,
	}
	credsResp, err = b.HandleRequest(req)
	if err != nil || (credsResp != nil && credsResp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, credsResp)
	}

	if !testCredsExist(t, credsResp, connURL) {
		t.Fatalf("Creds should exist")
	}
}

func testCredsExist(t *testing.T, resp *logical.Response, connURL string) bool {
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

	return returnedRows() == 2
}

const testRole = `
CREATE ROLE "{{name}}" WITH
  LOGIN
  PASSWORD '{{password}}'
  VALID UNTIL '{{expiration}}';
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "{{name}}";
`

const defaultRevocationSQL = `
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM {{name}};
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM {{name}};
REVOKE USAGE ON SCHEMA public FROM {{name}};

DROP ROLE IF EXISTS {{name}};
`
