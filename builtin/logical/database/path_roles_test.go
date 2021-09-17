package database

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/helper/namespace"
	postgreshelper "github.com/hashicorp/vault/helper/testhelpers/postgresql"
	"github.com/hashicorp/vault/sdk/logical"
)

var dataKeys = []string{"username", "password", "last_vault_rotation", "rotation_period"}

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

	cleanup, connURL := postgreshelper.PrepareTestContainer(t, "")
	defer cleanup()

	// create the database user
	createTestPGUser(t, connURL, dbUser, "password", testRoleStaticCreate)

	// Configure a connection
	data := map[string]interface{}{
		"connection_url":    connURL,
		"plugin_name":       "postgresql-database-plugin",
		"verify_connection": false,
		"allowed_roles":     []string{"plugin-role-test"},
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
		path     string
		expected map[string]interface{}
		err      error
	}{
		"basic": {
			account: map[string]interface{}{
				"username":        dbUser,
				"rotation_period": "5400s",
			},
			path: "plugin-role-test",
			expected: map[string]interface{}{
				"username":        dbUser,
				"rotation_period": float64(5400),
			},
		},
		"missing rotation period": {
			account: map[string]interface{}{
				"username": dbUser,
			},
			path: "plugin-role-test",
			err:  errors.New("rotation_period is required to create static accounts"),
		},
		"disallowed role config": {
			account: map[string]interface{}{
				"username":        dbUser,
				"rotation_period": "5400s",
			},
			path: "disallowed-role",
			err:  errors.New("\"disallowed-role\" is not an allowed role"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			data := map[string]interface{}{
				"name":                "plugin-role-test",
				"db_name":             "plugin-test",
				"rotation_statements": testRoleStaticUpdate,
			}

			for k, v := range tc.account {
				data[k] = v
			}

			path := "static-roles/" + tc.path

			req := &logical.Request{
				Operation: logical.CreateOperation,
				Path:      path,
				Storage:   config.StorageView,
				Data:      data,
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

			if tc.err != nil {
				if err == nil || (resp == nil || !resp.IsError()) {
					t.Fatal("expected error, got none")
				}
			}

			// Read the role
			data = map[string]interface{}{}
			req = &logical.Request{
				Operation: logical.ReadOperation,
				Path:      "static-roles/plugin-role-test",
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
				if v, ok := actual["last_vault_rotation"].(time.Time); !ok {
					t.Fatalf("expected last_vault_rotation to be set to time.Time type, got: %#v", v)
				}

				// delete these values before the comparison, since we can't know them in
				// advance
				delete(actual, "password")
				delete(actual, "last_vault_rotation")
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
				Path:      "static-roles/plugin-role-test",
				Storage:   config.StorageView,
			}
			resp, err = b.HandleRequest(namespace.RootContext(nil), req)
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("err:%s resp:%#v\n", err, resp)
			}
		})
	}
}

