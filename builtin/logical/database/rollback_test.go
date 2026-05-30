// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package database

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/namespace"
	postgreshelper "github.com/hashicorp/vault/helper/testhelpers/postgresql"
	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	databaseUser    = "postgres"
	defaultPassword = "secret"
	newPrivateKey   = "new-private-key-pem"
)

// Tests that the WAL rollback function rolls back the database password.
// The database password should be rolled back when:
//   - A WAL entry exists
//   - Password has been altered on the database
//   - Password has not been updated in storage
func TestBackend_RotateRootCredentials_WAL_rollback(t *testing.T) {
	_, sys := getClusterPostgresDB(t)

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

	cleanup, connURL := postgreshelper.PrepareTestContainer(t)
	defer cleanup()

	connURL = strings.ReplaceAll(connURL, "postgres:secret", "{{username}}:{{password}}")

	// Configure a connection to the database
	data := map[string]interface{}{
		"connection_url": connURL,
		"plugin_name":    "postgresql-database-plugin",
		"allowed_roles":  []string{"plugin-role-test"},
		"username":       databaseUser,
		"password":       defaultPassword,
	}
	resp, err := lb.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Create a role
	data = map[string]interface{}{
		"db_name":             "plugin-test",
		"creation_statements": testRole,
		"max_ttl":             "10m",
	}
	resp, err = lb.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/plugin-role-test",
		Storage:   config.StorageView,
		Data:      data,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Read credentials to verify this initially works
	credReq := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/plugin-role-test",
		Storage:   config.StorageView,
		Data:      make(map[string]interface{}),
	}
	credResp, err := lb.HandleRequest(context.Background(), credReq)
	if err != nil || (credResp != nil && credResp.IsError()) {
		t.Fatalf("err:%s resp:%v\n", err, credResp)
	}

	// Get a connection to the database plugin
	dbi, err := dbBackend.GetConnection(context.Background(),
		config.StorageView, "plugin-test")
	if err != nil {
		t.Fatal(err)
	}

	// Alter the database password so it no longer matches what is in storage
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	updateReq := v5.UpdateUserRequest{
		Username: databaseUser,
		Password: &v5.ChangePassword{
			NewPassword: "newSecret",
		},
	}
	_, err = dbi.database.UpdateUser(ctx, updateReq, false)
	if err != nil {
		t.Fatal(err)
	}

	// Clear the plugin connection to verify we're no longer able to connect
	err = dbBackend.ClearConnection("plugin-test")
	if err != nil {
		t.Fatal(err)
	}

	// Reading credentials should no longer work
	credResp, err = lb.HandleRequest(namespace.RootContext(nil), credReq)
	if err == nil {
		t.Fatalf("expected authentication to fail when reading credentials")
	}

	// Put a WAL entry that will be used for rolling back the database password
	walEntry := NewRotateRootCredentialsWALPasswordEntry("plugin-test", databaseUser, "newSecret", defaultPassword)
	_, err = framework.PutWAL(context.Background(), config.StorageView, rotateRootWALKey, walEntry)
	if err != nil {
		t.Fatal(err)
	}
	assertWALCount(t, config.StorageView, 1, rotateRootWALKey)

	// Trigger an immediate RollbackOperation so that the WAL rollback
	// function can use the WAL entry to roll back the database password
	_, err = lb.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.RollbackOperation,
		Path:      "",
		Storage:   config.StorageView,
		Data: map[string]interface{}{
			"immediate": true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertWALCount(t, config.StorageView, 0, rotateRootWALKey)

	// Reading credentials should work again after the database
	// password has been rolled back.
	credResp, err = lb.HandleRequest(namespace.RootContext(nil), credReq)
	if err != nil || (credResp != nil && credResp.IsError()) {
		t.Fatalf("err:%s resp:%v\n", err, credResp)
	}
}

