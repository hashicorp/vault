package redshift

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"

	"github.com/hashicorp/vault/sdk/database/newdbplugin"

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

The RotateRoot test is potentially destructive in that it will rotate your root
password on your Redshift cluster to an insecure, cleartext password defined in the
test method. Because of this, you must pass TEST_ROTATE_ROOT=1 to enable it explicitly.

Do not run this test suite against a production Redshift cluster.

Configuration:

		REDSHIFT_URL=my-redshift-url.region.redshift.amazonaws.com:5439/database-name
		REDSHIFT_USER=my-redshift-admin-user
		REDSHIFT_PASSWORD=my-redshift-admin-password
		VAULT_ACC=<unset || 1> # This must be set to run any of the tests in this test suite
		TEST_ROTATE_ROOT=<unset || 1> # This must be set to explicitly run the rotate root test
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
	errEmpty = errors.New("err: empty but required env value")

	if url = os.Getenv(keyRedshiftURL); url == "" {
		return "", "", "", "", errEmpty
	}

	if user = os.Getenv(keyRedshiftUser); url == "" {
		return "", "", "", "", errEmpty
	}

	if password = os.Getenv(keyRedshiftPassword); url == "" {
		return "", "", "", "", errEmpty
	}

	return interpolateConnectionURL(url, user, password), url, user, password, nil
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

	db := newRedshift(true)
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
		"max_open_connections": "5",
	}

	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
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

	db := newRedshift(true)
	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	usernameConfig := newdbplugin.UsernameMetadata{
		DisplayName: "test",
		RoleName:    "test",
	}

	// Test with no configured Creation Statement
	_, err = db.NewUser(context.Background(), newdbplugin.NewUserRequest{})
	if err == nil {
		t.Fatal("Expected error when no creation statement is provided")
	}

	req := newdbplugin.NewUserRequest{
		UsernameConfig: usernameConfig,
		Password:       "SuperSecurePa55w0rd!",
		Statements:     newdbplugin.Statements{Commands: []string{testRedshiftRole}},
		Expiration:     time.Now().Add(5 * time.Minute),
	}

	resp, err := db.NewUser(context.Background(), req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, url, resp.Username, req.Password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s\n%s:%s", err, resp.Username, req.Password)
	}

	req.Statements.Commands = []string{testRedshiftReadOnlyRole}
	resp, err = db.NewUser(context.Background(), req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Sleep to make sure we haven't expired if granularity is only down to the second
	time.Sleep(2 * time.Second)

	if err = testCredsExist(t, url, resp.Username, req.Password); err != nil {
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

	db := newRedshift(true)
	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	usernameConfig := newdbplugin.UsernameMetadata{
		DisplayName: "test",
		RoleName:    "test",
	}

	const password = "SuperSecurePa55w0rd!"
	newReq := newdbplugin.NewUserRequest{
		UsernameConfig: usernameConfig,
		Password:       password,
		Statements:     newdbplugin.Statements{Commands: []string{testRedshiftRole}},
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
	updateReq := newdbplugin.UpdateUserRequest{
		Username: username,
		Expiration: &newdbplugin.ChangeExpiration{
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
	newReq = newdbplugin.NewUserRequest{
		UsernameConfig: usernameConfig,
		Password:       password,
		Statements:     newdbplugin.Statements{Commands: []string{testRedshiftRole}},
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

	updateReq = newdbplugin.UpdateUserRequest{
		Username: username,
		Expiration: &newdbplugin.ChangeExpiration{
			NewExpiration: time.Now().Add(time.Minute),
			Statements: newdbplugin.Statements{
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

	db := newRedshift(true)
	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	usernameConfig := newdbplugin.UsernameMetadata{
		DisplayName: "test",
		RoleName:    "test",
	}

	const password = "SuperSecretPa55word!"
	newRequest := newdbplugin.NewUserRequest{
		UsernameConfig: usernameConfig,
		Statements:     newdbplugin.Statements{Commands: []string{testRedshiftRole}},
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
	_, err = db.DeleteUser(context.Background(), newdbplugin.DeleteUserRequest{
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
	req := newdbplugin.DeleteUserRequest{
		Username: username,
		Statements: newdbplugin.Statements{
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

	db := newRedshift(true)
	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	const password1 = "MyTemporaryRootPassword1!"
	const password2 = "MyTemporaryRootPassword2!"

	updateReq := newdbplugin.UpdateUserRequest{
		Username: dbUser,
		Password: &newdbplugin.ChangePassword{
			NewPassword: password1,
			Statements:  newdbplugin.Statements{},
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
	updateReq = newdbplugin.UpdateUserRequest{
		Username: dbUser,
		Password: &newdbplugin.ChangePassword{
			NewPassword: password2,
			Statements: newdbplugin.Statements{
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

func TestRedshift_UpdateUser_RootPassword(t *testing.T) {
	// Extra precaution is taken for rotating root creds because it's assumed that this
	// test will run against a live redshift cluster. It will try to set the password
	// back to the old password at the end, but only as a best effort.
	//
	// To run this test you must pass TEST_ROTATE_ROOT=1

	if os.Getenv(vaultACC) != "1" || os.Getenv("TEST_ROTATE_ROOT") != "1" {
		t.SkipNow()
	}

	connURL, url, adminUser, adminPassword, err := redshiftEnv()
	if err != nil {
		t.Fatal(err)
	}

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
		"username":       adminUser,
		"password":       adminPassword,
	}

	db := newRedshift(true)

	connProducer := db.SQLConnectionProducer

	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !connProducer.Initialized {
		t.Fatal("Database should be initialized")
	}

	const tempPassword = "MyTemporaryRootPassword1!"
	updateReq := newdbplugin.UpdateUserRequest{
		Username: adminUser,
		Password: &newdbplugin.ChangePassword{
			NewPassword: tempPassword,
		},
	}
	_, err = db.UpdateUser(context.Background(), updateReq)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if err = testCredsExist(t, url, adminUser, tempPassword); err != nil {
		t.Fatalf("Failed to test new admin user credentials after rotation: %s", err)
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Switch it back at the end to make the test somewhat less destructive.
	// We re-initialize the whole plugin, which is roughly what the real database
	// backend does for each operation.
	connectionDetails = map[string]interface{}{
		"connection_url": interpolateConnectionURL(url, adminUser, tempPassword),
		"username":       adminUser,
		"password":       tempPassword,
	}

	db = newRedshift(true)
	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	updateReq = newdbplugin.UpdateUserRequest{
		Username: adminUser,
		Password: &newdbplugin.ChangePassword{
			NewPassword: adminPassword,
		},
	}
	_, err = db.UpdateUser(context.Background(), updateReq)
	if err != nil {
		t.Fatalf("Failed to switch password back, err: %v", err)
	}

	if err = testCredsExist(t, url, adminUser, adminPassword); err != nil {
		t.Fatalf("Failed to reset admin user credentials after rotation: %s", err)
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
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
