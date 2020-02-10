package mongodb

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers/mongodb"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connURL))
	if err != nil {
		return err
	}
	return client.Ping(ctx, readpref.Primary())
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
	testCreateDBUser(t, connURL, "test", dbUser, startingPassword)
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

func testCreateDBUser(t testing.TB, connURL, db, username, password string) {
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
		t.Fatal(result.Err())
	}
}

func TestRedact(t *testing.T) {
	type testCase struct {
		input    dbplugin.DatabaseConfig
		expected dbplugin.DatabaseConfig
	}

	tests := map[string]testCase{
		"empty returns empty": {
			input:    dbplugin.DatabaseConfig{},
			expected: dbplugin.DatabaseConfig{},
		},
		"non-redacted are unaffected": {
			input: dbplugin.DatabaseConfig{
				ConnectionDetails: map[string]interface{}{
					"foo":        "bar",
					"some thing": "omgwtfbbq",
				},
			},
			expected: dbplugin.DatabaseConfig{
				ConnectionDetails: map[string]interface{}{
					"foo":        "bar",
					"some thing": "omgwtfbbq",
				},
			},
		},
		"tls_certificate_key": {
			input: dbplugin.DatabaseConfig{
				ConnectionDetails: map[string]interface{}{
					"tls_certificate_key": "foobar",
				},
			},
			expected: dbplugin.DatabaseConfig{
				ConnectionDetails: map[string]interface{}{
					"tls_certificate_key": "*****",
				},
			},
		},
		"tls_ca": {
			input: dbplugin.DatabaseConfig{
				ConnectionDetails: map[string]interface{}{
					"tls_ca": "foobar",
				},
			},
			expected: dbplugin.DatabaseConfig{
				ConnectionDetails: map[string]interface{}{
					"tls_ca": "*****",
				},
			},
		},
		"tls_certificate_key and tls_ca": {
			input: dbplugin.DatabaseConfig{
				ConnectionDetails: map[string]interface{}{
					"tls_certificate_key": "foobar",
					"tls_ca":              "foobar",
				},
			},
			expected: dbplugin.DatabaseConfig{
				ConnectionDetails: map[string]interface{}{
					"tls_certificate_key": "*****",
					"tls_ca":              "*****",
				},
			},
		},
		"mixed values": {
			input: dbplugin.DatabaseConfig{
				ConnectionDetails: map[string]interface{}{
					"foo":                 "bar",
					"some thing":          "omgwtfbbq",
					"tls_certificate_key": "foobar",
					"tls_ca":              "foobar",
				},
			},
			expected: dbplugin.DatabaseConfig{
				ConnectionDetails: map[string]interface{}{
					"foo":                 "bar",
					"some thing":          "omgwtfbbq",
					"tls_certificate_key": "*****",
					"tls_ca":              "*****",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := new()
			actual := m.Redact(test.input)
			if !reflect.DeepEqual(actual, test.expected) {
				t.Fatalf("Redaction failed:\nActual:\n%#v\nExpected:%#v", actual, test.expected)
			}
		})
	}
}

func TestGetTLSAuth(t *testing.T) {
	// certs := MakeCerts(t)
	cert := MakeCert(t, nil)

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
			tlsCAData: cert.CertPem,

			expectOpts: options.Client().
				SetTLSConfig(
					&tls.Config{
						RootCAs: appendToCertPool(t, x509.NewCertPool(), cert.CertPem),
					},
				),
			expectErr: false,
		},
		"good key": {
			username:   "unittest",
			tlsKeyData: mergePems(cert.CertPem, cert.KeyPem),

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

func mergePems(pems ...[]byte) []byte {
	return bytes.Join(pems, []byte("\n"))
}
