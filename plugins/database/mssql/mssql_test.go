package mssql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers/docker"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/ory/dockertest"
)

func prepareTestContainer(t *testing.T) (func(), string) {
	if os.Getenv("MSSQL_URL") != "" {
		return func() {}, os.Getenv("MSSQL_URL")
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	resource, err := pool.Run("mcr.microsoft.com/mssql/server", "2017-latest-ubuntu", []string{"ACCEPT_EULA=Y", "SA_PASSWORD=yourStrong(!)Password"})
	if err != nil {
		t.Fatalf("Could not start local MSSQL docker container: %s", err)
	}

	cleanup := func() {
		docker.CleanupResource(t, pool, resource)
	}

	retURL := fmt.Sprintf("sqlserver://sa:yourStrong(!)Password@127.0.0.1:%s", resource.GetPort("1433/tcp"))

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		var err error
		var db *sql.DB
		db, err = sql.Open("mssql", retURL)
		if err != nil {
			return err
		}
		defer db.Close()
		return db.Ping()
	}); err != nil {
		cleanup()
		t.Fatalf("Could not connect to MSSQL docker container: %s", err)
	}

	return cleanup, retURL
}

func TestMSSQL_Initialize(t *testing.T) {
	cleanup, connURL := prepareTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !db.Initialized {
		t.Fatal("Database should be initialized")
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Test decoding a string value for max_open_connections
	connectionDetails = map[string]interface{}{
		"connection_url":       connURL,
		"max_open_connections": "5",
	}

	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestMSSQL_CreateUser(t *testing.T) {
	cleanup, connURL := prepareTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	// Test with no configured Creation Statement
	_, _, err = db.CreateUser(context.Background(), dbplugin.Statements{}, usernameConfig, time.Now().Add(time.Minute))
	if err == nil {
		t.Fatal("Expected error when no creation statement is provided")
	}

	statements := dbplugin.Statements{
		Creation: []string{testMSSQLRole},
	}

	username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
}

func TestMSSQL_RotateRootCredentials(t *testing.T) {
	cleanup, connURL := prepareTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
		"username":       "sa",
		"password":       "yourStrong(!)Password",
	}

	db := new()

	connProducer := db.SQLConnectionProducer

	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !connProducer.Initialized {
		t.Fatal("Database should be initialized")
	}

	newConf, err := db.RotateRootCredentials(context.Background(), nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if newConf["password"] == "yourStrong(!)Password" {
		t.Fatal("password was not updated")
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestMSSQL_RevokeUser(t *testing.T) {
	cleanup, connURL := prepareTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		Creation: []string{testMSSQLRole},
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(2*time.Second))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test default revoke statements
	err = db.RevokeUser(context.Background(), statements, username)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err == nil {
		t.Fatal("Credentials were not revoked")
	}

	username, password, err = db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(2*time.Second))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test custom revoke statement
	statements.Revocation = []string{testMSSQLDrop}
	err = db.RevokeUser(context.Background(), statements, username)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err == nil {
		t.Fatal("Credentials were not revoked")
	}
}

func testCredsExist(t testing.TB, connURL, username, password string) error {
	// Log in with the new creds
	parts := strings.Split(connURL, "@")
	connURL = fmt.Sprintf("sqlserver://%s:%s@%s", username, password, parts[1])
	db, err := sql.Open("mssql", connURL)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}

const testMSSQLRole = `
CREATE LOGIN [{{name}}] WITH PASSWORD = '{{password}}';
CREATE USER [{{name}}] FOR LOGIN [{{name}}];
GRANT SELECT, INSERT, UPDATE, DELETE ON SCHEMA::dbo TO [{{name}}];`

const testMSSQLDrop = `
DROP USER [{{name}}];
DROP LOGIN [{{name}}];
`
