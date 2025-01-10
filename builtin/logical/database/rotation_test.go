// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Sectorbob/mlab-ns2/gae/ns/digest"
	"github.com/hashicorp/vault/builtin/logical/database/schedule"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers/mongodb"
	postgreshelper "github.com/hashicorp/vault/helper/testhelpers/postgresql"
	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/dbtxn"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	mongodbatlasapi "go.mongodb.org/atlas/mongodbatlas"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mockv5                       = "mockv5"
	dbUser                       = "vaultstatictest"
	dbUserDefaultPassword        = "password"
	testMinRotationWindowSeconds = 5
	testScheduleParseOptions     = cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow
)

func TestBackend_StaticRole_Rotation_basic(t *testing.T) {
	cluster, sys := getClusterPostgresDB(t)
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

	b.schedule = &TestSchedule{}

	cleanup, connURL := postgreshelper.PrepareTestContainer(t)
	t.Cleanup(cleanup)

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
		waitTime time.Duration
	}{
		"basic with rotation_period": {
			account: map[string]interface{}{
				"username":        dbUser,
				"rotation_period": "5400s",
			},
			path: "plugin-role-test-1",
			expected: map[string]interface{}{
				"username":        dbUser,
				"rotation_period": float64(5400),
			},
		},
		"rotation_schedule is set and expires": {
			account: map[string]interface{}{
				"username":          dbUser,
				"rotation_schedule": "*/10 * * * * *",
			},
			path: "plugin-role-test-2",
			expected: map[string]interface{}{
				"username":          dbUser,
				"rotation_schedule": "*/10 * * * * *",
			},
			waitTime: 20 * time.Second,
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
				Path:      "static-roles/" + tc.path,
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
				Path:      "static-creds/" + tc.path,
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

			// Verify username/password
			verifyPgConn(t, dbUser, password, connURL)

			// Re-read the creds, verifying they aren't changing on read
			data = map[string]interface{}{}
			req = &logical.Request{
				Operation: logical.ReadOperation,
				Path:      "static-creds/" + tc.path,
				Storage:   config.StorageView,
				Data:      data,
			}
			resp, err = b.HandleRequest(namespace.RootContext(nil), req)
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("err:%s resp:%#v\n", err, resp)
			}

			if username != resp.Data["username"].(string) || password != resp.Data["password"].(string) {
				t.Fatal("expected re-read username/password to match, but didn't")
			}

			// Trigger rotation
			data = map[string]interface{}{"name": "plugin-role-test"}
			req = &logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "rotate-role/" + tc.path,
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
				Path:      "static-creds/" + tc.path,
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

			// Verify new username/password
			verifyPgConn(t, username, newPassword, connURL)

			if tc.waitTime > 0 {
				time.Sleep(tc.waitTime)
				// Re-Read the creds after schedule expiration
				data = map[string]interface{}{}
				req = &logical.Request{
					Operation: logical.ReadOperation,
					Path:      "static-creds/" + tc.path,
					Storage:   config.StorageView,
					Data:      data,
				}
				resp, err = b.HandleRequest(namespace.RootContext(nil), req)
				if err != nil || (resp != nil && resp.IsError()) {
					t.Fatalf("err:%s resp:%#v\n", err, resp)
				}

				checkPassword := resp.Data["password"].(string)
				if newPassword == checkPassword {
					t.Fatalf("expected passwords to differ, got (%s)", checkPassword)
				}
			}
		})
	}
}

