package database

import (
	"context"
	"strings"
	"testing"
	"time"

	"database/sql"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	_ "github.com/lib/pq"
)

func TestBackend_StaticRole_Rotate_basic(t *testing.T) {
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
		"rotation_period":       "5400s",
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

	// verify username/password
	verifyPgConn(t, username, password, connURL)

	// re-read the creds, verifying they aren't changing on read
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

	if username != resp.Data["username"].(string) || password != resp.Data["password"].(string) {
		t.Fatal("expected re-read username/password to match, but didn't")
	}

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
		Path:      "static-creds/plugin-role-test",
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

// sanity check to make sure we don't allow an attempt of rotating credentials
// for non-static accounts, which doesn't make sense anyway, but doesn't hurt to
// verify we return an error
func TestBackend_StaticRole_Rotate_NonStaticError(t *testing.T) {
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
	// expect resp to be an error
	resp, _ = b.HandleRequest(namespace.RootContext(nil), req)
	if !resp.IsError() {
		t.Fatalf("expected error rotating non-static role")
	}

	if resp.Error().Error() != "no static role found for role name" {
		t.Fatalf("wrong error message: %s", err)
	}
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

//
// WAL testing
//
// First scenario, WAL contains a role name that does not exist.
func TestBackend_Static_QueueWAL_discard_role_not_found(t *testing.T) {
	cluster, sys := getCluster(t)
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

	assertWALCount(t, config.StorageView, 1)

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

	// verify empty queue
	if bd.credRotationQueue.Len() != 0 {
		t.Fatalf("expected zero queue items, got: %d", bd.credRotationQueue.Len())
	}

	// time.Sleep(15 * time.Second)
	assertWALCount(t, config.StorageView, 0)
}

// Second scenario, WAL contains a role name that does exist, but the role's
// LastVaultRotation is greater than the WAL has
func TestBackend_Static_QueueWAL_discard_role_newer_rotation_date(t *testing.T) {
	cluster, sys := getCluster(t)
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

	// save Now() to make sure rotation time is after this, as well as the WAL
	// time
	roleTime := time.Now()

	// create role
	data = map[string]interface{}{
		"name":                  roleName,
		"db_name":               "plugin-test",
		"creation_statements":   testRoleStaticCreate,
		"rotation_statements":   testRoleStaticUpdate,
		"revocation_statements": defaultRevocationSQL,
		"username":              "statictest",
		// low value here, to make sure the backend rotates this password at least
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

	// allow the first rotation to occur, setting LastVaultRotation
	time.Sleep(time.Second * 12)

	// cleanup the backend, then create a WAL for the role with a
	// LastVaultRotation of 1 hour ago, so that when we recreate the backend the
	// WAL will be read but discarded
	b.Cleanup(ctx)
	b = nil
	time.Sleep(time.Second * 3)

	// make a fake WAL entry with an older time
	oldRotationTime := roleTime.Add(time.Hour * -1)
	walPassword := "somejunkpassword"
	_, err = framework.PutWAL(ctx, config.StorageView, staticWALKey, &setCredentialsWAL{
		RoleName:          roleName,
		NewPassword:       walPassword,
		LastVaultRotation: oldRotationTime,
		Username:          "statictest",
	})
	if err != nil {
		t.Fatalf("error with PutWAL: %s", err)
	}

	assertWALCount(t, config.StorageView, 1)

	// reload backend
	lb, err = Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	b, ok = lb.(*databaseBackend)
	if !ok {
		t.Fatal("could not convert to db backend")
	}
	defer b.Cleanup(ctx)

	// allow enough time for populateQueue to work after boot
	time.Sleep(time.Second * 12)

	// populateQueue should have processed the entry
	assertWALCount(t, config.StorageView, 0)

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

	// grab password to verify it didn't change
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

// helper to assert the number of WAL entries is what we expect
func assertWALCount(t *testing.T, s logical.Storage, expected int) {
	var count int
	ctx := context.Background()
	keys, err := framework.ListWAL(ctx, s)
	if err != nil {
		t.Fatal("error listing WALs")
	}

	// loop through WAL keys and process any rotation ones
	for _, k := range keys {
		walEntry, _ := framework.GetWAL(ctx, s, k)
		if walEntry == nil {
			continue
		}

		if walEntry.Kind != staticWALKey {
			continue
		}
		count++
	}
	if expected != count {
		t.Fatalf("WAL count mismatch, expected (%d), got (%d)", expected, count)
	}
}

//
// end WAL testing
//

func TestBackend_StaticRole_Rotations_PostgreSQL(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup(context.Background())

	// allow initQueue to finish
	time.Sleep(5 * time.Second)
	bd := b.(*databaseBackend)
	if bd.credRotationQueue == nil {
		t.Fatal("database backend had no credential rotation queue")
	}

	// configure backend, add item and confirm length
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

	// create three static roles with different rotation periods
	testCases := []string{"65", "130", "5400"}
	for _, tc := range testCases {
		roleName := "plugin-static-role-" + tc
		data = map[string]interface{}{
			"name":                  roleName,
			"db_name":               "plugin-test",
			"creation_statements":   testRoleStaticCreate,
			"rotation_statements":   testRoleStaticUpdate,
			"revocation_statements": defaultRevocationSQL,
			"username":              "statictest" + tc,
			"rotation_period":       tc,
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

	// sleep to make sure the 65s role will be up for rotation by the time the
	// periodic function ticks
	time.Sleep(7 * time.Second)

	// sleep 75 to make sure the periodic func has time to actually run
	time.Sleep(75 * time.Second)
	pws = capturePasswords(t, b, config, testCases, pws)

	// sleep more, this should allow both sr65 and sr130 to rotate
	time.Sleep(140 * time.Second)
	pws = capturePasswords(t, b, config, testCases, pws)

	// verify all pws are as they should
	pass := true
	for k, v := range pws {
		switch {
		case k == "plugin-static-role-65":
			// expect all passwords to be different
			if v[0] == v[1] || v[1] == v[2] || v[0] == v[2] {
				pass = false
			}
		case k == "plugin-static-role-130":
			// expect the first two to be equal, but different from the third
			if v[0] != v[1] || v[0] == v[2] {
				pass = false
			}
		case k == "plugin-static-role-5400":
			// expect all passwords to be equal
			if v[0] != v[1] || v[1] != v[2] {
				pass = false
			}
		}
	}
	if !pass {
		t.Fatalf("password rotations did not match expected: %#v", pws)
	}
}

// capture the current passwords
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
