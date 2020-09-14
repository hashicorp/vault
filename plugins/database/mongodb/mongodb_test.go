package mongodb

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers/certhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/mongodb"
	"github.com/hashicorp/vault/sdk/database/newdbplugin"
	dbtesting "github.com/hashicorp/vault/sdk/database/newdbplugin/testing"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const mongoAdminRole = `{ "db": "admin", "roles": [ { "role": "readWrite" } ] }`

func TestMongoDB_Initialize(t *testing.T) {
	cleanup, connURL := mongodb.PrepareTestContainer(t, "latest")
	defer cleanup()

	db := new()
	defer dbtesting.AssertClose(t, db)

	config := map[string]interface{}{
		"connection_url": connURL,
	}

	// Make a copy since the original map could be modified by the Initialize call
	expectedConfig := copyConfig(config)

	req := newdbplugin.InitializeRequest{
		Config:           config,
		VerifyConnection: true,
	}

	resp := dbtesting.AssertInitialize(t, db, req)

	if !reflect.DeepEqual(resp.Config, expectedConfig) {
		t.Fatalf("Actual config: %#v\nExpected config: %#v", resp.Config, expectedConfig)
	}

	if !db.Initialized {
		t.Fatal("Database should be initialized")
	}
}

func TestMongoDB_CreateUser(t *testing.T) {
	cleanup, connURL := mongodb.PrepareTestContainer(t, "latest")
	defer cleanup()

	db := new()
	defer dbtesting.AssertClose(t, db)

	initReq := newdbplugin.InitializeRequest{
		Config: map[string]interface{}{
			"connection_url": connURL,
		},
		VerifyConnection: true,
	}
	dbtesting.AssertInitialize(t, db, initReq)

	password := "myreallysecurepassword"
	createReq := newdbplugin.NewUserRequest{
		UsernameConfig: newdbplugin.UsernameMetadata{
			DisplayName: "test",
			RoleName:    "test",
		},
		Statements: newdbplugin.Statements{
			Commands: []string{mongoAdminRole},
		},
		Password:   password,
		Expiration: time.Now().Add(time.Minute),
	}
	createResp := dbtesting.AssertNewUser(t, db, createReq)

	assertCredsExist(t, createResp.Username, password, connURL)
}

func TestMongoDB_CreateUser_writeConcern(t *testing.T) {
	cleanup, connURL := mongodb.PrepareTestContainer(t, "latest")
	defer cleanup()

	initReq := newdbplugin.InitializeRequest{
		Config: map[string]interface{}{
			"connection_url": connURL,
			"write_concern":  `{ "wmode": "majority", "wtimeout": 5000 }`,
		},
		VerifyConnection: true,
	}

	db := new()
	defer dbtesting.AssertClose(t, db)

	dbtesting.AssertInitialize(t, db, initReq)

	password := "myreallysecurepassword"
	createReq := newdbplugin.NewUserRequest{
		UsernameConfig: newdbplugin.UsernameMetadata{
			DisplayName: "test",
			RoleName:    "test",
		},
		Statements: newdbplugin.Statements{
			Commands: []string{mongoAdminRole},
		},
		Password:   password,
		Expiration: time.Now().Add(time.Minute),
	}
	createResp := dbtesting.AssertNewUser(t, db, createReq)

	assertCredsExist(t, createResp.Username, password, connURL)
}

func TestMongoDB_DeleteUser(t *testing.T) {
	cleanup, connURL := mongodb.PrepareTestContainer(t, "latest")
	defer cleanup()

	db := new()
	defer dbtesting.AssertClose(t, db)

	initReq := newdbplugin.InitializeRequest{
		Config: map[string]interface{}{
			"connection_url": connURL,
		},
		VerifyConnection: true,
	}
	dbtesting.AssertInitialize(t, db, initReq)

	password := "myreallysecurepassword"
	createReq := newdbplugin.NewUserRequest{
		UsernameConfig: newdbplugin.UsernameMetadata{
			DisplayName: "test",
			RoleName:    "test",
		},
		Statements: newdbplugin.Statements{
			Commands: []string{mongoAdminRole},
		},
		Password:   password,
		Expiration: time.Now().Add(time.Minute),
	}
	createResp := dbtesting.AssertNewUser(t, db, createReq)
	assertCredsExist(t, createResp.Username, password, connURL)

	// Test default revocation statement
	delReq := newdbplugin.DeleteUserRequest{
		Username: createResp.Username,
	}

	dbtesting.AssertDeleteUser(t, db, delReq)

	assertCredsDoNotExist(t, createResp.Username, password, connURL)
}

