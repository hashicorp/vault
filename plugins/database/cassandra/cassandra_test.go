package cassandra

import (
	"reflect"
	"regexp"
	"testing"
	"time"

	backoff "github.com/cenkalti/backoff/v3"
	"github.com/gocql/gocql"
	"github.com/hashicorp/vault/helper/testhelpers/cassandra"
	dbplugin "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	dbtesting "github.com/hashicorp/vault/sdk/database/dbplugin/v5/testing"
)

func getCassandra(t *testing.T, protocolVersion interface{}) (*Cassandra, func()) {
	host, cleanup := cassandra.PrepareTestContainer(t,
		cassandra.Version("latest"),
		cassandra.CopyFromTo(insecureFileMounts),
	)

	db := new()
	initReq := dbplugin.InitializeRequest{
		Config: map[string]interface{}{
			"hosts":            host.ConnectionURL(),
			"port":             host.Port,
			"username":         "cassandra",
			"password":         "cassandra",
			"protocol_version": protocolVersion,
			"connect_timeout":  "20s",
		},
		VerifyConnection: true,
	}

	expectedConfig := map[string]interface{}{
		"hosts":            host.ConnectionURL(),
		"port":             host.Port,
		"username":         "cassandra",
		"password":         "cassandra",
		"protocol_version": protocolVersion,
		"connect_timeout":  "20s",
	}

	initResp := dbtesting.AssertInitialize(t, db, initReq)
	if !reflect.DeepEqual(initResp.Config, expectedConfig) {
		t.Fatalf("Initialize response config actual: %#v\nExpected: %#v", initResp.Config, expectedConfig)
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
	createReq := dbplugin.NewUserRequest{
		UsernameConfig: dbplugin.UsernameMetadata{
			DisplayName: "test",
			RoleName:    "test",
		},
		Statements: dbplugin.Statements{
			Commands: []string{createUserStatements},
		},
		Password:   password,
		Expiration: time.Now().Add(1 * time.Minute),
	}

	createResp := dbtesting.AssertNewUser(t, db, createReq)

	expectedRegex := "^v_test_test_[a-zA-Z0-9]{20}_[0-9]{10}$"
	re := regexp.MustCompile(expectedRegex)
	if !re.MatchString(createResp.Username) {
		t.Fatalf("Generated username %q did not match regexp %q", createResp.Username, expectedRegex)
	}

	assertCreds(t, db.Hosts, db.Port, createResp.Username, password, 5*time.Second)
}

func TestMyCassandra_UpdateUserPassword(t *testing.T) {
	db, cleanup := getCassandra(t, 4)
	defer cleanup()

	password := "myreallysecurepassword"
	createReq := dbplugin.NewUserRequest{
		UsernameConfig: dbplugin.UsernameMetadata{
			DisplayName: "test",
			RoleName:    "test",
		},
		Statements: dbplugin.Statements{
			Commands: []string{createUserStatements},
		},
		Password:   password,
		Expiration: time.Now().Add(1 * time.Minute),
	}

	createResp := dbtesting.AssertNewUser(t, db, createReq)

	assertCreds(t, db.Hosts, db.Port, createResp.Username, password, 5*time.Second)

	newPassword := "somenewpassword"
	updateReq := dbplugin.UpdateUserRequest{
		Username: createResp.Username,
		Password: &dbplugin.ChangePassword{
			NewPassword: newPassword,
			Statements:  dbplugin.Statements{},
		},
		Expiration: nil,
	}

	dbtesting.AssertUpdateUser(t, db, updateReq)

	assertCreds(t, db.Hosts, db.Port, createResp.Username, newPassword, 5*time.Second)
}

func TestCassandra_DeleteUser(t *testing.T) {
	db, cleanup := getCassandra(t, 4)
	defer cleanup()

	password := "myreallysecurepassword"
	createReq := dbplugin.NewUserRequest{
		UsernameConfig: dbplugin.UsernameMetadata{
			DisplayName: "test",
			RoleName:    "test",
		},
		Statements: dbplugin.Statements{
			Commands: []string{createUserStatements},
		},
		Password:   password,
		Expiration: time.Now().Add(1 * time.Minute),
	}

	createResp := dbtesting.AssertNewUser(t, db, createReq)

	assertCreds(t, db.Hosts, db.Port, createResp.Username, password, 5*time.Second)

	deleteReq := dbplugin.DeleteUserRequest{
		Username: createResp.Username,
	}

	dbtesting.AssertDeleteUser(t, db, deleteReq)

	assertNoCreds(t, db.Hosts, db.Port, createResp.Username, password, 5*time.Second)
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

func assertNoCreds(t testing.TB, address string, port int, username, password string, timeout time.Duration) {
	t.Helper()

	op := func() error {
		// "Invert" the error so the backoff logic sees a failure to connect as a success
		err := connect(t, address, port, username, password)
		if err != nil {
			return nil
		}
		return err
	}
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = timeout
	bo.InitialInterval = 500 * time.Millisecond
	bo.MaxInterval = bo.InitialInterval
	bo.RandomizationFactor = 0.0

	err := backoff.Retry(op, bo)
	if err != nil {
		t.Fatalf("successfully connected after %s when it shouldn't", timeout)
	}
}

const createUserStatements = `CREATE USER '{{username}}' WITH PASSWORD '{{password}}' NOSUPERUSER;
GRANT ALL PERMISSIONS ON ALL KEYSPACES TO {{username}};`
