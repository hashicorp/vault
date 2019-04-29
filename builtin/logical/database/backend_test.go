package database

import (
	"context"
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers/docker"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/plugins/database/mongodb"
	"github.com/hashicorp/vault/plugins/database/postgresql"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"github.com/ory/dockertest"
	"gopkg.in/mgo.v2"
)

var (
	testImagePull sync.Once
)

func preparePostgresTestContainer(t *testing.T, s logical.Storage, b logical.Backend) (cleanup func(), retURL string) {
	t.Helper()
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
		// This will cause a validation to run
		resp, err := b.HandleRequest(namespace.RootContext(nil), &logical.Request{
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
			return fmt.Errorf("err:%#v resp:%#v", err, resp)
		}
		if resp == nil {
			t.Fatal("expected warning")
		}

		return nil
	}); err != nil {
		cleanup()
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
	vault.TestAddTestPlugin(t, cores[0].Core, "postgresql-database-plugin", consts.PluginTypeDatabase, "TestBackend_PluginMain_Postgres", []string{}, "")
	vault.TestAddTestPlugin(t, cores[0].Core, "mongodb-database-plugin", consts.PluginTypeDatabase, "TestBackend_PluginMain_Mongo", []string{}, "")

	return cluster, sys
}

func TestBackend_PluginMain_Postgres(t *testing.T) {
	if os.Getenv(pluginutil.PluginUnwrapTokenEnv) == "" {
		return
	}

	caPEM := os.Getenv(pluginutil.PluginCACertPEMEnv)
	if caPEM == "" {
		t.Fatal("CA cert not passed in")
	}

	args := []string{"--ca-cert=" + caPEM}

	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(args)

	postgresql.Run(apiClientMeta.GetTLSConfig())
}