func TestBackend_StaticRole_Updates(t *testing.T) {
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

	cleanup, connURL := postgreshelper.PrepareTestContainer(t, "")
	defer cleanup()

	// create the database user
	createTestPGUser(t, connURL, dbUser, "password", testRoleStaticCreate)

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
		"name":                "plugin-role-test-updates",
		"db_name":             "plugin-test",
		"rotation_statements": testRoleStaticUpdate,
		"username":            dbUser,
		"rotation_period":     "5400s",
	}

	req = &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "static-roles/plugin-role-test-updates",
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
		Path:      "static-roles/plugin-role-test-updates",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	rotation := resp.Data["rotation_period"].(float64)

	// capture the password to verify it doesn't change
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "static-creds/plugin-role-test-updates",
		Storage:   config.StorageView,
	}
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	username := resp.Data["username"].(string)
	password := resp.Data["password"].(string)
	if username == "" || password == "" {
		t.Fatalf("expected both username/password, got (%s), (%s)", username, password)
	}

	// update rotation_period
	updateData := map[string]interface{}{
		"name":            "plugin-role-test-updates",
		"db_name":         "plugin-test",
		"username":        dbUser,
		"rotation_period": "6400s",
	}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "static-roles/plugin-role-test-updates",
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
		Path:      "static-roles/plugin-role-test-updates",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	newRotation := resp.Data["rotation_period"].(float64)
	if newRotation == rotation {
		t.Fatalf("expected change in rotation, but got old value:  %#v", newRotation)
	}

	// re-capture the password to ensure it did not change
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "static-creds/plugin-role-test-updates",
		Storage:   config.StorageView,
	}
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	if username != resp.Data["username"].(string) {
		t.Fatalf("usernames dont match!: (%s) / (%s)", username, resp.Data["username"].(string))
	}
	if password != resp.Data["password"].(string) {
		t.Fatalf("passwords dont match!: (%s) / (%s)", password, resp.Data["password"].(string))
	}

	// verify that rotation_period is only required when creating
	updateData = map[string]interface{}{
		"name":                "plugin-role-test-updates",
		"db_name":             "plugin-test",
		"username":            dbUser,
		"rotation_statements": testRoleStaticUpdateRotation,
	}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "static-roles/plugin-role-test-updates",
		Storage:   config.StorageView,
		Data:      updateData,
	}

	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// verify updating static username returns an error
	updateData = map[string]interface{}{
		"name":     "plugin-role-test-updates",
		"db_name":  "plugin-test",
		"username": "statictestmodified",
	}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "static-roles/plugin-role-test-updates",
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

func TestBackend_StaticRole_Role_name_check(t *testing.T) {
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

	cleanup, connURL := postgreshelper.PrepareTestContainer(t, "")
	defer cleanup()

	// create the database user
	createTestPGUser(t, connURL, dbUser, "password", testRoleStaticCreate)

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

	// non-static role
	data = map[string]interface{}{
		"name":                  "plugin-role-test",
		"db_name":               "plugin-test",
		"creation_statements":   testRoleStaticCreate,
		"rotation_statements":   testRoleStaticUpdate,
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

	// create a static role with the same name, and expect failure
	// static role
	data = map[string]interface{}{
		"name":                  "plugin-role-test",
		"db_name":               "plugin-test",
		"creation_statements":   testRoleStaticCreate,
		"rotation_statements":   testRoleStaticUpdate,
		"revocation_statements": defaultRevocationSQL,
	}

	req = &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "static-roles/plugin-role-test",
		Storage:   config.StorageView,
		Data:      data,
	}

	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || !resp.IsError() {
		t.Fatalf("expected error, got none")
	}

	// repeat, with a static role first
	data = map[string]interface{}{
		"name":                "plugin-role-test-2",
		"db_name":             "plugin-test",
		"rotation_statements": testRoleStaticUpdate,
		"username":            dbUser,
		"rotation_period":     "1h",
	}

	req = &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "static-roles/plugin-role-test-2",
		Storage:   config.StorageView,
		Data:      data,
	}

	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// create a non-static role with the same name, and expect failure
	data = map[string]interface{}{
		"name":                  "plugin-role-test-2",
		"db_name":               "plugin-test",
		"creation_statements":   testRoleStaticCreate,
		"revocation_statements": defaultRevocationSQL,
		"default_ttl":           "5m",
		"max_ttl":               "10m",
	}

	req = &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "roles/plugin-role-test-2",
		Storage:   config.StorageView,
		Data:      data,
	}

	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || !resp.IsError() {
		t.Fatalf("expected error, got none")
	}
}

// func TestWALsStillTrackedAfterUpdate(t *testing.T) {
// 	ctx := context.Background()
// 	b, storage := getBackend(t, false)
// 	defer b.Cleanup(ctx)
// 	configureDBMount(t, b, storage)

