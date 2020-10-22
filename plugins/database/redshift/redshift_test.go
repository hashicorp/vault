package redshift

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
	dbplugin "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/helper/dbtxn"
	"github.com/lib/pq"
)

/*
To run these sets of acceptance tests, you must pre-configure a Redshift cluster
in AWS and ensure the machine running these tests has network access to it.

Once the redshift cluster is running, you can pass the admin username and password
as environment variables to be used to run these tests. Note that these tests
will create users on your redshift cluster and currently do not clean up after
themselves.

Do not run this test suite against a production Redshift cluster.

Configuration:

		REDSHIFT_URL=my-redshift-url.region.redshift.amazonaws.com:5439/database-name
		REDSHIFT_USER=my-redshift-admin-user
		REDSHIFT_PASSWORD=my-redshift-admin-password
		VAULT_ACC=<unset || 1> # This must be set to run any of the tests in this test suite
*/

var (
	keyRedshiftURL      = "REDSHIFT_URL"
	keyRedshiftUser     = "REDSHIFT_USER"
	keyRedshiftPassword = "REDSHIFT_PASSWORD"

	vaultACC = "VAULT_ACC"
)

func interpolateConnectionURL(url, user, password string) string {
	return fmt.Sprintf("postgres://%s:%s@%s", user, password, url)
}

func redshiftEnv() (connURL string, url string, user string, password string, errEmpty error) {
	if url = os.Getenv(keyRedshiftURL); url == "" {
		return "", "", "", "", fmt.Errorf("%s environment variable required", keyRedshiftURL)
	}

	if user = os.Getenv(keyRedshiftUser); url == "" {
		return "", "", "", "", fmt.Errorf("%s environment variable required", keyRedshiftUser)
	}

	if password = os.Getenv(keyRedshiftPassword); url == "" {
		return "", "", "", "", fmt.Errorf("%s environment variable required", keyRedshiftPassword)
	}

	connURL = interpolateConnectionURL(url, user, password)
	return connURL, url, user, password, nil
}

func TestRedshift_Initialize(t *testing.T) {
	if os.Getenv(vaultACC) != "1" {
		t.SkipNow()
	}

	connURL, _, _, _, err := redshiftEnv()
	if err != nil {
		t.Fatal(err)
	}

	connectionDetails := map[string]interface{}{
		"connection_url":       connURL,
		"max_open_connections": 5,
	}

	db := newRedshift()
	_, err = db.Init(context.Background(), connectionDetails, true)
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
		"max_open_connections": "73",
	}

	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if db.MaxOpenConnections != 73 {
		t.Fatalf("Expected max_open_connections to be set to 73, but was %d", db.MaxOpenConnections)
	}
}

