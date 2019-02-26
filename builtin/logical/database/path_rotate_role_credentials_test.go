package database

import (
        "context"
        "strings"
        "testing"

        "github.com/hashicorp/vault/helper/namespace"
        "github.com/hashicorp/vault/logical"

        "database/sql"

        _ "github.com/lib/pq"
)

func TestBackend_Static_Account_Rotate(t *testing.T) {
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

        data = map[string]interface{}{
                "name":                  "plugin-role-test",
                "db_name":               "plugin-test",
                "creation_statements":   testRoleStaticCreate,
                "rotation_statements":   testRoleStaticUpdate,
                "revocation_statements": defaultRevocationSQL,
                "username":              "statictest",
                "rotation_frequency":    "5400s",
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

        // Read the creds
        data = map[string]interface{}{}
        req = &logical.Request{
                Operation: logical.ReadOperation,
                Path:      "creds/plugin-role-test",
                Storage:   config.StorageView,
                Data:      data,
        }
        resp, err = b.HandleRequest(namespace.RootContext(nil), req)
        if err != nil || (resp != nil && resp.IsError()) {
                t.Fatalf("err:%s resp:%#v\n", err, resp)
        }

        username := resp.Data["username"].(string)
        password := resp.Data["password"].(string)
        if username == "" || password == "" {
                t.Fatalf("empty username (%s) or password (%s)", username, password)
        }

        // verify username/password
        verifyPgConn(t, username, password, connURL)

        // trigger rotation
        data = map[string]interface{}{"name": "plugin-role-test"}
        req = &logical.Request{
                Operation: logical.UpdateOperation,
                Path:      "rotate-role/plugin-role-test",
                Storage:   config.StorageView,
                Data:      data,
        }
        resp, err = b.HandleRequest(namespace.RootContext(nil), req)
        if err != nil || (resp != nil && resp.IsError()) {
                t.Fatalf("err:%s resp:%#v\n", err, resp)
        }

        if resp != nil {
                t.Fatalf("Expected empty response from rotate-role: (%#v)", resp)
        }

        // Re-Read the creds
        data = map[string]interface{}{}
        req = &logical.Request{
                Operation: logical.ReadOperation,
                Path:      "creds/plugin-role-test",
                Storage:   config.StorageView,
                Data:      data,
        }
        resp, err = b.HandleRequest(namespace.RootContext(nil), req)
        if err != nil || (resp != nil && resp.IsError()) {
                t.Fatalf("err:%s resp:%#v\n", err, resp)
        }

        newPassword := resp.Data["password"].(string)
        if password == newPassword {
                t.Fatalf("expected passwords to differ, got (%s)", newPassword)
        }

        // verify new username/password
        verifyPgConn(t, username, newPassword, connURL)
}

func verifyPgConn(t *testing.T, username, password, connURL string) {
        cURL := strings.Replace(connURL, "postgres:secret", username+":"+password, 1)
        db, err := sql.Open("postgres", cURL)
        if err != nil {
                t.Fatal(err)
        }
        if err := db.Ping(); err != nil {
                t.Fatal(err)
        }
        db.Close()
}
