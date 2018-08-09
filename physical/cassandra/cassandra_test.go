package cassandra

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/gocql/gocql"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/physical"
	"github.com/ory/dockertest"
)

func TestCassandraBackend(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipping in short mode")
	}

	cleanup, hosts := prepareCassandraTestContainer(t)
	defer cleanup()

	// Run vault tests
	logger := logging.NewVaultLogger(log.Debug)
	b, err := NewCassandraBackend(map[string]string{
		"hosts":            hosts,
		"protocol_version": "3",
	}, logger)

	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}

func TestCassandraBackendBuckets(t *testing.T) {
	expectations := map[string][]string{
		"":          {"."},
		"a":         {"."},
		"a/b":       {".", "a"},
		"a/b/c/d/e": {".", "a", "a/b", "a/b/c", "a/b/c/d"}}

	b := &CassandraBackend{}
	for input, expected := range expectations {
		actual := b.buckets(input)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("bad: %v expected: %v", actual, expected)
		}
	}
}

func prepareCassandraTestContainer(t *testing.T) (func(), string) {
	if os.Getenv("CASSANDRA_HOSTS") != "" {
		return func() {}, os.Getenv("CASSANDRA_HOSTS")
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("cassandra: failed to connect to docker: %s", err)
	}

	resource, err := pool.Run("cassandra", "3.11", []string{"CASSANDRA_BROADCAST_ADDRESS=127.0.0.1"})
	if err != nil {
		t.Fatalf("cassandra: could not start container: %s", err)
	}

	cleanup := func() {
		pool.Purge(resource)
	}

	setup := func() error {
		cluster := gocql.NewCluster("127.0.0.1")
		p, _ := strconv.Atoi(resource.GetPort("9042/tcp"))
		cluster.Port = p
		cluster.Timeout = 5 * time.Second
		sess, err := cluster.CreateSession()
		if err != nil {
			return err
		}
		defer sess.Close()

		// Create keyspace
		q := sess.Query(`CREATE KEYSPACE "vault" WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };`)
		if err := q.Exec(); err != nil {
			t.Fatalf("could not create cassandra keyspace: %v", err)
		}

		// Create table
		q = sess.Query(`CREATE TABLE "vault"."entries" (
		    bucket text,
		    key text,
		    value blob,
		    PRIMARY KEY (bucket, key)
		) WITH CLUSTERING ORDER BY (key ASC);`)
		if err := q.Exec(); err != nil {
			t.Fatalf("could not create cassandra table: %v", err)
		}

		return nil
	}
	if pool.Retry(setup); err != nil {
		cleanup()
		t.Fatalf("cassandra: could not setup container: %s", err)
	}

	return cleanup, fmt.Sprintf("127.0.0.1:%s", resource.GetPort("9042/tcp"))
}
