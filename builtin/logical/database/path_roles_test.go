// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package database

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/helper/namespace"
	postgreshelper "github.com/hashicorp/vault/helper/testhelpers/postgresql"
	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBackend_Roles_CredentialTypes(t *testing.T) {
	config := logical.TestBackendConfig()
	config.System = logical.TestSystemView()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		credentialType   v5.CredentialType
		credentialConfig map[string]string
	}
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		expectedResp map[string]interface{}
	}{
		{
			name: "role with invalid credential type",
			args: args{
				credentialType: v5.CredentialType(10),
			},
			wantErr: true,
		},
		{
			name: "role with invalid credential type and valid credential config",
			args: args{
				credentialType: v5.CredentialType(7),
				credentialConfig: map[string]string{
					"password_policy": "test-policy",
				},
			},
			wantErr: true,
		},
		{
			name: "role with password credential type",
			args: args{
				credentialType: v5.CredentialTypePassword,
			},
			expectedResp: map[string]interface{}{
				"credential_type":   v5.CredentialTypePassword.String(),
				"credential_config": nil,
			},
		},
		{
			name: "role with password credential type and configuration",
			args: args{
				credentialType: v5.CredentialTypePassword,
				credentialConfig: map[string]string{
					"password_policy": "test-policy",
				},
			},
			expectedResp: map[string]interface{}{
				"credential_type": v5.CredentialTypePassword.String(),
				"credential_config": map[string]interface{}{
					"password_policy": "test-policy",
				},
			},
		},
		{
			name: "role with rsa_private_key credential type and default configuration",
			args: args{
				credentialType: v5.CredentialTypeRSAPrivateKey,
			},
			expectedResp: map[string]interface{}{
				"credential_type": v5.CredentialTypeRSAPrivateKey.String(),
				"credential_config": map[string]interface{}{
					"key_bits": json.Number("2048"),
					"format":   "pkcs8",
				},
			},
		},
		{
			name: "role with rsa_private_key credential type and 2048 bit configuration",
			args: args{
				credentialType: v5.CredentialTypeRSAPrivateKey,
				credentialConfig: map[string]string{
					"key_bits": "2048",
				},
			},
			expectedResp: map[string]interface{}{
				"credential_type": v5.CredentialTypeRSAPrivateKey.String(),
				"credential_config": map[string]interface{}{
					"key_bits": json.Number("2048"),
					"format":   "pkcs8",
				},
			},
		},
		{
			name: "role with rsa_private_key credential type and 3072 bit configuration",
			args: args{
				credentialType: v5.CredentialTypeRSAPrivateKey,
				credentialConfig: map[string]string{
					"key_bits": "3072",
				},
			},
			expectedResp: map[string]interface{}{
				"credential_type": v5.CredentialTypeRSAPrivateKey.String(),
				"credential_config": map[string]interface{}{
					"key_bits": json.Number("3072"),
					"format":   "pkcs8",
				},
			},
		},
		{
			name: "role with rsa_private_key credential type and 4096 bit configuration",
			args: args{
				credentialType: v5.CredentialTypeRSAPrivateKey,
				credentialConfig: map[string]string{
					"key_bits": "4096",
				},
			},
			expectedResp: map[string]interface{}{
				"credential_type": v5.CredentialTypeRSAPrivateKey.String(),
				"credential_config": map[string]interface{}{
					"key_bits": json.Number("4096"),
					"format":   "pkcs8",
				},
			},
		},
		{
			name: "role with rsa_private_key credential type invalid key_bits configuration",
			args: args{
				credentialType: v5.CredentialTypeRSAPrivateKey,
				credentialConfig: map[string]string{
					"key_bits": "256",
				},
			},
			wantErr: true,
		},
		{
			name: "role with rsa_private_key credential type invalid format configuration",
			args: args{
				credentialType: v5.CredentialTypeRSAPrivateKey,
				credentialConfig: map[string]string{
					"format": "pkcs1",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &logical.Request{
				Operation: logical.CreateOperation,
				Path:      "roles/test",
				Storage:   config.StorageView,
				Data: map[string]interface{}{
					"db_name":             "test-database",
					"creation_statements": "CREATE USER {{name}}",
					"credential_type":     tt.args.credentialType.String(),
					"credential_config":   tt.args.credentialConfig,
				},
			}

			// Create the role
			resp, err := b.HandleRequest(context.Background(), req)
			if tt.wantErr {
				assert.True(t, resp.IsError(), "expected error")
				return
			}
			assert.False(t, resp.IsError())
			assert.Nil(t, err)

			// Read the role
			req.Operation = logical.ReadOperation
			resp, err = b.HandleRequest(context.Background(), req)
			assert.False(t, resp.IsError())
			assert.Nil(t, err)
			for k, v := range tt.expectedResp {
				assert.Equal(t, v, resp.Data[k])
			}

			// Delete the role
			req.Operation = logical.DeleteOperation
			resp, err = b.HandleRequest(context.Background(), req)
			assert.False(t, resp.IsError())
			assert.Nil(t, err)
		})
	}
}

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
		// use this field to check partial error strings, otherwise use err
		errContains string
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
		"missing required fields": {
			account: map[string]interface{}{
				"username": dbUser,
			},
			path: "plugin-role-test",
			err:  errors.New("one of rotation_schedule or rotation_period must be provided to create a static account"),
		},
		"rotation_period with rotation_schedule": {
			account: map[string]interface{}{
				"username":          dbUser,
				"rotation_period":   "5400s",
				"rotation_schedule": "* * * * *",
			},
			path: "plugin-role-test",
			err:  errors.New("mutually exclusive fields rotation_period and rotation_schedule were both specified; only one of them can be provided"),
		},
		"rotation window invalid with rotation_period": {
			account: map[string]interface{}{
				"username":        dbUser,
				"rotation_period": "5400s",
				"rotation_window": "3600s",
			},
			path: "disallowed-role",
			err:  errors.New("rotation_window is invalid with use of rotation_period"),
		},
		"happy path for rotation_schedule": {
			account: map[string]interface{}{
				"username":          dbUser,
				"rotation_schedule": "* * * * *",
			},
			path: "plugin-role-test",
			expected: map[string]interface{}{
				"username":          dbUser,
				"rotation_schedule": "* * * * *",
			},
		},
		"happy path for rotation_schedule and rotation_window": {
			account: map[string]interface{}{
				"username":          dbUser,
				"rotation_schedule": "* * * * *",
				"rotation_window":   "3600s",
			},
			path: "plugin-role-test",
			expected: map[string]interface{}{
				"username":          dbUser,
				"rotation_schedule": "* * * * *",
				"rotation_window":   float64(3600),
			},
		},
		"error parsing rotation_schedule": {
			account: map[string]interface{}{
				"username":          dbUser,
				"rotation_schedule": "foo",
			},
			path:        "plugin-role-test",
			errContains: "could not parse rotation_schedule",
		},
		"rotation_window invalid": {
			account: map[string]interface{}{
				"username":          dbUser,
				"rotation_schedule": "* * * * *",
				"rotation_window":   "59s",
			},
			path:        "plugin-role-test",
			errContains: "rotation_window is invalid",
		},
		"disallowed role config": {
			account: map[string]interface{}{
				"username":        dbUser,
				"rotation_period": "5400s",
			},
			path: "disallowed-role",
			err:  errors.New("\"disallowed-role\" is not an allowed role"),
		},
		"fails to parse cronSpec with seconds": {
			account: map[string]interface{}{
				"username":          dbUser,
				"rotation_schedule": "*/10 * * * * *",
			},
			path:        "plugin-role-test-1",
			errContains: "could not parse rotation_schedule",
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
			if tc.errContains != "" {
				if !strings.Contains(resp.Error().Error(), tc.errContains) {
					t.Fatalf("expected err message: (%s), got (%s), response error: (%s)", tc.err, err, resp.Error())
				}
				return
			} else if err != nil || (resp != nil && resp.IsError()) {
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
			dataKeys := []string{
				"username",
				"password",
				"last_vault_rotation",
				"rotation_period",
				"rotation_schedule",
				"rotation_window",
			}
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

func TestBackend_StaticRole_ReadCreds(t *testing.T) {
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
	createTestPGUser(t, connURL, dbUser, dbUserDefaultPassword, testRoleStaticCreate)

	verifyPgConn(t, dbUser, dbUserDefaultPassword, connURL)

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

	testCases := map[string]struct {
		account  map[string]interface{}
		path     string
		expected map[string]interface{}
	}{
		"happy path for rotation_period": {
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
		"happy path for rotation_schedule": {
			account: map[string]interface{}{
				"username":          dbUser,
				"rotation_schedule": "* * * * *",
			},
			path: "plugin-role-test",
			expected: map[string]interface{}{
				"username":          dbUser,
				"rotation_schedule": "* * * * *",
			},
		},
		"happy path for rotation_schedule and rotation_window": {
			account: map[string]interface{}{
				"username":          dbUser,
				"rotation_schedule": "* * * * *",
				"rotation_window":   "3600s",
			},
			path: "plugin-role-test",
			expected: map[string]interface{}{
				"username":          dbUser,
				"rotation_schedule": "* * * * *",
				"rotation_window":   float64(3600),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			data = map[string]interface{}{
				"name":                "plugin-role-test",
				"db_name":             "plugin-test",
				"rotation_statements": testRoleStaticUpdate,
				"username":            dbUser,
			}

			for k, v := range tc.account {
				data[k] = v
			}

			req = &logical.Request{
				Operation: logical.CreateOperation,
				Path:      "static-roles/plugin-role-test",
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
				Path:      "static-creds/plugin-role-test",
				Storage:   config.StorageView,
				Data:      data,
			}

			resp, err = b.HandleRequest(namespace.RootContext(nil), req)
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("err:%s resp:%#v\n", err, resp)
			}

			expected := tc.expected
			actual := make(map[string]interface{})
			dataKeys := []string{
				"username",
				"password",
				"last_vault_rotation",
				"rotation_period",
				"rotation_schedule",
				"rotation_window",
				"ttl",
			}
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
				if actual["ttl"] == "" {
					t.Fatalf("expected result to contain ttl, but none found")
				}
				if v, ok := actual["last_vault_rotation"].(time.Time); !ok {
					t.Fatalf("expected last_vault_rotation to be set to time.Time type, got: %#v", v)
				}

				// delete these values before the comparison, since we can't know them in
				// advance
				delete(actual, "password")
				delete(actual, "ttl")
				delete(actual, "last_vault_rotation")
				if diff := deep.Equal(expected, actual); diff != nil {
					t.Fatal(diff)
				}
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

func TestWALsStillTrackedAfterUpdate(t *testing.T) {
	ctx := context.Background()
	b, storage, mockDB := getBackend(t)
	defer b.Cleanup(ctx)
	configureDBMount(t, storage)

	createRole(t, b, storage, mockDB, "hashicorp")

	generateWALFromFailedRotation(t, b, storage, mockDB, "hashicorp")
	requireWALs(t, storage, 1)

	resp, err := b.HandleRequest(ctx, &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "static-roles/hashicorp",
		Storage:   storage,
		Data: map[string]interface{}{
			"username":        "hashicorp",
			"rotation_period": "600s",
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatal(resp, err)
	}
	walIDs := requireWALs(t, storage, 1)

	// Now when we trigger a manual rotate, it should use the WAL's new password
	// which will tell us that the in-memory structure still kept track of the
	// WAL in addition to it still being in storage.
	wal, err := b.findStaticWAL(ctx, storage, walIDs[0])
	if err != nil {
		t.Fatal(err)
	}
	rotateRole(t, b, storage, mockDB, "hashicorp")
	role, err := b.StaticRole(ctx, storage, "hashicorp")
	if err != nil {
		t.Fatal(err)
	}
	if role.StaticAccount.Password != wal.NewPassword {
		t.Fatal()
	}
	requireWALs(t, storage, 0)
}

func TestWALsDeletedOnRoleCreationFailed(t *testing.T) {
	ctx := context.Background()
	b, storage, mockDB := getBackend(t)
	defer b.Cleanup(ctx)
	configureDBMount(t, storage)

	for i := 0; i < 3; i++ {
		mockDB.On("UpdateUser", mock.Anything, mock.Anything).
			Return(v5.UpdateUserResponse{}, errors.New("forced error")).
			Once()
		resp, err := b.HandleRequest(ctx, &logical.Request{
			Operation: logical.CreateOperation,
			Path:      "static-roles/hashicorp",
			Storage:   storage,
			Data: map[string]interface{}{
				"username":        "hashicorp",
				"db_name":         "mockv5",
				"rotation_period": "5s",
			},
		})
		if err == nil {
			t.Fatal("expected error from DB")
		}
		if !strings.Contains(err.Error(), "forced error") {
			t.Fatal("expected forced error message", resp, err)
		}
	}

	requireWALs(t, storage, 0)
}

func TestWALsDeletedOnRoleDeletion(t *testing.T) {
	ctx := context.Background()
	b, storage, mockDB := getBackend(t)
	defer b.Cleanup(ctx)
	configureDBMount(t, storage)

	// Create the roles
	roleNames := []string{"hashicorp", "2"}
	for _, roleName := range roleNames {
		createRole(t, b, storage, mockDB, roleName)
	}

	// Fail to rotate the roles
	for _, roleName := range roleNames {
		generateWALFromFailedRotation(t, b, storage, mockDB, roleName)
	}

	// Should have 2 WALs hanging around
	requireWALs(t, storage, 2)

	// Delete one of the static roles
	resp, err := b.HandleRequest(ctx, &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      "static-roles/hashicorp",
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatal(resp, err)
	}

	// 1 WAL should be cleared by the delete
	requireWALs(t, storage, 1)
}

func TestIsInsideRotationWindow(t *testing.T) {
	for _, tc := range []struct {
		name         string
		expected     bool
		data         map[string]interface{}
		now          time.Time
		timeModifier func(t time.Time) time.Time
	}{
		{
			"always returns true for rotation_period type",
			true,
			map[string]interface{}{
				"rotation_period": "86400s",
			},
			time.Now(),
			nil,
		},
		{
			"always returns true for rotation_schedule when no rotation_window set",
			true,
			map[string]interface{}{
				"rotation_schedule": "0 0 */2 * * *",
			},
			time.Now(),
			nil,
		},
		{
			"returns true for rotation_schedule when inside rotation_window",
			true,
			map[string]interface{}{
				"rotation_schedule": "0 0 */2 * * *",
				"rotation_window":   "3600s",
			},
			time.Now(),
			func(t time.Time) time.Time {
				// set current time just inside window
				return t.Add(-3640 * time.Second)
			},
		},
		{
			"returns false for rotation_schedule when outside rotation_window",
			false,
			map[string]interface{}{
				"rotation_schedule": "0 0 */2 * * *",
				"rotation_window":   "3600s",
			},
			time.Now(),
			func(t time.Time) time.Time {
				// set current time just outside window
				return t.Add(-3560 * time.Second)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			b, s, mockDB := getBackend(t)
			defer b.Cleanup(ctx)
			configureDBMount(t, s)

			testTime := tc.now
			if tc.data["rotation_schedule"] != nil && tc.timeModifier != nil {
				rotationSchedule := tc.data["rotation_schedule"].(string)
				schedule, err := b.schedule.Parse(rotationSchedule)
				if err != nil {
					t.Fatalf("could not parse rotation_schedule: %s", err)
				}
				next1 := schedule.Next(tc.now) // the next rotation time we expect
				next2 := schedule.Next(next1)  // the next rotation time after that
				testTime = tc.timeModifier(next2)
			}

			tc.data["username"] = "hashicorp"
			tc.data["db_name"] = "mockv5"
			createRoleWithData(t, b, s, mockDB, "test-role", tc.data)
			role, err := b.StaticRole(ctx, s, "test-role")
			if err != nil {
				t.Fatal(err)
			}

			isInsideWindow := role.StaticAccount.IsInsideRotationWindow(testTime)
			if tc.expected != isInsideWindow {
				t.Fatalf("expected %t, got %t", tc.expected, isInsideWindow)
			}
		})
	}
}

func createRole(t *testing.T, b *databaseBackend, storage logical.Storage, mockDB *mockNewDatabase, roleName string) {
	t.Helper()
	mockDB.On("UpdateUser", mock.Anything, mock.Anything).
		Return(v5.UpdateUserResponse{}, nil).
		Once()
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "static-roles/" + roleName,
		Storage:   storage,
		Data: map[string]interface{}{
			"username":        roleName,
			"db_name":         "mockv5",
			"rotation_period": "86400s",
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatal(resp, err)
	}
}

func createRoleWithData(t *testing.T, b *databaseBackend, s logical.Storage, mockDB *mockNewDatabase, roleName string, data map[string]interface{}) {
	t.Helper()
	mockDB.On("UpdateUser", mock.Anything, mock.Anything).
		Return(v5.UpdateUserResponse{}, nil).
		Once()
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "static-roles/" + roleName,
		Storage:   s,
		Data:      data,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatal(resp, err)
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