func TestBackend_PluginMain_Mongo(t *testing.T) {
	if os.Getenv(pluginutil.PluginUnwrapTokenEnv) == "" {
		return
	}

	caPEM := os.Getenv(pluginutil.PluginCACertPEMEnv)
	if caPEM == "" {
		t.Fatal("CA cert not passed in")
	}

	args := []string{"--ca-cert=" + caPEM}

	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(args)

	err := mongodb.Run(apiClientMeta.GetTLSConfig())
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackend_RoleUpgrade(t *testing.T) {

	storage := &logical.InmemStorage{}
	backend := &databaseBackend{}

	roleExpected := &roleEntry{
		Statements: dbplugin.Statements{
			CreationStatements: "test",
			Creation:           []string{"test"},
		},
	}

	entry, err := logical.StorageEntryJSON("role/test", &roleEntry{
		Statements: dbplugin.Statements{
			CreationStatements: "test",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := storage.Put(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	role, err := backend.Role(context.Background(), storage, "test")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(role, roleExpected) {
		t.Fatalf("bad role %#v, %#v", role, roleExpected)
	}

	// Upgrade case
	badJSON := `{"statments":{"creation_statments":"test","revocation_statements":"","rollback_statements":"","renew_statements":""}}`
	entry = &logical.StorageEntry{
		Key:   "role/test",
		Value: []byte(badJSON),
	}
	if err := storage.Put(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	role, err = backend.Role(context.Background(), storage, "test")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(role, roleExpected) {
		t.Fatalf("bad role %#v, %#v", role, roleExpected)
	}

}

func TestBackend_config_connection(t *testing.T) {
	var resp *logical.Response
	var err error

	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys
	lb, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	b, ok := lb.(*databaseBackend)
	if !ok {
		t.Fatal("could not convert to database backend")
	}
	defer b.Cleanup(context.Background())

	// Test creation
	{
		configData := map[string]interface{}{
			"connection_url":    "sample_connection_url",
			"someotherdata":     "testing",
			"plugin_name":       "postgresql-database-plugin",
			"verify_connection": false,
			"allowed_roles":     []string{"*"},
			"name":              "plugin-test",
		}

		configReq := &logical.Request{
			Operation: logical.CreateOperation,
			Path:      "config/plugin-test",
			Storage:   config.StorageView,
			Data:      configData,
		}

		exists, err := b.connectionExistenceCheck()(context.Background(), configReq, &framework.FieldData{
			Raw:    configData,
			Schema: pathConfigurePluginConnection(b).Fields,
		})
		if err != nil {
			t.Fatal(err)
		}
		if exists {
			t.Fatal("expected not exists")
		}

		resp, err = b.HandleRequest(namespace.RootContext(nil), configReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%v resp:%#v\n", err, resp)
		}

		expected := map[string]interface{}{
			"plugin_name": "postgresql-database-plugin",
			"connection_details": map[string]interface{}{
				"connection_url": "sample_connection_url",
				"someotherdata":  "testing",
			},
			"allowed_roles":                      []string{"*"},
			"root_credentials_rotate_statements": []string{},
		}
		configReq.Operation = logical.ReadOperation
		resp, err = b.HandleRequest(namespace.RootContext(nil), configReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%s resp:%#v\n", err, resp)
		}

		delete(resp.Data["connection_details"].(map[string]interface{}), "name")
		if !reflect.DeepEqual(expected, resp.Data) {
			t.Fatalf("bad: expected:%#v\nactual:%#v\n", expected, resp.Data)
		}
	}

	// Test existence check and an update to a single connection detail parameter
	{
		configData := map[string]interface{}{
			"connection_url":    "sample_convection_url",
			"verify_connection": false,
			"name":              "plugin-test",
		}

		configReq := &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "config/plugin-test",
			Storage:   config.StorageView,
			Data:      configData,
		}

		exists, err := b.connectionExistenceCheck()(context.Background(), configReq, &framework.FieldData{
			Raw:    configData,
			Schema: pathConfigurePluginConnection(b).Fields,
		})
		if err != nil {
			t.Fatal(err)
		}
		if !exists {
			t.Fatal("expected exists")
		}

		resp, err = b.HandleRequest(namespace.RootContext(nil), configReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%v resp:%#v\n", err, resp)
		}

		expected := map[string]interface{}{
			"plugin_name": "postgresql-database-plugin",
			"connection_details": map[string]interface{}{
				"connection_url": "sample_convection_url",
				"someotherdata":  "testing",
			},
			"allowed_roles":                      []string{"*"},
			"root_credentials_rotate_statements": []string{},
		}
		configReq.Operation = logical.ReadOperation
		resp, err = b.HandleRequest(namespace.RootContext(nil), configReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%s resp:%#v\n", err, resp)
		}

		delete(resp.Data["connection_details"].(map[string]interface{}), "name")
		if !reflect.DeepEqual(expected, resp.Data) {
			t.Fatalf("bad: expected:%#v\nactual:%#v\n", expected, resp.Data)
		}
	}

	// Test an update to a non-details value
	{
		configData := map[string]interface{}{
			"verify_connection": false,
			"allowed_roles":     []string{"flu", "barre"},
			"name":              "plugin-test",
		}

		configReq := &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "config/plugin-test",
			Storage:   config.StorageView,
			Data:      configData,
		}

		resp, err = b.HandleRequest(namespace.RootContext(nil), configReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%v resp:%#v\n", err, resp)
		}

		expected := map[string]interface{}{
			"plugin_name": "postgresql-database-plugin",
			"connection_details": map[string]interface{}{
				"connection_url": "sample_convection_url",
				"someotherdata":  "testing",
			},
			"allowed_roles":                      []string{"flu", "barre"},
			"root_credentials_rotate_statements": []string{},
		}
		configReq.Operation = logical.ReadOperation
		resp, err = b.HandleRequest(namespace.RootContext(nil), configReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%s resp:%#v\n", err, resp)
		}

		delete(resp.Data["connection_details"].(map[string]interface{}), "name")
		if !reflect.DeepEqual(expected, resp.Data) {
			t.Fatalf("bad: expected:%#v\nactual:%#v\n", expected, resp.Data)
		}
	}

	req := &logical.Request{
		Operation: logical.ListOperation,
		Storage:   config.StorageView,
		Path:      "config/",
	}
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatal(err)
	}
	keys := resp.Data["keys"].([]string)
	key := keys[0]
	if key != "plugin-test" {
		t.Fatalf("bad key: %q", key)
	}
}

func TestBackend_BadConnectionString(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup(context.Background())

	cleanup, _ := preparePostgresTestContainer(t, config.StorageView, b)
	defer cleanup()

	respCheck := func(req *logical.Request) {
		t.Helper()
		resp, err := b.HandleRequest(namespace.RootContext(nil), req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if resp == nil || !resp.IsError() {
			t.Fatalf("expected error, resp:%#v", resp)
		}
		err = resp.Error()
		if strings.Contains(err.Error(), "localhost") {
			t.Fatalf("error should not contain connection info")
		}
	}

	// Configure a connection
	data := map[string]interface{}{
		"connection_url": "postgresql://:pw@[localhost",
		"plugin_name":    "postgresql-database-plugin",
		"allowed_roles":  []string{"plugin-role-test"},
	}
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	respCheck(req)

	time.Sleep(1 * time.Second)
}

func TestBackend_basic(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup(context.Background())

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
	resp, err := b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Create a role
	data = map[string]interface{}{
		"db_name":             "plugin-test",
		"creation_statements": testRole,
		"max_ttl":             "10m",
	}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/plugin-role-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
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
	credsResp, err := b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (credsResp != nil && credsResp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, credsResp)
	}
	// Update the role with no max ttl
	data = map[string]interface{}{
		"db_name":             "plugin-test",
		"creation_statements": testRole,
		"default_ttl":         "5m",
		"max_ttl":             0,
	}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/plugin-role-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
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
	credsResp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (credsResp != nil && credsResp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, credsResp)
	}
	// Test for #3812
	if credsResp.Secret.TTL != 5*time.Minute {
		t.Fatalf("unexpected TTL of %d", credsResp.Secret.TTL)
	}
	// Update the role with a max ttl
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
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Get creds and revoke when the role stays in existence
	{
		data = map[string]interface{}{}
		req = &logical.Request{
			Operation: logical.ReadOperation,
			Path:      "creds/plugin-role-test",
			Storage:   config.StorageView,
			Data:      data,
		}
		credsResp, err = b.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (credsResp != nil && credsResp.IsError()) {
			t.Fatalf("err:%s resp:%#v\n", err, credsResp)
		}
		// Test for #3812
		if credsResp.Secret.TTL != 5*time.Minute {
			t.Fatalf("unexpected TTL of %d", credsResp.Secret.TTL)
		}
		if !testCredsExist(t, credsResp, connURL) {
			t.Fatalf("Creds should exist")
		}

		// Revoke creds
		resp, err = b.HandleRequest(namespace.RootContext(nil), &logical.Request{
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

	// Get creds and revoke using embedded revocation data
	{
		data = map[string]interface{}{}
		req = &logical.Request{
			Operation: logical.ReadOperation,
			Path:      "creds/plugin-role-test",
			Storage:   config.StorageView,
			Data:      data,
		}
		credsResp, err = b.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (credsResp != nil && credsResp.IsError()) {
			t.Fatalf("err:%s resp:%#v\n", err, credsResp)
		}
		if !testCredsExist(t, credsResp, connURL) {
			t.Fatalf("Creds should exist")
		}

		// Delete role, forcing us to rely on embedded data
		req = &logical.Request{
			Operation: logical.DeleteOperation,
			Path:      "roles/plugin-role-test",
			Storage:   config.StorageView,
		}
		resp, err = b.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%s resp:%#v\n", err, resp)
		}

		// Revoke creds
		resp, err = b.HandleRequest(namespace.RootContext(nil), &logical.Request{
			Operation: logical.RevokeOperation,
			Storage:   config.StorageView,
			Secret: &logical.Secret{
				InternalData: map[string]interface{}{
					"secret_type":           "creds",
					"username":              credsResp.Data["username"],
					"role":                  "plugin-role-test",
					"db_name":               "plugin-test",
					"revocation_statements": nil,
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
}

func TestBackend_connectionCrud(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup(context.Background())

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
	resp, err := b.HandleRequest(namespace.RootContext(nil), req)
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
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Update the connection
	data = map[string]interface{}{
		"connection_url": connURL,
		"plugin_name":    "postgresql-database-plugin",
		"allowed_roles":  []string{"plugin-role-test"},
		"username":       "postgres",
		"password":       "secret",
	}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}
	if len(resp.Warnings) != 1 {
		t.Fatalf("expected warning about password in url %s, resp:%#v\n", connURL, resp)
	}

	req.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}
	if strings.Contains(resp.Data["connection_details"].(map[string]interface{})["connection_url"].(string), "secret") {
		t.Fatal("password should not be found in the connection url")
	}

	// Replace connection url with templated version
	req.Operation = logical.UpdateOperation
	connURL = strings.Replace(connURL, "postgres:secret", "{{username}}:{{password}}", -1)
	data["connection_url"] = connURL
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Read connection
	expected := map[string]interface{}{
		"plugin_name": "postgresql-database-plugin",
		"connection_details": map[string]interface{}{
			"username":       "postgres",
			"connection_url": connURL,
		},
		"allowed_roles":                      []string{"plugin-role-test"},
		"root_credentials_rotate_statements": []string(nil),
	}
	req.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	delete(resp.Data["connection_details"].(map[string]interface{}), "name")
	if diff := deep.Equal(resp.Data, expected); diff != nil {
		t.Fatal(diff)
	}

	// Reset Connection
	data = map[string]interface{}{}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "reset/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
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
	credsResp, err := b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (credsResp != nil && credsResp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, credsResp)
	}

	credCheckURL := dbutil.QueryHelper(connURL, map[string]string{
		"username": "postgres",
		"password": "secret",
	})
	if !testCredsExist(t, credsResp, credCheckURL) {
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
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Read connection
	req.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
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

	lb, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	b, ok := lb.(*databaseBackend)
	if !ok {
		t.Fatal("could not convert to db backend")
	}
	defer b.Cleanup(context.Background())

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
	resp, err := b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Test role creation
	{
		data = map[string]interface{}{
			"db_name":               "plugin-test",
			"creation_statements":   testRole,
			"revocation_statements": defaultRevocationSQL,
			"default_ttl":           "5m",
			"max_ttl":               "10m",
		}
		req = &logical.Request{
			Operation: logical.CreateOperation,
			Path:      "roles/plugin-role-test",
			Storage:   config.StorageView,
			Data:      data,
		}
		resp, err = b.HandleRequest(namespace.RootContext(nil), req)
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
		resp, err = b.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%s resp:%#v\n", err, resp)
		}

		expected := dbplugin.Statements{
			Creation:   []string{strings.TrimSpace(testRole)},
			Revocation: []string{strings.TrimSpace(defaultRevocationSQL)},
			Rollback:   []string{},
			Renewal:    []string{},
		}

		actual := dbplugin.Statements{
			Creation:   resp.Data["creation_statements"].([]string),
			Revocation: resp.Data["revocation_statements"].([]string),
			Rollback:   resp.Data["rollback_statements"].([]string),
			Renewal:    resp.Data["renew_statements"].([]string),
		}

		if diff := deep.Equal(expected, actual); diff != nil {
			t.Fatal(diff)
		}

		if diff := deep.Equal(resp.Data["db_name"], "plugin-test"); diff != nil {
			t.Fatal(diff)
		}
		if diff := deep.Equal(resp.Data["default_ttl"], float64(300)); diff != nil {
			t.Fatal(diff)
		}
		if diff := deep.Equal(resp.Data["max_ttl"], float64(600)); diff != nil {
			t.Fatal(diff)
		}
	}

	// Test role modification of TTL
	{
		data = map[string]interface{}{
			"name":    "plugin-role-test",
			"max_ttl": "7m",
		}
		req = &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "roles/plugin-role-test",
			Storage:   config.StorageView,
			Data:      data,
		}
		resp, err = b.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%v resp:%#v\n", err, resp)
		}

		// Read the role
		data = map[string]interface{}{}
		req = &logical.Request{
			Operation: logical.ReadOperation,
			Path:      "roles/plugin-role-test",
			Storage:   config.StorageView,
			Data:      data,
		}
		resp, err = b.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%s resp:%#v\n", err, resp)
		}

		expected := dbplugin.Statements{
			Creation:   []string{strings.TrimSpace(testRole)},
			Revocation: []string{strings.TrimSpace(defaultRevocationSQL)},
			Rollback:   []string{},
			Renewal:    []string{},
		}

		actual := dbplugin.Statements{
			Creation:   resp.Data["creation_statements"].([]string),
			Revocation: resp.Data["revocation_statements"].([]string),
			Rollback:   resp.Data["rollback_statements"].([]string),
			Renewal:    resp.Data["renew_statements"].([]string),
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("Statements did not match, expected %#v, got %#v", expected, actual)
		}

		if diff := deep.Equal(resp.Data["db_name"], "plugin-test"); diff != nil {
			t.Fatal(diff)
		}
		if diff := deep.Equal(resp.Data["default_ttl"], float64(300)); diff != nil {
			t.Fatal(diff)
		}
		if diff := deep.Equal(resp.Data["max_ttl"], float64(420)); diff != nil {
			t.Fatal(diff)
		}

	}

	// Test role modification of statements
	{
		data = map[string]interface{}{
			"name":                  "plugin-role-test",
			"creation_statements":   []string{testRole, testRole},
			"revocation_statements": []string{defaultRevocationSQL, defaultRevocationSQL},
			"rollback_statements":   testRole,
			"renew_statements":      defaultRevocationSQL,
		}
		req = &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "roles/plugin-role-test",
			Storage:   config.StorageView,
			Data:      data,
		}
		resp, err = b.HandleRequest(context.Background(), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%v resp:%#v\n", err, resp)
		}

		// Read the role
		data = map[string]interface{}{}
		req = &logical.Request{
			Operation: logical.ReadOperation,
			Path:      "roles/plugin-role-test",
			Storage:   config.StorageView,
			Data:      data,
		}
		resp, err = b.HandleRequest(context.Background(), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%s resp:%#v\n", err, resp)
		}

		expected := dbplugin.Statements{
			Creation:   []string{strings.TrimSpace(testRole), strings.TrimSpace(testRole)},
			Rollback:   []string{strings.TrimSpace(testRole)},
			Revocation: []string{strings.TrimSpace(defaultRevocationSQL), strings.TrimSpace(defaultRevocationSQL)},
			Renewal:    []string{strings.TrimSpace(defaultRevocationSQL)},
		}

		actual := dbplugin.Statements{
			Creation:   resp.Data["creation_statements"].([]string),
			Revocation: resp.Data["revocation_statements"].([]string),
			Rollback:   resp.Data["rollback_statements"].([]string),
			Renewal:    resp.Data["renew_statements"].([]string),
		}

		if diff := deep.Equal(expected, actual); diff != nil {
			t.Fatal(diff)
		}

		if diff := deep.Equal(resp.Data["db_name"], "plugin-test"); diff != nil {
			t.Fatal(diff)
		}
		if diff := deep.Equal(resp.Data["default_ttl"], float64(300)); diff != nil {
			t.Fatal(diff)
		}
		if diff := deep.Equal(resp.Data["max_ttl"], float64(420)); diff != nil {
			t.Fatal(diff)
		}
	}

	// Delete the role
	data = map[string]interface{}{}
	req = &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      "roles/plugin-role-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
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
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
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

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup(context.Background())

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
	resp, err := b.HandleRequest(namespace.RootContext(nil), req)
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
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
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
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
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
	credsResp, err := b.HandleRequest(namespace.RootContext(nil), req)
	if err == nil {
		t.Fatal("expected error because role is denied")
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
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
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
	credsResp, err = b.HandleRequest(namespace.RootContext(nil), req)
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
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
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
	credsResp, err = b.HandleRequest(namespace.RootContext(nil), req)
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
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
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
	credsResp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err == nil {
		t.Fatal("expected error because role is denied")
	}

	// Get creds from allowed role, should work.
	data = map[string]interface{}{}
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/allowed",
		Storage:   config.StorageView,
		Data:      data,
	}
	credsResp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (credsResp != nil && credsResp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, credsResp)
	}

	if !testCredsExist(t, credsResp, connURL) {
		t.Fatalf("Creds should exist")
	}
}

func TestBackend_RotateRootCredentials(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup(context.Background())

	cleanup, connURL := preparePostgresTestContainer(t, config.StorageView, b)
	defer cleanup()

	connURL = strings.Replace(connURL, "postgres:secret", "{{username}}:{{password}}", -1)

	// Configure a connection
	data := map[string]interface{}{
		"connection_url": connURL,
		"plugin_name":    "postgresql-database-plugin",
		"allowed_roles":  []string{"plugin-role-test"},
		"username":       "postgres",
		"password":       "secret",
	}
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err := b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Create a role
	data = map[string]interface{}{
		"db_name":             "plugin-test",
		"creation_statements": testRole,
		"max_ttl":             "10m",
	}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/plugin-role-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
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
	credsResp, err := b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (credsResp != nil && credsResp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, credsResp)
	}

	data = map[string]interface{}{}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "rotate-root/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (credsResp != nil && credsResp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, credsResp)
	}

	dbConfig, err := b.(*databaseBackend).DatabaseConfig(context.Background(), config.StorageView, "plugin-test")
	if err != nil {
		t.Fatalf("err: %#v", err)
	}
	if dbConfig.ConnectionDetails["password"].(string) == "secret" {
		t.Fatal("root credentials not rotated")
	}

	// Get creds to make sure it still works
	data = map[string]interface{}{}
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/plugin-role-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	credsResp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (credsResp != nil && credsResp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, credsResp)
	}
}