// TestBackend_StaticRole_Rotation_Schedule_ErrorRecover tests that failed
// rotations can successfully recover and that they do not occur outside of a
// rotation window.
func TestBackend_StaticRole_Rotation_Schedule_ErrorRecover(t *testing.T) {
	cluster, sys := getClusterPostgresDB(t)
	t.Cleanup(cluster.Cleanup)

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
		t.Fatal("could not convert to db backend")
	}
	defer b.Cleanup(context.Background())

	b.schedule = &TestSchedule{}

	cleanup, connURL := postgreshelper.PrepareTestContainer(t)
	t.Cleanup(cleanup)

	// create the database user
	createTestPGUser(t, connURL, dbUser, dbUserDefaultPassword, testRoleStaticCreate)
	verifyPgConn(t, dbUser, dbUserDefaultPassword, connURL)

	// Configure a connection
	connectionData := map[string]interface{}{
		"connection_url":    connURL,
		"plugin_name":       "postgresql-database-plugin",
		"verify_connection": false,
		"allowed_roles":     []string{"*"},
		"name":              "plugin-test",
	}
	configureConnection(t, b, config.StorageView, connectionData)

	// create the role that will rotate every 10th second
	// rotations will not be allowed after 5s
	data := map[string]interface{}{
		"name":                "plugin-role-test",
		"db_name":             "plugin-test",
		"rotation_statements": testRoleStaticUpdate,
		"rotation_schedule":   "*/10 * * * * *",
		"rotation_window":     "5s",
		"username":            dbUser,
	}
	req := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "static-roles/plugin-role-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err := b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Read the creds
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "static-creds/plugin-role-test",
		Storage:   config.StorageView,
	}
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	username := resp.Data["username"].(string)
	originalPassword := resp.Data["password"].(string)
	if username == "" || originalPassword == "" {
		t.Fatalf("empty username (%s) or password (%s)", username, originalPassword)
	}

	// Verify username/password
	verifyPgConn(t, dbUser, originalPassword, connURL)

	// Set invalid connection URL so we fail to rotate
	connectionData["connection_url"] = strings.Replace(connURL, "postgres:secret", "postgres:foo", 1)
	configureConnection(t, b, config.StorageView, connectionData)

	// determine next rotation schedules based on current test time
	rotationSchedule := data["rotation_schedule"].(string)
	schedule, err := b.schedule.Parse(rotationSchedule)
	if err != nil {
		t.Fatalf("could not parse rotation_schedule: %s", err)
	}
	next := schedule.Next(time.Now()) // the next rotation time we expect
	time.Sleep(next.Sub(time.Now()))

	// Re-Read the creds after schedule expiration
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "static-creds/plugin-role-test",
		Storage:   config.StorageView,
	}
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	checkPassword := resp.Data["password"].(string)
	if originalPassword != checkPassword {
		// should match because rotations should be failing
		t.Fatalf("expected passwords to match, got (%s)", checkPassword)
	}

	// wait until we are outside the rotation window so that rotations will not occur
	next = schedule.Next(time.Now()) // the next rotation time after now
	time.Sleep(next.Add(time.Second * 6).Sub(time.Now()))

	// reset to valid connection URL so we do not fail to rotate anymore
	connectionData["connection_url"] = connURL
	configureConnection(t, b, config.StorageView, connectionData)

	// we are outside a rotation window, Re-Read the creds
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "static-creds/plugin-role-test",
		Storage:   config.StorageView,
	}
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	checkPassword = resp.Data["password"].(string)
	if originalPassword != checkPassword {
		// should match because rotations should not occur outside the rotation window
		t.Fatalf("expected passwords to match, got (%s)", checkPassword)
	}
	// Verify new username/password
	verifyPgConn(t, username, checkPassword, connURL)

	// sleep until the next rotation time with a buffer to ensure we had time to rotate
	next = schedule.Next(time.Now()) // the next rotation time we expect
	time.Sleep(next.Add(time.Second * 5).Sub(time.Now()))

	// Re-Read the creds
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "static-creds/plugin-role-test",
		Storage:   config.StorageView,
	}
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	checkPassword = resp.Data["password"].(string)
	if originalPassword == checkPassword {
		// should differ because we slept until the next rotation time
		t.Fatalf("expected passwords to differ, got (%s)", checkPassword)
	}

	// Verify new username/password
	verifyPgConn(t, username, checkPassword, connURL)

	eventSender.Stop() // avoid race detector
	// check that we got a successful rotation event
	if len(eventSender.Events) == 0 {
		t.Fatal("Expected to have some events but got none")
	}
	// check that we got a rotate-fail event
	found := false
	for _, event := range eventSender.Events {
		if string(event.Type) == "database/rotate-fail" {
			found = true
			break
		}
	}
	assert.True(t, found)
	found = false
	for _, event := range eventSender.Events {
		if string(event.Type) == "database/rotate" {
			found = true
			break
		}
	}
	assert.True(t, found)
}

// Sanity check to make sure we don't allow an attempt of rotating credentials
// for non-static accounts, which doesn't make sense anyway, but doesn't hurt to
// verify we return an error
func TestBackend_StaticRole_Rotation_NonStaticError(t *testing.T) {
	cluster, sys := getClusterPostgresDB(t)
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

	cleanup, connURL := postgreshelper.PrepareTestContainer(t)
	defer cleanup()

	// create the database user
	createTestPGUser(t, connURL, dbUser, dbUserDefaultPassword, testRoleStaticCreate)

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

	// Verify username/password
	verifyPgConn(t, dbUser, dbUserDefaultPassword, connURL)
	// Trigger rotation
	data = map[string]interface{}{"name": "plugin-role-test"}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "rotate-role/plugin-role-test",
		Storage:   config.StorageView,
		Data:      data,
	}
	// expect resp to be an error
	resp, _ = b.HandleRequest(namespace.RootContext(nil), req)
	if !resp.IsError() {
		t.Fatalf("expected error rotating non-static role")
	}

	if resp.Error().Error() != "no static role found for role name" {
		t.Fatalf("wrong error message: %s", err)
	}
}

