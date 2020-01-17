package mysqlhelper

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/testhelpers/docker"
	"github.com/ory/dockertest"
)

func PrepareMySQLTestContainer(t *testing.T, legacy bool, pw string) (cleanup func(), retURL string) {
	if os.Getenv("MYSQL_URL") != "" {
		return func() {}, os.Getenv("MYSQL_URL")
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	imageVersion := "5.7"
	if legacy {
		imageVersion = "5.6"
	}

	resource, err := pool.Run("mysql", imageVersion, []string{"MYSQL_ROOT_PASSWORD=" + pw})
	if err != nil {
		t.Fatalf("Could not start local MySQL docker container: %s", err)
	}

	cleanup = func() {
		docker.CleanupResource(t, pool, resource)
	}

	retURL = fmt.Sprintf("root:%s@(localhost:%s)/mysql?parseTime=true", pw, resource.GetPort("3306/tcp"))

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		var err error
		var db *sql.DB
		db, err = sql.Open("mysql", retURL)
		if err != nil {
			return err
		}
		defer db.Close()
		return db.Ping()
	}); err != nil {
		cleanup()
		t.Fatalf("Could not connect to MySQL docker container: %s", err)
	}

	return
}

func TestCredsExist(t testing.TB, connURL, username, password string) error {
	// Log in with the new creds
	connURL = strings.Replace(connURL, "root:secret", fmt.Sprintf("%s:%s", username, password), 1)
	db, err := sql.Open("mysql", connURL)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}
