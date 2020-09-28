package cassandra

import (
	"context"
	"strings"
	"testing"
	"time"

	backoff "github.com/cenkalti/backoff/v3"
	"github.com/gocql/gocql"
	"github.com/hashicorp/vault/helper/testhelpers/cassandra"
	"github.com/hashicorp/vault/sdk/database/newdbplugin"
)

func getCassandra(t *testing.T, protocolVersion interface{}) (*Cassandra, func()) {
	cleanup, connURL := cassandra.PrepareTestContainer(t, "latest")
	pieces := strings.Split(connURL, ":")

	db := new()
	req := newdbplugin.InitializeRequest{
		Config: map[string]interface{}{
			"hosts":            connURL,
			"port":             pieces[1],
			"username":         "cassandra",
			"password":         "cassandra",
			"protocol_version": protocolVersion,
			"connect_timeout":  "20s",
		},
		VerifyConnection: true,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := db.Initialize(ctx, req)
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

	password := "myreallysecurepassword"
	req := newdbplugin.NewUserRequest{
		UsernameConfig: newdbplugin.UsernameMetadata{
			DisplayName: "test",
			RoleName:    "test",
		},
		Statements: newdbplugin.Statements{
			Commands: []string{createUserStatements},
		},
		Password:   password,
		Expiration: time.Now().Add(1 * time.Minute),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := db.NewUser(ctx, req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	assertCreds(t, db.Hosts, db.Port, resp.Username, password, 5*time.Second)
}

func TestMyCassandra_UpdateUserPassword(t *testing.T) {
	db, cleanup := getCassandra(t, 4)
	defer cleanup()

	password := "myreallysecurepassword"
	newUserReq := newdbplugin.NewUserRequest{
		UsernameConfig: newdbplugin.UsernameMetadata{
			DisplayName: "test",
			RoleName:    "test",
		},
		Statements: newdbplugin.Statements{
			Commands: []string{createUserStatements},
		},
		Password:   password,
		Expiration: time.Now().Add(1 * time.Minute),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	createResp, err := db.NewUser(ctx, newUserReq)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	assertCreds(t, db.Hosts, db.Port, createResp.Username, password, 5*time.Second)

	newPassword := "somenewpassword"
	updateReq := newdbplugin.UpdateUserRequest{
		Username: createResp.Username,
		Password: &newdbplugin.ChangePassword{
			NewPassword: newPassword,
			Statements:  newdbplugin.Statements{},
		},
		Expiration: nil,
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = db.UpdateUser(ctx, updateReq)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	assertCreds(t, db.Hosts, db.Port, createResp.Username, newPassword, 5*time.Second)
}

func TestCassandra_DeleteUser(t *testing.T) {
	db, cleanup := getCassandra(t, 4)
	defer cleanup()

	password := "myreallysecurepassword"
	createReq := newdbplugin.NewUserRequest{
		UsernameConfig: newdbplugin.UsernameMetadata{
			DisplayName: "test",
			RoleName:    "test",
		},
		Statements: newdbplugin.Statements{
			Commands: []string{createUserStatements},
		},
		Password:   password,
		Expiration: time.Now().Add(1 * time.Minute),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	createResp, err := db.NewUser(ctx, createReq)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	assertCreds(t, db.Hosts, db.Port, createResp.Username, password, 5*time.Second)

	deleteReq := newdbplugin.DeleteUserRequest{
		Username: createResp.Username,
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = db.DeleteUser(context.Background(), deleteReq)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	assertNoCreds(t, db.Hosts, db.Port, createResp.Username, password)
}

func assertCreds(t testing.TB, address string, port int, username, password string, timeout time.Duration) {
	t.Helper()
	op := func() error {
		return connect(t, address, port, username, password)
	}
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = timeout
	bo.InitialInterval = 500 * time.Millisecond
	bo.MaxInterval = bo.InitialInterval
	bo.RandomizationFactor = 0.0

	err := backoff.Retry(op, bo)
	if err != nil {
		t.Fatalf("failed to connect after %s: %s", timeout, err)
	}
}

func connect(t testing.TB, address string, port int, username, password string) error {
	t.Helper()
	clusterConfig := gocql.NewCluster(address)
	clusterConfig.Authenticator = gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}
	clusterConfig.ProtoVersion = 4
	clusterConfig.Port = port

	session, err := clusterConfig.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()
	return nil
}

func assertNoCreds(t testing.TB, address string, port int, username, password string) {
	t.Helper()
	clusterConfig := gocql.NewCluster(address)
	clusterConfig.Authenticator = gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}
	clusterConfig.ProtoVersion = 4
	clusterConfig.Port = port

	session, err := clusterConfig.CreateSession()
	if err != nil {
		return // Happy path
	}
	defer session.Close()
	t.Fatalf("able to make connection when credentials should not exist")
}

const createUserStatements = `CREATE USER '{{username}}' WITH PASSWORD '{{password}}' NOSUPERUSER;
GRANT ALL PERMISSIONS ON ALL KEYSPACES TO {{username}};`