func TestBackend_StaticRole_Rotation_Revoke_user(t *testing.T) {
	cluster, sys := getClusterPostgresDB(t)
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

	cleanup, connURL := postgreshelper.PrepareTestContainer(t)
	defer cleanup()

	// create the database user
	createTestPGUser(t, connURL, dbUser, dbUserDefaultPassword, testRoleStaticCreate)

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
		revoke          *bool
		expectVerifyErr bool
	}{
		// Default case: user does not specify, Vault leaves the database user
		// untouched, and the final connection check passes because the user still
		// exists
		"unset": {},
		// Revoke on delete. The final connection check should fail because the user
		// no longer exists
		"revoke": {
			revoke:          newBoolPtr(true),
			expectVerifyErr: true,
		},
		// Revoke false, final connection check should still pass
		"persist": {
			revoke: newBoolPtr(false),
		},
	}
	for k, tc := range testCases {
		t.Run(k, func(t *testing.T) {
			data = map[string]interface{}{
				"name":                "plugin-role-test",
				"db_name":             "plugin-test",
				"rotation_statements": testRoleStaticUpdate,
				"username":            dbUser,
				"rotation_period":     "5400s",
			}
			if tc.revoke != nil {
				data["revoke_user_on_delete"] = *tc.revoke
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

			username := resp.Data["username"].(string)
			password := resp.Data["password"].(string)
			if username == "" || password == "" {
				t.Fatalf("empty username (%s) or password (%s)", username, password)
			}

			// Verify username/password
			verifyPgConn(t, username, password, connURL)

			// delete the role, expect the default where the user is not destroyed
			// Read the creds
			req = &logical.Request{
				Operation: logical.DeleteOperation,
				Path:      "static-roles/plugin-role-test",
				Storage:   config.StorageView,
			}

			resp, err = b.HandleRequest(namespace.RootContext(nil), req)
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("err:%s resp:%#v\n", err, resp)
			}

			// Verify new username/password still work
			verifyPgConn(t, username, password, connURL)
		})
	}
}

func createTestPGUser(t *testing.T, connURL string, username, password, query string) {
	t.Helper()
	log.Printf("[TRACE] Creating test user")

	db, err := sql.Open("pgx", connURL)
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Start a transaction
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	m := map[string]string{
		"name":     username,
		"password": password,
	}
	if err := dbtxn.ExecuteTxQueryDirect(ctx, tx, m, query); err != nil {
		t.Fatal(err)
	}
	// Commit the transaction
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}

func verifyPgConn(t *testing.T, username, password, connURL string) {
	t.Helper()
	cURL := strings.Replace(connURL, "postgres:secret", username+":"+password, 1)
	db, err := sql.Open("pgx", cURL)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}
}

// WAL testing
//
// First scenario, WAL contains a role name that does not exist.
func TestBackend_StaticRole_Rotation_QueueWAL_discard_role_not_found(t *testing.T) {
	cluster, sys := getClusterPostgresDB(t)
	defer cluster.Cleanup()

	ctx := context.Background()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys

	_, err := framework.PutWAL(ctx, config.StorageView, staticWALKey, &setCredentialsWAL{
		RoleName: "doesnotexist",
	})
	if err != nil {
		t.Fatalf("error with PutWAL: %s", err)
	}

	assertWALCount(t, config.StorageView, 1, staticWALKey)

	b, err := Factory(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup(ctx)

	time.Sleep(5 * time.Second)
	bd := b.(*databaseBackend)
	if bd.credRotationQueue == nil {
		t.Fatal("database backend had no credential rotation queue")
	}

	// Verify empty queue
	if bd.credRotationQueue.Len() != 0 {
		t.Fatalf("expected zero queue items, got: %d", bd.credRotationQueue.Len())
	}

	assertWALCount(t, config.StorageView, 0, staticWALKey)
}

// Second scenario, WAL contains a role name that does exist, but the role's
// LastVaultRotation is greater than the WAL has
func TestBackend_StaticRole_Rotation_QueueWAL_discard_role_newer_rotation_date(t *testing.T) {
	cluster, sys := getClusterPostgresDB(t)
	defer cluster.Cleanup()

	ctx := context.Background()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys

	roleName := "test-discard-by-date"
	lb, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	b, ok := lb.(*databaseBackend)
	if !ok {
		t.Fatal("could not convert to db backend")
	}

	cleanup, connURL := postgreshelper.PrepareTestContainer(t)
	defer cleanup()

	// create the database user
	createTestPGUser(t, connURL, dbUser, dbUserDefaultPassword, testRoleStaticCreate)

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

	// Save Now() to make sure rotation time is after this, as well as the WAL
	// time
	roleTime := time.Now()

	// Create role
	data = map[string]interface{}{
		"name":                roleName,
		"db_name":             "plugin-test",
		"rotation_statements": testRoleStaticUpdate,
		"username":            dbUser,
		// Low value here, to make sure the backend rotates this password at least
		// once before we compare it to the WAL
		"rotation_period": "10s",
	}

	req = &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "static-roles/" + roleName,
		Storage:   config.StorageView,
		Data:      data,
	}

	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Allow the first rotation to occur, setting LastVaultRotation
	time.Sleep(time.Second * 12)

	// Cleanup the backend, then create a WAL for the role with a
	// LastVaultRotation of 1 hour ago, so that when we recreate the backend the
	// WAL will be read but discarded
	b.Cleanup(ctx)
	b = nil
	time.Sleep(time.Second * 3)

	// Make a fake WAL entry with an older time
	oldRotationTime := roleTime.Add(time.Hour * -1)
	walPassword := "somejunkpassword"
	_, err = framework.PutWAL(ctx, config.StorageView, staticWALKey, &setCredentialsWAL{
		RoleName:          roleName,
		NewPassword:       walPassword,
		LastVaultRotation: oldRotationTime,
		Username:          dbUser,
	})
	if err != nil {
		t.Fatalf("error with PutWAL: %s", err)
	}

	assertWALCount(t, config.StorageView, 1, staticWALKey)

	// Reload backend
	lb, err = Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	b, ok = lb.(*databaseBackend)
	if !ok {
		t.Fatal("could not convert to db backend")
	}
	defer b.Cleanup(ctx)

	// Allow enough time for populateQueue to work after boot
	time.Sleep(time.Second * 12)

	// PopulateQueue should have processed the entry
	assertWALCount(t, config.StorageView, 0, staticWALKey)

	// Read the role
	data = map[string]interface{}{}
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "static-roles/" + roleName,
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	lastVaultRotation := resp.Data["last_vault_rotation"].(time.Time)
	if !lastVaultRotation.After(oldRotationTime) {
		t.Fatal("last vault rotation time not greater than WAL time")
	}

	if !lastVaultRotation.After(roleTime) {
		t.Fatal("last vault rotation time not greater than role creation time")
	}

	// Grab password to verify it didn't change
	req = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "static-creds/" + roleName,
		Storage:   config.StorageView,
	}
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	password := resp.Data["password"].(string)
	if password == walPassword {
		t.Fatalf("expected password to not be changed by WAL, but was")
	}
}