// 	createRole(t, b, storage, "hashicorp")

// 	generateWALFromFailedRotation(t, b, storage, "hashicorp")
// 	requireWALs(t, storage, 1)

// 	_, err := b.HandleRequest(ctx, &logical.Request{
// 		Operation: logical.UpdateOperation,
// 		Path:      "static-roles/hashicorp",
// 		Storage:   storage,
// 		Data: map[string]interface{}{
// 			"username":        "hashicorp",
// 			"dn":              "uid=hashicorp,ou=users,dc=hashicorp,dc=com",
// 			"rotation_period": "600s",
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	requireWALs(t, storage, 1)

// 	// Check we've still got track of it in the queue as well
// 	item, err := b.credRotationQueue.PopByKey("hashicorp")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if wal, ok := item.Value.(string); !ok || wal == "" {
// 		t.Fatal("should have a WAL ID in the rotation queue")
// 	}
// }

// func TestWALsDeletedOnRoleCreationFailed(t *testing.T) {
// 	ctx := context.Background()
// 	b, storage := getBackend(t, true)
// 	defer b.Cleanup(ctx)
// 	configureDBMount(t, b, storage)

// 	for i := 0; i < 3; i++ {
// 		_, err := b.HandleRequest(ctx, &logical.Request{
// 			Operation: logical.CreateOperation,
// 			Path:      "static-roles/hashicorp",
// 			Storage:   storage,
// 			Data: map[string]interface{}{
// 				"username":        "hashicorp",
// 				"dn":              "uid=hashicorp,ou=users,dc=hashicorp,dc=com",
// 				"rotation_period": "5s",
// 			},
// 		})
// 		if err == nil {
// 			t.Fatal("expected error from OpenLDAP")
// 		}
// 	}

// 	requireWALs(t, storage, 0)
// }

// func TestWALsDeletedOnRoleDeletion(t *testing.T) {
// 	ctx := context.Background()
// 	b, storage := getBackend(t, false)
// 	defer b.Cleanup(ctx)
// 	configureDBMount(t, b, storage)

// 	// Create the roles
// 	roleNames := []string{"hashicorp", "2"}
// 	for _, roleName := range roleNames {
// 		createRole(t, b, storage, roleName)
// 	}

// 	// Fail to rotate the roles
// 	for _, roleName := range roleNames {
// 		generateWALFromFailedRotation(t, b, storage, roleName)
// 	}

// 	// Should have 2 WALs hanging around
// 	requireWALs(t, storage, 2)

// 	// Delete one of the static roles
// 	_, err := b.HandleRequest(ctx, &logical.Request{
// 		Operation: logical.DeleteOperation,
// 		Path:      "static-role/hashicorp",
// 		Storage:   storage,
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// 1 WAL should be cleared by the delete
// 	requireWALs(t, storage, 1)
// }

func configureDBMount(t *testing.T, b *databaseBackend, storage logical.Storage) {
	t.Helper()
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "config",
		Storage:   storage,
		Data: map[string]interface{}{
			"binddn":   "tester",
			"bindpass": "pa$$w0rd",
			"url":      "ldap://138.91.247.105",
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}
}

func createRole(t *testing.T, b *databaseBackend, storage logical.Storage, roleName string) {
	_, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "static-role/" + roleName,
		Storage:   storage,
		Data: map[string]interface{}{
			"username":        roleName,
			"dn":              "uid=hashicorp,ou=users,dc=hashicorp,dc=com",
			"rotation_period": "86400s",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

const testRoleStaticCreate = `
CREATE ROLE "{{name}}" WITH
  LOGIN
  PASSWORD '{{password}}';
`

const testRoleStaticUpdate = `
ALTER USER "{{name}}" WITH PASSWORD '{{password}}';
`

const testRoleStaticUpdateRotation = `
ALTER USER "{{name}}" WITH PASSWORD '{{password}}';GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "{{name}}";
`