func TestBackend_StaticRole_Rotations_PostgreSQL(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup(context.Background())

	// allow initQueue to finish
	time.Sleep(5 * time.Second)
	bd := b.(*databaseBackend)
	if bd.credRotationQueue == nil {
		t.Fatal("database backend had no credential rotation queue")
	}

	// configure backend, add item and confirm length
	cleanup, connURL := preparePostgresTestContainer(t, config.StorageView, b)
	defer cleanup()

	// Configure a connection
	data := map[string]interface{}{
		"connection_url":    connURL,
		"plugin_name":       "postgresql-database-plugin",
		"verify_connection": false,
		"allowed_roles":     []string{"*"},
		"name":              "plugin-test",
	}
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err := b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// create three static roles with different rotation periods
	testCases := []string{"65", "130", "5400"}
	for _, tc := range testCases {
		roleName := "plugin-static-role-" + tc
		data = map[string]interface{}{
			"name":                  roleName,
			"db_name":               "plugin-test",
			"creation_statements":   testRoleStaticCreate,
			"rotation_statements":   testRoleStaticUpdate,
			"revocation_statements": defaultRevocationSQL,
			"username":              "statictest" + tc,
			"rotation_period":       tc,
		}

		req = &logical.Request{
			Operation: logical.CreateOperation,
			Path:      "static-roles/" + roleName,
			Storage:   config.StorageView,
			Data:      data,
		}

		resp, err = b.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%s resp:%#v\n", err, resp)
		}
	}

	// verify the queue has 3 items in it
	if bd.credRotationQueue.Len() != 3 {
		t.Fatalf("expected 3 items in the rotation queue, got: (%d)", bd.credRotationQueue.Len())
	}

	// List the roles
	data = map[string]interface{}{}
	req = &logical.Request{
		Operation: logical.ListOperation,
		Path:      "static-roles/",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	keys := resp.Data["keys"].([]string)
	if len(keys) != 3 {
		t.Fatalf("expected 3 roles, got: (%d)", len(keys))
	}

	// capture initial passwords, before the periodic function is triggered
	pws := make(map[string][]string, 0)
	pws = capturePasswords(t, b, config, testCases, pws)

	// sleep to make sure the 65s role will be up for rotation by the time the
	// periodic function ticks
	time.Sleep(7 * time.Second)

	// sleep 75 to make sure the periodic func has time to actually run
	time.Sleep(75 * time.Second)
	pws = capturePasswords(t, b, config, testCases, pws)

	// sleep more, this should allow both sr65 and sr130 to rotate
	time.Sleep(140 * time.Second)
	pws = capturePasswords(t, b, config, testCases, pws)

	// verify all pws are as they should
	pass := true
	for k, v := range pws {
		switch {
		case k == "plugin-static-role-65":
			// expect all passwords to be different
			if v[0] == v[1] || v[1] == v[2] || v[0] == v[2] {
				pass = false
			}
		case k == "plugin-static-role-130":
			// expect the first two to be equal, but different from the third
			if v[0] != v[1] || v[0] == v[2] {
				pass = false
			}
		case k == "plugin-static-role-5400":
			// expect all passwords to be equal
			if v[0] != v[1] || v[1] != v[2] {
				pass = false
			}
		}
	}
	if !pass {
		t.Fatalf("password rotations did not match expected: %#v", pws)
	}
}