// Tests that the WAL rollback function does not roll back the database password.
// The database password should not be rolled back when:
//   - A WAL entry exists
//   - Password has not been altered on the database
//   - Password has not been updated in storage
func TestBackend_RotateRootCredentials_WAL_no_rollback_1(t *testing.T) {
	_, sys := getClusterPostgresDB(t)

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = sys

	lb, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	defer lb.Cleanup(context.Background())

	cleanup, connURL := postgreshelper.PrepareTestContainer(t)
	defer cleanup()

	connURL = strings.ReplaceAll(connURL, "postgres:secret", "{{username}}:{{password}}")

	// Configure a connection to the database
	data := map[string]interface{}{
		"connection_url": connURL,
		"plugin_name":    "postgresql-database-plugin",
		"allowed_roles":  []string{"plugin-role-test"},
		"username":       databaseUser,
		"password":       defaultPassword,
	}
	resp, err := lb.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Create a role
	data = map[string]interface{}{
		"db_name":             "plugin-test",
		"creation_statements": testRole,
		"max_ttl":             "10m",
	}
	resp, err = lb.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/plugin-role-test",
		Storage:   config.StorageView,
		Data:      data,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Read credentials to verify this initially works
	credReq := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/plugin-role-test",
		Storage:   config.StorageView,
		Data:      make(map[string]interface{}),
	}
	credResp, err := lb.HandleRequest(context.Background(), credReq)
	if err != nil || (credResp != nil && credResp.IsError()) {
		t.Fatalf("err:%s resp:%v\n", err, credResp)
	}

	// Put a WAL entry
	walEntry := NewRotateRootCredentialsWALPasswordEntry("plugin-test", databaseUser, "newSecret", defaultPassword)
	_, err = framework.PutWAL(context.Background(), config.StorageView, rotateRootWALKey, walEntry)
	if err != nil {
		t.Fatal(err)
	}
	assertWALCount(t, config.StorageView, 1, rotateRootWALKey)

	// Trigger an immediate RollbackOperation
	_, err = lb.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.RollbackOperation,
		Path:      "",
		Storage:   config.StorageView,
		Data: map[string]interface{}{
			"immediate": true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertWALCount(t, config.StorageView, 0, rotateRootWALKey)

	// Reading credentials should work
	credResp, err = lb.HandleRequest(namespace.RootContext(nil), credReq)
	if err != nil || (credResp != nil && credResp.IsError()) {
		t.Fatalf("err:%s resp:%v\n", err, credResp)
	}
}

