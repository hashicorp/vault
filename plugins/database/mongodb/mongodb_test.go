package mongodb

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/hashicorp/vault/helper/testhelpers/certhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/mongodb"
	dbplugin "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	dbtesting "github.com/hashicorp/vault/sdk/database/dbplugin/v5/testing"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	mongoAdminRole       = `{ "db": "admin", "roles": [ { "role": "readWrite" } ] }`
	mongoTestDBAdminRole = `{ "db": "test", "roles": [ { "role": "readWrite" } ] }`
)

func TestMongoDB_Initialize(t *testing.T) {
	cleanup, connURL := mongodb.PrepareTestContainer(t, "5.0.10")
	defer cleanup()

	db := new()
	defer dbtesting.AssertClose(t, db)

	config := map[string]interface{}{
		"connection_url": connURL,
	}

	// Make a copy since the original map could be modified by the Initialize call
	expectedConfig := copyConfig(config)

	req := dbplugin.InitializeRequest{
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

func TestNewUser_usernameTemplate(t *testing.T) {
	type testCase struct {
		usernameTemplate string

		newUserReq            dbplugin.NewUserRequest
		expectedUsernameRegex string
	}

	tests := map[string]testCase{
		"default username template": {
			usernameTemplate: "",

			newUserReq: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "token",
					RoleName:    "testrolenamewithmanycharacters",
				},
				Statements: dbplugin.Statements{
					Commands: []string{mongoAdminRole},
				},
				Password:   "98yq3thgnakjsfhjkl",
				Expiration: time.Now().Add(time.Minute),
			},

			expectedUsernameRegex: "^v-token-testrolenamewit-[a-zA-Z0-9]{20}-[0-9]{10}$",
		},
		"default username template with invalid chars": {
			usernameTemplate: "",

			newUserReq: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "a.bad.account",
					RoleName:    "a.bad.role",
				},
				Statements: dbplugin.Statements{
					Commands: []string{mongoAdminRole},
				},
				Password:   "98yq3thgnakjsfhjkl",
				Expiration: time.Now().Add(time.Minute),
			},

			expectedUsernameRegex: "^v-a-bad-account-a-bad-role-[a-zA-Z0-9]{20}-[0-9]{10}$",
		},
		"custom username template": {
			usernameTemplate: "{{random 2 | uppercase}}_{{unix_time}}_{{.RoleName | uppercase}}_{{.DisplayName | uppercase}}",

			newUserReq: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "token",
					RoleName:    "testrolenamewithmanycharacters",
				},
				Statements: dbplugin.Statements{
					Commands: []string{mongoAdminRole},
				},
				Password:   "98yq3thgnakjsfhjkl",
				Expiration: time.Now().Add(time.Minute),
			},

			expectedUsernameRegex: "^[A-Z0-9]{2}_[0-9]{10}_TESTROLENAMEWITHMANYCHARACTERS_TOKEN$",
		},
		"admin in test database username template": {
			usernameTemplate: "",

			newUserReq: dbplugin.NewUserRequest{
				UsernameConfig: dbplugin.UsernameMetadata{
					DisplayName: "token",
					RoleName:    "testrolenamewithmanycharacters",
				},
				Statements: dbplugin.Statements{
					Commands: []string{mongoTestDBAdminRole},
				},
				Password:   "98yq3thgnakjsfhjkl",
				Expiration: time.Now().Add(time.Minute),
			},

			expectedUsernameRegex: "^v-token-testrolenamewit-[a-zA-Z0-9]{20}-[0-9]{10}$",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cleanup, connURL := mongodb.PrepareTestContainer(t, "5.0.10")
			defer cleanup()

			if name == "admin in test database username template" {
				connURL = connURL + "/test?authSource=test"
			}

			db := new()
			defer dbtesting.AssertClose(t, db)

			initReq := dbplugin.InitializeRequest{
				Config: map[string]interface{}{
					"connection_url":    connURL,
					"username_template": test.usernameTemplate,
				},
				VerifyConnection: true,
			}
			dbtesting.AssertInitialize(t, db, initReq)

			ctx := context.Background()
			newUserResp, err := db.NewUser(ctx, test.newUserReq)
			require.NoError(t, err)
			require.Regexp(t, test.expectedUsernameRegex, newUserResp.Username)

			assertCredsExist(t, newUserResp.Username, test.newUserReq.Password, connURL)
		})
	}
}

func TestMongoDB_CreateUser(t *testing.T) {
	cleanup, connURL := mongodb.PrepareTestContainer(t, "5.0.10")
	defer cleanup()

	db := new()
	defer dbtesting.AssertClose(t, db)

	initReq := dbplugin.InitializeRequest{
		Config: map[string]interface{}{
			"connection_url": connURL,
		},
		VerifyConnection: true,
	}
	dbtesting.AssertInitialize(t, db, initReq)

	password := "myreallysecurepassword"
	createReq := dbplugin.NewUserRequest{
		UsernameConfig: dbplugin.UsernameMetadata{
			DisplayName: "test",
			RoleName:    "test",
		},
		Statements: dbplugin.Statements{
			Commands: []string{mongoAdminRole},
		},
		Password:   password,
		Expiration: time.Now().Add(time.Minute),
	}
	createResp := dbtesting.AssertNewUser(t, db, createReq)

	assertCredsExist(t, createResp.Username, password, connURL)
}

