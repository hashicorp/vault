package database

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"strings"
	"testing"
)

func TestBackend_RotateRootCredentials_WAL(t *testing.T) {
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

	connectionName := "plugin-test"
	roleName := "plugin-role-test"
	userName := "postgres"
	oldPassword := "secret"
	newPassword := "newSecret"

	// Create a database connection to postgres
	cleanup, connURL := preparePostgresTestContainer(t, config.StorageView, lb)
	defer cleanup()

	connURL = strings.Replace(connURL, "postgres:secret", "{{username}}:{{password}}", -1)
	data := map[string]interface{}{
		"connection_url": connURL,
		"plugin_name":    "postgresql-database-plugin",
		"allowed_roles":  []string{roleName},
		"username":       userName,
		"password":       oldPassword,
	}
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      fmt.Sprintf("config/%s", connectionName),
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err := lb.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Create a role
	data = map[string]interface{}{
		"db_name":             connectionName,
		"creation_statements": testRole,
		"max_ttl":             "10m",
	}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      fmt.Sprintf("roles/%s", roleName),
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = lb.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Read database credentials
	credReq := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      fmt.Sprintf("creds/%s", roleName),
		Storage:   config.StorageView,
		Data:      make(map[string]interface{}),
	}
	credRes, err := lb.HandleRequest(namespace.RootContext(nil), credReq)
	if err != nil || (credRes != nil && credRes.IsError()) {
		t.Fatalf("err:%s resp:%v\n", err, credRes)
	}

	// Set database root credentials to "newSecret"
	pluginConn, err := dbBackend.GetConnection(context.Background(),
		config.StorageView, connectionName)
	if err != nil {
		t.Fatal(err)
	}
	pluginConn.SetCredentials(context.Background(), dbplugin.Statements{}, dbplugin.StaticUserConfig{
		Username: userName,
		Password: newPassword,
	})
	dbBackend.ClearConnection(connectionName)

	// Reading credentials should not work because the credentials changes
	credRes, err = lb.HandleRequest(namespace.RootContext(nil), credReq)
	if err == nil {
		t.Fatalf("err:%s resp:%v\n", err, credRes)
	}

	// Put a rotateRootCredentialsWAL
	walEntry := &rotateRootCredentialsWAL{
		ConnectionName: connectionName,
		UserName:       userName,
		OldPassword:    oldPassword,
		NewPassword:    newPassword,
	}
	framework.PutWAL(context.Background(), config.StorageView, rootWALKey, walEntry)

	// Trigger an immediate rollback operation
	_, err = lb.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.RollbackOperation,
		Path:      "",
		Storage:   config.StorageView,
		Data: map[string]interface{}{
			"immediate": true,
		},
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Getting credentials should work because of the WAL rollback
	credRes, err = lb.HandleRequest(namespace.RootContext(nil), credReq)
	if err != nil || (credRes != nil && credRes.IsError()) {
		t.Fatalf("err:%s resp:%v\n", err, credRes)
	}
}
