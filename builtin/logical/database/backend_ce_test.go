// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package database

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/helper/namespace"
	postgreshelper "github.com/hashicorp/vault/helper/testhelpers/postgresql"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	_ "github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

// TestBackend_config_connection tests the configuration of a database connection
func TestBackend_config_connection(t *testing.T) {
	var resp *logical.Response
	var err error

	cluster, sys := getClusterPostgresDB(t)
	defer cluster.Cleanup()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys
	eventSender := logical.NewMockEventSender()
	config.EventsSender = eventSender
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
			"password_policy":                    "",
			"plugin_version":                     "",
			"verify_connection":                  false,
			"skip_static_role_import_rotation":   false,
			"rotation_schedule":                  "",
			"rotation_policy":                    "",
			"rotation_period":                    time.Duration(0).Seconds(),
			"rotation_window":                    time.Duration(0).Seconds(),
			"disable_automated_rotation":         false,
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
			"password_policy":                    "",
			"plugin_version":                     "",
			"verify_connection":                  false,
			"skip_static_role_import_rotation":   false,
			"rotation_schedule":                  "",
			"rotation_policy":                    "",
			"rotation_period":                    time.Duration(0).Seconds(),
			"rotation_window":                    time.Duration(0).Seconds(),
			"disable_automated_rotation":         false,
		}
		configReq.Operation = logical.ReadOperation
		resp, err = b.HandleRequest(namespace.RootContext(nil), configReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%s resp:%#v\n", err, resp)
		}

		delete(resp.Data["connection_details"].(map[string]interface{}), "name")
		delete(resp.Data, "AutomatedRotationParams")
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
			"password_policy":                    "",
			"plugin_version":                     "",
			"verify_connection":                  false,
			"skip_static_role_import_rotation":   false,
			"rotation_schedule":                  "",
			"rotation_policy":                    "",
			"rotation_period":                    time.Duration(0).Seconds(),
			"rotation_window":                    time.Duration(0).Seconds(),
			"disable_automated_rotation":         false,
		}
		configReq.Operation = logical.ReadOperation
		resp, err = b.HandleRequest(namespace.RootContext(nil), configReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%s resp:%#v\n", err, resp)
		}

		delete(resp.Data["connection_details"].(map[string]interface{}), "name")
		delete(resp.Data, "AutomatedRotationParams")
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
	assert.Equal(t, 3, len(eventSender.Events))
	assert.Equal(t, "database/config-write", string(eventSender.Events[0].Type))
	assert.Equal(t, "config/plugin-test", eventSender.Events[0].Event.Metadata.AsMap()["path"])
	assert.Equal(t, "plugin-test", eventSender.Events[0].Event.Metadata.AsMap()["name"])
	assert.Equal(t, "database/config-write", string(eventSender.Events[1].Type))
	assert.Equal(t, "config/plugin-test", eventSender.Events[1].Event.Metadata.AsMap()["path"])
	assert.Equal(t, "plugin-test", eventSender.Events[1].Event.Metadata.AsMap()["name"])
	assert.Equal(t, "database/config-write", string(eventSender.Events[2].Type))
	assert.Equal(t, "config/plugin-test", eventSender.Events[2].Event.Metadata.AsMap()["path"])
	assert.Equal(t, "plugin-test", eventSender.Events[2].Event.Metadata.AsMap()["name"])
}

