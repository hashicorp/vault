package mongodb

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers/mongodb"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"gopkg.in/mgo.v2"
)

const testMongoDBRole = `{ "db": "admin", "roles": [ { "role": "readWrite" } ] }`

const testMongoDBWriteConcern = `{ "wmode": "majority", "wtimeout": 5000 }`

func TestMongoDB_Initialize(t *testing.T) {
	cleanup, connURL := mongodb.PrepareTestContainer(t, "latest")
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
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

func TestMongoDB_CreateUser(t *testing.T) {
	cleanup, connURL := mongodb.PrepareTestContainer(t, "latest")
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		Creation: []string{testMongoDBRole},
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
}

func TestMongoDB_CreateUser_writeConcern(t *testing.T) {
	cleanup, connURL := mongodb.PrepareTestContainer(t, "latest")
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
		"write_concern":  testMongoDBWriteConcern,
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		Creation: []string{testMongoDBRole},
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
}

func TestMongoDB_RevokeUser(t *testing.T) {
	cleanup, connURL := mongodb.PrepareTestContainer(t, "latest")
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		Creation: []string{testMongoDBRole},
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

	// Test default revocation statement
	err = db.RevokeUser(context.Background(), statements, username)
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

func TestMongoDB_SetCredentials(t *testing.T) {
	cleanup, connURL := mongodb.PrepareTestContainer(t, "latest")
	defer cleanup()

	// The docker test method PrepareTestContainer defaults to a database "test"
	// if none is provided
	connURL = connURL + "/test"
	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// create the database user in advance, and test the connection
	dbUser := "testmongouser"
	startingPassword := "password"
	testCreateDBUser(t, connURL, dbUser, startingPassword)
	if err := testCredsExist(t, connURL, dbUser, startingPassword); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	newPassword, err := db.GenerateCredentials(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	usernameConfig := dbplugin.StaticUserConfig{
		Username: dbUser,
		Password: newPassword,
	}

	username, password, err := db.SetCredentials(context.Background(), dbplugin.Statements{}, usernameConfig)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
	// confirm the original creds used to set still work (should be the same)
	if err := testCredsExist(t, connURL, dbUser, newPassword); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	if (dbUser != username) || (newPassword != password) {
		t.Fatalf("username/password mismatch: (%s)/(%s) vs (%s)/(%s)", dbUser, username, newPassword, password)
	}
}

func testCreateDBUser(t testing.TB, connURL, username, password string) {
	dialInfo, err := parseMongoURL(connURL)
	if err != nil {
		t.Fatal(err)
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		t.Fatal(err)
	}
	session.SetSyncTimeout(1 * time.Minute)
	session.SetSocketTimeout(1 * time.Minute)
	mUser := mgo.User{
		Username: username,
		Password: password,
	}

	if err := session.DB(dialInfo.Database).UpsertUser(&mUser); err != nil {
		t.Fatal(err)
	}
}