// Helper to assert the number of WAL entries is what we expect
func assertWALCount(t *testing.T, s logical.Storage, expected int, key string) {
	t.Helper()

	var count int
	ctx := context.Background()
	keys, err := framework.ListWAL(ctx, s)
	if err != nil {
		t.Fatal("error listing WALs")
	}

	// Loop through WAL keys and process any rotation ones
	for _, k := range keys {
		walEntry, _ := framework.GetWAL(ctx, s, k)
		if walEntry == nil {
			continue
		}

		if walEntry.Kind != key {
			continue
		}
		count++
	}
	if expected != count {
		t.Fatalf("WAL count mismatch, expected (%d), got (%d)", expected, count)
	}
}

//
// End WAL testing
//

type userCreator func(t *testing.T, username, password string)

func TestBackend_StaticRole_Rotation_PostgreSQL(t *testing.T) {
	cleanup, connURL := postgreshelper.PrepareTestContainer(t)
	defer cleanup()
	uc := userCreator(func(t *testing.T, username, password string) {
		createTestPGUser(t, connURL, username, password, testRoleStaticCreate)
	})
	testBackend_StaticRole_Rotations(t, uc, map[string]interface{}{
		"connection_url": connURL,
		"plugin_name":    "postgresql-database-plugin",
	})
}

func TestBackend_StaticRole_Rotation_MongoDB(t *testing.T) {
	cleanup, connURL := mongodb.PrepareTestContainerWithDatabase(t, "5.0.10", "vaulttestdb")
	defer cleanup()

	uc := userCreator(func(t *testing.T, username, password string) {
		testCreateDBUser(t, connURL, "vaulttestdb", username, password)
	})
	testBackend_StaticRole_Rotations(t, uc, map[string]interface{}{
		"connection_url": connURL,
		"plugin_name":    "mongodb-database-plugin",
	})
}

func TestBackend_StaticRole_Rotation_MongoDBAtlas(t *testing.T) {
	// To get the project ID, connect to cloud.mongodb.com, go to the vault-test project and
	// look at Project Settings.
	projID := os.Getenv("VAULT_MONGODBATLAS_PROJECT_ID")
	// For the private and public key, go to Organization Access Manager on cloud.mongodb.com,
	// choose Create API Key, then create one using the defaults.  Then go back to the vault-test
	// project and add the API key to it, with permissions "Project Owner".
	privKey := os.Getenv("VAULT_MONGODBATLAS_PRIVATE_KEY")
	pubKey := os.Getenv("VAULT_MONGODBATLAS_PUBLIC_KEY")
	if projID == "" {
		t.Logf("Skipping MongoDB Atlas test because VAULT_MONGODBATLAS_PROJECT_ID not set")
		t.SkipNow()
	}

	transport := digest.NewTransport(pubKey, privKey)
	cl, err := transport.Client()
	if err != nil {
		t.Fatal(err)
	}

	api, err := mongodbatlasapi.New(cl)
	if err != nil {
		t.Fatal(err)
	}

	uc := userCreator(func(t *testing.T, username, password string) {
		// Delete the user in case it's still there from an earlier run, ignore
		// errors in case it's not.
		_, _ = api.DatabaseUsers.Delete(context.Background(), "admin", projID, username)

		req := &mongodbatlasapi.DatabaseUser{
			Username:     username,
			Password:     password,
			DatabaseName: "admin",
			Roles:        []mongodbatlasapi.Role{{RoleName: "atlasAdmin", DatabaseName: "admin"}},
		}
		_, _, err := api.DatabaseUsers.Create(context.Background(), projID, req)
		if err != nil {
			t.Fatal(err)
		}
	})
	testBackend_StaticRole_Rotations(t, uc, map[string]interface{}{
		"plugin_name": "mongodbatlas-database-plugin",
		"project_id":  projID,
		"private_key": privKey,
		"public_key":  pubKey,
	})
}

