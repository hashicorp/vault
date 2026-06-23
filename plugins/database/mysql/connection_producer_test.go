// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package mysql

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	paths "path"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers/certhelpers"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	dockertest "github.com/ory/dockertest/v3"
)

// TestInit_TLSRegistrationReusedForSameTLSConfig verifies re-init with the same TLS input refreshes registration without leaking state.
func TestInit_TLSRegistrationReusedForSameTLSConfig(t *testing.T) {
	ca := certhelpers.NewCert(
		t,
		certhelpers.CommonName("test-ca"),
		certhelpers.IsCA(true),
		certhelpers.SelfSign(),
	)

	registerCalls := 0
	deregisterCalls := 0
	registeredKeys := []string{}
	deregisteredKeys := []string{}
	registerTLSConfig := func(key string, cfg *tls.Config) error {
		registerCalls++
		registeredKeys = append(registeredKeys, key)
		return nil
	}
	deregisterTLSConfig := func(key string) {
		deregisterCalls++
		deregisteredKeys = append(deregisteredKeys, key)
	}

	p := &mySQLConnectionProducer{registerTLSConfig: registerTLSConfig, deregisterTLSConfig: deregisterTLSConfig}
	conf := map[string]interface{}{
		"connection_url": "{{username}}:{{password}}@tcp(localhost:3306)/test",
		"username":       "user",
		"password":       "pass",
		"tls_ca":         ca.Pem,
	}

	if _, err := p.Init(context.Background(), conf, false); err != nil {
		t.Fatalf("first init failed: %s", err)
	}
	if _, err := p.Init(context.Background(), conf, false); err != nil {
		t.Fatalf("second init failed: %s", err)
	}

	if registerCalls != 2 {
		t.Fatalf("expected register to be called twice, got %d", registerCalls)
	}
	if deregisterCalls != 1 {
		t.Fatalf("expected deregister to be called once, got %d", deregisterCalls)
	}
	if len(registeredKeys) != 2 || len(deregisteredKeys) != 1 {
		t.Fatalf("unexpected key tracking sizes, registered=%d deregistered=%d", len(registeredKeys), len(deregisteredKeys))
	}
	if deregisteredKeys[0] != registeredKeys[0] {
		t.Fatalf("expected deregistered key %q to match first registered key %q", deregisteredKeys[0], registeredKeys[0])
	}
	if deregisteredKeys[0] == registeredKeys[1] {
		t.Fatalf("expected deregistered key to differ from newly registered replacement key")
	}
}

// TestInit_TLSRegistrationRefreshedForChangedTLSConfig verifies TLS config re-init refreshes driver registration.
func TestInit_TLSRegistrationRefreshedForChangedTLSConfig(t *testing.T) {
	caA := certhelpers.NewCert(
		t,
		certhelpers.CommonName("test-ca-a"),
		certhelpers.IsCA(true),
		certhelpers.SelfSign(),
	)
	caB := certhelpers.NewCert(
		t,
		certhelpers.CommonName("test-ca-b"),
		certhelpers.IsCA(true),
		certhelpers.SelfSign(),
	)

	registerCalls := 0
	deregisterCalls := 0
	registeredKeys := []string{}
	deregisteredKeys := []string{}
	registerTLSConfig := func(key string, cfg *tls.Config) error {
		registerCalls++
		registeredKeys = append(registeredKeys, key)
		return nil
	}
	deregisterTLSConfig := func(key string) {
		deregisterCalls++
		deregisteredKeys = append(deregisteredKeys, key)
	}

	p := &mySQLConnectionProducer{registerTLSConfig: registerTLSConfig, deregisterTLSConfig: deregisterTLSConfig}
	base := map[string]interface{}{
		"connection_url": "{{username}}:{{password}}@tcp(localhost:3306)/test",
		"username":       "user",
		"password":       "pass",
	}

	confA := map[string]interface{}{}
	for k, v := range base {
		confA[k] = v
	}
	confA["tls_ca"] = caA.Pem

	confB := map[string]interface{}{}
	for k, v := range base {
		confB[k] = v
	}
	confB["tls_ca"] = caB.Pem

	if _, err := p.Init(context.Background(), confA, false); err != nil {
		t.Fatalf("first init failed: %s", err)
	}
	if _, err := p.Init(context.Background(), confB, false); err != nil {
		t.Fatalf("second init failed: %s", err)
	}

	if registerCalls != 2 {
		t.Fatalf("expected register to be called twice, got %d", registerCalls)
	}
	if deregisterCalls != 1 {
		t.Fatalf("expected deregister to be called once, got %d", deregisterCalls)
	}
	if len(registeredKeys) != 2 || len(deregisteredKeys) != 1 {
		t.Fatalf("unexpected key tracking sizes, registered=%d deregistered=%d", len(registeredKeys), len(deregisteredKeys))
	}
	if deregisteredKeys[0] != registeredKeys[0] {
		t.Fatalf("expected deregistered key %q to match first registered key %q", deregisteredKeys[0], registeredKeys[0])
	}
	if deregisteredKeys[0] == registeredKeys[1] {
		t.Fatalf("expected deregistered key to differ from newly registered replacement key")
	}
}