// copied from plugins/database/mongodb_test.go
func TestBackend_StaticRole_Rotations_MongoDB(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup(context.Background())

	// allow initQueue to finish
	time.Sleep(5 * time.Second)
	bd := b.(*databaseBackend)
	if bd.credRotationQueue == nil {
		t.Fatal("database backend had no credential rotation queue")
	}

	// configure backend, add item and confirm length
	cleanup, connURL := prepareMongoDBTestContainer(t)
	defer cleanup()

	// Configure a connection
	data := map[string]interface{}{
		"connection_url":    connURL,
		"plugin_name":       "mongodb-database-plugin",
		"verify_connection": false,
		"allowed_roles":     []string{"*"},
		"name":              "plugin-mongo-test",
	}

	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/plugin-mongo-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err := b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// create three static roles with different rotation periods
	testCases := []string{"65", "130", "5400"}
	for _, tc := range testCases {
		roleName := "plugin-static-role-" + tc
		data = map[string]interface{}{
			"name":                roleName,
			"db_name":             "plugin-mongo-test",
			"creation_statements": testMongoDBRole,
			"username":            "statictestMongo" + tc,
			"rotation_period":     tc,
		}

		req = &logical.Request{
			Operation: logical.CreateOperation,
			Path:      "static-roles/" + roleName,
			Storage:   config.StorageView,
			Data:      data,
		}

		resp, err = b.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%s resp:%#v\n", err, resp)
		}
	}

	// verify the queue has 3 items in it
	if bd.credRotationQueue.Len() != 3 {
		t.Fatalf("expected 3 items in the rotation queue, got: (%d)", bd.credRotationQueue.Len())
	}

	// List the roles
	data = map[string]interface{}{}
	req = &logical.Request{
		Operation: logical.ListOperation,
		Path:      "static-roles/",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	keys := resp.Data["keys"].([]string)
	if len(keys) != 3 {
		t.Fatalf("expected 3 roles, got: (%d)", len(keys))
	}

	// capture initial passwords, before the periodic function is triggered
	pws := make(map[string][]string, 0)
	pws = capturePasswords(t, b, config, testCases, pws)

	// sleep to make sure the 65s role will be up for rotation by the time the
	// periodic function ticks
	time.Sleep(7 * time.Second)

	// sleep 75 to make sure the periodic func has time to actually run
	time.Sleep(75 * time.Second)
	pws = capturePasswords(t, b, config, testCases, pws)

	// sleep more, this should allow both sr65 and sr130 to rotate
	time.Sleep(140 * time.Second)
	pws = capturePasswords(t, b, config, testCases, pws)

	// verify all pws are as they should
	pass := true
	for k, v := range pws {
		switch {
		case k == "plugin-static-role-65":
			// expect all passwords to be different
			if v[0] == v[1] || v[1] == v[2] || v[0] == v[2] {
				pass = false
			}
		case k == "plugin-static-role-130":
			// expect the first two to be equal, but different from the third
			if v[0] != v[1] || v[0] == v[2] {
				pass = false
			}
		case k == "plugin-static-role-5400":
			// expect all passwords to be equal
			if v[0] != v[1] || v[1] != v[2] {
				pass = false
			}
		}
	}
	if !pass {
		t.Fatalf("password rotations did not match expected: %#v", pws)
	}
}

