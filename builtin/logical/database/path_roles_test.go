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

func TestBackend_Static_Config(t *testing.T) {
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
        // connURL := "postgres://postgres:secret@localhost:32768/database?sslmode=disable"

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
                                "username":           "sa-test",
                                "rotation_frequency": "5400s",
                        },
                        expected: map[string]interface{}{
                                "username":           "sa-test",
                                "rotation_frequency": int64(5400000000000),
                        },
                },
                "missing rotation frequency": {
                        account: map[string]interface{}{
                                "username": "sa-test",
                        },
                        err: errors.New("rotation_frequency is a required field for static accounts"),
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
                                "username":           "sa-test",
                                "rotation_frequency": "5400s",
                                "password":           "somesecret123!!",
                        },
                        expected: map[string]interface{}{
                                "username":           "sa-test",
                                "rotation_frequency": int64(5400000000000),
                                "password":           "somesecret123!!",
                        },
                },
        }

        for name, tc := range testCases {
                t.Run(name, func(t *testing.T) {
                        data := map[string]interface{}{
                                "name":                  "plugin-role-test",
                                "db_name":               "plugin-test",
                                "creation_statements":   testStaticRoleCreate,
                                "revocation_statements": defaultRevocationSQL,
                                "default_ttl":           "5m",
                                "max_ttl":               "10m",
                                "static_account":        tc.account,
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
                        actual := resp.Data["static_account"]

                        if len(tc.expected) > 0 {
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

// const testStaticRole = `
// DO
// $do$
// BEGIN
//    IF NOT EXISTS (
//       SELECT
//       FROM   pg_catalog.pg_roles
//       WHERE  rolname = '{{name}}') THEN

//       CREATE ROLE "{{name}}" WITH LOGIN PASSWORD '{{password}}';
//    END IF;
// END
// $do$;
// `
var testStaticRoleCreate = []string{
        `
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = '{{name}}') THEN
        CREATE ROLE "{{name}}";
    END IF;
END
$$;
`,
        `
ALTER USER "{{name}}" WITH PASSWORD '{{password}}';
`,
}

// const testStaticRoleCreate = `
// CREATE ROLE "{{name}}" WITH
//   LOGIN
//   PASSWORD '{{password}}';
// GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "{{name}}";
// `

// const testStaticRoleUpdate = `
// ALTER USER "{{name}}" WITH PASSWORD '{{password}}';
// `
