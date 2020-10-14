package hana

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/database/helper/credsutil"
	"github.com/hashicorp/vault/sdk/database/newdbplugin"
)

func TestHANA_Initialize(t *testing.T) {
	if os.Getenv("HANA_URL") == "" || os.Getenv("VAULT_ACC") != "1" {
		t.SkipNow()
	}
	connURL := os.Getenv("HANA_URL")

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	initReq := newdbplugin.InitializeRequest{
		Config:           connectionDetails,
		VerifyConnection: true,
	}

	db := new()
	_, err := db.Initialize(context.Background(), initReq)
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
}

// this test will leave a lingering user on the system
func TestHANA_NewUser(t *testing.T) {
	if os.Getenv("HANA_URL") == "" || os.Getenv("VAULT_ACC") != "1" {
		t.SkipNow()
	}
	connURL := os.Getenv("HANA_URL")

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	initReq := newdbplugin.InitializeRequest{
		Config:           connectionDetails,
		VerifyConnection: true,
	}

	db := new()
	_, err := db.Initialize(context.Background(), initReq)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	password, err := credsutil.RandomAlphaNumeric(32, true)
	if err != nil {
		t.Fatalf("failed to generate password: %s", err)
	}
	password = strings.Replace(password, "-", "_", -1)

	req := newdbplugin.NewUserRequest{
		UsernameConfig: newdbplugin.UsernameMetadata{
			DisplayName: "test-test",
			RoleName:    "test-test",
		},
		Statements: newdbplugin.Statements{
			Commands: []string{},
		},
		Password:   password,
		Expiration: time.Now().Add(time.Hour),
	}

	// Test with no configured Creation Statement
	userResp, err := db.NewUser(context.Background(), req)
	if err == nil {
		t.Fatal("Expected error when no creation statement is provided")
	}

	// Add a statement command
	req.Statements.Commands = []string{testHANARole}

	userResp, err = db.NewUser(context.Background(), req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, connURL, userResp.Username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
}

func TestHANA_DeleteUser(t *testing.T) {
	if os.Getenv("HANA_URL") == "" || os.Getenv("VAULT_ACC") != "1" {
		t.SkipNow()
	}
	connURL := os.Getenv("HANA_URL")

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	initReq := newdbplugin.InitializeRequest{
		Config:           connectionDetails,
		VerifyConnection: true,
	}

	db := new()
	_, err := db.Initialize(context.Background(), initReq)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	password, err := credsutil.RandomAlphaNumeric(32, true)
	if err != nil {
		t.Fatalf("failed to generate password: %s", err)
	}
	password = strings.Replace(password, "-", "_", -1)

	newReq := newdbplugin.NewUserRequest{
		UsernameConfig: newdbplugin.UsernameMetadata{
			DisplayName: "test-test",
			RoleName:    "test-test",
		},
		Password: password,
		Statements: newdbplugin.Statements{
			Commands: []string{testHANARole},
		},
		Expiration: time.Now().Add(time.Hour),
	}

	// Test default revoke statements
	userResp, err := db.NewUser(context.Background(), newReq)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if err = testCredsExist(t, connURL, userResp.Username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	delReq := newdbplugin.DeleteUserRequest{
		Username: userResp.Username,
	}

	_, err = db.DeleteUser(context.Background(), delReq)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if err := testCredsExist(t, connURL, userResp.Username, password); err == nil {
		t.Fatal("Credentials were not revoked")
	}

	// Test custom revoke statement
	userResp, err = db.NewUser(context.Background(), newReq)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if err = testCredsExist(t, connURL, userResp.Username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	delReq.Statements.Commands = []string{testHANADrop}
	delReq.Username = userResp.Username
	_, err = db.DeleteUser(context.Background(), delReq)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if err := testCredsExist(t, connURL, userResp.Username, password); err == nil {
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
CREATE USER {{name}} PASSWORD {{password}} NO FORCE_FIRST_PASSWORD_CHANGE VALID UNTIL '{{expiration}}';`

const testHANADrop = `
DROP USER {{name}} CASCADE;`