// TestQueueTickIntervalKeyConfig tests the configuration of queueTickIntervalKey
// does not break on invalid values.
func TestQueueTickIntervalKeyConfig(t *testing.T) {
	t.Parallel()
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

	values := []string{"1", "0", "-1"}
	for _, v := range values {
		t.Run("test"+v, func(t *testing.T) {
			config := logical.TestBackendConfig()
			config.StorageView = &logical.InmemStorage{}
			config.System = sys
			config.Config[queueTickIntervalKey] = v

			// Rotation ticker starts running in Factory call
			b, err := Factory(context.Background(), config)
			require.Nil(t, err)
			b.Cleanup(context.Background())
		})
	}
}

func testBackend_StaticRole_Rotations(t *testing.T, createUser userCreator, opts map[string]interface{}) {
	// We need to set this value for the plugin to run, but it doesn't matter what we set it to.
	oldToken := os.Getenv(pluginutil.PluginUnwrapTokenEnv)
	os.Setenv(pluginutil.PluginUnwrapTokenEnv, "...")
	defer func() {
		if oldToken != "" {
			os.Setenv(pluginutil.PluginUnwrapTokenEnv, oldToken)
		} else {
			os.Unsetenv(pluginutil.PluginUnwrapTokenEnv)
		}
	}()

	cluster, sys := getClusterPostgresDB(t)
	defer cluster.Cleanup()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys
	// Change background task interval to 1s to give more margin
	// for it to successfully run during the sleeps below.
	config.Config[queueTickIntervalKey] = "1"

	// Rotation ticker starts running in Factory call
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup(context.Background())

	// allow initQueue to finish
	bd := b.(*databaseBackend)
	if bd.credRotationQueue == nil {
		t.Fatal("database backend had no credential rotation queue")
	}

	// Configure a connection
	data := map[string]interface{}{
		"verify_connection": false,
		"allowed_roles":     []string{"*"},
	}
	for k, v := range opts {
		data[k] = v
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

	testCases := []string{"10", "20", "100"}
	// Create database users ahead
	for _, tc := range testCases {
		createUser(t, "statictest"+tc, "test")
	}

	// create three static roles with different rotation periods
	for _, tc := range testCases {
		roleName := "plugin-static-role-" + tc
		data = map[string]interface{}{
			"name":            roleName,
			"db_name":         "plugin-test",
			"username":        "statictest" + tc,
			"rotation_period": tc,
		}

		req = &logical.Request{
			Operation: logical.CreateOperation,
			Path:      "static-roles/" + roleName,
			Storage:   config.StorageView,
			Data:      data,
		}

		resp, err = b.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%s resp:%#v\n", err, resp)
		}
	}

	// verify the queue has 3 items in it
	if bd.credRotationQueue.Len() != 3 {
		t.Fatalf("expected 3 items in the rotation queue, got: (%d)", bd.credRotationQueue.Len())
	}

	// List the roles
	data = map[string]interface{}{}
	req = &logical.Request{
		Operation: logical.ListOperation,
		Path:      "static-roles/",
		Storage:   config.StorageView,
		Data:      data,
	}
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	keys := resp.Data["keys"].([]string)
	if len(keys) != 3 {
		t.Fatalf("expected 3 roles, got: (%d)", len(keys))
	}

	// capture initial passwords, before the periodic function is triggered
	pws := make(map[string][]string, 0)
	pws = capturePasswords(t, b, config, testCases, pws)

	// sleep to make sure the periodic func has time to actually run
	time.Sleep(15 * time.Second)
	pws = capturePasswords(t, b, config, testCases, pws)

	// sleep more, this should allow both sr10 and sr20 to rotate
	time.Sleep(10 * time.Second)
	pws = capturePasswords(t, b, config, testCases, pws)

	// verify all pws are as they should
	pass := true
	for k, v := range pws {
		if len(v) < 3 {
			t.Fatalf("expected to find 3 passwords for (%s), only found (%d)", k, len(v))
		}
		switch {
		case k == "plugin-static-role-10":
			// expect all passwords to be different
			if v[0] == v[1] || v[1] == v[2] || v[0] == v[2] {
				pass = false
			}
		case k == "plugin-static-role-20":
			// expect the first two to be equal, but different from the third
			if v[0] != v[1] || v[0] == v[2] {
				pass = false
			}
		case k == "plugin-static-role-100":
			// expect all passwords to be equal
			if v[0] != v[1] || v[1] != v[2] {
				pass = false
			}
		default:
			t.Fatalf("unexpected password key: %v", k)
		}
	}
	if !pass {
		t.Fatalf("password rotations did not match expected: %#v", pws)
	}
}

