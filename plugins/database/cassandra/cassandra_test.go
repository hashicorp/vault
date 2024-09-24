// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cassandra

import (
	"context"
	"reflect"
	"testing"
	"time"

	backoff "github.com/cenkalti/backoff/v3"
	"github.com/gocql/gocql"
	"github.com/hashicorp/vault/helper/testhelpers/cassandra"
	dbplugin "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	dbtesting "github.com/hashicorp/vault/sdk/database/dbplugin/v5/testing"
	"github.com/stretchr/testify/require"
)

func getCassandra(t *testing.T, protocolVersion interface{}) (*Cassandra, func()) {
	host, cleanup := cassandra.PrepareTestContainer(t,
		cassandra.Version("3.11"),
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

func TestInitialize(t *testing.T) {
	t.Run("integer protocol version", func(t *testing.T) {
		// getCassandra performs an Initialize call
		db, cleanup := getCassandra(t, 4)
		t.Cleanup(cleanup)

		err := db.Close()
		if err != nil {
			t.Fatalf("err: %s", err)
		}
	})

	t.Run("string protocol version", func(t *testing.T) {
		// getCassandra performs an Initialize call
		db, cleanup := getCassandra(t, "4")
		t.Cleanup(cleanup)

		err := db.Close()
		if err != nil {
			t.Fatalf("err: %s", err)
		}
	})
}

func TestCreateUser(t *testing.T) {
	type testCase struct {
		// Config will have the hosts & port added to it during the test
		config                map[string]interface{}
		newUserReq            dbplugin.NewUserRequest
		expectErr             bool
		expectedUsernameRegex string
		assertCreds           func(t testing.TB, address string, port int, username, password string, sslOpts *gocql.SslOptions, timeout time.Duration)
	}

	tests := map[string]testCase{
		"default username_template": {
			config: map[string]interface{}{
				"username":         "cassandra",
				"password":         "cassandra",
				"protocol_version": "4",
				"connect_timeout":  "20s",
			},
			newUserReq: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "token",
					RoleName:    "mylongrolenamewithmanycharacters",
				},
				Statements: dbplugin.Statements{
					Commands: []string{createUserStatements},
				},
				Password:   "bfn985wjAHIh6t",
				Expiration: time.Now().Add(1 * time.Minute),
			},
			expectErr:             false,
			expectedUsernameRegex: `^v_token_mylongrolenamew_[a-z0-9]{20}_[0-9]{10}$`,
			assertCreds:           assertCreds,
		},
		"custom username_template": {
			config: map[string]interface{}{
				"username":          "cassandra",
				"password":          "cassandra",
				"protocol_version":  "4",
				"connect_timeout":   "20s",
				"username_template": `foo_{{random 20}}_{{.RoleName | replace "e" "3"}}_{{unix_time}}`,
			},
			newUserReq: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "token",
					RoleName:    "mylongrolenamewithmanycharacters",
				},
				Statements: dbplugin.Statements{
					Commands: []string{createUserStatements},
				},
				Password:   "bfn985wjAHIh6t",
				Expiration: time.Now().Add(1 * time.Minute),
			},
			expectErr:             false,
			expectedUsernameRegex: `^foo_[a-zA-Z0-9]{20}_mylongrol3nam3withmanycharact3rs_[0-9]{10}$`,
			assertCreds:           assertCreds,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			host, cleanup := cassandra.PrepareTestContainer(t,
				cassandra.Version("3.11"),
				cassandra.CopyFromTo(insecureFileMounts),
			)
			defer cleanup()

			db := new()

			config := test.config
			config["hosts"] = host.ConnectionURL()
			config["port"] = host.Port

			initReq := dbplugin.InitializeRequest{
				Config:           config,
				VerifyConnection: true,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			dbtesting.AssertInitialize(t, db, initReq)

			require.True(t, db.Initialized, "Database is not initialized")

			ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			newUserResp, err := db.NewUser(ctx, test.newUserReq)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}
			require.Regexp(t, test.expectedUsernameRegex, newUserResp.Username)
			test.assertCreds(t, db.Hosts, db.Port, newUserResp.Username, test.newUserReq.Password, nil, 5*time.Second)
		})
	}
}

func TestUpdateUserPassword(t *testing.T) {
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

	assertCreds(t, db.Hosts, db.Port, createResp.Username, password, nil, 5*time.Second)

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

	assertCreds(t, db.Hosts, db.Port, createResp.Username, newPassword, nil, 5*time.Second)
}

func TestDeleteUser(t *testing.T) {
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

	assertCreds(t, db.Hosts, db.Port, createResp.Username, password, nil, 5*time.Second)

	deleteReq := dbplugin.DeleteUserRequest{
		Username: createResp.Username,
	}

	dbtesting.AssertDeleteUser(t, db, deleteReq)

	assertNoCreds(t, db.Hosts, db.Port, createResp.Username, password, nil, 5*time.Second)
}

func assertCreds(t testing.TB, address string, port int, username, password string, sslOpts *gocql.SslOptions, timeout time.Duration) {
	t.Helper()
	op := func() error {
		return connect(t, address, port, username, password, sslOpts)
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

func connect(t testing.TB, address string, port int, username, password string, sslOpts *gocql.SslOptions) error {
	t.Helper()
	clusterConfig := gocql.NewCluster(address)
	clusterConfig.Authenticator = gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}
	clusterConfig.ProtoVersion = 4
	clusterConfig.Port = port
	clusterConfig.SslOpts = sslOpts

	session, err := clusterConfig.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()
	return nil
}

func assertNoCreds(t testing.TB, address string, port int, username, password string, sslOpts *gocql.SslOptions, timeout time.Duration) {
	t.Helper()

	op := func() error {
		// "Invert" the error so the backoff logic sees a failure to connect as a success
		err := connect(t, address, port, username, password, sslOpts)
		if err != nil {
			return nil
		}
		return nil
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
GRANT ALL PERMISSIONS ON ALL KEYSPACES TO '{{username}}';`
