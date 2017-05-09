package mongodb

import (
	"fmt"
	"os"
	"testing"
	"time"

	mgo "gopkg.in/mgo.v2"

	"strings"

	"github.com/calvn/dockertest"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/plugins/helper/database/connutil"
)

const testMongoDBRole = `{ "db": "admin", "roles": [ { "role": "readWrite" } ] }`

func prepareMongoDBTestContainer(t *testing.T) (cleanup func(), retURI string) {
	if os.Getenv("MONGODB_URI") != "" {
		return func() {}, os.Getenv("MONGODB_URI")
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

	retURI = fmt.Sprintf("mongodb://127.0.0.1:%s", resource.GetPort("27017/tcp"))

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		var err error
		dialInfo, err := connutil.ParseMongoURI(retURI)
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
	cleanup, connURI := prepareMongoDBTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"uri": connURI,
	}

	dbRaw, err := New()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	db := dbRaw.(*MongoDB)
	connProducer := db.ConnectionProducer.(*connutil.MongoDBConnectionProducer)

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
	cleanup, connURI := prepareMongoDBTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"uri": connURI,
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

	username, password, err := db.CreateUser(statements, "test", time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURI, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
}

func TestMongoDB_RevokeUser(t *testing.T) {
	cleanup, connURI := prepareMongoDBTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"uri": connURI,
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

	username, password, err := db.CreateUser(statements, "test", time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURI, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test default revocation statememt
	err = db.RevokeUser(statements, username)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, connURI, username, password); err == nil {
		t.Fatal("Credentials were not revoked")
	}
}

func testCredsExist(t testing.TB, connURI, username, password string) error {
	connURI = strings.Replace(connURI, "mongodb://", fmt.Sprintf("mongodb://%s:%s@", username, password), 1)
	dialInfo, err := connutil.ParseMongoURI(connURI)
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
