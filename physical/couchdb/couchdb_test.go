package couchdb

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/physical"
	"github.com/ory/dockertest"
)

func TestCouchDBBackend(t *testing.T) {
	cleanup, endpoint, username, password := prepareCouchdbDBTestContainer(t)
	defer cleanup()

	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewCouchDBBackend(map[string]string{
		"endpoint": endpoint,
		"username": username,
		"password": password,
	}, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}

func TestTransactionalCouchDBBackend(t *testing.T) {
	cleanup, endpoint, username, password := prepareCouchdbDBTestContainer(t)
	defer cleanup()

	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewTransactionalCouchDBBackend(map[string]string{
		"endpoint": endpoint,
		"username": username,
		"password": password,
	}, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}

func prepareCouchdbDBTestContainer(t *testing.T) (cleanup func(), retAddress, username, password string) {
	// If environment variable is set, assume caller wants to target a real
	// DynamoDB.
	if os.Getenv("COUCHDB_ENDPOINT") != "" {
		return func() {}, os.Getenv("COUCHDB_ENDPOINT"), os.Getenv("COUCHDB_USERNAME"), os.Getenv("COUCHDB_PASSWORD")
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	resource, err := pool.Run("couchdb", "1.6", []string{})
	if err != nil {
		t.Fatalf("Could not start local DynamoDB: %s", err)
	}

	retAddress = "http://localhost:" + resource.GetPort("5984/tcp")
	cleanup = func() {
		err := pool.Purge(resource)
		if err != nil {
			t.Fatalf("Failed to cleanup local DynamoDB: %s", err)
		}
	}

	// exponential backoff-retry, because the couchDB may not be able to accept
	// connections yet
	if err := pool.Retry(func() error {
		var err error
		resp, err := http.Get(retAddress)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("expected couchdb to return status code 200, got (%s) instead.", resp.Status)
		}
		return nil
	}); err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}

	dbName := fmt.Sprintf("vault-test-%d", time.Now().Unix())
	{
		req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%s", retAddress, dbName), nil)
		if err != nil {
			t.Fatalf("Could not create create database request: %q", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Could not create database: %q", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			bs, _ := ioutil.ReadAll(resp.Body)
			t.Fatalf("Failed to create database: %s %s\n", resp.Status, string(bs))
		}
	}
	{
		req, err := http.NewRequest("PUT", fmt.Sprintf("%s/_config/admins/admin", retAddress), strings.NewReader(`"admin"`))
		if err != nil {
			t.Fatalf("Could not create admin user request: %q", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Could not create admin user: %q", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			bs, _ := ioutil.ReadAll(resp.Body)
			t.Fatalf("Failed to create admin user: %s %s\n", resp.Status, string(bs))
		}
	}

	return cleanup, retAddress + "/" + dbName, "admin", "admin"
}
