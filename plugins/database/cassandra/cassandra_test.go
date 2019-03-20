package cassandra

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"fmt"

	"github.com/gocql/gocql"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/ory/dockertest"
)

func prepareCassandraTestContainer(t *testing.T) (func(), string, int) {
	if os.Getenv("CASSANDRA_HOST") != "" {
		return func() {}, os.Getenv("CASSANDRA_HOST"), 0
	}

	pool, err := dockertest.NewPool("")
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
		err := pool.Purge(resource)
		if err != nil {
			t.Fatalf("Failed to cleanup local container: %s", err)
		}
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
	return cleanup, address, port
}

func TestCassandra_Initialize(t *testing.T) {
	if os.Getenv("VAULT_ACC") == "" {
		t.SkipNow()
	}
	cleanup, address, port := prepareCassandraTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"hosts":            address,
		"port":             port,
		"username":         "cassandra",
		"password":         "cassandra",
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

	// test a string protocol
	connectionDetails = map[string]interface{}{
		"hosts":            address,
		"port":             strconv.Itoa(port),
		"username":         "cassandra",
		"password":         "cassandra",
		"protocol_version": "4",
	}

	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestCassandra_CreateUser(t *testing.T) {
	if os.Getenv("VAULT_ACC") == "" {
		t.SkipNow()
	}
	cleanup, address, port := prepareCassandraTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"hosts":            address,
		"port":             port,
		"username":         "cassandra",
		"password":         "cassandra",
		"protocol_version": 4,
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		Creation: []string{testCassandraRole},
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, address, port, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
}

func TestMyCassandra_RenewUser(t *testing.T) {
	if os.Getenv("VAULT_ACC") == "" {
		t.SkipNow()
	}
	cleanup, address, port := prepareCassandraTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"hosts":            address,
		"port":             port,
		"username":         "cassandra",
		"password":         "cassandra",
		"protocol_version": 4,
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		Creation: []string{testCassandraRole},
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, address, port, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	err = db.RenewUser(context.Background(), statements, username, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestCassandra_RevokeUser(t *testing.T) {
	if os.Getenv("VAULT_ACC") == "" {
		t.SkipNow()
	}
	cleanup, address, port := prepareCassandraTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"hosts":            address,
		"port":             port,
		"username":         "cassandra",
		"password":         "cassandra",
		"protocol_version": 4,
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		Creation: []string{testCassandraRole},
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, address, port, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test default revoke statements
	err = db.RevokeUser(context.Background(), statements, username)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, address, port, username, password); err == nil {
		t.Fatal("Credentials were not revoked")
	}
}

func TestCassandra_RotateRootCredentials(t *testing.T) {
	if os.Getenv("VAULT_ACC") == "" {
		t.SkipNow()
	}
	cleanup, address, port := prepareCassandraTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"hosts":            address,
		"port":             port,
		"username":         "cassandra",
		"password":         "cassandra",
		"protocol_version": 4,
	}

	db := new()

	connProducer := db.cassandraConnectionProducer

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
	if newConf["password"] == "cassandra" {
		t.Fatal("password was not updated")
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testCredsExist(t testing.TB, address string, port int, username, password string) error {
	clusterConfig := gocql.NewCluster(address)
	clusterConfig.Authenticator = gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}
	clusterConfig.ProtoVersion = 4
	clusterConfig.Port = port

	session, err := clusterConfig.CreateSession()
	if err != nil {
		return errwrap.Wrapf("error creating session: {{err}}", err)
	}
	defer session.Close()
	return nil
}

const testCassandraRole = `CREATE USER '{{username}}' WITH PASSWORD '{{password}}' NOSUPERUSER;
GRANT ALL PERMISSIONS ON ALL KEYSPACES TO {{username}};`