func TestMongoDB_UpdateUser_Password(t *testing.T) {
	cleanup, connURL := mongodb.PrepareTestContainer(t, "latest")
	defer cleanup()

	// The docker test method PrepareTestContainer defaults to a database "test"
	// if none is provided
	connURL = connURL + "/test"
	db := new()
	defer dbtesting.AssertClose(t, db)

	initReq := newdbplugin.InitializeRequest{
		Config: map[string]interface{}{
			"connection_url": connURL,
		},
		VerifyConnection: true,
	}
	dbtesting.AssertInitialize(t, db, initReq)

	// create the database user in advance, and test the connection
	dbUser := "testmongouser"
	startingPassword := "password"
	createDBUser(t, connURL, "test", dbUser, startingPassword)

	newPassword := "myreallysecurecredentials"

	updateReq := newdbplugin.UpdateUserRequest{
		Username: dbUser,
		Password: &newdbplugin.ChangePassword{
			NewPassword: newPassword,
		},
	}
	dbtesting.AssertUpdateUser(t, db, updateReq)

	assertCredsExist(t, dbUser, newPassword, connURL)
}

func TestGetTLSAuth(t *testing.T) {
	ca := certhelpers.NewCert(t,
		certhelpers.CommonName("certificate authority"),
		certhelpers.IsCA(true),
		certhelpers.SelfSign(),
	)
	cert := certhelpers.NewCert(t,
		certhelpers.CommonName("test cert"),
		certhelpers.Parent(ca),
	)

	type testCase struct {
		username   string
		tlsCAData  []byte
		tlsKeyData []byte

		expectOpts *options.ClientOptions
		expectErr  bool
	}

	tests := map[string]testCase{
		"no TLS data set": {
			expectOpts: nil,
			expectErr:  false,
		},
		"bad CA": {
			tlsCAData: []byte("foobar"),

			expectOpts: nil,
			expectErr:  true,
		},
		"bad key": {
			tlsKeyData: []byte("foobar"),

			expectOpts: nil,
			expectErr:  true,
		},
		"good ca": {
			tlsCAData: cert.Pem,

			expectOpts: options.Client().
				SetTLSConfig(
					&tls.Config{
						RootCAs: appendToCertPool(t, x509.NewCertPool(), cert.Pem),
					},
				),
			expectErr: false,
		},
		"good key": {
			username:   "unittest",
			tlsKeyData: cert.CombinedPEM(),

			expectOpts: options.Client().
				SetTLSConfig(
					&tls.Config{
						Certificates: []tls.Certificate{cert.TLSCert},
					},
				).
				SetAuth(options.Credential{
					AuthMechanism: "MONGODB-X509",
					Username:      "unittest",
				}),
			expectErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c := new()
			c.Username = test.username
			c.TLSCAData = test.tlsCAData
			c.TLSCertificateKeyData = test.tlsKeyData

			actual, err := c.getTLSAuth()
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}
			if !reflect.DeepEqual(actual, test.expectOpts) {
				t.Fatalf("Actual:\n%#v\nExpected:\n%#v", actual, test.expectOpts)
			}
		})
	}
}

func appendToCertPool(t *testing.T, pool *x509.CertPool, caPem []byte) *x509.CertPool {
	t.Helper()

	ok := pool.AppendCertsFromPEM(caPem)
	if !ok {
		t.Fatalf("Unable to append cert to cert pool")
	}
	return pool
}