// TestBackend_connectionCrud tests the full CRUD lifecycle of a database connection
func TestBackend_connectionCrud(t *testing.T) {
	t.Parallel()
	dbFactory := &singletonDBFactory{}
	cluster, sys := getClusterPostgresDBWithFactory(t, dbFactory.factory)
	defer cluster.Cleanup()

	dbFactory.sys = sys
	client := cluster.Cores[0].Client.Logical()

	cleanup, connURL := postgreshelper.PrepareTestContainer(t)
	defer cleanup()

	// Mount the database plugin.
	resp, err := client.Write("sys/mounts/database", map[string]interface{}{
		"type": "database",
	})
	if err != nil {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Configure a connection
	resp, err = client.Write("database/config/plugin-test", map[string]interface{}{
		"connection_url":    "test",
		"plugin_name":       "postgresql-database-plugin",
		"verify_connection": false,
	})
	if err != nil {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Configure a second connection to confirm below it doesn't get restarted.
	resp, err = client.Write("database/config/plugin-test-hana", map[string]interface{}{
		"connection_url":    "test",
		"plugin_name":       "hana-database-plugin",
		"verify_connection": false,
	})
	if err != nil {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Create a role
	resp, err = client.Write("database/roles/plugin-role-test", map[string]interface{}{
		"db_name":               "plugin-test",
		"creation_statements":   testRole,
		"revocation_statements": defaultRevocationSQL,
		"default_ttl":           "5m",
		"max_ttl":               "10m",
	})
	if err != nil {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Update the connection
	resp, err = client.Write("database/config/plugin-test", map[string]interface{}{
		"connection_url": connURL,
		"plugin_name":    "postgresql-database-plugin",
		"allowed_roles":  []string{"plugin-role-test"},
		"username":       "postgres",
		"password":       "secret",
	})
	if err != nil {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}
	if len(resp.Warnings) == 0 {
		t.Fatalf("expected warning about password in url %s, resp:%#v\n", connURL, resp)
	}

	resp, err = client.Read("database/config/plugin-test")
	if err != nil {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}
	returnedConnectionDetails := resp.Data["connection_details"].(map[string]interface{})
	if strings.Contains(returnedConnectionDetails["connection_url"].(string), "secret") {
		t.Fatal("password should not be found in the connection url")
	}
	// Covered by the filled out `expected` value below, but be explicit about this requirement.
	if _, exists := returnedConnectionDetails["password"]; exists {
		t.Fatal("password should NOT be found in the returned config")
	}

	// Replace connection url with templated version
	templatedConnURL := strings.ReplaceAll(connURL, "postgres:secret", "{{username}}:{{password}}")
	resp, err = client.Write("database/config/plugin-test", map[string]interface{}{
		"connection_url": templatedConnURL,
		"plugin_name":    "postgresql-database-plugin",
		"allowed_roles":  []string{"plugin-role-test"},
		"username":       "postgres",
		"password":       "secret",
	})
	if err != nil {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Read connection
	expected := map[string]interface{}{
		"plugin_name": "postgresql-database-plugin",
		"connection_details": map[string]interface{}{
			"username":       "postgres",
			"connection_url": templatedConnURL,
		},
		"allowed_roles":                      []any{"plugin-role-test"},
		"root_credentials_rotate_statements": []any{},
		"password_policy":                    "",
		"plugin_version":                     "",
		"verify_connection":                  false,
		"skip_static_role_import_rotation":   false,
		"rotation_schedule":                  "",
		"rotation_policy":                    "",
		"rotation_period":                    json.Number("0"),
		"rotation_window":                    json.Number("0"),
		"disable_automated_rotation":         false,
	}
	resp, err = client.Read("database/config/plugin-test")
	if err != nil {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	delete(resp.Data["connection_details"].(map[string]interface{}), "name")
	delete(resp.Data, "AutomatedRotationParams")
	if diff := deep.Equal(resp.Data, expected); diff != nil {
		t.Fatal(strings.Join(diff, "\n"))
	}

	// Test endpoints for reloading plugins.
	for _, reload := range []struct {
		path       string
		data       map[string]any
		checkCount bool
	}{
		{"database/reset/plugin-test", nil, false},
		{"database/reload/postgresql-database-plugin", nil, true},
		{"sys/plugins/reload/backend", map[string]any{
			"plugin": "postgresql-database-plugin",
		}, false},
	} {
		getConnectionID := func(name string) string {
			t.Helper()
			dbi := dbFactory.db.connections.Get(name)
			if dbi == nil {
				t.Fatal("no plugin-test dbi")
			}
			return dbi.ID()
		}
		initialID := getConnectionID("plugin-test")
		hanaID := getConnectionID("plugin-test-hana")
		resp, err = client.Write(reload.path, reload.data)
		if err != nil {
			t.Fatalf("err:%s resp:%#v\n", err, resp)
		}
		if initialID == getConnectionID("plugin-test") {
			t.Fatal("ID unchanged after connection reset")
		}
		if hanaID != getConnectionID("plugin-test-hana") {
			t.Fatal("hana plugin got restarted but shouldn't have been")
		}
		if reload.checkCount {
			actual, err := resp.Data["count"].(json.Number).Int64()
			if err != nil {
				t.Fatal(err)
			}
			if expected := 1; expected != int(actual) {
				t.Fatalf("expected %d but got %d", expected, resp.Data["count"].(int))
			}
			if expected := []any{"plugin-test"}; !reflect.DeepEqual(expected, resp.Data["connections"]) {
				t.Fatalf("expected %v but got %v", expected, resp.Data["connections"])
			}
		}
	}

	// Get creds
	credsResp, err := client.Read("database/creds/plugin-role-test")
	if err != nil {
		t.Fatalf("err:%s resp:%#v\n", err, credsResp)
	}

	credCheckURL := dbutil.QueryHelper(templatedConnURL, map[string]string{
		"username": "postgres",
		"password": "secret",
	})
	if !testCredsExist(t, credsResp.Data, credCheckURL) {
		t.Fatalf("Creds should exist")
	}

	// Delete Connection
	resp, err = client.Delete("database/config/plugin-test")
	if err != nil {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Read connection
	resp, err = client.Read("database/config/plugin-test")
	if err != nil {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Should be empty
	if resp != nil {
		t.Fatal("Expected response to be nil")
	}
}