// capture the current passwords
func capturePasswords(t *testing.T, b logical.Backend, config *logical.BackendConfig, testCases []string, pws map[string][]string) map[string][]string {
	new := make(map[string][]string, 0)
	for _, tc := range testCases {
		// Read the role
		roleName := "plugin-static-role-" + tc
		req := &logical.Request{
			Operation: logical.ReadOperation,
			Path:      "static-creds/" + roleName,
			Storage:   config.StorageView,
		}
		resp, err := b.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%s resp:%#v\n", err, resp)
		}

		username := resp.Data["username"].(string)
		password := resp.Data["password"].(string)
		if username == "" || password == "" {
			t.Fatalf("expected both username/password for (%s), got (%s), (%s)", roleName, username, password)
		}
		new[roleName] = append(new[roleName], password)
	}

	for k, v := range new {
		pws[k] = append(pws[k], v...)
	}

	return pws
}

func testCredsExist(t *testing.T, resp *logical.Response, connURL string) bool {
	t.Helper()
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
const testMongoDBRole = `{ "db": "admin", "roles": [ { "role": "readWrite" } ] }`

//
// WAL testing
//
// First scenario, WAL contains a role name that does not exist.
func TestBackend_Static_QueueWAL_discard_role_not_found(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

	ctx := context.Background()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys

	_, err := framework.PutWAL(ctx, config.StorageView, staticWALKey, &setCredentialsWAL{
		RoleName: "doesnotexist",
	})
	if err != nil {
		t.Fatalf("error with PutWAL: %s", err)
	}

	assertWALCount(t, config.StorageView, 1)

	b, err := Factory(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup(ctx)

	cleanup, _ := preparePostgresTestContainer(t, config.StorageView, b)
	defer cleanup()

	time.Sleep(1 * time.Second)
	bd := b.(*databaseBackend)
	if bd.credRotationQueue == nil {
		t.Fatal("database backend had no credential rotation queue")
	}

	// verify empty queue
	if bd.credRotationQueue.Len() != 0 {
		t.Fatalf("expected zero queue items, got: %d", bd.credRotationQueue.Len())
	}

	assertWALCount(t, config.StorageView, 0)
}

// Second scenario, WAL contains a role name that does exist, but the role's
// LastVaultRotation is greater than the WAL has
func TestBackend_Static_QueueWAL_discard_role_newer_rotation_date(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

	ctx := context.Background()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys

	roleName := "test-discard-by-date"
	lb, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	b, ok := lb.(*databaseBackend)
	if !ok {
		t.Fatal("could not convert to db backend")
	}

	cleanup, connURL := preparePostgresTestContainer(t, config.StorageView, b)
	defer cleanup()

	// Configure a connection
	data := map[string]interface{}{
		"connection_url":    connURL,
		"plugin_name":       "postgresql-database-plugin",
		"verify_connection": false,
		"allowed_roles":     []string{"*"},
		"name":              "plugin-test",
	}
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err := b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// save Now() to make sure rotation time is after this, as well as the WAL
	// time
	roleTime := time.Now()

	// create role
	data = map[string]interface{}{
		"name":                  roleName,
		"db_name":               "plugin-test",
		"creation_statements":   testRoleStaticCreate,
		"rotation_statements":   testRoleStaticUpdate,
		"revocation_statements": defaultRevocationSQL,
		"username":              "statictest",
		// low value here, to make sure the backend rotates this password at least
		// once before we compare it to the WAL
		"rotation_period": "10s",
	}

	req = &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "static-roles/" + roleName,
		Storage:   config.StorageView,
		Data:      data,
	}

	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// allow the first rotation to occur, setting LastVaultRotation
	time.Sleep(time.Second * 12)

	// cleanup the backend, then create a WAL for the role with a
	// LastVaultRotation of 1 hour ago, so that when we recreate the backend the
	// WAL will be read but discarded
	b.Cleanup(ctx)
	b = nil
	time.Sleep(time.Second * 3)

	// make a fake WAL entry with an older time
	oldRotationTime := roleTime.Add(time.Hour * -1)
	walPassword := "somejunkpassword"
	_, err = framework.PutWAL(ctx, config.StorageView, staticWALKey, &setCredentialsWAL{
		RoleName:          roleName,
		NewPassword:       walPassword,
		LastVaultRotation: oldRotationTime,
	})
	if err != nil {
		t.Fatalf("error with PutWAL: %s", err)
	}

	assertWALCount(t, config.StorageView, 1)

	// reload backend
	lb, err = Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	b, ok = lb.(*databaseBackend)
	if !ok {
		t.Fatal("could not convert to db backend")
	}
	defer b.Cleanup(ctx)

	// allow enough time for loadStaticWALs to work after boot
	time.Sleep(time.Second * 12)

	// loadStaticWALs should have processed the entry and removed the WAL log by now
	assertWALCount(t, config.StorageView, 0)

	// Read the role
	data = map[string]interface{}{}
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "static-roles/" + roleName,
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	lastVaultRotation := resp.Data["last_vault_rotation"].(time.Time)
	if !lastVaultRotation.After(oldRotationTime) {
		t.Fatal("last vault rotation time not greater than WAL time")
	}

	if !lastVaultRotation.After(roleTime) {
		t.Fatal("last vault rotation time not greater than role creation time")
	}

	// grab initial password to verify it doesn't change
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "static-creds/" + roleName,
		Storage:   config.StorageView,
	}
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	password := resp.Data["password"].(string)
	if password == walPassword {
		t.Fatalf("expected password to not be changed by WAL, but was")
	}
}

