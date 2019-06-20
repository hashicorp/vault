package database

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/helper/namespace"
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

	cleanup, connURL := preparePostgresTestContainer(t, config.StorageView, b)
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

	// Test static role creation scenarios. Uses a map, so there is no guaranteed
	// ordering, so each case cleans up by deleting the role
	testCases := map[string]struct {
		account  map[string]interface{}
		expected map[string]interface{}
		err      error
	}{
		"basic": {
			account: map[string]interface{}{
				"username":        dbUser,
				"rotation_period": "5400s",
			},
			expected: map[string]interface{}{
				"username":        dbUser,
				"rotation_period": float64(5400),
			},
		},
		"missing rotation period": {
			account: map[string]interface{}{
				"username": dbUser,
			},
			err: errors.New("rotation_period is required to create static accounts"),
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

			req := &logical.Request{
				Operation: logical.CreateOperation,
				Path:      "static-roles/plugin-role-test",
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

	cleanup, connURL := preparePostgresTestContainer(t, config.StorageView, b)
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

	cleanup, connURL := preparePostgresTestContainer(t, config.StorageView, b)
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
	if resp == nil || !resp.IsError() {
		t.Fatalf("expected error, got none")
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
