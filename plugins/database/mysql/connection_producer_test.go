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
	"github.com/ory/dockertest"
)

func TestInit_clientTLS(t *testing.T) {
	//t.Skip("Skipping this test because CircleCI can't mount the files we need without further investigation: " +
	//	"https://support.circleci.com/hc/en-us/articles/360007324514-How-can-I-mount-volumes-to-docker-containers-")

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
		certhelpers.Dns("localhost"),
		certhelpers.Parent(caCert),
	)
	clientCert := certhelpers.NewCert(t,
		certhelpers.CommonName("client"),
		certhelpers.Dns("client"),
		certhelpers.Parent(caCert),
	)

	certhelpers.WriteFile(t, paths.Join(confDir, "ca.pem"), caCert.CombinedPEM(), 0644)
	certhelpers.WriteFile(t, paths.Join(confDir, "server-cert.pem"), serverCert.Pem, 0644)
	certhelpers.WriteFile(t, paths.Join(confDir, "server-key.pem"), serverCert.PrivateKeyPEM(), 0644)
	certhelpers.WriteFile(t, paths.Join(confDir, "client.pem"), clientCert.CombinedPEM(), 0644)

	// //////////////////////////////////////////////////////
	// Set up MySQL config file
	rawConf := `
[mysqld]
ssl
ssl-ca=/etc/mysql/ca.pem
ssl-cert=/etc/mysql/server-cert.pem
ssl-key=/etc/mysql/server-key.pem`

	certhelpers.WriteFile(t, paths.Join(confDir, "my.cnf"), []byte(rawConf), 0644)

	// //////////////////////////////////////////////////////
	// Start MySQL container
	retURL, cleanup := startMySQLWithTLS(t, "5.7", confDir)
	defer cleanup()

	// //////////////////////////////////////////////////////
	// Set up x509 user
	mClient := connect(t, retURL)

	setUpX509User(t, mClient, clientCert)

	// //////////////////////////////////////////////////////
	// Test
	mysql := new(25, 25, 25)

	conf := map[string]interface{}{
		"connection_url":      retURL,
		//"tls_certificate_key": clientCert.CombinedPEM(),
		//"tls_ca":              caCert.Pem,
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

	results, err := stmt.QueryContext(ctx)

	if err != nil {
		t.Fatalf("Unable to execute MySQL query: %s", err)
	}

	expected := connStatus{
		AuthInfo: authInfo{
			AuthenticatedUsers: []user{
				{
					User: fmt.Sprintf("CN=%s", clientCert.Template.Subject.CommonName),
					DB:   "$external",
				},
			},
			AuthenticatedUserRoles: []role{
				{
					Role: "readWrite",
					DB:   "test",
				},
				{
					Role: "userAdminAnyDatabase",
					DB:   "admin",
				},
			},
		},
		Ok: 1,
	}

	actual := results

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Actual:%#v\nExpected:\n%#v", actual, expected)
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

	runOpts := &dockertest.RunOptions{
		Name:       containerName,
		Repository: "mysql",
		Tag:        version,
		Cmd:        []string{"--defaults-extra-file=/etc/mysql/my.cnf", "--auto-generate-certs=OFF"},
		Env:				[]string{"MYSQL_ROOT_PASSWORD=x509test"},
		// Mount the directory from local filesystem into the container
		Mounts: []string{
			fmt.Sprintf("%s:/etc/mysql", confDir),
		},
	}

	resource, err := pool.RunWithOptions(runOpts)
	if err != nil {
		t.Fatalf("Could not start local mysql docker container: %s", err)
	}
	resource.Expire(60)

	cleanup = func() {
		err := pool.Purge(resource)
		if err != nil {
			t.Fatalf("Failed to cleanup local container: %s", err)
		}
	}

	dsn := fmt.Sprintf("root:x509test@tcp(localhost:%s)/mysql", resource.GetPort("3306/tcp"))

	// exponential backoff-retry
	err = pool.Retry(func() error {
		var err error

		t.Logf("dsn: %s", dsn)
		db, err := sql.Open("mysql", dsn)
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
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Unable to make connection to MySQL: %s", err)
	}

	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping MySQL server: %s", err)
	}

	return db
}

func setUpX509User(t *testing.T, db *sql.DB, cert certhelpers.Certificate) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username := cert.Template.Subject.CommonName

	cmd := fmt.Sprintf("CREATE USER %s REQUIRE X509", username)

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

type connStatus struct {
	AuthInfo authInfo `bson:"authInfo"`
	Ok       int      `bson:"ok"`
}

type authInfo struct {
	AuthenticatedUsers     []user `bson:"authenticatedUsers"`
	AuthenticatedUserRoles roles  `bson:"authenticatedUserRoles"`
}

type user struct {
	User string `bson:"user"`
	DB   string `bson:"db"`
}

type role struct {
	Role string `bson:"role"`
	DB   string `bson:"db"`
}

type roles []role

func (r roles) Len() int           { return len(r) }
func (r roles) Less(i, j int) bool { return r[i].Role < r[j].Role }
func (r roles) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