// helper to assert the number of WAL entries is what we expect
func assertWALCount(t *testing.T, s logical.Storage, expected int) {
	var count int
	ctx := context.Background()
	keys, err := framework.ListWAL(ctx, s)
	if err != nil {
		t.Fatal("error listing WALs")
	}

	// loop through WAL keys and process any rotation ones
	for _, k := range keys {
		walEntry, _ := framework.GetWAL(ctx, s, k)
		if walEntry == nil {
			continue
		}

		if walEntry.Kind != staticWALKey {
			continue
		}
		count++
	}
	if expected != count {
		t.Fatalf("WAL count mismatch, expected (%d), got (%d)", expected, count)
	}
}

//
// end WAL testing
//

// copied from plugins/database/mongodb/mongodb_test.go
// maybe move to a new builtin/logical/database/dbplugin/helper
func prepareMongoDBTestContainer(t *testing.T) (cleanup func(), retURL string) {
	if os.Getenv("MONGODB_URL") != "" {
		return func() {}, os.Getenv("MONGODB_URL")
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	resource, err := pool.Run("mongo", "latest", []string{})
	if err != nil {
		t.Fatalf("Could not start local mongo docker container: %s", err)
	}

	cleanup = func() {
		err := pool.Purge(resource)
		if err != nil {
			t.Fatalf("Failed to cleanup local container: %s", err)
		}
	}

	retURL = fmt.Sprintf("mongodb://localhost:%s", resource.GetPort("27017/tcp"))

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		var err error
		dialInfo, err := parseMongoURL(retURL)
		if err != nil {
			return err
		}

		session, err := mgo.DialWithInfo(dialInfo)
		if err != nil {
			return err
		}
		defer session.Close()
		session.SetSyncTimeout(1 * time.Minute)
		session.SetSocketTimeout(1 * time.Minute)
		return session.Ping()
	}); err != nil {
		cleanup()
		t.Fatalf("Could not connect to mongo docker container: %s", err)
	}

	return
}

