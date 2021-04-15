package cassandra

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/hashicorp/vault/helper/testhelpers/docker"
)

type containerConfig struct {
	version    string
	copyFromTo map[string]string
	sslOpts    *gocql.SslOptions
}

type ContainerOpt func(*containerConfig)

func Version(version string) ContainerOpt {
	return func(cfg *containerConfig) {
		cfg.version = version
	}
}

func CopyFromTo(copyFromTo map[string]string) ContainerOpt {
	return func(cfg *containerConfig) {
		cfg.copyFromTo = copyFromTo
	}
}

func SslOpts(sslOpts *gocql.SslOptions) ContainerOpt {
	return func(cfg *containerConfig) {
		cfg.sslOpts = sslOpts
	}
}

func PrepareTestContainer(t *testing.T, opts ...ContainerOpt) (func(), string) {
	t.Helper()
	if os.Getenv("CASSANDRA_HOSTS") != "" {
		return func() {}, os.Getenv("CASSANDRA_HOSTS")
	}

	containerCfg := &containerConfig{
		version: "3.11",
	}

	for _, opt := range opts {
		opt(containerCfg)
	}

	copyFromTo := map[string]string{}
	for from, to := range containerCfg.copyFromTo {
		absFrom, err := filepath.Abs(from)
		if err != nil {
			t.Fatalf("Unable to get absolute path for file %s", from)
		}
		copyFromTo[absFrom] = to
	}

	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo:  "cassandra",
		ImageTag:   containerCfg.version,
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

		clusterConfig.SslOpts = containerCfg.sslOpts

		session, err := clusterConfig.CreateSession()
		if err != nil {
			return nil, fmt.Errorf("error creating session: %s", err)
		}
		defer session.Close()

		// Create keyspace
		query := session.Query(`CREATE KEYSPACE "vault" WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };`)
		if err := query.Exec(); err != nil {
			t.Fatalf("could not create cassandra keyspace: %v", err)
		}

		// Create table
		query = session.Query(`CREATE TABLE "vault"."entries" (
		    bucket text,
		    key text,
		    value blob,
		    PRIMARY KEY (bucket, key)
		) WITH CLUSTERING ORDER BY (key ASC);`)
		if err := query.Exec(); err != nil {
			t.Fatalf("could not create cassandra table: %v", err)
		}
		return cfg, nil
	})
	if err != nil {
		t.Fatalf("Could not start docker cassandra: %s", err)
	}
	return svc.Cleanup, svc.Config.Address()
}