// Tests that the WAL rollback function does not roll back the database password.
// The database password should not be rolled back when:
//   - A WAL entry exists
//   - Password has been altered on the database
//   - Password has been updated in storage
func TestBackend_RotateRootCredentials_WAL_no_rollback_2(t *testing.T) {
	_, sys := getClusterPostgresDB(t)

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

	cleanup, connURL := postgreshelper.PrepareTestContainer(t)
	defer cleanup()

	connURL = strings.ReplaceAll(connURL, "postgres:secret", "{{username}}:{{password}}")

	// Configure a connection to the database
	data := map[string]interface{}{
		"connection_url": connURL,
		"plugin_name":    "postgresql-database-plugin",
		"allowed_roles":  []string{"plugin-role-test"},
		"username":       databaseUser,
		"password":       defaultPassword,
	}
	resp, err := lb.HandleRequest(namespace.RootContext(nil), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/plugin-test",
		Storage:   config.StorageView,
		Data:      data,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Create a role
	data = map[string]interface{}{
		"db_name":             "plugin-test",
		"creation_statements": testRole,
		"max_ttl":             "10m",
	}
	resp, err = lb.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/plugin-role-test",
		Storage:   config.StorageView,
		Data:      data,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	// Read credentials to verify this initially works
	credReq := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/plugin-role-test",
		Storage:   config.StorageView,
		Data:      make(map[string]interface{}),
	}
	credResp, err := lb.HandleRequest(context.Background(), credReq)
	if err != nil || (credResp != nil && credResp.IsError()) {
		t.Fatalf("err:%s resp:%v\n", err, credResp)
	}

	// Get a connection to the database plugin
	dbi, err := dbBackend.GetConnection(context.Background(), config.StorageView, "plugin-test")
	if err != nil {
		t.Fatal(err)
	}

	// Alter the database password
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	updateReq := v5.UpdateUserRequest{
		Username: databaseUser,
		Password: &v5.ChangePassword{
			NewPassword: "newSecret",
		},
	}
	_, err = dbi.database.UpdateUser(ctx, updateReq, false)
	if err != nil {
		t.Fatal(err)
	}

	// Update storage with the new password
	dbConfig, err := dbBackend.DatabaseConfig(context.Background(), config.StorageView,
		"plugin-test")
	if err != nil {
		t.Fatal(err)
	}
	dbConfig.ConnectionDetails["password"] = "newSecret"
	entry, err := logical.StorageEntryJSON("config/plugin-test", dbConfig)
	if err != nil {
		t.Fatal(err)
	}
	err = config.StorageView.Put(context.Background(), entry)
	if err != nil {
		t.Fatal(err)
	}

	// Clear the plugin connection to verify we can connect to the database
	err = dbBackend.ClearConnection("plugin-test")
	if err != nil {
		t.Fatal(err)
	}

	// Reading credentials should work
	credResp, err = lb.HandleRequest(namespace.RootContext(nil), credReq)
	if err != nil || (credResp != nil && credResp.IsError()) {
		t.Fatalf("err:%s resp:%v\n", err, credResp)
	}

	// Put a WAL entry
	walEntry := NewRotateRootCredentialsWALPasswordEntry("plugin-test", databaseUser, "newSecret", defaultPassword)
	_, err = framework.PutWAL(context.Background(), config.StorageView, rotateRootWALKey, walEntry)
	if err != nil {
		t.Fatal(err)
	}
	assertWALCount(t, config.StorageView, 1, rotateRootWALKey)

	// Trigger an immediate RollbackOperation
	_, err = lb.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.RollbackOperation,
		Path:      "",
		Storage:   config.StorageView,
		Data: map[string]interface{}{
			"immediate": true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertWALCount(t, config.StorageView, 0, rotateRootWALKey)

	// Reading credentials should work
	credResp, err = lb.HandleRequest(namespace.RootContext(nil), credReq)
	if err != nil || (credResp != nil && credResp.IsError()) {
		t.Fatalf("err:%s resp:%v\n", err, credResp)
	}
}

// failingInitializeDatabase is a v5.Database mock whose Initialize always
// returns a fixed error. Used to simulate a persistent connection failure
// (wrong credentials, network refused) without a blocking timeout.
type failingInitializeDatabase struct {
	err error
}

func (d *failingInitializeDatabase) Initialize(_ context.Context, _ v5.InitializeRequest) (v5.InitializeResponse, error) {
	return v5.InitializeResponse{}, d.err
}

func (d *failingInitializeDatabase) NewUser(_ context.Context, _ v5.NewUserRequest) (v5.NewUserResponse, error) {
	return v5.NewUserResponse{}, nil
}

func (d *failingInitializeDatabase) UpdateUser(_ context.Context, _ v5.UpdateUserRequest) (v5.UpdateUserResponse, error) {
	return v5.UpdateUserResponse{}, nil
}

func (d *failingInitializeDatabase) DeleteUser(_ context.Context, _ v5.DeleteUserRequest) (v5.DeleteUserResponse, error) {
	return v5.DeleteUserResponse{}, nil
}
func (d *failingInitializeDatabase) Type() (string, error) { return mockV5Type, nil }
func (d *failingInitializeDatabase) Close() error          { return nil }

// TestWalRollback_InitializeTimeout_SkipsRollback verifies that a transient
// errDatabaseInitializeTimeout from GetConnection does not trigger
// rollbackDatabaseCredentials. The timeout only means the database was slow;
// it says nothing about whether credentials were already rotated successfully.
func TestWalRollback_InitializeTimeout_SkipsRollback(t *testing.T) {
	oldTimeout := databaseInitTimeout
	databaseInitTimeout = 25 * time.Millisecond
	defer func() { databaseInitTimeout = oldTimeout }()

	config := logical.TestBackendConfig()
	config.System = &systemViewWrapper{
		SystemView:     config.System,
		builtinFactory: newBlockingInitializeDatabase,
	}
	config.StorageView = &logical.InmemStorage{}

	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup(context.Background())

	entry, err := logical.StorageEntryJSON("config/mydb", &DatabaseConfig{
		AllowedRoles:     []string{"*"},
		VerifyConnection: true,
		PluginName:       mockV5Type,
		ConnectionDetails: map[string]interface{}{
			"password": "original-pass",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := config.StorageView.Put(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	walEntry := &rotateRootCredentialsWAL{
		ConnectionName: "mydb",
		UserName:       "root",
		NewPassword:    "new-pass", // != "original-pass": enters the credential-verification branch
		OldPassword:    "original-pass",
	}
	err = b.walRollback(context.Background(), &logical.Request{Storage: config.StorageView}, rotateRootWALKey, walEntry)
	if !errors.Is(err, errDatabaseInitializeTimeout) {
		t.Fatalf("expected errDatabaseInitializeTimeout to propagate, got: %v", err)
	}
}

// TestWalRollback_ConnectionFailed_TriggersRollback verifies that a
// non-timeout connection failure still reaches rollbackDatabaseCredentials,
// preserving the existing behavior for genuine authentication failures.
func TestWalRollback_ConnectionFailed_TriggersRollback(t *testing.T) {
	connErr := errors.New("connection refused: invalid credentials")
	config := logical.TestBackendConfig()
	config.System = &systemViewWrapper{
		SystemView: config.System,
		builtinFactory: func() (interface{}, error) {
			return &failingInitializeDatabase{err: connErr}, nil
		},
	}
	config.StorageView = &logical.InmemStorage{}

	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup(context.Background())

	entry, err := logical.StorageEntryJSON("config/mydb", &DatabaseConfig{
		AllowedRoles:     []string{"*"},
		VerifyConnection: true,
		PluginName:       mockV5Type,
		ConnectionDetails: map[string]interface{}{
			"password": "original-pass",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := config.StorageView.Put(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	walEntry := &rotateRootCredentialsWAL{
		ConnectionName: "mydb",
		UserName:       "root",
		NewPassword:    "new-pass",
		OldPassword:    "original-pass",
	}
	err = b.walRollback(context.Background(), &logical.Request{Storage: config.StorageView}, rotateRootWALKey, walEntry)
	// A non-timeout failure must not be swallowed: rollbackDatabaseCredentials is
	// called and its error (from a second failed GetConnectionWithConfig) propagates.
	if errors.Is(err, errDatabaseInitializeTimeout) {
		t.Fatal("timeout sentinel must not appear for a non-timeout connection failure")
	}
	if err == nil {
		t.Fatal("expected rollbackDatabaseCredentials to return an error")
	}
}

// configuredUpdateUserDatabase is a v5.Database mock whose Initialize always
// succeeds and whose UpdateUser returns a configurable error (nil for success).
type configuredUpdateUserDatabase struct {
	updateUserErr error
}

func (d *configuredUpdateUserDatabase) Initialize(_ context.Context, _ v5.InitializeRequest) (v5.InitializeResponse, error) {
	return v5.InitializeResponse{}, nil
}

func (d *configuredUpdateUserDatabase) NewUser(_ context.Context, _ v5.NewUserRequest) (v5.NewUserResponse, error) {
	return v5.NewUserResponse{}, nil
}

func (d *configuredUpdateUserDatabase) UpdateUser(_ context.Context, _ v5.UpdateUserRequest) (v5.UpdateUserResponse, error) {
	return v5.UpdateUserResponse{}, d.updateUserErr
}

func (d *configuredUpdateUserDatabase) DeleteUser(_ context.Context, _ v5.DeleteUserRequest) (v5.DeleteUserResponse, error) {
	return v5.DeleteUserResponse{}, nil
}

func (d *configuredUpdateUserDatabase) Type() (string, error) { return mockV5Type, nil }
func (d *configuredUpdateUserDatabase) Close() error          { return nil }

// generateTestRSAPrivateKeyPEM generates a 2048-bit RSA private key in
// PKCS#8 PEM format suitable for use with derivePublicKeyFromPrivateKeyPEM.
func generateTestRSAPrivateKeyPEM(t *testing.T) string {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	}))
}

// TestWalRollback_PrivateKey_RotationCompleted_NoRollback verifies that when
// the private key already stored matches the WAL new key, walRollback returns
// nil immediately — the rotation completed and the WAL simply wasn't deleted.
func TestWalRollback_PrivateKey_RotationCompleted_NoRollback(t *testing.T) {
	config := logical.TestBackendConfig()
	config.System = &systemViewWrapper{SystemView: config.System}
	config.StorageView = &logical.InmemStorage{}

	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup(context.Background())

	entry, err := logical.StorageEntryJSON("config/mydb", &DatabaseConfig{
		AllowedRoles: []string{"*"},
		PluginName:   mockV5Type,
		ConnectionDetails: map[string]interface{}{
			// Stored key matches WAL new key: rotation completed, WAL not yet GC'd.
			"private_key": newPrivateKey,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := config.StorageView.Put(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	walEntry := &rotateRootCredentialsWAL{
		ConnectionName: "mydb",
		UserName:       "root",
		NewPrivateKey:  newPrivateKey,
		OldPrivateKey:  "old-private-key-pem",
	}
	err = b.walRollback(context.Background(), &logical.Request{Storage: config.StorageView}, rotateRootWALKey, walEntry)
	if err != nil {
		t.Fatalf("expected no error when rotation already completed, got: %v", err)
	}
}

// TestWalRollback_PrivateKey_ConnectionFails_ReturnsError verifies that when
// connecting with the WAL new private key fails, the error propagates.
func TestWalRollback_PrivateKey_ConnectionFails_ReturnsError(t *testing.T) {
	connErr := errors.New("JWT token rejected: invalid key pair")
	config := logical.TestBackendConfig()
	config.System = &systemViewWrapper{
		SystemView: config.System,
		builtinFactory: func() (interface{}, error) {
			return &failingInitializeDatabase{err: connErr}, nil
		},
	}
	config.StorageView = &logical.InmemStorage{}

	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup(context.Background())

	oldPrivateKey := generateTestRSAPrivateKeyPEM(t)

	entry, err := logical.StorageEntryJSON("config/mydb", &DatabaseConfig{
		AllowedRoles:     []string{"*"},
		VerifyConnection: true,
		PluginName:       mockV5Type,
		ConnectionDetails: map[string]interface{}{
			// Stored key does not match new key: out-of-sync, triggers rollback path.
			"private_key": "old-private-key-pem",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := config.StorageView.Put(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	walEntry := &rotateRootCredentialsWAL{
		ConnectionName: "mydb",
		UserName:       "root",
		NewPrivateKey:  "new-private-key-pem",
		OldPrivateKey:  oldPrivateKey,
	}
	err = b.walRollback(context.Background(), &logical.Request{Storage: config.StorageView}, rotateRootWALKey, walEntry)
	if err == nil {
		t.Fatal("expected connection error to propagate, got nil")
	}
}

// TestWalRollback_PrivateKey_UpdateUserFails_ReturnsError verifies that a
// UpdateUser error propagates so the WAL framework can retry.
func TestWalRollback_PrivateKey_UpdateUserFails_ReturnsError(t *testing.T) {
	updateErr := errors.New("internal server error")

	config := logical.TestBackendConfig()
	config.System = &systemViewWrapper{
		SystemView: config.System,
		builtinFactory: func() (interface{}, error) {
			return &configuredUpdateUserDatabase{updateUserErr: updateErr}, nil
		},
	}
	config.StorageView = &logical.InmemStorage{}

	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup(context.Background())

	oldPrivateKey := generateTestRSAPrivateKeyPEM(t)

	entry, err := logical.StorageEntryJSON("config/mydb", &DatabaseConfig{
		AllowedRoles:     []string{"*"},
		VerifyConnection: true,
		PluginName:       mockV5Type,
		ConnectionDetails: map[string]interface{}{
			"private_key": "old-private-key-pem",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := config.StorageView.Put(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	walEntry := &rotateRootCredentialsWAL{
		ConnectionName: "mydb",
		UserName:       "root",
		NewPrivateKey:  "new-private-key-pem",
		OldPrivateKey:  oldPrivateKey,
	}
	err = b.walRollback(context.Background(), &logical.Request{Storage: config.StorageView}, rotateRootWALKey, walEntry)
	if err == nil {
		t.Fatal("expected UpdateUser error to propagate, got nil")
	}
}

// TestWalRollback_PrivateKey_RollbackSucceeds verifies the happy path: when the
// stored key is out-of-sync with the WAL new key and UpdateUser succeeds,
// walRollback returns nil indicating the rollback completed successfully.
func TestWalRollback_PrivateKey_RollbackSucceeds(t *testing.T) {
	config := logical.TestBackendConfig()
	config.System = &systemViewWrapper{
		SystemView: config.System,
		builtinFactory: func() (interface{}, error) {
			return &configuredUpdateUserDatabase{updateUserErr: nil}, nil
		},
	}
	config.StorageView = &logical.InmemStorage{}

	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup(context.Background())

	oldPrivateKey := generateTestRSAPrivateKeyPEM(t)

	entry, err := logical.StorageEntryJSON("config/mydb", &DatabaseConfig{
		AllowedRoles:     []string{"*"},
		VerifyConnection: true,
		PluginName:       mockV5Type,
		ConnectionDetails: map[string]interface{}{
			"private_key": "old-private-key-pem",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := config.StorageView.Put(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	walEntry := &rotateRootCredentialsWAL{
		ConnectionName: "mydb",
		UserName:       "root",
		NewPrivateKey:  "new-private-key-pem",
		OldPrivateKey:  oldPrivateKey,
	}
	err = b.walRollback(context.Background(), &logical.Request{Storage: config.StorageView}, rotateRootWALKey, walEntry)
	if err != nil {
		t.Fatalf("expected successful rollback, got: %v", err)
	}
}

// TestWalRollback_PrivateKey_SnowflakeJWTError_TreatsAsNoOp verifies the
// crash-before-UpdateUser safety path: when UpdateUser returns a Snowflake
// 390144 JWT error, the new private key was never registered with Snowflake,
// so the system is already consistent with the old key. walRollback must
// return nil to cleanly delete the WAL rather than retrying indefinitely.
func TestWalRollback_PrivateKey_SnowflakeJWTError_TreatsAsNoOp(t *testing.T) {
	// Simulate the error Snowflake returns when the JWT is signed with a key
	// it has never seen. The error crosses the gRPC plugin boundary as a plain
	// string, so it is matched via strings.Contains against the error code.
	jwtErr := errors.New("390144 (08001): JWT token is invalid")

	config := logical.TestBackendConfig()
	config.System = &systemViewWrapper{
		SystemView: config.System,
		builtinFactory: func() (interface{}, error) {
			return &configuredUpdateUserDatabase{updateUserErr: jwtErr}, nil
		},
	}
	config.StorageView = &logical.InmemStorage{}

	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup(context.Background())

	oldPrivateKey := generateTestRSAPrivateKeyPEM(t)

	entry, err := logical.StorageEntryJSON("config/mydb", &DatabaseConfig{
		AllowedRoles:     []string{"*"},
		VerifyConnection: true,
		PluginName:       mockV5Type,
		ConnectionDetails: map[string]interface{}{
			// Stored key does not match new key: rollback path is entered.
			"private_key": "old-private-key-pem",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := config.StorageView.Put(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	walEntry := &rotateRootCredentialsWAL{
		ConnectionName: "mydb",
		UserName:       "root",
		NewPrivateKey:  "new-private-key-pem",
		OldPrivateKey:  oldPrivateKey,
	}
	err = b.walRollback(context.Background(), &logical.Request{Storage: config.StorageView}, rotateRootWALKey, walEntry)
	if err != nil {
		t.Fatalf("expected 390144 JWT error to be treated as no-op (nil), got: %v", err)
	}
}
