package physical

import (
	"net"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/hashicorp/vault/helper/logformat"
	log "github.com/mgutz/logxi/v1"
	"github.com/ory-am/dockertest"
)

func TestCassandraBackend(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipping in short mode")
	}

	cid, hosts := prepareCassandraTestContainer(t)
	defer cleanupCassandraTestContainer(t, cid)

	// Run vault tests
	logger := logformat.NewVaultLogger(log.LevelTrace)
	b, err := NewBackend("cassandra", logger, map[string]string{
		"hosts": hosts,
	})

	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	testBackend(t, b)
	testBackend_ListPrefix(t, b)
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

func prepareCassandraTestContainer(t *testing.T) (dockertest.ContainerID, string) {
	if os.Getenv("CASSANDRA_HOSTS") != "" {
		return "", os.Getenv("CASSANDRA_HOSTS")
	}

	dockertest.Pull(dockertest.CassandraImageName)
	hosts := ""
	cid, connErr := dockertest.ConnectToCassandra("3.9", 90, time.Second, func(connAddress string) bool {
		host, _port, _ := net.SplitHostPort(connAddress)
		port, _ := strconv.Atoi(_port)

		cluster := gocql.NewCluster(host)
		cluster.Port = port
		cluster.Timeout = 5 * time.Second
		sess, err := cluster.CreateSession()
		if err != nil {
			return false
		}

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

		hosts = connAddress
		return true
	})

	if connErr != nil {
		t.Fatalf("could not connect to cassandra: %v", connErr)
	}

	return cid, hosts
}

func cleanupCassandraTestContainer(t *testing.T, cid dockertest.ContainerID) {
	if cid == "" {
		return
	}
	if err := cid.KillRemove(); err != nil {
		t.Fatal(err)
	}
}
