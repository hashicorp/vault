package couchbase

import (
	"context"
//	"errors"
//	"sync"
	"testing"
	"time"
	//	"fmt"
	"log"
	"os"

	"github.com/hashicorp/vault/sdk/database/dbplugin"
)

func prepareCouchbaseTestContainer(t *testing.T) (func(), string, int) {
	if os.Getenv("COUCHBASE_HOST") != "" {
		return func() {}, os.Getenv("COUCHBASE_HOST"), 0
	}

/*	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	cwd, _ := os.Getwd()
	cassandraMountPath := fmt.Sprintf("%s/test-fixtures/:/etc/cassandra/", cwd)

	ro := &dockertest.RunOptions{
		Repository: "cassandra",
		Tag:        "latest",
		Env:        []string{"CASSANDRA_BROADCAST_ADDRESS=127.0.0.1"},
		Mounts:     []string{cassandraMountPath},
	}
	resource, err := pool.RunWithOptions(ro)
	if err != nil {
		t.Fatalf("Could not start local cassandra docker container: %s", err)
	}

	cleanup := func() {
		docker.CleanupResource(t, pool, resource)
	}

	port, _ := strconv.Atoi(resource.GetPort("9042/tcp"))
	address := fmt.Sprintf("127.0.0.1:%d", port)

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		clusterConfig := gocql.NewCluster(address)
		clusterConfig.Authenticator = gocql.PasswordAuthenticator{
			Username: "cassandra",
			Password: "cassandra",
		}
		clusterConfig.ProtoVersion = 4
		clusterConfig.Port = port

		session, err := clusterConfig.CreateSession()
		if err != nil {
			return errwrap.Wrapf("error creating session: {{err}}", err)
		}
		defer session.Close()
		return nil
	}); err != nil {
		cleanup()
		t.Fatalf("Could not connect to cassandra docker container: %s", err)
	}
	return cleanup, address, port */
	return func() {}, "localhost", 0
}