// copied from plugins/database/mongodb/connection_producer.go
func parseMongoURL(rawURL string) (*mgo.DialInfo, error) {
	url, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	info := mgo.DialInfo{
		Addrs:    strings.Split(url.Host, ","),
		Database: strings.TrimPrefix(url.Path, "/"),
		Timeout:  10 * time.Second,
	}

	if url.User != nil {
		info.Username = url.User.Username()
		info.Password, _ = url.User.Password()
	}

	query := url.Query()
	for key, values := range query {
		var value string
		if len(values) > 0 {
			value = values[0]
		}

		switch key {
		case "authSource":
			info.Source = value
		case "authMechanism":
			info.Mechanism = value
		case "gssapiServiceName":
			info.Service = value
		case "replicaSet":
			info.ReplicaSetName = value
		case "maxPoolSize":
			poolLimit, err := strconv.Atoi(value)
			if err != nil {
				return nil, errors.New("bad value for maxPoolSize: " + value)
			}
			info.PoolLimit = poolLimit
		case "ssl":
			// Unfortunately, mgo doesn't support the ssl parameter in its MongoDB URI parsing logic, so we have to handle that
			// ourselves. See https://github.com/go-mgo/mgo/issues/84
			ssl, err := strconv.ParseBool(value)
			if err != nil {
				return nil, errors.New("bad value for ssl: " + value)
			}
			if ssl {
				info.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
					return tls.Dial("tcp", addr.String(), &tls.Config{})
				}
			}
		case "connect":
			if value == "direct" {
				info.Direct = true
				break
			}
			if value == "replicaSet" {
				break
			}
			fallthrough
		default:
			return nil, errors.New("unsupported connection URL option: " + key + "=" + value)
		}
	}

	return &info, nil
}
