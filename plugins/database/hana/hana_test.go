package hana

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/plugins/helper/database/connutil"
)

func TestHANA_Initialize(t *testing.T) {
	if os.Getenv("HANA_URL") == "" || os.Getenv("VAULT_ACC") != "1" {
		t.SkipNow()
	}
	connURL := os.Getenv("HANA_URL")

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	dbRaw, _ := New()
	db := dbRaw.(*HANA)

	err := db.Initialize(connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	connProducer := db.ConnectionProducer.(*connutil.SQLConnectionProducer)
	if !connProducer.Initialized {
		t.Fatal("Database should be initialized")
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

// this test will leave a lingering user on the system
func TestHANA_CreateUser(t *testing.T) {
	if os.Getenv("HANA_URL") == "" || os.Getenv("VAULT_ACC") != "1" {
		t.SkipNow()
	}
	connURL := os.Getenv("HANA_URL")

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	dbRaw, _ := New()
	db := dbRaw.(*HANA)

	err := db.Initialize(connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test-test",
		RoleName:    "test-test",
	}

	// Test with no configured Creation Statememt
	_, _, err = db.CreateUser(dbplugin.Statements{}, usernameConfig, time.Now().Add(time.Hour))
	if err == nil {
		t.Fatal("Expected error when no creation statement is provided")
	}

	statements := dbplugin.Statements{
		CreationStatements: testHANARole,
	}

	username, password, err := db.CreateUser(statements, usernameConfig, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
}

func TestHANA_RevokeUser(t *testing.T) {
	if os.Getenv("HANA_URL") == "" || os.Getenv("VAULT_ACC") != "1" {
		t.SkipNow()
	}
	connURL := os.Getenv("HANA_URL")

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	dbRaw, _ := New()
	db := dbRaw.(*HANA)

	err := db.Initialize(connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		CreationStatements: testHANARole,
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test-test",
		RoleName:    "test-test",
	}

	// Test default revoke statememts
	username, password, err := db.CreateUser(statements, usernameConfig, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if err = testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	err = db.RevokeUser(statements, username)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if err := testCredsExist(t, connURL, username, password); err == nil {
		t.Fatal("Credentials were not revoked")
	}

	// Test custom revoke statememt
	username, password, err = db.CreateUser(statements, usernameConfig, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if err = testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	statements.RevocationStatements = testHANADrop
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
	connURL = fmt.Sprintf("hdb://%s:%s@%s", username, password, parts[1])
	db, err := sql.Open("hdb", connURL)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}

const testHANARole = `
CREATE USER {{name}} PASSWORD {{password}} VALID UNTIL '{{expiration}}';`

const testHANADrop = `
DROP USER {{name}} CASCADE;`
