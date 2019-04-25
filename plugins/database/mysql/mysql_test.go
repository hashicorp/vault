package mysql

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
	"github.com/hashicorp/vault/sdk/database/helper/credsutil"
	"github.com/ory/dockertest"
)

func prepareMySQLTestContainer(t *testing.T, legacy bool) (cleanup func(), retURL string) {
	if os.Getenv("MYSQL_URL") != "" {
		return func() {}, os.Getenv("MYSQL_URL")
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	imageVersion := "5.7"
	if legacy {
		imageVersion = "5.6"
	}

	resource, err := pool.Run("mysql", imageVersion, []string{"MYSQL_ROOT_PASSWORD=secret"})
	if err != nil {
		t.Fatalf("Could not start local MySQL docker container: %s", err)
	}

	cleanup = func() {
		docker.CleanupResource(t, pool, resource)
	}

	retURL = fmt.Sprintf("root:secret@(localhost:%s)/mysql?parseTime=true", resource.GetPort("3306/tcp"))

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		var err error
		var db *sql.DB
		db, err = sql.Open("mysql", retURL)
		if err != nil {
			return err
		}
		defer db.Close()
		return db.Ping()
	}); err != nil {
		cleanup()
		t.Fatalf("Could not connect to MySQL docker container: %s", err)
	}

	return
}

func TestMySQL_Initialize(t *testing.T) {
	cleanup, connURL := prepareMySQLTestContainer(t, false)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := new(MetadataLen, MetadataLen, UsernameLen)
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

func TestMySQL_CreateUser(t *testing.T) {
	cleanup, connURL := prepareMySQLTestContainer(t, false)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := new(MetadataLen, MetadataLen, UsernameLen)
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test-long-displayname",
		RoleName:    "test-long-rolename",
	}

	// Test with no configured Creation Statement
	_, _, err = db.CreateUser(context.Background(), dbplugin.Statements{}, usernameConfig, time.Now().Add(time.Minute))
	if err == nil {
		t.Fatal("Expected error when no creation statement is provided")
	}

	statements := dbplugin.Statements{
		Creation: []string{testMySQLRoleWildCard},
	}

	username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test a second time to make sure usernames don't collide
	username, password, err = db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test with a manually prepare statement
	statements.Creation = []string{testMySQLRolePreparedStmt}

	username, password, err = db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

}

func TestMySQL_CreateUser_Legacy(t *testing.T) {
	cleanup, connURL := prepareMySQLTestContainer(t, true)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := new(credsutil.NoneLength, LegacyMetadataLen, LegacyUsernameLen)
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test-long-displayname",
		RoleName:    "test-long-rolename",
	}

	// Test with no configured Creation Statement
	_, _, err = db.CreateUser(context.Background(), dbplugin.Statements{}, usernameConfig, time.Now().Add(time.Minute))
	if err == nil {
		t.Fatal("Expected error when no creation statement is provided")
	}

	statements := dbplugin.Statements{
		Creation: []string{testMySQLRoleWildCard},
	}

	username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test a second time to make sure usernames don't collide
	username, password, err = db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
}

func TestMySQL_RotateRootCredentials(t *testing.T) {
	cleanup, connURL := prepareMySQLTestContainer(t, false)
	defer cleanup()

	connURL = strings.Replace(connURL, "root:secret", `{{username}}:{{password}}`, -1)

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
		"username":       "root",
		"password":       "secret",
	}

	db := new(MetadataLen, MetadataLen, UsernameLen)
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !db.Initialized {
		t.Fatal("Database should be initialized")
	}

	newConf, err := db.RotateRootCredentials(context.Background(), nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if newConf["password"] == "secret" {
		t.Fatal("password was not updated")
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestMySQL_RevokeUser(t *testing.T) {
	cleanup, connURL := prepareMySQLTestContainer(t, false)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := new(MetadataLen, MetadataLen, UsernameLen)
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		Creation: []string{testMySQLRoleWildCard},
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err != nil {
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

	statements.Creation = []string{testMySQLRoleWildCard}
	username, password, err = db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test custom revoke statements
	statements.Revocation = []string{testMySQLRevocationSQL}
	err = db.RevokeUser(context.Background(), statements, username)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err == nil {
		t.Fatal("Credentials were not revoked")
	}
}

func TestMySQL_SetCredentials(t *testing.T) {
	cleanup, connURL := prepareMySQLTestContainer(t, false)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := new(MetadataLen, MetadataLen, UsernameLen)
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	password, err := db.GenerateCredentials(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	usernameConfig := dbplugin.StaticUserConfig{
		Username: "test",
		Password: password,
	}

	// Test with no configured Creation Statement
	username, password, err := db.SetCredentials(context.Background(), dbplugin.Statements{}, usernameConfig)
	if err == nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		Creation: []string{testMySQLStaticCreateSQL},
	}
	// User should not exist, make sure we can create
	username, password, err = db.SetCredentials(context.Background(), statements, usernameConfig)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// call SetCredentials again, the user will already exist, password will
	// change. Without rotation statements, this should use the defaults
	newPassword, _ := db.GenerateCredentials(context.Background())
	usernameConfig.Password = newPassword
	username, password, err = db.SetCredentials(context.Background(), statements, usernameConfig)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if password != newPassword {
		t.Fatal("passwords should have changed")
	}

	if err := testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// generate a new password and supply owr own rotation statements
	newPassword2, _ := db.GenerateCredentials(context.Background())
	usernameConfig.Password = newPassword2
	statements.Rotation = []string{testMySQLStaticRotateSQL, testMySQLGrantSQL}
	username, password, err = db.SetCredentials(context.Background(), statements, usernameConfig)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if password != newPassword2 {
		t.Fatal("passwords should have changed")
	}

	if err := testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
}

func testCredsExist(t testing.TB, connURL, username, password string) error {
	// Log in with the new creds
	connURL = strings.Replace(connURL, "root:secret", fmt.Sprintf("%s:%s", username, password), 1)
	db, err := sql.Open("mysql", connURL)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}

const testMySQLRolePreparedStmt = `
CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';
set @grants=CONCAT("GRANT SELECT ON ", "*", ".* TO '{{name}}'@'%'");
PREPARE grantStmt from @grants;
EXECUTE grantStmt;
DEALLOCATE PREPARE grantStmt;
`
const testMySQLRoleWildCard = `
CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';
GRANT SELECT ON *.* TO '{{name}}'@'%';
`
const testMySQLRevocationSQL = `
REVOKE ALL PRIVILEGES, GRANT OPTION FROM '{{name}}'@'%'; 
DROP USER '{{name}}'@'%';
`

const testMySQLStaticCreateSQL = `
CREATE USER '{{username}}'@'%' IDENTIFIED BY '{{password}}';
GRANT SELECT ON *.* TO '{{username}}'@'%';
`
const testMySQLStaticRotateSQL = `
ALTER USER '{{username}}'@'%' IDENTIFIED BY '{{password}}';
`

const testMySQLGrantSQL = `
GRANT SELECT ON *.* TO '{{username}}'@'%';
`
