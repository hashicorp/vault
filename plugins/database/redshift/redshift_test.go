package redshift

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
	dbplugin "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	dbtesting "github.com/hashicorp/vault/sdk/database/dbplugin/v5/testing"
	"github.com/hashicorp/vault/sdk/helper/dbtxn"
	"github.com/stretchr/testify/require"
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
		"max_open_connections": 73,
	}

	db := newRedshift()
	resp := dbtesting.AssertInitialize(t, db, dbplugin.InitializeRequest{
		Config:           connectionDetails,
		VerifyConnection: true,
	})

	if !db.Initialized {
		t.Fatal("Database should be initialized")
	}
	expectedConfig := make(map[string]interface{})
	for k, v := range connectionDetails {
		expectedConfig[k] = v
	}
	if !reflect.DeepEqual(expectedConfig, resp.Config) {
		t.Fatalf("Expected config %+v, but was %v", expectedConfig, resp.Config)
	}
	if db.MaxOpenConnections != 73 {
		t.Fatalf("Expected max_open_connections to be set to 73, but was %d", db.MaxOpenConnections)
	}

	dbtesting.AssertClose(t, db)
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
	dbtesting.AssertInitialize(t, db, dbplugin.InitializeRequest{
		Config:           connectionDetails,
		VerifyConnection: true,
	})

	usernameConfig := dbplugin.UsernameMetadata{
		DisplayName: "test",
		RoleName:    "test",
	}

	const password = "SuperSecurePa55w0rd!"
	for _, commands := range [][]string{{testRedshiftRole}, {testRedshiftReadOnlyRole}} {
		resp := dbtesting.AssertNewUser(t, db, dbplugin.NewUserRequest{
			UsernameConfig: usernameConfig,
			Password:       password,
			Statements: dbplugin.Statements{
				Commands: commands,
			},
			Expiration: time.Now().Add(5 * time.Minute),
		})
		username := resp.Username

		if err = testCredsExist(t, url, username, password); err != nil {
			t.Fatalf("Could not connect with new credentials: %s\n%s:%s", err, username, password)
		}

		usernameRegex := regexp.MustCompile("^v-test-test-[a-zA-Z0-9]{20}-[0-9]{10}$")
		if !usernameRegex.Match([]byte(username)) {
			t.Fatalf("Expected username %q to match regex %q", username, usernameRegex.String())
		}
	}

	dbtesting.AssertClose(t, db)
}

func TestRedshift_NewUser_NoCreationStatement_ShouldError(t *testing.T) {
	if os.Getenv(vaultACC) != "1" {
		t.SkipNow()
	}

	connURL, _, _, _, err := redshiftEnv()
	if err != nil {
		t.Fatal(err)
	}

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := newRedshift()
	dbtesting.AssertInitialize(t, db, dbplugin.InitializeRequest{
		Config:           connectionDetails,
		VerifyConnection: true,
	})

	usernameConfig := dbplugin.UsernameMetadata{
		DisplayName: "test",
		RoleName:    "test",
	}

	const password = "SuperSecurePa55w0rd!"

	// Test with no configured Creation Statement
	_, err = db.NewUser(context.Background(), dbplugin.NewUserRequest{
		UsernameConfig: usernameConfig,
		Password:       password,
		Statements: dbplugin.Statements{
			Commands: []string{}, // Empty commands field here should cause error.
		},
		Expiration: time.Now().Add(5 * time.Minute),
	})
	if err == nil {
		t.Fatal("Expected error when no creation statement is provided")
	}

	dbtesting.AssertClose(t, db)
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
	dbtesting.AssertInitialize(t, db, dbplugin.InitializeRequest{
		Config:           connectionDetails,
		VerifyConnection: true,
	})

	usernameConfig := dbplugin.UsernameMetadata{
		DisplayName: "test",
		RoleName:    "test",
	}

	const password = "SuperSecurePa55w0rd!"
	const initialTTL = 2 * time.Second
	const longTTL = time.Minute
	for _, commands := range [][]string{{}, {defaultRenewSQL}} {
		newResp := dbtesting.AssertNewUser(t, db, dbplugin.NewUserRequest{
			UsernameConfig: usernameConfig,
			Password:       password,
			Statements:     dbplugin.Statements{Commands: []string{testRedshiftRole}},
			Expiration:     time.Now().Add(initialTTL),
		})
		username := newResp.Username

		if err = testCredsExist(t, url, username, password); err != nil {
			t.Fatalf("Could not connect with new credentials: %s", err)
		}

		dbtesting.AssertUpdateUser(t, db, dbplugin.UpdateUserRequest{
			Username: username,
			Expiration: &dbplugin.ChangeExpiration{
				NewExpiration: time.Now().Add(longTTL),
				Statements:    dbplugin.Statements{Commands: commands},
			},
		})

		// Sleep longer than the initial expiration time
		time.Sleep(initialTTL + time.Second)

		if err = testCredsExist(t, url, username, password); err != nil {
			t.Fatalf("Could not connect with new credentials: %s", err)
		}
	}

	dbtesting.AssertClose(t, db)
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
	dbtesting.AssertInitialize(t, db, dbplugin.InitializeRequest{
		Config:           connectionDetails,
		VerifyConnection: true,
	})

	const password1 = "MyTemporaryUserPassword1!"
	const password2 = "MyTemporaryUserPassword2!"

	for _, tc := range []struct {
		password string
		commands []string
	}{
		{password1, []string{}},
		{password2, []string{testRedshiftStaticRoleRotate}},
	} {
		dbtesting.AssertUpdateUser(t, db, dbplugin.UpdateUserRequest{
			Username: dbUser,
			Password: &dbplugin.ChangePassword{
				NewPassword: tc.password,
				Statements:  dbplugin.Statements{Commands: tc.commands},
			},
		})

		if err := testCredsExist(t, url, dbUser, tc.password); err != nil {
			t.Fatalf("Could not connect with new credentials: %s", err)
		}
	}

	dbtesting.AssertClose(t, db)
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
	dbtesting.AssertInitialize(t, db, dbplugin.InitializeRequest{
		Config:           connectionDetails,
		VerifyConnection: true,
	})

	usernameConfig := dbplugin.UsernameMetadata{
		DisplayName: "test",
		RoleName:    "test",
	}

	const password = "SuperSecretPa55word!"
	for _, commands := range [][]string{{}, {defaultRedshiftRevocationSQL}} {
		newResponse := dbtesting.AssertNewUser(t, db, dbplugin.NewUserRequest{
			UsernameConfig: usernameConfig,
			Statements:     dbplugin.Statements{Commands: []string{testRedshiftRole}},
			Password:       password,
			Expiration:     time.Now().Add(2 * time.Second),
		})
		username := newResponse.Username

		if err = testCredsExist(t, url, username, password); err != nil {
			t.Fatalf("Could not connect with new credentials: %s", err)
		}

		// Intentionally _not_ using dbtesting here as the call almost always takes longer than the 2s default timeout
		db.DeleteUser(context.Background(), dbplugin.DeleteUserRequest{
			Username:   username,
			Statements: dbplugin.Statements{Commands: commands},
		})

		if err := testCredsExist(t, url, username, password); err == nil {
			t.Fatal("Credentials were not revoked")
		}
	}

	dbtesting.AssertClose(t, db)
}