func testCreateDBUser(t testing.TB, connURL, db, username, password string) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connURL))
	if err != nil {
		t.Fatal(err)
	}

	createUserCmd := &createUserCommand{
		Username: username,
		Password: password,
		Roles:    []interface{}{},
	}
	result := client.Database(db).RunCommand(ctx, createUserCmd, nil)
	if result.Err() != nil {
		t.Fatal(result.Err())
	}
}

type createUserCommand struct {
	Username string        `bson:"createUser"`
	Password string        `bson:"pwd"`
	Roles    []interface{} `bson:"roles"`
}

// Demonstrates a bug fix for the credential rotation not releasing locks
func TestBackend_StaticRole_Rotation_LockRegression(t *testing.T) {
	cluster, sys := getClusterPostgresDB(t)
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

	cleanup, connURL := postgreshelper.PrepareTestContainer(t)
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

	createTestPGUser(t, connURL, dbUser, dbUserDefaultPassword, testRoleStaticCreate)
	data = map[string]interface{}{
		"name":                "plugin-role-test",
		"db_name":             "plugin-test",
		"rotation_statements": testRoleStaticUpdate,
		"username":            dbUser,
		"rotation_period":     "7s",
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
	for i := 0; i < 25; i++ {
		req = &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "static-roles/plugin-role-test",
			Storage:   config.StorageView,
			Data:      data,
		}

		resp, err = b.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%s resp:%#v\n", err, resp)
		}

		// sleeping is needed to trigger the deadlock, otherwise things are
		// processed too quickly to trigger the rotation lock on so few roles
		time.Sleep(500 * time.Millisecond)
	}
}

