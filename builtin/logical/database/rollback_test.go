package database

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/hashicorp/vault/helper/namespace"
)

const (
	databaseUser    = "postgres"
	defaultPassword = "secret"
)

// Tests that the WAL rollback function rolls back the database password.
// The database password should be rolled back when:
//  - A WAL entry exists
//  - Password has been altered on the database
//  - Password has not been updated in storage
func TestBackend_RotateRootCredentials_WAL_rollback(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys

	lb, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	dbBackend, ok := lb.(*databaseBackend)
	if !ok {
		t.Fatal("could not convert to db backend")
	}
	defer lb.Cleanup(context.Background())

	cleanup, connURL := preparePostgresTestContainer(t, config.StorageView, lb)
	defer cleanup()

	connURL = strings.Replace(connURL, "postgres:secret", "{{username}}:{{password}}", -1)

	// Configure a connection to the database
	data := map[string]interface{}{
		"connection_url": connURL,
		"plugin_name":    "postgresql-database-plugin",
		"allowed_roles":  []string{"plugin-role-test"},
		"username":       databaseUser,
		"password":       defaultPassword,
	}
	resp, err := lb.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Create a role
	data = map[string]interface{}{
		"db_name":             "plugin-test",
		"creation_statements": testRole,
		"max_ttl":             "10m",
	}
	resp, err = lb.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/plugin-role-test",
		Storage:   config.StorageView,
		Data:      data,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Read credentials to verify this initially works
	credReq := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/plugin-role-test",
		Storage:   config.StorageView,
		Data:      make(map[string]interface{}),
	}
	credResp, err := lb.HandleRequest(context.Background(), credReq)
	if err != nil || (credResp != nil && credResp.IsError()) {
		t.Fatalf("err:%s resp:%v\n", err, credResp)
	}

	// Get a connection to the database plugin
	pc, err := dbBackend.GetConnection(context.Background(),
		config.StorageView, "plugin-test")
	if err != nil {
		t.Fatal(err)
	}

	// Alter the database password so it no longer matches what is in storage
	err = changeUserPassword(context.Background(), pc.database, databaseUser, "newSecret", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Clear the plugin connection to verify we're no longer able to connect
	err = dbBackend.ClearConnection("plugin-test")
	if err != nil {
		t.Fatal(err)
	}

	// Reading credentials should no longer work
	credResp, err = lb.HandleRequest(namespace.RootContext(nil), credReq)
	if err == nil {
		t.Fatalf("expected authentication to fail when reading credentials")
	}

	// Put a WAL entry that will be used for rolling back the database password
	walEntry := &rotateRootCredentialsWAL{
		ConnectionName: "plugin-test",
		UserName:       databaseUser,
		OldPassword:    defaultPassword,
		NewPassword:    "newSecret",
	}
	_, err = framework.PutWAL(context.Background(), config.StorageView, rotateRootWALKey, walEntry)
	if err != nil {
		t.Fatal(err)
	}
	assertWALCount(t, config.StorageView, 1, rotateRootWALKey)

	// Trigger an immediate RollbackOperation so that the WAL rollback
	// function can use the WAL entry to roll back the database password
	_, err = lb.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.RollbackOperation,
		Path:      "",
		Storage:   config.StorageView,
		Data: map[string]interface{}{
			"immediate": true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertWALCount(t, config.StorageView, 0, rotateRootWALKey)

	// Reading credentials should work again after the database
	// password has been rolled back.
	credResp, err = lb.HandleRequest(namespace.RootContext(nil), credReq)
	if err != nil || (credResp != nil && credResp.IsError()) {
		t.Fatalf("err:%s resp:%v\n", err, credResp)
	}
}

// Tests that the WAL rollback function does not roll back the database password.
// The database password should not be rolled back when:
//  - A WAL entry exists
//  - Password has not been altered on the database
//  - Password has not been updated in storage
func TestBackend_RotateRootCredentials_WAL_no_rollback_1(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys

	lb, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	defer lb.Cleanup(context.Background())

	cleanup, connURL := preparePostgresTestContainer(t, config.StorageView, lb)
	defer cleanup()

	connURL = strings.Replace(connURL, "postgres:secret", "{{username}}:{{password}}", -1)

	// Configure a connection to the database
	data := map[string]interface{}{
		"connection_url": connURL,
		"plugin_name":    "postgresql-database-plugin",
		"allowed_roles":  []string{"plugin-role-test"},
		"username":       databaseUser,
		"password":       defaultPassword,
	}
	resp, err := lb.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Create a role
	data = map[string]interface{}{
		"db_name":             "plugin-test",
		"creation_statements": testRole,
		"max_ttl":             "10m",
	}
	resp, err = lb.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/plugin-role-test",
		Storage:   config.StorageView,
		Data:      data,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Read credentials to verify this initially works
	credReq := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/plugin-role-test",
		Storage:   config.StorageView,
		Data:      make(map[string]interface{}),
	}
	credResp, err := lb.HandleRequest(context.Background(), credReq)
	if err != nil || (credResp != nil && credResp.IsError()) {
		t.Fatalf("err:%s resp:%v\n", err, credResp)
	}

	// Put a WAL entry
	walEntry := &rotateRootCredentialsWAL{
		ConnectionName: "plugin-test",
		UserName:       databaseUser,
		OldPassword:    defaultPassword,
		NewPassword:    "newSecret",
	}
	_, err = framework.PutWAL(context.Background(), config.StorageView, rotateRootWALKey, walEntry)
	if err != nil {
		t.Fatal(err)
	}
	assertWALCount(t, config.StorageView, 1, rotateRootWALKey)

	// Trigger an immediate RollbackOperation
	_, err = lb.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.RollbackOperation,
		Path:      "",
		Storage:   config.StorageView,
		Data: map[string]interface{}{
			"immediate": true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertWALCount(t, config.StorageView, 0, rotateRootWALKey)

	// Reading credentials should work
	credResp, err = lb.HandleRequest(namespace.RootContext(nil), credReq)
	if err != nil || (credResp != nil && credResp.IsError()) {
		t.Fatalf("err:%s resp:%v\n", err, credResp)
	}
}

// Tests that the WAL rollback function does not roll back the database password.
// The database password should not be rolled back when:
//  - A WAL entry exists
//  - Password has been altered on the database
//  - Password has been updated in storage
func TestBackend_RotateRootCredentials_WAL_no_rollback_2(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys

	lb, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	dbBackend, ok := lb.(*databaseBackend)
	if !ok {
		t.Fatal("could not convert to db backend")
	}
	defer lb.Cleanup(context.Background())

	cleanup, connURL := preparePostgresTestContainer(t, config.StorageView, lb)
	defer cleanup()

	connURL = strings.Replace(connURL, "postgres:secret", "{{username}}:{{password}}", -1)

	// Configure a connection to the database
	data := map[string]interface{}{
		"connection_url": connURL,
		"plugin_name":    "postgresql-database-plugin",
		"allowed_roles":  []string{"plugin-role-test"},
		"username":       databaseUser,
		"password":       defaultPassword,
	}
	resp, err := lb.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Create a role
	data = map[string]interface{}{
		"db_name":             "plugin-test",
		"creation_statements": testRole,
		"max_ttl":             "10m",
	}
	resp, err = lb.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/plugin-role-test",
		Storage:   config.StorageView,
		Data:      data,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Read credentials to verify this initially works
	credReq := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/plugin-role-test",
		Storage:   config.StorageView,
		Data:      make(map[string]interface{}),
	}
	credResp, err := lb.HandleRequest(context.Background(), credReq)
	if err != nil || (credResp != nil && credResp.IsError()) {
		t.Fatalf("err:%s resp:%v\n", err, credResp)
	}

	// Get a connection to the database plugin
	pc, err := dbBackend.GetConnection(context.Background(), config.StorageView, "plugin-test")
	if err != nil {
		t.Fatal(err)
	}

	// Alter the database password
	err = changeUserPassword(context.Background(), pc.database, databaseUser, "newSecret", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Update storage with the new password
	dbConfig, err := dbBackend.DatabaseConfig(context.Background(), config.StorageView,
		"plugin-test")
	if err != nil {
		t.Fatal(err)
	}
	dbConfig.ConnectionDetails["password"] = "newSecret"
	entry, err := logical.StorageEntryJSON("config/plugin-test", dbConfig)
	if err != nil {
		t.Fatal(err)
	}
	err = config.StorageView.Put(context.Background(), entry)
	if err != nil {
		t.Fatal(err)
	}

	// Clear the plugin connection to verify we can connect to the database
	err = dbBackend.ClearConnection("plugin-test")
	if err != nil {
		t.Fatal(err)
	}

	// Reading credentials should work
	credResp, err = lb.HandleRequest(namespace.RootContext(nil), credReq)
	if err != nil || (credResp != nil && credResp.IsError()) {
		t.Fatalf("err:%s resp:%v\n", err, credResp)
	}

	// Put a WAL entry
	walEntry := &rotateRootCredentialsWAL{
		ConnectionName: "plugin-test",
		UserName:       databaseUser,
		OldPassword:    defaultPassword,
		NewPassword:    "newSecret",
	}
	_, err = framework.PutWAL(context.Background(), config.StorageView, rotateRootWALKey, walEntry)
	if err != nil {
		t.Fatal(err)
	}
	assertWALCount(t, config.StorageView, 1, rotateRootWALKey)

	// Trigger an immediate RollbackOperation
	_, err = lb.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.RollbackOperation,
		Path:      "",
		Storage:   config.StorageView,
		Data: map[string]interface{}{
			"immediate": true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertWALCount(t, config.StorageView, 0, rotateRootWALKey)

	// Reading credentials should work
	credResp, err = lb.HandleRequest(namespace.RootContext(nil), credReq)
	if err != nil || (credResp != nil && credResp.IsError()) {
		t.Fatalf("err:%s resp:%v\n", err, credResp)
	}
}
