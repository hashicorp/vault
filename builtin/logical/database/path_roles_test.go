package database

import (
        "context"
        "errors"
        "testing"

        "github.com/go-test/deep"
        "github.com/hashicorp/vault/helper/namespace"
        "github.com/hashicorp/vault/logical"
        "github.com/hashicorp/vault/logical/framework"
)

var dataKeys = []string{"username", "password", "last_vault_rotation", "rotation_frequency"}

func TestBackend_StaticRole_Config(t *testing.T) {
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

        // Test static role creation scenarios. Uses a map, so there is no guaranteed
        // ordering, so each case cleans up by deleting the role
        testCases := map[string]struct {
                account  map[string]interface{}
                expected map[string]interface{}
                err      error
        }{
                "normal": {},
                "basic": {
                        account: map[string]interface{}{
                                "username":           "statictest",
                                "rotation_frequency": "5400s",
                        },
                        expected: map[string]interface{}{
                                "username":           "statictest",
                                "rotation_frequency": int64(5400000000000),
                        },
                },
                "missing rotation frequency": {
                        account: map[string]interface{}{
                                "username": "statictest",
                        },
                        err: errors.New("rotation_frequency is required to create static accounts"),
                },
                "missing username frequency": {
                        account: map[string]interface{}{
                                "rotation_frequency": int64(5400000000000),
                        },
                        err: errors.New("username is a required field for static accounts"),
                },
                "missing all": {
                        account: map[string]interface{}{"fill": "stuff"},
                        err:     errors.New("username is a required field for static accounts"),
                },
                "with password": {
                        account: map[string]interface{}{
                                "username":           "statictest",
                                "rotation_frequency": "5400s",
                        },
                        expected: map[string]interface{}{
                                "username":           "statictest",
                                "rotation_frequency": int64(5400000000000),
                        },
                },
        }

        for name, tc := range testCases {
                t.Run(name, func(t *testing.T) {
                        data := map[string]interface{}{
                                "name":                  "plugin-role-test",
                                "db_name":               "plugin-test",
                                "creation_statements":   testRoleStaticCreate,
                                "rotation_statements":   testRoleStaticUpdate,
                                "revocation_statements": defaultRevocationSQL,
                                "default_ttl":           "5m",
                                "max_ttl":               "10m",
                        }

                        for k, v := range tc.account {
                                data[k] = v
                        }

                        req := &logical.Request{
                                Operation: logical.CreateOperation,
                                Path:      "roles/plugin-role-test",
                                Storage:   config.StorageView,
                                Data:      data,
                        }

                        exists, err := b.pathRoleExistenceCheck()(context.Background(), req, &framework.FieldData{
                                Raw:    data,
                                Schema: pathRoles(b).Fields,
                        })
                        if err != nil {
                                t.Fatal(err)
                        }
                        if exists {
                                t.Fatal("expected not exists")
                        }

                        resp, err = b.HandleRequest(namespace.RootContext(nil), req)
                        if err != nil || (resp != nil && resp.IsError()) {
                                if tc.err == nil {
                                        t.Fatalf("err:%s resp:%#v\n", err, resp)
                                }
                                if err != nil && tc.err.Error() == err.Error() {
                                        // errors match
                                        return
                                }
                                if err == nil && tc.err.Error() == resp.Error().Error() {
                                        // errors match
                                        return
                                }
                                t.Fatalf("expected err message: (%s), got (%s), response error: (%s)", tc.err, err, resp.Error())
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

                        expected := tc.expected
                        actual := make(map[string]interface{})
                        for _, key := range dataKeys {
                                if v, ok := resp.Data[key]; ok {
                                        actual[key] = v
                                }
                        }

                        if len(tc.expected) > 0 {
                                // verify a password is returned, but we don't care what it's value is
                                if actual["password"] == "" {
                                        t.Fatalf("expected result to contain password, but none found")
                                }
                                delete(actual, "password")
                                if diff := deep.Equal(expected, actual); diff != nil {
                                        t.Fatal(diff)
                                }
                        }

                        if len(tc.expected) == 0 && resp.Data["static_account"] != nil {
                                t.Fatalf("got unexpected static_account info: %#v", actual)
                        }

                        if diff := deep.Equal(resp.Data["db_name"], "plugin-test"); diff != nil {
                                t.Fatal(diff)
                        }

                        // Delete role for next run
                        req = &logical.Request{
                                Operation: logical.DeleteOperation,
                                Path:      "roles/plugin-role-test",
                                Storage:   config.StorageView,
                        }
                        resp, err = b.HandleRequest(namespace.RootContext(nil), req)
                        if err != nil || (resp != nil && resp.IsError()) {
                                t.Fatalf("err:%s resp:%#v\n", err, resp)
                        }
                })
        }
}

func TestBackend_StaticRole_Config_Update(t *testing.T) {
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
                "default_ttl":           "5m",
                "max_ttl":               "10m",
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

        rotation := resp.Data["rotation_frequency"].(int64)
        // update rotation_frequency
        updateData := map[string]interface{}{
                "name":               "plugin-role-test",
                "db_name":            "plugin-test",
                "username":           "statictest",
                "rotation_frequency": "6400s",
        }
        req = &logical.Request{
                Operation: logical.UpdateOperation,
                Path:      "roles/plugin-role-test",
                Storage:   config.StorageView,
                Data:      updateData,
        }

        resp, err = b.HandleRequest(namespace.RootContext(nil), req)
        if err != nil || (resp != nil && resp.IsError()) {
                t.Fatalf("err:%s resp:%#v\n", err, resp)
        }

        // re-read the role
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

        newRotation := resp.Data["rotation_frequency"].(int64)
        if newRotation == rotation {
                t.Fatalf("expected change in rotation, but got old value:  %#v", newRotation)
        }

        // verify that rotation_frequency is only required when creating
        updateData = map[string]interface{}{
                "name":                "plugin-role-test",
                "db_name":             "plugin-test",
                "username":            "statictest",
                "rotation_statements": testRoleStaticUpdateRotation,
        }
        req = &logical.Request{
                Operation: logical.UpdateOperation,
                Path:      "roles/plugin-role-test",
                Storage:   config.StorageView,
                Data:      updateData,
        }

        resp, err = b.HandleRequest(namespace.RootContext(nil), req)
        if err != nil || (resp != nil && resp.IsError()) {
                t.Fatalf("err:%s resp:%#v\n", err, resp)
        }

        // verify updating static username returns an error
        updateData = map[string]interface{}{
                "name":     "plugin-role-test",
                "db_name":  "plugin-test",
                "username": "statictestmodified",
        }
        req = &logical.Request{
                Operation: logical.UpdateOperation,
                Path:      "roles/plugin-role-test",
                Storage:   config.StorageView,
                Data:      updateData,
        }

        resp, err = b.HandleRequest(namespace.RootContext(nil), req)
        if err != nil || !resp.IsError() {
                t.Fatal("expected error on updating name")
        }
        err = resp.Error()
        if err.Error() != "cannot update static account username" {
                t.Fatalf("expected error on updating name, got: %s", err)
        }
}

const testRoleStaticCreate = `
CREATE ROLE "{{name}}" WITH
  LOGIN
  PASSWORD '{{password}}';
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "{{name}}";
`

const testRoleStaticUpdate = `
ALTER USER "{{name}}" WITH PASSWORD '{{password}}';
`

const testRoleStaticUpdateRotation = `
ALTER USER "{{name}}" WITH PASSWORD '{{password}}';GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "{{name}}";
`