func TestBackend_StaticRole_Rotation_Invalid_Role(t *testing.T) {
	cluster, sys := getClusterPostgresDB(t)
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

	cleanup, connURL := postgreshelper.PrepareTestContainer(t)
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

	data = map[string]interface{}{
		"name":                "plugin-role-test",
		"db_name":             "plugin-test",
		"rotation_statements": testRoleStaticUpdate,
		"username":            dbUser,
		"rotation_period":     "5400s",
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

	// Pop manually key to emulate a queue without existing key
	b.credRotationQueue.PopByKey("plugin-role-test")

	// Make sure queue is empty
	if b.credRotationQueue.Len() != 0 {
		t.Fatalf("expected queue length to be 0 but is %d", b.credRotationQueue.Len())
	}

	// Trigger rotation
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

	// Check if key is in queue
	if b.credRotationQueue.Len() != 1 {
		t.Fatalf("expected queue length to be 1 but is %d", b.credRotationQueue.Len())
	}
}

func TestRollsPasswordForwardsUsingWAL(t *testing.T) {
	ctx := context.Background()
	b, storage, mockDB := getBackend(t)
	defer b.Cleanup(ctx)
	configureDBMount(t, storage)
	createRole(t, b, storage, mockDB, "hashicorp")

	role, err := b.StaticRole(ctx, storage, "hashicorp")
	if err != nil {
		t.Fatal(err)
	}
	oldPassword := role.StaticAccount.Password

	generateWALFromFailedRotation(t, b, storage, mockDB, "hashicorp")

	walIDs := requireWALs(t, storage, 1)
	wal, err := b.findStaticWAL(ctx, storage, walIDs[0])
	if err != nil {
		t.Fatal(err)
	}
	role, err = b.StaticRole(ctx, storage, "hashicorp")
	if err != nil {
		t.Fatal(err)
	}
	// Role's password should still be the WAL's old password
	if role.StaticAccount.Password != oldPassword {
		t.Fatal(role.StaticAccount.Password, oldPassword)
	}

	rotateRole(t, b, storage, mockDB, "hashicorp")

	role, err = b.StaticRole(ctx, storage, "hashicorp")
	if err != nil {
		t.Fatal(err)
	}
	if role.StaticAccount.Password != wal.NewPassword {
		t.Fatal("role password", role.StaticAccount.Password, "WAL new password", wal.NewPassword)
	}
	// WAL should be cleared by the successful rotate
	requireWALs(t, storage, 0)
}

func TestStoredWALsCorrectlyProcessed(t *testing.T) {
	const walNewPassword = "new-password-from-wal"

	rotationPeriodData := map[string]interface{}{
		"username":        "hashicorp",
		"db_name":         mockv5,
		"rotation_period": "86400s",
	}

	for _, tc := range []struct {
		name         string
		shouldRotate bool
		wal          *setCredentialsWAL
		data         map[string]interface{}
	}{
		{
			"WAL is kept and used for roll forward",
			true,
			&setCredentialsWAL{
				RoleName:          "hashicorp",
				Username:          "hashicorp",
				NewPassword:       walNewPassword,
				LastVaultRotation: time.Now().Add(time.Hour),
			},
			rotationPeriodData,
		},
		{
			"zero-time WAL is discarded on load",
			false,
			&setCredentialsWAL{
				RoleName:          "hashicorp",
				Username:          "hashicorp",
				NewPassword:       walNewPassword,
				LastVaultRotation: time.Time{},
			},
			rotationPeriodData,
		},
		{
			"rotation_period empty-password WAL is kept but a new password is generated",
			true,
			&setCredentialsWAL{
				RoleName:          "hashicorp",
				Username:          "hashicorp",
				NewPassword:       "",
				LastVaultRotation: time.Now().Add(time.Hour),
			},
			rotationPeriodData,
		},
		{
			"rotation_schedule empty-password WAL is kept but a new password is generated",
			true,
			&setCredentialsWAL{
				RoleName:          "hashicorp",
				Username:          "hashicorp",
				NewPassword:       "",
				LastVaultRotation: time.Now().Add(time.Hour),
			},
			map[string]interface{}{
				"username":          "hashicorp",
				"db_name":           mockv5,
				"rotation_schedule": "*/10 * * * * *",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			config := logical.TestBackendConfig()
			storage := &logical.InmemStorage{}
			config.StorageView = storage
			b := Backend(config)
			defer b.Cleanup(ctx)
			mockDB := setupMockDB(b)
			if err := b.Setup(ctx, config); err != nil {
				t.Fatal(err)
			}
			b.credRotationQueue = queue.New()
			b.schedule = &TestSchedule{}
			configureDBMount(t, config.StorageView)
			createRoleWithData(t, b, config.StorageView, mockDB, tc.wal.RoleName, tc.data)
			role, err := b.StaticRole(ctx, config.StorageView, "hashicorp")
			if err != nil {
				t.Fatal(err)
			}
			initialPassword := role.StaticAccount.Password

			// Set up a WAL for our test case
			framework.PutWAL(ctx, config.StorageView, staticWALKey, tc.wal)
			requireWALs(t, config.StorageView, 1)
			// Reset the rotation queue to simulate startup memory state
			b.credRotationQueue = queue.New()

			// Now finish the startup process by populating the queue, which should discard the WAL
			b.initQueue(ctx, config)

			if tc.shouldRotate {
				requireWALs(t, storage, 1)
			} else {
				requireWALs(t, storage, 0)
			}

			// Run one tick
			mockDB.On("UpdateUser", mock.Anything, mock.Anything).
				Return(v5.UpdateUserResponse{}, nil).
				Once()
			b.rotateCredentials(ctx, storage)
			requireWALs(t, storage, 0)

			role, err = b.StaticRole(ctx, storage, "hashicorp")
			if err != nil {
				t.Fatal(err)
			}
			item, err := b.popFromRotationQueueByKey("hashicorp")
			if err != nil {
				t.Fatal(err)
			}

			nextRotationTime := role.StaticAccount.NextRotationTime()
			if tc.shouldRotate {
				if tc.wal.NewPassword != "" {
					// Should use WAL's new_password field
					if role.StaticAccount.Password != walNewPassword {
						t.Fatal()
					}
				} else {
					// Should rotate but ignore WAL's new_password field
					if role.StaticAccount.Password == initialPassword {
						t.Fatal()
					}
					if role.StaticAccount.Password == walNewPassword {
						t.Fatal()
					}
				}
				// Ensure the role was not promoted for early rotation
				assertPriorityUnchanged(t, item.Priority, nextRotationTime)
			} else {
				// Ensure the role was not promoted for early rotation
				assertPriorityUnchanged(t, item.Priority, nextRotationTime)
				if role.StaticAccount.Password != initialPassword {
					t.Fatal("password should not have been rotated yet")
				}
			}
		})
	}
}

func TestDeletesOlderWALsOnLoad(t *testing.T) {
	ctx := context.Background()
	b, storage, mockDB := getBackend(t)
	defer b.Cleanup(ctx)
	configureDBMount(t, storage)
	createRole(t, b, storage, mockDB, "hashicorp")

	// Create 4 WALs, with a clear winner for most recent.
	wal := &setCredentialsWAL{
		RoleName:          "hashicorp",
		Username:          "hashicorp",
		NewPassword:       "some-new-password",
		LastVaultRotation: time.Now(),
	}
	for i := 0; i < 3; i++ {
		_, err := framework.PutWAL(ctx, storage, staticWALKey, wal)
		if err != nil {
			t.Fatal(err)
		}
	}
	time.Sleep(2 * time.Second)
	// We expect this WAL to have the latest createdAt timestamp
	walID, err := framework.PutWAL(ctx, storage, staticWALKey, wal)
	if err != nil {
		t.Fatal(err)
	}
	requireWALs(t, storage, 4)

	walMap, err := b.loadStaticWALs(ctx, storage)
	if err != nil {
		t.Fatal(err)
	}
	if len(walMap) != 1 || walMap["hashicorp"] == nil || walMap["hashicorp"].walID != walID {
		t.Fatal()
	}
	requireWALs(t, storage, 1)
}

func generateWALFromFailedRotation(t *testing.T, b *databaseBackend, storage logical.Storage, mockDB *mockNewDatabase, roleName string) {
	t.Helper()
	mockDB.On("UpdateUser", mock.Anything, mock.Anything).
		Return(v5.UpdateUserResponse{}, errors.New("forced error")).
		Once()
	_, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "rotate-role/" + roleName,
		Storage:   storage,
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func rotateRole(t *testing.T, b *databaseBackend, storage logical.Storage, mockDB *mockNewDatabase, roleName string) {
	t.Helper()
	mockDB.On("UpdateUser", mock.Anything, mock.Anything).
		Return(v5.UpdateUserResponse{}, nil).
		Once()
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "rotate-role/" + roleName,
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatal(resp, err)
	}
}

// returns a slice of the WAL IDs in storage
func requireWALs(t *testing.T, storage logical.Storage, expectedCount int) []string {
	t.Helper()
	wals, err := storage.List(context.Background(), "wal/")
	if err != nil {
		t.Fatal(err)
	}
	if len(wals) != expectedCount {
		t.Fatal("expected WALs", expectedCount, "got", len(wals))
	}

	return wals
}

func getBackend(t *testing.T) (*databaseBackend, logical.Storage, *mockNewDatabase) {
	t.Helper()
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	// Create and init the backend ourselves instead of using a Factory because
	// the factory function kicks off threads that cause racy tests.
	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}
	b.schedule = &TestSchedule{}
	b.credRotationQueue = queue.New()
	b.populateQueue(context.Background(), config.StorageView)

	mockDB := setupMockDB(b)

	return b, config.StorageView, mockDB
}

