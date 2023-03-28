// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mysql

import (
	"context"
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
	caCert := certhelpers.NewCert(t,
		certhelpers.CommonName("test certificate authority"),
		certhelpers.IsCA(true),
		certhelpers.SelfSign(),
	)
	serverCert := certhelpers.NewCert(t,
		certhelpers.CommonName("server"),
		certhelpers.DNS("localhost"),
		certhelpers.Parent(caCert),
	)
	clientCert := certhelpers.NewCert(t,
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