// func TestMongoDB_RotateRootCredentials(t *testing.T) {
// 	cleanup, connURL := mongodb.PrepareTestContainer(t, "latest")
// 	defer cleanup()
//
// 	// Test to ensure that we can't rotate the root creds if no username has been specified
// 	testCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
//
// 	db := new()
// 	connDetailsWithoutUsername := map[string]interface{}{
// 		"connection_url": connURL,
// 	}
// 	_, err := db.Init(testCtx, connDetailsWithoutUsername, true)
// 	if err != nil {
// 		t.Fatalf("err: %s", err)
// 	}
//
// 	// Rotate credentials should fail because no username is specified
// 	cfg, err := db.RotateRootCredentials(testCtx, nil)
// 	if err == nil {
// 		t.Fatalf("successfully rotated root credentials when no username was present")
// 	}
// 	if !reflect.DeepEqual(cfg, connDetailsWithoutUsername) {
// 		t.Fatalf("expected connection details: %#v but were %#v", connDetailsWithoutUsername, cfg)
// 	}
//
// 	db.Close()
//
// 	// Reset the database object with new connection details
// 	username := "vault-test-admin"
// 	initialPassword := "myreallysecurepassword"
//
// 	db = new()
// 	connDetailsWithUsername := map[string]interface{}{
// 		"connection_url": connURL,
// 		"username":       username,
// 		"password":       initialPassword,
// 	}
// 	_, err = db.Init(testCtx, connDetailsWithUsername, true)
// 	if err != nil {
// 		t.Fatalf("err: %s", err)
// 	}
//
// 	// Create root user
// 	createAdminUser(t, connURL, username, initialPassword)
// 	initialURL := setUserPassOnURL(t, connURL, username, initialPassword)
//
// 	// Ensure the initial root user can connect
// 	err = assertConnection(testCtx, initialURL)
// 	if err != nil {
// 		t.Fatalf("%s", err)
// 	}
//
// 	// Rotate credentials
// 	newCfg, err := db.RotateRootCredentials(testCtx, nil)
// 	if err != nil {
// 		t.Fatalf("unexpected err rotating root credentials: %s", err)
// 	}
//
// 	// Ensure the initial root user can no longer connect
// 	err = assertConnection(testCtx, initialURL)
// 	if err == nil {
// 		t.Fatalf("connection with initial credentials succeeded when it shouldn't have")
// 	}
//
// 	// Ensure the new password can connect
// 	newURL := setUserPassOnURL(t, connURL, username, newCfg["password"].(string))
// 	err = assertConnection(testCtx, newURL)
// 	if err != nil {
// 		t.Fatalf("unexpected error pinging client with new credentials: %s", err)
// 	}
// }

func createAdminUser(t *testing.T, connURL, username, password string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := createClient(ctx, connURL, nil)
	if err != nil {
		t.Fatalf("Unable to make initial connection: %s", err)
	}

	createUserCmd := createUserCommand{
		Username: username,
		Password: password,
		Roles: []interface{}{
			"userAdminAnyDatabase",
			"dbAdminAnyDatabase",
			"readWriteAnyDatabase",
		},
	}

	result := client.Database("admin").RunCommand(ctx, createUserCmd, nil)
	err = result.Err()
	if err != nil {
		t.Fatalf("Unable to create admin user: %s", err)
	}

	assertCredsExist(t, username, password, connURL)
}

func createDBUser(t testing.TB, connURL, db, username, password string) {
	t.Helper()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connURL))
	if err != nil {
		t.Fatal(err)
	}

	createUserCmd := &createUserCommand{
		Username: username,
		Password: password,
		Roles:    []interface{}{},
	}
	result := client.Database(db).RunCommand(ctx, createUserCmd, nil)
	if result.Err() != nil {
		t.Fatalf("failed to create user in mongodb: %s", result.Err())
	}

	assertCredsExist(t, username, password, connURL)
}

func assertConnection(testCtx context.Context, connURL string) error {
	// Connect as initial root user and ensure the connection is successful
	client, err := createClient(testCtx, connURL, nil)
	if err != nil {
		return fmt.Errorf("unable to create client connection with initial root user: %w", err)
	}

	err = client.Ping(testCtx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping server with initial root user: %w", err)
	}
	client.Disconnect(testCtx)
	return nil
}

func assertCredsExist(t testing.TB, username, password, connURL string) {
	t.Helper()

	connURL = strings.Replace(connURL, "mongodb://", fmt.Sprintf("mongodb://%s:%s@", username, password), 1)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connURL))
	if err != nil {
		t.Fatalf("Failed to connect to mongo: %s", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		t.Fatalf("Failed to ping mongo with user %q: %s", username, err)
	}
}

func assertCredsDoNotExist(t testing.TB, username, password, connURL string) {
	t.Helper()

	connURL = strings.Replace(connURL, "mongodb://", fmt.Sprintf("mongodb://%s:%s@", username, password), 1)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connURL))
	if err != nil {
		return // Creds don't exist as expected
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return // Creds don't exist as expected
	}
	t.Fatalf("User %q exists and was able to authenticate", username)
}

func setUserPassOnURL(t *testing.T, connURL, username, password string) string {
	t.Helper()
	uri, err := url.Parse(connURL)
	if err != nil {
		t.Fatalf("unable to parse connection URL: %s", err)
	}

	uri.User = url.UserPassword(username, password)
	return uri.String()
}

func copyConfig(config map[string]interface{}) map[string]interface{} {
	newConfig := map[string]interface{}{}
	for k, v := range config {
		newConfig[k] = v
	}
	return newConfig
}
