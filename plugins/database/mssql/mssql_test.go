package mssql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	mssqlhelper "github.com/hashicorp/vault/helper/testhelpers/mssql"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/helper/dbtxn"
)

func TestMSSQL_Initialize(t *testing.T) {
	cleanup, connURL := mssqlhelper.PrepareMSSQLTestContainer(t)
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
	cleanup, connURL := mssqlhelper.PrepareMSSQLTestContainer(t)
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
	cleanup, connURL := mssqlhelper.PrepareMSSQLTestContainer(t)
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

func TestMSSQL_SetCredentials_missingArgs(t *testing.T) {
	type testCase struct {
		statements dbplugin.Statements
		userConfig dbplugin.StaticUserConfig
	}

	tests := map[string]testCase{
		"empty rotation statements": {
			statements: dbplugin.Statements{
				Rotation: nil,
			},
			userConfig: dbplugin.StaticUserConfig{
				Username: "testuser",
				Password: "password",
			},
		},
		"empty username": {
			statements: dbplugin.Statements{
				Rotation: []string{`
					ALTER LOGIN [{{username}}] WITH PASSWORD = '{{password}}';`,
				},
			},
			userConfig: dbplugin.StaticUserConfig{
				Username: "",
				Password: "password",
			},
		},
		"empty password": {
			statements: dbplugin.Statements{
				Rotation: []string{`
					ALTER LOGIN [{{username}}] WITH PASSWORD = '{{password}}';`,
				},
			},
			userConfig: dbplugin.StaticUserConfig{
				Username: "testuser",
				Password: "",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db := new()

			username, password, err := db.SetCredentials(context.Background(), test.statements, test.userConfig)
			if err == nil {
				t.Fatalf("expected err, got nil")
			}
			if username != "" {
				t.Fatalf("expected empty username, got [%s]", username)
			}
			if password != "" {
				t.Fatalf("expected empty password, got [%s]", password)
			}
		})
	}
}

func TestMSSQL_SetCredentials(t *testing.T) {
	type testCase struct {
		rotationStmts []string
	}

	tests := map[string]testCase{
		"empty rotation statements": {
			rotationStmts: []string{},
		}, "username rotation": {
			rotationStmts: []string{`
				ALTER LOGIN [{{username}}] WITH PASSWORD = '{{password}}';`,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cleanup, connURL := mssqlhelper.PrepareMSSQLTestContainer(t)
			defer cleanup()

			connectionDetails := map[string]interface{}{
				"connection_url": connURL,
			}

			db := new()

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			_, err := db.Init(ctx, connectionDetails, true)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			dbUser := "vaultstatictest"
			initPassword := "p4$sw0rd"
			createTestMSSQLUser(t, connURL, dbUser, initPassword, testMSSQLLogin)

			if err := testCredsExist(t, connURL, dbUser, initPassword); err != nil {
				t.Fatalf("Could not connect with initial credentials: %s", err)
			}

			statements := dbplugin.Statements{
				Rotation: test.rotationStmts,
			}

			newPassword, err := db.GenerateCredentials(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			usernameConfig := dbplugin.StaticUserConfig{
				Username: dbUser,
				Password: newPassword,
			}

			username, password, err := db.SetCredentials(ctx, statements, usernameConfig)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if err := testCredsExist(t, connURL, username, password); err != nil {
				t.Fatalf("Could not connect with new credentials: %s", err)
			}

			if err := testCredsExist(t, connURL, username, initPassword); err == nil {
				t.Fatalf("Should not be able to connect with initial credentials")
			}

		})
	}

}

func TestMSSQL_RevokeUser(t *testing.T) {
	cleanup, connURL := mssqlhelper.PrepareMSSQLTestContainer(t)
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

func createTestMSSQLUser(t *testing.T, connURL string, username, password, query string) {

	db, err := sql.Open("mssql", connURL)
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
	if err := dbtxn.ExecuteTxQuery(ctx, tx, m, query); err != nil {
		t.Fatal(err)
	}
	// Commit the transaction
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}

const testMSSQLRole = `
CREATE LOGIN [{{name}}] WITH PASSWORD = '{{password}}';
CREATE USER [{{name}}] FOR LOGIN [{{name}}];
GRANT SELECT, INSERT, UPDATE, DELETE ON SCHEMA::dbo TO [{{name}}];`

const testMSSQLDrop = `
DROP USER [{{name}}];
DROP LOGIN [{{name}}];
`

const testMSSQLLogin = `
CREATE LOGIN [{{name}}] WITH PASSWORD = '{{password}}';
`
