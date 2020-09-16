package cassandra

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/testhelpers/cassandra"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
)

func getCassandra(t *testing.T, protocolVersion interface{}) (*Cassandra, func()) {
	cleanup, connURL := cassandra.PrepareTestContainer(t, "latest")
	pieces := strings.Split(connURL, ":")

	connectionDetails := map[string]interface{}{
		"hosts":            connURL,
		"port":             pieces[1],
		"username":         "cassandra",
		"password":         "cassandra",
		"protocol_version": protocolVersion,
		"connect_timeout":  "20s",
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !db.Initialized {
		t.Fatal("Database should be initialized")
	}
	return db, cleanup
}

func TestCassandra_Initialize(t *testing.T) {
	db, cleanup := getCassandra(t, 4)
	defer cleanup()

	err := db.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	db, cleanup = getCassandra(t, "4")
	defer cleanup()
}

func TestCassandra_CreateUser(t *testing.T) {
	db, cleanup := getCassandra(t, 4)
	defer cleanup()

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

	if err := testCredsExist(db.Hosts, db.Port, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
}

func TestMyCassandra_RenewUser(t *testing.T) {
	db, cleanup := getCassandra(t, 4)
	defer cleanup()

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

	if err := testCredsExist(db.Hosts, db.Port, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	err = db.RenewUser(context.Background(), statements, username, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestCassandra_RevokeUser(t *testing.T) {
	db, cleanup := getCassandra(t, 4)
	defer cleanup()

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

	if err = testCredsExist(db.Hosts, db.Port, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test default revoke statements
	err = db.RevokeUser(context.Background(), statements, username)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(db.Hosts, db.Port, username, password); err == nil {
		t.Fatal("Credentials were not revoked")
	}
}

func TestCassandra_RotateRootCredentials(t *testing.T) {
	db, cleanup := getCassandra(t, 4)
	defer cleanup()

	if !db.cassandraConnectionProducer.Initialized {
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

func testCredsExist(address string, port int, username, password string) error {
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