// TestClose_DeregistersTLSConfig verifies Close deregisters any registered TLS config and is idempotent.
func TestClose_DeregistersTLSConfig(t *testing.T) {
	deregisterCalls := 0
	deregisteredKeys := []string{}
	deregisterTLSConfig := func(key string) {
		deregisterCalls++
		deregisteredKeys = append(deregisteredKeys, key)
	}

	testDB, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/test")
	if err != nil {
		t.Fatalf("failed to create test DB handle: %s", err)
	}

	p := &mySQLConnectionProducer{tlsConfigName: "test-tls-key", db: testDB, deregisterTLSConfig: deregisterTLSConfig}
	if err := p.Close(); err != nil {
		t.Fatalf("close failed: %s", err)
	}
	if err := p.Close(); err != nil {
		t.Fatalf("second close failed: %s", err)
	}

	if deregisterCalls != 1 {
		t.Fatalf("expected deregister to be called once, got %d", deregisterCalls)
	}
	if len(deregisteredKeys) != 1 || deregisteredKeys[0] != "test-tls-key" {
		t.Fatalf("expected deregistered key to be %q, got %v", "test-tls-key", deregisteredKeys)
	}
	if p.db != nil {
		t.Fatalf("expected db to be cleared")
	}
	if p.tlsConfigName != "" {
		t.Fatalf("expected tlsConfigName to be cleared")
	}
}

// TestInit_NoTLSClearsPreviousTLSRegistration verifies re-init without TLS clears previous TLS registration.
func TestInit_NoTLSClearsPreviousTLSRegistration(t *testing.T) {
	ca := certhelpers.NewCert(
		t,
		certhelpers.CommonName("test-ca"),
		certhelpers.IsCA(true),
		certhelpers.SelfSign(),
	)

	registerCalls := 0
	deregisterCalls := 0
	registeredKeys := []string{}
	deregisteredKeys := []string{}
	registerTLSConfig := func(key string, cfg *tls.Config) error {
		registerCalls++
		registeredKeys = append(registeredKeys, key)
		return nil
	}
	deregisterTLSConfig := func(key string) {
		deregisterCalls++
		deregisteredKeys = append(deregisteredKeys, key)
	}

	p := &mySQLConnectionProducer{registerTLSConfig: registerTLSConfig, deregisterTLSConfig: deregisterTLSConfig}
	withTLS := map[string]interface{}{
		"connection_url": "{{username}}:{{password}}@tcp(localhost:3306)/test",
		"username":       "user",
		"password":       "pass",
		"tls_ca":         ca.Pem,
	}
	withoutTLS := map[string]interface{}{
		"connection_url":      "{{username}}:{{password}}@tcp(localhost:3306)/test",
		"username":            "user",
		"password":            "pass",
		"tls_ca":              "",
		"tls_certificate_key": "",
	}

	if _, err := p.Init(context.Background(), withTLS, false); err != nil {
		t.Fatalf("init with tls failed: %s", err)
	}
	if _, err := p.Init(context.Background(), withoutTLS, false); err != nil {
		t.Fatalf("init without tls failed: %s", err)
	}

	if registerCalls != 1 {
		t.Fatalf("expected register to be called once, got %d", registerCalls)
	}
	if deregisterCalls != 1 {
		t.Fatalf("expected deregister to be called once, got %d", deregisterCalls)
	}
	if len(registeredKeys) != 1 || len(deregisteredKeys) != 1 {
		t.Fatalf("unexpected key tracking sizes, registered=%d deregistered=%d", len(registeredKeys), len(deregisteredKeys))
	}
	if deregisteredKeys[0] != registeredKeys[0] {
		t.Fatalf("expected deregistered key %q to match registered key %q", deregisteredKeys[0], registeredKeys[0])
	}
	if p.tlsConfigName != "" {
		t.Fatalf("expected tlsConfigName to be cleared")
	}
}