func TestRedshift_NewUser(t *testing.T) {
	if os.Getenv(vaultACC) != "1" {
		t.SkipNow()
	}

	connURL, url, _, _, err := redshiftEnv()
	if err != nil {
		t.Fatal(err)
	}

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := newRedshift()
	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	usernameConfig := dbplugin.UsernameMetadata{
		DisplayName: "test",
		RoleName:    "test",
	}

	// Test with no configured Creation Statement
	_, err = db.NewUser(context.Background(), dbplugin.NewUserRequest{})
	if err == nil {
		t.Fatal("Expected error when no creation statement is provided")
	}

	const password = "SuperSecurePa55w0rd!"
	req := dbplugin.NewUserRequest{
		UsernameConfig: usernameConfig,
		Password:       password,
		Statements: dbplugin.Statements{
			Commands: []string{testRedshiftRole},
		},
		Expiration: time.Now().Add(5 * time.Minute),
	}

	resp, err := db.NewUser(context.Background(), req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	username := resp.Username

	if err = testCredsExist(t, url, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s\n%s:%s", err, username, password)
	}

	req = dbplugin.NewUserRequest{
		UsernameConfig: usernameConfig,
		Password:       "SuperSecurePa55w0rd!",
		Statements: dbplugin.Statements{
			Commands: []string{testRedshiftReadOnlyRole},
		},
		Expiration: time.Now().Add(5 * time.Minute),
	}
	resp, err = db.NewUser(context.Background(), req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Sleep to make sure we haven't expired if granularity is only down to the second
	time.Sleep(2 * time.Second)

	if err = testCredsExist(t, url, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
}

func TestRedshift_UpdateUser_Expiration(t *testing.T) {
	if os.Getenv(vaultACC) != "1" {
		t.SkipNow()
	}

	connURL, url, _, _, err := redshiftEnv()
	if err != nil {
		t.Fatal(err)
	}

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := newRedshift()
	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	usernameConfig := dbplugin.UsernameMetadata{
		DisplayName: "test",
		RoleName:    "test",
	}

	const password = "SuperSecurePa55w0rd!"
	newReq := dbplugin.NewUserRequest{
		UsernameConfig: usernameConfig,
		Password:       password,
		Statements:     dbplugin.Statements{Commands: []string{testRedshiftRole}},
		Expiration:     time.Now().Add(5 * time.Second),
	}

	newResp, err := db.NewUser(context.Background(), newReq)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	username := newResp.Username

	if err = testCredsExist(t, url, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test default renew statement
	updateReq := dbplugin.UpdateUserRequest{
		Username: username,
		Expiration: &dbplugin.ChangeExpiration{
			NewExpiration: time.Now().Add(time.Minute),
		},
	}

	_, err = db.UpdateUser(context.Background(), updateReq)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Sleep longer than the initial expiration time
	time.Sleep(5 * time.Second)

	if err = testCredsExist(t, url, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Now test the same again with explicitly set renew statements
	newReq = dbplugin.NewUserRequest{
		UsernameConfig: usernameConfig,
		Password:       password,
		Statements:     dbplugin.Statements{Commands: []string{testRedshiftRole}},
		Expiration:     time.Now().Add(5 * time.Second),
	}

	newResp, err = db.NewUser(context.Background(), newReq)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	username = newResp.Username

	if err = testCredsExist(t, url, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	updateReq = dbplugin.UpdateUserRequest{
		Username: username,
		Expiration: &dbplugin.ChangeExpiration{
			NewExpiration: time.Now().Add(time.Minute),
			Statements: dbplugin.Statements{
				Commands: []string{defaultRenewSQL},
			},
		},
	}

	_, err = db.UpdateUser(context.Background(), updateReq)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Sleep longer than the initial expiration time
	time.Sleep(5 * time.Second)

	if err = testCredsExist(t, url, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

}

func TestRedshift_DeleteUser(t *testing.T) {
	if os.Getenv(vaultACC) != "1" {
		t.SkipNow()
	}

	connURL, url, _, _, err := redshiftEnv()
	if err != nil {
		t.Fatal(err)
	}

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := newRedshift()
	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	usernameConfig := dbplugin.UsernameMetadata{
		DisplayName: "test",
		RoleName:    "test",
	}

	const password = "SuperSecretPa55word!"
	newRequest := dbplugin.NewUserRequest{
		UsernameConfig: usernameConfig,
		Statements:     dbplugin.Statements{Commands: []string{testRedshiftRole}},
		Password:       password,
		Expiration:     time.Now().Add(2 * time.Second),
	}

	newResponse, err := db.NewUser(context.Background(), newRequest)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	username := newResponse.Username

	if err = testCredsExist(t, url, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test default revoke statements
	_, err = db.DeleteUser(context.Background(), dbplugin.DeleteUserRequest{
		Username: username,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, url, username, password); err == nil {
		t.Fatal("Credentials were not revoked")
	}

	newRequest.Expiration = time.Now().Add(2 * time.Second)
	newResponse, err = db.NewUser(context.Background(), newRequest)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	username = newResponse.Username

	if err = testCredsExist(t, url, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test custom revoke statements
	req := dbplugin.DeleteUserRequest{
		Username: username,
		Statements: dbplugin.Statements{
			Commands: []string{defaultRedshiftRevocationSQL},
		},
	}
	_, err = db.DeleteUser(context.Background(), req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, url, username, password); err == nil {
		t.Fatal("Credentials were not revoked")
	}
}

func TestRedshift_UpdateUser_Password(t *testing.T) {
	if os.Getenv(vaultACC) != "1" {
		t.SkipNow()
	}

	connURL, url, _, _, err := redshiftEnv()
	if err != nil {
		t.Fatal(err)
	}

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	// create the database user
	uid, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	dbUser := "vaultstatictest-" + fmt.Sprintf("%s", uid)
	createTestPGUser(t, connURL, dbUser, "1Password", testRoleStaticCreate)

	db := newRedshift()
	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	const password1 = "MyTemporaryUserPassword1!"
	const password2 = "MyTemporaryUserPassword2!"

	updateReq := dbplugin.UpdateUserRequest{
		Username: dbUser,
		Password: &dbplugin.ChangePassword{
			NewPassword: password1,
			Statements:  dbplugin.Statements{},
		},
	}

	// Test with no configured Rotation Statement
	_, err = db.UpdateUser(context.Background(), updateReq)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, url, dbUser, password1); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test with explicitly set rotation statement
	updateReq = dbplugin.UpdateUserRequest{
		Username: dbUser,
		Password: &dbplugin.ChangePassword{
			NewPassword: password2,
			Statements: dbplugin.Statements{
				Commands: []string{testRedshiftStaticRoleRotate},
			},
		},
	}
	_, err = db.UpdateUser(context.Background(), updateReq)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, url, dbUser, password2); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
}

func testCredsExist(t testing.TB, url, username, password string) error {
	t.Helper()

	connURL := interpolateConnectionURL(url, username, password)
	db, err := sql.Open("postgres", connURL)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}

const testRedshiftRole = `
CREATE USER "{{name}}" WITH PASSWORD '{{password}}' VALID UNTIL '{{expiration}}'; 
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "{{name}}";
`

const testRedshiftReadOnlyRole = `
CREATE USER "{{name}}" WITH
  PASSWORD '{{password}}'
  VALID UNTIL '{{expiration}}';
GRANT SELECT ON ALL TABLES IN SCHEMA public TO "{{name}}";
`

const defaultRedshiftRevocationSQL = `
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM "{{name}}";
REVOKE USAGE ON SCHEMA public FROM "{{name}}";

DROP USER IF EXISTS "{{name}}";
`

const testRedshiftStaticRole = `
CREATE USER "{{name}}" WITH
  PASSWORD '{{password}}';
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "{{name}}";
`

const testRoleStaticCreate = `
CREATE USER "{{name}}" WITH
  PASSWORD '{{password}}';
`

const testRedshiftStaticRoleRotate = `
ALTER USER "{{name}}" WITH PASSWORD '{{password}}';
`

// This is a copy of a test helper method also found in
// builtin/logical/database/rotation_test.go , and should be moved into a shared
// helper file in the future.
func createTestPGUser(t *testing.T, connURL string, username, password, query string) {
	t.Helper()
	conn, err := pq.ParseURL(connURL)
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("postgres", conn)
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
