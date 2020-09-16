package cassandra

import (
	"context"
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/hashicorp/vault/helper/testhelpers/docker"
	"os"
	"testing"
	"time"
)

func PrepareTestContainer(t *testing.T, version string) (func(), string) {
	t.Helper()
	if os.Getenv("CASSANDRA_HOSTS") != "" {
		return func() {}, os.Getenv("CASSANDRA_HOSTS")
	}

	if version == "" {
		version = "3.11"
	}

	var copyFromTo map[string]string
	cwd, _ := os.Getwd()
	fixturePath := fmt.Sprintf("%s/test-fixtures/", cwd)
	if _, err := os.Stat(fixturePath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			// If it doesn't exist, no biggie
			t.Fatal(err)
		}
	} else {
		copyFromTo = map[string]string{
			fixturePath: "/etc/cassandra",
		}
	}

	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo:  "cassandra",
		ImageTag:   version,
		Ports:      []string{"9042/tcp"},
		CopyFromTo: copyFromTo,
		Env:        []string{"CASSANDRA_BROADCAST_ADDRESS=127.0.0.1"},
	})
	if err != nil {
		t.Fatalf("Could not start docker cassandra: %s", err)
	}

	svc, err := runner.StartService(context.Background(), func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		cfg := docker.NewServiceHostPort(host, port)
		clusterConfig := gocql.NewCluster(cfg.Address())
		clusterConfig.Authenticator = gocql.PasswordAuthenticator{
			Username: "cassandra",
			Password: "cassandra",
		}
		clusterConfig.Timeout = 30 * time.Second
		clusterConfig.ProtoVersion = 4
		clusterConfig.Port = port

		session, err := clusterConfig.CreateSession()
		if err != nil {
			return nil, fmt.Errorf("error creating session: %s", err)
		}
		defer session.Close()

		// Create keyspace
		q := session.Query(`CREATE KEYSPACE "vault" WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };`)
		if err := q.Exec(); err != nil {
			t.Fatalf("could not create cassandra keyspace: %v", err)
		}

		// Create table
		q = session.Query(`CREATE TABLE "vault"."entries" (
		    bucket text,
		    key text,
		    value blob,
		    PRIMARY KEY (bucket, key)
		) WITH CLUSTERING ORDER BY (key ASC);`)
		if err := q.Exec(); err != nil {
			t.Fatalf("could not create cassandra table: %v", err)
		}
		return cfg, nil
	})
	if err != nil {
		t.Fatalf("Could not start docker cassandra: %s", err)
	}
	return svc.Cleanup, svc.Config.Address()
}