func testCredsExist(t testing.TB, url, username, password string) error {
	t.Helper()

	connURL := interpolateConnectionURL(url, username, password)
	db, err := sql.Open("pgx", connURL)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}

func TestRedshift_DefaultUsernameTemplate(t *testing.T) {
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
	dbtesting.AssertInitialize(t, db, dbplugin.InitializeRequest{
		Config:           connectionDetails,
		VerifyConnection: true,
	})

	usernameConfig := dbplugin.UsernameMetadata{
		DisplayName: "test",
		RoleName:    "test",
	}

	const password = "SuperSecurePa55w0rd!"
	for _, commands := range [][]string{{testRedshiftRole}, {testRedshiftReadOnlyRole}} {
		resp := dbtesting.AssertNewUser(t, db, dbplugin.NewUserRequest{
			UsernameConfig: usernameConfig,
			Password:       password,
			Statements: dbplugin.Statements{
				Commands: commands,
			},
			Expiration: time.Now().Add(5 * time.Minute),
		})
		username := resp.Username

		if resp.Username == "" {
			t.Fatalf("Missing username")
		}

		testCredsExist(t, url, username, password)

		require.Regexp(t, `^v-test-test-[a-z0-9]{20}-[0-9]{10}$`, resp.Username)
	}
	dbtesting.AssertClose(t, db)
}

func TestRedshift_CustomUsernameTemplate(t *testing.T) {
	if os.Getenv(vaultACC) != "1" {
		t.SkipNow()
	}

	connURL, url, _, _, err := redshiftEnv()
	if err != nil {
		t.Fatal(err)
	}

	connectionDetails := map[string]interface{}{
		"connection_url":    connURL,
		"username_template": "{{.DisplayName}}-{{random 10}}",
	}

	db := newRedshift()
	dbtesting.AssertInitialize(t, db, dbplugin.InitializeRequest{
		Config:           connectionDetails,
		VerifyConnection: true,
	})

	usernameConfig := dbplugin.UsernameMetadata{
		DisplayName: "test",
		RoleName:    "test",
	}

	const password = "SuperSecurePa55w0rd!"
	for _, commands := range [][]string{{testRedshiftRole}, {testRedshiftReadOnlyRole}} {
		resp := dbtesting.AssertNewUser(t, db, dbplugin.NewUserRequest{
			UsernameConfig: usernameConfig,
			Password:       password,
			Statements: dbplugin.Statements{
				Commands: commands,
			},
			Expiration: time.Now().Add(5 * time.Minute),
		})
		username := resp.Username

		if resp.Username == "" {
			t.Fatalf("Missing username")
		}

		testCredsExist(t, url, username, password)

		require.Regexp(t, `^test-[a-zA-Z0-9]{10}$`, resp.Username)
	}
	dbtesting.AssertClose(t, db)
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

	db, err := sql.Open("pgx", connURL)
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
	if err := dbtxn.ExecuteTxQueryDirect(ctx, tx, m, query); err != nil {
		t.Fatal(err)
	}
	// Commit the transaction
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}
