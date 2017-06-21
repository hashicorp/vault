package mongodb

import (
	"fmt"
	"os"
	"testing"
	"time"

	mgo "gopkg.in/mgo.v2"

	"strings"

	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	dockertest "gopkg.in/ory-am/dockertest.v3"
)

const testMongoDBRole = `{ "db": "admin", "roles": [ { "role": "readWrite" } ] }`

func prepareMongoDBTestContainer(t *testing.T) (cleanup func(), retURL string) {
	if os.Getenv("MONGODB_URL") != "" {
		return func() {}, os.Getenv("MONGODB_URL")
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	resource, err := pool.Run("mongo", "latest", []string{})
	if err != nil {
		t.Fatalf("Could not start local mongo docker container: %s", err)
	}

	cleanup = func() {
		err := pool.Purge(resource)
		if err != nil {
			t.Fatalf("Failed to cleanup local container: %s", err)
		}
	}

	retURL = fmt.Sprintf("mongodb://localhost:%s", resource.GetPort("27017/tcp"))

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		var err error
		dialInfo, err := parseMongoURL(retURL)
		if err != nil {
			return err
		}

		session, err := mgo.DialWithInfo(dialInfo)
		if err != nil {
			return err
		}
		session.SetSyncTimeout(1 * time.Minute)
		session.SetSocketTimeout(1 * time.Minute)
		return session.Ping()
	}); err != nil {
		t.Fatalf("Could not connect to mongo docker container: %s", err)
	}

	return
}

func TestMongoDB_Initialize(t *testing.T) {
	cleanup, connURL := prepareMongoDBTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	dbRaw, err := New()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	db := dbRaw.(*MongoDB)
	connProducer := db.ConnectionProducer.(*mongoDBConnectionProducer)

	err = db.Initialize(connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !connProducer.Initialized {
		t.Fatal("Database should be initialized")
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestMongoDB_CreateUser(t *testing.T) {
	cleanup, connURL := prepareMongoDBTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	dbRaw, err := New()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	db := dbRaw.(*MongoDB)
	err = db.Initialize(connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		CreationStatements: testMongoDBRole,
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	username, password, err := db.CreateUser(statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
}

func TestMongoDB_RevokeUser(t *testing.T) {
	cleanup, connURL := prepareMongoDBTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	dbRaw, err := New()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	db := dbRaw.(*MongoDB)
	err = db.Initialize(connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		CreationStatements: testMongoDBRole,
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	username, password, err := db.CreateUser(statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test default revocation statememt
	err = db.RevokeUser(statements, username)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, connURL, username, password); err == nil {
		t.Fatal("Credentials were not revoked")
	}
}

func testCredsExist(t testing.TB, connURL, username, password string) error {
	connURL = strings.Replace(connURL, "mongodb://", fmt.Sprintf("mongodb://%s:%s@", username, password), 1)
	dialInfo, err := parseMongoURL(connURL)
	if err != nil {
		return err
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return err
	}
	session.SetSyncTimeout(1 * time.Minute)
	session.SetSocketTimeout(1 * time.Minute)
	return session.Ping()
}