// TestEnsureTLSRegistration_ReplacementFailureKeepsPreviousRegistration verifies replacement failure preserves the prior TLS registration state.
func TestEnsureTLSRegistration_ReplacementFailureKeepsPreviousRegistration(t *testing.T) {
	registerCalls := 0
	deregisterCalls := 0
	registerTLSConfig := func(key string, cfg *tls.Config) error {
		registerCalls++
		return fmt.Errorf("forced register failure")
	}
	deregisterTLSConfig := func(key string) {
		deregisterCalls++
	}

	p := &mySQLConnectionProducer{tlsConfigName: "existing-tls-key", registerTLSConfig: registerTLSConfig, deregisterTLSConfig: deregisterTLSConfig}
	err := p.ensureTLSRegistration(&tls.Config{})
	if err == nil {
		t.Fatalf("expected ensureTLSRegistration to fail")
	}

	if registerCalls != 1 {
		t.Fatalf("expected register to be called once, got %d", registerCalls)
	}
	if deregisterCalls != 0 {
		t.Fatalf("expected deregister not to be called, got %d", deregisterCalls)
	}
	if p.tlsConfigName != "existing-tls-key" {
		t.Fatalf("expected tlsConfigName to remain unchanged, got %q", p.tlsConfigName)
	}
}

func Test_addTLStoDSN(t *testing.T) {
	type testCase struct {
		rootUrl        string
		tlsConfigName  string
		expectedResult string
	}

	tests := map[string]testCase{
		"no tls, no query string": {
			rootUrl:        "user:password@tcp(localhost:3306)/test",
			tlsConfigName:  "",
			expectedResult: "user:password@tcp(localhost:3306)/test",
		},
		"tls, no query string": {
			rootUrl:        "user:password@tcp(localhost:3306)/test",
			tlsConfigName:  "tlsTest101",
			expectedResult: "user:password@tcp(localhost:3306)/test?tls=tlsTest101",
		},
		"tls, query string": {
			rootUrl:        "user:password@tcp(localhost:3306)/test?foo=bar",
			tlsConfigName:  "tlsTest101",
			expectedResult: "user:password@tcp(localhost:3306)/test?tls=tlsTest101&foo=bar",
		},
		"tls, query string, ? in password": {
			rootUrl:        "user:pa?ssword?@tcp(localhost:3306)/test?foo=bar",
			tlsConfigName:  "tlsTest101",
			expectedResult: "user:pa?ssword?@tcp(localhost:3306)/test?tls=tlsTest101&foo=bar",
		},
		"tls, valid tls parameter in query string": {
			rootUrl:        "user:password@tcp(localhost:3306)/test?tls=true",
			tlsConfigName:  "",
			expectedResult: "user:password@tcp(localhost:3306)/test?tls=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tCase := mySQLConnectionProducer{
				ConnectionURL: test.rootUrl,
				tlsConfigName: test.tlsConfigName,
			}

			actual, err := tCase.addTLStoDSN()
			if err != nil {
				t.Fatalf("error occurred in test: %s", err)
			}
			if actual != test.expectedResult {
				t.Fatalf("generated: %s, expected: %s", actual, test.expectedResult)
			}
		})
	}
}