func TestMongoDB_CreateUser_writeConcern(t *testing.T) {
	cleanup, connURL := mongodb.PrepareTestContainer(t, "5.0.10")
	defer cleanup()

	initReq := dbplugin.InitializeRequest{
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
	createReq := dbplugin.NewUserRequest{
		UsernameConfig: dbplugin.UsernameMetadata{
			DisplayName: "test",
			RoleName:    "test",
		},
		Statements: dbplugin.Statements{
			Commands: []string{mongoAdminRole},
		},
		Password:   password,
		Expiration: time.Now().Add(time.Minute),
	}
	createResp := dbtesting.AssertNewUser(t, db, createReq)

	assertCredsExist(t, createResp.Username, password, connURL)
}

func TestMongoDB_DeleteUser(t *testing.T) {
	cleanup, connURL := mongodb.PrepareTestContainer(t, "5.0.10")
	defer cleanup()

	db := new()
	defer dbtesting.AssertClose(t, db)

	initReq := dbplugin.InitializeRequest{
		Config: map[string]interface{}{
			"connection_url": connURL,
		},
		VerifyConnection: true,
	}
	dbtesting.AssertInitialize(t, db, initReq)

	password := "myreallysecurepassword"
	createReq := dbplugin.NewUserRequest{
		UsernameConfig: dbplugin.UsernameMetadata{
			DisplayName: "test",
			RoleName:    "test",
		},
		Statements: dbplugin.Statements{
			Commands: []string{mongoAdminRole},
		},
		Password:   password,
		Expiration: time.Now().Add(time.Minute),
	}
	createResp := dbtesting.AssertNewUser(t, db, createReq)
	assertCredsExist(t, createResp.Username, password, connURL)

	// Test default revocation statement
	delReq := dbplugin.DeleteUserRequest{
		Username: createResp.Username,
	}

	dbtesting.AssertDeleteUser(t, db, delReq)

	assertCredsDoNotExist(t, createResp.Username, password, connURL)
}

func TestMongoDB_UpdateUser_Password(t *testing.T) {
	cleanup, connURL := mongodb.PrepareTestContainer(t, "5.0.10")
	defer cleanup()

	// The docker test method PrepareTestContainer defaults to a database "test"
	// if none is provided
	connURL = connURL + "/test"
	db := new()
	defer dbtesting.AssertClose(t, db)

	initReq := dbplugin.InitializeRequest{
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

	updateReq := dbplugin.UpdateUserRequest{
		Username: dbUser,
		Password: &dbplugin.ChangePassword{
			NewPassword: newPassword,
		},
	}
	dbtesting.AssertUpdateUser(t, db, updateReq)

	assertCredsExist(t, dbUser, newPassword, connURL)
}

func TestMongoDB_RotateRoot_NonAdminDB(t *testing.T) {
	cleanup, connURL := mongodb.PrepareTestContainer(t, "latest")
	defer cleanup()

	connURL = connURL + "/test?authSource=test"
	db := new()
	defer dbtesting.AssertClose(t, db)

	initReq := dbplugin.InitializeRequest{
		Config: map[string]interface{}{
			"connection_url": connURL,
		},
		VerifyConnection: true,
	}
	dbtesting.AssertInitialize(t, db, initReq)

	dbUser := "testmongouser"
	startingPassword := "password"
	createDBUser(t, connURL, "test", dbUser, startingPassword)

	newPassword := "myreallysecurecredentials"

	updateReq := dbplugin.UpdateUserRequest{
		Username: dbUser,
		Password: &dbplugin.ChangePassword{
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
			assertDeepEqual(t, test.expectOpts, actual)
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

var cmpClientOptionsOpts = cmp.Options{
	cmp.AllowUnexported(options.ClientOptions{}),

	cmp.AllowUnexported(tls.Config{}),
	cmpopts.IgnoreTypes(sync.Mutex{}, sync.RWMutex{}),

	// 'lazyCerts' has a func field which can't be compared.
	cmpopts.IgnoreFields(x509.CertPool{}, "lazyCerts"),
	cmp.AllowUnexported(x509.CertPool{}),
}

// Need a special comparison for ClientOptions because reflect.DeepEquals won't work in Go 1.16.
// See: https://github.com/golang/go/issues/45891
func assertDeepEqual(t *testing.T, a, b *options.ClientOptions) {
	t.Helper()

	if diff := cmp.Diff(a, b, cmpClientOptionsOpts); diff != "" {
		t.Fatalf("assertion failed: values are not equal\n--- expected\n+++ actual\n%v", diff)
	}
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

func copyConfig(config map[string]interface{}) map[string]interface{} {
	newConfig := map[string]interface{}{}
	for k, v := range config {
		newConfig[k] = v
	}
	return newConfig
}