func TestCouchbaseDB_Initialize(t *testing.T) {
	log.Printf("Testing Init()")
	cleanup, address, port := prepareCouchbaseTestContainer(t)
	defer cleanup()
	address = "couchbases://localhost"
	connectionDetails := map[string]interface{}{
		"hosts":            address,
		"port":             port,
		"username":         "Administrator",
		"password":         "Admin123",
		"tls":              true,
		"insecure_tls":     true,
		"base64pem":        Base64pemCA,
		"protocol_version": 4,
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
}

func TestCouchbaseDB_CreateUser(t *testing.T) {
	if os.Getenv("VAULT_ACC") == "" {
		t.SkipNow()
	}
	log.Printf("Testing CreateUser()")
	
	cleanup, address, port := prepareCouchbaseTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"hosts":            address,
		"port":             port,
		"username":         "Administrator",
		"password":         "Admin123",
		"protocol_version": 4,
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !db.Initialized {
		t.Fatal("Database should be initialized")
	}

	statements := dbplugin.Statements{
		Creation: []string{testCouchbaseRole},
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	db.Close()
	
	if err := testCredsExist(t, address, port, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	err = testRevokeUser(t, address, port, username)
	if err != nil {
		t.Fatalf("Could not revoke user: %s", username)
	}
}

func testCredsExist(t *testing.T, address string, port int, username string, password string) (error) {
	if os.Getenv("VAULT_ACC") == "" {
		t.SkipNow()
	}
	log.Printf("Testing testCredsExist()")
	cleanup, address, port := prepareCouchbaseTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"hosts":            address,
		"port":             port,
		"username":         username,
		"password":         password,
		"protocol_version": 4,
	}
	time.Sleep(1 * time.Second) // a brief pause to let couchbase finish creating the account
	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !db.Initialized {
		t.Fatal("Database should be initialized")
	}

	return nil
}

func testRevokeUser(t *testing.T, address string, port int, username string) (error) {
	if os.Getenv("VAULT_ACC") == "" {
		t.SkipNow()
	}
	log.Printf("Testing RevokeUser()")
	cleanup, address, port := prepareCouchbaseTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"hosts":            address,
		"port":             port,
		"username":         "Administrator",
		"password":         "Admin123",
		"protocol_version": 4,
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !db.Initialized {
		t.Fatal("Database should be initialized")
	}

	statements := dbplugin.Statements{
		Creation: []string{testCouchbaseRole},
		Revocation: []string{"foo"},
	}

	err = db.RevokeUser(context.Background(), statements, username)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	return nil
}

func TestCouchbaseDB_CreateUser_plusRole(t *testing.T) {
	if os.Getenv("VAULT_ACC") == "" {
		t.SkipNow()
	}
	log.Printf("Testing CreateUser_plusRole()")
	cleanup, address, port := prepareCouchbaseTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"hosts":            address,
		"port":             port,
		"username":         "Administrator",
		"password":         "Admin123",
		"protocol_version": 4,
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !db.Initialized {
		t.Fatal("Database should be initialized")
	}

	statements := dbplugin.Statements{
		Creation: []string{testCouchbaseRole},
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	db.Close()
	
	if err := testCredsExist(t, address, port, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	err = testRevokeUser(t, address, port, username)
	if err != nil {
		t.Fatalf("Could not revoke user: %s", username)
	}
}

func TestCouchbaseDB_RotateRootCredentials(t *testing.T) {
	if os.Getenv("VAULT_ACC") == "" {
		t.SkipNow()
	}
	log.Printf("Testing RotateRootCredentials()")
	cleanup, address, port := prepareCouchbaseTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"hosts":            address,
		"port":             port,
		"username":         "rotate-root",
		"password":         "rotate-rootpassword",
		"protocol_version": 4,
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !db.Initialized {
		t.Fatal("Database should be initialized")
	}

	statements := []string{""}

	/*usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}*/

	password, err := db.RotateRootCredentials(context.Background(), statements)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, address, port, db.Username, password["password"].(string)); err != nil {
		t.Fatalf("Could not connect with new RotatedRootcredentials: %s", err)
	}
	// Set password back
	testCouchbaseDB_SetCredentials(t, "rotate-root", "rotate-rootpassword")

	db.Close()
	
}

func testCouchbaseDB_SetCredentials(t *testing.T, username, password string) {
	
	cleanup, address, port := prepareCouchbaseTestContainer(t)
	defer cleanup()
	log.Printf("Testing SetCredentials()")

	connectionDetails := map[string]interface{}{
		"hosts":            address,
		"port":             port,
		"username":         "Administrator",
		"password":         "Admin123",
		"protocol_version": 4,
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !db.Initialized {
		t.Fatal("Database should be initialized")
	}

	statements := dbplugin.Statements{}

	// test that SetCredentials fails if the user does not exist...
	
	staticUser := dbplugin.StaticUserConfig{
		Username: "userThatDoesNotExist",
		Password: password,
	}

	_, _, err = db.SetCredentials(context.Background(), statements, staticUser)
	if err == nil {
		t.Fatalf("err: did not error on setting password for userThatDoesNotExist")
	}

	staticUser = dbplugin.StaticUserConfig{
		Username: username,
		Password: password,
	}

	username, password, err = db.SetCredentials(context.Background(), statements, staticUser)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	db.Close()
	
	if err := testCredsExist(t, address, port, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
}

func TestCouchbaseDB_SetCredentials(t *testing.T) {
	if os.Getenv("VAULT_ACC") == "" {
		t.SkipNow()
	}
	
	testCouchbaseDB_SetCredentials(t, "vault-edu", "password")
}

const testCouchbaseRole = `[{"name":"ro_admin"},{"name":"bucket_admin","bucket":"foo"}]`
const Base64pemCA = `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURBakNDQWVxZ0F3SUJBZ0lJRmhBcGRmdG5oeHd3RFFZSktvWklodmNOQVFFTEJRQXdKREVpTUNBR0ExVUUKQXhNWlEyOTFZMmhpWVhObElGTmxjblpsY2lCaE9UbG1OV05oWXpBZUZ3MHhNekF4TURFd01EQXdNREJhRncwMApPVEV5TXpFeU16VTVOVGxhTUNReElqQWdCZ05WQkFNVEdVTnZkV05vWW1GelpTQlRaWEoyWlhJZ1lUazVaalZqCllXTXdnZ0VpTUEwR0NTcUdTSWIzRFFFQkFRVUFBNElCRHdBd2dnRUtBb0lCQVFEWGV1a21CZWZtTEs1TGpXOEsKOW9rSUs4d1FMTjZJVnlLQ1NRelFJZXNURDRuZzZ0Z1F5bHJZL1Q1RlRHQURpUE9leDZpbEp4dDRnNzJteHVtcworZHZvaE1KaVpWOFFGaTNOeTVYdFc3S05mNUl1Nkk0djQvZmViSkYrdGNjTGgrUUNtaEtPR1F2VUdleDFlT0J5Ck8xVWxRWlVDbTFsZVNKVjRzUWhyWXZPR296THlHMkpUVjRESFNqQW1RbkxMQTNqTExWbHI4V1hOcEdEL1NsaWMKYWFPc0dvaEtudGdwU1AvSTdxMU5ESXhNaEtpY1dmUS9sN3dxQVgwMU1jWXFaNjErcDZuMkF2ZmY5NjJrUjF1aApWQWcwdDVZU3c0WG01RXd2L1hqS2ErZW1RVXg5TVB3TUloYldHVUtlZUJwbHFEbk5VemFDUWVQeTN0dFI2d09pCkZIdURBZ01CQUFHak9EQTJNQTRHQTFVZER3RUIvd1FFQXdJQ3BEQVRCZ05WSFNVRUREQUtCZ2dyQmdFRkJRY0QKQVRBUEJnTlZIUk1CQWY4RUJUQURBUUgvTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCQVFBNktJVW9vWjY3ayt1WApYWDlFY2lUOEQvejhUcDF4WXZ1aVliMlAwbGpZRjBPWU5jNC9sdkk5MUNGek5iMFBhQy9zbFhOYVcya3NIZnlaClJLN08wdk9tSEpVMW8wK0xUdlhJVjZsRld1N3k0aXB4Zy8zY0MwRFl5S0dad0xMbzZ5OUpIMmYvczc5SVFVZXoKZWZiUlhxZ05Ta2E4VUp0cWg5VEhBT0lYY09TUnFId2VvRG56a1NDaC9ZNnJDMWc1RlZWWjhuNDlJQUxMMG14ZAoya0VYYU9IUlhtSloyME9HWDgwZGdWT2Y2Z3lTRHN6ZWU1T3J3VDdESWNQMzdsQkQ0cytHZG9DM25HNktURDRwClVGc3FaSXBSM2x2cmF2V1Fxb1Z2UzJyNHpoN2IvZVZrNmlIOVcycWdxaTQrTzNyK25WcWpQdW1pQTJiNXNUNkkKbDQ5d1BNWTAKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=`
