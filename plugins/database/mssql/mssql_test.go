package mssql

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/plugins/helper/database/connutil"
)

var (
	testMSQLImagePull sync.Once
)

func TestMSSQL_Initialize(t *testing.T) {
	if os.Getenv("MSSQL_URL") == "" || os.Getenv("VAULT_ACC") != "1" {
		return
	}
	connURL := os.Getenv("MSSQL_URL")

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	dbRaw, _ := New()
	db := dbRaw.(*MSSQL)

	err := db.Initialize(connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	connProducer := db.ConnectionProducer.(*connutil.SQLConnectionProducer)
	if !connProducer.Initialized {
		t.Fatal("Database should be initalized")
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

	err = db.Initialize(connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestMSSQL_CreateUser(t *testing.T) {
	if os.Getenv("MSSQL_URL") == "" || os.Getenv("VAULT_ACC") != "1" {
		return
	}
	connURL := os.Getenv("MSSQL_URL")

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	dbRaw, _ := New()
	db := dbRaw.(*MSSQL)
	err := db.Initialize(connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	// Test with no configured Creation Statememt
	_, _, err = db.CreateUser(dbplugin.Statements{}, usernameConfig, time.Now().Add(time.Minute))
	if err == nil {
		t.Fatal("Expected error when no creation statement is provided")
	}

	statements := dbplugin.Statements{
		CreationStatements: testMSSQLRole,
	}

	username, password, err := db.CreateUser(statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
}

func TestMSSQL_RevokeUser(t *testing.T) {
	if os.Getenv("MSSQL_URL") == "" || os.Getenv("VAULT_ACC") != "1" {
		return
	}
	connURL := os.Getenv("MSSQL_URL")

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	dbRaw, _ := New()
	db := dbRaw.(*MSSQL)
	err := db.Initialize(connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		CreationStatements: testMSSQLRole,
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	username, password, err := db.CreateUser(statements, usernameConfig, time.Now().Add(2*time.Second))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test default revoke statememts
	err = db.RevokeUser(statements, username)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err == nil {
		t.Fatal("Credentials were not revoked")
	}

	username, password, err = db.CreateUser(statements, usernameConfig, time.Now().Add(2*time.Second))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test custom revoke statememt
	statements.RevocationStatements = testMSSQLDrop
	err = db.RevokeUser(statements, username)
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
