package redshift

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
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

func redshiftEnv() (url string, user string, password string, errEmpty error) {
	errEmpty = errors.New("err: empty but required env value")

	if url = os.Getenv(keyRedshiftURL); url == "" {
		return "", "", "", errEmpty
	}

	if user = os.Getenv(keyRedshiftUser); url == "" {
		return "", "", "", errEmpty
	}

	if password = os.Getenv(keyRedshiftPassword); url == "" {
		return "", "", "", errEmpty
	}

	url = fmt.Sprintf("postgres://%s:%s@%s", user, password, url)

	return url, user, password, nil
}

func TestPostgreSQL_Initialize(t *testing.T) {
	if os.Getenv(vaultACC) != "1" {
		t.SkipNow()
	}

	url, _, _, err := redshiftEnv()
	if err != nil {
		t.Fatal(err)
	}

	connectionDetails := map[string]interface{}{
		"connection_url":       url,
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
		"connection_url":       url,
		"max_open_connections": "5",
	}

	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

}

func TestPostgreSQL_CreateUser(t *testing.T) {
	if os.Getenv(vaultACC) != "1" {
		t.SkipNow()
	}

	url, _, _, err := redshiftEnv()
	if err != nil {
		t.Fatal(err)
	}

	connectionDetails := map[string]interface{}{
		"connection_url": url,
	}

	db := newRedshift(true)
	_, err = db.Init(context.Background(), connectionDetails, true)
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
		Creation: []string{testRedshiftRole},
	}

	username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, url, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s\n%s:%s", err, username, password)
	}

	statements.Creation = []string{testRedshiftReadOnlyRole}
	username, password, err = db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Sleep to make sure we haven't expired if granularity is only down to the second
	time.Sleep(2 * time.Second)

	if err = testCredsExist(t, url, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
}

func TestPostgreSQL_RenewUser(t *testing.T) {
	if os.Getenv(vaultACC) != "1" {
		t.SkipNow()
	}

	url, _, _, err := redshiftEnv()
	if err != nil {
		t.Fatal(err)
	}

	connectionDetails := map[string]interface{}{
		"connection_url": url,
	}

	db := newRedshift(true)
	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		Creation: []string{testRedshiftRole},
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(2*time.Second))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, url, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	err = db.RenewUser(context.Background(), statements, username, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Sleep longer than the initial expiration time
	time.Sleep(2 * time.Second)

	if err = testCredsExist(t, url, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
	statements.Renewal = []string{defaultRenewSQL}
	username, password, err = db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(2*time.Second))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, url, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	err = db.RenewUser(context.Background(), statements, username, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Sleep longer than the initial expiration time
	time.Sleep(2 * time.Second)

	if err = testCredsExist(t, url, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

}

func TestPostgreSQL_RevokeUser(t *testing.T) {
	if os.Getenv(vaultACC) != "1" {
		t.SkipNow()
	}

	url, _, _, err := redshiftEnv()
	if err != nil {
		t.Fatal(err)
	}

	connectionDetails := map[string]interface{}{
		"connection_url": url,
	}

	db := newRedshift(true)
	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		Creation: []string{testRedshiftRole},
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(2*time.Second))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, url, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test default revoke statements
	err = db.RevokeUser(context.Background(), statements, username)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, url, username, password); err == nil {
		t.Fatal("Credentials were not revoked")
	}

	username, password, err = db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(2*time.Second))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, url, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test custom revoke statements
	statements.Revocation = []string{defaultRedshiftRevocationSQL}
	err = db.RevokeUser(context.Background(), statements, username)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, url, username, password); err == nil {
		t.Fatal("Credentials were not revoked")
	}
}

func TestPostgresSQL_SetCredentials(t *testing.T) {
	if os.Getenv(vaultACC) != "1" {
		t.SkipNow()
	}

	url, _, _, err := redshiftEnv()
	if err != nil {
		t.Fatal(err)
	}

	connectionDetails := map[string]interface{}{
		"connection_url": url,
	}

	// create the database user
	uid, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	dbUser := "vaultstatictest-" + fmt.Sprintf("%s", uid)
	createTestPGUser(t, url, dbUser, "1Password", testRoleStaticCreate)

	db := newRedshift(true)
	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	password, err := db.GenerateCredentials(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	usernameConfig := dbplugin.StaticUserConfig{
		Username: dbUser,
		Password: password,
	}

	// Test with no configured Rotation Statement
	username, password, err := db.SetCredentials(context.Background(), dbplugin.Statements{}, usernameConfig)
	if err == nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		Rotation: []string{testRedshiftStaticRoleRotate},
	}
	// User should not exist, make sure we can create
	username, password, err = db.SetCredentials(context.Background(), statements, usernameConfig)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, url, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// call SetCredentials again, password will change
	newPassword, _ := db.GenerateCredentials(context.Background())
	usernameConfig.Password = newPassword
	username, password, err = db.SetCredentials(context.Background(), statements, usernameConfig)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if password != newPassword {
		t.Fatal("passwords should have changed")
	}

	if err := testCredsExist(t, url, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
}

func TestPostgreSQL_RotateRootCredentials(t *testing.T) {
	/*
		   Extra precaution is taken for rotating root creds because it's assumed that this
		   test will run against a live redshift cluster. This test must run last because
		   it is destructive.

			 To run this test you must pass TEST_ROTATE_ROOT=1
	*/
	if os.Getenv(vaultACC) != "1" || os.Getenv("TEST_ROTATE_ROOT") != "1" {
		t.SkipNow()
	}

	url, adminUser, adminPassword, err := redshiftEnv()
	if err != nil {
		t.Fatal(err)
	}

	connectionDetails := map[string]interface{}{
		"connection_url": url,
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

	newConf, err := db.RotateRootCredentials(context.Background(), nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	fmt.Printf("rotated root credentials, new user/pass:\nusername: %s\npassword: %s\n", newConf["username"], newConf["password"])

	if newConf["password"] == adminPassword {
		t.Fatal("password was not updated")
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testCredsExist(t testing.TB, connURL, username, password string) error {
	t.Helper()
	_, adminUser, adminPassword, err := redshiftEnv()
	if err != nil {
		return err
	}

	connURL = strings.Replace(connURL, fmt.Sprintf("%s:%s", adminUser, adminPassword), fmt.Sprintf("%s:%s", username, password), 1)
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