func setupMockDB(b *databaseBackend) *mockNewDatabase {
	mockDB := &mockNewDatabase{}
	mockDB.On("Initialize", mock.Anything, mock.Anything).Return(v5.InitializeResponse{}, nil)
	mockDB.On("Close").Return(nil)
	mockDB.On("Type").Return("mock", nil)
	dbw := databaseVersionWrapper{
		v5: mockDB,
	}

	dbi := &dbPluginInstance{
		database: dbw,
		id:       "foo-id",
		name:     mockv5,
	}
	b.connections.Put(mockv5, dbi)

	return mockDB
}

// configureDBMount puts config directly into storage to avoid the DB engine's
// plugin init code paths, allowing us to use a manually populated mock DB object.
func configureDBMount(t *testing.T, storage logical.Storage) {
	t.Helper()
	entry, err := logical.StorageEntryJSON(fmt.Sprintf("config/"+mockv5), &DatabaseConfig{
		AllowedRoles: []string{"*"},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = storage.Put(context.Background(), entry)
	if err != nil {
		t.Fatal(err)
	}
}

// capturePasswords captures the current passwords at the time of calling, and
// returns a map of username / passwords building off of the input map
func capturePasswords(t *testing.T, b logical.Backend, config *logical.BackendConfig, testCases []string, pws map[string][]string) map[string][]string {
	new := make(map[string][]string, 0)
	for _, tc := range testCases {
		// Read the role
		roleName := "plugin-static-role-" + tc
		req := &logical.Request{
			Operation: logical.ReadOperation,
			Path:      "static-creds/" + roleName,
			Storage:   config.StorageView,
		}
		resp, err := b.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%s resp:%#v\n", err, resp)
		}

		username := resp.Data["username"].(string)
		password := resp.Data["password"].(string)
		if username == "" || password == "" {
			t.Fatalf("expected both username/password for (%s), got (%s), (%s)", roleName, username, password)
		}
		new[roleName] = append(new[roleName], password)
	}

	for k, v := range new {
		pws[k] = append(pws[k], v...)
	}

	return pws
}

func configureConnection(t *testing.T, b *databaseBackend, s logical.Storage, data map[string]interface{}) {
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/" + data["name"].(string),
		Storage:   s,
		Data:      data,
	}
	resp, err := b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}
}

func newBoolPtr(b bool) *bool {
	v := b
	return &v
}

// assertPriorityUnchanged is a helper to verify that the priority is the
// expected value for a given rotation time
func assertPriorityUnchanged(t *testing.T, priority int64, nextRotationTime time.Time) {
	t.Helper()
	if priority != nextRotationTime.Unix() {
		t.Fatalf("expected next rotation at %s, but got %s", nextRotationTime, time.Unix(priority, 0).String())
	}
}

var _ schedule.Scheduler = &TestSchedule{}

type TestSchedule struct{}

func (d *TestSchedule) Parse(rotationSchedule string) (*cron.SpecSchedule, error) {
	parser := cron.NewParser(testScheduleParseOptions)
	schedule, err := parser.Parse(rotationSchedule)
	if err != nil {
		return nil, err
	}
	sched, ok := schedule.(*cron.SpecSchedule)
	if !ok {
		return nil, fmt.Errorf("invalid rotation schedule")
	}
	return sched, nil
}

func (d *TestSchedule) ValidateRotationWindow(s int) error {
	if s < testMinRotationWindowSeconds {
		return fmt.Errorf("rotation_window must be %d seconds or more", testMinRotationWindowSeconds)
	}
	return nil
}
