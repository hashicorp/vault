// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cassandra

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/hashicorp/vault/sdk/helper/docker"
)

type containerConfig struct {
	containerName string
	imageName     string
	version       string
	copyFromTo    map[string]string
	env           []string

	sslOpts *gocql.SslOptions
}

type ContainerOpt func(*containerConfig)

func ContainerName(name string) ContainerOpt {
	return func(cfg *containerConfig) {
		cfg.containerName = name
	}
}

func Image(imageName string, version string) ContainerOpt {
	return func(cfg *containerConfig) {
		cfg.imageName = imageName
		cfg.version = version

		// Reset the environment because there's a very good chance the default environment doesn't apply to the
		// non-default image being used
		cfg.env = nil
	}
}

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

func Env(keyValue string) ContainerOpt {
	return func(cfg *containerConfig) {
		cfg.env = append(cfg.env, keyValue)
	}
}

func SslOpts(sslOpts *gocql.SslOptions) ContainerOpt {
	return func(cfg *containerConfig) {
		cfg.sslOpts = sslOpts
	}
}

type Host struct {
	Name string
	Port string
}

func (h Host) ConnectionURL() string {
	return net.JoinHostPort(h.Name, h.Port)
}

func PrepareTestContainer(t *testing.T, opts ...ContainerOpt) (Host, func()) {
	t.Helper()

	// Skipping on ARM, as this image can't run on ARM architecture
	if strings.Contains(runtime.GOARCH, "arm") {
		t.Skip("Skipping, as this image is not supported on ARM architectures")
	}

	if os.Getenv("CASSANDRA_HOSTS") != "" {
		host, port, err := net.SplitHostPort(os.Getenv("CASSANDRA_HOSTS"))
		if err != nil {
			t.Fatalf("Failed to split host & port from CASSANDRA_HOSTS (%s): %s", os.Getenv("CASSANDRA_HOSTS"), err)
		}
		h := Host{
			Name: host,
			Port: port,
		}
		return h, func() {}
	}

	containerCfg := &containerConfig{
		imageName:     "docker.mirror.hashicorp.services/library/cassandra",
		containerName: "cassandra",
		version:       "3.11",
		env:           []string{"CASSANDRA_BROADCAST_ADDRESS=127.0.0.1"},
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

	runOpts := docker.RunOptions{
		ContainerName: containerCfg.containerName,
		ImageRepo:     containerCfg.imageName,
		ImageTag:      containerCfg.version,
		Ports:         []string{"9042/tcp"},
		CopyFromTo:    copyFromTo,
		Env:           containerCfg.env,
	}
	runner, err := docker.NewServiceRunner(runOpts)
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

	host, port, err := net.SplitHostPort(svc.Config.Address())
	if err != nil {
		t.Fatalf("Failed to split host & port from address (%s): %s", svc.Config.Address(), err)
	}
	h := Host{
		Name: host,
		Port: port,
	}
	return h, svc.Cleanup
}