func TestInit_clientTLS(t *testing.T) {
	t.Skip("Skipping this test because CircleCI can't mount the files we need without further investigation: " +
		"https://support.circleci.com/hc/en-us/articles/360007324514-How-can-I-mount-volumes-to-docker-containers-")

	// Set up temp directory so we can mount it to the docker container
	confDir := makeTempDir(t)
	defer os.RemoveAll(confDir)

	// Create certificates for MySQL authentication
	caCert := certhelpers.NewCert(
		t,
		certhelpers.CommonName("test certificate authority"),
		certhelpers.IsCA(true),
		certhelpers.SelfSign(),
	)
	serverCert := certhelpers.NewCert(
		t,
		certhelpers.CommonName("server"),
		certhelpers.DNS("localhost"),
		certhelpers.Parent(caCert),
	)
	clientCert := certhelpers.NewCert(
		t,
		certhelpers.CommonName("client"),
		certhelpers.DNS("client"),
		certhelpers.Parent(caCert),
	)

	writeFile(t, paths.Join(confDir, "ca.pem"), caCert.CombinedPEM(), 0o644)
	writeFile(t, paths.Join(confDir, "server-cert.pem"), serverCert.Pem, 0o644)
	writeFile(t, paths.Join(confDir, "server-key.pem"), serverCert.PrivateKeyPEM(), 0o644)
	writeFile(t, paths.Join(confDir, "client.pem"), clientCert.CombinedPEM(), 0o644)

	// //////////////////////////////////////////////////////
	// Set up MySQL config file
	rawConf := `
[mysqld]
ssl
ssl-ca=/etc/mysql/ca.pem
ssl-cert=/etc/mysql/server-cert.pem
ssl-key=/etc/mysql/server-key.pem`

	writeFile(t, paths.Join(confDir, "my.cnf"), []byte(rawConf), 0o644)

	// //////////////////////////////////////////////////////
	// Start MySQL container
	retURL, cleanup := startMySQLWithTLS(t, "5.7", confDir)
	defer cleanup()

	// //////////////////////////////////////////////////////
	// Set up x509 user
	mClient := connect(t, retURL)

	username := setUpX509User(t, mClient, clientCert)

	// //////////////////////////////////////////////////////
	// Test
	mysql := newMySQL(DefaultUserNameTemplate)

	conf := map[string]interface{}{
		"connection_url":      retURL,
		"username":            username,
		"tls_certificate_key": clientCert.CombinedPEM(),
		"tls_ca":              caCert.Pem,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := mysql.Init(ctx, conf, true)
	if err != nil {
		t.Fatalf("Unable to initialize mysql engine: %s", err)
	}

	// Initialization complete. The connection was established, but we need to ensure
	// that we're connected as the right user
	whoamiCmd := "SELECT CURRENT_USER()"

	client, err := mysql.getConnection(ctx)
	if err != nil {
		t.Fatalf("Unable to make connection to MySQL: %s", err)
	}
	stmt, err := client.Prepare(whoamiCmd)
	if err != nil {
		t.Fatalf("Unable to prepare MySQL statementL %s", err)
	}

	results := stmt.QueryRow()

	expected := fmt.Sprintf("%s@%%", username)

	var result string
	if err := results.Scan(&result); err != nil {
		t.Fatalf("result could not be scanned from result set: %s", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Actual:%#v\nExpected:\n%#v", result, expected)
	}
}

func makeTempDir(t *testing.T) (confDir string) {
	confDir, err := ioutil.TempDir(".", "mysql-test-data")
	if err != nil {
		t.Fatalf("Unable to make temp directory: %s", err)
	}
	// Convert the directory to an absolute path because docker needs it when mounting
	confDir, err = filepath.Abs(filepath.Clean(confDir))
	if err != nil {
		t.Fatalf("Unable to determine where temp directory is on absolute path: %s", err)
	}
	return confDir
}

func startMySQLWithTLS(t *testing.T, version string, confDir string) (retURL string, cleanup func()) {
	if os.Getenv("MYSQL_URL") != "" {
		return os.Getenv("MYSQL_URL"), func() {}
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}
	pool.MaxWait = 30 * time.Second

	containerName := "mysql-unit-test"

	// Remove previously running container if it is still running because cleanup failed
	err = pool.RemoveContainerByName(containerName)
	if err != nil {
		t.Fatalf("Unable to remove old running containers: %s", err)
	}

	username := "root"
	password := "x509test"

	runOpts := &dockertest.RunOptions{
		Name:       containerName,
		Repository: "mysql",
		Tag:        version,
		Cmd:        []string{"--defaults-extra-file=/etc/mysql/my.cnf", "--auto-generate-certs=OFF"},
		Env:        []string{fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", password)},
		// Mount the directory from local filesystem into the container
		Mounts: []string{
			fmt.Sprintf("%s:/etc/mysql", confDir),
		},
	}

	resource, err := pool.RunWithOptions(runOpts)
	if err != nil {
		t.Fatalf("Could not start local mysql docker container: %s", err)
	}
	resource.Expire(30)

	cleanup = func() {
		err := pool.Purge(resource)
		if err != nil {
			t.Fatalf("Failed to cleanup local container: %s", err)
		}
	}

	dsn := fmt.Sprintf("{{username}}:{{password}}@tcp(localhost:%s)/mysql", resource.GetPort("3306/tcp"))

	url := dbutil.QueryHelper(dsn, map[string]string{
		"username": username,
		"password": password,
	})
	// exponential backoff-retry
	err = pool.Retry(func() error {
		var err error

		db, err := sql.Open("mysql", url)
		if err != nil {
			t.Logf("err: %s", err)
			return err
		}
		defer db.Close()
		return db.Ping()
	})
	if err != nil {
		cleanup()
		t.Fatalf("Could not connect to mysql docker container: %s", err)
	}

	return dsn, cleanup
}

func connect(t *testing.T, dsn string) (db *sql.DB) {
	url := dbutil.QueryHelper(dsn, map[string]string{
		"username": "root",
		"password": "x509test",
	})

	db, err := sql.Open("mysql", url)
	if err != nil {
		t.Fatalf("Unable to make connection to MySQL: %s", err)
	}

	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping MySQL server: %s", err)
	}

	return db
}

func setUpX509User(t *testing.T, db *sql.DB, cert certhelpers.Certificate) (username string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username = cert.Template.Subject.CommonName

	cmds := []string{
		fmt.Sprintf("CREATE USER %s IDENTIFIED BY '' REQUIRE X509", username),
		fmt.Sprintf("GRANT ALL ON mysql.* TO '%s'@'%s' REQUIRE X509", username, "%"),
	}

	for _, cmd := range cmds {
		stmt, err := db.PrepareContext(ctx, cmd)
		if err != nil {
			t.Fatalf("Failed to prepare query: %s", err)
		}

		_, err = stmt.ExecContext(ctx)
		if err != nil {
			t.Fatalf("Failed to create x509 user in database: %s", err)
		}
		err = stmt.Close()
		if err != nil {
			t.Fatalf("Failed to close prepared statement: %s", err)
		}
	}

	return username
}

// ////////////////////////////////////////////////////////////////////////////
// Writing to file
// ////////////////////////////////////////////////////////////////////////////
func writeFile(t *testing.T, filename string, data []byte, perms os.FileMode) {
	t.Helper()

	err := ioutil.WriteFile(filename, data, perms)
	if err != nil {
		t.Fatalf("Unable to write to file [%s]: %s", filename, err)
	}
}
