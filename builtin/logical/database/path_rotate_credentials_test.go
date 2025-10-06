// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package database

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/namespace"
	postgreshelper "github.com/hashicorp/vault/helper/testhelpers/postgresql"
	"github.com/hashicorp/vault/sdk/logical"
)

// TestRotateRootPassword_Postgres tests database secrets root rotation for password credentials with Postgres.
func TestRotateRootPassword_Postgres(t *testing.T) {
	cleanup, connURL := postgreshelper.PrepareTestContainer(t)
	defer cleanup()

	cluster, sys := getClusterWithFactory(t, Factory)
	defer cluster.Cleanup()

	connURL = strings.Replace(connURL, "postgres:secret", "{{username}}:{{password}}", 1)

	testRotateRoot(t, sys, true, connURL, "postgresql-database-plugin", "postgres", "secret")
}

// TestRotateRootKeypair_Snowflake_Acc tests database secrets root rotation for private key credentials with Snowflake.
func TestRotateRootKeypair_Snowflake_Acc(t *testing.T) {
	// SNOWFLAKE_PRIVATE_KEY is the path to the private key file.
	privateKeyPath, ok := os.LookupEnv("VAULT_SNOWFLAKE_PRIVATE_KEY")
	if !ok {
		t.Skip("VAULT_SNOWFLAKE_PRIVATE_KEY not set, skipping test")
	}

	keyFile, err := os.ReadFile(privateKeyPath)
	if err != nil {
		t.Fatalf("failed to read private key file: %s", err)
	}

	connURL, ok := os.LookupEnv("VAULT_SNOWFLAKE_CONNECTION_URL")
	if !ok {
		t.Skip("VAULT_SNOWFLAKE_CONNECTION_URL not set, skipping test")
	}

	username, ok := os.LookupEnv("VAULT_SNOWFLAKE_USERNAME")
	if !ok {
		t.Skip("VAULT_SNOWFLAKE_USERNAME not set, skipping test")
	}

	cluster, sys := getClusterWithFactory(t, Factory)
	defer cluster.Cleanup()

	testRotateRoot(t, sys, false, connURL, "snowflake-database-plugin", username, string(keyFile))
}

// Helper function to run rotate root tests.
func testRotateRoot(t *testing.T, sys logical.SystemView, isPassword bool, connURL string, pluginName, username, credential string) {
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

	// Configure a connection
	configData := map[string]any{
		"connection_url":    connURL,
		"plugin_name":       pluginName,
		"verify_connection": false,
		"allowed_roles":     []string{"*"},
		"name":              "plugin-test",
		"username":          username,
	}
	if isPassword {
		configData["password"] = credential
	} else {
		configData["private_key"] = credential
	}
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      configData,
	}
	resp, err := b.HandleRequest(namespace.RootContext(context.Background()), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Create a dynamic role
	req = &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "roles/test",
		Storage:   config.StorageView,
		Data: map[string]interface{}{
			"db_name":             "plugin-test",
			"creation_statements": `CREATE USER "{{name}}"`,
		},
	}
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Read some creds to validate the connection
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/test",
		Storage:   config.StorageView,
	}
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Rotate the root credentials
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "rotate-root/plugin-test",
		Storage:   config.StorageView,
	}
	resp, err = b.HandleRequest(namespace.RootContext(context.Background()), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Read creds a second time to validate the root rotation still works
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/test",
		Storage:   config.StorageView,
	}
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Write back the original credential to ensure it no longer works
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      configData,
	}
	resp, err = b.HandleRequest(namespace.RootContext(context.Background()), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Read some more creds again and expect an error]
	// Note: For Snowflake, this step may fail if you are using the account's 2nd private key
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/test",
		Storage:   config.StorageView,
	}
	resp, err = b.HandleRequest(context.Background(), req)
	if err == nil || (resp != nil && !resp.IsError()) {
		t.Fatalf("expected authentication error but did not receive an error, resp:%#v\n", resp)
	}
}
