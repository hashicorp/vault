// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/docker"
)

const postgresVersion = "13.4-buster"

func defaultRunOpts(t *testing.T) docker.RunOptions {
	return docker.RunOptions{
		ContainerName: "postgres",
		ImageRepo:     "docker.mirror.hashicorp.services/postgres",
		ImageTag:      postgresVersion,
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_DB=database",
		},
		Ports:           []string{"5432/tcp"},
		DoNotAutoRemove: false,
		LogConsumer: func(s string) {
			if t.Failed() {
				t.Logf("container logs: %s", s)
			}
		},
	}
}

func PrepareTestContainer(t *testing.T) (func(), string) {
	_, cleanup, url, _ := prepareTestContainer(t, defaultRunOpts(t), "secret", true, false)

	return cleanup, url
}

// PrepareTestContainerWithVaultUser will setup a test container with a Vault
// admin user configured so that we can safely call rotate-root without
// rotating the root DB credentials
func PrepareTestContainerWithVaultUser(t *testing.T, ctx context.Context) (func(), string) {
	runner, cleanup, url, id := prepareTestContainer(t, defaultRunOpts(t), "secret", true, false)

	cmd := []string{"psql", "-U", "postgres", "-c", "CREATE USER vaultadmin WITH LOGIN PASSWORD 'vaultpass' SUPERUSER"}
	_, err := runner.RunCmdInBackground(ctx, id, cmd)
	if err != nil {
		t.Fatalf("Could not run command (%v) in container: %v", cmd, err)
	}

	return cleanup, url
}

func PrepareTestContainerWithSSL(t *testing.T, ctx context.Context, version string) (func(), string) {
	runOpts := defaultRunOpts(t)
	runOpts.Cmd = []string{"-c", "log_statement=all"}
	runner, cleanup, url, id := prepareTestContainer(t, runOpts, "secret", true, false)

	content := "echo 'hostssl all all all cert clientcert=verify-ca' > /var/lib/postgresql/data/pg_hba.conf"
	// Copy the ssl init script into the newly running container.
	buildCtx := docker.NewBuildContext()
	buildCtx["ssl-conf.sh"] = docker.PathContentsFromBytes([]byte(content))
	if err := runner.CopyTo(id, "/usr/local/bin", buildCtx); err != nil {
		t.Fatalf("Could not copy ssl init script into container: %v", err)
	}

	// run the ssl init script to overwrite the pg_hba.conf file and set it to
	// require SSL for each connection
	cmd := []string{"bash", "/usr/local/bin/ssl-conf.sh"}
	_, err := runner.RunCmdInBackground(ctx, id, cmd)
	if err != nil {
		t.Fatalf("Could not run command (%v) in container: %v", cmd, err)
	}

	// reload so the config changes take effect
	cmd = []string{"psql", "-U", "postgres", "-c", "SELECT pg_reload_conf()"}
	_, err = runner.RunCmdInBackground(ctx, id, cmd)
	if err != nil {
		t.Fatalf("Could not run command (%v) in container: %v", cmd, err)
	}

	return cleanup, url
}

func PrepareTestContainerWithPassword(t *testing.T, version, password string) (func(), string) {
	runOpts := defaultRunOpts(t)
	runOpts.Env = []string{
		"POSTGRES_PASSWORD=" + password,
		"POSTGRES_DB=database",
	}

	_, cleanup, url, _ := prepareTestContainer(t, runOpts, password, true, false)

	return cleanup, url
}

func PrepareTestContainerRepmgr(t *testing.T, name, version string, envVars []string) (*docker.Runner, func(), string, string) {
	runOpts := defaultRunOpts(t)
	runOpts.ImageRepo = "docker.mirror.hashicorp.services/bitnami/postgresql-repmgr"
	runOpts.ImageTag = version
	runOpts.Env = append(envVars,
		"REPMGR_PARTNER_NODES=psql-repl-node-0,psql-repl-node-1",
		"REPMGR_PRIMARY_HOST=psql-repl-node-0",
		"REPMGR_PASSWORD=repmgrpass",
		"POSTGRESQL_PASSWORD=secret")
	runOpts.DoNotAutoRemove = true

	return prepareTestContainer(t, runOpts, "secret", false, true)
}

func prepareTestContainer(t *testing.T, runOpts docker.RunOptions, password string, addSuffix, forceLocalAddr bool,
) (*docker.Runner, func(), string, string) {
	if os.Getenv("PG_URL") != "" {
		return nil, func() {}, "", os.Getenv("PG_URL")
	}

	if runOpts.ImageRepo == "bitnami/postgresql-repmgr" {
		runOpts.NetworkID = os.Getenv("POSTGRES_MULTIHOST_NET")
	}

	runner, err := docker.NewServiceRunner(runOpts)
	if err != nil {
		t.Fatalf("Could not start docker Postgres: %s", err)
	}

	svc, containerID, err := runner.StartNewService(context.Background(), addSuffix, forceLocalAddr, connectPostgres(password, runOpts.ImageRepo))
	if err != nil {
		t.Fatalf("Could not start docker Postgres: %s", err)
	}

	return runner, svc.Cleanup, svc.Config.URL().String(), containerID
}

func connectPostgres(password, repo string) docker.ServiceAdapter {
	return func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		u := url.URL{
			Scheme:   "postgres",
			User:     url.UserPassword("postgres", password),
			Host:     fmt.Sprintf("%s:%d", host, port),
			Path:     "postgres",
			RawQuery: "sslmode=disable",
		}

		db, err := sql.Open("pgx", u.String())
		if err != nil {
			return nil, err
		}
		defer db.Close()

		if err = db.Ping(); err != nil {
			return nil, err
		}
		return docker.NewServiceURL(u), nil
	}
}

func StopContainer(t *testing.T, ctx context.Context, runner *docker.Runner, containerID string) {
	if err := runner.Stop(ctx, containerID); err != nil {
		t.Fatalf("Could not stop docker Postgres: %s", err)
	}
}

func RestartContainer(t *testing.T, ctx context.Context, runner *docker.Runner, containerID string) {
	if err := runner.Restart(ctx, containerID); err != nil {
		t.Fatalf("Could not restart docker Postgres: %s", err)
	}
}
